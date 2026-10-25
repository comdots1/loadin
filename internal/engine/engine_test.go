package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestExpandBatchPattern(t *testing.T) {
	urls := ExpandBatchPattern("https://cdn.example.com/part_[01-04].zip")
	if len(urls) != 4 {
		t.Fatalf("expected 4 URLs, got %d", len(urls))
	}
	if urls[0] != "https://cdn.example.com/part_01.zip" || urls[3] != "https://cdn.example.com/part_04.zip" {
		t.Fatalf("unexpected expansion: %v", urls)
	}
}

type mockRangeRoundTripper struct {
	payload []byte
}

func (m *mockRangeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	total := len(m.payload)
	hdr := make(http.Header)
	hdr.Set("Accept-Ranges", "bytes")
	hdr.Set("Content-Type", "application/octet-stream")
	hdr.Set("Content-Disposition", `attachment; filename="test_payload.bin"`)
	hdr.Set("Server", "AetherDL-TestEdge/2.0")

	rangeHdr := req.Header.Get("Range")
	start, end := 0, total-1
	status := http.StatusOK

	if strings.HasPrefix(rangeHdr, "bytes=") {
		status = http.StatusPartialContent
		spec := strings.TrimPrefix(rangeHdr, "bytes=")
		parts := strings.SplitN(spec, "-", 2)
		start, _ = strconv.Atoi(parts[0])
		if parts[1] != "" {
			end, _ = strconv.Atoi(parts[1])
		}
		if end >= total {
			end = total - 1
		}
		hdr.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, total))
	}

	slice := m.payload[start : end+1]
	hdr.Set("Content-Length", strconv.Itoa(len(slice)))

	pr, pw := io.Pipe()
	go func(data []byte, startOffset int) {
		defer pw.Close()
		// Make second half stream gradually so worker #1 finishes early and steals half of worker #2!
		if startOffset > total/4 && len(data) > 64*1024 {
			chunk := 256 * 1024
			for off := 0; off < len(data); off += chunk {
				if req.Context().Err() != nil {
					return
				}
				endOff := off + chunk
				if endOff > len(data) {
					endOff = len(data)
				}
				if _, err := pw.Write(data[off:endOff]); err != nil {
					return
				}
				time.Sleep(12 * time.Millisecond)
			}
			return
		}
		_, _ = pw.Write(data)
	}(slice, start)

	return &http.Response{
		Status:        fmt.Sprintf("%d OK", status),
		StatusCode:    status,
		Header:        hdr,
		Body:          pr,
		ContentLength: int64(len(slice)),
		Request:       req,
	}, nil
}

func TestMultiSegmentDynamicSplitAndChecksum(t *testing.T) {
	// Generate a 6 MB deterministic payload so Dynamic Work-Stealing (>1MB remaining) triggers
	const payloadSize = 6 * 1024 * 1024
	payload := make([]byte, payloadSize)
	for i := range payload {
		payload[i] = byte((i * 31) & 0xFF)
	}
	sum := sha256.Sum256(payload)
	expectedSHA256 := hex.EncodeToString(sum[:])

	tmpDir, err := os.MkdirTemp("", "aetherdl_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	mockClient := &http.Client{
		Transport: &mockRangeRoundTripper{payload: payload},
	}

	task := &DownloadTask{
		ID:             "test-1",
		URL:            "https://cdn.aetherdl.dev/test_payload.bin",
		SaveDir:        tmpDir,
		MaxConnections: 2, // Start with 2 connections; worker #1 finishes fast and splits worker #2!
		ExpectedHash:   expectedSHA256,
		Status:         StatusQueued,
		client:         mockClient,
		globalLimiter:  NewRateLimiter(0),
	}

	task.Start()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		snap := task.Snapshot()
		if snap.Status == StatusCompleted {
			if snap.Checksum == nil || snap.Checksum.Verified == nil || !*snap.Checksum.Verified {
				t.Fatalf("expected verified SHA256 checksum, got %+v", snap.Checksum)
			}
			if snap.DynamicSplits < 1 {
				t.Fatalf("expected at least 1 dynamic work-stealing split, got 0")
			}
			t.Logf("✅ Dynamic Work-Stealing Splits: %d | Final Segments: %d | Verified SHA-256: %s",
				snap.DynamicSplits, len(snap.Segments), snap.Checksum.SHA256)
			return
		}
		if snap.Status == StatusError {
			t.Fatalf("task failed with error: %s", snap.ErrorMessage)
		}
		time.Sleep(40 * time.Millisecond)
	}
	t.Fatal("timeout waiting for download task to complete")
}

func TestEnsureBrowserExtensionFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "loadin_ext_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := EnsureBrowserExtensionFiles(tmpDir); err != nil {
		t.Fatalf("EnsureBrowserExtensionFiles failed: %v", err)
	}

	required := []string{
		"manifest.json",
		"background.js",
		"content.js",
		"popup.html",
		"popup.js",
		"firefox/manifest.json",
		"firefox/background.js",
		"firefox/content.js",
	}
	for _, rel := range required {
		full := tmpDir + string(os.PathSeparator) + rel
		info, err := os.Stat(full)
		if err != nil || info.Size() == 0 {
			t.Fatalf("expected non-empty extension file %s, err=%v", rel, err)
		}
	}
}
