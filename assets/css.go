package assets

import (
	"context"
	"fmt"
	stdhtml "html"
	"io"

	"github.com/a-h/templ"
)

// StyleSheet creates a <link> element for a CSS file
func (m *Manager) StyleSheet(path string, opts ...StyleOption) templ.Component {
	fmt.Println("StyleSheet", path)
	cfg := &styleConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<link rel="stylesheet" href="%s"`, stdhtml.EscapeString(url)); err != nil {
			return err
		}

		if cfg.media != "" {
			if _, err := fmt.Fprintf(w, ` media="%s"`, stdhtml.EscapeString(cfg.media)); err != nil {
				return err
			}
		}

		if cfg.integrity != "" {
			if _, err := fmt.Fprintf(w, ` integrity="%s"`, stdhtml.EscapeString(cfg.integrity)); err != nil {
				return err
			}
		}

		if cfg.crossOrigin != "" {
			if _, err := fmt.Fprintf(w, ` crossorigin="%s"`, stdhtml.EscapeString(cfg.crossOrigin)); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, `>`)
		return err
	})
}

// PreloadStyleSheet creates a <link rel="preload"> element for a CSS file
func (m *Manager) PreloadStyleSheet(path string, opts ...StyleOption) templ.Component {
	cfg := &styleConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<link rel="preload" as="style" href="%s"`, stdhtml.EscapeString(url)); err != nil {
			return err
		}

		if cfg.integrity != "" {
			if _, err := fmt.Fprintf(w, ` integrity="%s"`, stdhtml.EscapeString(cfg.integrity)); err != nil {
				return err
			}
		}

		if cfg.crossOrigin != "" {
			if _, err := fmt.Fprintf(w, ` crossorigin="%s"`, stdhtml.EscapeString(cfg.crossOrigin)); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, `>`)
		return err
	})
}

// InlineCSS creates a <style> element with inline CSS content
func InlineCSS(content string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<style>%s</style>`, content)
		return err
	})
}

// InlineCSSWithAttrs creates a <style> element with inline CSS and custom attributes
func InlineCSSWithAttrs(content string, attrs templ.Attributes) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<style`); err != nil {
			return err
		}

		for k, v := range attrs {
			if s, ok := v.(string); ok {
				if _, err := fmt.Fprintf(w, ` %s="%s"`, k, stdhtml.EscapeString(s)); err != nil {
					return err
				}
			} else if v == true {
				if _, err := fmt.Fprintf(w, ` %s`, k); err != nil {
					return err
				}
			}
		}

		_, err := fmt.Fprintf(w, `>%s</style>`, content)
		return err
	})
}
