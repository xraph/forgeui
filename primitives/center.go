package primitives

import g "github.com/maragudk/gomponents"

// Center creates a centered container using flexbox
// Centers both horizontally and vertically
func Center(children ...g.Node) g.Node {
	return Flex(
		FlexJustify("center"),
		FlexAlign("center"),
		FlexChildren(children...),
	)
}
