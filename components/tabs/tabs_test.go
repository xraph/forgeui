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

func TestTabListWithScrollable(t *testing.T) {
	tests := []struct {
		name    string
		opts    []tabs.TabListOption
		want    []string
		wantNot []string
	}{
		{
			name: "default TabList has w-full",
			opts: nil,
			want: []string{
				"w-full",
				`role="tablist"`,
			},
			wantNot: []string{
				"overflow-x-auto",
				"scroll-smooth",
				`data-scrollable="true"`,
			},
		},
		{
			name: "scrollable TabList has overflow classes and data attribute",
			opts: []tabs.TabListOption{tabs.WithScrollable()},
			want: []string{
				"overflow-x-auto",
				"scroll-smooth",
				"gap-1",
				`role="tablist"`,
				`data-scrollable="true"`,
			},
			wantNot: []string{
				"w-full",
			},
		},
		{
			name: "TabList with custom class",
			opts: []tabs.TabListOption{tabs.WithTabListClass("custom-tablist")},
			want: []string{
				"custom-tablist",
				`role="tablist"`,
			},
		},
		{
			name: "TabList with custom attrs",
			opts: []tabs.TabListOption{tabs.WithTabListAttrs(g.Attr("data-testid", "my-tablist"))},
			want: []string{
				`data-testid="my-tablist"`,
				`role="tablist"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.TabListWithOptions(tt.opts,
				tabs.Tab("tab1", g.Text("Tab 1")),
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

			for _, wantNot := range tt.wantNot {
				if strings.Contains(html, wantNot) {
					t.Errorf("output should not contain %q\ngot: %s", wantNot, html)
				}
			}
		})
	}
}

func TestTabListVariants(t *testing.T) {
	tests := []struct {
		name    string
		variant tabs.TabListVariant
		want    []string
		wantNot []string
	}{
		{
			name:    "default variant has boxed style",
			variant: tabs.TabListVariantDefault,
			want: []string{
				"bg-muted",
				"rounded-md",
				"text-muted-foreground",
			},
			wantNot: []string{
				"border-b",
			},
		},
		{
			name:    "underline variant has border",
			variant: tabs.TabListVariantUnderline,
			want: []string{
				"border-b",
				"border-border",
				"gap-6",
			},
			wantNot: []string{
				"bg-muted",
				"rounded-md",
			},
		},
		{
			name:    "pills variant has padding and gap",
			variant: tabs.TabListVariantPills,
			want: []string{
				"gap-2",
				"p-1",
			},
			wantNot: []string{
				"bg-muted",
				"border-b",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.TabListWithOptions(
				[]tabs.TabListOption{tabs.WithTabListVariant(tt.variant)},
				tabs.Tab("tab1", g.Text("Tab 1")),
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

			for _, wantNot := range tt.wantNot {
				if strings.Contains(html, wantNot) {
					t.Errorf("output should not contain %q\ngot: %s", wantNot, html)
				}
			}
		})
	}
}

func TestTabVariants(t *testing.T) {
	tests := []struct {
		name    string
		variant tabs.TabVariant
		want    []string
		wantNot []string
	}{
		{
			name:    "default variant has boxed style",
			variant: tabs.TabVariantDefault,
			want: []string{
				"rounded-sm",
				"data-[state=active]:bg-background",
				"data-[state=active]:shadow-sm",
			},
			wantNot: []string{
				"border-b-2",
				"rounded-full",
			},
		},
		{
			name:    "underline variant has bottom border and hover effects",
			variant: tabs.TabVariantUnderline,
			want: []string{
				"rounded-md",
				"border-b",
				"border-transparent",
				"hover:bg-muted/50",
				"data-[state=active]:border-b",
				"data-[state=active]:border-primary",
			},
			wantNot: []string{
				"rounded-sm",
				"shadow-sm",
			},
		},
		{
			name:    "pills variant has rounded full",
			variant: tabs.TabVariantPills,
			want: []string{
				"rounded-full",
				"data-[state=active]:bg-primary",
				"data-[state=active]:text-primary-foreground",
			},
			wantNot: []string{
				"border-b-2",
				"shadow-sm",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.Tab("tab1", g.Text("Tab 1"), tabs.WithTabVariant(tt.variant))
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

func TestTabWithCustomization(t *testing.T) {
	tests := []struct {
		name    string
		opts    []tabs.TabOption
		want    []string
		wantNot []string
	}{
		{
			name: "default Tab does not have flex-1 (natural width)",
			opts: nil,
			want: []string{
				`role="tab"`,
				`<button`,
				"inline-flex",
			},
			wantNot: []string{
				"flex-1",
			},
		},
		{
			name: "Tab with custom class",
			opts: []tabs.TabOption{tabs.WithTabClass("custom-tab-class")},
			want: []string{
				"custom-tab-class",
				`role="tab"`,
			},
			wantNot: []string{
				"flex-1",
			},
		},
		{
			name: "Tab with grow option has flex-1",
			opts: []tabs.TabOption{tabs.WithGrow()},
			want: []string{
				`role="tab"`,
				"inline-flex",
				"flex-1",
			},
		},
		{
			name: "Tab with shrink option (deprecated, no-op)",
			opts: []tabs.TabOption{tabs.WithShrink()},
			want: []string{
				`role="tab"`,
				"inline-flex",
			},
			wantNot: []string{
				"flex-1",
			},
		},
		{
			name: "Tab with custom attrs",
			opts: []tabs.TabOption{tabs.WithTabAttrs(g.Attr("data-testid", "my-tab"))},
			want: []string{
				`data-testid="my-tab"`,
				`role="tab"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.Tab("test-tab", g.Text("Test"), tt.opts...)
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

func TestTabAsLink(t *testing.T) {
	tests := []struct {
		name    string
		opts    []tabs.TabOption
		want    []string
		wantNot []string
	}{
		{
			name: "Tab with href renders as link",
			opts: []tabs.TabOption{tabs.WithHref("/overview")},
			want: []string{
				`<a`,
				`href="/overview"`,
				`role="tab"`,
				`:aria-selected`,
				`@click`,
			},
			wantNot: []string{
				`<button`,
				`type="button"`,
			},
		},
		{
			name: "Tab without href renders as button",
			opts: nil,
			want: []string{
				`<button`,
				`type="button"`,
				`role="tab"`,
			},
			wantNot: []string{
				`<a`,
				`href=`,
			},
		},
		{
			name: "Tab link with custom class and grow",
			opts: []tabs.TabOption{
				tabs.WithHref("/settings"),
				tabs.WithTabClass("link-tab"),
				tabs.WithGrow(),
			},
			want: []string{
				`<a`,
				`href="/settings"`,
				"link-tab",
				`role="tab"`,
				"flex-1",
			},
			wantNot: []string{
				`<button`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.Tab("test-tab", g.Text("Test"), tt.opts...)
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

func TestTabAccessibility(t *testing.T) {
	t.Run("Tab maintains Alpine.js directives", func(t *testing.T) {
		var buf bytes.Buffer

		node := tabs.Tab("tab1", g.Text("Tab 1"))
		if err := node.Render(&buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()
		requiredDirectives := []string{
			"@click",
			"@keydown.right.prevent",
			"@keydown.left.prevent",
			"x-init",
			"focusNextTab",
			"focusPrevTab",
		}

		for _, directive := range requiredDirectives {
			if !strings.Contains(html, directive) {
				t.Errorf("Tab missing required Alpine.js directive: %s", directive)
			}
		}
	})

	t.Run("Tab link maintains Alpine.js directives", func(t *testing.T) {
		var buf bytes.Buffer

		node := tabs.Tab("tab1", g.Text("Tab 1"), tabs.WithHref("/tab1"))
		if err := node.Render(&buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()
		requiredDirectives := []string{
			"@click",
			"@keydown.right.prevent",
			"@keydown.left.prevent",
			":aria-selected",
			":tabindex",
			`role="tab"`,
		}

		for _, directive := range requiredDirectives {
			if !strings.Contains(html, directive) {
				t.Errorf("Tab link missing required directive: %s", directive)
			}
		}
	})
}

func TestTabWithActive(t *testing.T) {
	tests := []struct {
		name    string
		opts    []tabs.TabOption
		want    []string
		wantNot []string
	}{
		{
			name: "Tab without active has no initial active state",
			opts: nil,
			want: []string{
				`role="tab"`,
				`:data-state`,
				`focusNextTab`,
				`focusPrevTab`,
			},
			wantNot: []string{
				`	activeTab = &#39;test-tab&#39;;`,
			},
		},
		{
			name: "Tab with active=true has initial active state",
			opts: []tabs.TabOption{tabs.WithActive(true)},
			want: []string{
				`x-init=`,
				`	activeTab = &#39;test-tab&#39;;`,
				`:data-state`,
				`:aria-selected`,
				`:tabindex`,
				`role="tab"`,
			},
		},
		{
			name: "Tab with active=false has no initial active state",
			opts: []tabs.TabOption{tabs.WithActive(false)},
			want: []string{
				`role="tab"`,
				`:data-state`,
			},
			wantNot: []string{
				`	activeTab = &#39;test-tab&#39;;`,
			},
		},
		{
			name: "Active tab as link",
			opts: []tabs.TabOption{
				tabs.WithActive(true),
				tabs.WithHref("/active"),
			},
			want: []string{
				`<a`,
				`href="/active"`,
				`x-init=`,
				`	activeTab = &#39;test-tab&#39;;`,
				`:aria-selected`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.Tab("test-tab", g.Text("Test"), tt.opts...)
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

func TestTabWithAttr(t *testing.T) {
	t.Run("WithAttr adds single attribute", func(t *testing.T) {
		var buf bytes.Buffer

		node := tabs.Tab("test-tab", g.Text("Test"),
			tabs.WithAttr(g.Attr("data-testid", "my-tab")),
		)
		if err := node.Render(&buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `data-testid="my-tab"`) {
			t.Errorf("Tab missing WithAttr attribute\ngot: %s", html)
		}
	})

	t.Run("WithAttr can be used multiple times", func(t *testing.T) {
		var buf bytes.Buffer

		node := tabs.Tab("test-tab", g.Text("Test"),
			tabs.WithAttr(g.Attr("data-testid", "tab1")),
			tabs.WithAttr(g.Attr("data-category", "nav")),
		)
		if err := node.Render(&buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `data-testid="tab1"`) {
			t.Errorf("Tab missing first WithAttr attribute\ngot: %s", html)
		}
		if !strings.Contains(html, `data-category="nav"`) {
			t.Errorf("Tab missing second WithAttr attribute\ngot: %s", html)
		}
	})

	t.Run("WithAttr and WithTabAttrs work together", func(t *testing.T) {
		var buf bytes.Buffer

		node := tabs.Tab("test-tab", g.Text("Test"),
			tabs.WithAttr(g.Attr("data-single", "one")),
			tabs.WithTabAttrs(
				g.Attr("data-multi1", "two"),
				g.Attr("data-multi2", "three"),
			),
		)
		if err := node.Render(&buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()
		expectedAttrs := []string{
			`data-single="one"`,
			`data-multi1="two"`,
			`data-multi2="three"`,
		}

		for _, attr := range expectedAttrs {
			if !strings.Contains(html, attr) {
				t.Errorf("Tab missing attribute %q\ngot: %s", attr, html)
			}
		}
	})
}

func TestTabWithTabAttrs(t *testing.T) {
	tests := []struct {
		name     string
		opts     []tabs.TabOption
		wantAll  []string
		wantNone []string
	}{
		{
			name: "WithTabAttrs adds single attribute",
			opts: []tabs.TabOption{
				tabs.WithTabAttrs(g.Attr("data-test", "value")),
			},
			wantAll: []string{
				`data-test="value"`,
				`role="tab"`,
			},
		},
		{
			name: "WithTabAttrs adds multiple attributes",
			opts: []tabs.TabOption{
				tabs.WithTabAttrs(
					g.Attr("data-id", "123"),
					g.Attr("data-name", "test-tab"),
					g.Attr("data-category", "navigation"),
				),
			},
			wantAll: []string{
				`data-id="123"`,
				`data-name="test-tab"`,
				`data-category="navigation"`,
			},
		},
		{
			name: "WithTabAttrs can be called multiple times",
			opts: []tabs.TabOption{
				tabs.WithTabAttrs(g.Attr("data-first", "1")),
				tabs.WithTabAttrs(g.Attr("data-second", "2"), g.Attr("data-third", "3")),
			},
			wantAll: []string{
				`data-first="1"`,
				`data-second="2"`,
				`data-third="3"`,
			},
		},
		{
			name: "WithTabAttrs works with other options",
			opts: []tabs.TabOption{
				tabs.WithHref("/test"),
				tabs.WithTabClass("custom"),
				tabs.WithTabAttrs(
					g.Attr("data-analytics", "track"),
					g.Attr("data-label", "Test Link"),
				),
				tabs.WithGrow(),
			},
			wantAll: []string{
				`<a`,
				`href="/test"`,
				"custom",
				`data-analytics="track"`,
				`data-label="Test Link"`,
				"flex-1",
			},
		},
		{
			name: "WithTabAttrs with ARIA attributes",
			opts: []tabs.TabOption{
				tabs.WithTabAttrs(
					g.Attr("aria-label", "Custom Label"),
					g.Attr("aria-describedby", "desc-123"),
					g.Attr("aria-controls", "panel-123"),
				),
			},
			wantAll: []string{
				`aria-label="Custom Label"`,
				`aria-describedby="desc-123"`,
				`aria-controls="panel-123"`,
			},
		},
		{
			name: "WithTabAttrs with empty call",
			opts: []tabs.TabOption{
				tabs.WithTabAttrs(),
			},
			wantAll: []string{
				`role="tab"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := tabs.Tab("test-tab", g.Text("Test"), tt.opts...)
			if err := node.Render(&buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			for _, want := range tt.wantAll {
				if !strings.Contains(html, want) {
					t.Errorf("output missing %q\ngot: %s", want, html)
				}
			}

			for _, wantNone := range tt.wantNone {
				if strings.Contains(html, wantNone) {
					t.Errorf("output should not contain %q\ngot: %s", wantNone, html)
				}
			}
		})
	}
}
