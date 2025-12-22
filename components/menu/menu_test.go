package menu_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui/components/menu"
	g "maragu.dev/gomponents"
)

func TestMenu(t *testing.T) {
	tests := []struct {
		name     string
		opts     []menu.MenuOption
		children []g.Node
		want     []string
	}{
		{
			name: "renders default menu",
			children: []g.Node{
				menu.Item("/", g.Text("Home")),
			},
			want: []string{
				`role="menu"`,
				`Home`,
			},
		},
		{
			name: "renders with custom class",
			opts: []menu.MenuOption{menu.WithClass("custom-class")},
			children: []g.Node{
				menu.Item("/", g.Text("Home")),
			},
			want: []string{
				`custom-class`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := menu.MenuWithOptions(tt.opts, tt.children...)
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

func TestMenuItem(t *testing.T) {
	tests := []struct {
		name string
		href string
		opts []menu.ItemOption
		want []string
	}{
		{
			name: "renders menu item with href",
			href: "/dashboard",
			opts: []menu.ItemOption{},
			want: []string{
				`href="/dashboard"`,
				`role="menuitem"`,
			},
		},
		{
			name: "renders active menu item",
			href: "/",
			opts: []menu.ItemOption{menu.Active()},
			want: []string{
				`bg-accent`,
			},
		},
		{
			name: "renders menu item without href",
			href: "",
			opts: []menu.ItemOption{},
			want: []string{
				`<button`,
				`role="menuitem"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := menu.Item(tt.href, g.Text("Item"), tt.opts...)
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

func TestMenuSection(t *testing.T) {
	var buf bytes.Buffer

	node := menu.Section("Main",
		menu.Item("/", g.Text("Home")),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Main") {
		t.Error("MenuSection missing label")
	}

	if !strings.Contains(html, "Home") {
		t.Error("MenuSection missing items")
	}
}

func TestMenuSeparator(t *testing.T) {
	var buf bytes.Buffer

	node := menu.Separator()
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `role="separator"`) {
		t.Error("Separator missing role attribute")
	}
}
