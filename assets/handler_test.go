package assets

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandler_ServeFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.css")

	content := []byte("body { color: red; }")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	handler := m.Handler()

	req := httptest.NewRequest(http.MethodGet, "/static/test.css", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != string(content) {
		t.Errorf("Expected body '%s', got '%s'", content, w.Body.String())
	}
}

func TestHandler_404(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	handler := m.Handler()

	req := httptest.NewRequest(http.MethodGet, "/static/nonexistent.css", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestHandler_CacheHeaders_Dev(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.css")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	handler := m.Handler()

	req := httptest.NewRequest(http.MethodGet, "/static/test.css", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	cacheControl := w.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "max-age=3600") {
		t.Errorf("Expected moderate cache in dev mode, got: %s", cacheControl)
	}
}

func TestHandler_CacheHeaders_Production_Fingerprinted(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.css")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	handler := m.Handler()

	// Request with fingerprinted URL
	req := httptest.NewRequest(http.MethodGet, "/static/test.abc12345.css", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	cacheControl := w.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "immutable") {
		t.Errorf("Expected immutable cache for fingerprinted assets, got: %s", cacheControl)
	}

	if !strings.Contains(cacheControl, "max-age=31536000") {
		t.Errorf("Expected 1-year cache for fingerprinted assets, got: %s", cacheControl)
	}
}

func TestHandler_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	handler := m.Handler()

	// Attempt path traversal
	req := httptest.NewRequest(http.MethodGet, "/static/../../../etc/passwd", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for path traversal, got %d", w.Code)
	}
}

func TestHandler_DirectoryRequest(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory
	subdir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	handler := m.Handler()

	req := httptest.NewRequest(http.MethodGet, "/static/css", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for directory request, got %d", w.Code)
	}
}
