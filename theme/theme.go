// Package theme provides a comprehensive theming system with color tokens,
// dark mode support, and CSS generation following shadcn/ui design patterns.
package theme

// Theme holds all design tokens for the UI system.
// It provides a complete set of color, spacing, typography, and shadow tokens
// that can be converted to CSS custom properties.
type Theme struct {
	Colors   ColorTokens
	Radius   RadiusTokens
	Spacing  SpacingTokens
	FontSize FontSizeTokens
	Shadow   ShadowTokens
}

// ColorTokens defines all color variables following shadcn/ui naming conventions.
// Colors are stored in OKLCH format without the "oklch()" wrapper for compatibility
// with Tailwind CSS custom properties. OKLCH provides better perceptual uniformity
// and more predictable color interpolation compared to HSL.
type ColorTokens struct {
	Background            string // Background color for the app
	Foreground            string // Foreground color for text
	Card                  string // Card background
	CardForeground        string // Card text color
	Popover               string // Popover background
	PopoverForeground     string // Popover text color
	Primary               string // Primary brand color
	PrimaryForeground     string // Text on primary color
	Secondary             string // Secondary color
	SecondaryForeground   string // Text on secondary color
	Muted                 string // Muted background
	MutedForeground       string // Muted text color
	Accent                string // Accent background
	AccentForeground      string // Text on accent color
	Destructive           string // Destructive/danger color
	DestructiveForeground string // Text on destructive color
	Success               string // Success/positive color
	Border                string // Border color
	Input                 string // Input border color
	Ring                  string // Focus ring color
	
	// Chart colors for data visualization
	Chart1  string // Chart color 1
	Chart2  string // Chart color 2
	Chart3  string // Chart color 3
	Chart4  string // Chart color 4
	Chart5  string // Chart color 5
	Chart6  string // Chart color 6
	Chart7  string // Chart color 7
	Chart8  string // Chart color 8
	Chart9  string // Chart color 9
	Chart10 string // Chart color 10
	Chart11 string // Chart color 11
	Chart12 string // Chart color 12
	
	// Sidebar-specific colors for independent theming
	Sidebar                   string // Sidebar background
	SidebarForeground         string // Sidebar text color
	SidebarPrimary            string // Sidebar primary color
	SidebarPrimaryForeground  string // Text on sidebar primary
	SidebarAccent             string // Sidebar accent background
	SidebarAccentForeground   string // Text on sidebar accent
	SidebarBorder             string // Sidebar border color
	SidebarRing               string // Sidebar focus ring color
}

// RadiusTokens defines border radius values for consistent rounded corners.
type RadiusTokens struct {
	SM   string // Small radius (e.g., "0.25rem")
	MD   string // Medium radius (e.g., "0.375rem")
	LG   string // Large radius (e.g., "0.5rem")
	XL   string // Extra large radius (e.g., "0.75rem")
	Full string // Full radius for circles (e.g., "9999px")
}

// SpacingTokens defines consistent spacing values throughout the UI.
type SpacingTokens struct {
	XS  string // Extra small spacing
	SM  string // Small spacing
	MD  string // Medium spacing
	LG  string // Large spacing
	XL  string // Extra large spacing
	XXL string // 2X large spacing
}

// FontSizeTokens defines the typography scale.
type FontSizeTokens struct {
	XS   string // Extra small text
	SM   string // Small text
	Base string // Base font size
	LG   string // Large text
	XL   string // Extra large text
	XXL  string // 2X large text
	XXXL string // 3X large text
}

// ShadowTokens defines elevation shadows for depth perception.
type ShadowTokens struct {
	SM string // Small shadow
	MD string // Medium shadow
	LG string // Large shadow
	XL string // Extra large shadow
}

// New creates a new Theme with default values if not provided.
// This is useful for creating custom themes that inherit defaults.
func New() Theme {
	return Theme{
		Colors:   ColorTokens{},
		Radius:   defaultRadius(),
		Spacing:  defaultSpacing(),
		FontSize: defaultFontSize(),
		Shadow:   defaultShadow(),
	}
}

// defaultRadius returns standard radius tokens matching shadcn/ui.
// Base radius is 0.5rem (8px) for a clean, modern look.
func defaultRadius() RadiusTokens {
	return RadiusTokens{
		SM:   "calc(0.5rem - 2px)",   // 0.375rem (6px)
		MD:   "0.5rem",                // 0.5rem (8px) - base
		LG:   "calc(0.5rem + 2px)",   // 0.625rem (10px)
		XL:   "0.75rem",               // 0.75rem (12px)
		Full: "9999px",
	}
}

// defaultSpacing returns standard spacing tokens.
func defaultSpacing() SpacingTokens {
	return SpacingTokens{
		XS:  "0.25rem",
		SM:  "0.5rem",
		MD:  "1rem",
		LG:  "1.5rem",
		XL:  "2rem",
		XXL: "3rem",
	}
}

// defaultFontSize returns standard font size tokens.
func defaultFontSize() FontSizeTokens {
	return FontSizeTokens{
		XS:   "0.75rem",
		SM:   "0.875rem",
		Base: "1rem",
		LG:   "1.125rem",
		XL:   "1.25rem",
		XXL:  "1.5rem",
		XXXL: "2rem",
	}
}

// defaultShadow returns standard shadow tokens matching shadcn/ui.
// These shadows provide subtle depth without being overpowering.
func defaultShadow() ShadowTokens {
	return ShadowTokens{
		SM: "0 1px 2px 0 rgb(0 0 0 / 0.05)",
		MD: "0 2px 4px -1px rgb(0 0 0 / 0.06), 0 4px 6px -1px rgb(0 0 0 / 0.10)",
		LG: "0 8px 16px -4px rgb(0 0 0 / 0.08), 0 4px 8px -2px rgb(0 0 0 / 0.08)",
		XL: "0 16px 24px -8px rgb(0 0 0 / 0.1), 0 8px 16px -4px rgb(0 0 0 / 0.08)",
	}
}

