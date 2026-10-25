// Loadin Content Script — Universal Multi-Browser Link & Media Interceptor
(function () {
  if (window.__loadinContentScriptInjected) return;
  window.__loadinContentScriptInjected = true;

  const ext = typeof browser !== "undefined" ? browser : typeof chrome !== "undefined" ? chrome : null;
  const ENDPOINT = "http://127.0.0.1:17890";
  const DL_EXTS = [
    ".zip", ".rar", ".7z", ".tar", ".gz", ".xz", ".tgz", ".bz2", ".zst", ".cab",
    ".exe", ".msi", ".msix", ".dmg", ".pkg", ".iso", ".img", ".vhd", ".vhdx", ".vmdk",
    ".apk", ".aab", ".ipa", ".deb", ".rpm", ".appimage", ".bin",
    ".mp4", ".mkv", ".webm", ".avi", ".mov", ".wmv", ".flv", ".ts", ".m4v",
    ".mp3", ".flac", ".wav", ".aac", ".ogg", ".m4a",
    ".pdf", ".epub", ".mobi", ".cbz", ".cbr", ".pptx", ".xlsx", ".docx", ".hwp", ".hwpx",
    ".psd", ".ai", ".blend", ".unitypackage", ".safetensors", ".gguf", ".onnx", ".pt", ".ckpt"
  ];

  let hoverBtn = null;
  let currentTargetUrl = "";
  let hideTimer = null;

  function sendInterceptViaBackground(payload) {
    return new Promise((resolve) => {
      if (ext && ext.runtime && ext.runtime.sendMessage) {
        try {
          ext.runtime.sendMessage(
            {
              type: "LOADIN_INTERCEPT_LINK",
              ...payload
            },
            (resp) => {
              if (resp && resp.intercepted) {
                resolve(true);
              } else {
                resolve(false);
              }
            }
          );
          return;
        } catch (_) {}
      }
      resolve(false);
    });
  }

  function ensureHoverButton() {
    if (hoverBtn) return hoverBtn;
    hoverBtn = document.createElement("div");
    hoverBtn.id = "loadin-floating-hover-btn";
    hoverBtn.style.cssText = [
      "position:fixed",
      "z-index:2147483647",
      "display:none",
      "align-items:center",
      "gap:7px",
      "padding:6px 12px",
      "border-radius:9px",
      "background:rgba(20,20,22,0.94)",
      "backdrop-filter:blur(16px)",
      "-webkit-backdrop-filter:blur(16px)",
      "border:1px solid rgba(10,132,255,0.5)",
      "color:#F5F5F7",
      "font-family:'Pretendard',-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif",
      "font-size:12px",
      "font-weight:600",
      "letter-spacing:-0.02em",
      "cursor:pointer",
      "user-select:none",
      "transition:opacity 0.14s ease, border-color 0.14s ease"
    ].join(";");

    hoverBtn.innerHTML = `
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#0A84FF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M12 4v12m0 0l-4-4m4 4l4-4M4 20h16"/>
      </svg>
      <span id="loadin-hover-label">Loadin으로 다운로드</span>
    `;

    hoverBtn.addEventListener("mouseenter", () => {
      if (hideTimer) clearTimeout(hideTimer);
      hoverBtn.style.borderColor = "#409CFF";
    });

    hoverBtn.addEventListener("mouseleave", () => {
      scheduleHide(300);
    });

    hoverBtn.addEventListener("click", async (e) => {
      e.preventDefault();
      e.stopPropagation();
      if (!currentTargetUrl) return;
      const label = hoverBtn.querySelector("#loadin-hover-label");
      label.textContent = "Loadin 연결 중...";

      const handledByBg = await sendInterceptViaBackground({
        url: currentTargetUrl,
        referer: location.href,
        userAgent: navigator.userAgent,
        source: "미디어 호버 버튼",
        force: true
      });

      if (handledByBg) {
        label.textContent = "다운로드 창 열림";
        setTimeout(() => hideHoverButton(), 900);
        return;
      }

      try {
        const res = await fetch(`${ENDPOINT}/api/browser/intercept`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            url: currentTargetUrl,
            referer: location.href,
            userAgent: navigator.userAgent,
            source: "미디어 호버 버튼"
          })
        });
        if (res.ok) {
          label.textContent = "다운로드 창 열림";
          setTimeout(() => hideHoverButton(), 900);
        } else {
          label.textContent = "엔진 연결 필요";
        }
      } catch (_) {
        location.href = `loadin://download?url=${encodeURIComponent(currentTargetUrl)}`;
        hideHoverButton();
      }
    });

    document.documentElement.appendChild(hoverBtn);
    return hoverBtn;
  }

  function isDownloadableLink(anchor) {
    if (!anchor || !anchor.href) return false;
    const raw = anchor.href;
    if (!raw.startsWith("http://") && !raw.startsWith("https://")) return false;
    if (raw.includes("127.0.0.1:17890") || raw.includes("localhost:17890")) return false;
    if (anchor.hasAttribute("download")) return true;
    const clean = raw.toLowerCase().split("?")[0].split("#")[0];
    return DL_EXTS.some((extName) => clean.endsWith(extName));
  }

  function showHoverNear(rect, url, customText) {
    if (hideTimer) clearTimeout(hideTimer);
    currentTargetUrl = url;
    const btn = ensureHoverButton();
    const label = btn.querySelector("#loadin-hover-label");
    label.textContent = customText || "Loadin으로 다운로드";
    btn.style.display = "inline-flex";

    const top = Math.max(8, rect.top - 36);
    const left = Math.min(window.innerWidth - 190, Math.max(8, rect.left));
    btn.style.top = `${top}px`;
    btn.style.left = `${left}px`;
  }

  function hideHoverButton() {
    if (hoverBtn) hoverBtn.style.display = "none";
  }

  function scheduleHide(ms) {
    if (hideTimer) clearTimeout(hideTimer);
    hideTimer = setTimeout(hideHoverButton, ms);
  }

  // Trigger native on-demand Loadin popup dialog when a download link is clicked
  document.addEventListener(
    "click",
    (e) => {
      if (e.defaultPrevented || e.button !== 0 || e.ctrlKey || e.metaKey || e.shiftKey || e.altKey) return;
      const target = e.target && typeof e.target.closest === "function" ? e.target : null;
      if (!target) return;
      const link = target.closest("a[href]");
      if (!isDownloadableLink(link)) return;

      const dlAttr = (link.getAttribute("download") || "").trim();

      // Send asynchronously via background worker; downloads.onCreated in background.js also acts as safety net
      sendInterceptViaBackground({
        url: link.href,
        fileName: dlAttr,
        referer: location.href,
        userAgent: navigator.userAgent,
        source: "다운로드 링크 클릭"
      }).then((intercepted) => {
        if (intercepted) {
          // Handled by Loadin popup
        }
      });
    },
    true
  );

  document.addEventListener(
    "mouseover",
    (e) => {
      const target = e.target && typeof e.target.closest === "function" ? e.target : null;
      if (!target) return;
      const media = target.closest("video, audio");
      if (media) {
        const src =
          media.currentSrc ||
          media.src ||
          (media.querySelector("source") && media.querySelector("source").src);
        if (src && (src.startsWith("http://") || src.startsWith("https://"))) {
          showHoverNear(media.getBoundingClientRect(), src, "Loadin 미디어 다운로드");
        }
      }
    },
    true
  );

  document.addEventListener(
    "mouseout",
    (e) => {
      const target = e.target && typeof e.target.closest === "function" ? e.target : null;
      if (target && target.closest("video, audio")) {
        scheduleHide(600);
      }
    },
    true
  );
})();
