// Package table provides table components for displaying tabular data.
//
// It includes both basic semantic HTML table elements and an interactive
// DataTable component with sorting, filtering, and pagination capabilities.
//
// # Basic Usage
//
//	table.Table(
//	    table.TableHeader(
//	        table.TableRow(
//	            table.TableHeaderCell(g.Text("Name")),
//	            table.TableHeaderCell(g.Text("Email")),
//	        ),
//	    ),
//	    table.TableBody(
//	        table.TableRow(
//	            table.TableCell(g.Text("John Doe")),
//	            table.TableCell(g.Text("john@example.com")),
//	        ),
//	    ),
//	)
package table

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Props defines the properties for the Table component.
type Props struct {
	Class    string
	Striped  bool
	Bordered bool
	Hover    bool
	Attrs    []g.Node
}

// Option is a functional option for configuring the Table component.
type Option func(*Props)

// WithClass adds additional CSS classes to the table.
func WithClass(class string) Option {
	return func(p *Props) {
		p.Class = class
	}
}

// Striped applies alternating row background colors.
func Striped() Option {
	return func(p *Props) {
		p.Striped = true
	}
}

// Bordered adds borders to all table cells.
func Bordered() Option {
	return func(p *Props) {
		p.Bordered = true
	}
}

// Hover enables hover effect on table rows.
func Hover() Option {
	return func(p *Props) {
		p.Hover = true
	}
}

// WithAttr adds custom HTML attributes to the table.
func WithAttr(attrs ...g.Node) Option {
	return func(p *Props) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// defaultProps returns the default table properties.
func defaultProps() *Props {
	return &Props{
		Hover: true, // Enable hover by default
	}
}

// Table creates a table component with responsive container.
//
// The table uses semantic HTML elements with proper accessibility attributes.
// It's wrapped in a responsive container that enables horizontal scrolling
// on mobile devices.
//
// Example:
//
//	table.Table(
//	    table.TableHeader(...),
//	    table.TableBody(...),
//	    table.WithClass("my-custom-table"),
//	)
func Table(opts ...Option) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := defaultProps()
		for _, opt := range opts {
			opt(props)
		}

		classes := "w-full caption-bottom text-sm"

		if props.Class != "" {
			classes += " " + props.Class
		}

		tableAttrs := []g.Node{
			html.Class(classes),
			g.Attr("role", "table"),
		}
		tableAttrs = append(tableAttrs, props.Attrs...)

		// Responsive wrapper
		return html.Div(
			g.Attr("data-slot", "table-container"),
			html.Class("relative w-full overflow-x-auto"),
			html.Table(
				g.Attr("data-slot", "table"),
				g.Group(tableAttrs),
				g.Group(children),
			),
		)
	}
}

// TableFooter creates a table footer (tfoot) element.
//
// Example:
//
//	table.TableFooter()(
//	    table.TableRow(
//	        table.TableCell(g.Text("Total")),
//	        table.TableCell(g.Text("$100")),
//	    ),
//	)
func TableFooter(opts ...Option) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &Props{}
		for _, opt := range opts {
			opt(props)
		}

		classes := "bg-muted/50 border-border border-t font-medium [&>tr]:last:border-b-0"

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{
			g.Attr("data-slot", "table-footer"),
			html.Class(classes),
		}
		attrs = append(attrs, props.Attrs...)

		return html.TFoot(
			g.Group(attrs),
			g.Group(children),
		)
	}
}

// TableCaption creates a table caption element.
//
// Example:
//
//	table.TableCaption()(g.Text("A list of your recent invoices."))
func TableCaption(opts ...Option) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &Props{}
		for _, opt := range opts {
			opt(props)
		}

		classes := "text-muted-foreground mt-4 text-sm"

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{
			g.Attr("data-slot", "table-caption"),
			html.Class(classes),
		}
		attrs = append(attrs, props.Attrs...)

		return html.Caption(
			g.Group(attrs),
			g.Group(children),
		)
	}
}
