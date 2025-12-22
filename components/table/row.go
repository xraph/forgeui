package table

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// RowProps defines the properties for the TableRow component.
type RowProps struct {
	Class     string
	Clickable bool
	OnClick   string
	Attrs     []g.Node
}

// RowOption is a functional option for configuring the TableRow component.
type RowOption func(*RowProps)

// WithRowClass adds additional CSS classes to the table row.
func WithRowClass(class string) RowOption {
	return func(p *RowProps) {
		p.Class = class
	}
}

// ClickableRow makes the row clickable with a cursor pointer.
func ClickableRow() RowOption {
	return func(p *RowProps) {
		p.Clickable = true
	}
}

// WithOnClick adds a click handler to the row.
func WithOnClick(handler string) RowOption {
	return func(p *RowProps) {
		p.OnClick = handler
		p.Clickable = true
	}
}

// WithRowAttr adds custom HTML attributes to the table row.
func WithRowAttr(attrs ...g.Node) RowOption {
	return func(p *RowProps) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// TableRow creates a table row (tr) element.
//
// Example:
//
//	table.TableRow(table.WithOnClick("handleRowClick()"))(
//	    table.TableCell(g.Text("Cell 1")),
//	    table.TableCell(g.Text("Cell 2")),
//	)
func TableRow(opts ...RowOption) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &RowProps{}
		for _, opt := range opts {
			opt(props)
		}

		classes := "border-border border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted"

		if props.Clickable {
			classes += " cursor-pointer"
		}

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{
			g.Attr("data-slot", "table-row"),
			html.Class(classes),
		}

		if props.OnClick != "" {
			attrs = append(attrs, g.Attr("onclick", props.OnClick))
		}

		attrs = append(attrs, props.Attrs...)

		return html.Tr(
			g.Group(attrs),
			g.Group(children),
		)
	}
}
