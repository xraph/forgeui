package modal

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
)

// Sheet creates a sheet container with compound component API.
//
// Sheet provides a more semantic API than the base Drawer component,
// similar to how Dialog wraps Modal.
//
// Example:
//
//	sheet.Sheet(
//	    sheet.SheetTrigger(button.Button(g.Text("Open Settings"))),
//	    sheet.SheetContent(
//	        forgeui.SideRight,
//	        sheet.SheetHeader(
//	            sheet.SheetTitle("Settings"),
//	            sheet.SheetDescription("Configure your preferences"),
//	        ),
//	        html.Div(/* settings content */),
//	        sheet.SheetFooter(
//	            sheet.SheetClose(button.Button(g.Text("Close"))),
//	        ),
//	    ),
//	)
func Sheet(children ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{
			"open": false,
		}),
		alpine.XOn("keydown.escape.window", "open = false"),
		g.Group(children),
	)
}

// SheetTrigger wraps an element to trigger opening the sheet when clicked.
func SheetTrigger(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = true"),
		child,
	)
}

// SheetContent creates the sheet content panel with specified side.
func SheetContent(side forgeui.Side, children ...g.Node) g.Node {
	return SheetContentWithOptions(SheetContentProps{Side: side}, children...)
}

// SheetContentProps defines configuration for sheet content
type SheetContentProps struct {
	Side                forgeui.Side
	Size                forgeui.Size
	CloseOnOutsideClick bool
	ShowClose           bool
	Class               string
}

// SheetContentWithOptions creates sheet content with custom props
func SheetContentWithOptions(props SheetContentProps, children ...g.Node) g.Node {
	// Apply defaults
	if props.Side == "" {
		props.Side = forgeui.SideRight
	}
	if props.Size == "" {
		props.Size = forgeui.SizeMD
	}
	if !props.CloseOnOutsideClick {
		props.CloseOnOutsideClick = true // Default to true for sheets
	}
	props.ShowClose = true // Always show close

	positionClass := getDrawerPositionClass(props.Side)
	sizeClass := getDrawerSizeClass(props.Side, props.Size)
	slideTransition := getDrawerSlideTransition(props.Side)

	var backdropClick g.Node
	if props.CloseOnOutsideClick {
		backdropClick = alpine.XOn("click", "open = false")
	}

	return html.Div(
		alpine.XShow("open"),
		alpine.XCloak(),
		html.Class("fixed inset-0 z-[60]"),
		g.Attr("aria-modal", "true"),
		g.Attr("role", "dialog"),

		// Backdrop
		html.Div(
			g.Group(alpine.XTransition(slideTransition)),
			html.Class("fixed inset-0 bg-background/80 backdrop-blur-sm transition-all"),
			g.If(props.CloseOnOutsideClick, backdropClick),
		),

		// Sheet panel
		html.Div(
			g.Group(alpine.XTransition(slideTransition)),
			alpine.XOn("click.stop", ""),
			g.Attr("x-trap.noscroll", "open"),
			html.Class("fixed "+positionClass+" "+sizeClass+" bg-background border border-border shadow-lg overflow-y-auto "+props.Class),

			// Close button
			closeButton(),

			// Content
			g.Group(children),
		),
	)
}

// SheetHeader creates a header section for the sheet.
func SheetHeader(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex flex-col gap-2 px-6 pt-6 pb-4 border-b border-border"),
		g.Group(children),
	)
}

// SheetTitle creates the sheet title element.
func SheetTitle(text string) g.Node {
	return html.H2(
		html.ID("sheet-title"),
		html.Class("text-lg font-semibold leading-none tracking-tight"),
		g.Text(text),
	)
}

// SheetDescription creates the sheet description element.
func SheetDescription(text string) g.Node {
	return html.P(
		html.Class("text-sm text-muted-foreground"),
		g.Text(text),
	)
}

// SheetBody creates the main content area of the sheet.
func SheetBody(children ...g.Node) g.Node {
	return html.Div(
		html.Class("px-6 py-4 flex-1"),
		g.Group(children),
	)
}

// SheetFooter creates a footer section for action buttons.
func SheetFooter(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex justify-end gap-3 px-6 pb-6 pt-4 border-t border-border mt-auto"),
		g.Group(children),
	)
}

// SheetClose wraps an element to close the sheet when clicked.
func SheetClose(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = false"),
		child,
	)
}

// Convenience functions for different sides
func SheetLeft(children ...g.Node) g.Node {
	return SheetContent(forgeui.SideLeft, children...)
}

func SheetRight(children ...g.Node) g.Node {
	return SheetContent(forgeui.SideRight, children...)
}

func SheetTop(children ...g.Node) g.Node {
	return SheetContent(forgeui.SideTop, children...)
}

func SheetBottom(children ...g.Node) g.Node {
	return SheetContent(forgeui.SideBottom, children...)
}

