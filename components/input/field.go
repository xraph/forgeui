package input

import (
	g "github.com/maragudk/gomponents"

	"github.com/xraph/forgeui/primitives"
)

// Field creates a complete form field with label, input, and optional description/error
func Field(labelText string, inputOpts []Option, children ...g.Node) g.Node {
	return primitives.VStack("2",
		append([]g.Node{
			primitives.Text(
				primitives.TextSize("text-sm"),
				primitives.TextWeight("font-medium"),
				primitives.TextChildren(g.Text(labelText)),
			),
			Input(inputOpts...),
		}, children...)...,
	)
}

// FormDescription creates a form field description
func FormDescription(text string) g.Node {
	return primitives.Text(
		primitives.TextSize("text-sm"),
		primitives.TextColor("text-muted-foreground"),
		primitives.TextChildren(g.Text(text)),
	)
}

// FormError creates a form field error message
func FormError(text string) g.Node {
	return primitives.Text(
		primitives.TextSize("text-sm"),
		primitives.TextWeight("font-medium"),
		primitives.TextColor("text-destructive"),
		primitives.TextChildren(g.Text(text)),
	)
}
