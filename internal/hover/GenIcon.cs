using System;
using System.Collections.Generic;
using System.Drawing;
using System.Drawing.Drawing2D;
using System.Drawing.Imaging;
using System.IO;

namespace LoadinIconGen
{
    public static class Program
    {
        public static int Main(string[] args)
        {
            if (args.Length < 1)
            {
                Console.Error.WriteLine("Usage: GenIcon.exe <output1.ico> [output2.ico ...]");
                return 1;
            }

            int[] sizes = new int[] { 256, 64, 48, 32, 16 };
            byte[] icoBytes = BuildMultiResolutionIco(sizes);

            foreach (string outPath in args)
            {
                string dir = Path.GetDirectoryName(Path.GetFullPath(outPath));
                if (!string.IsNullOrEmpty(dir) && !Directory.Exists(dir))
                {
                    Directory.CreateDirectory(dir);
                }
                File.WriteAllBytes(outPath, icoBytes);
            }
            return 0;
        }

        private static byte[] BuildMultiResolutionIco(int[] sizes)
        {
            List<byte[]> pngPayloads = new List<byte[]>();
            foreach (int sz in sizes)
            {
                using (Bitmap bmp = RenderLogoBitmap(sz))
                using (MemoryStream ms = new MemoryStream())
                {
                    bmp.Save(ms, ImageFormat.Png);
                    pngPayloads.Add(ms.ToArray());
                }
            }

            using (MemoryStream icoStream = new MemoryStream())
            using (BinaryWriter bw = new BinaryWriter(icoStream))
            {
                // ICONDIR header (6 bytes)
                bw.Write((ushort)0);               // Reserved
                bw.Write((ushort)1);               // Type: 1 = ICO
                bw.Write((ushort)sizes.Length);    // Image count

                int dataOffset = 6 + (16 * sizes.Length);
                for (int i = 0; i < sizes.Length; i++)
                {
                    int sz = sizes[i];
                    byte[] payload = pngPayloads[i];
                    bw.Write((byte)(sz >= 256 ? 0 : sz)); // Width
                    bw.Write((byte)(sz >= 256 ? 0 : sz)); // Height
                    bw.Write((byte)0);                    // Color count
                    bw.Write((byte)0);                    // Reserved
                    bw.Write((ushort)1);                  // Color planes
                    bw.Write((ushort)32);                 // Bits per pixel
                    bw.Write((uint)payload.Length);       // Image byte size
                    bw.Write((uint)dataOffset);           // Offset
                    dataOffset += payload.Length;
                }

                for (int i = 0; i < pngPayloads.Count; i++)
                {
                    bw.Write(pngPayloads[i]);
                }

                bw.Flush();
                return icoStream.ToArray();
            }
        }

        private static Bitmap RenderLogoBitmap(int size)
        {
            Bitmap bmp = new Bitmap(size, size, PixelFormat.Format32bppArgb);
            using (Graphics g = Graphics.FromImage(bmp))
            {
                g.SmoothingMode = SmoothingMode.AntiAlias;
                g.InterpolationMode = InterpolationMode.HighQualityBicubic;
                g.PixelOffsetMode = PixelOffsetMode.HighQuality;
                g.Clear(Color.Transparent);

                float pad = size * 0.04f;
                float boxSize = size - (pad * 2f);
                float radius = size * 0.24f;
                RectangleF rect = new RectangleF(pad, pad, boxSize, boxSize);

                using (GraphicsPath bgPath = RoundedRect(rect, radius))
                {
                    using (LinearGradientBrush bgBrush = new LinearGradientBrush(
                        rect,
                        Color.FromArgb(255, 48, 152, 255),
                        Color.FromArgb(255, 0, 102, 224),
                        LinearGradientMode.Vertical))
                    {
                        g.FillPath(bgBrush, bgPath);
                    }

                    // Subtle top inner highlight
                    using (Pen borderPen = new Pen(Color.FromArgb(95, 255, 255, 255), Math.Max(1f, size * 0.025f)))
                    {
                        g.DrawPath(borderPen, bgPath);
                    }
                }

                // Draw clean white download arrow & tray
                float stroke = Math.Max(1.6f, size * 0.095f);
                using (Pen whitePen = new Pen(Color.White, stroke))
                {
                    whitePen.StartCap = LineCap.Round;
                    whitePen.EndCap = LineCap.Round;
                    whitePen.LineJoin = LineJoin.Round;

                    float cx = size * 0.50f;
                    float arrowTop = size * 0.25f;
                    float arrowBottom = size * 0.60f;
                    float wingX = size * 0.16f;
                    float wingY = size * 0.46f;

                    // Vertical shaft
                    g.DrawLine(whitePen, cx, arrowTop, cx, arrowBottom);
                    // Arrow head wings
                    g.DrawLine(whitePen, cx - wingX, wingY, cx, arrowBottom);
                    g.DrawLine(whitePen, cx + wingX, wingY, cx, arrowBottom);

                    // Bottom dock bar
                    float barLeft = size * 0.28f;
                    float barRight = size * 0.72f;
                    float barY = size * 0.74f;
                    g.DrawLine(whitePen, barLeft, barY, barRight, barY);
                }
            }
            return bmp;
        }

        private static GraphicsPath RoundedRect(RectangleF bounds, float radius)
        {
            float d = radius * 2f;
            GraphicsPath path = new GraphicsPath();
            if (d <= 0f)
            {
                path.AddRectangle(bounds);
                return path;
            }
            path.AddArc(bounds.X, bounds.Y, d, d, 180, 90);
            path.AddArc(bounds.Right - d, bounds.Y, d, d, 270, 90);
            path.AddArc(bounds.Right - d, bounds.Bottom - d, d, d, 0, 90);
            path.AddArc(bounds.X, bounds.Bottom - d, d, d, 90, 90);
            path.CloseFigure();
            return path;
        }
    }
}
