package bridge

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func renderComponent(t *testing.T, c interface {
	Render(context.Context, *bytes.Buffer) error
}) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	return buf.String()
}

func TestBridgeScripts_BasicInline(t *testing.T) {
	comp := BridgeScripts(ScriptConfig{
		Endpoint: "/api/bridge",
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	// Should include bridge JS
	if !strings.Contains(output, "<script") {
		t.Error("output should contain <script> tags")
	}

	// Should include endpoint config
	if !strings.Contains(output, "/api/bridge") {
		t.Error("output should contain endpoint")
	}
}

func TestBridgeScripts_WithAlpine(t *testing.T) {
	comp := BridgeScripts(ScriptConfig{
		Endpoint:      "/api/bridge",
		IncludeAlpine: true,
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	// Should have at least 2 script tags (bridge + alpine)
	count := strings.Count(output, "<script")
	if count < 2 {
		t.Errorf("expected at least 2 script tags, got %d", count)
	}
}

func TestBridgeScripts_WithHTMX(t *testing.T) {
	comp := BridgeScripts(ScriptConfig{
		Endpoint:    "/api/bridge",
		IncludeHTMX: true,
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	// Should include HTMX bridge JS
	if !strings.Contains(output, "bridgeFnURL") || !strings.Contains(output, "BRIDGE_FN_ENDPOINT") {
		t.Error("output should contain HTMX bridge code")
	}
}

func TestBridgeScripts_WithCSRFToken(t *testing.T) {
	comp := BridgeScripts(ScriptConfig{
		Endpoint:  "/api/bridge",
		CSRFToken: "test-csrf-token-123",
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "test-csrf-token-123") {
		t.Error("output should contain CSRF token")
	}
	if !strings.Contains(output, "setCSRF") {
		t.Error("output should contain setCSRF call")
	}
}

func TestBridgeScripts_NoCSRF(t *testing.T) {
	comp := BridgeScripts(ScriptConfig{
		Endpoint: "/api/bridge",
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	// The bridge JS itself defines setCSRF, but the config section should not call it
	if strings.Contains(output, "window.bridge.setCSRF") {
		t.Error("output should not call window.bridge.setCSRF when no CSRF token")
	}
}

func TestBridgeScripts_HTMXEndpointConfig(t *testing.T) {
	comp := BridgeScripts(ScriptConfig{
		Endpoint:    "/custom/path",
		IncludeHTMX: true,
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "/custom/path/fn/") {
		t.Error("output should contain HTMX fn endpoint based on bridge endpoint")
	}
}

// --- BridgeScriptsExternal ---

func TestBridgeScriptsExternal_DefaultStaticPath(t *testing.T) {
	comp := BridgeScriptsExternal(ScriptConfig{
		Endpoint: "/api/bridge",
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "/static/js/forge-bridge.js") {
		t.Error("should use default /static path")
	}
}

func TestBridgeScriptsExternal_CustomStaticPath(t *testing.T) {
	comp := BridgeScriptsExternal(ScriptConfig{
		Endpoint:   "/api/bridge",
		StaticPath: "/custom/assets",
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "/custom/assets/js/forge-bridge.js") {
		t.Error("should use custom static path")
	}
}

func TestBridgeScriptsExternal_WithAlpineAndHTMX(t *testing.T) {
	comp := BridgeScriptsExternal(ScriptConfig{
		Endpoint:      "/api/bridge",
		IncludeAlpine: true,
		IncludeHTMX:   true,
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "alpine-bridge.js") {
		t.Error("should include Alpine bridge script")
	}
	if !strings.Contains(output, "htmx-bridge.js") {
		t.Error("should include HTMX bridge script")
	}
}

func TestBridgeScriptsExternal_SrcAttribute(t *testing.T) {
	comp := BridgeScriptsExternal(ScriptConfig{
		Endpoint: "/api/bridge",
	})

	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	output := buf.String()

	// External scripts should use src attribute
	if !strings.Contains(output, "src=") {
		t.Error("external scripts should use src attribute")
	}
}

// --- GetJS functions ---

func TestGetBridgeJS(t *testing.T) {
	js := GetBridgeJS()
	if js == "" {
		t.Error("GetBridgeJS() returned empty string")
	}
	if !strings.Contains(js, "ForgeBridge") {
		t.Error("bridge JS should contain ForgeBridge")
	}
}

func TestGetAlpineJS(t *testing.T) {
	js := GetAlpineJS()
	if js == "" {
		t.Error("GetAlpineJS() returned empty string")
	}
}

func TestGetHTMXBridgeJS(t *testing.T) {
	js := GetHTMXBridgeJS()
	if js == "" {
		t.Error("GetHTMXBridgeJS() returned empty string")
	}
	if !strings.Contains(js, "bridgeFnURL") {
		t.Error("HTMX bridge JS should contain bridgeFnURL")
	}
}

// --- ScriptTemplate ---

func TestScriptTemplate_ValidTemplate(t *testing.T) {
	comp, err := ScriptTemplate("var config = {{.Endpoint}};", struct {
		Endpoint string
	}{Endpoint: `"/api/bridge"`})

	if err != nil {
		t.Fatalf("ScriptTemplate() error = %v", err)
	}
	if comp == nil {
		t.Fatal("ScriptTemplate() returned nil component")
	}

	var buf bytes.Buffer
	if renderErr := comp.Render(context.Background(), &buf); renderErr != nil {
		t.Fatalf("Render error: %v", renderErr)
	}

	output := buf.String()
	if !strings.Contains(output, "/api/bridge") {
		t.Errorf("output = %q, want to contain endpoint", output)
	}
	if !strings.Contains(output, "<script") {
		t.Error("output should be wrapped in script tags")
	}
}

func TestScriptTemplate_InvalidTemplate(t *testing.T) {
	_, err := ScriptTemplate("{{.Invalid", nil)
	if err == nil {
		t.Error("expected error for invalid template syntax")
	}
}

func TestScriptTemplate_ExecutionError(t *testing.T) {
	_, err := ScriptTemplate("{{.Missing}}", struct{}{})
	// template.Execute with strict missing key may or may not error,
	// but at minimum it should not panic
	_ = err
}
