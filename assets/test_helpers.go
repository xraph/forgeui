package assets

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

// renderComponent renders a templ.Component to an HTML string for testing
func renderComponent(comp templ.Component) string {
	var buf bytes.Buffer
	if err := comp.Render(context.Background(), &buf); err != nil {
		return ""
	}

	return buf.String()
}
