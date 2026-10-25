package engine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SniffedLink represents a downloadable link discovered on a webpage.
type SniffedLink struct {
	URL       string `json:"url"`
	FileName  string `json:"fileName"`
	Extension string `json:"extension"`
	Category  string `json:"category"`
	Label     string `json:"label"`
}

var (
	hrefRegex    = regexp.MustCompile(`(?i)(?:href|src|data-src)\s*=\s*["']([^"'#]+)["']`)
	rangePattern = regexp.MustCompile(`\[(\d+)-(\d+)\]`)
)

// ExpandBatchPattern expands an IDM-style numeric pattern such as:
// https://example.com/archive_part[01-05].zip -> 5 URLs
func ExpandBatchPattern(pattern string) []string {
	pattern = strings.TrimSpace(pattern)
	loc := rangePattern.FindStringSubmatchIndex(pattern)
	if loc == nil {
		return []string{pattern}
	}

	startStr := pattern[loc[2]:loc[3]]
	endStr := pattern[loc[4]:loc[5]]
	startNum, err1 := strconv.Atoi(startStr)
	endNum, err2 := strconv.Atoi(endStr)
	if err1 != nil || err2 != nil || endNum < startNum || (endNum-startNum) > 500 {
		return []string{pattern}
	}

	width := len(startStr)
	prefix := pattern[:loc[0]]
	suffix := pattern[loc[1]:]

	var results []string
	formatStr := fmt.Sprintf("%%s%%0%dd%%s", width)
	if width == 1 || !strings.HasPrefix(startStr, "0") {
		formatStr = "%s%d%s"
	}
	for i := startNum; i <= endNum; i++ {
		results = append(results, fmt.Sprintf(formatStr, prefix, i, suffix))
	}
	return results
}

// SniffWebpageLinks fetches a webpage and extracts all downloadable links and media assets.
func SniffWebpageLinks(ctx context.Context, client *http.Client, pageURL string) ([]SniffedLink, error) {
	sniffCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(sniffCtx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	baseURL, err := url.Parse(resp.Request.URL.String())
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 4*1024*1024)) // Read up to 4MB HTML
	if err != nil {
		return nil, err
	}
	html := string(bodyBytes)

	matches := hrefRegex.FindAllStringSubmatch(html, -1)
	seen := make(map[string]bool)
	var links []SniffedLink

	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		rawHref := strings.TrimSpace(m[1])
		if rawHref == "" || strings.HasPrefix(rawHref, "javascript:") || strings.HasPrefix(rawHref, "mailto:") || strings.HasPrefix(rawHref, "data:") {
			continue
		}

		parsedRef, err := url.Parse(rawHref)
		if err != nil {
			continue
		}
		resolved := baseURL.ResolveReference(parsedRef)
		if resolved.Scheme != "http" && resolved.Scheme != "https" {
			continue
		}

		fullURL := resolved.String()
		if seen[fullURL] {
			continue
		}

		baseName := path.Base(resolved.Path)
		if unescaped, err := url.PathUnescape(baseName); err == nil && unescaped != "" {
			baseName = unescaped
		}
		ext := strings.ToLower(path.Ext(baseName))
		if ext == "" || ext == ".html" || ext == ".htm" || ext == ".php" || ext == ".asp" || ext == ".aspx" || ext == ".jsp" || ext == ".css" || ext == ".js" {
			continue
		}

		seen[fullURL] = true
		cat := CategorizeFile(baseName, "")
		links = append(links, SniffedLink{
			URL:       fullURL,
			FileName:  SanitizeFileName(baseName),
			Extension: strings.TrimPrefix(ext, "."),
			Category:  cat,
			Label:     baseName,
		})
	}

	return links, nil
}

// IsLikelyDownloadURL checks if a clipboard string looks like a direct file download URL.
func IsLikelyDownloadURL(text string) bool {
	text = strings.TrimSpace(text)
	if len(text) < 10 || len(text) > 2048 || strings.ContainsAny(text, " \t\r\n") {
		return false
	}
	u, err := url.Parse(text)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return false
	}
	ext := strings.ToLower(path.Ext(u.Path))
	switch ext {
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".xz", ".tgz", ".bz2",
		".exe", ".msi", ".dmg", ".iso", ".img", ".apk", ".deb", ".rpm",
		".mp4", ".mkv", ".webm", ".avi", ".mov", ".mp3", ".flac", ".wav", ".m4a",
		".pdf", ".epub", ".pptx", ".xlsx", ".docx":
		return true
	}
	return false
}
