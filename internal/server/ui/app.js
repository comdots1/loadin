// ============================================================================
// Loadin — Minimalist Controller, Persistent Theme & Multi-Language Engine
// ============================================================================

const I18N = {
  ko: {
    "nav.downloads": "다운로드",
    "nav.all": "전체 항목",
    "nav.active": "진행 중",
    "nav.completed": "완료됨",
    "nav.paused": "일시정지",
    "nav.categories": "카테고리",
    "nav.compressed": "압축 파일",
    "nav.programs": "프로그램",
    "nav.videos": "동영상",
    "nav.music": "오디오",
    "nav.documents": "문서",
    "hw.recommendedTitle": "권장 병렬 연결",
    "telemetry.speedLimit": "속도 제한",
    "telemetry.unlimited": "무제한",
    "telemetry.openFolder": "폴더 열기",
    "header.searchPlaceholder": "검색",
    "header.newDialog": "새 다운로드 창",
    "header.browserIntegration": "브라우저 연동",
    "header.linkScan": "링크 스캔",
    "header.batchAdd": "일괄 추가",
    "header.benchmark": "테스트",
    "clip.detected": "복사된 링크",
    "common.close": "닫기",
    "common.copy": "복사",
    "common.browse": "변경...",
    "common.cancel": "취소",
    "common.save": "저장",
    "composer.urlPlaceholder": "다운로드할 파일 주소(URL)를 붙여넣으세요...",
    "composer.probe": "분석",
    "composer.probing": "분석 중",
    "composer.options": "옵션",
    "composer.download": "다운로드",
    "probe.file": "파일",
    "probe.size": "크기",
    "probe.mode": "방식",
    "probe.recommended": "권장 설정",
    "composer.applyRecommended": "권장값 적용",
    "composer.fileName": "파일 이름",
    "composer.autoDetect": "자동 감지",
    "composer.parallelThreads": "병렬 스레드",
    "composer.saveLocation": "저장 위치",
    "composer.defaultFolder": "기본 다운로드 폴더",
    "composer.checksum": "체크섬 검증 (SHA-256 / MD5)",
    "composer.optional": "선택 사항",
    "composer.referer": "Referer 주소",
    "composer.cookie": "인증 Cookie 헤더",
    "empty.title": "다운로드 내역이 없습니다",
    "empty.desc": "상단 입력창에 주소를 붙여넣거나 웹 브라우저에서 다운로드 링크를 클릭하면 자동으로 최적의 스레드로 분할하여 빠르게 내려받습니다.",
    "empty.benchBtn": "64MB 가속 테스트 실행",
    "empty.browserBtn": "브라우저 연동 설정",
    "insp.details": "상세 정보",
    "insp.speed": "속도",
    "insp.eta": "남은 시간",
    "insp.segments": "세그먼트",
    "insp.splits": "동적 분할",
    "insp.speedHistory": "전송 속도 추이",
    "insp.segProgress": "세그먼트 진행도",
    "insp.threadState": "스레드별 상태",
    "insp.checksum": "체크섬 (클릭하여 복사)",
    "insp.logs": "로그",
    "browser.title": "브라우저 연동 및 다운로드 대화상자",
    "browser.desc": "평소에는 화면에 위젯을 띄워두지 않으며, 웹 브라우저에서 다운로드 링크를 클릭할 때만 깔끔한 다운로드 추가 창이 열립니다.",
    "browser.extBadge": "확장프로그램",
    "browser.extTitle": "브라우저 확장프로그램 설치",
    "browser.extDesc": "웹페이지에서 다운로드 링크를 클릭하면 Loadin의 새로운 다운로드 대화상자가 자동으로 나타나 파일 크기와 사양 맞춤 스레드를 안내합니다.",
    "browser.step1": "아래 [확장프로그램 폴더 열기]를 클릭합니다.",
    "browser.step2": "아래 주소 복사 버튼을 눌러 브라우저 주소창에 붙여넣고 개발자 모드를 켭니다.",
    "browser.step3": "[압축해제된 확장 프로그램을 로드합니다]를 클릭해 열린 browser-extension 폴더를 선택합니다.",
    "browser.openExtFolder": "확장프로그램 폴더 열기",
    "browser.sysBadge": "시스템 연동",
    "browser.sysTitle": "다운로드 대화상자 & 프로토콜 등록",
    "browser.sysDesc": "독립형 다운로드 대화상자를 바로 열어보거나, 브라우저에서 Loadin을 직접 호출할 수 있도록 시스템 프로토콜을 등록합니다.",
    "browser.openNewDialog": "새 다운로드 대화상자 열기",
    "browser.registerProtocol": "loadin:// 시스템 프로토콜 등록",
    "browser.bookmarkletLabel": "북마크바에 끌어다 놓아 사용할 수 있는 북마클릿:",
    "browser.bookmarkletBtn": "Loadin으로 다운로드",
    "bench.title": "분할 가속 테스트",
    "bench.card1Title": "동적 세그먼트 분할 시뮬레이션",
    "bench.card1Desc": "먼저 완료된 스레드가 느린 스레드의 남은 구간을 절반으로 나누어 이어받는 과정을 확인합니다.",
    "bench.card1Btn": "64MB 테스트 시작",
    "bench.card2Title": "32-스레드 대용량 스트림",
    "bench.card2Desc": "32개 병렬 연결과 실시간 SHA-256 무결성 검증 처리 속도를 측정합니다.",
    "bench.card2Btn": "256MB 테스트 시작",
    "bench.card3Title": "컴파일러 SDK 아카이브 (.zip)",
    "bench.card3Desc": "글로벌 CDN 서버에서 실제 대용량 압축 파일을 멀티 스레드로 내려받습니다.",
    "bench.card3Btn": "78MB 파일 다운로드",
    "bench.card4Title": "포터블 패키지 아카이브 (.7z.exe)",
    "bench.card4Desc": "외부 릴리즈 서버에서 리다이렉트 추적 및 병렬 다운로드를 수행합니다.",
    "bench.card4Btn": "55MB 파일 다운로드",
    "sniffer.title": "웹페이지 링크 스캔",
    "sniffer.placeholder": "스캔할 웹페이지 주소 입력...",
    "sniffer.scanBtn": "스캔",
    "sniffer.selectAll": "전체 선택",
    "sniffer.deselectAll": "선택 해제",
    "sniffer.emptyHint": "웹페이지 주소를 입력하고 [스캔]을 클릭하세요.",
    "sniffer.downloadSelected": "선택 항목 다운로드",
    "batch.title": "일괄 다운로드 추가",
    "batch.desc": "여러 개의 주소를 줄바꿈으로 입력하거나 범위 패턴(예: https://example.com/file_[01-10].zip)을 입력하세요.",
    "batch.addAll": "모두 추가",
    "settings.title": "환경설정",
    "settings.secStorage": "저장 위치 및 가속 엔진",
    "settings.secAppearance": "화면 및 언어",
    "settings.secAutomation": "시스템 연동 및 자동화",
    "settings.pickFolder": "폴더 변경...",
    "settings.resetFolder": "기본값",
    "settings.speedLimitDesc": "전체 다운로드 대역폭 상한 설정",
    "settings.browserInterceptTitle": "브라우저 다운로드 자동 연동",
    "settings.browserInterceptDesc": "확장프로그램에서 다운로드 링크 클릭 시 Loadin 추가 창 자동 표시",
    "settings.quitApp": "프로그램 완전 종료",
    "settings.hwRecTitle": "사양 맞춤 권장 스레드",
    "settings.themeTitle": "화면 테마",
    "settings.themeDesc": "블랙 또는 화이트 모드 선택 (백그라운드 전환 및 재실행 후에도 유지됨)",
    "settings.langTitle": "언어 (Language)",
    "settings.langDesc": "인터페이스 및 다운로드 대화상자 표시 언어",
    "settings.startupTitle": "윈도우 시작 시 자동 실행",
    "settings.startupDesc": "로그인 시 백그라운드 엔진을 자동으로 시작하여 브라우저 연동을 즉시 대기",
    "settings.saveDirTitle": "기본 저장 위치",
    "settings.saveDirDesc": "다운로드한 파일이 저장될 폴더",
    "settings.connTitle": "기본 스레드 수 (1 ~ 32)",
    "settings.connDesc": "파일당 분할 연결 개수",
    "settings.concurrentTitle": "동시 다운로드 작업 수",
    "settings.concurrentDesc": "한 번에 동시에 내려받을 최대 파일 수",
    "settings.autoCatTitle": "카테고리별 하위 폴더 자동 분류",
    "settings.autoCatDesc": "파일 종류에 따라 폴더별로 정리",
    "settings.clipboardTitle": "클립보드 주소 자동 감지",
    "settings.clipboardDesc": "다운로드 주소 복사 시 알림 배너 표시",
    "settings.soundTitle": "완료 효과음",
    "settings.soundDesc": "다운로드 완료 시 알림음 재생",
    "theme.black": "블랙 (Black)",
    "theme.white": "화이트 (White)",
  },
  en: {
    "nav.downloads": "Downloads",
    "nav.all": "All Items",
    "nav.active": "Active",
    "nav.completed": "Completed",
    "nav.paused": "Paused",
    "nav.categories": "Categories",
    "nav.compressed": "Archives",
    "nav.programs": "Programs",
    "nav.videos": "Videos",
    "nav.music": "Audio",
    "nav.documents": "Documents",
    "hw.recommendedTitle": "Recommended Threads",
    "telemetry.speedLimit": "Speed Limit",
    "telemetry.unlimited": "Unlimited",
    "telemetry.openFolder": "Open Folder",
    "header.searchPlaceholder": "Search",
    "header.newDialog": "New Download",
    "header.browserIntegration": "Browser Add-on",
    "header.linkScan": "Link Scanner",
    "header.batchAdd": "Batch Add",
    "header.benchmark": "Benchmark",
    "clip.detected": "Copied Link",
    "common.close": "Close",
    "common.copy": "Copy",
    "common.browse": "Browse...",
    "common.cancel": "Cancel",
    "common.save": "Save",
    "composer.urlPlaceholder": "Paste download URL here...",
    "composer.probe": "Analyze",
    "composer.probing": "Analyzing",
    "composer.options": "Options",
    "composer.download": "Download",
    "probe.file": "File",
    "probe.size": "Size",
    "probe.mode": "Mode",
    "probe.recommended": "Recommended",
    "composer.applyRecommended": "Apply Recommended",
    "composer.fileName": "File Name",
    "composer.autoDetect": "Auto-detect",
    "composer.parallelThreads": "Parallel Threads",
    "composer.saveLocation": "Save Folder",
    "composer.defaultFolder": "Default download folder",
    "composer.checksum": "Checksum (SHA-256 / MD5)",
    "composer.optional": "Optional",
    "composer.referer": "Referer URL",
    "composer.cookie": "Auth Cookie Header",
    "empty.title": "No downloads yet",
    "empty.desc": "Paste a file URL above or click any download link in your browser to accelerate with optimal parallel threads.",
    "empty.benchBtn": "Run 64MB Speed Test",
    "empty.browserBtn": "Browser Integration",
    "insp.details": "Inspector",
    "insp.speed": "Speed",
    "insp.eta": "ETA",
    "insp.segments": "Segments",
    "insp.splits": "Dynamic Splits",
    "insp.speedHistory": "Throughput History",
    "insp.segProgress": "Segment Progress",
    "insp.threadState": "Thread Status",
    "insp.checksum": "Checksums (Click to copy)",
    "insp.logs": "Logs",
    "browser.title": "Browser Integration & Download Dialog",
    "browser.desc": "No cluttered floating widget stays on screen. A sleek download dialog appears only when you click a download link.",
    "browser.extBadge": "Extension",
    "browser.extTitle": "Install Browser Extension",
    "browser.extDesc": "Automatically opens the Loadin Add Download dialog with hardware-matched threads whenever you click a download link.",
    "browser.step1": "Click [Open Extension Folder] below.",
    "browser.step2": "Click one of the copy buttons below, paste into your browser address bar, and enable Developer Mode.",
    "browser.step3": "Click [Load unpacked] and select the opened browser-extension folder.",
    "browser.openExtFolder": "Open Extension Folder",
    "browser.sysBadge": "System",
    "browser.sysTitle": "Download Dialog & Protocol Registration",
    "browser.sysDesc": "Test the standalone Add Download dialog or register the system protocol so browsers can launch Loadin directly.",
    "browser.openNewDialog": "Open New Download Dialog",
    "browser.registerProtocol": "Register loadin:// Protocol",
    "browser.bookmarkletLabel": "Drag this bookmarklet to your bookmarks bar:",
    "browser.bookmarkletBtn": "Download with Loadin",
    "bench.title": "Multi-Segment Acceleration Test",
    "bench.card1Title": "Dynamic Segment Split Simulation",
    "bench.card1Desc": "Watch fast threads automatically bisect and take over remaining work from slower threads.",
    "bench.card1Btn": "Start 64MB Test",
    "bench.card2Title": "32-Thread High-Load Stream",
    "bench.card2Desc": "Measures 32 parallel connections and real-time SHA-256 integrity verification speed.",
    "bench.card2Btn": "Start 256MB Test",
    "bench.card3Title": "Compiler SDK Archive (.zip)",
    "bench.card3Desc": "Downloads a real multi-megabyte archive from a global CDN using parallel streams.",
    "bench.card3Btn": "Download 78MB File",
    "bench.card4Title": "Portable Package Archive (.7z.exe)",
    "bench.card4Desc": "Follows release server redirects and performs multi-segment acceleration.",
    "bench.card4Btn": "Download 55MB File",
    "sniffer.title": "Webpage Link Scanner",
    "sniffer.placeholder": "Enter webpage URL to scan...",
    "sniffer.scanBtn": "Scan",
    "sniffer.selectAll": "Select All",
    "sniffer.deselectAll": "Deselect All",
    "sniffer.emptyHint": "Enter a webpage URL and click [Scan].",
    "sniffer.downloadSelected": "Download Selected",
    "batch.title": "Batch Add Downloads",
    "batch.desc": "Enter multiple URLs (one per line) or use a wildcard range pattern (e.g. https://example.com/file_[01-10].zip).",
    "batch.addAll": "Add All",
    "settings.title": "Settings",
    "settings.secStorage": "Storage & Acceleration Engine",
    "settings.secAppearance": "Appearance & Language",
    "settings.secAutomation": "System Integration & Automation",
    "settings.pickFolder": "Change Folder...",
    "settings.resetFolder": "Default",
    "settings.speedLimitDesc": "Maximum global download bandwidth limit",
    "settings.browserInterceptTitle": "Auto-Intercept Browser Downloads",
    "settings.browserInterceptDesc": "Show Loadin Add Download dialog when clicking links in browser",
    "settings.quitApp": "Quit Loadin Completely",
    "settings.hwRecTitle": "Hardware-Matched Threads",
    "settings.themeTitle": "Appearance Theme",
    "settings.themeDesc": "Choose Black or White mode (persists across background & restarts)",
    "settings.langTitle": "Language",
    "settings.langDesc": "Display language for main window and popup dialog",
    "settings.startupTitle": "Launch on Windows Startup",
    "settings.startupDesc": "Automatically start Loadin engine in the background when you sign in to Windows",
    "settings.saveDirTitle": "Default Save Folder",
    "settings.saveDirDesc": "Directory where downloaded files are stored",
    "settings.connTitle": "Default Threads (1 ~ 32)",
    "settings.connDesc": "Parallel connections per file",
    "settings.concurrentTitle": "Max Concurrent Downloads",
    "settings.concurrentDesc": "Maximum number of files downloaded simultaneously",
    "settings.autoCatTitle": "Auto-Categorize Subfolders",
    "settings.autoCatDesc": "Organize downloads into folders by file type",
    "settings.clipboardTitle": "Monitor Clipboard for Links",
    "settings.clipboardDesc": "Show quick banner when a download URL is copied",
    "settings.soundTitle": "Completion Sound",
    "settings.soundDesc": "Play chime when a download finishes",
    "theme.black": "Black",
    "theme.white": "White",
  },
  ja: {
    "nav.downloads": "ダウンロード",
    "nav.all": "すべての項目",
    "nav.active": "進行中",
    "nav.completed": "完了",
    "nav.paused": "一時停止",
    "nav.categories": "カテゴリ",
    "nav.compressed": "圧縮ファイル",
    "nav.programs": "プログラム",
    "nav.videos": "動画",
    "nav.music": "オーディオ",
    "nav.documents": "ドキュメント",
    "hw.recommendedTitle": "推奨並列スレッド",
    "telemetry.speedLimit": "速度制限",
    "telemetry.unlimited": "無制限",
    "telemetry.openFolder": "フォルダを開く",
    "header.searchPlaceholder": "検索",
    "header.newDialog": "新規ダウンロード窓",
    "header.browserIntegration": "ブラウザ連携",
    "header.linkScan": "リンクスキャン",
    "header.batchAdd": "一括追加",
    "header.benchmark": "テスト",
    "clip.detected": "コピーされたURL",
    "common.close": "閉じる",
    "common.copy": "コピー",
    "common.browse": "変更...",
    "common.cancel": "キャンセル",
    "common.save": "保存",
    "composer.urlPlaceholder": "ダウンロードするURLを貼り付けてください...",
    "composer.probe": "解析",
    "composer.probing": "解析中",
    "composer.options": "オプション",
    "composer.download": "ダウンロード",
    "probe.file": "ファイル",
    "probe.size": "サイズ",
    "probe.mode": "方式",
    "probe.recommended": "推奨設定",
    "composer.applyRecommended": "推奨値を適用",
    "composer.fileName": "ファイル名",
    "composer.autoDetect": "自動検出",
    "composer.parallelThreads": "並列スレッド",
    "composer.saveLocation": "保存先",
    "composer.defaultFolder": "既定のダウンロードフォルダ",
    "composer.checksum": "チェックサム検証 (SHA-256 / MD5)",
    "composer.optional": "任意",
    "composer.referer": "Referer URL",
    "composer.cookie": "認証 Cookie ヘッダー",
    "empty.title": "ダウンロード履歴がありません",
    "empty.desc": "上部の入力欄にURLを貼り付けるか、ブラウザでダウンロードリンクをクリックすると最適スレッドで高速取得します。",
    "empty.benchBtn": "64MB 加速テストを実行",
    "empty.browserBtn": "ブラウザ連携設定",
    "insp.details": "詳細情報",
    "insp.speed": "速度",
    "insp.eta": "残り時間",
    "insp.segments": "セグメント",
    "insp.splits": "動的分割",
    "insp.speedHistory": "転送速度推移",
    "insp.segProgress": "セグメント進捗",
    "insp.threadState": "スレッド別状態",
    "insp.checksum": "チェックサム (クリックでコピー)",
    "insp.logs": "ログ",
    "browser.title": "ブラウザ連携＆ダウンロードダイアログ",
    "browser.desc": "普段は邪魔なウィジェットを表示せず、ブラウザでダウンロードリンクをクリックした時だけダイアログが開きます。",
    "browser.extBadge": "拡張機能",
    "browser.extTitle": "ブラウザ拡張機能のインストール",
    "browser.extDesc": "WebページでリンクをクリックするとLoadinの追加ダイアログが自動で開き、スペック最適スレッドを提案します。",
    "browser.step1": "下の [拡張機能フォルダを開く] をクリックします。",
    "browser.step2": "下のURLコピーボタンを押してブラウザのアドレスバーに貼り付け、デベロッパーモードをオンにします。",
    "browser.step3": "[パッケージ化されていない拡張機能を読み込む] で browser-extension フォルダを選択します。",
    "browser.openExtFolder": "拡張機能フォルダを開く",
    "browser.sysBadge": "システム連携",
    "browser.sysTitle": "ダイアログ＆プロトコル登録",
    "browser.sysDesc": "独立型ダイアログを今すぐ開くか、ブラウザからLoadinを直接呼び出すシステムプロトコルを登録します。",
    "browser.openNewDialog": "新規ダウンロードダイアログを開く",
    "browser.registerProtocol": "loadin:// プロトコルを登録",
    "browser.bookmarkletLabel": "ブックマークバーにドラッグして使えるブックマークレット:",
    "browser.bookmarkletBtn": "Loadinでダウンロード",
    "bench.title": "マルチセグメント加速テスト",
    "bench.card1Title": "動的セグメント分割シミュレーション",
    "bench.card1Desc": "先に完了したスレッドが遅いスレッドの残り区間を半分引き継ぐ過程を確認します。",
    "bench.card1Btn": "64MB テスト開始",
    "bench.card2Title": "32スレッド大容量ストリーム",
    "bench.card2Desc": "32並列接続とリアルタイムSHA-256検証速度を測定します。",
    "bench.card2Btn": "256MB テスト開始",
    "bench.card3Title": "コンパイラ SDK アーカイブ (.zip)",
    "bench.card3Desc": "グローバルCDNから実際の大容量圧縮ファイルを並列取得します。",
    "bench.card3Btn": "78MB ファイルをダウンロード",
    "bench.card4Title": "ポータブルパッケージ (.7z.exe)",
    "bench.card4Desc": "外部リリースサーバーからリダイレクト追跡と並列取得を行います。",
    "bench.card4Btn": "55MB ファイルをダウンロード",
    "sniffer.title": "Webページリンクスキャン",
    "sniffer.placeholder": "スキャンするWebページのURLを入力...",
    "sniffer.scanBtn": "スキャン",
    "sniffer.selectAll": "すべて選択",
    "sniffer.deselectAll": "選択解除",
    "sniffer.emptyHint": "URLを入力して [スキャン] をクリックしてください。",
    "sniffer.downloadSelected": "選択項目をダウンロード",
    "batch.title": "一括ダウンロード追加",
    "batch.desc": "複数のURLを改行で入力するか、範囲パターン(例: https://example.com/file_[01-10].zip)を入力してください。",
    "batch.addAll": "すべて追加",
    "settings.title": "環境設定",
    "settings.secStorage": "保存先および加速エンジン",
    "settings.secAppearance": "外観および言語",
    "settings.secAutomation": "システム連携および自動化",
    "settings.pickFolder": "フォルダ変更...",
    "settings.resetFolder": "初期値",
    "settings.speedLimitDesc": "全体の最大ダウンロード帯域幅制限",
    "settings.browserInterceptTitle": "ブラウザダウンロード自動連携",
    "settings.browserInterceptDesc": "拡張機能でリンクをクリックした際に追加ダイアログを自動表示",
    "settings.quitApp": "プログラムを完全終了",
    "settings.hwRecTitle": "スペック推奨スレッド",
    "settings.themeTitle": "テーマ設定",
    "settings.themeDesc": "ブラックまたはホワイトモードを選択 (再起動後も維持されます)",
    "settings.langTitle": "言語 (Language)",
    "settings.langDesc": "メイン画面およびダイアログの表示言語",
    "settings.startupTitle": "Windows 起動時に自動実行",
    "settings.startupDesc": "サインイン時にバックグラウンドエンジンを自動起動し、ブラウザ連携を待機します",
    "settings.saveDirTitle": "既定の保存先",
    "settings.saveDirDesc": "ダウンロードしたファイルを保存するフォルダ",
    "settings.connTitle": "既定のスレッド数 (1 ~ 32)",
    "settings.connDesc": "1ファイルあたりの分割接続数",
    "settings.concurrentTitle": "同時ダウンロード数",
    "settings.concurrentDesc": "同時にダウンロードする最大ファイル数",
    "settings.autoCatTitle": "カテゴリ別サブフォルダ自動分類",
    "settings.autoCatDesc": "ファイルの種類ごとにフォルダを整理",
    "settings.clipboardTitle": "クリップボードURL自動検知",
    "settings.clipboardDesc": "URLコピー時に通知バナーを表示",
    "settings.soundTitle": "完了サウンド",
    "settings.soundDesc": "ダウンロード完了時に通知音を再生",
    "theme.black": "ブラック (Black)",
    "theme.white": "ホワイト (White)",
  },
};

const state = {
  snapshot: null,
  activeFilter: "all",
  searchQuery: "",
  selectedTaskId: null,
  knownCompletedIds: new Set(),
  initialLoadDone: false,
  sniffedLinks: [],
  sniffFilter: "all",
  theme: localStorage.getItem("loadin_theme") === "white" ? "white" : "black",
  language: localStorage.getItem("loadin_lang") || "ko",
  recommendedThreads: 16,
  hardwareInitialized: false,
};

function tStr(key, fallback = "") {
  const dict = I18N[state.language] || I18N.ko;
  return dict[key] || I18N.ko[key] || fallback || key;
}

function applyThemeToDOM(theme) {
  const normalized = theme === "white" ? "white" : "black";
  state.theme = normalized;
  if (document.documentElement.getAttribute("data-theme") !== normalized) {
    document.documentElement.setAttribute("data-theme", normalized);
  }
  try {
    localStorage.setItem("loadin_theme", normalized);
  } catch (_) {}
  const setThemeEl = document.getElementById("set-theme");
  if (setThemeEl && setThemeEl.value !== normalized) {
    setThemeEl.value = normalized;
  }
}

function applyLanguageToDOM(lang) {
  const normalized = lang === "en" || lang === "ja" ? lang : "ko";
  state.language = normalized;
  document.documentElement.setAttribute("lang", normalized);
  try {
    localStorage.setItem("loadin_lang", normalized);
  } catch (_) {}

  const headerLang = document.getElementById("header-lang-select");
  if (headerLang && headerLang.value !== normalized) {
    headerLang.value = normalized;
  }
  const setLang = document.getElementById("set-language");
  if (setLang && setLang.value !== normalized) {
    setLang.value = normalized;
  }

  document.querySelectorAll("[data-i18n]").forEach((node) => {
    const key = node.getAttribute("data-i18n");
    const val = tStr(key);
    if (val && node.textContent !== val) {
      node.textContent = val;
    }
  });

  document.querySelectorAll("[data-i18n-placeholder]").forEach((node) => {
    const key = node.getAttribute("data-i18n-placeholder");
    const val = tStr(key);
    if (val && node.getAttribute("placeholder") !== val) {
      node.setAttribute("placeholder", val);
    }
  });

  updateSliderRecTag();
  if (state.snapshot) {
    // Force refresh of task action button labels when language changes
    document.querySelectorAll(".task-card").forEach((c) => {
      c.dataset.lastStatus = "";
    });
    renderGlobalState(state.snapshot);
  }
}

function iconSvg(id, extraClass = "") {
  return `<svg class="svg-ic ${extraClass}"><use href="#${id}"/></svg>`;
}

const el = {
  countAll: document.getElementById("count-all"),
  countActive: document.getElementById("count-active"),
  countCompleted: document.getElementById("count-completed"),
  countPaused: document.getElementById("count-paused"),
  countCompressed: document.getElementById("count-compressed"),
  countPrograms: document.getElementById("count-programs"),
  countVideos: document.getElementById("count-videos"),
  countMusic: document.getElementById("count-music"),
  countDocuments: document.getElementById("count-documents"),

  sidebarHwCard: document.getElementById("sidebar-hw-card"),
  hwRecBadge: document.getElementById("hw-rec-badge"),
  hwSpecsText: document.getElementById("hw-specs-text"),

  globalWorkersBadge: document.getElementById("global-workers-badge"),
  globalSpeedText: document.getElementById("global-speed-text"),
  globalPeakText: document.getElementById("global-peak-text"),
  globalSparkline: document.getElementById("global-sparkline"),
  quickSpeedLimit: document.getElementById("quick-speed-limit"),

  diskDriveLabel: document.getElementById("disk-drive-label"),
  diskBarFill: document.getElementById("disk-bar-fill"),
  diskFreeText: document.getElementById("disk-free-text"),
  btnOpenSaveDir: document.getElementById("btn-open-save-dir"),

  searchInput: document.getElementById("search-input"),
  clipboardBanner: document.getElementById("clipboard-banner"),
  clipboardDetectedUrl: document.getElementById("clipboard-detected-url"),
  btnClipDownload: document.getElementById("btn-clip-download"),
  btnClipDismiss: document.getElementById("btn-clip-dismiss"),

  urlInput: document.getElementById("url-input"),
  btnProbeUrl: document.getElementById("btn-probe-url"),
  btnToggleAdvanced: document.getElementById("btn-toggle-advanced"),
  advChevron: document.getElementById("adv-chevron"),
  btnAddDownload: document.getElementById("btn-add-download"),

  probePreviewBar: document.getElementById("probe-preview-bar"),
  probeFilename: document.getElementById("probe-filename"),
  probeSize: document.getElementById("probe-size"),
  probeRange: document.getElementById("probe-range"),
  probeRecThreads: document.getElementById("probe-rec-threads"),

  composerAdvanced: document.getElementById("composer-advanced"),
  advHwReason: document.getElementById("adv-hw-reason"),
  btnApplyRecThreads: document.getElementById("btn-apply-rec-threads"),
  btnRecCountLabel: document.getElementById("btn-rec-count-label"),
  advFilename: document.getElementById("adv-filename"),
  advConnections: document.getElementById("adv-connections"),
  advConnVal: document.getElementById("adv-conn-val"),
  advConnRecTag: document.getElementById("adv-conn-rec-tag"),
  advSavedir: document.getElementById("adv-savedir"),
  btnPickAdvFolder: document.getElementById("btn-pick-adv-folder"),
  advHash: document.getElementById("adv-hash"),
  advReferer: document.getElementById("adv-referer"),
  advCookie: document.getElementById("adv-cookie"),

  emptyState: document.getElementById("empty-state"),
  tasksList: document.getElementById("tasks-list"),

  inspectorDrawer: document.getElementById("inspector-drawer"),
  btnCloseInspector: document.getElementById("btn-close-inspector"),
  inspFilename: document.getElementById("insp-filename"),
  inspSpeed: document.getElementById("insp-speed"),
  inspEta: document.getElementById("insp-eta"),
  inspSegmentsCount: document.getElementById("insp-segments-count"),
  inspSplits: document.getElementById("insp-splits"),
  inspPeakAvg: document.getElementById("insp-peak-avg"),
  inspSpeedCanvas: document.getElementById("insp-speed-canvas"),
  inspBytesProgress: document.getElementById("insp-bytes-progress"),
  inspByteMap: document.getElementById("insp-byte-map"),
  inspSegmentsList: document.getElementById("insp-segments-list"),
  inspHashBadge: document.getElementById("insp-hash-badge"),
  inspSha256: document.getElementById("insp-sha256"),
  inspMd5: document.getElementById("insp-md5"),
  inspServerTag: document.getElementById("insp-server-tag"),
  inspLogs: document.getElementById("insp-logs"),

  extInterceptCount: document.getElementById("ext-intercept-count"),
  protoStatusText: document.getElementById("proto-status-text"),
  setHwSummary: document.getElementById("set-hw-summary"),

  toastContainer: document.getElementById("toast-container"),
};

function setTextIfChanged(node, text) {
  if (!node) return;
  const str = String(text ?? "");
  if (node.textContent !== str) {
    node.textContent = str;
  }
}

function threadUnit() {
  if (state.language === "en") return "threads";
  if (state.language === "ja") return "スレッド";
  return "스레드";
}

function formatBytes(bytes) {
  if (bytes === undefined || bytes === null || bytes < 0) {
    return state.language === "en" ? "Unknown" : state.language === "ja" ? "不明" : "알 수 없음";
  }
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  const val = bytes / Math.pow(1024, i);
  return `${val.toFixed(i === 0 ? 0 : 2)} ${units[i]}`;
}

function formatSpeed(bps) {
  if (!bps || bps <= 0) return "0.00 MB/s";
  if (bps < 1024 * 1024) {
    return `${(bps / 1024).toFixed(1)} KB/s`;
  }
  return `${(bps / (1024 * 1024)).toFixed(2)} MB/s`;
}

function formatETA(sec) {
  if (sec === undefined || sec === null || sec < 0) return "--:--";
  if (sec === 0) return state.language === "en" ? "Done" : state.language === "ja" ? "完了" : "완료됨";
  if (sec > 86400) return "> 1d";
  const h = Math.floor(sec / 3600);
  const m = Math.floor((sec % 3600) / 60);
  const s = Math.floor(sec % 60);
  if (h > 0) {
    return `${h}h ${String(m).padStart(2, "0")}m`;
  }
  return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
}

function categoryIconId(cat) {
  switch (cat) {
    case "compressed": return "ic-archive";
    case "programs": return "ic-cpu";
    case "videos": return "ic-video";
    case "music": return "ic-music";
    case "documents": return "ic-file";
    default: return "ic-download";
  }
}

function statusBadgeInfo(status) {
  const isEn = state.language === "en";
  const isJa = state.language === "ja";
  switch (status) {
    case "running":
      return { text: isEn ? "Downloading" : isJa ? "ダウンロード中" : "다운로드 중", cls: "pill-running" };
    case "probing":
      return { text: isEn ? "Analyzing" : isJa ? "解析中" : "분석 중", cls: "pill-running" };
    case "verifying":
      return { text: isEn ? "Verifying" : isJa ? "検証中" : "검증 중", cls: "pill-split" };
    case "completed":
      return { text: isEn ? "Completed" : isJa ? "完了" : "완료됨", cls: "pill-completed" };
    case "paused":
      return { text: isEn ? "Paused" : isJa ? "一時停止" : "일시정지", cls: "pill-paused" };
    case "queued":
      return { text: isEn ? "Queued" : isJa ? "待機中" : "대기 중", cls: "pill-paused" };
    case "error":
      return { text: isEn ? "Error" : isJa ? "エラー" : "오류", cls: "pill-error" };
    default:
      return { text: status, cls: "pill-paused" };
  }
}

window.showToast = function(msg) {
  const t = document.createElement("div");
  t.className = "glass-toast";
  t.textContent = msg;
  el.toastContainer.appendChild(t);
  setTimeout(() => {
    t.style.opacity = "0";
    t.style.transition = "opacity 0.18s ease";
    setTimeout(() => t.remove(), 200);
  }, 2600);
};

function playCompletionChime() {
  if (state.snapshot && state.snapshot.settings && !state.snapshot.settings.soundEffects) {
    return;
  }
  try {
    const ctx = new (window.AudioContext || window.webkitAudioContext)();
    const notes = [523.25, 659.25, 783.99, 1046.5];
    notes.forEach((freq, idx) => {
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();
      osc.type = "sine";
      osc.frequency.setValueAtTime(freq, ctx.currentTime + idx * 0.065);
      gain.gain.setValueAtTime(0.001, ctx.currentTime + idx * 0.065);
      gain.gain.exponentialRampToValueAtTime(0.06, ctx.currentTime + idx * 0.065 + 0.03);
      gain.gain.exponentialRampToValueAtTime(0.0001, ctx.currentTime + idx * 0.065 + 0.6);
      osc.connect(gain);
      gain.connect(ctx.destination);
      osc.start(ctx.currentTime + idx * 0.065);
      osc.stop(ctx.currentTime + idx * 0.065 + 0.65);
    });
  } catch (_) {}
}

function drawSmoothAreaChart(canvas, values, strokeColor, fillTopColor) {
  if (!canvas) return;
  const dpr = window.devicePixelRatio || 1;
  const rect = canvas.getBoundingClientRect();
  const w = rect.width || canvas.width;
  const h = rect.height || canvas.height;
  if (w <= 0 || h <= 0) return;

  const targetW = Math.floor(w * dpr);
  const targetH = Math.floor(h * dpr);
  if (canvas.width !== targetW || canvas.height !== targetH) {
    canvas.width = targetW;
    canvas.height = targetH;
  }
  const ctx = canvas.getContext("2d");
  ctx.save();
  ctx.scale(dpr, dpr);
  ctx.clearRect(0, 0, w, h);

  const pts = values && values.length > 1 ? values : [0, 0, 0, 0];
  let maxVal = Math.max(...pts, 1024 * 100);
  maxVal *= 1.15;

  const stepX = w / (pts.length - 1);
  const coords = pts.map((v, i) => ({
    x: i * stepX,
    y: h - 3 - (Math.min(v, maxVal) / maxVal) * (h - 8),
  }));

  const grad = ctx.createLinearGradient(0, 0, 0, h);
  grad.addColorStop(0, fillTopColor);
  grad.addColorStop(1, "rgba(56, 189, 248, 0.0)");

  ctx.beginPath();
  ctx.moveTo(coords[0].x, h);
  coords.forEach((p, idx) => {
    if (idx === 0) {
      ctx.lineTo(p.x, p.y);
    } else {
      const prev = coords[idx - 1];
      const cx = (prev.x + p.x) / 2;
      ctx.bezierCurveTo(cx, prev.y, cx, p.y, p.x, p.y);
    }
  });
  ctx.lineTo(coords[coords.length - 1].x, h);
  ctx.closePath();
  ctx.fillStyle = grad;
  ctx.fill();

  ctx.beginPath();
  coords.forEach((p, idx) => {
    if (idx === 0) {
      ctx.moveTo(p.x, p.y);
    } else {
      const prev = coords[idx - 1];
      const cx = (prev.x + p.x) / 2;
      ctx.bezierCurveTo(cx, prev.y, cx, p.y, p.x, p.y);
    }
  });
  ctx.strokeStyle = strokeColor;
  ctx.lineWidth = 1.6;
  ctx.stroke();

  ctx.restore();
}

function renderGlobalState(snap) {
  state.snapshot = snap;
  const tasks = snap.tasks || [];

  // Synchronize persisted theme & language from server settings so reopening the window never resets them
  if (snap.settings) {
    if (snap.settings.theme && (snap.settings.theme === "black" || snap.settings.theme === "white")) {
      if (state.theme !== snap.settings.theme) {
        applyThemeToDOM(snap.settings.theme);
      }
    }
    if (snap.settings.language && ["ko", "en", "ja"].includes(snap.settings.language)) {
      if (state.language !== snap.settings.language) {
        applyLanguageToDOM(snap.settings.language);
        return;
      }
    }
  }

  const isEn = state.language === "en";
  const isJa = state.language === "ja";

  if (snap.hardware) {
    const hw = snap.hardware;
    state.recommendedThreads = hw.recommendedConnections || 16;
    setTextIfChanged(el.hwRecBadge, `${hw.recommendedConnections} ${threadUnit()}`);

    const coreLabel = isEn
      ? `${hw.logicalCores} Cores · ${hw.totalRamGB.toFixed(0)}GB RAM`
      : isJa
      ? `${hw.logicalCores}コア · ${hw.totalRamGB.toFixed(0)}GB RAM`
      : `${hw.logicalCores}코어 · ${hw.totalRamGB.toFixed(0)}GB RAM`;
    setTextIfChanged(el.hwSpecsText, coreLabel);

    const reasonText = isEn
      ? `Hardware Spec (${coreLabel}) → ${hw.recommendedConnections} parallel threads recommended`
      : isJa
      ? `システム仕様 (${coreLabel}) → ${hw.recommendedConnections} 並列スレッド推奨`
      : hw.recommendationReason;
    setTextIfChanged(el.advHwReason, reasonText);
    setTextIfChanged(el.btnRecCountLabel, hw.recommendedConnections);
    setTextIfChanged(
      el.setHwSummary,
      `${hw.cpuName} (${coreLabel}) → ${hw.recommendedConnections} ${threadUnit()}`
    );

    if (!state.hardwareInitialized) {
      state.hardwareInitialized = true;
      el.advConnections.value = hw.recommendedConnections;
      setTextIfChanged(el.advConnVal, hw.recommendedConnections);
      updateSliderRecTag();
    }
  }

  const interceptCnt = snap.browserInterceptCount || 0;
  setTextIfChanged(
    el.extInterceptCount,
    isEn ? `Intercepted: ${interceptCnt}` : isJa ? `連携受信 ${interceptCnt}件` : `연동 수신 ${interceptCnt}건`
  );
  setTextIfChanged(
    el.protoStatusText,
    snap.protocolRegistered
      ? isEn ? "Registered" : isJa ? "登録済み" : "레지스트리 등록됨"
      : isEn ? "Not Registered" : isJa ? "未登録" : "미등록"
  );

  tasks.forEach((t) => {
    if (t.status === "completed") {
      if (state.initialLoadDone && !state.knownCompletedIds.has(t.id)) {
        state.knownCompletedIds.add(t.id);
        playCompletionChime();
        showToast(
          isEn
            ? `Download Complete: ${t.fileName}`
            : isJa
            ? `ダウンロード完了: ${t.fileName}`
            : `다운로드 완료: ${t.fileName}`
        );
      } else {
        state.knownCompletedIds.add(t.id);
      }
    }
  });
  state.initialLoadDone = true;

  let cActive = 0, cCompleted = 0, cPaused = 0;
  let cComp = 0, cProg = 0, cVid = 0, cMus = 0, cDoc = 0;

  tasks.forEach((t) => {
    if (t.status === "running" || t.status === "probing" || t.status === "verifying") cActive++;
    else if (t.status === "completed") cCompleted++;
    else cPaused++;

    if (t.category === "compressed") cComp++;
    else if (t.category === "programs") cProg++;
    else if (t.category === "videos") cVid++;
    else if (t.category === "music") cMus++;
    else if (t.category === "documents") cDoc++;
  });

  setTextIfChanged(el.countAll, tasks.length);
  setTextIfChanged(el.countActive, cActive);
  setTextIfChanged(el.countCompleted, cCompleted);
  setTextIfChanged(el.countPaused, cPaused);
  setTextIfChanged(el.countCompressed, cComp);
  setTextIfChanged(el.countPrograms, cProg);
  setTextIfChanged(el.countVideos, cVid);
  setTextIfChanged(el.countMusic, cMus);
  setTextIfChanged(el.countDocuments, cDoc);

  setTextIfChanged(el.globalSpeedText, formatSpeed(snap.totalSpeedBps));
  setTextIfChanged(el.globalWorkersBadge, `${snap.activeWorkersCount || 0} ${threadUnit()}`);
  drawSmoothAreaChart(
    el.globalSparkline,
    snap.globalSpeedHistory,
    "#0A84FF",
    "rgba(10, 132, 255, 0.22)"
  );

  if (snap.diskSpace && snap.diskSpace.totalBytes > 0) {
    const usedRatio = 100 - (snap.diskSpace.freeBytes / snap.diskSpace.totalBytes) * 100;
    const driveName = snap.diskSpace.drive || "C:";
    setTextIfChanged(
      el.diskDriveLabel,
      isEn ? `Disk (${driveName})` : isJa ? `ディスク (${driveName})` : `디스크 (${driveName})`
    );
    el.diskBarFill.style.width = `${Math.min(100, Math.max(4, usedRatio)).toFixed(1)}%`;
    const freePrefix = isEn ? "Free" : isJa ? "空き" : "여유";
    setTextIfChanged(
      el.diskFreeText,
      `${freePrefix} ${formatBytes(snap.diskSpace.freeBytes)} / ${formatBytes(snap.diskSpace.totalBytes)}`
    );
  }

  if (snap.detectedClipboard) {
    el.clipboardBanner.classList.remove("hidden");
    setTextIfChanged(el.clipboardDetectedUrl, snap.detectedClipboard);
  } else {
    el.clipboardBanner.classList.add("hidden");
  }

  const filtered = tasks.filter((t) => {
    if (state.activeFilter === "active") {
      if (t.status !== "running" && t.status !== "probing" && t.status !== "verifying") return false;
    } else if (state.activeFilter === "completed") {
      if (t.status !== "completed") return false;
    } else if (state.activeFilter === "paused") {
      if (t.status !== "paused" && t.status !== "queued" && t.status !== "error") return false;
    } else if (state.activeFilter.startsWith("cat:")) {
      const cat = state.activeFilter.slice(4);
      if (t.category !== cat) return false;
    } else if (["compressed", "programs", "videos", "music", "documents"].includes(state.activeFilter)) {
      if (t.category !== state.activeFilter) return false;
    }

    if (state.searchQuery) {
      const q = state.searchQuery.toLowerCase();
      const matchName = (t.fileName || "").toLowerCase().includes(q);
      const matchUrl = (t.url || "").toLowerCase().includes(q);
      const matchHash = ((t.checksum && t.checksum.sha256) || "").toLowerCase().includes(q);
      if (!matchName && !matchUrl && !matchHash) return false;
    }
    return true;
  });

  if (filtered.length === 0) {
    el.emptyState.classList.remove("hidden");
    el.tasksList.innerHTML = "";
  } else {
    el.emptyState.classList.add("hidden");
    renderTasksListZeroFlicker(filtered);
  }

  if (state.selectedTaskId) {
    const sel = tasks.find((t) => t.id === state.selectedTaskId);
    if (sel) {
      renderInspectorZeroFlicker(sel);
    } else {
      el.inspectorDrawer.classList.add("hidden");
      state.selectedTaskId = null;
    }
  }
}

function createTaskCardElement(t) {
  const card = document.createElement("div");
  card.className = "task-card";
  card.dataset.id = t.id;

  card.innerHTML = `
    <div class="task-top-row">
      <div class="task-identity">
        <div class="task-cat-icon" data-role="cat-icon"></div>
        <div class="task-title-wrap">
          <div class="task-name-row">
            <span class="task-filename" data-role="filename"></span>
            <span class="task-pill" data-role="status-pill"></span>
            <span class="task-pill pill-split hidden" data-role="split-pill"></span>
          </div>
          <div class="task-sub-meta tabular">
            <span data-role="bytes-meta"></span>
            <span data-role="speed-meta"></span>
          </div>
        </div>
      </div>
      <div class="task-actions" data-role="actions"></div>
    </div>
    <div class="task-progress-Area">
      <div class="segment-ribbon" data-role="ribbon"></div>
    </div>
  `;

  card._refs = {
    catIcon: card.querySelector('[data-role="cat-icon"]'),
    filename: card.querySelector('[data-role="filename"]'),
    statusPill: card.querySelector('[data-role="status-pill"]'),
    splitPill: card.querySelector('[data-role="split-pill"]'),
    bytesMeta: card.querySelector('[data-role="bytes-meta"]'),
    speedMeta: card.querySelector('[data-role="speed-meta"]'),
    actions: card.querySelector('[data-role="actions"]'),
    ribbon: card.querySelector('[data-role="ribbon"]'),
  };

  card.addEventListener("click", (e) => {
    if (e.target.closest("button")) return;
    if (state.selectedTaskId === t.id) {
      state.selectedTaskId = null;
      el.inspectorDrawer.classList.add("hidden");
      card.classList.remove("selected");
      return;
    }
    state.selectedTaskId = t.id;
    el.inspectorDrawer.classList.remove("hidden");
    if (state.snapshot && state.snapshot.tasks) {
      const cur = state.snapshot.tasks.find((x) => x.id === t.id);
      if (cur) renderInspectorZeroFlicker(cur);
    }
    document.querySelectorAll(".task-card").forEach((c) => {
      c.classList.toggle("selected", c.dataset.id === t.id);
    });
  });

  return card;
}

function renderTasksListZeroFlicker(tasks) {
  const existingMap = new Map();
  Array.from(el.tasksList.children).forEach((node) => {
    existingMap.set(node.dataset.id, node);
  });

  const desiredIds = new Set(tasks.map((t) => t.id));
  existingMap.forEach((node, id) => {
    if (!desiredIds.has(id)) node.remove();
  });

  const isEn = state.language === "en";
  const isJa = state.language === "ja";

  tasks.forEach((t, index) => {
    let card = existingMap.get(t.id);
    if (!card) {
      card = createTaskCardElement(t);
      if (index < el.tasksList.children.length) {
        el.tasksList.insertBefore(card, el.tasksList.children[index]);
      } else {
        el.tasksList.appendChild(card);
      }
    }

    const r = card._refs;
    const isSelected = state.selectedTaskId === t.id;
    if (card.classList.contains("selected") !== isSelected) {
      card.classList.toggle("selected", isSelected);
    }

    const iconId = categoryIconId(t.category);
    if (card.dataset.catIcon !== iconId) {
      card.dataset.catIcon = iconId;
      r.catIcon.innerHTML = iconSvg(iconId);
    }

    setTextIfChanged(r.filename, t.fileName);
    if (r.filename.title !== t.fileName) r.filename.title = t.fileName;

    const badge = statusBadgeInfo(t.status);
    setTextIfChanged(r.statusPill, badge.text);
    const desiredPillClass = `task-pill ${badge.cls}`;
    if (r.statusPill.className !== desiredPillClass) {
      r.statusPill.className = desiredPillClass;
    }

    if (t.dynamicSplits > 0) {
      r.splitPill.classList.remove("hidden");
      setTextIfChanged(
        r.splitPill,
        isEn ? `Split +${t.dynamicSplits}` : isJa ? `分割 +${t.dynamicSplits}` : `분할 +${t.dynamicSplits}`
      );
    } else {
      r.splitPill.classList.add("hidden");
    }

    const pct = (t.progress || 0).toFixed(1);
    setTextIfChanged(
      r.bytesMeta,
      `${formatBytes(t.downloaded)} / ${formatBytes(t.totalSize)} (${pct}%)`
    );

    let speedInfoText = "";
    if (t.status === "running" || t.status === "probing") {
      const totalSegs = t.segments ? t.segments.length : t.maxConnections;
      const etaLabel = isEn ? "ETA" : isJa ? "残り" : "남은 시간";
      speedInfoText = `${formatSpeed(t.speedBps)} · ${etaLabel} ${formatETA(t.etaSeconds)} · ${t.activeWorkers}/${totalSegs} ${threadUnit()}`;
      r.speedMeta.style.color = "var(--accent-primary)";
    } else if (t.status === "completed") {
      const avgLbl = isEn ? "Avg" : isJa ? "平均" : "평균";
      const peakLbl = isEn ? "Peak" : isJa ? "最高" : "최고";
      speedInfoText = `${avgLbl} ${formatSpeed(t.avgSpeedBps)} · ${peakLbl} ${formatSpeed(t.peakSpeedBps)}`;
      r.speedMeta.style.color = "var(--text-secondary)";
    } else if (t.errorMessage) {
      speedInfoText = t.errorMessage;
      r.speedMeta.style.color = "var(--accent-rose)";
    } else {
      speedInfoText = `${t.maxConnections} ${threadUnit()}`;
      r.speedMeta.style.color = "var(--text-muted)";
    }
    setTextIfChanged(r.speedMeta, speedInfoText);

    const statusLangKey = `${t.status}:${state.language}`;
    if (card.dataset.lastStatus !== statusLangKey) {
      card.dataset.lastStatus = statusLangKey;
      const lblPause = isEn ? "Pause" : isJa ? "一時停止" : "일시정지";
      const lblResume = isEn ? "Resume" : isJa ? "再開" : "이어받기";
      const lblOpen = isEn ? "Open" : isJa ? "開く" : "열기";
      const lblFolderTip = isEn ? "Show in Folder" : isJa ? "フォルダを開く" : "저장 폴더 열기";
      const lblCopyLinkTip = isEn ? "Copy Download Link" : isJa ? "ダウンロードリンクをコピー" : "다운로드 링크 복사";
      const lblDetailsTip = isEn ? "Task Details" : isJa ? "詳細情報" : "상세 정보";
      const lblRemoveTip = isEn ? "Remove from List (Keep File)" : isJa ? "リストから削除 (ファイル保持)" : "목록에서 제거 (파일 유지)";
      const lblDeleteFile = isEn ? "Delete File" : isJa ? "ファイル削除" : "파일 삭제";
      const lblDeleteFileTip = isEn ? "Delete Task & File from Disk" : isJa ? "ファイルごと完全に削除" : "목록 및 실제 파일까지 바로 삭제";

      let actionBtns = "";
      if (t.status === "running" || t.status === "probing") {
        actionBtns += `<button class="glass-btn btn-sm" onclick="taskAction('${t.id}', 'pause')">${iconSvg("ic-pause")}<span>${lblPause}</span></button>`;
      } else if (t.status === "paused" || t.status === "error" || t.status === "queued") {
        actionBtns += `<button class="glass-btn primary-glow btn-sm" onclick="taskAction('${t.id}', 'resume')">${iconSvg("ic-play")}<span>${lblResume}</span></button>`;
      }
      if (t.status === "completed") {
        actionBtns += `<button class="glass-btn primary-glow btn-sm" onclick="taskAction('${t.id}', 'open')">${iconSvg("ic-external")}<span>${lblOpen}</span></button>`;
        actionBtns += `<button class="glass-btn-icon" onclick="taskAction('${t.id}', 'reveal')" title="${lblFolderTip}">${iconSvg("ic-folder")}</button>`;
      }
      actionBtns += `<button class="glass-btn-icon" onclick="copyTaskUrl('${t.id}')" title="${lblCopyLinkTip}">${iconSvg("ic-link")}</button>`;
      actionBtns += `<button class="glass-btn-icon" onclick="openInspectorFor('${t.id}')" title="${lblDetailsTip}">${iconSvg("ic-activity")}</button>`;
      actionBtns += `<button class="glass-btn-icon" onclick="taskAction('${t.id}', 'delete', false)" title="${lblRemoveTip}">${iconSvg("ic-trash")}</button>`;
      actionBtns += `<button class="glass-btn btn-sm btn-danger-file" onclick="taskAction('${t.id}', 'delete', true)" title="${lblDeleteFileTip}">${iconSvg("ic-trash")}<span>${lblDeleteFile}</span></button>`;
      r.actions.innerHTML = actionBtns;
    }

    const rawSegs =
      t.segments && t.segments.length > 0
        ? t.segments
        : [{ start: 0, current: t.downloaded, end: t.totalSize || 100, status: t.status, generation: 0 }];
    const segs = rawSegs.slice(0, 32);

    if (card.dataset.segCount !== String(segs.length)) {
      card.dataset.segCount = String(segs.length);
      r.ribbon.innerHTML = segs
        .map(() => `<div class="seg-slice"><div class="seg-slice-fill"></div></div>`)
        .join("");
    }

    const sliceNodes = r.ribbon.children;
    for (let i = 0; i < segs.length; i++) {
      const s = segs[i];
      const sliceEl = sliceNodes[i];
      if (!sliceEl) continue;
      const fillEl = sliceEl.firstElementChild;
      const span = Math.max(1, s.end - s.start + 1);
      const done = Math.max(0, Math.min(span, s.current - s.start));
      const fillPct =
        s.end >= 0
          ? Math.min(100, (done / span) * 100)
          : t.status === "completed"
          ? 100
          : 50;

      const flexVal = String(Math.max(1, span));
      if (sliceEl.style.flex !== flexVal) sliceEl.style.flex = flexVal;

      const widthStr = `${fillPct.toFixed(1)}%`;
      if (fillEl.style.width !== widthStr) fillEl.style.width = widthStr;

      const desiredFillCls =
        "seg-slice-fill" +
        (s.status === "completed" ? " done" : s.generation > 0 ? " split-gen" : "");
      if (fillEl.className !== desiredFillCls) fillEl.className = desiredFillCls;
    }
  });
}

function renderInspectorZeroFlicker(t) {
  const isEn = state.language === "en";
  const isJa = state.language === "ja";
  setTextIfChanged(el.inspFilename, t.fileName);
  setTextIfChanged(el.inspSpeed, formatSpeed(t.speedBps));
  setTextIfChanged(el.inspEta, formatETA(t.etaSeconds));
  const totalSegs = t.segments ? t.segments.length : t.maxConnections;
  setTextIfChanged(el.inspSegmentsCount, `${t.activeWorkers} / ${totalSegs}`);
  setTextIfChanged(
    el.inspSplits,
    isEn ? `${t.dynamicSplits || 0}` : isJa ? `${t.dynamicSplits || 0}回` : `${t.dynamicSplits || 0}회`
  );
  setTextIfChanged(
    el.inspPeakAvg,
    `${formatSpeed(t.peakSpeedBps)} / ${formatSpeed(t.avgSpeedBps)}`
  );

  drawSmoothAreaChart(
    el.inspSpeedCanvas,
    t.speedHistory,
    "#0A84FF",
    "rgba(10, 132, 255, 0.22)"
  );

  const pct = (t.progress || 0).toFixed(1);
  setTextIfChanged(
    el.inspBytesProgress,
    `${formatBytes(t.downloaded)} / ${formatBytes(t.totalSize)} (${pct}%)`
  );

  const total = t.totalSize > 0 ? t.totalSize : 1;
  const segs = t.segments || [];
  const structKey = `${t.id}:${segs.length}`;

  if (el.inspSegmentsList.dataset.structKey !== structKey) {
    el.inspSegmentsList.dataset.structKey = structKey;
    el.inspByteMap.innerHTML = segs.map(() => `<div class="byte-block"></div>`).join("");
    el.inspSegmentsList.innerHTML = segs
      .map(
        (s) => `
      <div class="seg-mini-card tabular">
        <div class="seg-mini-top">
          <span data-srole="title">#${s.id}</span>
          <span data-srole="speed"></span>
        </div>
        <div class="seg-mini-bar">
          <div class="seg-mini-fill" data-srole="fill"></div>
        </div>
        <div class="seg-mini-meta">
          <span data-srole="bytes"></span>
          <span data-srole="pct"></span>
        </div>
      </div>`
      )
      .join("");
  }

  const byteBlocks = el.inspByteMap.children;
  for (let i = 0; i < segs.length; i++) {
    const s = segs[i];
    const blk = byteBlocks[i];
    if (!blk) continue;
    if (s.end < 0) {
      blk.style.left = "0%";
      blk.style.width = `${pct}%`;
    } else {
      const leftPct = (s.start / total) * 100;
      const downloadedLen = Math.max(0, s.current - s.start);
      const widthPct = (downloadedLen / total) * 100;
      blk.style.left = `${leftPct.toFixed(2)}%`;
      blk.style.width = `${Math.max(0.4, widthPct).toFixed(2)}%`;
    }
    const cls =
      "byte-block" +
      (s.status === "completed" ? " completed" : s.generation > 0 ? " split" : "");
    if (blk.className !== cls) blk.className = cls;
  }

  const segCards = el.inspSegmentsList.children;
  for (let i = 0; i < segs.length; i++) {
    const s = segs[i];
    const cardNode = segCards[i];
    if (!cardNode) continue;

    const span = Math.max(1, s.end - s.start + 1);
    const done = Math.max(0, Math.min(span, s.current - s.start));
    const sPct =
      s.end >= 0
        ? Math.min(100, (done / span) * 100)
        : t.status === "completed"
        ? 100
        : 50;

    const titleNode = cardNode.querySelector('[data-srole="title"]');
    const speedNode = cardNode.querySelector('[data-srole="speed"]');
    const fillNode = cardNode.querySelector('[data-srole="fill"]');
    const bytesNode = cardNode.querySelector('[data-srole="bytes"]');
    const pctNode = cardNode.querySelector('[data-srole="pct"]');

    const titleTxt = s.generation > 0 ? `#${s.id} G${s.generation}` : `#${s.id}`;
    setTextIfChanged(titleNode, titleTxt);

    const doneStr = isEn ? "Done" : isJa ? "完了" : "완료";
    const speedStr = s.status === "completed" ? doneStr : formatSpeed(s.speedBps);
    setTextIfChanged(speedNode, speedStr);

    const fillColor =
      s.status === "completed"
        ? "var(--accent-emerald)"
        : s.generation > 0
        ? "var(--accent-purple)"
        : "var(--accent-primary)";
    speedNode.style.color = fillColor;
    fillNode.style.width = `${sPct.toFixed(1)}%`;
    fillNode.style.background = fillColor;

    setTextIfChanged(bytesNode, `${formatBytes(done)}`);
    setTextIfChanged(pctNode, `${sPct.toFixed(0)}%`);
  }

  if (t.checksum) {
    setTextIfChanged(el.inspSha256, t.checksum.sha256 || "-");
    setTextIfChanged(el.inspMd5, t.checksum.md5 || "-");
    if (t.checksum.verified === true) {
      setTextIfChanged(el.inspHashBadge, isEn ? "Verified" : isJa ? "検証済み" : "일치 검증됨");
      el.inspHashBadge.style.color = "var(--accent-emerald)";
    } else if (t.checksum.verified === false) {
      setTextIfChanged(el.inspHashBadge, isEn ? "Mismatch" : isJa ? "不一致" : "불일치");
      el.inspHashBadge.style.color = "var(--accent-rose)";
    } else {
      setTextIfChanged(el.inspHashBadge, isEn ? "Computed" : isJa ? "計算完了" : "계산 완료");
      el.inspHashBadge.style.color = "var(--accent-primary)";
    }
  } else {
    setTextIfChanged(
      el.inspSha256,
      t.status === "verifying"
        ? isEn ? "Verifying..." : isJa ? "検証中..." : "검증 중..."
        : isEn ? "Pending" : isJa ? "完了後に表示" : "완료 후 표시"
    );
    setTextIfChanged(el.inspMd5, "-");
    setTextIfChanged(el.inspHashBadge, isEn ? "Pending" : isJa ? "待機" : "대기");
    el.inspHashBadge.style.color = "var(--text-muted)";
  }

  setTextIfChanged(
    el.inspServerTag,
    t.server || (t.acceptRanges ? "HTTP 206 Range" : "Single")
  );

  const logs = (t.logs || []).slice(-25).reverse();
  const logKey = `${t.id}:${logs.length}:${logs[0] ? logs[0].time : ""}`;
  if (el.inspLogs.dataset.logKey !== logKey) {
    el.inspLogs.dataset.logKey = logKey;
    el.inspLogs.innerHTML = logs
      .map(
        (l) => `
      <div class="log-line">
        <span class="log-time">[${l.time}]</span>
        <span class="log-msg ${l.level}">${l.message}</span>
      </div>`
      )
      .join("");
  }
}

window.openInspectorFor = function(id) {
  if (state.selectedTaskId === id && !el.inspectorDrawer.classList.contains("hidden")) {
    el.inspectorDrawer.classList.add("hidden");
    state.selectedTaskId = null;
    return;
  }
  state.selectedTaskId = id;
  el.inspectorDrawer.classList.remove("hidden");
  if (state.snapshot && state.snapshot.tasks) {
    const t = state.snapshot.tasks.find((x) => x.id === id);
    if (t) renderInspectorZeroFlicker(t);
  }
};

window.copyTaskUrl = function(id) {
  if (!state.snapshot || !state.snapshot.tasks) return;
  const t = state.snapshot.tasks.find((x) => x.id === id);
  if (!t) return;
  const link = t.finalUrl || t.url || "";
  if (!link) return;
  navigator.clipboard.writeText(link);
  showToast(
    state.language === "en"
      ? "Download link copied to clipboard."
      : state.language === "ja"
      ? "ダウンロードリンクをコピーしました。"
      : "다운로드 링크가 클립보드에 복사되었습니다."
  );
};

window.taskAction = async function(id, action, deleteFiles = false) {
  await fetch("/api/tasks/action", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id, action, deleteFiles: !!deleteFiles }),
  });
  if (action === "delete" && deleteFiles) {
    showToast(
      state.language === "en"
        ? "Task and downloaded file deleted."
        : state.language === "ja"
        ? "ファイルとタスクを完全に削除しました。"
        : "목록 및 다운로드된 파일이 함께 삭제되었습니다."
    );
  }
};

function updateSliderRecTag() {
  const val = parseInt(el.advConnections.value, 10);
  const isEn = state.language === "en";
  const isJa = state.language === "ja";
  if (val === state.recommendedThreads) {
    setTextIfChanged(
      el.advConnRecTag,
      isEn ? "Optimal Spec" : isJa ? "スペック最適値" : "사양 최적값"
    );
    el.advConnRecTag.style.color = "var(--accent-emerald)";
  } else {
    setTextIfChanged(
      el.advConnRecTag,
      isEn
        ? `Rec: ${state.recommendedThreads} threads`
        : isJa
        ? `推奨: ${state.recommendedThreads} スレッド`
        : `추천: ${state.recommendedThreads} 스레드`
    );
    el.advConnRecTag.style.color = "var(--accent-purple)";
  }
}

function applyRecommendedThreads(count, note) {
  const target = count || state.recommendedThreads || 16;
  el.advConnections.value = target;
  setTextIfChanged(el.advConnVal, target);
  updateSliderRecTag();
  const isEn = state.language === "en";
  const isJa = state.language === "ja";
  showToast(
    note ||
      (isEn
        ? `Applied hardware-optimal ${target} threads`
        : isJa
        ? `スペック最適スレッド(${target}個)を適用しました`
        : `사양 맞춤 최적 스레드(${target}개) 적용됨`)
  );
}

async function persistPartialSettings(partial) {
  const current = (state.snapshot && state.snapshot.settings) || {};
  const payload = {
    ...current,
    theme: state.theme,
    language: state.language,
    ...partial,
  };
  try {
    const res = await fetch("/api/settings", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (res.ok && state.snapshot) {
      state.snapshot.settings = await res.json();
    }
  } catch (_) {}
}

async function probeCurrentUrl() {
  const url = el.urlInput.value.trim();
  if (!url) return;
  el.btnProbeUrl.textContent = tStr("composer.probing", "분석 중");
  try {
    const res = await fetch("/api/probe", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url, headers: collectCustomHeaders() }),
    });
    const data = await res.json();
    if (res.ok) {
      el.probePreviewBar.classList.remove("hidden");
      setTextIfChanged(el.probeFilename, data.fileName);
      setTextIfChanged(el.probeSize, formatBytes(data.totalSize));
      const isEn = state.language === "en";
      const isJa = state.language === "ja";
      setTextIfChanged(
        el.probeRange,
        data.acceptRanges
          ? isEn ? "Multi-Segment Ready" : isJa ? "マルチセグメント対応" : "멀티 세그먼트 지원"
          : isEn ? "Single Stream" : isJa ? "単一ストリーム" : "단일 스트림"
      );

      // Keep recommendedThreads strictly matched to the hardware profile's recommendedConnections
      const hwThreads =
        (state.snapshot && state.snapshot.hardware && state.snapshot.hardware.recommendedConnections) ||
        data.recommendedThreads ||
        state.recommendedThreads ||
        16;
      state.recommendedThreads = hwThreads;
      setTextIfChanged(el.probeRecThreads, `${hwThreads} ${threadUnit()}`);
      setTextIfChanged(el.btnRecCountLabel, hwThreads);
      el.advConnections.value = hwThreads;
      setTextIfChanged(el.advConnVal, hwThreads);
      updateSliderRecTag();

      if (!el.advFilename.value) {
        el.advFilename.placeholder = data.fileName;
      }
    } else {
      showToast(`Error: ${data.error || "Failed"}`);
    }
  } catch (e) {
    showToast("Failed to probe URL metadata.");
  } finally {
    el.btnProbeUrl.textContent = tStr("composer.probe", "분석");
  }
}

function collectCustomHeaders() {
  const h = {};
  const ref = el.advReferer.value.trim();
  const cook = el.advCookie.value.trim();
  if (ref) h["Referer"] = ref;
  if (cook) {
    if (cook.toLowerCase().startsWith("bearer ") || cook.toLowerCase().startsWith("basic ")) {
      h["Authorization"] = cook;
    } else {
      h["Cookie"] = cook;
    }
  }
  return h;
}

async function startDownloadFromInput(overrideUrl, overrideConn, overrideName) {
  const url = overrideUrl || el.urlInput.value.trim();
  if (!url) {
    showToast(
      state.language === "en"
        ? "Please enter a download URL."
        : state.language === "ja"
        ? "ダウンロードするURLを入力してください。"
        : "다운로드할 URL을 입력해주세요."
    );
    el.urlInput.focus();
    return;
  }

  const maxConnections =
    overrideConn || parseInt(el.advConnections.value, 10) || state.recommendedThreads || 16;
  const fileName = overrideName || el.advFilename.value.trim();
  const saveDir = el.advSavedir.value.trim();
  const expectedHash = el.advHash.value.trim();

  const res = await fetch("/api/tasks", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      url,
      fileName,
      saveDir,
      maxConnections,
      expectedHash,
      headers: collectCustomHeaders(),
    }),
  });
  const created = await res.json();
  if (res.ok && created.id) {
    el.urlInput.value = "";
    el.advFilename.value = "";
    el.advHash.value = "";
    el.probePreviewBar.classList.add("hidden");
    showToast(`${created.fileName} (${created.maxConnections} ${threadUnit()})`);
  } else {
    showToast(`${created.error || "Error"}`);
  }
}

async function launchHoverWidget() {
  try {
    const currentUrl = el.urlInput ? el.urlInput.value.trim() : "";
    const res = await fetch("/api/hover/launch", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url: currentUrl }),
    });
    if (res.ok) {
      showToast(
        state.language === "en"
          ? "Opened New Download dialog."
          : state.language === "ja"
          ? "新規ダウンロードダイアログを開きました。"
          : "새 다운로드 추가 창이 열렸습니다."
      );
    }
  } catch (_) {}
}

function initEvents() {
  document.querySelectorAll(".nav-item").forEach((btn) => {
    btn.addEventListener("click", () => {
      document.querySelectorAll(".nav-item").forEach((b) => b.classList.remove("active"));
      btn.classList.add("active");
      state.activeFilter = btn.dataset.filter;
      if (state.snapshot) renderGlobalState(state.snapshot);
    });
  });

  el.sidebarHwCard.addEventListener("click", () => {
    el.composerAdvanced.classList.remove("hidden");
    el.advChevron.textContent = "▴";
    applyRecommendedThreads(state.recommendedThreads);
  });

  el.btnApplyRecThreads.addEventListener("click", () => {
    applyRecommendedThreads(state.recommendedThreads);
  });

  document.getElementById("btn-set-apply-hw").addEventListener("click", () => {
    if (state.snapshot && state.snapshot.hardware) {
      const hw = state.snapshot.hardware;
      document.getElementById("set-connections").value = hw.recommendedConnections;
      document.getElementById("set-concurrent").value = hw.recommendedConcurrent;
      showToast(`${hw.recommendedConnections} ${threadUnit()}`);
    }
  });

  el.searchInput.addEventListener("input", () => {
    state.searchQuery = el.searchInput.value.trim();
    if (state.snapshot) renderGlobalState(state.snapshot);
  });

  el.urlInput.addEventListener("paste", () => {
    setTimeout(() => probeCurrentUrl(), 80);
  });
  el.urlInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") startDownloadFromInput();
  });
  el.btnProbeUrl.addEventListener("click", probeCurrentUrl);
  el.btnAddDownload.addEventListener("click", () => startDownloadFromInput());

  el.btnToggleAdvanced.addEventListener("click", () => {
    const hidden = el.composerAdvanced.classList.toggle("hidden");
    el.advChevron.textContent = hidden ? "▾" : "▴";
  });

  el.advConnections.addEventListener("input", () => {
    setTextIfChanged(el.advConnVal, el.advConnections.value);
    updateSliderRecTag();
  });

  el.btnPickAdvFolder.addEventListener("click", async () => {
    const curDir = el.advSavedir.value.trim() || (state.snapshot && state.snapshot.settings && state.snapshot.settings.defaultSaveDir) || "";
    const res = await fetch("/api/dialog/folder", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ initialDir: curDir, saveAsDefault: false }),
    });
    const data = await res.json();
    if (data.path) {
      el.advSavedir.value = data.path;
      showToast(
        state.language === "en"
          ? `Selected folder: ${data.path}`
          : state.language === "ja"
          ? `保存先を選択: ${data.path}`
          : `저장 위치 선택됨: ${data.path}`
      );
    }
  });

  document.getElementById("btn-set-pick-folder").addEventListener("click", async () => {
    const setDirInput = document.getElementById("set-savedir");
    const curDir = setDirInput.value.trim();
    const res = await fetch("/api/dialog/folder", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ initialDir: curDir, saveAsDefault: true }),
    });
    const data = await res.json();
    if (data.path) {
      setDirInput.value = data.path;
      el.advSavedir.placeholder = data.path;
      if (state.snapshot && state.snapshot.settings) {
        state.snapshot.settings.defaultSaveDir = data.path;
      }
      showToast(
        state.language === "en"
          ? `Default folder updated: ${data.path}`
          : state.language === "ja"
          ? `既定の保存先を変更しました: ${data.path}`
          : `기본 저장 폴더가 변경되었습니다: ${data.path}`
      );
    }
  });

  const btnSetOpenFolder = document.getElementById("btn-set-open-folder");
  if (btnSetOpenFolder) {
    btnSetOpenFolder.addEventListener("click", async () => {
      const targetPath = document.getElementById("set-savedir").value.trim();
      await fetch("/api/folder/open", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path: targetPath }),
      });
    });
  }

  const btnSetResetFolder = document.getElementById("btn-set-reset-folder");
  if (btnSetResetFolder) {
    btnSetResetFolder.addEventListener("click", async () => {
      const res = await fetch("/api/folder/default");
      const data = await res.json();
      if (data.path) {
        document.getElementById("set-savedir").value = data.path;
        el.advSavedir.placeholder = data.path;
        await persistPartialSettings({ defaultSaveDir: data.path });
        showToast(
          state.language === "en"
            ? `Reset to default folder: ${data.path}`
            : state.language === "ja"
            ? `初期フォルダに戻しました: ${data.path}`
            : `기본 다운로드 폴더로 초기화되었습니다: ${data.path}`
        );
      }
    });
  }

  el.btnOpenSaveDir.addEventListener("click", () => {
    window.taskAction("", "open-default-dir");
  });

  el.btnClipDownload.addEventListener("click", () => {
    const url = el.clipboardDetectedUrl.textContent;
    if (url) startDownloadFromInput(url);
  });
  el.btnClipDismiss.addEventListener("click", async () => {
    await fetch("/api/clipboard/dismiss", { method: "POST" });
    el.clipboardBanner.classList.add("hidden");
  });

  el.quickSpeedLimit.addEventListener("change", async () => {
    const limit = parseInt(el.quickSpeedLimit.value, 10) || 0;
    await persistPartialSettings({ globalSpeedLimitBps: limit });
    showToast(limit === 0 ? tStr("telemetry.unlimited") : `${formatSpeed(limit)}`);
  });

  document.getElementById("btn-launch-hover-widget").addEventListener("click", launchHoverWidget);
  document.getElementById("btn-modal-launch-hover").addEventListener("click", launchHoverWidget);

  document.getElementById("btn-resume-all").addEventListener("click", () => {
    window.taskAction("", "resume-all");
  });
  document.getElementById("btn-pause-all").addEventListener("click", () => {
    window.taskAction("", "pause-all");
  });
  document.getElementById("btn-clear-completed").addEventListener("click", () => {
    window.taskAction("", "clear-completed");
  });

  // Two-theme toggle (Black <-> White) persisted to backend settings + localStorage
  document.getElementById("btn-theme-cycle").addEventListener("click", async () => {
    const nextTheme = state.theme === "black" ? "white" : "black";
    applyThemeToDOM(nextTheme);
    await persistPartialSettings({ theme: nextTheme });
    const label = nextTheme === "white" ? tStr("theme.white") : tStr("theme.black");
    showToast(label);
  });

  // Header language selector persisted to backend settings + localStorage
  const headerLangSelect = document.getElementById("header-lang-select");
  if (headerLangSelect) {
    headerLangSelect.addEventListener("change", async () => {
      const nextLang = headerLangSelect.value;
      applyLanguageToDOM(nextLang);
      await persistPartialSettings({ language: nextLang });
    });
  }

  // One-click copy buttons for all supported browsers (Chrome, Edge, Whale, Brave, Opera, Firefox)
  document.querySelectorAll("[data-ext-url]").forEach((btn) => {
    btn.addEventListener("click", () => {
      const extUrl = btn.getAttribute("data-ext-url");
      if (!extUrl) return;
      navigator.clipboard.writeText(extUrl);
      showToast(
        state.language === "en"
          ? `Copied ${extUrl} — paste it in your browser address bar!`
          : state.language === "ja"
          ? `${extUrl} をコピーしました。ブラウザのアドレスバーに貼り付けてください。`
          : `${extUrl} 주소가 복사되었습니다. 브라우저 주소창에 붙여넣으세요!`
      );
    });
  });

  el.btnCloseInspector.addEventListener("click", () => {
    el.inspectorDrawer.classList.add("hidden");
    state.selectedTaskId = null;
  });

  [el.inspSha256, el.inspMd5].forEach((node) => {
    node.addEventListener("click", () => {
      const txt = node.textContent.trim();
      if (txt && txt.length >= 32 && !txt.includes(" ")) {
        navigator.clipboard.writeText(txt);
        showToast("Copied checksum.");
      }
    });
  });

  const openModal = (id) => document.getElementById(id).classList.remove("hidden");
  const closeModal = (id) => document.getElementById(id).classList.add("hidden");

  document.getElementById("btn-browser-modal").addEventListener("click", () => openModal("modal-browser"));
  document.getElementById("btn-empty-open-browser").addEventListener("click", () => openModal("modal-browser"));
  document.getElementById("btn-benchmark-modal").addEventListener("click", () => openModal("modal-benchmark"));
  document.getElementById("btn-sniffer-modal").addEventListener("click", () => openModal("modal-sniffer"));
  document.getElementById("btn-batch-modal").addEventListener("click", () => openModal("modal-batch"));
  document.getElementById("btn-settings-modal").addEventListener("click", () => {
    if (state.snapshot && state.snapshot.settings) {
      const s = state.snapshot.settings;
      document.getElementById("set-savedir").value = s.defaultSaveDir || "";
      document.getElementById("set-connections").value = s.defaultConnections || state.recommendedThreads;
      document.getElementById("set-concurrent").value = s.maxConcurrentTasks || 4;
      document.getElementById("set-speed-limit").value = String(s.globalSpeedLimitBps || 0);
      document.getElementById("set-autocat").checked = !!s.autoCategorize;
      document.getElementById("set-clipboard").checked = !!s.clipboardMonitor;
      document.getElementById("set-browser-intercept").checked = s.browserIntercept !== false;
      document.getElementById("set-sound").checked = !!s.soundEffects;
      document.getElementById("set-startup").checked = !!s.runOnStartup;
      document.getElementById("set-theme").value = s.theme === "white" ? "white" : "black";
      document.getElementById("set-language").value = ["ko", "en", "ja"].includes(s.language) ? s.language : state.language;
    }
    openModal("modal-settings");
  });

  // Live preview & immediate persistence when changing Theme or Language inside Settings modal
  document.getElementById("set-theme").addEventListener("change", async (e) => {
    const nextTheme = e.target.value === "white" ? "white" : "black";
    applyThemeToDOM(nextTheme);
    await persistPartialSettings({ theme: nextTheme });
  });

  document.getElementById("set-language").addEventListener("change", async (e) => {
    const nextLang = e.target.value || "ko";
    applyLanguageToDOM(nextLang);
    await persistPartialSettings({ language: nextLang });
  });

  const btnQuitApp = document.getElementById("btn-quit-app");
  if (btnQuitApp) {
    btnQuitApp.addEventListener("click", async () => {
      try {
        await fetch("/api/app/quit", { method: "POST" });
      } catch (_) {}
      window.close();
    });
  }

  document.querySelectorAll("[data-close]").forEach((btn) => {
    btn.addEventListener("click", () => closeModal(btn.dataset.close));
  });

  document.getElementById("btn-open-ext-folder").addEventListener("click", async () => {
    const res = await fetch("/api/browser/open-extension-folder", { method: "POST" });
    const data = await res.json();
    showToast(`${data.path || "browser-extension"}`);
  });

  document.getElementById("btn-register-protocol").addEventListener("click", async () => {
    const res = await fetch("/api/browser/register-protocol", { method: "POST" });
    const data = await res.json();
    if (data.ok) {
      setTextIfChanged(el.protoStatusText, "Registered");
      showToast("loadin:// protocol registered.");
    } else {
      showToast(`${data.error || "Error"}`);
    }
  });

  document.getElementById("btn-copy-bookmarklet").addEventListener("click", () => {
    const bm = document.getElementById("bookmarklet-anchor").getAttribute("href");
    navigator.clipboard.writeText(bm);
    showToast("Bookmarklet copied.");
  });

  document.getElementById("btn-empty-start-bench").addEventListener("click", () => {
    const benchUrl = `${window.location.origin}/api/sandbox-stream/Loadin_DynamicSplit_64MB.bin?sizeMB=64`;
    startDownloadFromInput(benchUrl, state.recommendedThreads, "Loadin_DynamicSplit_64MB.bin");
  });

  document.querySelectorAll("[data-bench]").forEach((btn) => {
    btn.addEventListener("click", () => {
      const type = btn.dataset.bench;
      closeModal("modal-benchmark");
      if (type === "sandbox-64") {
        const u = `${window.location.origin}/api/sandbox-stream/Loadin_DynamicSplit_64MB.bin?sizeMB=64`;
        startDownloadFromInput(u, state.recommendedThreads, "Loadin_DynamicSplit_64MB.bin");
      } else if (type === "sandbox-256") {
        const u = `${window.location.origin}/api/sandbox-stream/Loadin_HyperStress_256MB.iso?sizeMB=256`;
        startDownloadFromInput(u, 32, "Loadin_HyperStress_256MB.iso");
      } else if (type === "external-sdk") {
        startDownloadFromInput(
          "https://dl.google.com/go/go1.23.2.windows-amd64.zip",
          state.recommendedThreads,
          "sdk_archive_1.23.2_amd64.zip"
        );
      } else if (type === "external-pkg") {
        startDownloadFromInput(
          "https://github.com/git-for-windows/git/releases/download/v2.47.0.windows.2/PortableGit-2.47.0.2-64-bit.7z.exe",
          state.recommendedThreads,
          "PortableDevSuite_64bit.7z.exe"
        );
      }
    });
  });

  document.getElementById("btn-run-sniffer").addEventListener("click", async () => {
    const url = document.getElementById("sniffer-url-input").value.trim() || "https://go.dev/dl/";
    document.getElementById("sniffer-url-input").value = url;
    const container = document.getElementById("sniffer-results");
    container.innerHTML = `<div class="sniffer-placeholder">Scanning links...</div>`;
    try {
      const res = await fetch("/api/sniff", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });
      const data = await res.json();
      state.sniffedLinks = (data.links || []).map((l) => ({ ...l, checked: true }));
      renderSniffedLinks();
    } catch (e) {
      container.innerHTML = `<div class="sniffer-placeholder">Failed to scan URL.</div>`;
    }
  });

  document.querySelectorAll("[data-sfilter]").forEach((chip) => {
    chip.addEventListener("click", () => {
      document.querySelectorAll("[data-sfilter]").forEach((c) => c.classList.remove("active"));
      chip.classList.add("active");
      state.sniffFilter = chip.dataset.sfilter;
      renderSniffedLinks();
    });
  });

  document.getElementById("btn-sniff-select-all").addEventListener("click", () => {
    state.sniffedLinks.forEach((l) => (l.checked = true));
    renderSniffedLinks();
  });
  document.getElementById("btn-sniff-deselect-all").addEventListener("click", () => {
    state.sniffedLinks.forEach((l) => (l.checked = false));
    renderSniffedLinks();
  });

  document.getElementById("btn-sniff-download-selected").addEventListener("click", async () => {
    const selected = state.sniffedLinks.filter((l) => l.checked).map((l) => l.url);
    if (selected.length === 0) return;
    await fetch("/api/tasks/batch", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ urls: selected, maxConnections: state.recommendedThreads }),
    });
    closeModal("modal-sniffer");
    showToast(`Added ${selected.length} files.`);
  });

  document.getElementById("btn-submit-batch").addEventListener("click", async () => {
    const raw = document.getElementById("batch-urls-input").value.trim();
    if (!raw) return;
    const lines = raw
      .split(/\r?\n/)
      .map((s) => s.trim())
      .filter(Boolean);
    const conns = parseInt(document.getElementById("batch-connections").value, 10) || state.recommendedThreads;
    const res = await fetch("/api/tasks/batch", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ urls: lines, maxConnections: conns }),
    });
    const data = await res.json();
    closeModal("modal-batch");
    document.getElementById("batch-urls-input").value = "";
    showToast(`Added ${data.added || lines.length} tasks.`);
  });

  document.getElementById("btn-save-settings").addEventListener("click", async () => {
    const current = (state.snapshot && state.snapshot.settings) || {};
    const chosenTheme = document.getElementById("set-theme").value === "white" ? "white" : "black";
    const chosenLang = document.getElementById("set-language").value || "ko";
    const chosenSpeedLimit = parseInt(document.getElementById("set-speed-limit").value, 10) || 0;
    const chosenSaveDir = document.getElementById("set-savedir").value.trim();

    applyThemeToDOM(chosenTheme);
    applyLanguageToDOM(chosenLang);

    if (el.quickSpeedLimit) {
      el.quickSpeedLimit.value = String(chosenSpeedLimit);
    }
    if (chosenSaveDir) {
      el.advSavedir.placeholder = chosenSaveDir;
    }

    const payload = {
      ...current,
      defaultSaveDir: chosenSaveDir,
      defaultConnections: parseInt(document.getElementById("set-connections").value, 10) || state.recommendedThreads,
      maxConcurrentTasks: parseInt(document.getElementById("set-concurrent").value, 10) || 4,
      globalSpeedLimitBps: chosenSpeedLimit,
      autoCategorize: document.getElementById("set-autocat").checked,
      clipboardMonitor: document.getElementById("set-clipboard").checked,
      browserIntercept: document.getElementById("set-browser-intercept").checked,
      soundEffects: document.getElementById("set-sound").checked,
      runOnStartup: document.getElementById("set-startup").checked,
      theme: chosenTheme,
      language: chosenLang,
    };
    const res = await fetch("/api/settings", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (res.ok && state.snapshot) {
      state.snapshot.settings = await res.json();
    }
    closeModal("modal-settings");
    showToast(
      chosenLang === "en"
        ? "Settings saved."
        : chosenLang === "ja"
        ? "設定を保存しました。"
        : "환경설정이 저장되었습니다."
    );
  });

  window.addEventListener("keydown", (e) => {
    const key = e.key.toLowerCase();
    // Block browser-specific shortcuts (F12, Ctrl+Shift+I/J/C, Ctrl+U, Ctrl+P, F5, Ctrl+R)
    if (
      e.key === "F12" ||
      e.key === "F5" ||
      ((e.ctrlKey || e.metaKey) && e.shiftKey && ["i", "j", "c"].includes(key)) ||
      ((e.ctrlKey || e.metaKey) && ["u", "p", "r"].includes(key))
    ) {
      e.preventDefault();
      return;
    }

    if ((e.ctrlKey || e.metaKey) && key === "f") {
      e.preventDefault();
      el.searchInput.focus();
    } else if (e.key === "Escape") {
      document.querySelectorAll(".modal-backdrop").forEach((m) => m.classList.add("hidden"));
      const ctxMenu = document.getElementById("app-context-menu");
      if (ctxMenu) ctxMenu.classList.add("hidden");
    }
  });

  initCustomContextMenu();
}

function initCustomContextMenu() {
  const menu = document.getElementById("app-context-menu");
  if (!menu) return;

  const hideMenu = () => menu.classList.add("hidden");

  window.addEventListener("click", (e) => {
    if (!menu.contains(e.target)) hideMenu();
  });
  window.addEventListener("blur", hideMenu);
  window.addEventListener("resize", hideMenu);
  window.addEventListener("scroll", hideMenu, true);

  window.addEventListener("contextmenu", (e) => {
    e.preventDefault();

    const isEn = state.language === "en";
    const isJa = state.language === "ja";
    const L = (ko, en, ja) => (isEn ? en : isJa ? ja : ko);

    const items = [];
    const inputEl = e.target.closest("input[type='text'], input[type='url'], input[type='search'], input[type='number'], textarea");
    const taskCard = e.target.closest(".task-card");

    if (inputEl) {
      items.push({
        icon: "ic-copy",
        label: L("붙여넣기", "Paste", "貼り付け"),
        action: async () => {
          try {
            const text = await navigator.clipboard.readText();
            if (text) {
              const start = inputEl.selectionStart ?? inputEl.value.length;
              const end = inputEl.selectionEnd ?? inputEl.value.length;
              inputEl.value = inputEl.value.slice(0, start) + text + inputEl.value.slice(end);
              inputEl.dispatchEvent(new Event("input", { bubbles: true }));
              inputEl.focus();
            }
          } catch (_) {}
        },
      });
      if (inputEl.id === "url-input") {
        items.push({
          icon: "ic-download",
          label: L("붙여넣고 바로 다운로드", "Paste & Download Now", "貼り付けてすぐダウンロード"),
          action: async () => {
            try {
              const text = (await navigator.clipboard.readText()).trim();
              if (text) {
                inputEl.value = text;
                startDownloadFromInput(text);
              }
            } catch (_) {}
          },
        });
      }
      items.push({
        icon: "ic-link",
        label: L("복사", "Copy", "コピー"),
        action: () => {
          const sel = inputEl.value.substring(inputEl.selectionStart || 0, inputEl.selectionEnd || 0) || inputEl.value;
          if (sel) navigator.clipboard.writeText(sel);
        },
      });
      items.push({
        icon: "ic-check",
        label: L("전체 선택", "Select All", "すべて選択"),
        action: () => {
          inputEl.focus();
          inputEl.select();
        },
      });
    } else if (taskCard && taskCard.dataset.id) {
      const tid = taskCard.dataset.id;
      const t = (state.snapshot && state.snapshot.tasks || []).find((x) => x.id === tid);
      if (t) {
        if (t.status === "completed") {
          items.push({
            icon: "ic-external",
            label: L("파일 열기", "Open File", "ファイルを開く"),
            action: () => window.taskAction(tid, "open", false),
          });
          items.push({
            icon: "ic-folder",
            label: L("저장 폴더 열기", "Show in Folder", "保存フォルダを開く"),
            action: () => window.taskAction(tid, "reveal", false),
          });
        } else if (t.status === "running" || t.status === "probing") {
          items.push({
            icon: "ic-pause",
            label: L("일시정지", "Pause Download", "一時停止"),
            action: () => window.taskAction(tid, "pause", false),
          });
        } else {
          items.push({
            icon: "ic-play",
            label: L("이어받기", "Resume Download", "再開"),
            action: () => window.taskAction(tid, "resume", false),
          });
        }

        items.push({
          icon: "ic-link",
          label: L("다운로드 링크 복사", "Copy Download Link", "ダウンロードリンクをコピー"),
          action: () => window.copyTaskUrl(tid),
        });
        items.push({
          icon: "ic-activity",
          label: L("상세 정보 보기", "View Task Details", "詳細情報を表示"),
          action: () => window.openInspectorFor(tid),
        });
        items.push({ divider: true });
        items.push({
          icon: "ic-trash",
          label: L("목록에서 제거 (파일 유지)", "Remove from List (Keep File)", "リストから削除 (ファイル保持)"),
          action: () => window.taskAction(tid, "delete", false),
        });
        items.push({
          icon: "ic-trash",
          label: L("파일까지 바로 삭제", "Delete Task & File", "ファイルごと完全に削除"),
          danger: true,
          action: () => window.taskAction(tid, "delete", true),
        });
      }
    } else {
      items.push({
        icon: "ic-plus",
        label: L("새 다운로드 추가 창 열기", "New Download Dialog", "新規ダウンロード追加"),
        action: () => launchHoverWidget(),
      });
      items.push({
        icon: "ic-play",
        label: L("전체 이어받기", "Resume All Downloads", "すべて再開"),
        action: () => window.taskAction("", "resume-all", false),
      });
      items.push({
        icon: "ic-pause",
        label: L("전체 일시정지", "Pause All Downloads", "すべて一時停止"),
        action: () => window.taskAction("", "pause-all", false),
      });
      items.push({
        icon: "ic-check",
        label: L("완료 항목 정리", "Clear Completed Tasks", "完了項目を整理"),
        action: () => window.taskAction("", "clear-completed", false),
      });
      items.push({ divider: true });
      items.push({
        icon: "ic-folder",
        label: L("기본 다운로드 폴더 열기", "Open Downloads Folder", "保存先フォルダを開く"),
        action: () => window.taskAction("", "open-default-dir", false),
      });
      items.push({
        icon: "ic-settings",
        label: L("환경설정", "Settings", "環境設定"),
        action: () => document.getElementById("btn-settings-modal").click(),
      });
    }

    menu.innerHTML = items
      .map((item, idx) => {
        if (item.divider) return `<div class="ctx-menu-divider"></div>`;
        return `<button type="button" class="ctx-menu-item ${item.danger ? "ctx-danger" : ""}" data-cidx="${idx}">
          ${iconSvg(item.icon)}
          <span>${item.label}</span>
        </button>`;
      })
      .join("");

    menu.querySelectorAll(".ctx-menu-item").forEach((btn) => {
      btn.addEventListener("click", () => {
        const idx = parseInt(btn.getAttribute("data-cidx"), 10);
        hideMenu();
        if (items[idx] && typeof items[idx].action === "function") {
          items[idx].action();
        }
      });
    });

    menu.classList.remove("hidden");
    const rect = menu.getBoundingClientRect();
    const x = Math.min(e.clientX, window.innerWidth - rect.width - 10);
    const y = Math.min(e.clientY, window.innerHeight - rect.height - 10);
    menu.style.left = `${Math.max(8, x)}px`;
    menu.style.top = `${Math.max(8, y)}px`;
  });
}

function renderSniffedLinks() {
  const container = document.getElementById("sniffer-results");
  const filtered = state.sniffedLinks.filter(
    (l) => state.sniffFilter === "all" || l.category === state.sniffFilter
  );
  const selectedCount = state.sniffedLinks.filter((l) => l.checked).length;
  setTextIfChanged(
    document.getElementById("sniffer-selected-count"),
    `${selectedCount} / ${state.sniffedLinks.length}`
  );

  if (filtered.length === 0) {
    container.innerHTML = `<div class="sniffer-placeholder">No matching files.</div>`;
    return;
  }

  container.innerHTML = filtered
    .map(
      (l, idx) => `
    <label class="sniffer-row">
      <input type="checkbox" class="glass-checkbox" data-sidx="${idx}" ${l.checked ? "checked" : ""} />
      ${iconSvg(categoryIconId(l.category))}
      <strong style="flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">${l.fileName}</strong>
      <span class="task-pill pill-running tabular">.${l.extension}</span>
    </label>
  `
    )
    .join("");

  container.querySelectorAll("input[type=checkbox]").forEach((cb) => {
    cb.addEventListener("change", () => {
      const item = filtered[parseInt(cb.dataset.sidx, 10)];
      if (item) item.checked = cb.checked;
      const sel = state.sniffedLinks.filter((l) => l.checked).length;
      setTextIfChanged(
        document.getElementById("sniffer-selected-count"),
        `${sel} / ${state.sniffedLinks.length}`
      );
    });
  });
}

function connectTelemetryStream() {
  const es = new EventSource("/api/stream");
  es.onmessage = (evt) => {
    try {
      const snap = JSON.parse(evt.data);
      renderGlobalState(snap);
    } catch (_) {}
  };
  es.onerror = () => {
    es.close();
    setTimeout(connectTelemetryStream, 1500);
  };
}

document.addEventListener("DOMContentLoaded", () => {
  applyThemeToDOM(state.theme);
  applyLanguageToDOM(state.language);
  initEvents();
  connectTelemetryStream();
});
