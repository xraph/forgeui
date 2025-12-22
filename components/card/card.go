// Package card provides card components following shadcn/ui patterns.
// Cards are compound components with Header, Title, Description, Content, and Footer.
package card

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

// Props defines card configuration
type Props struct {
	Class string
	Attrs []g.Node
}

// Option is a functional option for configuring cards
type Option func(*Props)

// WithClass adds custom classes
func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

// WithAttrs adds custom attributes
func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Card creates a card container
// Use with Header, Title, Description, Content, and Footer components
//
// Example:
//
//	card.Card(
//	    card.Header(
//	        card.Title("Card Title"),
//	        card.Description("Card description"),
//	    ),
//	    card.Content(
//	        g.Text("Card content here"),
//	    ),
//	    card.Footer(
//	        button.Primary(g.Text("Action")),
//	    ),
//	)
func Card(children ...g.Node) g.Node {
	return CardWithOptions(nil, children...)
}

// CardWithOptions creates a card with options
func CardWithOptions(opts []Option, children ...g.Node) g.Node {
	props := &Props{}
	for _, opt := range opts {
		opt(props)
	}

	classes := cardCVA.Classes(nil)
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// Header creates a card header section
func Header(children ...g.Node) g.Node {
	return html.Div(
		html.Class("grid auto-rows-min grid-rows-[auto_auto] items-start gap-3 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto]"),
		g.Group(children),
	)
}

// Title creates a card title
func Title(text string, opts ...Option) g.Node {
	props := &Props{}
	for _, opt := range opts {
		opt(props)
	}

	classes := "text-lg leading-none font-semibold tracking-tight"
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.H3(
		g.Group(attrs),
		g.Text(text),
	)
}

// Description creates a card description
func Description(text string, opts ...Option) g.Node {
	props := &Props{}
	for _, opt := range opts {
		opt(props)
	}

	classes := "text-sm text-muted-foreground leading-relaxed"
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.P(
		g.Group(attrs),
		g.Text(text),
	)
}

// Content creates a card content section
func Content(children ...g.Node) g.Node {
	return html.Div(
		html.Class("px-6 pt-0"),
		g.Group(children),
	)
}

// Footer creates a card footer section
func Footer(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex items-center px-6 pt-0"),
		g.Group(children),
	)
}
