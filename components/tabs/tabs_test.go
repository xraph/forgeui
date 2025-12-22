package tabs_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui/components/tabs"
	g "maragu.dev/gomponents"
)

func TestTabs(t *testing.T) {
	tests := []struct {
		name     string
		opts     []tabs.Option
		children []g.Node
		want     []string
		wantNot  []string
	}{
		{
			name: "renders default tabs",
			children: []g.Node{
				tabs.TabList(
					tabs.Tab("tab1", g.Text("Tab 1")),
				),
			},
			want: []string{
				`x-data`,
				`activeTab`,
				`role="tablist"`,
			},
		},
		{
			name: "renders with default tab",
			opts: []tabs.Option{tabs.WithDefaultTab("custom")},
			children: []g.Node{
				tabs.TabList(
					tabs.Tab("custom", g.Text("Custom Tab")),
				),
			},
			want: []string{
				`activeTab`,
				`custom`,
			},
		},
		{
			name: "renders with custom class",
			opts: []tabs.Option{tabs.WithClass("custom-class")},
			children: []g.Node{
				tabs.TabList(),
			},
			want: []string{
				`custom-class`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.TabsWithOptions(tt.opts, tt.children...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, want := range tt.want {
				if !strings.Contains(html, want) {
					t.Errorf("output missing %q\ngot: %s", want, html)
				}
			}

			for _, wantNot := range tt.wantNot {
				if strings.Contains(html, wantNot) {
					t.Errorf("output should not contain %q\ngot: %s", wantNot, html)
				}
			}
		})
	}
}

func TestTabList(t *testing.T) {
	var buf bytes.Buffer

	node := tabs.TabList(
		tabs.Tab("tab1", g.Text("Tab 1")),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `role="tablist"`) {
		t.Error("TabList missing role attribute")
	}

	if !strings.Contains(html, "bg-muted") {
		t.Error("TabList missing expected classes")
	}
}

func TestTab(t *testing.T) {
	var buf bytes.Buffer

	node := tabs.Tab("test-tab", g.Text("Test"))
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `role="tab"`) {
		t.Error("Tab missing role attribute")
	}

	if !strings.Contains(html, `test-tab`) {
		t.Error("Tab missing id reference")
	}

	if !strings.Contains(html, "Test") {
		t.Error("Tab missing label")
	}
}

func TestTabPanel(t *testing.T) {
	var buf bytes.Buffer

	node := tabs.TabPanel("panel1", g.Text("Panel content"))
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `role="tabpanel"`) {
		t.Error("TabPanel missing role attribute")
	}

	if !strings.Contains(html, "Panel content") {
		t.Error("TabPanel missing content")
	}

	if !strings.Contains(html, "x-show") {
		t.Error("TabPanel missing x-show directive")
	}
}
