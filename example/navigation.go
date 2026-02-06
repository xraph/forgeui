package main

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/accordion"
	"github.com/xraph/forgeui/components/breadcrumb"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/components/menu"
	navbar "github.com/xraph/forgeui/components/navbar"
	"github.com/xraph/forgeui/components/pagination"
	"github.com/xraph/forgeui/components/tabs"
	"github.com/xraph/forgeui/primitives"
)

// NavigationDemo demonstrates all navigation components
func NavigationDemo() g.Node {
	return html.Div(
		html.Class("space-y-12"),

		// Page header
		html.Div(
			html.Class("space-y-2"),
			html.H1(
				html.Class("text-4xl font-bold tracking-tight"),
				g.Text("Navigation Components"),
			),
			html.P(
				html.Class("text-lg text-muted-foreground"),
				g.Text("Interactive navigation components with Alpine.js state management"),
			),
		),

		// Tabs - Basic
		componentSection(
			"Tabs",
			"Organize content into switchable panels with keyboard navigation",
			card.Card(
				card.Content(
					tabs.TabsWithOptions(
						[]tabs.Option{tabs.WithDefaultTab("overview")},
						tabs.TabList(
							tabs.Tab("overview", g.Text("Overview")),
							tabs.Tab("details", g.Text("Details")),
							tabs.Tab("settings", g.Text("Settings")),
						),
						tabs.TabPanel("overview",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Overview content goes here. This tab shows a summary of the information.")),
							),
						),
						tabs.TabPanel("details",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Detailed information is displayed in this panel.")),
							),
						),
						tabs.TabPanel("settings",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Configure your settings in this tab.")),
							),
						),
					),
				),
			),
		),

		// Tabs - Scrollable
		componentSection(
			"Tabs (Scrollable)",
			"Horizontal scrolling tabs - WithScrollable() auto-sizes tabs to content",
			card.Card(
				card.Content(
					tabs.TabsWithOptions(
						[]tabs.Option{tabs.WithDefaultTab("monday")},
						tabs.TabListWithOptions(
							[]tabs.TabListOption{tabs.WithScrollable()},
							tabs.Tab("monday", g.Text("Monday")),
							tabs.Tab("tuesday", g.Text("Tuesday")),
							tabs.Tab("wednesday", g.Text("Wednesday")),
							tabs.Tab("thursday", g.Text("Thursday")),
							tabs.Tab("friday", g.Text("Friday")),
							tabs.Tab("saturday", g.Text("Saturday")),
							tabs.Tab("sunday", g.Text("Sunday")),
						),
						tabs.TabPanel("monday",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Monday schedule and tasks.")),
							),
						),
						tabs.TabPanel("tuesday",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Tuesday schedule and tasks.")),
							),
						),
						tabs.TabPanel("wednesday",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Wednesday schedule and tasks.")),
							),
						),
						tabs.TabPanel("thursday",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Thursday schedule and tasks.")),
							),
						),
						tabs.TabPanel("friday",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Friday schedule and tasks.")),
							),
						),
						tabs.TabPanel("saturday",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Saturday schedule and tasks.")),
							),
						),
						tabs.TabPanel("sunday",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Sunday schedule and tasks.")),
							),
						),
					),
				),
			),
		),

		// Tabs - Link-Based
		componentSection(
			"Tabs (Link-Based)",
			"Tabs using links for better SEO and navigation",
			card.Card(
				card.Content(
					html.Div(
						html.Class("space-y-4"),
						html.P(
							html.Class("text-sm text-muted-foreground"),
							g.Text("These tabs render as <a> tags with href attributes. Great for navigation!"),
						),
						tabs.TabsWithOptions(
							[]tabs.Option{tabs.WithDefaultTab("home")},
							tabs.TabList(
								tabs.Tab("home", g.Text("Home"), tabs.WithHref("#home")),
								tabs.Tab("products", g.Text("Products"), tabs.WithHref("#products")),
								tabs.Tab("about", g.Text("About"), tabs.WithHref("#about")),
								tabs.Tab("contact", g.Text("Contact"), tabs.WithHref("#contact")),
							),
							tabs.TabPanel("home",
								html.Div(
									html.Class("py-4"),
									html.P(g.Text("Welcome to our home page.")),
								),
							),
							tabs.TabPanel("products",
								html.Div(
									html.Class("py-4"),
									html.P(g.Text("Browse our product catalog.")),
								),
							),
							tabs.TabPanel("about",
								html.Div(
									html.Class("py-4"),
									html.P(g.Text("Learn more about our company.")),
								),
							),
							tabs.TabPanel("contact",
								html.Div(
									html.Class("py-4"),
									html.P(g.Text("Get in touch with us.")),
								),
							),
						),
					),
				),
			),
		),

		// Tabs - Custom Styled
		componentSection(
			"Tabs (Custom Styled)",
			"Tabs with custom CSS classes for unique designs",
			card.Card(
				card.Content(
					tabs.TabsWithOptions(
						[]tabs.Option{tabs.WithDefaultTab("design")},
						tabs.TabListWithOptions(
							[]tabs.TabListOption{tabs.WithTabListClass("bg-primary/10")},
							tabs.Tab("design", g.Text("🎨 Design"), tabs.WithTabClass("font-bold")),
							tabs.Tab("develop", g.Text("💻 Develop"), tabs.WithTabClass("font-bold")),
							tabs.Tab("deploy", g.Text("🚀 Deploy"), tabs.WithTabClass("font-bold")),
						),
						tabs.TabPanel("design",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Design phase: Create mockups and prototypes.")),
							),
						),
						tabs.TabPanel("develop",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Development phase: Build and test features.")),
							),
						),
						tabs.TabPanel("deploy",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Deployment phase: Ship to production.")),
							),
						),
					),
				),
			),
		),

		// Tabs - With Active State and Attributes
		componentSection(
			"Tabs (Active & Attributes)",
			"Demonstrating WithActive() and WithAttr() options",
			card.Card(
				card.Content(
					html.Div(
						html.Class("space-y-4"),
						html.P(
							html.Class("text-sm text-muted-foreground"),
							g.Text("The middle tab is marked as active on initial render using WithActive(), and all tabs have custom data attributes using WithAttr()."),
						),
						tabs.TabsWithOptions(
							[]tabs.Option{tabs.WithDefaultTab("step2")},
							tabs.TabList(
								tabs.Tab("step1", g.Text("Step 1"),
									tabs.WithAttr(g.Attr("data-step", "1")),
									tabs.WithAttr(g.Attr("data-status", "completed")),
								),
								tabs.Tab("step2", g.Text("Step 2"),
									tabs.WithAttr(g.Attr("data-step", "2")),
									tabs.WithAttr(g.Attr("data-status", "current")),
								),
								tabs.Tab("step3", g.Text("Step 3"),
									tabs.WithAttr(g.Attr("data-step", "3")),
									tabs.WithAttr(g.Attr("data-status", "pending")),
								),
							),
							tabs.TabPanel("step1",
								html.Div(
									html.Class("py-4"),
									html.P(g.Text("Step 1: Initial setup completed ✓")),
								),
							),
							tabs.TabPanel("step2",
								html.Div(
									html.Class("py-4"),
									html.P(g.Text("Step 2: Currently in progress... This tab is active by default!")),
								),
							),
							tabs.TabPanel("step3",
								html.Div(
									html.Class("py-4"),
									html.P(g.Text("Step 3: Ready to start")),
								),
							),
						),
					),
				),
			),
		),

		// Tabs - Underline Variant (Vercel-style)
		componentSection(
			"Tabs (Underline Variant)",
			"Minimal tabs with bottom border like Vercel",
			card.Card(
				card.Content(
					tabs.TabsWithOptions(
						[]tabs.Option{tabs.WithDefaultTab("preview")},
						tabs.TabListWithOptions(
							[]tabs.TabListOption{tabs.WithTabListVariant(tabs.TabListVariantUnderline)},
							tabs.Tab("preview", g.Text("Preview"), tabs.WithTabVariant(tabs.TabVariantUnderline)),
							tabs.Tab("code", g.Text("Code"), tabs.WithTabVariant(tabs.TabVariantUnderline)),
							tabs.Tab("output", g.Text("Output"), tabs.WithTabVariant(tabs.TabVariantUnderline)),
						),
						tabs.TabPanel("preview",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Preview pane showing live component rendering.")),
							),
						),
						tabs.TabPanel("code",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Code editor with syntax highlighting.")),
							),
						),
						tabs.TabPanel("output",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Console output and logs.")),
							),
						),
					),
				),
			),
		),

		// Tabs - Pills Variant
		componentSection(
			"Tabs (Pills Variant)",
			"Rounded pill-style tabs",
			card.Card(
				card.Content(
					tabs.TabsWithOptions(
						[]tabs.Option{tabs.WithDefaultTab("all")},
						tabs.TabListWithOptions(
							[]tabs.TabListOption{tabs.WithTabListVariant(tabs.TabListVariantPills)},
							tabs.Tab("all", g.Text("All"), tabs.WithTabVariant(tabs.TabVariantPills)),
							tabs.Tab("active", g.Text("Active"), tabs.WithTabVariant(tabs.TabVariantPills)),
							tabs.Tab("archived", g.Text("Archived"), tabs.WithTabVariant(tabs.TabVariantPills)),
						),
						tabs.TabPanel("all",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("All items displayed here.")),
							),
						),
						tabs.TabPanel("active",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Only active items.")),
							),
						),
						tabs.TabPanel("archived",
							html.Div(
								html.Class("py-4"),
								html.P(g.Text("Archived items.")),
							),
						),
					),
				),
			),
		),

		// Accordion
		componentSection(
			"Accordion",
			"Collapsible sections with smooth animations using Alpine Collapse",
			card.Card(
				card.Content(
					accordion.AccordionWithOptions(
						[]accordion.Option{
							accordion.WithType(accordion.TypeSingle),
							accordion.WithCollapsible(),
						},
						accordion.Item("item1", "What is ForgeUI?",
							html.P(
								html.Class("text-sm text-muted-foreground"),
								g.Text("ForgeUI is a Go-based UI component library that combines the power of gomponents with shadcn/ui design patterns and Alpine.js for interactivity."),
							),
						),
						accordion.Item("item2", "How does it work?",
							html.P(
								html.Class("text-sm text-muted-foreground"),
								g.Text("ForgeUI components are built using gomponents for HTML generation, styled with Tailwind CSS following shadcn/ui patterns, and enhanced with Alpine.js for client-side interactivity."),
							),
						),
						accordion.Item("item3", "Is it production ready?",
							html.P(
								html.Class("text-sm text-muted-foreground"),
								g.Text("ForgeUI is actively being developed. Components are tested and follow best practices for accessibility and performance."),
							),
						),
					),
				),
			),
		),

		// Breadcrumb
		componentSection(
			"Breadcrumb",
			"Show the current page's location within a navigational hierarchy",
			card.Card(
				card.Content(
					primitives.VStack("4",
						breadcrumb.Breadcrumb(
							breadcrumb.Item("/", g.Text("Home")),
							breadcrumb.Item("/docs", g.Text("Documentation")),
							breadcrumb.Item("/docs/components", g.Text("Components")),
							breadcrumb.Page(g.Text("Navigation")),
						),
						html.Div(
							html.Class("text-sm text-muted-foreground"),
							g.Text("Basic breadcrumb with chevron separators"),
						),
					),
				),
			),
		),

		// Pagination
		componentSection(
			"Pagination",
			"Navigate through pages with dynamic page number buttons",
			card.Card(
				card.Content(
					primitives.VStack("6",
						html.Div(
							pagination.Pagination(
								pagination.WithCurrentPage(1),
								pagination.WithTotalPages(10),
								pagination.WithSiblingCount(1),
								pagination.WithShowPrevNext(true),
							),
						),
						html.Div(
							pagination.Pagination(
								pagination.WithCurrentPage(5),
								pagination.WithTotalPages(20),
								pagination.WithSiblingCount(2),
								pagination.WithShowFirstLast(true),
								pagination.WithShowPrevNext(true),
							),
						),
						html.Div(
							html.Class("text-sm text-muted-foreground text-center"),
							g.Text("Click page numbers to navigate. State is managed with Alpine.js."),
						),
					),
				),
			),
		),

		// Menu
		componentSection(
			"Menu",
			"Vertical navigation menu with sections and active states",
			html.Div(
				html.Class("grid grid-cols-1 md:grid-cols-2 gap-6"),
				card.Card(
					card.Header(
						card.Title("Simple Menu"),
					),
					card.Content(
						menu.Menu(
							menu.Item("/", g.Text("Home"), menu.Active()),
							menu.Item("/about", g.Text("About")),
							menu.Item("/services", g.Text("Services")),
							menu.Item("/contact", g.Text("Contact")),
						),
					),
				),
				card.Card(
					card.Header(
						card.Title("Menu with Sections"),
					),
					card.Content(
						menu.Menu(
							menu.Section("Main",
								menu.Item("/dashboard", g.Text("Dashboard"), menu.Active()),
								menu.Item("/analytics", g.Text("Analytics")),
							),
							menu.Separator(),
							menu.Section("Settings",
								menu.Item("/profile", g.Text("Profile")),
								menu.Item("/preferences", g.Text("Preferences")),
							),
						),
					),
				),
			),
		),

		// Navbar
		componentSection(
			"Navbar",
			"Responsive navigation bar with mobile menu drawer",
			card.Card(
				card.Content(
					html.Div(
						html.Class("border rounded-md overflow-hidden"),
						navbar.Navbar(
							navbar.NavbarBrand(
								g.Text("ForgeUI"),
							),
							navbar.NavbarMenu(
								menu.Item("/", g.Text("Home"), menu.Active()),
								menu.Item("/docs", g.Text("Docs")),
								menu.Item("/components", g.Text("Components")),
								menu.Item("/examples", g.Text("Examples")),
							),
							navbar.NavbarActions(
								button.Ghost(g.Text("Sign In")),
								button.Primary(g.Text("Get Started")),
							),
						),
					),
					html.Div(
						html.Class("mt-4 text-sm text-muted-foreground"),
						g.Text("Resize your browser to see mobile menu behavior."),
					),
				),
			),
		),

		// Sidebar
		componentSection(
			"Sidebar",
			"Collapsible sidebar with icon-only mode and content layout",
			card.Card(
				card.Content(
					html.Div(
						html.Class("border rounded-md h-96 overflow-hidden"),
						// Wrapper with Alpine data for demo (simulating page layout)
						html.Div(
							alpine.XData(map[string]any{
								"demoCollapsed": false,
								"demoMobile":    false,
							}),
							g.Attr("x-init", `
								demoMobile = window.innerWidth < 768;
								window.addEventListener('resize', () => {
									demoMobile = window.innerWidth < 768;
								});
							`),
							html.Class("relative h-full flex"),

							// Sidebar (absolute positioned within demo)
							html.Aside(
								html.Class("absolute top-0 bottom-0 left-0 z-30 flex flex-col border-r bg-background transition-all duration-300"),
								g.Attr(":class", "demoCollapsed ? 'w-16' : 'w-64'"),

								// Header
								html.Div(
									html.Class("flex items-center gap-2 border-b px-4 py-4 font-semibold"),
									g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
									html.Span(
										g.Attr("x-show", "!demoCollapsed"),
										html.Class("font-bold text-lg"),
										g.Text("App"),
									),
								),

								// Content
								html.Div(
									html.Class("flex-1 overflow-y-auto p-4"),
									// Dashboard section - hide label when collapsed
									html.Div(
										html.Class("flex flex-col gap-2 mb-4"),
										html.H4(
											html.Class("px-3 text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2"),
											g.Attr("x-show", "!demoCollapsed"),
											g.Text("Dashboard"),
										),
										html.A(
											g.Attr("href", "#dashboard"),
											html.Class("flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium bg-accent text-accent-foreground"),
											g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
											html.Span(
												html.Class("text-lg"),
												g.Text("📊"),
											),
											html.Span(
												g.Attr("x-show", "!demoCollapsed"),
												g.Text("Overview"),
											),
										),
										html.A(
											g.Attr("href", "#analytics"),
											html.Class("flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"),
											g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
											html.Span(
												html.Class("text-lg"),
												g.Text("📈"),
											),
											html.Span(
												g.Attr("x-show", "!demoCollapsed"),
												g.Text("Analytics"),
											),
										),
										html.A(
											g.Attr("href", "#reports"),
											html.Class("flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"),
											g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
											html.Span(
												html.Class("text-lg"),
												g.Text("📄"),
											),
											html.Span(
												g.Attr("x-show", "!demoCollapsed"),
												g.Text("Reports"),
											),
										),
									),
									// Separator
									html.Hr(
										html.Class("my-2 border-t border-border"),
									),
									// Settings section
									html.Div(
										html.Class("flex flex-col gap-2"),
										html.H4(
											html.Class("px-3 text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2"),
											g.Attr("x-show", "!demoCollapsed"),
											g.Text("Settings"),
										),
										html.A(
											g.Attr("href", "#profile"),
											html.Class("flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"),
											g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
											html.Span(
												html.Class("text-lg"),
												g.Text("👤"),
											),
											html.Span(
												g.Attr("x-show", "!demoCollapsed"),
												g.Text("Profile"),
											),
										),
										html.A(
											g.Attr("href", "#team"),
											html.Class("flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"),
											g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
											html.Span(
												html.Class("text-lg"),
												g.Text("👥"),
											),
											html.Span(
												g.Attr("x-show", "!demoCollapsed"),
												g.Text("Team"),
											),
										),
										html.A(
											g.Attr("href", "#billing"),
											html.Class("flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"),
											g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
											html.Span(
												html.Class("text-lg"),
												g.Text("💳"),
											),
											html.Span(
												g.Attr("x-show", "!demoCollapsed"),
												g.Text("Billing"),
											),
										),
									),
								),

								// Footer
								html.Div(
									html.Class("border-t p-4"),
									html.A(
										g.Attr("href", "#help"),
										html.Class("flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"),
										g.Attr(":class", "demoCollapsed ? 'justify-center' : ''"),
										html.Span(
											html.Class("text-lg"),
											g.Text("❓"),
										),
										html.Span(
											g.Attr("x-show", "!demoCollapsed"),
											g.Text("Help"),
										),
									),
								),

								// Toggle button
								html.Button(
									g.Attr("type", "button"),
									alpine.XOn("click", "demoCollapsed = !demoCollapsed"),
									html.Class("absolute -right-3 top-20 z-40 flex h-6 w-6 items-center justify-center rounded-full border bg-background shadow-md text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-transform"),
									g.Attr("aria-label", "Toggle sidebar"),
									g.El("svg",
										html.Class("h-4 w-4 transition-transform duration-200"),
										g.Attr(":class", "demoCollapsed ? 'rotate-180' : ''"),
										g.Attr("xmlns", "http://www.w3.org/2000/svg"),
										g.Attr("fill", "none"),
										g.Attr("viewBox", "0 0 24 24"),
										g.Attr("stroke", "currentColor"),
										g.El("path",
											g.Attr("stroke-linecap", "round"),
											g.Attr("stroke-linejoin", "round"),
											g.Attr("stroke-width", "2"),
											g.Attr("d", "m15 18-6-6 6-6"),
										),
									),
								),
							),

							// Content area that adjusts for sidebar
							html.Div(
								html.Class("flex-1 transition-all duration-300 h-full flex items-center justify-center bg-muted/20 p-8"),
								g.Attr(":class", "demoCollapsed ? 'ml-16' : 'ml-64'"),
								html.Div(
									html.Class("text-center space-y-2"),
									html.P(
										html.Class("text-sm text-muted-foreground"),
										g.Text("Click the toggle button to collapse/expand the sidebar."),
									),
									html.P(
										html.Class("text-sm text-muted-foreground"),
										g.Text("The content area automatically adjusts its margin."),
									),
									html.P(
										html.Class("text-sm text-muted-foreground font-medium"),
										g.Text("Toggle icon rotates to show state!"),
									),
								),
							),
						),
					),
				),
			),
		),

		// Feature comparison
		componentSection(
			"Features",
			"All navigation components include",
			html.Div(
				html.Class("grid grid-cols-1 md:grid-cols-3 gap-4"),
				featureCard("Alpine.js State", "Reactive state management for interactivity"),
				featureCard("Keyboard Navigation", "Full keyboard support for accessibility"),
				featureCard("Responsive Design", "Mobile-first with adaptive layouts"),
				featureCard("Smooth Animations", "Transitions using Alpine.js and Tailwind"),
				featureCard("ARIA Attributes", "Proper accessibility labels and roles"),
				featureCard("Modern Styling", "Following shadcn/ui design patterns"),
			),
		),
	)
}

func featureCard(title, description string) g.Node {
	return card.Card(
		card.Content(
			primitives.VStack("2",
				html.Div(
					html.Class("flex items-center gap-2"),
					html.Div(
						html.Class("h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center"),
						g.El("svg",
							html.Class("h-4 w-4 text-primary"),
							g.Attr("xmlns", "http://www.w3.org/2000/svg"),
							g.Attr("fill", "none"),
							g.Attr("viewBox", "0 0 24 24"),
							g.Attr("stroke", "currentColor"),
							g.El("path",
								g.Attr("stroke-linecap", "round"),
								g.Attr("stroke-linejoin", "round"),
								g.Attr("stroke-width", "2"),
								g.Attr("d", "M5 13l4 4L19 7"),
							),
						),
					),
					html.Span(
						html.Class("font-semibold"),
						g.Text(title),
					),
				),
				html.P(
					html.Class("text-sm text-muted-foreground"),
					g.Text(description),
				),
			),
		),
	)
}
