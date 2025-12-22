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
