package alpine

import (
	"bytes"
	"strings"
	"testing"

	"github.com/maragudk/gomponents/html"
)

func TestScripts(t *testing.T) {
	tests := []struct {
		name    string
		plugins []Plugin
		want    []string
	}{
		{
			name:    "no plugins",
			plugins: nil,
			want:    []string{AlpineCDN},
		},
		{
			name:    "single plugin",
			plugins: []Plugin{PluginFocus},
			want:    []string{"focus", AlpineCDN},
		},
		{
			name:    "multiple plugins",
			plugins: []Plugin{PluginFocus, PluginCollapse},
			want:    []string{"focus", "collapse", AlpineCDN},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			// Wrap in HTML element since g.Group can't be rendered directly
			node := html.Div(Scripts(tt.plugins...))
			if err := node.Render(&buf); err != nil {
				t.Fatalf("Render error: %v", err)
			}
			got := buf.String()

			// Check all expected URLs are present
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("Scripts() = %v, want to contain %v", got, w)
				}
			}

			// Check defer attribute is present
			if !strings.Contains(got, `defer=""`) {
				t.Errorf("Scripts() missing defer attribute")
			}

			// Check that Alpine CDN appears last
			alpineIndex := strings.Index(got, AlpineCDN)
			if alpineIndex == -1 {
				t.Errorf("Scripts() missing Alpine CDN")
			}

			// If plugins exist, they should appear before Alpine
			if len(tt.plugins) > 0 {
				for _, p := range tt.plugins {
					pluginIndex := strings.Index(got, string(p))
					if pluginIndex > alpineIndex {
						t.Errorf("Plugin %s appears after Alpine CDN (wrong order)", p)
					}
				}
			}
		})
	}
}

func TestScriptsWithVersion(t *testing.T) {
	var buf bytes.Buffer
	html.Div(ScriptsWithVersion("3.13.3", PluginFocus)).Render(&buf)
	got := buf.String()

	// Check version is in URL
	if !strings.Contains(got, "3.13.3") {
		t.Errorf("ScriptsWithVersion() = %v, want to contain version 3.13.3", got)
	}

	// Check plugin and Alpine are present
	if !strings.Contains(got, "focus") {
		t.Errorf("ScriptsWithVersion() missing Focus plugin")
	}
	if !strings.Contains(got, "alpinejs") {
		t.Errorf("ScriptsWithVersion() missing Alpine")
	}
}

func TestCloakCSS(t *testing.T) {
	var buf bytes.Buffer
	CloakCSS().Render(&buf)
	got := buf.String()

	if !strings.Contains(got, "[x-cloak]") {
		t.Errorf("CloakCSS() = %v, want to contain [x-cloak] selector", got)
	}
	if !strings.Contains(got, "display: none !important") {
		t.Errorf("CloakCSS() = %v, want to contain display: none", got)
	}
}

func TestScriptsWithNonce(t *testing.T) {
	var buf bytes.Buffer
	nonce := "random-nonce-value"
	html.Div(ScriptsWithNonce(nonce, PluginFocus)).Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `nonce="`+nonce+`"`) {
		t.Errorf("ScriptsWithNonce() = %v, want to contain nonce attribute", got)
	}
}

func TestPluginURL(t *testing.T) {
	tests := []struct {
		name   string
		plugin Plugin
		want   string
	}{
		{
			name:   "Focus plugin",
			plugin: PluginFocus,
			want:   "focus",
		},
		{
			name:   "Collapse plugin",
			plugin: PluginCollapse,
			want:   "collapse",
		},
		{
			name:   "Unknown plugin",
			plugin: Plugin("unknown"),
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PluginURL(tt.plugin)
			if tt.want == "" {
				if got != "" {
					t.Errorf("PluginURL() = %v, want empty string for unknown plugin", got)
				}
			} else {
				if !strings.Contains(got, tt.want) {
					t.Errorf("PluginURL() = %v, want to contain %v", got, tt.want)
				}
			}
		})
	}
}

func TestAllPlugins(t *testing.T) {
	plugins := AllPlugins()

	if len(plugins) == 0 {
		t.Error("AllPlugins() returned empty list")
	}

	// Check that common plugins are included
	found := make(map[Plugin]bool)
	for _, p := range plugins {
		found[p] = true
	}

	expectedPlugins := []Plugin{PluginFocus, PluginCollapse, PluginMask}
	for _, expected := range expectedPlugins {
		if !found[expected] {
			t.Errorf("AllPlugins() missing %s", expected)
		}
	}
}
