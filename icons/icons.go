// Package icons provides a flexible icon system with Lucide icon integration.
// Icons are rendered as inline SVG elements with customizable size, color, and stroke width.
//
//go:generate go run internal/generate/main.go
package icons

import (
	"context"
	"fmt"
	stdhtml "html"
	"io"

	"github.com/a-h/templ"

	"github.com/xraph/forgeui"
)

// Props defines icon configuration
type Props struct {
	Size        int     // Size in pixels (width and height)
	Color       string  // CSS color value
	StrokeWidth float64 // SVG stroke width
	Class       string  // Additional CSS classes
	Attrs       templ.Attributes
}

// Option is a functional option for configuring icons
type Option func(*Props)

// WithSize sets the icon size in pixels
func WithSize(size int) Option {
	return func(p *Props) { p.Size = size }
}

// WithColor sets the icon color
func WithColor(color string) Option {
	return func(p *Props) { p.Color = color }
}

// WithStrokeWidth sets the SVG stroke width
func WithStrokeWidth(width float64) Option {
	return func(p *Props) { p.StrokeWidth = width }
}

// WithClass adds custom CSS classes
func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs templ.Attributes) Option {
	return func(p *Props) {
		if p.Attrs == nil {
			p.Attrs = templ.Attributes{}
		}
		for k, v := range attrs {
			p.Attrs[k] = v
		}
	}
}

// defaultProps returns default icon properties
func defaultProps() *Props {
	return &Props{
		Size:        24,
		Color:       "currentColor",
		StrokeWidth: 2,
	}
}

// writeAttrs writes templ.Attributes as HTML attributes.
func writeAttrs(w io.Writer, attrs templ.Attributes) error {
	for k, v := range attrs {
		switch val := v.(type) {
		case string:
			if val == "" {
				if _, err := fmt.Fprintf(w, " %s", k); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(w, ` %s="%s"`, k, stdhtml.EscapeString(val)); err != nil {
					return err
				}
			}
		case bool:
			if val {
				if _, err := fmt.Fprintf(w, " %s", k); err != nil {
					return err
				}
			}
		default:
			if _, err := fmt.Fprintf(w, ` %s="%s"`, k, stdhtml.EscapeString(fmt.Sprint(val))); err != nil {
				return err
			}
		}
	}
	return nil
}

// Icon creates an icon wrapper around SVG content
// The pathData should be SVG path d attribute content
//
// Example:
//
//	icon := icons.Icon(
//	    "M5 12h14",  // SVG path data
//	    icons.WithSize(20),
//	    icons.WithColor("blue"),
//	)
func Icon(pathData string, opts ...Option) templ.Component {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := "inline-block shrink-0"
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w,
			`<svg class="%s" xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="0 0 24 24" fill="none" stroke="%s" stroke-width="%s" stroke-linecap="round" stroke-linejoin="round"`,
			stdhtml.EscapeString(classes),
			stdhtml.EscapeString(forgeui.ToString(props.Size)),
			stdhtml.EscapeString(forgeui.ToString(props.Size)),
			stdhtml.EscapeString(props.Color),
			stdhtml.EscapeString(forgeui.ToString(props.StrokeWidth)),
		); err != nil {
			return err
		}

		if err := writeAttrs(w, props.Attrs); err != nil {
			return err
		}

		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, `<path d="%s"></path>`, stdhtml.EscapeString(pathData)); err != nil {
			return err
		}

		_, err := io.WriteString(w, "</svg>")
		return err
	})
}

// MultiPathIcon creates an icon with multiple paths
func MultiPathIcon(paths []string, opts ...Option) templ.Component {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := "inline-block shrink-0"
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w,
			`<svg class="%s" xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="0 0 24 24" fill="none" stroke="%s" stroke-width="%s" stroke-linecap="round" stroke-linejoin="round"`,
			stdhtml.EscapeString(classes),
			stdhtml.EscapeString(forgeui.ToString(props.Size)),
			stdhtml.EscapeString(forgeui.ToString(props.Size)),
			stdhtml.EscapeString(props.Color),
			stdhtml.EscapeString(forgeui.ToString(props.StrokeWidth)),
		); err != nil {
			return err
		}

		if err := writeAttrs(w, props.Attrs); err != nil {
			return err
		}

		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}

		for _, pathData := range paths {
			if _, err := fmt.Fprintf(w, `<path d="%s"></path>`, stdhtml.EscapeString(pathData)); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, "</svg>")
		return err
	})
}
