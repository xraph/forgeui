package alpine

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
)

// XRoute creates an x-route attribute for defining a route path.
// Routes can include parameters using :param, wildcards using *path,
// optional segments using :param?, and more.
//
// See Pinecone Router documentation for full route syntax:
// https://github.com/pinecone-router/router
//
// Example:
//
//	alpine.XRoute("/")                    // Static route
//	alpine.XRoute("/users/:id")           // Named parameter
//	alpine.XRoute("/files/*path")         // Wildcard
//	alpine.XRoute("/profile/:id?")        // Optional parameter
//	alpine.XRoute("notfound")             // 404 route
func XRoute(path string, opts ...RouteOption) g.Node {
	cfg := &routeConfig{path: path}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.name != "" {
		// Use x-route:name syntax for named routes
		return g.Attr("x-route:"+cfg.name, path)
	}

	return g.Attr("x-route", path)
}

// RouteOption configures route behavior.
type RouteOption func(*routeConfig)

type routeConfig struct {
	path string
	name string
}

// WithRouteName sets a name for the route.
// Named routes can be helpful for identification and debugging.
//
// Example:
//
//	alpine.XRoute("/profile", WithRouteName("profile"))
func WithRouteName(name string) RouteOption {
	return func(cfg *routeConfig) {
		cfg.name = name
	}
}

// XTemplate creates an x-template attribute for loading external templates.
// Templates can be loaded from URLs or rendered inline.
//
// Example:
//
//	alpine.XTemplate("/views/home.html")
//	alpine.XTemplate("['/header.html', '/home.html']")
func XTemplate(url string, mods ...TemplateModifier) g.Node {
	cfg := &templateConfig{url: url}
	for _, mod := range mods {
		mod(cfg)
	}

	// Build attribute name with modifiers
	attrName := "x-template"
	if cfg.preload {
		attrName += ".preload"
	}
	if cfg.targetID != "" {
		attrName += ".target." + cfg.targetID
	}
	if cfg.interpolate {
		attrName += ".interpolate"
	}

	return g.Attr(attrName, url)
}

// XTemplateInline creates an x-template attribute with no value,
// indicating that the template content is defined inline as children
// of the template element.
//
// Example:
//
//	html.Template(
//	    alpine.XRoute("/"),
//	    alpine.XTemplateInline(),
//	    html.H1(g.Text("Welcome")),
//	)
func XTemplateInline(mods ...TemplateModifier) g.Node {
	cfg := &templateConfig{}
	for _, mod := range mods {
		mod(cfg)
	}

	// Build attribute name with modifiers
	attrName := "x-template"
	if cfg.targetID != "" {
		attrName += ".target." + cfg.targetID
	}

	return g.Attr(attrName, "")
}

// TemplateModifier configures template behavior.
type TemplateModifier func(*templateConfig)

type templateConfig struct {
	url         string
	preload     bool
	targetID    string
	interpolate bool
}

// Preload causes the template to be fetched after page load
// at low priority, without waiting for the route to be matched.
//
// Example:
//
//	alpine.XTemplate("/views/404.html", Preload())
func Preload() TemplateModifier {
	return func(cfg *templateConfig) {
		cfg.preload = true
	}
}

// TargetID specifies an element ID where the template should be rendered.
//
// Example:
//
//	alpine.XTemplate("/views/home.html", TargetID("app"))
func TargetID(id string) TemplateModifier {
	return func(cfg *templateConfig) {
		cfg.targetID = id
	}
}

// Interpolate enables named route params in template URLs.
// For example, on route /dynamic/:name, it will fetch /api/dynamic/[name].html
//
// Example:
//
//	alpine.XTemplate("/api/dynamic/:name.html", Interpolate())
func Interpolate() TemplateModifier {
	return func(cfg *templateConfig) {
		cfg.interpolate = true
	}
}

// XHandler creates an x-handler attribute for route handler functions.
// Handlers are executed when a route is matched, before templates are loaded.
//
// Example:
//
//	alpine.XHandler("myHandler")                  // Single handler
//	alpine.XHandler("[checkAuth, loadUser]")      // Multiple handlers
func XHandler(handler string, mods ...HandlerModifier) g.Node {
	cfg := &handlerConfig{handler: handler}
	for _, mod := range mods {
		mod(cfg)
	}

	attrName := "x-handler"
	if cfg.global {
		attrName += ".global"
	}

	return g.Attr(attrName, handler)
}

// HandlerModifier configures handler behavior.
type HandlerModifier func(*handlerConfig)

type handlerConfig struct {
	handler string
	global  bool
}

// Global marks a handler as global, meaning it will run for every route.
// Global handlers always run before route-specific handlers.
//
// Example:
//
//	html.Div(
//	    alpine.XData(map[string]any{...}),
//	    alpine.XHandler("[logger, analytics]", Global()),
//	    // ... routes
//	)
func Global() HandlerModifier {
	return func(cfg *handlerConfig) {
		cfg.global = true
	}
}

// RouterSettings generates an inline script that configures Pinecone Router settings.
// This should be called during alpine:init event or early in page load.
//
// Available settings:
//   - hash: Enable hash routing (default: false)
//   - basePath: Base path of the site (default: "/")
//   - targetID: Default target element ID for templates
//   - handleClicks: Intercept link clicks (default: true)
//   - preload: Preload all templates (default: false)
//   - pushState: Call history.pushState() (default: true)
//
// Example:
//
//	alpine.RouterSettings(map[string]any{
//	    "hash": false,
//	    "basePath": "/app",
//	    "targetID": "app",
//	})
func RouterSettings(settings map[string]any) g.Node {
	if len(settings) == 0 {
		return g.Group([]g.Node{})
	}

	var parts []string
	for key, value := range settings {
		var valStr string
		switch v := value.(type) {
		case string:
			valStr = fmt.Sprintf("%q", v)
		case bool:
			if v {
				valStr = "true"
			} else {
				valStr = "false"
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			valStr = fmt.Sprintf("%d", v)
		case float32, float64:
			valStr = fmt.Sprintf("%v", v)
		default:
			valStr = fmt.Sprintf("%v", v)
		}
		parts = append(parts, fmt.Sprintf("%s: %s", key, valStr))
	}

	settingsStr := strings.Join(parts, ", ")
	script := fmt.Sprintf("document.addEventListener('alpine:init', () => { window.PineconeRouter.settings({%s}); })", settingsStr)

	return g.El("script", g.Raw(script))
}

// NavigateTo generates a JavaScript expression for navigating to a path.
// Use this in event handlers or x-init directives.
//
// Example:
//
//	html.Button(
//	    alpine.XClick(alpine.NavigateTo("/dashboard")),
//	    g.Text("Go to Dashboard"),
//	)
//
//	html.Button(
//	    alpine.XClick(alpine.NavigateTo("'/users/' + userId")),
//	    g.Text("View User"),
//	)
func NavigateTo(path string) string {
	// Check if path starts with '/' - it's a literal path
	if strings.HasPrefix(path, "/") {
		return fmt.Sprintf("$router.navigate('%s')", path)
	}
	// Check if path is already quoted - it's a literal string with quotes
	if strings.HasPrefix(path, "'") || strings.HasPrefix(path, "\"") {
		return fmt.Sprintf("$router.navigate(%s)", path)
	}
	// Otherwise it's a JavaScript expression (variable, concatenation, etc.)
	return fmt.Sprintf("$router.navigate(%s)", path)
}

// RouterBack generates a JavaScript expression for navigating back in history.
//
// Example:
//
//	html.Button(
//	    alpine.XClick(RouterBack()),
//	    g.Text("Back"),
//	)
func RouterBack() string {
	return "$history.back()"
}

// RouterForward generates a JavaScript expression for navigating forward in history.
//
// Example:
//
//	html.Button(
//	    alpine.XClick(RouterForward()),
//	    g.Text("Forward"),
//	)
func RouterForward() string {
	return "$history.forward()"
}

// RouterCanGoBack generates a JavaScript expression to check if back navigation is possible.
//
// Example:
//
//	html.Button(
//	    alpine.XBindDisabled("!"+RouterCanGoBack()),
//	    alpine.XClick(RouterBack()),
//	    g.Text("Back"),
//	)
func RouterCanGoBack() string {
	return "$history.canGoBack()"
}

// RouterCanGoForward generates a JavaScript expression to check if forward navigation is possible.
//
// Example:
//
//	html.Button(
//	    alpine.XBindDisabled("!"+RouterCanGoForward()),
//	    alpine.XClick(RouterForward()),
//	    g.Text("Forward"),
//	)
func RouterCanGoForward() string {
	return "$history.canGoForward()"
}

// RouterParam generates a JavaScript expression to access a route parameter.
//
// Example:
//
//	html.H1(alpine.XText(alpine.RouterParam("id")))
//	// Equivalent to: alpine.XText("$params.id")
func RouterParam(name string) string {
	return fmt.Sprintf("$params.%s", name)
}

// RouterPath generates a JavaScript expression to access the current path.
//
// Example:
//
//	html.Span(alpine.XText(alpine.RouterPath()))
func RouterPath() string {
	return "$router.context.path"
}

// RouterLoading generates a JavaScript expression to check if the router is loading.
//
// Example:
//
//	html.Div(
//	    alpine.XShow(alpine.RouterLoading()),
//	    g.Text("Loading..."),
//	)
func RouterLoading() string {
	return "$router.loading"
}
