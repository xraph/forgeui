package forgeui

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Debug {
		t.Error("expected Debug to be false by default")
	}

	if cfg.Theme != "default" {
		t.Errorf("expected Theme to be 'default', got %v", cfg.Theme)
	}

	if cfg.StaticPath != "/static" {
		t.Errorf("expected StaticPath to be '/static', got %v", cfg.StaticPath)
	}

	if cfg.AssetsPath != "public" {
		t.Errorf("expected AssetsPath to be 'public', got %v", cfg.AssetsPath)
	}

	if cfg.DefaultSize != SizeMD {
		t.Errorf("expected DefaultSize to be SizeMD, got %v", cfg.DefaultSize)
	}

	if cfg.DefaultVariant != VariantDefault {
		t.Errorf("expected DefaultVariant to be VariantDefault, got %v", cfg.DefaultVariant)
	}

	if cfg.DefaultRadius != RadiusMD {
		t.Errorf("expected DefaultRadius to be RadiusMD, got %v", cfg.DefaultRadius)
	}
}

func TestConfigOptions(t *testing.T) {
	cfg := DefaultConfig()

	WithDebug(true)(cfg)

	if !cfg.Debug {
		t.Error("WithDebug(true) did not set Debug to true")
	}

	WithThemeName("dark")(cfg)

	if cfg.Theme != "dark" {
		t.Errorf("WithThemeName('dark') did not set Theme, got %v", cfg.Theme)
	}

	WithStaticPath("/assets")(cfg)

	if cfg.StaticPath != "/assets" {
		t.Errorf("WithStaticPath('/assets') did not set StaticPath, got %v", cfg.StaticPath)
	}

	WithDefaultSize(SizeLG)(cfg)

	if cfg.DefaultSize != SizeLG {
		t.Errorf("WithDefaultSize(SizeLG) did not set DefaultSize, got %v", cfg.DefaultSize)
	}
}

func TestNew(t *testing.T) {
	app := New(
		WithDebug(true),
		WithThemeName("custom"),
	)

	cfg := app.Config()
	if !cfg.Debug {
		t.Error("expected Debug to be true")
	}

	if cfg.Theme != "custom" {
		t.Errorf("expected Theme to be 'custom', got %v", cfg.Theme)
	}

	if !app.IsDev() {
		t.Error("IsDev() should return true when Debug is true")
	}
}
