package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	g "maragu.dev/gomponents"
)

func TestLogger(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	}

	middleware := Logger()
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	_, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRecovery(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		panic("test panic")
	}

	middleware := Recovery()
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	node, err := wrapped(ctx)
	if err == nil {
		t.Error("Expected error from panic recovery")
	}

	if node == nil {
		t.Error("Expected error page node")
	}

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestCORS(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	}

	middleware := CORS("*")
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	_, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header to be set")
	}
}

func TestCORS_Preflight(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	}

	middleware := CORS("*")
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodOptions, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	node, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if node != nil {
		t.Error("Expected nil node for OPTIONS request")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRequestID(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	}

	middleware := RequestID()
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
		values:         make(map[string]interface{}),
	}

	_, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	if ctx.GetString("request_id") == "" {
		t.Error("Expected request_id in context")
	}
}

func TestBasicAuth_Success(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("Authorized"), nil
	}

	middleware := BasicAuth("admin", "secret")
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	req.SetBasicAuth("admin", "secret")
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	node, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if node == nil {
		t.Error("Expected node to be returned")
	}
}

func TestBasicAuth_Failure(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("Authorized"), nil
	}

	middleware := BasicAuth("admin", "secret")
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	req.SetBasicAuth("admin", "wrong")
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	_, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRequireMethod(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	}

	middleware := RequireMethod(MethodGet, MethodPost)
	wrapped := middleware(handler)

	// Test allowed method
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	_, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test disallowed method
	req = httptest.NewRequest(MethodDelete, "/test", nil)
	w = httptest.NewRecorder()
	ctx = &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	_, err = wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestSetHeader(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	}

	middleware := SetHeader("X-Custom", "test-value")
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	_, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if w.Header().Get("X-Custom") != "test-value" {
		t.Error("Expected X-Custom header to be set")
	}
}

func TestChain(t *testing.T) {
	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	}

	m1 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.SetHeader("X-M1", "1")
			return next(ctx)
		}
	}

	m2 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.SetHeader("X-M2", "2")
			return next(ctx)
		}
	}

	middleware := Chain(m1, m2)
	wrapped := middleware(handler)

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := &PageContext{
		ResponseWriter: w,
		Request:        req,
	}

	_, err := wrapped(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if w.Header().Get("X-M1") != "1" {
		t.Error("Expected X-M1 header from first middleware")
	}

	if w.Header().Get("X-M2") != "2" {
		t.Error("Expected X-M2 header from second middleware")
	}
}

