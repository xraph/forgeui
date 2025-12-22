package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStyleSheet(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.css")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.StyleSheet("test.css")
	html := renderNode(node)

	if !strings.Contains(html, `rel="stylesheet"`) {
		t.Error("Expected stylesheet rel attribute")
	}

	if !strings.Contains(html, `href="/static/test.css"`) {
		t.Error("Expected correct href attribute")
	}
}

func TestStyleSheet_WithMedia(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.StyleSheet("test.css", WithMedia("print"))
	html := renderNode(node)

	if !strings.Contains(html, `media="print"`) {
		t.Error("Expected media attribute")
	}
}

func TestStyleSheet_WithIntegrity(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.StyleSheet("test.css", WithIntegrity("sha256-abc123"))
	html := renderNode(node)

	if !strings.Contains(html, `integrity="sha256-abc123"`) {
		t.Error("Expected integrity attribute")
	}
}

func TestStyleSheet_WithCrossOrigin(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.StyleSheet("test.css", WithCrossOrigin("anonymous"))
	html := renderNode(node)

	if !strings.Contains(html, `crossorigin="anonymous"`) {
		t.Error("Expected crossorigin attribute")
	}
}

func TestPreloadStyleSheet(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.PreloadStyleSheet("test.css")
	html := renderNode(node)

	if !strings.Contains(html, `rel="preload"`) {
		t.Error("Expected preload rel attribute")
	}

	if !strings.Contains(html, `as="style"`) {
		t.Error("Expected as=style attribute")
	}
}

func TestInlineCSS(t *testing.T) {
	content := "body { color: red; }"
	node := InlineCSS(content)
	html := renderNode(node)

	if !strings.Contains(html, "<style>") {
		t.Error("Expected style tag")
	}

	if !strings.Contains(html, content) {
		t.Error("Expected CSS content in style tag")
	}
}
