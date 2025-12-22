package plugin

import (
	"context"
	"net/http"
	"testing"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/theme"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Multi-type plugin that implements multiple interfaces
type multiPlugin struct {
	*PluginBase

	componentConstructors map[string]ComponentConstructor
	alpineScripts         []Script
	themes                map[string]theme.Theme
	middleware            func(http.Handler) http.Handler
	priority              int
}

func newMultiPlugin() *multiPlugin {
	return &multiPlugin{
		PluginBase: NewPluginBase(PluginInfo{
			Name:    "multi-plugin",
			Version: "1.0.0",
		}),
		componentConstructors: map[string]ComponentConstructor{
			"MultiComponent": func(props any, children ...g.Node) g.Node {
				return html.Div(g.Group(children))
			},
		},
		alpineScripts: []Script{
			{Name: "multi-script", URL: "https://example.com/multi.js", Priority: 20},
		},
		themes: map[string]theme.Theme{
			"multi": theme.DefaultLight(),
		},
		middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Multi", "true")
				next.ServeHTTP(w, r)
			})
		},
		priority: 25,
	}
}

func (m *multiPlugin) Components() map[string]ComponentConstructor {
	return m.componentConstructors
}

func (m *multiPlugin) CVAExtensions() map[string]*forgeui.CVA {
	return nil
}

func (m *multiPlugin) Scripts() []Script {
	return m.alpineScripts
}

func (m *multiPlugin) Directives() []AlpineDirective {
	return nil
}

func (m *multiPlugin) Stores() []AlpineStore {
	return nil
}

func (m *multiPlugin) Magics() []AlpineMagic {
	return nil
}

func (m *multiPlugin) AlpineComponents() []AlpineComponent {
	return nil
}

func (m *multiPlugin) Themes() map[string]theme.Theme {
	return m.themes
}

func (m *multiPlugin) DefaultTheme() string {
	return "multi"
}

func (m *multiPlugin) CSS() string {
	return ""
}

func (m *multiPlugin) Fonts() []theme.Font {
	return nil
}

func (m *multiPlugin) Middleware() func(http.Handler) http.Handler {
	return m.middleware
}

func (m *multiPlugin) Priority() int {
	return m.priority
}

func TestMultiTypePlugin(t *testing.T) {
	registry := NewRegistry()
	plugin := newMultiPlugin()

	err := registry.Register(plugin)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Should be registered in all type-specific registries
	if _, ok := registry.GetComponentPlugin("multi-plugin"); !ok {
		t.Error("plugin not found in component registry")
	}

	if _, ok := registry.GetAlpinePlugin("multi-plugin"); !ok {
		t.Error("plugin not found in alpine registry")
	}

	if _, ok := registry.GetThemePlugin("multi-plugin"); !ok {
		t.Error("plugin not found in theme registry")
	}

	// Should appear in all collection methods
	components := registry.CollectComponents()
	if len(components) != 1 {
		t.Errorf("expected 1 component, got %d", len(components))
	}

	scripts := registry.CollectScripts()
	if len(scripts) != 1 {
		t.Errorf("expected 1 script, got %d", len(scripts))
	}

	middleware := registry.CollectMiddleware()
	if len(middleware) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(middleware))
	}
}

func TestComplexPluginEcosystem(t *testing.T) {
	registry := NewRegistry()
	ctx := context.Background()

	// Register multiple plugins of different types
	compPlugin := NewComponentPluginBase(
		PluginInfo{Name: "components", Version: "1.0.0"},
		map[string]ComponentConstructor{
			"Comp1": func(props any, children ...g.Node) g.Node {
				return html.Div()
			},
		},
	)

	alpinePlugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "alpine",
			Version: "1.0.0",
		}),
	}

	themePlugin := NewThemePluginBase(
		PluginInfo{Name: "themes", Version: "1.0.0"},
		map[string]theme.Theme{"custom": theme.DefaultLight()},
		"custom",
	)

	mwPlugin := NewMiddlewarePluginBase(
		PluginInfo{Name: "middleware", Version: "1.0.0"},
		func(next http.Handler) http.Handler { return next },
		10,
	)

	multiP := newMultiPlugin()

	_ = registry.Register(compPlugin)
	_ = registry.Register(alpinePlugin)
	_ = registry.Register(themePlugin)
	_ = registry.Register(mwPlugin)
	_ = registry.Register(multiP)

	// Initialize all plugins
	err := registry.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	// Verify all collections
	if len(registry.CollectComponents()) != 2 {
		t.Errorf("expected 2 components, got %d", len(registry.CollectComponents()))
	}

	if len(registry.CollectScripts()) != 2 {
		t.Errorf("expected 2 scripts, got %d", len(registry.CollectScripts()))
	}

	if len(registry.CollectMiddleware()) != 2 {
		t.Errorf("expected 2 middleware, got %d", len(registry.CollectMiddleware()))
	}

	// Shutdown all plugins
	err = registry.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
}

func TestPluginDependenciesWithTypes(t *testing.T) {
	registry := NewRegistry()

	// Plugin A provides components
	pluginA := NewComponentPluginBase(
		PluginInfo{
			Name:    "plugin-a",
			Version: "1.0.0",
		},
		map[string]ComponentConstructor{
			"CompA": func(props any, children ...g.Node) g.Node {
				return html.Div()
			},
		},
	)

	// Plugin B depends on Plugin A
	pluginB := NewComponentPluginBase(
		PluginInfo{
			Name:    "plugin-b",
			Version: "1.0.0",
			Dependencies: []Dependency{
				{Name: "plugin-a", Version: ">=1.0.0"},
			},
		},
		map[string]ComponentConstructor{
			"CompB": func(props any, children ...g.Node) g.Node {
				return html.Div()
			},
		},
	)

	_ = registry.Register(pluginA)
	_ = registry.Register(pluginB)

	err := registry.ResolveDependencies()
	if err != nil {
		t.Fatalf("ResolveDependencies() error = %v", err)
	}

	// Initialize should work
	ctx := context.Background()

	err = registry.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
}

func TestPluginTypeCoexistence(t *testing.T) {
	registry := NewRegistry()

	// Multiple plugins of the same type should coexist
	comp1 := NewComponentPluginBase(
		PluginInfo{Name: "comp1", Version: "1.0.0"},
		map[string]ComponentConstructor{
			"C1": func(props any, children ...g.Node) g.Node { return html.Div() },
		},
	)

	comp2 := NewComponentPluginBase(
		PluginInfo{Name: "comp2", Version: "1.0.0"},
		map[string]ComponentConstructor{
			"C2": func(props any, children ...g.Node) g.Node { return html.Div() },
		},
	)

	_ = registry.Register(comp1)
	_ = registry.Register(comp2)

	components := registry.CollectComponents()
	if len(components) != 2 {
		t.Errorf("expected 2 components, got %d", len(components))
	}

	if _, ok := components["C1"]; !ok {
		t.Error("C1 not found")
	}

	if _, ok := components["C2"]; !ok {
		t.Error("C2 not found")
	}
}

func TestAssetCollectionImmutability(t *testing.T) {
	registry := NewRegistry()

	plugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "test",
			Version: "1.0.0",
		}),
	}

	_ = registry.Register(plugin)

	// Collect scripts
	scripts1 := registry.CollectScripts()
	scripts2 := registry.CollectScripts()

	// Should be separate slices
	if len(scripts1) != len(scripts2) {
		t.Error("script collections have different lengths")
	}

	// Modifying one shouldn't affect the other
	if len(scripts1) > 0 {
		scripts1[0].Name = "modified"
		if scripts2[0].Name == "modified" {
			t.Error("modifying collected scripts affected subsequent collections")
		}
	}
}
