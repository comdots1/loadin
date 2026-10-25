<div align="center">

# [✦ Loadin — High-Performance Native Download Accelerator](https://loadin.mooni.uk)

<p align="center">
  <strong>Next-Generation Windows Download Manager built with Pure Go 1.23 & Native Win32 Subsystems</strong><br>
  <em>Spatial Glassmorphism Interface · Dynamic Work-Stealing Acceleration · System Tray Suite</em>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go%201.23-00ADD8?style=flat-square&logo=go" alt="Go 1.23" />
  <img src="https://img.shields.io/badge/Platform-Windows%2010%20%7C%2011%20(x64)-0078D6?style=flat-square&logo=windows" alt="Windows" />
  <img src="https://img.shields.io/badge/Architecture-Zero--Merge%20Positional%20I%2FO-0A84FF?style=flat-square" alt="Architecture" />
  <img src="https://img.shields.io/badge/UI-Black%20%26%20White%20Glass-2C2C2E?style=flat-square" alt="Glassmorphism" />
  <img src="https://img.shields.io/badge/License-MIT-30D158?style=flat-square" alt="License" />
</p>

</div>

---

## Overview

**Loadin**은 고성능 시스템 프로그래밍 언어인 **Go 1.23**과 **Windows 네이티브 Win32/C# 서브시스템**으로 설계된 상용 소프트웨어급 초고속 다운로드 가속 엔진입니다.

다운로드 종료 시 파일 조각들을 결합하느라 시스템이 멈추는 **기존 다운로더들의 병합 지연을 0초로 단축**하고, 시스템 하드웨어 사양(CPU 논리 코어·물리 메모리)을 정밀 분석하여 최적의 멀티 스레드를 자동으로 구성합니다. 평소에는 불필요한 위젯 없이 백그라운드 및 **시스템 트레이**에서 가볍게 상주하며, 웹 브라우저에서 다운로드 링크를 클릭할 때만 전용 다운로드 추가 창이 즉시 반응합니다.

---

## 아키텍처 및 엔지니어링

### 1. 동적 세그먼트 워크 스틸링 (Dynamic Work-Stealing Segmentation)
- 고정 분할 방식의 한계를 극복하기 위해 **실시간 세그먼트 분할(Atomic Work-Stealing)** 알고리즘을 탑재했습니다.
- 전송 속도가 빠른 고루틴(Worker)이 작업을 마쳤을 때, 여전히 진행 중인 가장 큰 세그먼트(남은 용량 $\ge 1\text{ MB}$)를 찾아내어 **중간 지점을 원자적으로 분할(Atomic Bisect)하고 즉시 새로운 HTTP `Range: bytes=mid-end` 연결을 생성**합니다.
- 네트워크 변동성이나 대역폭 병목 현상이 발생해도 마지막 1바이트까지 모든 스레드가 100% 대역폭 효율을 유지합니다.

### 2. 병합 지연 0초 — 위치 지정 직접 I/O (`os.File.WriteAt`)
- 임시 분할 파일(`.part1`, `.part2` 등)을 내려받은 뒤 마지막 99%에서 디스크를 다시 읽고 합치는 병합 과정이 전혀 없습니다.
- 작업 시작 시 대상 파일의 크기만큼 디스크 공간을 선할당(`Truncate`)하며, 최대 32개의 스레드가 **Win32 오프셋 기반 위치 직접 I/O(`os.File.WriteAt`)**로 지정된 바이트 블록에 락-프리(Lock-Free)로 즉시 병렬 기록합니다.

### 3. 하드웨어 스펙 기반 최적 스레드 동기화 (Hardware Spec Advisor)
- 시스템 시작 시 **CPU 모델명, 논리 코어 수(`runtime.NumCPU()`), 물리 RAM(`GlobalMemoryStatusEx`)**을 복합 연산하여 컨텍스트 스위칭 오버헤드 없이 네트워크 인터페이스 카드를 극한으로 활용할 수 있는 권장 스레드 수를 산출합니다.
- 메인 대시보드 배지, 주소 분석 결과(URL Probe), 브라우저 연동 팝업 대화상자까지 **모든 환경에서 동일한 사양 맞춤 권장 스레드가 100% 일치하도록 보장**됩니다.

### 4. 윈도우 시스템 트레이 연동 (System Tray Integration)
- Windows `Shell_NotifyIconW` 네이티브 API를 활용한 전용 트레이 아이콘을 지원합니다.
- **실시간 속도 툴팁**: 마우스 호버 시 현재 진행 중인 다운로드 건수와 실시간 다운로드 속도(`Loadin — 24.50 MB/s (2개 다운로드 중)`)를 표시합니다.
- **원클릭 창 복원 & 컨텍스트 메뉴**:
  - 트레이 아이콘 좌클릭/더블클릭 시 메인 창 즉시 복원 및 최상단 포커스
  - 트레이 우클릭 컨텍스트 메뉴: 메인 창 열기, 새 다운로드 창, 기본 저장 폴더 열기, 모두 이어받기/일시정지, 프로그램 완전 종료

### 5. 독립형 브라우저 연동 스위트 (Browser Integration Suite)
- **온디맨드 네이티브 추가 창 ([AetherHover.exe](file:///c:/Users/Comdots/Documents/antigravity/beautiful-euclid/internal/platform/AetherHover.exe))**:
  - 화면을 가리는 불필요한 위젯을 상시 띄우지 않고, 브라우저에서 다운로드 링크를 클릭했을 때만 `560×412` 크기의 유려한 대화상자가 즉시 나타납니다.
- **원클릭 브라우저 확장 주소 복사**:
  - 브라우저 연동 창에서 `edge://extensions` 및 `chrome://extensions` 버튼을 클릭 한 번으로 클립보드에 복사할 수 있습니다.
- **시스템 프로토콜 (`loadin://`)**:
  - 관리자 권한 없이 `HKCU\Software\Classes\loadin` 레지스트리에 프로토콜을 등록하여 브라우저에서 Loadin을 즉시 호출합니다.

### 6. 다국어(i18n) & 테마 영구 동기화
- **2종 미니멀 테마**: 깊이감 있는 다크 그라파이트 **`Black` (`#121214`)** 과 청량한 데이라이트 **`White` (`#FFFFFF`)**
- **3개 국어 인터페이스**: **한국어 (Korean)**, **English**, **日本語 (Japanese)** 완벽 지원
- 백그라운드로 전환하거나 프로그램을 재실행해도 선택한 테마와 언어가 그대로 유지됩니다.

### 7. 단독 설치 마법사 & 클린 제거 엔진 ([LoadinSetup.exe](file:///c:/Users/Comdots/Documents/antigravity/beautiful-euclid/LoadinSetup.exe))
- 별도의 서드파티 인스톨러 없이 C# .NET 기반의 그래픽 설치 마법사를 탑재했습니다.
- 바탕화면/시작 메뉴 바로가기 생성, Windows 시작프로그램 자동 실행 등록(`-headless`), 프로토콜 및 브라우저 확장 프로그램 자동 구성을 제공합니다.
- **실행 중인 프로그램 자동 종료**: 삭제 프로그램(`Uninstall.exe` / `LoadinSetup.exe /uninstall`) 실행 시 본 프로그램 및 백그라운드 엔진을 안전하게 종료한 뒤 레지스트리와 설치 폴더를 완전히 정리합니다.

---

## 디자인 & 사용성

| 영역 | 기능 설명 |
| :--- | :--- |
| **미니멀 스크롤바** | 모든 목록 및 모달 창에 트랙이 투명하고 둥근 전용 캡슐형 커스텀 스크롤바 적용 |
| **환경설정 창 정비** | Win32 네이티브 `SHBrowseForFolderW` 기반 지연 없는 폴더 변경, 탐색기 열기, 기본값 초기화, 캡슐 토글 스위치 제공 |
| **전용 타이포그래피** | UI 레이아웃에 **Pretendard Variable**, 속도 및 바이트 수치에 **JetBrains Mono** 고정폭 폰트 적용 |
| **Zero-Flicker DOM** | 5Hz(200ms) 고빈도 텔레메트리 갱신 시에도 마우스 호버 및 포커스가 깜빡이지 않는 정밀 DOM 패칭 |

---

## 디렉터리 구조

```text
beautiful-euclid/
├── cmd/
├── internal/
│   ├── engine/           # 멀티 세그먼트 다운로드 엔진, 워크 스틸링, 무결성 검증
│   ├── hover/            # 온디맨드 다운로드 추가 팝업 (C# WinForms 네이티브 창)
│   ├── installer/        # 단독 설치 및 삭제 마법사 (LoadinSetup.cs)
│   ├── platform/         # Win32 API, DWM 글래스 프레임, 시스템 트레이, 폴더 선택
│   └── server/           # 임베디드 웹뷰 서버, SSE 스트림, 브라우저 인터셉터
│       └── ui/           # 공간형 인터페이스, 커스텀 벡터 아이콘, i18n
├── browser-extension/    # Chrome/Edge 브라우저 확장프로그램 (Manifest V3)
├── build.ps1             # 아이콘 생성, 호버 창, 엔진, 인스톨러 원클릭 빌드 스크립트
├── main.go               # Loadin 진입점, 싱글 인스턴스 핸들러, 시스템 트레이 구동
├── Loadin.exe            # 메인 실행 파일 (포터블 지원)
└── LoadinSetup.exe       # 상용 배포용 독립형 설치 프로그램
```

---

## Quick Start

### 1. 설치 프로그램을 통한 설치 (권장)
LoadinSetup.exe를 실행하여 원하는 설치 옵션을 선택하고 원클릭으로 설치를 완료합니다.

### 2. 포터블 직접 실행
```powershell
# 일반 데스크톱 모드로 실행
.\Loadin.exe

# 백그라운드 트레이 모드로 실행 (시작프로그램용)
.\Loadin.exe -headless

# 특정 URL 즉시 큐 추가
.\Loadin.exe -url "https://example.com/file.zip"
```

### 3. 소스코드 전체 빌드
PowerShell 환경에서 build.ps1을 실행하면 아이콘 생성, C# 서브시스템, Go 엔진, 설치 마법사가 순서대로 자동 빌드됩니다.
```powershell
powershell -ExecutionPolicy Bypass -File .\build.ps1
```

### 4. 단위 테스트 실행
```powershell
go test -v ./internal/engine
```
---

## 🛠️ 기술 스택 (Tech Stack)

- **Core Engine**: Go 1.23.2 (Goroutine Parallelism, Dynamic Atomic Splitting, Positional I/O)
- **Windows Subsystem**: Win32 Shell API (`shell32.dll`, `user32.dll`, `dwmapi.dll`), C# .NET 4.0
- **Frontend Layer**: HTML5 Semantic Architecture, Modern CSS3 Flex/Grid, Vanilla ES6+
- **Telemetry Streaming**: Server-Sent Events (SSE) 5Hz Real-Time Broadcasting
- **Typography**: Pretendard Variable, JetBrains Mono, Inter

---

<div align="center">
  <sub>Crafted for speed, responsiveness, and seamless native integration.</sub><br>
  <sub>Copyright © 2026 Loadin Project. All rights reserved.</sub>
</div>
