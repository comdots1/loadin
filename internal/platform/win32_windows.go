//go:build windows

package platform

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

//go:embed AetherHover.exe
var embeddedHoverExe []byte

//go:embed Loadin.ico
var embeddedAppIcon []byte

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	user32                   = syscall.NewLazyDLL("user32.dll")
	shell32                  = syscall.NewLazyDLL("shell32.dll")
	ole32                    = syscall.NewLazyDLL("ole32.dll")
	dwmapi                   = syscall.NewLazyDLL("dwmapi.dll")
	procGetDiskFreeSpaceExW  = kernel32.NewProc("GetDiskFreeSpaceExW")
	procGlobalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
	procGlobalLock           = kernel32.NewProc("GlobalLock")
	procGlobalUnlock         = kernel32.NewProc("GlobalUnlock")
	procOpenClipboard        = user32.NewProc("OpenClipboard")
	procCloseClipboard       = user32.NewProc("CloseClipboard")
	procGetClipboardData     = user32.NewProc("GetClipboardData")
	procIsClipboardFormatAvl = user32.NewProc("IsClipboardFormatAvailable")
	procEnumWindows          = user32.NewProc("EnumWindows")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procIsWindowVisible      = user32.NewProc("IsWindowVisible")
	procShowWindow           = user32.NewProc("ShowWindow")
	procIsZoomed             = user32.NewProc("IsZoomed")
	procPostMessageW         = user32.NewProc("PostMessageW")
	procSendMessageW         = user32.NewProc("SendMessageW")
	procSetForegroundWindow  = user32.NewProc("SetForegroundWindow")
	procLoadImageW           = user32.NewProc("LoadImageW")
	procSetClassLongPtrW     = user32.NewProc("SetClassLongPtrW")
	procDwmSetWindowAttr     = dwmapi.NewProc("DwmSetWindowAttribute")
	procShellExecuteW        = shell32.NewProc("ShellExecuteW")
	procSHBrowseForFolderW   = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
	procShellNotifyIconW     = shell32.NewProc("Shell_NotifyIconW")
	procCoInitializeEx       = ole32.NewProc("CoInitializeEx")
	procCoUninitialize       = ole32.NewProc("CoUninitialize")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")
	procGetModuleHandleW     = kernel32.NewProc("GetModuleHandleW")
	procRegisterClassExW     = user32.NewProc("RegisterClassExW")
	procCreateWindowExW      = user32.NewProc("CreateWindowExW")
	procDefWindowProcW       = user32.NewProc("DefWindowProcW")
	procGetMessageW          = user32.NewProc("GetMessageW")
	procTranslateMessage     = user32.NewProc("TranslateMessage")
	procDispatchMessageW     = user32.NewProc("DispatchMessageW")
	procCreatePopupMenu      = user32.NewProc("CreatePopupMenu")
	procAppendMenuW          = user32.NewProc("AppendMenuW")
	procSetMenuDefaultItem   = user32.NewProc("SetMenuDefaultItem")
	procTrackPopupMenu       = user32.NewProc("TrackPopupMenu")
	procDestroyMenu          = user32.NewProc("DestroyMenu")
	procGetCursorPos         = user32.NewProc("GetCursorPos")
	procIsIconic             = user32.NewProc("IsIconic")
	procRegisterWindowMsgW   = user32.NewProc("RegisterWindowMessageW")
)

const (
	cfUnicodeText  = 13
	createNoWindow = 0x08000000
	swShownormal   = 1
	swMaximize     = 3
	swMinimize     = 6
	swRestore      = 9
	wmClose        = 0x0010
	wmSetIcon      = 0x0080
	imageIcon      = 1
	lrLoadFromFile = 0x0010
)

// DiskSpaceInfo holds drive free and total capacity in bytes.
type DiskSpaceInfo struct {
	FreeBytes  uint64 `json:"freeBytes"`
	TotalBytes uint64 `json:"totalBytes"`
	Drive      string `json:"drive"`
}

// HardwareProfile contains detected PC hardware specifications and tailored thread recommendations.
type HardwareProfile struct {
	CPUName                string  `json:"cpuName"`
	LogicalCores           int     `json:"logicalCores"`
	TotalRAMBytes          uint64  `json:"totalRamBytes"`
	AvailRAMBytes          uint64  `json:"availRamBytes"`
	TotalRAMGB             float64 `json:"totalRamGB"`
	TierLabel              string  `json:"tierLabel"`
	RecommendedConnections int     `json:"recommendedConnections"`
	RecommendedConcurrent  int     `json:"recommendedConcurrent"`
	MaxRecommendedConns    int     `json:"maxRecommendedConns"`
	RecommendationReason   string  `json:"recommendationReason"`
}

type memoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var cachedHardware *HardwareProfile

// GetHardwareProfile inspects CPU logical cores, physical RAM via GlobalMemoryStatusEx,
// and processor model to compute optimal segment thread counts for this system.
func GetHardwareProfile() HardwareProfile {
	if cachedHardware != nil {
		// Refresh available RAM dynamically while keeping static CPU info cached
		var mem memoryStatusEx
		mem.dwLength = uint32(unsafe.Sizeof(mem))
		if ret, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem))); ret != 0 {
			cachedHardware.AvailRAMBytes = mem.ullAvailPhys
		}
		return *cachedHardware
	}

	cores := runtime.NumCPU()
	if cores < 1 {
		cores = 4
	}

	var mem memoryStatusEx
	mem.dwLength = uint32(unsafe.Sizeof(mem))
	var totalRAM, availRAM uint64 = 8 * 1024 * 1024 * 1024, 4 * 1024 * 1024 * 1024
	if ret, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem))); ret != 0 {
		totalRAM = mem.ullTotalPhys
		availRAM = mem.ullAvailPhys
	}
	ramGB := float64(totalRAM) / (1024 * 1024 * 1024)

	cpuName := detectCPUName()
	if cpuName == "" {
		cpuName = fmt.Sprintf("%d-Core x64 Processor", cores)
	}

	// Calculate optimal thread count based on logical cores and available memory
	var recConns, recConcurrent, maxConns int
	var tier string

	switch {
	case cores >= 16 && ramGB >= 24:
		tier = "울트라 멀티코어 워크스테이션"
		recConns = 24
		recConcurrent = 6
		maxConns = 32
	case cores >= 12 && ramGB >= 14:
		tier = "하이퍼포먼스 시스템"
		recConns = 16
		recConcurrent = 5
		maxConns = 32
	case cores >= 8 && ramGB >= 7:
		tier = "옵티멈 퍼포먼스 시스템"
		recConns = 12
		recConcurrent = 4
		maxConns = 24
	case cores >= 4:
		tier = "스탠다드 멀티스레드 시스템"
		recConns = 8
		recConcurrent = 3
		maxConns = 16
	default:
		tier = "라이트 시스템"
		recConns = 4
		recConcurrent = 2
		maxConns = 8
	}

	reason := fmt.Sprintf("%s (%d 논리 코어 · %.1f GB RAM) 환경에 맞춰 컨텍스트 스위칭 오버헤드 없이 대역폭을 극대화하는 파일당 %d 스레드(동시 %d개 작업)를 추천합니다.",
		cpuName, cores, ramGB, recConns, recConcurrent)

	 cachedHardware = &HardwareProfile{
		CPUName:                cpuName,
		LogicalCores:           cores,
		TotalRAMBytes:          totalRAM,
		AvailRAMBytes:          availRAM,
		TotalRAMGB:             ramGB,
		TierLabel:              tier,
		RecommendedConnections: recConns,
		RecommendedConcurrent:  recConcurrent,
		MaxRecommendedConns:    maxConns,
		RecommendationReason:   reason,
	}
	return *cachedHardware
}

// RecommendThreadsForFileSize returns the unified hardware-recommended thread count
// so browser download popups, URL probes, and hardware spec badges always stay 100% consistent.
func RecommendThreadsForFileSize(fileSize int64) (int, string) {
	hw := GetHardwareProfile()
	base := hw.RecommendedConnections

	if fileSize <= 0 {
		return base, fmt.Sprintf("시스템 사양(%d코어 · %.0fGB RAM) 기준 %d 스레드 권장", hw.LogicalCores, hw.TotalRAMGB, base)
	}
	mb := float64(fileSize) / (1024 * 1024)
	return base, fmt.Sprintf("시스템 사양(%d코어 · %.0fGB RAM) 및 파일 크기(%.1f MB) 맞춤 %d 스레드 권장", hw.LogicalCores, hw.TotalRAMGB, mb, base)
}

func detectCPUName() string {
	cmd := exec.Command("reg.exe", "query", `HKLM\HARDWARE\Description\System\CentralProcessor\0`, "/v", "ProcessorNameString")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if idx := strings.Index(line, "REG_SZ"); idx != -1 {
			return strings.TrimSpace(line[idx+len("REG_SZ"):])
		}
	}
	return ""
}

// RegisterBrowserProtocol registers the loadin:// custom URI scheme under HKCU\Software\Classes\loadin
// so browsers can launch Loadin directly without requiring Administrator privileges.
func RegisterBrowserProtocol() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, _ = filepath.Abs(exePath)

	commands := [][]string{
		{"add", `HKCU\Software\Classes\loadin`, "/ve", "/d", "URL:Loadin Protocol", "/f"},
		{"add", `HKCU\Software\Classes\loadin`, "/v", "URL Protocol", "/d", "", "/f"},
		{"add", `HKCU\Software\Classes\loadin\DefaultIcon`, "/ve", "/d", fmt.Sprintf(`"%s",0`, exePath), "/f"},
		{"add", `HKCU\Software\Classes\loadin\shell\open\command`, "/ve", "/d", fmt.Sprintf(`"%s" --url="%%1"`, exePath), "/f"},
	}

	for _, args := range commands {
		cmd := exec.Command("reg.exe", args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// IsBrowserProtocolRegistered checks whether HKCU\Software\Classes\loadin is registered.
func IsBrowserProtocolRegistered() bool {
	cmd := exec.Command("reg.exe", "query", `HKCU\Software\Classes\loadin\shell\open\command`)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	return cmd.Run() == nil
}

// SetStartupAutoRun registers or removes Loadin from Windows Startup (HKCU\Software\Microsoft\Windows\CurrentVersion\Run).
func SetStartupAutoRun(enabled bool) error {
	if !enabled {
		cmd := exec.Command("reg.exe", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "Loadin", "/f")
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
		_ = cmd.Run()
		return nil
	}
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, _ = filepath.Abs(exePath)
	val := fmt.Sprintf(`"%s" --headless`, exePath)
	cmd := exec.Command("reg.exe", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "Loadin", "/d", val, "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	return cmd.Run()
}

// IsStartupAutoRunEnabled checks if Loadin is registered in Windows Startup.
func IsStartupAutoRunEnabled() bool {
	cmd := exec.Command("reg.exe", "query", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "Loadin")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	return cmd.Run() == nil
}

// GetDiskSpace queries kernel32!GetDiskFreeSpaceExW for the given directory.
func GetDiskSpace(dirPath string) DiskSpaceInfo {
	abs, err := filepath.Abs(dirPath)
	if err != nil {
		abs = "C:\\"
	}
	vol := filepath.VolumeName(abs)
	if vol == "" {
		vol = "C:"
	}
	root := vol + "\\"

	rootPtr, err := syscall.UTF16PtrFromString(root)
	if err != nil {
		return DiskSpaceInfo{Drive: vol}
	}

	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	ret, _, _ := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(rootPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)
	if ret == 0 {
		return DiskSpaceInfo{Drive: vol}
	}
	return DiskSpaceInfo{
		FreeBytes:  freeBytesAvailable,
		TotalBytes: totalNumberOfBytes,
		Drive:      vol,
	}
}

// ReadClipboardText reads CF_UNICODETEXT directly from Windows User32 Clipboard with zero process spawn.
func ReadClipboardText() string {
	avail, _, _ := procIsClipboardFormatAvl.Call(cfUnicodeText)
	if avail == 0 {
		return ""
	}
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return ""
	}
	defer procCloseClipboard.Call()

	hMem, _, _ := procGetClipboardData.Call(cfUnicodeText)
	if hMem == 0 {
		return ""
	}
	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return ""
	}
	defer procGlobalUnlock.Call(hMem)

	const maxChars = 8192
	slice := unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), maxChars)
	for i, v := range slice {
		if v == 0 {
			return syscall.UTF16ToString(slice[:i])
		}
	}
	return syscall.UTF16ToString(slice)
}

// OpenFileDefault opens a file or folder with the default associated Windows application via ShellExecuteW.
func OpenFileDefault(targetPath string) error {
	abs, err := filepath.Abs(targetPath)
	if err != nil {
		return err
	}
	verb, _ := syscall.UTF16PtrFromString("open")
	file, err := syscall.UTF16PtrFromString(abs)
	if err != nil {
		return err
	}
	procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		0,
		0,
		swShownormal,
	)
	return nil
}

// ShowInExplorer opens Windows Explorer and highlights the specified file.
func ShowInExplorer(targetPath string) error {
	abs, err := filepath.Abs(targetPath)
	if err != nil {
		return err
	}
	if _, statErr := os.Stat(abs); statErr != nil {
		parent := filepath.Dir(abs)
		return OpenFileDefault(parent)
	}
	cmd := exec.Command("explorer.exe", "/select,"+abs)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	return cmd.Start()
}

type browseInfoW struct {
	hwndOwner      uintptr
	pidlRoot       uintptr
	pszDisplayName *uint16
	lpszTitle      *uint16
	ulFlags        uint32
	lpfn           uintptr
	lParam         uintptr
	iImage         int32
}

const (
	bifReturnOnlyFsDirs = 0x00000001
	bifEditBox          = 0x00000010
	bifNewDialogStyle   = 0x00000040
	bffmInitialized     = 1
	bffmSetSelectionW   = 0x0400 + 103
	coinitApartment     = 0x2
)

// PickFolderDialog opens a native Win32 Shell Folder Browser Dialog (SHBrowseForFolderW)
// directly parented to the main Loadin window so it opens instantaneously and always in front.
func PickFolderDialog(initialDir string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	procCoInitializeEx.Call(0, coinitApartment)
	defer procCoUninitialize.Call()

	hwndOwner := findMainWindowHandle()
	if hwndOwner != 0 {
		procSetForegroundWindow.Call(hwndOwner)
	}

	displayName := make([]uint16, 260)
	titlePtr, _ := syscall.UTF16PtrFromString("Loadin — 다운로드 저장 폴더를 선택하세요")

	cleanInit := strings.TrimSpace(initialDir)
	if cleanInit != "" {
		if abs, err := filepath.Abs(cleanInit); err == nil {
			cleanInit = abs
		}
	}
	var initPtr *uint16
	if cleanInit != "" {
		initPtr, _ = syscall.UTF16PtrFromString(cleanInit)
	}

	callback := syscall.NewCallback(func(hwnd uintptr, uMsg uint32, lParam uintptr, lpData uintptr) uintptr {
		if uMsg == bffmInitialized {
			procSetForegroundWindow.Call(hwnd)
			if lpData != 0 {
				procSendMessageW.Call(hwnd, bffmSetSelectionW, 1, lpData)
			}
		}
		return 0
	})

	bi := browseInfoW{
		hwndOwner:      hwndOwner,
		pszDisplayName: &displayName[0],
		lpszTitle:      titlePtr,
		ulFlags:        bifReturnOnlyFsDirs | bifNewDialogStyle | bifEditBox,
		lpfn:           callback,
		lParam:         uintptr(unsafe.Pointer(initPtr)),
	}

	pidl, _, _ := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return "", nil
	}
	defer procCoTaskMemFree.Call(pidl)

	pathBuf := make([]uint16, 520)
	ok, _, _ := procSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&pathBuf[0])))
	if ok == 0 {
		return "", nil
	}
	selected := strings.TrimSpace(syscall.UTF16ToString(pathBuf))
	if selected != "" {
		_ = os.MkdirAll(selected, 0755)
	}
	return selected, nil
}

// CloseAllAppWindows closes all open Loadin UI and dialog windows and removes the system tray icon.
func CloseAllAppWindows() {
	RemoveSystemTray()

	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		vis, _, _ := procIsWindowVisible.Call(hwnd)
		if vis == 0 {
			return 1
		}
		buf := make([]uint16, 256)
		length, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if length > 0 {
			title := syscall.UTF16ToString(buf[:length])
			if strings.Contains(title, "Loadin") && !strings.Contains(title, "Setup") && !strings.Contains(title, "Uninstall") && !strings.Contains(title, "제거") {
				procPostMessageW.Call(hwnd, wmClose, 0, 0)
			}
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)

	killDialog := exec.Command("taskkill.exe", "/F", "/IM", "Loadin_AddDialog.exe")
	killDialog.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	_ = killDialog.Run()
}

// DefaultDownloadsDir returns the user's native Windows Downloads directory.
func DefaultDownloadsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "Downloads")
	}
	dl := filepath.Join(home, "Downloads", "Loadin")
	_ = os.MkdirAll(dl, 0755)
	return dl
}

// ShowOrRestoreMainWindow brings the existing Loadin window to the foreground if open/minimized,
// or launches a new desktop app window if it was closed to the system tray.
func ShowOrRestoreMainWindow(url string) {
	if hwnd := findMainWindowHandle(); hwnd != 0 {
		iconic, _, _ := procIsIconic.Call(hwnd)
		if iconic != 0 {
			procShowWindow.Call(hwnd, swRestore)
		} else {
			procShowWindow.Call(hwnd, swShownormal)
		}
		procSetForegroundWindow.Call(hwnd)
		go StyleMainWindowGlass()
		return
	}
	_, _ = LaunchDesktopAppWindow(url)
}

// LaunchDesktopAppWindow opens the local UI in a sleek native desktop window
// and applies dark glass DWM window frame attributes.
func LaunchDesktopAppWindow(url string) (*exec.Cmd, error) {
	if hwnd := findMainWindowHandle(); hwnd != 0 {
		iconic, _, _ := procIsIconic.Call(hwnd)
		if iconic != 0 {
			procShowWindow.Call(hwnd, swRestore)
		} else {
			procShowWindow.Call(hwnd, swShownormal)
		}
		procSetForegroundWindow.Call(hwnd)
		go StyleMainWindowGlass()
		return nil, nil
	}

	candidates := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
	}

	userDataDir := filepath.Join(os.TempDir(), "Loadin_AppProfile")
	_ = os.MkdirAll(userDataDir, 0755)

	for _, bin := range candidates {
		if _, err := os.Stat(bin); err == nil {
			cmd := exec.Command(
				bin,
				"--app="+url,
				"--window-size=1360,860",
				"--user-data-dir="+userDataDir,
				"--disable-features=Translate,msEdgeSidebarV2",
				"--no-first-run",
				"--no-default-browser-check",
			)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
			if err := cmd.Start(); err == nil {
				go func() {
					for i := 0; i < 10; i++ {
						time.Sleep(300 * time.Millisecond)
						if StyleMainWindowGlass() {
							break
						}
					}
				}()
				return cmd, nil
			}
		}
	}

	_ = OpenFileDefault(url)
	return nil, nil
}

// ShowAddDownloadDialog extracts and launches the native frameless "새로운 다운로드" popup window
// on demand when a download link is clicked in the browser or triggered by the user.
func ShowAddDownloadDialog(port int, targetURL, referer, cookie, userAgent, fileName string) (*exec.Cmd, error) {
	if len(embeddedHoverExe) == 0 {
		return nil, fmt.Errorf("embedded dialog binary is empty")
	}
	dialogPath := filepath.Join(os.TempDir(), "Loadin_AddDialog.exe")
	_ = os.WriteFile(dialogPath, embeddedHoverExe, 0755)

	args := []string{"--port", strconv.Itoa(port)}
	if strings.TrimSpace(targetURL) != "" {
		args = append(args, "--url", strings.TrimSpace(targetURL))
	}
	if strings.TrimSpace(referer) != "" {
		args = append(args, "--referer", strings.TrimSpace(referer))
	}
	if strings.TrimSpace(cookie) != "" {
		args = append(args, "--cookie", strings.TrimSpace(cookie))
	}
	if strings.TrimSpace(userAgent) != "" {
		args = append(args, "--user-agent", strings.TrimSpace(userAgent))
	}
	if strings.TrimSpace(fileName) != "" {
		args = append(args, "--filename", strings.TrimSpace(fileName))
	}

	cmd := exec.Command(dialogPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// LaunchFloatingHoverWidget opens the on-demand native Add Download dialog.
func LaunchFloatingHoverWidget(port int) (*exec.Cmd, error) {
	return ShowAddDownloadDialog(port, "", "", "", "", "")
}

// findMainWindowHandle locates the top-level window handle for the Loadin window.
func findMainWindowHandle() uintptr {
	var foundHwnd uintptr
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		vis, _, _ := procIsWindowVisible.Call(hwnd)
		if vis == 0 {
			return 1
		}
		buf := make([]uint16, 256)
		length, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if length > 0 {
			title := syscall.UTF16ToString(buf[:length])
			if strings.Contains(title, "Loadin") {
				foundHwnd = hwnd
				return 0 // stop enumeration
			}
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	return foundHwnd
}

var currentWindowTheme = "black"

// StyleMainWindowGlass applies Windows DWM dark/light mode, rounded corners, caption colors,
// and sets the native Loadin window icon (replacing the browser icon).
func StyleMainWindowGlass(theme ...string) bool {
	if len(theme) > 0 && theme[0] != "" {
		currentWindowTheme = theme[0]
	}
	hwnd := findMainWindowHandle()
	if hwnd == 0 {
		return false
	}

	isWhite := currentWindowTheme == "white" || currentWindowTheme == "light-frost"
	dark := int32(1)
	borderColor := uint32(0x00323236)
	captionColor := uint32(0x00141212) // #121214 in BGR
	textColor := uint32(0x00F7F5F5)    // #F5F5F7 in BGR
	if isWhite {
		dark = int32(0)
		borderColor = uint32(0x00E0E0E5)
		captionColor = uint32(0x00F7F5F5) // #F5F5F7 in BGR
		textColor = uint32(0x001F1D1D)    // #1D1D1F in BGR
	}
	cornerRound := int32(2) // DWMWCP_ROUND

	procDwmSetWindowAttr.Call(hwnd, 20, uintptr(unsafe.Pointer(&dark)), 4)
	procDwmSetWindowAttr.Call(hwnd, 33, uintptr(unsafe.Pointer(&cornerRound)), 4)
	procDwmSetWindowAttr.Call(hwnd, 34, uintptr(unsafe.Pointer(&borderColor)), 4)
	procDwmSetWindowAttr.Call(hwnd, 35, uintptr(unsafe.Pointer(&captionColor)), 4)
	procDwmSetWindowAttr.Call(hwnd, 36, uintptr(unsafe.Pointer(&textColor)), 4)

	// Apply custom Loadin.ico to window titlebar and taskbar icon
	if len(embeddedAppIcon) > 0 {
		icoPath := filepath.Join(os.TempDir(), "Loadin_AppIcon.ico")
		if _, err := os.Stat(icoPath); err != nil {
			_ = os.WriteFile(icoPath, embeddedAppIcon, 0644)
		}
		if icoPtr, err := syscall.UTF16PtrFromString(icoPath); err == nil {
			hSmall, _, _ := procLoadImageW.Call(0, uintptr(unsafe.Pointer(icoPtr)), imageIcon, 16, 16, lrLoadFromFile)
			hBig, _, _ := procLoadImageW.Call(0, uintptr(unsafe.Pointer(icoPtr)), imageIcon, 32, 32, lrLoadFromFile)
			if hSmall != 0 {
				procSendMessageW.Call(hwnd, wmSetIcon, 0, hSmall)
				procSetClassLongPtrW.Call(hwnd, ^uintptr(33), hSmall) // GCLP_HICONSM = -34
			}
			if hBig != 0 {
				procSendMessageW.Call(hwnd, wmSetIcon, 1, hBig)
				procSetClassLongPtrW.Call(hwnd, ^uintptr(13), hBig) // GCLP_HICON = -14
			}
		}
	}
	return true
}

// ControlMainWindow executes minimize, maximize/restore, or close on the main studio window.
func ControlMainWindow(action string) bool {
	hwnd := findMainWindowHandle()
	if hwnd == 0 {
		return false
	}
	switch action {
	case "minimize":
		procShowWindow.Call(hwnd, swMinimize)
	case "maximize":
		zoomed, _, _ := procIsZoomed.Call(hwnd)
		if zoomed != 0 {
			procShowWindow.Call(hwnd, swRestore)
		} else {
			procShowWindow.Call(hwnd, swMaximize)
		}
	case "close":
		procPostMessageW.Call(hwnd, wmClose, 0, 0)
	default:
		return false
	}
	return true
}

// ============================================================================
// Native Win32 System Tray Icon & Context Menu (Shell_NotifyIconW)
// ============================================================================

const (
	nimAdd         = 0x00000000
	nimModify      = 0x00000001
	nimDelete      = 0x00000002
	nifMessage     = 0x00000001
	nifIcon        = 0x00000002
	nifTip         = 0x00000004
	wmTrayCallback = 0x8000 + 77
	wmCommand      = 0x0111
	wmNull         = 0x0000
	wmLButtonUp    = 0x0202
	wmLButtonDbl   = 0x0203
	wmRButtonUp    = 0x0205
	wmContextMenu  = 0x007B
	mfString       = 0x00000000
	mfSeparator    = 0x00000800
	tpmRightButton = 0x0002
	tpmBottomAlign = 0x0020

	trayMenuOpenMain   = 1001
	trayMenuNewDialog  = 1002
	trayMenuOpenFolder = 1003
	trayMenuResumeAll  = 1004
	trayMenuPauseAll   = 1005
	trayMenuQuit       = 1006
)

type trayGUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type notifyIconDataW struct {
	cbSize           uint32
	hWnd             uintptr
	uID              uint32
	uFlags           uint32
	uCallbackMessage uint32
	hIcon            uintptr
	szTip            [128]uint16
	dwState          uint32
	dwStateMask      uint32
	szInfo           [256]uint16
	uVersion         uint32
	szInfoTitle      [64]uint16
	dwInfoFlags      uint32
	guidItem         trayGUID
	hBalloonIcon     uintptr
}

type wndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type winPoint struct {
	X int32
	Y int32
}

type winMsg struct {
	HWND    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      winPoint
}

// TrayCallbacks defines the actions invoked by the Windows System Tray icon and its right-click menu.
type TrayCallbacks struct {
	OnOpenMainWindow func()
	OnNewDownload    func()
	OnOpenSaveFolder func()
	OnResumeAll      func()
	OnPauseAll       func()
	OnQuit           func()
	GetLanguage      func() string
}

var (
	trayMu               sync.Mutex
	trayHwnd             uintptr
	trayIconHandle       uintptr
	trayCallbacks        TrayCallbacks
	trayTaskbarCreatedWM uint32
	lastTrayTooltip      string
)

// StartSystemTray initializes a native Windows System Tray icon (Notification Area)
// with left-click window restore and a localized right-click popup context menu.
func StartSystemTray(cb TrayCallbacks) {
	trayMu.Lock()
	trayCallbacks = cb
	trayMu.Unlock()

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		if len(embeddedAppIcon) > 0 {
			icoPath := filepath.Join(os.TempDir(), "Loadin_AppIcon.ico")
			if _, err := os.Stat(icoPath); err != nil {
				_ = os.WriteFile(icoPath, embeddedAppIcon, 0644)
			}
			if icoPtr, err := syscall.UTF16PtrFromString(icoPath); err == nil {
				hIcon, _, _ := procLoadImageW.Call(0, uintptr(unsafe.Pointer(icoPtr)), imageIcon, 32, 32, lrLoadFromFile)
				trayIconHandle = hIcon
			}
		}

		if taskMsgPtr, err := syscall.UTF16PtrFromString("TaskbarCreated"); err == nil {
			msgID, _, _ := procRegisterWindowMsgW.Call(uintptr(unsafe.Pointer(taskMsgPtr)))
			trayTaskbarCreatedWM = uint32(msgID)
		}

		hInst, _, _ := procGetModuleHandleW.Call(0)
		className, _ := syscall.UTF16PtrFromString("LoadinTrayHiddenWndClass")
		wndProc := syscall.NewCallback(trayWndProc)

		wcx := wndClassExW{
			cbSize:        uint32(unsafe.Sizeof(wndClassExW{})),
			lpfnWndProc:   wndProc,
			hInstance:     hInst,
			hIcon:         trayIconHandle,
			lpszClassName: className,
			hIconSm:       trayIconHandle,
		}
		procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wcx)))

		wndTitle, _ := syscall.UTF16PtrFromString("LoadinTrayHost")
		hwnd, _, _ := procCreateWindowExW.Call(
			0,
			uintptr(unsafe.Pointer(className)),
			uintptr(unsafe.Pointer(wndTitle)),
			0,
			0, 0, 0, 0,
			0, 0, hInst, 0,
		)
		if hwnd == 0 {
			return
		}

		trayMu.Lock()
		trayHwnd = hwnd
		trayMu.Unlock()

		addOrUpdateTrayIcon(nimAdd, "Loadin — 다운로드 매니저 (실행 중)")

		var msg winMsg
		for {
			ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if int32(ret) <= 0 {
				break
			}
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}()
}

func addOrUpdateTrayIcon(action uint32, tooltip string) {
	trayMu.Lock()
	hwnd := trayHwnd
	hIcon := trayIconHandle
	trayMu.Unlock()
	if hwnd == 0 {
		return
	}

	var nid notifyIconDataW
	nid.cbSize = uint32(unsafe.Sizeof(nid))
	nid.hWnd = hwnd
	nid.uID = 1
	nid.uFlags = nifMessage | nifIcon | nifTip
	nid.uCallbackMessage = wmTrayCallback
	nid.hIcon = hIcon

	utf16Tip, _ := syscall.UTF16FromString(tooltip)
	for i := 0; i < len(utf16Tip) && i < len(nid.szTip)-1; i++ {
		nid.szTip[i] = utf16Tip[i]
	}

	procShellNotifyIconW.Call(uintptr(action), uintptr(unsafe.Pointer(&nid)))
}

// UpdateSystemTrayTooltip updates the hover text on the Windows System Tray icon.
func UpdateSystemTrayTooltip(tooltip string) {
	trayMu.Lock()
	if tooltip == lastTrayTooltip || trayHwnd == 0 {
		trayMu.Unlock()
		return
	}
	lastTrayTooltip = tooltip
	trayMu.Unlock()

	addOrUpdateTrayIcon(nimModify, tooltip)
}

// RemoveSystemTray removes the Loadin icon from the Windows System Tray on exit.
func RemoveSystemTray() {
	trayMu.Lock()
	hwnd := trayHwnd
	trayHwnd = 0
	trayMu.Unlock()
	if hwnd == 0 {
		return
	}
	var nid notifyIconDataW
	nid.cbSize = uint32(unsafe.Sizeof(nid))
	nid.hWnd = hwnd
	nid.uID = 1
	procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&nid)))
}

func trayWndProc(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	if trayTaskbarCreatedWM != 0 && msg == trayTaskbarCreatedWM {
		addOrUpdateTrayIcon(nimAdd, "Loadin — 다운로드 매니저")
		return 0
	}

	switch msg {
	case wmTrayCallback:
		lowEvent := uint32(lParam & 0xFFFF)
		switch lowEvent {
		case wmLButtonUp, wmLButtonDbl:
			if trayCallbacks.OnOpenMainWindow != nil {
				go trayCallbacks.OnOpenMainWindow()
			}
		case wmRButtonUp, wmContextMenu:
			showTrayContextMenu(hwnd)
		}
		return 0

	case wmCommand:
		cmdID := uint32(wParam & 0xFFFF)
		switch cmdID {
		case trayMenuOpenMain:
			if trayCallbacks.OnOpenMainWindow != nil {
				go trayCallbacks.OnOpenMainWindow()
			}
		case trayMenuNewDialog:
			if trayCallbacks.OnNewDownload != nil {
				go trayCallbacks.OnNewDownload()
			}
		case trayMenuOpenFolder:
			if trayCallbacks.OnOpenSaveFolder != nil {
				go trayCallbacks.OnOpenSaveFolder()
			}
		case trayMenuResumeAll:
			if trayCallbacks.OnResumeAll != nil {
				go trayCallbacks.OnResumeAll()
			}
		case trayMenuPauseAll:
			if trayCallbacks.OnPauseAll != nil {
				go trayCallbacks.OnPauseAll()
			}
		case trayMenuQuit:
			if trayCallbacks.OnQuit != nil {
				go trayCallbacks.OnQuit()
			}
		}
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func showTrayContextMenu(hwnd uintptr) {
	hMenu, _, _ := procCreatePopupMenu.Call()
	if hMenu == 0 {
		return
	}
	defer procDestroyMenu.Call(hMenu)

	lang := "ko"
	if trayCallbacks.GetLanguage != nil {
		if l := trayCallbacks.GetLanguage(); l != "" {
			lang = l
		}
	}

	lblOpen := "Loadin 열기"
	lblNew := "새 다운로드 추가..."
	lblFolder := "기본 저장 폴더 열기"
	lblResume := "모든 다운로드 이어받기"
	lblPause := "모든 다운로드 일시정지"
	lblQuit := "Loadin 완전 종료"

	if lang == "en" {
		lblOpen = "Open Loadin"
		lblNew = "New Download..."
		lblFolder = "Open Download Folder"
		lblResume = "Resume All Downloads"
		lblPause = "Pause All Downloads"
		lblQuit = "Quit Loadin"
	} else if lang == "ja" {
		lblOpen = "Loadin を開く"
		lblNew = "新規ダウンロード追加..."
		lblFolder = "保存先フォルダを開く"
		lblResume = "すべて再開"
		lblPause = "すべて一時停止"
		lblQuit = "Loadin を完全終了"
	}

	appendMenuStr := func(id uintptr, text string) {
		ptr, _ := syscall.UTF16PtrFromString(text)
		procAppendMenuW.Call(hMenu, mfString, id, uintptr(unsafe.Pointer(ptr)))
	}

	appendMenuStr(trayMenuOpenMain, lblOpen)
	procSetMenuDefaultItem.Call(hMenu, trayMenuOpenMain, 0)
	appendMenuStr(trayMenuNewDialog, lblNew)
	appendMenuStr(trayMenuOpenFolder, lblFolder)
	procAppendMenuW.Call(hMenu, mfSeparator, 0, 0)
	appendMenuStr(trayMenuResumeAll, lblResume)
	appendMenuStr(trayMenuPauseAll, lblPause)
	procAppendMenuW.Call(hMenu, mfSeparator, 0, 0)
	appendMenuStr(trayMenuQuit, lblQuit)

	var pt winPoint
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	// Required by Win32 so clicking outside the popup menu dismisses it cleanly
	procSetForegroundWindow.Call(hwnd)
	procTrackPopupMenu.Call(hMenu, tpmRightButton|tpmBottomAlign, uintptr(pt.X), uintptr(pt.Y), 0, hwnd, 0)
	procPostMessageW.Call(hwnd, wmNull, 0, 0)
}


