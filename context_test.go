package forgeui

import (
	"context"
	"testing"
)

func TestWithTheme(t *testing.T) {
	ctx := context.Background()
	theme := "dark"

	ctx = WithTheme(ctx, theme)

	retrieved, ok := ThemeFromContext(ctx)
	if !ok {
		t.Error("ThemeFromContext() returned false, expected true")
	}

	if retrieved != theme {
		t.Errorf("ThemeFromContext() = %v, want %v", retrieved, theme)
	}
}

func TestThemeFromContext_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := ThemeFromContext(ctx)
	if ok {
		t.Error("ThemeFromContext() returned true for empty context, expected false")
	}
}

func TestWithConfig(t *testing.T) {
	ctx := context.Background()
	cfg := &Config{
		Debug: true,
		Theme: "custom",
	}

	ctx = WithConfig(ctx, cfg)

	retrieved, ok := ConfigFromContext(ctx)
	if !ok {
		t.Error("ConfigFromContext() returned false, expected true")
	}

	if retrieved != cfg {
		t.Error("ConfigFromContext() returned different config")
	}

	if !retrieved.Debug {
		t.Error("Config Debug should be true")
	}
}

func TestConfigFromContext_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := ConfigFromContext(ctx)
	if ok {
		t.Error("ConfigFromContext() returned true for empty context, expected false")
	}
}
