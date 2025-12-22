package list

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ItemVariant represents the visual style of the list item.
type ItemVariant string

const (
	ItemSimple   ItemVariant = "simple"   // Simple text item
	ItemDetailed ItemVariant = "detailed" // Item with subtitle
	ItemCard     ItemVariant = "card"     // Card-style item with border
)

// ItemProps defines the properties for the ListItem component.
type ItemProps struct {
	Variant   ItemVariant
	Class     string
	Clickable bool
	OnClick   string
	Attrs     []g.Node
}

// ItemOption is a functional option for configuring the ListItem component.
type ItemOption func(*ItemProps)

// WithItemVariant sets the list item variant.
func WithItemVariant(variant ItemVariant) ItemOption {
	return func(p *ItemProps) {
		p.Variant = variant
	}
}

// Simple creates a simple list item (default).
func Simple() ItemOption {
	return func(p *ItemProps) {
		p.Variant = ItemSimple
	}
}

// Detailed creates a list item with space for subtitle/description.
func Detailed() ItemOption {
	return func(p *ItemProps) {
		p.Variant = ItemDetailed
	}
}

// CardStyle creates a card-style list item with border and padding.
func CardStyle() ItemOption {
	return func(p *ItemProps) {
		p.Variant = ItemCard
	}
}

// Clickable makes the list item clickable with hover effects.
func Clickable() ItemOption {
	return func(p *ItemProps) {
		p.Clickable = true
	}
}

// WithItemOnClick adds a click handler to the list item.
func WithItemOnClick(handler string) ItemOption {
	return func(p *ItemProps) {
		p.OnClick = handler
		p.Clickable = true
	}
}

// WithItemClass adds additional CSS classes to the list item.
func WithItemClass(class string) ItemOption {
	return func(p *ItemProps) {
		p.Class = class
	}
}

// WithItemAttr adds custom HTML attributes to the list item.
func WithItemAttr(attrs ...g.Node) ItemOption {
	return func(p *ItemProps) {
		p.Attrs = append(p.Attrs, attrs...)
	}
}

// ListItem creates a list item component.
//
// Example:
//
//	// Simple item
//	list.ListItem()(g.Text("Simple item"))
//
//	// Item with icon
//	list.ListItem()(
//	    icons.Check(),
//	    g.Text("Completed task"),
//	)
//
//	// Detailed item
//	list.ListItem(list.Detailed())(
//	    html.Div(
//	        html.Class("font-medium"),
//	        g.Text("Item title"),
//	    ),
//	    html.Div(
//	        html.Class("text-sm text-muted-foreground"),
//	        g.Text("Item description"),
//	    ),
//	)
//
//	// Card-style clickable item
//	list.ListItem(
//	    list.CardStyle(),
//	    list.WithItemOnClick("handleClick()"),
//	)(
//	    icons.User(),
//	    g.Text("User Name"),
//	    badge.Badge(g.Text("Admin")),
//	)
func ListItem(opts ...ItemOption) func(...g.Node) g.Node {
	return func(children ...g.Node) g.Node {
		props := &ItemProps{
			Variant: ItemSimple,
		}
		for _, opt := range opts {
			opt(props)
		}

		classes := ""

		switch props.Variant {
		case ItemSimple:
			classes = "flex items-center gap-2 py-1"
		case ItemDetailed:
			classes = "flex flex-col gap-1 py-2"
		case ItemCard:
			classes = "flex items-center gap-3 p-4 rounded-md border border-border bg-card"
		}

		if props.Clickable {
			classes += " cursor-pointer hover:bg-accent/50 transition-colors"
		}

		if props.Class != "" {
			classes += " " + props.Class
		}

		attrs := []g.Node{html.Class(classes)}

		if props.OnClick != "" {
			attrs = append(attrs, g.Attr("onclick", props.OnClick))
		}

		attrs = append(attrs, props.Attrs...)

		return html.Li(
			g.Group(attrs),
			g.Group(children),
		)
	}
}
