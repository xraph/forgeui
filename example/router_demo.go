package main

import (
	"net/http"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/icons"
	"github.com/xraph/forgeui/primitives"
)

// handleRouterDemo demonstrates Pinecone Router integration with Alpine.js
func handleRouterDemo(w http.ResponseWriter, r *http.Request) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("Pinecone Router Demo - ForgeUI")),
			html.Link(
				html.Rel("stylesheet"),
				html.Href("https://cdn.jsdelivr.net/npm/tailwindcss@3/dist/tailwind.min.css"),
			),
			alpine.CloakCSS(),
		),
		html.Body(
			html.Class("bg-gray-50 min-h-screen"),
			routerDemoApp(),
			// Load Alpine.js with Router plugin (router must load before Alpine)
			alpine.Scripts(alpine.PluginRouter),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

func routerDemoApp() g.Node {
	return html.Div(
		html.Class("container mx-auto px-4 py-8 max-w-6xl"),
		alpine.XData(map[string]any{
			"currentUser": nil,
			"users": alpine.RawJS(`[
				{ id: 1, name: 'Alice Johnson', email: 'alice@example.com', role: 'Admin' },
				{ id: 2, name: 'Bob Smith', email: 'bob@example.com', role: 'Developer' },
				{ id: 3, name: 'Charlie Brown', email: 'charlie@example.com', role: 'Designer' }
			]`),
			"posts": alpine.RawJS(`[
				{ id: 101, title: 'Getting Started with Alpine.js', author: 'Alice' },
				{ id: 102, title: 'Understanding Client-Side Routing', author: 'Bob' },
				{ id: 103, title: 'Modern UI Design Patterns', author: 'Charlie' }
			]`),
			"homeHandler": alpine.RawJS(`function(context) {
				console.log('Home route loaded', context);
			}`),
			"userHandler": alpine.RawJS(`function(context) {
				const userId = parseInt(context.params.id);
				this.currentUser = this.users.find(u => u.id === userId) || null;
				if (!this.currentUser) {
					this.$router.navigate('/404');
				}
			}`),
			"usersHandler": alpine.RawJS(`function(context) {
				console.log('Users list loaded');
			}`),
			"postsHandler": alpine.RawJS(`function(context) {
				console.log('Posts loaded');
			}`),
			"postHandler": alpine.RawJS(`function(context) {
				const postId = parseInt(context.params.id);
				const post = this.posts.find(p => p.id === postId);
				if (!post) {
					this.$router.navigate('/404');
				}
			}`),
			"aboutHandler": alpine.RawJS(`function(context) {
				console.log('About page loaded');
			}`),
			"notFoundHandler": alpine.RawJS(`function(context) {
				console.error('404: Page not found -', context.path);
			}`),
		}),

		// Header with navigation
		routerDemoHeader(),

		// Main content area
		html.Div(
			html.ID("app"),
			html.Class("mt-8"),

			// Home route
			html.Template(
				alpine.XRoute("/", alpine.WithRouteName("home")),
				alpine.XTemplateInline(alpine.TargetID("app")),
				alpine.XHandler("homeHandler"),
				routerHomeView(),
			),

			// Users list route
			html.Template(
				alpine.XRoute("/users", alpine.WithRouteName("users")),
				alpine.XTemplateInline(alpine.TargetID("app")),
				alpine.XHandler("usersHandler"),
				routerUsersListView(),
			),

			// User detail route (with parameter)
			html.Template(
				alpine.XRoute("/users/:id"),
				alpine.XTemplateInline(alpine.TargetID("app")),
				alpine.XHandler("userHandler"),
				routerUserDetailView(),
			),

			// Posts list route
			html.Template(
				alpine.XRoute("/posts", alpine.WithRouteName("posts")),
				alpine.XTemplateInline(alpine.TargetID("app")),
				alpine.XHandler("postsHandler"),
				routerPostsListView(),
			),

			// Post detail route (with parameter)
			html.Template(
				alpine.XRoute("/posts/:id"),
				alpine.XTemplateInline(alpine.TargetID("app")),
				alpine.XHandler("postHandler"),
				routerPostDetailView(),
			),

			// About route
			html.Template(
				alpine.XRoute("/about", alpine.WithRouteName("about")),
				alpine.XTemplateInline(alpine.TargetID("app")),
				alpine.XHandler("aboutHandler"),
				routerAboutView(),
			),

			// 404 route
			html.Template(
				alpine.XRoute("notfound"),
				alpine.XTemplateInline(alpine.TargetID("app")),
				alpine.XHandler("notFoundHandler"),
				router404View(),
			),
		),

		// Loading indicator
		html.Div(
			alpine.XShow(alpine.RouterLoading()),
			html.Class("fixed top-4 right-4 bg-blue-500 text-white px-4 py-2 rounded-lg shadow-lg"),
			primitives.HStack("2",
				html.Div(html.Class("animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent")),
				g.Text("Loading..."),
			),
		),
	)
}

func routerDemoHeader() g.Node {
	return html.Header(
		html.Class("bg-white shadow-sm rounded-lg p-6 mb-8"),
		primitives.VStack("6",
			// Title
			primitives.VStack("2",
				html.H1(
					html.Class("text-3xl font-bold text-gray-900"),
					g.Text("Pinecone Router Demo"),
				),
				html.P(
					html.Class("text-gray-600"),
					g.Text("Client-side routing with Alpine.js - Navigate without page reloads"),
				),
			),

			// Navigation buttons
			primitives.HStack("3",
				g.Group([]g.Node{
					button.Button(
						g.Text("Home"),
						button.WithAttrs(alpine.XClick(alpine.NavigateTo("/"))),
						button.WithVariant(forgeui.VariantDefault),
					),
					button.Button(
						g.Text("Users"),
						button.WithAttrs(alpine.XClick(alpine.NavigateTo("/users"))),
						button.WithVariant(forgeui.VariantOutline),
					),
					button.Button(
						g.Text("Posts"),
						button.WithAttrs(alpine.XClick(alpine.NavigateTo("/posts"))),
						button.WithVariant(forgeui.VariantOutline),
					),
					button.Button(
						g.Text("About"),
						button.WithAttrs(alpine.XClick(alpine.NavigateTo("/about"))),
						button.WithVariant(forgeui.VariantOutline),
					),
				}),
			),

			// Browser navigation controls
			html.Div(
				html.Class("flex items-center gap-2 pt-4 border-t"),
				button.Button(
					primitives.HStack("2",
						icons.ArrowLeft(icons.WithSize(16)),
						g.Text("Back"),
					),
					button.WithAttrs(
						alpine.XClick(alpine.RouterBack()),
						alpine.XBindDisabled("!"+alpine.RouterCanGoBack()),
					),
					button.WithVariant(forgeui.VariantGhost),
					button.WithSize(forgeui.SizeSM),
				),
				button.Button(
					primitives.HStack("2",
						g.Text("Forward"),
						icons.ArrowRight(icons.WithSize(16)),
					),
					button.WithAttrs(
						alpine.XClick(alpine.RouterForward()),
						alpine.XBindDisabled("!"+alpine.RouterCanGoForward()),
					),
					button.WithVariant(forgeui.VariantGhost),
					button.WithSize(forgeui.SizeSM),
				),
				html.Div(
					html.Class("ml-auto flex items-center gap-2 text-sm text-gray-500"),
					g.Text("Current path:"),
					html.Code(
						html.Class("px-2 py-1 bg-gray-100 rounded"),
						alpine.XText(alpine.RouterPath()),
					),
				),
			),
		),
	)
}

func routerHomeView() g.Node {
	return card.CardWithOptions(
		[]card.Option{card.WithClass("animate-fade-in")},
		card.Header(
			card.Title("Welcome to Pinecone Router"),
			card.Description("A feature-packed router for Alpine.js applications"),
		),
		card.Content(
			primitives.VStack("6",
				html.P(
					html.Class("text-gray-700 leading-relaxed"),
					g.Text("This demo showcases client-side routing with Pinecone Router. Navigate using the buttons above or try these features:"),
				),

				primitives.VStack("3",
					routerFeatureItem("Dynamic Routes", "Visit user profiles with parameters like /users/1"),
					routerFeatureItem("Browser History", "Use back/forward buttons - they work!"),
					routerFeatureItem("404 Handling", "Try navigating to /invalid-page"),
					routerFeatureItem("Loading States", "Watch the loading indicator in the top-right"),
				),

				html.Div(
					html.Class("grid grid-cols-1 md:grid-cols-3 gap-4 pt-4"),
					routerQuickLinkCard("Users", "Browse user profiles", "/users", "users"),
					routerQuickLinkCard("Posts", "Read blog posts", "/posts", "file-text"),
					routerQuickLinkCard("About", "Learn more", "/about", "info"),
				),
			),
		),
	)
}

func routerUsersListView() g.Node {
	return card.CardWithOptions(
		[]card.Option{card.WithClass("animate-fade-in")},
		card.Header(
			card.Title("Users"),
			card.Description("Click on a user to view their profile"),
		),
		card.Content(
			html.Div(
				html.Class("divide-y divide-gray-200"),
				html.Template(
					g.Group(alpine.XForKeyed("user in users", "user.id")),
					html.Div(
						html.Class("py-4 cursor-pointer hover:bg-gray-50 px-4 -mx-4 rounded-lg transition-colors"),
						alpine.XClick(alpine.NavigateTo("'/users/' + user.id")),
						primitives.HStack("4",
							html.Div(
								html.Class("flex-shrink-0 w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-bold"),
								alpine.XText("user.name.charAt(0)"),
							),
							primitives.VStack("1",
								html.Div(
									html.Class("font-semibold text-gray-900"),
									alpine.XText("user.name"),
								),
								html.Div(
									html.Class("text-sm text-gray-500"),
									alpine.XText("user.email"),
								),
							),
							html.Div(
								html.Class("ml-auto"),
								html.Span(
									alpine.XText("user.role"),
									html.Class("inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold"),
								),
							),
						),
					),
				),
			),
		),
	)
}

func routerUserDetailView() g.Node {
	return g.Group([]g.Node{
		// Show loading state while user data is being fetched
		html.Div(
			alpine.XIf("!currentUser"),
			html.Class("text-center py-12"),
			g.Text("Loading user..."),
		),

		// Show user details when loaded
		html.Div(
			alpine.XIf("currentUser"),
			card.CardWithOptions(
				[]card.Option{card.WithClass("animate-fade-in")},
				card.Header(
					primitives.HStack("4",
						html.Div(
							html.Class("flex-shrink-0 w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-bold text-2xl"),
							alpine.XText("currentUser.name.charAt(0)"),
						),
						primitives.VStack("1",
							html.H3(
								html.Class("text-lg leading-none font-semibold tracking-tight"),
								alpine.XText("currentUser.name"),
							),
							html.P(
								html.Class("text-sm text-muted-foreground leading-relaxed"),
								alpine.XText("currentUser.email"),
							),
						),
						html.Div(
							html.Class("ml-auto"),
							html.Span(
								alpine.XText("currentUser.role"),
								html.Class("inline-flex items-center rounded-md bg-primary text-primary-foreground px-2.5 py-0.5 text-xs font-semibold"),
							),
						),
					),
				),
				card.Content(
					primitives.VStack("6",
						html.Div(
							html.Class("grid grid-cols-2 gap-4"),
							routerInfoItem("User ID", alpine.XText("currentUser.id")),
							routerInfoItem("Role", alpine.XText("currentUser.role")),
							routerInfoItem("Email", alpine.XText("currentUser.email")),
							routerInfoItem("Status", g.Text("Active")),
						),

						html.Div(
							html.Class("flex gap-3 pt-4 border-t"),
							button.Button(
								g.Text("Back to Users"),
								button.WithAttrs(alpine.XClick(alpine.NavigateTo("/users"))),
								button.WithVariant(forgeui.VariantOutline),
							),
							button.Button(
								g.Text("Edit Profile"),
								button.WithVariant(forgeui.VariantDefault),
							),
						),
					),
				),
			),
		),
	})
}

func routerPostsListView() g.Node {
	return card.CardWithOptions(
		[]card.Option{card.WithClass("animate-fade-in")},
		card.Header(
			card.Title("Blog Posts"),
			card.Description("Recent articles and tutorials"),
		),
		card.Content(
			html.Div(
				html.Class("space-y-4"),
				html.Template(
					g.Group(alpine.XForKeyed("post in posts", "post.id")),
					html.Div(
						html.Class("p-4 border border-gray-200 rounded-lg cursor-pointer hover:border-blue-500 hover:shadow-md transition-all"),
						alpine.XClick(alpine.NavigateTo("'/posts/' + post.id")),
						primitives.VStack("2",
							html.H3(
								html.Class("font-semibold text-lg text-gray-900"),
								alpine.XText("post.title"),
							),
							html.Div(
								html.Class("flex items-center gap-2 text-sm text-gray-500"),
								icons.User(icons.WithSize(14)),
								alpine.XText("post.author"),
							),
						),
					),
				),
			),
		),
	)
}

func routerPostDetailView() g.Node {
	return card.CardWithOptions(
		[]card.Option{card.WithClass("animate-fade-in")},
		card.Header(
			primitives.VStack("2",
				html.H3(
					html.Class("text-lg leading-none font-semibold tracking-tight"),
					alpine.XText("posts.find(p => p.id == $params.id)?.title || 'Post Not Found'"),
				),
				html.Div(
					html.Class("flex items-center gap-2 text-sm text-muted-foreground"),
					icons.User(icons.WithSize(14)),
					alpine.XText("posts.find(p => p.id == $params.id)?.author || 'Unknown'"),
				),
			),
		),
		card.Content(
			primitives.VStack("4",
				html.P(
					html.Class("text-gray-700 leading-relaxed"),
					g.Text("This is where the full blog post content would appear. In a real application, you would fetch the post data from an API based on the route parameter."),
				),
				html.P(
					html.Class("text-sm text-gray-500"),
					g.Text("Post ID from URL: "),
					html.Code(
						html.Class("px-2 py-1 bg-gray-100 rounded"),
						alpine.XText("$params.id"),
					),
				),
				button.Button(
					g.Text("Back to Posts"),
					button.WithAttrs(alpine.XClick(alpine.NavigateTo("/posts"))),
					button.WithVariant(forgeui.VariantOutline),
				),
			),
		),
	)
}

func routerAboutView() g.Node {
	return card.CardWithOptions(
		[]card.Option{card.WithClass("animate-fade-in")},
		card.Header(
			card.Title("About Pinecone Router"),
			card.Description("Feature-packed routing for Alpine.js"),
		),
		card.Content(
			primitives.VStack("6",
				html.P(
					html.Class("text-gray-700 leading-relaxed"),
					g.Text("Pinecone Router is a powerful client-side router designed specifically for Alpine.js. It provides a seamless single-page application experience with features like:"),
				),

				primitives.VStack("3",
					routerFeatureItem("Named Parameters", "Extract values from URLs (/users/:id)"),
					routerFeatureItem("Wildcards", "Match multiple segments (/files/*path)"),
					routerFeatureItem("Route Handlers", "Execute code before rendering templates"),
					routerFeatureItem("Template Loading", "Load external HTML templates dynamically"),
					routerFeatureItem("History API", "Full browser navigation support"),
					routerFeatureItem("Magic Helpers", "Access $router, $params, and $history in Alpine"),
				),

				html.Div(
					html.Class("bg-blue-50 border border-blue-200 rounded-lg p-4"),
					primitives.VStack("2",
						html.Div(
							html.Class("font-semibold text-blue-900"),
							g.Text("ForgeUI Integration"),
						),
						html.P(
							html.Class("text-sm text-blue-800"),
							g.Text("This demo shows how Pinecone Router integrates seamlessly with ForgeUI's component system and Alpine.js helpers."),
						),
					),
				),
			),
		),
	)
}

func router404View() g.Node {
	return card.CardWithOptions(
		[]card.Option{card.WithClass("animate-fade-in text-center")},
		card.Content(
			primitives.VStack("6",
				html.Div(
					html.Class("text-6xl font-bold text-gray-300"),
					g.Text("404"),
				),
				primitives.VStack("2",
					html.H2(
						html.Class("text-2xl font-bold text-gray-900"),
						g.Text("Page Not Found"),
					),
					html.P(
						html.Class("text-gray-600"),
						g.Text("The page you're looking for doesn't exist."),
					),
				),
				html.Div(
					html.Class("text-sm text-gray-500"),
					g.Text("Attempted path: "),
					html.Code(
						html.Class("px-2 py-1 bg-gray-100 rounded"),
						alpine.XText(alpine.RouterPath()),
					),
				),
				button.Button(
					g.Text("Go Home"),
					button.WithAttrs(alpine.XClick(alpine.NavigateTo("/"))),
					button.WithVariant(forgeui.VariantDefault),
				),
			),
		),
	)
}

// Helper components

func routerFeatureItem(title, description string) g.Node {
	return primitives.HStack("3",
		icons.CheckCircle(icons.WithSize(20), icons.WithClass("text-green-500 flex-shrink-0")),
		primitives.VStack("0.5",
			html.Span(html.Class("font-medium text-gray-900"), g.Text(title)),
			html.Span(html.Class("text-sm text-gray-600"), g.Text(description)),
		),
	)
}

func routerQuickLinkCard(title, description, path, iconName string) g.Node {
	var iconNode g.Node
	switch iconName {
	case "users":
		iconNode = icons.Users(icons.WithSize(24), icons.WithClass("text-blue-500"))
	case "file-text":
		iconNode = icons.FileText(icons.WithSize(24), icons.WithClass("text-blue-500"))
	case "info":
		iconNode = icons.Info(icons.WithSize(24), icons.WithClass("text-blue-500"))
	default:
		iconNode = icons.Info(icons.WithSize(24), icons.WithClass("text-blue-500"))
	}
	
	return html.Div(
		html.Class("p-4 border border-gray-200 rounded-lg cursor-pointer hover:border-blue-500 hover:shadow-md transition-all"),
		alpine.XClick(alpine.NavigateTo(path)),
		primitives.VStack("3",
			iconNode,
			primitives.VStack("1",
				html.Div(html.Class("font-semibold text-gray-900"), g.Text(title)),
				html.Div(html.Class("text-sm text-gray-600"), g.Text(description)),
			),
		),
	)
}

func routerInfoItem(label string, valueNode g.Node) g.Node {
	return primitives.VStack("1",
		html.Div(html.Class("text-sm font-medium text-gray-500"), g.Text(label)),
		html.Div(html.Class("text-base text-gray-900"), valueNode),
	)
}
