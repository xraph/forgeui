package toast

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/icons"
)

// Container renders the toast container that holds all toasts.
//
// This should be added once to your layout, typically near the end of the body.
//
// Example:
//
//	html.Body(
//	    // ... page content ...
//	    toast.Container(),
//	)
func Container() g.Node {
	return html.Div(
		html.Class("fixed z-50 pointer-events-none"),
		alpine.XData(nil),
		alpine.XBind("class", `{
			'top-4 right-4': $store.toasts.position === 'top-right',
			'top-4 left-4': $store.toasts.position === 'top-left',
			'top-4 left-1/2 -translate-x-1/2': $store.toasts.position === 'top-center',
			'bottom-4 right-4': $store.toasts.position === 'bottom-right',
			'bottom-4 left-4': $store.toasts.position === 'bottom-left',
			'bottom-4 left-1/2 -translate-x-1/2': $store.toasts.position === 'bottom-center'
		}`),
		g.El("template",
			g.Group(alpine.XForKeyed("toast in $store.toasts.items", "toast.id")),
			ToastItem(),
		),
	)
}

// ToastItem renders a single toast notification.
func ToastItem() g.Node {
	return html.Div(
		html.Class("pointer-events-auto mb-4 w-96 rounded-lg shadow-lg transition-all duration-300"),
		alpine.XShow("toast.visible"),
		alpine.XBind("class", `{
			'bg-blue-50 border border-blue-200 text-blue-900': toast.variant === 'info',
			'bg-green-50 border border-green-200 text-green-900': toast.variant === 'success',
			'bg-yellow-50 border border-yellow-200 text-yellow-900': toast.variant === 'warning',
			'bg-red-50 border border-red-200 text-red-900': toast.variant === 'error'
		}`),
		html.Div(
			html.Class("flex items-start p-4"),
			
			// Icon
			html.Div(
				html.Class("flex-shrink-0"),
				alpine.XShow("toast.variant === 'info'"),
				icons.Info(icons.WithSize(20), icons.WithClass("text-blue-500")),
			),
			html.Div(
				html.Class("flex-shrink-0"),
				alpine.XShow("toast.variant === 'success'"),
				icons.CheckCircle(icons.WithSize(20), icons.WithClass("text-green-500")),
			),
			html.Div(
				html.Class("flex-shrink-0"),
				alpine.XShow("toast.variant === 'warning'"),
				icons.AlertCircle(icons.WithSize(20), icons.WithClass("text-yellow-500")),
			),
			html.Div(
				html.Class("flex-shrink-0"),
				alpine.XShow("toast.variant === 'error'"),
				icons.XCircle(icons.WithSize(20), icons.WithClass("text-red-500")),
			),
			
			// Message
			html.Div(
				html.Class("ml-3 flex-1"),
				html.P(
					html.Class("text-sm font-medium"),
					alpine.XText("toast.message"),
				),
			),
			
			// Close button
			html.Button(
				html.Type("button"),
				html.Class("ml-4 flex-shrink-0 rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"),
				alpine.XClick("$store.toasts.dismiss(toast.id)"),
				g.Attr("aria-label", "Close"),
				icons.X(icons.WithSize(16)),
			),
		),
	)
}

// ShowScript returns JavaScript to show a toast.
// Useful for triggering toasts from server-side events.
//
// Example:
//
//	htmx.TriggerEventWithDetail(w, "showToast", map[string]any{
//	    "message": "Item saved!",
//	    "variant": "success",
//	})
func ShowScript(message, variant string, timeout int) string {
	return fmt.Sprintf("$store.toasts.show('%s', '%s', %d)", message, variant, timeout)
}

// InfoScript returns JavaScript to show an info toast.
func InfoScript(message string) string {
	return fmt.Sprintf("$store.toasts.info('%s')", message)
}

// SuccessScript returns JavaScript to show a success toast.
func SuccessScript(message string) string {
	return fmt.Sprintf("$store.toasts.success('%s')", message)
}

// WarningScript returns JavaScript to show a warning toast.
func WarningScript(message string) string {
	return fmt.Sprintf("$store.toasts.warning('%s')", message)
}

// ErrorScript returns JavaScript to show an error toast.
func ErrorScript(message string) string {
	return fmt.Sprintf("$store.toasts.error('%s')", message)
}

