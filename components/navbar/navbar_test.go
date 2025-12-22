package navbar_test

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"github.com/xraph/forgeui/components/menu"
	navbar "github.com/xraph/forgeui/components/navbar"
)

func TestNavbar(t *testing.T) {
	tests := []struct {
		name string
		opts []navbar.NavbarOption
		want []string
	}{
		{
			name: "renders default navbar",
			opts: []navbar.NavbarOption{},
			want: []string{
				`x-data`,
				`mobileMenuOpen`,
				`z-40`,
			},
		},
		{
			name: "renders fixed navbar",
			opts: []navbar.NavbarOption{navbar.WithFixed()},
			want: []string{
				`fixed`,
				`top-0`,
			},
		},
		{
			name: "renders sticky navbar",
			opts: []navbar.NavbarOption{navbar.WithSticky()},
			want: []string{
				`sticky`,
				`top-0`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			node := navbar.NavbarWithOptions(tt.opts,
				navbar.NavbarBrand(g.Text("Brand")),
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
	node := navbar.NavbarBrand(g.Text("My App"))
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
		navbar.NavbarMenu(
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
