package primitives

import "github.com/a-h/templ"

// Center creates a centered container using flexbox.
// Centers both horizontally and vertically.
func Center(children ...templ.Component) templ.Component {
	return Flex(
		FlexJustify("center"),
		FlexAlign("center"),
		FlexChildren(children...),
	)
}
