package table

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// BodyProps defines the properties for the TableBody component.
type BodyProps struct {
	Class string
	Attrs []g.Node
}

// BodyOption is a functional option for configuring the TableBody component.
type BodyOption func(*BodyProps)

// WithBodyClass adds additional CSS classes to the table body.
func WithBodyClass(class string) BodyOption {
	return func(p *BodyProps) {
		p.Class = class
	}
}

// WithBodyAttr adds custom HTML attributes to the table body.
func WithBodyAttr(attrs ...g.Node) BodyOption {
	return func(p *BodyProps) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// TableBody creates a table body (tbody) element.
//
// Example:
//
//	table.TableBody()(
//	    table.TableRow()(
//	        table.TableCell(g.Text("Data 1")),
//	        table.TableCell(g.Text("Data 2")),
//	    ),
//	)
func TableBody(opts ...BodyOption) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &BodyProps{}
		for _, opt := range opts {
			opt(props)
		}

		classes := "[&_tr:last-child]:border-0 border-border"

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{
			g.Attr("data-slot", "table-body"),
			html.Class(classes),
		}
		attrs = append(attrs, props.Attrs...)

		return html.TBody(
			g.Group(attrs),
			g.Group(children),
		)
	}
}
