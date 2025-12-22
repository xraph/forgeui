package plugin

import (
	"context"
	"testing"
)

func TestAlpinePluginBase(t *testing.T) {
	plugin := NewAlpinePluginBase(PluginInfo{
		Name:    "test-alpine",
		Version: "1.0.0",
	})

	if plugin.Name() != "test-alpine" {
		t.Errorf("expected name 'test-alpine', got %s", plugin.Name())
	}

	// Test default implementations
	if len(plugin.Scripts()) != 0 {
		t.Error("expected empty scripts by default")
	}
	if len(plugin.Directives()) != 0 {
		t.Error("expected empty directives by default")
	}
	if len(plugin.Stores()) != 0 {
		t.Error("expected empty stores by default")
	}
	if len(plugin.Magics()) != 0 {
		t.Error("expected empty magics by default")
	}
	if len(plugin.AlpineComponents()) != 0 {
		t.Error("expected empty components by default")
	}
}

type mockAlpinePlugin struct {
	*AlpinePluginBase
}

func (m *mockAlpinePlugin) Scripts() []Script {
	return []Script{
		{Name: "test-lib", URL: "https://example.com/test.js", Priority: 10},
	}
}

func (m *mockAlpinePlugin) Directives() []AlpineDirective {
	return []AlpineDirective{
		{Name: "test", Definition: "(el) => {}"},
	}
}

func (m *mockAlpinePlugin) Stores() []AlpineStore {
	return []AlpineStore{
		{Name: "test", InitialState: map[string]any{"count": 0}},
	}
}

func (m *mockAlpinePlugin) Magics() []AlpineMagic {
	return []AlpineMagic{
		{Name: "test", Definition: "(el) => ({})"},
	}
}

func (m *mockAlpinePlugin) AlpineComponents() []AlpineComponent {
	return []AlpineComponent{
		{Name: "test", Definition: "() => ({})"},
	}
}

func TestAlpinePluginRegistration(t *testing.T) {
	registry := NewRegistry()

	plugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "alpine-test",
			Version: "1.0.0",
		}),
	}

	err := registry.Register(plugin)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Should be retrievable as AlpinePlugin
	ap, ok := registry.GetAlpinePlugin("alpine-test")
	if !ok {
		t.Error("alpine plugin not found in registry")
	}

	if ap.Name() != "alpine-test" {
		t.Errorf("expected name 'alpine-test', got %s", ap.Name())
	}
}

func TestCollectScripts(t *testing.T) {
	registry := NewRegistry()

	plugin1 := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "plugin1",
			Version: "1.0.0",
		}),
	}

	_ = registry.Register(plugin1)

	scripts := registry.CollectScripts()
	if len(scripts) != 1 {
		t.Errorf("expected 1 script, got %d", len(scripts))
	}

	if scripts[0].Name != "test-lib" {
		t.Errorf("expected script name 'test-lib', got %s", scripts[0].Name)
	}
}

type priorityPlugin struct {
	*AlpinePluginBase
	scripts []Script
}

func (p *priorityPlugin) Scripts() []Script {
	return p.scripts
}

func TestCollectScriptsPriority(t *testing.T) {
	registry := NewRegistry()

	p1 := &priorityPlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{Name: "p1", Version: "1.0.0"}),
		scripts:          []Script{{Name: "script1", Priority: 50}},
	}

	p2 := &priorityPlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{Name: "p2", Version: "1.0.0"}),
		scripts:          []Script{{Name: "script2", Priority: 10}},
	}

	_ = registry.Register(p1)
	_ = registry.Register(p2)

	scripts := registry.CollectScripts()
	if len(scripts) != 2 {
		t.Fatalf("expected 2 scripts, got %d", len(scripts))
	}

	// Should be sorted by priority (lower first)
	if scripts[0].Name != "script2" {
		t.Errorf("expected script2 first (priority 10), got %s", scripts[0].Name)
	}
	if scripts[1].Name != "script1" {
		t.Errorf("expected script1 second (priority 50), got %s", scripts[1].Name)
	}
}

func TestCollectDirectives(t *testing.T) {
	registry := NewRegistry()

	plugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "test",
			Version: "1.0.0",
		}),
	}

	_ = registry.Register(plugin)

	directives := registry.CollectDirectives()
	if len(directives) != 1 {
		t.Errorf("expected 1 directive, got %d", len(directives))
	}

	if directives[0].Name != "test" {
		t.Errorf("expected directive name 'test', got %s", directives[0].Name)
	}
}

func TestCollectStores(t *testing.T) {
	registry := NewRegistry()

	plugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "test",
			Version: "1.0.0",
		}),
	}

	_ = registry.Register(plugin)

	stores := registry.CollectStores()
	if len(stores) != 1 {
		t.Errorf("expected 1 store, got %d", len(stores))
	}

	if stores[0].Name != "test" {
		t.Errorf("expected store name 'test', got %s", stores[0].Name)
	}

	if stores[0].InitialState["count"] != 0 {
		t.Errorf("expected count 0, got %v", stores[0].InitialState["count"])
	}
}

func TestCollectMagics(t *testing.T) {
	registry := NewRegistry()

	plugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "test",
			Version: "1.0.0",
		}),
	}

	_ = registry.Register(plugin)

	magics := registry.CollectMagics()
	if len(magics) != 1 {
		t.Errorf("expected 1 magic, got %d", len(magics))
	}

	if magics[0].Name != "test" {
		t.Errorf("expected magic name 'test', got %s", magics[0].Name)
	}
}

func TestCollectAlpineComponents(t *testing.T) {
	registry := NewRegistry()

	plugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "test",
			Version: "1.0.0",
		}),
	}

	_ = registry.Register(plugin)

	components := registry.CollectAlpineComponents()
	if len(components) != 1 {
		t.Errorf("expected 1 component, got %d", len(components))
	}

	if components[0].Name != "test" {
		t.Errorf("expected component name 'test', got %s", components[0].Name)
	}
}

func TestAlpinePluginLifecycle(t *testing.T) {
	registry := NewRegistry()
	ctx := context.Background()

	plugin := &mockAlpinePlugin{
		AlpinePluginBase: NewAlpinePluginBase(PluginInfo{
			Name:    "test",
			Version: "1.0.0",
		}),
	}

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

func TestScriptStruct(t *testing.T) {
	script := Script{
		Name:        "test",
		URL:         "https://example.com/test.js",
		Priority:    10,
		Defer:       true,
		Async:       false,
		Module:      true,
		Integrity:   "sha384-...",
		Crossorigin: "anonymous",
	}

	if script.Name != "test" {
		t.Errorf("expected name 'test', got %s", script.Name)
	}
	if script.Priority != 10 {
		t.Errorf("expected priority 10, got %d", script.Priority)
	}
	if !script.Defer {
		t.Error("expected defer true")
	}
	if !script.Module {
		t.Error("expected module true")
	}
}

