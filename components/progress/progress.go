// Package progress provides progress bar components.
package progress

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Value int // 0-100
	Class string
	Attrs []g.Node
}

type Option func(*Props)

func WithValue(value int) Option {
	return func(p *Props) {
		if value < 0 {
			value = 0
		}
		if value > 100 {
			value = 100
		}
		p.Value = value
	}
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Progress creates a progress bar
func Progress(opts ...Option) g.Node {
	props := &Props{
		Value: 0,
	}

	for _, opt := range opts {
		opt(props)
	}

	containerClasses := forgeui.CN(
		"relative h-2 w-full overflow-hidden rounded-full bg-secondary shadow-sm",
		props.Class,
	)

	containerAttrs := []g.Node{
		html.Class(containerClasses),
		g.Attr("role", "progressbar"),
		g.Attr("aria-valuemin", "0"),
		g.Attr("aria-valuemax", "100"),
		g.Attr("aria-valuenow", fmt.Sprintf("%d", props.Value)),
	}
	containerAttrs = append(containerAttrs, props.Attrs...)

	indicatorStyle := fmt.Sprintf("width: %d%%", props.Value)

	return html.Div(
		g.Group(containerAttrs),
		html.Div(
			html.Class("h-full w-full flex-1 bg-primary transition-all duration-500 ease-out"),
			html.StyleAttr(indicatorStyle),
		),
	)
}
