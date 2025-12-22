package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScript(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.js")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.Script("test.js")
	html := renderNode(node)

	if !strings.Contains(html, `<script`) {
		t.Error("Expected script tag")
	}

	if !strings.Contains(html, `src="/static/test.js"`) {
		t.Error("Expected correct src attribute")
	}
}

func TestScript_WithDefer(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.Script("test.js", WithDefer())
	html := renderNode(node)

	if !strings.Contains(html, `defer`) {
		t.Error("Expected defer attribute")
	}
}

func TestScript_WithAsync(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.Script("test.js", WithAsync())
	html := renderNode(node)

	if !strings.Contains(html, `async`) {
		t.Error("Expected async attribute")
	}
}

func TestScript_WithModule(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.Script("test.js", WithModule())
	html := renderNode(node)

	if !strings.Contains(html, `type="module"`) {
		t.Error("Expected type=module attribute")
	}
}

func TestScript_WithScriptIntegrity(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.Script("test.js", WithScriptIntegrity("sha256-abc123"))
	html := renderNode(node)

	if !strings.Contains(html, `integrity="sha256-abc123"`) {
		t.Error("Expected integrity attribute")
	}
}

func TestScript_WithScriptCrossOrigin(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.Script("test.js", WithScriptCrossOrigin("anonymous"))
	html := renderNode(node)

	if !strings.Contains(html, `crossorigin="anonymous"`) {
		t.Error("Expected crossorigin attribute")
	}
}

func TestScript_WithNoModule(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.Script("test.js", WithNoModule())
	html := renderNode(node)

	if !strings.Contains(html, `nomodule`) {
		t.Error("Expected nomodule attribute")
	}
}

func TestPreloadScript(t *testing.T) {
	tmpDir := t.TempDir()

	m := NewManager(Config{
		PublicDir: tmpDir,
		IsDev:     true,
	})

	node := m.PreloadScript("test.js")
	html := renderNode(node)

	if !strings.Contains(html, `rel="preload"`) {
		t.Error("Expected preload rel attribute")
	}

	if !strings.Contains(html, `as="script"`) {
		t.Error("Expected as=script attribute")
	}
}

func TestInlineScript(t *testing.T) {
	content := "console.log('test');"
	node := InlineScript(content)
	html := renderNode(node)

	if !strings.Contains(html, "<script>") {
		t.Error("Expected script tag")
	}

	if !strings.Contains(html, content) {
		t.Error("Expected JavaScript content in script tag")
	}
}
