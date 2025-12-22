package modal

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// AlertDialog creates an alert dialog for critical confirmations.
//
// Unlike regular dialogs, alert dialogs:
//   - Cannot be closed by clicking outside (by design)
//   - Require explicit user action
//   - Are typically used for destructive actions
//   - Auto-focus on the primary action button
//
// Example:
//
//	alertDialog.AlertDialog(
//	    alertDialog.AlertDialogTrigger(
//	        button.Button(g.Text("Delete Account"), button.WithVariant("destructive")),
//	    ),
//	    alertDialog.AlertDialogContent(
//	        alertDialog.AlertDialogHeader(
//	            alertDialog.AlertDialogTitle("Are you absolutely sure?"),
//	            alertDialog.AlertDialogDescription("This action cannot be undone."),
//	        ),
//	        alertDialog.AlertDialogFooter(
//	            alertDialog.AlertDialogCancel(button.Button(g.Text("Cancel"))),
//	            alertDialog.AlertDialogAction(button.Button(g.Text("Delete"), button.WithVariant("destructive"))),
//	        ),
//	    ),
//	)
func AlertDialog(children ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{
			"open": false,
		}),
		alpine.XOn("keydown.escape.window", "open = false"),
		g.Group(children),
	)
}

// AlertDialogTrigger wraps an element to trigger opening the alert dialog.
func AlertDialogTrigger(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = true"),
		child,
	)
}

// AlertDialogContent creates the alert dialog content panel.
//
// By design, alert dialogs:
// - Cannot close on outside click
// - Always show close button
// - Use medium size by default
func AlertDialogContent(children ...g.Node) g.Node {
	return AlertDialogContentWithSize(forgeui.SizeMD, children...)
}

// AlertDialogContentWithSize creates alert dialog content with custom size
func AlertDialogContentWithSize(size forgeui.Size, children ...g.Node) g.Node {
	if size == "" {
		size = forgeui.SizeMD
	}

	sizeClass := getSizeClass(size)

	return html.Div(
		alpine.XShow("open"),
		alpine.XCloak(),
		html.Class("fixed inset-0 z-[60] overflow-y-auto"),
		g.Attr("aria-modal", "true"),
		g.Attr("role", "alertdialog"),
		g.Attr("aria-labelledby", "alert-dialog-title"),
		g.Attr("aria-describedby", "alert-dialog-description"),

		// Backdrop - no click handler (cannot close by clicking outside)
		html.Div(
			g.Group(alpine.XTransition(animation.FadeIn())),
			html.Class("fixed inset-0 bg-background/80 backdrop-blur-sm transition-all"),
		),

		// Content container
		html.Div(
			html.Class("fixed inset-0 flex items-center justify-center p-4"),

			// Alert dialog panel
			html.Div(
				g.Group(alpine.XTransition(animation.ScaleIn())),
				alpine.XOn("click.stop", ""),
				g.Attr("x-trap.noscroll", "open"),
				html.Class("relative bg-background rounded-lg shadow-lg w-full "+sizeClass+" border border-border"),

				// Close button (top right)
				closeButton(),

				// Content
				g.Group(children),
			),
		),
	)
}

// AlertDialogHeader creates a header section for the alert dialog.
func AlertDialogHeader(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex flex-col gap-2 px-6 pt-6 pb-4"),
		g.Group(children),
	)
}

// AlertDialogTitle creates the alert dialog title element.
func AlertDialogTitle(text string) g.Node {
	return html.H2(
		html.ID("alert-dialog-title"),
		html.Class("text-lg font-semibold leading-none tracking-tight"),
		g.Text(text),
	)
}

// AlertDialogDescription creates the alert dialog description element.
func AlertDialogDescription(text string) g.Node {
	return html.P(
		html.ID("alert-dialog-description"),
		html.Class("text-sm text-muted-foreground"),
		g.Text(text),
	)
}

// AlertDialogFooter creates a footer section for action buttons.
//
// Typically contains Cancel and Action buttons.
func AlertDialogFooter(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex justify-end gap-3 px-6 pb-6 pt-4"),
		g.Group(children),
	)
}

// AlertDialogCancel wraps a cancel button to close the alert dialog.
//
// This is the safe action - typically an outline or ghost button.
func AlertDialogCancel(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = false"),
		child,
	)
}

// AlertDialogAction wraps a confirm/action button.
//
// This should be the primary action button, often with a destructive variant.
// The button itself handles the actual action; this wrapper only provides
// the proper HTML structure.
func AlertDialogAction(child g.Node) g.Node {
	return html.Div(
		// The actual action is handled by the button's own click handler
		// We just wrap it here for consistency
		g.Attr("x-ref", "action"),
		child,
	)
}
