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

func TestManagerWithEmbedFS(t *testing.T) {
	cfg := Config{
		PublicDir:  "testdata",
		IsDev:      false,
		FileSystem: testFS,
	}

	m := NewManager(cfg)

	if m.publicDir != "testdata" {
		t.Errorf("Expected publicDir to be 'testdata', got '%s'", m.publicDir)
	}

	if m.fileSystem == nil {
		t.Error("Expected fileSystem to be set")
	}
}

func TestEmbeddedHandler(t *testing.T) {
	m := NewManager(Config{
		PublicDir:  "",
		IsDev:      true,
		FileSystem: testFS,
	})

	handler := m.Handler()

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
	m := NewManager(Config{
		PublicDir:  "",
		IsDev:      true,
		FileSystem: testFS,
	})

	handler := m.Handler()

	req := httptest.NewRequest(http.MethodGet, "/static/nonexistent.txt", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestEmbeddedHandler_PathTraversal(t *testing.T) {
	m := NewManager(Config{
		PublicDir:  "",
		IsDev:      true,
		FileSystem: testFS,
	})

	handler := m.Handler()

	req := httptest.NewRequest(http.MethodGet, "/static/../../../etc/passwd", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for path traversal, got %d", w.Code)
	}
}

func TestEmbeddedHandler_CacheHeaders_Production(t *testing.T) {
	m := NewManager(Config{
		PublicDir:  "",
		IsDev:      false,
		FileSystem: testFS,
	})

	handler := m.Handler()

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

func TestFingerprintAllWithEmbedFS(t *testing.T) {
	m := NewManager(Config{
		PublicDir:  "",
		IsDev:      false,
		FileSystem: testFS,
	})

	if err := m.FingerprintAll(); err != nil {
		t.Fatalf("FingerprintAll failed: %v", err)
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
