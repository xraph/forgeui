// Package modal provides modal dialog components with Alpine.js state management.
//
// Modal is the foundation for overlay components including Dialog, AlertDialog,
// Drawer, and Sheet. It provides:
//   - Backdrop with fade transition
//   - Content with scale-in transition
//   - Configurable sizes (sm, md, lg, xl, full)
//   - Close on escape (configurable)
//   - Close on outside click (configurable)
//   - Optional close button
//   - Focus trap using Alpine Focus plugin
//
// Basic usage:
//
//	modal.Modal(
//	    modal.ModalProps{
//	        Size: forgeui.SizeMD,
//	        CloseOnEscape: true,
//	        CloseOnOutsideClick: true,
//	    },
//	    button.Button(g.Text("Open Modal")),
//	    html.Div(
//	        modal.ModalHeader("Confirm Action", "Are you sure?"),
//	        html.P(g.Text("This action cannot be undone.")),
//	        modal.ModalFooter(
//	            button.Button(g.Text("Cancel")),
//	            button.Button(g.Text("Confirm")),
//	        ),
//	    ),
//	)
package modal

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// ModalProps defines configuration for a modal dialog
type ModalProps struct {
	// Size controls the maximum width of the modal
	Size forgeui.Size

	// CloseOnEscape enables closing the modal when Escape key is pressed
	CloseOnEscape bool

	// CloseOnOutsideClick enables closing when clicking the backdrop
	CloseOnOutsideClick bool

	// ShowClose displays a close button in the top-right corner
	ShowClose bool

	// Class adds additional CSS classes to the modal content
	Class string
}

// defaultModalProps returns default modal properties
func defaultModalProps() ModalProps {
	return ModalProps{
		Size:                forgeui.SizeMD,
		CloseOnEscape:       true,
		CloseOnOutsideClick: false,
		ShowClose:           true,
	}
}

// Modal creates a modal dialog with Alpine.js state management.
//
// The modal consists of three main parts:
//  1. Trigger element - clicks to open the modal
//  2. Backdrop - semi-transparent overlay
//  3. Content - the actual modal panel
//
// Example:
//
//	modal.Modal(
//	    modal.ModalProps{Size: forgeui.SizeLG},
//	    button.Button(g.Text("Open")),
//	    html.Div(g.Text("Modal content")),
//	)
func Modal(props ModalProps, trigger g.Node, content ...g.Node) g.Node {
	// Apply defaults
	if props.Size == "" {
		props.Size = forgeui.SizeMD
	}

	// Build size classes
	sizeClass := getSizeClass(props.Size)

	// Alpine.js state management
	alpineData := alpine.XData(map[string]any{
		"open": false,
	})

	// Keyboard handler for Escape key
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

		// Trigger slot - wraps trigger with click handler
		html.Div(
			alpine.XOn("click", "open = true"),
			trigger,
		),

		// Modal overlay
		html.Div(
			alpine.XShow("open"),
			alpine.XCloak(),
			html.Class("fixed inset-0 z-[60] overflow-y-auto"),
			g.Attr("aria-modal", "true"),
			g.Attr("role", "dialog"),
			g.Attr("aria-labelledby", "modal-title"),

		// Backdrop
		html.Div(
			g.Group(alpine.XTransition(animation.FadeIn())),
			html.Class("fixed inset-0 bg-background/80 backdrop-blur-sm transition-all"),
			g.If(props.CloseOnOutsideClick, backdropClick),
		),

			// Content container - centers the modal
			html.Div(
				html.Class("fixed inset-0 flex items-center justify-center p-4"),

				// Modal panel
				html.Div(
					g.Group(alpine.XTransition(animation.ScaleIn())),
					alpine.XOn("click.stop", ""), // Prevent close on content click
					g.Attr("x-trap.noscroll", "open"), // Focus trap with Alpine Focus plugin
					html.Class(fmt.Sprintf(
						"relative bg-background rounded-lg shadow-lg w-full %s %s border border-border",
						sizeClass,
						props.Class,
					)),

					// Close button
					g.If(props.ShowClose, closeButton()),

					// Content
					g.Group(content),
				),
			),
		),
	)
}

// ModalHeader creates a header section for the modal with title and optional description.
//
// Example:
//
//	modal.ModalHeader("Delete Account", "This action cannot be undone.")
func ModalHeader(title, description string) g.Node {
	return html.Div(
		html.Class("px-6 pt-6 pb-4"),
		html.H2(
			html.ID("modal-title"),
			html.Class("text-lg font-semibold leading-none tracking-tight"),
			g.Text(title),
		),
		g.If(description != "", html.P(
			html.Class("text-sm text-muted-foreground mt-2"),
			g.Text(description),
		)),
	)
}

// ModalFooter creates a footer section for action buttons.
//
// Example:
//
//	modal.ModalFooter(
//	    button.Button(g.Text("Cancel")),
//	    button.Button(g.Text("Confirm")),
//	)
func ModalFooter(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex justify-end gap-3 px-6 pb-6 pt-4 border-t border-border"),
		g.Group(children),
	)
}

// ModalClose wraps a child element to close the modal when clicked.
//
// Example:
//
//	modal.ModalClose(button.Button(g.Text("Close")))
func ModalClose(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = false"),
		child,
	)
}

// closeButton creates the X button in the top-right corner
func closeButton() g.Node {
	return html.Button(
		alpine.XOn("click", "open = false"),
		html.Class("absolute right-4 top-4 rounded-sm opacity-70 hover:opacity-100 transition-opacity focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"),
		g.Attr("aria-label", "Close"),
		g.Attr("type", "button"),
		// X icon using SVG
		html.SVG(
			g.Attr("xmlns", "http://www.w3.org/2000/svg"),
			g.Attr("width", "16"),
			g.Attr("height", "16"),
			g.Attr("viewBox", "0 0 24 24"),
			g.Attr("fill", "none"),
			g.Attr("stroke", "currentColor"),
			g.Attr("stroke-width", "2"),
			g.Attr("stroke-linecap", "round"),
			g.Attr("stroke-linejoin", "round"),
		g.El("path", g.Attr("d", "M18 6 6 18")),
		g.El("path", g.Attr("d", "m6 6 12 12")),
		),
	)
}

// getSizeClass returns the appropriate max-width class for the modal size
func getSizeClass(size forgeui.Size) string {
	switch size {
	case forgeui.SizeSM:
		return "max-w-sm"
	case forgeui.SizeMD:
		return "max-w-md"
	case forgeui.SizeLG:
		return "max-w-lg"
	case forgeui.SizeXL:
		return "max-w-xl"
	case forgeui.SizeFull:
		return "max-w-[calc(100vw-2rem)] max-h-[calc(100vh-2rem)]"
	default:
		return "max-w-md"
	}
}

