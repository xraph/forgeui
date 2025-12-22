// Package popover provides floating content anchored to a trigger element.
//
// Popovers display rich content in a floating panel positioned relative
// to a trigger element. Unlike tooltips, popovers can contain interactive
// content and are opened by clicking rather than hovering.
//
// Basic usage:
//
//	popover.Popover(
//	    popover.PopoverProps{
//	        Position: forgeui.PositionBottom,
//	        Align:    forgeui.AlignStart,
//	    },
//	    button.Button(g.Text("Open Popover")),
//	    html.Div(g.Text("Popover content")),
//	)
package popover

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// PopoverProps defines configuration for a popover
type PopoverProps struct {
	// Position specifies where the popover appears relative to trigger
	Position forgeui.Position

	// Align specifies how the popover aligns with the trigger
	Align forgeui.Align

	// ShowArrow displays a pointer arrow
	ShowArrow bool

	// Offset is the distance in pixels from the trigger
	Offset int

	// CloseOnOutsideClick enables closing when clicking outside
	CloseOnOutsideClick bool

	// Class adds additional CSS classes
	Class string
}

// defaultPopoverProps returns default popover properties
func defaultPopoverProps() PopoverProps {
	return PopoverProps{
		Position:            forgeui.PositionBottom,
		Align:               forgeui.AlignCenter,
		ShowArrow:           false,
		Offset:              8,
		CloseOnOutsideClick: true,
	}
}

// Popover creates a popover with floating content anchored to a trigger.
//
// The popover uses Alpine.js for state management and can be positioned
// on any side of the trigger element (top, right, bottom, left) with
// different alignments (start, center, end).
//
// Example:
//
//	popover.Popover(
//	    popover.PopoverProps{
//	        Position: forgeui.PositionRight,
//	        Align:    forgeui.AlignStart,
//	    },
//	    button.Button(g.Text("Settings")),
//	    html.Div(
//	        html.H4(g.Text("Quick Settings")),
//	        html.P(g.Text("Configure your options here")),
//	    ),
//	)
func Popover(props PopoverProps, trigger g.Node, content ...g.Node) g.Node {
	// Apply defaults
	if props.Position == "" {
		props.Position = forgeui.PositionBottom
	}
	if props.Align == "" {
		props.Align = forgeui.AlignCenter
	}
	if props.Offset == 0 {
		props.Offset = 8
	}

	// Get positioning classes
	positionClass := getPopoverPositionClass(props.Position, props.Align)
	slideTransition := getPopoverTransition(props.Position)

	// Alpine state
	alpineData := alpine.XData(map[string]any{
		"open": false,
	})

	// Click outside to close
	var clickOutside g.Node
	if props.CloseOnOutsideClick {
		clickOutside = alpine.XOn("click.outside", "open = false")
	}

	return html.Div(
		alpineData,
		html.Class("relative inline-block"),

		// Trigger
		html.Div(
			alpine.XOn("click", "open = !open"),
			trigger,
		),

		// Popover content
		html.Div(
			alpine.XShow("open"),
			alpine.XCloak(),
			g.Group(alpine.XTransition(slideTransition)),
			g.If(props.CloseOnOutsideClick, clickOutside),
			html.Class(fmt.Sprintf(
				"absolute z-50 %s bg-popover text-popover-foreground border border-border rounded-md shadow-md p-4 min-w-[200px] %s",
				positionClass,
				props.Class,
			)),
			g.Attr("role", "dialog"),
			g.Attr("aria-modal", "false"),

			// Arrow (if enabled)
			g.If(props.ShowArrow, popoverArrow(props.Position)),

			// Content
			g.Group(content),
		),
	)
}

// PopoverClose wraps an element to close the popover when clicked.
func PopoverClose(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = false"),
		child,
	)
}

// popoverArrow creates the arrow pointer
func popoverArrow(position forgeui.Position) g.Node {
	arrowClass := getArrowPositionClass(position)

	return html.Div(
		html.Class(fmt.Sprintf(
			"absolute w-2 h-2 bg-popover border-border rotate-45 %s",
			arrowClass,
		)),
	)
}

// getPopoverPositionClass returns positioning classes based on position and align
func getPopoverPositionClass(position forgeui.Position, align forgeui.Align) string {
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

// getArrowPositionClass returns arrow positioning classes
func getArrowPositionClass(position forgeui.Position) string {
	switch position {
	case forgeui.PositionTop:
		return "top-full left-1/2 -translate-x-1/2 -mt-1 border-b border-r"
	case forgeui.PositionRight:
		return "right-full top-1/2 -translate-y-1/2 -mr-1 border-l border-t"
	case forgeui.PositionBottom:
		return "bottom-full left-1/2 -translate-x-1/2 -mb-1 border-t border-l"
	case forgeui.PositionLeft:
		return "left-full top-1/2 -translate-y-1/2 -ml-1 border-r border-b"
	default:
		return "bottom-full left-1/2 -translate-x-1/2 -mb-1 border-t border-l"
	}
}

// getPopoverTransition returns the appropriate transition based on position
func getPopoverTransition(position forgeui.Position) *animation.Transition {
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

