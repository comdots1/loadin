// Loadin Universal Browser Connector (Chrome, Edge, Whale, Brave, Opera, Vivaldi, Firefox)
const ext = typeof browser !== "undefined" ? browser : chrome;
const LOADIN_ENDPOINT = "http://127.0.0.1:17890";

const INTERCEPT_EXTENSIONS = [
  ".zip", ".rar", ".7z", ".tar", ".gz", ".xz", ".tgz", ".bz2", ".zst", ".cab",
  ".exe", ".msi", ".msix", ".dmg", ".pkg", ".iso", ".img", ".vhd", ".vhdx", ".vmdk",
  ".apk", ".aab", ".ipa", ".deb", ".rpm", ".appimage", ".bin",
  ".mp4", ".mkv", ".webm", ".avi", ".mov", ".wmv", ".flv", ".ts", ".m4v",
  ".mp3", ".flac", ".wav", ".aac", ".ogg", ".m4a",
  ".pdf", ".epub", ".mobi", ".cbz", ".cbr", ".pptx", ".xlsx", ".docx", ".hwp", ".hwpx",
  ".psd", ".ai", ".blend", ".unitypackage", ".safetensors", ".gguf", ".onnx", ".pt", ".ckpt"
];

const INTERCEPT_MIMES = [
  "application/octet-stream",
  "application/x-msdownload",
  "application/x-msdos-program",
  "application/x-msi",
  "application/zip",
  "application/x-zip-compressed",
  "application/x-rar-compressed",
  "application/x-7z-compressed",
  "application/x-tar",
  "application/gzip",
  "application/x-bzip2",
  "application/x-iso9660-image",
  "application/vnd.android.package-archive"
];

// Deduplication cache to prevent double-popup when both content script & downloads API fire
const recentInterceptMap = new Map();

function isRecentlyIntercepted(url) {
  if (!url) return false;
  const now = Date.now();
  for (const [k, ts] of recentInterceptMap.entries()) {
    if (now - ts > 3000) recentInterceptMap.delete(k);
  }
  const prev = recentInterceptMap.get(url);
  if (prev && now - prev < 3000) {
    return true;
  }
  recentInterceptMap.set(url, now);
  return false;
}

function getStorageConfig(keys) {
  return new Promise((resolve) => {
    try {
      const maybePromise = ext.storage.local.get(keys, (res) => {
        if (res) resolve(res);
      });
      if (maybePromise && typeof maybePromise.then === "function") {
        maybePromise.then(resolve).catch(() => resolve({}));
      }
    } catch (_) {
      resolve({});
    }
  });
}

function getCookiesHeader(targetUrl) {
  return new Promise((resolve) => {
    try {
      if (!ext.cookies || !ext.cookies.getAll) {
        resolve("");
        return;
      }
      const maybePromise = ext.cookies.getAll({ url: targetUrl }, (cookies) => {
        if (Array.isArray(cookies)) {
          resolve(cookies.map((c) => `${c.name}=${c.value}`).join("; "));
        }
      });
      if (maybePromise && typeof maybePromise.then === "function") {
        maybePromise
          .then((cookies) => {
            if (Array.isArray(cookies)) {
              resolve(cookies.map((c) => `${c.name}=${c.value}`).join("; "));
            } else {
              resolve("");
            }
          })
          .catch(() => resolve(""));
      }
    } catch (_) {
      resolve("");
    }
  });
}

function setupContextMenus() {
  try {
    ext.contextMenus.removeAll(() => {
      ext.contextMenus.create({
        id: "loadin-download-link",
        title: "Loadin으로 고속 다운로드",
        contexts: ["link", "video", "audio", "image"]
      });
      ext.contextMenus.create({
        id: "loadin-sniff-page",
        title: "이 페이지의 모든 파일 Loadin으로 스캔",
        contexts: ["page"]
      });
    });
  } catch (_) {}
}

ext.runtime.onInstalled.addListener(() => {
  try {
    ext.storage.local.set({
      autoIntercept: true,
      minSizeMB: 0
    });
  } catch (_) {}
  setupContextMenus();
});

if (ext.runtime.onStartup) {
  ext.runtime.onStartup.addListener(() => {
    setupContextMenus();
  });
}

ext.contextMenus.onClicked.addListener(async (info, tab) => {
  if (info.menuItemId === "loadin-download-link") {
    const targetUrl = info.linkUrl || info.srcUrl;
    if (!targetUrl) return;
    const cookie = await getCookiesHeader(targetUrl);
    await sendToLoadin({
      url: targetUrl,
      referer: info.pageUrl || (tab && tab.url) || "",
      cookie,
      userAgent: navigator.userAgent || "",
      source: "우클릭 컨텍스트 메뉴"
    });
  } else if (info.menuItemId === "loadin-sniff-page") {
    const pageUrl = info.pageUrl || (tab && tab.url);
    if (!pageUrl) return;
    try {
      await fetch(`${LOADIN_ENDPOINT}/api/browser/intercept`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          url: pageUrl,
          referer: pageUrl,
          source: "페이지 링크 스캔"
        })
      });
    } catch (_) {
      try {
        ext.tabs.create({ url: `loadin://download?url=${encodeURIComponent(pageUrl)}` });
      } catch (_) {}
    }
  }
});

// Listen for messages from content.js (bypasses strict webpage CSP & HTTPS mixed-content blocks)
ext.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (!msg || !msg.type) return false;

  if (msg.type === "LOADIN_INTERCEPT_LINK") {
    (async () => {
      try {
        const cfg = await getStorageConfig(["autoIntercept"]);
        if (cfg.autoIntercept === false && !msg.force) {
          sendResponse({ ok: false, intercepted: false });
          return;
        }
        const pingRes = await fetch(`${LOADIN_ENDPOINT}/api/browser/ping`, { method: "GET" });
        if (!pingRes.ok) {
          sendResponse({ ok: false, intercepted: false });
          return;
        }
        const pingData = await pingRes.json();
        if (pingData.browserIntercept === false && !msg.force) {
          sendResponse({ ok: false, intercepted: false });
          return;
        }

        if (isRecentlyIntercepted(msg.url)) {
          sendResponse({ ok: true, intercepted: true });
          return;
        }

        const cookie = await getCookiesHeader(msg.url);
        const ok = await sendToLoadin({
          url: msg.url,
          fileName: msg.fileName || "",
          referer: msg.referer || (sender && sender.tab && sender.tab.url) || "",
          cookie,
          userAgent: msg.userAgent || navigator.userAgent || "",
          source: msg.source || "브라우저 링크 클릭"
        });
        sendResponse({ ok, intercepted: ok });
      } catch (_) {
        sendResponse({ ok: false, intercepted: false });
      }
    })();
    return true; // keep message channel open for async response
  }

  return false;
});

// Intercept native browser downloads across Chrome, Edge, Whale, Brave, Opera, Vivaldi, Firefox
if (ext.downloads && ext.downloads.onCreated) {
  ext.downloads.onCreated.addListener(async (downloadItem) => {
    try {
      const url = downloadItem.finalUrl || downloadItem.url || "";
      if (!url.startsWith("http://") && !url.startsWith("https://")) return;
      if (url.includes("127.0.0.1:17890") || url.includes("localhost:17890")) return;

      const cfg = await getStorageConfig(["autoIntercept", "minSizeMB"]);
      if (cfg.autoIntercept === false) return;

      const minBytes = (cfg.minSizeMB ?? 0) * 1024 * 1024;
      const lowerUrl = url.toLowerCase().split("?")[0].split("#")[0];
      const lowerFn = (downloadItem.filename || "").toLowerCase();
      const lowerMime = (downloadItem.mime || "").toLowerCase();

      const hasMatchExt = INTERCEPT_EXTENSIONS.some(
        (e) => lowerUrl.endsWith(e) || lowerFn.endsWith(e)
      );
      const hasMatchMime =
        INTERCEPT_MIMES.some((m) => lowerMime.includes(m)) ||
        lowerMime.startsWith("video/") ||
        lowerMime.startsWith("audio/");

      if (downloadItem.fileSize > 0 && minBytes > 0 && downloadItem.fileSize < minBytes && !hasMatchExt) {
        return;
      }
      if (downloadItem.fileSize <= 0 && !hasMatchExt && !hasMatchMime) {
        return;
      }

      // Check if Loadin desktop engine is online and browserIntercept is enabled
      const pingRes = await fetch(`${LOADIN_ENDPOINT}/api/browser/ping`, { method: "GET" });
      if (!pingRes.ok) return;
      const pingData = await pingRes.json();
      if (pingData.browserIntercept === false) return;

      // Cancel native browser download and hand over to Loadin popup dialog
      try {
        ext.downloads.cancel(downloadItem.id, () => {
          if (ext.downloads.erase) {
            ext.downloads.erase({ id: downloadItem.id });
          }
        });
      } catch (_) {}

      if (isRecentlyIntercepted(url)) {
        return;
      }

      const cookie = await getCookiesHeader(url);
      await sendToLoadin({
        url,
        fileName: downloadItem.filename ? downloadItem.filename.split(/[\\/]/).pop() : "",
        referer: downloadItem.referrer || "",
        cookie,
        userAgent: navigator.userAgent || "",
        totalSize: downloadItem.fileSize > 0 ? downloadItem.fileSize : 0,
        mime: downloadItem.mime || "",
        source: "브라우저 자동 연동"
      });
    } catch (_) {
      // If Loadin is not running, let the browser handle the download normally
    }
  });
}

async function sendToLoadin(payload) {
  try {
    const res = await fetch(`${LOADIN_ENDPOINT}/api/browser/intercept`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    return res.ok;
  } catch (_) {
    return false;
  }
}
