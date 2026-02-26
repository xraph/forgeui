// Package forgeui is an SSR-first UI framework for Go
//
// ForgeUI provides a full-stack web application framework built on templ
// with Alpine.js interactivity, Tailwind CSS styling, and templui components.
//
// Key Features:
//   - SSR-first with progressive enhancement
//   - Type-safe templ templates
//   - templui component library (shadcn-inspired)
//   - Alpine.js and HTMX integration helpers
//   - Functional options pattern for flexible configuration
//   - Plugin system for extensibility
//
// Basic usage:
//
//	app := forgeui.New(
//	    forgeui.WithDebug(true),
//	    forgeui.WithThemeName("default"),
//	)
//
//	app.Get("/", func(ctx *router.PageContext) (templ.Component, error) {
//	    return HomePage(), nil
//	})
//
//	http.ListenAndServe(":3000", app.Handler())
//
// Subpackages:
//   - plugin: Plugin system for extending ForgeUI functionality
//   - router: HTTP routing with layouts, loaders, and middleware
//   - theme: Theme system with dark mode support
//   - alpine: Alpine.js integration helpers
//   - htmx: HTMX attribute helpers
//   - assets: Asset management and fingerprinting
//   - bridge: Go-JS RPC bridge
//   - primitives: Low-level layout primitives
//   - components: templui-based UI component library
//   - icons: Icon system with Lucide icons
package forgeui

// Version is the current version of ForgeUI
const Version = "0.0.3"
