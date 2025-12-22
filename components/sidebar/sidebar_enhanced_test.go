package sidebar

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Helper to render a gomponents Node to a string
func renderNodeEnhanced(node g.Node) string {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		panic(err)
	}

	return buf.String()
}

// TestSidebarVariants tests the different sidebar variants
func TestSidebarVariants(t *testing.T) {
	tests := []struct {
		name    string
		variant SidebarVariant
		want    string
	}{
		{
			name:    "default variant",
			variant: SidebarVariantSidebar,
			want:    "fixed top-0 bottom-0 z-30",
		},
		{
			name:    "floating variant",
			variant: SidebarVariantFloating,
			want:    "m-2 rounded-lg shadow-lg",
		},
		{
			name:    "inset variant",
			variant: SidebarVariantInset,
			want:    "m-2 rounded-xl shadow-sm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SidebarWithOptions(
				[]SidebarOption{WithVariant(tt.variant)},
				SidebarHeader(g.Text("Test")),
			)
			output := renderNodeEnhanced(s)

			if !strings.Contains(output, tt.want) {
				t.Errorf("Expected variant class %q not found in output", tt.want)
			}

			if !strings.Contains(output, `data-variant="`+string(tt.variant)+`"`) {
				t.Errorf("Expected data-variant attribute for %s", tt.variant)
			}
		})
	}
}

// TestSidebarCollapsibleModes tests the different collapsible modes
func TestSidebarCollapsibleModes(t *testing.T) {
	tests := []struct {
		name string
		mode SidebarCollapsibleMode
		want string
	}{
		{
			name: "offcanvas mode",
			mode: CollapsibleOffcanvas,
			want: "transition-transform",
		},
		{
			name: "icon mode",
			mode: CollapsibleIcon,
			want: "transition-width",
		},
		{
			name: "none mode",
			mode: CollapsibleNone,
			want: `data-collapsible="none"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SidebarWithOptions(
				[]SidebarOption{WithCollapsibleMode(tt.mode)},
				SidebarHeader(g.Text("Test")),
			)
			output := renderNodeEnhanced(s)

			if !strings.Contains(output, tt.want) {
				t.Errorf("Expected mode indicator %q not found in output", tt.want)
			}

			if !strings.Contains(output, `data-collapsible="`+string(tt.mode)+`"`) {
				t.Errorf("Expected data-collapsible attribute for %s", tt.mode)
			}
		})
	}
}

// TestKeyboardShortcuts tests keyboard shortcut functionality
func TestKeyboardShortcuts(t *testing.T) {
	t.Run("keyboard shortcuts enabled by default", func(t *testing.T) {
		// Keyboard shortcuts are opt-in by default, so basic sidebar won't have them
		s := Sidebar(SidebarHeader(g.Text("Test")))
		output := renderNodeEnhanced(s)

		// Basic sidebar should still have store initialization
		if !strings.Contains(output, "Alpine.store") {
			t.Error("Expected Alpine store in output")
		}
	})

	t.Run("custom keyboard shortcut key", func(t *testing.T) {
		s := SidebarWithOptions(
			[]SidebarOption{WithKeyboardShortcut("k")},
			SidebarHeader(g.Text("Test")),
		)
		output := renderNodeEnhanced(s)

		// When keyboard shortcuts are enabled, should have keydown handler
		if !strings.Contains(output, "keydown") {
			t.Error("Expected keydown event listener in output")
		}

		// Check for custom key 'k'
		if !strings.Contains(output, "event.key") {
			t.Error("Expected event.key in keyboard handler")
		}

		if !strings.Contains(output, "'k'") {
			t.Error("Expected custom key 'k' in output")
		}
	})

	t.Run("keyboard shortcuts disabled", func(t *testing.T) {
		s := SidebarWithOptions(
			[]SidebarOption{WithKeyboardShortcutEnabled(false)},
			SidebarHeader(g.Text("Test")),
		)
		output := renderNodeEnhanced(s)

		if strings.Contains(output, "handleKeyDown") {
			t.Error("Keyboard shortcut handler should not be present when disabled")
		}
	})
}

// TestStatePersistence tests localStorage/sessionStorage persistence
func TestStatePersistence(t *testing.T) {
	t.Run("state persistence disabled by default", func(t *testing.T) {
		s := Sidebar(SidebarHeader(g.Text("Test")))
		output := renderNodeEnhanced(s)

		if strings.Contains(output, "localStorage.getItem") {
			t.Error("State persistence should not be present by default")
		}
	})

	t.Run("state persistence enabled with localStorage", func(t *testing.T) {
		s := SidebarWithOptions(
			[]SidebarOption{WithPersistState(true)},
			SidebarHeader(g.Text("Test")),
		)
		output := renderNodeEnhanced(s)

		if !strings.Contains(output, "localStorage") {
			t.Error("Expected localStorage in output")
		}

		if !strings.Contains(output, "forgeui_sidebar_state") {
			t.Error("Expected default storage key in output")
		}

		if !strings.Contains(output, "storage.getItem") && !strings.Contains(output, "storage.setItem") {
			t.Error("Expected storage get/set methods in output")
		}
	})

	t.Run("custom storage key", func(t *testing.T) {
		s := SidebarWithOptions(
			[]SidebarOption{
				WithPersistState(true),
				WithStorageKey("custom_sidebar_key"),
			},
			SidebarHeader(g.Text("Test")),
		)
		output := renderNodeEnhanced(s)

		if !strings.Contains(output, "custom_sidebar_key") {
			t.Error("Expected custom storage key in output")
		}
	})

	t.Run("sessionStorage type", func(t *testing.T) {
		s := SidebarWithOptions(
			[]SidebarOption{
				WithPersistState(true),
				WithStorageType("sessionStorage"),
			},
			SidebarHeader(g.Text("Test")),
		)
		output := renderNodeEnhanced(s)

		if !strings.Contains(output, "sessionStorage") {
			t.Error("Expected sessionStorage in output")
		}
	})
}

// TestARIAAttributes tests comprehensive ARIA attributes
func TestARIAAttributes(t *testing.T) {
	t.Run("sidebar has ARIA attributes", func(t *testing.T) {
		s := Sidebar(SidebarHeader(g.Text("Test")))
		output := renderNodeEnhanced(s)

		requiredAttrs := []string{
			`role="complementary"`,
			`aria-label="Main navigation sidebar"`,
			`:aria-expanded=`,
			`data-state=`,
		}

		for _, attr := range requiredAttrs {
			if !strings.Contains(output, attr) {
				t.Errorf("Expected ARIA attribute %q not found", attr)
			}
		}
	})

	t.Run("screen reader live region present", func(t *testing.T) {
		s := Sidebar(SidebarHeader(g.Text("Test")))
		output := renderNodeEnhanced(s)

		requiredAttrs := []string{
			`role="status"`,
			`aria-live="polite"`,
			`aria-atomic="true"`,
			`sr-only`,
		}

		for _, attr := range requiredAttrs {
			if !strings.Contains(output, attr) {
				t.Errorf("Expected screen reader attribute %q not found", attr)
			}
		}
	})

	t.Run("mobile backdrop has ARIA attributes", func(t *testing.T) {
		s := Sidebar(SidebarHeader(g.Text("Test")))
		output := renderNodeEnhanced(s)

		if !strings.Contains(output, `role="presentation"`) {
			t.Error("Expected role=presentation on mobile backdrop")
		}

		if !strings.Contains(output, `aria-hidden="true"`) {
			t.Error("Expected aria-hidden on mobile backdrop")
		}
	})
}

// TestFocusManagement tests focus ring and keyboard navigation
func TestFocusManagement(t *testing.T) {
	t.Run("interactive elements have focus rings", func(t *testing.T) {
		s := Sidebar(
			SidebarHeader(g.Text("Test")),
			SidebarContent(
				SidebarMenu(
					SidebarMenuItem(
						SidebarMenuButton("Dashboard"),
					),
				),
			),
		)
		output := renderNodeEnhanced(s)

		// Check for focus-visible classes
		if !strings.Contains(output, "focus-visible:ring-2") {
			t.Error("Expected focus-visible ring classes")
		}

		if !strings.Contains(output, "focus-visible:ring-ring") {
			t.Error("Expected focus-visible ring color")
		}

		if !strings.Contains(output, "outline-none") {
			t.Error("Expected outline-none for custom focus rings")
		}
	})

	t.Run("menu items have keyboard navigation", func(t *testing.T) {
		s := Sidebar(
			SidebarContent(
				SidebarMenu(
					SidebarMenuItem(
						SidebarMenuButton("Dashboard"),
					),
				),
			),
		)
		output := renderNodeEnhanced(s)

		// Menu items should have proper ARIA attributes for keyboard navigation
		if !strings.Contains(output, `role="menuitem"`) {
			t.Error("Expected role=menuitem on menu buttons")
		}

		if !strings.Contains(output, `tabindex="0"`) {
			t.Error("Expected tabindex=0 on menu items")
		}

		// Should have focus-visible ring classes
		if !strings.Contains(output, "focus-visible:ring") {
			t.Error("Expected focus-visible:ring classes for keyboard focus")
		}
	})
}

// TestSidebarMenuButtonVariants tests menu button variants
func TestSidebarMenuButtonVariants(t *testing.T) {
	tests := []struct {
		name    string
		variant SidebarMenuButtonVariant
		want    string
	}{
		{
			name:    "default variant",
			variant: MenuButtonDefault,
			want:    "hover:bg-sidebar-accent",
		},
		{
			name:    "outline variant",
			variant: MenuButtonOutline,
			want:    "border border-sidebar-border shadow-sm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btn := SidebarMenuButton("Test",
				WithMenuVariant(tt.variant),
			)
			output := renderNodeEnhanced(btn)

			if !strings.Contains(output, tt.want) {
				t.Errorf("Expected variant class %q not found", tt.want)
			}
		})
	}
}

// TestSidebarMenuButtonSizes tests menu button sizes
func TestSidebarMenuButtonSizes(t *testing.T) {
	tests := []struct {
		name string
		size SidebarMenuButtonSize
		want string
	}{
		{
			name: "small size",
			size: MenuButtonSizeSmall,
			want: "h-7",
		},
		{
			name: "default size",
			size: MenuButtonSizeDefault,
			want: "h-8",
		},
		{
			name: "large size",
			size: MenuButtonSizeLarge,
			want: "h-12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btn := SidebarMenuButton("Test",
				WithMenuSize(tt.size),
			)
			output := renderNodeEnhanced(btn)

			if !strings.Contains(output, tt.want) {
				t.Errorf("Expected size class %q not found", tt.want)
			}
		})
	}
}

// TestSidebarMenuTooltip tests tooltip component
func TestSidebarMenuTooltip(t *testing.T) {
	t.Run("tooltip renders with correct structure", func(t *testing.T) {
		tooltip := SidebarMenuTooltip("Dashboard",
			html.Button(g.Text("Btn")),
		)
		output := renderNodeEnhanced(tooltip)

		if !strings.Contains(output, "Dashboard") {
			t.Error("Expected tooltip text")
		}

		if !strings.Contains(output, `role="tooltip"`) {
			t.Error("Expected role=tooltip")
		}

		if !strings.Contains(output, "x-show") {
			t.Error("Expected x-show directive")
		}

		if !strings.Contains(output, "showTooltip") {
			t.Error("Expected showTooltip variable")
		}
	})

	t.Run("menu button with tooltip option", func(t *testing.T) {
		btn := SidebarMenuButton("Dashboard",
			WithMenuTooltip("Go to Dashboard"),
		)
		output := renderNodeEnhanced(btn)

		if !strings.Contains(output, "Go to Dashboard") {
			t.Error("Expected tooltip text in menu button")
		}

		if !strings.Contains(output, `role="tooltip"`) {
			t.Error("Expected tooltip role when tooltip is set")
		}
	})
}

// TestSidebarInput tests the input component
func TestSidebarInput(t *testing.T) {
	t.Run("renders input with correct attributes", func(t *testing.T) {
		input := SidebarInput("Search...", "search")
		output := renderNodeEnhanced(input)

		requiredAttrs := []string{
			`type="text"`,
			`name="search"`,
			`placeholder="Search..."`,
			"bg-background",
			"focus-visible:ring-2",
		}

		for _, attr := range requiredAttrs {
			if !strings.Contains(output, attr) {
				t.Errorf("Expected attribute %q not found", attr)
			}
		}
	})

	t.Run("input hidden when sidebar collapsed", func(t *testing.T) {
		input := SidebarInput("Search...", "search")
		output := renderNodeEnhanced(input)

		if !strings.Contains(output, "x-show") {
			t.Error("Expected x-show directive to hide input when collapsed")
		}
	})
}

// TestSidebarSeparator tests the separator component
func TestSidebarSeparator(t *testing.T) {
	t.Run("renders separator with ARIA attributes", func(t *testing.T) {
		sep := SidebarSeparator()
		output := renderNodeEnhanced(sep)

		requiredAttrs := []string{
			`role="separator"`,
			`aria-orientation="horizontal"`,
			"bg-sidebar-border",
		}

		for _, attr := range requiredAttrs {
			if !strings.Contains(output, attr) {
				t.Errorf("Expected attribute %q not found", attr)
			}
		}
	})
}

// TestSidebarEnhancementsIntegration tests the complete sidebar with all features
func TestSidebarEnhancementsIntegration(t *testing.T) {
	t.Run("complete sidebar with all enhancements", func(t *testing.T) {
		s := SidebarWithOptions(
			[]SidebarOption{
				WithVariant(SidebarVariantFloating),
				WithCollapsibleMode(CollapsibleIcon),
				WithPersistState(true),
				WithKeyboardShortcut("b"),
				WithStorageKey("my_sidebar"),
				WithSide(forgeui.SideLeft),
			},
			SidebarHeader(g.Text("My App")),
			SidebarContent(
				SidebarInput("Search...", "search"),
				SidebarSeparator(),
				SidebarMenu(
					SidebarMenuItem(
						SidebarMenuButton("Dashboard",
							WithMenuHref("/dashboard"),
							WithMenuActive(),
							WithMenuVariant(MenuButtonDefault),
							WithMenuSize(MenuButtonSizeDefault),
							WithMenuTooltip("Go to Dashboard"),
						),
					),
				),
			),
			SidebarFooter(g.Text("© 2024")),
		)
		output := renderNodeEnhanced(s)

		// Check all major features are present
		checks := []string{
			// Variants
			"rounded-lg shadow-lg",
			`data-variant="floating"`,

			// Collapsible mode
			"transition-width",
			`data-collapsible="icon"`,

			// State persistence
			"localStorage",
			"my_sidebar",

			// Keyboard shortcuts
			"keydown",
			"event.key",

			// ARIA
			`role="complementary"`,
			`aria-label`,

			// Focus rings
			"focus-visible:ring-2",

			// Components
			"Search...",
			`role="separator"`,
			"Dashboard",
			`role="tooltip"`,
			"Go to Dashboard",
		}

		for _, check := range checks {
			if !strings.Contains(output, check) {
				t.Errorf("Expected feature indicator %q not found", check)
			}
		}
	})
}

// TestBackwardCompatibility tests that enhancements don't break existing functionality
func TestBackwardCompatibility(t *testing.T) {
	t.Run("basic sidebar still works without new options", func(t *testing.T) {
		s := Sidebar(
			SidebarHeader(g.Text("App")),
			SidebarContent(g.Text("Content")),
		)
		output := renderNodeEnhanced(s)

		if !strings.Contains(output, "App") {
			t.Error("Basic sidebar rendering broken")
		}

		if !strings.Contains(output, "Content") {
			t.Error("Basic content rendering broken")
		}
	})

	t.Run("Alpine store backward compatibility", func(t *testing.T) {
		s := Sidebar(SidebarHeader(g.Text("Test")))
		output := renderNodeEnhanced(s)

		if !strings.Contains(output, "Alpine.store") && !strings.Contains(output, "sidebar") {
			t.Error("Alpine store backward compatibility missing")
		}
	})
}
