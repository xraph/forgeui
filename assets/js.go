package assets

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Script creates a <script> element for a JavaScript file
func (m *Manager) Script(path string, opts ...ScriptOption) g.Node {
	cfg := &scriptConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	attrs := []g.Node{
		g.Attr("src", url),
	}

	if cfg.defer_ {
		attrs = append(attrs, g.Attr("defer"))
	}

	if cfg.async {
		attrs = append(attrs, g.Attr("async"))
	}

	if cfg.module {
		attrs = append(attrs, g.Attr("type", "module"))
	}

	if cfg.noModule {
		attrs = append(attrs, g.Attr("nomodule"))
	}

	if cfg.integrity != "" {
		attrs = append(attrs, g.Attr("integrity", cfg.integrity))
	}

	if cfg.crossOrigin != "" {
		attrs = append(attrs, g.Attr("crossorigin", cfg.crossOrigin))
	}

	return html.Script(attrs...)
}

// PreloadScript creates a <link rel="preload"> element for a JavaScript file
func (m *Manager) PreloadScript(path string, opts ...ScriptOption) g.Node {
	cfg := &scriptConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	attrs := []g.Node{
		g.Attr("rel", "preload"),
		g.Attr("as", "script"),
		g.Attr("href", url),
	}

	if cfg.integrity != "" {
		attrs = append(attrs, g.Attr("integrity", cfg.integrity))
	}

	if cfg.crossOrigin != "" {
		attrs = append(attrs, g.Attr("crossorigin", cfg.crossOrigin))
	}

	return html.Link(attrs...)
}

// InlineScript creates a <script> element with inline JavaScript content
func InlineScript(content string) g.Node {
	return html.Script(g.Raw(content))
}

// InlineScriptWithAttrs creates a <script> element with inline JavaScript and custom attributes
func InlineScriptWithAttrs(content string, attrs ...g.Node) g.Node {
	allAttrs := append([]g.Node{}, attrs...)
	allAttrs = append(allAttrs, g.Raw(content))
	return html.Script(allAttrs...)
}

