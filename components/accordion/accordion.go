// Package accordion provides collapsible accordion components following shadcn/ui patterns.
// Accordion uses Alpine.js Collapse plugin for smooth animations.
package accordion

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
)

// Type defines accordion behavior
type Type string

const (
	TypeSingle   Type = "single"   // Only one item can be open at a time
	TypeMultiple Type = "multiple" // Multiple items can be open
)

// AccordionProps defines accordion configuration
type AccordionProps struct {
	Type         Type
	Collapsible  bool // Allow closing all items (for single type)
	DefaultValue []string
	Class        string
	Attrs        []g.Node
}

// Option is a functional option for configuring accordion
type Option func(*AccordionProps)

// WithType sets the accordion type (single or multiple)
func WithType(t Type) Option {
	return func(p *AccordionProps) { p.Type = t }
}

// WithCollapsible allows closing all items
func WithCollapsible() Option {
	return func(p *AccordionProps) { p.Collapsible = true }
}

// WithDefaultValue sets initially open items
func WithDefaultValue(items ...string) Option {
	return func(p *AccordionProps) { p.DefaultValue = items }
}

// WithClass adds custom classes
func WithClass(class string) Option {
	return func(p *AccordionProps) { p.Class = class }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs ...g.Node) Option {
	return func(p *AccordionProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultProps returns default property values
func defaultProps() *AccordionProps {
	return &AccordionProps{
		Type:         TypeSingle,
		Collapsible:  false,
		DefaultValue: []string{},
	}
}

// Accordion creates an accordion container
//
// Example:
//
//	accordion.Accordion(
//	    accordion.WithType(accordion.TypeMultiple),
//	    accordion.Item("item1", "What is ForgeUI?",
//	        g.Text("A Go UI library...")),
//	    accordion.Item("item2", "How to install?",
//	        g.Text("Run go get...")),
//	)
func Accordion(children ...g.Node) g.Node {
	return AccordionWithOptions(nil, children...)
}

// AccordionWithOptions creates accordion with custom options
func AccordionWithOptions(opts []Option, children ...g.Node) g.Node {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := "w-full"
	if props.Class != "" {
		classes += " " + props.Class
	}

	// Build initial state based on type
	var initialState map[string]any
	if props.Type == TypeSingle {
		defaultOpen := ""
		if len(props.DefaultValue) > 0 {
			defaultOpen = props.DefaultValue[0]
		}
		initialState = map[string]any{
			"openItem":    defaultOpen,
			"type":        "single",
			"collapsible": props.Collapsible,
		}
	} else {
		initialState = map[string]any{
			"openItems": props.DefaultValue,
			"type":      "multiple",
		}
	}

	attrs := []g.Node{
		html.Class(classes),
		alpine.XData(initialState),
		g.Attr("data-orientation", "vertical"),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// Item creates an accordion item with trigger and content
//
// Example:
//
//	accordion.Item("item1", "Question?", g.Text("Answer."))
func Item(id, title string, content ...g.Node) g.Node {
	return html.Div(
		html.Class("border-border border-b"),
		g.Attr("data-state", ""),
		g.Attr("x-init", fmt.Sprintf(`
			$watch('type === "single" ? openItem : openItems', (val) => {
				const isOpen = type === 'single' ? val === '%s' : val.includes('%s');
				$el.setAttribute('data-state', isOpen ? 'open' : 'closed');
			});
		`, id, id)),

		// Trigger
		html.H3(
			html.Class("flex"),
			html.Button(
				g.Attr("type", "button"),
				html.Class("flex flex-1 items-center justify-between py-4 font-medium transition-all hover:underline text-sm [&[data-state=open]>svg]:rotate-180"),
				g.Attr("x-on:click", fmt.Sprintf(`
					if (type === 'single') {
						if (openItem === '%s') {
							if (collapsible) openItem = '';
						} else {
							openItem = '%s';
						}
					} else {
						if (openItems.includes('%s')) {
							openItems = openItems.filter(i => i !== '%s');
						} else {
							openItems = [...openItems, '%s'];
						}
					}
				`, id, id, id, id, id)),
				g.Attr(":data-state", fmt.Sprintf(
					"type === 'single' ? (openItem === '%s' ? 'open' : 'closed') : (openItems.includes('%s') ? 'open' : 'closed')",
					id, id,
				)),
				g.Text(title),

				// Chevron icon
				g.El("svg",
					html.Class("h-4 w-4 shrink-0 transition-transform duration-200"),
					g.Attr("xmlns", "http://www.w3.org/2000/svg"),
					g.Attr("width", "24"),
					g.Attr("height", "24"),
					g.Attr("viewBox", "0 0 24 24"),
					g.Attr("fill", "none"),
					g.Attr("stroke", "currentColor"),
					g.Attr("stroke-width", "2"),
					g.Attr("stroke-linecap", "round"),
					g.Attr("stroke-linejoin", "round"),
					g.El("path", g.Attr("d", "m6 9 6 6 6-6")),
				),
			),
		),

		// Content with collapse animation
		html.Div(
			g.Attr("x-show", fmt.Sprintf(
				"type === 'single' ? openItem === '%s' : openItems.includes('%s')",
				id, id,
			)),
			g.Attr("x-collapse", ""),
			g.Attr("role", "region"),
			html.Class("overflow-hidden text-sm transition-all"),

			html.Div(
				html.Class("pb-4 pt-0"),
				g.Group(content),
			),
		),
	)
}
