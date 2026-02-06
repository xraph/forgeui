package assets

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
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

func TestHandler_MIMETypes(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		content      string
		expectedMIME string
	}{
		{
			name:         "CSS file",
			filename:     "styles.css",
			content:      "body { color: red; }",
			expectedMIME: "text/css; charset=utf-8",
		},
		{
			name:         "JavaScript file",
			filename:     "script.js",
			content:      "console.log('test');",
			expectedMIME: "text/javascript; charset=utf-8",
		},
		{
			name:         "JSON file",
			filename:     "data.json",
			content:      `{"key": "value"}`,
			expectedMIME: "application/json",
		},
		{
			name:         "PNG image",
			filename:     "image.png",
			content:      "fake-png-data",
			expectedMIME: "image/png",
		},
		{
			name:         "SVG image",
			filename:     "icon.svg",
			content:      "<svg></svg>",
			expectedMIME: "image/svg+xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test file
			testFile := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			m := NewManager(Config{
				PublicDir: tmpDir,
				IsDev:     true,
			})

			handler := m.Handler()

			req := httptest.NewRequest(http.MethodGet, "/static/"+tt.filename, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedMIME {
				t.Errorf("Expected Content-Type '%s', got '%s'", tt.expectedMIME, contentType)
			}
		})
	}
}

func TestHandler_CustomFileSystem(t *testing.T) {
	// Create an in-memory filesystem (simulating embed.FS usage)
	memFS := fstest.MapFS{
		"css/output.css": {
			Data: []byte("body { color: blue; }"),
		},
		"js/forge-bridge.js": {
			Data: []byte("console.log('bridge loaded');"),
		},
	}

	// Create manager with a dummy public dir (won't be used)
	m := NewManager(Config{
		PublicDir: "/nonexistent/path",
		IsDev:     true,
	})

	// Set custom filesystem - this is the key for library usage
	m.SetFileSystem(memFS)

	handler := m.Handler()

	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedMIME string
		expectedBody string
	}{
		{
			name:         "CSS file from custom FS",
			path:         "/static/css/output.css",
			expectedCode: http.StatusOK,
			expectedMIME: "text/css; charset=utf-8",
			expectedBody: "body { color: blue; }",
		},
		{
			name:         "JS file from custom FS",
			path:         "/static/js/forge-bridge.js",
			expectedCode: http.StatusOK,
			expectedMIME: "text/javascript; charset=utf-8",
			expectedBody: "console.log('bridge loaded');",
		},
		{
			name:         "Non-existent file",
			path:         "/static/missing.css",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				if contentType != tt.expectedMIME {
					t.Errorf("Expected Content-Type '%s', got '%s'", tt.expectedMIME, contentType)
				}

				if w.Body.String() != tt.expectedBody {
					t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
				}
			}
		})
	}
}

func TestHandler_WithEmbedFS(t *testing.T) {
	// Simulate using forgeui as a library with embed.FS
	// This tests the real-world scenario where files don't exist on disk
	embedFS := fstest.MapFS{
		"dist/css/styles.css": {
			Data: []byte(".container { max-width: 1200px; }"),
		},
		"dist/js/app.js": {
			Data: []byte("const app = { init: () => {} };"),
		},
	}

	// Create a subdirectory FS to serve only the "dist" folder
	distFS, err := fs.Sub(embedFS, "dist")
	if err != nil {
		t.Fatalf("Failed to create sub filesystem: %v", err)
	}

	m := NewManager(Config{
		PublicDir: "dist", // This path doesn't exist on disk
		IsDev:     false,
	})

	m.SetFileSystem(distFS)

	handler := m.Handler()

	// Test CSS file
	req := httptest.NewRequest(http.MethodGet, "/static/css/styles.css", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for CSS, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "text/css; charset=utf-8" {
		t.Errorf("Expected CSS MIME type, got '%s'", ct)
	}

	// Test JS file
	req = httptest.NewRequest(http.MethodGet, "/static/js/app.js", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for JS, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "text/javascript; charset=utf-8" {
		t.Errorf("Expected JS MIME type, got '%s'", ct)
	}
}
