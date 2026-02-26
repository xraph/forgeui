package primitives

import (
	"context"
	"io"

	"github.com/a-h/templ"

	"github.com/xraph/forgeui"
)

// BoxProps defines properties for the Box component
type BoxProps struct {
	As         string           // HTML tag (div, span, section, article, etc.)
	Class      string           // custom classes
	P          string           // padding
	M          string           // margin
	Bg         string           // background
	Rounded    string           // border-radius
	Shadow     string           // box-shadow
	W          string           // width
	H          string           // height
	Children   []templ.Component
	Attributes templ.Attributes
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
func WithPadding(padding string) BoxOption {
	return func(props *BoxProps) { props.P = padding }
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

// WithChildren adds child components
func WithChildren(children ...templ.Component) BoxOption {
	return func(p *BoxProps) { p.Children = append(p.Children, children...) }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs templ.Attributes) BoxOption {
	return func(p *BoxProps) {
		if p.Attributes == nil {
			p.Attributes = templ.Attributes{}
		}
		for k, v := range attrs {
			p.Attributes[k] = v
		}
	}
}

// Box creates a polymorphic container element.
// It's the most basic primitive for layout.
func Box(opts ...BoxOption) templ.Component {
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

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := writeOpenTag(w, props.As, classes, props.Attributes); err != nil {
			return err
		}
		if err := renderChildren(ctx, w, props.Children); err != nil {
			return err
		}
		return writeCloseTag(w, props.As)
	})
}
