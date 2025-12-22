// Package menu provides menu and navigation components following shadcn/ui patterns.
// Menu components support icons, badges, and active states.
package menu

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// MenuProps defines menu configuration
type MenuProps struct {
	Class string
	Attrs []g.Node
}

// MenuOption is a functional option for configuring menus
type MenuOption func(*MenuProps)

// WithClass adds custom classes to menu
func WithClass(class string) MenuOption {
	return func(p *MenuProps) { p.Class = class }
}

// WithAttrs adds custom attributes to menu
func WithAttrs(attrs ...g.Node) MenuOption {
	return func(p *MenuProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultMenuProps returns default menu property values
func defaultMenuProps() *MenuProps {
	return &MenuProps{}
}

// Menu creates a vertical navigation menu
//
// Example:
//
//	menu.Menu(
//	    menu.Item("/", g.Text("Home"), menu.Active()),
//	    menu.Item("/about", g.Text("About")),
//	)
func Menu(children ...g.Node) g.Node {
	return MenuWithOptions(nil, children...)
}

// MenuWithOptions creates menu with custom options
func MenuWithOptions(opts []MenuOption, children ...g.Node) g.Node {
	props := defaultMenuProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := "flex flex-col gap-1"
	if props.Class != "" {
		classes += " " + props.Class
	}

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("role", "menu"),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Nav(
		g.Group(attrs),
		g.Group(children),
	)
}

// ItemProps defines menu item configuration
type ItemProps struct {
	Href    string
	Active  bool
	Icon    g.Node
	Badge   g.Node
	Variant string // "default", "ghost"
	Class   string
	Attrs   []g.Node
}

// ItemOption is a functional option for configuring menu items
type ItemOption func(*ItemProps)

// WithHref sets the link href
func WithHref(href string) ItemOption {
	return func(p *ItemProps) { p.Href = href }
}

// Active marks the item as active
func Active() ItemOption {
	return func(p *ItemProps) { p.Active = true }
}

// WithIcon adds an icon to the item
func WithIcon(icon g.Node) ItemOption {
	return func(p *ItemProps) { p.Icon = icon }
}

// WithBadge adds a badge to the item
func WithBadge(badge g.Node) ItemOption {
	return func(p *ItemProps) { p.Badge = badge }
}

// WithVariant sets the menu item variant
func WithVariant(variant string) ItemOption {
	return func(p *ItemProps) { p.Variant = variant }
}

// WithItemClass adds custom classes to item
func WithItemClass(class string) ItemOption {
	return func(p *ItemProps) { p.Class = class }
}

// WithItemAttrs adds custom attributes to item
func WithItemAttrs(attrs ...g.Node) ItemOption {
	return func(p *ItemProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultItemProps returns default item property values
func defaultItemProps() *ItemProps {
	return &ItemProps{
		Variant: "default",
	}
}

// Item creates a menu item (link or button)
//
// Example:
//
//	menu.Item("/dashboard", g.Text("Dashboard"), menu.Active(), menu.WithIcon(icon))
func Item(href string, label g.Node, opts ...ItemOption) g.Node {
	props := defaultItemProps()

	props.Href = href
	for _, opt := range opts {
		opt(props)
	}

	baseClasses := "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors"

	var variantClasses string

	if props.Active {
		variantClasses = "bg-accent text-accent-foreground"
	} else {
		switch props.Variant {
		case "ghost":
			variantClasses = "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
		default:
			variantClasses = "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
		}
	}

	classes := baseClasses + " " + variantClasses
	if props.Class != "" {
		classes += " " + props.Class
	}

	attrs := []g.Node{
		g.Attr("role", "menuitem"),
		html.Class(classes),
	}
	attrs = append(attrs, props.Attrs...)

	content := []g.Node{}
	if props.Icon != nil {
		content = append(content, html.Span(
			html.Class("shrink-0"),
			props.Icon,
		))
	}

	content = append(content, html.Span(
		html.Class("flex-1"),
		label,
	))
	if props.Badge != nil {
		content = append(content, html.Span(
			html.Class("shrink-0"),
			props.Badge,
		))
	}

	if href != "" {
		return html.A(
			g.Attr("href", href),
			g.Group(attrs),
			g.Group(content),
		)
	}

	return html.Button(
		g.Attr("type", "button"),
		g.Group(attrs),
		g.Group(content),
	)
}

// SectionProps defines menu section configuration
type SectionProps struct {
	Label string
	Class string
	Attrs []g.Node
}

// SectionOption is a functional option for configuring menu sections
type SectionOption func(*SectionProps)

// WithSectionClass adds custom classes to section
func WithSectionClass(class string) SectionOption {
	return func(p *SectionProps) { p.Class = class }
}

// WithSectionAttrs adds custom attributes to section
func WithSectionAttrs(attrs ...g.Node) SectionOption {
	return func(p *SectionProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Section creates a grouped menu section with label
//
// Example:
//
//	menu.Section("Main",
//	    menu.Item("/", g.Text("Home")),
//	    menu.Item("/about", g.Text("About")),
//	)
func Section(label string, children ...g.Node) g.Node {
	return SectionWithOptions(label, nil, children...)
}

// SectionWithOptions creates section with custom options
func SectionWithOptions(label string, opts []SectionOption, children ...g.Node) g.Node {
	props := &SectionProps{Label: label}
	for _, opt := range opts {
		opt(props)
	}

	classes := "flex flex-col gap-1"
	if props.Class != "" {
		classes += " " + props.Class
	}

	return html.Div(
		html.Class("flex flex-col gap-2"),
		g.If(label != "", html.H4(
			html.Class("px-3 text-xs font-semibold text-muted-foreground uppercase tracking-wider"),
			g.Text(label),
		)),
		html.Div(
			html.Class(classes),
			g.Group(children),
		),
	)
}

// Separator creates a menu separator/divider
//
// Example:
//
//	menu.Separator()
func Separator() g.Node {
	return html.Hr(
		html.Class("my-2 border-t border-border"),
		g.Attr("role", "separator"),
	)
}
