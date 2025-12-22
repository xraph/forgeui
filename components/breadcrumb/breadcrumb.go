// Package breadcrumb provides breadcrumb navigation components following shadcn/ui patterns.
// Breadcrumbs show the current page's location within a navigational hierarchy.
package breadcrumb

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// BreadcrumbProps defines breadcrumb configuration
type BreadcrumbProps struct {
	Separator g.Node
	MaxItems  int // 0 means no limit
	Class     string
	Attrs     []g.Node
}

// Option is a functional option for configuring breadcrumbs
type Option func(*BreadcrumbProps)

// WithSeparator sets custom separator (default is chevron)
func WithSeparator(sep g.Node) Option {
	return func(p *BreadcrumbProps) { p.Separator = sep }
}

// WithMaxItems sets maximum items to display (middle items collapse to ellipsis)
func WithMaxItems(maxItems int) Option {
	return func(p *BreadcrumbProps) { p.MaxItems = maxItems }
}

// WithClass adds custom classes
func WithClass(class string) Option {
	return func(p *BreadcrumbProps) { p.Class = class }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs ...g.Node) Option {
	return func(p *BreadcrumbProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultProps returns default property values
func defaultProps() *BreadcrumbProps {
	return &BreadcrumbProps{
		Separator: chevronIcon(),
		MaxItems:  0,
	}
}

// chevronIcon returns the default chevron separator
func chevronIcon() g.Node {
	return g.El("svg",
		html.Class("h-4 w-4"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("width", "24"),
		g.Attr("height", "24"),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", "currentColor"),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
		g.El("path", g.Attr("d", "m9 18 6-6-6-6")),
	)
}

// Breadcrumb creates a breadcrumb navigation container
//
// Example:
//
//	breadcrumb.Breadcrumb(
//	    breadcrumb.Item("/", g.Text("Home")),
//	    breadcrumb.Item("/docs", g.Text("Docs")),
//	    breadcrumb.Page(g.Text("Components")),
//	)
func Breadcrumb(children ...g.Node) g.Node {
	return BreadcrumbWithOptions(nil, children...)
}

// BreadcrumbWithOptions creates breadcrumb with custom options
func BreadcrumbWithOptions(opts []Option, children ...g.Node) g.Node {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := ""
	if props.Class != "" {
		classes = props.Class
	}

	attrs := []g.Node{
		g.Attr("aria-label", "Breadcrumb"),
	}
	if classes != "" {
		attrs = append(attrs, html.Class(classes))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Nav(
		g.Group(attrs),
		html.Ol(
			html.Class("flex flex-wrap items-center gap-1.5 break-words text-sm text-muted-foreground sm:gap-2.5"),
			g.Group(children),
		),
	)
}

// Item creates a breadcrumb item with a link
//
// Example:
//
//	breadcrumb.Item("/docs", g.Text("Documentation"))
func Item(href string, label g.Node) g.Node {
	return html.Li(
		html.Class("inline-flex items-center gap-1.5"),
		html.A(
			g.Attr("href", href),
			html.Class("transition-colors hover:text-foreground"),
			label,
		),
		Separator(),
	)
}

// Link creates a breadcrumb link without the separator (for custom layouts)
//
// Example:
//
//	breadcrumb.Link("/docs", g.Text("Documentation"))
func Link(href string, label g.Node) g.Node {
	return html.A(
		g.Attr("href", href),
		html.Class("transition-colors hover:text-foreground"),
		label,
	)
}

// Page creates the current page breadcrumb (not a link)
//
// Example:
//
//	breadcrumb.Page(g.Text("Current Page"))
func Page(label g.Node) g.Node {
	return html.Li(
		html.Class("inline-flex items-center gap-1.5"),
		html.Span(
			g.Attr("role", "link"),
			g.Attr("aria-disabled", "true"),
			g.Attr("aria-current", "page"),
			html.Class("font-normal text-foreground"),
			label,
		),
	)
}

// Separator creates a breadcrumb separator
//
// Example:
//
//	breadcrumb.Separator()
func Separator() g.Node {
	return SeparatorWithIcon(chevronIcon())
}

// SeparatorWithIcon creates a separator with custom icon
//
// Example:
//
//	breadcrumb.SeparatorWithIcon(g.Text("/"))
func SeparatorWithIcon(icon g.Node) g.Node {
	return html.Li(
		g.Attr("role", "presentation"),
		g.Attr("aria-hidden", "true"),
		html.Class("inline-flex items-center"),
		icon,
	)
}

// Ellipsis creates a breadcrumb ellipsis for collapsed items
//
// Example:
//
//	breadcrumb.Ellipsis()
func Ellipsis() g.Node {
	return html.Li(
		html.Class("inline-flex items-center gap-1.5"),
		html.Span(
			g.Attr("role", "presentation"),
			g.Attr("aria-hidden", "true"),
			html.Class("flex h-9 w-9 items-center justify-center"),
			g.El("svg",
				html.Class("h-4 w-4"),
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("width", "24"),
				g.Attr("height", "24"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("fill", "none"),
				g.Attr("stroke", "currentColor"),
				g.Attr("stroke-width", "2"),
				g.Attr("stroke-linecap", "round"),
				g.Attr("stroke-linejoin", "round"),
				g.El("circle", g.Attr("cx", "12"), g.Attr("cy", "12"), g.Attr("r", "1")),
				g.El("circle", g.Attr("cx", "19"), g.Attr("cy", "12"), g.Attr("r", "1")),
				g.El("circle", g.Attr("cx", "5"), g.Attr("cy", "12"), g.Attr("r", "1")),
			),
			g.Text("More"),
		),
		Separator(),
	)
}
