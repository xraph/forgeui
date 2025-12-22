// Package spinner provides loading spinner components.
package spinner

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Size  forgeui.Size
	Class string
	Attrs []g.Node
}

type Option func(*Props)

func WithSize(size forgeui.Size) Option {
	return func(p *Props) { p.Size = size }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Spinner creates a loading spinner
func Spinner(opts ...Option) g.Node {
	props := &Props{
		Size: forgeui.SizeMD,
	}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := "animate-spin rounded-full border-2 border-current border-t-transparent"

	sizeClasses := map[forgeui.Size]string{
		forgeui.SizeSM: "size-4",
		forgeui.SizeMD: "size-6",
		forgeui.SizeLG: "size-8",
		forgeui.SizeXL: "size-12",
	}

	sizeClass := sizeClasses[props.Size]
	if sizeClass == "" {
		sizeClass = sizeClasses[forgeui.SizeMD]
	}

	classes := forgeui.CN(baseClasses, sizeClass, props.Class)

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("role", "status"),
		g.Attr("aria-label", "Loading"),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(g.Group(attrs))
}
