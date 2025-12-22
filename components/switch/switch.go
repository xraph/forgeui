// Package switch provides toggle switch components.
package switchc

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Name     string
	ID       string
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

// Switch creates a toggle switch (styled checkbox)
func Switch(opts ...Option) g.Node {
	props := &Props{}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := forgeui.CN(
		"peer",
		"inline-flex",
		"h-6",
		"w-11",
		"shrink-0",
		"cursor-pointer",
		"items-center",
		"rounded-full",
		"border-2",
		"border-transparent",
		"transition-colors",
		"shadow-sm",
		"outline-none",
		"ring-offset-background",
		"focus-visible:outline-none",
		"focus-visible:ring-2",
		"focus-visible:ring-ring",
		"focus-visible:ring-offset-2",
		"disabled:cursor-not-allowed",
		"disabled:opacity-50",
		"bg-input",
		"checked:bg-primary",
		"data-[state=checked]:bg-primary",
	)

	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{
		html.Type("checkbox"),
		html.Class(classes),
		g.Attr("role", "switch"),
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
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
