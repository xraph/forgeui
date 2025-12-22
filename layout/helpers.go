package layout

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/theme"
)

// Theme injects theme system (styles and scripts) into the layout
// Supports both light and dark themes
func Theme(lightTheme, darkTheme *theme.Theme) g.Node {
	if lightTheme == nil && darkTheme == nil {
		return g.Group(nil)
	}

	nodes := []g.Node{}

	// Add theme head content
	if lightTheme != nil && darkTheme != nil {
		nodes = append(nodes, theme.HeadContent(*lightTheme, *darkTheme))
		nodes = append(nodes, theme.StyleTag(*lightTheme, *darkTheme))
	} else if lightTheme != nil {
		// Use same theme for both light and dark if only one provided
		nodes = append(nodes, theme.HeadContent(*lightTheme, *lightTheme))
		nodes = append(nodes, theme.StyleTag(*lightTheme, *lightTheme))
	}

	// Add Tailwind config script
	nodes = append(nodes, theme.TailwindConfigScript())

	return g.Group(nodes)
}

// ThemeSingle injects a single theme into the layout (uses same theme for light and dark)
func ThemeSingle(t *theme.Theme) g.Node {
	if t == nil {
		return g.Group(nil)
	}

	return g.Group([]g.Node{
		theme.HeadContent(*t, *t),
		theme.StyleTag(*t, *t),
		theme.TailwindConfigScript(),
	})
}

// DarkModeScript injects the dark mode toggle script into the body
func DarkModeScript() g.Node {
	return theme.DarkModeScript()
}

// Alpine injects Alpine.js cloak CSS into the head
// This prevents flash of unstyled content before Alpine initializes
func Alpine() g.Node {
	return alpine.CloakCSS()
}

// AlpineScripts injects Alpine.js scripts at the end of body
// Accepts optional plugins
func AlpineScripts(plugins ...alpine.Plugin) g.Node {
	return alpine.Scripts(plugins...)
}

// HotReload injects the hot reload script for development mode
// This should only be used in development
func HotReload() g.Node {
	return g.Raw(`
<script>
	// ForgeUI Hot Reload
	(function() {
		const source = new EventSource('/_forgeui/reload');
		source.onmessage = function(e) {
			if (e.data === 'reload') {
				console.log('[ForgeUI] Reloading...');
				window.location.reload();
			}
		};
		source.onerror = function() {
			console.log('[ForgeUI] Hot reload disconnected');
		};
	})();
</script>
`)
}

// BridgeClient injects the ForgeUI bridge client script
func BridgeClient() g.Node {
	return html.Script(
		html.Src("/static/js/forge-bridge.js"),
		html.Type("module"),
	)
}

// AlpineBridgeClient injects the Alpine.js bridge client script
func AlpineBridgeClient() g.Node {
	return html.Script(
		html.Src("/static/js/alpine-bridge.js"),
		html.Type("module"),
	)
}

// Favicon adds a favicon link to the head
func Favicon(href string) g.Node {
	return html.Link(
		html.Rel("icon"),
		html.Href(href),
	)
}

// OpenGraph adds Open Graph meta tags for social media sharing
func OpenGraph(title, description, image, url string) g.Node {
	return g.Group([]g.Node{
		html.Meta(g.Attr("property", "og:title"), html.Content(title)),
		html.Meta(g.Attr("property", "og:description"), html.Content(description)),
		g.If(image != "", html.Meta(g.Attr("property", "og:image"), html.Content(image))),
		g.If(url != "", html.Meta(g.Attr("property", "og:url"), html.Content(url))),
		html.Meta(g.Attr("property", "og:type"), html.Content("website")),
	})
}

// TwitterCard adds Twitter Card meta tags
func TwitterCard(title, description, image string) g.Node {
	return g.Group([]g.Node{
		html.Meta(html.Name("twitter:card"), html.Content("summary_large_image")),
		html.Meta(html.Name("twitter:title"), html.Content(title)),
		html.Meta(html.Name("twitter:description"), html.Content(description)),
		g.If(image != "", html.Meta(html.Name("twitter:image"), html.Content(image))),
	})
}

// GoogleFonts adds Google Fonts preconnect and stylesheet links
func GoogleFonts(families ...string) g.Node {
	if len(families) == 0 {
		return g.Group(nil)
	}

	// Build font URL
	fontURL := "https://fonts.googleapis.com/css2?"
	for i, family := range families {
		if i > 0 {
			fontURL += "&"
		}
		fontURL += "family=" + family
	}
	fontURL += "&display=swap"

	return g.Group([]g.Node{
		html.Link(html.Rel("preconnect"), html.Href("https://fonts.googleapis.com")),
		html.Link(html.Rel("preconnect"), html.Href("https://fonts.gstatic.com"), g.Attr("crossorigin", "")),
		html.Link(html.Rel("stylesheet"), html.Href(fontURL)),
	})
}

// Preload adds a preload link for critical resources
func Preload(href, as string) g.Node {
	return html.Link(
		html.Rel("preload"),
		html.Href(href),
		g.Attr("as", as),
	)
}

// PreloadFont adds a preload link for fonts
func PreloadFont(href string) g.Node {
	return html.Link(
		html.Rel("preload"),
		html.Href(href),
		g.Attr("as", "font"),
		g.Attr("type", "font/woff2"),
		g.Attr("crossorigin", ""),
	)
}

// ProviderScriptUtilities injects the provider utility scripts
func ProviderScriptUtilities() g.Node {
	// This would import from primitives package
	// For now, return empty to avoid circular dependency
	// Users can manually add primitives.ProviderScriptUtilities() in their layouts
	return g.Group(nil)
}
