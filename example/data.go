package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/badge"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/emptystate"
	"github.com/xraph/forgeui/components/list"
	"github.com/xraph/forgeui/components/table"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

// Sample data for tables
var userData = []map[string]any{
	{"name": "John Doe", "email": "john@example.com", "role": "Admin", "status": "active"},
	{"name": "Jane Smith", "email": "jane@example.com", "role": "User", "status": "active"},
	{"name": "Bob Johnson", "email": "bob@example.com", "role": "User", "status": "inactive"},
	{"name": "Alice Williams", "email": "alice@example.com", "role": "Editor", "status": "active"},
	{"name": "Charlie Brown", "email": "charlie@example.com", "role": "User", "status": "active"},
	{"name": "Diana Prince", "email": "diana@example.com", "role": "Admin", "status": "active"},
	{"name": "Eve Adams", "email": "eve@example.com", "role": "Editor", "status": "inactive"},
	{"name": "Frank Miller", "email": "frank@example.com", "role": "User", "status": "active"},
}

func handleData(w http.ResponseWriter, r *http.Request) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Data Components")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			alpine.CloakCSS(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),
			html.StyleEl(g.Raw(`
				@layer base {
					* {
						@apply border-border;
					}
				}
			`)),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground"),
			alpine.XData(map[string]any{"darkMode": false}),
			alpine.XBind("class", `darkMode ? 'dark' : ''`),

			// Header
			html.Div(
				html.Class("border-b"),
				html.Div(
					html.Class("container mx-auto px-4 py-4"),
					html.Div(
						html.Class("flex flex-row items-center justify-between"),
						html.Div(
							html.H1(
								html.Class("text-2xl font-bold"),
								g.Text("ForgeUI - Data Components"),
							),
							html.P(
								html.Class("text-sm text-muted-foreground"),
								g.Text("Tables, lists, and data display components"),
							),
						),
						html.Div(
							html.Class("flex flex-row gap-2"),
							html.A(
								html.Href("/"),
								button.Outline(g.Text("Home")),
							),
							button.Button(
								g.Group([]g.Node{
									alpine.XOn("click", "darkMode = !darkMode"),
									alpine.XText("darkMode ? '☀️' : '🌙'"),
								}),
								button.WithClass("size-9"),
							),
						),
					),
				),
			),

			// Main content
			primitives.Box(
				primitives.WithClass("container mx-auto px-4 py-12"),
				primitives.WithChildren(
					primitives.Box(
						primitives.WithClass("space-y-12"),
						primitives.WithChildren(
							// Basic Table Section
							sectionBox(
								"Basic Table",
								"Simple semantic HTML table with hover effects",
								table.Table()(
									table.TableHeader()(
										table.TableRow()(
											table.TableHeaderCell()(g.Text("Name")),
											table.TableHeaderCell()(g.Text("Email")),
											table.TableHeaderCell()(g.Text("Role")),
											table.TableHeaderCell()(g.Text("Status")),
										),
									),
									table.TableBody()(
										g.Group(renderTableRows(userData[:4])),
									),
								),
							),

							// DataTable with Sorting Section
							sectionBox(
								"Data Table with Sorting",
								"Interactive table with sortable columns",
								table.DataTable(
									table.WithColumns(
										table.Column{Key: "name", Label: "Name", Sortable: true},
										table.Column{Key: "email", Label: "Email", Sortable: true},
										table.Column{Key: "role", Label: "Role", Sortable: true},
										table.Column{Key: "status", Label: "Status", Sortable: true},
									),
									table.WithData(userData),
								),
							),

							// DataTable with Filtering Section
							sectionBox(
								"Data Table with Filtering",
								"Table with column filters and sorting",
								table.DataTable(
									table.WithColumns(
										table.Column{Key: "name", Label: "Name", Sortable: true, Filterable: true},
										table.Column{Key: "email", Label: "Email"},
										table.Column{
											Key:        "role",
											Label:      "Role",
											Filterable: true,
											FilterOptions: []table.FilterOption{
												{Value: "Admin", Label: "Admin"},
												{Value: "User", Label: "User"},
												{Value: "Editor", Label: "Editor"},
											},
										},
										table.Column{
											Key:        "status",
											Label:      "Status",
											Filterable: true,
											FilterOptions: []table.FilterOption{
												{Value: "active", Label: "Active"},
												{Value: "inactive", Label: "Inactive"},
											},
										},
									),
									table.WithData(userData),
								),
							),

							// DataTable with Pagination Section
							sectionBox(
								"Data Table with Pagination",
								"Table with page navigation",
								table.DataTable(
									table.WithColumns(
										table.Column{Key: "name", Label: "Name", Sortable: true},
										table.Column{Key: "email", Label: "Email"},
										table.Column{Key: "role", Label: "Role"},
										table.Column{Key: "status", Label: "Status"},
									),
									table.WithData(userData),
									table.WithPagination(),
									table.WithPageSize(4),
								),
							),

							// Lists Section
							sectionBox(
								"Lists",
								"Various list styles and variants",
								primitives.Box(
									primitives.WithClass("grid grid-cols-1 md:grid-cols-3 gap-6"),
									primitives.WithChildren(
										// Bullet list
										html.Div(
											html.H3(
												html.Class("font-semibold mb-2"),
												g.Text("Bullet List"),
											),
											list.List()(
												list.ListItem()(g.Text("First item")),
												list.ListItem()(g.Text("Second item")),
												list.ListItem()(g.Text("Third item")),
											),
										),

										// Ordered list
										html.Div(
											html.H3(
												html.Class("font-semibold mb-2"),
												g.Text("Ordered List"),
											),
											list.OrderedList()(
												list.ListItem()(g.Text("Step 1")),
												list.ListItem()(g.Text("Step 2")),
												list.ListItem()(g.Text("Step 3")),
											),
										),

										// List with badges
										html.Div(
											html.H3(
												html.Class("font-semibold mb-2"),
												g.Text("List with Badges"),
											),
											list.List(list.None())(
												list.ListItem(list.CardStyle())(
													g.Text("Admin"),
													badge.Badge("Active", badge.WithVariant(forgeui.VariantDefault)),
												),
												list.ListItem(list.CardStyle())(
													g.Text("User"),
													badge.Badge("Pending", badge.WithVariant(forgeui.VariantSecondary)),
												),
												list.ListItem(list.CardStyle())(
													g.Text("Editor"),
													badge.Badge("Blocked", badge.WithVariant(forgeui.VariantDestructive)),
												),
											),
										),
									),
								),
							),

							// Empty State Section
							sectionBox(
								"Empty State",
								"Component for 'no data' scenarios",
								emptystate.EmptyState(
									emptystate.WithIcon(
										html.Div(
											html.Class("w-16 h-16 mx-auto text-muted-foreground/50"),
											g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" /></svg>`),
										),
									),
									emptystate.WithTitle("No data found"),
									emptystate.WithDescription("We couldn't find any items. Try adjusting your filters or add new data."),
									emptystate.WithAction(
										button.Primary(g.Text("Add Item")),
									),
								),
							),
						),
					),
				),
			),

			// Alpine scripts
			alpine.Scripts(),
		),
	)

	_ = page.Render(w)
}

// Helper function to render table rows
func renderTableRows(data []map[string]any) []g.Node {
	rows := make([]g.Node, len(data))
	for i, row := range data {
		rows[i] = table.TableRow()(
			table.TableCell()(g.Text(row["name"].(string))),
			table.TableCell()(g.Text(row["email"].(string))),
			table.TableCell()(g.Text(row["role"].(string))),
			table.TableCell()(renderStatusBadge(row["status"].(string))),
		)
	}
	return rows
}

// Helper function to render status badge
func renderStatusBadge(status string) g.Node {
	variant := forgeui.VariantDefault
	if status == "inactive" {
		variant = forgeui.VariantSecondary
	}
	return badge.Badge(status, badge.WithVariant(variant))
}

// Helper function to create a section box
func sectionBox(title, description string, content g.Node) g.Node {
	return primitives.Box(
		primitives.WithClass("bg-card text-card-foreground rounded-xl border p-6 space-y-4"),
		primitives.WithChildren(
			html.Div(
				html.H2(
					html.Class("text-xl font-semibold"),
					g.Text(title),
				),
				html.P(
					html.Class("text-sm text-muted-foreground"),
					g.Text(description),
				),
			),
			content,
		),
	)
}

// Dashboard data generators

// DashboardMetric represents a metric card data
type DashboardMetric struct {
	Title       string
	Value       string
	Change      string
	ChangeType  string // "up" or "down"
	Description string
}

// GetDashboardMetrics returns sample metrics for the dashboard
func GetDashboardMetrics() []DashboardMetric {
	return []DashboardMetric{
		{
			Title:       "Total Revenue",
			Value:       "$1,250.00",
			Change:      "+12.5%",
			ChangeType:  "up",
			Description: "Trending up this month",
		},
		{
			Title:       "New Customers",
			Value:       "1,234",
			Change:      "-20%",
			ChangeType:  "down",
			Description: "Down 20% this period",
		},
		{
			Title:       "Active Accounts",
			Value:       "45,678",
			Change:      "+12.5%",
			ChangeType:  "up",
			Description: "Strong user retention",
		},
		{
			Title:       "Growth Rate",
			Value:       "4.5%",
			Change:      "+4.5%",
			ChangeType:  "up",
			Description: "Steady performance increase",
		},
	}
}

// ChartDataPoint represents a single data point in a chart
type ChartDataPoint struct {
	Label string
	Value float64
}

// GetVisitorChartData returns time-series data for visitor chart
func GetVisitorChartData() []ChartDataPoint {
	return []ChartDataPoint{
		{Label: "Apr 3", Value: 186},
		{Label: "Apr 9", Value: 305},
		{Label: "Apr 15", Value: 237},
		{Label: "Apr 21", Value: 273},
		{Label: "Apr 27", Value: 209},
		{Label: "May 3", Value: 214},
		{Label: "May 9", Value: 286},
		{Label: "May 15", Value: 305},
		{Label: "May 21", Value: 237},
		{Label: "May 28", Value: 273},
		{Label: "Jun 3", Value: 209},
		{Label: "Jun 9", Value: 314},
		{Label: "Jun 15", Value: 286},
		{Label: "Jun 21", Value: 305},
		{Label: "Jun 29", Value: 237},
	}
}

// GetMonthlyComparisonData returns bar chart data for monthly comparisons
func GetMonthlyComparisonData() []ChartDataPoint {
	return []ChartDataPoint{
		{Label: "Jan", Value: 186},
		{Label: "Feb", Value: 305},
		{Label: "Mar", Value: 237},
		{Label: "Apr", Value: 273},
		{Label: "May", Value: 209},
		{Label: "Jun", Value: 214},
	}
}

// Transaction represents a transaction record
type Transaction struct {
	ID       string
	Customer string
	Email    string
	Amount   string
	Status   string
	Date     string
}

// GetRecentTransactions returns sample transaction data
func GetRecentTransactions() []map[string]any {
	return []map[string]any{
		{
			"id":       "TXN-001",
			"customer": "Olivia Martin",
			"email":    "olivia.martin@email.com",
			"amount":   "$1,999.00",
			"status":   "completed",
			"date":     "2024-06-23",
		},
		{
			"id":       "TXN-002",
			"customer": "Jackson Lee",
			"email":    "jackson.lee@email.com",
			"amount":   "$39.00",
			"status":   "completed",
			"date":     "2024-06-23",
		},
		{
			"id":       "TXN-003",
			"customer": "Isabella Nguyen",
			"email":    "isabella.nguyen@email.com",
			"amount":   "$299.00",
			"status":   "pending",
			"date":     "2024-06-22",
		},
		{
			"id":       "TXN-004",
			"customer": "William Kim",
			"email":    "will@email.com",
			"amount":   "$99.00",
			"status":   "completed",
			"date":     "2024-06-22",
		},
		{
			"id":       "TXN-005",
			"customer": "Sofia Davis",
			"email":    "sofia.davis@email.com",
			"amount":   "$39.00",
			"status":   "failed",
			"date":     "2024-06-21",
		},
		{
			"id":       "TXN-006",
			"customer": "Liam Johnson",
			"email":    "liam@email.com",
			"amount":   "$199.00",
			"status":   "completed",
			"date":     "2024-06-21",
		},
		{
			"id":       "TXN-007",
			"customer": "Emma Wilson",
			"email":    "emma.wilson@email.com",
			"amount":   "$49.00",
			"status":   "completed",
			"date":     "2024-06-20",
		},
		{
			"id":       "TXN-008",
			"customer": "Noah Brown",
			"email":    "noah@email.com",
			"amount":   "$149.00",
			"status":   "pending",
			"date":     "2024-06-20",
		},
	}
}
