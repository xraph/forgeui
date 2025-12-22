package primitives

import g "github.com/maragudk/gomponents"

// Spacer creates a flexible spacer that fills available space
// Useful for pushing elements apart in flex layouts
func Spacer() g.Node {
	return Box(
		WithClass("flex-1"),
	)
}
