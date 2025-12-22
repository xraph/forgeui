package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
)

// TestAdvancedIntegration tests the full integration of groups, layouts, loaders, and metadata
func TestAdvancedIntegration(t *testing.T) {
	r := New()

	// Register layouts
	r.RegisterLayout("main", func(ctx *PageContext, content g.Node) g.Node {
		return g.El("div", g.Text("<main>"), content, g.Text("</main>"))
	})

	// Create loader
	loader := func(ctx context.Context, params Params) (any, error) {
		return "loaded:" + params["id"], nil
	}

	// Create middleware
	addHeader := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.SetHeader("X-Test", "true")
			return next(ctx)
		}
	}

	// Create route group with all features
	api := r.Group("/api", GroupLayout("main"), GroupMiddleware(addHeader))
	api.Get("/users/:id", func(ctx *PageContext) (g.Node, error) {
		data := ctx.LoaderData().(string)
		meta := ctx.GetMeta()
		return g.Text(data + ":" + meta.Title), nil
	}).Loader(loader).Title("User Page")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/api/users/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check header from middleware
	if w.Header().Get("X-Test") != "true" {
		t.Error("Expected X-Test header to be set")
	}

	// Check output (layout + loader + meta)
	output := w.Body.String()
	if !strings.Contains(output, "loaded:123:User Page") {
		t.Errorf("Expected output to contain 'loaded:123:User Page', got '%s'", output)
	}
}

// TestAdvancedNestedGroupsWithLoaders tests nested groups with loaders
func TestAdvancedNestedGroupsWithLoaders(t *testing.T) {
	r := New()

	// Parent loader
	parentLoader := func(ctx context.Context, params Params) (any, error) {
		return "parent", nil
	}

	// Child loader
	childLoader := func(ctx context.Context, params Params) (any, error) {
		return "child", nil
	}

	// Create nested groups
	parent := r.Group("/parent")
	parent.Get("/test", func(ctx *PageContext) (g.Node, error) {
		return g.Text(ctx.LoaderData().(string)), nil
	}).Loader(parentLoader)

	child := parent.Group("/child")
	child.Get("/test", func(ctx *PageContext) (g.Node, error) {
		return g.Text(ctx.LoaderData().(string)), nil
	}).Loader(childLoader)

	// Test parent route
	req1 := httptest.NewRequest(MethodGet, "/parent/test", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Body.String() != "parent" {
		t.Errorf("Expected 'parent', got '%s'", w1.Body.String())
	}

	// Test child route
	req2 := httptest.NewRequest(MethodGet, "/parent/child/test", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Body.String() != "child" {
		t.Errorf("Expected 'child', got '%s'", w2.Body.String())
	}
}

// TestAdvancedErrorHandling tests error handling with loaders and custom error pages
func TestAdvancedErrorHandling(t *testing.T) {
	r := New()

	// Set custom error page
	r.SetErrorPage(404, func(ctx *PageContext) (g.Node, error) {
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		return g.Text("Custom 404"), nil
	})

	// Create loader that returns 404
	loader := func(ctx context.Context, params Params) (any, error) {
		return nil, Error404("Not found")
	}

	// Create route with loader
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		return g.Text("should not reach"), nil
	}).Loader(loader)

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	if w.Body.String() != "Custom 404" {
		t.Errorf("Expected 'Custom 404', got '%s'", w.Body.String())
	}
}

// TestAdvancedLayoutWithMetadata tests layout rendering with full metadata
func TestAdvancedLayoutWithMetadata(t *testing.T) {
	r := New()

	// Register layout that uses metadata
	r.RegisterLayout("seo", func(ctx *PageContext, content g.Node) g.Node {
		nodes := []g.Node{g.Text("<head>")}
		if ctx.Meta != nil {
			nodes = append(nodes, g.Group(ctx.Meta.MetaTags()))
		}
		nodes = append(nodes, g.Text("</head><body>"), content, g.Text("</body>"))
		return g.El("div", nodes...)
	})

	// Create route with full metadata
	r.Get("/page", func(ctx *PageContext) (g.Node, error) {
		return g.Text("Content"), nil
	}).SetLayout("seo").Meta(RouteMeta{
		Title:       "SEO Page",
		Description: "A page with SEO",
		Keywords:    []string{"go", "web"},
		OGImage:     "https://example.com/image.jpg",
	})

	// Test the route
	req := httptest.NewRequest(MethodGet, "/page", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	output := w.Body.String()

	// Check for meta tags
	if !strings.Contains(output, "A page with SEO") {
		t.Error("Expected description in output")
	}
	if !strings.Contains(output, "go, web") {
		t.Error("Expected keywords in output")
	}
	if !strings.Contains(output, "og:image") {
		t.Error("Expected OG image tag")
	}
	if !strings.Contains(output, "Content") {
		t.Error("Expected content in output")
	}
}

// TestAdvancedMultipleMiddleware tests multiple middleware in groups and routes
func TestAdvancedMultipleMiddleware(t *testing.T) {
	r := New()

	// Create multiple middleware
	mw1 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.Set("mw1", true)
			return next(ctx)
		}
	}

	mw2 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.Set("mw2", true)
			return next(ctx)
		}
	}

	mw3 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.Set("mw3", true)
			return next(ctx)
		}
	}

	// Global middleware
	r.Use(mw1)

	// Group middleware
	group := r.Group("/group", GroupMiddleware(mw2))

	// Route middleware
	route := group.Get("/test", func(ctx *PageContext) (g.Node, error) {
		if _, ok := ctx.Get("mw1"); !ok {
			t.Error("Expected mw1 to be set")
		}
		if _, ok := ctx.Get("mw2"); !ok {
			t.Error("Expected mw2 to be set")
		}
		if _, ok := ctx.Get("mw3"); !ok {
			t.Error("Expected mw3 to be set")
		}
		return g.Text("ok"), nil
	})
	route.WithMiddleware(mw3)

	// Test the route
	req := httptest.NewRequest(MethodGet, "/group/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Body.String() != "ok" {
		t.Errorf("Expected 'ok', got '%s'", w.Body.String())
	}
}

// TestAdvancedLoaderWithLayout tests that loaders execute before layouts
func TestAdvancedLoaderWithLayout(t *testing.T) {
	r := New()

	// Register layout that uses loaded data
	r.RegisterLayout("data-layout", func(ctx *PageContext, content g.Node) g.Node {
		data := ctx.LoaderData().(string)
		return g.El("div", g.Text("["+data+"]"), content)
	})

	// Create loader
	loader := func(ctx context.Context, params Params) (any, error) {
		return "loaded", nil
	}

	// Create route
	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		return g.Text("content"), nil
	}).Loader(loader).SetLayout("data-layout")

	// Test the route
	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	output := w.Body.String()
	if !strings.Contains(output, "[loaded]") || !strings.Contains(output, "content") {
		t.Errorf("Expected output to contain '[loaded]' and 'content', got '%s'", output)
	}
}

// TestAdvancedGroupLayoutOverride tests that route layout overrides group layout
func TestAdvancedGroupLayoutOverride(t *testing.T) {
	r := New()

	// Register layouts
	r.RegisterLayout("group-layout", func(ctx *PageContext, content g.Node) g.Node {
		return g.El("div", g.Text("GROUP:"), content)
	})
	r.RegisterLayout("route-layout", func(ctx *PageContext, content g.Node) g.Node {
		return g.El("div", g.Text("ROUTE:"), content)
	})

	// Create group with layout
	group := r.Group("/group", GroupLayout("group-layout"))

	// Route without override
	group.Get("/default", func(ctx *PageContext) (g.Node, error) {
		return g.Text("content"), nil
	})

	// Route with override
	group.Get("/override", func(ctx *PageContext) (g.Node, error) {
		return g.Text("content"), nil
	}).SetLayout("route-layout")

	// Test default
	req1 := httptest.NewRequest(MethodGet, "/group/default", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if !strings.Contains(w1.Body.String(), "GROUP:") || !strings.Contains(w1.Body.String(), "content") {
		t.Errorf("Expected output to contain 'GROUP:' and 'content', got '%s'", w1.Body.String())
	}

	// Test override
	req2 := httptest.NewRequest(MethodGet, "/group/override", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if !strings.Contains(w2.Body.String(), "ROUTE:") || !strings.Contains(w2.Body.String(), "content") {
		t.Errorf("Expected output to contain 'ROUTE:' and 'content', got '%s'", w2.Body.String())
	}
}

