package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "manifest.json")

	// Create a test manifest
	manifestContent := `{
  "app.css": "app.abc12345.css",
  "app.js": "app.def67890.js"
}`

	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}

	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("Failed to load manifest: %v", err)
	}

	if len(manifest) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(manifest))
	}

	if manifest["app.css"] != "app.abc12345.css" {
		t.Errorf("Unexpected value for app.css: %s", manifest["app.css"])
	}
}

func TestManifest_Save(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "subdir", "manifest.json")

	manifest := Manifest{
		"app.css": "app.abc12345.css",
		"app.js":  "app.def67890.js",
	}

	if err := manifest.Save(manifestPath); err != nil {
		t.Fatalf("Failed to save manifest: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Error("Manifest file was not created")
	}

	// Load and verify content
	loaded, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("Failed to load saved manifest: %v", err)
	}

	if len(loaded) != len(manifest) {
		t.Errorf("Expected %d entries, got %d", len(manifest), len(loaded))
	}

	for k, v := range manifest {
		if loaded[k] != v {
			t.Errorf("Mismatch for key %s: expected %s, got %s", k, v, loaded[k])
		}
	}
}

func TestManifest_Get(t *testing.T) {
	manifest := Manifest{
		"app.css": "app.abc12345.css",
	}

	fp, ok := manifest.Get("app.css")
	if !ok {
		t.Error("Expected to find app.css in manifest")
	}

	if fp != "app.abc12345.css" {
		t.Errorf("Expected 'app.abc12345.css', got '%s'", fp)
	}

	_, ok = manifest.Get("nonexistent.css")
	if ok {
		t.Error("Expected not to find nonexistent.css in manifest")
	}
}

func TestManifest_Set(t *testing.T) {
	manifest := make(Manifest)

	manifest.Set("app.css", "app.abc12345.css")

	if manifest["app.css"] != "app.abc12345.css" {
		t.Error("Failed to set manifest entry")
	}
}

func TestGenerateManifest(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"app.css",
		"main.js",
		"css/style.css",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)

		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	manifest, err := GenerateManifest(m)
	if err != nil {
		t.Fatalf("Failed to generate manifest: %v", err)
	}

	if len(manifest) != len(testFiles) {
		t.Errorf("Expected %d entries, got %d", len(testFiles), len(manifest))
	}

	for _, file := range testFiles {
		if _, ok := manifest[file]; !ok {
			t.Errorf("Missing manifest entry for %s", file)
		}
	}
}
