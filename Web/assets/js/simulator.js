/**
 * Loadin Dynamic Work-Stealing vs Legacy Downloader Interactive Simulator
 */
(function() {
  const btnStart = document.getElementById("simStartBtn");
  const btnReset = document.getElementById("simResetBtn");
  const simSpeedEl = document.getElementById("simLiveSpeed");
  const logTerminal = document.getElementById("simLogTerminal");

  if (!btnStart || !btnReset) return;

  let isRunning = false;
  let timerId = null;
  let progressLegacy = [0, 0, 0, 0];
  let progressLoadin = [0, 0, 0, 0];
  let loadinStolen = false;
  let legacyMerging = false;
  let legacyDone = false;
  let loadinDone = false;

  const legacyFills = [
    document.getElementById("legacy-fill-1"),
    document.getElementById("legacy-fill-2"),
    document.getElementById("legacy-fill-3"),
    document.getElementById("legacy-fill-4")
  ];

  const loadinFills = [
    document.getElementById("loadin-fill-1"),
    document.getElementById("loadin-fill-2"),
    document.getElementById("loadin-fill-3"),
    document.getElementById("loadin-fill-4")
  ];

  const legacyStatusText = document.getElementById("legacy-status-text");
  const loadinStatusText = document.getElementById("loadin-status-text");

  function appendLog(msg, type = "info") {
    if (!logTerminal) return;
    const now = new Date();
    const timeStr = now.toTimeString().split(" ")[0] + "." + Math.floor(now.getMilliseconds() / 100);
    const div = document.createElement("div");
    div.className = "log-entry-" + type;
    div.textContent = `[${timeStr}] ${msg}`;
    logTerminal.appendChild(div);
    logTerminal.scrollTop = logTerminal.scrollHeight;
  }

  function resetSim() {
    isRunning = false;
    clearInterval(timerId);
    timerId = null;

    progressLegacy = [0, 0, 0, 0];
    progressLoadin = [0, 0, 0, 0];
    loadinStolen = false;
    legacyMerging = false;
    legacyDone = false;
    loadinDone = false;

    legacyFills.forEach(el => {
      if (el) {
        el.style.width = "0%";
        el.className = "thread-bar-fill";
      }
    });

    loadinFills.forEach(el => {
      if (el) {
        el.style.width = "0%";
        el.className = "thread-bar-fill";
      }
    });

    if (legacyStatusText) legacyStatusText.textContent = "준비 (대기 중)";
    if (loadinStatusText) loadinStatusText.textContent = "준비 (대기 중)";
    if (simSpeedEl) simSpeedEl.textContent = "0.0 MB/s";

    btnStart.textContent = "시뮬레이션 시작";
    btnStart.disabled = false;
  }

  function runSim() {
    if (isRunning) return;
    resetSim();
    isRunning = true;
    btnStart.disabled = true;

    appendLog("1,024 MB 대용량 파일 다운로드 세그먼트 할당 시작...", "info");
    appendLog("4개 병렬 HTTP Range 연결 생성됨 (각 256 MB)", "info");

    if (legacyStatusText) legacyStatusText.textContent = "다운로드 중...";
    if (loadinStatusText) loadinStatusText.textContent = "고속 가속 스트리밍 중...";

    let tick = 0;

    timerId = setInterval(() => {
      tick++;

      // 1. Legacy Downloader Progress
      if (!legacyDone) {
        // Threads 1, 2, 4 progress normally
        if (progressLegacy[0] < 100) progressLegacy[0] += 2.8;
        if (progressLegacy[1] < 100) progressLegacy[1] += 2.5;
        if (progressLegacy[3] < 100) progressLegacy[3] += 2.6;

        // Thread 3 gets throttled / lagged after 45%
        if (progressLegacy[2] < 45) {
          progressLegacy[2] += 2.0;
        } else if (progressLegacy[2] < 100) {
          progressLegacy[2] += 0.45; // Bottleneck!
          if (legacyFills[2]) legacyFills[2].classList.add("stalled");
        }

        // Cap at 100
        for (let i = 0; i < 4; i++) {
          if (progressLegacy[i] >= 100) {
            progressLegacy[i] = 100;
            if (legacyFills[i]) legacyFills[i].classList.add("complete");
          }
          if (legacyFills[i]) legacyFills[i].style.width = progressLegacy[i] + "%";
        }

        const legacyAll100 = progressLegacy.every(p => p >= 100);
        if (legacyAll100 && !legacyMerging) {
          legacyMerging = true;
          if (legacyStatusText) {
            legacyStatusText.textContent = "⚠️ 파일 조각 디스크 병합 중 (지연 발생)...";
            legacyStatusText.style.color = "var(--accent-orange)";
          }
          appendLog("⚠️ 기존 방식: 임시 파일 조각 합치는 중... 디스크 멈춤 발생", "warn");
          setTimeout(() => {
            legacyDone = true;
            if (legacyStatusText) {
              legacyStatusText.textContent = "완료 (병합 지연 +8.5초 소요)";
              legacyStatusText.style.color = "var(--text-muted)";
            }
            appendLog("기존 다운로더 완료: 총 소요 23.4초", "info");
          }, 4500);
        }
      }

      // 2. Loadin Engine Progress (Work-Stealing)
      if (!loadinDone) {
        if (progressLoadin[0] < 100) progressLoadin[0] += 3.4;
        if (progressLoadin[1] < 100) progressLoadin[1] += 3.2;
        if (progressLoadin[3] < 100) progressLoadin[3] += 3.3;

        // Thread 3 slow link simulation
        if (progressLoadin[2] < 40) {
          progressLoadin[2] += 1.6;
        } else {
          progressLoadin[2] += 0.8;
        }

        // Work-Stealing Trigger!
        if (progressLoadin[0] >= 100 && !loadinStolen) {
          loadinStolen = true;
          appendLog("⚡ [Work-Stealing] Thread-1 완료! Thread-3의 잔여 세그먼트 원자적 분할 (Atomic Bisect) 감지!", "stealing");
          appendLog("⚡ Thread-1 -> Range: bytes=188M-256M 동적 재할당 완료", "stealing");

          if (loadinFills[0]) {
            loadinFills[0].classList.add("stolen");
            loadinFills[0].style.width = "40%";
            progressLoadin[0] = 40; // Starts aiding Thread 3
          }
        }

        if (loadinStolen) {
          if (progressLoadin[0] < 100) progressLoadin[0] += 4.2;
          if (progressLoadin[2] < 100) progressLoadin[2] += 3.6;
        }

        for (let i = 0; i < 4; i++) {
          if (progressLoadin[i] >= 100) {
            progressLoadin[i] = 100;
            if (loadinFills[i]) loadinFills[i].className = "thread-bar-fill complete";
          }
          if (loadinFills[i]) loadinFills[i].style.width = progressLoadin[i] + "%";
        }

        const loadinAll100 = progressLoadin.every(p => p >= 100);
        if (loadinAll100 && !loadinDone) {
          loadinDone = true;
          if (loadinStatusText) {
            loadinStatusText.textContent = "✓ 완료! Positional I/O 병합 지연 0.00초";
            loadinStatusText.style.color = "var(--accent-green)";
          }
          appendLog("✓ [Loadin Engine] os.File.WriteAt 락프리 즉시 저장 완료! 병합 지연 0.00초", "success");
        }
      }

      // Update Live Speed
      if (!loadinDone) {
        const speed = (78.4 + Math.random() * 8.6).toFixed(1);
        if (simSpeedEl) simSpeedEl.textContent = speed + " MB/s";
      } else {
        if (simSpeedEl) simSpeedEl.textContent = "0.0 MB/s";
      }

      if (loadinDone && legacyDone) {
        clearInterval(timerId);
        timerId = null;
        isRunning = false;
        btnStart.textContent = "다시 시뮬레이션";
        btnStart.disabled = false;
      }
    }, 120);
  }

  btnStart.addEventListener("click", runSim);
  btnReset.addEventListener("click", resetSim);
})();
