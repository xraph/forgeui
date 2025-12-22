// Package badge provides badge/tag components with variant support.
package badge

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

var badgeCVA = forgeui.NewCVA(
	"inline-flex",
	"items-center",
	"justify-center",
	"rounded-md",
	"border",
	"px-2.5",
	"py-0.5",
	"text-xs",
	"font-semibold",
	"w-fit",
	"whitespace-nowrap",
	"shrink-0",
	"[&>svg]:size-3",
	"gap-1",
	"[&>svg]:pointer-events-none",
	"transition-colors",
	"focus-visible:outline-none",
	"focus-visible:ring-2",
	"focus-visible:ring-ring",
	"focus-visible:ring-offset-2",
	"overflow-hidden",
).Variant("variant", map[string][]string{
	"default": {
		"border-transparent",
		"bg-primary",
		"text-primary-foreground",
		"shadow-sm",
		"hover:bg-primary/80",
	},
	"secondary": {
		"border-transparent",
		"bg-secondary",
		"text-secondary-foreground",
		"hover:bg-secondary/80",
	},
	"destructive": {
		"border-transparent",
		"bg-destructive",
		"text-destructive-foreground",
		"shadow-sm",
		"hover:bg-destructive/80",
	},
	"outline": {
		"text-foreground",
		"hover:bg-accent",
		"hover:text-accent-foreground",
	},
}).Default("variant", "default")

type Props struct {
	Variant forgeui.Variant
	Class   string
	Attrs   []g.Node
}

type Option func(*Props)

func WithVariant(v forgeui.Variant) Option {
	return func(p *Props) { p.Variant = v }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Badge creates a badge component
func Badge(text string, opts ...Option) g.Node {
	props := &Props{
		Variant: forgeui.VariantDefault,
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := badgeCVA.Classes(map[string]string{
		"variant": string(props.Variant),
	})

	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.Span(
		g.Group(attrs),
		g.Text(text),
	)
}
