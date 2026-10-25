/**
 * Loadin Hardware Spec Advisor
 * Detects visitor's client hardware environment (CPU Cores, RAM, Platform)
 * and calculates the exact recommended threads and buffer size used by Loadin's engine.
 */
(function() {
  const cpuChip = document.getElementById("detectedCpuCores");
  const ramChip = document.getElementById("detectedRamSize");
  const platformChip = document.getElementById("detectedPlatform");
  const recommendedThreads = document.getElementById("recommendedThreads");
  const recommendedBuffer = document.getElementById("recommendedBuffer");
  const customCoreSlider = document.getElementById("customCoreSlider");
  const sliderCoreVal = document.getElementById("sliderCoreVal");

  // Detection
  const cores = navigator.hardwareConcurrency || 8;
  const memory = (navigator.deviceMemory || 16);
  const platform = navigator.userAgent.indexOf("Win") !== -1 ? "Windows 11 / 10 x64" : (navigator.userAgent.indexOf("Mac") !== -1 ? "macOS" : "Linux / Other");

  function calculateOptimal(c) {
    let threads = 8;
    let buffer = "4 MB";
    let speedup = "3.8x ~ 5.2x";

    if (c <= 2) {
      threads = 4;
      buffer = "2 MB";
      speedup = "2.5x ~ 3.2x";
    } else if (c <= 4) {
      threads = 8;
      buffer = "4 MB";
      speedup = "3.8x ~ 4.6x";
    } else if (c <= 8) {
      threads = 12;
      buffer = "8 MB";
      speedup = "5.0x ~ 6.5x";
    } else if (c <= 16) {
      threads = 16;
      buffer = "12 MB";
      speedup = "7.2x ~ 9.0x";
    } else {
      threads = 24;
      buffer = "16 MB";
      speedup = "9.5x ~ 12.0x";
    }

    return { threads, buffer, speedup };
  }

  // Populate UI
  if (cpuChip) cpuChip.textContent = `${cores} Cores`;
  if (ramChip) ramChip.textContent = `${memory} GB+`;
  if (platformChip) platformChip.textContent = platform;

  const currentOpt = calculateOptimal(cores);
  if (recommendedThreads) recommendedThreads.textContent = `${currentOpt.threads} 스레드`;
  if (recommendedBuffer) recommendedBuffer.textContent = `버퍼 ${currentOpt.buffer} · 예상 가속 ${currentOpt.speedup}`;

  if (customCoreSlider) {
    customCoreSlider.value = cores;
    if (sliderCoreVal) sliderCoreVal.textContent = `${cores} Cores`;

    customCoreSlider.addEventListener("input", function(e) {
      const val = parseInt(e.target.value, 10);
      if (sliderCoreVal) sliderCoreVal.textContent = `${val} Cores`;
      const opt = calculateOptimal(val);
      if (recommendedThreads) recommendedThreads.textContent = `${opt.threads} 스레드`;
      if (recommendedBuffer) recommendedBuffer.textContent = `버퍼 ${opt.buffer} · 예상 가속 ${opt.speedup}`;
    });
  }
})();
