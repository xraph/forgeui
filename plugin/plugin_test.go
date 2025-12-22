package plugin

import (
	"context"
	"testing"
)

// mockPlugin is a test plugin implementation
type mockPlugin struct {
	*PluginBase
	initCalled     bool
	shutdownCalled bool
	initError      error
	shutdownError  error
}

func newMockPlugin(name, version string, deps []Dependency) *mockPlugin {
	return &mockPlugin{
		PluginBase: NewPluginBase(PluginInfo{
			Name:         name,
			Version:      version,
			Description:  "Test plugin",
			Dependencies: deps,
		}),
	}
}

func (m *mockPlugin) Init(ctx context.Context, r *Registry) error {
	m.initCalled = true
	return m.initError
}

func (m *mockPlugin) Shutdown(ctx context.Context) error {
	m.shutdownCalled = true
	return m.shutdownError
}

func TestPluginBase(t *testing.T) {
	info := PluginInfo{
		Name:        "test-plugin",
		Version:     "1.0.0",
		Description: "Test plugin",
		Author:      "Test Author",
		License:     "MIT",
		Homepage:    "https://example.com",
		Repository:  "https://github.com/test/plugin",
		Tags:        []string{"test", "example"},
		Dependencies: []Dependency{
			{Name: "dep1", Version: ">=1.0.0"},
		},
	}

	base := NewPluginBase(info)

	if base.Name() != "test-plugin" {
		t.Errorf("expected name 'test-plugin', got %s", base.Name())
	}

	if base.Version() != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", base.Version())
	}

	if base.Description() != "Test plugin" {
		t.Errorf("expected description 'Test plugin', got %s", base.Description())
	}

	deps := base.Dependencies()
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(deps))
	}

	// Test default Init and Shutdown
	ctx := context.Background()
	if err := base.Init(ctx, nil); err != nil {
		t.Errorf("default Init should not error: %v", err)
	}

	if err := base.Shutdown(ctx); err != nil {
		t.Errorf("default Shutdown should not error: %v", err)
	}
}

func TestDependency(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		version    string
		optional   bool
		want       bool
	}{
		{"empty constraint", "", "1.0.0", false, true},
		{"wildcard", "*", "1.0.0", false, true},
		{"exact match", "=1.0.0", "1.0.0", false, true},
		{"exact mismatch", "=1.0.0", "1.0.1", false, false},
		{"greater than", ">1.0.0", "1.0.1", false, true},
		{"greater than fail", ">1.0.0", "1.0.0", false, false},
		{"greater or equal", ">=1.0.0", "1.0.0", false, true},
		{"greater or equal pass", ">=1.0.0", "1.0.1", false, true},
		{"less than", "<2.0.0", "1.0.0", false, true},
		{"less than fail", "<2.0.0", "2.0.0", false, false},
		{"less or equal", "<=2.0.0", "2.0.0", false, true},
		{"tilde patch", "~1.2.3", "1.2.4", false, true},
		{"tilde minor fail", "~1.2.3", "1.3.0", false, false},
		{"caret minor", "^1.2.3", "1.3.0", false, true},
		{"caret major fail", "^1.2.3", "2.0.0", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := Dependency{
				Name:     "test",
				Version:  tt.constraint,
				Optional: tt.optional,
			}

			got := dep.Satisfies(tt.version)
			if got != tt.want {
				t.Errorf("Satisfies(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

