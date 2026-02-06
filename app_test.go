package forgeui

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestApp_BridgeScriptHandler(t *testing.T) {
	tests := []struct {
		name           string
		basePath       string
		requestPath    string
		expectStatus   int
		expectContent  string
		expectMIME     string
	}{
		{
			name:          "forge-bridge.js without base path",
			basePath:      "",
			requestPath:   "/static/js/forge-bridge.js",
			expectStatus:  http.StatusOK,
			expectContent: "ForgeBridge", // Part of the embedded JS
			expectMIME:    "text/javascript; charset=utf-8",
		},
		{
			name:          "alpine-bridge.js without base path",
			basePath:      "",
			requestPath:   "/static/js/alpine-bridge.js",
			expectStatus:  http.StatusOK,
			expectContent: "Alpine", // Part of the embedded JS
			expectMIME:    "text/javascript; charset=utf-8",
		},
		{
			name:          "forge-bridge.js with base path",
			basePath:      "/api/identity/ui",
			requestPath:   "/api/identity/ui/static/js/forge-bridge.js",
			expectStatus:  http.StatusOK,
			expectContent: "ForgeBridge",
			expectMIME:    "text/javascript; charset=utf-8",
		},
		{
			name:          "alpine-bridge.js with base path",
			basePath:      "/api/identity/ui",
			requestPath:   "/api/identity/ui/static/js/alpine-bridge.js",
			expectStatus:  http.StatusOK,
			expectContent: "Alpine",
			expectMIME:    "text/javascript; charset=utf-8",
		},
		{
			name:         "non-existent script",
			basePath:     "",
			requestPath:  "/static/js/unknown.js",
			expectStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []AppOption{WithBridge()}
			if tt.basePath != "" {
				opts = append(opts, WithBasePath(tt.basePath))
			}

			app := New(opts...)
			handler := app.Handler()

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}

			if tt.expectStatus == http.StatusOK {
				body := w.Body.String()
				if !strings.Contains(body, tt.expectContent) {
					t.Errorf("Expected body to contain '%s', got: %s", tt.expectContent, body[:100])
				}

				contentType := w.Header().Get("Content-Type")
				if contentType != tt.expectMIME {
					t.Errorf("Expected Content-Type '%s', got '%s'", tt.expectMIME, contentType)
				}

				cacheControl := w.Header().Get("Cache-Control")
				if !strings.Contains(cacheControl, "public") {
					t.Errorf("Expected Cache-Control to contain 'public', got: %s", cacheControl)
				}
			}
		})
	}
}

func TestApp_BridgeScriptHandler_NoBridge(t *testing.T) {
	// When bridge is not enabled, scripts should 404
	app := New() // No WithBridge()
	handler := app.Handler()

	req := httptest.NewRequest(http.MethodGet, "/static/js/forge-bridge.js", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 when bridge disabled, got %d", w.Code)
	}
}

func TestApp_BasePath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		wantCall string
		wantJS   string
	}{
		{
			name:     "no base path",
			basePath: "",
			wantCall: "/api/bridge/call",
			wantJS:   "/static/js/forge-bridge.js",
		},
		{
			name:     "with base path",
			basePath: "/api/identity/ui",
			wantCall: "/api/identity/ui/bridge/call",
			wantJS:   "/api/identity/ui/static/js/forge-bridge.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []AppOption{WithBridge()}
			if tt.basePath != "" {
				opts = append(opts, WithBasePath(tt.basePath))
			}

			app := New(opts...)

			if got := app.BridgeCallPath(); got != tt.wantCall {
				t.Errorf("BridgeCallPath() = %v, want %v", got, tt.wantCall)
			}

			if got := app.BridgeScriptPath(); got != tt.wantJS {
				t.Errorf("BridgeScriptPath() = %v, want %v", got, tt.wantJS)
			}
		})
	}
}
