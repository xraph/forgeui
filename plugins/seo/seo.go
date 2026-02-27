// Package seo provides SEO meta tag and structured data management for ForgeUI.
package seo

import (
	"context"
	"encoding/json"
	"fmt"
	stdhtml "html"
	"io"
	"maps"
	"strings"

	"github.com/a-h/templ"

	"github.com/xraph/forgeui/plugin"
)

// SEO plugin.
type SEO struct {
	*plugin.PluginBase
}

// New creates a new SEO plugin.
func New() *SEO {
	return &SEO{
		PluginBase: plugin.NewPluginBase(plugin.PluginInfo{
			Name:        "seo",
			Version:     "1.0.0",
			Description: "SEO meta tags and structured data",
			Author:      "ForgeUI",
			License:     "MIT",
		}),
	}
}

// Init initializes the SEO plugin.
func (s *SEO) Init(ctx context.Context, registry *plugin.Registry) error {
	return nil
}

// Shutdown cleanly shuts down the plugin.
func (s *SEO) Shutdown(ctx context.Context) error {
	return nil
}

// MetaTags contains SEO meta tag configuration.
type MetaTags struct {
	Title       string
	Description string
	Keywords    []string
	Author      string
	Canonical   string

	OGTitle       string
	OGDescription string
	OGImage       string
	OGType        string
	OGURL         string
	OGSiteName    string

	TwitterCard        string
	TwitterSite        string
	TwitterCreator     string
	TwitterTitle       string
	TwitterDescription string
	TwitterImage       string

	Robots    string
	GoogleBot string
	BingBot   string
}

// DefaultMetaTags returns default meta tags.
func DefaultMetaTags() MetaTags {
	return MetaTags{
		OGType:      "website",
		TwitterCard: "summary_large_image",
		Robots:      "index,follow",
	}
}

// MetaTagsNode generates meta tag nodes.
func MetaTagsNode(tags MetaTags) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if tags.Title != "" {
			if _, err := fmt.Fprintf(w, `<title>%s</title>`, stdhtml.EscapeString(tags.Title)); err != nil {
				return err
			}
		}

		if tags.Description != "" {
			if _, err := fmt.Fprintf(w, `<meta name="description" content="%s">`, stdhtml.EscapeString(tags.Description)); err != nil {
				return err
			}
		}

		if len(tags.Keywords) > 0 {
			keywords := strings.Join(tags.Keywords, ", ")
			if _, err := fmt.Fprintf(w, `<meta name="keywords" content="%s">`, stdhtml.EscapeString(keywords)); err != nil {
				return err
			}
		}

		if tags.Author != "" {
			if _, err := fmt.Fprintf(w, `<meta name="author" content="%s">`, stdhtml.EscapeString(tags.Author)); err != nil {
				return err
			}
		}

		if tags.Canonical != "" {
			if _, err := fmt.Fprintf(w, `<link rel="canonical" href="%s">`, stdhtml.EscapeString(tags.Canonical)); err != nil {
				return err
			}
		}

		// Open Graph tags
		ogTags := []struct{ property, content string }{
			{"og:title", tags.OGTitle},
			{"og:description", tags.OGDescription},
			{"og:image", tags.OGImage},
			{"og:type", tags.OGType},
			{"og:url", tags.OGURL},
			{"og:site_name", tags.OGSiteName},
		}
		for _, og := range ogTags {
			if og.content != "" {
				if _, err := fmt.Fprintf(w, `<meta property="%s" content="%s">`, og.property, stdhtml.EscapeString(og.content)); err != nil {
					return err
				}
			}
		}

		// Twitter Card tags
		twitterTags := []struct{ name, content string }{
			{"twitter:card", tags.TwitterCard},
			{"twitter:site", tags.TwitterSite},
			{"twitter:creator", tags.TwitterCreator},
			{"twitter:title", tags.TwitterTitle},
			{"twitter:description", tags.TwitterDescription},
			{"twitter:image", tags.TwitterImage},
		}
		for _, tw := range twitterTags {
			if tw.content != "" {
				if _, err := fmt.Fprintf(w, `<meta name="%s" content="%s">`, tw.name, stdhtml.EscapeString(tw.content)); err != nil {
					return err
				}
			}
		}

		// Robots meta tags
		robotsTags := []struct{ name, content string }{
			{"robots", tags.Robots},
			{"googlebot", tags.GoogleBot},
			{"bingbot", tags.BingBot},
		}
		for _, rt := range robotsTags {
			if rt.content != "" {
				if _, err := fmt.Fprintf(w, `<meta name="%s" content="%s">`, rt.name, stdhtml.EscapeString(rt.content)); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// StructuredData contains JSON-LD structured data.
type StructuredData struct {
	Type string
	Data map[string]any
}

// JSONLDNode generates JSON-LD structured data node.
func JSONLDNode(data StructuredData) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		schema := map[string]any{
			"@context": "https://schema.org",
			"@type":    data.Type,
		}

		maps.Copy(schema, data.Data)

		jsonData, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			jsonData = []byte("{}")
		}

		_, werr := fmt.Fprintf(w, `<script type="application/ld+json">%s</script>`, string(jsonData))
		return werr
	})
}

// OrganizationSchema creates Organization structured data.
func OrganizationSchema(name, url, logo string) StructuredData {
	return StructuredData{
		Type: "Organization",
		Data: map[string]any{
			"name": name,
			"url":  url,
			"logo": logo,
		},
	}
}

// ArticleSchema creates Article structured data.
func ArticleSchema(headline, description, image, datePublished, author string) StructuredData {
	return StructuredData{
		Type: "Article",
		Data: map[string]any{
			"headline":      headline,
			"description":   description,
			"image":         image,
			"datePublished": datePublished,
			"author": map[string]any{
				"@type": "Person",
				"name":  author,
			},
		},
	}
}

// BreadcrumbSchema creates BreadcrumbList structured data.
func BreadcrumbSchema(items []BreadcrumbItem) StructuredData {
	itemList := make([]map[string]any, len(items))
	for i, item := range items {
		itemList[i] = map[string]any{
			"@type":    "ListItem",
			"position": i + 1,
			"name":     item.Name,
			"item":     item.URL,
		}
	}

	return StructuredData{
		Type: "BreadcrumbList",
		Data: map[string]any{
			"itemListElement": itemList,
		},
	}
}

// BreadcrumbItem represents a breadcrumb item.
type BreadcrumbItem struct {
	Name string
	URL  string
}

// SitemapLink generates a sitemap link tag.
func SitemapLink(url string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<link rel="sitemap" type="application/xml" href="%s">`, stdhtml.EscapeString(url))
		return err
	})
}

// AlternateLink generates an alternate language link tag.
func AlternateLink(lang, url string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<link rel="alternate" hreflang="%s" href="%s">`, stdhtml.EscapeString(lang), stdhtml.EscapeString(url))
		return err
	})
}

// RSSLink generates an RSS feed link tag.
func RSSLink(title, url string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<link rel="alternate" type="application/rss+xml" title="%s" href="%s">`, stdhtml.EscapeString(title), stdhtml.EscapeString(url))
		return err
	})
}

// Preconnect generates a preconnect link tag for external domains.
func Preconnect(url string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<link rel="preconnect" href="%s">`, stdhtml.EscapeString(url))
		return err
	})
}

// DNSPrefetch generates a DNS prefetch link tag.
func DNSPrefetch(url string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<link rel="dns-prefetch" href="%s">`, stdhtml.EscapeString(url))
		return err
	})
}

// GenerateSitemap generates a basic sitemap structure.
func GenerateSitemap(urls []SitemapURL) string {
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	var sitemapSb strings.Builder
	for _, url := range urls {
		fmt.Fprintf(&sitemapSb, `
  <url>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>%s</changefreq>
    <priority>%.1f</priority>
  </url>`, url.Loc, url.LastMod, url.ChangeFreq, url.Priority)
	}

	sitemap += sitemapSb.String()

	sitemap += `
</urlset>`

	return sitemap
}

// SitemapURL represents a URL in a sitemap.
type SitemapURL struct {
	Loc        string
	LastMod    string
	ChangeFreq string
	Priority   float64
}
