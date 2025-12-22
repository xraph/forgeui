package table

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// HeaderProps defines the properties for the TableHeader component.
type HeaderProps struct {
	Class  string
	Sticky bool
	Attrs  []g.Node
}

// HeaderOption is a functional option for configuring the TableHeader component.
type HeaderOption func(*HeaderProps)

// WithHeaderClass adds additional CSS classes to the table header.
func WithHeaderClass(class string) HeaderOption {
	return func(p *HeaderProps) {
		p.Class = class
	}
}

// StickyHeader makes the table header sticky on scroll.
func StickyHeader() HeaderOption {
	return func(p *HeaderProps) {
		p.Sticky = true
	}
}

// WithHeaderAttr adds custom HTML attributes to the table header.
func WithHeaderAttr(attrs ...g.Node) HeaderOption {
	return func(p *HeaderProps) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// TableHeader creates a table header (thead) element.
//
// Example:
//
//	table.TableHeader(
//	    table.TableRow(
//	        table.TableHeaderCell(g.Text("Column 1")),
//	        table.TableHeaderCell(g.Text("Column 2")),
//	    ),
//	    table.StickyHeader(),
//	)
func TableHeader(opts ...HeaderOption) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &HeaderProps{}
		for _, opt := range opts {
			opt(props)
		}

		classes := "[&_tr]:border-b border-border"

		if props.Sticky {
			classes += " sticky top-0 z-10"
		}

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{
			g.Attr("data-slot", "table-header"),
			html.Class(classes),
		}
		attrs = append(attrs, props.Attrs...)

		return html.THead(
			g.Group(attrs),
			g.Group(children),
		)
	}
}
