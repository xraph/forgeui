// Package button provides button components with multiple variants,
// sizes, and states following shadcn/ui design patterns.
package button

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

// Props defines button component configuration
type Props struct {
	Variant  forgeui.Variant
	Size     forgeui.Size
	Type     string // button, submit, reset
	Disabled bool
	Loading  bool
	Class    string
	Attrs    []g.Node
}

// Option is a functional option for configuring the button
type Option func(*Props)

// WithVariant sets the visual variant
func WithVariant(v forgeui.Variant) Option {
	return func(p *Props) { p.Variant = v }
}

// WithSize sets the button size
func WithSize(s forgeui.Size) Option {
	return func(p *Props) { p.Size = s }
}

// WithType sets the button type attribute
func WithType(t string) Option {
	return func(p *Props) { p.Type = t }
}

// WithClass adds custom classes
func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

// Disabled sets the disabled state
func Disabled() Option {
	return func(p *Props) { p.Disabled = true }
}

// Loading sets the loading state (also disables the button)
func Loading() Option {
	return func(p *Props) { p.Loading = true }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultProps returns default button properties
func defaultProps() *Props {
	return &Props{
		Variant: forgeui.VariantDefault,
		Size:    forgeui.SizeMD,
		Type:    "button",
	}
}

// Button creates a button component
//
// Example:
//
//	btn := button.Button(
//	    g.Text("Click me"),
//	    button.WithVariant(forgeui.VariantPrimary),
//	    button.WithSize(forgeui.SizeLG),
//	)
func Button(children g.Node, opts ...Option) g.Node {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := buttonCVA.Classes(map[string]string{
		"variant": string(props.Variant),
		"size":    string(props.Size),
	})

	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
		html.Type(props.Type),
	}

	if props.Disabled || props.Loading {
		attrs = append(attrs, g.Attr("disabled", ""))
	}

	if props.Loading {
		attrs = append(attrs, g.Attr("aria-busy", "true"))
	}

	attrs = append(attrs, props.Attrs...)

	// Add spinner if loading
	content := children
	if props.Loading {
		content = g.Group([]g.Node{
			loadingSpinner(),
			children,
		})
	}

	return html.Button(
		g.Group(attrs),
		content,
	)
}

// loadingSpinner creates a simple SVG spinner for loading state
func loadingSpinner() g.Node {
	return g.El("svg",
		html.Class("animate-spin -ml-1 mr-2 h-4 w-4"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("fill", "none"),
		g.Attr("viewBox", "0 0 24 24"),
		g.El("circle",
			html.Class("opacity-25"),
			g.Attr("cx", "12"),
			g.Attr("cy", "12"),
			g.Attr("r", "10"),
			g.Attr("stroke", "currentColor"),
			g.Attr("stroke-width", "4"),
		),
		g.El("path",
			html.Class("opacity-75"),
			g.Attr("fill", "currentColor"),
			g.Attr("d", "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"),
		),
	)
}

// Convenience constructors for common variants

// Primary creates a primary button (alias for default variant)
func Primary(children g.Node, opts ...Option) g.Node {
	return Button(children, append([]Option{WithVariant(forgeui.VariantDefault)}, opts...)...)
}

// Secondary creates a secondary button
func Secondary(children g.Node, opts ...Option) g.Node {
	return Button(children, append([]Option{WithVariant(forgeui.VariantSecondary)}, opts...)...)
}

// Destructive creates a destructive/danger button
func Destructive(children g.Node, opts ...Option) g.Node {
	return Button(children, append([]Option{WithVariant(forgeui.VariantDestructive)}, opts...)...)
}

// Outline creates an outline button
func Outline(children g.Node, opts ...Option) g.Node {
	return Button(children, append([]Option{WithVariant(forgeui.VariantOutline)}, opts...)...)
}

// Ghost creates a ghost button
func Ghost(children g.Node, opts ...Option) g.Node {
	return Button(children, append([]Option{WithVariant(forgeui.VariantGhost)}, opts...)...)
}

// Link creates a link-styled button
func Link(children g.Node, opts ...Option) g.Node {
	return Button(children, append([]Option{WithVariant(forgeui.VariantLink)}, opts...)...)
}

// IconButton creates an icon-only button
func IconButton(children g.Node, opts ...Option) g.Node {
	return Button(children, append([]Option{WithSize(forgeui.SizeIcon)}, opts...)...)
}
