// Package label provides label components for form fields.
package label

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	For   string
	Class string
	Attrs []g.Node
}

type Option func(*Props)

func WithFor(htmlFor string) Option {
	return func(p *Props) { p.For = htmlFor }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Label creates a form label
func Label(text string, opts ...Option) g.Node {
	props := &Props{}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{html.Class(classes)}

	if props.For != "" {
		attrs = append(attrs, html.For(props.For))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Label(
		g.Group(attrs),
		g.Text(text),
	)
}
