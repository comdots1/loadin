$ErrorActionPreference = "Stop"

$cscExe = "C:\Windows\Microsoft.NET\Framework64\v4.0.30319\csc.exe"
if (Test-Path $cscExe) {
    Write-Host "1. Generating Multi-Resolution Loadin App Icons (Loadin.ico & favicon.ico)..." -ForegroundColor Cyan
    & $cscExe /nologo /target:exe /optimize+ /out:internal\hover\GenIcon.exe internal\hover\GenIcon.cs
    & .\internal\hover\GenIcon.exe internal\platform\Loadin.ico internal\server\ui\favicon.ico

    Write-Host "2. Compiling Native Add Download Dialog (AetherHover.exe)..." -ForegroundColor Cyan
    & $cscExe /nologo /target:winexe /optimize+ /win32icon:internal\platform\Loadin.ico /out:internal\platform\AetherHover.exe internal\hover\AetherHover.cs
}

$goExe = "go"
if (-not (Get-Command $goExe -ErrorAction SilentlyContinue)) {
    $candidates = @(
        "C:\Users\Comdots\.gemini\antigravity\brain\385d27c9-d2b4-4900-a756-85013364b92f\scratch\go-sdk\go\bin\go.exe",
        "$env:TEMP\go-sdk\go\bin\go.exe",
        "$env:LOCALAPPDATA\go-toolchain\go\bin\go.exe"
    )
    foreach ($c in $candidates) {
        if (Test-Path $c -ErrorAction SilentlyContinue) {
            $goExe = $c
            break
        }
    }
}

Write-Host "3. Compiling Loadin Native Windows Executable (Loadin.exe)..." -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path "internal\engine\extassets" | Out-Null
Copy-Item -Force "browser-extension\manifest.json" "internal\engine\extassets\manifest.json"
Copy-Item -Force "browser-extension\background.js" "internal\engine\extassets\background.js"
Copy-Item -Force "browser-extension\content.js" "internal\engine\extassets\content.js"
Copy-Item -Force "browser-extension\popup.html" "internal\engine\extassets\popup.html"
Copy-Item -Force "browser-extension\popup.js" "internal\engine\extassets\popup.js"

$env:CGO_ENABLED = "0"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:GOCACHE = "$env:TEMP\aetherdl-gocache"
$env:GOPATH = "$env:TEMP\aetherdl-gopath"

& $goExe build -trimpath -ldflags="-s -w -H=windowsgui" -o Loadin.exe .
$item = Get-Item Loadin.exe
Write-Host "   -> Built Loadin.exe: $($item.FullName) ($([Math]::Round($item.Length / 1MB, 2)) MB)" -ForegroundColor Green

if (Test-Path $cscExe) {
    Write-Host "4. Compiling Commercial Windows Setup Installer (LoadinSetup.exe)..." -ForegroundColor Cyan
    & $cscExe /nologo /target:winexe /optimize+ `
        /win32icon:internal\platform\Loadin.ico `
        /resource:Loadin.exe,Loadin.exe `
        /resource:internal\platform\Loadin.ico,Loadin.ico `
        /resource:browser-extension\manifest.json,manifest.json `
        /resource:browser-extension\background.js,background.js `
        /resource:browser-extension\content.js,content.js `
        /resource:browser-extension\popup.html,popup.html `
        /resource:browser-extension\popup.js,popup.js `
        /out:LoadinSetup.exe `
        internal\installer\LoadinSetup.cs

    $setupItem = Get-Item LoadinSetup.exe
    Write-Host "   -> Built LoadinSetup.exe: $($setupItem.FullName) ($([Math]::Round($setupItem.Length / 1MB, 2)) MB)" -ForegroundColor Green
}
