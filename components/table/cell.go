package table

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// Alignment represents the text alignment in a table cell.
type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)

// CellProps defines the properties for table cell components.
type CellProps struct {
	Class string
	Align Alignment
	Width string
	Attrs []g.Node
}

// CellOption is a functional option for configuring table cells.
type CellOption func(*CellProps)

// WithCellClass adds additional CSS classes to the cell.
func WithCellClass(class string) CellOption {
	return func(p *CellProps) {
		p.Class = class
	}
}

// WithAlign sets the text alignment for the cell.
func WithAlign(align Alignment) CellOption {
	return func(p *CellProps) {
		p.Align = align
	}
}

// WithWidth sets the width of the cell.
func WithWidth(width string) CellOption {
	return func(p *CellProps) {
		p.Width = width
	}
}

// WithCellAttr adds custom HTML attributes to the cell.
func WithCellAttr(attrs ...g.Node) CellOption {
	return func(p *CellProps) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// TableCell creates a table data cell (td) element.
//
// Example:
//
//	table.TableCell(
//	    table.WithAlign(table.AlignRight),
//	    table.WithWidth("200px"),
//	)(g.Text("Cell content"))
func TableCell(opts ...CellOption) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &CellProps{
			Align: AlignLeft, // Default alignment
		}
		for _, opt := range opts {
			opt(props)
		}

		classes := "p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0"

		// Add alignment class
		switch props.Align {
		case AlignCenter:
			classes += " text-center"
		case AlignRight:
			classes += " text-right"
		}

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{
			g.Attr("data-slot", "table-cell"),
			html.Class(classes),
		}

		if props.Width != "" {
			attrs = append(attrs, g.Attr("style", "width: "+props.Width))
		}

		attrs = append(attrs, props.Attrs...)

		return html.Td(
			g.Group(attrs),
			g.Group(children),
		)
	}
}

// TableHeaderCell creates a table header cell (th) element.
//
// Example:
//
//	table.TableHeaderCell(
//	    table.WithAlign(table.AlignCenter),
//	)(g.Text("Column Header"))
func TableHeaderCell(opts ...CellOption) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &CellProps{
			Align: AlignLeft, // Default alignment
		}
		for _, opt := range opts {
			opt(props)
		}

		classes := "h-10 px-2 text-left align-middle font-medium text-foreground whitespace-nowrap [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]"

		// Add alignment class
		switch props.Align {
		case AlignCenter:
			classes += " text-center"
		case AlignRight:
			classes += " text-right"
		}

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{
			g.Attr("data-slot", "table-head"),
			html.Class(classes),
		}

		if props.Width != "" {
			attrs = append(attrs, g.Attr("style", "width: "+props.Width))
		}

		attrs = append(attrs, props.Attrs...)

		return html.Th(
			g.Group(attrs),
			g.Group(children),
		)
	}
}

