package modal

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// DrawerProps defines configuration for a drawer component
type DrawerProps struct {
	// Side specifies which edge the drawer slides in from
	Side forgeui.Side

	// Size controls the width (for left/right) or height (for top/bottom)
	Size forgeui.Size

	// CloseOnEscape enables closing when Escape key is pressed
	CloseOnEscape bool

	// CloseOnOutsideClick enables closing when clicking the backdrop
	CloseOnOutsideClick bool

	// ShowClose displays a close button
	ShowClose bool

	// Class adds additional CSS classes
	Class string
}

// Drawer creates a drawer that slides in from a screen edge.
//
// Drawers are panels that slide in from the edge of the screen,
// commonly used for navigation menus, filters, or forms.
//
// Example:
//
//	drawer.Drawer(
//	    drawer.DrawerProps{
//	        Side: forgeui.SideRight,
//	        Size: forgeui.SizeMD,
//	    },
//	    button.Button(g.Text("Open Drawer")),
//	    html.Div(
//	        drawer.DrawerHeader("Settings", "Configure your preferences"),
//	        html.Div(g.Text("Drawer content")),
//	    ),
//	)
func Drawer(props DrawerProps, trigger g.Node, content ...g.Node) g.Node {
	// Apply defaults
	if props.Side == "" {
		props.Side = forgeui.SideRight
	}

	if props.Size == "" {
		props.Size = forgeui.SizeMD
	}

	// Get positioning classes
	positionClass := getDrawerPositionClass(props.Side)
	sizeClass := getDrawerSizeClass(props.Side, props.Size)
	slideTransition := getDrawerSlideTransition(props.Side)

	// Alpine.js state management
	alpineData := alpine.XData(map[string]any{
		"open": false,
	})

	// Keyboard handler
	var escapeHandler g.Node
	if props.CloseOnEscape {
		escapeHandler = alpine.XOn("keydown.escape.window", "open = false")
	}

	// Backdrop click handler
	var backdropClick g.Node
	if props.CloseOnOutsideClick {
		backdropClick = alpine.XOn("click", "open = false")
	}

	return html.Div(
		alpineData,
		g.If(props.CloseOnEscape, escapeHandler),

		// Trigger
		html.Div(
			alpine.XOn("click", "open = true"),
			trigger,
		),

		// Drawer overlay
		html.Div(
			alpine.XShow("open"),
			alpine.XCloak(),
			html.Class("fixed inset-0 z-[60]"),
			g.Attr("aria-modal", "true"),
			g.Attr("role", "dialog"),

			// Backdrop
			html.Div(
				g.Group(alpine.XTransition(animation.FadeIn())),
				html.Class("fixed inset-0 bg-background/80 backdrop-blur-sm transition-all"),
				g.If(props.CloseOnOutsideClick, backdropClick),
			),

			// Drawer panel
			html.Div(
				g.Group(alpine.XTransition(slideTransition)),
				alpine.XOn("click.stop", ""),
				g.Attr("x-trap.noscroll", "open"),
				html.Class("fixed "+positionClass+" "+sizeClass+" bg-background border border-border shadow-lg overflow-y-auto "+props.Class),

				// Close button
				g.If(props.ShowClose, closeButton()),

				// Content
				g.Group(content),
			),
		),
	)
}

// DrawerHeader creates a header section for the drawer.
func DrawerHeader(title, description string) g.Node {
	return html.Div(
		html.Class("px-6 pt-6 pb-4 border-b border-border"),
		html.H2(
			html.Class("text-lg font-semibold leading-none tracking-tight"),
			g.Text(title),
		),
		g.If(description != "", html.P(
			html.Class("text-sm text-muted-foreground mt-2"),
			g.Text(description),
		)),
	)
}

// DrawerBody creates the main content area of the drawer.
func DrawerBody(children ...g.Node) g.Node {
	return html.Div(
		html.Class("px-6 py-4"),
		g.Group(children),
	)
}

// DrawerFooter creates a footer section for action buttons.
func DrawerFooter(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex justify-end gap-3 px-6 pb-6 pt-4 border-t border-border mt-auto"),
		g.Group(children),
	)
}

// DrawerClose wraps an element to close the drawer when clicked.
func DrawerClose(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = false"),
		child,
	)
}

// getDrawerPositionClass returns positioning classes based on side
func getDrawerPositionClass(side forgeui.Side) string {
	switch side {
	case forgeui.SideTop:
		return "top-0 left-0 right-0"
	case forgeui.SideRight:
		return "top-0 right-0 bottom-0"
	case forgeui.SideBottom:
		return "bottom-0 left-0 right-0"
	case forgeui.SideLeft:
		return "top-0 left-0 bottom-0"
	default:
		return "top-0 right-0 bottom-0"
	}
}

// getDrawerSizeClass returns size classes based on side and size
func getDrawerSizeClass(side forgeui.Side, size forgeui.Size) string {
	if side == forgeui.SideTop || side == forgeui.SideBottom {
		// Height for top/bottom drawers
		switch size {
		case forgeui.SizeSM:
			return "h-1/4"
		case forgeui.SizeMD:
			return "h-1/3"
		case forgeui.SizeLG:
			return "h-1/2"
		case forgeui.SizeXL:
			return "h-2/3"
		case forgeui.SizeFull:
			return "h-full"
		default:
			return "h-1/3"
		}
	} else {
		// Width for left/right drawers
		switch size {
		case forgeui.SizeSM:
			return "w-80"
		case forgeui.SizeMD:
			return "w-96"
		case forgeui.SizeLG:
			return "w-[32rem]"
		case forgeui.SizeXL:
			return "w-[40rem]"
		case forgeui.SizeFull:
			return "w-full"
		default:
			return "w-96"
		}
	}
}

// getDrawerSlideTransition returns the appropriate slide transition based on side
func getDrawerSlideTransition(side forgeui.Side) *animation.Transition {
	switch side {
	case forgeui.SideTop:
		return animation.SlideInFromTop()
	case forgeui.SideRight:
		return animation.SlideInFromRight()
	case forgeui.SideBottom:
		return animation.SlideInFromBottom()
	case forgeui.SideLeft:
		return animation.SlideInFromLeft()
	default:
		return animation.SlideInFromRight()
	}
}
