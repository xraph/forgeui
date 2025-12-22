// Package slider provides range slider components.
package slider

import (
	"strconv"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Name     string
	ID       string
	Min      int
	Max      int
	Value    int
	Step     int
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

func WithMin(min int) Option {
	return func(p *Props) { p.Min = min }
}

func WithMax(max int) Option {
	return func(p *Props) { p.Max = max }
}

func WithValue(value int) Option {
	return func(p *Props) { p.Value = value }
}

func WithStep(step int) Option {
	return func(p *Props) { p.Step = step }
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

// Slider creates a range slider input
func Slider(opts ...Option) g.Node {
	props := &Props{
		Min:  0,
		Max:  100,
		Step: 1,
	}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := forgeui.CN(
		"relative",
		"flex",
		"w-full",
		"touch-none",
		"select-none",
		"items-center",
		"transition-[color,box-shadow]",
		"outline-none",
		"focus-visible:ring-ring/50",
		"focus-visible:ring-[3px]",
		"aria-invalid:ring-destructive/20",
		"dark:aria-invalid:ring-destructive/40",
		"disabled:cursor-not-allowed",
		"disabled:opacity-50",
	)

	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{
		html.Type("range"),
		html.Class(classes),
		g.Attr("min", strconv.Itoa(props.Min)),
		g.Attr("max", strconv.Itoa(props.Max)),
		g.Attr("step", strconv.Itoa(props.Step)),
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Value > 0 {
		attrs = append(attrs, html.Value(strconv.Itoa(props.Value)))
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
