// Package icons provides a flexible icon system with Lucide icon integration.
// Icons are rendered as inline SVG elements with customizable size, color, and stroke width.
//
//go:generate go run internal/generate/main.go
package icons

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

// Props defines icon configuration
type Props struct {
	Size        int     // Size in pixels (width and height)
	Color       string  // CSS color value
	StrokeWidth float64 // SVG stroke width
	Class       string  // Additional CSS classes
	Attrs       []g.Node
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
func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultProps returns default icon properties
func defaultProps() *Props {
	return &Props{
		Size:        24,
		Color:       "currentColor",
		StrokeWidth: 2,
	}
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
func Icon(pathData string, opts ...Option) g.Node {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := "inline-block shrink-0"
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("width", forgeui.ToString(props.Size)),
		g.Attr("height", forgeui.ToString(props.Size)),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", props.Color),
		g.Attr("stroke-width", forgeui.ToString(props.StrokeWidth)),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
	}

	attrs = append(attrs, props.Attrs...)

	return g.El("svg",
		g.Group(attrs),
		g.El("path", g.Attr("d", pathData)),
	)
}

// MultiPathIcon creates an icon with multiple paths
func MultiPathIcon(paths []string, opts ...Option) g.Node {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := "inline-block shrink-0"
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("width", forgeui.ToString(props.Size)),
		g.Attr("height", forgeui.ToString(props.Size)),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", props.Color),
		g.Attr("stroke-width", forgeui.ToString(props.StrokeWidth)),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
	}

	attrs = append(attrs, props.Attrs...)

	// Create path elements
	pathNodes := make([]g.Node, len(paths))
	for i, pathData := range paths {
		pathNodes[i] = g.El("path", g.Attr("d", pathData))
	}

	return g.El("svg",
		g.Group(attrs),
		g.Group(pathNodes),
	)
}
