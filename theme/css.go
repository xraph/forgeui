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
	b.WriteString(fmt.Sprintf("  --background: %s;\n", colors.Background))
	b.WriteString(fmt.Sprintf("  --foreground: %s;\n", colors.Foreground))

	// Card colors
	b.WriteString(fmt.Sprintf("  --card: %s;\n", colors.Card))
	b.WriteString(fmt.Sprintf("  --card-foreground: %s;\n", colors.CardForeground))

	// Popover colors
	b.WriteString(fmt.Sprintf("  --popover: %s;\n", colors.Popover))
	b.WriteString(fmt.Sprintf("  --popover-foreground: %s;\n", colors.PopoverForeground))

	// Primary colors
	b.WriteString(fmt.Sprintf("  --primary: %s;\n", colors.Primary))
	b.WriteString(fmt.Sprintf("  --primary-foreground: %s;\n", colors.PrimaryForeground))

	// Secondary colors
	b.WriteString(fmt.Sprintf("  --secondary: %s;\n", colors.Secondary))
	b.WriteString(fmt.Sprintf("  --secondary-foreground: %s;\n", colors.SecondaryForeground))

	// Muted colors
	b.WriteString(fmt.Sprintf("  --muted: %s;\n", colors.Muted))
	b.WriteString(fmt.Sprintf("  --muted-foreground: %s;\n", colors.MutedForeground))

	// Accent colors
	b.WriteString(fmt.Sprintf("  --accent: %s;\n", colors.Accent))
	b.WriteString(fmt.Sprintf("  --accent-foreground: %s;\n", colors.AccentForeground))

	// Destructive colors
	b.WriteString(fmt.Sprintf("  --destructive: %s;\n", colors.Destructive))
	b.WriteString(fmt.Sprintf("  --destructive-foreground: %s;\n", colors.DestructiveForeground))

	// Success colors
	if colors.Success != "" {
		b.WriteString(fmt.Sprintf("  --success: %s;\n", colors.Success))
	}

	// Border and input colors
	b.WriteString(fmt.Sprintf("  --border: %s;\n", colors.Border))
	b.WriteString(fmt.Sprintf("  --input: %s;\n", colors.Input))
	b.WriteString(fmt.Sprintf("  --ring: %s;\n", colors.Ring))

	// Chart colors for data visualization (1-12)
	if colors.Chart1 != "" {
		b.WriteString(fmt.Sprintf("  --chart-1: %s;\n", colors.Chart1))
	}
	if colors.Chart2 != "" {
		b.WriteString(fmt.Sprintf("  --chart-2: %s;\n", colors.Chart2))
	}
	if colors.Chart3 != "" {
		b.WriteString(fmt.Sprintf("  --chart-3: %s;\n", colors.Chart3))
	}
	if colors.Chart4 != "" {
		b.WriteString(fmt.Sprintf("  --chart-4: %s;\n", colors.Chart4))
	}
	if colors.Chart5 != "" {
		b.WriteString(fmt.Sprintf("  --chart-5: %s;\n", colors.Chart5))
	}
	if colors.Chart6 != "" {
		b.WriteString(fmt.Sprintf("  --chart-6: %s;\n", colors.Chart6))
	}
	if colors.Chart7 != "" {
		b.WriteString(fmt.Sprintf("  --chart-7: %s;\n", colors.Chart7))
	}
	if colors.Chart8 != "" {
		b.WriteString(fmt.Sprintf("  --chart-8: %s;\n", colors.Chart8))
	}
	if colors.Chart9 != "" {
		b.WriteString(fmt.Sprintf("  --chart-9: %s;\n", colors.Chart9))
	}
	if colors.Chart10 != "" {
		b.WriteString(fmt.Sprintf("  --chart-10: %s;\n", colors.Chart10))
	}
	if colors.Chart11 != "" {
		b.WriteString(fmt.Sprintf("  --chart-11: %s;\n", colors.Chart11))
	}
	if colors.Chart12 != "" {
		b.WriteString(fmt.Sprintf("  --chart-12: %s;\n", colors.Chart12))
	}

	// Sidebar colors for independent theming
	if colors.Sidebar != "" {
		b.WriteString(fmt.Sprintf("  --sidebar: %s;\n", colors.Sidebar))
		b.WriteString(fmt.Sprintf("  --sidebar-foreground: %s;\n", colors.SidebarForeground))
		b.WriteString(fmt.Sprintf("  --sidebar-primary: %s;\n", colors.SidebarPrimary))
		b.WriteString(fmt.Sprintf("  --sidebar-primary-foreground: %s;\n", colors.SidebarPrimaryForeground))
		b.WriteString(fmt.Sprintf("  --sidebar-accent: %s;\n", colors.SidebarAccent))
		b.WriteString(fmt.Sprintf("  --sidebar-accent-foreground: %s;\n", colors.SidebarAccentForeground))
		b.WriteString(fmt.Sprintf("  --sidebar-border: %s;\n", colors.SidebarBorder))
		b.WriteString(fmt.Sprintf("  --sidebar-ring: %s;\n", colors.SidebarRing))
	}

	return b.String()
}

// generateRadiusVars generates CSS custom properties for radius tokens.
func generateRadiusVars(radius RadiusTokens) string {
	var b strings.Builder

	if radius.SM != "" {
		b.WriteString(fmt.Sprintf("  --radius-sm: %s;\n", radius.SM))
	}
	if radius.MD != "" {
		b.WriteString(fmt.Sprintf("  --radius-md: %s;\n", radius.MD))
	}
	if radius.LG != "" {
		b.WriteString(fmt.Sprintf("  --radius-lg: %s;\n", radius.LG))
		// Default radius for components (using LG as per modern shadcn/ui)
		b.WriteString(fmt.Sprintf("  --radius: %s;\n", radius.LG))
	}
	if radius.XL != "" {
		b.WriteString(fmt.Sprintf("  --radius-xl: %s;\n", radius.XL))
	}
	if radius.Full != "" {
		b.WriteString(fmt.Sprintf("  --radius-full: %s;\n", radius.Full))
	}

	return b.String()
}

// generateSpacingVars generates CSS custom properties for spacing tokens.
func generateSpacingVars(spacing SpacingTokens) string {
	var b strings.Builder

	if spacing.XS != "" {
		b.WriteString(fmt.Sprintf("  --spacing-xs: %s;\n", spacing.XS))
	}
	if spacing.SM != "" {
		b.WriteString(fmt.Sprintf("  --spacing-sm: %s;\n", spacing.SM))
	}
	if spacing.MD != "" {
		b.WriteString(fmt.Sprintf("  --spacing-md: %s;\n", spacing.MD))
		// Default spacing
		b.WriteString(fmt.Sprintf("  --spacing: %s;\n", spacing.MD))
	}
	if spacing.LG != "" {
		b.WriteString(fmt.Sprintf("  --spacing-lg: %s;\n", spacing.LG))
	}
	if spacing.XL != "" {
		b.WriteString(fmt.Sprintf("  --spacing-xl: %s;\n", spacing.XL))
	}
	if spacing.XXL != "" {
		b.WriteString(fmt.Sprintf("  --spacing-xxl: %s;\n", spacing.XXL))
	}

	return b.String()
}

// generateFontSizeVars generates CSS custom properties for font size tokens.
func generateFontSizeVars(fontSize FontSizeTokens) string {
	var b strings.Builder

	if fontSize.XS != "" {
		b.WriteString(fmt.Sprintf("  --font-size-xs: %s;\n", fontSize.XS))
	}
	if fontSize.SM != "" {
		b.WriteString(fmt.Sprintf("  --font-size-sm: %s;\n", fontSize.SM))
	}
	if fontSize.Base != "" {
		b.WriteString(fmt.Sprintf("  --font-size-base: %s;\n", fontSize.Base))
	}
	if fontSize.LG != "" {
		b.WriteString(fmt.Sprintf("  --font-size-lg: %s;\n", fontSize.LG))
	}
	if fontSize.XL != "" {
		b.WriteString(fmt.Sprintf("  --font-size-xl: %s;\n", fontSize.XL))
	}
	if fontSize.XXL != "" {
		b.WriteString(fmt.Sprintf("  --font-size-xxl: %s;\n", fontSize.XXL))
	}
	if fontSize.XXXL != "" {
		b.WriteString(fmt.Sprintf("  --font-size-xxxl: %s;\n", fontSize.XXXL))
	}

	return b.String()
}

// generateShadowVars generates CSS custom properties for shadow tokens.
func generateShadowVars(shadow ShadowTokens) string {
	var b strings.Builder

	if shadow.SM != "" {
		b.WriteString(fmt.Sprintf("  --shadow-sm: %s;\n", shadow.SM))
	}
	if shadow.MD != "" {
		b.WriteString(fmt.Sprintf("  --shadow-md: %s;\n", shadow.MD))
		// Default shadow
		b.WriteString(fmt.Sprintf("  --shadow: %s;\n", shadow.MD))
	}
	if shadow.LG != "" {
		b.WriteString(fmt.Sprintf("  --shadow-lg: %s;\n", shadow.LG))
	}
	if shadow.XL != "" {
		b.WriteString(fmt.Sprintf("  --shadow-xl: %s;\n", shadow.XL))
	}

	return b.String()
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

