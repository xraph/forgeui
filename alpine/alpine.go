// Package alpine provides Alpine.js integration helpers for ForgeUI.
//
// Alpine.js is a lightweight (~15kb) JavaScript framework that provides
// reactive and declarative interactivity to server-rendered HTML.
//
// # Basic Usage
//
// Use Alpine directive helpers to add interactivity to your components:
//
//	html.Div(
//	    alpine.XData(map[string]any{
//	        "count": 0,
//	        "message": "Hello, Alpine!",
//	    }),
//	    html.Button(
//	        alpine.XClick("count++"),
//	        g.Text("Increment"),
//	    ),
//	    html.P(
//	        html.Class("text-xl font-bold"),
//	        alpine.XText("'Count: ' + count"),
//	    ),
//	)
//
// # State Management
//
// Use x-data to define component state:
//
//	alpine.XData(map[string]any{
//	    "open": false,
//	    "items": []any{},
//	})
//
// # Event Handling
//
// Use x-on (or @ shorthand) for event listeners:
//
//	alpine.XClick("doSomething()")
//	alpine.XSubmit("handleSubmit()")
//	alpine.XKeydown("escape", "close()")
//
// # Conditional Rendering
//
// Use x-show for CSS-based visibility or x-if for DOM removal:
//
//	alpine.XShow("isVisible")
//	alpine.XIf("shouldRender")
//
// # List Rendering
//
// Use x-for to iterate over arrays:
//
//	html.Template(
//	    g.Group(alpine.XForKeyed("item in items", "item.id")),
//	    html.Li(alpine.XText("item.name")),
//	)
//
// # Two-Way Binding
//
// Use x-model for form inputs:
//
//	html.Input(
//	    html.Type("text"),
//	    alpine.XModel("name"),
//	)
//
// # Global Stores
//
// Create global reactive state with Alpine stores:
//
//	alpine.RegisterStores(
//	    Store{
//	        Name: "notifications",
//	        State: map[string]any{"items": []any{}},
//	        Methods: `
//	            add(msg) { this.items.push(msg); },
//	            clear() { this.items = []; }
//	        `,
//	    },
//	)
//
// # Plugins
//
// Load Alpine plugins for additional functionality:
//
//	alpine.Scripts(
//	    alpine.PluginFocus,    // Focus management
//	    alpine.PluginCollapse, // Height transitions
//	    alpine.PluginMask,     // Input masking
//	)
//
// # Transitions
//
// Add smooth transitions with x-transition:
//
//	html.Div(
//	    alpine.XShow("open"),
//	    g.Group(alpine.XTransition(animation.FadeIn())),
//	    g.Text("Content"),
//	)
//
// # Routing
//
// Use Pinecone Router for client-side navigation (requires PluginRouter):
//
//	html.Div(
//	    alpine.XData(map[string]any{}),
//	    // Static route with inline template
//	    html.Template(
//	        alpine.XRoute("/"),
//	        alpine.XTemplateInline(),
//	        html.H1(g.Text("Home")),
//	    ),
//	    // Dynamic route with parameters
//	    html.Template(
//	        alpine.XRoute("/users/:id"),
//	        alpine.XTemplate("/views/user.html", TargetID("app")),
//	        alpine.XHandler("loadUser"),
//	    ),
//	    // 404 route
//	    html.Template(
//	        alpine.XRoute("notfound"),
//	        alpine.XTemplate("/views/404.html", Preload()),
//	    ),
//	)
//	
//	// Navigation buttons
//	html.Button(
//	    alpine.XClick(alpine.NavigateTo("/dashboard")),
//	    g.Text("Dashboard"),
//	)
//	html.Button(
//	    alpine.XClick(alpine.RouterBack()),
//	    alpine.XBindDisabled("!"+alpine.RouterCanGoBack()),
//	    g.Text("Back"),
//	)
//
// Load router plugin:
//
//	alpine.Scripts(alpine.PluginRouter)
//
// Router provides magic helpers in Alpine expressions:
//   - $router: Access to PineconeRouter object
//   - $params: Access to route parameters (e.g., $params.id)
//   - $history: Navigation history operations
//
// # Best Practices
//
// 1. Keep component state small and focused
// 2. Use stores for global state
// 3. Prefer x-show over x-if when element will toggle frequently
// 4. Always load plugins BEFORE Alpine core
// 5. Use XCloak() to prevent flash of unstyled content
// 6. Be careful with XHtml() to avoid XSS vulnerabilities
//
// # Plugin Load Order
//
// CRITICAL: Alpine plugins must be loaded BEFORE Alpine core:
//
//	// Correct
//	alpine.Scripts(alpine.PluginFocus) // Loads plugin, then Alpine
//
//	// Wrong - will not work
//	// Load Alpine first, then try to add plugin
//
// The Scripts() function handles this automatically.
package alpine
