package primitives

import "github.com/a-h/templ"

// Container creates a responsive container with max-width constraints.
// Commonly used for page layouts.
func Container(children ...templ.Component) templ.Component {
	return Box(
		WithClass("container mx-auto px-4"),
		WithChildren(children...),
	)
}
