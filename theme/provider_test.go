package theme_test

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
	"github.com/xraph/forgeui/theme"
)

func TestProvider(t *testing.T) {
	var buf bytes.Buffer
	node := theme.Provider()
	
	// Provider returns g.Group, wrap it for testing
	testNode := g.El("div", node)
	
	if err := testNode.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for theme style tag
	if !strings.Contains(html, `data-forgeui-theme`) {
		t.Error("Provider should include theme style tag")
	}

	// Check for CSS variables
	if !strings.Contains(html, "--background:") {
		t.Error("Provider should include CSS variables")
	}

	// Check for dark mode script
	if !strings.Contains(html, `data-theme-script`) {
		t.Error("Provider should include dark mode script")
	}
}

func TestProviderWithThemes(t *testing.T) {
	var buf bytes.Buffer
	light := theme.BlueLight()
	dark := theme.BlueDark()
	
	node := theme.ProviderWithThemes(light, dark)
	
	// Provider returns g.Group, wrap it for testing
	testNode := g.El("div", node)
	
	if err := testNode.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, ":root {") {
		t.Error("Provider should include :root selector")
	}

	if !strings.Contains(html, ".dark {") {
		t.Error("Provider should include .dark selector")
	}
}

func TestStyleTag(t *testing.T) {
	var buf bytes.Buffer
	light := theme.DefaultLight()
	dark := theme.DefaultDark()
	
	node := theme.StyleTag(light, dark)
	
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "<style") {
		t.Error("StyleTag should render a style element")
	}

	if !strings.Contains(html, `data-forgeui-theme`) {
		t.Error("StyleTag should have data-forgeui-theme attribute")
	}
}

func TestHeadContent(t *testing.T) {
	var buf bytes.Buffer
	light := theme.DefaultLight()
	dark := theme.DefaultDark()
	
	node := theme.HeadContent(light, dark)
	
	// HeadContent returns g.Group, wrap it for testing
	testNode := g.El("div", node)
	
	if err := testNode.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for meta tags
	expectedMeta := []string{
		`charset="utf-8"`,
		`name="viewport"`,
		`name="theme-color"`,
		`name="color-scheme"`,
	}

	for _, meta := range expectedMeta {
		if !strings.Contains(html, meta) {
			t.Errorf("HeadContent should contain %s", meta)
		}
	}
}

func TestHTMLWrapper(t *testing.T) {
	var buf bytes.Buffer
	light := theme.DefaultLight()
	dark := theme.DefaultDark()
	
	node := theme.HTMLWrapper(light, dark)
	
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "<html") {
		t.Error("HTMLWrapper should render an html element")
	}

	if !strings.Contains(html, `lang="en"`) {
		t.Error("HTMLWrapper should have lang attribute")
	}
}

