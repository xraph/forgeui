package bridge

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// End-to-end integration test
func TestBridge_Integration(t *testing.T) {
	// Create bridge
	b := New(WithCSRF(false))

	// Register functions
	_ = b.Register("echo", func(ctx Context, input struct {
		Message string `json:"message"`
	}) (struct {
		Echo string `json:"echo"`
	}, error) {
		return struct {
			Echo string `json:"echo"`
		}{Echo: input.Message}, nil
	})

	_ = b.Register("add", func(ctx Context, input struct {
		A int `json:"a"`
		B int `json:"b"`
	}) (struct {
		Sum int `json:"sum"`
	}, error) {
		return struct {
			Sum int `json:"sum"`
		}{Sum: input.A + input.B}, nil
	}, WithDescription("Add two numbers"))

	// Test HTTP handler
	t.Run("HTTP Single Request", func(t *testing.T) {
		handler := NewHTTPHandler(b)

		req := Request{
			JSONRPC: "2.0",
			ID:      "1",
			Method:  "echo",
			Params:  json.RawMessage(`{"message":"Hello"}`),
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/bridge", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		}

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Error != nil {
			t.Errorf("unexpected error: %v", resp.Error)
		}
	})

	// Test batch request
	t.Run("HTTP Batch Request", func(t *testing.T) {
		handler := NewHTTPHandler(b)

		batch := BatchRequest{
			{JSONRPC: "2.0", ID: "1", Method: "add", Params: json.RawMessage(`{"a":5,"b":3}`)},
			{JSONRPC: "2.0", ID: "2", Method: "add", Params: json.RawMessage(`{"a":10,"b":20}`)},
		}

		body, _ := json.Marshal(batch)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/bridge", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, httpReq)

		var responses BatchResponse
		_ = json.Unmarshal(w.Body.Bytes(), &responses)

		if len(responses) != 2 {
			t.Errorf("len(responses) = %d, want 2", len(responses))
		}

		for i, resp := range responses {
			if resp.Error != nil {
				t.Errorf("response[%d] has error: %v", i, resp.Error)
			}
		}
	})

	// Test introspection
	t.Run("Introspection", func(t *testing.T) {
		handler := NewIntrospectionHandler(b)

		httpReq := httptest.NewRequest(http.MethodGet, "/api/bridge/functions", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
		}

		var result map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &result)

		count, ok := result["count"].(float64)
		if !ok || count != 2 {
			t.Errorf("function count = %v, want 2", count)
		}
	})

	// Test hooks
	t.Run("Hooks", func(t *testing.T) {
		var called atomic.Bool

		done := make(chan struct{})

		b.GetHooks().Register(BeforeCall, func(ctx Context, data HookData) {
			called.Store(true)
			close(done)
		})

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		ctx := NewContext(req)

		_, _ = b.Call(ctx, "echo", json.RawMessage(`{"message":"test"}`))

		// Wait for async hook to complete
		select {
		case <-done:
			if !called.Load() {
				t.Error("Hook was not called")
			}
		case <-time.After(1 * time.Second):
			t.Error("Hook timed out")
		}
	})

	// Test caching
	t.Run("Cache", func(t *testing.T) {
		cache := NewMemoryCache()

		key := generateCacheKey("test", json.RawMessage(`{"a":1}`))
		cache.Set(key, "result1", 1*time.Minute)

		val, ok := cache.Get(key)
		if !ok || val != "result1" {
			t.Errorf("cache Get() = %v, %v, want result1, true", val, ok)
		}
	})

	// Test rate limiting
	t.Run("RateLimiter", func(t *testing.T) {
		rl := NewRateLimiter(2, 2)

		// First 2 should pass
		if !rl.Allow("test") {
			t.Error("first request should pass")
		}

		if !rl.Allow("test") {
			t.Error("second request should pass")
		}

		// Third should fail
		if rl.Allow("test") {
			t.Error("third request should fail")
		}
	})
}

// Test full workflow with authentication
func TestBridge_AuthWorkflow(t *testing.T) {
	b := New(WithCSRF(false))

	// Register protected function
	_ = b.Register("protected", func(ctx Context, input struct{}) (struct {
		Message string `json:"message"`
	}, error) {
		return struct {
			Message string `json:"message"`
		}{Message: "Success"}, nil
	}, RequireAuth())

	handler := NewHTTPHandler(b)

	// Without authentication
	t.Run("Without Auth", func(t *testing.T) {
		req := Request{
			JSONRPC: "2.0",
			ID:      "1",
			Method:  "protected",
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/bridge", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, httpReq)

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp.Error == nil {
			t.Error("expected error for unauthenticated request")
		}

		if resp.Error.Code != ErrCodeUnauthorized {
			t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeUnauthorized)
		}
	})
}
