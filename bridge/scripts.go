package bridge

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/a-h/templ"
)

//go:embed client/forge-bridge.js
var bridgeJS string

//go:embed client/alpine-bridge.js
var alpineJS string

// ScriptConfig configures the bridge scripts
type ScriptConfig struct {
	Endpoint      string
	CSRFToken     string
	IncludeAlpine bool
	StaticPath    string // Base path for static assets (e.g., "/static" or "/api/identity/ui/static")
}

// BridgeScripts returns script tags for the bridge client
func BridgeScripts(config ScriptConfig) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		// Add bridge client script
		if _, err := fmt.Fprintf(w, `<script type="text/javascript">%s</script>`, bridgeJS); err != nil {
			return err
		}

		// Add Alpine integration if requested
		if config.IncludeAlpine {
			if _, err := fmt.Fprintf(w, `<script type="text/javascript">%s</script>`, alpineJS); err != nil {
				return err
			}
		}

		// Add configuration script
		configScript := fmt.Sprintf(`
window.BRIDGE_ENDPOINT = %q;
`, config.Endpoint)

		if config.CSRFToken != "" {
			configScript += fmt.Sprintf(`
// Set CSRF token if bridge is initialized
if (window.ForgeBridge && window.bridge) {
  window.bridge.setCSRF(%q);
}
`, config.CSRFToken)
		}

		_, err := fmt.Fprintf(w, `<script type="text/javascript">%s</script>`, configScript)
		return err
	})
}

// BridgeScriptsExternal returns script tags referencing external files
func BridgeScriptsExternal(config ScriptConfig) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		// Use provided static path or default to "/static"
		staticPath := config.StaticPath
		if staticPath == "" {
			staticPath = "/static"
		}

		// Add bridge client script
		if _, err := fmt.Fprintf(w, `<script type="text/javascript" src="%s/js/forge-bridge.js"></script>`, staticPath); err != nil {
			return err
		}

		// Add Alpine integration if requested
		if config.IncludeAlpine {
			if _, err := fmt.Fprintf(w, `<script type="text/javascript" src="%s/js/alpine-bridge.js"></script>`, staticPath); err != nil {
				return err
			}
		}

		// Add configuration script
		configScript := fmt.Sprintf(`
window.BRIDGE_ENDPOINT = %q;
`, config.Endpoint)

		if config.CSRFToken != "" {
			configScript += fmt.Sprintf(`
if (window.ForgeBridge && window.bridge) {
  window.bridge.setCSRF(%q);
}
`, config.CSRFToken)
		}

		_, err := fmt.Fprintf(w, `<script type="text/javascript">%s</script>`, configScript)
		return err
	})
}

// GetBridgeJS returns the bridge JavaScript code
func GetBridgeJS() string {
	return bridgeJS
}

// GetAlpineJS returns the Alpine integration JavaScript code
func GetAlpineJS() string {
	return alpineJS
}

// ScriptTemplate generates a custom script with configuration
func ScriptTemplate(tmpl string, data any) (templ.Component, error) {
	t, err := template.New("script").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	writer := &strings.Builder{}
	if err := t.Execute(writer, data); err != nil {
		return nil, err
	}

	buf := writer.String()

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<script type="text/javascript">%s</script>`, buf)
		return err
	}), nil
}
