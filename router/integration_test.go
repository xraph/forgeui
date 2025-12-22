package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// TestIntegration_FullRequestLifecycle tests a complete HTTP request through the router
func TestIntegration_FullRequestLifecycle(t *testing.T) {
	r := New()

	// Add global middleware
	r.Use(Logger())
	r.Use(RequestID())

	// Register routes
	r.Get("/", func(ctx *PageContext) (g.Node, error) {
		return html.H1(g.Text("Home")), nil
	})

	r.Get("/users/:id", func(ctx *PageContext) (g.Node, error) {
		id := ctx.Param("id")
		return html.Div(
			html.H1(g.Text("User Profile")),
			html.P(g.Textf("User ID: %s", id)),
		), nil
	})

	r.Post("/users", func(ctx *PageContext) (g.Node, error) {
		return html.P(g.Text("User created")), nil
	})

	// Test GET /
	req := httptest.NewRequest(MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET / - Expected status 200, got %d", w.Code)
	}

	// Test GET /users/123
	req = httptest.NewRequest(MethodGet, "/users/123", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /users/123 - Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty body")
	}

	// Test POST /users
	req = httptest.NewRequest(MethodPost, "/users", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /users - Expected status 200, got %d", w.Code)
	}
}

// TestIntegration_MiddlewareChain tests middleware execution order
func TestIntegration_MiddlewareChain(t *testing.T) {
	r := New()

	order := []string{}

	m1 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			order = append(order, "m1-before")
			node, err := next(ctx)
			order = append(order, "m1-after")
			return node, err
		}
	}

	m2 := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			order = append(order, "m2-before")
			node, err := next(ctx)
			order = append(order, "m2-after")
			return node, err
		}
	}

	r.Use(m1, m2)

	r.Get("/test", func(ctx *PageContext) (g.Node, error) {
		order = append(order, "handler")
		return g.Text("OK"), nil
	})

	req := httptest.NewRequest(MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Fatalf("Expected %d calls, got %d", len(expected), len(order))
	}

	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("Step %d: expected %s, got %s", i, exp, order[i])
		}
	}
}

// TestIntegration_RouteSpecificMiddleware tests route-specific middleware
func TestIntegration_RouteSpecificMiddleware(t *testing.T) {
	r := New()

	authMiddleware := func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.Set("authenticated", true)
			return next(ctx)
		}
	}

	// Public route
	r.Get("/public", func(ctx *PageContext) (g.Node, error) {
		_, ok := ctx.Get("authenticated")
		if ok {
			t.Error("Public route should not have authentication")
		}
		return g.Text("Public"), nil
	})

	// Protected route
	route := r.Get("/protected", func(ctx *PageContext) (g.Node, error) {
		_, ok := ctx.Get("authenticated")
		if !ok {
			t.Error("Protected route should have authentication")
		}
		return g.Text("Protected"), nil
	})
	route.WithMiddleware(authMiddleware)

	// Test public route
	req := httptest.NewRequest(MethodGet, "/public", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Test protected route
	req = httptest.NewRequest(MethodGet, "/protected", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

// TestIntegration_OverlappingRoutes tests route priority
func TestIntegration_OverlappingRoutes(t *testing.T) {
	r := New()

	// Register routes in non-priority order
	r.Get("/users/:id", func(ctx *PageContext) (g.Node, error) {
		return g.Text("User by ID"), nil
	})

	r.Get("/users/new", func(ctx *PageContext) (g.Node, error) {
		return g.Text("New User Form"), nil
	})

	// Static route should match first
	req := httptest.NewRequest(MethodGet, "/users/new", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if body != "New User Form" {
		t.Errorf("Expected 'New User Form', got %q", body)
	}

	// Parameter route should match
	req = httptest.NewRequest(MethodGet, "/users/123", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body = w.Body.String()
	if body != "User by ID" {
		t.Errorf("Expected 'User by ID', got %q", body)
	}
}

// TestIntegration_ErrorHandling tests error handling through the stack
func TestIntegration_ErrorHandling(t *testing.T) {
	customError := func(ctx *PageContext, err error) g.Node {
		return html.Div(
			html.H1(g.Text("Custom Error")),
			html.P(g.Text(err.Error())),
		)
	}

	r := New(WithErrorHandler(customError))

	r.Get("/error", func(ctx *PageContext) (g.Node, error) {
		return nil, http.ErrAbortHandler
	})

	req := httptest.NewRequest(MethodGet, "/error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if body == "" {
		t.Error("Expected error page to be rendered")
	}
}

// TestIntegration_ContextValues tests context value propagation
func TestIntegration_ContextValues(t *testing.T) {
	r := New()

	// Middleware that sets a value
	r.Use(func(next PageHandler) PageHandler {
		return func(ctx *PageContext) (g.Node, error) {
			ctx.Set("user_id", 42)
			ctx.Set("username", "john")
			return next(ctx)
		}
	})

	r.Get("/profile", func(ctx *PageContext) (g.Node, error) {
		userID := ctx.GetInt("user_id")
		username := ctx.GetString("username")

		if userID != 42 {
			t.Errorf("Expected user_id=42, got %d", userID)
		}

		if username != "john" {
			t.Errorf("Expected username=john, got %s", username)
		}

		return g.Text("Profile"), nil
	})

	req := httptest.NewRequest(MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

// TestIntegration_NamedRoutesURLGeneration tests URL generation
func TestIntegration_NamedRoutesURLGeneration(t *testing.T) {
	r := New()

	route := r.Get("/users/:id/posts/:postId", func(ctx *PageContext) (g.Node, error) {
		return g.Text("Post"), nil
	})

	r.Name("user.post", route)

	// Generate URL
	url := r.URL("user.post", 123, 456)
	expected := "/users/123/posts/456"

	if url != expected {
		t.Errorf("URL() = %q, want %q", url, expected)
	}

	// Test that generated URL actually works
	req := httptest.NewRequest(MethodGet, url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Generated URL should be valid, got status %d", w.Code)
	}
}

// TestIntegration_HTTPMethods tests all HTTP method handlers
func TestIntegration_HTTPMethods(t *testing.T) {
	r := New()

	handler := func(ctx *PageContext) (g.Node, error) {
		return g.Text(ctx.Method()), nil
	}

	r.Get("/test", handler)
	r.Post("/test", handler)
	r.Put("/test", handler)
	r.Patch("/test", handler)
	r.Delete("/test", handler)
	r.Options("/test", handler)
	r.Head("/test", handler)

	methods := []string{
		MethodGet, MethodPost, MethodPut, MethodPatch,
		MethodDelete, MethodOptions, MethodHead,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
			}
		})
	}
}

