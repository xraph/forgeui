package router

import (
	"testing"
)

func TestRouteMatch_Static(t *testing.T) {
	route := newRoute("/users", MethodGet, nil)

	tests := []struct {
		path        string
		shouldMatch bool
	}{
		{"/users", true},
		{"/users/", true}, // Trailing slash normalized
		{"/user", false},
		{"/users/123", false},
		{"/", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			_, ok := route.Match(tt.path)
			if ok != tt.shouldMatch {
				t.Errorf("Match(%q) = %v, want %v", tt.path, ok, tt.shouldMatch)
			}
		})
	}
}

func TestRouteMatch_Parameters(t *testing.T) {
	route := newRoute("/users/:id", MethodGet, nil)

	params, ok := route.Match("/users/123")
	if !ok {
		t.Fatal("Expected route to match /users/123")
	}

	if params["id"] != "123" {
		t.Errorf("Expected id=123, got %s", params["id"])
	}
}

func TestRouteMatch_MultipleParameters(t *testing.T) {
	route := newRoute("/users/:userId/posts/:postId", MethodGet, nil)

	params, ok := route.Match("/users/42/posts/99")
	if !ok {
		t.Fatal("Expected route to match")
	}

	if params["userId"] != "42" {
		t.Errorf("Expected userId=42, got %s", params["userId"])
	}

	if params["postId"] != "99" {
		t.Errorf("Expected postId=99, got %s", params["postId"])
	}
}

func TestRouteMatch_Wildcard(t *testing.T) {
	route := newRoute("/files/*filepath", MethodGet, nil)

	params, ok := route.Match("/files/docs/readme.md")
	if !ok {
		t.Fatal("Expected route to match")
	}

	if params["filepath"] != "docs/readme.md" {
		t.Errorf("Expected filepath=docs/readme.md, got %s", params["filepath"])
	}
}

func TestRoutePriority(t *testing.T) {
	tests := []struct {
		pattern  string
		priority int
	}{
		{"/users", 0},        // Static - highest priority
		{"/users/:id", 10},   // Parameter - medium priority
		{"/files/*path", 20}, // Wildcard - lowest priority
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			route := newRoute(tt.pattern, MethodGet, nil)
			if route.priority != tt.priority {
				t.Errorf("Expected priority %d, got %d", tt.priority, route.priority)
			}
		})
	}
}

func TestRouteURL(t *testing.T) {
	tests := []struct {
		pattern  string
		params   []any
		expected string
	}{
		{"/users", nil, "/users"},
		{"/users/:id", []any{123}, "/users/123"},
		{"/users/:userId/posts/:postId", []any{42, 99}, "/users/42/posts/99"},
		{"/files/*path", []any{"docs/readme.md"}, "/files/docs/readme.md"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			route := newRoute(tt.pattern, MethodGet, nil)

			url := route.URL(tt.params...)
			if url != tt.expected {
				t.Errorf("URL() = %q, want %q", url, tt.expected)
			}
		})
	}
}

func TestRouteURLMap(t *testing.T) {
	route := newRoute("/users/:userId/posts/:postId", MethodGet, nil)

	params := map[string]any{
		"userId": 42,
		"postId": 99,
	}

	url := route.URLMap(params)
	expected := "/users/42/posts/99"

	if url != expected {
		t.Errorf("URLMap() = %q, want %q", url, expected)
	}
}
