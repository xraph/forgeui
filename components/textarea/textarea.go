// Package textarea provides textarea components for multi-line text input.
package textarea

import (
	"strconv"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Name        string
	ID          string
	Placeholder string
	Value       string
	Rows        int
	Required    bool
	Disabled    bool
	Class       string
	Attrs       []g.Node
}

type Option func(*Props)

func WithName(name string) Option {
	return func(p *Props) { p.Name = name }
}

func WithID(id string) Option {
	return func(p *Props) { p.ID = id }
}

func WithPlaceholder(placeholder string) Option {
	return func(p *Props) { p.Placeholder = placeholder }
}

func WithValue(value string) Option {
	return func(p *Props) { p.Value = value }
}

func WithRows(rows int) Option {
	return func(p *Props) { p.Rows = rows }
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

// Textarea creates a textarea field
func Textarea(opts ...Option) g.Node {
	props := &Props{
		Rows: 3,
	}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := forgeui.CN(
		"placeholder:text-muted-foreground",
		"selection:bg-primary",
		"selection:text-primary-foreground",
		"border-input",
		"min-h-[80px]",
		"w-full",
		"min-w-0",
		"rounded-md",
		"border",
		"bg-background",
		"px-3",
		"py-2",
		"text-sm",
		"shadow-sm",
		"transition-colors",
		"outline-none",
		"resize-y",
		"ring-offset-background",
		"focus-visible:outline-none",
		"focus-visible:ring-2",
		"focus-visible:ring-ring",
		"focus-visible:ring-offset-2",
		"disabled:cursor-not-allowed",
		"disabled:opacity-50",
	)

	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{html.Class(classes)}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Placeholder != "" {
		attrs = append(attrs, html.Placeholder(props.Placeholder))
	}

	if props.Rows > 0 {
		attrs = append(attrs, html.Rows(strconv.Itoa(props.Rows)))
	}

	if props.Required {
		attrs = append(attrs, g.Attr("required", ""))
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("disabled", ""))
	}

	attrs = append(attrs, props.Attrs...)

	children := []g.Node{}
	if props.Value != "" {
		children = append(children, g.Text(props.Value))
	}

	return html.Textarea(
		g.Group(attrs),
		g.Group(children),
	)
}
