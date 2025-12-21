package menu

import (
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// NavbarProps defines navbar configuration
type NavbarProps struct {
	Fixed       bool
	Sticky      bool
	Transparent bool
	Class       string
	Attrs       []g.Node
}

// NavbarOption is a functional option for configuring navbar
type NavbarOption func(*NavbarProps)

// WithFixed makes navbar fixed position
func WithFixed() NavbarOption {
	return func(p *NavbarProps) { p.Fixed = true }
}

// WithSticky makes navbar sticky position
func WithSticky() NavbarOption {
	return func(p *NavbarProps) { p.Sticky = true }
}

// WithTransparent makes navbar background transparent
func WithTransparent() NavbarOption {
	return func(p *NavbarProps) { p.Transparent = true }
}

// WithNavbarClass adds custom classes to navbar
func WithNavbarClass(class string) NavbarOption {
	return func(p *NavbarProps) { p.Class = class }
}

// WithNavbarAttrs adds custom attributes to navbar
func WithNavbarAttrs(attrs ...g.Node) NavbarOption {
	return func(p *NavbarProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// defaultNavbarProps returns default navbar property values
func defaultNavbarProps() *NavbarProps {
	return &NavbarProps{}
}

// Navbar creates a responsive navigation bar with mobile menu
//
// Example:
//
//	menu.Navbar(
//	    menu.NavbarBrand(g.Text("My App")),
//	    menu.NavbarMenu(
//	        menu.Item("/", g.Text("Home"), menu.Active()),
//	        menu.Item("/about", g.Text("About")),
//	    ),
//	    menu.NavbarActions(
//	        button.Ghost(g.Text("Sign In")),
//	    ),
//	)
func Navbar(children ...g.Node) g.Node {
	return NavbarWithOptions(nil, children...)
}

// NavbarWithOptions creates navbar with custom options
func NavbarWithOptions(opts []NavbarOption, children ...g.Node) g.Node {
	props := defaultNavbarProps()
	for _, opt := range opts {
		opt(props)
	}

	positionClass := ""
	if props.Fixed {
		positionClass = "fixed top-0 left-0 right-0"
	} else if props.Sticky {
		positionClass = "sticky top-0"
	}

	bgClass := "bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
	if props.Transparent {
		bgClass = "bg-transparent"
	}

	classes := positionClass + " z-40 w-full border-b " + bgClass
	if props.Class != "" {
		classes += " " + props.Class
	}

	attrs := []g.Node{
		html.Class(classes),
		alpine.XData(map[string]any{
			"mobileMenuOpen": false,
		}),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Nav(
		g.Group(attrs),
		html.Div(
			html.Class("container flex h-16 items-center justify-between px-4"),
			g.Group(children),
		),
	)
}

// NavbarBrand creates the brand/logo section
//
// Example:
//
//	menu.NavbarBrand(
//	    html.Img(g.Attr("src", "/logo.svg"), g.Attr("alt", "Logo")),
//	    g.Text("My App"),
//	)
func NavbarBrand(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex items-center gap-2 font-semibold"),
		g.Group(children),
	)
}

// NavbarMenu creates the navigation links section
// On mobile, this appears in the drawer
//
// Example:
//
//	menu.NavbarMenu(
//	    menu.Item("/", g.Text("Home")),
//	    menu.Item("/docs", g.Text("Docs")),
//	)
func NavbarMenu(children ...g.Node) g.Node {
	return g.Group([]g.Node{
		// Desktop menu
		html.Div(
			html.Class("hidden md:flex md:items-center md:gap-6"),
			html.Div(
				html.Class("flex items-center gap-1"),
				g.Group(children),
			),
		),

		// Mobile menu toggle button
		html.Button(
			g.Attr("type", "button"),
			alpine.XOn("click", "mobileMenuOpen = !mobileMenuOpen"),
			html.Class("md:hidden inline-flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground focus:outline-none focus:ring-2 focus:ring-ring"),
			g.Attr("aria-label", "Toggle menu"),

			// Hamburger icon
			html.Span(
				g.Attr("x-show", "!mobileMenuOpen"),
				g.El("svg",
					html.Class("h-6 w-6"),
					g.Attr("xmlns", "http://www.w3.org/2000/svg"),
					g.Attr("fill", "none"),
					g.Attr("viewBox", "0 0 24 24"),
					g.Attr("stroke", "currentColor"),
					g.El("path",
						g.Attr("stroke-linecap", "round"),
						g.Attr("stroke-linejoin", "round"),
						g.Attr("stroke-width", "2"),
						g.Attr("d", "M4 6h16M4 12h16M4 18h16"),
					),
				),
			),

			// Close icon
			html.Span(
				g.Attr("x-show", "mobileMenuOpen"),
				g.Attr("x-cloak", ""),
				g.El("svg",
					html.Class("h-6 w-6"),
					g.Attr("xmlns", "http://www.w3.org/2000/svg"),
					g.Attr("fill", "none"),
					g.Attr("viewBox", "0 0 24 24"),
					g.Attr("stroke", "currentColor"),
					g.El("path",
						g.Attr("stroke-linecap", "round"),
						g.Attr("stroke-linejoin", "round"),
						g.Attr("stroke-width", "2"),
						g.Attr("d", "M6 18L18 6M6 6l12 12"),
					),
				),
			),
		),

		// Mobile menu drawer
		html.Div(
			g.Attr("x-show", "mobileMenuOpen"),
			alpine.XOn("click", "mobileMenuOpen = false"),
			g.Group(alpine.XTransition(animation.FadeIn())),
			html.Class("fixed inset-0 z-50 bg-background/80 backdrop-blur-sm md:hidden"),
		),

		html.Div(
			g.Attr("x-show", "mobileMenuOpen"),
			g.Group(alpine.XTransition(animation.SlideInFromRight())),
			html.Class("fixed right-0 top-0 bottom-0 z-50 w-64 bg-background border-l md:hidden"),

			html.Div(
				html.Class("flex flex-col h-full"),

				// Close button
				html.Div(
					html.Class("flex items-center justify-end p-4 border-b"),
					html.Button(
						g.Attr("type", "button"),
						alpine.XOn("click", "mobileMenuOpen = false"),
						html.Class("inline-flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground"),
						g.Attr("aria-label", "Close menu"),
						g.El("svg",
							html.Class("h-6 w-6"),
							g.Attr("xmlns", "http://www.w3.org/2000/svg"),
							g.Attr("fill", "none"),
							g.Attr("viewBox", "0 0 24 24"),
							g.Attr("stroke", "currentColor"),
							g.El("path",
								g.Attr("stroke-linecap", "round"),
								g.Attr("stroke-linejoin", "round"),
								g.Attr("stroke-width", "2"),
								g.Attr("d", "M6 18L18 6M6 6l12 12"),
							),
						),
					),
				),

				// Menu items
				html.Div(
					html.Class("flex-1 overflow-y-auto p-4"),
					html.Div(
						html.Class("flex flex-col gap-1"),
						g.Group(children),
					),
				),
			),
		),
	})
}

// NavbarActions creates the right-side actions section
//
// Example:
//
//	menu.NavbarActions(
//	    button.Ghost(g.Text("Sign In")),
//	    button.Primary(g.Text("Sign Up")),
//	)
func NavbarActions(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex items-center gap-2"),
		g.Group(children),
	)
}

