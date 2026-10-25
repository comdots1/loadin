const ext = typeof browser !== "undefined" ? browser : chrome;
const ENDPOINT = "http://127.0.0.1:17890";

function getStorage(keys) {
  return new Promise((resolve) => {
    try {
      const p = ext.storage.local.get(keys, (res) => {
        if (res) resolve(res);
      });
      if (p && typeof p.then === "function") {
        p.then(resolve).catch(() => resolve({}));
      }
    } catch (_) {
      resolve({});
    }
  });
}

document.addEventListener("DOMContentLoaded", async () => {
  const statusEl = document.getElementById("conn-status");
  const toggleEl = document.getElementById("toggle-intercept");
  const minSizeEl = document.getElementById("min-size");
  const directUrlEl = document.getElementById("direct-url");

  const cfg = await getStorage(["autoIntercept", "minSizeMB"]);
  toggleEl.checked = cfg.autoIntercept !== false;
  minSizeEl.value = String(cfg.minSizeMB ?? 0);

  toggleEl.addEventListener("change", () => {
    ext.storage.local.set({ autoIntercept: toggleEl.checked });
  });
  minSizeEl.addEventListener("change", () => {
    ext.storage.local.set({ minSizeMB: parseInt(minSizeEl.value, 10) || 0 });
  });

  try {
    const res = await fetch(`${ENDPOINT}/api/browser/ping`);
    if (res.ok) {
      statusEl.textContent = "● 엔진 연결됨";
      statusEl.classList.remove("offline");
    } else {
      throw new Error();
    }
  } catch (_) {
    statusEl.textContent = "○ 오프라인";
    statusEl.classList.add("offline");
  }

  document.getElementById("btn-send").addEventListener("click", async () => {
    const url = directUrlEl.value.trim();
    if (!url) return;
    await fetch(`${ENDPOINT}/api/browser/intercept`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url, source: "확장프로그램 팝업" })
    });
    directUrlEl.value = "";
    window.close();
  });

  document.getElementById("btn-open-ui").addEventListener("click", () => {
    ext.tabs.create({ url: ENDPOINT });
  });
});
