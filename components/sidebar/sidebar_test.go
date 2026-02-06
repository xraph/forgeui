package sidebar

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Helper function to render a node to string
func renderNode(t *testing.T, node g.Node) string {
	t.Helper()
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	return buf.String()
}

// Helper function to check if string contains substring
func assertContains(t *testing.T, str, substr string) {
	t.Helper()

	if !strings.Contains(str, substr) {
		t.Errorf("Expected string to contain %q", substr)
	}
}

func TestSidebar(t *testing.T) {
	t.Run("renders basic sidebar with provider", func(t *testing.T) {
		sidebar := Sidebar(
			SidebarHeader(g.Text("My App")),
			SidebarContent(g.Text("Content")),
		)

		rendered := renderNode(t, sidebar)

		// Should have provider wrapper
		assertContains(t, rendered, `data-provider="sidebar"`)

		// Should have sidebar state
		assertContains(t, rendered, `collapsed`)
		assertContains(t, rendered, `mobileOpen`)
		assertContains(t, rendered, `isMobile`)

		// Should have content
		assertContains(t, rendered, "My App")
		assertContains(t, rendered, "Content")
	})

	t.Run("renders with custom options", func(t *testing.T) {
		sidebar := SidebarWithOptions(
			[]SidebarOption{
				WithDefaultCollapsed(true),
				WithCollapsible(false),
				WithSidebarClass("custom-sidebar"),
			},
			SidebarContent(g.Text("Content")),
		)

		rendered := renderNode(t, sidebar)

		assertContains(t, rendered, `data-provider="sidebar"`)
		assertContains(t, rendered, `custom-sidebar`)
	})

	t.Run("includes Alpine store integration", func(t *testing.T) {
		sidebar := Sidebar()

		rendered := renderNode(t, sidebar)

		// Should have Alpine store for state management
		assertContains(t, rendered, `Alpine.store('sidebar'`)
		assertContains(t, rendered, `collapsed`)
		assertContains(t, rendered, `mobileOpen`)
		assertContains(t, rendered, `$store.sidebar`)
	})

	t.Run("includes mobile backdrop and toggle", func(t *testing.T) {
		sidebar := Sidebar()

		rendered := renderNode(t, sidebar)

		// Should have backdrop
		assertContains(t, rendered, `backdrop-blur`)

		// Should have mobile toggle button
		assertContains(t, rendered, `Open sidebar`)
	})
}

func TestSidebarHeader(t *testing.T) {
	t.Run("renders header with children", func(t *testing.T) {
		header := SidebarHeader(
			html.Img(g.Attr("src", "/logo.svg")),
			g.Text("My App"),
		)

		rendered := renderNode(t, header)

		assertContains(t, rendered, "/logo.svg")
		assertContains(t, rendered, "My App")
	})

	t.Run("has responsive visibility", func(t *testing.T) {
		header := SidebarHeader(g.Text("Title"))

		rendered := renderNode(t, header)

		// Should adjust padding/layout when collapsed via :class
		assertContains(t, rendered, `:class`)
		assertContains(t, rendered, `collapsed`)
		assertContains(t, rendered, `isMobile`)
	})

	t.Run("renders with custom class", func(t *testing.T) {
		header := SidebarHeaderWithOptions(
			[]SidebarHeaderOption{WithSidebarHeaderClass("custom-header bg-primary")},
			g.Text("Brand"),
		)

		rendered := renderNode(t, header)

		assertContains(t, rendered, "custom-header")
		assertContains(t, rendered, "bg-primary")
		assertContains(t, rendered, "flex items-center") // Base classes preserved
	})
}

func TestSidebarContent(t *testing.T) {
	t.Run("renders content area", func(t *testing.T) {
		content := SidebarContent(
			g.Text("Navigation items"),
		)

		rendered := renderNode(t, content)

		assertContains(t, rendered, "Navigation items")
		assertContains(t, rendered, "overflow-auto")
		assertContains(t, rendered, `data-slot="sidebar-content"`)
	})
}

func TestSidebarFooter(t *testing.T) {
	t.Run("renders footer section", func(t *testing.T) {
		footer := SidebarFooter(
			g.Text("© 2024"),
		)

		rendered := renderNode(t, footer)

		assertContains(t, rendered, "© 2024")
		assertContains(t, rendered, "border-t")
	})

	t.Run("renders with custom class", func(t *testing.T) {
		footer := SidebarFooterWithOptions(
			[]SidebarFooterOption{WithSidebarFooterClass("custom-footer bg-secondary")},
			g.Text("© 2024"),
		)

		rendered := renderNode(t, footer)

		assertContains(t, rendered, "custom-footer")
		assertContains(t, rendered, "bg-secondary")
		assertContains(t, rendered, "border-t") // Base classes preserved
	})
}

func TestSidebarToggle(t *testing.T) {
	t.Run("renders toggle button", func(t *testing.T) {
		toggle := SidebarToggle()

		rendered := renderNode(t, toggle)

		assertContains(t, rendered, "Toggle sidebar")
		assertContains(t, rendered, `x-show`)
		assertContains(t, rendered, `collapsible`)
		assertContains(t, rendered, `collapsed`)
	})
}

func TestSidebarGroup(t *testing.T) {
	t.Run("renders group container", func(t *testing.T) {
		group := SidebarGroup(
			SidebarGroupLabel("Navigation"),
			g.Text("Items"),
		)

		rendered := renderNode(t, group)

		assertContains(t, rendered, "Navigation")
		assertContains(t, rendered, "Items")
	})
}

func TestSidebarGroupCollapsible(t *testing.T) {
	t.Run("renders collapsible group", func(t *testing.T) {
		group := SidebarGroupCollapsible(
			[]SidebarGroupOption{
				WithGroupKey("projects"),
				WithGroupDefaultOpen(true),
			},
			SidebarGroupLabel("Projects"),
			g.Text("Project items"),
		)

		rendered := renderNode(t, group)

		assertContains(t, rendered, `projects_open`)
		assertContains(t, rendered, "Projects")
		assertContains(t, rendered, "Project items")
	})
}

func TestSidebarGroupLabel(t *testing.T) {
	t.Run("renders group label", func(t *testing.T) {
		label := SidebarGroupLabel("Settings")

		rendered := renderNode(t, label)

		assertContains(t, rendered, "Settings")
		assertContains(t, rendered, `data-slot="sidebar-group-label"`)
		assertContains(t, rendered, "text-xs")
		assertContains(t, rendered, "font-medium")
	})

	t.Run("hides when collapsed", func(t *testing.T) {
		label := SidebarGroupLabel("Settings")

		rendered := renderNode(t, label)

		assertContains(t, rendered, `:class`)
		assertContains(t, rendered, `collapsed`)
		assertContains(t, rendered, `isMobile`)
	})
}

func TestSidebarMenu(t *testing.T) {
	t.Run("renders menu container", func(t *testing.T) {
		menu := SidebarMenu(
			SidebarMenuItem(g.Text("Item 1")),
			SidebarMenuItem(g.Text("Item 2")),
		)

		rendered := renderNode(t, menu)

		assertContains(t, rendered, "Item 1")
		assertContains(t, rendered, "Item 2")
		assertContains(t, rendered, "flex-col")
	})
}

func TestSidebarMenuItem(t *testing.T) {
	t.Run("renders menu item container", func(t *testing.T) {
		item := SidebarMenuItem(
			g.Text("Menu content"),
		)

		rendered := renderNode(t, item)

		assertContains(t, rendered, "Menu content")
		assertContains(t, rendered, "group/menu-item")
	})
}

func TestSidebarMenuButton(t *testing.T) {
	t.Run("renders as link by default", func(t *testing.T) {
		button := SidebarMenuButton(
			"Dashboard",
			WithMenuHref("/dashboard"),
		)

		rendered := renderNode(t, button)

		assertContains(t, rendered, "Dashboard")
		assertContains(t, rendered, `href="/dashboard"`)
		assertContains(t, rendered, "<a")
	})

	t.Run("renders as button when specified", func(t *testing.T) {
		button := SidebarMenuButton(
			"Action",
			WithMenuAsButton(),
		)

		rendered := renderNode(t, button)

		assertContains(t, rendered, "Action")
		assertContains(t, rendered, `type="button"`)
		assertContains(t, rendered, "<button")
	})

	t.Run("renders with icon", func(t *testing.T) {
		button := SidebarMenuButton(
			"Dashboard",
			WithMenuIcon(html.Span(g.Text("icon"))),
		)

		rendered := renderNode(t, button)

		assertContains(t, rendered, "icon")
		assertContains(t, rendered, "Dashboard")
	})

	t.Run("renders with badge", func(t *testing.T) {
		button := SidebarMenuButton(
			"Messages",
			WithMenuBadge(SidebarMenuBadge("5")),
		)

		rendered := renderNode(t, button)

		assertContains(t, rendered, "Messages")
		assertContains(t, rendered, "5")
	})

	t.Run("applies active state", func(t *testing.T) {
		button := SidebarMenuButton(
			"Active",
			WithMenuActive(),
		)

		rendered := renderNode(t, button)

		assertContains(t, rendered, "bg-sidebar-accent")
		assertContains(t, rendered, "text-sidebar-accent-foreground")
	})

	t.Run("hides label when collapsed", func(t *testing.T) {
		button := SidebarMenuButton("Label")

		rendered := renderNode(t, button)

		assertContains(t, rendered, `x-show`)
		assertContains(t, rendered, `$store.sidebar`)
		assertContains(t, rendered, `collapsed`)
		assertContains(t, rendered, `isMobile`)
	})
}

func TestSidebarMenuBadge(t *testing.T) {
	t.Run("renders badge with text", func(t *testing.T) {
		badge := SidebarMenuBadge("12")

		rendered := renderNode(t, badge)

		assertContains(t, rendered, "12")
		assertContains(t, rendered, "bg-sidebar-primary")
	})
}

func TestSidebarInset(t *testing.T) {
	t.Run("renders inset with provider integration", func(t *testing.T) {
		inset := SidebarInset(
			html.Main(g.Text("Main content")),
		)

		rendered := renderNode(t, inset)

		assertContains(t, rendered, "Main content")
		assertContains(t, rendered, `:class`)
		// Should use provider value for margin calculation
		assertContains(t, rendered, `md:ml-64`)
		assertContains(t, rendered, `md:ml-16`)
	})
}

func TestSidebarInsetHeader(t *testing.T) {
	t.Run("renders inset header with default classes", func(t *testing.T) {
		header := SidebarInsetHeader(
			g.Text("Dashboard"),
		)

		rendered := renderNode(t, header)

		assertContains(t, rendered, "Dashboard")
		assertContains(t, rendered, "sticky top-0")
		assertContains(t, rendered, "border-b")
	})

	t.Run("renders with custom class", func(t *testing.T) {
		header := SidebarInsetHeaderWithOptions(
			[]SidebarInsetHeaderOption{WithSidebarInsetHeaderClass("custom-inset-header shadow-lg")},
			g.Text("Dashboard"),
		)

		rendered := renderNode(t, header)

		assertContains(t, rendered, "custom-inset-header")
		assertContains(t, rendered, "shadow-lg")
		assertContains(t, rendered, "sticky top-0") // Base classes preserved
		assertContains(t, rendered, "border-b")
	})
}

func TestSidebarLayoutContent(t *testing.T) {
	t.Run("renders layout content with provider integration", func(t *testing.T) {
		content := SidebarLayoutContent(
			html.Div(g.Text("Page content")),
		)

		rendered := renderNode(t, content)

		assertContains(t, rendered, "Page content")
		assertContains(t, rendered, `:class`)
		assertContains(t, rendered, `ml-64`)
		assertContains(t, rendered, `ml-16`)
		assertContains(t, rendered, `ml-0`)
	})
}

func TestSidebarTrigger(t *testing.T) {
	t.Run("renders trigger button", func(t *testing.T) {
		trigger := SidebarTrigger()

		rendered := renderNode(t, trigger)

		assertContains(t, rendered, "Toggle sidebar")
		assertContains(t, rendered, `type="button"`)
		// Should use provider method
		assertContains(t, rendered, `@click`)
	})
}

func TestSidebarIntegration(t *testing.T) {
	t.Run("complete sidebar with all components", func(t *testing.T) {
		sidebar := Sidebar(
			SidebarHeader(g.Text("App Name")),
			SidebarContent(
				SidebarGroup(
					SidebarGroupLabel("Main"),
					SidebarMenu(
						SidebarMenuItem(
							SidebarMenuButton(
								"Dashboard",
								WithMenuHref("/dashboard"),
								WithMenuActive(),
							),
						),
					),
				),
			),
			SidebarFooter(g.Text("Footer")),
			SidebarToggle(),
		)

		rendered := renderNode(t, sidebar)

		// Verify provider integration
		assertContains(t, rendered, `data-provider="sidebar"`)

		// Verify all sections
		assertContains(t, rendered, "App Name")
		assertContains(t, rendered, "Main")
		assertContains(t, rendered, "Dashboard")
		assertContains(t, rendered, "/dashboard")
		assertContains(t, rendered, "Footer")

		// Verify toggle
		assertContains(t, rendered, "Toggle sidebar")
	})
}
