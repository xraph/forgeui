package assets

import (
	"strings"
	"sync"
	"testing"

	g "github.com/maragudk/gomponents"
)

func TestNewScriptManager(t *testing.T) {
	sm := NewScriptManager()
	if sm == nil {
		t.Fatal("NewScriptManager returned nil")
	}
	if sm.Count() != 0 {
		t.Errorf("New manager should have 0 scripts, got %d", sm.Count())
	}
}

func TestScriptManager_Add(t *testing.T) {
	sm := NewScriptManager()

	// Add a script
	sm.Add(ScriptEntry{
		Path:     "/js/app.js",
		Priority: 50,
		Position: "body",
	})

	if sm.Count() != 1 {
		t.Errorf("Expected 1 script, got %d", sm.Count())
	}
}

func TestScriptManager_AddDuplicates(t *testing.T) {
	sm := NewScriptManager()

	// Add same script twice
	sm.Add(ScriptEntry{
		Path:     "/js/app.js",
		Priority: 50,
	})
	sm.Add(ScriptEntry{
		Path:     "/js/app.js",
		Priority: 100, // Different priority
	})

	if sm.Count() != 1 {
		t.Errorf("Duplicate scripts should be ignored, got %d scripts", sm.Count())
	}
}

func TestScriptManager_AddInlineDuplicates(t *testing.T) {
	sm := NewScriptManager()

	// Inline scripts should not be deduplicated (they might have different content)
	sm.AddInline("console.log('a')", 50, "body")
	sm.AddInline("console.log('b')", 60, "body")

	if sm.Count() != 2 {
		t.Errorf("Expected 2 inline scripts, got %d", sm.Count())
	}
}

func TestScriptManager_AddWithPriority(t *testing.T) {
	sm := NewScriptManager()

	attrs := map[string]string{"defer": ""}
	sm.AddWithPriority("/js/app.js", 75, attrs)

	if sm.Count() != 1 {
		t.Errorf("Expected 1 script, got %d", sm.Count())
	}
}

func TestScriptManager_AddInline(t *testing.T) {
	sm := NewScriptManager()

	sm.AddInline("console.log('test')", 50, "head")

	if sm.Count() != 1 {
		t.Errorf("Expected 1 script, got %d", sm.Count())
	}
}

func TestScriptManager_PrioritySorting(t *testing.T) {
	sm := NewScriptManager()

	// Add scripts in random order
	sm.Add(ScriptEntry{Path: "/js/analytics.js", Priority: 90, Position: "body"})
	sm.Add(ScriptEntry{Path: "/js/framework.js", Priority: 10, Position: "body"})
	sm.Add(ScriptEntry{Path: "/js/app.js", Priority: 50, Position: "body"})
	sm.Add(ScriptEntry{Path: "/js/library.js", Priority: 30, Position: "body"})

	nodes := sm.Render("body")
	if len(nodes) != 4 {
		t.Fatalf("Expected 4 script nodes, got %d", len(nodes))
	}

	// Render and check order
	html := renderNodes(nodes)

	// Check that scripts appear in priority order
	frameworkPos := strings.Index(html, "/js/framework.js")
	libraryPos := strings.Index(html, "/js/library.js")
	appPos := strings.Index(html, "/js/app.js")
	analyticsPos := strings.Index(html, "/js/analytics.js")

	if frameworkPos == -1 || libraryPos == -1 || appPos == -1 || analyticsPos == -1 {
		t.Fatal("Not all scripts were rendered")
	}

	if !(frameworkPos < libraryPos && libraryPos < appPos && appPos < analyticsPos) {
		t.Errorf("Scripts not in priority order: framework=%d, library=%d, app=%d, analytics=%d",
			frameworkPos, libraryPos, appPos, analyticsPos)
	}
}

func TestScriptManager_PositionFiltering(t *testing.T) {
	sm := NewScriptManager()

	// Add scripts to different positions
	sm.Add(ScriptEntry{Path: "/js/critical.js", Priority: 10, Position: "head"})
	sm.Add(ScriptEntry{Path: "/js/app.js", Priority: 50, Position: "body"})
	sm.Add(ScriptEntry{Path: "/js/defer.js", Priority: 60, Position: "body"})

	headNodes := sm.Render("head")
	if len(headNodes) != 1 {
		t.Errorf("Expected 1 head script, got %d", len(headNodes))
	}

	bodyNodes := sm.Render("body")
	if len(bodyNodes) != 2 {
		t.Errorf("Expected 2 body scripts, got %d", len(bodyNodes))
	}
}

func TestScriptManager_InlineScriptRendering(t *testing.T) {
	sm := NewScriptManager()

	content := "console.log('inline test');"
	sm.AddInline(content, 50, "body")

	nodes := sm.Render("body")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 script node, got %d", len(nodes))
	}

	html := renderNodes(nodes)

	if !strings.Contains(html, content) {
		t.Errorf("Inline script content not found in output: %s", html)
	}

	// Should not have src attribute
	if strings.Contains(html, "src=") {
		t.Error("Inline script should not have src attribute")
	}
}

func TestScriptManager_ExternalScriptRendering(t *testing.T) {
	sm := NewScriptManager()

	sm.Add(ScriptEntry{
		Path:     "/js/app.js",
		Priority: 50,
		Position: "body",
	})

	nodes := sm.Render("body")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 script node, got %d", len(nodes))
	}

	html := renderNodes(nodes)

	if !strings.Contains(html, `src="/js/app.js"`) {
		t.Errorf("Script src not found in output: %s", html)
	}
}

func TestScriptManager_ScriptAttributes(t *testing.T) {
	sm := NewScriptManager()

	sm.Add(ScriptEntry{
		Path:     "/js/module.js",
		Priority: 50,
		Position: "body",
		Attrs: map[string]string{
			"type":  "module",
			"defer": "",
		},
	})

	nodes := sm.Render("body")
	html := renderNodes(nodes)

	if !strings.Contains(html, `type="module"`) {
		t.Error("Script should have type=module attribute")
	}

	if !strings.Contains(html, "defer") {
		t.Error("Script should have defer attribute")
	}
}

func TestScriptManager_InlineScriptAttributes(t *testing.T) {
	sm := NewScriptManager()

	sm.Add(ScriptEntry{
		Inline:   true,
		Content:  "console.log('test');",
		Priority: 50,
		Position: "body",
		Attrs: map[string]string{
			"type": "module",
		},
	})

	nodes := sm.Render("body")
	html := renderNodes(nodes)

	if !strings.Contains(html, `type="module"`) {
		t.Error("Inline script should have type=module attribute")
	}

	if !strings.Contains(html, "console.log('test')") {
		t.Error("Inline script content missing")
	}
}

func TestScriptManager_DefaultPosition(t *testing.T) {
	sm := NewScriptManager()

	// Add script without specifying position
	sm.Add(ScriptEntry{
		Path:     "/js/app.js",
		Priority: 50,
	})

	// Should default to "body"
	bodyNodes := sm.Render("body")
	if len(bodyNodes) != 1 {
		t.Errorf("Script should default to body position, got %d body scripts", len(bodyNodes))
	}

	headNodes := sm.Render("head")
	if len(headNodes) != 0 {
		t.Errorf("Should have no head scripts, got %d", len(headNodes))
	}
}

func TestScriptManager_Clear(t *testing.T) {
	sm := NewScriptManager()

	sm.Add(ScriptEntry{Path: "/js/app.js", Priority: 50})
	sm.Add(ScriptEntry{Path: "/js/lib.js", Priority: 30})

	if sm.Count() != 2 {
		t.Errorf("Expected 2 scripts before clear, got %d", sm.Count())
	}

	sm.Clear()

	if sm.Count() != 0 {
		t.Errorf("Expected 0 scripts after clear, got %d", sm.Count())
	}
}

func TestScriptManager_CountByPosition(t *testing.T) {
	sm := NewScriptManager()

	sm.Add(ScriptEntry{Path: "/js/critical.js", Priority: 10, Position: "head"})
	sm.Add(ScriptEntry{Path: "/js/app.js", Priority: 50, Position: "body"})
	sm.Add(ScriptEntry{Path: "/js/defer.js", Priority: 60, Position: "body"})

	headCount := sm.CountByPosition("head")
	if headCount != 1 {
		t.Errorf("Expected 1 head script, got %d", headCount)
	}

	bodyCount := sm.CountByPosition("body")
	if bodyCount != 2 {
		t.Errorf("Expected 2 body scripts, got %d", bodyCount)
	}

	footerCount := sm.CountByPosition("footer")
	if footerCount != 0 {
		t.Errorf("Expected 0 footer scripts, got %d", footerCount)
	}
}

func TestScriptManager_ConcurrentAccess(t *testing.T) {
	sm := NewScriptManager()

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			sm.Add(ScriptEntry{
				Path:     "/js/script" + string(rune(n)) + ".js",
				Priority: n,
				Position: "body",
			})
		}(i)
	}

	wg.Wait()

	// Should have all scripts
	if sm.Count() < 1 {
		t.Error("Concurrent adds should result in scripts being added")
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sm.Render("body")
			_ = sm.Count()
			_ = sm.CountByPosition("body")
		}()
	}

	wg.Wait()
}

func TestScriptManager_EmptyRender(t *testing.T) {
	sm := NewScriptManager()

	nodes := sm.Render("body")
	if len(nodes) != 0 {
		t.Errorf("Empty manager should render 0 nodes, got %d", len(nodes))
	}
}

// Helper function to render nodes to HTML string
func renderNodes(nodes []g.Node) string {
	var sb strings.Builder
	for _, node := range nodes {
		_ = node.Render(&sb)
	}
	return sb.String()
}

