package theme_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui/theme"
	g "maragu.dev/gomponents"
)

func TestDefaultFontConfig(t *testing.T) {
	config := theme.DefaultFontConfig()

	if config.Sans.Family != "system-ui" {
		t.Error("DefaultFontConfig should use system-ui for sans")
	}

	if config.Mono.Family != "ui-monospace" {
		t.Error("DefaultFontConfig should use ui-monospace for mono")
	}

	if config.BaseFontSize != "16px" {
		t.Error("DefaultFontConfig should have 16px base font size")
	}
}

func TestInterFontConfig(t *testing.T) {
	config := theme.InterFontConfig()

	if config.Sans.Family != "Inter" {
		t.Error("InterFontConfig should use Inter for sans")
	}

	if len(config.Sans.Weights) == 0 {
		t.Error("InterFontConfig should specify font weights")
	}
}

func TestGenerateGoogleFontsURL(t *testing.T) {
	fonts := []theme.Font{
		{
			Family:  "Inter",
			Weights: []int{400, 600, 700},
			Display: "swap",
		},
	}

	url := theme.GenerateGoogleFontsURL(fonts)

	if url == "" {
		t.Error("GenerateGoogleFontsURL should return non-empty URL")
	}

	// Check for Google Fonts base URL
	if !strings.Contains(url, "fonts.googleapis.com") {
		t.Error("URL should use fonts.googleapis.com")
	}

	// Check for font family
	if !strings.Contains(url, "Inter") {
		t.Error("URL should contain font family name")
	}

	// Check for weights
	if !strings.Contains(url, "400") {
		t.Error("URL should contain font weights")
	}

	// Check for display strategy
	if !strings.Contains(url, "display=swap") {
		t.Error("URL should contain display strategy")
	}
}

func TestGenerateGoogleFontsURLMultipleFonts(t *testing.T) {
	fonts := []theme.Font{
		{Family: "Inter", Weights: []int{400, 600}},
		{Family: "JetBrains Mono", Weights: []int{400}},
	}

	url := theme.GenerateGoogleFontsURL(fonts)

	if !strings.Contains(url, "Inter") {
		t.Error("URL should contain Inter")
	}

	if !strings.Contains(url, "JetBrains") {
		t.Error("URL should contain JetBrains Mono")
	}
}

func TestFontLink(t *testing.T) {
	var buf bytes.Buffer

	fonts := []theme.Font{
		{Family: "Inter", Weights: []int{400, 600}},
	}

	node := theme.FontLink(fonts...)

	// FontLink returns g.Group, so wrap it in html.Div for testing
	testNode := g.El("div", node)

	if err := testNode.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for preconnect links
	if !strings.Contains(html, `rel="preconnect"`) {
		t.Error("FontLink should include preconnect links")
	}

	// Check for stylesheet link
	if !strings.Contains(html, `rel="stylesheet"`) {
		t.Error("FontLink should include stylesheet link")
	}

	// Check for Google Fonts URLs
	if !strings.Contains(html, "fonts.googleapis.com") {
		t.Error("FontLink should link to Google Fonts")
	}
}

func TestGenerateFontFaceCSS(t *testing.T) {
	font := theme.Font{
		Family:  "CustomFont",
		URL:     "/fonts/custom.woff2",
		Weights: []int{400, 700},
		Display: "swap",
	}

	css := theme.GenerateFontFaceCSS(font)

	if css == "" {
		t.Error("GenerateFontFaceCSS should return non-empty CSS")
	}

	// Check for @font-face rule
	if !strings.Contains(css, "@font-face") {
		t.Error("CSS should contain @font-face rule")
	}

	// Check for font-family
	if !strings.Contains(css, "font-family: 'CustomFont'") {
		t.Error("CSS should contain font-family declaration")
	}

	// Check for weights
	if !strings.Contains(css, "font-weight: 400") {
		t.Error("CSS should contain font-weight declarations")
	}

	// Check for display strategy
	if !strings.Contains(css, "font-display: swap") {
		t.Error("CSS should contain font-display declaration")
	}
}

func TestGenerateFontCSS(t *testing.T) {
	config := theme.DefaultFontConfig()
	css := theme.GenerateFontCSS(config)

	if css == "" {
		t.Error("GenerateFontCSS should return non-empty CSS")
	}

	// Check for CSS variables
	expectedVars := []string{
		"--font-sans:",
		"--font-serif:",
		"--font-mono:",
	}

	for _, v := range expectedVars {
		if !strings.Contains(css, v) {
			t.Errorf("CSS should contain %s variable", v)
		}
	}

	// Check for body styles
	if !strings.Contains(css, "body {") {
		t.Error("CSS should contain body styles")
	}

	// Check for code styles
	if !strings.Contains(css, "code, pre, kbd, samp {") {
		t.Error("CSS should contain code element styles")
	}
}

func TestFontStyleTag(t *testing.T) {
	var buf bytes.Buffer

	config := theme.DefaultFontConfig()

	node := theme.FontStyleTag(config)

	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for style tag
	if !strings.Contains(html, "<style") {
		t.Error("FontStyleTag should render a style element")
	}

	// Check for data attribute
	if !strings.Contains(html, `data-forgeui-fonts`) {
		t.Error("FontStyleTag should have data-forgeui-fonts attribute")
	}

	// Check for CSS variables
	if !strings.Contains(html, "--font-sans:") {
		t.Error("Style tag should contain font CSS variables")
	}
}
