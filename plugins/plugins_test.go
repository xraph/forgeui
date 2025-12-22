package plugins

import (
	"context"
	"testing"

	"github.com/xraph/forgeui/plugin"
	"github.com/xraph/forgeui/plugins/analytics"
)

func TestAllPlugins(t *testing.T) {
	plugins := AllPlugins()
	if len(plugins) == 0 {
		t.Error("AllPlugins should return at least one plugin")
	}

	// Should return 7 plugins
	if len(plugins) != 7 {
		t.Errorf("Expected 7 plugins, got %d", len(plugins))
	}
}

func TestEssentialPlugins(t *testing.T) {
	plugins := EssentialPlugins()
	if len(plugins) == 0 {
		t.Error("EssentialPlugins should return at least one plugin")
	}

	// Should return 3 plugins
	if len(plugins) != 3 {
		t.Errorf("Expected 3 plugins, got %d", len(plugins))
	}
}

func TestDataVisualizationPlugins(t *testing.T) {
	plugins := DataVisualizationPlugins()
	if len(plugins) == 0 {
		t.Error("DataVisualizationPlugins should return at least one plugin")
	}

	// Should return 2 plugins
	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(plugins))
	}
}

func TestPluginInitialization(t *testing.T) {
	ctx := context.Background()
	registry := plugin.NewRegistry()

	// Register essential plugins
	for _, p := range EssentialPlugins() {
		if pluginType, ok := p.(plugin.Plugin); ok {
			if err := registry.Register(pluginType); err != nil {
				t.Fatalf("Failed to register plugin: %v", err)
			}
		}
	}

	// Initialize all plugins
	if err := registry.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize plugins: %v", err)
	}

	// Verify plugins are registered
	if registry.Count() != 3 {
		t.Errorf("Expected 3 registered plugins, got %d", registry.Count())
	}

	// Shutdown plugins
	if err := registry.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to shutdown plugins: %v", err)
	}
}

func TestPluginConstructors(t *testing.T) {
	// Test that all constructor functions work
	tests := []struct {
		name string
		fn   func() plugin.Plugin
	}{
		{"Toast", func() plugin.Plugin { return NewToast() }},
		{"Sortable", func() plugin.Plugin { return NewSortable() }},
		{"Charts", func() plugin.Plugin { return NewCharts() }},
		{"Analytics", func() plugin.Plugin { return NewAnalytics(analytics.DefaultConfig()) }},
		{"SEO", func() plugin.Plugin { return NewSEO() }},
		{"HTMX", func() plugin.Plugin { return NewHTMX() }},
		{"CorporateTheme", func() plugin.Plugin { return NewCorporateTheme() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fn()
			if p == nil {
				t.Errorf("%s constructor returned nil", tt.name)
			}

			// Verify plugin has a name
			if p.Name() == "" {
				t.Errorf("%s plugin has empty name", tt.name)
			}

			// Verify plugin has a version
			if p.Version() == "" {
				t.Errorf("%s plugin has empty version", tt.name)
			}
		})
	}
}

func TestPluginIntegration(t *testing.T) {
	ctx := context.Background()
	registry := plugin.NewRegistry()

	// Register all plugins
	for _, p := range AllPlugins() {
		if pluginType, ok := p.(plugin.Plugin); ok {
			if err := registry.Register(pluginType); err != nil {
				t.Fatalf("Failed to register plugin: %v", err)
			}
		}
	}

	// Initialize
	if err := registry.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Verify all plugins are registered
	if registry.Count() != 7 {
		t.Errorf("Expected 7 plugins, got %d", registry.Count())
	}

	// Collect scripts from all plugins
	scripts := registry.CollectScripts()
	if len(scripts) == 0 {
		t.Error("Expected scripts from plugins")
	}

	// Shutdown
	if err := registry.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to shutdown: %v", err)
	}
}
