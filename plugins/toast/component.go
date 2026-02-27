package toast

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"

	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/icons"
)

// Container renders the toast container that holds all toasts.
func Container() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="fixed z-50 pointer-events-none"`); err != nil {
			return err
		}
		writeAlpineAttrs(w, alpine.XData(nil))
		writeAlpineAttrs(w, alpine.XBind("class", `{
			'top-4 right-4': $store.toasts.position === 'top-right',
			'top-4 left-4': $store.toasts.position === 'top-left',
			'top-4 left-1/2 -translate-x-1/2': $store.toasts.position === 'top-center',
			'bottom-4 right-4': $store.toasts.position === 'bottom-right',
			'bottom-4 left-4': $store.toasts.position === 'bottom-left',
			'bottom-4 left-1/2 -translate-x-1/2': $store.toasts.position === 'bottom-center'
		}`))
		if _, err := io.WriteString(w, `>`); err != nil {
			return err
		}

		// Template for x-for
		if _, err := io.WriteString(w, `<template`); err != nil {
			return err
		}
		writeAlpineAttrs(w, alpine.XForKeyed("toast in $store.toasts.items", "toast.id"))
		if _, err := io.WriteString(w, `>`); err != nil {
			return err
		}

		// Render toast item inline
		if err := renderToastItem(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `</template></div>`); err != nil {
			return err
		}

		return nil
	})
}

// ToastItem renders a single toast notification.
func ToastItem() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return renderToastItem(ctx, w)
	})
}

func renderToastItem(ctx context.Context, w io.Writer) error {
	// Toast wrapper
	if _, err := io.WriteString(w, `<div class="pointer-events-auto mb-4 w-96 rounded-lg shadow-lg transition-all duration-300"`); err != nil {
		return err
	}
	writeAlpineAttrs(w, alpine.XShow("toast.visible"))
	writeAlpineAttrs(w, alpine.XBind("class", `{
		'bg-blue-50 border border-blue-200 text-blue-900': toast.variant === 'info',
		'bg-green-50 border border-green-200 text-green-900': toast.variant === 'success',
		'bg-yellow-50 border border-yellow-200 text-yellow-900': toast.variant === 'warning',
		'bg-red-50 border border-red-200 text-red-900': toast.variant === 'error'
	}`))
	if _, err := io.WriteString(w, `>`); err != nil {
		return err
	}

	// Inner flex container
	if _, err := io.WriteString(w, `<div class="flex items-start p-4">`); err != nil {
		return err
	}

	// Icons for each variant
	iconVariants := []struct {
		variant string
		icon    func(...icons.Option) templ.Component
		color   string
	}{
		{"info", icons.Info, "text-blue-500"},
		{"success", icons.CheckCircle, "text-green-500"},
		{"warning", icons.AlertCircle, "text-yellow-500"},
		{"error", icons.XCircle, "text-red-500"},
	}

	for _, iv := range iconVariants {
		if _, err := io.WriteString(w, `<div class="flex-shrink-0"`); err != nil {
			return err
		}
		writeAlpineAttrs(w, alpine.XShow(fmt.Sprintf("toast.variant === '%s'", iv.variant)))
		if _, err := io.WriteString(w, `>`); err != nil {
			return err
		}
		if err := iv.icon(icons.WithSize(20), icons.WithClass(iv.color)).Render(ctx, w); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `</div>`); err != nil {
			return err
		}
	}

	// Message
	if _, err := io.WriteString(w, `<div class="ml-3 flex-1"><p class="text-sm font-medium"`); err != nil {
		return err
	}
	writeAlpineAttrs(w, alpine.XText("toast.message"))
	if _, err := io.WriteString(w, `></p></div>`); err != nil {
		return err
	}

	// Close button
	if _, err := io.WriteString(w, `<button type="button" class="ml-4 flex-shrink-0 rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500" aria-label="Close"`); err != nil {
		return err
	}
	writeAlpineAttrs(w, alpine.XClick("$store.toasts.dismiss(toast.id)"))
	if _, err := io.WriteString(w, `>`); err != nil {
		return err
	}
	if err := icons.X(icons.WithSize(16)).Render(ctx, w); err != nil {
		return err
	}
	if _, err := io.WriteString(w, `</button>`); err != nil {
		return err
	}

	// Close inner flex and outer div
	if _, err := io.WriteString(w, `</div></div>`); err != nil {
		return err
	}

	return nil
}

// writeAlpineAttrs writes templ.Attributes as HTML attributes.
func writeAlpineAttrs(w io.Writer, attrs templ.Attributes) {
	for k, v := range attrs {
		if s, ok := v.(string); ok {
			_, _ = fmt.Fprintf(w, ` %s="%s"`, k, s)
		} else if v == true {
			_, _ = fmt.Fprintf(w, ` %s`, k)
		}
	}
}

// ShowScript returns JavaScript to show a toast.
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
