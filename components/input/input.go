// Package input provides input components for forms.
package input

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Type        string // text, email, password, etc.
	Name        string
	ID          string
	Placeholder string
	Value       string
	Required    bool
	Disabled    bool
	Invalid     bool // Sets aria-invalid for error styling
	Variant     forgeui.Variant
	Class       string
	Attrs       []g.Node
}

type Option func(*Props)

func WithType(t string) Option {
	return func(p *Props) { p.Type = t }
}

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

func Required() Option {
	return func(p *Props) { p.Required = true }
}

func Disabled() Option {
	return func(p *Props) { p.Disabled = true }
}

// Invalid sets aria-invalid="true" for error state styling
func Invalid() Option {
	return func(p *Props) { p.Invalid = true }
}

func WithVariant(v forgeui.Variant) Option {
	return func(p *Props) { p.Variant = v }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Input creates an input field following shadcn/ui patterns
//
// Example:
//
//	input.Input(
//	    input.WithPlaceholder("Enter email"),
//	    input.WithType("email"),
//	    input.WithName("email"),
//	)
//
// With error state:
//
//	input.Input(
//	    input.WithPlaceholder("Enter email"),
//	    input.Invalid(), // Shows error styling via aria-invalid
//	)
func Input(opts ...Option) g.Node {
	props := &Props{
		Type:    "text",
		Variant: forgeui.VariantDefault,
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := inputCVA.Classes(map[string]string{
		"variant": string(props.Variant),
	})

	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
		html.Type(props.Type),
		g.Attr("data-slot", "input"),
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Placeholder != "" {
		attrs = append(attrs, html.Placeholder(props.Placeholder))
	}

	if props.Value != "" {
		attrs = append(attrs, html.Value(props.Value))
	}

	if props.Required {
		attrs = append(attrs, g.Attr("required", ""))
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("disabled", ""))
	}

	if props.Invalid {
		attrs = append(attrs, g.Attr("aria-invalid", "true"))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Input(g.Group(attrs))
}
