# ForgeUI Assets Package

> Production-ready static asset management with fingerprinting, caching, and embedded filesystem support

---

## Features

- **Content-Based Fingerprinting**: Automatic cache busting using SHA256 hashes
- **Intelligent Caching**: Immutable cache for fingerprinted assets (1 year), moderate cache for dev
- **Security**: Path traversal protection, proper MIME types, optional SRI
- **Embedded FS Support**: Single-binary deployments with `embed.FS`
- **Manifest System**: Pre-computed hashes for production builds
- **Dev/Prod Modes**: Zero-config development, optimized production
- **Thread-Safe**: Concurrent asset serving with RWMutex protection

---

## Quick Start

### Basic Usage

```go
package main

import (
    "net/http"
    "github.com/xraph/forgeui"
)

func main() {
    // Initialize ForgeUI app
    app := forgeui.New(
        forgeui.WithDebug(true),              // Dev mode
        forgeui.WithAssetPublicDir("public"),  // Assets directory
    )

    // Serve static files through asset pipeline
    http.Handle("/static/", app.Assets.Handler())

    http.ListenAndServe(":8080", nil)
}
```

### In HTML Templates

```go
// CSS Stylesheets
app.Assets.StyleSheet("css/app.css")
app.Assets.StyleSheet("css/print.css", assets.WithMedia("print"))

// JavaScript
app.Assets.Script("js/app.js", assets.WithDefer())
app.Assets.Script("js/module.js", assets.WithModule())

// Inline CSS/JS
assets.InlineCSS("body { margin: 0; }")
assets.InlineScript("console.log('hello');")
```

---

## Development vs Production

### Development Mode

```go
app := forgeui.New(forgeui.WithDebug(true))
```

- **No fingerprinting**: URLs remain simple (`/static/app.css`)
- **Moderate caching**: `Cache-Control: public, max-age=3600`
- **Fast iteration**: No build step required

### Production Mode

```go
app := forgeui.New(forgeui.WithDebug(false))
```

- **Automatic fingerprinting**: URLs include content hash (`/static/app.abc12345.css`)
- **Immutable caching**: `Cache-Control: public, max-age=31536000, immutable`
- **Optimal performance**: 1-year browser cache

---

## Manifest System

For production builds, pre-compute fingerprints to eliminate runtime I/O:

### Generate Manifest

```go
package main

import (
    "github.com/xraph/forgeui/assets"
)

func main() {
    m := assets.NewManager(assets.Config{
        PublicDir: "public",
        IsDev:     false,
    })

    // Generate manifest for all assets
    manifest, err := assets.GenerateManifest(m)
    if err != nil {
        panic(err)
    }

    // Save to file
    if err := manifest.Save("dist/manifest.json"); err != nil {
        panic(err)
    }
}
```

### Use Manifest

```go
app := forgeui.New(
    forgeui.WithAssetManifest("dist/manifest.json"),
)
```

**Manifest Format** (`manifest.json`):
```json
{
  "css/app.css": "css/app.abc12345.css",
  "js/main.js": "js/main.def67890.js"
}
```

---

## Embedded Filesystem

For single-binary deployments:

```go
package main

import (
    "embed"
    "net/http"
    "github.com/xraph/forgeui/assets"
)

//go:embed static/*
var staticFS embed.FS

func main() {
    m := assets.NewEmbeddedManager(staticFS, assets.Config{
        PublicDir: "static",
        IsDev:     false,
    })

    // Pre-generate fingerprints
    m.FingerprintAllEmbedded()

    // Serve from embedded FS
    http.Handle("/static/", m.EmbeddedHandler())
    
    http.ListenAndServe(":8080", nil)
}
```

---

## API Reference

### Manager

```go
type Manager struct {
    // ...
}

func NewManager(cfg Config) *Manager

// URL returns the URL for an asset (with fingerprint in production)
func (m *Manager) URL(path string) string

// Handler returns an http.Handler for serving static files
func (m *Manager) Handler() http.Handler

// StyleSheet creates a <link> element for CSS
func (m *Manager) StyleSheet(path string, opts ...StyleOption) g.Node

// Script creates a <script> element for JavaScript
func (m *Manager) Script(path string, opts ...ScriptOption) g.Node

// FingerprintAll generates fingerprints for all assets
func (m *Manager) FingerprintAll() error

// SaveManifest writes fingerprint mappings to a file
func (m *Manager) SaveManifest(path string) error
```

### Configuration

```go
type Config struct {
    PublicDir  string // Source directory (default: "public")
    OutputDir  string // Output directory (default: "dist")
    IsDev      bool   // Development mode
    Manifest   string // Manifest file path
}
```

### Style Options

```go
assets.WithMedia("print")              // media attribute
assets.WithPreload()                   // preload the stylesheet
assets.WithIntegrity("sha256-...")     // SRI hash
assets.WithCrossOrigin("anonymous")    // CORS
```

### Script Options

```go
assets.WithDefer()                     // defer attribute
assets.WithAsync()                     // async attribute
assets.WithModule()                    // type="module"
assets.WithScriptIntegrity("sha256-...") // SRI hash
assets.WithScriptCrossOrigin("anonymous") // CORS
assets.WithNoModule()                  // nomodule fallback
```

---

## Examples

### Preloading Assets

```go
// Preload critical CSS
app.Assets.PreloadStyleSheet("css/critical.css")
app.Assets.StyleSheet("css/critical.css")

// Preload critical JS
app.Assets.PreloadScript("js/critical.js")
app.Assets.Script("js/critical.js", assets.WithDefer())
```

### Media Queries

```go
// Screen styles
app.Assets.StyleSheet("css/screen.css", assets.WithMedia("screen"))

// Print styles
app.Assets.StyleSheet("css/print.css", assets.WithMedia("print"))
```

### ES Modules

```go
// Modern browsers
app.Assets.Script("js/app.js", assets.WithModule())

// Legacy fallback
app.Assets.Script("js/app.legacy.js", assets.WithNoModule())
```

### Subresource Integrity (SRI)

```go
app.Assets.StyleSheet(
    "css/app.css",
    assets.WithIntegrity("sha256-abc123..."),
    assets.WithCrossOrigin("anonymous"),
)
```

---

## Security

### Path Traversal Protection

All paths are validated to prevent directory traversal attacks:

```go
// ✅ Valid
app.Assets.URL("css/app.css")
app.Assets.URL("js/vendor/lib.js")

// ❌ Blocked
app.Assets.URL("../../../etc/passwd")
app.Assets.URL("../config.json")
```

### MIME Type Detection

Proper `Content-Type` headers are set automatically based on file extension to prevent XSS attacks.

### Optional SRI

Subresource Integrity (SRI) hashes can be added for CDN fallbacks:

```go
app.Assets.StyleSheet(
    "css/app.css",
    assets.WithIntegrity("sha256-..."),
)
```

---

## Performance

### Caching Strategy

| Asset Type | Cache-Control | Duration |
|------------|---------------|----------|
| Fingerprinted (prod) | `public, max-age=31536000, immutable` | 1 year |
| Non-fingerprinted | `public, max-age=3600` | 1 hour |
| Dev mode | `public, max-age=3600` | 1 hour |

### Thread Safety

The asset manager uses `sync.RWMutex` for thread-safe concurrent access:
- Multiple concurrent reads
- Exclusive writes
- Zero-allocation fingerprint lookups

### Manifest Pre-loading

Loading a manifest eliminates runtime I/O in production:
- No disk reads during request handling
- Instant URL generation
- Reduced latency

---

## Testing

Run tests with coverage:

```bash
go test ./assets/... -cover
```

Expected coverage: **>85%**

---

## License

MIT License - see LICENSE file for details

