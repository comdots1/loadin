package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// TaskStatus represents the lifecycle state of a download task.
type TaskStatus string

const (
	StatusQueued    TaskStatus = "queued"
	StatusProbing   TaskStatus = "probing"
	StatusRunning   TaskStatus = "running"
	StatusPaused    TaskStatus = "paused"
	StatusVerifying TaskStatus = "verifying"
	StatusCompleted TaskStatus = "completed"
	StatusError     TaskStatus = "error"
)

const (
	// MinSplitSize is the minimum remaining bytes (1 MB) required for Dynamic Work-Stealing to split an active segment.
	MinSplitSize int64 = 1024 * 1024
	// MaxRetriesPerSegment is the retry threshold with exponential backoff per segment.
	MaxRetriesPerSegment = 6
)

// SegmentState is the serializable and live-inspectable state of a single byte-range segment.
type SegmentState struct {
	ID         int    `json:"id"`
	Start      int64  `json:"start"`
	Current    int64  `json:"current"`
	End        int64  `json:"end"`
	SpeedBps   int64  `json:"speedBps"`
	Status     string `json:"status"` // "streaming", "completed", "retrying", "paused", "split"
	Retries    int    `json:"retries"`
	Generation int    `json:"generation"`
	 currentAtomic atomic.Int64
	 endAtomic     atomic.Int64
	 speedAtomic   atomic.Int64
	 statusMu      sync.RWMutex
}

func (s *SegmentState) getStatus() string {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.Status
}

func (s *SegmentState) setStatus(st string) {
	s.statusMu.Lock()
	s.Status = st
	s.statusMu.Unlock()
}

// LogEntry represents a timestamped diagnostic event in the task's live inspector log.
type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"` // "info", "split", "warn", "success", "error"
	Message string `json:"message"`
}

// TaskSnapshot is the JSON-friendly view sent at 5Hz to the Glassmorphism UI.
type TaskSnapshot struct {
	ID             string          `json:"id"`
	URL            string          `json:"url"`
	FinalURL       string          `json:"finalUrl"`
	FileName       string          `json:"fileName"`
	SaveDir        string          `json:"saveDir"`
	TargetPath     string          `json:"targetPath"`
	Category       string          `json:"category"`
	TotalSize      int64           `json:"totalSize"`
	Downloaded     int64           `json:"downloaded"`
	Progress       float64         `json:"progress"`
	SpeedBps       int64           `json:"speedBps"`
	PeakSpeedBps   int64           `json:"peakSpeedBps"`
	AvgSpeedBps    int64           `json:"avgSpeedBps"`
	ETASeconds     int64           `json:"etaSeconds"`
	Status         TaskStatus      `json:"status"`
	ErrorMessage   string          `json:"errorMessage,omitempty"`
	AcceptRanges   bool            `json:"acceptRanges"`
	MaxConnections int             `json:"maxConnections"`
	ActiveWorkers  int             `json:"activeWorkers"`
	DynamicSplits  int             `json:"dynamicSplits"`
	ContentType    string          `json:"contentType"`
	Server         string          `json:"server"`
	ETag           string          `json:"etag"`
	ExpectedHash   string          `json:"expectedHash,omitempty"`
	Checksum       *ChecksumResult `json:"checksum,omitempty"`
	CreatedAt      int64           `json:"createdAt"`
	CompletedAt    int64           `json:"completedAt,omitempty"`
	SpeedHistory   []int64         `json:"speedHistory"`
	Segments       []SegmentState  `json:"segments"`
	Logs           []LogEntry      `json:"logs"`
}

// DownloadTask manages a single multi-segment download with IDM-style dynamic work stealing.
type DownloadTask struct {
	mu             sync.RWMutex
	ID             string
	URL            string
	FinalURL       string
	FileName       string
	SaveDir        string
	TargetPath     string
	PartPath       string
	MetaPath       string
	Category       string
	TotalSize      int64
	AcceptRanges   bool
	MaxConnections int
	Headers        map[string]string
	ExpectedHash   string
	Checksum       *ChecksumResult
	ContentType    string
	Server         string
	ETag           string
	Status         TaskStatus
	ErrorMessage   string
	CreatedAt      int64
	CompletedAt    int64
	DynamicSplits  int

	downloadedAtomic atomic.Int64
	speedAtomic      atomic.Int64
	peakSpeedAtomic  atomic.Int64
	activeWorkers    atomic.Int32
	startedAt        time.Time
	accumulatedSec   float64

	segments     []*SegmentState
	speedHistory []int64
	logs         []LogEntry

	client        *http.Client
	globalLimiter *RateLimiter
	taskLimiter   *RateLimiter
	cancelFunc    context.CancelFunc
	doneChan      chan struct{}
	onStateChange func()
}

// PersistedTaskMeta is saved to disk so tasks survive app restarts and support instant resume.
type PersistedTaskMeta struct {
	ID             string            `json:"id"`
	URL            string            `json:"url"`
	FinalURL       string            `json:"finalUrl"`
	FileName       string            `json:"fileName"`
	SaveDir        string            `json:"saveDir"`
	TargetPath     string            `json:"targetPath"`
	PartPath       string            `json:"partPath"`
	Category       string            `json:"category"`
	TotalSize      int64             `json:"totalSize"`
	Downloaded     int64             `json:"downloaded"`
	AcceptRanges   bool              `json:"acceptRanges"`
	MaxConnections int               `json:"maxConnections"`
	Headers        map[string]string `json:"headers"`
	ExpectedHash   string            `json:"expectedHash"`
	Checksum       *ChecksumResult   `json:"checksum"`
	ContentType    string            `json:"contentType"`
	Server         string            `json:"server"`
	ETag           string            `json:"etag"`
	Status         TaskStatus        `json:"status"`
	ErrorMessage   string            `json:"errorMessage"`
	CreatedAt      int64             `json:"createdAt"`
	CompletedAt    int64             `json:"completedAt"`
	DynamicSplits  int               `json:"dynamicSplits"`
	PeakSpeedBps   int64             `json:"peakSpeedBps"`
	Segments       []SegmentState    `json:"segments"`
	Logs           []LogEntry        `json:"logs"`
}

func (t *DownloadTask) addLog(level, msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	entry := LogEntry{
		Time:    time.Now().Format("15:04:05.000"),
		Level:   level,
		Message: msg,
	}
	t.logs = append(t.logs, entry)
	if len(t.logs) > 80 {
		t.logs = t.logs[len(t.logs)-80:]
	}
}

// Snapshot returns an immutable, thread-safe snapshot of the task for UI telemetry.
func (t *DownloadTask) Snapshot() TaskSnapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()

	dl := t.downloadedAtomic.Load()
	total := t.TotalSize
	if total > 0 && dl > total {
		dl = total
	}
	var progress float64
	if total > 0 {
		progress = (float64(dl) / float64(total)) * 100.0
		if progress > 100.0 {
			progress = 100.0
		}
	} else if t.Status == StatusCompleted {
		progress = 100.0
	}

	speed := t.speedAtomic.Load()
	peak := t.peakSpeedAtomic.Load()

	var eta int64 = -1
	if t.Status == StatusRunning && speed > 0 && total > dl {
		eta = (total - dl) / speed
	} else if t.Status == StatusCompleted {
		eta = 0
	}

	var avgSpeed int64
	elapsed := t.accumulatedSec
	if t.Status == StatusRunning && !t.startedAt.IsZero() {
		elapsed += time.Since(t.startedAt).Seconds()
	}
	if elapsed > 0.2 && dl > 0 {
		avgSpeed = int64(float64(dl) / elapsed)
	}

	segs := make([]SegmentState, len(t.segments))
	for i, s := range t.segments {
		segs[i] = SegmentState{
			ID:         s.ID,
			Start:      s.Start,
			Current:    s.currentAtomic.Load(),
			End:        s.endAtomic.Load(),
			SpeedBps:   s.speedAtomic.Load(),
			Status:     s.getStatus(),
			Retries:    s.Retries,
			Generation: s.Generation,
		}
	}

	hist := make([]int64, len(t.speedHistory))
	copy(hist, t.speedHistory)

	logsCopy := make([]LogEntry, len(t.logs))
	copy(logsCopy, t.logs)

	return TaskSnapshot{
		ID:             t.ID,
		URL:            t.URL,
		FinalURL:       t.FinalURL,
		FileName:       t.FileName,
		SaveDir:        t.SaveDir,
		TargetPath:     t.TargetPath,
		Category:       t.Category,
		TotalSize:      total,
		Downloaded:     dl,
		Progress:       progress,
		SpeedBps:       speed,
		PeakSpeedBps:   peak,
		AvgSpeedBps:    avgSpeed,
		ETASeconds:     eta,
		Status:         t.Status,
		ErrorMessage:   t.ErrorMessage,
		AcceptRanges:   t.AcceptRanges,
		MaxConnections: t.MaxConnections,
		ActiveWorkers:  int(t.activeWorkers.Load()),
		DynamicSplits:  t.DynamicSplits,
		ContentType:    t.ContentType,
		Server:         t.Server,
		ETag:           t.ETag,
		ExpectedHash:   t.ExpectedHash,
		Checksum:       t.Checksum,
		CreatedAt:      t.CreatedAt,
		CompletedAt:    t.CompletedAt,
		SpeedHistory:   hist,
		Segments:       segs,
		Logs:           logsCopy,
	}
}

// ExportMeta serializes the task for crash-safe persistence.
func (t *DownloadTask) ExportMeta() PersistedTaskMeta {
	snap := t.Snapshot()
	st := snap.Status
	if st == StatusRunning || st == StatusProbing || st == StatusVerifying {
		st = StatusPaused
	}
	return PersistedTaskMeta{
		ID:             snap.ID,
		URL:            snap.URL,
		FinalURL:       snap.FinalURL,
		FileName:       snap.FileName,
		SaveDir:        snap.SaveDir,
		TargetPath:     snap.TargetPath,
		PartPath:       t.PartPath,
		Category:       snap.Category,
		TotalSize:      snap.TotalSize,
		Downloaded:     snap.Downloaded,
		AcceptRanges:   snap.AcceptRanges,
		MaxConnections: snap.MaxConnections,
		Headers:        t.Headers,
		ExpectedHash:   snap.ExpectedHash,
		Checksum:       snap.Checksum,
		ContentType:    snap.ContentType,
		Server:         snap.Server,
		ETag:           snap.ETag,
		Status:         st,
		ErrorMessage:   snap.ErrorMessage,
		CreatedAt:      snap.CreatedAt,
		CompletedAt:    snap.CompletedAt,
		DynamicSplits:  snap.DynamicSplits,
		PeakSpeedBps:   snap.PeakSpeedBps,
		Segments:       snap.Segments,
		Logs:           snap.Logs,
	}
}

// Start launches the download task asynchronously.
func (t *DownloadTask) Start() {
	t.mu.Lock()
	if t.Status == StatusRunning || t.Status == StatusProbing || t.Status == StatusVerifying {
		t.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.cancelFunc = cancel
	t.doneChan = make(chan struct{})
	t.Status = StatusProbing
	t.ErrorMessage = ""
	t.startedAt = time.Now()
	t.mu.Unlock()

	go t.run(ctx)
}

// Pause gracefully halts all active HTTP segment streams while preserving exact byte offsets.
func (t *DownloadTask) Pause() {
	t.mu.Lock()
	if t.Status != StatusRunning && t.Status != StatusProbing && t.Status != StatusQueued {
		t.mu.Unlock()
		return
	}
	cancel := t.cancelFunc
	done := t.doneChan
	t.Status = StatusPaused
	if !t.startedAt.IsZero() {
		t.accumulatedSec += time.Since(t.startedAt).Seconds()
		t.startedAt = time.Time{}
	}
	t.speedAtomic.Store(0)
	for _, s := range t.segments {
		s.speedAtomic.Store(0)
		if s.getStatus() == "streaming" {
			s.setStatus("paused")
		}
	}
	t.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if done != nil {
		select {
		case <-done:
		case <-time.After(2 * time.Second):
		}
	}
	t.addLog("warn", "[일시정지] 다운로드가 일시정지되었습니다. (세그먼트 오프셋 안전 저장됨)")
	t.saveMetaFile()
	if t.onStateChange != nil {
		t.onStateChange()
	}
}

func (t *DownloadTask) run(ctx context.Context) {
	defer close(t.doneChan)

	// Step 1: Probe metadata if not yet initialized
	t.mu.RLock()
	needProbe := t.TotalSize == 0 && len(t.segments) == 0
	rawURL := t.URL
	headers := t.Headers
	t.mu.RUnlock()

	if needProbe {
		t.addLog("info", fmt.Sprintf("[프로빙] 원격 서버 메타데이터 및 Range 지원 확인 중: %s", rawURL))
		probe, err := ProbeURL(ctx, t.client, rawURL, headers)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			t.fail(fmt.Sprintf("서버 연결 실패: %v", err))
			return
		}

		t.mu.Lock()
		t.FinalURL = probe.FinalURL
		if t.FileName == "" || strings.HasPrefix(t.FileName, "download_") {
			t.FileName = probe.FileName
		}
		if t.Category == "" || t.Category == "general" {
			t.Category = probe.Category
		}
		t.TotalSize = probe.TotalSize
		t.AcceptRanges = probe.AcceptRanges
		t.ContentType = probe.ContentType
		t.Server = probe.Server
		t.ETag = probe.ETag

		// Prepare file paths
		_ = os.MkdirAll(t.SaveDir, 0755)
		t.TargetPath = resolveUniqueFilePath(t.SaveDir, t.FileName)
		t.PartPath = t.TargetPath + ".aether.part"
		t.MetaPath = t.TargetPath + ".aether.meta"
		t.mu.Unlock()

		t.addLog("info", fmt.Sprintf("[분석완료] 크기=%s, 멀티 세그먼트(Range)=%v, 서버=%s",
			formatBytesSimple(probe.TotalSize), probe.AcceptRanges, nonEmpty(probe.Server, "CDN")))
	}

	// Ensure paths are set even if restored or pre-probed
	t.mu.Lock()
	_ = os.MkdirAll(t.SaveDir, 0755)
	if t.TargetPath == "" {
		t.TargetPath = resolveUniqueFilePath(t.SaveDir, t.FileName)
	}
	if t.PartPath == "" {
		t.PartPath = t.TargetPath + ".aether.part"
	}
	if t.MetaPath == "" {
		t.MetaPath = t.TargetPath + ".aether.meta"
	}
	t.Status = StatusRunning
	t.startedAt = time.Now()
	t.mu.Unlock()

	if t.onStateChange != nil {
		t.onStateChange()
	}

	// Step 2: Open or pre-allocate .aether.part file
	file, err := os.OpenFile(t.PartPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.fail(fmt.Sprintf("디스크 파일 생성 실패: %v", err))
		return
	}
	defer file.Close()

	// Initialize segments if first run
	t.mu.Lock()
	if len(t.segments) == 0 {
		if t.AcceptRanges && t.TotalSize > 0 {
			// Pre-allocate file size on disk for zero-fragmentation positional WriteAt
			_ = file.Truncate(t.TotalSize)

			numConns := t.MaxConnections
			if numConns < 1 {
				numConns = 8
			}
			// Don't create tiny segments smaller than 256KB
			maxReasonable := int(t.TotalSize / (256 * 1024))
			if maxReasonable < 1 {
				maxReasonable = 1
			}
			if numConns > maxReasonable {
				numConns = maxReasonable
			}

			chunkSize := t.TotalSize / int64(numConns)
			for i := 0; i < numConns; i++ {
				start := int64(i) * chunkSize
				end := start + chunkSize - 1
				if i == numConns-1 {
					end = t.TotalSize - 1
				}
				seg := &SegmentState{
					ID:         i + 1,
					Start:      start,
					Current:    start,
					End:        end,
					Status:     "streaming",
					Generation: 0,
				}
				seg.currentAtomic.Store(start)
				seg.endAtomic.Store(end)
				t.segments = append(t.segments, seg)
			}
		} else {
			// Single stream fallback for non-Range or chunked transfer servers
			end := t.TotalSize - 1
			if t.TotalSize <= 0 {
				end = -1
			}
			seg := &SegmentState{
				ID:      1,
				Start:   0,
				Current: 0,
				End:     end,
				Status:  "streaming",
			}
			seg.currentAtomic.Store(0)
			seg.endAtomic.Store(end)
			t.segments = append(t.segments, seg)
		}
	} else {
		// Recompute already downloaded bytes from persisted segments
		var already int64
		for _, s := range t.segments {
			cur := s.currentAtomic.Load()
			end := s.endAtomic.Load()
			if cur > s.Start {
				doneBytes := cur - s.Start
				if end >= 0 && cur > end+1 {
					doneBytes = (end + 1) - s.Start
				}
				already += doneBytes
			}
			if end >= 0 && cur > end {
				s.setStatus("completed")
			} else {
				s.setStatus("streaming")
			}
		}
		t.downloadedAtomic.Store(already)
	}
	initialSegments := make([]*SegmentState, len(t.segments))
	copy(initialSegments, t.segments)
	t.mu.Unlock()

	t.addLog("info", fmt.Sprintf("[가속시작] %d개 병렬 세그먼트 엔진 가동 (Dynamic Work-Stealing 활성화)", len(initialSegments)))

	// Step 3: Start telemetry & speed sampling ticker (every 250ms)
	telemetryCtx, stopTelemetry := context.WithCancel(ctx)
	go t.runTelemetryTicker(telemetryCtx)

	// Step 4: Spawn worker goroutines with Dynamic Segment Stealing
	var wg sync.WaitGroup
	errChan := make(chan error, 64)

	var spawnSegmentWorker func(seg *SegmentState)
	spawnSegmentWorker = func(seg *SegmentState) {
		wg.Add(1)
		t.activeWorkers.Add(1)
		go func(s *SegmentState) {
			defer func() {
				s.speedAtomic.Store(0)
				t.activeWorkers.Add(-1)
				wg.Done()
			}()

			if err := t.downloadSegment(ctx, file, s); err != nil {
				if !errors.Is(err, context.Canceled) {
					s.setStatus("error")
					select {
					case errChan <- err:
					default:
					}
				}
				return
			}

			s.setStatus("completed")

			// IDM Dynamic Work-Stealing: When this worker finishes early, try to split the largest remaining segment!
			for ctx.Err() == nil {
				newSeg := t.tryStealSegmentWork()
				if newSeg == nil {
					break
				}
				if err := t.downloadSegment(ctx, file, newSeg); err != nil {
					if !errors.Is(err, context.Canceled) {
						newSeg.setStatus("error")
						select {
						case errChan <- err:
						default:
						}
					}
					return
				}
				newSeg.setStatus("completed")
			}
		}(seg)
	}

	for _, seg := range initialSegments {
		cur := seg.currentAtomic.Load()
		end := seg.endAtomic.Load()
		if end >= 0 && cur > end {
			seg.setStatus("completed")
			continue
		}
		spawnSegmentWorker(seg)
	}

	wg.Wait()
	stopTelemetry()

	if ctx.Err() != nil {
		// Paused or cancelled
		return
	}

	select {
	case workerErr := <-errChan:
		t.fail(fmt.Sprintf("세그먼트 다운로드 오류: %v", workerErr))
		return
	default:
	}

	// Close part file before rename & checksum verification
	_ = file.Sync()
	_ = file.Close()

	// Step 5: Rename .aether.part -> Final TargetPath
	if err := os.Rename(t.PartPath, t.TargetPath); err != nil {
		// If target exists, remove and retry
		_ = os.Remove(t.TargetPath)
		if err2 := os.Rename(t.PartPath, t.TargetPath); err2 != nil {
			t.fail(fmt.Sprintf("파일 최종 병합/이름변경 실패: %v", err2))
			return
		}
	}
	_ = os.Remove(t.MetaPath)

	// Step 6: Compute SHA-256 & MD5 Checksums
	t.mu.Lock()
	t.Status = StatusVerifying
	expectedHash := t.ExpectedHash
	targetPath := t.TargetPath
	t.mu.Unlock()

	t.addLog("info", "[검증] 무결성 해시(SHA-256 / MD5) 고속 스트리밍 검증 중...")
	if t.onStateChange != nil {
		t.onStateChange()
	}

	checksum, hashErr := ComputeFileChecksums(targetPath, expectedHash)
	if hashErr == nil {
		t.mu.Lock()
		t.Checksum = checksum
		t.mu.Unlock()
		if checksum.Verified != nil {
			if *checksum.Verified {
				t.addLog("success", fmt.Sprintf("[검증완료] 해시 일치 검증 성공! SHA-256: %s...", checksum.SHA256[:16]))
			} else {
				t.addLog("error", fmt.Sprintf("[경고] 해시 불일치! 계산된 SHA-256: %s", checksum.SHA256))
			}
		} else {
			t.addLog("info", fmt.Sprintf("[SHA-256] %s", checksum.SHA256))
		}
	}

	t.mu.Lock()
	if !t.startedAt.IsZero() {
		t.accumulatedSec += time.Since(t.startedAt).Seconds()
		t.startedAt = time.Time{}
	}
	if t.TotalSize <= 0 {
		t.TotalSize = t.downloadedAtomic.Load()
	} else {
		t.downloadedAtomic.Store(t.TotalSize)
	}
	t.speedAtomic.Store(0)
	t.Status = StatusCompleted
	t.CompletedAt = time.Now().Unix()
	t.mu.Unlock()

	t.addLog("success", fmt.Sprintf("[완료] 다운로드 완료! 저장 위치: %s", t.TargetPath))
	if t.onStateChange != nil {
		t.onStateChange()
	}
}

// tryStealSegmentWork implements Dynamic File Segmentation (Work-Stealing):
// Finds the active segment with the largest remaining byte count >= MinSplitSize (1 MB)
// and atomically splits its tail half into a brand-new segment.
func (t *DownloadTask) tryStealSegmentWork() *SegmentState {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.AcceptRanges || t.TotalSize <= 0 {
		return nil
	}
	if len(t.segments) >= 64 {
		// Cap total dynamic segments at 64 for visual clarity & overhead balance
		return nil
	}

	var victim *SegmentState
	var maxRemaining int64

	for _, s := range t.segments {
		if s.getStatus() != "streaming" {
			continue
		}
		cur := s.currentAtomic.Load()
		end := s.endAtomic.Load()
		rem := end - cur
		if rem > maxRemaining {
			maxRemaining = rem
			victim = s
		}
	}

	if victim == nil || maxRemaining < MinSplitSize*2 {
		return nil
	}

	cur := victim.currentAtomic.Load()
	oldEnd := victim.endAtomic.Load()
	rem := oldEnd - cur
	if rem < MinSplitSize*2 {
		return nil
	}

	mid := cur + (rem / 2)
	// Atomically shrink the victim's end boundary to mid - 1
	victim.endAtomic.Store(mid - 1)
	victim.End = mid - 1

	newID := len(t.segments) + 1
	newSeg := &SegmentState{
		ID:         newID,
		Start:      mid,
		Current:    mid,
		End:        oldEnd,
		Status:     "streaming",
		Generation: victim.Generation + 1,
	}
	newSeg.currentAtomic.Store(mid)
	newSeg.endAtomic.Store(oldEnd)

	t.segments = append(t.segments, newSeg)
	t.DynamicSplits++

	// Log outside lock? Wait: we already hold t.mu.Lock(), so append directly to t.logs
	entry := LogEntry{
		Time:  time.Now().Format("15:04:05.000"),
		Level: "split",
		Message: fmt.Sprintf("[동적분할] 세그먼트 #%d 생성 (세그먼트 #%d의 남은 %s 절반 분할 가속)",
			newID, victim.ID, formatBytesSimple(rem)),
	}
	t.logs = append(t.logs, entry)
	if len(t.logs) > 80 {
		t.logs = t.logs[len(t.logs)-80:]
	}

	return newSeg
}

func (t *DownloadTask) downloadSegment(ctx context.Context, file *os.File, seg *SegmentState) error {
	for attempt := 0; attempt <= MaxRetriesPerSegment; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		cur := seg.currentAtomic.Load()
		end := seg.endAtomic.Load()
		if end >= 0 && cur > end {
			return nil
		}

		if attempt > 0 {
			seg.Retries = attempt
			seg.setStatus("retrying")
			backoff := time.Duration(attempt*attempt*400) * time.Millisecond
			if backoff > 5*time.Second {
				backoff = 5 * time.Second
			}
			t.addLog("warn", fmt.Sprintf("[재시도] 세그먼트 #%d 재연결 시도 (%d/%d)...", seg.ID, attempt, MaxRetriesPerSegment))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		seg.setStatus("streaming")
		err := t.streamSegmentOnce(ctx, file, seg)
		if err == nil {
			return nil
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		if attempt == MaxRetriesPerSegment {
			return fmt.Errorf("세그먼트 #%d 최대 재시도 초과: %w", seg.ID, err)
		}
	}
	return nil
}

func (t *DownloadTask) streamSegmentOnce(ctx context.Context, file *os.File, seg *SegmentState) error {
	cur := seg.currentAtomic.Load()
	end := seg.endAtomic.Load()
	if end >= 0 && cur > end {
		return nil
	}

	t.mu.RLock()
	targetURL := t.FinalURL
	if targetURL == "" {
		targetURL = t.URL
	}
	acceptRanges := t.AcceptRanges
	headers := t.Headers
	t.mu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("Accept", "*/*")
	for k, v := range headers {
		if strings.TrimSpace(k) != "" && strings.TrimSpace(v) != "" {
			req.Header.Set(k, v)
		}
	}

	if acceptRanges && end >= 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", cur, end))
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	bufPtr := BufferPool.Get().(*[]byte)
	defer BufferPool.Put(bufPtr)
	buf := *bufPtr

	var segWindowBytes int64
	lastRateTick := time.Now()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		// Check if Dynamic Work-Stealing shrunk our end offset!
		dynamicEnd := seg.endAtomic.Load()
		curPos := seg.currentAtomic.Load()
		if dynamicEnd >= 0 && curPos > dynamicEnd {
			return nil
		}

		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			// Re-verify dynamicEnd right before writing so we never overwrite a stolen half
			dynamicEnd = seg.endAtomic.Load()
			curPos = seg.currentAtomic.Load()
			writeLen := int64(n)
			if dynamicEnd >= 0 {
				remainingAllowed := (dynamicEnd - curPos) + 1
				if remainingAllowed <= 0 {
					return nil
				}
				if writeLen > remainingAllowed {
					writeLen = remainingAllowed
				}
			}

			// Enforce global & task bandwidth token-bucket limiters
			if t.globalLimiter != nil {
				if err := t.globalLimiter.Wait(ctx, int(writeLen)); err != nil {
					return err
				}
			}
			if t.taskLimiter != nil {
				if err := t.taskLimiter.Wait(ctx, int(writeLen)); err != nil {
					return err
				}
			}

			// Lock-free positional write at exact byte offset
			if _, wErr := file.WriteAt(buf[:writeLen], curPos); wErr != nil {
				return wErr
			}

			seg.currentAtomic.Add(writeLen)
			t.downloadedAtomic.Add(writeLen)
			segWindowBytes += writeLen

			now := time.Now()
			if elapsed := now.Sub(lastRateTick).Seconds(); elapsed >= 0.25 {
				seg.speedAtomic.Store(int64(float64(segWindowBytes) / elapsed))
				segWindowBytes = 0
				lastRateTick = now
			}

			if dynamicEnd >= 0 && (curPos+writeLen) > dynamicEnd {
				return nil
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				// Verify we actually reached the current dynamicEnd
				dynamicEnd = seg.endAtomic.Load()
				curPos = seg.currentAtomic.Load()
				if dynamicEnd >= 0 && curPos <= dynamicEnd {
					return fmt.Errorf("조기 연결 종료 (남은 바이트: %d)", (dynamicEnd-curPos)+1)
				}
				return nil
			}
			return readErr
		}
	}
}

func (t *DownloadTask) runTelemetryTicker(ctx context.Context) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	lastBytes := t.downloadedAtomic.Load()
	lastTime := time.Now()
	tickCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			curBytes := t.downloadedAtomic.Load()
			elapsed := now.Sub(lastTime).Seconds()
			if elapsed > 0 {
				instantSpeed := int64(float64(curBytes-lastBytes) / elapsed)
				if instantSpeed < 0 {
					instantSpeed = 0
				}
				// Smooth EMA (70% new, 30% previous) for jitter-free UI readout
				prev := t.speedAtomic.Load()
				smoothed := instantSpeed
				if prev > 0 {
					smoothed = int64(float64(instantSpeed)*0.7 + float64(prev)*0.3)
				}
				t.speedAtomic.Store(smoothed)
				if smoothed > t.peakSpeedAtomic.Load() {
					t.peakSpeedAtomic.Store(smoothed)
				}

				t.mu.Lock()
				t.speedHistory = append(t.speedHistory, smoothed)
				if len(t.speedHistory) > 60 {
					t.speedHistory = t.speedHistory[len(t.speedHistory)-60:]
				}
				t.mu.Unlock()
			}
			lastBytes = curBytes
			lastTime = now
			tickCount++

			// Persist segment progress every 2 seconds (8 ticks)
			if tickCount%8 == 0 {
				t.saveMetaFile()
			}
		}
	}
}

func (t *DownloadTask) saveMetaFile() {
	t.mu.RLock()
	metaPath := t.MetaPath
	t.mu.RUnlock()
	if metaPath == "" {
		return
	}
	meta := t.ExportMeta()
	if data, err := json.MarshalIndent(meta, "", "  "); err == nil {
		_ = os.WriteFile(metaPath, data, 0644)
	}
}

func (t *DownloadTask) fail(msg string) {
	t.mu.Lock()
	t.Status = StatusError
	t.ErrorMessage = msg
	t.speedAtomic.Store(0)
	if !t.startedAt.IsZero() {
		t.accumulatedSec += time.Since(t.startedAt).Seconds()
		t.startedAt = time.Time{}
	}
	t.mu.Unlock()
	t.addLog("error", "[오류] "+msg)
	t.saveMetaFile()
	if t.onStateChange != nil {
		t.onStateChange()
	}
}

func resolveUniqueFilePath(dir, fileName string) string {
	clean := SanitizeFileName(fileName)
	candidate := filepath.Join(dir, clean)
	if _, err := os.Stat(candidate); os.IsNotExist(err) {
		return candidate
	}
	ext := filepath.Ext(clean)
	base := strings.TrimSuffix(clean, ext)
	for i := 1; i < 1000; i++ {
		alt := filepath.Join(dir, fmt.Sprintf("%s (%d)%s", base, i, ext))
		if _, err := os.Stat(alt); os.IsNotExist(err) {
			return alt
		}
	}
	return filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, time.Now().Unix(), ext))
}

func formatBytesSimple(b int64) string {
	if b < 0 {
		return "알 수 없음"
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func nonEmpty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
