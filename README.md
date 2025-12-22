# ForgeUI

ForgeUI is a comprehensive SSR-first UI framework for Go, built on gomponents with Tailwind CSS styling and shadcn-inspired design patterns. It provides everything you need to build modern, interactive web applications entirely in Go.

## Features

### Core Framework
- ✅ **SSR-First**: Pure Go component rendering with zero client-side dependencies required
- ✅ **Type-Safe**: Full Go type safety with functional options pattern
- ✅ **CVA**: Class Variance Authority for flexible variant management
- ✅ **Tailwind CSS**: Utility-first CSS styling with built-in processing
- ✅ **35+ Components**: Production-ready UI components

### Frontend Integration
- ✅ **Alpine.js Integration**: Directives, stores, magic helpers, and plugins
- ✅ **HTMX Support**: Complete HTMX attribute helpers and server-side utilities
- ✅ **Icons**: 1600+ Lucide icons with customization options
- ✅ **Animation System**: Tailwind animations, transitions, and keyframes

### Backend Features
- ✅ **Router**: Production-ready HTTP routing with middleware support
- ✅ **Bridge**: Go-JavaScript RPC bridge for calling Go functions from client-side
- ✅ **Plugin System**: Extensible plugin architecture with dependency management
- ✅ **Theme System**: Customizable themes with CSS variables and color tokens

### Developer Tools
- ✅ **Assets Pipeline**: Built-in esbuild, Tailwind CSS, and file fingerprinting
- ✅ **Dev Server**: Hot-reload development server with file watching
- ✅ **CLI**: Command-line tools for project scaffolding and management
- ✅ **Layout Helpers**: Page builder with meta tags, scripts, and structured layouts

## Why ForgeUI?

- 🚀 **Go All The Way**: Write your entire frontend in Go - no JavaScript required
- 🎯 **Type Safety**: Catch errors at compile time, not runtime
- ⚡ **SSR Performance**: Server-rendered HTML with optional progressive enhancement
- 🎨 **Beautiful by Default**: shadcn-inspired design that looks great out of the box
- 🔧 **Full Stack**: Router, RPC, assets, themes - everything you need
- 📦 **Zero Config**: Works out of the box with sensible defaults
- 🔌 **Extensible**: Plugin system for adding custom functionality
- 🎭 **Progressive**: Start with pure SSR, add Alpine.js/HTMX as needed

## Installation

```bash
go get github.com/xraph/forgeui
```

## Quick Start

### 1. Create Your First Application

```go
package main

import (
    "log"
    "net/http"
    
    g "maragu.dev/gomponents"
    "maragu.dev/gomponents/html"
    
    "github.com/xraph/forgeui"
    "github.com/xraph/forgeui/components/button"
    "github.com/xraph/forgeui/components/card"
    "github.com/xraph/forgeui/primitives"
    "github.com/xraph/forgeui/router"
)

func main() {
    // Create ForgeUI app
    app := forgeui.New(
        forgeui.WithDebug(true),
        forgeui.WithThemeName("default"),
    )
    
    // Register routes
    app.Get("/", HomePage)
    
    // Start server
    log.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", app.Router())
}

func HomePage(ctx *router.PageContext) (g.Node, error) {
    return html.Html(
        html.Head(
            html.Meta(html.Charset("UTF-8")),
            html.Title(g.Text("ForgeUI App")),
            html.Link(html.Rel("stylesheet"), html.Href("/static/styles.css")),
        ),
        html.Body(
            primitives.Container(
                primitives.VStack("8",
                    card.Card(
                        card.Header(
                            card.Title("Welcome to ForgeUI"),
                            card.Description("Build beautiful UIs with Go"),
                        ),
                        card.Content(
                            primitives.VStack("4",
                                primitives.Text(
                                    primitives.TextChildren(
                                        g.Text("ForgeUI provides type-safe, composable UI components."),
                                    ),
                                ),
                                button.Primary(
                                    g.Text("Get Started"),
                                    button.WithSize(forgeui.SizeLG),
                                ),
                            ),
                        ),
                    ),
                ),
            ),
        ),
    ), nil
}
```

### 2. Run Your Application

```bash
go run main.go
```

Visit `http://localhost:8080` to see your app!

### 3. Add Interactivity (Optional)

Add Alpine.js for client-side interactivity:

```go
import (
    "github.com/xraph/forgeui/alpine"
)

func HomePage(ctx *router.PageContext) (g.Node, error) {
    return html.Html(
        html.Head(
            html.Title(g.Text("Interactive App")),
            alpine.Scripts(), // Add Alpine.js
        ),
        html.Body(
            html.Div(
                alpine.XData(`{count: 0}`),
                html.Button(
                    alpine.XOn("click", "count++"),
                    g.Text("Increment"),
                ),
                html.Span(alpine.XText("count")),
            ),
        ),
    ), nil
}
```

## Components

### Layout Primitives
- ✅ **Box** - Polymorphic container with responsive props
- ✅ **Flex** - Flexbox layout with direction and alignment
- ✅ **Grid** - CSS Grid layout with responsive columns
- ✅ **Stack** - VStack/HStack for vertical/horizontal layouts
- ✅ **Center** - Centered container with responsive centering
- ✅ **Container** - Responsive container with max-width
- ✅ **Text** - Typography primitive with semantic HTML
- ✅ **Spacer** - Flexible spacer for layouts

### Buttons & Actions
- ✅ **Button** - All variants (Primary, Secondary, Destructive, Outline, Ghost, Link)
- ✅ **Button Group** - Grouped buttons with gap control
- ✅ **Icon Button** - Compact icon-only buttons

### Content Display
- ✅ **Card** - Compound component (Header, Title, Description, Content, Footer)
- ✅ **Badge** - Status indicators and labels
- ✅ **Avatar** - User avatars with image and fallback
- ✅ **Alert** - Alert messages with variants (Info, Success, Warning, Error)
- ✅ **Separator** - Horizontal/vertical dividers
- ✅ **Empty State** - Empty state placeholders
- ✅ **List** - List containers with list items

### Navigation
- ✅ **Navbar** - Navigation bar component
- ✅ **Breadcrumb** - Breadcrumb navigation
- ✅ **Tabs** - Tabbed navigation and content
- ✅ **Menu** - Menu and menu items
- ✅ **Sidebar** - Collapsible sidebar navigation
- ✅ **Pagination** - Page navigation controls

### Forms
- ✅ **Form** - Form wrapper with validation helpers
- ✅ **Label** - Accessible form labels
- ✅ **Input** - Text inputs with variants and validation states
- ✅ **Input Group** - Input with icons and addons
- ✅ **Textarea** - Multi-line text input
- ✅ **Checkbox** - Checkbox inputs with labels
- ✅ **Radio** - Radio buttons with radio groups
- ✅ **Switch** - Toggle switches
- ✅ **Select** - Native select dropdowns
- ✅ **Slider** - Range sliders

### Overlays
- ✅ **Modal** - Modal dialogs
- ✅ **Dialog** - Dialog component
- ✅ **Alert Dialog** - Confirmation dialogs
- ✅ **Drawer** - Slide-out panels
- ✅ **Sheet** - Bottom/side sheets
- ✅ **Dropdown** - Dropdown menus
- ✅ **Context Menu** - Right-click context menus
- ✅ **Popover** - Floating popovers
- ✅ **Tooltip** - Hover tooltips
- ✅ **Toast** - Toast notifications with toaster

### Feedback
- ✅ **Spinner** - Loading spinners
- ✅ **Skeleton** - Loading placeholders
- ✅ **Progress** - Progress bars

### Data Display
- ✅ **Table** - Data tables with header, body, rows, and cells
- ✅ **DataTable** - Advanced tables with sorting, filtering, and pagination

## Usage Examples

### Button Variants

```go
// Primary button
button.Primary(g.Text("Save"))

// Destructive button
button.Destructive(
    g.Text("Delete"),
    button.WithSize(forgeui.SizeLG),
)

// Icon button
button.IconButton(
    g.Text("×"),
    button.WithVariant(forgeui.VariantGhost),
)

// Button group
button.Group(
    []button.GroupOption{button.WithGap("2")},
    button.Primary(g.Text("Save")),
    button.Secondary(g.Text("Cancel")),
)
```

### Card Component

```go
card.Card(
    card.Header(
        card.Title("Settings"),
        card.Description("Manage your account settings"),
    ),
    card.Content(
        // Your content here
    ),
    card.Footer(
        button.Primary(g.Text("Save Changes")),
    ),
)
```

### Form Example

```go
form.Form(
    []form.Option{
        form.WithAction("/submit"),
        form.WithMethod("POST"),
    },
    input.Field(
        "Email",
        []input.Option{
            input.WithType("email"),
            input.WithName("email"),
            input.WithPlaceholder("Enter your email"),
            input.Required(),
        },
        input.FormDescription("We'll never share your email."),
    ),
    checkbox.Checkbox(
        checkbox.WithName("subscribe"),
        checkbox.WithID("subscribe"),
    ),
    button.Primary(
        g.Text("Submit"),
        button.WithType("submit"),
    ),
)
```

### Layout with Primitives

```go
primitives.Container(
    primitives.VStack("8",
        primitives.Text(
            primitives.TextAs("h1"),
            primitives.TextSize("text-4xl"),
            primitives.TextWeight("font-bold"),
            primitives.TextChildren(g.Text("Dashboard")),
        ),
        primitives.Grid(
            primitives.GridCols(1),
            primitives.GridColsMD(2),
            primitives.GridColsLG(3),
            primitives.GridGap("6"),
            primitives.GridChildren(
                card.Card(/* ... */),
                card.Card(/* ... */),
                card.Card(/* ... */),
            ),
        ),
    ),
)
```

## Router

Production-ready HTTP routing with pattern matching and middleware:

```go
app := forgeui.New()

// Static routes
app.Get("/", HomePage)
app.Get("/about", AboutPage)

// Path parameters
app.Get("/users/:id", UserProfile)
app.Get("/users/:userId/posts/:postId", PostDetail)

// Wildcards
app.Get("/files/*filepath", ServeFile)

// Middleware
app.Use(router.Logger())
app.Use(router.Recovery())
app.Use(router.CORS("*"))

// Route-specific middleware
route := app.Get("/admin", AdminDashboard)
route.WithMiddleware(AuthMiddleware)

// Named routes for URL generation
app.Router().Name("user.post", route)
url := app.Router().URL("user.post", userID, postID)
```

### PageContext

Access request data with rich context utilities:

```go
func UserProfile(ctx *router.PageContext) (g.Node, error) {
    // Path parameters
    id := ctx.Param("id")
    userID, _ := ctx.ParamInt("id")
    
    // Query parameters
    query := ctx.Query("q")
    page, _ := ctx.QueryInt("page")
    
    // Headers and cookies
    auth := ctx.Header("Authorization")
    cookie, _ := ctx.Cookie("session")
    
    // Context values (set by middleware)
    userID := ctx.GetInt("user_id")
    
    return html.Div(/* ... */), nil
}
```

See [router/README.md](router/README.md) for complete documentation.

## Bridge - Go to JavaScript RPC

Call Go functions directly from JavaScript using JSON-RPC 2.0:

### Server Side

```go
b := bridge.New()

// Register Go functions
b.Register("createUser", func(ctx bridge.Context, input struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}) (struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}, error) {
    user := createUserInDB(input.Name, input.Email)
    return struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }{ID: user.ID, Name: user.Name}, nil
})

// Function options
b.Register("adminAction", handler,
    bridge.RequireAuth(),
    bridge.RequireRoles("admin"),
    bridge.WithRateLimit(10),
    bridge.WithTimeout(10*time.Second),
)

// Enable bridge on router
bridge.EnableBridge(app.Router(), b)
```

### Client Side

```go
// Include bridge scripts
html.Head(
    bridge.BridgeScripts(bridge.ScriptConfig{
        Endpoint: "/api/bridge",
        CSRFToken: csrfToken,
        IncludeAlpine: true,
    }),
)

// Alpine.js integration
html.Button(
    g.Attr("@click", "result = await $bridge.call('createUser', {name, email})"),
    g.Text("Create User"),
)
```

**Features:**
- HTTP (JSON-RPC 2.0), WebSocket, and SSE transports
- Built-in authentication, authorization, and rate limiting
- Automatic parameter validation
- CSRF protection
- Caching support
- Alpine.js magic helpers

See [bridge/README.md](bridge/README.md) for complete documentation.

## HTMX Integration

Complete HTMX support with type-safe attribute helpers:

```go
// Include HTMX
html.Head(
    htmx.Scripts(),
    htmx.IndicatorCSS(),
)

// HTMX attributes
html.Button(
    htmx.HxGet("/api/users"),
    htmx.HxTarget("#user-list"),
    htmx.HxSwap("innerHTML"),
    g.Text("Load Users"),
)

// Advanced triggers
html.Input(
    htmx.HxTriggerDebounce("keyup", "500ms"),
    htmx.HxGet("/search"),
)

// Server-side detection
func handler(w http.ResponseWriter, r *http.Request) {
    if htmx.IsHTMX(r) {
        // Return partial HTML
        renderPartial(w)
    } else {
        // Return full page
        renderFullPage(w)
    }
}

// Response headers
htmx.TriggerEvent(w, "refresh")
htmx.SetHTMXRedirect(w, "/login")
```

See [htmx/README.md](htmx/README.md) for complete documentation.

## Alpine.js Integration

Seamless Alpine.js integration with directives, stores, and plugins:

```go
// Include Alpine.js
html.Head(
    alpine.Scripts(),
)

// Alpine directives
html.Div(
    alpine.XData(`{count: 0}`),
    html.Button(
        alpine.XOn("click", "count++"),
        g.Text("Increment"),
    ),
    html.Span(alpine.XText("count")),
)

// Alpine stores
alpine.Store("app", map[string]any{
    "user": nil,
    "isLoggedIn": false,
})

// Custom directives
alpine.Directive("click-outside", `(el, {expression}, {evaluate}) => {
    // directive implementation
}`)
```

## Icons

1600+ Lucide icons with full customization:

```go
// Basic usage
icons.Check()
icons.Search()
icons.User()

// Customization
icons.Check(
    icons.WithSize(24),
    icons.WithColor("green"),
    icons.WithStrokeWidth(2),
    icons.WithClass("text-green-500"),
)

// In buttons
button.Primary(
    g.Group([]g.Node{
        icons.Plus(icons.WithSize(16)),
        g.Text("Add Item"),
    }),
)

// Loading spinner
icons.Loader(
    icons.WithSize(20),
    icons.WithClass("animate-spin"),
)
```

See [icons/README.md](icons/README.md) for complete documentation.

## Plugin System

Extensible plugin architecture with dependency management:

```go
// Create a plugin
type MyPlugin struct {
    *plugin.PluginBase
}

func New() *MyPlugin {
    return &MyPlugin{
        PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
            Name:    "my-plugin",
            Version: "1.0.0",
            Dependencies: []plugin.Dependency{
                {Name: "other-plugin", Version: ">=1.0.0"},
            },
        }),
    }
}

func (p *MyPlugin) Init(ctx context.Context, r *plugin.Registry) error {
    // Initialize plugin
    return nil
}

// Use plugins
registry := plugin.NewRegistry()
registry.Use(
    myplugin.New(),
    otherplugin.New(),
)
registry.Initialize(ctx)
defer registry.Shutdown(ctx)
```

**Plugin Types:**
- Component plugins (extend UI components)
- Alpine plugins (scripts, directives, stores)
- Theme plugins (custom themes and fonts)
- Middleware plugins (HTTP middleware)

See [plugin/README.md](plugin/README.md) for complete documentation.

## Assets Pipeline

Built-in asset management with Tailwind CSS, esbuild, and fingerprinting:

```go
app := forgeui.New()

// Configure assets
app.WithAssets(assets.Config{
    Enabled:      true,
    OutputDir:    "./dist",
    Minify:       true,
    Fingerprint:  true,
    Tailwind: assets.TailwindConfig{
        Input:     "./styles/input.css",
        Output:    "./dist/styles.css",
        ConfigFile: "./tailwind.config.js",
    },
    ESBuild: assets.ESBuildConfig{
        EntryPoints: []string{"./js/app.js"},
        Outdir:     "./dist/js",
    },
})

// Development server with hot-reload
if err := app.StartDevServer(":3000"); err != nil {
    log.Fatal(err)
}
```

**Features:**
- Tailwind CSS compilation
- JavaScript bundling with esbuild
- File fingerprinting for cache busting
- Hot-reload development server
- Asset manifest generation

See [assets/README.md](assets/README.md) for complete documentation.

## Theme System

Customizable themes with CSS variables:

```go
// Use built-in theme
app := forgeui.New(
    forgeui.WithThemeName("default"),
)

// Create custom theme
customTheme := theme.Theme{
    Colors: theme.ColorTokens{
        Primary:   "#3b82f6",
        Secondary: "#64748b",
        Background: "#ffffff",
        Foreground: "#0f172a",
    },
    Radius: theme.RadiusTokens{
        Default: "0.5rem",
        Button:  "0.375rem",
        Card:    "0.75rem",
    },
}

app.RegisterTheme("custom", customTheme)
```

## CLI Tools

Command-line tools for project management:

```bash
# Create new project
forgeui new myproject

# Create new component
forgeui create component MyComponent

# Create new page
forgeui create page HomePage

# Run development server
forgeui dev

# Build for production
forgeui build
```

## Architecture

ForgeUI follows a layered architecture designed for scalability and maintainability:

```text
┌─────────────────────────────────────┐
│  Application Layer                  │
│  (Your App, Examples, CLI)          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Integration Layer                  │
│  (Router, Bridge, HTMX, Alpine)     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Components Layer                   │
│  (UI Components, Icons, Layouts)    │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Foundation Layer                   │
│  (Primitives, CVA, Theme, Assets)   │
└─────────────────────────────────────┘
```

### Design Principles
- **Composability**: All components are composable gomponents nodes
- **Type Safety**: Full Go type checking at compile time
- **Minimal Dependencies**: Core components have zero client-side dependencies
- **Progressive Enhancement**: Add interactivity with Alpine.js and HTMX
- **Unidirectional Flow**: Higher layers depend on lower layers, never the reverse

## Testing

ForgeUI includes comprehensive test coverage across all packages:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run specific package tests
go test ./components/button/...
go test ./router/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Structure

Components and features include:
- **Unit tests**: Testing individual functions and components
- **Integration tests**: Testing component interactions
- **HTTP tests**: Using `httptest` for router and bridge testing
- **Golden files**: Snapshot testing for component output

## Configuration

### Application Configuration

```go
app := forgeui.New(
    // Core options
    forgeui.WithDebug(true),
    forgeui.WithThemeName("default"),
    forgeui.WithStaticPath("/static"),
    forgeui.WithDefaultSize(forgeui.SizeMD),
    forgeui.WithDefaultVariant(forgeui.VariantDefault),
    forgeui.WithDefaultRadius(forgeui.RadiusMD),
    
    // Assets configuration
    forgeui.WithAssets(assets.Config{
        Enabled:     true,
        OutputDir:   "./dist",
        Minify:      true,
        Fingerprint: true,
    }),
)

// Initialize the application
if err := app.Initialize(context.Background()); err != nil {
    log.Fatal(err)
}

// Start HTTP server
http.ListenAndServe(":8080", app.Router())
```

### Environment-Specific Config

```go
var cfg forgeui.Config

if os.Getenv("ENV") == "production" {
    cfg = forgeui.Config{
        Debug:      false,
        StaticPath: "/static",
        Assets: assets.Config{
            Enabled:     true,
            Minify:      true,
            Fingerprint: true,
        },
    }
} else {
    cfg = forgeui.Config{
        Debug:      true,
        StaticPath: "/static",
        Assets: assets.Config{
            Enabled:    true,
            Minify:     false,
            Fingerprint: false,
        },
    }
}

app := forgeui.NewWithConfig(cfg)
```

## Examples

The `example/` directory contains complete working examples:

```bash
cd example
go run main.go
```

Visit `http://localhost:8080` to see:
- **Component Showcase**: All UI components with variations
- **Dashboard Demo**: Real-world dashboard layout
- **Interactive Examples**: Alpine.js and HTMX integration
- **Bridge Demo**: Go-JavaScript RPC examples
- **Assets Demo**: Asset pipeline in action

## Documentation

Each package includes comprehensive documentation:

- [Router](router/README.md) - HTTP routing system
- [Bridge](bridge/README.md) - Go-JavaScript RPC bridge
- [HTMX](htmx/README.md) - HTMX integration
- [Icons](icons/README.md) - Icon system
- [Plugin](plugin/README.md) - Plugin architecture
- [Assets](assets/README.md) - Asset management
- [Components](components/) - Individual component docs

## Project Structure

```text
forgeui/
├── alpine/           # Alpine.js integration
├── animation/        # Animation and transition utilities
├── assets/          # Asset pipeline (CSS, JS, Tailwind)
├── bridge/          # Go-JavaScript RPC bridge
├── cli/             # Command-line tools
├── components/      # UI components
│   ├── button/
│   ├── card/
│   ├── form/
│   └── ...
├── example/         # Example applications
├── htmx/           # HTMX integration
├── icons/          # Icon system (Lucide)
├── layout/         # Layout helpers
├── plugin/         # Plugin system
├── primitives/     # Layout primitives
├── router/         # HTTP router
├── theme/          # Theme system
└── ...
```

## Contributing

ForgeUI welcomes contributions! Follow these guidelines:

### Code Style
1. **Follow Go conventions**: Use `gofmt`, `golangci-lint`
2. **Use CVA for variants**: All component variants should use CVA
3. **Functional options pattern**: Use for component configuration
4. **Documentation**: Document all exported functions with examples
5. **Type safety**: Leverage Go's type system

### Component Development
1. **Respect architecture layers**: Components should only import from lower layers
2. **Test thoroughly**: Aim for 80%+ coverage
3. **Include examples**: Add examples to component documentation
4. **Accessibility**: Follow ARIA guidelines
5. **Mobile-first**: Ensure responsive design

### Running Tests

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run linter
golangci-lint run
```

### Pull Request Process
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass
5. Update documentation
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Best Practices

### Component Composition

```go
// ✅ Good: Compose components
button.Primary(
    g.Group([]g.Node{
        icons.Plus(icons.WithSize(16)),
        g.Text("Add Item"),
    }),
)

// ❌ Bad: Don't nest unrelated components
```

### Error Handling

```go
// ✅ Good: Return errors from handlers
func UserProfile(ctx *router.PageContext) (g.Node, error) {
    user, err := getUser(ctx.Param("id"))
    if err != nil {
        return nil, err
    }
    return renderUser(user), nil
}

// ❌ Bad: Don't panic or ignore errors
```

### Context Usage

```go
// ✅ Good: Pass context through the stack
func handler(ctx *router.PageContext) (g.Node, error) {
    user := ctx.GetString("user_id")
    // Use user...
}

// ❌ Bad: Don't use global variables
```

### Performance

```go
// ✅ Good: Reuse nodes when possible
var headerNode = html.Header(/* ... */)

// ✅ Good: Use streaming for large responses
// ✅ Good: Enable asset fingerprinting in production
// ✅ Good: Minify CSS and JS in production
```

## Performance

ForgeUI is designed for performance:

### Server-Side Rendering
- **Fast**: Pure Go rendering with minimal overhead
- **Efficient**: No JavaScript runtime required for basic pages
- **Scalable**: Handles thousands of requests per second

### Asset Pipeline
- **Optimized**: Built-in minification and fingerprinting
- **Cached**: Long-lived browser caching with cache busting
- **Small**: Tailwind CSS purging removes unused styles

### Progressive Enhancement
- **Lightweight**: Start with 0KB of JavaScript
- **Optional**: Add Alpine.js (15KB) or HTMX (14KB) only when needed
- **Lazy**: Load components on demand

### Benchmarks

Typical performance on modest hardware (4 core, 8GB RAM):

- **Simple page render**: ~0.5ms
- **Complex dashboard**: ~2ms
- **Component with children**: ~0.1ms per child
- **Throughput**: 10,000+ req/s

## Production Checklist

Before deploying to production:

- [ ] Set `WithDebug(false)`
- [ ] Enable asset minification
- [ ] Enable asset fingerprinting
- [ ] Use HTTPS (TLS)
- [ ] Set up proper logging
- [ ] Configure CORS if needed
- [ ] Set up rate limiting
- [ ] Enable CSRF protection
- [ ] Configure proper caching headers
- [ ] Test with production data
- [ ] Set up monitoring and alerts

## Troubleshooting

### Common Issues

**Components not rendering?**
- Check if you've initialized the app with `app.Initialize(ctx)`
- Verify Tailwind CSS is included in your HTML

**Styles not applying?**
- Ensure Tailwind CSS is properly configured
- Check if the asset pipeline is enabled
- Verify the CSS file is being served

**HTMX not working?**
- Include `htmx.Scripts()` in your HTML head
- Check browser console for errors
- Verify HTMX attributes are correct

**Alpine.js not working?**
- Include `alpine.Scripts()` in your HTML head
- Check for JavaScript errors in console
- Verify Alpine syntax is correct

## Community & Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/xraph/forgeui/issues)
- **Discussions**: [Ask questions and share ideas](https://github.com/xraph/forgeui/discussions)
- **Examples**: Check the `example/` directory for working code

## Roadmap

Planned features and improvements:

- [ ] More component variants and customization options
- [ ] Server-side validation helpers
- [ ] Form builder with automatic validation
- [ ] More built-in plugins
- [ ] WebSocket utilities
- [ ] Admin panel template
- [ ] CLI improvements (scaffolding, generators)
- [ ] Performance monitoring utilities
- [ ] i18n support

## License

MIT License - See [LICENSE](LICENSE) file for details

## Credits & Inspiration

ForgeUI stands on the shoulders of giants:

- **[gomponents](https://maragu.dev/gomponents)** - The foundation for Go HTML components
- **[shadcn/ui](https://ui.shadcn.com/)** - Design inspiration and component patterns
- **[Tailwind CSS](https://tailwindcss.com/)** - Utility-first CSS framework
- **[Alpine.js](https://alpinejs.dev/)** - Lightweight JavaScript framework
- **[HTMX](https://htmx.org/)** - High power tools for HTML
- **[Lucide](https://lucide.dev/)** - Beautiful icon library

Made with ❤️ for the Go community