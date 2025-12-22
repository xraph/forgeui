// Package analytics provides event tracking integration for ForgeUI.
//
// The analytics plugin hooks into ForgeUI's lifecycle to track page views,
// clicks, form submissions, and custom events.
//
// # Basic Usage
//
//	registry := plugin.NewRegistry()
//	registry.Use(analytics.New(analytics.Config{
//	    Provider: "plausible",
//	    Domain: "example.com",
//	}))
//
// # Features
//
//   - Page view tracking
//   - Click event tracking
//   - Form submission tracking
//   - Custom event emission
//   - Provider adapters (Google Analytics, Plausible, etc.)
//   - Privacy-friendly options (no cookies)
package analytics

import (
	"context"
	"fmt"

	"github.com/xraph/forgeui/plugin"
)

// Analytics plugin implements plugin lifecycle hooks.
type Analytics struct {
	*plugin.PluginBase
	config Config
}

// Config holds analytics configuration.
type Config struct {
	// Provider is the analytics provider ("ga4", "plausible", "matomo", "custom")
	Provider string

	// Domain for the analytics (required for some providers)
	Domain string

	// TrackingID is the tracking ID or measurement ID
	TrackingID string

	// TrackPageViews enables automatic page view tracking
	TrackPageViews bool

	// TrackClicks enables automatic click tracking
	TrackClicks bool

	// TrackForms enables automatic form submission tracking
	TrackForms bool

	// CustomScript for custom analytics implementations
	CustomScript string
}

// DefaultConfig returns default analytics configuration.
func DefaultConfig() Config {
	return Config{
		Provider:       "plausible",
		TrackPageViews: true,
		TrackClicks:    false,
		TrackForms:     false,
	}
}

// New creates a new Analytics plugin.
func New(config Config) *Analytics {
	return &Analytics{
		PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
			Name:        "analytics",
			Version:     "1.0.0",
			Description: "Event tracking and analytics integration",
			Author:      "ForgeUI",
			License:     "MIT",
		}),
		config: config,
	}
}

// Init initializes the analytics plugin.
func (a *Analytics) Init(ctx context.Context, registry *plugin.Registry) error {
	// Register hooks for tracking
	hooks := registry.Hooks()
	
	if a.config.TrackPageViews {
		// Hook into page renders to track page views
		hooks.On(plugin.HookAfterRender, func(hctx *plugin.HookContext) error {
			return a.trackPageView(hctx.Data)
		})
	}

	return nil
}

// Shutdown cleanly shuts down the plugin.
func (a *Analytics) Shutdown(ctx context.Context) error {
	return nil
}

// trackPageView is called after rendering to track page views.
func (a *Analytics) trackPageView(data any) error {
	// This would be called with page context
	// Implementation depends on provider
	return nil
}

// Scripts returns analytics tracking scripts.
func (a *Analytics) Scripts() []plugin.Script {
	scripts := []plugin.Script{}

	switch a.config.Provider {
	case "ga4":
		if a.config.TrackingID != "" {
			scripts = append(scripts, plugin.Script{
				Name: "google-analytics",
				URL:  fmt.Sprintf("https://www.googletagmanager.com/gtag/js?id=%s", a.config.TrackingID),
				Priority: 10,
				Defer:   true,
			})
			scripts = append(scripts, plugin.Script{
				Name:     "ga4-init",
				Inline:   a.ga4InitScript(),
				Priority: 11,
			})
		}

	case "plausible":
		if a.config.Domain != "" {
			scripts = append(scripts, plugin.Script{
				Name:     "plausible",
				URL:      "https://plausible.io/js/script.js",
				Priority: 10,
				Defer:    true,
			})
		}

	case "matomo":
		if a.config.TrackingID != "" {
			scripts = append(scripts, plugin.Script{
				Name:     "matomo",
				Inline:   a.matomoScript(),
				Priority: 10,
			})
		}

	case "custom":
		if a.config.CustomScript != "" {
			scripts = append(scripts, plugin.Script{
				Name:     "custom-analytics",
				Inline:   a.config.CustomScript,
				Priority: 10,
			})
		}
	}

	return scripts
}

// ga4InitScript returns GA4 initialization script.
func (a *Analytics) ga4InitScript() string {
	return fmt.Sprintf(`
		window.dataLayer = window.dataLayer || [];
		function gtag(){dataLayer.push(arguments);}
		gtag('js', new Date());
		gtag('config', '%s');
	`, a.config.TrackingID)
}

// matomoScript returns Matomo tracking script.
func (a *Analytics) matomoScript() string {
	return fmt.Sprintf(`
		var _paq = window._paq = window._paq || [];
		_paq.push(['trackPageView']);
		_paq.push(['enableLinkTracking']);
		(function() {
			var u="%s";
			_paq.push(['setTrackerUrl', u+'matomo.php']);
			_paq.push(['setSiteId', '%s']);
			var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
			g.async=true; g.src=u+'matomo.js'; s.parentNode.insertBefore(g,s);
		})();
	`, a.config.Domain, a.config.TrackingID)
}

// Directives returns analytics directives.
func (a *Analytics) Directives() []plugin.AlpineDirective {
	return nil
}

// Stores returns analytics stores.
func (a *Analytics) Stores() []plugin.AlpineStore {
	return []plugin.AlpineStore{
		{
			Name: "analytics",
			InitialState: map[string]any{
				"provider": a.config.Provider,
			},
			Methods: `
				trackEvent(name, props = {}) {
					switch(this.provider) {
						case 'ga4':
							if (typeof gtag !== 'undefined') {
								gtag('event', name, props);
							}
							break;
						case 'plausible':
							if (typeof plausible !== 'undefined') {
								plausible(name, { props });
							}
							break;
						case 'matomo':
							if (typeof _paq !== 'undefined') {
								_paq.push(['trackEvent', name, props.category || '', props.action || '', props.value || '']);
							}
							break;
					}
				},
				
				trackPageView(path) {
					switch(this.provider) {
						case 'ga4':
							if (typeof gtag !== 'undefined') {
								gtag('event', 'page_view', { page_path: path });
							}
							break;
						case 'plausible':
							if (typeof plausible !== 'undefined') {
								plausible('pageview');
							}
							break;
						case 'matomo':
							if (typeof _paq !== 'undefined') {
								_paq.push(['setCustomUrl', path]);
								_paq.push(['trackPageView']);
							}
							break;
					}
				}
			`,
		},
	}
}

// Magics returns analytics magic properties.
func (a *Analytics) Magics() []plugin.AlpineMagic {
	return nil
}

// AlpineComponents returns analytics components.
func (a *Analytics) AlpineComponents() []plugin.AlpineComponent {
	return nil
}

