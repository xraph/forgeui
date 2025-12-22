package bridge

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPHandler_SingleRequest(t *testing.T) {
	b := New(WithCSRF(false)) // Disable CSRF for testing

	// Register test function
	b.Register("greet", func(ctx Context, input struct {
		Name string `json:"name"`
	}) (struct {
		Message string `json:"message"`
	}, error) {
		return struct {
			Message string `json:"message"`
		}{Message: "Hello, " + input.Name}, nil
	})

	handler := NewHTTPHandler(b)

	// Create request
	req := Request{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "greet",
		Params:  json.RawMessage(`{"name":"World"}`),
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/bridge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("response has error: %v", resp.Error)
	}

	if resp.ID != "1" {
		t.Errorf("response ID = %v, want 1", resp.ID)
	}
}

func TestHTTPHandler_BatchRequest(t *testing.T) {
	b := New(WithCSRF(false))

	// Register test functions
	b.Register("double", func(ctx Context, input struct {
		N int `json:"n"`
	}) (struct {
		Result int `json:"result"`
	}, error) {
		return struct {
			Result int `json:"result"`
		}{Result: input.N * 2}, nil
	})

	handler := NewHTTPHandler(b)

	// Create batch request
	batchReq := BatchRequest{
		{JSONRPC: "2.0", ID: "1", Method: "double", Params: json.RawMessage(`{"n":5}`)},
		{JSONRPC: "2.0", ID: "2", Method: "double", Params: json.RawMessage(`{"n":10}`)},
	}

	body, err := json.Marshal(batchReq)
	if err != nil {
		t.Fatalf("Failed to marshal batch request: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/bridge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	var responses BatchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &responses); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responses) != 2 {
		t.Errorf("len(responses) = %d, want 2", len(responses))
	}
}

func TestHTTPHandler_MethodNotFound(t *testing.T) {
	b := New(WithCSRF(false))
	handler := NewHTTPHandler(b)

	req := Request{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "nonexistent",
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/bridge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httpReq)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Error == nil {
		t.Error("expected error for method not found")
	}

	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeMethodNotFound)
	}
}

func TestHTTPHandler_InvalidRequest(t *testing.T) {
	b := New(WithCSRF(false))
	handler := NewHTTPHandler(b)

	// Send invalid JSON
	httpReq := httptest.NewRequest(http.MethodPost, "/api/bridge", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httpReq)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Error == nil {
		t.Error("expected error for invalid request")
	}
}

func TestHTTPHandler_MethodNotAllowed(t *testing.T) {
	b := New(WithCSRF(false))
	handler := NewHTTPHandler(b)

	httpReq := httptest.NewRequest(http.MethodGet, "/api/bridge", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, httpReq)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Error == nil {
		t.Error("expected error for GET request")
	}
}

func TestHTTPHandler_CORS(t *testing.T) {
	b := New(
		WithCSRF(false),
		WithCORS(true),
		WithAllowedOrigins("http://localhost:3000"),
	)
	handler := NewHTTPHandler(b)

	// OPTIONS request
	httpReq := httptest.NewRequest(http.MethodOptions, "/api/bridge", nil)
	httpReq.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", w.Code, http.StatusOK)
	}

	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Errorf("CORS origin = %s, want http://localhost:3000", origin)
	}
}
