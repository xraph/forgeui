package bridge

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// testTemplComponent creates a simple templ.Component for testing
func testTemplComponent(text string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, "<div>%s</div>", text)
		return err
	})
}

func setupHTMXBridge(t *testing.T) *Bridge {
	t.Helper()
	b := New(WithCSRF(false))

	// SigInputOutput: func(Context, Input) (Output, error)
	err := b.Register("test.getData", func(ctx Context, params testInput) (testOutput, error) {
		return testOutput{Result: "got " + params.Name}, nil
	})
	if err != nil {
		t.Fatalf("Register getData: %v", err)
	}

	// SigOutput: func(Context) (Output, error) - no input
	err = b.Register("test.getAll", func(ctx Context) (testOutput, error) {
		return testOutput{Result: "all data"}, nil
	})
	if err != nil {
		t.Fatalf("Register getAll: %v", err)
	}

	// SigVoid: func(Context) error
	err = b.Register("test.ping", func(ctx Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Register ping: %v", err)
	}

	// SigInputOnly: func(Context, Input) error
	err = b.Register("test.action", func(ctx Context, params testInput) error {
		if params.Name == "" {
			return fmt.Errorf("name required")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Register action: %v", err)
	}

	// Returns templ.Component directly
	err = b.Register("test.htmlDirect", func(ctx Context) (templ.Component, error) {
		return testTemplComponent("hello from bridge"), nil
	})
	if err != nil {
		t.Fatalf("Register htmlDirect: %v", err)
	}

	// Returns data with a renderer
	err = b.Register("test.htmlRendered", func(ctx Context) (testOutput, error) {
		return testOutput{Result: "rendered data"}, nil
	}, WithRenderer(func(data testOutput) templ.Component {
		return testTemplComponent("rendered: " + data.Result)
	}))
	if err != nil {
		t.Fatalf("Register htmlRendered: %v", err)
	}

	// Function with HTMX headers
	err = b.Register("test.withHeaders", func(ctx Context) error {
		return nil
	}, WithHTMXTrigger("itemUpdated"), WithHTMXReswap("none"))
	if err != nil {
		t.Fatalf("Register withHeaders: %v", err)
	}

	// Function restricted to POST only
	err = b.Register("test.postOnly", func(ctx Context, params testInput) error {
		return nil
	}, WithHTTPMethod("POST"), WithLaxValidation())
	if err != nil {
		t.Fatalf("Register postOnly: %v", err)
	}

	return b
}

func TestHTMXHandler_GetNoInput(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.getAll", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	// Should return JSON for non-HTMX request
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Errorf("Content-Type = %s, want application/json", w.Header().Get("Content-Type"))
	}
}

func TestHTMXHandler_GetWithQueryParams(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	// test.getData has input, defaults to POST only
	// Override with explicit GET
	b.mu.Lock()
	fn := b.functions["test.getData"]
	fn.AllowedMethods = []string{"GET", "POST"}
	b.mu.Unlock()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.getData?name=world", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	body := w.Body.String()
	if !strings.Contains(body, "got world") {
		t.Errorf("body = %q, want to contain 'got world'", body)
	}
}

func TestHTMXHandler_PostFormData(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	body := strings.NewReader("name=formtest")
	r := httptest.NewRequest(http.MethodPost, "/api/bridge/fn/test.postOnly", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestHTMXHandler_PostJSON(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	body := strings.NewReader(`{"name":"jsontest"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/bridge/fn/test.postOnly", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestHTMXHandler_HTMLDirect(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.htmlDirect", nil)
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Content-Type = %s, want text/html", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()
	if !strings.Contains(body, "hello from bridge") {
		t.Errorf("body = %q, want to contain 'hello from bridge'", body)
	}
}

func TestHTMXHandler_HTMLRendered(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.htmlRendered", nil)
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "rendered: rendered data") {
		t.Errorf("body = %q, want to contain 'rendered: rendered data'", body)
	}
}

func TestHTMXHandler_VoidFunction(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.ping", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHTMXHandler_HTMXHeaders(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.withHeaders", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Header().Get("HX-Trigger") != "itemUpdated" {
		t.Errorf("HX-Trigger = %q, want 'itemUpdated'", w.Header().Get("HX-Trigger"))
	}

	if w.Header().Get("HX-Reswap") != "none" {
		t.Errorf("HX-Reswap = %q, want 'none'", w.Header().Get("HX-Reswap"))
	}
}

func TestHTMXHandler_NotFound(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/nonexistent", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHTMXHandler_MethodNotAllowed(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	// test.postOnly is POST only, try GET
	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.postOnly", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestHTMXHandler_EmptyFunctionName(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHTMXHandler_AutoDetectMethods(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	// test.getAll has no input → should accept GET
	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.getAll", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("GET getAll: status = %d, want %d", w.Code, http.StatusOK)
	}

	// test.action has input → should default to POST only, reject GET
	r = httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.action", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET action: status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestHTMXHandler_FallbackJSONForHTMX(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler()

	// test.getAll returns data, no renderer, HTMX request → <pre> JSON fallback
	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.getAll", nil)
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<pre>") {
		t.Errorf("body = %q, want to contain '<pre>'", body)
	}
	if !strings.Contains(body, "all data") {
		t.Errorf("body = %q, want to contain 'all data'", body)
	}
}

func TestHTMXHandler_CustomPrefix(t *testing.T) {
	b := setupHTMXBridge(t)
	handler := b.HTMXHandler("/custom/prefix/")

	r := httptest.NewRequest(http.MethodGet, "/custom/prefix/test.ping", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

// --- isMethodAllowed unit tests ---

func TestIsMethodAllowed(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	tests := []struct {
		name    string
		fn      *Function
		method  string
		allowed bool
	}{
		{
			name:    "explicit GET allowed",
			fn:      &Function{AllowedMethods: []string{"GET"}},
			method:  "GET",
			allowed: true,
		},
		{
			name:    "explicit POST not in list",
			fn:      &Function{AllowedMethods: []string{"GET"}},
			method:  "POST",
			allowed: false,
		},
		{
			name:    "case insensitive match",
			fn:      &Function{AllowedMethods: []string{"POST"}},
			method:  "post",
			allowed: true,
		},
		{
			name:    "auto-detect no input → GET allowed",
			fn:      &Function{HasInput: false},
			method:  "GET",
			allowed: true,
		},
		{
			name:    "auto-detect no input → POST allowed",
			fn:      &Function{HasInput: false},
			method:  "POST",
			allowed: true,
		},
		{
			name:    "auto-detect with input → GET rejected",
			fn:      &Function{HasInput: true},
			method:  "GET",
			allowed: false,
		},
		{
			name:    "auto-detect with input → POST allowed",
			fn:      &Function{HasInput: true},
			method:  "POST",
			allowed: true,
		},
		{
			name:    "DELETE not in any default",
			fn:      &Function{HasInput: false},
			method:  "DELETE",
			allowed: false,
		},
		{
			name:    "multiple explicit methods",
			fn:      &Function{AllowedMethods: []string{"GET", "POST", "PUT"}},
			method:  "PUT",
			allowed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := h.isMethodAllowed(tt.fn, tt.method)
			if got != tt.allowed {
				t.Errorf("isMethodAllowed() = %v, want %v", got, tt.allowed)
			}
		})
	}
}

// --- errorToStatus unit tests ---

func TestErrorToStatus(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	tests := []struct {
		name     string
		err      *Error
		wantCode int
	}{
		{"unauthorized", NewError(ErrCodeUnauthorized, ""), http.StatusUnauthorized},
		{"forbidden", NewError(ErrCodeForbidden, ""), http.StatusForbidden},
		{"not found", NewError(ErrCodeMethodNotFound, ""), http.StatusNotFound},
		{"invalid params", NewError(ErrCodeInvalidParams, ""), http.StatusBadRequest},
		{"bad request", NewError(ErrCodeBadRequest, ""), http.StatusBadRequest},
		{"rate limit", NewError(ErrCodeRateLimit, ""), http.StatusTooManyRequests},
		{"timeout", NewError(ErrCodeTimeout, ""), http.StatusGatewayTimeout},
		{"internal", NewError(ErrCodeInternal, ""), http.StatusInternalServerError},
		{"unknown code", NewError(99999, ""), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := h.errorToStatus(tt.err)
			if got != tt.wantCode {
				t.Errorf("errorToStatus() = %d, want %d", got, tt.wantCode)
			}
		})
	}
}

// --- setHTMXHeaders unit tests ---

func TestSetHTMXHeaders_All(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	fn := &Function{
		HTMXTriggers: []string{"created", "refreshList"},
		HTMXRedirect: "/dashboard",
		HTMXReswap:   "outerHTML",
		HTMXRetarget: "#main",
	}

	w := httptest.NewRecorder()
	h.setHTMXHeaders(w, fn)

	if got := w.Header().Get("HX-Trigger"); got != "created, refreshList" {
		t.Errorf("HX-Trigger = %q, want %q", got, "created, refreshList")
	}
	if got := w.Header().Get("HX-Redirect"); got != "/dashboard" {
		t.Errorf("HX-Redirect = %q, want %q", got, "/dashboard")
	}
	if got := w.Header().Get("HX-Reswap"); got != "outerHTML" {
		t.Errorf("HX-Reswap = %q, want %q", got, "outerHTML")
	}
	if got := w.Header().Get("HX-Retarget"); got != "#main" {
		t.Errorf("HX-Retarget = %q, want %q", got, "#main")
	}
}

func TestSetHTMXHeaders_Empty(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	fn := &Function{} // no HTMX config

	w := httptest.NewRecorder()
	h.setHTMXHeaders(w, fn)

	if got := w.Header().Get("HX-Trigger"); got != "" {
		t.Errorf("HX-Trigger = %q, want empty", got)
	}
	if got := w.Header().Get("HX-Redirect"); got != "" {
		t.Errorf("HX-Redirect = %q, want empty", got)
	}
	if got := w.Header().Get("HX-Reswap"); got != "" {
		t.Errorf("HX-Reswap = %q, want empty", got)
	}
	if got := w.Header().Get("HX-Retarget"); got != "" {
		t.Errorf("HX-Retarget = %q, want empty", got)
	}
}

// --- renderResponse unit tests ---

func TestRenderResponse_VoidResult(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	fn := &Function{}
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	h.renderResponse(w, r, fn, nil)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", w.Body.String())
	}
}

func TestRenderResponse_HTMLDirect(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	fn := &Function{ReturnsHTML: true}
	comp := testTemplComponent("test content")

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	h.renderResponse(w, r, fn, comp)

	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Content-Type = %s, want text/html", w.Header().Get("Content-Type"))
	}
	if !strings.Contains(w.Body.String(), "test content") {
		t.Errorf("body = %q, want to contain 'test content'", w.Body.String())
	}
}

func TestRenderResponse_WithRenderer(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	fn := &Function{
		Renderer: func(data any) templ.Component {
			d := data.(testOutput)
			return testTemplComponent("rendered:" + d.Result)
		},
	}

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	h.renderResponse(w, r, fn, testOutput{Result: "mydata"})

	if !strings.Contains(w.Body.String(), "rendered:mydata") {
		t.Errorf("body = %q, want to contain 'rendered:mydata'", w.Body.String())
	}
}

func TestRenderResponse_NonHTMX_JSON(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	fn := &Function{}

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No HX-Request header
	w := httptest.NewRecorder()

	h.renderResponse(w, r, fn, testOutput{Result: "json data"})

	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Errorf("Content-Type = %s, want application/json", w.Header().Get("Content-Type"))
	}
	if !strings.Contains(w.Body.String(), "json data") {
		t.Errorf("body = %q, want to contain 'json data'", w.Body.String())
	}
}

func TestRenderResponse_HTMX_NoRenderer_PreFallback(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/api/bridge/fn/")

	fn := &Function{}

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	h.renderResponse(w, r, fn, testOutput{Result: "fallback"})

	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Content-Type = %s, want text/html", w.Header().Get("Content-Type"))
	}
	body := w.Body.String()
	if !strings.Contains(body, "<pre>") {
		t.Errorf("body = %q, want to contain '<pre>'", body)
	}
	if !strings.Contains(body, "fallback") {
		t.Errorf("body = %q, want to contain 'fallback'", body)
	}
}

// --- CORS tests ---

func TestHTMXHandler_CORSPreflight(t *testing.T) {
	b := New(WithCSRF(false), WithCORS(true), WithAllowedOrigins("https://example.com"))

	err := b.Register("test.cors", func(ctx Context) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodOptions, "/api/bridge/fn/test.cors", nil)
	r.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Errorf("ACAO = %q, want %q", got, "https://example.com")
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("ACAC = %q, want %q", got, "true")
	}
}

func TestHTMXHandler_CORSDisallowedOrigin(t *testing.T) {
	b := New(WithCSRF(false), WithCORS(true), WithAllowedOrigins("https://allowed.com"))

	err := b.Register("test.cors2", func(ctx Context) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.cors2", nil)
	r.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("ACAO = %q, want empty (origin not allowed)", got)
	}
}

func TestHTMXHandler_CORSWildcard(t *testing.T) {
	b := New(WithCSRF(false), WithCORS(true), WithAllowedOrigins("*"))

	err := b.Register("test.corsWild", func(ctx Context) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.corsWild", nil)
	r.Header.Set("Origin", "https://any-origin.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://any-origin.com" {
		t.Errorf("ACAO = %q, want %q", got, "https://any-origin.com")
	}
}

func TestHTMXHandler_CORSNoOriginHeader(t *testing.T) {
	b := New(WithCSRF(false), WithCORS(true), WithAllowedOrigins("*"))

	err := b.Register("test.corsNoOrigin", func(ctx Context) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.corsNoOrigin", nil)
	// No Origin header
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("ACAO = %q, want empty (no Origin header)", got)
	}
}

// --- isHTMXRequest test ---

func TestIsHTMXRequest(t *testing.T) {
	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	if isHTMXRequest(r1) {
		t.Error("expected false for regular request")
	}

	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	r2.Header.Set("HX-Request", "true")
	if !isHTMXRequest(r2) {
		t.Error("expected true for HTMX request")
	}

	r3 := httptest.NewRequest(http.MethodGet, "/", nil)
	r3.Header.Set("HX-Request", "false")
	if isHTMXRequest(r3) {
		t.Error("expected false for HX-Request=false")
	}
}

// --- NewHTMXHandler tests ---

func TestNewHTMXHandler_DefaultPrefix(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "")

	if h.prefix != "/api/bridge/fn/" {
		t.Errorf("prefix = %q, want %q", h.prefix, "/api/bridge/fn/")
	}
}

func TestNewHTMXHandler_AddTrailingSlash(t *testing.T) {
	b := New(WithCSRF(false))
	h := NewHTMXHandler(b, "/custom/path")

	if h.prefix != "/custom/path/" {
		t.Errorf("prefix = %q, want %q", h.prefix, "/custom/path/")
	}
}

// --- Execution error response tests ---

func TestHTMXHandler_FunctionError(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.internalErr", func(ctx Context) error {
		return fmt.Errorf("internal failure")
	})
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.internalErr", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestHTMXHandler_BridgeErrorMapping(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.forbidden", func(ctx Context) error {
		return NewError(ErrCodeForbidden, "not allowed")
	})
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.forbidden", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

// --- POST with all HTMX headers test ---

func TestHTMXHandler_AllHTMXHeaders(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.allHeaders", func(ctx Context) error {
		return nil
	},
		WithHTMXTrigger("ev1", "ev2"),
		WithHTMXRedirect("/new-page"),
		WithHTMXReswap("innerHTML"),
		WithHTMXRetarget("#target"),
	)
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.allHeaders", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("HX-Trigger"); got != "ev1, ev2" {
		t.Errorf("HX-Trigger = %q, want %q", got, "ev1, ev2")
	}
	if got := w.Header().Get("HX-Redirect"); got != "/new-page" {
		t.Errorf("HX-Redirect = %q, want %q", got, "/new-page")
	}
	if got := w.Header().Get("HX-Reswap"); got != "innerHTML" {
		t.Errorf("HX-Reswap = %q, want %q", got, "innerHTML")
	}
	if got := w.Header().Get("HX-Retarget"); got != "#target" {
		t.Errorf("HX-Retarget = %q, want %q", got, "#target")
	}
}

// --- Trailing slash in function name ---

func TestHTMXHandler_TrailingSlashFunctionName(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.slash", func(ctx Context) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.slash/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (trailing slash should be trimmed)", w.Code, http.StatusOK)
	}
}

// --- Auth check tests ---

func TestHTMXHandler_AuthRequired_NoUser(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.authRequired", func(ctx Context) error {
		return nil
	}, RequireAuth())
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.authRequired", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHTMXHandler_AuthRequired_WithRoles(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.adminOnly", func(ctx Context) error {
		return nil
	}, RequireAuth(), RequireRoles("admin"))
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	// Without any user middleware, auth always fails for RequireAuth functions
	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.adminOnly", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d (no user = unauthorized)", w.Code, http.StatusUnauthorized)
	}
}

// --- CSRF check tests ---

func TestHTMXHandler_CSRFRequired_PostFails(t *testing.T) {
	b := New(WithCSRF(true)) // CSRF enabled

	err := b.Register("test.csrfPost", func(ctx Context, params testInput) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	body := strings.NewReader(`{"name":"test"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/bridge/fn/test.csrfPost", body)
	r.Header.Set("Content-Type", "application/json")
	// No CSRF token
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d (CSRF should fail)", w.Code, http.StatusForbidden)
	}
}

func TestHTMXHandler_CSRFSkipped_ForGET(t *testing.T) {
	b := New(WithCSRF(true)) // CSRF enabled

	err := b.Register("test.csrfGet", func(ctx Context) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.csrfGet", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	// GET should skip CSRF check
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (GET skips CSRF)", w.Code, http.StatusOK)
	}
}

// --- Rate limiting tests ---

func TestHTMXHandler_RateLimited(t *testing.T) {
	b := New(WithCSRF(false), WithDefaultRateLimit(1))

	err := b.Register("test.rateLimit", func(ctx Context) error {
		return nil
	}, WithRateLimit(1)) // flag this function for rate limiting
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	// The rate limiter uses a token bucket with burst = rate*2 = 2.
	// New bucket starts with burst-1 = 1 token, plus the first call is free.
	// So we need to exhaust: first call (free) + 1 token = 2 allowed, third blocked.
	for i := range 3 {
		r := httptest.NewRequest(http.MethodGet, "/api/bridge/fn/test.rateLimit", nil)
		r.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if i < 2 {
			if w.Code != http.StatusOK {
				t.Errorf("request %d: status = %d, want %d", i+1, w.Code, http.StatusOK)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("request %d: status = %d, want %d", i+1, w.Code, http.StatusTooManyRequests)
			}
		}
	}
}

// --- Parse error tests ---

func TestHTMXHandler_InvalidJSONBody(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.parseErr", func(ctx Context, params testInput) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	handler := b.HTMXHandler()

	body := strings.NewReader(`{invalid json}`)
	r := httptest.NewRequest(http.MethodPost, "/api/bridge/fn/test.parseErr", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
