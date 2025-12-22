package toast

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
)

// ToasterProps defines toaster container configuration
type ToasterProps struct {
	Position  string // "top-left", "top-right", "top-center", "bottom-left", "bottom-right", "bottom-center"
	MaxToasts int    // Maximum number of visible toasts (0 = unlimited)
	Class     string
}

// ToasterOption is a functional option for configuring the toaster
type ToasterOption func(*ToasterProps)

// WithPosition sets the toast container position
func WithPosition(pos string) ToasterOption {
	return func(p *ToasterProps) { p.Position = pos }
}

// WithMaxToasts sets the maximum number of visible toasts
func WithMaxToasts(maxToasts int) ToasterOption {
	return func(p *ToasterProps) { p.MaxToasts = maxToasts }
}

// WithToasterClass adds custom classes to the toaster
func WithToasterClass(class string) ToasterOption {
	return func(p *ToasterProps) { p.Class = class }
}

// Toaster creates the toast notification container.
//
// Place this component in your page layout (typically in the body)
// to display toast notifications. Toasts are managed via Alpine store.
//
// Example in your layout:
//
//	html.Body(
//	    toast.RegisterToastStore(), // Register the Alpine store
//	    toast.Toaster(),            // Toast container
//	    // ... rest of your content
//	)
//
// Trigger toasts from anywhere:
//
//	button.Button(
//	    g.Text("Show Success"),
//	    button.WithAttr(
//	        alpine.XOn("click", `
//	            $store.toast.add({
//	                title: 'Success!',
//	                description: 'Your action completed',
//	                variant: 'success',
//	                duration: 3000
//	            })
//	        `),
//	    ),
//	)
func Toaster(opts ...ToasterOption) g.Node {
	props := &ToasterProps{
		Position:  "bottom-right",
		MaxToasts: 0, // unlimited by default
	}

	for _, opt := range opts {
		opt(props)
	}

	// Determine positioning classes based on position
	positionClass := getToasterPositionClass(props.Position)

	return html.Div(
		g.Attr("x-data", fmt.Sprintf(`{
			maxToasts: %d,
			get visibleToasts() {
				if (this.maxToasts > 0) {
					return $store.toast.items.slice(0, this.maxToasts);
				}
				return $store.toast.items;
			}
		}`, props.MaxToasts)),
		html.Class(fmt.Sprintf("%s z-[100] flex flex-col gap-2 p-4 max-h-screen overflow-hidden %s", positionClass, props.Class)),
		g.Attr("aria-live", "polite"),
		g.Attr("aria-atomic", "false"),
		g.Attr("@keydown.escape.window", "$store.toast.clear()"),

		// Toast list
		g.El("template",
			alpine.XFor("toast in visibleToasts"),
			g.Attr(":key", "toast.id"),
			html.Div(
				g.Attr("x-data", `{
					show: false,
					progress: 100,
					timer: null,
					init() {
						this.$nextTick(() => {
							this.show = true;
							if (this.toast.duration > 0) {
								const interval = 50;
								const decrement = (interval / this.toast.duration) * 100;
								this.timer = setInterval(() => {
									this.progress -= decrement;
									if (this.progress <= 0) {
										clearInterval(this.timer);
										this.close();
									}
								}, interval);
							}
						});
					},
					close() {
						this.show = false;
						setTimeout(() => $store.toast.remove(this.toast.id), 300);
					},
					toast: $el.closest('[x-for]').__x.$data.toast
				}`),
				alpine.XShow("show"),
				g.Group(alpine.XTransition(&alpine.Transition{
					Enter:      "transition-all duration-200 ease-out",
					EnterStart: "opacity-0 scale-95 translate-y-1",
					EnterEnd:   "opacity-100 scale-100 translate-y-0",
					Leave:      "transition-all duration-150 ease-in",
					LeaveStart: "opacity-100 scale-100 translate-y-0",
					LeaveEnd:   "opacity-0 scale-95 translate-y-1",
				})),
				html.Class("relative flex gap-3 rounded-md border p-4 shadow-lg w-full max-w-md"),
				g.Attr(":class", `{
				'bg-card border-border text-foreground': toast.variant === 'success',
				'bg-destructive/10 border-destructive/50 text-destructive': toast.variant === 'error' || toast.variant === 'destructive',
				'bg-card border-border text-foreground': toast.variant === 'warning',
				'bg-card border-border text-foreground': !toast.variant || toast.variant === 'default'
			}`),
				g.Attr("role", "alert"),

				// Icon (if variant is set)
				html.Div(
					alpine.XIf("toast.variant"),
					html.Class("flex-shrink-0"),
					html.SVG(
						g.Attr("xmlns", "http://www.w3.org/2000/svg"),
						g.Attr("width", "20"),
						g.Attr("height", "20"),
						g.Attr("viewBox", "0 0 24 24"),
						g.Attr("fill", "none"),
						g.Attr("stroke", "currentColor"),
						g.Attr("stroke-width", "2"),
						g.Attr("stroke-linecap", "round"),
						g.Attr("stroke-linejoin", "round"),
						html.Class("text-current"),
						g.El("path", g.Attr(":d", `
							toast.variant === 'success' ? 'M22 11.08V12a10 10 0 1 1-5.93-9.14 M22 4 12 14.01l-3-3' :
							toast.variant === 'error' || toast.variant === 'destructive' ? 'm21.73 18-8-8 8-8M2.27 6l8 8-8 8' :
							toast.variant === 'warning' ? 'm21.73 18-8-8 8-8M2.27 6l8 8-8 8' :
							'M12 16v-4 M12 8h.01'
						`)),
					),
				),

				// Content
				html.Div(
					html.Class("flex-1"),

					// Title
					html.Div(
						html.Class("font-semibold text-sm"),
						alpine.XText("toast.title"),
					),

					// Description
					html.Div(
						alpine.XIf("toast.description"),
						html.Class("text-sm text-muted-foreground mt-1"),
						alpine.XText("toast.description"),
					),

					// Action button
					html.Button(
						alpine.XIf("toast.action"),
						g.Attr("type", "button"),
						g.Attr("@click", "toast.action.onClick"),
						html.Class("mt-2 text-sm font-medium underline underline-offset-4 hover:no-underline"),
						alpine.XText("toast.action?.label"),
					),
				),

				// Close button
				html.Button(
					g.Attr("type", "button"),
					alpine.XOn("click", "close()"),
					html.Class("flex-shrink-0 rounded-sm opacity-70 hover:opacity-100 transition-opacity"),
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
				html.Div(
					alpine.XIf("toast.showProgress && toast.duration > 0"),
					html.Class("absolute bottom-0 left-0 right-0 h-1 bg-muted/20 rounded-b-md overflow-hidden"),
					html.Div(
						g.Attr(":style", "`width: ${progress}%`"),
						html.Class("h-full bg-current transition-all duration-50"),
					),
				),
			),
		),
	)
}

// RegisterToastStore returns the Alpine store registration script.
//
// Include this once in your page (typically in the head or body)
// to register the toast Alpine store.
//
// Example:
//
//	html.Body(
//	    toast.RegisterToastStore(),
//	    toast.Toaster(),
//	    // ... content
//	)
func RegisterToastStore() g.Node {
	return alpine.RegisterStores(alpine.Store{
		Name: "toast",
		State: map[string]any{
			"items": []any{},
		},
		Methods: `
			add(toast) {
				const id = Date.now() + Math.random();
				this.items.push({
					...toast,
					id: id,
					variant: toast.variant || 'default',
					duration: toast.duration || 5000,
					showProgress: toast.showProgress !== false
				});
				
				// Auto-remove after duration (if > 0)
				if (toast.duration && toast.duration > 0) {
					setTimeout(() => this.remove(id), toast.duration);
				}
			},
			remove(id) {
				this.items = this.items.filter(t => t.id !== id);
			},
			clear() {
				this.items = [];
			}
		`,
	})
}

// ToastSuccess is a convenience helper to show a success toast.
//
// Returns an Alpine expression to trigger a success toast.
// Use with alpine.XOn:
//
//	button.Button(
//	    g.Text("Save"),
//	    button.WithAttr(alpine.XOn("click", toast.ToastSuccess("Saved successfully!"))),
//	)
func ToastSuccess(title string) string {
	return fmt.Sprintf("$store.toast.add({title: '%s', variant: 'success', duration: 3000})", title)
}

// ToastError is a convenience helper to show an error toast.
func ToastError(title string) string {
	return fmt.Sprintf("$store.toast.add({title: '%s', variant: 'error', duration: 5000})", title)
}

// ToastWarning is a convenience helper to show a warning toast.
func ToastWarning(title string) string {
	return fmt.Sprintf("$store.toast.add({title: '%s', variant: 'warning', duration: 4000})", title)
}

// ToastInfo is a convenience helper to show an info toast.
func ToastInfo(title string) string {
	return fmt.Sprintf("$store.toast.add({title: '%s', variant: 'default', duration: 3000})", title)
}

// getToasterPositionClass returns CSS classes for positioning the toaster
func getToasterPositionClass(position string) string {
	switch position {
	case "top-left":
		return "fixed top-0 left-0"
	case "top-right":
		return "fixed top-0 right-0"
	case "top-center":
		return "fixed top-0 left-1/2 -translate-x-1/2"
	case "bottom-left":
		return "fixed bottom-0 left-0"
	case "bottom-center":
		return "fixed bottom-0 left-1/2 -translate-x-1/2"
	case "bottom-right":
		fallthrough
	default:
		return "fixed bottom-0 right-0"
	}
}
