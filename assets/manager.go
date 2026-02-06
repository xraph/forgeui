package assets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"sync"
)

// Manager handles static assets with caching and fingerprinting
type Manager struct {
	publicDir    string
	outputDir    string
	staticPath   string
	fingerprints map[string]string
	mu           sync.RWMutex
	isDev        bool
	manifest     map[string]string
	pipeline     *Pipeline
	devServer    *DevServer
	fileSystem   fs.FS // Filesystem abstraction for serving files
}

// Config defines configuration options for asset management
type Config struct {
	// PublicDir is the source directory for static assets (e.g., "public")
	PublicDir string

	// OutputDir is the output directory for processed assets (e.g., "dist")
	OutputDir string

	// StaticPath is the URL path prefix for static assets (e.g., "/static")
	StaticPath string

	// IsDev enables development mode (no fingerprinting)
	IsDev bool

	// Manifest is the path to a manifest file for production builds
	Manifest string

	// FileSystem is an optional custom filesystem (e.g., embed.FS)
	// If nil, os.DirFS(PublicDir) will be used
	FileSystem fs.FS
}

// NewManager creates a new asset manager with the given configuration
func NewManager(cfg Config) *Manager {
	if cfg.PublicDir == "" {
		cfg.PublicDir = "public"
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = "dist"
	}

	// Default static path to "/static" if not provided
	staticPath := cfg.StaticPath
	if staticPath == "" {
		staticPath = "/static"
	}

	// Ensure leading slash
	if staticPath[0] != '/' {
		staticPath = "/" + staticPath
	}

	// Ensure trailing slash for consistent URL building
	if staticPath[len(staticPath)-1] != '/' {
		staticPath += "/"
	}

	// Use custom filesystem if provided, otherwise default to os.DirFS
	fileSystem := cfg.FileSystem
	if fileSystem == nil {
		fileSystem = os.DirFS(cfg.PublicDir)
	}

	m := &Manager{
		publicDir:    cfg.PublicDir,
		outputDir:    cfg.OutputDir,
		staticPath:   staticPath,
		fingerprints: make(map[string]string),
		isDev:        cfg.IsDev,
		manifest:     make(map[string]string),
		fileSystem:   fileSystem,
	}

	// Load manifest if exists
	if cfg.Manifest != "" {
		_ = m.loadManifest(cfg.Manifest)
	}

	return m
}

// URL returns the URL for an asset, with fingerprint in production
func (m *Manager) URL(path string) string {
	if m.isDev {
		return m.staticPath + path
	}

	// Check manifest first
	m.mu.RLock()

	if fp, ok := m.manifest[path]; ok {
		m.mu.RUnlock()
		return m.staticPath + fp
	}

	// Check cached fingerprints
	if fp, ok := m.fingerprints[path]; ok {
		m.mu.RUnlock()
		return m.staticPath + fp
	}

	m.mu.RUnlock()

	// Generate fingerprint
	fp := m.fingerprint(path)
	m.mu.Lock()
	m.fingerprints[path] = fp
	m.mu.Unlock()

	return m.staticPath + fp
}

// IsDev returns whether the manager is in development mode
func (m *Manager) IsDev() bool {
	return m.isDev
}

// PublicDir returns the configured public directory
func (m *Manager) PublicDir() string {
	return m.publicDir
}

// SetFileSystem allows setting a custom filesystem for serving assets.
// This is useful when using embed.FS or other fs.FS implementations.
func (m *Manager) SetFileSystem(fsys fs.FS) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fileSystem = fsys
}

// loadManifest loads asset mappings from a manifest file
func (m *Manager) loadManifest(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var manifest map[string]string
	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	m.mu.Lock()
	m.manifest = manifest
	m.mu.Unlock()

	return nil
}

// SaveManifest writes the current fingerprint mappings to a manifest file
func (m *Manager) SaveManifest(path string) error {
	m.mu.RLock()

	data := make(map[string]string)
	maps.Copy(data, m.fingerprints)

	m.mu.RUnlock()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0600)
}

// Pipeline returns the asset pipeline for this manager.
// Creates a new pipeline if one doesn't exist.
func (m *Manager) Pipeline() *Pipeline {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pipeline == nil {
		m.pipeline = NewPipeline(PipelineConfig{
			InputDir:  m.publicDir,
			OutputDir: m.outputDir,
			IsDev:     m.isDev,
		}, m)
	}

	return m.pipeline
}

// Build runs the asset pipeline to process all assets.
// This is typically used for production builds.
func (m *Manager) Build(ctx context.Context) error {
	pipeline := m.Pipeline()

	// Add default processors if none exist
	if pipeline.ProcessorCount() == 0 {
		// Add Tailwind CSS processor
		tailwind := NewTailwindProcessor()
		pipeline.AddProcessor(tailwind)

		// Add ESBuild processor (optional, only if entry points exist)
		esbuild := NewESBuildProcessor()
		pipeline.AddProcessor(esbuild)
	}

	return pipeline.Build(ctx)
}

// StartDevServer starts the development server with hot reload.
// This watches for file changes and automatically rebuilds assets.
func (m *Manager) StartDevServer(ctx context.Context) error {
	// Check if already running (with lock)
	m.mu.Lock()

	if m.devServer != nil {
		m.mu.Unlock()
		return errors.New("dev server already running")
	}

	m.mu.Unlock()

	// Get pipeline (this will acquire its own lock)
	pipeline := m.Pipeline()

	// Add default processors if none exist
	if pipeline.ProcessorCount() == 0 {
		// Add Tailwind CSS processor
		tailwind := NewTailwindProcessor().WithVerbose(false)
		pipeline.AddProcessor(tailwind)

		// Add ESBuild processor (optional)
		esbuild := NewESBuildProcessor().WithVerbose(false)
		pipeline.AddProcessor(esbuild)
	}

	devServer, err := NewDevServer(pipeline)
	if err != nil {
		return fmt.Errorf("failed to create dev server: %w", err)
	}

	devServer.SetVerbose(false)

	// Set dev server (with lock)
	m.mu.Lock()
	m.devServer = devServer
	m.mu.Unlock()

	return devServer.Start(ctx)
}

// SSEHandler returns the Server-Sent Events handler for hot reload.
// Mount this at /_forgeui/reload in your HTTP server.
// Returns nil if dev server is not running.
func (m *Manager) SSEHandler() any {
	m.mu.RLock()
	ds := m.devServer
	m.mu.RUnlock()

	if ds == nil {
		return nil
	}

	return ds.SSEHandler()
}

// HotReloadScript returns the client-side JavaScript for hot reload.
// Include this in your HTML during development.
func (m *Manager) HotReloadScript() string {
	m.mu.RLock()
	ds := m.devServer
	m.mu.RUnlock()

	if ds == nil {
		return ""
	}

	return ds.HotReloadScript()
}

// StopDevServer stops the development server if running
func (m *Manager) StopDevServer() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.devServer == nil {
		return nil
	}

	err := m.devServer.Close()
	m.devServer = nil

	return err
}
