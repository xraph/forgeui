package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPageContext_Param(t *testing.T) {
	ctx := &PageContext{
		Params: Params{"id": "123", "name": "test"},
	}

	if ctx.Param("id") != "123" {
		t.Errorf("Expected id=123, got %s", ctx.Param("id"))
	}

	if ctx.Param("name") != "test" {
		t.Errorf("Expected name=test, got %s", ctx.Param("name"))
	}

	if ctx.Param("missing") != "" {
		t.Error("Expected empty string for missing param")
	}
}

func TestPageContext_ParamInt(t *testing.T) {
	ctx := &PageContext{
		Params: Params{"id": "123", "invalid": "abc"},
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		t.Fatalf("ParamInt() error = %v", err)
	}

	if id != 123 {
		t.Errorf("Expected 123, got %d", id)
	}

	_, err = ctx.ParamInt("invalid")
	if err == nil {
		t.Error("Expected error for invalid integer")
	}
}

func TestPageContext_Query(t *testing.T) {
	req := httptest.NewRequest(MethodGet, "/test?q=golang&page=2", nil)
	ctx := &PageContext{
		Request: req,
	}

	if ctx.Query("q") != "golang" {
		t.Errorf("Expected q=golang, got %s", ctx.Query("q"))
	}

	if ctx.Query("page") != "2" {
		t.Errorf("Expected page=2, got %s", ctx.Query("page"))
	}

	if ctx.Query("missing") != "" {
		t.Error("Expected empty string for missing query param")
	}
}

func TestPageContext_QueryDefault(t *testing.T) {
	req := httptest.NewRequest(MethodGet, "/test?page=2", nil)
	ctx := &PageContext{
		Request: req,
	}

	if ctx.QueryDefault("page", "1") != "2" {
		t.Error("Expected actual value when present")
	}

	if ctx.QueryDefault("missing", "default") != "default" {
		t.Error("Expected default value when missing")
	}
}

func TestPageContext_QueryInt(t *testing.T) {
	req := httptest.NewRequest(MethodGet, "/test?page=2&invalid=abc", nil)
	ctx := &PageContext{
		Request: req,
	}

	page, err := ctx.QueryInt("page")
	if err != nil {
		t.Fatalf("QueryInt() error = %v", err)
	}

	if page != 2 {
		t.Errorf("Expected 2, got %d", page)
	}

	_, err = ctx.QueryInt("invalid")
	if err == nil {
		t.Error("Expected error for invalid integer")
	}
}

func TestPageContext_QueryBool(t *testing.T) {
	tests := []struct {
		url      string
		key      string
		expected bool
	}{
		{"/test?active=true", "active", true},
		{"/test?active=1", "active", true},
		{"/test?active=yes", "active", true},
		{"/test?active=false", "active", false},
		{"/test?active=0", "active", false},
		{"/test", "active", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req := httptest.NewRequest(MethodGet, tt.url, nil)
			ctx := &PageContext{Request: req}

			result := ctx.QueryBool(tt.key)
			if result != tt.expected {
				t.Errorf("QueryBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPageContext_Header(t *testing.T) {
	req := httptest.NewRequest(MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")

	ctx := &PageContext{Request: req}

	if ctx.Header("Authorization") != "Bearer token123" {
		t.Error("Expected Authorization header")
	}
}

func TestPageContext_SetHeader(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := &PageContext{ResponseWriter: w}

	ctx.SetHeader("X-Custom", "value")

	if w.Header().Get("X-Custom") != "value" {
		t.Error("Expected X-Custom header to be set")
	}
}

func TestPageContext_Cookie(t *testing.T) {
	req := httptest.NewRequest(MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})

	ctx := &PageContext{Request: req}

	cookie, err := ctx.Cookie("session")
	if err != nil {
		t.Fatalf("Cookie() error = %v", err)
	}

	if cookie.Value != "abc123" {
		t.Errorf("Expected cookie value abc123, got %s", cookie.Value)
	}
}

func TestPageContext_SetCookie(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := &PageContext{ResponseWriter: w}

	cookie := &http.Cookie{
		Name:  "session",
		Value: "xyz789",
	}

	ctx.SetCookie(cookie)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Expected 1 cookie, got %d", len(cookies))
	}

	if cookies[0].Value != "xyz789" {
		t.Errorf("Expected cookie value xyz789, got %s", cookies[0].Value)
	}
}

func TestPageContext_SetGet(t *testing.T) {
	ctx := &PageContext{
		values: make(map[string]interface{}),
	}

	ctx.Set("user_id", 123)
	ctx.Set("username", "john")

	val, ok := ctx.Get("user_id")
	if !ok {
		t.Error("Expected user_id to exist")
	}

	if val.(int) != 123 {
		t.Errorf("Expected 123, got %v", val)
	}

	_, ok = ctx.Get("missing")
	if ok {
		t.Error("Expected missing key to not exist")
	}
}

func TestPageContext_GetString(t *testing.T) {
	ctx := &PageContext{
		values: make(map[string]interface{}),
	}

	ctx.Set("name", "test")
	ctx.Set("number", 123)

	if ctx.GetString("name") != "test" {
		t.Error("Expected 'test'")
	}

	if ctx.GetString("number") != "" {
		t.Error("Expected empty string for non-string value")
	}

	if ctx.GetString("missing") != "" {
		t.Error("Expected empty string for missing key")
	}
}

func TestPageContext_GetInt(t *testing.T) {
	ctx := &PageContext{
		values: make(map[string]interface{}),
	}

	ctx.Set("count", 42)
	ctx.Set("name", "test")

	if ctx.GetInt("count") != 42 {
		t.Error("Expected 42")
	}

	if ctx.GetInt("name") != 0 {
		t.Error("Expected 0 for non-int value")
	}

	if ctx.GetInt("missing") != 0 {
		t.Error("Expected 0 for missing key")
	}
}

func TestPageContext_Method(t *testing.T) {
	req := httptest.NewRequest(MethodPost, "/test", nil)
	ctx := &PageContext{Request: req}

	if ctx.Method() != MethodPost {
		t.Errorf("Expected POST, got %s", ctx.Method())
	}
}

func TestPageContext_Path(t *testing.T) {
	req := httptest.NewRequest(MethodGet, "/users/123", nil)
	ctx := &PageContext{Request: req}

	if ctx.Path() != "/users/123" {
		t.Errorf("Expected /users/123, got %s", ctx.Path())
	}
}

func TestPageContext_ClientIP(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		remoteAddr string
		expected string
	}{
		{
			name:     "X-Forwarded-For",
			headers:  map[string]string{"X-Forwarded-For": "1.2.3.4"},
			expected: "1.2.3.4",
		},
		{
			name:     "X-Real-IP",
			headers:  map[string]string{"X-Real-IP": "5.6.7.8"},
			expected: "5.6.7.8",
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: "9.10.11.12:1234",
			expected:   "9.10.11.12:1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(MethodGet, "/test", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			if tt.remoteAddr != "" {
				req.RemoteAddr = tt.remoteAddr
			}

			ctx := &PageContext{Request: req}

			if ctx.ClientIP() != tt.expected {
				t.Errorf("ClientIP() = %q, want %q", ctx.ClientIP(), tt.expected)
			}
		})
	}
}

