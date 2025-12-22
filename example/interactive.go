package main

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
	"github.com/xraph/forgeui/components/card"
)

// InteractiveDemoPage creates a page demonstrating Alpine.js integration and animations.
func InteractiveDemoPage() g.Node {
	return html.Div(
		html.Class("space-y-12"),

		// Counter Example
		CounterExample(),

		// Modal Example
		ModalExample(),

		// Dropdown Example
		DropdownExample(),

		// Toast Example
		ToastExample(),

		// Tabs Example
		TabsExample(),
	)
}

// CounterExample demonstrates basic Alpine.js reactivity.
func CounterExample() g.Node {
	return html.Div(
		html.Class("space-y-4"),

		html.H2(
			html.Class("text-2xl font-bold"),
			g.Text("Interactive Counter"),
		),

		html.P(
			html.Class("text-muted-foreground"),
			g.Text("Basic Alpine.js reactivity with state management"),
		),

		card.Card(
			card.Content(
				html.Div(
					// Initialize Alpine state
					alpine.XData(map[string]any{
						"count": 0,
					}),
					html.Class("space-y-4"),

					// Display count with conditional styling
					html.Div(
						html.Class("text-center p-8"),
						html.P(
							alpine.XBindClass("count > 10 ? 'text-green-500' : 'text-gray-900'"),
							html.Class("text-6xl font-bold"),
							alpine.XText("count"),
						),
						html.P(
							html.Class("text-sm text-muted-foreground mt-2"),
							alpine.XText("count === 0 ? 'Start counting!' : count > 10 ? 'Over 10!' : 'Keep going'"),
						),
					),

					// Control buttons
					html.Div(
						html.Class("flex gap-2 justify-center"),
						html.Button(
							html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 bg-primary text-primary-foreground hover:bg-primary/90"),
							alpine.XClick("count++"),
							g.Text("Increment"),
						),
						html.Button(
							html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 bg-secondary text-secondary-foreground hover:bg-secondary/80"),
							alpine.XClick("count--"),
							g.Text("Decrement"),
						),
						html.Button(
							html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 border bg-background hover:bg-accent hover:text-accent-foreground"),
							alpine.XClick("count = 0"),
							g.Text("Reset"),
						),
					),

					// Conditional message with fade animation
					html.Div(
						alpine.XShow("count >= 20"),
						g.Group(alpine.XTransition(animation.FadeIn())),
						html.Class("mt-4 p-4 bg-green-100 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-md text-center"),
						html.P(
							html.Class("text-green-800 dark:text-green-200 font-medium"),
							g.Text("🎉 You've reached 20!"),
						),
					),
				),
			),
		),
	)
}

// ModalExample demonstrates animated modal with backdrop.
func ModalExample() g.Node {
	return html.Div(
		html.Class("space-y-4"),

		html.H2(
			html.Class("text-2xl font-bold"),
			g.Text("Animated Modal"),
		),

		html.P(
			html.Class("text-muted-foreground"),
			g.Text("Modal with fade backdrop and scale animation"),
		),

		card.Card(
			card.Content(
				html.Div(
					// Modal state
					alpine.XData(map[string]any{
						"open": false,
					}),
					html.Class("space-y-4"),

					// Trigger button
					html.Button(
						html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 bg-primary text-primary-foreground hover:bg-primary/90"),
						alpine.XClick("open = true"),
						g.Text("Open Modal"),
					),

					// Modal container
					html.Div(
						alpine.XShow("open"),
						html.Class("fixed inset-0 z-50 flex items-center justify-center p-4"),
						alpine.XOn("keydown.escape.window", "open = false"),

						// Backdrop
						html.Div(
							alpine.XShow("open"),
							g.Group(alpine.XTransition(animation.FadeIn())),
							html.Class("fixed inset-0 bg-black/50"),
							alpine.XOn("click", "open = false"),
						),

						// Modal content
						html.Div(
							alpine.XShow("open"),
							g.Group(alpine.XTransition(animation.ScaleIn())),
							html.Class("relative bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6 space-y-4 z-10"),
							alpine.XOn("click.stop", ""), // Prevent backdrop click

							// Header
							html.Div(
								html.Class("flex items-center justify-between"),
								html.H3(
									html.Class("text-lg font-semibold"),
									g.Text("Modal Title"),
								),
								html.Button(
									html.Class("text-2xl leading-none p-0 h-8 w-8 hover:bg-accent hover:text-accent-foreground rounded-md"),
									alpine.XClick("open = false"),
									g.Text("×"),
								),
							),

							// Content
							html.Div(
								html.Class("space-y-2"),
								html.P(
									html.Class("text-muted-foreground"),
									g.Text("This is an animated modal with a fade backdrop and scale-in animation."),
								),
								html.P(
									html.Class("text-sm text-muted-foreground"),
									g.Text("Press Escape or click outside to close."),
								),
							),

							// Footer
							html.Div(
								html.Class("flex gap-2 justify-end"),
								html.Button(
									html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 border bg-background hover:bg-accent hover:text-accent-foreground"),
									alpine.XClick("open = false"),
									g.Text("Cancel"),
								),
								html.Button(
									html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 bg-primary text-primary-foreground hover:bg-primary/90"),
									alpine.XClick("open = false"),
									g.Text("Confirm"),
								),
							),
						),
					),
				),
			),
		),
	)
}

// DropdownExample demonstrates dropdown with slide animation.
func DropdownExample() g.Node {
	return html.Div(
		html.Class("space-y-4"),

		html.H2(
			html.Class("text-2xl font-bold"),
			g.Text("Animated Dropdown"),
		),

		html.P(
			html.Class("text-muted-foreground"),
			g.Text("Dropdown menu with slide-down animation"),
		),

		card.Card(
			card.Content(
				html.Div(
					alpine.XData(map[string]any{
						"open": false,
					}),
					html.Class("relative inline-block"),
					alpine.XOn("click.outside", "open = false"),

					// Trigger
					html.Button(
						html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 border bg-background hover:bg-accent hover:text-accent-foreground"),
						alpine.XClick("open = !open"),
						g.Text("Open Menu"),
					),

					// Dropdown content
					html.Div(
						alpine.XShow("open"),
						g.Group(alpine.XTransition(animation.SlideUp())),
						html.Class("absolute left-0 mt-2 w-56 bg-white dark:bg-gray-800 rounded-md shadow-lg border dark:border-gray-700 py-1 z-10"),

						html.A(
							html.Href("#"),
							html.Class("block px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"),
							g.Text("Profile"),
						),
						html.A(
							html.Href("#"),
							html.Class("block px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"),
							g.Text("Settings"),
						),
						html.Hr(html.Class("my-1 border-gray-200 dark:border-gray-700")),
						html.A(
							html.Href("#"),
							html.Class("block px-4 py-2 text-sm text-red-600 hover:bg-gray-100 dark:hover:bg-gray-700"),
							g.Text("Logout"),
						),
					),
				),
			),
		),
	)
}

// ToastExample demonstrates toast notifications.
func ToastExample() g.Node {
	return html.Div(
		html.Class("space-y-4"),

		html.H2(
			html.Class("text-2xl font-bold"),
			g.Text("Toast Notifications"),
		),

		html.P(
			html.Class("text-muted-foreground"),
			g.Text("Animated toast notifications sliding from bottom"),
		),

		card.Card(
			card.Content(
				html.Div(
					alpine.XData(map[string]any{
						"showToast": false,
					}),
					html.Class("space-y-4"),

					html.Button(
						html.Class("inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all shrink-0 h-9 px-4 py-2 bg-primary text-primary-foreground hover:bg-primary/90"),
						alpine.XClick("showToast = true; setTimeout(() => showToast = false, 3000)"),
						g.Text("Show Toast"),
					),

					// Toast container (fixed position)
					html.Div(
						alpine.XShow("showToast"),
						g.Group(alpine.XTransition(animation.SlideInFromBottom())),
						html.Class("fixed bottom-4 right-4 bg-white dark:bg-gray-800 rounded-md shadow-lg border dark:border-gray-700 p-4 max-w-sm z-50"),

						html.Div(
							html.Class("flex items-start gap-3"),
							html.Div(
								html.Class("flex-1"),
								html.H4(
									html.Class("font-semibold text-sm"),
									g.Text("Success"),
								),
								html.P(
									html.Class("text-sm text-muted-foreground mt-1"),
									g.Text("Your action was completed successfully!"),
								),
							),
							html.Button(
								html.Class("text-xl leading-none p-0 h-6 w-6 hover:bg-accent hover:text-accent-foreground rounded-md"),
								alpine.XClick("showToast = false"),
								g.Text("×"),
							),
						),
					),
				),
			),
		),
	)
}

// TabsExample demonstrates tabs with conditional content.
func TabsExample() g.Node {
	return html.Div(
		html.Class("space-y-4"),

		html.H2(
			html.Class("text-2xl font-bold"),
			g.Text("Interactive Tabs"),
		),

		html.P(
			html.Class("text-muted-foreground"),
			g.Text("Tabs with fade transitions between content"),
		),

		card.Card(
			card.Content(
				html.Div(
					alpine.XData(map[string]any{
						"activeTab": "tab1",
					}),
					html.Class("space-y-4"),

					// Tab buttons
					html.Div(
						html.Class("flex gap-2 border-b dark:border-gray-700"),

						html.Button(
							alpine.XBindClass("activeTab === 'tab1' ? 'border-b-2 border-primary text-primary' : 'text-muted-foreground'"),
							html.Class("px-4 py-2 -mb-px transition-colors"),
							alpine.XClick("activeTab = 'tab1'"),
							g.Text("Tab 1"),
						),
						html.Button(
							alpine.XBindClass("activeTab === 'tab2' ? 'border-b-2 border-primary text-primary' : 'text-muted-foreground'"),
							html.Class("px-4 py-2 -mb-px transition-colors"),
							alpine.XClick("activeTab = 'tab2'"),
							g.Text("Tab 2"),
						),
						html.Button(
							alpine.XBindClass("activeTab === 'tab3' ? 'border-b-2 border-primary text-primary' : 'text-muted-foreground'"),
							html.Class("px-4 py-2 -mb-px transition-colors"),
							alpine.XClick("activeTab = 'tab3'"),
							g.Text("Tab 3"),
						),
					),

					// Tab content
					html.Div(
						html.Class("pt-4"),

						// Tab 1
						html.Div(
							alpine.XShow("activeTab === 'tab1'"),
							g.Group(alpine.XTransition(animation.FadeIn())),
							html.Class("space-y-2"),
							html.H3(html.Class("font-semibold"), g.Text("Tab 1 Content")),
							html.P(html.Class("text-sm text-muted-foreground"), g.Text("This is the content for tab 1.")),
						),

						// Tab 2
						html.Div(
							alpine.XShow("activeTab === 'tab2'"),
							g.Group(alpine.XTransition(animation.FadeIn())),
							html.Class("space-y-2"),
							html.H3(html.Class("font-semibold"), g.Text("Tab 2 Content")),
							html.P(html.Class("text-sm text-muted-foreground"), g.Text("This is the content for tab 2.")),
						),

						// Tab 3
						html.Div(
							alpine.XShow("activeTab === 'tab3'"),
							g.Group(alpine.XTransition(animation.FadeIn())),
							html.Class("space-y-2"),
							html.H3(html.Class("font-semibold"), g.Text("Tab 3 Content")),
							html.P(html.Class("text-sm text-muted-foreground"), g.Text("This is the content for tab 3.")),
						),
					),
				),
			),
		),
	)
}
