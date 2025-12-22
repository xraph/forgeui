package forgeui

import (
	"context"
	"testing"
)

func TestApp_New(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("New() returned nil")
	}

	cfg := app.Config()
	if cfg == nil {
		t.Fatal("Config() returned nil")
	}
}

func TestApp_NewWithOptions(t *testing.T) {
	app := New(
		WithDebug(true),
		WithThemeName("dark"),
		WithStaticPath("/assets"),
		WithDefaultSize(SizeLG),
	)

	cfg := app.Config()
	if !cfg.Debug {
		t.Error("expected Debug to be true")
	}

	if cfg.Theme != "dark" {
		t.Errorf("expected Theme to be 'dark', got %v", cfg.Theme)
	}

	if cfg.StaticPath != "/assets" {
		t.Errorf("expected StaticPath to be '/assets', got %v", cfg.StaticPath)
	}

	if cfg.DefaultSize != SizeLG {
		t.Errorf("expected DefaultSize to be SizeLG, got %v", cfg.DefaultSize)
	}
}

func TestApp_IsDev(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
		want  bool
	}{
		{
			name:  "debug mode",
			debug: true,
			want:  true,
		},
		{
			name:  "production mode",
			debug: false,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(WithDebug(tt.debug))
			if got := app.IsDev(); got != tt.want {
				t.Errorf("IsDev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_Initialize(t *testing.T) {
	app := New()

	err := app.Initialize(context.Background())
	if err != nil {
		t.Errorf("Initialize() error = %v, want nil", err)
	}
}

func TestApp_Config(t *testing.T) {
	app := New(WithThemeName("custom"))

	cfg1 := app.Config()
	cfg2 := app.Config()

	// Should return the same config instance
	if cfg1 != cfg2 {
		t.Error("Config() should return the same instance")
	}

	if cfg1.Theme != "custom" {
		t.Errorf("Config.Theme = %v, want 'custom'", cfg1.Theme)
	}
}
