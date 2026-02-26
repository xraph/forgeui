package alpine

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

const (
	// AlpineCDN is the default CDN URL for Alpine.js
	AlpineCDN = "https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"
)

// Scripts returns a templ.Component that renders script tags for Alpine.js and any requested plugins.
//
// IMPORTANT: Plugins MUST be loaded BEFORE Alpine.js core.
// This function ensures correct loading order.
//
// Example (in .templ files):
//
//	@alpine.Scripts(alpine.PluginFocus, alpine.PluginCollapse)
func Scripts(plugins ...Plugin) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		// Plugin scripts must load BEFORE Alpine core
		for _, p := range plugins {
			if url := pluginURLs[p]; url != "" {
				if _, err := fmt.Fprintf(w, `<script defer src="%s"></script>`, url); err != nil {
					return err
				}
			}
		}

		// Alpine.js core (must be last)
		_, err := fmt.Fprintf(w, `<script defer src="%s"></script>`, AlpineCDN)
		return err
	})
}

// ScriptsWithVersion returns script tags with a specific Alpine.js version.
func ScriptsWithVersion(version string, plugins ...Plugin) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		for _, p := range plugins {
			if pluginURLs[p] != "" {
				versionedURL := fmt.Sprintf("https://cdn.jsdelivr.net/npm/@alpinejs/%s@%s/dist/cdn.min.js", p, version)
				if _, err := fmt.Fprintf(w, `<script defer src="%s"></script>`, versionedURL); err != nil {
					return err
				}
			}
		}

		alpineURL := fmt.Sprintf("https://cdn.jsdelivr.net/npm/alpinejs@%s/dist/cdn.min.js", version)
		_, err := fmt.Fprintf(w, `<script defer src="%s"></script>`, alpineURL)
		return err
	})
}

// ScriptsImmediate returns script tags for Alpine.js WITHOUT the defer attribute.
func ScriptsImmediate(plugins ...Plugin) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		for _, p := range plugins {
			if url := pluginURLs[p]; url != "" {
				if _, err := fmt.Fprintf(w, `<script src="%s"></script>`, url); err != nil {
					return err
				}
			}
		}

		_, err := fmt.Fprintf(w, `<script src="%s"></script>`, AlpineCDN)
		return err
	})
}

// CloakCSS returns a templ.Component with CSS to prevent flash of unstyled content.
func CloakCSS() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<style>[x-cloak] { display: none !important; }</style>`)
		return err
	})
}

// ScriptsWithNonce adds a nonce attribute to script tags for Content Security Policy.
func ScriptsWithNonce(nonce string, plugins ...Plugin) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		for _, p := range plugins {
			if url := pluginURLs[p]; url != "" {
				if _, err := fmt.Fprintf(w, `<script defer src="%s" nonce="%s"></script>`, url, nonce); err != nil {
					return err
				}
			}
		}

		_, err := fmt.Fprintf(w, `<script defer src="%s" nonce="%s"></script>`, AlpineCDN, nonce)
		return err
	})
}
