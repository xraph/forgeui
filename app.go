package forgeui

import (
	"context"

	"github.com/xraph/forgeui/assets"
	"github.com/xraph/forgeui/router"
)

// App is the main ForgeUI application
// It holds configuration and will be extended with registry, router, and assets in later phases
type App struct {
	config *Config
	Assets *assets.Manager
	router *router.Router
	// registry will be added in Phase 11
}

// New creates a new ForgeUI application with the provided configuration options
func New(opts ...ConfigOption) *App {
	config := DefaultConfig()
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

	return &App{
		config: config,
		Assets: assetManager,
		router: router.New(),
	}
}

// Config returns the application configuration
func (a *App) Config() *Config {
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

// Use adds global middleware (convenience method)
func (a *App) Use(middleware ...router.Middleware) *App {
	a.router.Use(middleware...)
	return a
}
