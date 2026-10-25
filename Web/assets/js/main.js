/**
 * Loadin Official Web — Main Interactive Controller
 * Theme Switcher, i18n Language Toggle, OS Tabs, Copy Clipboard, Mobile Nav, FAQ Accordion
 */
(function () {
  // 1. Theme Management (Matches Loadin Desktop App)
  const themeBtn = document.getElementById("themeToggleBtn");
  const htmlRoot = document.documentElement;

  function applyTheme(theme) {
    htmlRoot.setAttribute("data-theme", theme);
    try {
      localStorage.setItem("loadin_theme", theme);
    } catch (_) {}

    if (themeBtn) {
      themeBtn.innerHTML = theme === "white" 
        ? `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>`
        : `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>`;
    }
  }

  const savedTheme = localStorage.getItem("loadin_theme") || "black";
  applyTheme(savedTheme);

  if (themeBtn) {
    themeBtn.addEventListener("click", () => {
      const current = htmlRoot.getAttribute("data-theme") || "black";
      applyTheme(current === "black" ? "white" : "black");
    });
  }

  // 2. Language Management (KO / EN)
  const langBtn = document.getElementById("langToggleBtn");
  let currentLang = localStorage.getItem("loadin_lang") || "ko";

  const i18nData = {
    ko: {
      navFeatures: "주요 기능",
      navSimulator: "가속 시뮬레이터",
      navAdvisor: "사양 최적화",
      navBenchmark: "성능 비교",
      navExtension: "브라우저 연동",
      navDownload: "다운로드",
      navFaq: "FAQ",
      navCtaBtn: "다운로드",
      heroTitlePrefix: "다운로드의 패러다임을 바꾸다.",
      heroTitleHighlight: "병합 지연 0초",
      heroTitleSuffix: "의 초고속 가속 엔진.",
      heroDesc: "동적 세그먼트 워크 스틸링과 위치 직접 I/O(WriteAt)로 99% 멈춤 없이 100% 대역폭을 유지합니다. Windows와 macOS를 모두 네이티브로 지원합니다.",
      btnInstall: "Windows 설치 프로그램 (6.8 MB)",
      btnMacApp: "macOS 앱 번들 다운로드",
      btnPortable: "무설치 포터블 (Loadin.exe)",
      heroMeta1: "Windows 10/11 & macOS",
      heroMeta2: "초경량 네이티브 바이너리",
      heroMeta3: "MIT 라이선스",
      simTitle: "실시간 동적 워크 스틸링(Work-Stealing) 가속 시뮬레이터",
      simDesc: "스레드가 멈추면 남은 조각을 원자적으로 분할(Bisect)하여 즉시 가속하는 Loadin의 알고리즘을 직접 확인해 보세요.",
      advisorTitle: "내 PC 사양 기반 최적 스레드 자동 동기화",
      advisorDesc: "CPU 논리 코어와 물리 RAM을 분석하여 네트워크 카드를 극한으로 활용할 수 있는 스레드와 버퍼를 자동 추천합니다.",
      dlHubTitle: "릴리즈 다운로드",
      dlHubDesc: "사용 중인 운영체제(Windows / macOS)와 환경에 맞춰 적합한 패키지를 선택하세요."
    },
    en: {
      navFeatures: "Features",
      navSimulator: "Simulator",
      navAdvisor: "Spec Advisor",
      navBenchmark: "Benchmark",
      navExtension: "Browser Extension",
      navDownload: "Download",
      navFaq: "FAQ",
      navCtaBtn: "Download",
      heroTitlePrefix: "Redefining High Speed.",
      heroTitleHighlight: "Zero-Merge Delay",
      heroTitleSuffix: "Native Download Accelerator.",
      heroDesc: "Dynamic Work-Stealing Segmentation and Positional Direct I/O ensure instant disk writes without the 99% freeze. Fully native on Windows and macOS.",
      btnInstall: "Windows Installer (6.8 MB)",
      btnMacApp: "Download macOS App Bundle",
      btnPortable: "Portable Binary (Loadin.exe)",
      heroMeta1: "Windows 10/11 & macOS",
      heroMeta2: "Ultra-lightweight Native",
      heroMeta3: "MIT Open Source",
      simTitle: "Live Dynamic Work-Stealing Acceleration Simulator",
      simDesc: "Watch how fast workers atomically bisect lagging segments in real-time, eliminating bottlenecks entirely.",
      advisorTitle: "Hardware Spec Advisor & Auto-Tuning",
      advisorDesc: "Accurately analyzes your CPU logical cores & physical memory to scale optimal HTTP range streams without context-switch overhead.",
      dlHubTitle: "Release Downloads",
      dlHubDesc: "Choose the package tailored for your operating system (Windows or macOS)."
    }
  };

  function applyLang(lang) {
    currentLang = lang;
    try {
      localStorage.setItem("loadin_lang", lang);
    } catch (_) {}

    if (langBtn) {
      langBtn.textContent = lang === "ko" ? "EN" : "한국어";
    }

    const dict = i18nData[lang];
    if (!dict) return;

    document.querySelectorAll("[data-i18n]").forEach(el => {
      const key = el.getAttribute("data-i18n");
      if (dict[key]) {
        el.textContent = dict[key];
      }
    });
  }

  applyLang(currentLang);

  if (langBtn) {
    langBtn.addEventListener("click", () => {
      applyLang(currentLang === "ko" ? "en" : "ko");
    });
  }

  // 3. Header Scroll Effect
  const header = document.querySelector(".site-header");
  window.addEventListener("scroll", () => {
    if (window.scrollY > 20) {
      header?.classList.add("scrolled");
    } else {
      header?.classList.remove("scrolled");
    }
  }, { passive: true });

  // 4. One-Click Copy Buttons
  document.querySelectorAll("[data-copy]").forEach(btn => {
    btn.addEventListener("click", () => {
      const textToCopy = btn.getAttribute("data-copy");
      if (!textToCopy) return;

      navigator.clipboard.writeText(textToCopy).then(() => {
        const originalText = btn.innerHTML;
        btn.innerHTML = `<span style="color:var(--accent-green);font-weight:700;">✓ 복사됨!</span>`;
        setTimeout(() => {
          btn.innerHTML = originalText;
        }, 2000);
      });
    });
  });

  // 5. OS Download Tabs (Windows vs macOS vs All)
  const isMac = navigator.userAgent.indexOf("Mac") !== -1 || navigator.platform.indexOf("Mac") !== -1;
  const osTabs = document.querySelectorAll(".os-tab-btn");
  const winPanel = document.getElementById("dl-panel-windows");
  const macPanel = document.getElementById("dl-panel-mac");

  function switchOsTab(targetOs) {
    osTabs.forEach(t => {
      const match = t.getAttribute("data-os") === targetOs;
      if (match) {
        t.classList.add("active");
        t.setAttribute("aria-selected", "true");
      } else {
        t.classList.remove("active");
        t.setAttribute("aria-selected", "false");
      }
    });

    if (targetOs === "mac") {
      if (winPanel) winPanel.style.display = "none";
      if (macPanel) macPanel.style.display = "grid";
    } else if (targetOs === "all") {
      if (winPanel) winPanel.style.display = "grid";
      if (macPanel) macPanel.style.display = "grid";
    } else {
      if (winPanel) winPanel.style.display = "grid";
      if (macPanel) macPanel.style.display = "none";
    }
  }

  osTabs.forEach(tab => {
    tab.addEventListener("click", () => {
      const os = tab.getAttribute("data-os");
      switchOsTab(os);
    });
  });

  // Auto-select macOS tab if user is on Mac
  if (isMac) {
    switchOsTab("mac");
    // Also adjust hero primary CTA to Mac App
    const heroPrimary = document.getElementById("heroPrimaryCta");
    if (heroPrimary) {
      heroPrimary.href = "/downloads/Loadin-macOS.zip";
      heroPrimary.setAttribute("download", "Loadin-macOS.zip");
      const textEl = heroPrimary.querySelector(".cta-text");
      if (textEl) textEl.textContent = "macOS 앱 번들 다운로드 (Universal)";
    }
  }

  // 6. Browser Extension Tab Switcher
  const extTabs = document.querySelectorAll(".ext-tab-btn");
  const extStepEdge = document.getElementById("ext-steps-edge");
  const extStepChrome = document.getElementById("ext-steps-chrome");

  extTabs.forEach(tab => {
    tab.addEventListener("click", () => {
      extTabs.forEach(t => t.classList.remove("active"));
      tab.classList.add("active");
      const target = tab.getAttribute("data-tab");

      if (target === "edge") {
        if (extStepEdge) extStepEdge.style.display = "flex";
        if (extStepChrome) extStepChrome.style.display = "none";
      } else {
        if (extStepEdge) extStepEdge.style.display = "none";
        if (extStepChrome) extStepChrome.style.display = "flex";
      }
    });
  });

  // 7. FAQ Accordion
  document.querySelectorAll(".faq-item").forEach(item => {
    item.addEventListener("click", () => {
      const ans = item.querySelector(".faq-answer");
      const icon = item.querySelector(".faq-icon");
      if (!ans) return;

      const isExpanded = ans.style.display === "block";
      ans.style.display = isExpanded ? "none" : "block";
      if (icon) {
        icon.style.transform = isExpanded ? "rotate(0deg)" : "rotate(180deg)";
      }
    });
  });

  // 8. Mobile Navigation Toggle
  const mobileToggle = document.getElementById("mobileMenuBtn");
  const navLinks = document.querySelector(".nav-links");

  if (mobileToggle && navLinks) {
    mobileToggle.addEventListener("click", () => {
      const isOpen = navLinks.classList.toggle("mobile-open");
      mobileToggle.setAttribute("aria-expanded", isOpen);
      if (isOpen) {
        navLinks.style.display = "flex";
        navLinks.style.flexDirection = "column";
        navLinks.style.position = "absolute";
        navLinks.style.top = "72px";
        navLinks.style.left = "0";
        navLinks.style.right = "0";
        navLinks.style.background = "var(--bg-surface)";
        navLinks.style.padding = "24px";
        navLinks.style.borderBottom = "1px solid var(--border-medium)";
        navLinks.style.boxShadow = "var(--shadow-lg)";
      } else {
        navLinks.style.display = "";
      }
    });
  }
})();
