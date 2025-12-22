package radio

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

// GroupProps defines radio group configuration
type GroupProps struct {
	Name     string
	Value    string
	Disabled bool
	Class    string
	Attrs    []g.Node
}

// GroupOption is a functional option for configuring radio groups
type GroupOption func(*GroupProps)

// WithGroupName sets the group name (required for radio buttons)
func WithGroupName(name string) GroupOption {
	return func(p *GroupProps) { p.Name = name }
}

// WithGroupValue sets the selected value
func WithGroupValue(value string) GroupOption {
	return func(p *GroupProps) { p.Value = value }
}

// WithGroupDisabled disables all radio buttons in the group
func WithGroupDisabled() GroupOption {
	return func(p *GroupProps) { p.Disabled = true }
}

// WithGroupClass adds custom classes
func WithGroupClass(class string) GroupOption {
	return func(p *GroupProps) { p.Class = class }
}

// WithGroupAttrs adds custom attributes
func WithGroupAttrs(attrs ...g.Node) GroupOption {
	return func(p *GroupProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Group creates a radio button group container
//
// Example:
//
//	group := radio.Group(
//	    radio.WithGroupName("size"),
//	    radio.WithGroupValue("medium"),
//	    radio.Radio("small", "Small", radio.WithValue("small")),
//	    radio.Radio("medium", "Medium", radio.WithValue("medium")),
//	    radio.Radio("large", "Large", radio.WithValue("large")),
//	)
func Group(opts []GroupOption, children ...g.Node) g.Node {
	props := &GroupProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN("space-y-2", props.Class)

	attrs := []g.Node{
		html.Class(classes),
		html.Role("radiogroup"),
	}

	if props.Name != "" {
		attrs = append(attrs, g.Attr("data-name", props.Name))
	}

	if props.Value != "" {
		attrs = append(attrs, g.Attr("data-value", props.Value))
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("aria-disabled", "true"))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// GroupItem creates a radio button within a group with label
//
// Example:
//
//	item := radio.GroupItem(
//	    "size-small",
//	    "Small",
//	    radio.WithValue("small"),
//	    radio.WithName("size"),
//	)
func GroupItem(id, labelText string, opts ...Option) g.Node {
	// Add ID to options
	allOpts := append([]Option{WithID(id)}, opts...)
	radioBtn := Radio(allOpts...)

	return html.Div(
		html.Class("flex items-center space-x-2"),
		radioBtn,
		html.Label(
			html.For(id),
			html.Class("text-sm cursor-pointer"),
			g.Text(labelText),
		),
	)
}

