package alpine

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

// XRoute creates an x-route attribute for defining a route path.
func XRoute(path string, opts ...RouteOption) templ.Attributes {
	cfg := &routeConfig{path: path}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.name != "" {
		return templ.Attributes{"x-route:" + cfg.name: path}
	}

	return templ.Attributes{"x-route": path}
}

// RouteOption configures route behavior.
type RouteOption func(*routeConfig)

type routeConfig struct {
	path string
	name string
}

// WithRouteName sets a name for the route.
func WithRouteName(name string) RouteOption {
	return func(cfg *routeConfig) {
		cfg.name = name
	}
}

// XTemplate creates an x-template attribute for loading external templates.
func XTemplate(url string, mods ...TemplateModifier) templ.Attributes {
	cfg := &templateConfig{url: url}
	for _, mod := range mods {
		mod(cfg)
	}

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

	return templ.Attributes{attrName: url}
}

// XTemplateInline creates an x-template attribute with no value,
// indicating that the template content is defined inline.
func XTemplateInline(mods ...TemplateModifier) templ.Attributes {
	cfg := &templateConfig{}
	for _, mod := range mods {
		mod(cfg)
	}

	attrName := "x-template"
	if cfg.targetID != "" {
		attrName += ".target." + cfg.targetID
	}

	return templ.Attributes{attrName: ""}
}

// TemplateModifier configures template behavior.
type TemplateModifier func(*templateConfig)

type templateConfig struct {
	url         string
	preload     bool
	targetID    string
	interpolate bool
}

// Preload causes the template to be fetched after page load at low priority.
func Preload() TemplateModifier {
	return func(cfg *templateConfig) {
		cfg.preload = true
	}
}

// TargetID specifies an element ID where the template should be rendered.
func TargetID(id string) TemplateModifier {
	return func(cfg *templateConfig) {
		cfg.targetID = id
	}
}

// Interpolate enables named route params in template URLs.
func Interpolate() TemplateModifier {
	return func(cfg *templateConfig) {
		cfg.interpolate = true
	}
}

// XHandler creates an x-handler attribute for route handler functions.
func XHandler(handler string, mods ...HandlerModifier) templ.Attributes {
	cfg := &handlerConfig{handler: handler}
	for _, mod := range mods {
		mod(cfg)
	}

	attrName := "x-handler"
	if cfg.global {
		attrName += ".global"
	}

	return templ.Attributes{attrName: handler}
}

// HandlerModifier configures handler behavior.
type HandlerModifier func(*handlerConfig)

type handlerConfig struct {
	handler string
	global  bool
}

// Global marks a handler as global, meaning it will run for every route.
func Global() HandlerModifier {
	return func(cfg *handlerConfig) {
		cfg.global = true
	}
}

// RouterSettingsJS generates the inline JavaScript code for configuring Pinecone Router.
// Use this in a templ script block or with templ.Raw().
func RouterSettingsJS(settings map[string]any) string {
	if len(settings) == 0 {
		return ""
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
	return fmt.Sprintf("document.addEventListener('alpine:init', () => { window.PineconeRouter.settings({%s}); })", settingsStr)
}

// NavigateTo generates a JavaScript expression for navigating to a path.
func NavigateTo(path string) string {
	if strings.HasPrefix(path, "/") {
		return fmt.Sprintf("$router.navigate('%s')", path)
	}
	if strings.HasPrefix(path, "'") || strings.HasPrefix(path, "\"") {
		return fmt.Sprintf("$router.navigate(%s)", path)
	}
	return fmt.Sprintf("$router.navigate(%s)", path)
}

// RouterBack generates a JavaScript expression for navigating back in history.
func RouterBack() string {
	return "$history.back()"
}

// RouterForward generates a JavaScript expression for navigating forward in history.
func RouterForward() string {
	return "$history.forward()"
}

// RouterCanGoBack generates a JavaScript expression to check if back navigation is possible.
func RouterCanGoBack() string {
	return "$history.canGoBack()"
}

// RouterCanGoForward generates a JavaScript expression to check if forward navigation is possible.
func RouterCanGoForward() string {
	return "$history.canGoForward()"
}

// RouterParam generates a JavaScript expression to access a route parameter.
func RouterParam(name string) string {
	return fmt.Sprintf("$params.%s", name)
}

// RouterPath generates a JavaScript expression to access the current path.
func RouterPath() string {
	return "$router.context.path"
}

// RouterLoading generates a JavaScript expression to check if the router is loading.
func RouterLoading() string {
	return "$router.loading"
}
