// Package forgeui is an SSR-first UI component library for Go
//
// ForgeUI provides type-safe, composable UI components built on gomponents
// with Alpine.js interactivity, Tailwind CSS styling, and shadcn-inspired design.
//
// Key Features:
//   - SSR-first with progressive enhancement
//   - Pure Go API with full type safety
//   - CVA (Class Variance Authority) for variant management
//   - Functional options pattern for flexible configuration
//   - Comprehensive component library
//   - Plugin system for extensibility
//
// Basic usage:
//
//	app := forgeui.New(
//	    forgeui.WithDebug(true),
//	    forgeui.WithThemeName("default"),
//	)
//
//	if err := app.Initialize(context.Background()); err != nil {
//	    log.Fatal(err)
//	}
//
// Subpackages:
//   - plugin: Plugin system for extending ForgeUI functionality
//   - router: HTTP routing with layouts, loaders, and middleware
//   - theme: Theme system with dark mode support
//   - alpine: Alpine.js integration helpers
//   - animation: Animation and transition utilities
//   - assets: Asset management and fingerprinting
//   - primitives: Low-level layout primitives
//   - components: UI component library
//   - icons: Icon system with Lucide icons
package forgeui

// Version is the current version of ForgeUI
const Version = "0.0.2"
