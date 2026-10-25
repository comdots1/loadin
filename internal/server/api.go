package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aetherdl/aetherdl/internal/engine"
	"github.com/aetherdl/aetherdl/internal/platform"
)

//go:embed ui/*
var uiFS embed.FS

// Server exposes the embedded Glassmorphism UI, Browser Connector Bridge, and high-frequency telemetry API.
type Server struct {
	manager *engine.Manager
	mux     *http.ServeMux
	addr    string
}

// NewServer creates and wires up the HTTP + SSE server.
func NewServer(mgr *engine.Manager, addr string) *Server {
	s := &Server{
		manager: mgr,
		mux:     http.NewServeMux(),
		addr:    addr,
	}
	s.registerRoutes()
	return s
}

// Handler returns the root HTTP handler with CORS middleware for browser extension / bookmarklet calls.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		s.mux.ServeHTTP(w, r)
	})
}

func (s *Server) registerRoutes() {
	subUI, err := fs.Sub(uiFS, "ui")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(subUI))

	s.mux.HandleFunc("/api/state", s.handleState)
	s.mux.HandleFunc("/api/stream", s.handleSSEStream)
	s.mux.HandleFunc("/api/hardware", s.handleHardware)
	s.mux.HandleFunc("/api/probe", s.handleProbe)
	s.mux.HandleFunc("/api/tasks", s.handleTasks)
	s.mux.HandleFunc("/api/tasks/item", s.handleTaskItem)
	s.mux.HandleFunc("/api/tasks/batch", s.handleBatchTasks)
	s.mux.HandleFunc("/api/tasks/action", s.handleTaskAction)
	s.mux.HandleFunc("/api/sniff", s.handleSniff)
	s.mux.HandleFunc("/api/settings", s.handleSettings)
	s.mux.HandleFunc("/api/dialog/folder", s.handleFolderDialog)
	s.mux.HandleFunc("/api/folder/open", s.handleOpenFolderPath)
	s.mux.HandleFunc("/api/folder/default", s.handleDefaultFolderPath)
	s.mux.HandleFunc("/api/app/quit", s.handleAppQuit)
	s.mux.HandleFunc("/api/clipboard/dismiss", s.handleDismissClipboard)
	s.mux.HandleFunc("/api/sandbox-stream/", s.handleSandboxStream)

	// Browser Extension & Bookmarklet Integration Endpoints
	s.mux.HandleFunc("/api/browser/ping", s.handleBrowserPing)
	s.mux.HandleFunc("/api/browser/intercept", s.handleBrowserIntercept)
	s.mux.HandleFunc("/api/browser/add", s.handleBrowserBookmarkletAdd)
	s.mux.HandleFunc("/api/browser/register-protocol", s.handleRegisterProtocol)
	s.mux.HandleFunc("/api/browser/open-extension-folder", s.handleOpenExtensionFolder)
	s.mux.HandleFunc("/api/hover/launch", s.handleLaunchHoverWidget)
	s.mux.HandleFunc("/api/window/control", s.handleWindowControl)

	s.mux.Handle("/", fileServer)
}

func (s *Server) enginePort() int {
	port := 17890
	if idx := strings.LastIndex(s.addr, ":"); idx != -1 {
		if p, err := strconv.Atoi(s.addr[idx+1:]); err == nil && p > 0 {
			port = p
		}
	}
	return port
}

func (s *Server) handleLaunchHoverWidget(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL      string `json:"url"`
		Referer  string `json:"referer"`
		FileName string `json:"fileName"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	_, err := platform.ShowAddDownloadDialog(s.enginePort(), req.URL, req.Referer, "", "", req.FileName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleWindowControl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false})
		return
	}
	ok := platform.ControlMainWindow(req.Action)
	writeJSON(w, http.StatusOK, map[string]any{"ok": ok})
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.manager.GlobalSnapshot())
}

func (s *Server) handleHardware(w http.ResponseWriter, r *http.Request) {
	var sizeBytes int64
	if q := r.URL.Query().Get("size"); q != "" {
		sizeBytes, _ = strconv.ParseInt(q, 10, 64)
	}
	hw := platform.GetHardwareProfile()
	recThreads, note := platform.RecommendThreadsForFileSize(sizeBytes)
	writeJSON(w, http.StatusOK, map[string]any{
		"hardware":           hw,
		"recommendedThreads": recThreads,
		"recommendationNote": note,
	})
}

func (s *Server) handleBrowserPing(w http.ResponseWriter, r *http.Request) {
	settings := s.manager.GetSettings()
	writeJSON(w, http.StatusOK, map[string]any{
		"online":           true,
		"version":          "2.1.0",
		"browserIntercept": settings.BrowserIntercept,
	})
}

func (s *Server) handleBrowserIntercept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		URL       string `json:"url"`
		FileName  string `json:"fileName"`
		Referer   string `json:"referer"`
		Cookie    string `json:"cookie"`
		UserAgent string `json:"userAgent"`
		TotalSize int64  `json:"totalSize"`
		Source    string `json:"source"`
		Direct    bool   `json:"direct"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "유효하지 않은 URL입니다."})
		return
	}

	s.manager.IncrementBrowserIntercept()

	// By default, when a download link is clicked in the browser, pop up the on-demand
	// native "새 다운로드 추가" dialog so the user can confirm/customize file name, folder, and threads.
	if !req.Direct {
		if _, err := platform.ShowAddDownloadDialog(
			s.enginePort(),
			req.URL,
			req.Referer,
			req.Cookie,
			req.UserAgent,
			req.FileName,
		); err == nil {
			writeJSON(w, http.StatusOK, map[string]any{
				"ok":           true,
				"dialogOpened": true,
				"url":          req.URL,
			})
			return
		}
	}

	headers := make(map[string]string)
	if req.Referer != "" {
		headers["Referer"] = req.Referer
	}
	if req.Cookie != "" {
		headers["Cookie"] = req.Cookie
	}
	if req.UserAgent != "" {
		headers["User-Agent"] = req.UserAgent
	}

	recThreads, _ := platform.RecommendThreadsForFileSize(req.TotalSize)
	source := req.Source
	if source == "" {
		source = "브라우저 자동 연동"
	}

	snap, err := s.manager.AddTask(engine.CreateTaskRequest{
		URL:            req.URL,
		FileName:       req.FileName,
		MaxConnections: recThreads,
		Headers:        headers,
		TotalSize:      req.TotalSize,
		Source:         source,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, snap)
}

func (s *Server) handleBrowserBookmarkletAdd(w http.ResponseWriter, r *http.Request) {
	targetURL := strings.TrimSpace(r.URL.Query().Get("url"))
	referer := strings.TrimSpace(r.URL.Query().Get("referer"))
	if targetURL != "" {
		s.manager.IncrementBrowserIntercept()
		if _, err := platform.ShowAddDownloadDialog(s.enginePort(), targetURL, referer, "", "", ""); err != nil {
			headers := map[string]string{}
			if referer != "" {
				headers["Referer"] = referer
			}
			_, _ = s.manager.AddTask(engine.CreateTaskRequest{
				URL:     targetURL,
				Headers: headers,
				Source:  "북마클릿 연동",
			})
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, `<!DOCTYPE html><html lang="ko"><head><meta charset="utf-8"><title>Loadin</title>
<style>body{background:#161618;color:#F5F5F7;font-family:-apple-system,BlinkMacSystemFont,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}
.box{padding:24px 32px;border-radius:14px;background:#1C1C1E;border:1px solid rgba(255,255,255,0.1);text-align:center}</style></head>
<body><div class="box"><h3>Loadin 다운로드 추가 창을 열었습니다</h3><p style="color:#98989D;font-size:13px">잠시 후 이 탭이 자동으로 닫힙니다...</p></div>
<script>setTimeout(()=>window.close(),700);</script></body></html>`)
}

func (s *Server) handleRegisterProtocol(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := platform.RegisterBrowserProtocol()
	reg := s.manager.RefreshProtocolStatus()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"ok":         false,
			"registered": reg,
			"error":      err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"registered": reg,
	})
}

func (s *Server) handleOpenExtensionFolder(w http.ResponseWriter, r *http.Request) {
	extDir := s.manager.ExtensionDir()
	_ = platform.OpenFileDefault(extDir)
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":   true,
		"path": extDir,
	})
}

// handleSSEStream streams live telemetry snapshots at 5Hz (every 200ms) via Server-Sent Events.
func (s *Server) handleSSEStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sendSnap := func() error {
		snap := s.manager.GlobalSnapshot()
		data, err := json.Marshal(snap)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "data: %s\n\n", data)
		if err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	_ = sendSnap()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if err := sendSnap(); err != nil {
				return
			}
		}
	}
}

func (s *Server) handleProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "유효한 URL을 입력해주세요."})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()

	res, err := engine.ProbeURL(ctx, s.manager.HTTPClient(), strings.TrimSpace(req.URL), req.Headers)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req engine.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "유효한 URL을 입력해주세요."})
		return
	}

	snap, err := s.manager.AddTask(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, snap)
}

func (s *Server) handleTaskItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}
	snap, ok := s.manager.GetTask(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, snap)
}

func (s *Server) handleBatchTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		URLs           []string          `json:"urls"`
		Pattern        string            `json:"pattern"`
		SaveDir        string            `json:"saveDir"`
		MaxConnections int               `json:"maxConnections"`
		Headers        map[string]string `json:"headers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "잘못된 요청 형식입니다."})
		return
	}

	var allURLs []string
	if strings.TrimSpace(req.Pattern) != "" {
		allURLs = append(allURLs, engine.ExpandBatchPattern(req.Pattern)...)
	}
	for _, u := range req.URLs {
		u = strings.TrimSpace(u)
		if u != "" {
			allURLs = append(allURLs, engine.ExpandBatchPattern(u)...)
		}
	}

	added := 0
	for _, u := range allURLs {
		if _, err := s.manager.AddTask(engine.CreateTaskRequest{
			URL:            u,
			SaveDir:        req.SaveDir,
			MaxConnections: req.MaxConnections,
			Headers:        req.Headers,
		}); err == nil {
			added++
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"added": added,
	})
}

func (s *Server) handleTaskAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID          string `json:"id"`
		Action      string `json:"action"`
		DeleteFiles bool   `json:"deleteFiles"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	switch req.Action {
	case "pause":
		s.manager.PauseTask(req.ID)
	case "resume":
		s.manager.ResumeTask(req.ID)
	case "delete":
		s.manager.DeleteTask(req.ID, req.DeleteFiles)
	case "pause-all":
		s.manager.PauseAll()
	case "resume-all":
		s.manager.ResumeAll()
	case "clear-completed":
		s.manager.ClearCompleted()
	case "open":
		if task, ok := s.manager.GetTask(req.ID); ok && task.TargetPath != "" {
			_ = platform.OpenFileDefault(task.TargetPath)
		}
	case "reveal":
		if task, ok := s.manager.GetTask(req.ID); ok && task.TargetPath != "" {
			_ = platform.ShowInExplorer(task.TargetPath)
		}
	case "open-default-dir":
		_ = platform.OpenFileDefault(s.manager.GetSettings().DefaultSaveDir)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unknown action"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleSniff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "스캔할 웹페이지 URL을 입력해주세요."})
		return
	}

	links, err := engine.SniffWebpageLinks(r.Context(), s.manager.HTTPClient(), strings.TrimSpace(req.URL))
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"links": links})
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, s.manager.GetSettings())
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req engine.Settings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid settings"})
		return
	}
	updated := s.manager.UpdateSettings(req)
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleFolderDialog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		InitialDir    string `json:"initialDir"`
		SaveAsDefault bool   `json:"saveAsDefault"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	current := strings.TrimSpace(req.InitialDir)
	if current == "" {
		current = s.manager.GetSettings().DefaultSaveDir
	}
	picked, err := platform.PickFolderDialog(current)
	if err != nil || picked == "" {
		writeJSON(w, http.StatusOK, map[string]string{"path": ""})
		return
	}
	if req.SaveAsDefault {
		curSettings := s.manager.GetSettings()
		curSettings.DefaultSaveDir = picked
		s.manager.UpdateSettings(curSettings)
	}
	writeJSON(w, http.StatusOK, map[string]string{"path": picked})
}

func (s *Server) handleOpenFolderPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	target := strings.TrimSpace(req.Path)
	if target == "" {
		target = s.manager.GetSettings().DefaultSaveDir
	}
	_ = os.MkdirAll(target, 0755)
	_ = platform.OpenFileDefault(target)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "path": target})
}

func (s *Server) handleDefaultFolderPath(w http.ResponseWriter, r *http.Request) {
	defDir := platform.DefaultDownloadsDir()
	writeJSON(w, http.StatusOK, map[string]string{"path": defDir})
}

func (s *Server) handleAppQuit(w http.ResponseWriter, r *http.Request) {
	s.manager.PauseAll()
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	go func() {
		time.Sleep(120 * time.Millisecond)
		platform.CloseAllAppWindows()
		os.Exit(0)
	}()
}

func (s *Server) handleDismissClipboard(w http.ResponseWriter, r *http.Request) {
	s.manager.DismissClipboard()
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleSandboxStream(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/api/sandbox-stream/")
	if fileName == "" {
		fileName = "Loadin_Benchmark_64MB.bin"
	}

	sizeMB := int64(64)
	if q := r.URL.Query().Get("sizeMB"); q != "" {
		if parsed, err := strconv.ParseInt(q, 10, 64); err == nil && parsed > 0 && parsed <= 2048 {
			sizeMB = parsed
		}
	}
	totalSize := sizeMB * 1024 * 1024

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	w.Header().Set("Server", "Loadin-Stream/2.1")
	w.Header().Set("ETag", fmt.Sprintf(`"loadin-bench-%d"`, totalSize))

	if r.Method == http.MethodHead {
		w.Header().Set("Content-Length", strconv.FormatInt(totalSize, 10))
		w.WriteHeader(http.StatusOK)
		return
	}

	var start, end int64 = 0, totalSize - 1
	isPartial := false

	if rangeHdr := r.Header.Get("Range"); strings.HasPrefix(rangeHdr, "bytes=") {
		isPartial = true
		spec := strings.TrimPrefix(rangeHdr, "bytes=")
		parts := strings.SplitN(spec, "-", 2)
		if len(parts) == 2 {
			if sVal, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
				start = sVal
			}
			if parts[1] != "" {
				if eVal, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					end = eVal
				}
			}
		}
	}

	if start < 0 || start >= totalSize || end < start {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	if end >= totalSize {
		end = totalSize - 1
	}

	contentLen := (end - start) + 1
	if isPartial {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
		w.Header().Set("Content-Length", strconv.FormatInt(contentLen, 10))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.Header().Set("Content-Length", strconv.FormatInt(totalSize, 10))
		w.WriteHeader(http.StatusOK)
	}

	if contentLen <= 16 {
		buf := make([]byte, contentLen)
		for i := int64(0); i < contentLen; i++ {
			buf[i] = byte((start + i) & 0xFF)
		}
		_, _ = w.Write(buf)
		return
	}

	chunk := make([]byte, 64*1024)
	for i := range chunk {
		chunk[i] = byte(i & 0xFF)
	}

	streamVariant := (start / (1024 * 1024)) % 3
	sleepDur := 12 * time.Millisecond
	if streamVariant == 0 {
		sleepDur = 6 * time.Millisecond
	} else if streamVariant == 2 {
		sleepDur = 22 * time.Millisecond
	}

	remaining := contentLen
	for remaining > 0 {
		if r.Context().Err() != nil {
			return
		}
		toWrite := int64(len(chunk))
		if remaining < toWrite {
			toWrite = remaining
		}
		_, err := w.Write(chunk[:toWrite])
		if err != nil {
			return
		}
		remaining -= toWrite
		time.Sleep(sleepDur)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
