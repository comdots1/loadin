using System;
using System.Diagnostics;
using System.Drawing;
using System.Drawing.Drawing2D;
using System.IO;
using System.Net;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
using System.Windows.Forms;
using Microsoft.Win32;

namespace LoadinInstaller
{
    public class SetupForm : Form
    {
        private delegate bool EnumWindowsProc(IntPtr hWnd, IntPtr lParam);

        [DllImport("user32.dll")]
        private static extern bool ReleaseCapture();

        [DllImport("user32.dll")]
        private static extern IntPtr SendMessage(IntPtr hWnd, int Msg, int wParam, int lParam);

        [DllImport("user32.dll")]
        private static extern bool PostMessage(IntPtr hWnd, uint Msg, IntPtr wParam, IntPtr lParam);

        [DllImport("user32.dll")]
        private static extern bool EnumWindows(EnumWindowsProc lpEnumFunc, IntPtr lParam);

        [DllImport("user32.dll", CharSet = CharSet.Unicode)]
        private static extern int GetWindowText(IntPtr hWnd, StringBuilder lpString, int nMaxCount);

        [DllImport("user32.dll")]
        private static extern bool IsWindowVisible(IntPtr hWnd);

        [DllImport("dwmapi.dll")]
        private static extern int DwmSetWindowAttribute(IntPtr hwnd, int attr, ref int attrValue, int attrSize);

        private const uint WM_CLOSE = 0x0010;

        private TextBox txtInstallDir;
        private Button btnBrowse;
        private CheckBox chkDesktopShortcut;
        private CheckBox chkStartMenuShortcut;
        private CheckBox chkStartupRun;
        private CheckBox chkRegisterProtocol;
        private CheckBox chkLaunchAfter;
        private ProgressBar progressBar;
        private Label lblStatus;
        private Button btnUninstall;
        private Button btnInstall;
        private Button btnCancel;
        private ComboBox cmbLang;
        private Label lblTitle;
        private Label lblSubtitle;
        private Label lblDirCaption;

        private string lang = "ko";

        public SetupForm()
        {
            Text = "Loadin Setup";
            FormBorderStyle = FormBorderStyle.None;
            StartPosition = FormStartPosition.CenterScreen;
            Size = new Size(580, 470);
            BackColor = Color.FromArgb(18, 18, 20);
            ForeColor = Color.FromArgb(245, 245, 247);
            DoubleBuffered = true;

            try
            {
                Icon appIcon = ExtractEmbeddedIcon();
                if (appIcon != null) Icon = appIcon;
            }
            catch { }

            Font titleFont = CreateFont(13f, FontStyle.Bold);
            Font subFont = CreateFont(9.2f, FontStyle.Regular);
            Font bodyFont = CreateFont(9.5f, FontStyle.Regular);
            Font btnFont = CreateFont(9.8f, FontStyle.Bold);

            lblTitle = new Label
            {
                Text = "Loadin 설치 마법사",
                Font = titleFont,
                ForeColor = Color.FromArgb(245, 245, 247),
                BackColor = Color.Transparent,
                Location = new Point(62, 20),
                Size = new Size(340, 26)
            };
            lblTitle.MouseDown += DragWindow;
            Controls.Add(lblTitle);

            lblSubtitle = new Label
            {
                Text = "고속 멀티 세그먼트 다운로드 매니저 (v2.1 Commercial Edition)",
                Font = subFont,
                ForeColor = Color.FromArgb(152, 152, 157),
                BackColor = Color.Transparent,
                Location = new Point(62, 46),
                Size = new Size(380, 20)
            };
            lblSubtitle.MouseDown += DragWindow;
            Controls.Add(lblSubtitle);

            cmbLang = new ComboBox
            {
                DropDownStyle = ComboBoxStyle.DropDownList,
                FlatStyle = FlatStyle.Flat,
                BackColor = Color.FromArgb(34, 34, 38),
                ForeColor = Color.FromArgb(245, 245, 247),
                Font = subFont,
                Location = new Point(435, 22),
                Size = new Size(92, 26)
            };
            cmbLang.Items.AddRange(new object[] { "한국어", "English", "日本語" });
            cmbLang.SelectedIndex = 0;
            cmbLang.SelectedIndexChanged += (s, e) =>
            {
                lang = cmbLang.SelectedIndex == 1 ? "en" : (cmbLang.SelectedIndex == 2 ? "ja" : "ko");
                ApplyLanguage();
            };
            Controls.Add(cmbLang);

            Button btnCloseTop = new Button
            {
                Text = "✕",
                Font = CreateFont(9.5f, FontStyle.Regular),
                ForeColor = Color.FromArgb(160, 160, 168),
                BackColor = Color.Transparent,
                FlatStyle = FlatStyle.Flat,
                Location = new Point(536, 18),
                Size = new Size(28, 28),
                Cursor = Cursors.Hand,
                TabStop = false
            };
            btnCloseTop.FlatAppearance.BorderSize = 0;
            btnCloseTop.FlatAppearance.MouseOverBackColor = Color.FromArgb(255, 69, 58);
            btnCloseTop.Click += (s, e) => Close();
            Controls.Add(btnCloseTop);

            lblDirCaption = new Label
            {
                Text = "설치 폴더 위치",
                Font = CreateFont(9.2f, FontStyle.Bold),
                ForeColor = Color.FromArgb(165, 165, 172),
                BackColor = Color.Transparent,
                Location = new Point(28, 90),
                Size = new Size(300, 20)
            };
            Controls.Add(lblDirCaption);

            string defaultDir = Path.Combine(
                Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
                "Programs",
                "Loadin"
            );

            Panel dirPanel = new Panel
            {
                Location = new Point(28, 114),
                Size = new Size(428, 36),
                BackColor = Color.FromArgb(28, 28, 32)
            };
            dirPanel.Paint += (s, e) =>
            {
                using (Pen p = new Pen(Color.FromArgb(52, 52, 58), 1))
                {
                    e.Graphics.DrawRectangle(p, 0, 0, dirPanel.Width - 1, dirPanel.Height - 1);
                }
            };

            txtInstallDir = new TextBox
            {
                BorderStyle = BorderStyle.None,
                BackColor = Color.FromArgb(28, 28, 32),
                ForeColor = Color.FromArgb(245, 245, 247),
                Font = bodyFont,
                Location = new Point(10, 9),
                Size = new Size(408, 20),
                Text = defaultDir
            };
            dirPanel.Controls.Add(txtInstallDir);
            Controls.Add(dirPanel);

            btnBrowse = new Button
            {
                Text = "찾아보기...",
                Font = subFont,
                ForeColor = Color.FromArgb(245, 245, 247),
                BackColor = Color.FromArgb(44, 44, 48),
                FlatStyle = FlatStyle.Flat,
                Location = new Point(464, 114),
                Size = new Size(88, 36),
                Cursor = Cursors.Hand
            };
            btnBrowse.FlatAppearance.BorderColor = Color.FromArgb(64, 64, 70);
            btnBrowse.Click += OnBrowseFolder;
            Controls.Add(btnBrowse);

            chkDesktopShortcut = CreateOptionCheck("바탕화면에 Loadin 바로가기 만들기", 170, true, bodyFont);
            chkStartMenuShortcut = CreateOptionCheck("시작 메뉴 프로그램 목록에 Loadin 추가", 204, true, bodyFont);
            chkStartupRun = CreateOptionCheck("윈도우 시작 시 자동 실행 (백그라운드 브라우저 연동 대기)", 238, true, bodyFont);
            chkRegisterProtocol = CreateOptionCheck("웹 브라우저 연동 프로토콜(loadin://) 및 확장프로그램 설치", 272, true, bodyFont);
            chkLaunchAfter = CreateOptionCheck("설치 완료 후 Loadin 바로 실행", 306, true, bodyFont);

            Controls.Add(chkDesktopShortcut);
            Controls.Add(chkStartMenuShortcut);
            Controls.Add(chkStartupRun);
            Controls.Add(chkRegisterProtocol);
            Controls.Add(chkLaunchAfter);

            lblStatus = new Label
            {
                Text = "설치 준비가 완료되었습니다. [설치 시작]을 클릭하세요.",
                Font = subFont,
                ForeColor = Color.FromArgb(152, 152, 157),
                BackColor = Color.Transparent,
                Location = new Point(28, 352),
                Size = new Size(524, 20)
            };
            Controls.Add(lblStatus);

            progressBar = new ProgressBar
            {
                Location = new Point(28, 376),
                Size = new Size(524, 10),
                Style = ProgressBarStyle.Continuous,
                Minimum = 0,
                Maximum = 100,
                Value = 0
            };
            Controls.Add(progressBar);

            btnUninstall = new Button
            {
                Text = "프로그램 제거",
                Font = bodyFont,
                ForeColor = Color.FromArgb(255, 85, 75),
                BackColor = Color.FromArgb(40, 26, 28),
                FlatStyle = FlatStyle.Flat,
                Location = new Point(28, 408),
                Size = new Size(118, 40),
                Cursor = Cursors.Hand
            };
            btnUninstall.FlatAppearance.BorderColor = Color.FromArgb(95, 45, 45);
            btnUninstall.Click += (s, e) =>
            {
                RunUninstall();
                Close();
            };
            Controls.Add(btnUninstall);

            btnCancel = new Button
            {
                Text = "취소",
                Font = bodyFont,
                ForeColor = Color.FromArgb(215, 215, 220),
                BackColor = Color.FromArgb(42, 42, 46),
                FlatStyle = FlatStyle.Flat,
                Location = new Point(332, 408),
                Size = new Size(96, 40),
                Cursor = Cursors.Hand
            };
            btnCancel.FlatAppearance.BorderColor = Color.FromArgb(62, 62, 68);
            btnCancel.Click += (s, e) => Close();
            Controls.Add(btnCancel);

            btnInstall = new Button
            {
                Text = "설치 시작",
                Font = btnFont,
                ForeColor = Color.White,
                BackColor = Color.FromArgb(10, 132, 255),
                FlatStyle = FlatStyle.Flat,
                Location = new Point(436, 408),
                Size = new Size(116, 40),
                Cursor = Cursors.Hand
            };
            btnInstall.FlatAppearance.BorderSize = 0;
            btnInstall.FlatAppearance.MouseOverBackColor = Color.FromArgb(40, 150, 255);
            btnInstall.Click += OnRunInstall;
            Controls.Add(btnInstall);

            MouseDown += DragWindow;
        }

        private CheckBox CreateOptionCheck(string text, int y, bool isChecked, Font font)
        {
            return new CheckBox
            {
                Text = text,
                Checked = isChecked,
                Font = font,
                ForeColor = Color.FromArgb(235, 235, 240),
                BackColor = Color.Transparent,
                Location = new Point(30, y),
                Size = new Size(520, 26),
                Cursor = Cursors.Hand
            };
        }

        private void ApplyLanguage()
        {
            if (lang == "en")
            {
                lblTitle.Text = "Loadin Setup Wizard";
                lblSubtitle.Text = "High-Performance Multi-Segment Download Manager (v2.1)";
                lblDirCaption.Text = "Installation Folder";
                btnBrowse.Text = "Browse...";
                chkDesktopShortcut.Text = "Create Loadin shortcut on Desktop";
                chkStartMenuShortcut.Text = "Add Loadin to Windows Start Menu";
                chkStartupRun.Text = "Launch automatically on Windows Startup (Background Engine)";
                chkRegisterProtocol.Text = "Register browser protocol (loadin://) & unpack browser extension";
                chkLaunchAfter.Text = "Launch Loadin when installation finishes";
                lblStatus.Text = "Ready to install. Click [Install Now] to begin.";
                btnUninstall.Text = "Uninstall";
                btnCancel.Text = "Cancel";
                btnInstall.Text = "Install Now";
            }
            else if (lang == "ja")
            {
                lblTitle.Text = "Loadin セットアップ";
                lblSubtitle.Text = "高速マルチセグメント・ダウンロードマネージャー (v2.1)";
                lblDirCaption.Text = "インストール先フォルダ";
                btnBrowse.Text = "参照...";
                chkDesktopShortcut.Text = "デスクトップにショートカットを作成";
                chkStartMenuShortcut.Text = "スタートメニューに Loadin を追加";
                chkStartupRun.Text = "Windows 起動時に自動実行 (バックグラウンド待機)";
                chkRegisterProtocol.Text = "ブラウザ連携プロトコル (loadin://) と拡張機能を展開";
                chkLaunchAfter.Text = "インストール完了後に Loadin を起動";
                lblStatus.Text = "インストールの準備ができました。[インストール] をクリックしてください。";
                btnUninstall.Text = "アンインストール";
                btnCancel.Text = "キャンセル";
                btnInstall.Text = "インストール";
            }
            else
            {
                lblTitle.Text = "Loadin 설치 마법사";
                lblSubtitle.Text = "고속 멀티 세그먼트 다운로드 매니저 (v2.1 Commercial Edition)";
                lblDirCaption.Text = "설치 폴더 위치";
                btnBrowse.Text = "찾아보기...";
                chkDesktopShortcut.Text = "바탕화면에 Loadin 바로가기 만들기";
                chkStartMenuShortcut.Text = "시작 메뉴 프로그램 목록에 Loadin 추가";
                chkStartupRun.Text = "윈도우 시작 시 자동 실행 (백그라운드 브라우저 연동 대기)";
                chkRegisterProtocol.Text = "웹 브라우저 연동 프로토콜(loadin://) 및 확장프로그램 설치";
                chkLaunchAfter.Text = "설치 완료 후 Loadin 바로 실행";
                lblStatus.Text = "설치 준비가 완료되었습니다. [설치 시작]을 클릭하세요.";
                btnUninstall.Text = "프로그램 제거";
                btnCancel.Text = "취소";
                btnInstall.Text = "설치 시작";
            }
        }

        private void OnBrowseFolder(object sender, EventArgs e)
        {
            using (FolderBrowserDialog fbd = new FolderBrowserDialog())
            {
                fbd.SelectedPath = txtInstallDir.Text;
                if (fbd.ShowDialog(this) == DialogResult.OK && !string.IsNullOrWhiteSpace(fbd.SelectedPath))
                {
                    txtInstallDir.Text = fbd.SelectedPath;
                }
            }
        }

        private void OnRunInstall(object sender, EventArgs e)
        {
            string targetDir = txtInstallDir.Text.Trim();
            if (string.IsNullOrEmpty(targetDir))
            {
                return;
            }

            btnInstall.Enabled = false;
            btnBrowse.Enabled = false;

            try
            {
                progressBar.Value = 10;
                lblStatus.Text = lang == "en" ? "Stopping running Loadin instances..." : "실행 중인 Loadin 본 프로그램 종료 중...";
                Application.DoEvents();

                TerminateRunningLoadin();

                Directory.CreateDirectory(targetDir);
                string extDir = Path.Combine(targetDir, "browser-extension");
                Directory.CreateDirectory(extDir);

                progressBar.Value = 30;
                lblStatus.Text = lang == "en" ? "Extracting Loadin.exe & core components..." : "Loadin 실행 파일 및 구성 요소 설치 중...";
                Application.DoEvents();

                string exePath = Path.Combine(targetDir, "Loadin.exe");
                string icoPath = Path.Combine(targetDir, "Loadin.ico");
                string uninstPath = Path.Combine(targetDir, "Uninstall.exe");

                ExtractResourceToFile("Loadin.exe", exePath);
                ExtractResourceToFile("Loadin.ico", icoPath);

                // Copy current installer binary as Uninstall.exe
                string selfPath = Application.ExecutablePath;
                if (!string.Equals(Path.GetFullPath(selfPath), Path.GetFullPath(uninstPath), StringComparison.OrdinalIgnoreCase))
                {
                    File.Copy(selfPath, uninstPath, true);
                }

                progressBar.Value = 55;
                lblStatus.Text = lang == "en" ? "Installing browser extension files..." : "브라우저 확장프로그램 파일 설치 중...";
                Application.DoEvents();

                ExtractResourceToFile("manifest.json", Path.Combine(extDir, "manifest.json"));
                ExtractResourceToFile("background.js", Path.Combine(extDir, "background.js"));
                ExtractResourceToFile("content.js", Path.Combine(extDir, "content.js"));
                ExtractResourceToFile("popup.html", Path.Combine(extDir, "popup.html"));
                ExtractResourceToFile("popup.js", Path.Combine(extDir, "popup.js"));

                progressBar.Value = 75;
                lblStatus.Text = lang == "en" ? "Configuring shortcuts and Windows Registry..." : "바로가기 및 윈도우 레지스트리 등록 중...";
                Application.DoEvents();

                if (chkDesktopShortcut.Checked)
                {
                    string desktopDir = Environment.GetFolderPath(Environment.SpecialFolder.DesktopDirectory);
                    CreateShortcut(Path.Combine(desktopDir, "Loadin.lnk"), exePath, targetDir, icoPath, "Loadin Download Manager");
                }

                if (chkStartMenuShortcut.Checked)
                {
                    string programsDir = Environment.GetFolderPath(Environment.SpecialFolder.Programs);
                    string startFolder = Path.Combine(programsDir, "Loadin");
                    Directory.CreateDirectory(startFolder);
                    CreateShortcut(Path.Combine(startFolder, "Loadin.lnk"), exePath, targetDir, icoPath, "Loadin Download Manager");
                    CreateShortcut(Path.Combine(startFolder, "Uninstall Loadin.lnk"), uninstPath, targetDir, icoPath, "Uninstall Loadin", "/uninstall");
                }

                // Windows Startup Auto-Run Registry
                using (RegistryKey runKey = Registry.CurrentUser.CreateSubKey(@"Software\Microsoft\Windows\CurrentVersion\Run"))
                {
                    if (runKey != null)
                    {
                        if (chkStartupRun.Checked)
                        {
                            runKey.SetValue("Loadin", "\"" + exePath + "\" -headless", RegistryValueKind.String);
                        }
                        else
                        {
                            runKey.DeleteValue("Loadin", false);
                        }
                    }
                }

                // Protocol Registration (loadin://)
                if (chkRegisterProtocol.Checked)
                {
                    using (RegistryKey protoKey = Registry.CurrentUser.CreateSubKey(@"Software\Classes\loadin"))
                    {
                        if (protoKey != null)
                        {
                            protoKey.SetValue("", "URL:Loadin Download Protocol");
                            protoKey.SetValue("URL Protocol", "");
                            using (RegistryKey iconKey = protoKey.CreateSubKey("DefaultIcon"))
                            {
                                iconKey.SetValue("", "\"" + exePath + "\",0");
                            }
                            using (RegistryKey cmdKey = protoKey.CreateSubKey(@"shell\open\command"))
                            {
                                cmdKey.SetValue("", "\"" + exePath + "\" -url \"%1\"");
                            }
                        }
                    }
                }

                // Windows "Installed Apps" Uninstaller Registration
                using (RegistryKey uninstKey = Registry.CurrentUser.CreateSubKey(@"Software\Microsoft\Windows\CurrentVersion\Uninstall\Loadin"))
                {
                    if (uninstKey != null)
                    {
                        uninstKey.SetValue("DisplayName", "Loadin Download Manager");
                        uninstKey.SetValue("DisplayVersion", "2.1.0");
                        uninstKey.SetValue("Publisher", "Loadin Project");
                        uninstKey.SetValue("DisplayIcon", exePath);
                        uninstKey.SetValue("InstallLocation", targetDir);
                        uninstKey.SetValue("UninstallString", "\"" + uninstPath + "\" /uninstall");
                        uninstKey.SetValue("NoModify", 1, RegistryValueKind.DWord);
                        uninstKey.SetValue("NoRepair", 1, RegistryValueKind.DWord);
                    }
                }

                // Persist initial settings with runOnStartup and language
                string appDataDir = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData), "Loadin");
                Directory.CreateDirectory(appDataDir);

                progressBar.Value = 100;
                lblStatus.ForeColor = Color.FromArgb(48, 209, 88);
                lblStatus.Text = lang == "en"
                    ? "Installation complete!"
                    : (lang == "ja" ? "インストールが完了しました！" : "Loadin 설치가 완료되었습니다!");
                Application.DoEvents();

                if (chkLaunchAfter.Checked && File.Exists(exePath))
                {
                    Process.Start(new ProcessStartInfo
                    {
                        FileName = exePath,
                        WorkingDirectory = targetDir,
                        UseShellExecute = true
                    });
                }

                MessageBox.Show(
                    this,
                    lang == "en"
                        ? "Loadin has been installed successfully.\n\nPath: " + targetDir
                        : (lang == "ja"
                            ? "Loadin のインストールが正常に完了しました。\n\nインストール先: " + targetDir
                            : "Loadin 설치가 성공적으로 완료되었습니다.\n\n설치 경로: " + targetDir),
                    "Loadin Setup",
                    MessageBoxButtons.OK,
                    MessageBoxIcon.Information
                );
                Close();
            }
            catch (Exception ex)
            {
                btnInstall.Enabled = true;
                btnBrowse.Enabled = true;
                MessageBox.Show(this, "설치 중 오류가 발생했습니다:\n" + ex.Message, "Loadin Setup Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        private static void ExtractResourceToFile(string resourceName, string outPath)
        {
            Assembly asm = Assembly.GetExecutingAssembly();
            using (Stream s = asm.GetManifestResourceStream(resourceName))
            {
                if (s == null)
                {
                    throw new FileNotFoundException("Embedded resource not found: " + resourceName);
                }
                using (FileStream fs = new FileStream(outPath, FileMode.Create, FileAccess.Write))
                {
                    s.CopyTo(fs);
                }
            }
        }

        private static Icon ExtractEmbeddedIcon()
        {
            Assembly asm = Assembly.GetExecutingAssembly();
            using (Stream s = asm.GetManifestResourceStream("Loadin.ico"))
            {
                if (s != null) return new Icon(s);
            }
            return null;
        }

        private static void CreateShortcut(string lnkPath, string targetPath, string workDir, string iconPath, string desc, string args = "")
        {
            try
            {
                Type shellType = Type.GetTypeFromProgID("WScript.Shell");
                if (shellType == null) return;
                dynamic shell = Activator.CreateInstance(shellType);
                dynamic shortcut = shell.CreateShortcut(lnkPath);
                shortcut.TargetPath = targetPath;
                shortcut.WorkingDirectory = workDir;
                if (!string.IsNullOrEmpty(args)) shortcut.Arguments = args;
                if (!string.IsNullOrEmpty(iconPath)) shortcut.IconLocation = iconPath;
                shortcut.Description = desc;
                shortcut.Save();
            }
            catch { }
        }

        protected override void OnHandleCreated(EventArgs e)
        {
            base.OnHandleCreated(e);
            try
            {
                int dark = 1;
                DwmSetWindowAttribute(Handle, 20, ref dark, sizeof(int));
                int roundPref = 2; // DWMWCP_ROUND
                DwmSetWindowAttribute(Handle, 33, ref roundPref, sizeof(int));
            }
            catch { }
        }

        protected override void OnPaint(PaintEventArgs e)
        {
            base.OnPaint(e);
            Graphics g = e.Graphics;
            g.SmoothingMode = SmoothingMode.AntiAlias;

            // Draw Loadin squircle logo
            Rectangle logoRect = new Rectangle(22, 20, 32, 32);
            using (GraphicsPath path = RoundedRect(logoRect, 8))
            using (LinearGradientBrush b = new LinearGradientBrush(logoRect, Color.FromArgb(48, 152, 255), Color.FromArgb(0, 102, 224), LinearGradientMode.Vertical))
            {
                g.FillPath(b, path);
            }
            using (Pen whitePen = new Pen(Color.White, 2.2f))
            {
                whitePen.StartCap = LineCap.Round;
                whitePen.EndCap = LineCap.Round;
                whitePen.LineJoin = LineJoin.Round;
                g.DrawLine(whitePen, 38, 27, 38, 39);
                g.DrawLine(whitePen, 33, 35, 38, 39);
                g.DrawLine(whitePen, 43, 35, 38, 39);
                g.DrawLine(whitePen, 31, 44, 45, 44);
            }

            // Header divider
            using (Pen divPen = new Pen(Color.FromArgb(36, 36, 40), 1))
            {
                g.DrawLine(divPen, 0, 74, Width, 74);
            }

            // Outer border
            using (Pen borderPen = new Pen(Color.FromArgb(52, 52, 58), 1))
            {
                g.DrawRectangle(borderPen, 0, 0, Width - 1, Height - 1);
            }
        }

        private static GraphicsPath RoundedRect(Rectangle bounds, int radius)
        {
            int d = radius * 2;
            GraphicsPath path = new GraphicsPath();
            path.AddArc(bounds.X, bounds.Y, d, d, 180, 90);
            path.AddArc(bounds.Right - d, bounds.Y, d, d, 270, 90);
            path.AddArc(bounds.Right - d, bounds.Bottom - d, d, d, 0, 90);
            path.AddArc(bounds.X, bounds.Bottom - d, d, d, 90, 90);
            path.CloseFigure();
            return path;
        }

        private void DragWindow(object sender, MouseEventArgs e)
        {
            if (e.Button == MouseButtons.Left)
            {
                ReleaseCapture();
                SendMessage(Handle, 0xA1, 0x2, 0);
            }
        }

        private static Font CreateFont(float size, FontStyle style)
        {
            string[] candidates = { "Pretendard", "Malgun Gothic", "Segoe UI" };
            foreach (string name in candidates)
            {
                try
                {
                    Font f = new Font(name, size, style);
                    if (f.Name.IndexOf(name, StringComparison.OrdinalIgnoreCase) >= 0) return f;
                }
                catch { }
            }
            return new Font(SystemFonts.MessageBoxFont.FontFamily, size, style);
        }

        public static void TerminateRunningLoadin()
        {
            // 1. Gracefully request the running Loadin engine to pause tasks, close its desktop window, and exit
            try
            {
                HttpWebRequest req = (HttpWebRequest)WebRequest.Create("http://127.0.0.1:17890/api/app/quit");
                req.Method = "POST";
                req.Timeout = 700;
                req.ReadWriteTimeout = 700;
                req.ContentLength = 0;
                using (HttpWebResponse resp = (HttpWebResponse)req.GetResponse()) { }
            }
            catch { }

            // 2. Close any visible Loadin main desktop window or popup dialog via Win32 WM_CLOSE
            try
            {
                EnumWindows((hWnd, lParam) =>
                {
                    if (!IsWindowVisible(hWnd)) return true;
                    StringBuilder sb = new StringBuilder(256);
                    if (GetWindowText(hWnd, sb, sb.Capacity) > 0)
                    {
                        string title = sb.ToString();
                        if (title.IndexOf("Loadin", StringComparison.OrdinalIgnoreCase) >= 0 &&
                            title.IndexOf("Setup", StringComparison.OrdinalIgnoreCase) < 0 &&
                            title.IndexOf("Uninstall", StringComparison.OrdinalIgnoreCase) < 0 &&
                            title.IndexOf("제거", StringComparison.OrdinalIgnoreCase) < 0 &&
                            title.IndexOf("설치", StringComparison.OrdinalIgnoreCase) < 0)
                        {
                            PostMessage(hWnd, WM_CLOSE, IntPtr.Zero, IntPtr.Zero);
                        }
                    }
                    return true;
                }, IntPtr.Zero);
            }
            catch { }

            // 3. Force-terminate any remaining Loadin engine and helper dialog processes
            int currentPid = Process.GetCurrentProcess().Id;
            string[] targetNames = { "Loadin", "Loadin_AddDialog", "AetherHover", "AetherDL" };
            foreach (string procName in targetNames)
            {
                try
                {
                    foreach (Process proc in Process.GetProcessesByName(procName))
                    {
                        if (proc.Id == currentPid) continue;
                        try
                        {
                            proc.Kill();
                            proc.WaitForExit(1500);
                        }
                        catch { }
                    }
                }
                catch { }
            }

            try
            {
                ProcessStartInfo psi = new ProcessStartInfo
                {
                    FileName = "taskkill.exe",
                    Arguments = "/F /IM Loadin.exe /IM Loadin_AddDialog.exe /T",
                    CreateNoWindow = true,
                    UseShellExecute = false,
                    WindowStyle = ProcessWindowStyle.Hidden
                };
                using (Process p = Process.Start(psi))
                {
                    if (p != null) p.WaitForExit(1200);
                }
            }
            catch { }

            Thread.Sleep(250);
        }

        public static void RunUninstall()
        {
            // Immediately terminate running Loadin main program & windows when uninstaller is launched
            TerminateRunningLoadin();

            DialogResult dr = MessageBox.Show(
                "실행 중인 Loadin 프로그램을 종료했습니다.\n\nLoadin 다운로드 매니저를 시스템에서 완전히 제거하시겠습니까?\n(기존에 다운로드한 파일들은 삭제되지 않습니다.)",
                "Loadin 제거 (Uninstall)",
                MessageBoxButtons.YesNo,
                MessageBoxIcon.Question
            );
            if (dr != DialogResult.Yes) return;

            try
            {
                TerminateRunningLoadin();

                string installLocation = null;
                using (RegistryKey uninstKey = Registry.CurrentUser.OpenSubKey(@"Software\Microsoft\Windows\CurrentVersion\Uninstall\Loadin"))
                {
                    if (uninstKey != null)
                    {
                        installLocation = uninstKey.GetValue("InstallLocation") as string;
                    }
                }
                if (string.IsNullOrEmpty(installLocation))
                {
                    installLocation = Path.Combine(
                        Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
                        "Programs",
                        "Loadin"
                    );
                }

                using (RegistryKey runKey = Registry.CurrentUser.CreateSubKey(@"Software\Microsoft\Windows\CurrentVersion\Run"))
                {
                    if (runKey != null) runKey.DeleteValue("Loadin", false);
                }
                Registry.CurrentUser.DeleteSubKeyTree(@"Software\Classes\loadin", false);
                Registry.CurrentUser.DeleteSubKeyTree(@"Software\Microsoft\Windows\CurrentVersion\Uninstall\Loadin", false);

                string desktopLnk = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.DesktopDirectory), "Loadin.lnk");
                if (File.Exists(desktopLnk)) File.Delete(desktopLnk);

                string startFolder = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.Programs), "Loadin");
                if (Directory.Exists(startFolder)) Directory.Delete(startFolder, true);

                // Delete installed files in installLocation
                if (!string.IsNullOrEmpty(installLocation) && Directory.Exists(installLocation))
                {
                    string selfPath = Path.GetFullPath(Application.ExecutablePath);
                    string fullInstallDir = Path.GetFullPath(installLocation);

                    // Delete Loadin.exe, Loadin.ico, browser-extension right away
                    try
                    {
                        string mainExe = Path.Combine(fullInstallDir, "Loadin.exe");
                        if (File.Exists(mainExe)) File.Delete(mainExe);
                        string icoFile = Path.Combine(fullInstallDir, "Loadin.ico");
                        if (File.Exists(icoFile)) File.Delete(icoFile);
                        string extDir = Path.Combine(fullInstallDir, "browser-extension");
                        if (Directory.Exists(extDir)) Directory.Delete(extDir, true);
                    }
                    catch { }

                    // Schedule background cleanup of the install directory after Uninstall.exe exits
                    try
                    {
                        Process.Start(new ProcessStartInfo
                        {
                            FileName = "cmd.exe",
                            Arguments = "/C ping 127.0.0.1 -n 2 > nul & rmdir /S /Q \"" + fullInstallDir + "\"",
                            CreateNoWindow = true,
                            UseShellExecute = false,
                            WindowStyle = ProcessWindowStyle.Hidden
                        });
                    }
                    catch { }
                }

                MessageBox.Show("Loadin 본 프로그램 종료 및 제거가 완료되었습니다.", "Loadin Uninstall", MessageBoxButtons.OK, MessageBoxIcon.Information);
            }
            catch (Exception ex)
            {
                MessageBox.Show("제거 중 오류 발생: " + ex.Message, "Loadin Uninstall", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        [STAThread]
        public static void Main(string[] args)
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);

            string exeName = Path.GetFileNameWithoutExtension(Application.ExecutablePath) ?? "";
            if (exeName.IndexOf("Uninstall", StringComparison.OrdinalIgnoreCase) >= 0)
            {
                RunUninstall();
                return;
            }

            foreach (string a in args)
            {
                if (string.Equals(a, "/uninstall", StringComparison.OrdinalIgnoreCase) ||
                    string.Equals(a, "--uninstall", StringComparison.OrdinalIgnoreCase))
                {
                    RunUninstall();
                    return;
                }
            }

            Application.Run(new SetupForm());
        }
    }
}
