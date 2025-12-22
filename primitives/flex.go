package primitives

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

// FlexProps defines properties for the Flex component
type FlexProps struct {
	Direction string // row, col, row-reverse, col-reverse
	Wrap      string // wrap, nowrap, wrap-reverse
	Justify   string // start, end, center, between, around, evenly
	Align     string // start, end, center, stretch, baseline
	Gap       string // gap size
	Class     string
	Children  []g.Node
	Attrs     []g.Node
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

// FlexChildren adds child nodes
func FlexChildren(children ...g.Node) FlexOption {
	return func(p *FlexProps) { p.Children = append(p.Children, children...) }
}

// FlexAttrs adds custom attributes
func FlexAttrs(attrs ...g.Node) FlexOption {
	return func(p *FlexProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Flex creates a flexbox container
func Flex(opts ...FlexOption) g.Node {
	props := &FlexProps{
		Direction: "row",
		Wrap:      "nowrap",
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := []string{"flex"}

	// Direction
	switch props.Direction {
	case "col", "column":
		classes = append(classes, "flex-col")
	case "row-reverse":
		classes = append(classes, "flex-row-reverse")
	case "col-reverse", "column-reverse":
		classes = append(classes, "flex-col-reverse")
		// "row" is default, no extra class needed
	}

	// Wrap
	switch props.Wrap {
	case "wrap":
		classes = append(classes, "flex-wrap")
	case "wrap-reverse":
		classes = append(classes, "flex-wrap-reverse")
	case "nowrap":
		classes = append(classes, "flex-nowrap")
	}

	// Justify
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

	// Align
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

	// Gap
	if props.Gap != "" {
		classes = append(classes, "gap-"+props.Gap)
	}

	// Custom class
	if props.Class != "" {
		classes = append(classes, props.Class)
	}

	attrs := []g.Node{html.Class(forgeui.CN(classes...))}
	attrs = append(attrs, props.Attrs...)

	return html.Div(g.Group(attrs), g.Group(props.Children))
}
