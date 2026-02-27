package theme

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Font defines a font family with its properties.
type Font struct {
	Family      string   // Font family name (e.g., "Inter", "Roboto")
	Weights     []int    // Font weights to load (e.g., 400, 500, 600, 700)
	WeightRange [2]int   // Weight range for variable fonts (e.g., [100, 900])
	Variable    bool     // True for variable fonts (uses WeightRange instead of Weights)
	Styles      []string // Font styles (e.g., "normal", "italic")
	URL         string   // Custom font URL or path (for self-hosted or remote fonts)
	Format      string   // Font format (e.g., "woff2", "woff", "truetype")
	Display     string   // font-display strategy (swap, block, fallback, optional)
	Fallback    []string // Fallback font stack
	Preload     bool     // Whether to generate <link rel="preload"> for this font
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
func GenerateGoogleFontsURL(fonts []Font) string {
	if len(fonts) == 0 {
		return ""
	}

	var families []string

	display := "swap" // Default display strategy

	for _, font := range fonts {
		if font.URL != "" {
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

// FontLink returns link tags for loading fonts.
// It automatically uses Google Fonts for standard fonts or custom URLs.
func FontLink(fonts ...Font) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		googleFonts := []Font{}

		for _, font := range fonts {
			if font.URL != "" {
				// Custom font URL
				if _, err := fmt.Fprintf(w, `<link rel="stylesheet" href="%s">`, font.URL); err != nil {
					return err
				}
			} else {
				googleFonts = append(googleFonts, font)
			}
		}

		if len(googleFonts) > 0 {
			url := GenerateGoogleFontsURL(googleFonts)
			if url != "" {
				if _, err := io.WriteString(w, `<link rel="preconnect" href="https://fonts.googleapis.com">`); err != nil {
					return err
				}
				if _, err := io.WriteString(w, `<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>`); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(w, `<link rel="stylesheet" href="%s">`, url); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GenerateFontFaceCSS generates @font-face CSS rules for custom fonts.
// For variable fonts (Variable: true), it generates a single @font-face block
// with a weight range. For non-variable fonts, it generates one block per weight.
func GenerateFontFaceCSS(font Font) string {
	if font.URL == "" {
		return ""
	}

	format := font.Format
	if format == "" {
		format = "woff2"
	}

	styles := font.Styles
	if len(styles) == 0 {
		styles = []string{"normal"}
	}

	var b strings.Builder

	if font.Variable {
		// Variable font: single @font-face with weight range
		for _, style := range styles {
			b.WriteString("@font-face {\n")
			fmt.Fprintf(&b, "  font-family: '%s';\n", font.Family)
			fmt.Fprintf(&b, "  src: url('%s') format('%s');\n", font.URL, format)
			if font.WeightRange[0] > 0 && font.WeightRange[1] > 0 {
				fmt.Fprintf(&b, "  font-weight: %d %d;\n", font.WeightRange[0], font.WeightRange[1])
			}
			fmt.Fprintf(&b, "  font-style: %s;\n", style)
			if font.Display != "" {
				fmt.Fprintf(&b, "  font-display: %s;\n", font.Display)
			}
			b.WriteString("}\n\n")
		}
	} else {
		// Non-variable font: one @font-face per weight
		weights := font.Weights
		if len(weights) == 0 {
			weights = []int{400}
		}

		for _, weight := range weights {
			for _, style := range styles {
				b.WriteString("@font-face {\n")
				fmt.Fprintf(&b, "  font-family: '%s';\n", font.Family)
				fmt.Fprintf(&b, "  font-weight: %d;\n", weight)
				fmt.Fprintf(&b, "  font-style: %s;\n", style)
				fmt.Fprintf(&b, "  src: url('%s') format('%s');\n", font.URL, format)
				if font.Display != "" {
					fmt.Fprintf(&b, "  font-display: %s;\n", font.Display)
				}
				b.WriteString("}\n\n")
			}
		}
	}

	return b.String()
}

// GenerateFontCSS generates CSS for font configuration.
func GenerateFontCSS(config FontConfig) string {
	var b strings.Builder

	b.WriteString(":root {\n")

	sansFallback := formatFontStack(config.Sans.Family, config.Sans.Fallback)
	serifFallback := formatFontStack(config.Serif.Family, config.Serif.Fallback)
	monoFallback := formatFontStack(config.Mono.Family, config.Mono.Fallback)

	fmt.Fprintf(&b, "  --font-sans: %s;\n", sansFallback)
	fmt.Fprintf(&b, "  --font-serif: %s;\n", serifFallback)
	fmt.Fprintf(&b, "  --font-mono: %s;\n", monoFallback)

	if config.BaseFontSize != "" {
		fmt.Fprintf(&b, "  --font-size-base: %s;\n", config.BaseFontSize)
	}

	b.WriteString("}\n\n")

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
		if strings.Contains(f, " ") {
			quoted[i] = fmt.Sprintf(`'%s'`, f)
		} else {
			quoted[i] = f
		}
	}

	return strings.Join(quoted, ", ")
}

// fontMIME returns the MIME type for a font format string.
func fontMIME(format string) string {
	switch strings.ToLower(format) {
	case "woff2":
		return "font/woff2"
	case "woff":
		return "font/woff"
	case "truetype", "ttf":
		return "font/ttf"
	case "opentype", "otf":
		return "font/otf"
	default:
		return "font/woff2"
	}
}

// FontPreloadLinks returns <link rel="preload"> tags for fonts that have Preload: true.
// Place these in the <head> before stylesheets for optimal font loading.
func FontPreloadLinks(fonts ...Font) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		for _, f := range fonts {
			if !f.Preload || f.URL == "" {
				continue
			}
			format := f.Format
			if format == "" {
				format = "woff2"
			}
			mime := fontMIME(format)
			if _, err := fmt.Fprintf(w, `<link rel="preload" href="%s" as="font" type="%s" crossorigin="anonymous">`, f.URL, mime); err != nil {
				return err
			}
		}
		return nil
	})
}

// FontPreloadLinksFromConfig returns <link rel="preload"> tags for all fonts
// in a FontConfig that have Preload: true.
func FontPreloadLinksFromConfig(config FontConfig) templ.Component {
	var fonts []Font
	for _, f := range []Font{config.Sans, config.Serif, config.Mono} {
		if f.Preload && f.URL != "" {
			fonts = append(fonts, f)
		}
	}
	return FontPreloadLinks(fonts...)
}

// GeistFontConfig returns a font configuration using Geist variable fonts.
// fontsBasePath is the URL path to the fonts directory (e.g., "/static/fonts").
func GeistFontConfig(fontsBasePath string) FontConfig {
	return FontConfig{
		Sans: Font{
			Family:      "Geist",
			Variable:    true,
			WeightRange: [2]int{100, 900},
			URL:         fontsBasePath + "/geist/geist-variable.woff2",
			Format:      "woff2",
			Display:     "swap",
			Preload:     true,
			Fallback:    []string{"ui-sans-serif", "system-ui", "sans-serif"},
		},
		Mono: Font{
			Family:      "Geist Mono",
			Variable:    true,
			WeightRange: [2]int{100, 900},
			URL:         fontsBasePath + "/geist/geist-mono-variable.woff2",
			Format:      "woff2",
			Display:     "swap",
			Preload:     true,
			Fallback:    []string{"ui-monospace", "monospace"},
		},
		Body:         "sans",
		Heading:      "sans",
		Code:         "mono",
		BaseFontSize: "16px",
	}
}

// FontStyleTag returns a <style> tag with font configuration CSS.
func FontStyleTag(config FontConfig) templ.Component {
	css := GenerateFontCSS(config)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<style data-forgeui-fonts>%s</style>`, css)
		return err
	})
}
