// Package corporate provides a professional theme preset for business applications.
//
// The corporate theme offers a conservative color palette, formal typography,
// and accessibility-focused design suitable for enterprise applications.
//
// # Basic Usage
//
//	registry := plugin.NewRegistry()
//	registry.Use(corporate.New())
//
// # Features
//
//   - Conservative color palette (blues, grays)
//   - Formal typography (system fonts)
//   - High contrast for accessibility (WCAG AAA)
//   - Print-friendly styles
//   - Professional components
package corporate

import (
	"context"

	"github.com/xraph/forgeui/plugin"
	"github.com/xraph/forgeui/theme"
)

// Corporate theme plugin.
type Corporate struct {
	*plugin.PluginBase
}

// New creates a new Corporate theme plugin.
func New() *Corporate {
	return &Corporate{
		PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
			Name:        "corporate-theme",
			Version:     "1.0.0",
			Description: "Professional theme for business applications",
			Author:      "ForgeUI",
			License:     "MIT",
		}),
	}
}

// Init initializes the corporate theme plugin.
func (c *Corporate) Init(ctx context.Context, registry *plugin.Registry) error {
	return nil
}

// Shutdown cleanly shuts down the plugin.
func (c *Corporate) Shutdown(ctx context.Context) error {
	return nil
}

// Theme returns the corporate theme configuration.
// Returns a light theme suitable for corporate applications.
func (c *Corporate) Theme() *theme.Theme {
	return &theme.Theme{
		Colors: theme.ColorTokens{
			// Conservative blue and gray palette (OKLCH format)
			Background:            "1 0 0",               // Pure white
			Foreground:            "0.1 0 0",             // Near black
			Card:                  "0.98 0 0",            // Light gray
			CardForeground:        "0.1 0 0",             // Near black
			Popover:               "1 0 0",               // Pure white
			PopoverForeground:     "0.1 0 0",             // Near black
			Primary:               "0.4 0.2 264",         // Professional blue
			PrimaryForeground:     "1 0 0",               // Pure white
			Secondary:             "0.5 0 0",             // Neutral gray
			SecondaryForeground:   "1 0 0",               // Pure white
			Muted:                 "0.96 0 0",            // Light gray
			MutedForeground:       "0.5 0 0",             // Neutral gray
			Accent:                "0.92 0 0",            // Light gray
			AccentForeground:      "0.1 0 0",             // Near black
			Destructive:           "0.577 0.245 27.325",  // Red
			DestructiveForeground: "1 0 0",               // Pure white
			Success:               "0.508 0.118 165.612", // Green
			Border:                "0.85 0 0",            // Border gray
			Input:                 "0.92 0 0",            // Input border
			Ring:                  "0.4 0.2 264",         // Professional blue
		},
		Radius: theme.RadiusTokens{
			SM:   "0.125rem", // Conservative corners
			MD:   "0.25rem",
			LG:   "0.375rem",
			XL:   "0.5rem",
			Full: "9999px",
		},
		Spacing:  theme.SpacingTokens{},
		FontSize: theme.FontSizeTokens{},
		Shadow:   theme.ShadowTokens{},
	}
}

// CustomFonts returns additional font configurations.
func (c *Corporate) CustomFonts() []string {
	return []string{
		`-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`,
	}
}

// CustomCSS returns additional CSS for the corporate theme.
func (c *Corporate) CustomCSS() string {
	return `
		/* Corporate Theme Additions */
		
		/* Print styles for reports */
		@media print {
			body {
				background: white;
				color: black;
			}
			
			.no-print {
				display: none !important;
			}
			
			a[href]:after {
				content: " (" attr(href) ")";
			}
		}
		
		/* Professional table styles */
		table.corporate {
			border-collapse: collapse;
			width: 100%;
		}
		
		table.corporate th {
			background-color: #F3F4F6;
			font-weight: 600;
			text-align: left;
			padding: 0.75rem;
			border-bottom: 2px solid #D1D5DB;
		}
		
		table.corporate td {
			padding: 0.75rem;
			border-bottom: 1px solid #E5E7EB;
		}
		
		table.corporate tbody tr:hover {
			background-color: #F9FAFB;
		}
		
		/* Formal button styles */
		.btn-corporate {
			font-weight: 500;
			letter-spacing: 0.025em;
			text-transform: none;
			padding: 0.5rem 1.5rem;
		}
		
		/* Professional card styles */
		.card-corporate {
			box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1);
			border: 1px solid #E5E7EB;
		}
		
		/* Letterhead styles */
		.letterhead {
			padding: 2rem;
			border-bottom: 3px solid #1E40AF;
			margin-bottom: 2rem;
		}
		
		/* Report header */
		.report-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 2rem;
			padding-bottom: 1rem;
			border-bottom: 2px solid #E5E7EB;
		}
		
		/* Professional focus styles for accessibility */
		*:focus-visible {
			outline: 3px solid #1E40AF;
			outline-offset: 2px;
		}
		
		/* High contrast mode support */
		@media (prefers-contrast: high) {
			body {
				border: 2px solid currentColor;
			}
			
			button, a {
				border: 2px solid currentColor;
			}
		}
	`
}

// Scripts returns any additional scripts for the theme.
func (c *Corporate) Scripts() []plugin.Script {
	return nil
}

// Directives returns theme-specific directives.
func (c *Corporate) Directives() []plugin.AlpineDirective {
	return nil
}

// Stores returns theme stores.
func (c *Corporate) Stores() []plugin.AlpineStore {
	return nil
}

// Magics returns theme magic properties.
func (c *Corporate) Magics() []plugin.AlpineMagic {
	return nil
}

// AlpineComponents returns theme components.
func (c *Corporate) AlpineComponents() []plugin.AlpineComponent {
	return nil
}
