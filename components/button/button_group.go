package button

import (
	g "maragu.dev/gomponents"

	"github.com/xraph/forgeui/primitives"
)

// GroupProps defines button group configuration
type GroupProps struct {
	Gap   string
	Class string
}

// GroupOption is a functional option for configuring button groups
type GroupOption func(*GroupProps)

// WithGap sets the gap between buttons
func WithGap(gap string) GroupOption {
	return func(p *GroupProps) { p.Gap = gap }
}

// WithGroupClass adds custom classes to the button group
func WithGroupClass(class string) GroupOption {
	return func(p *GroupProps) { p.Class = class }
}

// Group creates a button group container
// Useful for grouping related buttons together
//
// Example:
//
//	btnGroup := button.Group(
//	    button.WithGap("2"),
//	    button.Primary(g.Text("Save")),
//	    button.Secondary(g.Text("Cancel")),
//	)
func Group(opts []GroupOption, children ...g.Node) g.Node {
	props := &GroupProps{
		Gap: "2",
	}

	for _, opt := range opts {
		opt(props)
	}

	flexOpts := []primitives.FlexOption{
		primitives.FlexDirection("row"),
		primitives.FlexAlign("center"),
		primitives.FlexGap(props.Gap),
		primitives.FlexChildren(children...),
	}

	if props.Class != "" {
		flexOpts = append(flexOpts, primitives.FlexClass(props.Class))
	}

	return primitives.Flex(flexOpts...)
}
