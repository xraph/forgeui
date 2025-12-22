package main

import (
	"net/http"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

func handleAssetsDemo(w http.ResponseWriter, r *http.Request, app *forgeui.App) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Assets Pipeline Demo")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			alpine.CloakCSS(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),

			// Demonstrate asset pipeline - CSS
			app.Assets.StyleSheet("css/custom.css"),
			
			// NEW: Hot reload script in development mode
			g.If(app.IsDev(), g.Raw(app.Assets.HotReloadScript())),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground"),
			g.Attr("x-data", ""),
			theme.DarkModeScript(),

			primitives.Container(
				primitives.VStack("8",
					// Header
					html.Div(
						html.Class("py-8"),
						primitives.VStack("4",
							html.H1(
								html.Class("text-4xl font-bold"),
								g.Text("Assets Pipeline Demo"),
							),
							html.P(
								html.Class("text-lg text-muted-foreground"),
								g.Text("ForgeUI's production-ready asset management system"),
							),
							html.Div(
								html.Class("flex gap-2"),
								html.A(
									html.Href("/"),
									html.Class("text-sm text-primary hover:underline"),
									g.Text("← Back to Components"),
								),
							),
						),
					),

					// Hot Reload Demo (Phase 16)
					g.If(app.IsDev(), card.Card(
						card.Header(
							card.Title("🔥 Hot Reload Active"),
							card.Description("Phase 16: Assets Pipeline with Development Server"),
						),
						card.Content(
							primitives.VStack("4",
								html.Div(
									html.Class("bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded-lg p-4"),
									primitives.VStack("2",
										html.P(
											html.Class("font-semibold text-green-800 dark:text-green-200"),
											g.Text("✅ Hot reload is enabled!"),
										),
										html.P(
											html.Class("text-sm text-green-700 dark:text-green-300"),
											g.Text("Edit any .go file and watch this page reload automatically."),
										),
										html.P(
											html.Class("text-xs text-green-600 dark:text-green-400 font-mono"),
											g.Text("Watching: *.go, *.css, *.js files"),
										),
									),
								),
								featureItem("⚡", "File Watcher", "Monitors Go, CSS, and JS files for changes"),
								featureItem("🔄", "Auto Rebuild", "Rebuilds assets when files change (500ms debounce)"),
								featureItem("🌐", "SSE Endpoint", "Server-Sent Events at /_forgeui/reload"),
								featureItem("📦", "Tailwind Processing", "Scans Go files for classes, generates optimized CSS"),
							),
						),
					)),

					// Features Overview
					card.Card(
						card.Header(
							card.Title("Asset Pipeline Features"),
							card.Description("Built-in features for production-ready asset serving"),
						),
						card.Content(
							primitives.VStack("4",
								featureItem("🔒", "Content-Based Fingerprinting", "Automatic cache busting with SHA256 hashes"),
								featureItem("⚡", "Intelligent Caching", "Immutable cache for fingerprinted assets (1 year)"),
								featureItem("🛡️", "Security", "Path traversal protection and proper MIME types"),
								featureItem("📦", "Embedded FS Support", "Single-binary deployments with embed.FS"),
								featureItem("🎯", "Manifest Support", "Pre-computed hashes for production builds"),
								featureItem("🔧", "Dev/Prod Modes", "No fingerprinting in dev, automatic in production"),
							),
						),
					),

					// Code Examples
					card.Card(
						card.Header(
							card.Title("Usage Examples"),
							card.Description("How to use the asset pipeline in your application"),
						),
						card.Content(
							primitives.VStack("6",
								codeExample(
									"Initialize App",
									`app := forgeui.New(
  forgeui.WithDebug(true),
  forgeui.WithAssetPublicDir("public"),
)

// Serve assets
http.Handle("/static/", app.Assets.Handler())`,
								),
								codeExample(
									"CSS Stylesheets",
									`// In your HTML head
app.Assets.StyleSheet("css/app.css")
app.Assets.StyleSheet("css/print.css", 
  assets.WithMedia("print"),
)`,
								),
								codeExample(
									"JavaScript Scripts",
									`// In your HTML
app.Assets.Script("js/app.js", 
  assets.WithDefer(),
)
app.Assets.Script("js/module.js", 
  assets.WithModule(),
)`,
								),
								codeExample(
									"Hot Reload (Phase 16)",
									`// Start dev server with hot reload
if app.IsDev() {
  go app.Assets.StartDevServer(context.Background())
}

// Add SSE endpoint
http.HandleFunc("/_forgeui/reload", 
  app.Assets.SSEHandler().(func(http.ResponseWriter, *http.Request)))

// Inject hot reload script in HTML
html.Head(
  g.If(app.IsDev(), g.Raw(app.Assets.HotReloadScript())),
)`,
								),
								codeExample(
									"Production Build",
									`// Build assets with Tailwind processing
ctx := context.Background()
if err := app.Assets.Build(ctx); err != nil {
  log.Fatal(err)
}

// Load manifest on startup
app := forgeui.New(
  forgeui.WithAssetManifest("dist/manifest.json"),
)`,
								),
							),
						),
					),

					// Demo Custom CSS
					html.Div(
						html.Class("custom-banner"),
						html.H3(
							html.Class("text-2xl font-bold mb-2"),
							g.Text("Custom CSS Loaded via Asset Pipeline"),
						),
						html.P(
							g.Text("This banner is styled with CSS served through the asset pipeline (css/custom.css)"),
						),
					),

					// Benefits
					card.Card(
						card.Header(
							card.Title("Production Benefits"),
						),
						card.Content(
							html.Ul(
								html.Class("list-disc pl-6 space-y-2"),
								html.Li(g.Text("Automatic cache busting prevents stale assets")),
								html.Li(g.Text("1-year immutable cache reduces bandwidth")),
								html.Li(g.Text("Zero-config in development mode")),
								html.Li(g.Text("Single binary deployments with embed.FS")),
								html.Li(g.Text("Thread-safe concurrent asset serving")),
								html.Li(g.Text("Security-hardened against path traversal")),
							),
						),
					),
				),
			),

			// Demonstrate asset pipeline - JavaScript
			app.Assets.Script("js/demo.js"),

			// Alpine.js scripts (using CDN for convenience)
			alpine.Scripts(),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

func featureItem(icon, title, description string) g.Node {
	return html.Div(
		html.Class("flex gap-4"),
		html.Div(
			html.Class("text-3xl"),
			g.Text(icon),
		),
		html.Div(
			html.Class("flex-1"),
			html.H3(
				html.Class("font-semibold text-lg"),
				g.Text(title),
			),
			html.P(
				html.Class("text-muted-foreground"),
				g.Text(description),
			),
		),
	)
}

func codeExample(title, code string) g.Node {
	return html.Div(
		html.H4(
			html.Class("font-semibold mb-2"),
			g.Text(title),
		),
		html.Pre(
			html.Class("bg-muted p-4 rounded-lg overflow-x-auto text-sm"),
			html.Code(
				g.Text(code),
			),
		),
	)
}

