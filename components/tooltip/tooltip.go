// Package tooltip provides small informational popups on hover or focus.
//
// Tooltips display brief text information when hovering over or focusing
// an element. They are non-interactive and auto-dismiss when the trigger
// loses hover/focus.
//
// Basic usage:
//
//	tooltip.Tooltip(
//	    tooltip.TooltipProps{Position: forgeui.PositionTop},
//	    button.Button(g.Text("?")),
//	    "This is helpful information",
//	)
package tooltip

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// TooltipProps defines configuration for a tooltip
type TooltipProps struct {
	// Position specifies where the tooltip appears relative to trigger
	Position forgeui.Position

	// Delay is the hover delay before showing (in milliseconds)
	Delay int

	// Class adds additional CSS classes
	Class string
}

// defaultTooltipProps returns default tooltip properties
func defaultTooltipProps() TooltipProps {
	return TooltipProps{
		Position: forgeui.PositionTop,
		Delay:    200,
	}
}

// Tooltip creates a tooltip that shows on hover/focus.
//
// Tooltips are lightweight, non-interactive popups for displaying
// brief information. They automatically show on hover/focus and
// hide when the trigger loses hover/focus.
//
// Example:
//
//	tooltip.Tooltip(
//	    tooltip.TooltipProps{
//	        Position: forgeui.PositionRight,
//	        Delay:    300,
//	    },
//	    html.Button(g.Text("Help")),
//	    "Click here for assistance",
//	)
func Tooltip(props TooltipProps, trigger g.Node, content string) g.Node {
	// Apply defaults
	if props.Position == "" {
		props.Position = forgeui.PositionTop
	}
	if props.Delay == 0 {
		props.Delay = 200
	}

	// Get positioning classes
	positionClass := getTooltipPositionClass(props.Position)
	slideTransition := getTooltipTransition(props.Position)

	// Alpine state with delay
	alpineData := alpine.XData(map[string]any{
		"show":  false,
		"timer": nil,
	})

	return html.Div(
		alpineData,
		html.Class("relative inline-block"),

		// Trigger with hover/focus handlers
		html.Div(
			alpine.XOn("mouseenter", fmt.Sprintf("timer = setTimeout(() => show = true, %d)", props.Delay)),
			alpine.XOn("mouseleave", "clearTimeout(timer); show = false"),
			alpine.XOn("focus", fmt.Sprintf("timer = setTimeout(() => show = true, %d)", props.Delay)),
			alpine.XOn("blur", "clearTimeout(timer); show = false"),
			trigger,
		),

		// Tooltip content
		html.Div(
			alpine.XShow("show"),
			alpine.XCloak(),
			g.Group(alpine.XTransition(slideTransition)),
			html.Class(fmt.Sprintf(
				"absolute z-50 %s bg-primary text-primary-foreground px-3 py-1.5 text-sm rounded-md shadow-md whitespace-nowrap pointer-events-none %s",
				positionClass,
				props.Class,
			)),
			g.Attr("role", "tooltip"),
			g.Text(content),
		),
	)
}

// TooltipProvider wraps children with tooltip context for default settings.
//
// This is useful when you want to set default delay or other props for
// multiple tooltips in a section.
func TooltipProvider(children ...g.Node) g.Node {
	return html.Div(
		g.Group(children),
	)
}

// getTooltipPositionClass returns positioning classes based on position
func getTooltipPositionClass(position forgeui.Position) string {
	switch position {
	case forgeui.PositionTop:
		return "bottom-full mb-2 left-1/2 -translate-x-1/2"
	case forgeui.PositionRight:
		return "left-full ml-2 top-1/2 -translate-y-1/2"
	case forgeui.PositionBottom:
		return "top-full mt-2 left-1/2 -translate-x-1/2"
	case forgeui.PositionLeft:
		return "right-full mr-2 top-1/2 -translate-y-1/2"
	default:
		return "bottom-full mb-2 left-1/2 -translate-x-1/2"
	}
}

// getTooltipTransition returns the appropriate transition based on position
func getTooltipTransition(position forgeui.Position) *animation.Transition {
	// Tooltips use faster, simpler animations than popovers
	switch position {
	case forgeui.PositionTop:
		return animation.FadeIn()
	case forgeui.PositionRight:
		return animation.FadeIn()
	case forgeui.PositionBottom:
		return animation.FadeIn()
	case forgeui.PositionLeft:
		return animation.FadeIn()
	default:
		return animation.FadeIn()
	}
}

