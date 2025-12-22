// Package htmxplugin provides an HTMX plugin wrapper for ForgeUI's plugin system.
//
// This allows HTMX to be registered as a plugin, enabling other plugins
// to declare it as a dependency and ensuring proper initialization order.
//
// # Basic Usage
//
//	registry := plugin.NewRegistry()
//	registry.Use(htmxplugin.New())
package htmxplugin

import (
	"context"

	"github.com/xraph/forgeui/plugin"
)

// HTMXPlugin wraps HTMX functionality as a ForgeUI plugin.
type HTMXPlugin struct {
	*plugin.PluginBase
	version    string
	extensions []string
}

// New creates a new HTMX plugin.
func New(extensions ...string) *HTMXPlugin {
	return &HTMXPlugin{
		PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
			Name:        "htmx",
			Version:     "1.0.0",
			Description: "HTMX integration plugin wrapper",
			Author:      "ForgeUI",
			License:     "MIT",
		}),
		version:    "2.0.3",
		extensions: extensions,
	}
}

// Init initializes the HTMX plugin.
func (h *HTMXPlugin) Init(ctx context.Context, registry *plugin.Registry) error {
	return nil
}

// Shutdown cleanly shuts down the plugin.
func (h *HTMXPlugin) Shutdown(ctx context.Context) error {
	return nil
}

// Scripts returns HTMX library and extension scripts.
func (h *HTMXPlugin) Scripts() []plugin.Script {
	scripts := []plugin.Script{
		{
			Name:     "htmx",
			URL:      "https://unpkg.com/htmx.org@" + h.version,
			Priority: 50,
			Defer:    false,
		},
	}

	// Add extension scripts
	for _, ext := range h.extensions {
		scripts = append(scripts, plugin.Script{
			Name:     "htmx-ext-" + ext,
			URL:      "https://unpkg.com/htmx-ext-" + ext + "@2.0.0/ext/" + ext + ".js",
			Priority: 51,
			Defer:    false,
		})
	}

	return scripts
}

// Directives returns HTMX directives (none, as HTMX uses attributes directly).
func (h *HTMXPlugin) Directives() []plugin.AlpineDirective {
	return nil
}

// Stores returns HTMX stores (none needed).
func (h *HTMXPlugin) Stores() []plugin.AlpineStore {
	return nil
}

// Magics returns HTMX magic properties (none).
func (h *HTMXPlugin) Magics() []plugin.AlpineMagic {
	return nil
}

// AlpineComponents returns HTMX Alpine components (none).
func (h *HTMXPlugin) AlpineComponents() []plugin.AlpineComponent {
	return nil
}

