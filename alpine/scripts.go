package alpine

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

const (
	// AlpineCDN is the default CDN URL for Alpine.js
	// Using 3.x.x for automatic minor/patch updates
	AlpineCDN = "https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"
)

// Scripts returns script tags for Alpine.js and any requested plugins.
//
// IMPORTANT: Plugins MUST be loaded BEFORE Alpine.js core.
// This function ensures correct loading order.
//
// Example:
//
//	alpine.Scripts(alpine.PluginFocus, alpine.PluginCollapse)
//
// Generates:
//
//	<script defer src="...focus@3.x.x..."></script>
//	<script defer src="...collapse@3.x.x..."></script>
//	<script defer src="...alpinejs@3.x.x..."></script>
func Scripts(plugins ...Plugin) g.Node {
	scripts := make([]g.Node, 0, len(plugins)+1)

	// Plugin scripts must load BEFORE Alpine core
	for _, p := range plugins {
		if url := pluginURLs[p]; url != "" {
			scripts = append(scripts, html.Script(
				g.Attr("defer", ""),
				g.Attr("src", url),
			))
		}
	}

	// Alpine.js core (must be last)
	scripts = append(scripts, html.Script(
		g.Attr("defer", ""),
		g.Attr("src", AlpineCDN),
	))

	return g.Group(scripts)
}

// ScriptsWithVersion returns script tags with a specific Alpine.js version.
// Useful for pinning to a specific version in production.
//
// Example:
//
//	alpine.ScriptsWithVersion("3.13.3", alpine.PluginFocus)
func ScriptsWithVersion(version string, plugins ...Plugin) g.Node {
	scripts := make([]g.Node, 0, len(plugins)+1)

	// Plugin scripts
	for _, p := range plugins {
		if url := pluginURLs[p]; url != "" {
			// Replace version in URL
			versionedURL := fmt.Sprintf("https://cdn.jsdelivr.net/npm/@alpinejs/%s@%s/dist/cdn.min.js", p, version)
			scripts = append(scripts, html.Script(
				g.Attr("defer", ""),
				g.Attr("src", versionedURL),
			))
		}
	}

	// Alpine.js core with specific version
	alpineURL := fmt.Sprintf("https://cdn.jsdelivr.net/npm/alpinejs@%s/dist/cdn.min.js", version)
	scripts = append(scripts, html.Script(
		g.Attr("defer", ""),
		g.Attr("src", alpineURL),
	))

	return g.Group(scripts)
}

// ScriptsImmediate returns script tags for Alpine.js WITHOUT the defer attribute.
//
// Use this when placing scripts at the END of the <body> tag with inline store
// registrations. The scripts will execute immediately in document order.
//
// IMPORTANT: Only use this when Alpine stores are registered inline in the body.
// If stores are not needed or scripts are in <head>, use Scripts() instead.
//
// Example:
//
//	html.Body(
//	    toast.RegisterToastStore(),  // Inline store registration
//	    // ... content
//	    alpine.ScriptsImmediate(alpine.PluginFocus),  // Execute immediately
//	)
func ScriptsImmediate(plugins ...Plugin) g.Node {
	scripts := make([]g.Node, 0, len(plugins)+1)

	// Plugin scripts must load BEFORE Alpine core
	for _, p := range plugins {
		if url := pluginURLs[p]; url != "" {
			scripts = append(scripts, html.Script(
				g.Attr("src", url),
			))
		}
	}

	// Alpine.js core (must be last)
	scripts = append(scripts, html.Script(
		g.Attr("src", AlpineCDN),
	))

	return g.Group(scripts)
}

// CloakCSS returns a style tag with CSS to prevent flash of unstyled content.
// Add this to your <head> section when using x-cloak.
//
// Example:
//
//	html.Head(
//	    alpine.CloakCSS(),
//	    // ... other head elements
//	)
func CloakCSS() g.Node {
	return html.StyleEl(
		g.Raw("[x-cloak] { display: none !important; }"),
	)
}

// CSPNonce adds a nonce attribute to script tags for Content Security Policy.
// Use this when you have CSP enabled.
//
// Example:
//
//	alpine.ScriptsWithNonce("random-nonce-value", alpine.PluginFocus)
func ScriptsWithNonce(nonce string, plugins ...Plugin) g.Node {
	scripts := make([]g.Node, 0, len(plugins)+1)

	// Plugin scripts
	for _, p := range plugins {
		if url := pluginURLs[p]; url != "" {
			scripts = append(scripts, html.Script(
				g.Attr("defer", ""),
				g.Attr("src", url),
				g.Attr("nonce", nonce),
			))
		}
	}

	// Alpine.js core
	scripts = append(scripts, html.Script(
		g.Attr("defer", ""),
		g.Attr("src", AlpineCDN),
		g.Attr("nonce", nonce),
	))

	return g.Group(scripts)
}
