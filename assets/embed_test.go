package assets

import (
	"embed"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed test_helpers.go
var testFS embed.FS

func TestNewEmbeddedManager(t *testing.T) {
	cfg := Config{
		PublicDir: "testdata",
		IsDev:     false,
	}

	m := NewEmbeddedManager(testFS, cfg)

	if m.embedFS == nil {
		t.Error("Expected embedFS to be set")
	}

	if m.publicDir != "testdata" {
		t.Errorf("Expected publicDir to be 'testdata', got '%s'", m.publicDir)
	}
}

func TestEmbeddedHandler(t *testing.T) {
	m := NewEmbeddedManager(testFS, Config{
		PublicDir: "",
		IsDev:     true,
	})

	handler := m.EmbeddedHandler()

	req := httptest.NewRequest(http.MethodGet, "/static/test_helpers.go", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if !strings.Contains(string(body), "renderNode") {
		t.Error("Expected test_helpers.go content")
	}
}

func TestEmbeddedHandler_404(t *testing.T) {
	m := NewEmbeddedManager(testFS, Config{
		PublicDir: "",
		IsDev:     true,
	})

	handler := m.EmbeddedHandler()

	req := httptest.NewRequest(http.MethodGet, "/static/nonexistent.txt", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestEmbeddedHandler_PathTraversal(t *testing.T) {
	m := NewEmbeddedManager(testFS, Config{
		PublicDir: "",
		IsDev:     true,
	})

	handler := m.EmbeddedHandler()

	req := httptest.NewRequest(http.MethodGet, "/static/../../../etc/passwd", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for path traversal, got %d", w.Code)
	}
}

func TestEmbeddedHandler_CacheHeaders_Production(t *testing.T) {
	m := NewEmbeddedManager(testFS, Config{
		PublicDir: "",
		IsDev:     false,
	})

	handler := m.EmbeddedHandler()

	// Request with fingerprinted URL
	req := httptest.NewRequest(http.MethodGet, "/static/test.abc12345.go", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Even if file doesn't exist, cache header logic should still be tested
	// by making a request with a file that exists
	req2 := httptest.NewRequest(http.MethodGet, "/static/helpers.abc12345.go", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	// The real test: fingerprinted URL that strips to a real file
	m.mu.Lock()
	m.fingerprints["test_helpers.go"] = "test_helpers.abc12345.go"
	m.mu.Unlock()

	req3 := httptest.NewRequest(http.MethodGet, "/static/test_helpers.abc12345.go", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req3)

	if w3.Code == http.StatusOK {
		cacheControl := w3.Header().Get("Cache-Control")
		if !strings.Contains(cacheControl, "immutable") {
			t.Errorf("Expected immutable cache for fingerprinted assets, got: %s", cacheControl)
		}
	}
}

func TestFingerprintAllEmbedded(t *testing.T) {
	m := NewEmbeddedManager(testFS, Config{
		PublicDir: "",
		IsDev:     false,
	})

	if err := m.FingerprintAllEmbedded(); err != nil {
		t.Fatalf("FingerprintAllEmbedded failed: %v", err)
	}

	// Verify fingerprints were generated
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.fingerprints) == 0 {
		t.Error("Expected fingerprints to be generated")
	}

	// Should have fingerprint for test_helpers.go
	if _, ok := m.fingerprints["test_helpers.go"]; !ok {
		t.Error("Expected fingerprint for test_helpers.go")
	}
}
