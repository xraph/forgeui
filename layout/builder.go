// Package layout provides composable layout builders for ForgeUI applications.
//
// The layout package is designed for building root layouts that define the complete
// HTML structure with head and body sections. Child layouts should use standard
// gomponents to wrap content with structural elements only.
package layout

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Builder represents a layout builder for creating complete HTML documents
type Builder struct {
	headContent []g.Node
	bodyContent []g.Node
	bodyAttrs   []g.Node
	lang        string
}

// Build creates a new layout builder and returns the complete HTML document
func Build(head g.Node, body g.Node) g.Node {
	return html.Doctype(
		html.HTML(
			html.Lang("en"),
			head,
			body,
		),
	)
}

// BuildWithLang creates a new layout builder with custom language
func BuildWithLang(lang string, head g.Node, body g.Node) g.Node {
	return html.Doctype(
		html.HTML(
			html.Lang(lang),
			head,
			body,
		),
	)
}

// Head creates a head section with the given content
func Head(content ...g.Node) g.Node {
	return html.Head(
		html.Meta(html.Charset("utf-8")),
		g.Group(content),
	)
}

// Body creates a body section with the given content
func Body(content ...g.Node) g.Node {
	return html.Body(
		g.Group(content),
	)
}

// Title creates a title element
func Title(title string) g.Node {
	return html.TitleEl(g.Text(title))
}

// Meta creates a meta tag
func Meta(name, content string) g.Node {
	return html.Meta(
		html.Name(name),
		html.Content(content),
	)
}

// Description creates a meta description tag
func Description(description string) g.Node {
	return Meta("description", description)
}

// Viewport creates a viewport meta tag
func Viewport(content string) g.Node {
	return Meta("viewport", content)
}

// Class adds CSS classes to an element
func Class(classes string) g.Node {
	return html.Class(classes)
}

// Styles creates a style section in the head
func Styles(styles ...g.Node) g.Node {
	return g.Group(styles)
}

// Scripts creates a script section (typically at end of body)
func Scripts(scripts ...g.Node) g.Node {
	return g.Group(scripts)
}

// Charset creates a charset meta tag
func Charset(charset string) g.Node {
	return html.Meta(html.Charset(charset))
}

// Link creates a link element
func Link(rel, href string, attrs ...g.Node) g.Node {
	nodes := []g.Node{
		html.Rel(rel),
		html.Href(href),
	}
	nodes = append(nodes, attrs...)

	return html.Link(nodes...)
}

// StyleSheet creates a stylesheet link
func StyleSheet(href string) g.Node {
	return Link("stylesheet", href)
}

// Script creates a script element
func Script(src string, attrs ...g.Node) g.Node {
	nodes := []g.Node{
		html.Src(src),
	}
	nodes = append(nodes, attrs...)

	return html.Script(nodes...)
}

// InlineScript creates an inline script element
func InlineScript(content string) g.Node {
	return html.Script(g.Raw(content))
}

// InlineStyle creates an inline style element
func InlineStyle(content string) g.Node {
	return html.StyleEl(g.Raw(content))
}
