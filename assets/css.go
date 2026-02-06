package assets

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// StyleSheet creates a <link> element for a CSS file
func (m *Manager) StyleSheet(path string, opts ...StyleOption) g.Node {
	fmt.Println("StyleSheet", path)
	cfg := &styleConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	attrs := []g.Node{
		html.Rel("stylesheet"),
		html.Href(url),
	}

	if cfg.media != "" {
		attrs = append(attrs, g.Attr("media", cfg.media))
	}

	if cfg.integrity != "" {
		attrs = append(attrs, g.Attr("integrity", cfg.integrity))
	}

	if cfg.crossOrigin != "" {
		attrs = append(attrs, g.Attr("crossorigin", cfg.crossOrigin))
	}

	return html.Link(attrs...)
}

// PreloadStyleSheet creates a <link rel="preload"> element for a CSS file
func (m *Manager) PreloadStyleSheet(path string, opts ...StyleOption) g.Node {
	cfg := &styleConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	attrs := []g.Node{
		html.Rel("preload"),
		html.As("style"),
		html.Href(url),
	}

	if cfg.integrity != "" {
		attrs = append(attrs, g.Attr("integrity", cfg.integrity))
	}

	if cfg.crossOrigin != "" {
		attrs = append(attrs, g.Attr("crossorigin", cfg.crossOrigin))
	}

	return html.Link(attrs...)
}

// InlineCSS creates a <style> element with inline CSS content
func InlineCSS(content string) g.Node {
	return html.StyleEl(g.Raw(content))
}

// InlineCSSWithAttrs creates a <style> element with inline CSS and custom attributes
func InlineCSSWithAttrs(content string, attrs ...g.Node) g.Node {
	allAttrs := append([]g.Node{}, attrs...)
	allAttrs = append(allAttrs, g.Raw(content))

	return html.StyleEl(allAttrs...)
}
