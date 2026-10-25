using System;
using System.Collections.Generic;
using System.Drawing;
using System.Drawing.Drawing2D;
using System.Globalization;
using System.IO;
using System.Net;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading;
using System.Windows.Forms;
using System.Runtime.InteropServices;

namespace LoadinPopup
{
    static class Program
    {
        [STAThread]
        static void Main(string[] args)
        {
            string port = "17890";
            string targetUrl = "";
            string referer = "";
            string cookie = "";
            string userAgent = "";
            string initialFileName = "";

            for (int i = 0; i < args.Length; i++)
            {
                if (args[i] == "--port" && i + 1 < args.Length) port = args[++i];
                else if (args[i] == "--url" && i + 1 < args.Length) targetUrl = args[++i];
                else if (args[i] == "--referer" && i + 1 < args.Length) referer = args[++i];
                else if (args[i] == "--cookie" && i + 1 < args.Length) cookie = args[++i];
                else if (args[i] == "--user-agent" && i + 1 < args.Length) userAgent = args[++i];
                else if (args[i] == "--filename" && i + 1 < args.Length) initialFileName = args[++i];
                else if (i == 0 && !args[i].StartsWith("-")) port = args[i];
            }

            if (string.IsNullOrWhiteSpace(targetUrl))
            {
                try
                {
                    if (Clipboard.ContainsText())
                    {
                        string clip = Clipboard.GetText().Trim();
                        if (clip.StartsWith("http://") || clip.StartsWith("https://"))
                            targetUrl = clip;
                    }
                }
                catch { }
            }

            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new AddDownloadDialog(port, targetUrl, referer, cookie, userAgent, initialFileName));
        }
    }

    public class AddDownloadDialog : Form
    {
        [DllImport("Gdi32.dll", EntryPoint = "CreateRoundRectRgn")]
        private static extern IntPtr CreateRoundRectRgn(int nLeftRect, int nTopRect, int nRightRect, int nBottomRect, int nWidthEllipse, int nHeightEllipse);

        [DllImport("dwmapi.dll")]
        private static extern int DwmSetWindowAttribute(IntPtr hwnd, int attr, ref int attrValue, int attrSize);

        [DllImport("user32.dll")]
        private static extern bool ReleaseCapture();

        [DllImport("user32.dll")]
        private static extern IntPtr SendMessage(IntPtr hWnd, int Msg, int wParam, int lParam);

        private const int WM_NCLBUTTONDOWN = 0xA1;
        private const int HT_CAPTION = 0x2;

        private readonly string baseUrl;
        private readonly string refererUrl;
        private readonly string cookieHeader;
        private readonly string userAgentHeader;
        private bool isExpandedSize = false;
        private bool isWhiteTheme = false;
        private string lang = "ko";

        private TextBox txtUrl;
        private TextBox txtFileName;
        private TextBox txtSaveDir;

        private string probeSizeText = "...";
        private bool acceptRanges = true;
        private string diskFreeText = "";
        private int selectedThreads = 16;
        private int recommendedThreads = 16;
        private readonly int[] threadOptions = new int[] { 4, 8, 12, 16, 24, 32 };

        // Live Download Progress Mode State
        private bool isDownloadingMode = false;
        private string activeTaskId = "";
        private string dlFileName = "";
        private string dlSaveDir = "";
        private string dlStatus = "running";
        private long dlDownloaded = 0;
        private long dlTotalSize = 0;
        private double dlProgress = 0.0;
        private long dlSpeedBps = 0;
        private long dlEtaSeconds = -1;
        private int dlActiveWorkers = 0;
        private int dlMaxConns = 16;
        private int dlDynamicSplits = 0;
        private string dlErrorMsg = "";
        private float[] dlSegmentPcts = new float[0];
        private System.Windows.Forms.Timer pollTimer;
        private int isPollingInFlight = 0;

        // Unified Caption Control Rectangles (- □ ×)
        private Rectangle rectMinBtn;
        private Rectangle rectMaxBtn;
        private Rectangle rectCloseBtn;
        private int hoverCaptionBtn = 0;

        // Action & Segmented Control Rectangles
        private Rectangle rectBrowseBtn;
        private Rectangle rectSegmentTrack;
        private Rectangle[] rectThreadPills;
        private Rectangle rectCancelBtn;
        private Rectangle rectQueueBtn;
        private Rectangle rectStartBtn;
        private int hoverActionBtn = 0;
        private int hoverThreadIdx = -1;

        private readonly Font fontTitle = new Font("Malgun Gothic", 10.2f, FontStyle.Bold);
        private readonly Font fontLabel = new Font("Malgun Gothic", 9.0f, FontStyle.Bold);
        private readonly Font fontBody = new Font("Malgun Gothic", 9.6f, FontStyle.Regular);
        private readonly Font fontSmall = new Font("Malgun Gothic", 8.8f, FontStyle.Regular);
        private readonly Font fontSmallBold = new Font("Malgun Gothic", 8.8f, FontStyle.Bold);
        private readonly Font fontPctLarge = new Font("Malgun Gothic", 13.0f, FontStyle.Bold);
        private readonly Font fontMetricVal = new Font("Malgun Gothic", 10.5f, FontStyle.Bold);

        private Color BgSurface { get { return isWhiteTheme ? Color.FromArgb(245, 245, 247) : Color.FromArgb(20, 20, 22); } }
        private Color BgHeader { get { return isWhiteTheme ? Color.FromArgb(255, 255, 255) : Color.FromArgb(18, 18, 20); } }
        private Color BgInput { get { return isWhiteTheme ? Color.FromArgb(255, 255, 255) : Color.FromArgb(34, 34, 38); } }
        private Color BorderSubtle { get { return isWhiteTheme ? Color.FromArgb(216, 216, 222) : Color.FromArgb(56, 56, 62); } }
        private Color TextPrimary { get { return isWhiteTheme ? Color.FromArgb(29, 29, 31) : Color.FromArgb(245, 245, 247); } }
        private Color TextSecondary { get { return isWhiteTheme ? Color.FromArgb(110, 110, 115) : Color.FromArgb(152, 152, 157); } }
        private static readonly Color AccentBlue = Color.FromArgb(10, 132, 255);
        private static readonly Color AccentBlueHover = Color.FromArgb(48, 150, 255);
        private static readonly Color AccentGreen = Color.FromArgb(48, 209, 88);
        private static readonly Color AccentAmber = Color.FromArgb(255, 159, 10);
        private static readonly Color AccentRed = Color.FromArgb(255, 69, 58);

        public AddDownloadDialog(string port, string url, string referer, string cookie, string userAgent, string fileName)
        {
            this.baseUrl = "http://127.0.0.1:" + port;
            this.refererUrl = referer ?? "";
            this.cookieHeader = cookie ?? "";
            this.userAgentHeader = userAgent ?? "";

            this.Text = "Loadin";
            this.FormBorderStyle = FormBorderStyle.None;
            this.StartPosition = FormStartPosition.CenterScreen;
            this.TopMost = true;
            this.DoubleBuffered = true;
            this.KeyPreview = true;
            this.BackColor = BgSurface;
            this.Size = new Size(560, 412);

            try
            {
                this.Icon = Icon.ExtractAssociatedIcon(Application.ExecutablePath);
            }
            catch { }

            ComputeLayoutRects();
            ApplyRoundedRegion(14);
            InitializeInputs(url, fileName);

            this.MouseDown += OnDialogMouseDown;
            this.MouseMove += OnDialogMouseMove;
            this.MouseClick += OnDialogMouseClick;
            this.MouseLeave += (s, e) =>
            {
                hoverCaptionBtn = 0;
                hoverActionBtn = 0;
                hoverThreadIdx = -1;
                this.Invalidate();
            };
            this.KeyDown += (s, e) =>
            {
                if (e.KeyCode == Keys.Escape)
                {
                    if (!isDownloadingMode) DismissClipboard();
                    this.Close();
                }
                else if (e.KeyCode == Keys.Enter && !isDownloadingMode)
                {
                    SubmitDownload(true);
                }
            };
            this.FormClosing += (s, e) =>
            {
                if (pollTimer != null)
                {
                    pollTimer.Stop();
                    pollTimer.Dispose();
                    pollTimer = null;
                }
            };

            ThreadPool.QueueUserWorkItem(_ => FetchDefaultsAndProbe(url));
        }

        private string T(string ko, string en, string ja)
        {
            if (lang == "en") return en;
            if (lang == "ja") return ja;
            return ko;
        }

        private void ComputeLayoutRects()
        {
            int w = this.Width;
            int h = this.Height;

            rectCloseBtn = new Rectangle(w - 42, 9, 32, 26);
            rectMaxBtn = new Rectangle(w - 76, 9, 32, 26);
            rectMinBtn = new Rectangle(w - 110, 9, 32, 26);

            rectBrowseBtn = new Rectangle(w - 108, 214, 84, 32);

            rectSegmentTrack = new Rectangle(24, 292, w - 48, 34);
            rectThreadPills = new Rectangle[threadOptions.Length];
            int segW = (rectSegmentTrack.Width - 6) / threadOptions.Length;
            for (int i = 0; i < threadOptions.Length; i++)
            {
                rectThreadPills[i] = new Rectangle(rectSegmentTrack.Left + 3 + i * segW, rectSegmentTrack.Top + 3, segW - 2, 28);
            }

            if (isDownloadingMode)
            {
                rectCancelBtn = new Rectangle(24, h - 54, 114, 36);
                rectQueueBtn = new Rectangle(w - 294, h - 54, 116, 36);
                rectStartBtn = new Rectangle(w - 168, h - 54, 144, 36);
            }
            else
            {
                rectCancelBtn = new Rectangle(24, h - 54, 88, 36);
                rectQueueBtn = new Rectangle(w - 270, h - 54, 116, 36);
                rectStartBtn = new Rectangle(w - 144, h - 54, 120, 36);
            }

            if (txtUrl != null) txtUrl.Width = w - 64;
            if (txtFileName != null) txtFileName.Width = w - 224;
            if (txtSaveDir != null) txtSaveDir.Width = w - 154;
        }

        private void ApplyRoundedRegion(int radius)
        {
            try
            {
                int dark = isWhiteTheme ? 0 : 1;
                DwmSetWindowAttribute(this.Handle, 20, ref dark, sizeof(int));
                int cornerPref = 2; // DWMWCP_ROUND
                int hr = DwmSetWindowAttribute(this.Handle, 33, ref cornerPref, sizeof(int));
                if (hr == 0)
                {
                    this.Region = null;
                    return;
                }
            }
            catch { }
            this.Region = Region.FromHrgn(CreateRoundRectRgn(0, 0, this.Width + 1, this.Height + 1, radius, radius));
        }

        private void InitializeInputs(string initialUrl, string initialFileName)
        {
            txtUrl = CreateCleanTextBox(32, 80, this.Width - 64, initialUrl);
            txtUrl.TextChanged += (s, e) =>
            {
                string u = txtUrl.Text.Trim();
                if (u.StartsWith("http://") || u.StartsWith("https://"))
                {
                    ThreadPool.QueueUserWorkItem(_ => ProbeUrlMetadata(u));
                }
            };

            string guessedName = initialFileName;
            if (string.IsNullOrEmpty(guessedName) && !string.IsNullOrEmpty(initialUrl))
            {
                try
                {
                    Uri u = new Uri(initialUrl);
                    guessedName = Path.GetFileName(Uri.UnescapeDataString(u.AbsolutePath));
                }
                catch { }
            }
            if (string.IsNullOrEmpty(guessedName)) guessedName = "download_file";

            txtFileName = CreateCleanTextBox(32, 150, this.Width - 224, guessedName);

            string defaultDir = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), "Downloads", "Loadin");
            txtSaveDir = CreateCleanTextBox(32, 221, this.Width - 154, defaultDir);

            this.Controls.Add(txtUrl);
            this.Controls.Add(txtFileName);
            this.Controls.Add(txtSaveDir);
        }

        private TextBox CreateCleanTextBox(int x, int y, int width, string text)
        {
            var tb = new TextBox();
            tb.BorderStyle = BorderStyle.None;
            tb.BackColor = BgInput;
            tb.ForeColor = TextPrimary;
            tb.Font = fontBody;
            tb.Location = new Point(x, y);
            tb.Width = width;
            tb.Text = text ?? "";
            return tb;
        }

        private void RefreshThemeColors()
        {
            this.BackColor = BgSurface;
            if (txtUrl != null) { txtUrl.BackColor = BgInput; txtUrl.ForeColor = TextPrimary; }
            if (txtFileName != null) { txtFileName.BackColor = BgInput; txtFileName.ForeColor = TextPrimary; }
            if (txtSaveDir != null) { txtSaveDir.BackColor = BgInput; txtSaveDir.ForeColor = TextPrimary; }
            ApplyRoundedRegion(14);
            this.Invalidate();
        }

        private void FetchDefaultsAndProbe(string url)
        {
            try
            {
                var req = (HttpWebRequest)WebRequest.Create(baseUrl + "/api/state");
                req.Timeout = 1500;
                using (var resp = (HttpWebResponse)req.GetResponse())
                using (var reader = new StreamReader(resp.GetResponseStream(), Encoding.UTF8))
                {
                    string json = reader.ReadToEnd();
                    string saveDir = ExtractJsonString(json, "defaultSaveDir");
                    string themeStr = ExtractJsonString(json, "theme");
                    string langStr = ExtractJsonString(json, "language");
                    int recConns = (int)ExtractJsonLong(json, "recommendedConnections");
                    long freeBytes = ExtractJsonLong(json, "freeBytes");

                    this.BeginInvoke((MethodInvoker)delegate
                    {
                        if (!string.IsNullOrEmpty(langStr)) lang = langStr;
                        isWhiteTheme = (themeStr == "white" || themeStr == "light-frost");
                        RefreshThemeColors();

                        if (!string.IsNullOrEmpty(saveDir)) txtSaveDir.Text = saveDir;
                        if (recConns > 0)
                        {
                            recommendedThreads = recConns;
                            selectedThreads = recConns;
                        }
                        if (freeBytes > 0)
                        {
                            diskFreeText = FormatBytes(freeBytes);
                        }
                        this.Invalidate();
                    });
                }
            }
            catch { }

            if (!string.IsNullOrWhiteSpace(url))
            {
                ProbeUrlMetadata(url.Trim());
            }
            else
            {
                this.BeginInvoke((MethodInvoker)delegate
                {
                    probeSizeText = T("대기 중", "Ready", "待機中");
                    this.Invalidate();
                });
            }
        }

        private void ProbeUrlMetadata(string targetUrl)
        {
            this.BeginInvoke((MethodInvoker)delegate
            {
                probeSizeText = T("분석 중...", "Probing...", "解析中...");
                this.Invalidate();
            });

            try
            {
                var req = (HttpWebRequest)WebRequest.Create(baseUrl + "/api/probe");
                req.Method = "POST";
                req.ContentType = "application/json; charset=utf-8";
                req.Timeout = 8000;
                string body = "{\"url\":\"" + EscapeJson(targetUrl) + "\"}";
                byte[] bytes = Encoding.UTF8.GetBytes(body);
                using (var s = req.GetRequestStream())
                {
                    s.Write(bytes, 0, bytes.Length);
                }
                using (var resp = (HttpWebResponse)req.GetResponse())
                using (var reader = new StreamReader(resp.GetResponseStream(), Encoding.UTF8))
                {
                    string json = reader.ReadToEnd();
                    string fn = ExtractJsonString(json, "fileName");
                    long totalSize = ExtractJsonLong(json, "totalSize");
                    bool ranges = json.Contains("\"acceptRanges\":true");

                    this.BeginInvoke((MethodInvoker)delegate
                    {
                        if (!string.IsNullOrEmpty(fn)) txtFileName.Text = fn;
                        probeSizeText = totalSize > 0 ? FormatBytes(totalSize) : T("크기 미상", "Unknown", "不明");
                        acceptRanges = ranges;
                        this.Invalidate();
                    });
                }
            }
            catch
            {
                this.BeginInvoke((MethodInvoker)delegate
                {
                    probeSizeText = T("다운로드 가능", "Ready", "準備完了");
                    this.Invalidate();
                });
            }
        }

        private void SubmitDownload(bool autoStart)
        {
            string url = txtUrl.Text.Trim();
            if (string.IsNullOrEmpty(url)) return;

            string createdTaskId = "";
            string createdFileName = txtFileName.Text.Trim();
            string createdSaveDir = txtSaveDir.Text.Trim();
            long createdTotalSize = 0;

            try
            {
                var req = (HttpWebRequest)WebRequest.Create(baseUrl + "/api/tasks");
                req.Method = "POST";
                req.ContentType = "application/json; charset=utf-8";
                req.Timeout = 4000;

                StringBuilder hdr = new StringBuilder("{");
                bool first = true;
                if (!string.IsNullOrEmpty(refererUrl))
                {
                    hdr.AppendFormat("\"Referer\":\"{0}\"", EscapeJson(refererUrl));
                    first = false;
                }
                if (!string.IsNullOrEmpty(cookieHeader))
                {
                    if (!first) hdr.Append(",");
                    hdr.AppendFormat("\"Cookie\":\"{0}\"", EscapeJson(cookieHeader));
                    first = false;
                }
                if (!string.IsNullOrEmpty(userAgentHeader))
                {
                    if (!first) hdr.Append(",");
                    hdr.AppendFormat("\"User-Agent\":\"{0}\"", EscapeJson(userAgentHeader));
                }
                hdr.Append("}");

                string payload = string.Format(
                    "{{\"url\":\"{0}\",\"fileName\":\"{1}\",\"saveDir\":\"{2}\",\"maxConnections\":{3},\"autoStart\":{4},\"headers\":{5},\"source\":\"Loadin Popup\"}}",
                    EscapeJson(url),
                    EscapeJson(createdFileName),
                    EscapeJson(createdSaveDir),
                    selectedThreads,
                    autoStart ? "true" : "false",
                    hdr.ToString()
                );

                byte[] bytes = Encoding.UTF8.GetBytes(payload);
                using (var s = req.GetRequestStream())
                {
                    s.Write(bytes, 0, bytes.Length);
                }
                using (var resp = (HttpWebResponse)req.GetResponse())
                using (var reader = new StreamReader(resp.GetResponseStream(), Encoding.UTF8))
                {
                    string respJson = reader.ReadToEnd();
                    createdTaskId = ExtractJsonString(respJson, "id");
                    string fn = ExtractJsonString(respJson, "fileName");
                    string sd = ExtractJsonString(respJson, "saveDir");
                    createdTotalSize = ExtractJsonLong(respJson, "totalSize");
                    if (!string.IsNullOrEmpty(fn)) createdFileName = fn;
                    if (!string.IsNullOrEmpty(sd)) createdSaveDir = sd;
                }
            }
            catch { }

            DismissClipboard();

            if (!autoStart || string.IsNullOrEmpty(createdTaskId))
            {
                this.Close();
                return;
            }

            EnterLiveDownloadMode(createdTaskId, createdFileName, createdSaveDir, createdTotalSize);
        }

        private void EnterLiveDownloadMode(string taskId, string fileName, string saveDir, long totalSize)
        {
            this.activeTaskId = taskId;
            this.dlFileName = fileName;
            this.dlSaveDir = saveDir;
            this.dlTotalSize = totalSize;
            this.dlMaxConns = selectedThreads;
            this.dlStatus = "running";
            this.isDownloadingMode = true;

            if (txtUrl != null) txtUrl.Visible = false;
            if (txtFileName != null) txtFileName.Visible = false;
            if (txtSaveDir != null) txtSaveDir.Visible = false;

            ComputeLayoutRects();
            this.Invalidate();

            pollTimer = new System.Windows.Forms.Timer();
            pollTimer.Interval = 200;
            pollTimer.Tick += (s, e) =>
            {
                if (Interlocked.CompareExchange(ref isPollingInFlight, 1, 0) == 0)
                {
                    ThreadPool.QueueUserWorkItem(_ => PollTaskProgress());
                }
            };
            pollTimer.Start();
            ThreadPool.QueueUserWorkItem(_ => PollTaskProgress());
        }

        private void PollTaskProgress()
        {
            try
            {
                if (string.IsNullOrEmpty(activeTaskId) || this.IsDisposed) return;

                var req = (HttpWebRequest)WebRequest.Create(baseUrl + "/api/tasks/item?id=" + Uri.EscapeDataString(activeTaskId));
                req.Timeout = 1200;
                using (var resp = (HttpWebResponse)req.GetResponse())
                using (var reader = new StreamReader(resp.GetResponseStream(), Encoding.UTF8))
                {
                    string json = reader.ReadToEnd();
                    string st = ExtractJsonString(json, "status");
                    string fn = ExtractJsonString(json, "fileName");
                    string sd = ExtractJsonString(json, "saveDir");
                    long downloaded = ExtractJsonLong(json, "downloaded");
                    long totalSize = ExtractJsonLong(json, "totalSize");
                    double prog = ExtractJsonDouble(json, "progress");
                    long speed = ExtractJsonLong(json, "speedBps");
                    long eta = ExtractJsonLong(json, "etaSeconds");
                    int workers = (int)ExtractJsonLong(json, "activeWorkers");
                    int maxConns = (int)ExtractJsonLong(json, "maxConnections");
                    int splits = (int)ExtractJsonLong(json, "dynamicSplits");
                    string errMsg = ExtractJsonString(json, "errorMessage");
                    float[] segPcts = ExtractSegmentProgress(json);

                    if (!this.IsDisposed && this.IsHandleCreated)
                    {
                        this.BeginInvoke((MethodInvoker)delegate
                        {
                            if (!string.IsNullOrEmpty(st)) dlStatus = st;
                            if (!string.IsNullOrEmpty(fn)) dlFileName = fn;
                            if (!string.IsNullOrEmpty(sd)) dlSaveDir = sd;
                            dlDownloaded = downloaded;
                            dlTotalSize = totalSize;
                            dlProgress = prog;
                            dlSpeedBps = speed;
                            dlEtaSeconds = eta;
                            dlActiveWorkers = workers;
                            if (maxConns > 0) dlMaxConns = maxConns;
                            dlDynamicSplits = splits;
                            dlErrorMsg = errMsg;
                            if (segPcts != null && segPcts.Length > 0) dlSegmentPcts = segPcts;

                            if (dlStatus == "completed")
                            {
                                dlProgress = 100.0;
                                if (dlTotalSize > 0) dlDownloaded = dlTotalSize;
                                if (pollTimer != null) pollTimer.Stop();
                            }
                            this.Invalidate();
                        });
                    }
                }
            }
            catch (WebException wex)
            {
                var httpResp = wex.Response as HttpWebResponse;
                if (httpResp != null && httpResp.StatusCode == HttpStatusCode.NotFound)
                {
                    if (!this.IsDisposed && this.IsHandleCreated)
                    {
                        this.BeginInvoke((MethodInvoker)delegate { this.Close(); });
                    }
                }
            }
            catch { }
            finally
            {
                Interlocked.Exchange(ref isPollingInFlight, 0);
            }
        }

        private void SendTaskAction(string action, bool deleteFiles, bool closeAfter)
        {
            string tid = activeTaskId;
            if (string.IsNullOrEmpty(tid)) return;

            ThreadPool.QueueUserWorkItem(_ =>
            {
                try
                {
                    var req = (HttpWebRequest)WebRequest.Create(baseUrl + "/api/tasks/action");
                    req.Method = "POST";
                    req.ContentType = "application/json; charset=utf-8";
                    req.Timeout = 2000;
                    string payload = string.Format(
                        "{{\"id\":\"{0}\",\"action\":\"{1}\",\"deleteFiles\":{2}}}",
                        EscapeJson(tid), EscapeJson(action), deleteFiles ? "true" : "false"
                    );
                    byte[] bytes = Encoding.UTF8.GetBytes(payload);
                    using (var s = req.GetRequestStream())
                    {
                        s.Write(bytes, 0, bytes.Length);
                    }
                    using (var resp = req.GetResponse()) { }
                }
                catch { }

                if (closeAfter && !this.IsDisposed && this.IsHandleCreated)
                {
                    this.BeginInvoke((MethodInvoker)delegate { this.Close(); });
                }
                else
                {
                    PollTaskProgress();
                }
            });
        }

        private void DismissClipboard()
        {
            try
            {
                var req = (HttpWebRequest)WebRequest.Create(baseUrl + "/api/clipboard/dismiss");
                req.Method = "POST";
                req.ContentLength = 0;
                req.Timeout = 800;
                using (var resp = req.GetResponse()) { }
            }
            catch { }
        }

        private void OnDialogMouseDown(object sender, MouseEventArgs e)
        {
            if (e.Button == MouseButtons.Left && e.Y <= 44 && hoverCaptionBtn == 0)
            {
                ReleaseCapture();
                SendMessage(this.Handle, WM_NCLBUTTONDOWN, HT_CAPTION, 0);
            }
        }

        private void OnDialogMouseMove(object sender, MouseEventArgs e)
        {
            int cap = 0;
            if (rectMinBtn.Contains(e.Location)) cap = 1;
            else if (rectMaxBtn.Contains(e.Location)) cap = 2;
            else if (rectCloseBtn.Contains(e.Location)) cap = 3;

            int act = 0;
            if (!isDownloadingMode && rectBrowseBtn.Contains(e.Location)) act = 1;
            else if (rectCancelBtn.Contains(e.Location)) act = 2;
            else if (rectQueueBtn.Contains(e.Location)) act = 3;
            else if (rectStartBtn.Contains(e.Location)) act = 4;

            int tIdx = -1;
            if (!isDownloadingMode)
            {
                for (int i = 0; i < rectThreadPills.Length; i++)
                {
                    if (rectThreadPills[i].Contains(e.Location))
                    {
                        tIdx = i;
                        break;
                    }
                }
            }

            if (cap != hoverCaptionBtn || act != hoverActionBtn || tIdx != hoverThreadIdx)
            {
                hoverCaptionBtn = cap;
                hoverActionBtn = act;
                hoverThreadIdx = tIdx;
                this.Cursor = (cap > 0 || act > 0 || tIdx >= 0) ? Cursors.Hand : Cursors.Default;
                this.Invalidate();
            }
        }

        private void OnDialogMouseClick(object sender, MouseEventArgs e)
        {
            if (e.Button != MouseButtons.Left) return;

            if (rectMinBtn.Contains(e.Location))
            {
                this.WindowState = FormWindowState.Minimized;
                return;
            }
            if (rectMaxBtn.Contains(e.Location))
            {
                isExpandedSize = !isExpandedSize;
                this.Size = isExpandedSize ? new Size(680, 412) : new Size(560, 412);
                ComputeLayoutRects();
                ApplyRoundedRegion(14);
                this.Invalidate();
                return;
            }
            if (rectCloseBtn.Contains(e.Location))
            {
                // Closing via top-right × sends active download to background without cancelling it
                if (!isDownloadingMode) DismissClipboard();
                this.Close();
                return;
            }

            if (isDownloadingMode)
            {
                if (dlStatus == "completed")
                {
                    if (rectCancelBtn.Contains(e.Location))
                    {
                        this.Close();
                        return;
                    }
                    if (rectQueueBtn.Contains(e.Location))
                    {
                        SendTaskAction("reveal", false, true);
                        return;
                    }
                    if (rectStartBtn.Contains(e.Location))
                    {
                        SendTaskAction("open", false, true);
                        return;
                    }
                }
                else
                {
                    if (rectCancelBtn.Contains(e.Location))
                    {
                        // Cancel & delete active download task
                        if (pollTimer != null) pollTimer.Stop();
                        SendTaskAction("delete", true, true);
                        return;
                    }
                    if (rectQueueBtn.Contains(e.Location))
                    {
                        // Toggle Pause / Resume
                        if (dlStatus == "paused" || dlStatus == "error")
                            SendTaskAction("resume", false, false);
                        else
                            SendTaskAction("pause", false, false);
                        return;
                    }
                    if (rectStartBtn.Contains(e.Location))
                    {
                        // Send to background (close popup window while engine continues downloading)
                        this.Close();
                        return;
                    }
                }
                return;
            }

            if (rectCancelBtn.Contains(e.Location))
            {
                DismissClipboard();
                this.Close();
                return;
            }
            if (rectBrowseBtn.Contains(e.Location))
            {
                using (var fbd = new FolderBrowserDialog())
                {
                    fbd.SelectedPath = txtSaveDir.Text;
                    if (fbd.ShowDialog(this) == DialogResult.OK)
                    {
                        txtSaveDir.Text = fbd.SelectedPath;
                    }
                }
                return;
            }
            for (int i = 0; i < rectThreadPills.Length; i++)
            {
                if (rectThreadPills[i].Contains(e.Location))
                {
                    selectedThreads = threadOptions[i];
                    this.Invalidate();
                    return;
                }
            }
            if (rectQueueBtn.Contains(e.Location))
            {
                SubmitDownload(false);
                return;
            }
            if (rectStartBtn.Contains(e.Location))
            {
                SubmitDownload(true);
                return;
            }
        }

        protected override void OnPaint(PaintEventArgs e)
        {
            Graphics g = e.Graphics;
            g.SmoothingMode = SmoothingMode.AntiAlias;
            g.TextRenderingHint = System.Drawing.Text.TextRenderingHint.ClearTypeGridFit;

            using (var bgBrush = new SolidBrush(BgSurface))
            {
                g.FillRectangle(bgBrush, this.ClientRectangle);
            }

            using (var headerBrush = new SolidBrush(BgHeader))
            using (var linePen = new Pen(BorderSubtle, 1f))
            {
                g.FillRectangle(headerBrush, 0, 0, this.Width, 44);
                g.DrawLine(linePen, 0, 44, this.Width, 44);
            }

            string headerTitle = T("Loadin — 새로운 다운로드", "Loadin — New Download", "Loadin — 新規ダウンロード");
            if (isDownloadingMode)
            {
                if (dlStatus == "completed")
                    headerTitle = T("Loadin — 다운로드 완료", "Loadin — Download Complete", "Loadin — ダウンロード完了");
                else if (dlStatus == "paused")
                    headerTitle = T("Loadin — 일시정지됨", "Loadin — Paused", "Loadin — 一時停止");
                else if (dlStatus == "error")
                    headerTitle = T("Loadin — 다운로드 오류", "Loadin — Download Error", "Loadin — エラー");
                else
                    headerTitle = T("Loadin — 다운로드 진행 중", "Loadin — Downloading", "Loadin — ダウンロード中");
            }

            using (var iconBg = new SolidBrush(dlStatus == "completed" && isDownloadingMode ? AccentGreen : AccentBlue))
            using (var whitePen = new Pen(Color.White, 1.8f))
            using (var titleBrush = new SolidBrush(TextPrimary))
            {
                Rectangle iconRect = new Rectangle(18, 12, 20, 20);
                g.FillPath(iconBg, RoundedRect(iconRect, 5));
                g.DrawLine(whitePen, 28, 15, 28, 27);
                g.DrawLine(whitePen, 24, 23, 28, 27);
                g.DrawLine(whitePen, 32, 23, 28, 27);
                g.DrawString(headerTitle, fontTitle, titleBrush, new PointF(46, 12));
            }

            PaintCaptionButtons(g);

            if (isDownloadingMode)
            {
                PaintLiveDownloadView(g);
            }
            else
            {
                PaintSetupFormView(g);
            }

            using (var framePen = new Pen(BorderSubtle, 1f))
            {
                g.DrawPath(framePen, RoundedRect(new Rectangle(0, 0, this.Width - 1, this.Height - 1), 14));
            }
        }

        private void PaintSetupFormView(Graphics g)
        {
            PaintFieldBox(g, T("다운로드 주소 (URL)", "Download URL", "ダウンロード URL"), new Rectangle(24, 74, this.Width - 48, 32));

            int fnBoxW = this.Width - 208;
            PaintFieldBox(g, T("파일 이름", "File Name", "ファイル名"), new Rectangle(24, 144, fnBoxW, 32));

            int metaLeft = 24 + fnBoxW + 8;
            Rectangle metaRect = new Rectangle(metaLeft, 144, this.Width - metaLeft - 24, 32);
            using (var metaBg = new SolidBrush(BgInput))
            using (var metaPen = new Pen(BorderSubtle, 1f))
            using (var valBrush = new SolidBrush(TextPrimary))
            using (var subBrush = new SolidBrush(TextSecondary))
            {
                g.DrawString(T("파일 크기", "File Size", "サイズ"), fontLabel, subBrush, new PointF(metaLeft, 123));
                g.FillPath(metaBg, RoundedRect(metaRect, 7));
                g.DrawPath(metaPen, RoundedRect(metaRect, 7));
                StringFormat sf = new StringFormat { Alignment = StringAlignment.Center, LineAlignment = StringAlignment.Center };
                g.DrawString(probeSizeText, fontSmallBold, valBrush, metaRect, sf);
            }

            string dirLabel = string.IsNullOrEmpty(diskFreeText)
                ? T("저장 위치", "Save Location", "保存先")
                : T("저장 위치 (여유 ", "Save Location (Free ", "保存先 (空き ") + diskFreeText + ")";
            PaintFieldBox(g, dirLabel, new Rectangle(24, 214, this.Width - 140, 32));

            using (var browseBg = new SolidBrush(hoverActionBtn == 1 ? (isWhiteTheme ? Color.FromArgb(230, 230, 236) : Color.FromArgb(52, 52, 58)) : BgInput))
            using (var browseBorder = new Pen(BorderSubtle, 1f))
            using (var txtBrush = new SolidBrush(TextPrimary))
            {
                g.FillPath(browseBg, RoundedRect(rectBrowseBtn, 7));
                g.DrawPath(browseBorder, RoundedRect(rectBrowseBtn, 7));
                StringFormat sf = new StringFormat { Alignment = StringAlignment.Center, LineAlignment = StringAlignment.Center };
                g.DrawString(T("변경...", "Browse...", "変更..."), fontSmallBold, txtBrush, rectBrowseBtn, sf);
            }

            using (var labelBrush = new SolidBrush(TextSecondary))
            using (var recBrush = new SolidBrush(AccentBlue))
            {
                g.DrawString(T("병렬 스레드 수", "Parallel Threads", "並列スレッド数"), fontLabel, labelBrush, new PointF(24, 270));
                StringFormat rightSf = new StringFormat { Alignment = StringAlignment.Far };
                string recText = string.Format(
                    T("사양 권장: {0} 스레드", "Recommended: {0} Threads", "推奨: {0} スレッド"),
                    recommendedThreads);
                g.DrawString(recText, fontSmallBold, recBrush, new PointF(this.Width - 24, 270), rightSf);
            }

            using (var trackBg = new SolidBrush(BgInput))
            using (var trackBorder = new Pen(BorderSubtle, 1f))
            {
                g.FillPath(trackBg, RoundedRect(rectSegmentTrack, 8));
                g.DrawPath(trackBorder, RoundedRect(rectSegmentTrack, 8));
            }

            StringFormat centerSf = new StringFormat { Alignment = StringAlignment.Center, LineAlignment = StringAlignment.Center };
            for (int i = 0; i < threadOptions.Length; i++)
            {
                int tVal = threadOptions[i];
                bool isSelected = (selectedThreads == tVal);
                bool isRec = (recommendedThreads == tVal);
                bool isHovered = (hoverThreadIdx == i);

                Rectangle r = rectThreadPills[i];
                if (isSelected)
                {
                    using (var selBg = new SolidBrush(AccentBlue))
                    {
                        g.FillPath(selBg, RoundedRect(r, 6));
                    }
                }
                else if (isHovered)
                {
                    using (var hovBg = new SolidBrush(isWhiteTheme ? Color.FromArgb(232, 232, 238) : Color.FromArgb(50, 50, 56)))
                    {
                        g.FillPath(hovBg, RoundedRect(r, 6));
                    }
                }

                using (var fg = new SolidBrush(isSelected ? Color.White : TextPrimary))
                {
                    string recSuffix = isRec ? T(" (권장)", " (Rec)", " (推奨)") : "";
                    string label = tVal.ToString() + recSuffix;
                    g.DrawString(label, isSelected ? fontSmallBold : fontSmall, fg, r, centerSf);
                }
            }

            using (var linePen = new Pen(BorderSubtle, 1f))
            {
                g.DrawLine(linePen, 0, this.Height - 68, this.Width, this.Height - 68);
            }

            PaintActionButton(g, rectCancelBtn, T("취소", "Cancel", "キャンセル"), hoverActionBtn == 2, false);
            PaintActionButton(g, rectQueueBtn, T("대기열 추가", "Add to Queue", "キューに追加"), hoverActionBtn == 3, false);
            PaintActionButton(g, rectStartBtn, T("다운로드", "Download", "ダウンロード"), hoverActionBtn == 4, true);
        }

        private void PaintLiveDownloadView(Graphics g)
        {
            int padX = 24;
            int contentW = this.Width - padX * 2;

            // 1. File Summary Card (Y: 58..114)
            Rectangle fileCard = new Rectangle(padX, 58, contentW, 56);
            using (var cardBg = new SolidBrush(BgInput))
            using (var cardPen = new Pen(BorderSubtle, 1f))
            using (var primaryBrush = new SolidBrush(TextPrimary))
            using (var secBrush = new SolidBrush(TextSecondary))
            {
                g.FillPath(cardBg, RoundedRect(fileCard, 9));
                g.DrawPath(cardPen, RoundedRect(fileCard, 9));

                StringFormat ellipsisSf = new StringFormat
                {
                    Trimming = StringTrimming.EllipsisPath,
                    FormatFlags = StringFormatFlags.NoWrap
                };

                Rectangle fnRect = new Rectangle(fileCard.Left + 14, fileCard.Top + 9, fileCard.Width - 120, 20);
                Rectangle dirRect = new Rectangle(fileCard.Left + 14, fileCard.Top + 31, fileCard.Width - 120, 18);
                g.DrawString(string.IsNullOrEmpty(dlFileName) ? "download_file" : dlFileName, fontTitle, primaryBrush, fnRect, ellipsisSf);
                g.DrawString(dlSaveDir ?? "", fontSmall, secBrush, dirRect, ellipsisSf);

                // Status Pill Badge on right side of file card
                string statusText = T("다운로드 중", "Downloading", "ダウンロード中");
                Color badgeColor = AccentBlue;
                if (dlStatus == "completed")
                {
                    statusText = T("완료됨", "Completed", "完了");
                    badgeColor = AccentGreen;
                }
                else if (dlStatus == "paused")
                {
                    statusText = T("일시정지", "Paused", "一時停止");
                    badgeColor = AccentAmber;
                }
                else if (dlStatus == "verifying")
                {
                    statusText = T("무결성 검증", "Verifying", "検証中");
                    badgeColor = AccentBlue;
                }
                else if (dlStatus == "error")
                {
                    statusText = T("오류 발생", "Error", "エラー");
                    badgeColor = AccentRed;
                }

                Rectangle badgeRect = new Rectangle(fileCard.Right - 96, fileCard.Top + 15, 82, 26);
                using (var badgeBg = new SolidBrush(Color.FromArgb(isWhiteTheme ? 28 : 42, badgeColor)))
                using (var badgeFg = new SolidBrush(badgeColor))
                {
                    g.FillPath(badgeBg, RoundedRect(badgeRect, 6));
                    StringFormat csf = new StringFormat { Alignment = StringAlignment.Center, LineAlignment = StringAlignment.Center };
                    g.DrawString(statusText, fontSmallBold, badgeFg, badgeRect, csf);
                }
            }

            // 2. Main Progress & Multi-Segment Ribbon (Y: 126..210)
            double clampedProg = Math.Max(0.0, Math.Min(100.0, dlProgress));
            Color barColor = dlStatus == "completed" ? AccentGreen : (dlStatus == "paused" ? AccentAmber : AccentBlue);

            using (var secBrush = new SolidBrush(TextSecondary))
            using (var pctBrush = new SolidBrush(TextPrimary))
            {
                g.DrawString(T("전체 진행률", "Overall Progress", "全体の進捗"), fontLabel, secBrush, new PointF(padX, 128));
                StringFormat rightSf = new StringFormat { Alignment = StringAlignment.Far };
                g.DrawString(string.Format("{0:F1}%", clampedProg), fontPctLarge, pctBrush, new PointF(this.Width - padX, 122), rightSf);
            }

            Rectangle mainTrack = new Rectangle(padX, 150, contentW, 12);
            using (var trackBg = new SolidBrush(BgInput))
            using (var trackPen = new Pen(BorderSubtle, 1f))
            {
                g.FillPath(trackBg, RoundedRect(mainTrack, 6));
                g.DrawPath(trackPen, RoundedRect(mainTrack, 6));
            }

            int fillW = (int)Math.Round(mainTrack.Width * (clampedProg / 100.0));
            if (fillW > 0)
            {
                fillW = Math.Max(8, Math.Min(mainTrack.Width, fillW));
                Rectangle fillRect = new Rectangle(mainTrack.Left, mainTrack.Top, fillW, mainTrack.Height);
                using (var fillBrush = new SolidBrush(barColor))
                {
                    g.FillPath(fillBrush, RoundedRect(fillRect, 6));
                }
            }

            // Multi-Segment Thread Ribbon
            Rectangle segRibbon = new Rectangle(padX, 170, contentW, 10);
            using (var ribbonBg = new SolidBrush(BgInput))
            using (var ribbonPen = new Pen(BorderSubtle, 1f))
            {
                g.FillPath(ribbonBg, RoundedRect(segRibbon, 4));
                g.DrawPath(ribbonPen, RoundedRect(segRibbon, 4));
            }

            int segCount = (dlSegmentPcts != null && dlSegmentPcts.Length > 0) ? dlSegmentPcts.Length : Math.Max(1, dlMaxConns);
            float slotW = (float)(segRibbon.Width - 4) / segCount;
            using (var segBrush = new SolidBrush(barColor))
            {
                for (int i = 0; i < segCount; i++)
                {
                    float pct = (dlSegmentPcts != null && i < dlSegmentPcts.Length)
                        ? dlSegmentPcts[i]
                        : (float)(clampedProg / 100.0);
                    pct = Math.Max(0f, Math.Min(1f, pct));
                    float sx = segRibbon.Left + 2 + i * slotW;
                    float maxSw = Math.Max(1.5f, slotW - 1.5f);
                    float sw = maxSw * pct;
                    if (sw >= 1.0f)
                    {
                        g.FillRectangle(segBrush, sx, segRibbon.Top + 2, sw, segRibbon.Height - 4);
                    }
                }
            }

            // 3. 2x2 Live Telemetry Grid (Y: 192..314)
            int gap = 10;
            int cellW = (contentW - gap) / 2;
            int cellH = 54;

            string speedStr = dlStatus == "completed"
                ? T("완료됨", "Completed", "完了")
                : (dlSpeedBps > 0 ? FormatBytes(dlSpeedBps) + "/s" : (dlStatus == "paused" ? T("일시정지", "Paused", "一時停止") : "0 B/s"));

            string etaStr = dlStatus == "completed"
                ? "00:00"
                : FormatEta(dlEtaSeconds);

            string sizeStr = dlTotalSize > 0
                ? string.Format("{0} / {1}", FormatBytes(dlDownloaded), FormatBytes(dlTotalSize))
                : FormatBytes(dlDownloaded);

            string threadsStr = string.Format(
                T("{0} / {1} 스레드 (분할 {2}회)", "{0} / {1} Threads ({2} splits)", "{0} / {1} スレッド (分割 {2})"),
                dlStatus == "completed" ? dlMaxConns : dlActiveWorkers,
                dlMaxConns,
                dlDynamicSplits);

            PaintMetricCell(g, new Rectangle(padX, 192, cellW, cellH),
                T("전송 속도", "Transfer Speed", "転送速度"), speedStr, true);
            PaintMetricCell(g, new Rectangle(padX + cellW + gap, 192, cellW, cellH),
                T("남은 시간", "Estimated Time", "残り時間"), etaStr, false);
            PaintMetricCell(g, new Rectangle(padX, 192 + cellH + gap, cellW, cellH),
                T("받은 용량 / 전체 크기", "Downloaded / Total", "受信済み / 合計"), sizeStr, false);
            PaintMetricCell(g, new Rectangle(padX + cellW + gap, 192 + cellH + gap, cellW, cellH),
                T("병렬 스레드 연결", "Active Connections", "並列スレッド接続"), threadsStr, false);

            // 4. Background Hint Caption
            using (var hintBrush = new SolidBrush(TextSecondary))
            {
                string hint = !string.IsNullOrEmpty(dlErrorMsg)
                    ? dlErrorMsg
                    : (dlStatus == "completed"
                        ? T("파일 다운로드 및 병합이 완료되었습니다.", "File download and assembly completed.", "ダウンロードが完了しました。")
                        : T("창을 닫거나 백그라운드로 전환해도 다운로드가 중단 없이 계속 진행됩니다.",
                            "Closing this window keeps the download running seamlessly in the background.",
                            "ウィンドウを閉じてもバックグラウンドでダウンロードが継続されます。"));
                g.DrawString(hint, fontSmall, hintBrush, new PointF(padX, 318));
            }

            using (var linePen = new Pen(BorderSubtle, 1f))
            {
                g.DrawLine(linePen, 0, this.Height - 68, this.Width, this.Height - 68);
            }

            if (dlStatus == "completed")
            {
                PaintActionButton(g, rectCancelBtn, T("창 닫기", "Close", "閉じる"), hoverActionBtn == 2, false);
                PaintActionButton(g, rectQueueBtn, T("폴더 열기", "Open Folder", "フォルダを開く"), hoverActionBtn == 3, false);
                PaintActionButton(g, rectStartBtn, T("파일 열기", "Open File", "ファイルを開く"), hoverActionBtn == 4, true);
            }
            else
            {
                string pauseResumeLabel = (dlStatus == "paused" || dlStatus == "error")
                    ? T("이어받기", "Resume", "再開")
                    : T("일시정지", "Pause", "一時停止");

                PaintActionButton(g, rectCancelBtn, T("다운로드 취소", "Cancel Task", "中止"), hoverActionBtn == 2, false);
                PaintActionButton(g, rectQueueBtn, pauseResumeLabel, hoverActionBtn == 3, false);
                PaintActionButton(g, rectStartBtn, T("백그라운드 전환", "Background", "バックグラウンド"), hoverActionBtn == 4, true);
            }
        }

        private void PaintMetricCell(Graphics g, Rectangle rect, string label, string value, bool highlight)
        {
            using (var bg = new SolidBrush(BgInput))
            using (var pen = new Pen(BorderSubtle, 1f))
            using (var lblBrush = new SolidBrush(TextSecondary))
            using (var valBrush = new SolidBrush(highlight ? AccentBlue : TextPrimary))
            {
                g.FillPath(bg, RoundedRect(rect, 8));
                g.DrawPath(pen, RoundedRect(rect, 8));
                g.DrawString(label, fontSmall, lblBrush, new PointF(rect.Left + 12, rect.Top + 7));
                g.DrawString(value, fontMetricVal, valBrush, new PointF(rect.Left + 12, rect.Top + 26));
            }
        }

        private void PaintCaptionButtons(Graphics g)
        {
            Color iconColor = isWhiteTheme ? Color.FromArgb(60, 60, 67) : Color.FromArgb(210, 210, 215);
            Color hoverBg = isWhiteTheme ? Color.FromArgb(228, 228, 234) : Color.FromArgb(44, 44, 48);

            if (hoverCaptionBtn == 1)
            {
                using (var b = new SolidBrush(hoverBg))
                    g.FillPath(b, RoundedRect(rectMinBtn, 6));
            }
            using (var pen = new Pen(iconColor, 1.5f))
            {
                int cy = rectMinBtn.Top + rectMinBtn.Height / 2;
                g.DrawLine(pen, rectMinBtn.Left + 10, cy, rectMinBtn.Right - 10, cy);
            }

            if (hoverCaptionBtn == 2)
            {
                using (var b = new SolidBrush(hoverBg))
                    g.FillPath(b, RoundedRect(rectMaxBtn, 6));
            }
            using (var pen = new Pen(iconColor, 1.4f))
            {
                g.DrawRectangle(pen, rectMaxBtn.Left + 11, rectMaxBtn.Top + 8, 10, 10);
            }

            if (hoverCaptionBtn == 3)
            {
                using (var b = new SolidBrush(Color.FromArgb(232, 17, 35)))
                    g.FillPath(b, RoundedRect(rectCloseBtn, 6));
            }
            using (var pen = new Pen(hoverCaptionBtn == 3 ? Color.White : iconColor, 1.5f))
            {
                g.DrawLine(pen, rectCloseBtn.Left + 11, rectCloseBtn.Top + 8, rectCloseBtn.Right - 11, rectCloseBtn.Bottom - 8);
                g.DrawLine(pen, rectCloseBtn.Right - 11, rectCloseBtn.Top + 8, rectCloseBtn.Left + 11, rectCloseBtn.Bottom - 8);
            }
        }

        private void PaintFieldBox(Graphics g, string label, Rectangle boxRect)
        {
            using (var labelBrush = new SolidBrush(TextSecondary))
            using (var boxBg = new SolidBrush(BgInput))
            using (var boxPen = new Pen(BorderSubtle, 1f))
            {
                g.DrawString(label, fontLabel, labelBrush, new PointF(boxRect.Left, boxRect.Top - 21));
                g.FillPath(boxBg, RoundedRect(boxRect, 7));
                g.DrawPath(boxPen, RoundedRect(boxRect, 7));
            }
        }

        private void PaintActionButton(Graphics g, Rectangle rect, string text, bool hovered, bool isPrimary)
        {
            StringFormat sf = new StringFormat { Alignment = StringAlignment.Center, LineAlignment = StringAlignment.Center };
            if (isPrimary)
            {
                Color c = hovered ? AccentBlueHover : AccentBlue;
                using (var b = new SolidBrush(c))
                using (var fg = new SolidBrush(Color.White))
                {
                    g.FillPath(b, RoundedRect(rect, 7));
                    g.DrawString(text, fontTitle, fg, rect, sf);
                }
            }
            else
            {
                Color bg = hovered
                    ? (isWhiteTheme ? Color.FromArgb(230, 230, 236) : Color.FromArgb(52, 52, 58))
                    : BgInput;
                using (var b = new SolidBrush(bg))
                using (var p = new Pen(BorderSubtle, 1f))
                using (var fg = new SolidBrush(TextPrimary))
                {
                    g.FillPath(b, RoundedRect(rect, 7));
                    g.DrawPath(p, RoundedRect(rect, 7));
                    g.DrawString(text, fontBody, fg, rect, sf);
                }
            }
        }

        private static GraphicsPath RoundedRect(Rectangle bounds, int radius)
        {
            int d = radius * 2;
            GraphicsPath path = new GraphicsPath();
            path.AddArc(bounds.Left, bounds.Top, d, d, 180, 90);
            path.AddArc(bounds.Right - d, bounds.Top, d, d, 270, 90);
            path.AddArc(bounds.Right - d, bounds.Bottom - d, d, d, 0, 90);
            path.AddArc(bounds.Left, bounds.Bottom - d, d, d, 90, 90);
            path.CloseFigure();
            return path;
        }

        private static long ExtractJsonLong(string json, string key)
        {
            var m = Regex.Match(json, "\"" + key + "\"\\s*:\\s*(-?\\d+)");
            if (m.Success)
            {
                long v;
                if (long.TryParse(m.Groups[1].Value, out v)) return v;
            }
            return 0;
        }

        private static double ExtractJsonDouble(string json, string key)
        {
            var m = Regex.Match(json, "\"" + key + "\"\\s*:\\s*(-?\\d+(?:\\.\\d+)?)");
            if (m.Success)
            {
                double v;
                if (double.TryParse(m.Groups[1].Value, NumberStyles.Float, CultureInfo.InvariantCulture, out v))
                    return v;
            }
            return 0.0;
        }

        private static string ExtractJsonString(string json, string key)
        {
            var m = Regex.Match(json, "\"" + key + "\"\\s*:\\s*\"((?:[^\"\\\\]|\\\\.)*)\"");
            if (m.Success) return Regex.Unescape(m.Groups[1].Value);
            return "";
        }

        private static float[] ExtractSegmentProgress(string json)
        {
            int idx = json.IndexOf("\"segments\"");
            if (idx < 0) return new float[0];
            string sub = json.Substring(idx);
            var matches = Regex.Matches(sub, "\\{[^{}]*\"start\"\\s*:\\s*(\\d+)\\s*,\\s*\"end\"\\s*:\\s*(-?\\d+)\\s*,\\s*\"current\"\\s*:\\s*(\\d+)[^{}]*\\}");
            if (matches.Count == 0) return new float[0];

            List<float> list = new List<float>(matches.Count);
            foreach (Match m in matches)
            {
                long start, end, current;
                long.TryParse(m.Groups[1].Value, out start);
                long.TryParse(m.Groups[2].Value, out end);
                long.TryParse(m.Groups[3].Value, out current);
                long len = end - start + 1;
                if (len > 0)
                {
                    float p = (float)(current - start) / (float)len;
                    list.Add(Math.Max(0f, Math.Min(1f, p)));
                }
                else
                {
                    list.Add(current > start ? 1f : 0f);
                }
            }
            return list.ToArray();
        }

        private static string FormatEta(long sec)
        {
            if (sec < 0 || sec > 86400 * 7) return "--:--";
            long h = sec / 3600;
            long m = (sec % 3600) / 60;
            long s = sec % 60;
            if (h > 0) return string.Format("{0:D2}:{1:D2}:{2:D2}", h, m, s);
            return string.Format("{0:D2}:{1:D2}", m, s);
        }

        private static string FormatBytes(long b)
        {
            if (b <= 0) return "0 B";
            string[] u = { "B", "KB", "MB", "GB", "TB" };
            int i = (int)Math.Floor(Math.Log(b) / Math.Log(1024));
            if (i >= u.Length) i = u.Length - 1;
            double v = b / Math.Pow(1024, i);
            return string.Format(i == 0 ? "{0:F0} {1}" : "{0:F2} {1}", v, u[i]);
        }

        private static string EscapeJson(string s)
        {
            return (s ?? "").Replace("\\", "\\\\").Replace("\"", "\\\"");
        }
    }
}
