// Package dropdown provides dropdown menu components with keyboard navigation.
//
// Dropdowns display a list of options or actions when clicked. They support
// positioning, alignment, keyboard navigation, and click-outside-to-close.
//
// Basic usage:
//
//	dropdown.Dropdown(
//	    dropdown.DropdownProps{
//	        Position: forgeui.PositionBottom,
//	        Align:    forgeui.AlignStart,
//	    },
//	    button.Button(g.Text("Options")),
//	    html.Div(g.Text("Dropdown content")),
//	)
package dropdown

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// DropdownProps defines configuration for a dropdown
type DropdownProps struct {
	// Position specifies where the dropdown appears relative to trigger
	Position forgeui.Position

	// Align specifies how the dropdown aligns with the trigger
	Align forgeui.Align

	// Class adds additional CSS classes
	Class string
}

// defaultDropdownProps returns default dropdown properties
func defaultDropdownProps() DropdownProps {
	return DropdownProps{
		Position: forgeui.PositionBottom,
		Align:    forgeui.AlignStart,
	}
}

// Dropdown creates a basic dropdown with positioning.
//
// The dropdown toggles open/closed on click and closes when clicking
// outside. It supports positioning and alignment relative to the trigger.
//
// Example:
//
//	dropdown.Dropdown(
//	    dropdown.DropdownProps{
//	        Position: forgeui.PositionBottom,
//	        Align:    forgeui.AlignEnd,
//	    },
//	    button.Button(g.Text("Menu")),
//	    html.Div(
//	        html.A(html.Href("#"), g.Text("Option 1")),
//	        html.A(html.Href("#"), g.Text("Option 2")),
//	    ),
//	)
func Dropdown(props DropdownProps, trigger g.Node, content ...g.Node) g.Node {
	// Apply defaults
	if props.Position == "" {
		props.Position = forgeui.PositionBottom
	}
	if props.Align == "" {
		props.Align = forgeui.AlignStart
	}

	// Get positioning classes
	positionClass := getDropdownPositionClass(props.Position, props.Align)
	slideTransition := getDropdownTransition(props.Position)

	// Alpine state
	alpineData := alpine.XData(map[string]any{
		"open": false,
	})

	return html.Div(
		alpineData,
		alpine.XOn("keydown.escape", "open = false"),
		html.Class("relative inline-block"),

		// Trigger
		html.Div(
			alpine.XOn("click", "open = !open"),
			trigger,
		),

		// Dropdown content
		html.Div(
			alpine.XShow("open"),
			alpine.XCloak(),
			g.Group(alpine.XTransition(slideTransition)),
			alpine.XOn("click.outside", "open = false"),
			html.Class(fmt.Sprintf(
				"absolute z-50 %s bg-popover text-popover-foreground border border-border rounded-md shadow-md p-1 min-w-[200px] %s",
				positionClass,
				props.Class,
			)),
			g.Attr("role", "menu"),
			g.Attr("aria-orientation", "vertical"),

			// Content
			g.Group(content),
		),
	)
}

// DropdownClose wraps an element to close the dropdown when clicked.
func DropdownClose(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = false"),
		child,
	)
}

// getDropdownPositionClass returns positioning classes based on position and align
func getDropdownPositionClass(position forgeui.Position, align forgeui.Align) string {
	var classes string

	switch position {
	case forgeui.PositionTop:
		classes = "bottom-full mb-2"
	case forgeui.PositionRight:
		classes = "left-full ml-2"
	case forgeui.PositionBottom:
		classes = "top-full mt-2"
	case forgeui.PositionLeft:
		classes = "right-full mr-2"
	}

	switch align {
	case forgeui.AlignStart:
		if position == forgeui.PositionTop || position == forgeui.PositionBottom {
			classes += " left-0"
		} else {
			classes += " top-0"
		}
	case forgeui.AlignCenter:
		if position == forgeui.PositionTop || position == forgeui.PositionBottom {
			classes += " left-1/2 -translate-x-1/2"
		} else {
			classes += " top-1/2 -translate-y-1/2"
		}
	case forgeui.AlignEnd:
		if position == forgeui.PositionTop || position == forgeui.PositionBottom {
			classes += " right-0"
		} else {
			classes += " bottom-0"
		}
	}

	return classes
}

// getDropdownTransition returns the appropriate transition based on position
func getDropdownTransition(position forgeui.Position) *animation.Transition {
	switch position {
	case forgeui.PositionTop:
		return animation.SlideUp()
	case forgeui.PositionRight:
		return animation.SlideRight()
	case forgeui.PositionBottom:
		return animation.SlideDown()
	case forgeui.PositionLeft:
		return animation.SlideLeft()
	default:
		return animation.SlideDown()
	}
}

