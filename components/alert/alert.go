// Package alert provides alert components with different variants.
package alert

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
)

var alertCVA = forgeui.NewCVA(
	"relative",
	"w-full",
	"rounded-lg",
	"border",
	"px-4",
	"py-4",
	"text-sm",
	"shadow-sm",
	"grid",
	"has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr]",
	"grid-cols-[0_1fr]",
	"has-[>svg]:gap-x-3",
	"gap-y-1",
	"items-start",
	"[&>svg]:size-4",
	"[&>svg]:translate-y-0.5",
	"[&>svg]:text-current",
).Variant("variant", map[string][]string{
	"default": {"bg-card", "text-card-foreground"},
	"destructive": {
		"text-destructive",
		"bg-card",
		"border-destructive/50",
		"[&>svg]:text-current",
		"*:data-[slot=alert-description]:text-destructive/90",
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

// Alert creates an alert container
func Alert(opts []Option, children ...g.Node) g.Node {
	props := &Props{
		Variant: forgeui.VariantDefault,
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := alertCVA.Classes(map[string]string{
		"variant": string(props.Variant),
	})

	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("role", "alert"),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// AlertTitle creates an alert title
func AlertTitle(text string) g.Node {
	return html.H5(
		html.Class("col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight"),
		g.Attr("data-slot", "alert-title"),
		g.Text(text),
	)
}

// AlertDescription creates an alert description
func AlertDescription(text string) g.Node {
	return html.Div(
		html.Class("text-muted-foreground col-start-2 grid justify-items-start gap-1 text-sm [&_p]:leading-relaxed"),
		g.Attr("data-slot", "alert-description"),
		g.Text(text),
	)
}
