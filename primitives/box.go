package primitives

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

// BoxProps defines properties for the Box component
type BoxProps struct {
	As       string // HTML tag (div, span, section, article, etc.)
	Class    string
	P        string // padding
	M        string // margin
	Bg       string // background
	Rounded  string // border-radius
	Shadow   string // box-shadow
	W        string // width
	H        string // height
	Children []g.Node
	Attrs    []g.Node
}

// BoxOption is a functional option for configuring Box
type BoxOption func(*BoxProps)

// WithAs sets the HTML element type
func WithAs(tag string) BoxOption {
	return func(p *BoxProps) { p.As = tag }
}

// WithClass adds custom classes
func WithClass(class string) BoxOption {
	return func(p *BoxProps) { p.Class = class }
}

// WithPadding sets padding classes
func WithPadding(p string) BoxOption {
	return func(props *BoxProps) { props.P = p }
}

// WithMargin sets margin classes
func WithMargin(m string) BoxOption {
	return func(props *BoxProps) { props.M = m }
}

// WithBackground sets background classes
func WithBackground(bg string) BoxOption {
	return func(p *BoxProps) { p.Bg = bg }
}

// WithRounded sets border-radius classes
func WithRounded(rounded string) BoxOption {
	return func(p *BoxProps) { p.Rounded = rounded }
}

// WithShadow sets box-shadow classes
func WithShadow(shadow string) BoxOption {
	return func(p *BoxProps) { p.Shadow = shadow }
}

// WithWidth sets width classes
func WithWidth(w string) BoxOption {
	return func(p *BoxProps) { p.W = w }
}

// WithHeight sets height classes
func WithHeight(h string) BoxOption {
	return func(p *BoxProps) { p.H = h }
}

// WithChildren adds child nodes
func WithChildren(children ...g.Node) BoxOption {
	return func(p *BoxProps) { p.Children = append(p.Children, children...) }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs ...g.Node) BoxOption {
	return func(p *BoxProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Box creates a polymorphic container element
// It's the most basic primitive for layout
func Box(opts ...BoxOption) g.Node {
	props := &BoxProps{
		As: "div",
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		props.P,
		props.M,
		props.Bg,
		props.Rounded,
		props.Shadow,
		props.W,
		props.H,
		props.Class,
	)

	attrs := []g.Node{}
	if classes != "" {
		attrs = append(attrs, html.Class(classes))
	}

	attrs = append(attrs, props.Attrs...)

	return g.El(props.As, g.Group(attrs), g.Group(props.Children))
}
