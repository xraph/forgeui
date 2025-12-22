package modal

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// Dialog creates a dialog container with compound component API.
//
// Dialog provides a more semantic API than the base Modal component,
// with separate components for trigger, content, header, footer, etc.
//
// Example:
//
//	dialog.Dialog(
//	    dialog.DialogTrigger(button.Button(g.Text("Edit Profile"))),
//	    dialog.DialogContent(
//	        dialog.DialogHeader(
//	            dialog.DialogTitle("Edit Profile"),
//	            dialog.DialogDescription("Make changes here."),
//	        ),
//	        html.Div(/* form fields */),
//	        dialog.DialogFooter(
//	            dialog.DialogClose(button.Button(g.Text("Cancel"))),
//	            button.Button(g.Text("Save")),
//	        ),
//	    ),
//	)
func Dialog(children ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{
			"open": false,
		}),
		alpine.XOn("keydown.escape.window", "open = false"),
		g.Group(children),
	)
}

// DialogTrigger wraps an element to trigger opening the dialog when clicked.
func DialogTrigger(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = true"),
		child,
	)
}

// DialogContent creates the dialog content panel with backdrop.
//
// Options can be passed to customize size and behavior:
//
//	DialogContent(opts ...DialogOption)
func DialogContent(children ...g.Node) g.Node {
	return DialogContentWithOptions(DialogContentProps{}, children...)
}

// DialogContentProps defines configuration for dialog content
type DialogContentProps struct {
	Size                forgeui.Size
	CloseOnOutsideClick bool
	ShowClose           bool
	Class               string
}

// DialogContentWithOptions creates dialog content with custom props
func DialogContentWithOptions(props DialogContentProps, children ...g.Node) g.Node {
	// Apply defaults
	if props.Size == "" {
		props.Size = forgeui.SizeMD
	}

	if !props.CloseOnOutsideClick {
		props.CloseOnOutsideClick = false // Default to false for dialogs
	}

	props.ShowClose = true // Always show close for dialogs

	sizeClass := getSizeClass(props.Size)

	var backdropClick g.Node
	if props.CloseOnOutsideClick {
		backdropClick = alpine.XOn("click", "open = false")
	}

	return html.Div(
		alpine.XShow("open"),
		alpine.XCloak(),
		html.Class("fixed inset-0 z-[60] overflow-y-auto"),
		g.Attr("aria-modal", "true"),
		g.Attr("role", "dialog"),

		// Backdrop
		html.Div(
			g.Group(alpine.XTransition(animation.FadeIn())),
			html.Class("fixed inset-0 bg-background/80 backdrop-blur-sm transition-all"),
			g.If(props.CloseOnOutsideClick, backdropClick),
		),

		// Content container
		html.Div(
			html.Class("fixed inset-0 flex items-center justify-center p-4"),

			// Dialog panel
			html.Div(
				g.Group(alpine.XTransition(animation.ScaleIn())),
				alpine.XOn("click.stop", ""),
				g.Attr("x-trap.noscroll", "open"),
				html.Class("relative bg-background rounded-lg shadow-lg w-full "+sizeClass+" "+props.Class+" border border-border"),

				// Close button
				closeButton(),

				// Content
				g.Group(children),
			),
		),
	)
}

// DialogHeader creates a header section for the dialog.
func DialogHeader(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex flex-col gap-2 px-6 pt-6 pb-4"),
		g.Group(children),
	)
}

// DialogTitle creates the dialog title element.
func DialogTitle(text string) g.Node {
	return html.H2(
		html.ID("dialog-title"),
		html.Class("text-lg font-semibold leading-none tracking-tight"),
		g.Text(text),
	)
}

// DialogDescription creates the dialog description element.
func DialogDescription(text string) g.Node {
	return html.P(
		html.Class("text-sm text-muted-foreground"),
		g.Text(text),
	)
}

// DialogBody creates the main content area of the dialog.
func DialogBody(children ...g.Node) g.Node {
	return html.Div(
		html.Class("px-6 py-4"),
		g.Group(children),
	)
}

// DialogFooter creates a footer section for action buttons.
func DialogFooter(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex justify-end gap-3 px-6 pb-6 pt-4 border-t border-border"),
		g.Group(children),
	)
}

// DialogClose wraps an element to close the dialog when clicked.
func DialogClose(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = false"),
		child,
	)
}

// Dialog size options for convenience
func DialogSM(children ...g.Node) g.Node {
	return DialogContentWithOptions(DialogContentProps{Size: forgeui.SizeSM}, children...)
}

func DialogLG(children ...g.Node) g.Node {
	return DialogContentWithOptions(DialogContentProps{Size: forgeui.SizeLG}, children...)
}

func DialogXL(children ...g.Node) g.Node {
	return DialogContentWithOptions(DialogContentProps{Size: forgeui.SizeXL}, children...)
}

func DialogFull(children ...g.Node) g.Node {
	return DialogContentWithOptions(DialogContentProps{Size: forgeui.SizeFull}, children...)
}
