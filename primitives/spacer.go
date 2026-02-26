package primitives

import "github.com/a-h/templ"

// Spacer creates a flexible spacer that fills available space.
// Useful for pushing elements apart in flex layouts.
func Spacer() templ.Component {
	return Box(
		WithClass("flex-1"),
	)
}
