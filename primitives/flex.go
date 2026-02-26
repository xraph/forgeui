package primitives

import (
	"context"
	"io"

	"github.com/a-h/templ"

	"github.com/xraph/forgeui"
)

// FlexProps defines properties for the Flex component
type FlexProps struct {
	Direction  string // row, col, row-reverse, col-reverse
	Wrap       string // wrap, nowrap, wrap-reverse
	Justify    string // start, end, center, between, around, evenly
	Align      string // start, end, center, stretch, baseline
	Gap        string // gap size
	Class      string
	Children   []templ.Component
	Attributes templ.Attributes
}

// FlexOption is a functional option for configuring Flex
type FlexOption func(*FlexProps)

// FlexDirection sets the flex direction
func FlexDirection(direction string) FlexOption {
	return func(p *FlexProps) { p.Direction = direction }
}

// FlexWrap sets the flex wrap
func FlexWrap(wrap string) FlexOption {
	return func(p *FlexProps) { p.Wrap = wrap }
}

// FlexJustify sets the justify-content
func FlexJustify(justify string) FlexOption {
	return func(p *FlexProps) { p.Justify = justify }
}

// FlexAlign sets the align-items
func FlexAlign(align string) FlexOption {
	return func(p *FlexProps) { p.Align = align }
}

// FlexGap sets the gap
func FlexGap(gap string) FlexOption {
	return func(p *FlexProps) { p.Gap = gap }
}

// FlexClass adds custom classes
func FlexClass(class string) FlexOption {
	return func(p *FlexProps) { p.Class = class }
}

// FlexChildren adds child components
func FlexChildren(children ...templ.Component) FlexOption {
	return func(p *FlexProps) { p.Children = append(p.Children, children...) }
}

// FlexAttrs adds custom attributes
func FlexAttrs(attrs templ.Attributes) FlexOption {
	return func(p *FlexProps) {
		if p.Attributes == nil {
			p.Attributes = templ.Attributes{}
		}
		for k, v := range attrs {
			p.Attributes[k] = v
		}
	}
}

// flexClasses computes the CSS classes for a Flex component.
func flexClasses(props *FlexProps) string {
	classes := []string{"flex"}

	switch props.Direction {
	case "col", "column":
		classes = append(classes, "flex-col")
	case "row-reverse":
		classes = append(classes, "flex-row-reverse")
	case "col-reverse", "column-reverse":
		classes = append(classes, "flex-col-reverse")
	}

	switch props.Wrap {
	case "wrap":
		classes = append(classes, "flex-wrap")
	case "wrap-reverse":
		classes = append(classes, "flex-wrap-reverse")
	case "nowrap":
		classes = append(classes, "flex-nowrap")
	}

	switch props.Justify {
	case "start":
		classes = append(classes, "justify-start")
	case "end":
		classes = append(classes, "justify-end")
	case "center":
		classes = append(classes, "justify-center")
	case "between":
		classes = append(classes, "justify-between")
	case "around":
		classes = append(classes, "justify-around")
	case "evenly":
		classes = append(classes, "justify-evenly")
	}

	switch props.Align {
	case "start":
		classes = append(classes, "items-start")
	case "end":
		classes = append(classes, "items-end")
	case "center":
		classes = append(classes, "items-center")
	case "stretch":
		classes = append(classes, "items-stretch")
	case "baseline":
		classes = append(classes, "items-baseline")
	}

	if props.Gap != "" {
		classes = append(classes, "gap-"+props.Gap)
	}

	if props.Class != "" {
		classes = append(classes, props.Class)
	}

	return forgeui.CN(classes...)
}

// Flex creates a flexbox container.
func Flex(opts ...FlexOption) templ.Component {
	props := &FlexProps{
		Direction: "row",
		Wrap:      "nowrap",
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := flexClasses(props)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := writeOpenTag(w, "div", classes, props.Attributes); err != nil {
			return err
		}
		if err := renderChildren(ctx, w, props.Children); err != nil {
			return err
		}
		return writeCloseTag(w, "div")
	})
}
