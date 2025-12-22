// Package radio provides radio button components.
package radio

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/primitives"
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

// Radio creates a radio button input
func Radio(opts ...Option) g.Node {
	props := &Props{}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := forgeui.CN(
		"aspect-square",
		"size-4",
		"rounded-full",
		"border",
		"border-input",
		"text-primary",
		"dark:bg-input/30",
		"data-[state=checked]:bg-primary",
		"data-[state=checked]:text-primary-foreground",
		"dark:data-[state=checked]:bg-primary",
		"data-[state=checked]:border-primary",
		"focus-visible:border-ring",
		"focus-visible:ring-ring/50",
		"aria-invalid:ring-destructive/20",
		"dark:aria-invalid:ring-destructive/40",
		"aria-invalid:border-destructive",
		"shadow-xs",
		"transition-shadow",
		"outline-none",
		"focus-visible:ring-[3px]",
		"disabled:cursor-not-allowed",
		"disabled:opacity-50",
	)

	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{
		html.Type("radio"),
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

// RadioGroup creates a group of radio buttons
func RadioGroup(name string, options []RadioGroupOption) g.Node {
	items := make([]g.Node, 0, len(options))

	for _, opt := range options {
		items = append(items, primitives.Flex(
			primitives.FlexAlign("center"),
			primitives.FlexGap("2"),
			primitives.FlexChildren(
				func() g.Node {
					opts := []Option{
						WithName(name),
						WithID(opt.ID),
						WithValue(opt.Value),
					}
					if opt.Checked {
						opts = append(opts, Checked())
					}
					return Radio(opts...)
				}(),
				primitives.Text(
					primitives.TextAs("label"),
					primitives.TextSize("text-sm"),
					primitives.TextClass("cursor-pointer"),
					primitives.TextAttrs(html.For(opt.ID)),
					primitives.TextChildren(g.Text(opt.Label)),
				),
			),
		))
	}

	return primitives.VStack("2", items...)
}

// RadioGroupOption defines a radio button option in a group
type RadioGroupOption struct {
	ID      string
	Value   string
	Label   string
	Checked bool
}
