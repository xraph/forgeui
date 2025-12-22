// Package pagination provides pagination components following shadcn/ui patterns.
// Pagination uses Alpine.js for state management and page navigation.
package pagination

import (
	"fmt"
	"strconv"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
)

// PaginationProps defines pagination configuration
type PaginationProps struct {
	CurrentPage   int
	TotalPages    int
	SiblingCount  int  // Pages shown on each side of current page
	ShowFirstLast bool // Show first/last page buttons
	ShowPrevNext  bool // Show previous/next buttons
	OnPageChange  string
	Class         string
	Attrs         []g.Node
}

// Option is a functional option for configuring pagination
type Option func(*PaginationProps)

// WithCurrentPage sets the current active page
func WithCurrentPage(page int) Option {
	return func(p *PaginationProps) { p.CurrentPage = page }
}

// WithTotalPages sets the total number of pages
func WithTotalPages(total int) Option {
	return func(p *PaginationProps) { p.TotalPages = total }
}

// WithSiblingCount sets pages shown around current page
func WithSiblingCount(count int) Option {
	return func(p *PaginationProps) { p.SiblingCount = count }
}

// WithShowFirstLast enables first/last page buttons
func WithShowFirstLast(show bool) Option {
	return func(p *PaginationProps) { p.ShowFirstLast = show }
}

// WithShowPrevNext enables previous/next buttons
func WithShowPrevNext(show bool) Option {
	return func(p *PaginationProps) { p.ShowPrevNext = show }
}

// WithOnPageChange sets the callback for page changes (Alpine expression)
func WithOnPageChange(expr string) Option {
	return func(p *PaginationProps) { p.OnPageChange = expr }
}

// WithClass adds custom classes
func WithClass(class string) Option {
	return func(p *PaginationProps) { p.Class = class }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs ...g.Node) Option {
	return func(p *PaginationProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultProps returns default property values
func defaultProps() *PaginationProps {
	return &PaginationProps{
		CurrentPage:   1,
		TotalPages:    1,
		SiblingCount:  1,
		ShowFirstLast: false,
		ShowPrevNext:  true,
	}
}

// Pagination creates a pagination component
//
// Example:
//
//	pagination.Pagination(
//	    pagination.WithCurrentPage(3),
//	    pagination.WithTotalPages(10),
//	    pagination.WithSiblingCount(1),
//	)
func Pagination(opts ...Option) g.Node {
	props := defaultProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := ""
	if props.Class != "" {
		classes = props.Class
	}

	attrs := []g.Node{
		g.Attr("role", "navigation"),
		g.Attr("aria-label", "Pagination"),
		alpine.XData(map[string]any{
			"currentPage": props.CurrentPage,
			"totalPages":  props.TotalPages,
		}),
	}
	if classes != "" {
		attrs = append(attrs, html.Class(classes))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Nav(
		g.Group(attrs),
		html.Ul(
			html.Class("flex flex-row items-center gap-1"),
			g.Group(buildPaginationItems(props)),
		),
	)
}

// buildPaginationItems generates all pagination items
func buildPaginationItems(props *PaginationProps) []g.Node {
	var items []g.Node

	// First button
	if props.ShowFirstLast {
		items = append(items, FirstButton())
	}

	// Previous button
	if props.ShowPrevNext {
		items = append(items, PrevButton())
	}

	// Page numbers with ellipsis
	items = append(items, buildPageNumbers(props)...)

	// Next button
	if props.ShowPrevNext {
		items = append(items, NextButton())
	}

	// Last button
	if props.ShowFirstLast {
		items = append(items, LastButton())
	}

	return items
}

// buildPageNumbers generates page number buttons with ellipsis
func buildPageNumbers(props *PaginationProps) []g.Node {
	var items []g.Node

	total := props.TotalPages
	sibling := props.SiblingCount

	// Always show first page
	items = append(items, PageButton(1))

	// Calculate range around current page
	// For dynamic current page, we'll show a reasonable range
	// and use Alpine to conditionally render based on currentPage

	if total <= 7 {
		// Show all pages if total is small
		for i := 2; i < total; i++ {
			items = append(items, PageButton(i))
		}
	} else {
		// Show ellipsis pattern
		// [1] ... [4] [5] [6] ... [10]

		// Left ellipsis (shown when currentPage > 3)
		items = append(items, html.Li(
			g.Attr("x-show", "currentPage > "+(strconv.Itoa(2+sibling))),
			html.Class("flex items-center"),
			Ellipsis(),
		))

		// Middle pages (shown based on currentPage)
		for i := 2; i < total; i++ {
			items = append(items, html.Li(
				g.Attr("x-show", fmt.Sprintf(
					"Math.abs(currentPage - %d) <= %d",
					i, sibling,
				)),
				html.Class("flex items-center"),
				PageButton(i),
			))
		}

		// Right ellipsis (shown when currentPage < totalPages - 2)
		items = append(items, html.Li(
			g.Attr("x-show", fmt.Sprintf("currentPage < totalPages - %d", 1+sibling)),
			html.Class("flex items-center"),
			Ellipsis(),
		))
	}

	// Always show last page (if > 1)
	if total > 1 {
		items = append(items, PageButton(total))
	}

	return items
}

// FirstButton creates a "go to first page" button
func FirstButton() g.Node {
	return html.Li(
		html.Button(
			g.Attr("type", "button"),
			g.Attr(":disabled", "currentPage === 1"),
			alpine.XOn("click", "currentPage = 1"),
			html.Class("inline-flex items-center justify-center gap-1 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 px-4"),
			g.Attr("aria-label", "Go to first page"),
			g.Text("First"),
		),
	)
}

// LastButton creates a "go to last page" button
func LastButton() g.Node {
	return html.Li(
		html.Button(
			g.Attr("type", "button"),
			g.Attr(":disabled", "currentPage === totalPages"),
			alpine.XOn("click", "currentPage = totalPages"),
			html.Class("inline-flex items-center justify-center gap-1 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 px-4"),
			g.Attr("aria-label", "Go to last page"),
			g.Text("Last"),
		),
	)
}

// PrevButton creates a "previous page" button
func PrevButton() g.Node {
	return html.Li(
		html.Button(
			g.Attr("type", "button"),
			g.Attr(":disabled", "currentPage === 1"),
			alpine.XOn("click", "if (currentPage > 1) currentPage--"),
			html.Class("inline-flex items-center justify-center gap-1 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 px-4"),
			g.Attr("aria-label", "Go to previous page"),
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
				g.El("path", g.Attr("d", "m15 18-6-6 6-6")),
			),
			g.Text("Previous"),
		),
	)
}

// NextButton creates a "next page" button
func NextButton() g.Node {
	return html.Li(
		html.Button(
			g.Attr("type", "button"),
			g.Attr(":disabled", "currentPage === totalPages"),
			alpine.XOn("click", "if (currentPage < totalPages) currentPage++"),
			html.Class("inline-flex items-center justify-center gap-1 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 px-4"),
			g.Attr("aria-label", "Go to next page"),
			g.Text("Next"),
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
				g.El("path", g.Attr("d", "m9 18 6-6-6-6")),
			),
		),
	)
}

// PageButton creates a page number button
func PageButton(page int) g.Node {
	return html.Button(
		g.Attr("type", "button"),
		alpine.XOn("click", fmt.Sprintf("currentPage = %d", page)),
		g.Attr(":aria-current", fmt.Sprintf("currentPage === %d ? 'page' : undefined", page)),
		html.Class("inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9"),
		g.Attr(":class", fmt.Sprintf("currentPage === %d ? 'border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground' : ''", page)),
		g.Text(strconv.Itoa(page)),
	)
}

// Ellipsis creates an ellipsis indicator for collapsed pages
func Ellipsis() g.Node {
	return html.Span(
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
	)
}
