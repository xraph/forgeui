// Package skeleton provides skeleton loading placeholders.
package skeleton

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Width  string
	Height string
	Class  string
	Attrs  []g.Node
}

type Option func(*Props)

func WithWidth(w string) Option {
	return func(p *Props) { p.Width = w }
}

func WithHeight(h string) Option {
	return func(p *Props) { p.Height = h }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Skeleton creates a skeleton loading placeholder
func Skeleton(opts ...Option) g.Node {
	props := &Props{}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := "animate-pulse rounded-md bg-muted"

	classes := forgeui.CN(baseClasses, props.Width, props.Height, props.Class)

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.Div(g.Group(attrs))
}
