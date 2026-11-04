<div align="center">

# 🗜️ JustCompress

**Fast, beautiful image compression — right in your browser or powered by a Go backend.**

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![JavaScript](https://img.shields.io/badge/JavaScript-ES6+-F7DF1E?logo=javascript&logoColor=black)](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.x-38B2AC?logo=tailwindcss&logoColor=white)](https://tailwindcss.com)

[Live Demo](https://anomg0x.github.io/JustCompress) · [Report Bug](https://github.com/anomg0x/JustCompress/issues) · [Request Feature](https://github.com/anomg0x/JustCompress/issues)

</div>

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🖼️ **Dual Engine** | Client-side (browser) compression **or** server-side Go backend — switch with one toggle |
| ⚡ **Instant Preview** | Real-time side-by-side comparison of original vs compressed image |
| 🎚️ **Quality Control** | Adjustable quality slider (10–100%) for fine-grained compression |
| 🔒 **Lossy & Lossless** | Choose between lossy (smaller files) or lossless (perfect quality) modes |
| 📐 **Smart Resize** | Optional max-width/max-height resizing with preserved aspect ratio (Go backend) |
| 🌙 **Dark Theme** | Beautiful dark UI built with Tailwind CSS |
| 📦 **Drag & Drop** | Simply drag your images into the browser — no clicks needed |
| 💾 **One-Click Download** | Download compressed images instantly |
| 🖥️ **CLI Tool** | Compress images from the command line with the Go CLI |

---

## 📸 Supported Formats

| Format | Read | Write | Notes |
|--------|:----:|:-----:|-------|
| **JPEG** | ✅ | ✅ | Full lossy compression with quality control |
| **PNG** | ✅ | ✅ | Lossless with configurable compression level |
| **WebP** | ✅ | — | Decoded and re-encoded as JPEG/PNG |
| **GIF** | ✅ | — | First frame extracted and compressed |

---

## 🚀 Quick Start

### Option 1 — Browser Only (No Install)

Just open [`index.html`](index.html) in your browser or visit the [live demo](https://anomg0x.github.io/JustCompress). Works offline, no server needed.

### Option 2 — With Go Backend

```bash
# Clone the repo
git clone https://github.com/anomg0x/JustCompress.git
cd JustCompress/server

# Install dependencies & build
go mod download
go build -o justcompress-server .

# Run the server
./justcompress-server
```

Open **http://localhost:8080** and toggle the **Compression Engine** switch to "Go Server (Backend)".

### Option 3 — CLI Tool

```bash
cd JustCompress/server
go build -o justcompress-cli ./cmd/cli/

# Compress a JPEG with 60% quality
./justcompress-cli -input photo.jpg -quality 60

# Convert PNG to JPEG at 80% quality
./justcompress-cli -input screenshot.png -format jpeg -quality 80

# Resize and compress
./justcompress-cli -input photo.jpg -max-width 1920 -max-height 1080 -quality 75
```

---

## 🛠️ CLI Reference

```
Usage:
  justcompress-cli [flags]

Flags:
  -input        string   Input image path (required)
  -output       string   Output image path (auto-generated if empty)
  -quality      int      Compression quality, 1-100 (default: 75)
  -format       string   Output format: jpeg, png (auto-detect if empty)
  -lossless     bool     Use lossless compression (default: false)
  -max-width    uint     Max width in pixels, 0 = no resize (default: 0)
  -max-height   uint     Max height in pixels, 0 = no resize (default: 0)
```

**Examples:**

```bash
# Basic compression
justcompress-cli -input vacation.jpg -quality 50

# Lossless PNG optimization
justcompress-cli -input diagram.png -lossless

# Custom output path
justcompress-cli -input photo.jpg -output ./out/photo_small.jpg -quality 60

# Resize large photo for web
justcompress-cli -input camera_raw.jpg -max-width 1280 -quality 70
```

---

## 🏗️ Project Structure

```
JustCompress/
├── index.html              # Frontend — main HTML page
├── script.js               # Frontend — compression logic & UI
├── style.css               # Frontend — custom styles
├── server/                 # Go backend
│   ├── main.go             # HTTP server entry point
│   ├── handlers/
│   │   └── handlers.go     # API route handlers
│   ├── compressor/
│   │   └── compressor.go   # Image compression engine
│   ├── cmd/cli/
│   │   └── main.go         # CLI tool entry point
│   ├── go.mod
│   └── go.sum
├── LICENSE
└── README.md
```

---

## 🔌 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/compress` | Compress an uploaded image |
| `GET` | `/api/health` | Server health check |

### `POST /api/compress`

**Form Data:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `image` | file | — | Image file (required) |
| `quality` | int | `75` | Compression quality (1–100) |
| `format` | string | `jpeg` | Output format: `jpeg`, `png` |
| `lossless` | string | `false` | Lossless mode: `true` / `false` |
| `metadata` | string | `false` | Return JSON metadata instead of binary |

**Response** — compressed image binary with headers:

```
X-Original-Size: 2458624
X-Compressed-Size: 184320
X-Savings: 92.5%
X-Duration: 45.2ms
```

---

## 🧰 Tech Stack

- **Frontend:** HTML5, JavaScript (ES6+), [Tailwind CSS](https://tailwindcss.com)
- **Backend:** [Go](https://go.dev) 1.27+
- **Libraries:**
  - [`golang.org/x/image`](https://pkg.go.dev/golang.org/x/image) — WebP decoding
  - [`github.com/nfnt/resize`](https://github.com/nfnt/resize) — High-quality image resizing (Lanczos3)
- **No external dependencies** on the frontend — works offline

---

## 🤝 Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

Distributed under the **MIT License**. See [`LICENSE`](LICENSE) for details.

---

<div align="center">

**Built with ❤️ by [AnomG](https://github.com/anomg0x)**

⭐ Star this repo if you find it useful!

</div>
