package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscover(t *testing.T) {
	r := NewRegistry()

	// Create a temporary directory
	tmpDir := t.TempDir()

	// Test discovering from empty directory
	err := r.Discover(tmpDir)
	if err != nil {
		t.Errorf("Discover() error = %v", err)
	}

	if r.Count() != 0 {
		t.Errorf("expected 0 plugins, got %d", r.Count())
	}
}

func TestDiscoverNonexistent(t *testing.T) {
	r := NewRegistry()

	err := r.Discover("/nonexistent/directory")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestDiscoverSafe(t *testing.T) {
	r := NewRegistry()

	// Create a temporary directory
	tmpDir := t.TempDir()

	// Test discovering from empty directory
	errs := r.DiscoverSafe(tmpDir)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d", len(errs))
	}

	if r.Count() != 0 {
		t.Errorf("expected 0 plugins, got %d", r.Count())
	}
}

func TestDiscoverSafeNonexistent(t *testing.T) {
	r := NewRegistry()

	errs := r.DiscoverSafe("/nonexistent/directory")
	if len(errs) == 0 {
		t.Error("expected errors for nonexistent directory")
	}
}

func TestDiscoverIgnoresNonSoFiles(t *testing.T) {
	r := NewRegistry()

	// Create a temporary directory with non-.so files
	tmpDir := t.TempDir()

	// Create some test files
	_ = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)

	err := r.Discover(tmpDir)
	if err != nil {
		t.Errorf("Discover() error = %v", err)
	}

	if r.Count() != 0 {
		t.Errorf("expected 0 plugins (non-.so files should be ignored), got %d", r.Count())
	}
}

func TestDiscoverIgnoresDirectories(t *testing.T) {
	r := NewRegistry()

	// Create a temporary directory with subdirectories
	tmpDir := t.TempDir()
	_ = os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	err := r.Discover(tmpDir)
	if err != nil {
		t.Errorf("Discover() error = %v", err)
	}

	if r.Count() != 0 {
		t.Errorf("expected 0 plugins, got %d", r.Count())
	}
}

