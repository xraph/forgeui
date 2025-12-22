package table

import (
	"encoding/json"
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
)

// Column represents a column configuration for DataTable.
type Column struct {
	Key           string
	Label         string
	Sortable      bool
	Filterable    bool
	FilterOptions []FilterOption
	Width         string
	Align         Alignment
}

// FilterOption represents a filter option for a column.
type FilterOption struct {
	Value string
	Label string
}

// DataTableProps defines the properties for the DataTable component.
type DataTableProps struct {
	Columns        []Column
	Data           []map[string]any
	PageSize       int
	ShowPagination bool
	Class          string
	Attrs          []g.Node
}

// DataTableOption is a functional option for configuring the DataTable component.
type DataTableOption func(*DataTableProps)

// WithColumns sets the columns for the data table.
func WithColumns(columns ...Column) DataTableOption {
	return func(p *DataTableProps) {
		p.Columns = columns
	}
}

// WithData sets the data for the data table.
func WithData(data []map[string]any) DataTableOption {
	return func(p *DataTableProps) {
		p.Data = data
	}
}

// WithPageSize sets the number of items per page.
func WithPageSize(size int) DataTableOption {
	return func(p *DataTableProps) {
		p.PageSize = size
	}
}

// WithPagination enables pagination for the data table.
func WithPagination() DataTableOption {
	return func(p *DataTableProps) {
		p.ShowPagination = true
	}
}

// WithDataTableClass adds additional CSS classes to the data table container.
func WithDataTableClass(class string) DataTableOption {
	return func(p *DataTableProps) {
		p.Class = class
	}
}

// WithDataTableAttr adds custom HTML attributes to the data table container.
func WithDataTableAttr(attrs ...g.Node) DataTableOption {
	return func(p *DataTableProps) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// defaultDataTableProps returns the default data table properties.
func defaultDataTableProps() *DataTableProps {
	return &DataTableProps{
		PageSize:       10,
		ShowPagination: false,
	}
}

// DataTable creates an interactive data table with sorting, filtering, and pagination.
//
// Example:
//
//	table.DataTable(
//	    table.WithColumns(
//	        table.Column{Key: "name", Label: "Name", Sortable: true},
//	        table.Column{Key: "email", Label: "Email"},
//	        table.Column{
//	            Key: "status",
//	            Label: "Status",
//	            Filterable: true,
//	            FilterOptions: []table.FilterOption{
//	                {Value: "active", Label: "Active"},
//	                {Value: "inactive", Label: "Inactive"},
//	            },
//	        },
//	    ),
//	    table.WithData(userData),
//	    table.WithPagination(),
//	    table.WithPageSize(20),
//	)
func DataTable(opts ...DataTableOption) g.Node {
	props := defaultDataTableProps()
	for _, opt := range opts {
		opt(props)
	}

	// Serialize data to JSON for Alpine
	dataJSON, _ := json.Marshal(props.Data)

	// Build Alpine data with methods
	alpineMethods := map[string]any{
		"data":          string(dataJSON),
		"sortColumn":    "",
		"sortDirection": "asc",
		"filters":       map[string]string{},
		"currentPage":   1,
		"pageSize":      props.PageSize,
		"rawData":       []any{},

		"init": alpine.RawJS(`function() {
			this.rawData = JSON.parse(this.data);
		}`),

		"sortBy": alpine.RawJS(`function(column) {
			if (this.sortColumn === column) {
				this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
			} else {
				this.sortColumn = column;
				this.sortDirection = 'asc';
			}
			this.currentPage = 1;
		}`),

		"setFilter": alpine.RawJS(`function(column, value) {
			if (value === '') {
				delete this.filters[column];
			} else {
				this.filters[column] = value;
			}
			this.currentPage = 1;
		}`),

		"clearFilters": alpine.RawJS(`function() {
			this.filters = {};
			this.currentPage = 1;
		}`),

		"get filteredData": alpine.RawJS(`function() {
			let result = [...this.rawData];
			
			Object.keys(this.filters).forEach(key => {
				const filterValue = this.filters[key];
				if (filterValue) {
					result = result.filter(row => {
						const cellValue = row[key];
						if (typeof cellValue === 'string') {
							return cellValue.toLowerCase().includes(filterValue.toLowerCase());
						}
						return String(cellValue) === String(filterValue);
					});
				}
			});
			
			return result;
		}`),

		"get sortedData": alpine.RawJS(`function() {
			if (!this.sortColumn) {
				return this.filteredData;
			}
			
			return [...this.filteredData].sort((a, b) => {
				const aVal = a[this.sortColumn];
				const bVal = b[this.sortColumn];
				
				if (aVal === bVal) return 0;
				
				const comparison = aVal < bVal ? -1 : 1;
				return this.sortDirection === 'asc' ? comparison : -comparison;
			});
		}`),

		"get paginatedData": alpine.RawJS(`function() {
			if (!this.showPagination) {
				return this.sortedData;
			}
			
			const start = (this.currentPage - 1) * this.pageSize;
			const end = start + this.pageSize;
			return this.sortedData.slice(start, end);
		}`),

		"get totalPages": alpine.RawJS(`function() {
			return Math.ceil(this.sortedData.length / this.pageSize);
		}`),

		"goToPage": alpine.RawJS(`function(page) {
			if (page >= 1 && page <= this.totalPages) {
				this.currentPage = page;
			}
		}`),

		"showPagination": props.ShowPagination,
	}

	classes := "w-full"
	if props.Class != "" {
		classes += " " + props.Class
	}

	containerAttrs := []g.Node{
		html.Class(classes),
		alpine.XData(alpineMethods),
	}
	containerAttrs = append(containerAttrs, props.Attrs...)

	return html.Div(
		g.Group(containerAttrs),

		// Filters (if any columns are filterable)
		g.If(hasFilterableColumns(props.Columns),
			html.Div(
				html.Class("flex items-center gap-2 mb-4"),
				g.Group(renderFilters(props.Columns)),
				// Clear filters button
				html.Button(
					html.Type("button"),
					html.Class("inline-flex items-center justify-center rounded-md text-sm font-medium h-9 px-4 py-2 border-border border bg-background hover:bg-accent hover:text-accent-foreground"),
					alpine.XOn("click", "clearFilters()"),
					alpine.XShow("Object.keys(filters).length > 0"),
					g.Text("Clear Filters"),
				),
			),
		),

		// Table
		Table()(
			// Header
			TableHeader()(
				TableRow()(
					g.Group(renderDataTableHeaders(props.Columns)),
				),
			),

			// Body
			TableBody()(
				// Use Alpine x-for to render rows
				g.El("template",
					g.Group(alpine.XForKeyed("row in paginatedData", "row")),
					g.Group([]g.Node{
						TableRow()(
							g.Group(renderDataTableCells(props.Columns)),
						),
					}),
				),
			),
		),

		// Pagination (if enabled)
		g.If(props.ShowPagination,
			html.Div(
				html.Class("flex items-center justify-between px-2 py-4"),
				html.Div(
					html.Class("text-sm text-muted-foreground"),
					alpine.XText("'Showing ' + ((currentPage - 1) * pageSize + 1) + ' to ' + Math.min(currentPage * pageSize, sortedData.length) + ' of ' + sortedData.length + ' results'"),
				),
				html.Div(
					html.Class("flex items-center gap-2"),
					// Previous button
					html.Button(
						html.Type("button"),
						html.Class("inline-flex items-center justify-center rounded-md text-sm font-medium h-9 px-4 py-2 border-border border bg-background hover:bg-accent hover:text-accent-foreground disabled:opacity-50 disabled:pointer-events-none"),
						alpine.XBind("disabled", "currentPage === 1"),
						alpine.XOn("click", "goToPage(currentPage - 1)"),
						g.Text("Previous"),
					),
					// Page info
					html.Span(
						html.Class("text-sm"),
						alpine.XText("'Page ' + currentPage + ' of ' + totalPages"),
					),
					// Next button
					html.Button(
						html.Type("button"),
						html.Class("inline-flex items-center justify-center rounded-md text-sm font-medium h-9 px-4 py-2 border bg-background hover:bg-accent hover:text-accent-foreground disabled:opacity-50 disabled:pointer-events-none"),
						alpine.XBind("disabled", "currentPage === totalPages"),
						alpine.XOn("click", "goToPage(currentPage + 1)"),
						g.Text("Next"),
					),
				),
			),
		),
	)
}

// renderDataTableHeaders generates the header cells for the data table.
func renderDataTableHeaders(columns []Column) []g.Node {
	headers := make([]g.Node, len(columns))

	for i, col := range columns {
		cellOpts := []CellOption{}

		if col.Width != "" {
			cellOpts = append(cellOpts, WithWidth(col.Width))
		}
		if col.Align != "" {
			cellOpts = append(cellOpts, WithAlign(col.Align))
		}

		var content []g.Node

		if col.Sortable {
			// Sortable header with click handler
			content = []g.Node{
				html.Button(
					html.Type("button"),
					html.Class("flex items-center gap-2 font-medium hover:text-foreground"),
					alpine.XOn("click", fmt.Sprintf("sortBy('%s')", col.Key)),
					g.Text(col.Label),
					// Sort indicator
					html.Span(
						html.Class("ml-1"),
						alpine.XShow(fmt.Sprintf("sortColumn === '%s'", col.Key)),
						alpine.XText(fmt.Sprintf("sortDirection === 'asc' ? '↑' : '↓'")),
					),
				),
			}
		} else {
			content = []g.Node{g.Text(col.Label)}
		}

		headers[i] = TableHeaderCell(cellOpts...)(g.Group(content))
	}

	return headers
}

// renderDataTableCells generates the data cells for the data table rows.
func renderDataTableCells(columns []Column) []g.Node {
	cells := make([]g.Node, len(columns))

	for i, col := range columns {
		cellOpts := []CellOption{}

		if col.Width != "" {
			cellOpts = append(cellOpts, WithWidth(col.Width))
		}
		if col.Align != "" {
			cellOpts = append(cellOpts, WithAlign(col.Align))
		}

		cells[i] = TableCell(cellOpts...)(
			alpine.XText(fmt.Sprintf("row.%s", col.Key)),
		)
	}

	return cells
}

// hasFilterableColumns checks if any columns are filterable.
func hasFilterableColumns(columns []Column) bool {
	for _, col := range columns {
		if col.Filterable {
			return true
		}
	}
	return false
}

// renderFilters generates filter controls for filterable columns.
func renderFilters(columns []Column) []g.Node {
	var filters []g.Node

	for _, col := range columns {
		if !col.Filterable {
			continue
		}

		if len(col.FilterOptions) > 0 {
			// Dropdown filter for columns with predefined options
			filters = append(filters,
				html.Div(
					html.Class("flex flex-col gap-1"),
					html.Label(
						html.Class("text-sm font-medium"),
						g.Text(col.Label),
					),
					html.Select(
						html.Class("h-9 rounded-md border-border border border-input bg-background px-3 py-1 text-sm"),
						alpine.XModel(fmt.Sprintf("filters['%s']", col.Key)),
						html.Option(
							html.Value(""),
							g.Text("All"),
						),
						g.Group(renderFilterOptions(col.FilterOptions)),
					),
				),
			)
		} else {
			// Text input filter for free-form filtering
			filters = append(filters,
				html.Div(
					html.Class("flex flex-col gap-1"),
					html.Label(
						html.Class("text-sm font-medium"),
						g.Text(col.Label),
					),
					html.Input(
						html.Type("text"),
						html.Class("h-9 rounded-md border-border border border-input bg-background px-3 py-1 text-sm"),
						html.Placeholder(fmt.Sprintf("Filter %s...", col.Label)),
						alpine.XModel(fmt.Sprintf("filters['%s']", col.Key)),
					),
				),
			)
		}
	}

	return filters
}

// renderFilterOptions generates option elements for filter dropdowns.
func renderFilterOptions(options []FilterOption) []g.Node {
	result := make([]g.Node, len(options))
	for i, opt := range options {
		result[i] = html.Option(
			html.Value(opt.Value),
			g.Text(opt.Label),
		)
	}
	return result
}
