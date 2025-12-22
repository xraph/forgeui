// Package form provides form wrapper and validation helpers.
package form

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Action string
	Method string
	Class  string
	Attrs  []g.Node
}

type Option func(*Props)

func WithAction(action string) Option {
	return func(p *Props) { p.Action = action }
}

func WithMethod(method string) Option {
	return func(p *Props) { p.Method = method }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Form creates a form wrapper
func Form(opts []Option, children ...g.Node) g.Node {
	props := &Props{
		Method: "POST",
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN("space-y-6", props.Class)

	attrs := []g.Node{html.Class(classes)}

	if props.Action != "" {
		attrs = append(attrs, html.Action(props.Action))
	}

	if props.Method != "" {
		attrs = append(attrs, html.Method(props.Method))
	}

	attrs = append(attrs, props.Attrs...)

	return html.FormEl(
		g.Group(attrs),
		g.Group(children),
	)
}
