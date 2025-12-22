// Package animation provides smooth transitions and animations for ForgeUI components.
//
// This package integrates with Alpine.js's transition system and provides preset
// transitions following shadcn/ui design patterns.
//
// # Basic Usage
//
// Use preset transitions with Alpine.js directives:
//
//	html.Div(
//	    alpine.XShow("open"),
//	    g.Group(alpine.XTransition(animation.FadeIn())),
//	    g.Text("Content"),
//	)
//
// # Preset Transitions
//
// The package provides common transition presets:
//
//	animation.FadeIn()          // Fade in from transparent
//	animation.FadeOut()         // Fade out to transparent
//	animation.ScaleIn()         // Scale and fade in
//	animation.ScaleOut()        // Scale and fade out
//	animation.SlideUp()         // Slide up into view
//	animation.SlideDown()       // Slide down into view
//	animation.SlideLeft()       // Slide from right to left
//	animation.SlideRight()      // Slide from left to right
//	animation.Collapse()        // Smooth height collapse (requires Collapse plugin)
//
// # Custom Transitions
//
// Build custom transitions with the fluent API:
//
//	custom := animation.NewTransition().
//	    Enter("transition-all duration-300").
//	    EnterStart("opacity-0 scale-90").
//	    EnterEnd("opacity-100 scale-100").
//	    Leave("transition-all duration-200").
//	    LeaveStart("opacity-100 scale-100").
//	    LeaveEnd("opacity-0 scale-90").
//	    Build()
//
// # Modal Example
//
// Animated modal with backdrop:
//
//	html.Div(
//	    alpine.XData(map[string]any{"open": false}),
//
//	    // Backdrop with fade
//	    html.Div(
//	        alpine.XShow("open"),
//	        g.Group(alpine.XTransition(animation.FadeIn())),
//	        html.Class("fixed inset-0 bg-black/50"),
//	    ),
//
//	    // Modal with scale
//	    html.Div(
//	        alpine.XShow("open"),
//	        g.Group(alpine.XTransition(animation.ScaleIn())),
//	        html.Class("fixed inset-0 flex items-center justify-center"),
//	        // ... modal content
//	    ),
//	)
//
// # Dropdown Example
//
// Animated dropdown menu:
//
//	html.Div(
//	    alpine.XData(map[string]any{"open": false}),
//
//	    // Trigger
//	    html.Button(alpine.XClick("open = !open"), g.Text("Menu")),
//
//	    // Dropdown content
//	    html.Div(
//	        alpine.XShow("open"),
//	        g.Group(alpine.XTransition(animation.SlideDown())),
//	        html.Class("absolute mt-2 w-48 rounded-md shadow-lg"),
//	        // ... menu items
//	    ),
//	)
//
// # Toast Example
//
// Slide in toast notification:
//
//	html.Div(
//	    alpine.XShow("visible"),
//	    g.Group(alpine.XTransition(animation.SlideInFromBottom())),
//	    html.Class("fixed bottom-4 right-4 p-4 bg-white rounded-md shadow-lg"),
//	    g.Text("Notification message"),
//	)
//
// # Best Practices
//
// 1. Match transition durations to content type:
//   - Fast (150ms): Small UI changes, tooltips
//   - Medium (200-300ms): Modals, dropdowns, most UI
//   - Slow (500ms): Large content, page transitions
//
// 2. Use appropriate easings:
//   - ease-out: Elements entering the view
//   - ease-in: Elements leaving the view
//   - ease-in-out: Elements transforming in place
//
// 3. Keep transitions subtle:
//   - Small scale changes (0.95 - 1.0)
//   - Short distances for slides (10-20px)
//   - Combine transforms for richer effects
//
// 4. Consider accessibility:
//   - Respect prefers-reduced-motion
//   - Provide skip mechanisms for long animations
//   - Ensure content is accessible during transitions
package animation
