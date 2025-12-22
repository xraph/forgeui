// Package checkbox provides checkbox components.
package checkbox

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Name     string
	ID       string
	Value    string
	Checked  bool
	Required bool
	Disabled bool
	Class    string
	Attrs    []g.Node
}

type Option func(*Props)

func WithName(name string) Option {
	return func(p *Props) { p.Name = name }
}

func WithID(id string) Option {
	return func(p *Props) { p.ID = id }
}

func WithValue(value string) Option {
	return func(p *Props) { p.Value = value }
}

func Checked() Option {
	return func(p *Props) { p.Checked = true }
}

func Required() Option {
	return func(p *Props) { p.Required = true }
}

func Disabled() Option {
	return func(p *Props) { p.Disabled = true }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Checkbox creates a checkbox input
func Checkbox(opts ...Option) g.Node {
	props := &Props{}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := forgeui.CN(
		"peer",
		"border-input",
		"bg-background",
		"data-[state=checked]:bg-primary",
		"data-[state=checked]:text-primary-foreground",
		"data-[state=checked]:border-primary",
		"size-4",
		"shrink-0",
		"rounded-sm",
		"border",
		"shadow-sm",
		"transition-colors",
		"outline-none",
		"ring-offset-background",
		"focus-visible:outline-none",
		"focus-visible:ring-2",
		"focus-visible:ring-ring",
		"focus-visible:ring-offset-2",
		"disabled:cursor-not-allowed",
		"disabled:opacity-50",
	)

	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{
		html.Type("checkbox"),
		html.Class(classes),
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Value != "" {
		attrs = append(attrs, html.Value(props.Value))
	}

	if props.Checked {
		attrs = append(attrs, g.Attr("checked", ""))
	}

	if props.Required {
		attrs = append(attrs, g.Attr("required", ""))
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("disabled", ""))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Input(g.Group(attrs))
}
