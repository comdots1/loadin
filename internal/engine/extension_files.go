package engine

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed extassets/*
var embeddedExtFS embed.FS

const firefoxManifestJSON = `{
  "manifest_version": 3,
  "name": "Loadin Browser Connector",
  "version": "2.2.0",
  "description": "Firefox, Zen, Waterfox 등 Gecko 기반 브라우저에서 다운로드 링크 클릭 시 Loadin 다운로드 창을 자동으로 연결합니다.",
  "browser_specific_settings": {
    "gecko": {
      "id": "connector@loadin.app",
      "strict_min_version": "109.0"
    }
  },
  "permissions": [
    "downloads",
    "contextMenus",
    "storage",
    "tabs",
    "cookies"
  ],
  "host_permissions": [
    "http://127.0.0.1:17890/*",
    "http://localhost:17890/*",
    "<all_urls>"
  ],
  "background": {
    "scripts": ["background.js"]
  },
  "content_scripts": [
    {
      "matches": ["<all_urls>"],
      "js": ["content.js"],
      "run_at": "document_idle",
      "all_frames": true
    }
  ],
  "action": {
    "default_title": "Loadin 브라우저 연동",
    "default_popup": "popup.html"
  }
}
`

// EnsureBrowserExtensionFiles writes all required browser extension files (manifest.json,
// background.js, content.js, popup.html, popup.js, plus firefox/ subfolder) to targetDir
// so unpacking the extension in Chrome, Edge, Whale, Brave, Opera, Vivaldi, or Firefox never fails.
func EnsureBrowserExtensionFiles(targetDir string) error {
	if targetDir == "" {
		return nil
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	entries, err := fs.ReadDir(embeddedExtFS, "extassets")
	if err != nil {
		return err
	}

	firefoxDir := filepath.Join(targetDir, "firefox")
	_ = os.MkdirAll(firefoxDir, 0755)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		data, err := embeddedExtFS.ReadFile("extassets/" + name)
		if err != nil {
			continue
		}
		_ = os.WriteFile(filepath.Join(targetDir, name), data, 0644)
		if name != "manifest.json" {
			_ = os.WriteFile(filepath.Join(firefoxDir, name), data, 0644)
		}
	}

	_ = os.WriteFile(filepath.Join(firefoxDir, "manifest.json"), []byte(firefoxManifestJSON), 0644)
	return nil
}

// ResolveAndRepairExtensionDirs repairs all known Loadin browser-extension directories
// and returns the primary path to open for the user.
func ResolveAndRepairExtensionDirs(configDir string) string {
	candidates := []string{}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates, filepath.Join(exeDir, "browser-extension"))
	}

	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		progLoadinDir := filepath.Join(localAppData, "Programs", "Loadin")
		if st, err := os.Stat(progLoadinDir); err == nil && st.IsDir() {
			candidates = append(candidates, filepath.Join(progLoadinDir, "browser-extension"))
		}
	}

	if cwdExt, err := filepath.Abs("browser-extension"); err == nil {
		if st, err := os.Stat(cwdExt); err == nil && st.IsDir() {
			candidates = append(candidates, cwdExt)
		}
	}

	if configDir != "" {
		candidates = append(candidates, filepath.Join(configDir, "browser-extension"))
	}

	primary := ""
	seen := map[string]bool{}
	for _, dir := range candidates {
		clean := filepath.Clean(dir)
		if seen[clean] {
			continue
		}
		seen[clean] = true
		if err := EnsureBrowserExtensionFiles(clean); err == nil && primary == "" {
			primary = clean
		}
	}

	if primary == "" && configDir != "" {
		primary = filepath.Join(configDir, "browser-extension")
		_ = EnsureBrowserExtensionFiles(primary)
	}
	return primary
}
