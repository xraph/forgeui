package forgeui

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"

	"github.com/xraph/forgeui/assets"
	"github.com/xraph/forgeui/bridge"
	"github.com/xraph/forgeui/router"
	"github.com/xraph/forgeui/theme"
)

// App is the main ForgeUI application with enhanced features
type App struct {
	config     *AppConfig
	Assets     *assets.Manager
	router     *router.Router
	bridge     *bridge.Bridge
	lightTheme *theme.Theme
	darkTheme  *theme.Theme
	staticPath string
	cssBuilt   bool // true when CSS was compiled via Tailwind CLI
}

// New creates a new ForgeUI application with enhanced initialization
func New(opts ...AppOption) *App {
	config := DefaultAppConfig()
	for _, opt := range opts {
		opt(config)
	}

	// Build full static path with base path if provided
	staticPath := config.StaticPath
	if config.BasePath != "" {
		staticPath = config.BasePath + config.StaticPath
	}

	// Initialize asset manager
	assetManager := assets.NewManager(assets.Config{
		PublicDir:  config.AssetPublicDir,
		OutputDir:  config.AssetOutputDir,
		StaticPath: staticPath,
		IsDev:      config.Debug,
		Manifest:   config.AssetManifest,
		FileSystem: config.AssetFileSystem,
	})

	// Initialize router
	r := router.New()
	if config.DefaultLayout != "" {
		r.SetDefaultLayout(config.DefaultLayout)
	}

	// Initialize bridge if enabled
	var b *bridge.Bridge

	if config.EnableBridge {
		if config.BridgeConfig != nil {
			b = bridge.New(func(c *bridge.Config) {
				*c = *config.BridgeConfig
			})
		} else {
			b = bridge.New()
		}
	}

	app := &App{
		config:     config,
		Assets:     assetManager,
		router:     r,
		bridge:     b,
		lightTheme: config.LightTheme,
		darkTheme:  config.DarkTheme,
		staticPath: staticPath,
	}

	// Set app reference in router for PageContext
	r.SetApp(app)

	return app
}

// Config returns the application configuration
func (a *App) Config() *AppConfig {
	return a.config
}

// IsDev returns true if the application is in debug/development mode
func (a *App) IsDev() bool {
	return a.config.Debug
}

// Router returns the application's HTTP router
func (a *App) Router() *router.Router {
	return a.router
}

// Bridge returns the bridge system (may be nil if not enabled)
func (a *App) Bridge() *bridge.Bridge {
	return a.bridge
}

// HasBridge returns true if bridge system is enabled
func (a *App) HasBridge() bool {
	return a.bridge != nil
}

// BridgeCallPath returns the full bridge call endpoint path
func (a *App) BridgeCallPath() string {
	if a.config.BasePath != "" {
		return a.config.BasePath + "/bridge/call"
	}
	return "/api/bridge/call"
}

// BridgeStreamPath returns the full bridge stream endpoint path
func (a *App) BridgeStreamPath() string {
	if a.config.BasePath != "" {
		return a.config.BasePath + "/bridge/stream/"
	}
	return "/api/bridge/stream/"
}

// BridgeScriptPath returns the full path to the bridge JavaScript file
func (a *App) BridgeScriptPath() string {
	return a.staticPath + "/js/forge-bridge.js"
}

// AlpineBridgeScriptPath returns the full path to the Alpine bridge JavaScript file
func (a *App) AlpineBridgeScriptPath() string {
	return a.staticPath + "/js/alpine-bridge.js"
}

// HotReloadPath returns the full hot reload SSE endpoint path
func (a *App) HotReloadPath() string {
	if a.config.BasePath != "" {
		return a.config.BasePath + "/_forgeui/reload"
	}
	return "/_forgeui/reload"
}

// BridgeScripts returns properly configured bridge script tags as a templ.Component.
// This respects the BasePath configuration.
func (a *App) BridgeScripts(includeAlpine bool, csrfToken ...string) templ.Component {
	if !a.HasBridge() {
		return templ.NopComponent
	}

	token := ""
	if len(csrfToken) > 0 {
		token = csrfToken[0]
	}

	return bridge.BridgeScriptsExternal(bridge.ScriptConfig{
		Endpoint:      a.BridgeCallPath(),
		CSRFToken:     token,
		IncludeAlpine: includeAlpine,
		StaticPath:    a.staticPath,
	})
}

// bridgeScriptHandler returns an http.Handler that serves embedded bridge client scripts
func (a *App) bridgeScriptHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the script name from the URL path
		// e.g., "/api/identity/ui/static/js/forge-bridge.js" -> "forge-bridge.js"
		path := strings.TrimPrefix(r.URL.Path, a.staticPath+"/js/")

		var content string
		var contentType string

		switch path {
		case "forge-bridge.js":
			content = bridge.GetBridgeJS()
			contentType = "text/javascript; charset=utf-8"
		case "alpine-bridge.js":
			content = bridge.GetAlpineJS()
			contentType = "text/javascript; charset=utf-8"
		default:
			http.NotFound(w, r)
			return
		}

		// Set appropriate headers
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=3600")

		// Serve the embedded script
		_, _ = w.Write([]byte(content))
	})
}

// Theme returns the light theme
func (a *App) Theme() *theme.Theme {
	return a.lightTheme
}

// LightTheme returns the light theme
func (a *App) LightTheme() *theme.Theme {
	return a.lightTheme
}

// DarkTheme returns the dark theme
func (a *App) DarkTheme() *theme.Theme {
	return a.darkTheme
}

// RegisterLayout registers a named layout
func (a *App) RegisterLayout(name string, fn router.LayoutFunc, opts ...router.LayoutOption) {
	a.router.RegisterLayout(name, fn, opts...)
}

// Initialize prepares the application (plugins, router, assets, etc.)
// This will be expanded in later phases as more systems are added
func (a *App) Initialize(ctx context.Context) error {
	// Auto-build CSS if themes are configured
	if a.lightTheme != nil && a.darkTheme != nil {
		if err := a.buildThemeCSS(ctx); err != nil {
			// Non-fatal: fall back to CDN mode
			fmt.Printf("[ForgeUI] CSS build skipped: %v (using CDN fallback)\n", err)
		}
	}

	// Initialize asset manager for production
	if !a.IsDev() && a.config.AssetManifest == "" {
		// In production mode without a manifest, pre-generate all fingerprints
		_ = a.Assets.FingerprintAll()
	}

	return nil
}

// buildThemeCSS generates input CSS from themes and compiles it with Tailwind CLI.
// If Tailwind CLI is not available, falls back gracefully (CDN mode).
func (a *App) buildThemeCSS(ctx context.Context) error {
	outputDir := a.config.AssetOutputDir

	// Ensure output directory exists
	cssDir := filepath.Join(outputDir, "css")
	if err := os.MkdirAll(cssDir, 0755); err != nil {
		return fmt.Errorf("failed to create CSS output directory: %w", err)
	}

	// Determine input CSS path
	inputCSS := a.config.InputCSS
	if inputCSS == "" {
		// Generate input CSS from themes
		content := theme.GenerateInputCSS(*a.lightTheme, *a.darkTheme)
		inputPath := filepath.Join(outputDir, "_forgeui_input.css")

		if err := os.WriteFile(inputPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write input CSS: %w", err)
		}

		inputCSS = inputPath

		defer func() { _ = os.Remove(inputPath) }()
	}

	// Create and configure Tailwind processor
	tp := assets.NewTailwindProcessor().
		WithInputCSS(inputCSS).
		WithThemes(a.lightTheme, a.darkTheme).
		WithOutputCSS("css/app.css")

	if a.config.Verbose {
		tp.WithVerbose(true)
	}

	// Process CSS
	cfg := assets.ProcessorConfig{
		InputDir:  a.config.AssetPublicDir,
		OutputDir: outputDir,
		IsDev:     a.IsDev(),
		Minify:    !a.IsDev(),
	}

	if err := tp.Process(ctx, cfg); err != nil {
		return err
	}

	a.cssBuilt = true

	return nil
}

// StaticPath returns the full URL path prefix for static assets.
func (a *App) StaticPath() string {
	return a.staticPath
}

// CSSPath returns the URL path to the compiled CSS stylesheet.
// This respects the BasePath and StaticPath configuration.
func (a *App) CSSPath() string {
	return a.staticPath + "/css/app.css"
}

// FontPreloadLinks returns a templ.Component that renders <link rel="preload">
// tags for all configured fonts that have Preload: true.
// Returns a nop component if no font config is set.
func (a *App) FontPreloadLinks() templ.Component {
	if a.config.FontConfig == nil {
		return templ.NopComponent
	}
	return theme.FontPreloadLinksFromConfig(*a.config.FontConfig)
}

// FontFaceCSS returns the @font-face CSS rules for all configured fonts.
// Returns an empty string if no font config is set.
func (a *App) FontFaceCSS() string {
	if a.config.FontConfig == nil {
		return ""
	}
	var b strings.Builder
	for _, f := range []theme.Font{a.config.FontConfig.Sans, a.config.FontConfig.Serif, a.config.FontConfig.Mono} {
		if css := theme.GenerateFontFaceCSS(f); css != "" {
			b.WriteString(css)
		}
	}
	return b.String()
}

// IsCDNMode returns true if compiled CSS was not built (e.g., Tailwind CLI not available).
// When in CDN mode, layouts should include the Tailwind CDN script and inline theme styles.
func (a *App) IsCDNMode() bool {
	return !a.cssBuilt
}

// ThemeStylesheet returns a <link rel="stylesheet"> templ.Component for the compiled CSS.
// Returns a nop component if CSS was not compiled (CDN mode).
func (a *App) ThemeStylesheet() templ.Component {
	if a.IsCDNMode() {
		return templ.NopComponent
	}

	return a.Assets.StyleSheet("css/app.css")
}

// Handler returns an http.Handler that serves the entire application
// This includes static assets, bridge endpoints, and routed pages
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()

	// Build bridge path with base path if provided
	bridgeCallPath := "/api/bridge/call"
	bridgeStreamPath := "/api/bridge/stream/"
	hotReloadPath := "/_forgeui/reload"

	if a.config.BasePath != "" {
		bridgeCallPath = a.config.BasePath + "/bridge/call"
		bridgeStreamPath = a.config.BasePath + "/bridge/stream/"
		hotReloadPath = a.config.BasePath + "/_forgeui/reload"
	}

	// Serve embedded bridge client scripts if bridge is enabled
	if a.HasBridge() {
		mux.Handle(a.staticPath+"/js/forge-bridge.js", a.bridgeScriptHandler())
		mux.Handle(a.staticPath+"/js/alpine-bridge.js", a.bridgeScriptHandler())
	}

	// Serve static assets
	mux.Handle(a.staticPath+"/", a.Assets.Handler())

	// Serve bridge endpoints if enabled
	if a.HasBridge() {
		mux.Handle(bridgeCallPath, a.bridge.Handler())
		mux.Handle(bridgeStreamPath, a.bridge.StreamHandler())
	}

	// Serve SSE endpoint for hot reload in dev mode
	if a.IsDev() {
		if handler := a.Assets.SSEHandler(); handler != nil {
			mux.Handle(hotReloadPath, handler.(http.Handler))
		}
	}

	// Serve all other requests through router.
	// When BasePath is set, mux routes (static, bridge, SSE) already include the
	// basePath prefix and match before this catch-all. Page routes, however, are
	// registered as relative paths (e.g., "/", "/health"). Strip the basePath
	// so the router sees relative paths that match its compiled patterns.
	if a.config.BasePath != "" {
		mux.Handle("/", http.StripPrefix(a.config.BasePath, a.router))
	} else {
		mux.Handle("/", a.router)
	}

	return mux
}

// Use adds global middleware (convenience method)
func (a *App) Use(middleware ...router.Middleware) *App {
	a.router.Use(middleware...)
	return a
}

// Get registers a GET route (convenience method)
func (a *App) Get(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Get(pattern, handler)
}

// Post registers a POST route (convenience method)
func (a *App) Post(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Post(pattern, handler)
}

// Put registers a PUT route (convenience method)
func (a *App) Put(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Put(pattern, handler)
}

// Patch registers a PATCH route (convenience method)
func (a *App) Patch(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Patch(pattern, handler)
}

// Delete registers a DELETE route (convenience method)
func (a *App) Delete(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Delete(pattern, handler)
}

// Page creates a new PageBuilder for fluent page registration
func (a *App) Page(pattern string) *PageBuilder {
	return &PageBuilder{
		app:     a,
		pattern: pattern,
		method:  "GET",
	}
}

// Group creates a new route group
func (a *App) Group(prefix string) *router.Group {
	return a.router.Group(prefix)
}
