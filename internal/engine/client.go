package engine

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aetherdl/aetherdl/internal/platform"
)

const (
	// BufferSize is 128KB per read syscall — optimal sweet spot for NVMe + TCP window scaling
	BufferSize = 128 * 1024
	// DefaultUserAgent mimics a modern desktop browser so CDNs never throttle or block requests
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 AetherDL/2.0"
)

// BufferPool provides zero-allocation reusable 128KB byte slices for segment workers.
var BufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, BufferSize)
		return &buf
	},
}

// NewOptimizedHTTPClient builds a high-throughput HTTP client with IPv4-priority Happy Eyeballs
// and large 4MB TCP receive buffers.
func NewOptimizedHTTPClient() *http.Client {
	baseDialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		// Try IPv4 first with fast timeout to avoid 60s blackhole stalls on broken IPv6 routes
		ipv4Ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
		conn, err := baseDialer.DialContext(ipv4Ctx, "tcp4", addr)
		cancel()
		if err != nil {
			// Fallback to standard dual-stack
			conn, err = baseDialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
		}
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			_ = tcpConn.SetNoDelay(true)
			_ = tcpConn.SetReadBuffer(4 * 1024 * 1024)
		}
		return conn, nil
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          256,
		MaxIdleConnsPerHost:   64,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		DisableCompression:    true, // Critical for accurate byte-range offsets
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   0, // Streaming body reads have no fixed total deadline
	}
}

// ProbeResult contains metadata discovered when inspecting a target URL.
type ProbeResult struct {
	FinalURL           string `json:"finalUrl"`
	FileName           string `json:"fileName"`
	TotalSize          int64  `json:"totalSize"`
	AcceptRanges       bool   `json:"acceptRanges"`
	ContentType        string `json:"contentType"`
	ETag               string `json:"etag"`
	LastModified       string `json:"lastModified"`
	Server             string `json:"server"`
	Category           string `json:"category"`
	StatusMessage      string `json:"statusMessage"`
	RecommendedThreads int    `json:"recommendedThreads"`
	RecommendationNote string `json:"recommendationNote"`
}

// ProbeURL inspects a remote resource using a 1-byte Range GET (with HEAD fallback)
// to reliably detect exact file size, Range support, redirects, and filename.
func ProbeURL(ctx context.Context, client *http.Client, rawURL string, headers map[string]string) (*ProbeResult, error) {
	probeCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(probeCtx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("Range", "bytes=0-0")
	req.Header.Set("Accept", "*/*")
	for k, v := range headers {
		if strings.TrimSpace(k) != "" && strings.TrimSpace(v) != "" {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Retry without Range header via HEAD in case server rejects Range: bytes=0-0
		reqHead, err2 := http.NewRequestWithContext(probeCtx, http.MethodHead, rawURL, nil)
		if err2 == nil {
			reqHead.Header.Set("User-Agent", DefaultUserAgent)
			for k, v := range headers {
				if strings.TrimSpace(k) != "" {
					reqHead.Header.Set(k, v)
				}
			}
			if respHead, errH := client.Do(reqHead); errH == nil {
				defer respHead.Body.Close()
				if respHead.StatusCode < 400 {
					resp = respHead
				}
			}
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP %d (%s)", resp.StatusCode, resp.Status)
		}
	}

	finalURL := resp.Request.URL.String()
	var totalSize int64 = -1
	acceptRanges := false

	// Check Content-Range: bytes 0-0/1234567
	if cr := resp.Header.Get("Content-Range"); cr != "" {
		acceptRanges = true
		if slashIdx := strings.LastIndex(cr, "/"); slashIdx != -1 {
			sizeStr := strings.TrimSpace(cr[slashIdx+1:])
			if sizeStr != "*" {
				if parsed, err := strconv.ParseInt(sizeStr, 10, 64); err == nil && parsed > 0 {
					totalSize = parsed
				}
			}
		}
	}

	if resp.StatusCode == http.StatusPartialContent {
		acceptRanges = true
	}
	if strings.EqualFold(resp.Header.Get("Accept-Ranges"), "bytes") {
		acceptRanges = true
	}

	if totalSize <= 0 && resp.ContentLength > 0 {
		// Note: if status is 206 and Content-Range wasn't parsed, ContentLength is 1
		if resp.StatusCode != http.StatusPartialContent {
			totalSize = resp.ContentLength
		}
	}

	// If server returned 206 with Content-Range: bytes 0-0/* or missing total, do a quick HEAD check
	if totalSize <= 0 {
		if reqHead, errH := http.NewRequestWithContext(probeCtx, http.MethodHead, finalURL, nil); errH == nil {
			reqHead.Header.Set("User-Agent", DefaultUserAgent)
			if respHead, errDo := client.Do(reqHead); errDo == nil {
				if respHead.ContentLength > 0 {
					totalSize = respHead.ContentLength
				}
				respHead.Body.Close()
			}
		}
	}

	fileName := extractFileName(resp, finalURL, rawURL)
	contentType := resp.Header.Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = strings.TrimSpace(contentType[:idx])
	}

	// Ensure extension if Content-Type provides one and filename lacks extension
	if path.Ext(fileName) == "" && contentType != "" {
		if ext := guessExtFromMime(contentType); ext != "" {
			fileName += ext
		}
	}

	cat := CategorizeFile(fileName, contentType)
	recThreads, recNote := platform.RecommendThreadsForFileSize(totalSize)
	if !acceptRanges {
		recThreads = 1
		recNote = "단일 스트림 서버이므로 1 스레드로 연결됩니다."
	}

	return &ProbeResult{
		FinalURL:           finalURL,
		FileName:           SanitizeFileName(fileName),
		TotalSize:          totalSize,
		AcceptRanges:       acceptRanges && totalSize > 0,
		ContentType:        contentType,
		ETag:               resp.Header.Get("ETag"),
		LastModified:       resp.Header.Get("Last-Modified"),
		Server:             resp.Header.Get("Server"),
		Category:           cat,
		StatusMessage:      resp.Status,
		RecommendedThreads: recThreads,
		RecommendationNote: recNote,
	}, nil
}

func extractFileName(resp *http.Response, finalURL, origURL string) string {
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if _, params, err := mime.ParseMediaType(cd); err == nil {
			if fn, ok := params["filename*"]; ok && fn != "" {
				// Strip UTF-8'' prefix if still present
				fn = strings.TrimPrefix(fn, "UTF-8''")
				fn = strings.TrimPrefix(fn, "utf-8''")
				if unescaped, err := url.QueryUnescape(fn); err == nil && unescaped != "" {
					return unescaped
				}
				return fn
			}
			if fn, ok := params["filename"]; ok && fn != "" {
				if unescaped, err := url.QueryUnescape(fn); err == nil && unescaped != "" {
					return unescaped
				}
				return fn
			}
		}
	}

	for _, candidate := range []string{finalURL, origURL} {
		if u, err := url.Parse(candidate); err == nil {
			// Check query params like ?filename= or ?file= or ?response-content-disposition=
			q := u.Query()
			for _, key := range []string{"filename", "file", "name", "title"} {
				if val := q.Get(key); val != "" {
					return val
				}
			}
			base := path.Base(u.Path)
			if base != "" && base != "/" && base != "." {
				if unescaped, err := url.PathUnescape(base); err == nil && unescaped != "" {
					return unescaped
				}
				return base
			}
		}
	}
	return fmt.Sprintf("download_%d.bin", time.Now().Unix())
}

// SanitizeFileName removes illegal Windows filesystem characters.
func SanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	replacer := strings.NewReplacer(
		"<", "_", ">", "_", ":", "_", "\"", "'",
		"/", "_", "\\", "_", "|", "_", "?", "", "*", "_",
	)
	clean := replacer.Replace(name)
	if clean == "" || clean == "." || clean == ".." {
		return fmt.Sprintf("file_%d.bin", time.Now().Unix())
	}
	return clean
}

// CategorizeFile maps a filename or MIME type into an IDM-style smart category.
func CategorizeFile(fileName, contentType string) string {
	ext := strings.ToLower(path.Ext(fileName))
	switch ext {
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".zst", ".tgz", ".cab":
		return "compressed"
	case ".exe", ".msi", ".msix", ".dmg", ".pkg", ".deb", ".rpm", ".appimage", ".apk", ".bat", ".iso", ".img":
		return "programs"
	case ".mp4", ".mkv", ".webm", ".avi", ".mov", ".wmv", ".flv", ".m4v", ".ts", ".m3u8":
		return "videos"
	case ".mp3", ".flac", ".wav", ".aac", ".ogg", ".m4a", ".opus", ".wma":
		return "music"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md", ".csv", ".epub", ".json", ".xml":
		return "documents"
	case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".svg", ".bmp", ".ico", ".psd", ".ai", ".fig":
		return "images"
	}

	ct := strings.ToLower(contentType)
	switch {
	case strings.HasPrefix(ct, "video/"):
		return "videos"
	case strings.HasPrefix(ct, "audio/"):
		return "music"
	case strings.HasPrefix(ct, "image/"):
		return "images"
	case strings.Contains(ct, "zip") || strings.Contains(ct, "compressed") || strings.Contains(ct, "tar"):
		return "compressed"
	case strings.Contains(ct, "pdf") || strings.Contains(ct, "document") || strings.Contains(ct, "text/"):
		return "documents"
	case strings.Contains(ct, "executable") || strings.Contains(ct, "msdownload") || strings.Contains(ct, "octet-stream"):
		return "programs"
	}
	return "general"
}

func guessExtFromMime(contentType string) string {
	switch strings.ToLower(contentType) {
	case "application/zip":
		return ".zip"
	case "application/pdf":
		return ".pdf"
	case "video/mp4":
		return ".mp4"
	case "audio/mpeg":
		return ".mp3"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "text/html":
		return ".html"
	case "application/json":
		return ".json"
	}
	return ""
}
