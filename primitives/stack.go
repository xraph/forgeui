package primitives

import g "github.com/maragudk/gomponents"

// VStack creates a vertical stack (flex column)
// This is a convenience wrapper around Flex for vertical layouts
func VStack(gap string, children ...g.Node) g.Node {
	return Flex(
		FlexDirection("col"),
		FlexGap(gap),
		FlexChildren(children...),
	)
}

// HStack creates a horizontal stack (flex row)
// This is a convenience wrapper around Flex for horizontal layouts
func HStack(gap string, children ...g.Node) g.Node {
	return Flex(
		FlexDirection("row"),
		FlexAlign("center"),
		FlexGap(gap),
		FlexChildren(children...),
	)
}
