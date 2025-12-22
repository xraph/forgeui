package assets

import "testing"

func TestStyleOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []StyleOption
		test func(*testing.T, *styleConfig)
	}{
		{
			name: "WithMedia sets media",
			opts: []StyleOption{WithMedia("print")},
			test: func(t *testing.T, cfg *styleConfig) {
				if cfg.media != "print" {
					t.Errorf("Expected media 'print', got '%s'", cfg.media)
				}
			},
		},
		{
			name: "WithPreload sets preload",
			opts: []StyleOption{WithPreload()},
			test: func(t *testing.T, cfg *styleConfig) {
				if !cfg.preload {
					t.Error("Expected preload to be true")
				}
			},
		},
		{
			name: "WithIntegrity sets integrity",
			opts: []StyleOption{WithIntegrity("sha256-abc123")},
			test: func(t *testing.T, cfg *styleConfig) {
				if cfg.integrity != "sha256-abc123" {
					t.Errorf("Expected integrity 'sha256-abc123', got '%s'", cfg.integrity)
				}
			},
		},
		{
			name: "WithCrossOrigin sets crossOrigin",
			opts: []StyleOption{WithCrossOrigin("anonymous")},
			test: func(t *testing.T, cfg *styleConfig) {
				if cfg.crossOrigin != "anonymous" {
					t.Errorf("Expected crossOrigin 'anonymous', got '%s'", cfg.crossOrigin)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &styleConfig{}
			for _, opt := range tt.opts {
				opt(cfg)
			}
			tt.test(t, cfg)
		})
	}
}

func TestScriptOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []ScriptOption
		test func(*testing.T, *scriptConfig)
	}{
		{
			name: "WithDefer sets defer",
			opts: []ScriptOption{WithDefer()},
			test: func(t *testing.T, cfg *scriptConfig) {
				if !cfg.defer_ {
					t.Error("Expected defer to be true")
				}
			},
		},
		{
			name: "WithAsync sets async",
			opts: []ScriptOption{WithAsync()},
			test: func(t *testing.T, cfg *scriptConfig) {
				if !cfg.async {
					t.Error("Expected async to be true")
				}
			},
		},
		{
			name: "WithModule sets module",
			opts: []ScriptOption{WithModule()},
			test: func(t *testing.T, cfg *scriptConfig) {
				if !cfg.module {
					t.Error("Expected module to be true")
				}
			},
		},
		{
			name: "WithScriptIntegrity sets integrity",
			opts: []ScriptOption{WithScriptIntegrity("sha256-def456")},
			test: func(t *testing.T, cfg *scriptConfig) {
				if cfg.integrity != "sha256-def456" {
					t.Errorf("Expected integrity 'sha256-def456', got '%s'", cfg.integrity)
				}
			},
		},
		{
			name: "WithScriptCrossOrigin sets crossOrigin",
			opts: []ScriptOption{WithScriptCrossOrigin("use-credentials")},
			test: func(t *testing.T, cfg *scriptConfig) {
				if cfg.crossOrigin != "use-credentials" {
					t.Errorf("Expected crossOrigin 'use-credentials', got '%s'", cfg.crossOrigin)
				}
			},
		},
		{
			name: "WithNoModule sets noModule",
			opts: []ScriptOption{WithNoModule()},
			test: func(t *testing.T, cfg *scriptConfig) {
				if !cfg.noModule {
					t.Error("Expected noModule to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &scriptConfig{}
			for _, opt := range tt.opts {
				opt(cfg)
			}
			tt.test(t, cfg)
		})
	}
}

