package router

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
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

// MetaTags generates HTML meta tags from the metadata as a templ.Component.
func (m *RouteMeta) MetaTags() templ.Component {
	if m == nil {
		return templ.ComponentFunc(func(_ context.Context, _ io.Writer) error {
			return nil
		})
	}

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		// Description
		if m.Description != "" {
			if _, err := fmt.Fprintf(w, `<meta name="description" content="%s"/>`, m.Description); err != nil {
				return err
			}
		}

		// Keywords
		if len(m.Keywords) > 0 {
			if _, err := fmt.Fprintf(w, `<meta name="keywords" content="%s"/>`, strings.Join(m.Keywords, ", ")); err != nil {
				return err
			}
		}

		// Robots
		if m.NoIndex {
			if _, err := io.WriteString(w, `<meta name="robots" content="noindex, nofollow"/>`); err != nil {
				return err
			}
		}

		// Open Graph - Title
		if m.Title != "" {
			if _, err := fmt.Fprintf(w, `<meta property="og:title" content="%s"/>`, m.Title); err != nil {
				return err
			}
		}

		// Open Graph - Description
		if m.Description != "" {
			if _, err := fmt.Fprintf(w, `<meta property="og:description" content="%s"/>`, m.Description); err != nil {
				return err
			}
		}

		// Open Graph - Image
		if m.OGImage != "" {
			if _, err := fmt.Fprintf(w, `<meta property="og:image" content="%s"/>`, m.OGImage); err != nil {
				return err
			}
		}

		// Open Graph - Type
		ogType := m.OGType
		if ogType == "" {
			ogType = "website"
		}
		if _, err := fmt.Fprintf(w, `<meta property="og:type" content="%s"/>`, ogType); err != nil {
			return err
		}

		// Canonical URL
		if m.CanonicalURL != "" {
			if _, err := fmt.Fprintf(w, `<link rel="canonical" href="%s"/>`, m.CanonicalURL); err != nil {
				return err
			}
		}

		// Twitter Card
		if m.Description != "" || m.OGImage != "" {
			if _, err := io.WriteString(w, `<meta name="twitter:card" content="summary_large_image"/>`); err != nil {
				return err
			}

			if m.Title != "" {
				if _, err := fmt.Fprintf(w, `<meta name="twitter:title" content="%s"/>`, m.Title); err != nil {
					return err
				}
			}

			if m.Description != "" {
				if _, err := fmt.Fprintf(w, `<meta name="twitter:description" content="%s"/>`, m.Description); err != nil {
					return err
				}
			}

			if m.OGImage != "" {
				if _, err := fmt.Fprintf(w, `<meta name="twitter:image" content="%s"/>`, m.OGImage); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
