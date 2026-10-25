package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aetherdl/aetherdl/internal/engine"
	"github.com/aetherdl/aetherdl/internal/platform"
	"github.com/aetherdl/aetherdl/internal/server"
)

const defaultPort = 17890

func main() {
	portFlag := flag.Int("port", defaultPort, "Local engine HTTP/SSE port")
	headlessFlag := flag.Bool("headless", false, "Run engine server without opening desktop window")
	addURLFlag := flag.String("url", "", "Immediately queue a URL for accelerated download")
	flag.Parse()

	rawTarget := strings.TrimSpace(*addURLFlag)
	if rawTarget == "" && flag.NArg() > 0 {
		rawTarget = strings.TrimSpace(flag.Arg(0))
	}
	targetURL := parseProtocolURL(rawTarget)

	addr := fmt.Sprintf("127.0.0.1:%d", *portFlag)
	appURL := fmt.Sprintf("http://%s", addr)

	// Single-instance check: if Loadin engine is already listening on defaultPort,
	// forward any CLI / protocol URL to it and bring up the desktop window.
	if isEngineRunning(appURL) {
		if targetURL != "" {
			_ = forwardURLToRunningInstance(appURL, targetURL)
			return
		}
		if !*headlessFlag {
			platform.ShowOrRestoreMainWindow(appURL)
		}
		return
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to bind local engine port: %v\n", err)
			os.Exit(1)
		}
		addr = listener.Addr().String()
		appURL = fmt.Sprintf("http://%s", addr)
	}

	mgr := engine.NewManager()

	srv := server.NewServer(mgr, addr)
	httpServer := &http.Server{
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		_ = httpServer.Serve(listener)
	}()

	fmt.Printf("Loadin Engine running at %s\n", appURL)

	// Register native Windows System Tray icon and right-click context menu
	platform.StartSystemTray(platform.TrayCallbacks{
		OnOpenMainWindow: func() {
			platform.ShowOrRestoreMainWindow(appURL)
		},
		OnNewDownload: func() {
			_, _ = platform.ShowAddDownloadDialog(*portFlag, "", "", "", "", "")
		},
		OnOpenSaveFolder: func() {
			saveDir := mgr.GetSettings().DefaultSaveDir
			_ = os.MkdirAll(saveDir, 0755)
			_ = platform.OpenFileDefault(saveDir)
		},
		OnResumeAll: func() {
			mgr.ResumeAll()
		},
		OnPauseAll: func() {
			mgr.PauseAll()
		},
		OnQuit: func() {
			mgr.PauseAll()
			platform.CloseAllAppWindows()
			os.Exit(0)
		},
		GetLanguage: func() string {
			return mgr.GetSettings().Language
		},
	})

	if targetURL != "" {
		_, _ = platform.ShowAddDownloadDialog(*portFlag, targetURL, "", "", "", "")
	} else if !*headlessFlag {
		_, _ = platform.LaunchDesktopAppWindow(appURL)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	mgr.PauseAll()
	platform.RemoveSystemTray()
}

// parseProtocolURL extracts the target download URL from a loadin:// URI scheme if present.
func parseProtocolURL(raw string) string {
	raw = strings.TrimSpace(strings.Trim(raw, `"`))
	if raw == "" {
		return ""
	}
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "loadin://") || strings.HasPrefix(lower, "aetherdl://") {
		if u, err := url.Parse(raw); err == nil {
			if qURL := u.Query().Get("url"); qURL != "" {
				return qURL
			}
		}
		trimmed := raw
		if idx := strings.Index(trimmed, "://"); idx != -1 {
			trimmed = trimmed[idx+3:]
		}
		trimmed = strings.TrimPrefix(trimmed, "download/?url=")
		trimmed = strings.TrimPrefix(trimmed, "download?url=")
		if unescaped, err := url.QueryUnescape(trimmed); err == nil && strings.HasPrefix(unescaped, "http") {
			return unescaped
		}
		if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
			return trimmed
		}
	}
	return raw
}

func isEngineRunning(baseURL string) bool {
	client := &http.Client{Timeout: 600 * time.Millisecond}
	resp, err := client.Get(baseURL + "/api/state")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func forwardURLToRunningInstance(baseURL, downloadURL string) error {
	payload, _ := json.Marshal(map[string]any{
		"url":    downloadURL,
		"source": "브라우저 프로토콜 연동",
	})
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(baseURL+"/api/browser/intercept", "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
