package primitives

import g "maragu.dev/gomponents"

// Spacer creates a flexible spacer that fills available space
// Useful for pushing elements apart in flex layouts
func Spacer() g.Node {
	return Box(
		WithClass("flex-1"),
	)
}
