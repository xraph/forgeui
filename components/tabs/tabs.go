// Package tabs provides tab components following shadcn/ui patterns.
// Tabs use Alpine.js for state management and support keyboard navigation.
package tabs

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

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

// Tabs creates a tabs container with Alpine.js state management
//
// Example:
//
//	tabs.Tabs(
//	    tabs.WithDefaultTab("overview"),
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
func TabList(children ...g.Node) g.Node {
	return html.Div(
		g.Attr("role", "tablist"),
		html.Class("inline-flex h-10 items-center justify-center rounded-md bg-muted p-1 text-muted-foreground w-full"),
		g.Group(children),
	)
}

// Tab creates an individual tab button
//
// Example:
//
//	tabs.Tab("overview", g.Text("Overview"))
func Tab(id string, label g.Node) g.Node {
	return html.Button(
		g.Attr("type", "button"),
		g.Attr("role", "tab"),
		g.Attr(":aria-selected", fmt.Sprintf("activeTab === '%s'", id)),
		g.Attr(":tabindex", fmt.Sprintf("activeTab === '%s' ? 0 : -1", id)),
		alpine.XOn("click", fmt.Sprintf("activeTab = '%s'", id)),
		alpine.XOn("keydown.right.prevent", "focusNextTab($el)"),
		alpine.XOn("keydown.left.prevent", "focusPrevTab($el)"),
		g.Attr("x-init", `
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
		`),
		html.Class(fmt.Sprintf(
			"inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-sm flex-1",
		)),
		g.Attr(":data-state", fmt.Sprintf("activeTab === '%s' ? 'active' : 'inactive'", id)),
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

