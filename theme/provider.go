package theme

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ProviderProps defines configuration for the theme provider.
type ProviderProps struct {
	LightTheme       Theme  // Light theme (required)
	DarkTheme        Theme  // Dark theme (required)
	DefaultTheme     string // Default theme if no preference ("light" or "dark")
	EnableSystemSync bool   // Sync with system preference
	StorageKey       string // localStorage key for theme preference
}

// ProviderOption is a functional option for configuring the Provider.
type ProviderOption func(*ProviderProps)

// WithDefaultTheme sets the default theme when no preference is stored.
func WithDefaultTheme(theme string) ProviderOption {
	return func(p *ProviderProps) {
		p.DefaultTheme = theme
	}
}

// WithStorageKey sets a custom localStorage key for theme persistence.
func WithStorageKey(key string) ProviderOption {
	return func(p *ProviderProps) {
		p.StorageKey = key
	}
}

// DisableSystemSync disables automatic syncing with system theme preference.
func DisableSystemSync() ProviderOption {
	return func(p *ProviderProps) {
		p.EnableSystemSync = false
	}
}

// Provider wraps the application with theme support.
// It injects CSS custom properties and dark mode script into the page.
//
// Usage:
//
//	Provider(children...)                           // Uses default light/dark themes
//	Provider(children..., WithDefaultTheme("dark")) // Starts with dark theme
func Provider(children ...g.Node) g.Node {
	return ProviderWithThemes(DefaultLight(), DefaultDark(), children...)
}

// ProviderWithThemes wraps the application with custom light and dark themes.
//
// Usage:
//
//	ProviderWithThemes(theme.BlueLight(), theme.BlueDark(), children...)
func ProviderWithThemes(lightTheme, darkTheme Theme, children ...g.Node) g.Node {
	// Generate CSS for both themes
	css := GenerateCSS(lightTheme, darkTheme)

	return g.Group([]g.Node{
		// Inject theme CSS
		g.El("style",
			g.Attr("data-forgeui-theme", ""),
			g.Raw(css),
		),
		// Add dark mode initialization script
		DarkModeScript(),
		// Add theme cloak styles
		ThemeStyleCloak(),
		// Render children
		g.Group(children),
	})
}

// ProviderWithOptions provides fine-grained control over theme configuration.
//
// Usage:
//
//	ProviderWithOptions(
//	    children,
//	    WithDefaultTheme("dark"),
//	    WithStorageKey("app-theme"),
//	)
func ProviderWithOptions(children []g.Node, opts ...ProviderOption) g.Node {
	props := &ProviderProps{
		LightTheme:       DefaultLight(),
		DarkTheme:        DefaultDark(),
		DefaultTheme:     "light",
		EnableSystemSync: true,
		StorageKey:       "theme",
	}

	for _, opt := range opts {
		opt(props)
	}

	// Generate CSS for both themes
	css := GenerateCSS(props.LightTheme, props.DarkTheme)

	return g.Group([]g.Node{
		// Inject theme CSS
		g.El("style",
			g.Attr("data-forgeui-theme", ""),
			g.Raw(css),
		),
		// Add dark mode initialization script
		DarkModeScriptWithDefault(props.DefaultTheme),
		// Add theme cloak styles
		ThemeStyleCloak(),
		// Render children
		g.Group(children),
	})
}

// HTMLWrapper wraps the entire HTML document with proper theme attributes.
// This is useful for server-rendered applications.
//
// Usage:
//
//	HTMLWrapper(
//	    theme.DefaultLight(),
//	    theme.DefaultDark(),
//	    html.Head(...),
//	    html.Body(...),
//	)
func HTMLWrapper(lightTheme, darkTheme Theme, children ...g.Node) g.Node {
	return html.HTML(
		html.Lang("en"),
		// Class will be set by dark mode script
		html.Class(""),
		g.Group(children),
	)
}

// HeadContent returns common <head> content for theme support.
// Include this in your page's <head> section.
//
// It includes:
//   - Theme color meta tag (adapts to theme)
//   - Color scheme meta tag
//   - Viewport meta tag
func HeadContent(lightTheme, darkTheme Theme) g.Node {
	return g.Group([]g.Node{
		html.Meta(g.Attr("charset", "utf-8")),
		html.Meta(
			g.Attr("name", "viewport"),
			g.Attr("content", "width=device-width, initial-scale=1"),
		),
		// Theme color for browser chrome
		html.Meta(
			g.Attr("name", "theme-color"),
			g.Attr("content", "hsl("+lightTheme.Colors.Background+")"),
			g.Attr("media", "(prefers-color-scheme: light)"),
		),
		html.Meta(
			g.Attr("name", "theme-color"),
			g.Attr("content", "hsl("+darkTheme.Colors.Background+")"),
			g.Attr("media", "(prefers-color-scheme: dark)"),
		),
		// Color scheme
		html.Meta(
			g.Attr("name", "color-scheme"),
			g.Attr("content", "light dark"),
		),
	})
}

// StyleTag returns just the theme CSS as a style tag.
// Useful if you want to manually place the theme CSS.
func StyleTag(lightTheme, darkTheme Theme) g.Node {
	css := GenerateCSS(lightTheme, darkTheme)

	return g.El("style",
		g.Attr("data-forgeui-theme", ""),
		g.Raw(css),
	)
}
