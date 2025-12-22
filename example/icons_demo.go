package main

import (
	"net/http"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/badge"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/icons"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

// handleIcons demonstrates the icon system
func handleIcons(w http.ResponseWriter, r *http.Request) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Icons Library")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			alpine.CloakCSS(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground"),
			g.Attr("x-data", ""), // Initialize Alpine.js
			theme.DarkModeScript(),

			// Header
			primitives.Box(
				primitives.WithBackground("bg-card border-b border-border"),
				primitives.WithPadding("py-4"),
				primitives.WithChildren(
					primitives.Container(
						primitives.HStack("4",
							html.Div(
								html.Class("flex-1 flex items-center gap-2"),
								icons.Menu(icons.WithSize(24)),
								primitives.Text(
									primitives.TextAs("span"),
									primitives.TextSize("text-lg"),
									primitives.TextWeight("font-semibold"),
									primitives.TextChildren(g.Text("ForgeUI Icons")),
								),
							),
							primitives.HStack("2",
								html.A(
									html.Href("/"),
									button.Ghost(
										g.Group([]g.Node{
											icons.ChevronLeft(icons.WithSize(16)),
											g.Text("Back"),
										}),
									),
								),
								theme.SimpleToggle(),
							),
						),
					),
				),
			),

			// Main Content
			primitives.Container(
				primitives.Box(
					primitives.WithClass("py-12 space-y-12"),
					primitives.WithChildren(
						// Hero
						primitives.VStack("6",
							primitives.Text(
								primitives.TextAs("h1"),
								primitives.TextSize("text-4xl md:text-5xl"),
								primitives.TextWeight("font-bold"),
								primitives.TextClass("tracking-tight"),
								primitives.TextChildren(g.Text("Icons Library")),
							),
							primitives.Text(
								primitives.TextSize("text-lg"),
								primitives.TextColor("text-muted-foreground"),
								primitives.TextChildren(g.Text("25+ beautiful Lucide icons with customizable size, color, and stroke width")),
							),
							primitives.HStack("2",
								badge.Badge("25+ Icons", badge.WithVariant(forgeui.VariantDefault)),
								badge.Badge("Lucide", badge.WithVariant(forgeui.VariantSecondary)),
								badge.Badge("Customizable", badge.WithVariant(forgeui.VariantOutline)),
							),
						),

						// Features Overview
						card.Card(
							card.Header(
								card.Title("Icon System Features"),
								card.Description("Flexible and powerful icon system built on Lucide icons"),
							),
							card.Content(
								primitives.Grid(
									primitives.GridCols(1),
									primitives.GridColsMD(2),
									primitives.GridColsLG(3),
									primitives.GridGap("4"),
									primitives.GridChildren(
										featureBox(
											icons.CheckCircle(icons.WithSize(24), icons.WithColor("rgb(34, 197, 94)")),
											"Customizable",
											"Size, color, and stroke width options",
										),
										featureBox(
											icons.CheckCircle(icons.WithSize(24), icons.WithColor("rgb(34, 197, 94)")),
											"Inline SVG",
											"No HTTP requests, perfect scaling",
										),
										featureBox(
											icons.CheckCircle(icons.WithSize(24), icons.WithColor("rgb(34, 197, 94)")),
											"Type Safe",
											"Full Go type checking and validation",
										),
										featureBox(
											icons.CheckCircle(icons.WithSize(24), icons.WithColor("rgb(34, 197, 94)")),
											"Accessible",
											"Proper ARIA attributes support",
										),
										featureBox(
											icons.CheckCircle(icons.WithSize(24), icons.WithColor("rgb(34, 197, 94)")),
											"CSS Animatable",
											"Works with transitions and animations",
										),
										featureBox(
											icons.CheckCircle(icons.WithSize(24), icons.WithColor("rgb(34, 197, 94)")),
											"Easy to Use",
											"Simple functional API with options",
										),
									),
								),
							),
						),

						// Navigation Icons
						iconSection(
							"Navigation",
							"Icons for wayfinding and directional cues",
							[]iconDemo{
								{Icon: icons.ChevronUp(), Name: "ChevronUp", Code: "icons.ChevronUp()"},
								{Icon: icons.ChevronDown(), Name: "ChevronDown", Code: "icons.ChevronDown()"},
								{Icon: icons.ChevronLeft(), Name: "ChevronLeft", Code: "icons.ChevronLeft()"},
								{Icon: icons.ChevronRight(), Name: "ChevronRight", Code: "icons.ChevronRight()"},
								{Icon: icons.Home(), Name: "Home", Code: "icons.Home()"},
								{Icon: icons.Menu(), Name: "Menu", Code: "icons.Menu()"},
								{Icon: icons.ExternalLink(), Name: "ExternalLink", Code: "icons.ExternalLink()"},
							},
						),

						// Action Icons
						iconSection(
							"Actions",
							"Icons for user actions and interactions",
							[]iconDemo{
								{Icon: icons.Plus(), Name: "Plus", Code: "icons.Plus()"},
								{Icon: icons.Minus(), Name: "Minus", Code: "icons.Minus()"},
								{Icon: icons.Check(), Name: "Check", Code: "icons.Check()"},
								{Icon: icons.X(), Name: "X", Code: "icons.X()"},
								{Icon: icons.Pencil(), Name: "Pencil", Code: "icons.Pencil()"},
								{Icon: icons.Trash(), Name: "Trash", Code: "icons.Trash()"},
								{Icon: icons.Copy(), Name: "Copy", Code: "icons.Copy()"},
								{Icon: icons.Download(), Name: "Download", Code: "icons.Download()"},
								{Icon: icons.Upload(), Name: "Upload", Code: "icons.Upload()"},
								{Icon: icons.Search(), Name: "Search", Code: "icons.Search()"},
							},
						),

						// Status Icons
						iconSection(
							"Status & Feedback",
							"Icons for displaying status and feedback to users",
							[]iconDemo{
								{Icon: icons.AlertCircle(), Name: "AlertCircle", Code: "icons.AlertCircle()"},
								{Icon: icons.Info(), Name: "Info", Code: "icons.Info()"},
								{Icon: icons.CheckCircle(), Name: "CheckCircle", Code: "icons.CheckCircle()"},
								{Icon: icons.XCircle(), Name: "XCircle", Code: "icons.XCircle()"},
								{Icon: icons.Loader(), Name: "Loader", Code: "icons.Loader()"},
							},
						),

						// User & Communication Icons
						iconSection(
							"User & Communication",
							"Icons for user profiles and communication",
							[]iconDemo{
								{Icon: icons.User(), Name: "User", Code: "icons.User()"},
								{Icon: icons.Mail(), Name: "Mail", Code: "icons.Mail()"},
								{Icon: icons.Eye(), Name: "Eye", Code: "icons.Eye()"},
								{Icon: icons.EyeOff(), Name: "EyeOff", Code: "icons.EyeOff()"},
							},
						),

						// Utility Icons
						iconSection(
							"Utility",
							"General purpose utility icons",
							[]iconDemo{
								{Icon: icons.Settings(), Name: "Settings", Code: "icons.Settings()"},
								{Icon: icons.Calendar(), Name: "Calendar", Code: "icons.Calendar()"},
								{Icon: icons.Clock(), Name: "Clock", Code: "icons.Clock()"},
							},
						),

						// Customization Examples
						card.Card(
							card.Header(
								card.Title("Customization Options"),
								card.Description("Icons can be customized with various options"),
							),
							card.Content(
								primitives.VStack("8",
									// Size variations
									customizationDemo(
										"Size Variations",
										"Control icon size with WithSize option",
										primitives.HStack("4",
											iconWithLabel(icons.Check(icons.WithSize(16)), "16px"),
											iconWithLabel(icons.Check(icons.WithSize(24)), "24px (default)"),
											iconWithLabel(icons.Check(icons.WithSize(32)), "32px"),
											iconWithLabel(icons.Check(icons.WithSize(48)), "48px"),
										),
										`icons.Check(icons.WithSize(16))
icons.Check(icons.WithSize(24))
icons.Check(icons.WithSize(32))
icons.Check(icons.WithSize(48))`,
									),

									// Color variations
									customizationDemo(
										"Color Variations",
										"Set custom colors with WithColor option",
										primitives.HStack("4",
											iconWithLabel(Heart(icons.WithSize(32), icons.WithColor("rgb(239, 68, 68)")), "Red"),
											iconWithLabel(Heart(icons.WithSize(32), icons.WithColor("rgb(34, 197, 94)")), "Green"),
											iconWithLabel(Heart(icons.WithSize(32), icons.WithColor("rgb(59, 130, 246)")), "Blue"),
											iconWithLabel(Heart(icons.WithSize(32), icons.WithColor("rgb(168, 85, 247)")), "Purple"),
										),
										`Heart(icons.WithColor("rgb(239, 68, 68)"))
Heart(icons.WithColor("rgb(34, 197, 94)"))
Heart(icons.WithColor("rgb(59, 130, 246)"))
Heart(icons.WithColor("rgb(168, 85, 247)"))`,
									),

									// Stroke width variations
									customizationDemo(
										"Stroke Width Variations",
										"Adjust line thickness with WithStrokeWidth",
										primitives.HStack("4",
											iconWithLabel(Circle(icons.WithSize(32), icons.WithStrokeWidth(1)), "Thin (1)"),
											iconWithLabel(Circle(icons.WithSize(32), icons.WithStrokeWidth(2)), "Normal (2)"),
											iconWithLabel(Circle(icons.WithSize(32), icons.WithStrokeWidth(3)), "Bold (3)"),
										),
										`Circle(icons.WithStrokeWidth(1))
Circle(icons.WithStrokeWidth(2))
Circle(icons.WithStrokeWidth(3))`,
									),

									// Custom classes
									customizationDemo(
										"Custom Classes",
										"Add Tailwind classes for hover effects and animations",
										primitives.HStack("4",
											button.Outline(
												g.Group([]g.Node{
													icons.Plus(
														icons.WithSize(16),
														icons.WithClass("transition-transform group-hover:rotate-90"),
													),
													g.Text("Hover me"),
												}),
												button.WithClass("group"),
											),
											button.Outline(
												icons.Loader(
													icons.WithSize(16),
													icons.WithClass("animate-spin"),
												),
												button.Disabled(),
											),
										),
										`icons.Plus(
  icons.WithClass("transition-transform group-hover:rotate-90")
)

icons.Loader(
  icons.WithClass("animate-spin")
)`,
									),
								),
							),
						),

						// Usage Examples
						card.Card(
							card.Header(
								card.Title("Usage Examples"),
								card.Description("Common patterns for using icons in your application"),
							),
							card.Content(
								primitives.VStack("6",
									// In buttons
									usageExample(
										"Icons in Buttons",
										primitives.HStack("2",
											button.Primary(
												g.Group([]g.Node{
													icons.Plus(icons.WithSize(16)),
													g.Text("Add Item"),
												}),
											),
											button.Secondary(
												g.Group([]g.Node{
													icons.Download(icons.WithSize(16)),
													g.Text("Download"),
												}),
											),
											button.Destructive(
												g.Group([]g.Node{
													icons.Trash(icons.WithSize(16)),
													g.Text("Delete"),
												}),
											),
											button.IconButton(
												icons.Settings(icons.WithSize(16)),
											),
										),
										`button.Primary(
  g.Group([]g.Node{
    icons.Plus(icons.WithSize(16)),
    g.Text("Add Item"),
  }),
)`,
									),

									// In alerts
									usageExample(
										"Icons in Alerts",
										primitives.VStack("3",
											html.Div(
												html.Class("flex items-center gap-3 p-4 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg"),
												icons.Info(icons.WithSize(20), icons.WithColor("rgb(59, 130, 246)")),
												html.P(
													html.Class("text-sm text-blue-800 dark:text-blue-200"),
													g.Text("This is an informational message"),
												),
											),
											html.Div(
												html.Class("flex items-center gap-3 p-4 bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded-lg"),
												icons.CheckCircle(icons.WithSize(20), icons.WithColor("rgb(34, 197, 94)")),
												html.P(
													html.Class("text-sm text-green-800 dark:text-green-200"),
													g.Text("Operation completed successfully"),
												),
											),
											html.Div(
												html.Class("flex items-center gap-3 p-4 bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded-lg"),
												icons.XCircle(icons.WithSize(20), icons.WithColor("rgb(239, 68, 68)")),
												html.P(
													html.Class("text-sm text-red-800 dark:text-red-200"),
													g.Text("An error occurred"),
												),
											),
										),
										`icons.Info(icons.WithSize(20), icons.WithColor("rgb(59, 130, 246)"))`,
									),

									// With text
									usageExample(
										"Icons with Text",
										primitives.VStack("3",
											html.Div(
												html.Class("flex items-center gap-2"),
												icons.User(icons.WithSize(16)),
												g.Text("John Doe"),
											),
											html.Div(
												html.Class("flex items-center gap-2"),
												icons.Mail(icons.WithSize(16)),
												g.Text("john@example.com"),
											),
											html.Div(
												html.Class("flex items-center gap-2"),
												icons.Calendar(icons.WithSize(16)),
												g.Text("Dec 20, 2025"),
											),
										),
										`html.Div(
  html.Class("flex items-center gap-2"),
  icons.User(icons.WithSize(16)),
  g.Text("John Doe"),
)`,
									),
								),
							),
						),

						// Best Practices
						card.Card(
							card.Header(
								card.Title("Best Practices"),
							),
							card.Content(
								html.Div(
									html.Class("space-y-4"),
									bestPracticeItem("Use appropriate sizes", "16px for inline text, 20-24px for standalone icons, 32px+ for hero sections"),
									bestPracticeItem("Inherit color by default", "Use currentColor (default) to match text color automatically"),
									bestPracticeItem("Add ARIA labels", "For icon-only buttons, add descriptive aria-label attributes"),
									bestPracticeItem("Consistent stroke width", "Stick to default stroke width (2) for visual consistency"),
									bestPracticeItem("Semantic usage", "Choose icons that clearly represent their action or meaning"),
									bestPracticeItem("Performance", "Icons are inline SVG with no HTTP requests and perfect scaling"),
								),
							),
						),
					),
				),
			),

			// Alpine.js scripts
			alpine.Scripts(),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

// Helper function for feature boxes
func featureBox(icon g.Node, title, description string) g.Node {
	return html.Div(
		html.Class("flex flex-col items-start gap-2 p-4 border rounded-lg"),
		icon,
		html.Div(
			html.Class("space-y-1"),
			html.H4(
				html.Class("font-semibold text-sm"),
				g.Text(title),
			),
			html.P(
				html.Class("text-xs text-muted-foreground"),
				g.Text(description),
			),
		),
	)
}

// Helper type for icon demos
type iconDemo struct {
	Icon g.Node
	Name string
	Code string
}

// Helper function for icon sections
func iconSection(title, description string, demos []iconDemo) g.Node {
	iconNodes := make([]g.Node, len(demos))
	for i, demo := range demos {
		iconNodes[i] = iconCard(demo.Icon, demo.Name, demo.Code)
	}

	return card.Card(
		card.Header(
			card.Title(title),
			card.Description(description),
		),
		card.Content(
			primitives.Grid(
				primitives.GridCols(2),
				primitives.GridColsMD(3),
				primitives.GridColsLG(5),
				primitives.GridGap("4"),
				primitives.GridChildren(iconNodes...),
			),
		),
	)
}

// Helper function for individual icon cards
func iconCard(icon g.Node, name, code string) g.Node {
	return html.Div(
		html.Class("flex flex-col items-center gap-3 p-4 border rounded-lg hover:bg-accent transition-colors cursor-pointer group"),
		html.Div(
			html.Class("text-foreground group-hover:scale-110 transition-transform"),
			icon,
		),
		html.Div(
			html.Class("text-center space-y-1"),
			html.P(
				html.Class("text-sm font-medium"),
				g.Text(name),
			),
			html.P(
				html.Class("text-xs text-muted-foreground font-mono"),
				g.Text(code),
			),
		),
	)
}

// Helper function for icon with label
func iconWithLabel(icon g.Node, label string) g.Node {
	return html.Div(
		html.Class("flex flex-col items-center gap-2"),
		icon,
		html.Span(
			html.Class("text-xs text-muted-foreground"),
			g.Text(label),
		),
	)
}

// Helper function for customization demos
func customizationDemo(title, description string, demo g.Node, code string) g.Node {
	return html.Div(
		html.Class("space-y-3"),
		html.Div(
			html.Class("space-y-1"),
			html.H4(
				html.Class("font-semibold"),
				g.Text(title),
			),
			html.P(
				html.Class("text-sm text-muted-foreground"),
				g.Text(description),
			),
		),
		html.Div(
			html.Class("flex items-center justify-center p-6 border rounded-lg bg-muted/30"),
			demo,
		),
		html.Pre(
			html.Class("p-3 bg-muted rounded-lg text-xs overflow-x-auto"),
			html.Code(g.Text(code)),
		),
	)
}

// Helper function for usage examples
func usageExample(title string, demo g.Node, code string) g.Node {
	return html.Div(
		html.Class("space-y-3"),
		html.H4(
			html.Class("font-semibold"),
			g.Text(title),
		),
		html.Div(
			html.Class("p-4 border rounded-lg"),
			demo,
		),
		html.Pre(
			html.Class("p-3 bg-muted rounded-lg text-xs overflow-x-auto"),
			html.Code(g.Text(code)),
		),
	)
}

// Helper function for best practice items
func bestPracticeItem(title, description string) g.Node {
	return html.Div(
		html.Class("flex gap-3"),
		icons.Check(icons.WithSize(20), icons.WithColor("rgb(34, 197, 94)")),
		html.Div(
			html.Class("space-y-1"),
			html.H5(
				html.Class("font-semibold text-sm"),
				g.Text(title),
			),
			html.P(
				html.Class("text-sm text-muted-foreground"),
				g.Text(description),
			),
		),
	)
}

// Heart icon helper (custom icon for color demo)
func Heart(opts ...icons.Option) g.Node {
	return icons.MultiPathIcon([]string{
		"M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z",
	}, opts...)
}

// Circle icon helper (custom icon for stroke demo)
func Circle(opts ...icons.Option) g.Node {
	return icons.Icon("M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10Z", opts...)
}
