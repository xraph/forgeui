package alpine

import (
	"bytes"
	"strings"
	"testing"
)

func TestXRoute(t *testing.T) {
	tests := []struct {
		name string
		path string
		opts []RouteOption
		want string
	}{
		{
			name: "static route",
			path: "/",
			opts: nil,
			want: `x-route="/"`,
		},
		{
			name: "route with param",
			path: "/users/:id",
			opts: nil,
			want: `x-route="/users/:id"`,
		},
		{
			name: "wildcard route",
			path: "/files/*path",
			opts: nil,
			want: `x-route="/files/*path"`,
		},
		{
			name: "optional param",
			path: "/profile/:id?",
			opts: nil,
			want: `x-route="/profile/:id?"`,
		},
		{
			name: "notfound route",
			path: "notfound",
			opts: nil,
			want: `x-route="notfound"`,
		},
		{
			name: "named route",
			path: "/profile",
			opts: []RouteOption{WithRouteName("profile")},
			want: `x-route:profile="/profile"`,
		},
		{
			name: "complex path",
			path: "/movies/:title.(mp4|mov)",
			opts: nil,
			want: `x-route="/movies/:title.(mp4|mov)"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := XRoute(tt.path, tt.opts...).Render(&buf); err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := buf.String()

			if !strings.Contains(got, tt.want) {
				t.Errorf("XRoute() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestXTemplate(t *testing.T) {
	tests := []struct {
		name string
		url  string
		mods []TemplateModifier
		want string
	}{
		{
			name: "simple template",
			url:  "/views/home.html",
			mods: nil,
			want: `x-template="/views/home.html"`,
		},
		{
			name: "template with preload",
			url:  "/views/404.html",
			mods: []TemplateModifier{Preload()},
			want: `x-template.preload="/views/404.html"`,
		},
		{
			name: "template with target",
			url:  "/views/profile.html",
			mods: []TemplateModifier{TargetID("app")},
			want: `x-template.target.app="/views/profile.html"`,
		},
		{
			name: "template with interpolate",
			url:  "/api/dynamic/:name.html",
			mods: []TemplateModifier{Interpolate()},
			want: `x-template.interpolate="/api/dynamic/:name.html"`,
		},
		{
			name: "template with multiple modifiers",
			url:  "/views/user.html",
			mods: []TemplateModifier{Preload(), TargetID("main")},
			want: `x-template.preload.target.main="/views/user.html"`,
		},
		{
			name: "array of templates",
			url:  "['/header.html', '/home.html']",
			mods: nil,
			want: `x-template="[`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := XTemplate(tt.url, tt.mods...).Render(&buf); err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := buf.String()

			if !strings.Contains(got, tt.want) {
				t.Errorf("XTemplate() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestXTemplateInline(t *testing.T) {
	tests := []struct {
		name string
		mods []TemplateModifier
		want string
	}{
		{
			name: "inline template",
			mods: nil,
			want: `x-template=""`,
		},
		{
			name: "inline template with target",
			mods: []TemplateModifier{TargetID("app")},
			want: `x-template.target.app=""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := XTemplateInline(tt.mods...).Render(&buf); err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := buf.String()

			if !strings.Contains(got, tt.want) {
				t.Errorf("XTemplateInline() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestXHandler(t *testing.T) {
	tests := []struct {
		name    string
		handler string
		mods    []HandlerModifier
		want    string
	}{
		{
			name:    "single handler",
			handler: "myHandler",
			mods:    nil,
			want:    `x-handler="myHandler"`,
		},
		{
			name:    "multiple handlers",
			handler: "[checkAuth, loadUser]",
			mods:    nil,
			want:    `x-handler="[checkAuth, loadUser]"`,
		},
		{
			name:    "global handler",
			handler: "[logger, analytics]",
			mods:    []HandlerModifier{Global()},
			want:    `x-handler.global="[logger, analytics]"`,
		},
		{
			name:    "arrow function handler",
			handler: "(ctx) => console.log(ctx)",
			mods:    nil,
			want:    `x-handler="(ctx) =&gt; console.log(ctx)"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := XHandler(tt.handler, tt.mods...).Render(&buf); err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := buf.String()

			if !strings.Contains(got, tt.want) {
				t.Errorf("XHandler() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestRouterSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings map[string]any
		want     []string
	}{
		{
			name:     "empty settings",
			settings: map[string]any{},
			want:     []string{},
		},
		{
			name: "hash routing",
			settings: map[string]any{
				"hash": true,
			},
			want: []string{"<script>", "hash: true", "PineconeRouter.settings", "</script>"},
		},
		{
			name: "base path",
			settings: map[string]any{
				"basePath": "/app",
			},
			want: []string{"<script>", `basePath: "/app"`, "</script>"},
		},
		{
			name: "multiple settings",
			settings: map[string]any{
				"hash":     false,
				"basePath": "/blog",
				"targetID": "main",
			},
			want: []string{"<script>", "hash: false", `basePath: "/blog"`, `targetID: "main"`, "</script>"},
		},
		{
			name: "boolean settings",
			settings: map[string]any{
				"handleClicks": true,
				"pushState":    false,
			},
			want: []string{"<script>", "handleClicks: true", "pushState: false", "</script>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := RouterSettings(tt.settings).Render(&buf); err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			got := buf.String()

			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("RouterSettings() = %v, want to contain %v", got, want)
				}
			}
		})
	}
}

func TestNavigateTo(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "simple path",
			path: "/home",
			want: "$router.navigate('/home')",
		},
		{
			name: "path with parameter",
			path: "/users/123",
			want: "$router.navigate('/users/123')",
		},
		{
			name: "dynamic path with concatenation",
			path: "'/users/' + userId",
			want: "$router.navigate('/users/' + userId)",
		},
		{
			name: "path with expression",
			path: "user.profileUrl",
			want: "$router.navigate(user.profileUrl)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NavigateTo(tt.path)
			if got != tt.want {
				t.Errorf("NavigateTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterBack(t *testing.T) {
	got := RouterBack()
	want := "$history.back()"

	if got != want {
		t.Errorf("RouterBack() = %v, want %v", got, want)
	}
}

func TestRouterForward(t *testing.T) {
	got := RouterForward()
	want := "$history.forward()"

	if got != want {
		t.Errorf("RouterForward() = %v, want %v", got, want)
	}
}

func TestRouterCanGoBack(t *testing.T) {
	got := RouterCanGoBack()
	want := "$history.canGoBack()"

	if got != want {
		t.Errorf("RouterCanGoBack() = %v, want %v", got, want)
	}
}

func TestRouterCanGoForward(t *testing.T) {
	got := RouterCanGoForward()
	want := "$history.canGoForward()"

	if got != want {
		t.Errorf("RouterCanGoForward() = %v, want %v", got, want)
	}
}

func TestRouterParam(t *testing.T) {
	tests := []struct {
		name  string
		param string
		want  string
	}{
		{
			name:  "simple param",
			param: "id",
			want:  "$params.id",
		},
		{
			name:  "param with underscore",
			param: "user_id",
			want:  "$params.user_id",
		},
		{
			name:  "nested param",
			param: "data.name",
			want:  "$params.data.name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RouterParam(tt.param)
			if got != tt.want {
				t.Errorf("RouterParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterPath(t *testing.T) {
	got := RouterPath()
	want := "$router.context.path"

	if got != want {
		t.Errorf("RouterPath() = %v, want %v", got, want)
	}
}

func TestRouterLoading(t *testing.T) {
	got := RouterLoading()
	want := "$router.loading"

	if got != want {
		t.Errorf("RouterLoading() = %v, want %v", got, want)
	}
}

func TestRouteOptionWithName(t *testing.T) {
	cfg := &routeConfig{}
	WithRouteName("profile")(cfg)

	if cfg.name != "profile" {
		t.Errorf("WithRouteName() name = %v, want profile", cfg.name)
	}
}

func TestTemplateModifiers(t *testing.T) {
	t.Run("Preload", func(t *testing.T) {
		cfg := &templateConfig{}
		Preload()(cfg)

		if !cfg.preload {
			t.Error("Preload() should set preload to true")
		}
	})

	t.Run("TargetID", func(t *testing.T) {
		cfg := &templateConfig{}
		TargetID("app")(cfg)

		if cfg.targetID != "app" {
			t.Errorf("TargetID() targetID = %v, want app", cfg.targetID)
		}
	})

	t.Run("Interpolate", func(t *testing.T) {
		cfg := &templateConfig{}
		Interpolate()(cfg)

		if !cfg.interpolate {
			t.Error("Interpolate() should set interpolate to true")
		}
	})
}

func TestHandlerModifiers(t *testing.T) {
	t.Run("Global", func(t *testing.T) {
		cfg := &handlerConfig{}
		Global()(cfg)

		if !cfg.global {
			t.Error("Global() should set global to true")
		}
	})
}
