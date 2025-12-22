package primitives

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

// TextProps defines properties for the Text component
type TextProps struct {
	As       string // HTML tag (p, span, div, h1-h6)
	Size     string // text size class
	Weight   string // font weight class
	Color    string // text color class
	Align    string // text alignment
	Class    string
	Children []g.Node
	Attrs    []g.Node
}

// TextOption is a functional option for configuring Text
type TextOption func(*TextProps)

// TextAs sets the HTML element type
func TextAs(tag string) TextOption {
	return func(p *TextProps) { p.As = tag }
}

// TextSize sets the text size
func TextSize(size string) TextOption {
	return func(p *TextProps) { p.Size = size }
}

// TextWeight sets the font weight
func TextWeight(weight string) TextOption {
	return func(p *TextProps) { p.Weight = weight }
}

// TextColor sets the text color
func TextColor(color string) TextOption {
	return func(p *TextProps) { p.Color = color }
}

// TextAlign sets the text alignment
func TextAlign(align string) TextOption {
	return func(p *TextProps) { p.Align = align }
}

// TextClass adds custom classes
func TextClass(class string) TextOption {
	return func(p *TextProps) { p.Class = class }
}

// TextChildren adds child nodes
func TextChildren(children ...g.Node) TextOption {
	return func(p *TextProps) { p.Children = append(p.Children, children...) }
}

// TextAttrs adds custom attributes
func TextAttrs(attrs ...g.Node) TextOption {
	return func(p *TextProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Text creates a typography primitive
func Text(opts ...TextOption) g.Node {
	props := &TextProps{
		As: "p",
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		props.Size,
		props.Weight,
		props.Color,
		props.Align,
		props.Class,
	)

	attrs := []g.Node{}
	if classes != "" {
		attrs = append(attrs, html.Class(classes))
	}
	attrs = append(attrs, props.Attrs...)

	return g.El(props.As, g.Group(attrs), g.Group(props.Children))
}
