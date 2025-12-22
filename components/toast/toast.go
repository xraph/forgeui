// Package toast provides notification system with queue management.
//
// Toasts are temporary notifications that appear on screen, typically
// at the top or bottom. They support:
//   - Multiple variants (default, success, error, warning)
//   - Auto-dismiss with timer
//   - Manual dismiss
//   - Stacking/queuing
//   - Progress bar
//   - Action buttons
//
// Basic usage:
//
//	// In your page, include the Toaster container:
//	html.Body(
//	    toast.Toaster(),
//	    // ... rest of your content
//	)
//
//	// Trigger a toast notification:
//	button.Button(
//	    g.Text("Show Toast"),
//	    button.WithAttr(
//	        alpine.XOn("click", "$store.toast.add({title: 'Success', variant: 'success', duration: 3000})"),
//	    ),
//	)
package toast

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// ToastProps defines configuration for a toast notification
type ToastProps struct {
	// Variant specifies the visual style
	Variant forgeui.Variant

	// Duration is the auto-dismiss time in milliseconds (0 = no auto-dismiss)
	Duration int

	// Title is the main toast message
	Title string

	// Description is optional additional text
	Description string

	// Action is an optional action button
	Action *ToastAction

	// ShowProgress displays a progress bar for auto-dismiss
	ShowProgress bool

	// Class adds additional CSS classes
	Class string
}

// ToastAction defines an action button configuration
type ToastAction struct {
	Label   string
	OnClick string // Alpine expression
}

// Toast creates a single toast notification.
//
// Note: Typically you won't call this directly. Instead, use the
// Toaster container and trigger toasts via Alpine store:
//
//	$store.toast.add({title: "Success", variant: "success", duration: 3000})
func Toast(props ToastProps) g.Node {
	// Apply defaults
	if props.Variant == "" {
		props.Variant = forgeui.VariantDefault
	}
	if props.Duration == 0 {
		props.Duration = 5000 // 5 seconds default
	}

	// Get variant classes
	variantClasses := getToastVariantClasses(props.Variant)

	return html.Div(
		g.Attr("x-data", fmt.Sprintf(`{
			show: false,
			progress: 100,
			timer: null,
			init() {
				this.show = true;
				if (%d > 0) {
					const interval = 50;
					const decrement = (interval / %d) * 100;
					this.timer = setInterval(() => {
						this.progress -= decrement;
						if (this.progress <= 0) {
							clearInterval(this.timer);
							this.close();
						}
					}, interval);
				}
			},
			close() {
				this.show = false;
				setTimeout(() => this.$el.remove(), 300);
			}
		}`, props.Duration, props.Duration)),
		alpine.XShow("show"),
		g.Group(alpine.XTransition(animation.SlideInFromBottom())),
		html.Class(fmt.Sprintf(
			"relative flex gap-3 rounded-md border p-4 shadow-lg w-full max-w-md %s %s",
			variantClasses,
			props.Class,
		)),
		g.Attr("role", "alert"),
		g.Attr("aria-live", "polite"),

		// Icon (if variant is set)
		g.If(props.Variant != "", toastIcon(props.Variant)),

		// Content
		html.Div(
			html.Class("flex-1"),

			// Title
			html.Div(
				html.Class("font-semibold text-sm"),
				g.Text(props.Title),
			),

		// Description
		g.If(props.Description != "", html.Div(
			html.Class("text-sm text-muted-foreground mt-1"),
			g.Text(props.Description),
		)),

		// Action button (if provided)
		renderToastAction(props.Action),
		),

		// Close button
		html.Button(
			g.Attr("type", "button"),
			alpine.XOn("click", "close()"),
			html.Class("rounded-sm opacity-70 hover:opacity-100 transition-opacity"),
			g.Attr("aria-label", "Close"),
			// X icon
			html.SVG(
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("width", "16"),
				g.Attr("height", "16"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("fill", "none"),
				g.Attr("stroke", "currentColor"),
				g.Attr("stroke-width", "2"),
				g.El("path", g.Attr("d", "M18 6 6 18")),
				g.El("path", g.Attr("d", "m6 6 12 12")),
			),
		),

		// Progress bar
		g.If(props.ShowProgress && props.Duration > 0, html.Div(
			html.Class("absolute bottom-0 left-0 right-0 h-1 bg-muted/20 rounded-b-md overflow-hidden"),
			html.Div(
				g.Attr(":style", "`width: ${progress}%`"),
				html.Class("h-full bg-current transition-all duration-50"),
			),
		)),
	)
}

// renderToastAction renders the action button if provided
func renderToastAction(action *ToastAction) g.Node {
	if action == nil {
		return g.Text("")
	}
	return html.Button(
		g.Attr("type", "button"),
		alpine.XOn("click", action.OnClick),
		html.Class("mt-2 text-sm font-medium underline underline-offset-4 hover:no-underline"),
		g.Text(action.Label),
	)
}

// toastIcon returns the appropriate icon for each variant
func toastIcon(variant forgeui.Variant) g.Node {
	var path string
	iconClass := "flex-shrink-0 h-5 w-5 text-current"

	switch variant {
	case "success":
		// Checkmark circle
		path = "M22 11.08V12a10 10 0 1 1-5.93-9.14 M22 4 12 14.01l-3-3"
	case "error", forgeui.VariantDestructive:
		// X circle
		path = "m21.73 18-8-8 8-8M2.27 6l8 8-8 8"
	case "warning":
		// Alert triangle
		path = "m21.73 18-8-8 8-8M2.27 6l8 8-8 8"
	default:
		// Info circle
		path = "M12 16v-4 M12 8h.01"
	}

	return html.SVG(
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("width", "20"),
		g.Attr("height", "20"),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", "currentColor"),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
		html.Class(iconClass),
		g.El("path", g.Attr("d", path)),
	)
}

// getToastVariantClasses returns styling classes for each variant
func getToastVariantClasses(variant forgeui.Variant) string {
	// Use theme colors for all variants to ensure proper light/dark mode support
	switch variant {
	case "success":
		return "bg-card border-border text-foreground"
	case "error", forgeui.VariantDestructive:
		return "bg-destructive/10 border-destructive/50 text-destructive"
	case "warning":
		return "bg-card border-border text-foreground"
	default:
		return "bg-card border-border text-foreground"
	}
}

