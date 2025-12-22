package plugin

import (
	"context"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"github.com/xraph/forgeui"
)

// Mock component constructor
func mockComponentConstructor(props any, children ...g.Node) g.Node {
	return html.Div(
		html.Class("mock-component"),
		g.Group(children),
	)
}

func TestComponentPluginBase(t *testing.T) {
	components := map[string]ComponentConstructor{
		"MockComponent": mockComponentConstructor,
	}

	plugin := NewComponentPluginBase(
		PluginInfo{
			Name:    "test-components",
			Version: "1.0.0",
		},
		components,
	)

	if plugin.Name() != "test-components" {
		t.Errorf("expected name 'test-components', got %s", plugin.Name())
	}

	comps := plugin.Components()
	if len(comps) != 1 {
		t.Errorf("expected 1 component, got %d", len(comps))
	}

	if _, ok := comps["MockComponent"]; !ok {
		t.Error("MockComponent not found in components")
	}
}

func TestComponentPluginBaseWithCVA(t *testing.T) {
	components := map[string]ComponentConstructor{
		"Button": mockComponentConstructor,
	}

	cva := map[string]*forgeui.CVA{
		"Button": forgeui.NewCVA("btn").
			Variant("size", map[string][]string{
				"sm": {"btn-sm"},
				"lg": {"btn-lg"},
			}),
	}

	plugin := NewComponentPluginBaseWithCVA(
		PluginInfo{Name: "buttons", Version: "1.0.0"},
		components,
		cva,
	)

	cvaExts := plugin.CVAExtensions()
	if len(cvaExts) != 1 {
		t.Errorf("expected 1 CVA extension, got %d", len(cvaExts))
	}

	if _, ok := cvaExts["Button"]; !ok {
		t.Error("Button CVA not found")
	}
}

func TestComponentPluginAddComponent(t *testing.T) {
	plugin := NewComponentPluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		nil,
	)

	plugin.AddComponent("NewComponent", mockComponentConstructor)

	comps := plugin.Components()
	if len(comps) != 1 {
		t.Errorf("expected 1 component after adding, got %d", len(comps))
	}

	if _, ok := comps["NewComponent"]; !ok {
		t.Error("NewComponent not found after adding")
	}
}

func TestComponentPluginAddCVA(t *testing.T) {
	plugin := NewComponentPluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		nil,
	)

	cva := forgeui.NewCVA("test")
	plugin.AddCVA("TestComponent", cva)

	cvaExts := plugin.CVAExtensions()
	if len(cvaExts) != 1 {
		t.Errorf("expected 1 CVA after adding, got %d", len(cvaExts))
	}

	if _, ok := cvaExts["TestComponent"]; !ok {
		t.Error("TestComponent CVA not found after adding")
	}
}

func TestComponentPluginRegistration(t *testing.T) {
	registry := NewRegistry()

	plugin := NewComponentPluginBase(
		PluginInfo{Name: "components", Version: "1.0.0"},
		map[string]ComponentConstructor{
			"Test": mockComponentConstructor,
		},
	)

	err := registry.Register(plugin)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Should be retrievable as ComponentPlugin
	cp, ok := registry.GetComponentPlugin("components")
	if !ok {
		t.Error("component plugin not found in registry")
	}

	if cp.Name() != "components" {
		t.Errorf("expected name 'components', got %s", cp.Name())
	}
}

func TestCollectComponents(t *testing.T) {
	registry := NewRegistry()

	plugin1 := NewComponentPluginBase(
		PluginInfo{Name: "plugin1", Version: "1.0.0"},
		map[string]ComponentConstructor{
			"Component1": mockComponentConstructor,
			"Component2": mockComponentConstructor,
		},
	)

	plugin2 := NewComponentPluginBase(
		PluginInfo{Name: "plugin2", Version: "1.0.0"},
		map[string]ComponentConstructor{
			"Component3": mockComponentConstructor,
		},
	)

	_ = registry.Register(plugin1)
	_ = registry.Register(plugin2)

	components := registry.CollectComponents()
	if len(components) != 3 {
		t.Errorf("expected 3 components, got %d", len(components))
	}

	if _, ok := components["Component1"]; !ok {
		t.Error("Component1 not found")
	}
	if _, ok := components["Component2"]; !ok {
		t.Error("Component2 not found")
	}
	if _, ok := components["Component3"]; !ok {
		t.Error("Component3 not found")
	}
}

func TestComponentPluginLifecycle(t *testing.T) {
	registry := NewRegistry()
	ctx := context.Background()

	plugin := NewComponentPluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		map[string]ComponentConstructor{
			"Test": mockComponentConstructor,
		},
	)

	_ = registry.Register(plugin)

	err := registry.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	err = registry.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
}

