package bridge

import (
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
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
}

// BridgeScripts returns script tags for the bridge client
func BridgeScripts(config ScriptConfig) g.Node {
	scripts := []g.Node{}

	// Add bridge client script
	scripts = append(scripts, html.Script(
		g.Attr("type", "text/javascript"),
		g.Raw(bridgeJS),
	))

	// Add Alpine integration if requested
	if config.IncludeAlpine {
		scripts = append(scripts, html.Script(
			g.Attr("type", "text/javascript"),
			g.Raw(alpineJS),
		))
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

	scripts = append(scripts, html.Script(
		g.Attr("type", "text/javascript"),
		g.Raw(configScript),
	))

	return g.Group(scripts)
}

// BridgeScriptsExternal returns script tags referencing external files
func BridgeScriptsExternal(config ScriptConfig) g.Node {
	scripts := []g.Node{}

	// Add bridge client script
	scripts = append(scripts, html.Script(
		g.Attr("type", "text/javascript"),
		g.Attr("src", "/static/js/forge-bridge.js"),
	))

	// Add Alpine integration if requested
	if config.IncludeAlpine {
		scripts = append(scripts, html.Script(
			g.Attr("type", "text/javascript"),
			g.Attr("src", "/static/js/alpine-bridge.js"),
		))
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

	scripts = append(scripts, html.Script(
		g.Attr("type", "text/javascript"),
		g.Raw(configScript),
	))

	return g.Group(scripts)
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
func ScriptTemplate(tmpl string, data any) (g.Node, error) {
	t, err := template.New("script").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	var buf string
	writer := &strings.Builder{}
	if err := t.Execute(writer, data); err != nil {
		return nil, err
	}
	buf = writer.String()

	return html.Script(
		g.Attr("type", "text/javascript"),
		g.Raw(buf),
	), nil
}

