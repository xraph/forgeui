package breadcrumb_test

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/xraph/forgeui/components/breadcrumb"
)

func TestBreadcrumb(t *testing.T) {
	tests := []struct {
		name     string
		opts     []breadcrumb.Option
		children []g.Node
		want     []string
	}{
		{
			name: "renders default breadcrumb",
			children: []g.Node{
				breadcrumb.Item("/", g.Text("Home")),
				breadcrumb.Page(g.Text("Current")),
			},
			want: []string{
				`aria-label="Breadcrumb"`,
				`Home`,
				`Current`,
			},
		},
		{
			name: "renders with custom class",
			opts: []breadcrumb.Option{breadcrumb.WithClass("custom-class")},
			children: []g.Node{
				breadcrumb.Item("/", g.Text("Home")),
			},
			want: []string{
				`custom-class`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := breadcrumb.BreadcrumbWithOptions(tt.opts, tt.children...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, want := range tt.want {
				if !strings.Contains(html, want) {
					t.Errorf("output missing %q\ngot: %s", want, html)
				}
			}
		})
	}
}

func TestBreadcrumbItem(t *testing.T) {
	var buf bytes.Buffer
	node := breadcrumb.Item("/docs", g.Text("Documentation"))
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `href="/docs"`) {
		t.Error("BreadcrumbItem missing href")
	}
	if !strings.Contains(html, "Documentation") {
		t.Error("BreadcrumbItem missing label")
	}
}

func TestBreadcrumbPage(t *testing.T) {
	var buf bytes.Buffer
	node := breadcrumb.Page(g.Text("Current Page"))
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `aria-current="page"`) {
		t.Error("BreadcrumbPage missing aria-current")
	}
	if !strings.Contains(html, "Current Page") {
		t.Error("BreadcrumbPage missing label")
	}
}

func TestBreadcrumbSeparator(t *testing.T) {
	var buf bytes.Buffer
	node := breadcrumb.Separator()
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `role="presentation"`) {
		t.Error("Separator missing role attribute")
	}
	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("Separator missing aria-hidden")
	}
}

