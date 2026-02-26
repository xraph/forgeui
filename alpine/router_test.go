package alpine

import (
	"strings"
	"testing"
)

func TestXRoute(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		opts    []RouteOption
		wantKey string
		wantVal string
	}{
		{
			name:    "static route",
			path:    "/",
			opts:    nil,
			wantKey: "x-route",
			wantVal: "/",
		},
		{
			name:    "route with param",
			path:    "/users/:id",
			opts:    nil,
			wantKey: "x-route",
			wantVal: "/users/:id",
		},
		{
			name:    "wildcard route",
			path:    "/files/*path",
			opts:    nil,
			wantKey: "x-route",
			wantVal: "/files/*path",
		},
		{
			name:    "optional param",
			path:    "/profile/:id?",
			opts:    nil,
			wantKey: "x-route",
			wantVal: "/profile/:id?",
		},
		{
			name:    "notfound route",
			path:    "notfound",
			opts:    nil,
			wantKey: "x-route",
			wantVal: "notfound",
		},
		{
			name:    "named route",
			path:    "/profile",
			opts:    []RouteOption{WithRouteName("profile")},
			wantKey: "x-route:profile",
			wantVal: "/profile",
		},
		{
			name:    "complex path",
			path:    "/movies/:title.(mp4|mov)",
			opts:    nil,
			wantKey: "x-route",
			wantVal: "/movies/:title.(mp4|mov)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := XRoute(tt.path, tt.opts...)

			if v, ok := attrs[tt.wantKey]; !ok || v != tt.wantVal {
				t.Errorf("XRoute() = %v, want %s=%s", attrs, tt.wantKey, tt.wantVal)
			}
		})
	}
}

func TestXTemplate(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		mods    []TemplateModifier
		wantKey string
		wantVal string
	}{
		{
			name:    "simple template",
			url:     "/views/home.html",
			mods:    nil,
			wantKey: "x-template",
			wantVal: "/views/home.html",
		},
		{
			name:    "template with preload",
			url:     "/views/404.html",
			mods:    []TemplateModifier{Preload()},
			wantKey: "x-template.preload",
			wantVal: "/views/404.html",
		},
		{
			name:    "template with target",
			url:     "/views/profile.html",
			mods:    []TemplateModifier{TargetID("app")},
			wantKey: "x-template.target.app",
			wantVal: "/views/profile.html",
		},
		{
			name:    "template with interpolate",
			url:     "/api/dynamic/:name.html",
			mods:    []TemplateModifier{Interpolate()},
			wantKey: "x-template.interpolate",
			wantVal: "/api/dynamic/:name.html",
		},
		{
			name:    "template with multiple modifiers",
			url:     "/views/user.html",
			mods:    []TemplateModifier{Preload(), TargetID("main")},
			wantKey: "x-template.preload.target.main",
			wantVal: "/views/user.html",
		},
		{
			name:    "array of templates",
			url:     "['/header.html', '/home.html']",
			mods:    nil,
			wantKey: "x-template",
			wantVal: "['/header.html', '/home.html']",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := XTemplate(tt.url, tt.mods...)

			if v, ok := attrs[tt.wantKey]; !ok || v != tt.wantVal {
				t.Errorf("XTemplate() = %v, want %s=%s", attrs, tt.wantKey, tt.wantVal)
			}
		})
	}
}

func TestXTemplateInline(t *testing.T) {
	tests := []struct {
		name    string
		mods    []TemplateModifier
		wantKey string
	}{
		{
			name:    "inline template",
			mods:    nil,
			wantKey: "x-template",
		},
		{
			name:    "inline template with target",
			mods:    []TemplateModifier{TargetID("app")},
			wantKey: "x-template.target.app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := XTemplateInline(tt.mods...)

			if v, ok := attrs[tt.wantKey]; !ok || v != "" {
				t.Errorf("XTemplateInline() = %v, want %s=\"\"", attrs, tt.wantKey)
			}
		})
	}
}

func TestXHandler(t *testing.T) {
	tests := []struct {
		name    string
		handler string
		mods    []HandlerModifier
		wantKey string
		wantVal string
	}{
		{
			name:    "single handler",
			handler: "myHandler",
			mods:    nil,
			wantKey: "x-handler",
			wantVal: "myHandler",
		},
		{
			name:    "multiple handlers",
			handler: "[checkAuth, loadUser]",
			mods:    nil,
			wantKey: "x-handler",
			wantVal: "[checkAuth, loadUser]",
		},
		{
			name:    "global handler",
			handler: "[logger, analytics]",
			mods:    []HandlerModifier{Global()},
			wantKey: "x-handler.global",
			wantVal: "[logger, analytics]",
		},
		{
			name:    "arrow function handler",
			handler: "(ctx) => console.log(ctx)",
			mods:    nil,
			wantKey: "x-handler",
			wantVal: "(ctx) => console.log(ctx)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := XHandler(tt.handler, tt.mods...)

			if v, ok := attrs[tt.wantKey]; !ok || v != tt.wantVal {
				t.Errorf("XHandler() = %v, want %s=%s", attrs, tt.wantKey, tt.wantVal)
			}
		})
	}
}

func TestRouterSettingsJS(t *testing.T) {
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
			want: []string{"hash: true", "PineconeRouter.settings"},
		},
		{
			name: "base path",
			settings: map[string]any{
				"basePath": "/app",
			},
			want: []string{`basePath: "/app"`, "PineconeRouter.settings"},
		},
		{
			name: "boolean settings",
			settings: map[string]any{
				"handleClicks": true,
				"pushState":    false,
			},
			want: []string{"handleClicks: true", "pushState: false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RouterSettingsJS(tt.settings)

			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("RouterSettingsJS() = %v, want to contain %v", got, want)
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
