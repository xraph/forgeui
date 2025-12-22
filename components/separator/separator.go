// Package separator provides horizontal and vertical separators.
package separator

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Orientation string // horizontal or vertical
	Class       string
	Attrs       []g.Node
}

type Option func(*Props)

func Vertical() Option {
	return func(p *Props) { p.Orientation = "vertical" }
}

func Horizontal() Option {
	return func(p *Props) { p.Orientation = "horizontal" }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Separator creates a separator line
func Separator(opts ...Option) g.Node {
	props := &Props{
		Orientation: "horizontal",
	}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := "shrink-0 bg-border"
	var orientationClass string

	if props.Orientation == "vertical" {
		orientationClass = "h-full w-[1px]"
	} else {
		orientationClass = "h-[1px] w-full"
	}

	classes := forgeui.CN(baseClasses, orientationClass, props.Class)

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("role", "separator"),
		g.Attr("aria-orientation", props.Orientation),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(g.Group(attrs))
}
