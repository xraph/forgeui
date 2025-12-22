// Package avatar provides avatar components with image and fallback support.
package avatar

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
)

type Props struct {
	Src      string
	Alt      string
	Fallback string
	Size     forgeui.Size
	Class    string
	Attrs    []g.Node
}

type Option func(*Props)

func WithSrc(src string) Option {
	return func(p *Props) { p.Src = src }
}

func WithAlt(alt string) Option {
	return func(p *Props) { p.Alt = alt }
}

func WithFallback(fallback string) Option {
	return func(p *Props) { p.Fallback = fallback }
}

func WithSize(size forgeui.Size) Option {
	return func(p *Props) { p.Size = size }
}

func WithClass(class string) Option {
	return func(p *Props) { p.Class = class }
}

func WithAttrs(attrs ...g.Node) Option {
	return func(p *Props) { p.Attrs = append(p.Attrs, attrs...) }
}

// Avatar creates an avatar component
func Avatar(opts ...Option) g.Node {
	props := &Props{
		Size: forgeui.SizeMD,
	}

	for _, opt := range opts {
		opt(props)
	}

	baseClasses := "relative flex shrink-0 overflow-hidden rounded-full"

	sizeClasses := map[forgeui.Size]string{
		forgeui.SizeSM: "size-8",
		forgeui.SizeMD: "size-10",
		forgeui.SizeLG: "size-12",
		forgeui.SizeXL: "size-16",
	}

	sizeClass := sizeClasses[props.Size]
	if sizeClass == "" {
		sizeClass = sizeClasses[forgeui.SizeMD]
	}

	classes := forgeui.CN(baseClasses, sizeClass, props.Class)

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	children := []g.Node{}

	if props.Src != "" {
		children = append(children, html.Img(
			html.Class("aspect-square h-full w-full"),
			html.Src(props.Src),
			html.Alt(props.Alt),
		))
	} else if props.Fallback != "" {
		children = append(children, html.Span(
			html.Class("flex h-full w-full items-center justify-center rounded-full bg-muted"),
			g.Text(props.Fallback),
		))
	}

	return html.Span(
		g.Group(attrs),
		g.Group(children),
	)
}
