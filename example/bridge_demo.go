package main

import (
	"fmt"
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/bridge"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/icons"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

// handleBridgeDemo demonstrates the Go↔JS bridge functionality (Phase 21)
func handleBridgeDemo(w http.ResponseWriter, r *http.Request, app *forgeui.App, b *bridge.Bridge) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Bridge System Demo")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			alpine.CloakCSS(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),

			// Bridge client script (without Alpine integration - that loads after Alpine.js)
			bridge.BridgeScripts(bridge.ScriptConfig{
				Endpoint:      "/api/bridge/call",
				IncludeAlpine: false,
			}),

			// Hot reload in dev mode
			g.If(app.IsDev(), g.Raw(app.Assets.HotReloadScript())),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground"),
			g.Attr("x-data", `{
				greeting: '',
				sum: null,
				userData: null,
				notifications: [],
				loading: false,
				
				async callGreeting() {
					this.loading = true;
					try {
						const result = await $go('Greet', { name: 'World' });
						this.greeting = result;
					} catch (err) {
						alert('Error: ' + err.message);
					} finally {
						this.loading = false;
					}
				},
				
				async calculateSum() {
					this.loading = true;
					try {
						const result = await $go('Add', { a: 10, b: 32 });
						this.sum = result;
					} catch (err) {
						alert('Error: ' + err.message);
					} finally {
						this.loading = false;
					}
				},
				
				async fetchUser() {
					this.loading = true;
					try {
						const result = await $go('GetUser', { id: 42 });
						this.userData = result;
					} catch (err) {
						alert('Error: ' + err.message);
					} finally {
						this.loading = false;
					}
				},
				
				async startStream() {
					this.loading = true;
					try {
						const result = await $go('GetNotifications', {});
						this.notifications = result;
					} catch (err) {
						alert('Error: ' + err.message);
					} finally {
						this.loading = false;
					}
				}
			}`),
			theme.DarkModeScript(),

			primitives.Container(
				primitives.VStack("8",
					// Header
					html.Div(
						html.Class("py-8"),
						primitives.VStack("4",
							html.H1(
								html.Class("text-4xl font-bold"),
								g.Text("Bridge System Demo"),
							),
							html.P(
								html.Class("text-lg text-muted-foreground"),
								g.Text("Phase 21: Go↔JavaScript Bridge for Server-Side Function Execution"),
							),
							html.Div(
								html.Class("flex gap-2"),
								html.A(
									html.Href("/"),
									html.Class("text-sm text-primary hover:underline flex items-center gap-1"),
									icons.ChevronLeft(icons.WithSize(14)),
									g.Text("Back to Components"),
								),
							),
						),
					),

					// Features Overview
					card.Card(
						card.Header(
							card.Title("Bridge System Features"),
							card.Description("Call Go functions directly from JavaScript"),
						),
						card.Content(
							primitives.VStack("4",
								bridgeFeatureItem(icons.ExternalLink(icons.WithSize(20)), "RPC Bridge", "Call Go functions from JavaScript with $go()"),
								bridgeFeatureItem(icons.Loader(icons.WithSize(20)), "Real-time Streams", "Server-Sent Events for streaming data"),
								bridgeFeatureItem(icons.CheckCircle(icons.WithSize(20)), "Secure by Default", "CSRF protection, rate limiting, authentication"),
								bridgeFeatureItem(icons.Plus(icons.WithSize(20)), "High Performance", "Request batching and caching support"),
								bridgeFeatureItem(icons.CheckCircle(icons.WithSize(20)), "Type Safety", "Automatic parameter validation and type conversion"),
								bridgeFeatureItem(icons.Settings(icons.WithSize(20)), "Alpine Integration", "Magic helpers: $go, $goBatch, $goStream"),
							),
						),
					),

					// Simple Function Call Demo
					card.Card(
						card.Header(
							card.Title("Simple Function Call"),
							card.Description("Call Go function from JavaScript"),
						),
						card.Content(
							primitives.VStack("4",
								html.Div(
									html.Class("space-y-2"),
									html.P(
										html.Class("text-sm text-muted-foreground"),
										g.Text("Click the button to call the Go Greet() function"),
									),
									html.Button(
										html.Class("inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2"),
										g.Attr("@click", "callGreeting()"),
										g.Attr(":disabled", "loading"),
										g.Text("Call Greet()"),
									),
									html.Div(
										g.Attr("x-show", "greeting"),
										html.Class("mt-4 p-4 bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded-lg"),
										html.P(
											html.Class("font-semibold text-green-800 dark:text-green-200"),
											g.Attr("x-text", "greeting"),
										),
									),
								),
								bridgeCodeExample(
									"Go Function",
									`func Greet(ctx bridge.Context, params struct {
  Name string `+"`json:\"name\"`"+`
}) (string, error) {
  return fmt.Sprintf("Hello, %s! 👋", params.Name), nil
}

// Register function
b.Register("Greet", Greet)`,
								),
								bridgeCodeExample(
									"JavaScript Call",
									`// Using the bridge
const result = await $go('Greet', { name: 'World' });
console.log(result); // "Hello, World! 👋"`,
								),
							),
						),
					),

					// Math Function Demo
					card.Card(
						card.Header(
							card.Title("Math Operations"),
							card.Description("Server-side calculations"),
						),
						card.Content(
							primitives.VStack("4",
								html.Div(
									html.Class("space-y-2"),
									html.P(
										html.Class("text-sm text-muted-foreground"),
										g.Text("Calculate 10 + 32 on the server"),
									),
									html.Button(
										html.Class("inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2"),
										g.Attr("@click", "calculateSum()"),
										g.Attr(":disabled", "loading"),
										g.Text("Calculate Sum"),
									),
									html.Div(
										g.Attr("x-show", "sum !== null"),
										html.Class("mt-4 p-4 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg"),
										html.P(
											html.Class("text-blue-800 dark:text-blue-200"),
											html.Span(g.Text("Result: ")),
											html.Span(
												html.Class("font-bold text-xl"),
												g.Attr("x-text", "sum"),
											),
										),
									),
								),
								bridgeCodeExample(
									"Go Function",
									`func Add(ctx bridge.Context, params struct {
  A int `+"`json:\"a\"`"+`
  B int `+"`json:\"b\"`"+`
}) (int, error) {
  return params.A + params.B, nil
}`,
								),
							),
						),
					),

					// Data Fetching Demo
					card.Card(
						card.Header(
							card.Title("Data Fetching"),
							card.Description("Fetch structured data from Go"),
						),
						card.Content(
							primitives.VStack("4",
								html.Div(
									html.Class("space-y-2"),
									html.Button(
										html.Class("inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2"),
										g.Attr("@click", "fetchUser()"),
										g.Attr(":disabled", "loading"),
										g.Text("Fetch User Data"),
									),
									html.Div(
										g.Attr("x-show", "userData"),
										html.Class("mt-4 p-4 bg-muted rounded-lg"),
										html.Pre(
											html.Class("text-sm overflow-x-auto"),
											html.Code(
												g.Attr("x-text", "JSON.stringify(userData, null, 2)"),
											),
										),
									),
								),
								bridgeCodeExample(
									"Go Function",
									`type User struct {
  ID    int    `+"`json:\"id\"`"+`
  Name  string `+"`json:\"name\"`"+`
  Email string `+"`json:\"email\"`"+`
}

func GetUser(ctx bridge.Context, params struct {
  ID int `+"`json:\"id\"`"+`
}) (*User, error) {
  return &User{
    ID:    params.ID,
    Name:  "John Doe",
    Email: "john@example.com",
  }, nil
}`,
								),
							),
						),
					),

					// List/Array Demo
					card.Card(
						card.Header(
							card.Title("Array Data"),
							card.Description("Fetching array of notifications"),
						),
						card.Content(
							primitives.VStack("4",
								html.Div(
									html.Class("space-y-2"),
									html.Button(
										html.Class("inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2"),
										g.Attr("@click", "startStream()"),
										g.Attr(":disabled", "loading"),
										g.Text("Fetch Notifications"),
									),
									html.Div(
										g.Attr("x-show", "notifications.length > 0"),
										html.Class("mt-4 space-y-2"),
										g.El("template",
											g.Attr("x-for", "notification in notifications"),
											g.Attr(":key", "notification"),
											html.Div(
												html.Class("p-3 bg-purple-50 dark:bg-purple-950 border border-purple-200 dark:border-purple-800 rounded-lg"),
												html.P(
													html.Class("text-sm text-purple-800 dark:text-purple-200"),
													g.Attr("x-text", "notification"),
												),
											),
										),
									),
								),
								bridgeCodeExample(
									"Go Function",
									`func GetNotifications(ctx bridge.Context, _ struct{}) ([]string, error) {
  return []string{
    "🚀 Starting process...",
    "⚙️  Processing data...",
    "📊 Analyzing results...",
    "✅ Complete!",
  }, nil
}`,
								),
								bridgeCodeExample(
									"JavaScript Call",
									`const notifications = await $go('GetNotifications', {});
notifications.forEach(n => console.log(n));`,
								),
							),
						),
					),

					// Architecture
					card.Card(
						card.Header(
							card.Title("How It Works"),
						),
						card.Content(
							html.Div(
								html.Class("space-y-4"),
								html.Ol(
									html.Class("list-decimal pl-6 space-y-2 text-sm"),
									html.Li(g.Text("Register Go functions with the bridge")),
									html.Li(g.Text("Bridge analyzes function signatures and parameters")),
									html.Li(g.Text("Client calls $go() with function name and parameters")),
									html.Li(g.Text("Bridge validates, authenticates, and rate-limits request")),
									html.Li(g.Text("Function executes on server with timeout protection")),
									html.Li(g.Text("Result serialized to JSON and returned to client")),
								),
								html.Div(
									html.Class("mt-6 p-4 bg-yellow-50 dark:bg-yellow-950 border border-yellow-200 dark:border-yellow-800 rounded-lg"),
									html.P(
										html.Class("text-sm text-yellow-800 dark:text-yellow-200"),
										html.Strong(g.Text("Security: ")),
										g.Text("CSRF protection, rate limiting, and authentication are built-in and enforced automatically."),
									),
								),
							),
						),
					),
				),
			),

			// Bridge Alpine integration (loads with defer BEFORE Alpine.js)
			html.Script(
				g.Attr("defer", ""),
				g.Raw(bridge.GetAlpineJS()),
			),

			// Alpine.js scripts (loads AFTER bridge plugin via defer)
			alpine.Scripts(),

			// Manually register bridge plugin when Alpine initializes
			html.Script(
				g.Attr("defer", ""),
				g.Raw(`
					// Register bridge plugin before Alpine starts
					document.addEventListener('alpine:init', () => {
						if (window.AlpineBridgePlugin && window.Alpine) {
							window.Alpine.plugin(window.AlpineBridgePlugin);
						}
					});
				`),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

// Demo bridge functions

// Greet returns a greeting message
func Greet(ctx bridge.Context, params struct {
	Name string `json:"name"`
}) (string, error) {
	return fmt.Sprintf("Hello, %s! 👋", params.Name), nil
}

// Add calculates the sum of two numbers
func Add(ctx bridge.Context, params struct {
	A int `json:"a"`
	B int `json:"b"`
}) (int, error) {
	fmt.Println("Add function called with params:", params)
	return params.A + params.B, nil
}

// User represents a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUser fetches user data
func GetUser(ctx bridge.Context, params struct {
	ID int `json:"id"`
}) (*User, error) {
	return &User{
		ID:    params.ID,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil
}

// GetNotifications returns a list of notifications
func GetNotifications(ctx bridge.Context, _ struct{}) ([]string, error) {
	return []string{
		"🚀 Starting process...",
		"⚙️  Processing data...",
		"📊 Analyzing results...",
		"✅ Complete!",
	}, nil
}

// Helper functions for bridge demo

// bridgeFeatureItem creates a feature list item with icon and description
func bridgeFeatureItem(icon g.Node, title, description string) g.Node {
	return html.Div(
		html.Class("flex items-start gap-3"),
		html.Div(
			html.Class("text-primary"),
			icon,
		),
		html.Div(
			html.Class("flex-1"),
			html.H4(
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

// bridgeCodeExample creates a code example card
func bridgeCodeExample(title, code string) g.Node {
	return html.Div(
		html.Class("space-y-2"),
		html.H4(
			html.Class("text-sm font-semibold text-muted-foreground"),
			g.Text(title),
		),
		html.Pre(
			html.Class("p-4 bg-muted rounded-lg text-xs overflow-x-auto"),
			html.Code(g.Text(code)),
		),
	)
}
