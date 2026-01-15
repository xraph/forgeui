package main

import (
	"context"
	"log"
	"net/http"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/bridge"
	"github.com/xraph/forgeui/components/alert"
	"github.com/xraph/forgeui/components/avatar"
	"github.com/xraph/forgeui/components/badge"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/components/checkbox"
	"github.com/xraph/forgeui/components/form"
	"github.com/xraph/forgeui/components/input"
	"github.com/xraph/forgeui/components/label"
	"github.com/xraph/forgeui/components/progress"
	"github.com/xraph/forgeui/components/radio"
	"github.com/xraph/forgeui/components/separator"
	"github.com/xraph/forgeui/components/skeleton"
	"github.com/xraph/forgeui/components/spinner"
	"github.com/xraph/forgeui/icons"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

func main() {
	// Run the modern ForgeUI demo with Next.js-inspired API
	// RunModernDemo()

	mainLegacy()
}

// mainLegacy contains the original demo implementation
// Uncomment this and swap with main() to run the legacy demo
func mainLegacy() {
	// Initialize ForgeUI App with asset management
	app := forgeui.New(
		forgeui.WithDebug(true),
		forgeui.WithAssetPublicDir("example/static"),
	)

	// NEW: Initialize Bridge System (Phase 21)
	b := bridge.New(
		bridge.WithTimeout(30), // 30 seconds
		bridge.WithCSRF(false), // Disable for demo
	)

	// Register bridge functions
	_ = b.Register("Greet", Greet)
	_ = b.Register("Add", Add)
	_ = b.Register("GetUser", GetUser)
	_ = b.Register("GetNotifications", GetNotifications)

	log.Println("🌉 Bridge system initialized with 4 functions")

	// NEW: Start dev server with hot reload in development mode
	if app.IsDev() {
		go func() {
			log.Println("🔥 Starting dev server with hot reload...")
			if err := app.Assets.StartDevServer(context.Background()); err != nil {
				log.Printf("Dev server error: %v\n", err)
			}
		}()
	}

	// Serve static files through asset pipeline
	// In development mode: no fingerprinting, moderate caching
	// In production mode: automatic fingerprinting, immutable caching
	http.Handle("/static/", app.Assets.Handler())

	// NEW: Bridge HTTP endpoints
	http.Handle("/api/bridge/call", b.Handler())
	http.Handle("/api/bridge/stream/", b.StreamHandler())

	// NEW: SSE endpoint for hot reload
	if app.IsDev() {
		if handler := app.Assets.SSEHandler(); handler != nil {
			// Use http.Handle instead of http.HandleFunc to avoid type assertion
			// SSEHandler() returns http.HandlerFunc wrapped in interface{}
			http.Handle("/_forgeui/reload", handler.(http.Handler))
		}
	}

	// Main demo page
	http.HandleFunc("/", handleIndex)

	// Complete components showcase
	http.HandleFunc("/showcase", handleComponentsShowcase)

	// Interactive demo page
	http.HandleFunc("/interactive", handleInteractive)

	// Overlays demo page
	http.HandleFunc("/overlays", handleOverlays)

	// Navigation demo page
	http.HandleFunc("/navigation", handleNavigation)

	// Theme demo page
	http.HandleFunc("/theme", handleTheme)

	// Data components demo page
	http.HandleFunc("/data", handleData)

	// Icons demo page - demonstrates icon system
	http.HandleFunc("/icons", handleIcons)

	// Dashboard demo page - comprehensive sidebar and dashboard example
	http.HandleFunc("/dashboard", handleDashboard)

	// Assets demo page - demonstrates asset pipeline features
	http.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
		handleAssetsDemo(w, r, app)
	})

	// NEW: Bridge demo page - demonstrates Go↔JS bridge (Phase 21)
	http.HandleFunc("/bridge", func(w http.ResponseWriter, r *http.Request) {
		handleBridgeDemo(w, r, app, b)
	})

	log.Println("🚀 ForgeUI Example running at http://localhost:8080")
	log.Println("📄 Static Components: http://localhost:8080")
	log.Println("✨ Complete Showcase: http://localhost:8080/showcase")
	log.Println("🎮 Interactive Demo: http://localhost:8080/interactive")
	log.Println("📊 Dashboard Demo: http://localhost:8080/dashboard")
	log.Println("🔲 Overlay Components: http://localhost:8080/overlays")
	log.Println("🧭 Navigation Components: http://localhost:8080/navigation")
	log.Println("🎨 Theme System: http://localhost:8080/theme")
	log.Println("📊 Data Components: http://localhost:8080/data")
	log.Println("🎯 Icons Library: http://localhost:8080/icons")
	log.Println("📦 Assets Pipeline Demo: http://localhost:8080/assets")
	log.Println("🌉 Bridge System Demo: http://localhost:8080/bridge")

	if app.IsDev() {
		log.Println("🔥 Hot reload enabled - edit Go files to see changes instantly!")
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Get ForgeUI app from context or create new one (for demo purposes, we'll access global)
	// In production, pass app through context or handler closure
	app := forgeui.New(forgeui.WithDebug(true))

	page := html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Component Library Demo")),
			// Note: Using CDN for Tailwind CSS for simplicity in demo
			// In production, use: app.Assets.StyleSheet("css/tailwind.css")
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
			// NEW: Hot reload script in development mode
			g.If(app.IsDev(), g.Raw(app.Assets.HotReloadScript())),
		),
		html.Body(
			html.Class("bg-background text-foreground"),
			g.Attr("x-data", ""), // Initialize Alpine.js
			theme.DarkModeScript(),

			// Header with Theme Toggle
			primitives.Box(
				primitives.WithBackground("bg-card border-b border-border"),
				primitives.WithPadding("py-4"),
				primitives.WithChildren(
					primitives.Container(
						primitives.HStack("4",
							html.Div(
								html.Class("flex-1 flex items-center"),
								primitives.Text(
									primitives.TextAs("span"),
									primitives.TextSize("text-lg"),
									primitives.TextWeight("font-semibold"),
									primitives.TextChildren(g.Text("ForgeUI")),
								),
							),
							// Navigation Links
							primitives.HStack("4",
								html.A(
									html.Href("/"),
									html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors flex items-center gap-1"),
									icons.Home(icons.WithSize(16)),
									g.Text("Home"),
								),
								html.A(
									html.Href("/showcase"),
									html.Class("text-sm font-medium text-primary underline-offset-4 hover:underline flex items-center gap-1"),
									icons.Menu(icons.WithSize(16)),
									g.Text("Showcase"),
								),
								html.A(
									html.Href("/interactive"),
									html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors flex items-center gap-1"),
									icons.Plus(icons.WithSize(16)),
									g.Text("Interactive"),
								),
								html.A(
									html.Href("/data"),
									html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors"),
									g.Text("Data"),
								),
								html.A(
									html.Href("/icons"),
									html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors"),
									g.Text("Icons"),
								),
								html.A(
									html.Href("/theme"),
									html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors"),
									g.Text("Theme"),
								),
								html.A(
									html.Href("/assets"),
									html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors"),
									g.Text("Assets"),
								),
								html.A(
									html.Href("/bridge"),
									html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors"),
									g.Text("Bridge"),
								),
							),
							// Theme Toggle
							theme.SimpleToggle(),
						),
					),
				),
			),

			// Hero Section
			primitives.Box(
				primitives.WithBackground("bg-gradient-to-r from-blue-500 to-purple-600"),
				primitives.WithPadding("py-24 md:py-32"),
				primitives.WithChildren(
					primitives.Container(
						primitives.VStack("8",
							primitives.Text(
								primitives.TextAs("h1"),
								primitives.TextSize("text-5xl md:text-6xl lg:text-7xl"),
								primitives.TextWeight("font-bold"),
								primitives.TextColor("text-white"),
								primitives.TextAlign("text-center"),
								primitives.TextClass("tracking-tight"),
								primitives.TextChildren(g.Text("ForgeUI Component Library")),
							),
							primitives.Text(
								primitives.TextSize("text-xl md:text-2xl"),
								primitives.TextColor("text-white/90"),
								primitives.TextAlign("text-center"),
								primitives.TextClass("max-w-3xl mx-auto"),
								primitives.TextChildren(g.Text("SSR-first UI components for Go with Tailwind CSS")),
							),
							primitives.Center(
								button.Primary(
									g.Group([]g.Node{
										icons.ChevronRight(icons.WithSize(16)),
										g.Text("Get Started"),
									}),
									button.WithSize(forgeui.SizeLG),
								),
							),
						),
					),
				),
			),

			// Main Content
			primitives.Box(
				primitives.WithPadding("py-16 md:py-24"),
				primitives.WithChildren(
					primitives.Container(
						primitives.VStack("16",

							// Buttons Section
							componentSection(
								"Buttons",
								"Multiple variants and sizes for different use cases",
								primitives.VStack("4",
									primitives.HStack("2",
										button.Primary(g.Text("Primary")),
										button.Secondary(g.Text("Secondary")),
										button.Destructive(g.Text("Destructive")),
										button.Outline(g.Text("Outline")),
										button.Ghost(g.Text("Ghost")),
										button.Link(g.Text("Link")),
									),
									primitives.HStack("2",
										button.Primary(g.Text("Small"), button.WithSize(forgeui.SizeSM)),
										button.Primary(g.Text("Default")),
										button.Primary(g.Text("Large"), button.WithSize(forgeui.SizeLG)),
										button.IconButton(g.Text("×")),
									),
									button.Group(
										[]button.GroupOption{button.WithGap("2")},
										button.Primary(g.Text("Save")),
										button.Secondary(g.Text("Cancel")),
									),
								),
							),

							// Cards Section
							componentSection(
								"Cards",
								"Compound components for content containers",
								primitives.Grid(
									primitives.GridCols(1),
									primitives.GridColsMD(2),
									primitives.GridColsLG(3),
									primitives.GridGap("6"),
									primitives.GridChildren(
										card.Card(
											card.Header(
												card.Title("Getting Started"),
												card.Description("Quick introduction to ForgeUI"),
											),
											card.Content(
												primitives.Text(
													primitives.TextSize("text-sm"),
													primitives.TextChildren(g.Text("ForgeUI is a type-safe UI library for Go that renders server-side HTML with gomponents.")),
												),
											),
											card.Footer(
												button.Primary(g.Text("Learn More"), button.WithSize(forgeui.SizeSM)),
											),
										),
										card.Card(
											card.Header(
												card.Title("Components"),
												card.Description("25+ production-ready components"),
											),
											card.Content(
												primitives.VStack("2",
													badge.Badge("Primitives", badge.WithVariant(forgeui.VariantSecondary)),
													badge.Badge("Forms", badge.WithVariant(forgeui.VariantDefault)),
													badge.Badge("Overlays", badge.WithVariant(forgeui.VariantOutline)),
												),
											),
										),
										card.Card(
											card.Header(
												card.Title("Type Safe"),
												card.Description("Full Go type checking"),
											),
											card.Content(
												primitives.VStack("2",
													primitives.HStack("2",
														avatar.Avatar(
															avatar.WithFallback("FU"),
															avatar.WithSize(forgeui.SizeSM),
														),
														primitives.Text(
															primitives.TextSize("text-sm"),
															primitives.TextChildren(g.Text("Built with Go generics")),
														),
													),
													progress.Progress(progress.WithValue(85)),
												),
											),
										),
									),
								),
							),

							// Badges and Avatars
							componentSection(
								"Badges & Avatars",
								"Visual indicators and user representations",
								primitives.VStack("4",
									primitives.HStack("2",
										badge.Badge("Default"),
										badge.Badge("Secondary", badge.WithVariant(forgeui.VariantSecondary)),
										badge.Badge("Destructive", badge.WithVariant(forgeui.VariantDestructive)),
										badge.Badge("Outline", badge.WithVariant(forgeui.VariantOutline)),
									),
									primitives.HStack("2",
										avatar.Avatar(
											avatar.WithFallback("AB"),
											avatar.WithSize(forgeui.SizeSM),
										),
										avatar.Avatar(
											avatar.WithFallback("CD"),
											avatar.WithSize(forgeui.SizeMD),
										),
										avatar.Avatar(
											avatar.WithFallback("EF"),
											avatar.WithSize(forgeui.SizeLG),
										),
										avatar.Avatar(
											avatar.WithFallback("GH"),
											avatar.WithSize(forgeui.SizeXL),
										),
									),
								),
							),

							// Alerts Section
							componentSection(
								"Alerts",
								"Display important messages",
								primitives.VStack("4",
									alert.Alert(
										nil,
										g.Group([]g.Node{
											icons.Info(icons.WithSize(20), icons.WithClass("mr-2")),
											alert.AlertTitle("Heads up!"),
										}),
										alert.AlertDescription("You can add components to your app using ForgeUI."),
									),
									alert.Alert(
										[]alert.Option{alert.WithVariant(forgeui.VariantDestructive)},
										g.Group([]g.Node{
											icons.AlertCircle(icons.WithSize(20), icons.WithClass("mr-2")),
											alert.AlertTitle("Error"),
										}),
										alert.AlertDescription("Your session has expired. Please log in again."),
									),
								),
							),

							// Form Example
							componentSection(
								"Form Components",
								"Complete form controls with validation",
								card.Card(
									card.Header(
										card.Title("Create Account"),
										card.Description("Enter your details to get started"),
									),
									card.Content(
										form.Form(
											[]form.Option{form.WithAction("/submit")},

											primitives.VStack("4",
												// Name input
												primitives.VStack("2",
													label.Label("Name", label.WithFor("name")),
													input.Input(
														input.WithID("name"),
														input.WithName("name"),
														input.WithPlaceholder("John Doe"),
														input.Required(),
													),
												),

												// Email input
												primitives.VStack("2",
													label.Label("Email", label.WithFor("email")),
													input.Input(
														input.WithType("email"),
														input.WithID("email"),
														input.WithName("email"),
														input.WithPlaceholder("john@example.com"),
														input.Required(),
													),
													input.FormDescription("We'll never share your email."),
												),

												// Checkboxes
												primitives.VStack("2",
													primitives.HStack("2",
														checkbox.Checkbox(
															checkbox.WithID("terms"),
															checkbox.WithName("terms"),
														),
														label.Label("Accept terms and conditions", label.WithFor("terms")),
													),
													primitives.HStack("2",
														checkbox.Checkbox(
															checkbox.WithID("newsletter"),
															checkbox.WithName("newsletter"),
														),
														label.Label("Subscribe to newsletter", label.WithFor("newsletter")),
													),
												),

												// Radio buttons
												primitives.VStack("2",
													label.Label("Account Type"),
													radio.RadioGroup("account_type", []radio.RadioGroupOption{
														{ID: "personal", Value: "personal", Label: "Personal", Checked: true},
														{ID: "business", Value: "business", Label: "Business"},
													}),
												),

												separator.Separator(),

												button.Primary(
													g.Text("Create Account"),
													button.WithType("submit"),
													button.WithClass("w-full"),
												),
											),
										),
									),
								),
							),

							// Loading States
							componentSection(
								"Loading States",
								"Spinners, skeletons, and progress indicators",
								primitives.VStack("6",
									primitives.HStack("4",
										spinner.Spinner(spinner.WithSize(forgeui.SizeSM)),
										spinner.Spinner(spinner.WithSize(forgeui.SizeMD)),
										spinner.Spinner(spinner.WithSize(forgeui.SizeLG)),
									),
									primitives.VStack("2",
										skeleton.Skeleton(
											skeleton.WithHeight("h-4"),
											skeleton.WithWidth("w-full"),
										),
										skeleton.Skeleton(
											skeleton.WithHeight("h-4"),
											skeleton.WithWidth("w-3/4"),
										),
										skeleton.Skeleton(
											skeleton.WithHeight("h-4"),
											skeleton.WithWidth("w-1/2"),
										),
									),
									primitives.VStack("2",
										primitives.Text(
											primitives.TextSize("text-sm"),
											primitives.TextWeight("font-medium"),
											primitives.TextChildren(g.Text("Upload Progress")),
										),
										progress.Progress(progress.WithValue(65)),
									),
								),
							),
						)),
				),
			),

			// Footer
			primitives.Box(
				primitives.WithBackground("bg-gray-900"),
				primitives.WithPadding("py-16 md:py-20"),
				primitives.WithChildren(
					primitives.Container(
						primitives.VStack("6",
							primitives.Text(
								primitives.TextSize("text-2xl"),
								primitives.TextWeight("font-bold"),
								primitives.TextColor("text-white"),
								primitives.TextChildren(g.Text("ForgeUI")),
							),
							primitives.Text(
								primitives.TextSize("text-sm"),
								primitives.TextColor("text-gray-400"),
								primitives.TextChildren(g.Text("Built with ❤️ using Go and Tailwind CSS")),
							),
							primitives.Text(
								primitives.TextSize("text-xs"),
								primitives.TextColor("text-gray-500"),
								primitives.TextChildren(g.Text("© 2024 ForgeUI. All rights reserved.")),
							),
						),
					),
				),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

func handleInteractive(w http.ResponseWriter, r *http.Request) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("ForgeUI - Interactive Demo")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),
			alpine.CloakCSS(),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground"),
			g.Attr("x-data", ""), // Initialize Alpine.js
			theme.DarkModeScript(),
			html.Div(
				html.Class("container mx-auto px-4 py-12"),
				html.Div(
					html.Class("mb-8 space-y-2"),
					html.H1(
						html.Class("text-4xl font-bold"),
						g.Text("ForgeUI Interactive Demo"),
					),
					html.P(
						html.Class("text-lg text-muted-foreground"),
						g.Text("Alpine.js integration with animations"),
					),
					html.Div(
						html.Class("flex gap-2 mt-4"),
						html.A(
							html.Href("/"),
							html.Class("text-sm text-primary hover:underline"),
							g.Text("← Back to Components"),
						),
					),
				),
				InteractiveDemoPage(),
			),
			// Alpine.js scripts (must be at end of body)
			alpine.Scripts(),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

func componentSection(title, description string, content g.Node) g.Node {
	return primitives.VStack("8",
		primitives.VStack("3",
			primitives.Text(
				primitives.TextAs("h2"),
				primitives.TextSize("text-3xl md:text-4xl"),
				primitives.TextWeight("font-bold"),
				primitives.TextClass("tracking-tight"),
				primitives.TextChildren(g.Text(title)),
			),
			primitives.Text(
				primitives.TextSize("text-lg md:text-xl"),
				primitives.TextColor("text-muted-foreground"),
				primitives.TextClass("leading-relaxed"),
				primitives.TextChildren(g.Text(description)),
			),
		),
		content,
	)
}

func handleOverlays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := OverlaysDemo().Render(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleNavigation(w http.ResponseWriter, r *http.Request) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("ForgeUI - Navigation Components")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),
			alpine.CloakCSS(),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground antialiased"),
			g.Attr("x-data", ""), // Initialize Alpine.js
			theme.DarkModeScript(),
			html.Main(
				html.Class("container mx-auto px-4 py-12"),
				html.Div(
					html.Class("mb-8 space-y-2"),
					html.H1(
						html.Class("text-4xl font-bold"),
						g.Text("ForgeUI Navigation Demo"),
					),
					html.P(
						html.Class("text-lg text-muted-foreground"),
						g.Text("Navigation components with Alpine.js state management"),
					),
					html.Div(
						html.Class("flex gap-2 mt-4"),
						html.A(
							html.Href("/"),
							html.Class("text-sm text-primary hover:underline"),
							g.Text("← Back to Components"),
						),
						html.A(
							html.Href("/interactive"),
							html.Class("text-sm text-primary hover:underline"),
							g.Text("Interactive Demo"),
						),
						html.A(
							html.Href("/overlays"),
							html.Class("text-sm text-primary hover:underline"),
							g.Text("Overlay Components"),
						),
					),
				),
				NavigationDemo(),
			),
			// Alpine.js scripts with Collapse plugin for accordion
			alpine.Scripts(alpine.PluginCollapse),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

func handleTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ThemeDemo().Render(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
