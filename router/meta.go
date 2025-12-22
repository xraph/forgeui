package router

import (
	"strings"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// RouteMeta contains SEO and page metadata
type RouteMeta struct {
	Title        string
	Description  string
	Keywords     []string
	OGImage      string
	OGType       string
	CanonicalURL string
	NoIndex      bool
}

// Meta sets the metadata for a route
func (r *Route) Meta(meta RouteMeta) *Route {
	r.Metadata = &meta
	return r
}

// Title sets just the title metadata
func (r *Route) Title(title string) *Route {
	if r.Metadata == nil {
		r.Metadata = &RouteMeta{}
	}
	r.Metadata.Title = title
	return r
}

// Description sets just the description metadata
func (r *Route) Description(desc string) *Route {
	if r.Metadata == nil {
		r.Metadata = &RouteMeta{}
	}
	r.Metadata.Description = desc
	return r
}

// Keywords sets the keywords metadata
func (r *Route) Keywords(keywords ...string) *Route {
	if r.Metadata == nil {
		r.Metadata = &RouteMeta{}
	}
	r.Metadata.Keywords = keywords
	return r
}

// OGImage sets the Open Graph image
func (r *Route) OGImage(url string) *Route {
	if r.Metadata == nil {
		r.Metadata = &RouteMeta{}
	}
	r.Metadata.OGImage = url
	return r
}

// NoIndex marks the page as no-index for search engines
func (r *Route) NoIndex() *Route {
	if r.Metadata == nil {
		r.Metadata = &RouteMeta{}
	}
	r.Metadata.NoIndex = true
	return r
}

// MetaTags generates HTML meta tags from the metadata
func (m *RouteMeta) MetaTags() []g.Node {
	if m == nil {
		return nil
	}

	tags := make([]g.Node, 0)

	// Description
	if m.Description != "" {
		tags = append(tags,
			html.Meta(
				g.Attr("name", "description"),
				g.Attr("content", m.Description),
			),
		)
	}

	// Keywords
	if len(m.Keywords) > 0 {
		tags = append(tags,
			html.Meta(
				g.Attr("name", "keywords"),
				g.Attr("content", strings.Join(m.Keywords, ", ")),
			),
		)
	}

	// Robots
	if m.NoIndex {
		tags = append(tags,
			html.Meta(
				g.Attr("name", "robots"),
				g.Attr("content", "noindex, nofollow"),
			),
		)
	}

	// Open Graph - Title
	if m.Title != "" {
		tags = append(tags,
			html.Meta(
				g.Attr("property", "og:title"),
				g.Attr("content", m.Title),
			),
		)
	}

	// Open Graph - Description
	if m.Description != "" {
		tags = append(tags,
			html.Meta(
				g.Attr("property", "og:description"),
				g.Attr("content", m.Description),
			),
		)
	}

	// Open Graph - Image
	if m.OGImage != "" {
		tags = append(tags,
			html.Meta(
				g.Attr("property", "og:image"),
				g.Attr("content", m.OGImage),
			),
		)
	}

	// Open Graph - Type
	ogType := m.OGType
	if ogType == "" {
		ogType = "website"
	}
	tags = append(tags,
		html.Meta(
			g.Attr("property", "og:type"),
			g.Attr("content", ogType),
		),
	)

	// Canonical URL
	if m.CanonicalURL != "" {
		tags = append(tags,
			html.Link(
				g.Attr("rel", "canonical"),
				g.Attr("href", m.CanonicalURL),
			),
		)
	}

	// Twitter Card
	if m.Description != "" || m.OGImage != "" {
		tags = append(tags,
			html.Meta(
				g.Attr("name", "twitter:card"),
				g.Attr("content", "summary_large_image"),
			),
		)

		if m.Title != "" {
			tags = append(tags,
				html.Meta(
					g.Attr("name", "twitter:title"),
					g.Attr("content", m.Title),
				),
			)
		}

		if m.Description != "" {
			tags = append(tags,
				html.Meta(
					g.Attr("name", "twitter:description"),
					g.Attr("content", m.Description),
				),
			)
		}

		if m.OGImage != "" {
			tags = append(tags,
				html.Meta(
					g.Attr("name", "twitter:image"),
					g.Attr("content", m.OGImage),
				),
			)
		}
	}

	return tags
}

