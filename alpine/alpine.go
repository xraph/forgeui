// Package alpine provides Alpine.js integration helpers for ForgeUI.
//
// Alpine.js is a lightweight (~15kb) JavaScript framework that provides
// reactive and declarative interactivity to server-rendered HTML.
//
// # Basic Usage
//
// Use Alpine directive helpers to add interactivity to your templ components:
//
//	<div { alpine.XData(map[string]any{"count": 0, "message": "Hello, Alpine!"})... }>
//	    <button { alpine.XClick("count++")... }>Increment</button>
//	    <p class="text-xl font-bold" { alpine.XText("'Count: ' + count")... }></p>
//	</div>
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
//	<template { alpine.XForKeyed("item in items", "item.id")... }>
//	    <li { alpine.XText("item.name")... }></li>
//	</template>
//
// # Two-Way Binding
//
// Use x-model for form inputs:
//
//	<input type="text" { alpine.XModel("name")... }/>
//
// # Global Stores
//
// Create global reactive state with Alpine stores:
//
//	@alpine.RegisterStores(
//	    alpine.Store{
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
//	@alpine.Scripts(alpine.PluginFocus, alpine.PluginCollapse, alpine.PluginMask)
//
// # Transitions
//
// Add smooth transitions with x-transition:
//
//	<div { alpine.XShow("open")... } { alpine.XTransition(myTransition)... }>
//	    Content
//	</div>
//
// # Routing
//
// Use Pinecone Router for client-side navigation (requires PluginRouter):
//
//	<div { alpine.XData(map[string]any{})... }>
//	    <template { alpine.XRoute("/")... } { alpine.XTemplateInline()... }>
//	        <h1>Home</h1>
//	    </template>
//	</div>
//
// Load router plugin:
//
//	@alpine.Scripts(alpine.PluginRouter)
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
//	@alpine.Scripts(alpine.PluginFocus)
//
// The Scripts() function handles this automatically.
package alpine
