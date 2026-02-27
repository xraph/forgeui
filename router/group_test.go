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

func TestRouteGroup(t *testing.T) {
	r := New()

	// Create a route group
	api := r.Group("/api")
	api.Get("/users", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("users"), nil
	})

	// Test the grouped route
	req := httptest.NewRequest(MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "users" {
		t.Errorf("Expected 'users', got '%s'", w.Body.String())
	}
}

func TestRouteGroupWithMiddleware(t *testing.T) {
	r := New()

	// Middleware that adds a header
	addHeader := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			ctx.SetHeader("X-Group", "true")
			return next(ctx)
		}
	}

	// Create a route group with middleware
	api := r.Group("/api", GroupMiddleware(addHeader))
	api.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("test"), nil
	})

	// Test the grouped route
	req := httptest.NewRequest(MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Group") != "true" {
		t.Error("Expected X-Group header to be set")
	}
}

func TestNestedRouteGroups(t *testing.T) {
	r := New()

	// Create nested groups
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/users", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("v1 users"), nil
	})

	// Test the nested grouped route
	req := httptest.NewRequest(MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "v1 users" {
		t.Errorf("Expected 'v1 users', got '%s'", w.Body.String())
	}
}

func TestRouteGroupWithLayout(t *testing.T) {
	r := New()

	// Register a layout
	r.RegisterLayout("api", func(ctx *PageContext, content templ.Component) templ.Component {
		return templ.ComponentFunc(func(tCtx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, "<div>API: "); err != nil {
				return err
			}
			if err := content.Render(tCtx, w); err != nil {
				return err
			}
			_, err := io.WriteString(w, "</div>")
			return err
		})
	})

	// Create a route group with layout
	api := r.Group("/api", GroupLayout("api"))
	api.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("content"), nil
	})

	// Test the grouped route
	req := httptest.NewRequest(MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "API: ") || !strings.Contains(w.Body.String(), "content") {
		t.Errorf("Expected output to contain 'API: ' and 'content', got '%s'", w.Body.String())
	}
}

func TestRouteGroupInheritance(t *testing.T) {
	r := New()

	// Middleware that adds a header
	middleware1 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			ctx.SetHeader("X-Parent", "true")
			return next(ctx)
		}
	}

	middleware2 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (templ.Component, error) {
			ctx.SetHeader("X-Child", "true")
			return next(ctx)
		}
	}

	// Create nested groups with middleware
	parent := r.Group("/parent", GroupMiddleware(middleware1))
	child := parent.Group("/child", GroupMiddleware(middleware2))
	child.Get("/test", func(ctx *PageContext) (templ.Component, error) {
		return templ.Raw("test"), nil
	})

	// Test the nested grouped route
	req := httptest.NewRequest(MethodGet, "/parent/child/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Both middleware should have been applied
	if w.Header().Get("X-Parent") != "true" {
		t.Error("Expected X-Parent header from parent group")
	}

	if w.Header().Get("X-Child") != "true" {
		t.Error("Expected X-Child header from child group")
	}
}
