package accordion_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui/components/accordion"
	g "maragu.dev/gomponents"
)

func TestAccordion(t *testing.T) {
	tests := []struct {
		name     string
		opts     []accordion.Option
		children []g.Node
		want     []string
	}{
		{
			name: "renders default accordion",
			children: []g.Node{
				accordion.Item("item1", "Title", g.Text("Content")),
			},
			want: []string{
				`x-data`,
				`openItem`,
				`type`,
				`single`,
			},
		},
		{
			name: "renders with multiple type",
			opts: []accordion.Option{accordion.WithType(accordion.TypeMultiple)},
			children: []g.Node{
				accordion.Item("item1", "Title 1", g.Text("Content 1")),
			},
			want: []string{
				`openItems`,
				`multiple`,
			},
		},
		{
			name: "renders with collapsible",
			opts: []accordion.Option{accordion.WithCollapsible()},
			children: []g.Node{
				accordion.Item("item1", "Title", g.Text("Content")),
			},
			want: []string{
				`collapsible`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := accordion.AccordionWithOptions(tt.opts, tt.children...)
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

func TestAccordionItem(t *testing.T) {
	var buf bytes.Buffer

	node := accordion.Item("test-item", "Test Title", g.Text("Test Content"))
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Test Title") {
		t.Error("AccordionItem missing title")
	}

	if !strings.Contains(html, "Test Content") {
		t.Error("AccordionItem missing content")
	}

	if !strings.Contains(html, "x-collapse") {
		t.Error("AccordionItem missing collapse directive")
	}

	if !strings.Contains(html, "test-item") {
		t.Error("AccordionItem missing id reference")
	}
}
