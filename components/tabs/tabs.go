// Package tabs provides tab components following shadcn/ui patterns.
// Tabs use Alpine.js for state management and support keyboard navigation.
package tabs

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
)

// Orientation defines tab layout direction
type Orientation string

const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// TabsProps defines tabs configuration
type TabsProps struct {
	DefaultTab  string
	Orientation Orientation
	Class       string
	Attrs       []g.Node
}

// Option is a functional option for configuring tabs
type Option func(*TabsProps)

// WithDefaultTab sets the initially active tab
func WithDefaultTab(id string) Option {
	return func(p *TabsProps) { p.DefaultTab = id }
}

// WithOrientation sets the tab orientation
func WithOrientation(o Orientation) Option {
	return func(p *TabsProps) { p.Orientation = o }
}

// WithClass adds custom classes
func WithClass(class string) Option {
	return func(p *TabsProps) { p.Class = class }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs ...g.Node) Option {
	return func(p *TabsProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultProps returns default property values
func defaultProps() *TabsProps {
	return &TabsProps{
		Orientation: OrientationHorizontal,
	}
}

// TabListVariant defines the visual style of the tab list
type TabListVariant string

const (
	TabListVariantDefault   TabListVariant = "default"   // Boxed style with background
	TabListVariantUnderline TabListVariant = "underline" // Minimal with bottom border
	TabListVariantPills     TabListVariant = "pills"     // Rounded pill style
)

// TabListProps defines tab list configuration
type TabListProps struct {
	Variant    TabListVariant
	Scrollable bool
	Class      string
	Attrs      []g.Node
}

// TabListOption is a functional option for configuring tab list
type TabListOption func(*TabListProps)

// WithTabListVariant sets the visual style of the tab list
func WithTabListVariant(variant TabListVariant) TabListOption {
	return func(p *TabListProps) { p.Variant = variant }
}

// WithScrollable enables horizontal scrolling for overflow tabs
func WithScrollable() TabListOption {
	return func(p *TabListProps) { p.Scrollable = true }
}

// WithTabListClass adds custom classes to the tab list
func WithTabListClass(class string) TabListOption {
	return func(p *TabListProps) { p.Class = class }
}

// WithTabListAttrs adds custom attributes to the tab list
func WithTabListAttrs(attrs ...g.Node) TabListOption {
	return func(p *TabListProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultTabListProps returns default tab list properties
func defaultTabListProps() *TabListProps {
	return &TabListProps{
		Variant: TabListVariantDefault,
	}
}

// TabVariant defines the visual style of individual tabs
type TabVariant string

const (
	TabVariantDefault   TabVariant = "default"   // Boxed style matching default list
	TabVariantUnderline TabVariant = "underline" // Minimal with bottom border
	TabVariantPills     TabVariant = "pills"     // Rounded pill style
)

// TabProps defines individual tab configuration
type TabProps struct {
	Variant TabVariant
	Href    string
	Class   string
	Attrs   []g.Node
	Shrink  bool // Deprecated: No longer needed, tabs default to natural width
	Grow    bool // Add flex-1 for equal width distribution
	Active  bool
}

// TabOption is a functional option for configuring individual tabs
type TabOption func(*TabProps)

// WithTabVariant sets the visual style of the tab
func WithTabVariant(variant TabVariant) TabOption {
	return func(p *TabProps) { p.Variant = variant }
}

// WithHref makes the tab render as a link
func WithHref(href string) TabOption {
	return func(p *TabProps) { p.Href = href }
}

// WithTabClass adds custom classes to the tab
func WithTabClass(class string) TabOption {
	return func(p *TabProps) { p.Class = class }
}

// WithTabAttrs adds custom attributes to the tab
func WithTabAttrs(attrs ...g.Node) TabOption {
	return func(p *TabProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// WithAttr adds a single custom attribute to the tab
func WithAttr(attr g.Node) TabOption {
	return func(p *TabProps) { p.Attrs = append(p.Attrs, attr) }
}

// WithShrink removes flex-1 class for natural width sizing
// Deprecated: No longer needed. Tabs now default to natural width.
// Use WithGrow() if you want equal width distribution.
func WithShrink() TabOption {
	return func(p *TabProps) { p.Shrink = true }
}

// WithGrow adds flex-1 class for equal width distribution
// By default, tabs size naturally to their content.
func WithGrow() TabOption {
	return func(p *TabProps) { p.Grow = true }
}

// WithActive marks the tab as initially active
func WithActive(active bool) TabOption {
	return func(p *TabProps) { p.Active = active }
}

// defaultTabProps returns default tab properties
func defaultTabProps() *TabProps {
	return &TabProps{
		Variant: TabVariantDefault,
	}
}

// Tabs creates a tabs container with Alpine.js state management
//
// IMPORTANT: Use WithDefaultTab to set the initially active tab
//
// Example:
//
//	tabs.TabsWithOptions(
//	    []tabs.Option{tabs.WithDefaultTab("overview")},
//	    tabs.TabList(
//	        tabs.Tab("overview", g.Text("Overview")),
//	        tabs.Tab("details", g.Text("Details")),
//	    ),
//	    tabs.TabPanel("overview", g.Text("Overview content")),
//	    tabs.TabPanel("details", g.Text("Details content")),
//	)
func Tabs(children ...g.Node) g.Node {
	return TabsWithOptions(nil, children...)
}

// TabsWithOptions creates tabs with custom options
func TabsWithOptions(opts []Option, children ...g.Node) g.Node {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	defaultTab := props.DefaultTab
	if defaultTab == "" {
		defaultTab = "tab1"
	}

	classes := "w-full"
	if props.Class != "" {
		classes += " " + props.Class
	}

	orientation := string(props.Orientation)

	attrs := []g.Node{
		html.Class(classes),
		alpine.XData(map[string]any{
			"activeTab": defaultTab,
		}),
		g.Attr("data-orientation", orientation),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// TabList creates a wrapper for tab buttons
//
// Example:
//
//	tabs.TabList(
//	    tabs.Tab("tab1", g.Text("Tab 1")),
//	    tabs.Tab("tab2", g.Text("Tab 2")),
//	)
//
// With scrollable option for many tabs:
//
//	tabs.TabList(
//	    tabs.Tab("tab1", g.Text("Tab 1")),
//	    tabs.Tab("tab2", g.Text("Tab 2")),
//	    tabs.WithScrollable(),
//	)
//
// With equal width distribution:
//
//	tabs.TabList(
//	    tabs.Tab("tab1", g.Text("Tab 1"), tabs.WithGrow()),
//	    tabs.Tab("tab2", g.Text("Tab 2"), tabs.WithGrow()),
//	)
func TabList(children ...g.Node) g.Node {
	return TabListWithOptions(nil, children...)
}

// TabListWithOptions creates a tab list with custom options
func TabListWithOptions(opts []TabListOption, children ...g.Node) g.Node {
	props := defaultTabListProps()
	for _, opt := range opts {
		opt(props)
	}

	// Base classes for all variants
	baseClasses := "inline-flex items-center"

	// Variant-specific styles
	switch props.Variant {
	case TabListVariantUnderline:
		baseClasses += " h-10 gap-6 border-b border-border"
	case TabListVariantPills:
		baseClasses += " h-10 gap-2 p-1"
	default: // TabListVariantDefault
		baseClasses += " h-10 justify-center rounded-md bg-muted p-1 text-muted-foreground"
	}

	// Scrollable handling
	if props.Scrollable {
		if props.Variant == TabListVariantDefault {
			baseClasses += " overflow-x-auto scroll-smooth gap-1"
		} else {
			baseClasses += " overflow-x-auto scroll-smooth"
		}
	} else {
		baseClasses += " w-full"
	}

	classes := baseClasses
	if props.Class != "" {
		classes = forgeui.CN(baseClasses, props.Class)
	}

	attrs := []g.Node{
		g.Attr("role", "tablist"),
		html.Class(classes),
	}

	// Add scrollable data attribute so child tabs can detect it
	if props.Scrollable {
		attrs = append(attrs, g.Attr("data-scrollable", "true"))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// Tab creates an individual tab button or link
// Tabs default to natural width (sized to content). Use WithGrow() for equal width distribution.
//
// Example:
//
//	tabs.Tab("overview", g.Text("Overview"))
//
// With custom options:
//
//	tabs.Tab("overview", g.Text("Overview"),
//	    tabs.WithHref("/overview"),
//	    tabs.WithTabClass("custom-class"),
//	    tabs.WithGrow(),  // Equal width distribution
//	    tabs.WithActive(true),
//	    tabs.WithAttr(g.Attr("data-testid", "overview-tab")),
//	)
func Tab(id string, label g.Node, opts ...TabOption) g.Node {
	props := defaultTabProps()
	for _, opt := range opts {
		opt(props)
	}

	// Base classes common to all variants (natural width by default)
	baseClasses := "inline-flex items-center justify-center whitespace-nowrap text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"

	// Add flex-1 for equal width distribution if requested
	if props.Grow {
		baseClasses += " flex-1"
	}

	// Variant-specific styles
	switch props.Variant {
	case TabVariantUnderline:
		baseClasses += " rounded-md px-3 py-2.5 border-b border-transparent hover:text-foreground hover:bg-muted/50 data-[state=active]:border-b data-[state=active]:border-primary data-[state=active]:text-foreground"
	case TabVariantPills:
		baseClasses += " rounded-full px-3 py-1.5 hover:bg-muted/50 data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
	default: // TabVariantDefault
		baseClasses += " rounded-sm px-3 py-1.5 data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-sm"
	}

	classes := baseClasses
	if props.Class != "" {
		classes = forgeui.CN(baseClasses, props.Class)
	}

	// Build x-init directive for keyboard navigation
	xInitContent := `
		focusNextTab = (el) => {
			const tabs = Array.from(el.parentElement.querySelectorAll('[role=tab]'));
			const idx = tabs.indexOf(el);
			const next = tabs[(idx + 1) % tabs.length];
			if (next) next.click();
		};
		focusPrevTab = (el) => {
			const tabs = Array.from(el.parentElement.querySelectorAll('[role=tab]'));
			const idx = tabs.indexOf(el);
			const prev = tabs[(idx - 1 + tabs.length) % tabs.length];
			if (prev) prev.click();
		};
	`

	// If active, set the activeTab on init
	if props.Active {
		xInitContent += fmt.Sprintf("\n\t\tactiveTab = '%s';", id)
	}

	// Common attributes for both button and link
	commonAttrs := []g.Node{
		g.Attr("role", "tab"),
		g.Attr(":aria-selected", fmt.Sprintf("activeTab === '%s'", id)),
		g.Attr(":tabindex", fmt.Sprintf("activeTab === '%s' ? 0 : -1", id)),
		alpine.XOn("click", fmt.Sprintf("activeTab = '%s'", id)),
		alpine.XOn("keydown.right.prevent", "focusNextTab($el)"),
		alpine.XOn("keydown.left.prevent", "focusPrevTab($el)"),
		g.Attr("x-init", xInitContent),
		html.Class(classes),
		g.Attr(":data-state", fmt.Sprintf("activeTab === '%s' ? 'active' : 'inactive'", id)),
	}

	commonAttrs = append(commonAttrs, props.Attrs...)

	// Render as link if href is provided, otherwise as button
	if props.Href != "" {
		return html.A(
			g.Attr("href", props.Href),
			g.Group(commonAttrs),
			label,
		)
	}

	return html.Button(
		g.Attr("type", "button"),
		g.Group(commonAttrs),
		label,
	)
}

// TabPanel creates a content panel for a tab
//
// Example:
//
//	tabs.TabPanel("overview", g.Text("Overview content here"))
func TabPanel(id string, content ...g.Node) g.Node {
	return html.Div(
		g.Attr("role", "tabpanel"),
		g.Attr(":aria-hidden", fmt.Sprintf("activeTab !== '%s'", id)),
		g.Attr("x-show", fmt.Sprintf("activeTab === '%s'", id)),
		g.Attr("x-transition:enter", "transition ease-out duration-200"),
		g.Attr("x-transition:enter-start", "opacity-0"),
		g.Attr("x-transition:enter-end", "opacity-100"),
		html.Class("mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"),
		g.Attr("tabindex", "0"),
		g.Group(content),
	)
}
