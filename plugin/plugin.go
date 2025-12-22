// Package plugin provides the ForgeUI plugin system.
//
// Plugins extend ForgeUI functionality through a common interface.
// The registry manages plugin lifecycle, dependencies, and hooks.
//
// Basic usage:
//
//	registry := plugin.NewRegistry()
//	registry.Use(
//	    myplugin.New(),
//	    anotherplugin.New(),
//	)
//
//	if err := registry.Initialize(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer registry.Shutdown(ctx)
package plugin

import "context"

// Plugin is the base interface all plugins must implement.
type Plugin interface {
	// Name returns the unique plugin identifier
	Name() string

	// Version returns the plugin version (semver)
	Version() string

	// Description returns a human-readable description
	Description() string

	// Dependencies returns required plugin dependencies
	Dependencies() []Dependency

	// Init initializes the plugin
	Init(ctx context.Context, registry *Registry) error

	// Shutdown cleanly shuts down the plugin
	Shutdown(ctx context.Context) error
}

// Dependency represents a plugin dependency with version constraints.
type Dependency struct {
	Name     string
	Version  string // semver constraint, e.g., ">=1.0.0"
	Optional bool
}

// PluginInfo contains plugin metadata.
type PluginInfo struct {
	Name         string
	Version      string
	Description  string
	Author       string
	License      string
	Homepage     string
	Repository   string
	Tags         []string
	Dependencies []Dependency
}

// PluginBase provides a base implementation for plugins.
// Embed this in custom plugins to inherit default implementations.
type PluginBase struct {
	info PluginInfo
}

// NewPluginBase creates a new PluginBase with the given info.
func NewPluginBase(info PluginInfo) *PluginBase {
	return &PluginBase{info: info}
}

// Name returns the plugin name.
func (p *PluginBase) Name() string {
	return p.info.Name
}

// Version returns the plugin version.
func (p *PluginBase) Version() string {
	return p.info.Version
}

// Description returns the plugin description.
func (p *PluginBase) Description() string {
	return p.info.Description
}

// Dependencies returns the plugin dependencies.
func (p *PluginBase) Dependencies() []Dependency {
	return p.info.Dependencies
}

// Init is the default implementation (override in plugins).
func (p *PluginBase) Init(ctx context.Context, r *Registry) error {
	return nil
}

// Shutdown is the default implementation (override in plugins).
func (p *PluginBase) Shutdown(ctx context.Context) error {
	return nil
}

