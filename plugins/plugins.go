// Package plugins provides convenience imports for all built-in ForgeUI plugins.
//
// # Built-in Plugins
//
//   - Toast: Notification system with Alpine store
//   - Sortable: Drag-and-drop list reordering
//   - Charts: Data visualization (Line, Bar, Pie, Area, Doughnut)
//   - Analytics: Event tracking integration
//   - SEO: Meta tags and structured data
//   - HTMX: HTMX plugin wrapper
//   - Corporate: Professional theme preset
//
// # Quick Start
//
//	import (
//	    "github.com/xraph/forgeui/plugin"
//	    "github.com/xraph/forgeui/plugins/toast"
//	    "github.com/xraph/forgeui/plugins/charts"
//	)
//
//	func main() {
//	    registry := plugin.NewRegistry()
//	    registry.Use(
//	        toast.New(),
//	        charts.New(),
//	    )
//	    registry.Initialize(context.Background())
//	    // ... rest of app
//	}
package plugins

import (
	"github.com/xraph/forgeui/plugins/analytics"
	"github.com/xraph/forgeui/plugins/charts"
	"github.com/xraph/forgeui/plugins/htmxplugin"
	"github.com/xraph/forgeui/plugins/seo"
	"github.com/xraph/forgeui/plugins/sortable"
	"github.com/xraph/forgeui/plugins/themes/corporate"
	"github.com/xraph/forgeui/plugins/toast"
)

// AllPlugins returns all built-in plugins.
// Useful for quickly setting up a full-featured ForgeUI application.
func AllPlugins() []any {
	return []any{
		toast.New(),
		sortable.New(),
		charts.New(),
		analytics.New(analytics.DefaultConfig()),
		seo.New(),
		htmxplugin.New(),
		corporate.New(),
	}
}

// EssentialPlugins returns essential plugins for most applications.
func EssentialPlugins() []any {
	return []any{
		toast.New(),
		htmxplugin.New(),
		seo.New(),
	}
}

// DataVisualizationPlugins returns plugins for data-heavy applications.
func DataVisualizationPlugins() []any {
	return []any{
		charts.New(),
		sortable.New(),
	}
}

// Re-export plugin constructors for convenience
var (
	// Toast notification system
	NewToast = toast.New

	// Sortable drag-and-drop
	NewSortable = sortable.New

	// Charts data visualization
	NewCharts = charts.New

	// Analytics event tracking
	NewAnalytics = analytics.New

	// SEO meta tags
	NewSEO = seo.New

	// HTMX wrapper
	NewHTMX = htmxplugin.New

	// Corporate theme
	NewCorporateTheme = corporate.New
)
