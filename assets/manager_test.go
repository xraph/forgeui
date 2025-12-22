package assets

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestNewManager(t *testing.T) {
	cfg := Config{
		PublicDir: "testdata",
		OutputDir: "dist",
		IsDev:     false,
	}

	m := NewManager(cfg)

	if m.publicDir != "testdata" {
		t.Errorf("Expected publicDir to be 'testdata', got '%s'", m.publicDir)
	}

	if m.outputDir != "dist" {
		t.Errorf("Expected outputDir to be 'dist', got '%s'", m.outputDir)
	}

	if m.isDev {
		t.Error("Expected isDev to be false")
	}
}

func TestNewManager_Defaults(t *testing.T) {
	m := NewManager(Config{})

	if m.publicDir != "public" {
		t.Errorf("Expected default publicDir to be 'public', got '%s'", m.publicDir)
	}

	if m.outputDir != "dist" {
		t.Errorf("Expected default outputDir to be 'dist', got '%s'", m.outputDir)
	}
}

func TestManager_URL_Dev(t *testing.T) {
	m := NewManager(Config{
		IsDev: true,
	})

	url := m.URL("app.css")
	expected := "/static/app.css"

	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}
}

func TestManager_URL_Production(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.css")
	content := []byte("body { color: red; }")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	url := m.URL("test.css")

	// Should contain fingerprint
	if url == "/static/test.css" {
		t.Error("Expected URL to contain fingerprint in production mode")
	}

	// Should match pattern: /static/test.{hash}.css
	if len(url) <= len("/static/test.css") {
		t.Errorf("URL doesn't appear to have fingerprint: %s", url)
	}
}

func TestManager_URL_Concurrent(t *testing.T) {
	// Test thread-safety
	tmpDir := t.TempDir()

	// Create test files
	for i := 1; i <= 3; i++ {
		testFile := filepath.Join(tmpDir, "test"+string(rune(i+48))+".css")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.URL("test1.css")
			_ = m.URL("test2.css")
			_ = m.URL("test3.css")
		}()
	}

	wg.Wait()
}

func TestManager_IsDev(t *testing.T) {
	tests := []struct {
		name     string
		isDev    bool
		expected bool
	}{
		{"development mode", true, true},
		{"production mode", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(Config{IsDev: tt.isDev})
			if m.IsDev() != tt.expected {
				t.Errorf("Expected IsDev() to return %v, got %v", tt.expected, m.IsDev())
			}
		})
	}
}

func TestManager_PublicDir(t *testing.T) {
	m := NewManager(Config{PublicDir: "assets"})
	if m.PublicDir() != "assets" {
		t.Errorf("Expected PublicDir() to return 'assets', got '%s'", m.PublicDir())
	}
}

func TestManager_SaveManifest(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "manifest.json")

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	// Add some fingerprints
	m.mu.Lock()
	m.fingerprints["app.css"] = "app.abc12345.css"
	m.fingerprints["app.js"] = "app.def67890.js"
	m.mu.Unlock()

	// Save manifest
	if err := m.SaveManifest(manifestPath); err != nil {
		t.Fatalf("Failed to save manifest: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Error("Manifest file was not created")
	}

	// Load manifest in a new manager
	m2 := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
		Manifest:  manifestPath,
	})

	// Check that fingerprints were loaded
	url := m2.URL("app.css")
	expected := "/static/app.abc12345.css"
	if url != expected {
		t.Errorf("Expected URL from manifest '%s', got '%s'", expected, url)
	}
}

