package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFingerprint(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create a test file with known content
	testFile := filepath.Join(tmpDir, "test.css")
	content := []byte("body { color: red; }")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	fp := m.fingerprint("test.css")

	// Check format: should be test.{hash}.css
	if !m.isFingerprinted(fp) {
		t.Errorf("Fingerprint doesn't match expected format: %s", fp)
	}

	// Verify extension is preserved
	if filepath.Ext(fp) != ".css" {
		t.Errorf("Expected extension .css, got %s", filepath.Ext(fp))
	}
}

func TestFingerprint_Consistency(t *testing.T) {
	// Same content should produce same fingerprint
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.css")
	content := []byte("body { color: blue; }")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	fp1 := m.fingerprint("test.css")
	fp2 := m.fingerprint("test.css")

	if fp1 != fp2 {
		t.Errorf("Fingerprints should be consistent: %s != %s", fp1, fp2)
	}
}

func TestFingerprint_DifferentContent(t *testing.T) {
	// Different content should produce different fingerprints
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "test1.css")
	file2 := filepath.Join(tmpDir, "test2.css")

	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     false,
	})

	fp1 := m.fingerprint("test1.css")
	fp2 := m.fingerprint("test2.css")

	// Extract just the hash part for comparison
	if fp1 == fp2 {
		t.Error("Different content should produce different fingerprints")
	}
}

func TestStripFingerprint(t *testing.T) {
	m := NewManager(Config{})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "fingerprinted css",
			input:    "app.abc12345.css",
			expected: "app.css",
		},
		{
			name:     "fingerprinted js",
			input:    "bundle.def67890.js",
			expected: "bundle.js",
		},
		{
			name:     "nested path",
			input:    "css/main.abc12345.css",
			expected: "css/main.css",
		},
		{
			name:     "non-fingerprinted",
			input:    "app.css",
			expected: "app.css",
		},
		{
			name:     "multiple dots",
			input:    "app.min.abc12345.js",
			expected: "app.min.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.stripFingerprint(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsFingerprinted(t *testing.T) {
	m := NewManager(Config{})

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "fingerprinted",
			path:     "app.abc12345.css",
			expected: true,
		},
		{
			name:     "not fingerprinted",
			path:     "app.css",
			expected: false,
		},
		{
			name:     "wrong hash length",
			path:     "app.abc.css",
			expected: false,
		},
		{
			name:     "invalid characters in hash",
			path:     "app.xxxxxxxx.css",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.isFingerprinted(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for path: %s", tt.expected, result, tt.path)
			}
		})
	}
}

func TestIsValidPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "valid relative path",
			path:     "css/app.css",
			expected: true,
		},
		{
			name:     "valid simple path",
			path:     "app.css",
			expected: true,
		},
		{
			name:     "path traversal with ..",
			path:     "../etc/passwd",
			expected: false,
		},
		{
			name:     "path traversal in middle",
			path:     "css/../../etc/passwd",
			expected: false,
		},
		{
			name:     "absolute path",
			path:     "/etc/passwd",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for path: %s", tt.expected, result, tt.path)
			}
		})
	}
}

func TestFingerprintAll(t *testing.T) {
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

	if err := m.FingerprintAll(); err != nil {
		t.Fatalf("FingerprintAll failed: %v", err)
	}

	// Verify all files have fingerprints
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, file := range testFiles {
		if _, ok := m.fingerprints[file]; !ok {
			t.Errorf("Missing fingerprint for file: %s", file)
		}
	}
}

