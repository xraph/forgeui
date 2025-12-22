package router

import (
	"fmt"
	"strings"
)

// URL generates a URL for a named route with the given parameters
// Example: router.URL("user", 123) -> "/users/123"
func (r *Router) URL(name string, params ...any) string {
	r.mu.RLock()
	route, ok := r.namedRoutes[name]
	r.mu.RUnlock()

	if !ok {
		return ""
	}

	return route.URL(params...)
}

// URL generates a URL from this route's pattern with the given parameters
func (r *Route) URL(params ...any) string {
	pattern := r.Pattern
	paramIdx := 0

	// Replace parameters in pattern
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	result := make([]string, 0, len(segments))

	for _, segment := range segments {
		switch {
		case strings.HasPrefix(segment, ":"):
			// Named parameter - replace with provided value
			if paramIdx < len(params) {
				result = append(result, fmt.Sprint(params[paramIdx]))
				paramIdx++
			} else {
				// Not enough parameters provided
				result = append(result, segment)
			}
		case strings.HasPrefix(segment, "*"):
			// Wildcard - replace with provided value
			if paramIdx < len(params) {
				result = append(result, fmt.Sprint(params[paramIdx]))
				paramIdx++
			} else {
				result = append(result, segment)
			}
		default:
			// Static segment
			result = append(result, segment)
		}
	}

	if len(result) == 0 {
		return "/"
	}

	return "/" + strings.Join(result, "/")
}

// URLMap generates a URL from this route's pattern with named parameters
// Example: route.URLMap(map[string]interface{}{"id": 123, "action": "edit"})
func (r *Route) URLMap(params map[string]any) string {
	pattern := r.Pattern

	// Replace parameters in pattern
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	result := make([]string, 0, len(segments))

	for _, segment := range segments {
		if after, ok := strings.CutPrefix(segment, ":"); ok {
			// Named parameter
			paramName := after
			if val, ok := params[paramName]; ok {
				result = append(result, fmt.Sprint(val))
			} else {
				result = append(result, segment)
			}
		} else if after, ok := strings.CutPrefix(segment, "*"); ok {
			// Wildcard
			paramName := after
			if paramName == "" {
				paramName = "wildcard"
			}

			if val, ok := params[paramName]; ok {
				result = append(result, fmt.Sprint(val))
			} else {
				result = append(result, segment)
			}
		} else {
			// Static segment
			result = append(result, segment)
		}
	}

	if len(result) == 0 {
		return "/"
	}

	return "/" + strings.Join(result, "/")
}
