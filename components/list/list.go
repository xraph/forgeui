// Package list provides list components for displaying items in ordered
// or unordered lists with various styling options.
//
// # Basic Usage
//
//	list.List(
//	    list.ListItem(g.Text("Item 1")),
//	    list.ListItem(g.Text("Item 2")),
//	    list.ListItem(g.Text("Item 3")),
//	)
//
// # With Icons
//
//	list.List(
//	    list.Icons(),
//	    list.ListItem(
//	        icons.Check(),
//	        g.Text("Completed task"),
//	    ),
//	)
package list

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ListVariant represents the visual style of the list.
type ListVariant string

const (
	VariantBullets ListVariant = "bullets" // Default bullets
	VariantNone    ListVariant = "none"    // No markers
	VariantIcons   ListVariant = "icons"   // Icon markers
)

// Props defines the properties for the List component.
type Props struct {
	Variant ListVariant
	Class   string
	Spaced  bool
	Attrs   []g.Node
}

// Option is a functional option for configuring the List component.
type Option func(*Props)

// WithVariant sets the list variant.
func WithVariant(variant ListVariant) Option {
	return func(p *Props) {
		p.Variant = variant
	}
}

// Bullets displays the list with default bullet points.
func Bullets() Option {
	return func(p *Props) {
		p.Variant = VariantBullets
	}
}

// None displays the list without markers.
func None() Option {
	return func(p *Props) {
		p.Variant = VariantNone
	}
}

// Icons displays the list with icon markers (provided by ListItem).
func Icons() Option {
	return func(p *Props) {
		p.Variant = VariantIcons
	}
}

// Spaced adds extra spacing between list items.
func Spaced() Option {
	return func(p *Props) {
		p.Spaced = true
	}
}

// WithClass adds additional CSS classes to the list.
func WithClass(class string) Option {
	return func(p *Props) {
		p.Class = class
	}
}

// WithAttr adds custom HTML attributes to the list.
func WithAttr(attrs ...g.Node) Option {
	return func(p *Props) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// defaultProps returns the default list properties.
func defaultProps() *Props {
	return &Props{
		Variant: VariantBullets,
	}
}

// List creates an unordered list component.
//
// Example:
//
//	list.List(list.Spaced())(
//	    list.ListItem(g.Text("First item")),
//	    list.ListItem(g.Text("Second item")),
//	)
func List(opts ...Option) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := defaultProps()
		for _, opt := range opts {
			opt(props)
		}

		classes := ""

		switch props.Variant {
		case VariantBullets:
			classes = "list-disc pl-5"
		case VariantNone:
			classes = "list-none"
		case VariantIcons:
			classes = "list-none space-y-2"
		}

		if props.Spaced && props.Variant != VariantIcons {
			classes += " space-y-2"
		}

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{html.Class(classes)}
		attrs = append(attrs, props.Attrs...)

		return html.Ul(
			g.Group(attrs),
			g.Group(children),
		)
	}
}

// OrderedList creates an ordered (numbered) list component.
//
// Example:
//
//	list.OrderedList(list.Spaced())(
//	    list.ListItem(g.Text("Step 1")),
//	    list.ListItem(g.Text("Step 2")),
//	)
func OrderedList(opts ...Option) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &Props{
			Variant: VariantBullets, // Default for numbered
		}
		for _, opt := range opts {
			opt(props)
		}

		classes := "list-decimal pl-5"

		if props.Spaced {
			classes += " space-y-2"
		}

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{html.Class(classes)}
		attrs = append(attrs, props.Attrs...)

		return html.Ol(
			g.Group(attrs),
			g.Group(children),
		)
	}
}

