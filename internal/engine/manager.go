package engine

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aetherdl/aetherdl/internal/platform"
)

// Settings holds user-configurable engine preferences.
type Settings struct {
	DefaultSaveDir      string `json:"defaultSaveDir"`
	DefaultConnections  int    `json:"defaultConnections"`  // 1..32
	MaxConcurrentTasks  int    `json:"maxConcurrentTasks"`  // 1..10
	GlobalSpeedLimitBps int64  `json:"globalSpeedLimitBps"` // 0 = unlimited
	AutoCategorize      bool   `json:"autoCategorize"`      // Sort into category subfolders
	ClipboardMonitor    bool   `json:"clipboardMonitor"`    // Watch clipboard for download URLs
	BrowserIntercept    bool   `json:"browserIntercept"`    // Accept browser extension intercepts
	SoundEffects        bool   `json:"soundEffects"`        // Play subtle glass chime on completion
	RunOnStartup        bool   `json:"runOnStartup"`        // Launch Loadin automatically at Windows Startup
	Theme               string `json:"theme"`               // "black" or "white"
	Language            string `json:"language"`            // "ko", "en", "ja"
}

// GlobalStateSnapshot is broadcast at 5Hz over SSE to keep the Glassmorphism UI in sync.
type GlobalStateSnapshot struct {
	TotalSpeedBps       int64                    `json:"totalSpeedBps"`
	PeakGlobalSpeedBps  int64                    `json:"peakGlobalSpeedBps"`
	ActiveTasksCount    int                      `json:"activeTasksCount"`
	ActiveWorkersCount  int                      `json:"activeWorkersCount"`
	GlobalSpeedHistory  []int64                  `json:"globalSpeedHistory"`
	DetectedClipboard   string                   `json:"detectedClipboard,omitempty"`
	BrowserInterceptCnt int                      `json:"browserInterceptCount"`
	ProtocolRegistered  bool                     `json:"protocolRegistered"`
	ExtensionDir        string                   `json:"extensionDir"`
	Hardware            platform.HardwareProfile `json:"hardware"`
	DiskSpace           platform.DiskSpaceInfo   `json:"diskSpace"`
	Settings            Settings                 `json:"settings"`
	Tasks               []TaskSnapshot           `json:"tasks"`
}

// CreateTaskRequest defines parameters for adding a new download task.
type CreateTaskRequest struct {
	URL            string            `json:"url"`
	FileName       string            `json:"fileName,omitempty"`
	SaveDir        string            `json:"saveDir,omitempty"`
	Category       string            `json:"category,omitempty"`
	MaxConnections int               `json:"maxConnections,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	ExpectedHash   string            `json:"expectedHash,omitempty"`
	AutoStart      *bool             `json:"autoStart,omitempty"`
	TotalSize      int64             `json:"totalSize,omitempty"`
	AcceptRanges   bool              `json:"acceptRanges,omitempty"`
	Source         string            `json:"source,omitempty"`
}

// Manager coordinates all download tasks, queue scheduling, persistence, and clipboard/browser sniffing.
type Manager struct {
	mu                  sync.RWMutex
	tasks               map[string]*DownloadTask
	taskOrder           []string
	settings            Settings
	stateFilePath       string
	extensionDir        string
	client              *http.Client
	globalLimiter       *RateLimiter
	globalSpeedHistory  []int64
	peakGlobalSpeed     int64
	lastClipboard       string
	detectedClipboard   string
	browserInterceptCnt int
	protocolRegistered  bool
}

// NewManager creates and restores the Loadin download engine manager.
func NewManager() *Manager {
	defaultDir := platform.DefaultDownloadsDir()
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		appData = defaultDir
	}
	configDir := filepath.Join(appData, "Loadin")
	_ = os.MkdirAll(configDir, 0755)

	hw := platform.GetHardwareProfile()

	extDir := ResolveAndRepairExtensionDirs(configDir)

	s := Settings{
		DefaultSaveDir:      defaultDir,
		DefaultConnections:  hw.RecommendedConnections,
		MaxConcurrentTasks:  hw.RecommendedConcurrent,
		GlobalSpeedLimitBps: 0,
		AutoCategorize:      true,
		ClipboardMonitor:    true,
		BrowserIntercept:    true,
		SoundEffects:        true,
		RunOnStartup:        platform.IsStartupAutoRunEnabled(),
		Theme:               "black",
		Language:            "ko",
	}

	m := &Manager{
		tasks:              make(map[string]*DownloadTask),
		settings:           s,
		stateFilePath:      filepath.Join(configDir, "loadin_state.json"),
		extensionDir:       extDir,
		client:             NewOptimizedHTTPClient(),
		globalLimiter:      NewRateLimiter(0),
		globalSpeedHistory: make([]int64, 48),
		protocolRegistered: platform.IsBrowserProtocolRegistered(),
	}

	m.loadState()
	go m.backgroundLoop()
	return m
}

// HTTPClient returns the shared optimized HTTP client.
func (m *Manager) HTTPClient() *http.Client {
	return m.client
}

// ExtensionDir returns the path to the browser extension folder, ensuring all files are present.
func (m *Manager) ExtensionDir() string {
	m.mu.RLock()
	dir := m.extensionDir
	m.mu.RUnlock()
	_ = EnsureBrowserExtensionFiles(dir)
	return dir
}

// IncrementBrowserIntercept records a browser-triggered download event.
func (m *Manager) IncrementBrowserIntercept() {
	m.mu.Lock()
	m.browserInterceptCnt++
	m.mu.Unlock()
}

// RefreshProtocolStatus updates cached protocol registration state.
func (m *Manager) RefreshProtocolStatus() bool {
	reg := platform.IsBrowserProtocolRegistered()
	m.mu.Lock()
	m.protocolRegistered = reg
	m.mu.Unlock()
	return reg
}

// GetSettings returns the current engine configuration.
func (m *Manager) GetSettings() Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings
}

// UpdateSettings updates and persists engine configuration at runtime.
func (m *Manager) UpdateSettings(newSettings Settings) Settings {
	m.mu.Lock()
	if newSettings.DefaultSaveDir != "" {
		m.settings.DefaultSaveDir = newSettings.DefaultSaveDir
		_ = os.MkdirAll(newSettings.DefaultSaveDir, 0755)
	}
	if newSettings.DefaultConnections >= 1 && newSettings.DefaultConnections <= 32 {
		m.settings.DefaultConnections = newSettings.DefaultConnections
	}
	if newSettings.MaxConcurrentTasks >= 1 && newSettings.MaxConcurrentTasks <= 16 {
		m.settings.MaxConcurrentTasks = newSettings.MaxConcurrentTasks
	}
	if newSettings.GlobalSpeedLimitBps >= 0 {
		m.settings.GlobalSpeedLimitBps = newSettings.GlobalSpeedLimitBps
		m.globalLimiter.SetLimit(newSettings.GlobalSpeedLimitBps)
	}
	m.settings.AutoCategorize = newSettings.AutoCategorize
	m.settings.ClipboardMonitor = newSettings.ClipboardMonitor
	m.settings.BrowserIntercept = newSettings.BrowserIntercept
	m.settings.SoundEffects = newSettings.SoundEffects
	if m.settings.RunOnStartup != newSettings.RunOnStartup {
		m.settings.RunOnStartup = newSettings.RunOnStartup
		_ = platform.SetStartupAutoRun(newSettings.RunOnStartup)
	}
	if newSettings.Theme == "white" || newSettings.Theme == "light-frost" {
		m.settings.Theme = "white"
	} else if newSettings.Theme != "" {
		m.settings.Theme = "black"
	}
	if newSettings.Language == "en" || newSettings.Language == "ja" || newSettings.Language == "ko" {
		m.settings.Language = newSettings.Language
	}
	res := m.settings
	m.mu.Unlock()

	go platform.StyleMainWindowGlass(res.Theme)
	m.saveState()
	m.scheduleQueue()
	return res
}

// DismissClipboard clears the currently detected clipboard URL banner.
func (m *Manager) DismissClipboard() {
	m.mu.Lock()
	m.detectedClipboard = ""
	m.mu.Unlock()
}

// AddTask registers a new download task and schedules it.
func (m *Manager) AddTask(req CreateTaskRequest) (*TaskSnapshot, error) {
	m.mu.Lock()
	rawURL := strings.TrimSpace(req.URL)
	if m.detectedClipboard == rawURL {
		m.detectedClipboard = ""
	}

	cat := req.Category
	if cat == "" && req.FileName != "" {
		cat = CategorizeFile(req.FileName, "")
	}
	if cat == "" {
		cat = CategorizeFile(rawURL, "")
	}

	saveDir := req.SaveDir
	if saveDir == "" {
		saveDir = m.settings.DefaultSaveDir
		if m.settings.AutoCategorize && cat != "" && cat != "general" {
			saveDir = filepath.Join(saveDir, strings.Title(cat))
		}
	}
	_ = os.MkdirAll(saveDir, 0755)

	id := generateTaskID()
	fileName := SanitizeFileName(req.FileName)
	if req.FileName == "" {
		if probe, err := ProbeURL(context.Background(), m.client, rawURL, req.Headers); err == nil {
			fileName = probe.FileName
			if req.TotalSize == 0 {
				req.TotalSize = probe.TotalSize
				req.AcceptRanges = probe.AcceptRanges
			}
			if cat == "" || cat == "general" {
				cat = probe.Category
				if req.SaveDir == "" && m.settings.AutoCategorize && cat != "" && cat != "general" {
					saveDir = filepath.Join(m.settings.DefaultSaveDir, strings.Title(cat))
					_ = os.MkdirAll(saveDir, 0755)
				}
			}
		} else {
			fileName = "download_" + id[:6]
		}
	}

	conns := req.MaxConnections
	if conns <= 0 {
		if req.TotalSize > 0 {
			rec, _ := platform.RecommendThreadsForFileSize(req.TotalSize)
			conns = rec
		} else {
			conns = m.settings.DefaultConnections
		}
	}
	if conns > 32 {
		conns = 32
	}

	targetPath := resolveUniqueFilePath(saveDir, fileName)
	task := &DownloadTask{
		ID:             id,
		URL:            rawURL,
		FileName:       filepath.Base(targetPath),
		SaveDir:        saveDir,
		TargetPath:     targetPath,
		PartPath:       targetPath + ".aether.part",
		MetaPath:       targetPath + ".aether.meta",
		Category:       cat,
		TotalSize:      req.TotalSize,
		AcceptRanges:   req.AcceptRanges,
		MaxConnections: conns,
		Headers:        req.Headers,
		ExpectedHash:   strings.TrimSpace(req.ExpectedHash),
		Status:         StatusQueued,
		CreatedAt:      time.Now().Unix(),
		client:         m.client,
		globalLimiter:  m.globalLimiter,
		onStateChange: func() {
			m.saveState()
			m.scheduleQueue()
		},
	}

	autoStart := true
	if req.AutoStart != nil {
		autoStart = *req.AutoStart
	}
	if !autoStart {
		task.Status = StatusPaused
	}

	m.tasks[id] = task
	m.taskOrder = append([]string{id}, m.taskOrder...)
	m.mu.Unlock()

	if req.Source != "" {
		task.addLog("info", "[브라우저연동] 요청 수신 ("+req.Source+"): "+task.FileName)
	} else {
		task.addLog("info", "[대기열등록] 다운로드 큐에 등록됨 ("+task.FileName+")")
	}
	m.saveState()
	m.scheduleQueue()

	snap := task.Snapshot()
	return &snap, nil
}

// PauseTask pauses a specific task.
func (m *Manager) PauseTask(id string) bool {
	m.mu.RLock()
	task, ok := m.tasks[id]
	m.mu.RUnlock()
	if !ok {
		return false
	}
	task.Pause()
	m.saveState()
	m.scheduleQueue()
	return true
}

// ResumeTask queues or immediately starts a paused/errored task.
func (m *Manager) ResumeTask(id string) bool {
	m.mu.RLock()
	task, ok := m.tasks[id]
	m.mu.RUnlock()
	if !ok {
		return false
	}
	task.mu.Lock()
	if task.Status == StatusCompleted || task.Status == StatusRunning || task.Status == StatusProbing {
		task.mu.Unlock()
		return true
	}
	task.Status = StatusQueued
	task.ErrorMessage = ""
	task.mu.Unlock()

	task.addLog("info", "[재개] 다운로드 이어받기 요청됨")
	m.scheduleQueue()
	return true
}

// DeleteTask stops and removes a task, optionally deleting the downloaded file from disk.
func (m *Manager) DeleteTask(id string, deleteFiles bool) bool {
	m.mu.Lock()
	task, ok := m.tasks[id]
	if !ok {
		m.mu.Unlock()
		return false
	}
	delete(m.tasks, id)
	newOrder := make([]string, 0, len(m.taskOrder))
	for _, tid := range m.taskOrder {
		if tid != id {
			newOrder = append(newOrder, tid)
		}
	}
	m.taskOrder = newOrder
	m.mu.Unlock()

	task.Pause()
	_ = os.Remove(task.PartPath)
	_ = os.Remove(task.MetaPath)
	if deleteFiles && task.TargetPath != "" {
		_ = os.Remove(task.TargetPath)
	}

	m.saveState()
	m.scheduleQueue()
	return true
}

// PauseAll pauses all running/queued tasks.
func (m *Manager) PauseAll() {
	m.mu.RLock()
	list := make([]*DownloadTask, 0, len(m.tasks))
	for _, t := range m.tasks {
		list = append(list, t)
	}
	m.mu.RUnlock()

	for _, t := range list {
		st := t.Snapshot().Status
		if st == StatusRunning || st == StatusProbing || st == StatusQueued {
			t.Pause()
		}
	}
	m.saveState()
}

// ResumeAll resumes all paused or errored tasks.
func (m *Manager) ResumeAll() {
	m.mu.RLock()
	list := make([]*DownloadTask, 0, len(m.tasks))
	for _, t := range m.tasks {
		list = append(list, t)
	}
	m.mu.RUnlock()

	for _, t := range list {
		t.mu.Lock()
		if t.Status == StatusPaused || t.Status == StatusError {
			t.Status = StatusQueued
			t.ErrorMessage = ""
		}
		t.mu.Unlock()
	}
	m.scheduleQueue()
	m.saveState()
}

// ClearCompleted removes all completed tasks from the list (leaving files intact on disk).
func (m *Manager) ClearCompleted() int {
	m.mu.Lock()
	removed := 0
	newOrder := make([]string, 0, len(m.taskOrder))
	for _, id := range m.taskOrder {
		t, ok := m.tasks[id]
		if !ok {
			continue
		}
		if t.Snapshot().Status == StatusCompleted {
			delete(m.tasks, id)
			removed++
		} else {
			newOrder = append(newOrder, id)
		}
	}
	m.taskOrder = newOrder
	m.mu.Unlock()
	if removed > 0 {
		m.saveState()
	}
	return removed
}

// GetTask returns a single task snapshot.
func (m *Manager) GetTask(id string) (*TaskSnapshot, bool) {
	m.mu.RLock()
	t, ok := m.tasks[id]
	m.mu.RUnlock()
	if !ok {
		return nil, false
	}
	snap := t.Snapshot()
	return &snap, true
}

// GlobalSnapshot compiles the entire application telemetry snapshot.
func (m *Manager) GlobalSnapshot() GlobalStateSnapshot {
	m.mu.RLock()
	order := make([]string, len(m.taskOrder))
	copy(order, m.taskOrder)
	taskMap := make(map[string]*DownloadTask, len(m.tasks))
	for k, v := range m.tasks {
		taskMap[k] = v
	}
	settings := m.settings
	hist := make([]int64, len(m.globalSpeedHistory))
	copy(hist, m.globalSpeedHistory)
	peak := m.peakGlobalSpeed
	clip := m.detectedClipboard
	bCnt := m.browserInterceptCnt
	protoReg := m.protocolRegistered
	extDir := m.extensionDir
	m.mu.RUnlock()

	snapshots := make([]TaskSnapshot, 0, len(order))
	var totalSpeed int64
	activeTasks := 0
	activeWorkers := 0

	for _, id := range order {
		if t, ok := taskMap[id]; ok {
			s := t.Snapshot()
			snapshots = append(snapshots, s)
			if s.Status == StatusRunning || s.Status == StatusProbing || s.Status == StatusVerifying {
				activeTasks++
				totalSpeed += s.SpeedBps
				activeWorkers += s.ActiveWorkers
			}
		}
	}

	return GlobalStateSnapshot{
		TotalSpeedBps:       totalSpeed,
		PeakGlobalSpeedBps:  peak,
		ActiveTasksCount:    activeTasks,
		ActiveWorkersCount:  activeWorkers,
		GlobalSpeedHistory:  hist,
		DetectedClipboard:   clip,
		BrowserInterceptCnt: bCnt,
		ProtocolRegistered:  protoReg,
		ExtensionDir:        extDir,
		Hardware:            platform.GetHardwareProfile(),
		DiskSpace:           platform.GetDiskSpace(settings.DefaultSaveDir),
		Settings:            settings,
		Tasks:               snapshots,
	}
}

func (m *Manager) scheduleQueue() {
	m.mu.RLock()
	maxConc := m.settings.MaxConcurrentTasks
	if maxConc <= 0 {
		maxConc = 4
	}
	running := 0
	var queued []*DownloadTask

	for i := len(m.taskOrder) - 1; i >= 0; i-- {
		id := m.taskOrder[i]
		t, ok := m.tasks[id]
		if !ok {
			continue
		}
		t.mu.RLock()
		st := t.Status
		t.mu.RUnlock()
		if st == StatusRunning || st == StatusProbing || st == StatusVerifying {
			running++
		} else if st == StatusQueued {
			queued = append(queued, t)
		}
	}
	m.mu.RUnlock()

	for _, qt := range queued {
		if running >= maxConc {
			break
		}
		qt.Start()
		running++
	}
}

func (m *Manager) backgroundLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	m.lastClipboard = strings.TrimSpace(platform.ReadClipboardText())

	for range ticker.C {
		m.mu.RLock()
		watchClip := m.settings.ClipboardMonitor
		m.mu.RUnlock()

		var currentTotalSpeed int64
		var activeRunning int
		m.mu.RLock()
		lang := m.settings.Language
		for _, t := range m.tasks {
			spd := t.speedAtomic.Load()
			currentTotalSpeed += spd
			t.mu.RLock()
			if t.Status == StatusRunning || t.Status == StatusProbing || t.Status == StatusVerifying {
				activeRunning++
			}
			t.mu.RUnlock()
		}
		m.mu.RUnlock()

		if activeRunning > 0 {
			mbps := float64(currentTotalSpeed) / (1024.0 * 1024.0)
			if lang == "en" {
				platform.UpdateSystemTrayTooltip(fmt.Sprintf("Loadin — %.2f MB/s (%d active)", mbps, activeRunning))
			} else if lang == "ja" {
				platform.UpdateSystemTrayTooltip(fmt.Sprintf("Loadin — %.2f MB/s (%d件進行中)", mbps, activeRunning))
			} else {
				platform.UpdateSystemTrayTooltip(fmt.Sprintf("Loadin — %.2f MB/s (%d개 다운로드 중)", mbps, activeRunning))
			}
		} else {
			if lang == "en" {
				platform.UpdateSystemTrayTooltip("Loadin — Download Manager (Running)")
			} else if lang == "ja" {
				platform.UpdateSystemTrayTooltip("Loadin — ダウンロードマネージャー (実行中)")
			} else {
				platform.UpdateSystemTrayTooltip("Loadin — 다운로드 매니저 (실행 중)")
			}
		}

		m.mu.Lock()
		m.globalSpeedHistory = append(m.globalSpeedHistory, currentTotalSpeed)
		if len(m.globalSpeedHistory) > 48 {
			m.globalSpeedHistory = m.globalSpeedHistory[len(m.globalSpeedHistory)-48:]
		}
		if currentTotalSpeed > m.peakGlobalSpeed {
			m.peakGlobalSpeed = currentTotalSpeed
		}
		m.mu.Unlock()

		if watchClip {
			clip := strings.TrimSpace(platform.ReadClipboardText())
			if clip != "" && clip != m.lastClipboard {
				m.lastClipboard = clip
				if IsLikelyDownloadURL(clip) {
					m.mu.Lock()
					alreadyExists := false
					for _, t := range m.tasks {
						if t.URL == clip {
							alreadyExists = true
							break
						}
					}
					if !alreadyExists {
						m.detectedClipboard = clip
					}
					m.mu.Unlock()
				}
			}
		}
	}
}

type persistedAppState struct {
	Settings Settings            `json:"settings"`
	Tasks    []PersistedTaskMeta `json:"tasks"`
}

func (m *Manager) saveState() {
	m.mu.RLock()
	statePath := m.stateFilePath
	settings := m.settings
	order := make([]string, len(m.taskOrder))
	copy(order, m.taskOrder)
	tasks := make([]*DownloadTask, 0, len(order))
	for _, id := range order {
		if t, ok := m.tasks[id]; ok {
			tasks = append(tasks, t)
		}
	}
	m.mu.RUnlock()

	metas := make([]PersistedTaskMeta, 0, len(tasks))
	for _, t := range tasks {
		metas = append(metas, t.ExportMeta())
	}

	payload := persistedAppState{
		Settings: settings,
		Tasks:    metas,
	}
	if data, err := json.MarshalIndent(payload, "", "  "); err == nil {
		_ = os.WriteFile(statePath, data, 0644)
	}
}

func (m *Manager) loadState() {
	data, err := os.ReadFile(m.stateFilePath)
	if err != nil {
		return
	}
	var payload persistedAppState
	if err := json.Unmarshal(data, &payload); err != nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if payload.Settings.DefaultSaveDir != "" {
		m.settings = payload.Settings
		if m.settings.Theme == "white" || m.settings.Theme == "light-frost" {
			m.settings.Theme = "white"
		} else {
			m.settings.Theme = "black"
		}
		if m.settings.Language == "" {
			m.settings.Language = "ko"
		}
		m.settings.RunOnStartup = platform.IsStartupAutoRunEnabled()
		m.globalLimiter.SetLimit(m.settings.GlobalSpeedLimitBps)
		go platform.StyleMainWindowGlass(m.settings.Theme)
	}

	sort.SliceStable(payload.Tasks, func(i, j int) bool {
		return payload.Tasks[i].CreatedAt > payload.Tasks[j].CreatedAt
	})

	for _, meta := range payload.Tasks {
		t := &DownloadTask{
			ID:             meta.ID,
			URL:            meta.URL,
			FinalURL:       meta.FinalURL,
			FileName:       meta.FileName,
			SaveDir:        meta.SaveDir,
			TargetPath:     meta.TargetPath,
			PartPath:       meta.PartPath,
			MetaPath:       meta.TargetPath + ".aether.meta",
			Category:       meta.Category,
			TotalSize:      meta.TotalSize,
			AcceptRanges:   meta.AcceptRanges,
			MaxConnections: meta.MaxConnections,
			Headers:        meta.Headers,
			ExpectedHash:   meta.ExpectedHash,
			Checksum:       meta.Checksum,
			ContentType:    meta.ContentType,
			Server:         meta.Server,
			ETag:           meta.ETag,
			Status:         meta.Status,
			ErrorMessage:   meta.ErrorMessage,
			CreatedAt:      meta.CreatedAt,
			CompletedAt:    meta.CompletedAt,
			DynamicSplits:  meta.DynamicSplits,
			logs:           meta.Logs,
			client:         m.client,
			globalLimiter:  m.globalLimiter,
			onStateChange: func() {
				m.saveState()
				m.scheduleQueue()
			},
		}
		t.downloadedAtomic.Store(meta.Downloaded)
		t.peakSpeedAtomic.Store(meta.PeakSpeedBps)

		for _, s := range meta.Segments {
			seg := &SegmentState{
				ID:         s.ID,
				Start:      s.Start,
				Current:    s.Current,
				End:        s.End,
				Status:     s.Status,
				Retries:    s.Retries,
				Generation: s.Generation,
			}
			seg.currentAtomic.Store(s.Current)
			seg.endAtomic.Store(s.End)
			t.segments = append(t.segments, seg)
		}

		m.tasks[t.ID] = t
		m.taskOrder = append(m.taskOrder, t.ID)
	}
}

func generateTaskID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
