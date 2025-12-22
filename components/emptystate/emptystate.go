// Package emptystate provides an EmptyState component for displaying
// "no data" or "empty results" scenarios with optional icons, text, and actions.
//
// # Basic Usage
//
//	emptystate.EmptyState(
//	    emptystate.WithTitle("No data available"),
//	    emptystate.WithDescription("Get started by adding your first item"),
//	)
//
// # With Icon and Action
//
//	emptystate.EmptyState(
//	    emptystate.WithIcon(icons.Database()),
//	    emptystate.WithTitle("No users found"),
//	    emptystate.WithDescription("Create a new user to get started"),
//	    emptystate.WithAction(
//	        button.Primary(g.Text("Add User")),
//	    ),
//	)
package emptystate

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Props defines the properties for the EmptyState component.
type Props struct {
	Icon        g.Node
	Title       string
	Description string
	Action      g.Node
	Class       string
	Attrs       []g.Node
}

// Option is a functional option for configuring the EmptyState component.
type Option func(*Props)

// WithIcon sets the icon to display at the top of the empty state.
func WithIcon(icon g.Node) Option {
	return func(p *Props) {
		p.Icon = icon
	}
}

// WithTitle sets the title text for the empty state.
func WithTitle(title string) Option {
	return func(p *Props) {
		p.Title = title
	}
}

// WithDescription sets the description text for the empty state.
func WithDescription(description string) Option {
	return func(p *Props) {
		p.Description = description
	}
}

// WithAction sets the action button or component for the empty state.
func WithAction(action g.Node) Option {
	return func(p *Props) {
		p.Action = action
	}
}

// WithClass adds additional CSS classes to the empty state container.
func WithClass(class string) Option {
	return func(p *Props) {
		p.Class = class
	}
}

// WithAttr adds custom HTML attributes to the empty state container.
func WithAttr(attrs ...g.Node) Option {
	return func(p *Props) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// EmptyState creates an empty state component for displaying "no data" scenarios.
//
// The component is centered and includes optional icon, title, description,
// and action button. It's commonly used in tables, lists, and search results
// when no data is available.
//
// Example:
//
//	emptystate.EmptyState(
//	    emptystate.WithIcon(
//	        html.Div(
//	            html.Class("w-16 h-16 mx-auto text-muted-foreground/50"),
//	            icons.FolderOpen(icons.WithSize("64")),
//	        ),
//	    ),
//	    emptystate.WithTitle("No items found"),
//	    emptystate.WithDescription("Try adjusting your search or filter criteria"),
//	    emptystate.WithAction(
//	        button.Outline(
//	            g.Text("Clear Filters"),
//	            button.WithAttr(g.Attr("onclick", "clearFilters()")),
//	        ),
//	    ),
//	)
func EmptyState(opts ...Option) g.Node {
	props := &Props{}
	for _, opt := range opts {
		opt(props)
	}

	classes := "flex flex-col items-center justify-center py-12 px-4 text-center"

	if props.Class != "" {
		classes += " " + props.Class
	}

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	var content []g.Node

	// Icon (if provided)
	if props.Icon != nil {
		content = append(content,
			html.Div(
				html.Class("mb-4 text-muted-foreground/60"),
				props.Icon,
			),
		)
	}

	// Title
	if props.Title != "" {
		content = append(content,
			html.H3(
				html.Class("text-2xl font-bold tracking-tight mb-2"),
				g.Text(props.Title),
			),
		)
	}

	// Description
	if props.Description != "" {
		content = append(content,
			html.P(
				html.Class("text-sm text-muted-foreground max-w-sm mb-4"),
				g.Text(props.Description),
			),
		)
	}

	// Action (if provided)
	if props.Action != nil {
		content = append(content,
			html.Div(
				html.Class("mt-2"),
				props.Action,
			),
		)
	}

	return html.Div(
		g.Group(attrs),
		g.Group(content),
	)
}

