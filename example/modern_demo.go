package main

import (
	"log"
	"net/http"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/bridge"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/icons"
	"github.com/xraph/forgeui/layout"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/router"
	"github.com/xraph/forgeui/theme"
)

// ModernDemo demonstrates the new ForgeUI API with:
// - New() initialization
// - Fluent Page builder API
// - Root layout pattern with nested layouts
// - Auto-initialized bridge system
// - Framework-agnostic design

func RunModernDemo() {
	// Initialize ForgeUI with modern API
	lightTheme := theme.DefaultLight()
	darkTheme := theme.DefaultDark()

	app := forgeui.New(
		forgeui.WithDev(true),
		forgeui.WithAssets("example/static"),
		forgeui.WithBridge(
			bridge.WithTimeout(30),
			bridge.WithCSRF(false),
		),
		forgeui.WithThemes(&lightTheme, &darkTheme),
		forgeui.WithDefaultLayout("root"),
	)

	// Register bridge functions
	if app.HasBridge() {
		_ = app.Bridge().Register("Greet", Greet)
		_ = app.Bridge().Register("GetUser", GetUser)
	}

	// Register root layout - defines HTML structure for ALL pages
	app.RegisterLayout("root", RootLayout)

	// Register dashboard layout - inherits from root, adds sidebar
	// The parent layout (root) will wrap this layout
	app.RegisterLayout("dashboard", DashboardLayout, router.WithParentLayout("root"))

	// Register pages using fluent API
	app.Page("/modern").
		Handler(ModernHomePage).
		Meta("Modern ForgeUI", "Next.js-inspired API for Go").
		Register()

	app.Page("/modern/about").
		Handler(AboutPage).
		Meta("About", "Learn about ForgeUI").
		Register()

	// Dashboard pages with nested layout
	dashboardGroup := app.Group("/modern/dashboard").
		Layout("dashboard")

	dashboardGroup.Get("/", DashboardHomePage).
		WithMeta(&router.RouteMeta{Title: "Dashboard", Description: "Your dashboard"})

	dashboardGroup.Get("/settings", SettingsPage).
		WithMeta(&router.RouteMeta{Title: "Settings", Description: "Manage your settings"})

	// API endpoint without layout
	app.Page("/modern/api/status").
		NoLayout().
		Handler(APIStatusHandler).
		Register()

	// Start server - ForgeUI integrates with any framework
	log.Println("🚀 ForgeUI Modern Demo running at http://localhost:8080")
	log.Println("📄 Home: http://localhost:8080/modern")
	log.Println("📖 About: http://localhost:8080/modern/about")
	log.Println("📊 Dashboard: http://localhost:8080/modern/dashboard")
	log.Println("⚙️  Settings: http://localhost:8080/modern/dashboard/settings")
	log.Println("🔌 API Status: http://localhost:8080/modern/api/status")

	if err := http.ListenAndServe(":8080", app.Handler()); err != nil {
		log.Fatal(err)
	}
}

// RootLayout - THE ONLY layout with full HTML structure
// This is the single source of truth for head/body configuration
func RootLayout(ctx *router.PageContext, content g.Node) g.Node {
	// Type assert to get app (interface{} to avoid circular dependency)
	app := ctx.App().(*forgeui.App)

	return layout.Build(
		layout.Head(
			layout.Meta("viewport", "width=device-width, initial-scale=1"),
			layout.Charset("utf-8"),

			// Single source of truth for theme
			layout.Theme(app.LightTheme(), app.DarkTheme()),

			// Alpine.js cloak CSS
			layout.Alpine(),

			// Tailwind CSS (using CDN for demo simplicity)
			// In production, use: app.Assets.StyleSheet("css/tailwind.css")
			html.Script(
				html.Src("https://cdn.tailwindcss.com"),
			),

			// Page-specific metadata
			g.If(ctx.Meta != nil, g.Group([]g.Node{
				layout.Title(ctx.Meta.Title),
				layout.Description(ctx.Meta.Description),
			})),
		),

		layout.Body(
			layout.Class("min-h-screen bg-background text-foreground antialiased"),
			g.Attr("x-data", ""), // Alpine initialization

			// Dark mode script
			layout.DarkModeScript(),

			// Child content rendered here (from child layouts or pages)
			content,

			// Scripts at end of body
			layout.Scripts(
				layout.AlpineScripts(alpine.PluginCollapse),
			),

			// Auto-inject based on app configuration (disabled for now to avoid errors)
			// g.If(app.IsDev(), layout.HotReload()),
			g.If(app.HasBridge(), layout.BridgeClient()),
		),
	)
}

// DashboardLayout - ONLY wraps content with sidebar, inherits root's HTML structure
func DashboardLayout(ctx *router.PageContext, content g.Node) g.Node {
	return html.Div(
		html.Class("flex min-h-screen"),

		// Sidebar
		html.Aside(
			html.Class("w-64 bg-card border-r border-border p-4"),
			html.Nav(
				html.Class("space-y-2"),
				html.H2(
					html.Class("text-lg font-semibold mb-4"),
					g.Text("Dashboard"),
				),
				NavLink("/modern/dashboard", "Home", icons.Home(icons.WithSize(18))),
				NavLink("/modern/dashboard/settings", "Settings", icons.Settings(icons.WithSize(18))),
			),
		),

		// Main content area
		html.Main(
			html.Class("flex-1 p-8"),
			content,
		),
	)
}

// NavLink creates a navigation link
func NavLink(href, label string, iconNode g.Node) g.Node {
	return html.A(
		html.Href(href),
		html.Class("flex items-center gap-2 px-3 py-2 rounded-md hover:bg-accent transition-colors"),
		iconNode,
		g.Text(label),
	)
}

// ModernHomePage - Simple handler returning content only
func ModernHomePage(ctx *router.PageContext) (g.Node, error) {
	return primitives.Container(
		primitives.VStack("8",
			// Hero section
			primitives.VStack("4",
				primitives.Text(
					primitives.TextAs("h1"),
					primitives.TextSize("text-5xl"),
					primitives.TextWeight("font-bold"),
					primitives.TextChildren(g.Text("Modern ForgeUI API")),
				),
				primitives.Text(
					primitives.TextSize("text-xl"),
					primitives.TextColor("text-muted-foreground"),
					primitives.TextChildren(g.Text("Next.js-inspired declarative API for Go")),
				),
			),

			// Features
			primitives.Grid(
				primitives.GridCols(1),
				primitives.GridColsMD(2),
				primitives.GridColsLG(3),
				primitives.GridGap("6"),
				primitives.GridChildren(
					FeatureCard(
						"Fluent Page Builder",
						"Declarative page registration with chained methods",
						"plus",
					),
					FeatureCard(
						"Root Layout Pattern",
						"Single source of truth for HTML structure",
						"layout",
					),
					FeatureCard(
						"Framework Agnostic",
						"Works with net/http, Fiber, Echo, Gin, Chi",
						"code",
					),
				),
			),

			// Actions
			primitives.HStack("4",
				button.Primary(
					g.Group([]g.Node{
						icons.ChevronRight(icons.WithSize(16)),
						g.Text("Get Started"),
					}),
					button.WithSize(forgeui.SizeLG),
				),
				button.Secondary(
					g.Group([]g.Node{
						icons.Book(icons.WithSize(16)),
						g.Text("Documentation"),
					}),
					button.WithSize(forgeui.SizeLG),
				),
			),
		),
	), nil
}

// FeatureCard creates a feature card
func FeatureCard(title, description, icon string) g.Node {
	return card.Card(
		card.Header(
			html.Div(
				html.Class("flex items-center gap-2"),
				icons.Icon(icon, icons.WithSize(24), icons.WithClass("text-primary")),
				card.Title(title),
			),
		),
		card.Content(
			primitives.Text(
				primitives.TextSize("text-sm"),
				primitives.TextColor("text-muted-foreground"),
				primitives.TextChildren(g.Text(description)),
			),
		),
	)
}

// AboutPage demonstrates simple page handler
func AboutPage(ctx *router.PageContext) (g.Node, error) {
	return primitives.Container(
		primitives.VStack("6",
			primitives.Text(
				primitives.TextAs("h1"),
				primitives.TextSize("text-4xl"),
				primitives.TextWeight("font-bold"),
				primitives.TextChildren(g.Text("About ForgeUI")),
			),
			primitives.Text(
				primitives.TextSize("text-lg"),
				primitives.TextColor("text-muted-foreground"),
				primitives.TextChildren(g.Text("A modern UI library for Go with SSR-first design.")),
			),
		),
	), nil
}

// DashboardHomePage demonstrates dashboard page
func DashboardHomePage(ctx *router.PageContext) (g.Node, error) {
	return primitives.VStack("6",
		primitives.Text(
			primitives.TextAs("h1"),
			primitives.TextSize("text-3xl"),
			primitives.TextWeight("font-bold"),
			primitives.TextChildren(g.Text("Dashboard Home")),
		),
		primitives.Text(
			primitives.TextChildren(g.Text("Welcome to your dashboard!")),
		),
	), nil
}

// SettingsPage demonstrates nested dashboard page
func SettingsPage(ctx *router.PageContext) (g.Node, error) {
	return primitives.VStack("6",
		primitives.Text(
			primitives.TextAs("h1"),
			primitives.TextSize("text-3xl"),
			primitives.TextWeight("font-bold"),
			primitives.TextChildren(g.Text("Settings")),
		),
		primitives.Text(
			primitives.TextChildren(g.Text("Manage your application settings.")),
		),
	), nil
}

// APIStatusHandler demonstrates API endpoint without layout
func APIStatusHandler(ctx *router.PageContext) (g.Node, error) {
	ctx.SetHeader("Content-Type", "application/json")
	return g.Raw(`{"status":"ok","version":"1.0.0"}`), nil
}
