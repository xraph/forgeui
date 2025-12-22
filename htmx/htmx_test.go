package htmx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHxGet(t *testing.T) {
	node := HxGet("/api/users")
	html := renderNode(node)
	if html != `hx-get="/api/users"` {
		t.Errorf("Expected hx-get attribute, got: %s", html)
	}
}

func TestHxPost(t *testing.T) {
	node := HxPost("/api/create")
	html := renderNode(node)
	if html != `hx-post="/api/create"` {
		t.Errorf("Expected hx-post attribute, got: %s", html)
	}
}

func TestHxPut(t *testing.T) {
	node := HxPut("/api/update")
	html := renderNode(node)
	if html != `hx-put="/api/update"` {
		t.Errorf("Expected hx-put attribute, got: %s", html)
	}
}

func TestHxPatch(t *testing.T) {
	node := HxPatch("/api/patch")
	html := renderNode(node)
	if html != `hx-patch="/api/patch"` {
		t.Errorf("Expected hx-patch attribute, got: %s", html)
	}
}

func TestHxDelete(t *testing.T) {
	node := HxDelete("/api/delete")
	html := renderNode(node)
	if html != `hx-delete="/api/delete"` {
		t.Errorf("Expected hx-delete attribute, got: %s", html)
	}
}

func TestHxTarget(t *testing.T) {
	node := HxTarget("#results")
	html := renderNode(node)
	if html != `hx-target="#results"` {
		t.Errorf("Expected hx-target attribute, got: %s", html)
	}
}

func TestHxSwap(t *testing.T) {
	tests := []struct {
		name     string
		strategy string
		expected string
	}{
		{"innerHTML", "innerHTML", `hx-swap="innerHTML"`},
		{"outerHTML", "outerHTML", `hx-swap="outerHTML"`},
		{"beforebegin", "beforebegin", `hx-swap="beforebegin"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := HxSwap(tt.strategy)
			html := renderNode(node)
			if html != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, html)
			}
		})
	}
}

func TestHxSwapConvenience(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() interface{ Render(io.Writer) error }
		expected string
	}{
		{"innerHTML", func() interface{ Render(io.Writer) error } { return HxSwapInnerHTML() }, `hx-swap="innerHTML"`},
		{"outerHTML", func() interface{ Render(io.Writer) error } { return HxSwapOuterHTML() }, `hx-swap="outerHTML"`},
		{"beforebegin", func() interface{ Render(io.Writer) error } { return HxSwapBeforeBegin() }, `hx-swap="beforebegin"`},
		{"afterbegin", func() interface{ Render(io.Writer) error } { return HxSwapAfterBegin() }, `hx-swap="afterbegin"`},
		{"beforeend", func() interface{ Render(io.Writer) error } { return HxSwapBeforeEnd() }, `hx-swap="beforeend"`},
		{"afterend", func() interface{ Render(io.Writer) error } { return HxSwapAfterEnd() }, `hx-swap="afterend"`},
		{"delete", func() interface{ Render(io.Writer) error } { return HxSwapDelete() }, `hx-swap="delete"`},
		{"none", func() interface{ Render(io.Writer) error } { return HxSwapNone() }, `hx-swap="none"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.fn()
			html := renderNode(node)
			if html != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, html)
			}
		})
	}
}

func TestHxBoost(t *testing.T) {
	node := HxBoost(true)
	html := renderNode(node)
	if html != `hx-boost="true"` {
		t.Errorf("Expected hx-boost=true, got: %s", html)
	}

	node = HxBoost(false)
	html = renderNode(node)
	if html != `hx-boost="false"` {
		t.Errorf("Expected hx-boost=false, got: %s", html)
	}
}

func TestHxTrigger(t *testing.T) {
	node := HxTrigger("click")
	html := renderNode(node)
	if html != `hx-trigger="click"` {
		t.Errorf("Expected hx-trigger, got: %s", html)
	}
}

func TestHxTriggerDebounce(t *testing.T) {
	node := HxTriggerDebounce("keyup", "500ms")
	html := renderNode(node)
	if !strings.Contains(html, "keyup changed delay:500ms") {
		t.Errorf("Expected debounced trigger, got: %s", html)
	}
}

func TestHxTriggerThrottle(t *testing.T) {
	node := HxTriggerThrottle("scroll", "1s")
	html := renderNode(node)
	if !strings.Contains(html, "scroll throttle:1s") {
		t.Errorf("Expected throttled trigger, got: %s", html)
	}
}

func TestHxIndicator(t *testing.T) {
	node := HxIndicator("#spinner")
	html := renderNode(node)
	if html != `hx-indicator="#spinner"` {
		t.Errorf("Expected hx-indicator, got: %s", html)
	}
}

func TestHxConfirm(t *testing.T) {
	node := HxConfirm("Are you sure?")
	html := renderNode(node)
	if html != `hx-confirm="Are you sure?"` {
		t.Errorf("Expected hx-confirm, got: %s", html)
	}
}

func TestIsHTMX(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	
	if IsHTMX(req) {
		t.Error("Expected non-HTMX request")
	}

	req.Header.Set("HX-Request", "true")
	if !IsHTMX(req) {
		t.Error("Expected HTMX request")
	}
}

func TestIsHTMXBoosted(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	
	if IsHTMXBoosted(req) {
		t.Error("Expected non-boosted request")
	}

	req.Header.Set("HX-Boosted", "true")
	if !IsHTMXBoosted(req) {
		t.Error("Expected boosted request")
	}
}

func TestHTMXTarget(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("HX-Target", "main")

	target := HTMXTarget(req)
	if target != "main" {
		t.Errorf("Expected target 'main', got: %s", target)
	}
}

func TestHTMXTrigger(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("HX-Trigger", "button-1")

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

	header := w.Header().Get("HX-Trigger")
	if !strings.Contains(header, "event1") {
		t.Errorf("Expected HX-Trigger header with event1, got: %s", header)
	}
}

func TestSetHTMXRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	SetHTMXRedirect(w, "/login")

	header := w.Header().Get("HX-Redirect")
	if header != "/login" {
		t.Errorf("Expected HX-Redirect=/login, got: %s", header)
	}
}

func TestSetHTMXRefresh(t *testing.T) {
	w := httptest.NewRecorder()
	SetHTMXRefresh(w)

	header := w.Header().Get("HX-Refresh")
	if header != "true" {
		t.Errorf("Expected HX-Refresh=true, got: %s", header)
	}
}

func TestTriggerEvent(t *testing.T) {
	w := httptest.NewRecorder()
	TriggerEvent(w, "myEvent")

	header := w.Header().Get("HX-Trigger")
	if header != "myEvent" {
		t.Errorf("Expected HX-Trigger=myEvent, got: %s", header)
	}
}

func TestTriggerEventWithDetail(t *testing.T) {
	w := httptest.NewRecorder()
	TriggerEventWithDetail(w, "showMessage", map[string]any{
		"text": "Hello",
	})

	header := w.Header().Get("HX-Trigger")
	if !strings.Contains(header, "showMessage") || !strings.Contains(header, "Hello") {
		t.Errorf("Expected HX-Trigger with showMessage and Hello, got: %s", header)
	}
}

func TestMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Middleware(handler)
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("HX-Request", "true")
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
	node := Scripts()
	html := renderNode(node)
	
	if !strings.Contains(html, "htmx.org") {
		t.Errorf("Expected htmx.org CDN URL, got: %s", html)
	}
}

func TestScriptsWithVersion(t *testing.T) {
	node := Scripts("1.9.0")
	html := renderNode(node)
	
	if !strings.Contains(html, "1.9.0") {
		t.Errorf("Expected version 1.9.0, got: %s", html)
	}
}

// Helper function to render a node to string
func renderNode(node interface{ Render(io.Writer) error }) string {
	var sb strings.Builder
	if err := node.Render(&sb); err != nil {
		return ""
	}
	// Trim leading/trailing whitespace - gomponents attributes render with leading space
	return strings.TrimSpace(sb.String())
}

