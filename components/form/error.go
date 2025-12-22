package form

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

// ErrorProps defines error message configuration
type ErrorProps struct {
	Class string
	Attrs []g.Node
}

// ErrorOption is a functional option for configuring error messages
type ErrorOption func(*ErrorProps)

// WithErrorClass adds custom classes
func WithErrorClass(class string) ErrorOption {
	return func(p *ErrorProps) { p.Class = class }
}

// WithErrorAttrs adds custom attributes
func WithErrorAttrs(attrs ...g.Node) ErrorOption {
	return func(p *ErrorProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Error creates a standalone error message component
//
// Example:
//
//	err := form.Error(
//	    "Invalid email address",
//	    form.WithErrorClass("mt-2"),
//	)
func Error(text string, opts ...ErrorOption) g.Node {
	if text == "" {
		return nil
	}

	props := &ErrorProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		"text-sm",
		"font-medium",
		"text-destructive",
		props.Class,
	)

	attrs := []g.Node{
		html.Class(classes),
		html.Role("alert"),
		g.Attr("aria-live", "polite"),
	}

	attrs = append(attrs, props.Attrs...)

	return html.P(
		g.Group(attrs),
		g.Text(text),
	)
}
