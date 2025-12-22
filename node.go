package forgeui

import (
	"strings"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// Node wraps gomponents.Node with a fluent API for building HTML elements
type Node struct {
	tag      string
	attrs    []g.Node
	children []g.Node
	classes  []string
}

// El creates a new element node with the specified HTML tag
func El(tag string) *Node {
	return &Node{
		tag:      tag,
		attrs:    make([]g.Node, 0),
		children: make([]g.Node, 0),
		classes:  make([]string, 0),
	}
}

// Class adds CSS classes to the element
// Multiple calls accumulate classes
func (n *Node) Class(classes ...string) *Node {
	n.classes = append(n.classes, classes...)
	return n
}

// Attr adds an HTML attribute to the element
func (n *Node) Attr(key, value string) *Node {
	n.attrs = append(n.attrs, g.Attr(key, value))
	return n
}

// Children adds child nodes to the element
func (n *Node) Children(children ...g.Node) *Node {
	n.children = append(n.children, children...)
	return n
}

// Build renders the node to a gomponents.Node
// This is the final step that produces the actual renderable element
func (n *Node) Build() g.Node {
	attrs := make([]g.Node, 0, len(n.attrs)+1)

	// Add classes if any
	if len(n.classes) > 0 {
		attrs = append(attrs, html.Class(strings.Join(n.classes, " ")))
	}

	// Add other attributes
	attrs = append(attrs, n.attrs...)

	// Build and return the element
	return g.El(n.tag, g.Group(attrs), g.Group(n.children))
}
