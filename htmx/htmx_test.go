package htmx

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestHxGet(t *testing.T) {
	attrs := HxGet("/api/users")
	assertAttr(t, attrs, "hx-get", "/api/users")
}

func TestHxPost(t *testing.T) {
	attrs := HxPost("/api/create")
	assertAttr(t, attrs, "hx-post", "/api/create")
}

func TestHxPut(t *testing.T) {
	attrs := HxPut("/api/update")
	assertAttr(t, attrs, "hx-put", "/api/update")
}

func TestHxPatch(t *testing.T) {
	attrs := HxPatch("/api/patch")
	assertAttr(t, attrs, "hx-patch", "/api/patch")
}

func TestHxDelete(t *testing.T) {
	attrs := HxDelete("/api/delete")
	assertAttr(t, attrs, "hx-delete", "/api/delete")
}

func TestHxTarget(t *testing.T) {
	attrs := HxTarget("#results")
	assertAttr(t, attrs, "hx-target", "#results")
}

func TestHxSwap(t *testing.T) {
	tests := []struct {
		name     string
		strategy string
	}{
		{"innerHTML", "innerHTML"},
		{"outerHTML", "outerHTML"},
		{"beforebegin", "beforebegin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := HxSwap(tt.strategy)
			assertAttr(t, attrs, "hx-swap", tt.strategy)
		})
	}
}

func TestHxSwapConvenience(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() templ.Attributes
		expected string
	}{
		{"innerHTML", HxSwapInnerHTML, "innerHTML"},
		{"outerHTML", HxSwapOuterHTML, "outerHTML"},
		{"beforebegin", HxSwapBeforeBegin, "beforebegin"},
		{"afterbegin", HxSwapAfterBegin, "afterbegin"},
		{"beforeend", HxSwapBeforeEnd, "beforeend"},
		{"afterend", HxSwapAfterEnd, "afterend"},
		{"delete", HxSwapDelete, "delete"},
		{"none", HxSwapNone, "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := tt.fn()
			assertAttr(t, attrs, "hx-swap", tt.expected)
		})
	}
}

func TestHxBoost(t *testing.T) {
	attrs := HxBoost(true)
	assertAttr(t, attrs, "hx-boost", "true")

	attrs = HxBoost(false)
	assertAttr(t, attrs, "hx-boost", "false")
}

func TestHxTrigger(t *testing.T) {
	attrs := HxTrigger("click")
	assertAttr(t, attrs, "hx-trigger", "click")
}

func TestHxTriggerDebounce(t *testing.T) {
	attrs := HxTriggerDebounce("keyup", "500ms")
	v := attrs["hx-trigger"]
	if !strings.Contains(v.(string), "keyup changed delay:500ms") {
		t.Errorf("Expected debounced trigger, got: %s", v)
	}
}

func TestHxTriggerThrottle(t *testing.T) {
	attrs := HxTriggerThrottle("scroll", "1s")
	v := attrs["hx-trigger"]
	if !strings.Contains(v.(string), "scroll throttle:1s") {
		t.Errorf("Expected throttled trigger, got: %s", v)
	}
}

func TestHxIndicator(t *testing.T) {
	attrs := HxIndicator("#spinner")
	assertAttr(t, attrs, "hx-indicator", "#spinner")
}

func TestHxConfirm(t *testing.T) {
	attrs := HxConfirm("Are you sure?")
	assertAttr(t, attrs, "hx-confirm", "Are you sure?")
}

func TestIsHTMX(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	if IsHTMX(req) {
		t.Error("Expected non-HTMX request")
	}

	req.Header.Set("Hx-Request", "true")

	if !IsHTMX(req) {
		t.Error("Expected HTMX request")
	}
}

func TestIsHTMXBoosted(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	if IsHTMXBoosted(req) {
		t.Error("Expected non-boosted request")
	}

	req.Header.Set("Hx-Boosted", "true")

	if !IsHTMXBoosted(req) {
		t.Error("Expected boosted request")
	}
}

func TestHTMXTarget(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Hx-Target", "main")

	target := HTMXTarget(req)
	if target != "main" {
		t.Errorf("Expected target 'main', got: %s", target)
	}
}

func TestHTMXTrigger(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Hx-Trigger", "button-1")

	trigger := HTMXTrigger(req)
	if trigger != "button-1" {
		t.Errorf("Expected trigger 'button-1', got: %s", trigger)
	}
}

func TestSetHTMXTrigger(t *testing.T) {
	w := httptest.NewRecorder()
	SetHTMXTrigger(w, map[string]any{
		"event1": "data",
	})

	header := w.Header().Get("Hx-Trigger")
	if !strings.Contains(header, "event1") {
		t.Errorf("Expected HX-Trigger header with event1, got: %s", header)
	}
}

func TestSetHTMXRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	SetHTMXRedirect(w, "/login")

	header := w.Header().Get("Hx-Redirect")
	if header != "/login" {
		t.Errorf("Expected HX-Redirect=/login, got: %s", header)
	}
}

func TestSetHTMXRefresh(t *testing.T) {
	w := httptest.NewRecorder()
	SetHTMXRefresh(w)

	header := w.Header().Get("Hx-Refresh")
	if header != "true" {
		t.Errorf("Expected HX-Refresh=true, got: %s", header)
	}
}

func TestTriggerEvent(t *testing.T) {
	w := httptest.NewRecorder()
	TriggerEvent(w, "myEvent")

	header := w.Header().Get("Hx-Trigger")
	if header != "myEvent" {
		t.Errorf("Expected HX-Trigger=myEvent, got: %s", header)
	}
}

func TestTriggerEventWithDetail(t *testing.T) {
	w := httptest.NewRecorder()
	TriggerEventWithDetail(w, "showMessage", map[string]any{
		"text": "Hello",
	})

	header := w.Header().Get("Hx-Trigger")
	if !strings.Contains(header, "showMessage") || !strings.Contains(header, "Hello") {
		t.Errorf("Expected HX-Trigger with showMessage and Hello, got: %s", header)
	}
}

func TestMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Hx-Request", "true")

	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", w.Code)
	}
}

func TestStopPolling(t *testing.T) {
	w := httptest.NewRecorder()
	StopPolling(w)

	if w.Code != 286 {
		t.Errorf("Expected status 286, got: %d", w.Code)
	}
}

func TestScripts(t *testing.T) {
	comp := Scripts()
	html := renderComponent(comp)

	if !strings.Contains(html, "htmx.org") {
		t.Errorf("Expected htmx.org CDN URL, got: %s", html)
	}
}

func TestScriptsWithVersion(t *testing.T) {
	comp := Scripts("1.9.0")
	html := renderComponent(comp)

	if !strings.Contains(html, "1.9.0") {
		t.Errorf("Expected version 1.9.0, got: %s", html)
	}
}

// Helper function to render a templ.Component to string
func renderComponent(comp templ.Component) string {
	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}

// Helper function to assert a templ.Attributes key-value pair
func assertAttr(t *testing.T, attrs templ.Attributes, key, expected string) {
	t.Helper()
	v, ok := attrs[key]
	if !ok {
		t.Errorf("Expected attribute %q to exist", key)
		return
	}
	if v != expected {
		t.Errorf("Expected %s=%q, got %q", key, expected, v)
	}
}
