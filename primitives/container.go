package primitives

import g "github.com/maragudk/gomponents"

// Container creates a responsive container with max-width constraints
// Commonly used for page layouts
func Container(children ...g.Node) g.Node {
	return Box(
		WithClass("container mx-auto px-4"),
		WithChildren(children...),
	)
}
