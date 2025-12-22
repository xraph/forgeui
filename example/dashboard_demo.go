package main

import (
	"fmt"
	"net/http"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/breadcrumb"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/components/separator"
	"github.com/xraph/forgeui/components/sidebar"
	"github.com/xraph/forgeui/components/table"
	"github.com/xraph/forgeui/icons"
	"github.com/xraph/forgeui/theme"
)

// handleDashboard renders the comprehensive dashboard demo
func handleDashboard(w http.ResponseWriter, r *http.Request) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Dashboard")),
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

			// Sidebar with all patterns
			dashboardSidebar(),

			// Main content area
			sidebar.SidebarInset(
				// Header with breadcrumb
				sidebar.SidebarInsetHeader(
					sidebar.SidebarTriggerDesktop(),
					separator.Separator(separator.Vertical(), separator.WithClass("mr-2 h-4")),
					breadcrumb.Breadcrumb(
						breadcrumb.Item("/", g.Text("Dashboard")),
						breadcrumb.Item("/analytics", g.Text("Analytics")),
						breadcrumb.Page(g.Text("Overview")),
					),
					// Theme toggle in header
					html.Div(
						html.Class("ml-auto"),
						button.Button(
							g.Group([]g.Node{
								alpine.XOn("click", "darkMode = !darkMode"),
								alpine.XText("darkMode ? '☀️' : '🌙'"),
							}),
							button.WithVariant(forgeui.VariantGhost),
							button.WithSize(forgeui.SizeIcon),
						),
					),
				),

				// Main content
				html.Main(
					html.Class("flex flex-1 flex-col gap-4 p-4 md:p-6"),

					// Metrics cards
					html.Div(
						html.Class("grid gap-4 md:grid-cols-2 lg:grid-cols-4"),
						g.Group(renderMetricCards()),
					),

					// Charts section
					html.Div(
						html.Class("grid gap-4 md:grid-cols-2 lg:grid-cols-7"),

						// Visitor chart (larger)
						html.Div(
							html.Class("lg:col-span-4"),
							card.Card(
								card.Header(
									card.Title("Total Visitors"),
									card.Description("Total for the last 3 months"),
								),
								card.Content(
									html.Div(
										html.Class("h-[300px] flex items-center justify-center text-muted-foreground"),
										g.Text("Chart visualization would go here"),
										html.Div(
											html.Class("text-xs mt-2"),
											g.Text("(Using plugins/charts for actual implementation)"),
										),
									),
								),
							),
						),

						// Recent activity (smaller)
						html.Div(
							html.Class("lg:col-span-3"),
							card.Card(
								card.Header(
									card.Title("Recent Activity"),
									card.Description("Latest system events"),
								),
								card.Content(
									html.Div(
										html.Class("space-y-4"),
										activityItem("New user registered", "2 minutes ago", icons.User()),
										activityItem("Payment received", "15 minutes ago", icons.DollarSign()),
										activityItem("New order placed", "1 hour ago", icons.Box()),
										activityItem("System updated", "3 hours ago", icons.Settings()),
									),
								),
							),
						),
					),

					// Transactions table
					card.Card(
						card.Header(
							card.Title("Recent Transactions"),
							card.Description("A list of recent transactions from your store"),
						),
						card.Content(
							table.DataTable(
								table.WithColumns(
									table.Column{Key: "id", Label: "Transaction ID", Sortable: true},
									table.Column{Key: "customer", Label: "Customer", Sortable: true},
									table.Column{Key: "email", Label: "Email"},
									table.Column{Key: "amount", Label: "Amount", Sortable: true},
									table.Column{Key: "status", Label: "Status", Sortable: true, Filterable: true, FilterOptions: []table.FilterOption{
										{Value: "completed", Label: "Completed"},
										{Value: "pending", Label: "Pending"},
										{Value: "failed", Label: "Failed"},
									}},
									table.Column{Key: "date", Label: "Date", Sortable: true},
								),
								table.WithData(GetRecentTransactions()),
								table.WithPagination(),
								table.WithPageSize(5),
							),
						),
					),
				),
			),

			// Alpine scripts with Collapse plugin
			alpine.Scripts(alpine.PluginCollapse),
		),
	)

	_ = page.Render(w)
}

// dashboardSidebar creates the enhanced sidebar with all patterns
func dashboardSidebar() g.Node {
	return sidebar.SidebarWithOptions(
		[]sidebar.SidebarOption{
			sidebar.WithDefaultCollapsed(false),
			sidebar.WithCollapsible(true),
			sidebar.WithCollapsibleMode(sidebar.CollapsibleIcon), // Use icon mode to show icons when collapsed
		},
		// Header with logo - icon always visible, text hidden when collapsed
		sidebar.SidebarHeader(
			html.Div(
				html.Class("flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground shrink-0"),
				icons.LayoutDashboard(icons.WithSize(20)),
			),
			html.Span(
				html.Class("font-bold text-lg"),
				g.Attr("x-data", "{}"),
				g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
				g.Text("Acme Inc."),
			),
		),

		// Main content
		sidebar.SidebarContent(
			// Platform section
			sidebar.SidebarGroup(
				sidebar.SidebarGroupLabel("Platform"),
				sidebar.SidebarMenu(
					sidebar.SidebarMenuItem(
						sidebar.SidebarMenuButton(
							"Dashboard",
							sidebar.WithMenuHref("/dashboard"),
							sidebar.WithMenuIcon(icons.LayoutDashboard(icons.WithSize(20))),
							sidebar.WithMenuActive(),
						),
					),
					sidebar.SidebarMenuItem(
						sidebar.SidebarMenuButton(
							"Analytics",
							sidebar.WithMenuHref("/analytics"),
							sidebar.WithMenuIcon(icons.Activity(icons.WithSize(20))),
						),
					),
					sidebar.SidebarMenuItem(
						sidebar.SidebarMenuButton(
							"Customers",
							sidebar.WithMenuHref("/customers"),
							sidebar.WithMenuIcon(icons.Users(icons.WithSize(20))),
							sidebar.WithMenuBadge(sidebar.SidebarMenuBadge("12")),
						),
					),
					sidebar.SidebarMenuItem(
						sidebar.SidebarMenuButton(
							"Orders",
							sidebar.WithMenuHref("/orders"),
							sidebar.WithMenuIcon(icons.Box(icons.WithSize(20))),
							sidebar.WithMenuBadge(sidebar.SidebarMenuBadge("3")),
						),
					),
				),
			),

			// Projects section (collapsible)
			sidebar.SidebarGroupCollapsible(
				[]sidebar.SidebarGroupOption{
					sidebar.WithGroupKey("projects"),
					sidebar.WithGroupDefaultOpen(true),
				},
				sidebar.SidebarGroupLabelCollapsible("projects", "Projects", icons.FolderKanban(icons.WithSize(16))),
				sidebar.SidebarGroupContent("projects",
					sidebar.SidebarMenu(
						sidebar.SidebarMenuItem(
							sidebar.SidebarMenuButton(
								"Design Engineering",
								sidebar.WithMenuHref("/projects/design"),
								sidebar.WithMenuIcon(icons.FolderKanban(icons.WithSize(20))),
							),
							sidebar.SidebarMenuAction(
								icons.EllipsisVertical(icons.WithSize(16)),
								"More options",
							),
							// Submenu
							sidebar.SidebarMenuSub(
								sidebar.SidebarMenuSubItem(
									sidebar.SidebarMenuSubButton("Overview", "/projects/design/overview", false),
								),
								sidebar.SidebarMenuSubItem(
									sidebar.SidebarMenuSubButton("Tasks", "/projects/design/tasks", false),
								),
								sidebar.SidebarMenuSubItem(
									sidebar.SidebarMenuSubButton("Settings", "/projects/design/settings", false),
								),
							),
						),
						sidebar.SidebarMenuItem(
							sidebar.SidebarMenuButton(
								"Sales & Marketing",
								sidebar.WithMenuHref("/projects/sales"),
								sidebar.WithMenuIcon(icons.FolderKanban(icons.WithSize(20))),
							),
							sidebar.SidebarMenuAction(
								icons.EllipsisVertical(icons.WithSize(16)),
								"More options",
							),
						),
						sidebar.SidebarMenuItem(
							sidebar.SidebarMenuButton(
								"Travel",
								sidebar.WithMenuHref("/projects/travel"),
								sidebar.WithMenuIcon(icons.FolderKanban(icons.WithSize(20))),
							),
							sidebar.SidebarMenuAction(
								icons.EllipsisVertical(icons.WithSize(16)),
								"More options",
							),
						),
					),
				),
			),

			// Settings section
			sidebar.SidebarGroup(
				sidebar.SidebarGroupLabel("Settings"),
				sidebar.SidebarMenu(
					sidebar.SidebarMenuItem(
						sidebar.SidebarMenuButton(
							"Account",
							sidebar.WithMenuHref("/settings/account"),
							sidebar.WithMenuIcon(icons.User(icons.WithSize(20))),
						),
					),
					sidebar.SidebarMenuItem(
						sidebar.SidebarMenuButton(
							"Billing",
							sidebar.WithMenuHref("/settings/billing"),
							sidebar.WithMenuIcon(icons.CreditCard(icons.WithSize(20))),
						),
					),
					sidebar.SidebarMenuItem(
						sidebar.SidebarMenuButton(
							"Notifications",
							sidebar.WithMenuHref("/settings/notifications"),
							sidebar.WithMenuIcon(icons.Bell(icons.WithSize(20))),
							sidebar.WithMenuBadge(sidebar.SidebarMenuBadge("5")),
						),
					),
				),
			),
		),

		// Footer with user profile
		sidebar.SidebarFooter(
			sidebar.SidebarMenu(
				sidebar.SidebarMenuItem(
					// User dropdown
					html.Div(
						alpine.XData(map[string]any{"userMenuOpen": false}),
						html.Class("relative"),
						sidebar.SidebarMenuButton(
							"John Doe",
							sidebar.WithMenuIcon(
								html.Div(
									html.Class("flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground"),
									g.Text("JD"),
								),
							),
							sidebar.WithMenuAsButton(),
							sidebar.WithMenuAttrs(
								alpine.XOn("click", "userMenuOpen = !userMenuOpen"),
							),
						),
						// Dropdown menu
						html.Div(
							g.Attr("x-show", "userMenuOpen"),
							alpine.XOn("click.away", "userMenuOpen = false"),
							g.Attr("x-transition", ""),
							html.Class("absolute bottom-full left-0 mb-2 w-56 rounded-md border bg-popover p-1 shadow-md"),
							html.Div(
								html.Class("px-2 py-1.5 text-sm font-semibold"),
								g.Text("john@example.com"),
							),
							separator.Separator(separator.WithClass("my-1")),
							dropdownMenuItem("/profile", "Profile", icons.User(icons.WithSize(16))),
							dropdownMenuItem("/settings", "Settings", icons.Settings(icons.WithSize(16))),
							dropdownMenuItem("/help", "Help", icons.Info(icons.WithSize(16))),
							separator.Separator(separator.WithClass("my-1")),
							dropdownMenuItem("/logout", "Log out", icons.LogOut(icons.WithSize(16))),
						),
					),
				),
			),
		),

		// Toggle button
		sidebar.SidebarToggle(),

		// Rail for better UX
		sidebar.SidebarRail(),
	)
}

// renderMetricCards creates the metric cards
func renderMetricCards() []g.Node {
	metrics := GetDashboardMetrics()
	cards := make([]g.Node, len(metrics))

	for i, metric := range metrics {
		cards[i] = metricCard(metric)
	}

	return cards
}

// metricCard creates a single metric card
func metricCard(metric DashboardMetric) g.Node {
	// Determine icon based on title
	var metricIcon g.Node
	switch metric.Title {
	case "Total Revenue":
		metricIcon = icons.DollarSign(icons.WithSize(16), icons.WithClass("text-muted-foreground"))
	case "New Customers":
		metricIcon = icons.Users(icons.WithSize(16), icons.WithClass("text-muted-foreground"))
	case "Active Accounts":
		metricIcon = icons.CreditCard(icons.WithSize(16), icons.WithClass("text-muted-foreground"))
	case "Growth Rate":
		metricIcon = icons.Activity(icons.WithSize(16), icons.WithClass("text-muted-foreground"))
	default:
		metricIcon = icons.Activity(icons.WithSize(16), icons.WithClass("text-muted-foreground"))
	}

	// Trend icon
	var trendIcon g.Node
	var trendColor string
	if metric.ChangeType == "up" {
		trendIcon = icons.TrendingUp(icons.WithSize(16))
		trendColor = "text-success"
	} else {
		trendIcon = icons.TrendingDown(icons.WithSize(16))
		trendColor = "text-destructive"
	}

	return card.Card(
		card.Header(
			html.Div(
				html.Class("flex flex-row items-center justify-between space-y-0 pb-2"),
				card.Title(metric.Title),
				metricIcon,
			),
		),
		card.Content(
			html.Div(
				html.Class("text-2xl font-bold"),
				g.Text(metric.Value),
			),
			html.Div(
				html.Class("flex items-center gap-1 text-xs text-muted-foreground mt-1"),
				html.Span(
					html.Class(fmt.Sprintf("flex items-center gap-1 font-medium %s", trendColor)),
					trendIcon,
					g.Text(metric.Change),
				),
				html.Span(
					g.Text(metric.Description),
				),
			),
		),
	)
}

// activityItem creates an activity list item
func activityItem(title, time string, icon g.Node) g.Node {
	return html.Div(
		html.Class("flex items-start gap-3"),
		html.Div(
			html.Class("flex h-9 w-9 items-center justify-center rounded-full bg-primary/10 text-primary"),
			icon,
		),
		html.Div(
			html.Class("flex-1 space-y-1"),
			html.P(
				html.Class("text-sm font-medium leading-none"),
				g.Text(title),
			),
			html.P(
				html.Class("text-xs text-muted-foreground"),
				g.Text(time),
			),
		),
	)
}

// dropdownMenuItem creates a dropdown menu item
func dropdownMenuItem(href, label string, icon g.Node) g.Node {
	return html.A(
		g.Attr("href", href),
		html.Class("flex items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent cursor-pointer"),
		icon,
		g.Text(label),
	)
}
