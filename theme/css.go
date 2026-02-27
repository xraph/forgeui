package theme

import (
	"fmt"
	"strings"
)

// GenerateCSS generates CSS custom properties for the given theme.
// It creates both :root variables for light mode and .dark variables for dark mode.
// The generated CSS can be injected into a <style> tag.
func GenerateCSS(lightTheme, darkTheme Theme) string {
	var b strings.Builder

	// Generate :root (light mode) variables
	b.WriteString(":root {\n")
	b.WriteString(generateColorVars(lightTheme.Colors))
	b.WriteString(generateRadiusVars(lightTheme.Radius))
	b.WriteString(generateSpacingVars(lightTheme.Spacing))
	b.WriteString(generateFontSizeVars(lightTheme.FontSize))
	b.WriteString(generateShadowVars(lightTheme.Shadow))
	b.WriteString("}\n\n")

	// Generate .dark (dark mode) variables
	b.WriteString(".dark {\n")
	b.WriteString(generateColorVars(darkTheme.Colors))
	// Dark mode typically only changes colors, not spacing/fonts
	// But we allow full override if needed
	b.WriteString("}\n")

	return b.String()
}

// GenerateLightCSS generates CSS custom properties for a light theme only.
// Useful when you don't need dark mode support.
func GenerateLightCSS(theme Theme) string {
	var b strings.Builder

	b.WriteString(":root {\n")
	b.WriteString(generateColorVars(theme.Colors))
	b.WriteString(generateRadiusVars(theme.Radius))
	b.WriteString(generateSpacingVars(theme.Spacing))
	b.WriteString(generateFontSizeVars(theme.FontSize))
	b.WriteString(generateShadowVars(theme.Shadow))
	b.WriteString("}\n")

	return b.String()
}

// generateColorVars generates CSS custom properties for color tokens in OKLCH format.
func generateColorVars(colors ColorTokens) string {
	var b strings.Builder

	// Core colors
	fmt.Fprintf(&b,"  --background: %s;\n", colors.Background)
	fmt.Fprintf(&b,"  --foreground: %s;\n", colors.Foreground)

	// Card colors
	fmt.Fprintf(&b,"  --card: %s;\n", colors.Card)
	fmt.Fprintf(&b,"  --card-foreground: %s;\n", colors.CardForeground)

	// Popover colors
	fmt.Fprintf(&b,"  --popover: %s;\n", colors.Popover)
	fmt.Fprintf(&b,"  --popover-foreground: %s;\n", colors.PopoverForeground)

	// Primary colors
	fmt.Fprintf(&b,"  --primary: %s;\n", colors.Primary)
	fmt.Fprintf(&b,"  --primary-foreground: %s;\n", colors.PrimaryForeground)

	// Secondary colors
	fmt.Fprintf(&b,"  --secondary: %s;\n", colors.Secondary)
	fmt.Fprintf(&b,"  --secondary-foreground: %s;\n", colors.SecondaryForeground)

	// Muted colors
	fmt.Fprintf(&b,"  --muted: %s;\n", colors.Muted)
	fmt.Fprintf(&b,"  --muted-foreground: %s;\n", colors.MutedForeground)

	// Accent colors
	fmt.Fprintf(&b,"  --accent: %s;\n", colors.Accent)
	fmt.Fprintf(&b,"  --accent-foreground: %s;\n", colors.AccentForeground)

	// Destructive colors
	fmt.Fprintf(&b,"  --destructive: %s;\n", colors.Destructive)
	fmt.Fprintf(&b,"  --destructive-foreground: %s;\n", colors.DestructiveForeground)

	// Success colors
	if colors.Success != "" {
		fmt.Fprintf(&b,"  --success: %s;\n", colors.Success)
	}

	// Border and input colors
	fmt.Fprintf(&b,"  --border: %s;\n", colors.Border)
	fmt.Fprintf(&b,"  --input: %s;\n", colors.Input)
	fmt.Fprintf(&b,"  --ring: %s;\n", colors.Ring)

	// Chart colors for data visualization (1-12)
	if colors.Chart1 != "" {
		fmt.Fprintf(&b,"  --chart-1: %s;\n", colors.Chart1)
	}

	if colors.Chart2 != "" {
		fmt.Fprintf(&b,"  --chart-2: %s;\n", colors.Chart2)
	}

	if colors.Chart3 != "" {
		fmt.Fprintf(&b,"  --chart-3: %s;\n", colors.Chart3)
	}

	if colors.Chart4 != "" {
		fmt.Fprintf(&b,"  --chart-4: %s;\n", colors.Chart4)
	}

	if colors.Chart5 != "" {
		fmt.Fprintf(&b,"  --chart-5: %s;\n", colors.Chart5)
	}

	if colors.Chart6 != "" {
		fmt.Fprintf(&b,"  --chart-6: %s;\n", colors.Chart6)
	}

	if colors.Chart7 != "" {
		fmt.Fprintf(&b,"  --chart-7: %s;\n", colors.Chart7)
	}

	if colors.Chart8 != "" {
		fmt.Fprintf(&b,"  --chart-8: %s;\n", colors.Chart8)
	}

	if colors.Chart9 != "" {
		fmt.Fprintf(&b,"  --chart-9: %s;\n", colors.Chart9)
	}

	if colors.Chart10 != "" {
		fmt.Fprintf(&b,"  --chart-10: %s;\n", colors.Chart10)
	}

	if colors.Chart11 != "" {
		fmt.Fprintf(&b,"  --chart-11: %s;\n", colors.Chart11)
	}

	if colors.Chart12 != "" {
		fmt.Fprintf(&b,"  --chart-12: %s;\n", colors.Chart12)
	}

	// Sidebar colors for independent theming
	if colors.Sidebar != "" {
		fmt.Fprintf(&b,"  --sidebar: %s;\n", colors.Sidebar)
		fmt.Fprintf(&b,"  --sidebar-foreground: %s;\n", colors.SidebarForeground)
		fmt.Fprintf(&b,"  --sidebar-primary: %s;\n", colors.SidebarPrimary)
		fmt.Fprintf(&b,"  --sidebar-primary-foreground: %s;\n", colors.SidebarPrimaryForeground)
		fmt.Fprintf(&b,"  --sidebar-accent: %s;\n", colors.SidebarAccent)
		fmt.Fprintf(&b,"  --sidebar-accent-foreground: %s;\n", colors.SidebarAccentForeground)
		fmt.Fprintf(&b,"  --sidebar-border: %s;\n", colors.SidebarBorder)
		fmt.Fprintf(&b,"  --sidebar-ring: %s;\n", colors.SidebarRing)
	}

	return b.String()
}

// generateRadiusVars generates CSS custom properties for radius tokens.
func generateRadiusVars(radius RadiusTokens) string {
	var b strings.Builder

	if radius.SM != "" {
		fmt.Fprintf(&b,"  --radius-sm: %s;\n", radius.SM)
	}

	if radius.MD != "" {
		fmt.Fprintf(&b,"  --radius-md: %s;\n", radius.MD)
	}

	if radius.LG != "" {
		fmt.Fprintf(&b,"  --radius-lg: %s;\n", radius.LG)
		// Default radius for components (using LG as per modern shadcn/ui)
		fmt.Fprintf(&b,"  --radius: %s;\n", radius.LG)
	}

	if radius.XL != "" {
		fmt.Fprintf(&b,"  --radius-xl: %s;\n", radius.XL)
	}

	if radius.Full != "" {
		fmt.Fprintf(&b,"  --radius-full: %s;\n", radius.Full)
	}

	return b.String()
}

// generateSpacingVars generates CSS custom properties for spacing tokens.
func generateSpacingVars(spacing SpacingTokens) string {
	var b strings.Builder

	if spacing.XS != "" {
		fmt.Fprintf(&b,"  --spacing-xs: %s;\n", spacing.XS)
	}

	if spacing.SM != "" {
		fmt.Fprintf(&b,"  --spacing-sm: %s;\n", spacing.SM)
	}

	if spacing.MD != "" {
		fmt.Fprintf(&b,"  --spacing-md: %s;\n", spacing.MD)
		// Default spacing
		fmt.Fprintf(&b,"  --spacing: %s;\n", spacing.MD)
	}

	if spacing.LG != "" {
		fmt.Fprintf(&b,"  --spacing-lg: %s;\n", spacing.LG)
	}

	if spacing.XL != "" {
		fmt.Fprintf(&b,"  --spacing-xl: %s;\n", spacing.XL)
	}

	if spacing.XXL != "" {
		fmt.Fprintf(&b,"  --spacing-xxl: %s;\n", spacing.XXL)
	}

	return b.String()
}

// generateFontSizeVars generates CSS custom properties for font size tokens.
func generateFontSizeVars(fontSize FontSizeTokens) string {
	var b strings.Builder

	if fontSize.XS != "" {
		fmt.Fprintf(&b,"  --font-size-xs: %s;\n", fontSize.XS)
	}

	if fontSize.SM != "" {
		fmt.Fprintf(&b,"  --font-size-sm: %s;\n", fontSize.SM)
	}

	if fontSize.Base != "" {
		fmt.Fprintf(&b,"  --font-size-base: %s;\n", fontSize.Base)
	}

	if fontSize.LG != "" {
		fmt.Fprintf(&b,"  --font-size-lg: %s;\n", fontSize.LG)
	}

	if fontSize.XL != "" {
		fmt.Fprintf(&b,"  --font-size-xl: %s;\n", fontSize.XL)
	}

	if fontSize.XXL != "" {
		fmt.Fprintf(&b,"  --font-size-xxl: %s;\n", fontSize.XXL)
	}

	if fontSize.XXXL != "" {
		fmt.Fprintf(&b,"  --font-size-xxxl: %s;\n", fontSize.XXXL)
	}

	return b.String()
}

// generateShadowVars generates CSS custom properties for shadow tokens.
func generateShadowVars(shadow ShadowTokens) string {
	var b strings.Builder

	if shadow.SM != "" {
		fmt.Fprintf(&b,"  --shadow-sm: %s;\n", shadow.SM)
	}

	if shadow.MD != "" {
		fmt.Fprintf(&b,"  --shadow-md: %s;\n", shadow.MD)
		// Default shadow
		fmt.Fprintf(&b,"  --shadow: %s;\n", shadow.MD)
	}

	if shadow.LG != "" {
		fmt.Fprintf(&b,"  --shadow-lg: %s;\n", shadow.LG)
	}

	if shadow.XL != "" {
		fmt.Fprintf(&b,"  --shadow-xl: %s;\n", shadow.XL)
	}

	return b.String()
}

// GenerateInputCSS generates a complete Tailwind CSS v4 input.css from the given
// light and dark themes. The output follows the templui pattern:
//
//   - @import "tailwindcss" (v4 import)
//   - @custom-variant dark (dark mode using .dark class)
//   - @theme inline { ... } (maps Tailwind --color-* to CSS variables)
//   - :root { ... } (light theme OKLCH values)
//   - .dark { ... } (dark theme OKLCH values)
//   - @layer base { ... } (default body/border styles)
//
// The generated CSS can be used as input to the Tailwind v4 CLI:
//
//	npx @tailwindcss/cli -i input.css -o output.css
func GenerateInputCSS(lightTheme, darkTheme Theme) string {
	var b strings.Builder

	// Tailwind v4 import
	b.WriteString("@import \"tailwindcss\";\n\n")

	// Dark mode custom variant (uses .dark class)
	b.WriteString("@custom-variant dark (&:where(.dark, .dark *));\n\n")

	// @theme inline — maps Tailwind's --color-* namespace to our --* CSS variables
	b.WriteString("@theme inline {\n")
	b.WriteString("  --breakpoint-3xl: 1600px;\n")
	b.WriteString("  --breakpoint-4xl: 2000px;\n")
	// Radius mappings (relative to base --radius)
	b.WriteString("  --radius-sm: calc(var(--radius) - 4px);\n")
	b.WriteString("  --radius-md: calc(var(--radius) - 2px);\n")
	b.WriteString("  --radius-lg: var(--radius);\n")
	b.WriteString("  --radius-xl: calc(var(--radius) + 4px);\n")
	// Color token mappings
	b.WriteString(generateThemeColorMappings())
	b.WriteString("}\n\n")

	// :root — light theme values (wrapped in oklch())
	b.WriteString(":root {\n")
	fmt.Fprintf(&b,"  --radius: %s;\n", resolveRadius(lightTheme.Radius))
	b.WriteString(generateOKLCHColorVars(lightTheme.Colors))
	b.WriteString("}\n\n")

	// .dark — dark theme values (wrapped in oklch())
	b.WriteString(".dark {\n")
	b.WriteString(generateOKLCHColorVars(darkTheme.Colors))
	b.WriteString("}\n\n")

	// @layer base — default body/border styles
	b.WriteString("@layer base {\n")
	b.WriteString("  * {\n")
	b.WriteString("    @apply border-border;\n")
	b.WriteString("  }\n")
	b.WriteString("  body {\n")
	b.WriteString("    @apply bg-background text-foreground;\n")
	b.WriteString("  }\n")
	b.WriteString("}\n")

	return b.String()
}

// generateThemeColorMappings generates the @theme inline color mappings.
// These map Tailwind's --color-* namespace to our CSS variable names.
func generateThemeColorMappings() string {
	var b strings.Builder

	// Core colors
	mappings := []struct{ tw, css string }{
		{"background", "background"},
		{"foreground", "foreground"},
		{"card", "card"},
		{"card-foreground", "card-foreground"},
		{"popover", "popover"},
		{"popover-foreground", "popover-foreground"},
		{"primary", "primary"},
		{"primary-foreground", "primary-foreground"},
		{"secondary", "secondary"},
		{"secondary-foreground", "secondary-foreground"},
		{"muted", "muted"},
		{"muted-foreground", "muted-foreground"},
		{"accent", "accent"},
		{"accent-foreground", "accent-foreground"},
		{"destructive", "destructive"},
		{"destructive-foreground", "destructive-foreground"},
		{"success", "success"},
		{"border", "border"},
		{"input", "input"},
		{"ring", "ring"},
		{"chart-1", "chart-1"},
		{"chart-2", "chart-2"},
		{"chart-3", "chart-3"},
		{"chart-4", "chart-4"},
		{"chart-5", "chart-5"},
		{"chart-6", "chart-6"},
		{"chart-7", "chart-7"},
		{"chart-8", "chart-8"},
		{"chart-9", "chart-9"},
		{"chart-10", "chart-10"},
		{"chart-11", "chart-11"},
		{"chart-12", "chart-12"},
		{"sidebar", "sidebar"},
		{"sidebar-foreground", "sidebar-foreground"},
		{"sidebar-primary", "sidebar-primary"},
		{"sidebar-primary-foreground", "sidebar-primary-foreground"},
		{"sidebar-accent", "sidebar-accent"},
		{"sidebar-accent-foreground", "sidebar-accent-foreground"},
		{"sidebar-border", "sidebar-border"},
		{"sidebar-ring", "sidebar-ring"},
	}

	for _, m := range mappings {
		fmt.Fprintf(&b,"  --color-%s: var(--%s);\n", m.tw, m.css)
	}

	return b.String()
}

// generateOKLCHColorVars generates CSS custom properties with oklch() wrapped values.
// Unlike generateColorVars which outputs raw values, this wraps each value with oklch()
// for use in compiled Tailwind v4 CSS where values must be complete color definitions.
func generateOKLCHColorVars(colors ColorTokens) string {
	var b strings.Builder

	oklch := func(name, value string) {
		if value != "" {
			fmt.Fprintf(&b,"  --%s: oklch(%s);\n", name, value)
		}
	}

	oklch("background", colors.Background)
	oklch("foreground", colors.Foreground)
	oklch("card", colors.Card)
	oklch("card-foreground", colors.CardForeground)
	oklch("popover", colors.Popover)
	oklch("popover-foreground", colors.PopoverForeground)
	oklch("primary", colors.Primary)
	oklch("primary-foreground", colors.PrimaryForeground)
	oklch("secondary", colors.Secondary)
	oklch("secondary-foreground", colors.SecondaryForeground)
	oklch("muted", colors.Muted)
	oklch("muted-foreground", colors.MutedForeground)
	oklch("accent", colors.Accent)
	oklch("accent-foreground", colors.AccentForeground)
	oklch("destructive", colors.Destructive)
	oklch("destructive-foreground", colors.DestructiveForeground)
	oklch("success", colors.Success)
	oklch("border", colors.Border)
	oklch("input", colors.Input)
	oklch("ring", colors.Ring)

	// Chart colors
	for i, c := range []string{
		colors.Chart1, colors.Chart2, colors.Chart3, colors.Chart4,
		colors.Chart5, colors.Chart6, colors.Chart7, colors.Chart8,
		colors.Chart9, colors.Chart10, colors.Chart11, colors.Chart12,
	} {
		oklch(fmt.Sprintf("chart-%d", i+1), c)
	}

	// Sidebar colors
	oklch("sidebar", colors.Sidebar)
	oklch("sidebar-foreground", colors.SidebarForeground)
	oklch("sidebar-primary", colors.SidebarPrimary)
	oklch("sidebar-primary-foreground", colors.SidebarPrimaryForeground)
	oklch("sidebar-accent", colors.SidebarAccent)
	oklch("sidebar-accent-foreground", colors.SidebarAccentForeground)
	oklch("sidebar-border", colors.SidebarBorder)
	oklch("sidebar-ring", colors.SidebarRing)

	return b.String()
}

// resolveRadius returns the base --radius value from RadiusTokens.
// Defaults to "0.625rem" if no LG radius is configured.
func resolveRadius(radius RadiusTokens) string {
	if radius.LG != "" {
		return radius.LG
	}

	return "0.625rem"
}

// GenerateTailwindConfig generates a Tailwind CSS configuration object
// that extends the default theme with custom properties using OKLCH color format.
// This is useful for integrating with Tailwind CSS configuration files.
func GenerateTailwindConfig() string {
	return `module.exports = {
  darkMode: ["class"],
  theme: {
    extend: {
      colors: {
        background: "oklch(var(--background))",
        foreground: "oklch(var(--foreground))",
        card: {
          DEFAULT: "oklch(var(--card))",
          foreground: "oklch(var(--card-foreground))",
        },
        popover: {
          DEFAULT: "oklch(var(--popover))",
          foreground: "oklch(var(--popover-foreground))",
        },
        primary: {
          DEFAULT: "oklch(var(--primary))",
          foreground: "oklch(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "oklch(var(--secondary))",
          foreground: "oklch(var(--secondary-foreground))",
        },
        muted: {
          DEFAULT: "oklch(var(--muted))",
          foreground: "oklch(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "oklch(var(--accent))",
          foreground: "oklch(var(--accent-foreground))",
        },
        destructive: {
          DEFAULT: "oklch(var(--destructive))",
          foreground: "oklch(var(--destructive-foreground))",
        },
        success: "oklch(var(--success))",
        border: "oklch(var(--border))",
        input: "oklch(var(--input))",
        ring: "oklch(var(--ring))",
        chart: {
          "1": "oklch(var(--chart-1))",
          "2": "oklch(var(--chart-2))",
          "3": "oklch(var(--chart-3))",
          "4": "oklch(var(--chart-4))",
          "5": "oklch(var(--chart-5))",
          "6": "oklch(var(--chart-6))",
          "7": "oklch(var(--chart-7))",
          "8": "oklch(var(--chart-8))",
          "9": "oklch(var(--chart-9))",
          "10": "oklch(var(--chart-10))",
          "11": "oklch(var(--chart-11))",
          "12": "oklch(var(--chart-12))",
        },
        sidebar: {
          DEFAULT: "oklch(var(--sidebar))",
          foreground: "oklch(var(--sidebar-foreground))",
          primary: "oklch(var(--sidebar-primary))",
          "primary-foreground": "oklch(var(--sidebar-primary-foreground))",
          accent: "oklch(var(--sidebar-accent))",
          "accent-foreground": "oklch(var(--sidebar-accent-foreground))",
          border: "oklch(var(--sidebar-border))",
          ring: "oklch(var(--sidebar-ring))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
    },
  },
}`
}
