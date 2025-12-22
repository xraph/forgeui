# Built-in Plugins

ForgeUI includes a collection of high-quality plugins that extend functionality with minimal setup.

## Available Plugins

### 🍞 Toast Notifications
**Package:** `github.com/xraph/forgeui/plugins/toast`

Beautiful notification toasts with Alpine.js store integration.

**Features:**
- Multiple variants (info, success, warning, error)
- Auto-dismiss with configurable timeout
- Queue management
- Position options

**Usage:**
```go
registry.Use(toast.New(
    toast.WithPosition("top-right"),
    toast.WithMaxVisible(5),
))
```

```go
// In your templates
toast.Container()

// Trigger from JavaScript
$store.toasts.success('Item saved!')
```

---

### 🎯 Sortable Lists
**Package:** `github.com/xraph/forgeui/plugins/sortable`

Drag-and-drop list reordering powered by SortableJS.

**Features:**
- Smooth animations
- Handle selectors
- Server synchronization
- Callback hooks

**Usage:**
```go
registry.Use(sortable.New())
```

```html
<ul x-sortable="{ animation: 150, onUpdate: handleUpdate }">
  <li data-id="1">Item 1</li>
  <li data-id="2">Item 2</li>
</ul>
```

---

### 📊 Charts
**Package:** `github.com/xraph/forgeui/plugins/charts`

Data visualization with Chart.js integration.

**Chart Types:**
- LineChart - Time series and trends
- BarChart - Categorical comparisons
- PieChart - Proportional data
- AreaChart - Filled line charts
- DoughnutChart - Hollow pie charts

**Usage:**
```go
registry.Use(charts.New())
```

```go
charts.LineChart(charts.LineChartData{
    Labels: []string{"Jan", "Feb", "Mar"},
    Datasets: []charts.DatasetConfig{
        {
            Label: "Sales",
            Data: []float64{12, 19, 3},
            BorderColor: "rgb(59, 130, 246)",
        },
    },
})
```

---

### 📈 Analytics
**Package:** `github.com/xraph/forgeui/plugins/analytics`

Event tracking integration for multiple providers.

**Supported Providers:**
- Google Analytics 4 (GA4)
- Plausible (privacy-friendly)
- Matomo (self-hosted)
- Custom implementations

**Usage:**
```go
registry.Use(analytics.New(analytics.Config{
    Provider: "plausible",
    Domain: "example.com",
    TrackPageViews: true,
}))
```

```javascript
// Track custom events
$store.analytics.trackEvent('button_click', { button: 'signup' })
```

---

### 🔍 SEO
**Package:** `github.com/xraph/forgeui/plugins/seo`

Meta tags, Open Graph, Twitter Cards, and structured data.

**Features:**
- Open Graph tags
- Twitter Card tags
- JSON-LD structured data
- Canonical URLs
- Robots meta tags
- Sitemap generation helpers

**Usage:**
```go
registry.Use(seo.New())
```

```go
seo.MetaTagsNode(seo.MetaTags{
    Title: "My Page",
    Description: "Page description",
    OGImage: "/og-image.jpg",
    TwitterCard: "summary_large_image",
})

seo.JSONLDNode(seo.ArticleSchema(
    "Article Title",
    "Description",
    "/image.jpg",
    "2024-12-20",
    "Author Name",
))
```

---

### ⚡ HTMX Plugin Wrapper
**Package:** `github.com/xraph/forgeui/plugins/htmxplugin`

Wraps HTMX as a plugin for dependency management.

**Usage:**
```go
registry.Use(htmxplugin.New(
    "sse",        // Server-Sent Events extension
    "ws",         // WebSocket extension
    "json-enc",   // JSON encoding extension
))
```

---

### 🏢 Corporate Theme
**Package:** `github.com/xraph/forgeui/plugins/themes/corporate`

Professional theme for business applications.

**Features:**
- Conservative blue/gray palette
- High contrast for accessibility (WCAG AAA)
- Print-friendly styles
- Professional typography
- Corporate components (letterhead, reports)

**Usage:**
```go
registry.Use(corporate.New())
```

---

## Quick Start

### Install All Plugins

```go
import (
    "github.com/xraph/forgeui/plugin"
    "github.com/xraph/forgeui/plugins"
)

func main() {
    registry := plugin.NewRegistry()
    
    // Option 1: Use all plugins
    for _, p := range plugins.AllPlugins() {
        if plugin, ok := p.(plugin.Plugin); ok {
            registry.Register(plugin)
        }
    }
    
    // Option 2: Use essential plugins only
    for _, p := range plugins.EssentialPlugins() {
        if plugin, ok := p.(plugin.Plugin); ok {
            registry.Register(plugin)
        }
    }
    
    // Option 3: Pick and choose
    registry.Use(
        plugins.NewToast(),
        plugins.NewCharts(),
        plugins.NewSEO(),
    )
    
    registry.Initialize(context.Background())
}
```

### Using Plugin Presets

```go
// Essential plugins (toast, htmx, seo)
plugins.EssentialPlugins()

// Data visualization plugins (charts, sortable)
plugins.DataVisualizationPlugins()

// All built-in plugins
plugins.AllPlugins()
```

## Plugin Configuration

Each plugin supports configuration options:

```go
// Toast with custom options
toast.New(
    toast.WithPosition("bottom-center"),
    toast.WithMaxVisible(3),
    toast.WithDefaultTimeout(3000),
)

// Analytics with provider config
analytics.New(analytics.Config{
    Provider: "ga4",
    TrackingID: "G-XXXXXXXXXX",
    TrackPageViews: true,
    TrackClicks: false,
})

// HTMX with extensions
htmxplugin.New("sse", "ws", "preload")
```

## Integration Example

Complete application setup:

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/xraph/forgeui"
    "github.com/xraph/forgeui/plugin"
    "github.com/xraph/forgeui/plugins"
    "github.com/xraph/forgeui/plugins/toast"
    "github.com/xraph/forgeui/plugins/charts"
    "github.com/xraph/forgeui/plugins/analytics"
    "github.com/xraph/forgeui/plugins/seo"
)

func main() {
    // Create plugin registry
    registry := plugin.NewRegistry()
    
    // Register plugins
    registry.Use(
        toast.New(toast.WithPosition("top-right")),
        charts.New(),
        analytics.New(analytics.Config{
            Provider: "plausible",
            Domain: "example.com",
        }),
        seo.New(),
    )
    
    // Initialize plugins
    if err := registry.Initialize(context.Background()); err != nil {
        panic(err)
    }
    defer registry.Shutdown(context.Background())
    
    // Create ForgeUI app
    app := forgeui.New(forgeui.Config{
        Name: "My App",
        Dev:  true,
    })
    
    // Collect scripts from plugins
    for _, script := range registry.CollectScripts() {
        // Add scripts to your page head
    }
    
    // Start server
    http.ListenAndServe(":8080", app)
}
```

## Creating Custom Plugins

See the [Plugin Development Guide](../plugin/README.md) for creating your own plugins.

## Plugin Architecture

Built-in plugins follow these patterns:

- **Plugin Interface**: All plugins implement `plugin.Plugin`
- **Alpine Integration**: Plugins can provide stores, directives, and magic properties
- **Component Plugins**: Provide reusable UI components
- **Theme Plugins**: Extend the theme system
- **Lifecycle Hooks**: Plugins hook into ForgeUI's lifecycle events

## Testing

Run plugin tests:

```bash
go test ./plugins/...
```

Run specific plugin tests:

```bash
go test ./plugins/toast
go test ./plugins/charts
```

## Performance

Built-in plugins are optimized for:

- **Lazy loading**: Scripts load only when needed
- **CDN delivery**: External libraries from CDN
- **Tree shaking**: Use only plugins you need
- **Minimal overhead**: Lightweight plugin system

## Browser Support

All plugins support modern browsers:
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## License

All built-in plugins are MIT licensed.

