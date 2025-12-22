package theme_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui/theme"
)

func TestDarkModeScript(t *testing.T) {
	var buf bytes.Buffer
	node := theme.DarkModeScript()
	
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for script tag
	if !strings.Contains(html, "<script") {
		t.Error("DarkModeScript should render a script element")
	}

	// Check for data attribute
	if !strings.Contains(html, `data-theme-script`) {
		t.Error("DarkModeScript should have data-theme-script attribute")
	}

	// Check for localStorage
	if !strings.Contains(html, "localStorage.getItem") {
		t.Error("Script should read from localStorage")
	}

	// Check for dark class toggle
	if !strings.Contains(html, "classList.toggle") {
		t.Error("Script should toggle dark class")
	}

	// Check for system preference
	if !strings.Contains(html, "prefers-color-scheme") {
		t.Error("Script should check system preference")
	}
}

func TestDarkModeScriptWithDefault(t *testing.T) {
	var buf bytes.Buffer
	node := theme.DarkModeScriptWithDefault("dark")
	
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "dark") {
		t.Error("Script should contain default theme")
	}
}

func TestThemeScript(t *testing.T) {
	var buf bytes.Buffer
	node := theme.ThemeScript()
	
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for theme manager attribute
	if !strings.Contains(html, `data-theme-manager`) {
		t.Error("ThemeScript should have data-theme-manager attribute")
	}

	// Check for exposed functions
	expectedFunctions := []string{
		"getTheme",
		"setTheme",
		"toggleTheme",
		"initTheme",
	}

	for _, fn := range expectedFunctions {
		if !strings.Contains(html, fn) {
			t.Errorf("Script should contain %s function", fn)
		}
	}

	// Check for window.forgeui namespace
	if !strings.Contains(html, "window.forgeui") {
		t.Error("Script should expose window.forgeui namespace")
	}

	// Check for storage event listener
	if !strings.Contains(html, "addEventListener('storage'") {
		t.Error("Script should listen for storage events")
	}
}

func TestThemeStyleCloak(t *testing.T) {
	var buf bytes.Buffer
	node := theme.ThemeStyleCloak()
	
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for style tag
	if !strings.Contains(html, "<style") {
		t.Error("ThemeStyleCloak should render a style element")
	}

	// Check for x-cloak styles
	if !strings.Contains(html, "[x-cloak]") {
		t.Error("Should contain x-cloak styles")
	}

	// Check for theme-loading styles
	if !strings.Contains(html, "[data-theme-loading]") {
		t.Error("Should contain data-theme-loading styles")
	}
}

