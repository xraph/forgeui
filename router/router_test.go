package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	g "maragu.dev/gomponents"
)

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("Expected non-nil router")
	}

	if r.routes == nil {
		t.Error("Expected routes to be initialized")
	}

	if r.namedRoutes == nil {
		t.Error("Expected namedRoutes to be initialized")
	}
}

func TestRouter_Get(t *testing.T) {
	r := New()

	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("Hello"), nil
	}

	route := r.Get("/test", handler)

	if route == nil {
		t.Fatal("Expected non-nil route")
	}

	if route.Pattern != "/test" {
		t.Errorf("Expected pattern /test, got %s", route.Pattern)
	}

	if route.Method != MethodGet {
		t.Errorf("Expected method GET, got %s", route.Method)
	}
}

func TestRouter_ServeHTTP_Success(t *testing.T) {
	r := New()

	r.Get("/hello", func(ctx *PageContext) (g.Node, error) {
		return g.Text("Hello, World!"), nil
	})

	req := httptest.NewRequest(MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %q", body)
	}
}

func TestRouter_ServeHTTP_NotFound(t *testing.T) {
	r := New()

	r.Get("/hello", func(ctx *PageContext) (g.Node, error) {
		return g.Text("Hello"), nil
	})

	req := httptest.NewRequest(MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestRouter_ServeHTTP_MethodMismatch(t *testing.T) {
	r := New()

	r.Get("/hello", func(ctx *PageContext) (g.Node, error) {
		return g.Text("Hello"), nil
	})

	req := httptest.NewRequest(MethodPost, "/hello", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestRouter_Parameters(t *testing.T) {
	r := New()

	r.Get("/users/:id", func(ctx *PageContext) (g.Node, error) {
		id := ctx.Param("id")
		return g.Text("User ID: " + id), nil
	})

	req := httptest.NewRequest(MethodGet, "/users/123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()

	expected := "User ID: 123"
	if body != expected {
		t.Errorf("Expected %q, got %q", expected, body)
	}
}

func TestRouter_QueryParams(t *testing.T) {
	r := New()

	r.Get("/search", func(ctx *PageContext) (g.Node, error) {
		q := ctx.Query("q")
		return g.Text("Search: " + q), nil
	})

	req := httptest.NewRequest(MethodGet, "/search?q=golang", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body := w.Body.String()

	expected := "Search: golang"
	if body != expected {
		t.Errorf("Expected %q, got %q", expected, body)
	}
}

func TestRouter_Middleware(t *testing.T) {
	r := New()

	// Add middleware that adds a header
	r.Use(func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.SetHeader("X-Test", "middleware")
			return next(ctx)
		}
	})

	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		return g.Text("OK"), nil
	})

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Header().Get("X-Test") != "middleware" {
		t.Error("Expected X-Test header to be set by middleware")
	}
}

func TestRouter_NamedRoutes(t *testing.T) {
	r := New()

	route := r.Get("/users/:id", func(ctx *PageContext) (g.Node, error) {
		return g.Text("User"), nil
	})

	r.Name("user", route)

	url := r.URL("user", 123)
	expected := "/users/123"

	if url != expected {
		t.Errorf("URL() = %q, want %q", url, expected)
	}
}

func TestRouter_Match(t *testing.T) {
	r := New()

	r.Match([]string{MethodGet, MethodPost}, "/api", func(ctx *PageContext) (g.Node, error) {
		return g.Text("API"), nil
	})

	// Test GET
	req := httptest.NewRequest(MethodGet, "/api", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /api should work, got status %d", w.Code)
	}

	// Test POST
	req = httptest.NewRequest(MethodPost, "/api", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /api should work, got status %d", w.Code)
	}

	// Test PUT (should not match)
	req = httptest.NewRequest(MethodPut, "/api", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("PUT /api should not match, got status %d", w.Code)
	}
}

func TestRouter_BasePath(t *testing.T) {
	r := New(WithBasePath("/api/v1"))

	r.Get("/users", func(ctx *PageContext) (g.Node, error) {
		return g.Text("Users"), nil
	})

	req := httptest.NewRequest(MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRouter_CustomNotFound(t *testing.T) {
	customHandler := func(ctx *PageContext) (g.Node, error) {
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		return g.Text("Custom 404"), nil
	}

	r := New(WithNotFound(customHandler))

	req := httptest.NewRequest(MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body := w.Body.String()
	if body != "Custom 404" {
		t.Errorf("Expected 'Custom 404', got %q", body)
	}
}
