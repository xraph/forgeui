package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestRouteMeta(t *testing.T) {
	r := New()

	// Create a route with metadata
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		meta := ctx.GetMeta()
		if meta == nil {
			t.Error("Expected metadata to be set")
		}
		if meta.Title != "Test Page" {
			t.Errorf("Expected title 'Test Page', got '%s'", meta.Title)
		}
		return g.Text("content"), nil
	}).Meta(RouteMeta{
		Title:       "Test Page",
		Description: "A test page",
	})

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRouteMetaTitle(t *testing.T) {
	r := New()

	// Create a route with title
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		meta := ctx.GetMeta()
		if meta.Title != "My Title" {
			t.Errorf("Expected title 'My Title', got '%s'", meta.Title)
		}
		return g.Text("content"), nil
	}).Title("My Title")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestRouteMetaDescription(t *testing.T) {
	r := New()

	// Create a route with description
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		meta := ctx.GetMeta()
		if meta.Description != "My Description" {
			t.Errorf("Expected description 'My Description', got '%s'", meta.Description)
		}
		return g.Text("content"), nil
	}).Description("My Description")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestRouteMetaKeywords(t *testing.T) {
	r := New()

	// Create a route with keywords
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		meta := ctx.GetMeta()
		if len(meta.Keywords) != 2 {
			t.Errorf("Expected 2 keywords, got %d", len(meta.Keywords))
		}
		return g.Text("content"), nil
	}).Keywords("go", "web")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestRouteMetaNoIndex(t *testing.T) {
	r := New()

	// Create a route with NoIndex
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		meta := ctx.GetMeta()
		if !meta.NoIndex {
			t.Error("Expected NoIndex to be true")
		}
		return g.Text("content"), nil
	}).NoIndex()

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestMetaTagsGeneration(t *testing.T) {
	meta := &RouteMeta{
		Title:       "Test Page",
		Description: "A test page description",
		Keywords:    []string{"go", "web", "framework"},
		OGImage:     "https://example.com/image.jpg",
		OGType:      "website",
	}

	tags := meta.MetaTags()
	if tags == nil {
		t.Error("Expected meta tags to be generated")
	}

	// Render tags and check output
	var buf strings.Builder
	for _, tag := range tags {
		_ = tag.Render(&buf)
	}
	output := buf.String()

	// Check for various meta tags
	if !strings.Contains(output, "A test page description") {
		t.Error("Expected description meta tag")
	}
	if !strings.Contains(output, "go, web, framework") {
		t.Error("Expected keywords meta tag")
	}
	if !strings.Contains(output, "og:title") {
		t.Error("Expected Open Graph title")
	}
	if !strings.Contains(output, "og:image") {
		t.Error("Expected Open Graph image")
	}
	if !strings.Contains(output, "twitter:card") {
		t.Error("Expected Twitter card")
	}
}

func TestMetaTagsWithNoIndex(t *testing.T) {
	meta := &RouteMeta{
		NoIndex: true,
	}

	tags := meta.MetaTags()

	// Render tags and check output
	var buf strings.Builder
	for _, tag := range tags {
		_ = tag.Render(&buf)
	}
	output := buf.String()

	if !strings.Contains(output, "noindex") {
		t.Error("Expected noindex robots tag")
	}
}

func TestMetaTagsWithCanonicalURL(t *testing.T) {
	meta := &RouteMeta{
		CanonicalURL: "https://example.com/page",
	}

	tags := meta.MetaTags()

	// Render tags and check output
	var buf strings.Builder
	for _, tag := range tags {
		_ = tag.Render(&buf)
	}
	output := buf.String()

	if !strings.Contains(output, "canonical") {
		t.Error("Expected canonical link tag")
	}
	if !strings.Contains(output, "https://example.com/page") {
		t.Error("Expected canonical URL in output")
	}
}

func TestMetaTagsNil(t *testing.T) {
	var meta *RouteMeta
	tags := meta.MetaTags()

	if tags != nil {
		t.Error("Expected MetaTags to return nil when meta is nil")
	}
}

func TestMetaInLayout(t *testing.T) {
	r := New()

	// Register layout that uses meta
	r.RegisterLayout("with-meta", func(ctx *PageContext, content g.Node) g.Node {
		title := "Default"
		if ctx.Meta != nil && ctx.Meta.Title != "" {
			title = ctx.Meta.Title
		}
		return g.El("div", g.Text(title+": "), content)
	})

	// Create route with meta and layout
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		return g.Text("content"), nil
	}).Title("Page Title").SetLayout("with-meta")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "Page Title: content") {
		t.Errorf("Expected layout to use meta title, got '%s'", w.Body.String())
	}
}

