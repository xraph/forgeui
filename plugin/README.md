# Plugin System

The ForgeUI plugin system provides a flexible and extensible architecture for adding functionality to your ForgeUI applications.

## Features

- **Plugin Interface**: Simple interface for creating plugins
- **Dependency Management**: Semver-based dependency resolution
- **Lifecycle Management**: Automatic initialization and shutdown in dependency order
- **Hook System**: Extensibility points for intercepting rendering and lifecycle events
- **Thread-Safe**: Concurrent-safe registry operations
- **Plugin Discovery**: Optional dynamic plugin loading (experimental)

## Quick Start

### Creating a Plugin

```go
package myplugin

import (
    "context"
    "github.com/xraph/forgeui/plugin"
)

type MyPlugin struct {
    *plugin.PluginBase
}

func New() *MyPlugin {
    return &MyPlugin{
        PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
            Name:        "my-plugin",
            Version:     "1.0.0",
            Description: "My awesome plugin",
            Author:      "Your Name",
            License:     "MIT",
        }),
    }
}

func (p *MyPlugin) Init(ctx context.Context, r *plugin.Registry) error {
    // Initialize your plugin here
    return nil
}

func (p *MyPlugin) Shutdown(ctx context.Context) error {
    // Clean up resources here
    return nil
}
```

### Using Plugins

```go
package main

import (
    "context"
    "log"
    
    "github.com/xraph/forgeui/plugin"
    "yourapp/plugins/myplugin"
)

func main() {
    ctx := context.Background()
    
    // Create registry
    registry := plugin.NewRegistry()
    
    // Register plugins
    registry.Use(
        myplugin.New(),
        // Add more plugins...
    )
    
    // Initialize all plugins
    if err := registry.Initialize(ctx); err != nil {
        log.Fatal(err)
    }
    defer registry.Shutdown(ctx)
    
    // Your application code here...
}
```

## Dependencies

Plugins can declare dependencies on other plugins:

```go
func New() *MyPlugin {
    return &MyPlugin{
        PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
            Name:    "my-plugin",
            Version: "1.0.0",
            Dependencies: []plugin.Dependency{
                {Name: "other-plugin", Version: ">=1.0.0"},
                {Name: "optional-plugin", Version: ">=2.0.0", Optional: true},
            },
        }),
    }
}
```

### Version Constraints

The plugin system supports standard semver constraints:

- `=1.2.3` - Exact version
- `>1.2.3` - Greater than
- `>=1.2.3` - Greater than or equal
- `<1.2.3` - Less than
- `<=1.2.3` - Less than or equal
- `~1.2.3` - Patch updates (>=1.2.3 and <1.3.0)
- `^1.2.3` - Minor updates (>=1.2.3 and <2.0.0)

## Hooks

Plugins can register hooks to intercept various events:

```go
func (p *MyPlugin) Init(ctx context.Context, r *plugin.Registry) error {
    // Register a hook
    r.Hooks().On(plugin.HookBeforeRender, func(hctx *plugin.HookContext) error {
        // Do something before rendering
        return nil
    })
    
    return nil
}
```

### Available Hooks

**Lifecycle Hooks:**
- `before_init` - Before plugin initialization
- `after_init` - After plugin initialization
- `before_shutdown` - Before plugin shutdown
- `after_shutdown` - After plugin shutdown

**Render Hooks:**
- `before_render` - Before page render
- `after_render` - After page render
- `before_head` - Before `<head>` content
- `after_head` - After `<head>` content
- `before_body` - Before `<body>` content
- `after_body` - After `</body>` (scripts area)
- `before_scripts` - Before script tags
- `after_scripts` - After script tags

## Plugin Discovery (Experimental)

The plugin system supports dynamic plugin loading from `.so` files:

```go
registry := plugin.NewRegistry()

// Discover plugins from directory
if err := registry.Discover("./plugins"); err != nil {
    log.Printf("Discovery error: %v", err)
}

// Or use safe discovery (continues on error)
errs := registry.DiscoverSafe("./plugins")
for _, err := range errs {
    log.Printf("Discovery error: %v", err)
}
```

**Important Limitations:**
- Only supported on Linux, FreeBSD, and macOS
- Must be built with the same Go version as the main program
- Must use the same versions of all dependencies
- Cannot be unloaded once loaded

For production use, consider statically linking plugins instead.

## Plugin Types

ForgeUI provides specialized plugin types for common use cases:

### Component Plugins

Component plugins extend ForgeUI with new UI components:

```go
package chartplugin

import (
    "github.com/a-h/templ"
    "github.com/xraph/forgeui/plugin"
)

type ChartPlugin struct {
    *plugin.ComponentPluginBase
}

func New() *ChartPlugin {
    return &ChartPlugin{
        ComponentPluginBase: plugin.NewComponentPluginBase(
            plugin.PluginInfo{
                Name:    "charts",
                Version: "1.0.0",
            },
            map[string]plugin.ComponentConstructor{
                "LineChart": lineChartConstructor,
                "BarChart":  barChartConstructor,
            },
        ),
    }
}

func lineChartConstructor(props any, children ...templ.Component) templ.Component {
    opts := props.(*ChartOptions)
    // Return a templ component that renders the chart
    return LineChart(opts)
}
```

**Using component plugins:**

```go
registry := plugin.NewRegistry()
registry.Use(chartplugin.New())
registry.Initialize(ctx)

// Collect all components from plugins
components := registry.CollectComponents()
lineChart := components["LineChart"]
```

### Alpine.js Plugins

Alpine plugins extend Alpine.js with scripts, directives, stores, and magic properties:

```go
package sortableplugin

import "github.com/xraph/forgeui/plugin"

type SortablePlugin struct {
    *plugin.AlpinePluginBase
}

func New() *SortablePlugin {
    return &SortablePlugin{
        AlpinePluginBase: plugin.NewAlpinePluginBase(plugin.PluginInfo{
            Name:    "sortable",
            Version: "1.0.0",
        }),
    }
}

func (p *SortablePlugin) Scripts() []plugin.Script {
    return []plugin.Script{
        {
            Name:     "sortablejs",
            URL:      "https://cdn.jsdelivr.net/npm/sortablejs@1.15.0/Sortable.min.js",
            Priority: 10,
            Defer:    true,
        },
    }
}

func (p *SortablePlugin) Directives() []plugin.AlpineDirective {
    return []plugin.AlpineDirective{
        {
            Name: "sortable",
            Definition: `
                (el, { expression, modifiers }, { evaluate }) => {
                    let options = expression ? evaluate(expression) : {};
                    new Sortable(el, options);
                }
            `,
        },
    }
}

func (p *SortablePlugin) Stores() []plugin.AlpineStore {
    return []plugin.AlpineStore{
        {
            Name: "sortable",
            InitialState: map[string]any{
                "items": []any{},
            },
            Methods: `
                add(item) {
                    this.items.push(item);
                },
                remove(id) {
                    this.items = this.items.filter(i => i.id !== id);
                }
            `,
        },
    }
}
```

**Collecting Alpine assets:**

```go
registry := plugin.NewRegistry()
registry.Use(sortableplugin.New())
registry.Initialize(ctx)

// Collect all Alpine assets
scripts := registry.CollectScripts()        // All external scripts
directives := registry.CollectDirectives()  // All custom directives
stores := registry.CollectStores()          // All Alpine stores
magics := registry.CollectMagics()          // All magic properties
components := registry.CollectAlpineComponents() // All Alpine.data components
```

### Theme Plugins

Theme plugins provide custom themes and fonts:

```go
package corporatetheme

import (
    "github.com/xraph/forgeui/plugin"
    "github.com/xraph/forgeui/theme"
)

type CorporateThemePlugin struct {
    *plugin.ThemePluginBase
}

func New() *CorporateThemePlugin {
    themes := map[string]theme.Theme{
        "corporate-light": {
            Colors: theme.ColorTokens{
                Primary:   "#003366",
                Secondary: "#0066CC",
                // ... more colors
            },
            Radius: theme.RadiusTokens{
                Default: "0.25rem",
                // ... more radii
            },
        },
        "corporate-dark": {
            Colors: theme.ColorTokens{
                Primary:   "#0066CC",
                Secondary: "#003366",
                // ... more colors
            },
            Radius: theme.RadiusTokens{
                Default: "0.25rem",
            },
        },
    }

    fonts := []theme.Font{
        {Family: "Inter", Weights: []int{400, 600, 700}},
        {Family: "JetBrains Mono", Weights: []int{400, 500}},
    }

    return &CorporateThemePlugin{
        ThemePluginBase: plugin.NewThemePluginBaseWithFonts(
            plugin.PluginInfo{
                Name:    "corporate-theme",
                Version: "1.0.0",
            },
            themes,
            "corporate-light",
            fonts,
        ),
    }
}

func (p *CorporateThemePlugin) CSS() string {
    return `
        :root {
            --corporate-accent: #FF6600;
        }
        .corporate-button {
            background: var(--corporate-accent);
        }
    `
}
```

### Middleware Plugins

Middleware plugins provide HTTP middleware with priority-based execution:

```go
package htmxplugin

import (
    "context"
    "net/http"
    "github.com/xraph/forgeui/plugin"
)

type HTMXPlugin struct {
    *plugin.MiddlewarePluginBase
}

func New() *HTMXPlugin {
    middleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Detect HTMX requests
            if r.Header.Get("HX-Request") == "true" {
                ctx := context.WithValue(r.Context(), htmxRequestKey, true)
                r = r.WithContext(ctx)
            }
            next.ServeHTTP(w, r)
        })
    }

    return &HTMXPlugin{
        MiddlewarePluginBase: plugin.NewMiddlewarePluginBase(
            plugin.PluginInfo{
                Name:    "htmx",
                Version: "1.0.0",
            },
            middleware,
            10, // Low priority = executes early
        ),
    }
}
```

**Using middleware plugins:**

```go
registry := plugin.NewRegistry()
registry.Use(htmxplugin.New())
registry.Initialize(ctx)

// Collect middleware in priority order
middlewares := registry.CollectMiddleware()

// Chain middleware
handler := yourHandler
for i := len(middlewares) - 1; i >= 0; i-- {
    handler = middlewares[i].Middleware()(handler)
}
```

### Multi-Type Plugins

A plugin can implement multiple interfaces:

```go
type SuperPlugin struct {
    *plugin.PluginBase
}

// Implements ComponentPlugin
func (p *SuperPlugin) Components() map[string]plugin.ComponentConstructor {
    return map[string]plugin.ComponentConstructor{
        "SuperComponent": superComponentConstructor,
    }
}

func (p *SuperPlugin) CVAExtensions() map[string]*forgeui.CVA {
    return nil
}

// Implements AlpinePlugin
func (p *SuperPlugin) Scripts() []plugin.Script {
    return []plugin.Script{
        {Name: "super-script", URL: "https://example.com/super.js"},
    }
}

func (p *SuperPlugin) Directives() []plugin.AlpineDirective {
    return []plugin.AlpineDirective{
        {Name: "super", Definition: "(el) => {}"},
    }
}

// ... implement other AlpinePlugin methods

// Implements ThemePlugin
func (p *SuperPlugin) Themes() map[string]theme.Theme {
    return map[string]theme.Theme{
        "super": theme.DefaultLight(),
    }
}

// ... implement other ThemePlugin methods
```

## Advanced Usage

### Type-Specific Retrieval

```go
// Get specific plugin types
if cp, ok := registry.GetComponentPlugin("charts"); ok {
    components := cp.Components()
}

if ap, ok := registry.GetAlpinePlugin("sortable"); ok {
    scripts := ap.Scripts()
}

if tp, ok := registry.GetThemePlugin("corporate-theme"); ok {
    themes := tp.Themes()
}
```

### Accessing Other Plugins

Plugins can access other plugins through the registry:

```go
func (p *MyPlugin) Init(ctx context.Context, r *plugin.Registry) error {
    // Get another plugin
    if other, ok := r.Get("other-plugin"); ok {
        // Use the other plugin
    }
    
    return nil
}
```

### Hook Context

Hooks receive a context with useful information:

```go
r.Hooks().On(plugin.HookBeforeRender, func(hctx *plugin.HookContext) error {
    // Access context
    ctx := hctx.Context
    
    // Access data
    pluginName := hctx.Data["plugin"].(string)
    
    // Access nodes (for render hooks)
    nodes := hctx.Nodes
    
    return nil
})
```

## Best Practices

1. **Keep plugins focused**: Each plugin should do one thing well
2. **Document dependencies**: Clearly document what your plugin depends on
3. **Handle errors gracefully**: Return clear error messages from Init/Shutdown
4. **Clean up resources**: Always clean up in Shutdown
5. **Use hooks sparingly**: Only hook into events you actually need
6. **Test thoroughly**: Write tests for your plugin's functionality
7. **Version carefully**: Follow semantic versioning for your plugins

## Testing

Example test for a plugin:

```go
func TestMyPlugin(t *testing.T) {
    ctx := context.Background()
    registry := plugin.NewRegistry()
    
    p := New()
    if err := registry.Register(p); err != nil {
        t.Fatalf("Register failed: %v", err)
    }
    
    if err := registry.Initialize(ctx); err != nil {
        t.Fatalf("Initialize failed: %v", err)
    }
    defer registry.Shutdown(ctx)
    
    // Test your plugin functionality...
}
```

## API Reference

### Plugin Interface

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Dependencies() []Dependency
    Init(ctx context.Context, registry *Registry) error
    Shutdown(ctx context.Context) error
}
```

### Registry Methods

**Basic Operations:**
- `NewRegistry() *Registry` - Create a new registry
- `Register(p Plugin) error` - Register a plugin
- `Get(name string) (Plugin, bool)` - Get a plugin by name
- `Has(name string) bool` - Check if a plugin exists
- `All() []Plugin` - Get all plugins
- `Count() int` - Get plugin count
- `Use(plugins ...Plugin) *Registry` - Register multiple plugins (chainable)
- `Unregister(name string) error` - Remove a plugin
- `Initialize(ctx context.Context) error` - Initialize all plugins
- `Shutdown(ctx context.Context) error` - Shutdown all plugins
- `Hooks() *HookManager` - Get the hook manager

**Type-Specific Retrieval:**
- `GetComponentPlugin(name string) (ComponentPlugin, bool)` - Get a component plugin
- `GetAlpinePlugin(name string) (AlpinePlugin, bool)` - Get an Alpine plugin
- `GetThemePlugin(name string) (ThemePlugin, bool)` - Get a theme plugin

**Asset Collection:**
- `CollectComponents() map[string]ComponentConstructor` - Collect all component constructors
- `CollectScripts() []Script` - Collect all scripts (sorted by priority)
- `CollectDirectives() []AlpineDirective` - Collect all Alpine directives
- `CollectStores() []AlpineStore` - Collect all Alpine stores
- `CollectMagics() []AlpineMagic` - Collect all Alpine magic properties
- `CollectAlpineComponents() []AlpineComponent` - Collect all Alpine.data components
- `CollectMiddleware() []MiddlewarePlugin` - Collect all middleware (sorted by priority)

### HookManager Methods

- `On(hook string, fn HookFunc)` - Register a hook handler
- `Off(hook string)` - Remove all handlers for a hook
- `Trigger(hook string, ctx *HookContext) error` - Execute hook handlers
- `Has(hook string) bool` - Check if a hook has handlers
- `Count(hook string) int` - Get handler count for a hook

## Examples

See the `example/` directory for complete examples of:
- Creating custom plugins
- Using dependencies
- Registering hooks
- Plugin discovery

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code coverage remains above 80%
- Documentation is updated
- Examples are provided for new features

