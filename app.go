package forgeui

import (
	"context"
	"net/http"

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
}

// New creates a new ForgeUI application with enhanced initialization
func New(opts ...AppOption) *App {
	config := DefaultAppConfig()
	for _, opt := range opts {
		opt(config)
	}

	// Initialize asset manager
	assetManager := assets.NewManager(assets.Config{
		PublicDir:  config.AssetPublicDir,
		OutputDir:  config.AssetOutputDir,
		StaticPath: config.StaticPath,
		IsDev:      config.Debug,
		Manifest:   config.AssetManifest,
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
		staticPath: config.StaticPath,
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
	// Will be expanded in later phases:
	// - Phase 11: Initialize plugin registry

	// Phase 15: Initialize asset manager
	if !a.IsDev() && a.config.AssetManifest == "" {
		// In production mode without a manifest, pre-generate all fingerprints
		_ = a.Assets.FingerprintAll()
	}

	return nil
}

// Handler returns an http.Handler that serves the entire application
// This includes static assets, bridge endpoints, and routed pages
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()

	// Serve static assets
	mux.Handle(a.staticPath+"/", a.Assets.Handler())

	// Serve bridge endpoints if enabled
	if a.HasBridge() {
		mux.Handle("/api/bridge/call", a.bridge.Handler())
		mux.Handle("/api/bridge/stream/", a.bridge.StreamHandler())
	}

	// Serve SSE endpoint for hot reload in dev mode
	if a.IsDev() {
		if handler := a.Assets.SSEHandler(); handler != nil {
			mux.Handle("/_forgeui/reload", handler.(http.Handler))
		}
	}

	// Serve all other requests through router
	mux.Handle("/", a.router)

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
