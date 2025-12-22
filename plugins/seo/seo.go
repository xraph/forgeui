// Package seo provides SEO meta tag and structured data management for ForgeUI.
//
// The SEO plugin generates Open Graph, Twitter Card, and JSON-LD structured data
// to improve search engine visibility and social media sharing.
//
// # Basic Usage
//
//	registry := plugin.NewRegistry()
//	registry.Use(seo.New())
//
// # Features
//
//   - Open Graph tags
//   - Twitter Card tags
//   - JSON-LD structured data
//   - Canonical URL management
//   - Robots meta tags
package seo

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

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
	// Basic meta tags
	Title       string
	Description string
	Keywords    []string
	Author      string
	Canonical   string

	// Open Graph
	OGTitle       string
	OGDescription string
	OGImage       string
	OGType        string
	OGURL         string
	OGSiteName    string

	// Twitter Card
	TwitterCard        string // summary, summary_large_image, app, player
	TwitterSite        string // @username
	TwitterCreator     string // @username
	TwitterTitle       string
	TwitterDescription string
	TwitterImage       string

	// Robots
	Robots    string // index,follow or noindex,nofollow
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
func MetaTagsNode(tags MetaTags) g.Node {
	nodes := []g.Node{}

	// Basic meta tags
	if tags.Title != "" {
		nodes = append(nodes, html.TitleEl(g.Text(tags.Title)))
	}

	if tags.Description != "" {
		nodes = append(nodes,
			html.Meta(html.Name("description"), html.Content(tags.Description)),
		)
	}

	if len(tags.Keywords) > 0 {
		keywords := ""

		var keywordsSb116 strings.Builder

		for i, kw := range tags.Keywords {
			if i > 0 {
				keywordsSb116.WriteString(", ")
			}

			keywordsSb116.WriteString(kw)
		}

		keywords += keywordsSb116.String()

		nodes = append(nodes,
			html.Meta(html.Name("keywords"), html.Content(keywords)),
		)
	}

	if tags.Author != "" {
		nodes = append(nodes,
			html.Meta(html.Name("author"), html.Content(tags.Author)),
		)
	}

	if tags.Canonical != "" {
		nodes = append(nodes,
			html.Link(html.Rel("canonical"), html.Href(tags.Canonical)),
		)
	}

	// Open Graph tags
	if tags.OGTitle != "" {
		nodes = append(nodes,
			html.Meta(g.Attr("property", "og:title"), html.Content(tags.OGTitle)),
		)
	}

	if tags.OGDescription != "" {
		nodes = append(nodes,
			html.Meta(g.Attr("property", "og:description"), html.Content(tags.OGDescription)),
		)
	}

	if tags.OGImage != "" {
		nodes = append(nodes,
			html.Meta(g.Attr("property", "og:image"), html.Content(tags.OGImage)),
		)
	}

	if tags.OGType != "" {
		nodes = append(nodes,
			html.Meta(g.Attr("property", "og:type"), html.Content(tags.OGType)),
		)
	}

	if tags.OGURL != "" {
		nodes = append(nodes,
			html.Meta(g.Attr("property", "og:url"), html.Content(tags.OGURL)),
		)
	}

	if tags.OGSiteName != "" {
		nodes = append(nodes,
			html.Meta(g.Attr("property", "og:site_name"), html.Content(tags.OGSiteName)),
		)
	}

	// Twitter Card tags
	if tags.TwitterCard != "" {
		nodes = append(nodes,
			html.Meta(html.Name("twitter:card"), html.Content(tags.TwitterCard)),
		)
	}

	if tags.TwitterSite != "" {
		nodes = append(nodes,
			html.Meta(html.Name("twitter:site"), html.Content(tags.TwitterSite)),
		)
	}

	if tags.TwitterCreator != "" {
		nodes = append(nodes,
			html.Meta(html.Name("twitter:creator"), html.Content(tags.TwitterCreator)),
		)
	}

	if tags.TwitterTitle != "" {
		nodes = append(nodes,
			html.Meta(html.Name("twitter:title"), html.Content(tags.TwitterTitle)),
		)
	}

	if tags.TwitterDescription != "" {
		nodes = append(nodes,
			html.Meta(html.Name("twitter:description"), html.Content(tags.TwitterDescription)),
		)
	}

	if tags.TwitterImage != "" {
		nodes = append(nodes,
			html.Meta(html.Name("twitter:image"), html.Content(tags.TwitterImage)),
		)
	}

	// Robots meta tags
	if tags.Robots != "" {
		nodes = append(nodes,
			html.Meta(html.Name("robots"), html.Content(tags.Robots)),
		)
	}

	if tags.GoogleBot != "" {
		nodes = append(nodes,
			html.Meta(html.Name("googlebot"), html.Content(tags.GoogleBot)),
		)
	}

	if tags.BingBot != "" {
		nodes = append(nodes,
			html.Meta(html.Name("bingbot"), html.Content(tags.BingBot)),
		)
	}

	return g.Group(nodes)
}

// StructuredData contains JSON-LD structured data.
type StructuredData struct {
	Type string         // Schema.org type (e.g., "Organization", "Article")
	Data map[string]any // Structured data fields
}

// JSONLDNode generates JSON-LD structured data node.
func JSONLDNode(data StructuredData) g.Node {
	schema := map[string]any{
		"@context": "https://schema.org",
		"@type":    data.Type,
	}

	maps.Copy(schema, data.Data)

	jsonData, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		// Fallback to empty object on marshal error
		jsonData = []byte("{}")
	}

	return html.Script(
		html.Type("application/ld+json"),
		g.Raw(string(jsonData)),
	)
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
func SitemapLink(url string) g.Node {
	return html.Link(
		html.Rel("sitemap"),
		html.Type("application/xml"),
		html.Href(url),
	)
}

// AlternateLink generates an alternate language link tag.
func AlternateLink(lang, url string) g.Node {
	return html.Link(
		html.Rel("alternate"),
		g.Attr("hreflang", lang),
		html.Href(url),
	)
}

// RSSLink generates an RSS feed link tag.
func RSSLink(title, url string) g.Node {
	return html.Link(
		html.Rel("alternate"),
		html.Type("application/rss+xml"),
		g.Attr("title", title),
		html.Href(url),
	)
}

// Preconnect generates a preconnect link tag for external domains.
func Preconnect(url string) g.Node {
	return html.Link(
		html.Rel("preconnect"),
		html.Href(url),
	)
}

// DNSPrefetch generates a DNS prefetch link tag.
func DNSPrefetch(url string) g.Node {
	return html.Link(
		html.Rel("dns-prefetch"),
		html.Href(url),
	)
}

// GenerateSitemap generates a basic sitemap structure (not a full implementation).
func GenerateSitemap(urls []SitemapURL) string {
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	var sitemapSb368 strings.Builder
	for _, url := range urls {
		sitemapSb368.WriteString(fmt.Sprintf(`
  <url>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>%s</changefreq>
    <priority>%.1f</priority>
  </url>`, url.Loc, url.LastMod, url.ChangeFreq, url.Priority))
	}

	sitemap += sitemapSb368.String()

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
