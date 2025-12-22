package form

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

// DescriptionProps defines description configuration
type DescriptionProps struct {
	Class string
	Attrs []g.Node
}

// DescriptionOption is a functional option for configuring descriptions
type DescriptionOption func(*DescriptionProps)

// WithDescriptionClass adds custom classes
func WithDescriptionClass(class string) DescriptionOption {
	return func(p *DescriptionProps) { p.Class = class }
}

// WithDescriptionAttrs adds custom attributes
func WithDescriptionAttrs(attrs ...g.Node) DescriptionOption {
	return func(p *DescriptionProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Description creates a standalone helper text component
//
// Example:
//
//	desc := form.Description(
//	    "Enter your email address for password reset",
//	    form.WithDescriptionClass("mt-2"),
//	)
func Description(text string, opts ...DescriptionOption) g.Node {
	if text == "" {
		return nil
	}

	props := &DescriptionProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		"text-sm",
		"text-muted-foreground",
		props.Class,
	)

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.P(
		g.Group(attrs),
		g.Text(text),
	)
}

