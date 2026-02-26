package assets

import (
	"context"
	"fmt"
	stdhtml "html"
	"io"

	"github.com/a-h/templ"
)

// Script creates a <script> element for a JavaScript file
func (m *Manager) Script(path string, opts ...ScriptOption) templ.Component {
	cfg := &scriptConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<script src="%s"`, stdhtml.EscapeString(url)); err != nil {
			return err
		}

		if cfg.defer_ {
			if _, err := io.WriteString(w, ` defer`); err != nil {
				return err
			}
		}

		if cfg.async {
			if _, err := io.WriteString(w, ` async`); err != nil {
				return err
			}
		}

		if cfg.module {
			if _, err := io.WriteString(w, ` type="module"`); err != nil {
				return err
			}
		}

		if cfg.noModule {
			if _, err := io.WriteString(w, ` nomodule`); err != nil {
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

		_, err := io.WriteString(w, `></script>`)
		return err
	})
}

// PreloadScript creates a <link rel="preload"> element for a JavaScript file
func (m *Manager) PreloadScript(path string, opts ...ScriptOption) templ.Component {
	cfg := &scriptConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	url := m.URL(path)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<link rel="preload" as="script" href="%s"`, stdhtml.EscapeString(url)); err != nil {
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

// InlineScript creates a <script> element with inline JavaScript content
func InlineScript(content string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<script>%s</script>`, content)
		return err
	})
}

// InlineScriptWithAttrs creates a <script> element with inline JavaScript and custom attributes
func InlineScriptWithAttrs(content string, attrs templ.Attributes) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<script`); err != nil {
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

		_, err := fmt.Fprintf(w, `>%s</script>`, content)
		return err
	})
}
