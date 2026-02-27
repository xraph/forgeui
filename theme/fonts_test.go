package theme_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/xraph/forgeui/theme"
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

	if err := node.Render(context.Background(), &buf); err != nil {
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

	if err := node.Render(context.Background(), &buf); err != nil {
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

func TestGenerateFontFaceCSSVariable(t *testing.T) {
	font := theme.Font{
		Family:      "Geist",
		Variable:    true,
		WeightRange: [2]int{100, 900},
		URL:         "/fonts/geist/geist-variable.woff2",
		Format:      "woff2",
		Display:     "swap",
	}

	css := theme.GenerateFontFaceCSS(font)

	if css == "" {
		t.Fatal("GenerateFontFaceCSS should return non-empty CSS for variable font")
	}

	if !strings.Contains(css, "@font-face") {
		t.Error("CSS should contain @font-face rule")
	}

	if !strings.Contains(css, "font-family: 'Geist'") {
		t.Error("CSS should contain font-family declaration")
	}

	// Variable font should have weight range, not individual weight
	if !strings.Contains(css, "font-weight: 100 900") {
		t.Errorf("CSS should contain weight range '100 900', got:\n%s", css)
	}

	if strings.Contains(css, "font-weight: 400") {
		t.Error("Variable font CSS should NOT contain individual weight like 400")
	}

	if !strings.Contains(css, "format('woff2')") {
		t.Error("CSS should contain format declaration")
	}

	if !strings.Contains(css, "font-display: swap") {
		t.Error("CSS should contain font-display declaration")
	}
}

func TestFontPreloadLinks(t *testing.T) {
	var buf bytes.Buffer

	fonts := []theme.Font{
		{
			Family:  "Geist",
			URL:     "/fonts/geist/geist-variable.woff2",
			Format:  "woff2",
			Preload: true,
		},
		{
			Family:  "Geist Mono",
			URL:     "/fonts/geist/geist-mono-variable.woff2",
			Format:  "woff2",
			Preload: true,
		},
		{
			Family:  "NoPreload",
			URL:     "/fonts/no-preload.woff2",
			Preload: false,
		},
	}

	comp := theme.FontPreloadLinks(fonts...)
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Should have two preload links
	if strings.Count(html, `rel="preload"`) != 2 {
		t.Errorf("expected 2 preload links, got %d in:\n%s", strings.Count(html, `rel="preload"`), html)
	}

	// Check for correct href values
	if !strings.Contains(html, `href="/fonts/geist/geist-variable.woff2"`) {
		t.Error("should contain href for geist-variable.woff2")
	}

	if !strings.Contains(html, `href="/fonts/geist/geist-mono-variable.woff2"`) {
		t.Error("should contain href for geist-mono-variable.woff2")
	}

	// Check font MIME type
	if !strings.Contains(html, `type="font/woff2"`) {
		t.Error("should contain type='font/woff2'")
	}

	// Check crossorigin
	if !strings.Contains(html, `crossorigin="anonymous"`) {
		t.Error("should contain crossorigin='anonymous'")
	}

	// NoPreload font should NOT appear
	if strings.Contains(html, "no-preload") {
		t.Error("font with Preload: false should not appear in preload links")
	}
}

func TestFontPreloadLinksFromConfig(t *testing.T) {
	var buf bytes.Buffer

	config := theme.GeistFontConfig("/static/fonts")

	comp := theme.FontPreloadLinksFromConfig(config)
	if err := comp.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if strings.Count(html, `rel="preload"`) != 2 {
		t.Errorf("expected 2 preload links from GeistFontConfig, got %d", strings.Count(html, `rel="preload"`))
	}

	if !strings.Contains(html, "/static/fonts/geist/geist-variable.woff2") {
		t.Error("should contain geist-variable.woff2 path")
	}

	if !strings.Contains(html, "/static/fonts/geist/geist-mono-variable.woff2") {
		t.Error("should contain geist-mono-variable.woff2 path")
	}
}

func TestGeistFontConfig(t *testing.T) {
	config := theme.GeistFontConfig("/assets/fonts")

	if config.Sans.Family != "Geist" {
		t.Error("GeistFontConfig should use Geist for sans")
	}

	if !config.Sans.Variable {
		t.Error("Geist sans should be a variable font")
	}

	if config.Sans.WeightRange != [2]int{100, 900} {
		t.Errorf("Geist sans weight range should be [100, 900], got %v", config.Sans.WeightRange)
	}

	if config.Sans.URL != "/assets/fonts/geist/geist-variable.woff2" {
		t.Errorf("Geist sans URL should include base path, got %s", config.Sans.URL)
	}

	if !config.Sans.Preload {
		t.Error("Geist sans should have Preload: true")
	}

	if config.Mono.Family != "Geist Mono" {
		t.Error("GeistFontConfig should use Geist Mono for mono")
	}

	if !config.Mono.Variable {
		t.Error("Geist Mono should be a variable font")
	}

	if config.Mono.URL != "/assets/fonts/geist/geist-mono-variable.woff2" {
		t.Errorf("Geist mono URL should include base path, got %s", config.Mono.URL)
	}

	if config.Body != "sans" {
		t.Error("GeistFontConfig body should be 'sans'")
	}

	if config.Code != "mono" {
		t.Error("GeistFontConfig code should be 'mono'")
	}
}
