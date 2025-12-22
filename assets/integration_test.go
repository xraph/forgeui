package assets

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIntegration_DevMode tests the complete flow in development mode
func TestIntegration_DevMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testCSS := filepath.Join(tmpDir, "app.css")

	content := []byte("body { color: red; }")
	if err := os.WriteFile(testCSS, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create manager in dev mode
	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	// Get URL (should not be fingerprinted)
	url := m.URL("app.css")
	if url != "/static/app.css" {
		t.Errorf("Expected non-fingerprinted URL in dev mode, got: %s", url)
	}

	// Test serving the file
	handler := m.Handler()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if string(body) != string(content) {
		t.Errorf("Expected body '%s', got '%s'", content, body)
	}

	// Check cache headers (moderate caching in dev)
	cacheControl := w.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "max-age=3600") {
		t.Errorf("Expected moderate cache, got: %s", cacheControl)
	}
}

// TestIntegration_ProductionMode tests the complete flow in production mode
func TestIntegration_ProductionMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testCSS := filepath.Join(tmpDir, "app.css")

	content := []byte("body { color: blue; }")
	if err := os.WriteFile(testCSS, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create manager in production mode
	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	// Get URL (should be fingerprinted)
	url := m.URL("app.css")
	if url == "/static/app.css" {
		t.Error("Expected fingerprinted URL in production mode")
	}

	// Should match pattern: /static/app.{hash}.css
	if !strings.Contains(url, "/static/app.") || !strings.HasSuffix(url, ".css") {
		t.Errorf("URL doesn't match expected pattern: %s", url)
	}

	// Test serving the fingerprinted file
	handler := m.Handler()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if string(body) != string(content) {
		t.Errorf("Expected body '%s', got '%s'", content, body)
	}

	// Check cache headers (immutable in production)
	cacheControl := w.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "immutable") {
		t.Errorf("Expected immutable cache, got: %s", cacheControl)
	}

	if !strings.Contains(cacheControl, "max-age=31536000") {
		t.Errorf("Expected 1-year cache, got: %s", cacheControl)
	}
}

// TestIntegration_ManifestFlow tests manifest generation and loading
func TestIntegration_ManifestFlow(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"app.css":   "body { color: red; }",
		"main.js":   "console.log('test');",
		"style.css": "h1 { font-size: 2rem; }",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	// Create manager and generate manifest
	m1 := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	manifest, err := GenerateManifest(m1)
	if err != nil {
		t.Fatalf("Failed to generate manifest: %v", err)
	}

	// Save manifest
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if err := manifest.Save(manifestPath); err != nil {
		t.Fatalf("Failed to save manifest: %v", err)
	}

	// Create new manager with manifest
	m2 := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
		Manifest:  manifestPath,
	})

	// Verify URLs match manifest
	for originalFile, fingerprintedFile := range manifest {
		url := m2.URL(originalFile)

		expectedURL := "/static/" + fingerprintedFile
		if url != expectedURL {
			t.Errorf("For %s, expected URL %s, got %s", originalFile, expectedURL, url)
		}
	}

	// Test serving via manifest
	handler := m2.Handler()
	for originalFile := range manifest {
		url := m2.URL(originalFile)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Failed to serve %s: status %d", originalFile, w.Code)
		}

		expectedContent := files[originalFile]

		body, _ := io.ReadAll(w.Body)
		if string(body) != expectedContent {
			t.Errorf("Content mismatch for %s", originalFile)
		}
	}
}

// TestIntegration_NestedPaths tests handling of nested directories
func TestIntegration_NestedPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	cssDir := filepath.Join(tmpDir, "css")
	jsDir := filepath.Join(tmpDir, "js")

	if err := os.MkdirAll(cssDir, 0755); err != nil {
		t.Fatalf("Failed to create css dir: %v", err)
	}

	if err := os.MkdirAll(jsDir, 0755); err != nil {
		t.Fatalf("Failed to create js dir: %v", err)
	}

	// Create files in nested directories
	files := map[string]string{
		"css/main.css": "body { margin: 0; }",
		"js/bundle.js": "console.log('app');",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", name, err)
		}
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	// Test URL generation for nested paths
	for file := range files {
		url := m.URL(file)

		// Should contain directory structure
		if !strings.Contains(url, filepath.Dir(file)) {
			t.Errorf("URL %s doesn't preserve directory structure", url)
		}

		// Test serving
		handler := m.Handler()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Failed to serve %s: status %d", file, w.Code)
		}
	}
}

// TestIntegration_ConcurrentAccess tests thread-safety
func TestIntegration_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple test files
	for i := 1; i <= 5; i++ {
		name := filepath.Join(tmpDir, "file"+string(rune(i+48))+".css")
		if err := os.WriteFile(name, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	handler := m.Handler()

	// Simulate concurrent requests
	done := make(chan bool)
	numGoroutines := 20

	for i := range numGoroutines {
		go func(id int) {
			fileNum := (id % 5) + 1
			file := "file" + string(rune(fileNum+48)) + ".css"

			// Get URL multiple times
			for range 10 {
				url := m.URL(file)

				// Make request
				req := httptest.NewRequest(http.MethodGet, url, nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("Goroutine %d: Failed to serve %s", id, file)
				}
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for range numGoroutines {
		<-done
	}
}
