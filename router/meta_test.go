package router

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestRouteMeta(t *testing.T) {
	r := New()

	// Create a route with metadata
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		meta := ctx.GetMeta()
		if meta == nil {
			t.Error("Expected metadata to be set")
			return templ.Raw("content"), nil
		}

		if meta.Title != "Test Page" {
			t.Errorf("Expected title 'Test Page', got '%s'", meta.Title)
		}

		return templ.Raw("content"), nil
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
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		meta := ctx.GetMeta()
		if meta.Title != "My Title" {
			t.Errorf("Expected title 'My Title', got '%s'", meta.Title)
		}

		return templ.Raw("content"), nil
	}).Title("My Title")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestRouteMetaDescription(t *testing.T) {
	r := New()

	// Create a route with description
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		meta := ctx.GetMeta()
		if meta.Description != "My Description" {
			t.Errorf("Expected description 'My Description', got '%s'", meta.Description)
		}

		return templ.Raw("content"), nil
	}).Description("My Description")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestRouteMetaKeywords(t *testing.T) {
	r := New()

	// Create a route with keywords
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		meta := ctx.GetMeta()
		if len(meta.Keywords) != 2 {
			t.Errorf("Expected 2 keywords, got %d", len(meta.Keywords))
		}

		return templ.Raw("content"), nil
	}).Keywords("go", "web")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestRouteMetaNoIndex(t *testing.T) {
	r := New()

	// Create a route with NoIndex
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		meta := ctx.GetMeta()
		if !meta.NoIndex {
			t.Error("Expected NoIndex to be true")
		}

		return templ.Raw("content"), nil
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

	comp := meta.MetaTags()
	if comp == nil {
		t.Error("Expected meta tags to be generated")
	}

	// Render component and check output
	var buf strings.Builder
	_ = comp.Render(context.Background(), &buf)

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

	comp := meta.MetaTags()

	// Render component and check output
	var buf strings.Builder
	_ = comp.Render(context.Background(), &buf)

	output := buf.String()

	if !strings.Contains(output, "noindex") {
		t.Error("Expected noindex robots tag")
	}
}

func TestMetaTagsWithCanonicalURL(t *testing.T) {
	meta := &RouteMeta{
		CanonicalURL: "https://example.com/page",
	}

	comp := meta.MetaTags()

	// Render component and check output
	var buf strings.Builder
	_ = comp.Render(context.Background(), &buf)

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

	comp := meta.MetaTags()

	// MetaTags on nil meta returns a nop component that renders nothing
	var buf strings.Builder
	_ = comp.Render(context.Background(), &buf)

	if buf.Len() != 0 {
		t.Errorf("Expected empty output for nil meta, got %q", buf.String())
	}
}

func TestMetaInLayout(t *testing.T) {
	r := New()

	// Register layout that uses meta
	r.RegisterLayout("with-meta", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			title := "Default"
			if ctx.Meta != nil && ctx.Meta.Title != "" {
				title = ctx.Meta.Title
			}
			io.WriteString(w, "<div>"+title+": ")
			if err := content.Render(tCtx, w); err != nil {
				return err
			}
			io.WriteString(w, "</div>")
			return nil
		})
	})

	// Create route with meta and layout
	r.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("content"), nil
	}).Title("Page Title").SetLayout("with-meta")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "Page Title: content") {
		t.Errorf("Expected layout to use meta title, got '%s'", w.Body.String())
	}
}
