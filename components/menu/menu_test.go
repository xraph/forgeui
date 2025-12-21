package menu_test

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"github.com/xraph/forgeui/components/menu"
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

func TestNavbar(t *testing.T) {
	tests := []struct {
		name string
		opts []menu.NavbarOption
		want []string
	}{
		{
			name: "renders default navbar",
			opts: []menu.NavbarOption{},
			want: []string{
				`x-data`,
				`mobileMenuOpen`,
				`z-40`,
			},
		},
		{
			name: "renders fixed navbar",
			opts: []menu.NavbarOption{menu.WithFixed()},
			want: []string{
				`fixed`,
				`top-0`,
			},
		},
		{
			name: "renders sticky navbar",
			opts: []menu.NavbarOption{menu.WithSticky()},
			want: []string{
				`sticky`,
				`top-0`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := menu.NavbarWithOptions(tt.opts,
				menu.NavbarBrand(g.Text("Brand")),
			)
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

func TestNavbarBrand(t *testing.T) {
	var buf bytes.Buffer
	node := menu.NavbarBrand(g.Text("My App"))
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "My App") {
		t.Error("NavbarBrand missing content")
	}
}

func TestNavbarMenu(t *testing.T) {
	var buf bytes.Buffer
	// NavbarMenu returns a g.Group, so we need to wrap it in a container to render
	node := html.Div(
		menu.NavbarMenu(
			menu.Item("/", g.Text("Home")),
		),
	)
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	htmlStr := buf.String()
	if !strings.Contains(htmlStr, "Home") {
		t.Error("NavbarMenu missing items")
	}
	if !strings.Contains(htmlStr, "mobileMenuOpen") {
		t.Error("NavbarMenu missing mobile menu state")
	}
}

func TestSidebar(t *testing.T) {
	tests := []struct {
		name string
		opts []menu.SidebarOption
		want []string
	}{
		{
			name: "renders default sidebar",
			opts: []menu.SidebarOption{},
			want: []string{
				`x-data`,
				`collapsed`,
				`mobileOpen`,
				`z-30`,
			},
		},
		{
			name: "renders collapsed sidebar",
			opts: []menu.SidebarOption{menu.WithDefaultCollapsed(true)},
			want: []string{
				`collapsed`,
				`true`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			// Sidebar returns a g.Group, so we need to wrap it in a container to render
			node := html.Div(
				menu.SidebarWithOptions(tt.opts,
					menu.SidebarHeader(g.Text("App")),
				),
			)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			htmlStr := buf.String()

			for _, want := range tt.want {
				if !strings.Contains(htmlStr, want) {
					t.Errorf("output missing %q\ngot: %s", want, htmlStr)
				}
			}
		})
	}
}

func TestSidebarComponents(t *testing.T) {
	tests := []struct {
		name string
		node g.Node
		want []string
	}{
		{
			name: "sidebar header",
			node: menu.SidebarHeader(g.Text("My App")),
			want: []string{
				`My App`,
				`border-b`,
			},
		},
		{
			name: "sidebar content",
			node: menu.SidebarContent(g.Text("Content")),
			want: []string{
				`Content`,
				`flex-1`,
			},
		},
		{
			name: "sidebar footer",
			node: menu.SidebarFooter(g.Text("Footer")),
			want: []string{
				`Footer`,
				`border-t`,
			},
		},
		{
			name: "sidebar toggle",
			node: menu.SidebarToggle(),
			want: []string{
				`collapsed`,
				`x-show`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tt.node.Render(&buf); err != nil {
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

