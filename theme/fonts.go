package theme

import (
	"fmt"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Font defines a font family with its properties.
type Font struct {
	Family   string   // Font family name (e.g., "Inter", "Roboto")
	Weights  []int    // Font weights to load (e.g., 400, 500, 600, 700)
	Styles   []string // Font styles (e.g., "normal", "italic")
	URL      string   // Custom font URL (if not using Google Fonts)
	Display  string   // font-display strategy (swap, block, fallback, optional)
	Fallback []string // Fallback font stack
}

// FontConfig defines typography configuration for the theme.
type FontConfig struct {
	Sans         Font   // Sans-serif font
	Serif        Font   // Serif font
	Mono         Font   // Monospace font
	Body         string // Body font family
	Heading      string // Heading font family
	Code         string // Code font family
	BaseFontSize string // Base font size (e.g., "16px")
}

// DefaultFontConfig returns a default font configuration using system fonts.
func DefaultFontConfig() FontConfig {
	return FontConfig{
		Sans: Font{
			Family:   "system-ui",
			Fallback: []string{"-apple-system", "BlinkMacSystemFont", "Segoe UI", "Roboto", "sans-serif"},
			Display:  "swap",
		},
		Serif: Font{
			Family:   "Georgia",
			Fallback: []string{"Cambria", "Times New Roman", "Times", "serif"},
			Display:  "swap",
		},
		Mono: Font{
			Family:   "ui-monospace",
			Fallback: []string{"SFMono-Regular", "Menlo", "Monaco", "Consolas", "monospace"},
			Display:  "swap",
		},
		Body:         "sans",
		Heading:      "sans",
		Code:         "mono",
		BaseFontSize: "16px",
	}
}

// InterFontConfig returns a font configuration using Inter from Google Fonts.
// Inter is a popular choice for modern web applications.
func InterFontConfig() FontConfig {
	return FontConfig{
		Sans: Font{
			Family:   "Inter",
			Weights:  []int{400, 500, 600, 700},
			Fallback: []string{"system-ui", "sans-serif"},
			Display:  "swap",
		},
		Mono: Font{
			Family:   "JetBrains Mono",
			Weights:  []int{400, 500, 600},
			Fallback: []string{"ui-monospace", "monospace"},
			Display:  "swap",
		},
		Body:         "sans",
		Heading:      "sans",
		Code:         "mono",
		BaseFontSize: "16px",
	}
}

// GenerateGoogleFontsURL generates a Google Fonts URL for the given fonts.
//
// Example:
//
//	fonts := []Font{
//	    {Family: "Inter", Weights: []int{400, 600, 700}},
//	    {Family: "JetBrains Mono", Weights: []int{400, 500}},
//	}
//	url := GenerateGoogleFontsURL(fonts)
//	// Returns: https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&family=JetBrains+Mono:wght@400;500&display=swap
func GenerateGoogleFontsURL(fonts []Font) string {
	if len(fonts) == 0 {
		return ""
	}

	var families []string

	display := "swap" // Default display strategy

	for _, font := range fonts {
		if font.URL != "" {
			// Skip fonts with custom URLs
			continue
		}

		family := strings.ReplaceAll(font.Family, " ", "+")

		if len(font.Weights) > 0 {
			weights := make([]string, len(font.Weights))
			for i, w := range font.Weights {
				weights[i] = strconv.Itoa(w)
			}

			family += ":wght@" + strings.Join(weights, ";")
		}

		families = append(families, family)

		if font.Display != "" {
			display = font.Display
		}
	}

	if len(families) == 0 {
		return ""
	}

	return "https://fonts.googleapis.com/css2?" +
		"family=" + strings.Join(families, "&family=") +
		"&display=" + display
}

// FontLink returns a <link> tag for loading fonts.
// It automatically uses Google Fonts for standard fonts or custom URLs.
func FontLink(fonts ...Font) g.Node {
	googleFonts := []Font{}
	customFonts := []g.Node{}

	for _, font := range fonts {
		if font.URL != "" {
			// Custom font URL
			customFonts = append(customFonts, html.Link(
				g.Attr("rel", "stylesheet"),
				g.Attr("href", font.URL),
			))
		} else {
			// Google Font
			googleFonts = append(googleFonts, font)
		}
	}

	nodes := []g.Node{}

	// Add Google Fonts link if any
	if len(googleFonts) > 0 {
		url := GenerateGoogleFontsURL(googleFonts)
		if url != "" {
			nodes = append(nodes,
				// Preconnect for faster loading
				html.Link(
					g.Attr("rel", "preconnect"),
					g.Attr("href", "https://fonts.googleapis.com"),
				),
				html.Link(
					g.Attr("rel", "preconnect"),
					g.Attr("href", "https://fonts.gstatic.com"),
					g.Attr("crossorigin", ""),
				),
				// Font stylesheet
				html.Link(
					g.Attr("rel", "stylesheet"),
					g.Attr("href", url),
				),
			)
		}
	}

	// Add custom font links
	nodes = append(nodes, customFonts...)

	return g.Group(nodes)
}

// GenerateFontFaceCSS generates @font-face CSS rules for custom fonts.
// This is useful for self-hosted fonts.
//
// Example:
//
//	font := Font{
//	    Family: "CustomFont",
//	    URL:    "/fonts/custom-font.woff2",
//	    Weights: []int{400, 700},
//	    Display: "swap",
//	}
//	css := GenerateFontFaceCSS(font)
func GenerateFontFaceCSS(font Font) string {
	if font.URL == "" {
		return ""
	}

	var b strings.Builder

	weights := font.Weights
	if len(weights) == 0 {
		weights = []int{400} // Default to normal weight
	}

	styles := font.Styles
	if len(styles) == 0 {
		styles = []string{"normal"}
	}

	for _, weight := range weights {
		for _, style := range styles {
			b.WriteString("@font-face {\n")
			b.WriteString(fmt.Sprintf("  font-family: '%s';\n", font.Family))
			b.WriteString(fmt.Sprintf("  font-weight: %d;\n", weight))
			b.WriteString(fmt.Sprintf("  font-style: %s;\n", style))
			b.WriteString(fmt.Sprintf("  src: url('%s') format('woff2');\n", font.URL))

			if font.Display != "" {
				b.WriteString(fmt.Sprintf("  font-display: %s;\n", font.Display))
			}

			b.WriteString("}\n\n")
		}
	}

	return b.String()
}

// GenerateFontCSS generates CSS for font configuration.
// This includes CSS custom properties for font families.
func GenerateFontCSS(config FontConfig) string {
	var b strings.Builder

	b.WriteString(":root {\n")

	// Font families
	sansFallback := formatFontStack(config.Sans.Family, config.Sans.Fallback)
	serifFallback := formatFontStack(config.Serif.Family, config.Serif.Fallback)
	monoFallback := formatFontStack(config.Mono.Family, config.Mono.Fallback)

	b.WriteString(fmt.Sprintf("  --font-sans: %s;\n", sansFallback))
	b.WriteString(fmt.Sprintf("  --font-serif: %s;\n", serifFallback))
	b.WriteString(fmt.Sprintf("  --font-mono: %s;\n", monoFallback))

	// Base font size
	if config.BaseFontSize != "" {
		b.WriteString(fmt.Sprintf("  --font-size-base: %s;\n", config.BaseFontSize))
	}

	b.WriteString("}\n\n")

	// Body styles
	b.WriteString("body {\n")

	switch config.Body {
	case "sans":
		b.WriteString("  font-family: var(--font-sans);\n")
	case "serif":
		b.WriteString("  font-family: var(--font-serif);\n")
	case "mono":
		b.WriteString("  font-family: var(--font-mono);\n")
	}

	if config.BaseFontSize != "" {
		b.WriteString("  font-size: var(--font-size-base);\n")
	}

	b.WriteString("}\n\n")

	// Heading styles
	if config.Heading != config.Body {
		b.WriteString("h1, h2, h3, h4, h5, h6 {\n")

		switch config.Heading {
		case "sans":
			b.WriteString("  font-family: var(--font-sans);\n")
		case "serif":
			b.WriteString("  font-family: var(--font-serif);\n")
		case "mono":
			b.WriteString("  font-family: var(--font-mono);\n")
		}

		b.WriteString("}\n\n")
	}

	// Code styles
	b.WriteString("code, pre, kbd, samp {\n")
	b.WriteString("  font-family: var(--font-mono);\n")
	b.WriteString("}\n")

	return b.String()
}

// formatFontStack formats a font family with its fallback stack.
func formatFontStack(family string, fallbacks []string) string {
	all := append([]string{family}, fallbacks...)
	quoted := make([]string, len(all))

	for i, f := range all {
		// Quote font names with spaces
		if strings.Contains(f, " ") {
			quoted[i] = fmt.Sprintf(`'%s'`, f)
		} else {
			quoted[i] = f
		}
	}

	return strings.Join(quoted, ", ")
}

// FontStyleTag returns a <style> tag with font configuration CSS.
func FontStyleTag(config FontConfig) g.Node {
	css := GenerateFontCSS(config)

	return g.El("style",
		g.Attr("data-forgeui-fonts", ""),
		g.Raw(css),
	)
}
