package forgeui

import (
	"net/http"

	"github.com/xraph/forgeui/assets"
	"github.com/xraph/forgeui/bridge"
	"github.com/xraph/forgeui/router"
	"github.com/xraph/forgeui/theme"
)

// AppConfig holds enhanced ForgeUI application configuration
type AppConfig struct {
	// Debug enables debug/development mode
	Debug bool

	// AssetPublicDir is the source directory for static assets
	AssetPublicDir string

	// AssetOutputDir is the output directory for processed assets
	AssetOutputDir string

	// AssetManifest is the path to a manifest file for production builds
	AssetManifest string

	// Bridge configuration (optional)
	BridgeConfig *bridge.Config
	EnableBridge bool

	// Theme configuration
	LightTheme *theme.Theme
	DarkTheme  *theme.Theme

	// DefaultLayout is the default layout name for all pages
	DefaultLayout string

	// StaticPath is the URL path for static assets
	StaticPath string
}

// DefaultAppConfig returns sensible defaults for ForgeUI application
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		Debug:          false,
		AssetPublicDir: "public",
		AssetOutputDir: "dist",
		AssetManifest:  "",
		EnableBridge:   false,
		BridgeConfig:   nil,
		LightTheme:     nil,
		DarkTheme:      nil,
		DefaultLayout:  "",
		StaticPath:     "/static",
	}
}

// AppOption is a functional option for configuring the App
type AppOption func(*AppConfig)

// WithDev enables or disables development mode
func WithDev(dev bool) AppOption {
	return func(c *AppConfig) { c.Debug = dev }
}

// WithAssets sets the asset directories
func WithAssets(publicDir string, opts ...string) AppOption {
	return func(c *AppConfig) {
		c.AssetPublicDir = publicDir
		if len(opts) > 0 {
			c.AssetOutputDir = opts[0]
		}
	}
}

// WithBridge enables and configures the bridge system
func WithBridge(opts ...bridge.ConfigOption) AppOption {
	return func(c *AppConfig) {
		c.EnableBridge = true
		bridgeConfig := bridge.DefaultConfig()
		for _, opt := range opts {
			opt(bridgeConfig)
		}
		c.BridgeConfig = bridgeConfig
	}
}

// WithThemes sets the light and dark themes
func WithThemes(light, dark *theme.Theme) AppOption {
	return func(c *AppConfig) {
		c.LightTheme = light
		c.DarkTheme = dark
	}
}

// WithDefaultLayout sets the default layout for all pages
func WithDefaultLayout(layout string) AppOption {
	return func(c *AppConfig) {
		c.DefaultLayout = layout
	}
}

// WithAppStaticPath sets the URL path for static assets
func WithAppStaticPath(path string) AppOption {
	return func(c *AppConfig) {
		c.StaticPath = path
	}
}

// EnhancedApp is the enhanced ForgeUI application with modern API
type EnhancedApp struct {
	config     *AppConfig
	Assets     *assets.Manager
	router     *router.Router
	bridge     *bridge.Bridge
	lightTheme *theme.Theme
	darkTheme  *theme.Theme
	staticPath string
}

// NewApp creates a new ForgeUI application with enhanced initialization
func NewApp(opts ...AppOption) *EnhancedApp {
	config := DefaultAppConfig()
	for _, opt := range opts {
		opt(config)
	}

	// Initialize asset manager
	assetManager := assets.NewManager(assets.Config{
		PublicDir: config.AssetPublicDir,
		OutputDir: config.AssetOutputDir,
		IsDev:     config.Debug,
		Manifest:  config.AssetManifest,
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

	app := &EnhancedApp{
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

// IsDev returns true if the application is in debug/development mode
func (a *EnhancedApp) IsDev() bool {
	return a.config.Debug
}

// Router returns the application's HTTP router
func (a *EnhancedApp) Router() *router.Router {
	return a.router
}

// Bridge returns the bridge system (may be nil if not enabled)
func (a *EnhancedApp) Bridge() *bridge.Bridge {
	return a.bridge
}

// HasBridge returns true if bridge system is enabled
func (a *EnhancedApp) HasBridge() bool {
	return a.bridge != nil
}

// Theme returns the light theme
func (a *EnhancedApp) Theme() *theme.Theme {
	return a.lightTheme
}

// LightTheme returns the light theme
func (a *EnhancedApp) LightTheme() *theme.Theme {
	return a.lightTheme
}

// DarkTheme returns the dark theme
func (a *EnhancedApp) DarkTheme() *theme.Theme {
	return a.darkTheme
}

// RegisterLayout registers a named layout
func (a *EnhancedApp) RegisterLayout(name string, fn router.LayoutFunc, opts ...router.LayoutOption) {
	a.router.RegisterLayout(name, fn, opts...)
}

// Handler returns an http.Handler that serves the entire application
// This includes static assets, bridge endpoints, and routed pages
func (a *EnhancedApp) Handler() http.Handler {
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
func (a *EnhancedApp) Use(middleware ...router.Middleware) *EnhancedApp {
	a.router.Use(middleware...)
	return a
}

// Get registers a GET route (convenience method)
func (a *EnhancedApp) Get(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Get(pattern, handler)
}

// Post registers a POST route (convenience method)
func (a *EnhancedApp) Post(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Post(pattern, handler)
}

// Put registers a PUT route (convenience method)
func (a *EnhancedApp) Put(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Put(pattern, handler)
}

// Patch registers a PATCH route (convenience method)
func (a *EnhancedApp) Patch(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Patch(pattern, handler)
}

// Delete registers a DELETE route (convenience method)
func (a *EnhancedApp) Delete(pattern string, handler router.PageHandler) *router.Route {
	return a.router.Delete(pattern, handler)
}

// Page creates a new PageBuilder for fluent page registration
func (a *EnhancedApp) Page(pattern string) *PageBuilder {
	return &PageBuilder{
		app:     a,
		pattern: pattern,
		method:  "GET",
	}
}

// Group creates a new route group
func (a *EnhancedApp) Group(prefix string) *router.Group {
	return a.router.Group(prefix)
}
