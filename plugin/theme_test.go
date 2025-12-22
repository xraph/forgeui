package plugin

import (
	"context"
	"testing"

	"github.com/xraph/forgeui/theme"
)

func TestThemePluginBase(t *testing.T) {
	themes := map[string]theme.Theme{
		"light": theme.DefaultLight(),
		"dark":  theme.DefaultDark(),
	}

	plugin := NewThemePluginBase(
		PluginInfo{Name: "themes", Version: "1.0.0"},
		themes,
		"light",
	)

	if plugin.Name() != "themes" {
		t.Errorf("expected name 'themes', got %s", plugin.Name())
	}

	th := plugin.Themes()
	if len(th) != 2 {
		t.Errorf("expected 2 themes, got %d", len(th))
	}

	if plugin.DefaultTheme() != "light" {
		t.Errorf("expected default theme 'light', got %s", plugin.DefaultTheme())
	}
}

func TestThemePluginBaseWithFonts(t *testing.T) {
	themes := map[string]theme.Theme{
		"custom": theme.DefaultLight(),
	}

	fonts := []theme.Font{
		{Family: "Inter", Weights: []int{400, 600, 700}},
		{Family: "JetBrains Mono", Weights: []int{400, 500}},
	}

	plugin := NewThemePluginBaseWithFonts(
		PluginInfo{Name: "custom-theme", Version: "1.0.0"},
		themes,
		"custom",
		fonts,
	)

	f := plugin.Fonts()
	if len(f) != 2 {
		t.Errorf("expected 2 fonts, got %d", len(f))
	}

	if f[0].Family != "Inter" {
		t.Errorf("expected first font 'Inter', got %s", f[0].Family)
	}
}

func TestThemePluginAddTheme(t *testing.T) {
	plugin := NewThemePluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		nil,
		"",
	)

	plugin.AddTheme("new-theme", theme.DefaultLight())

	themes := plugin.Themes()
	if len(themes) != 1 {
		t.Errorf("expected 1 theme after adding, got %d", len(themes))
	}

	if _, ok := themes["new-theme"]; !ok {
		t.Error("new-theme not found after adding")
	}
}

func TestThemePluginAddFont(t *testing.T) {
	plugin := NewThemePluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		nil,
		"",
	)

	font := theme.Font{Family: "Roboto", Weights: []int{400, 700}}
	plugin.AddFont(font)

	fonts := plugin.Fonts()
	if len(fonts) != 1 {
		t.Errorf("expected 1 font after adding, got %d", len(fonts))
	}

	if fonts[0].Family != "Roboto" {
		t.Errorf("expected font 'Roboto', got %s", fonts[0].Family)
	}
}

func TestThemePluginCSS(t *testing.T) {
	plugin := NewThemePluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		nil,
		"",
	)

	// Default implementation returns empty string
	css := plugin.CSS()
	if css != "" {
		t.Errorf("expected empty CSS by default, got %s", css)
	}
}

type customThemePlugin struct {
	*ThemePluginBase
}

func (c *customThemePlugin) CSS() string {
	return ":root { --custom-color: red; }"
}

func TestThemePluginCustomCSS(t *testing.T) {
	plugin := &customThemePlugin{
		ThemePluginBase: NewThemePluginBase(
			PluginInfo{Name: "custom", Version: "1.0.0"},
			nil,
			"",
		),
	}

	css := plugin.CSS()
	if css == "" {
		t.Error("expected custom CSS, got empty string")
	}

	if css != ":root { --custom-color: red; }" {
		t.Errorf("unexpected CSS: %s", css)
	}
}

func TestThemePluginRegistration(t *testing.T) {
	registry := NewRegistry()

	plugin := NewThemePluginBase(
		PluginInfo{Name: "themes", Version: "1.0.0"},
		map[string]theme.Theme{
			"light": theme.DefaultLight(),
		},
		"light",
	)

	err := registry.Register(plugin)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Should be retrievable as ThemePlugin
	tp, ok := registry.GetThemePlugin("themes")
	if !ok {
		t.Error("theme plugin not found in registry")
	}

	if tp.Name() != "themes" {
		t.Errorf("expected name 'themes', got %s", tp.Name())
	}
}

func TestThemePluginLifecycle(t *testing.T) {
	registry := NewRegistry()
	ctx := context.Background()

	plugin := NewThemePluginBase(
		PluginInfo{Name: "test", Version: "1.0.0"},
		map[string]theme.Theme{
			"test": theme.DefaultLight(),
		},
		"test",
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

