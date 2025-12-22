package assets

import (
	"bytes"

	g "maragu.dev/gomponents"
)

// renderNode renders a gomponents node to an HTML string for testing
func renderNode(node g.Node) string {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return ""
	}
	return buf.String()
}

