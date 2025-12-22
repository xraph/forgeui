package router

import (
	"regexp"
	"strings"
)

// Route represents a single route in the router
type Route struct {
	Pattern    string
	Method     string
	Handler    PageHandler
	Middleware []Middleware
	Name       string
	Layout     string
	LoaderFn   LoaderFunc
	Metadata   *RouteMeta

	// Internal fields for matching
	regex      *regexp.Regexp
	paramNames []string
	priority   int // Lower is higher priority (static > param > wildcard)
}

// newRoute creates a new route with the given pattern and method
func newRoute(pattern, method string, handler PageHandler) *Route {
	r := &Route{
		Pattern:    pattern,
		Method:     method,
		Handler:    handler,
		Middleware: make([]Middleware, 0),
	}
	r.compile()
	return r
}

// compile parses the pattern and builds the regex for matching
func (r *Route) compile() {
	// Handle root path
	if r.Pattern == "/" {
		r.regex = regexp.MustCompile("^/$")
		r.priority = 0 // Highest priority for exact root
		return
	}

	// Remove trailing slash for consistency
	pattern := strings.TrimSuffix(r.Pattern, "/")
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	
	regexParts := make([]string, 0, len(segments))
	paramNames := make([]string, 0)
	priority := 0

	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			// Named parameter: :id -> capture group
			paramName := strings.TrimPrefix(segment, ":")
			paramNames = append(paramNames, paramName)
			regexParts = append(regexParts, "([^/]+)")
			priority += 10 // Parameters have lower priority than static
		} else if strings.HasPrefix(segment, "*") {
			// Wildcard: *path -> capture everything
			paramName := strings.TrimPrefix(segment, "*")
			if paramName == "" {
				paramName = "wildcard"
			}
			paramNames = append(paramNames, paramName)
			regexParts = append(regexParts, "(.+)")
			priority += 20 // Wildcards have lowest priority
		} else {
			// Static segment
			regexParts = append(regexParts, regexp.QuoteMeta(segment))
			// Static segments don't increase priority (priority = 0 is highest)
		}
	}

	// Build the final regex
	regexStr := "^/" + strings.Join(regexParts, "/") + "$"
	r.regex = regexp.MustCompile(regexStr)
	r.paramNames = paramNames
	r.priority = priority
}

// Match checks if the given path matches this route and extracts parameters
func (r *Route) Match(path string) (params Params, ok bool) {
	// Normalize path
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}

	matches := r.regex.FindStringSubmatch(path)
	if matches == nil {
		return nil, false
	}

	// Extract parameters
	params = make(Params)
	for i, name := range r.paramNames {
		if i+1 < len(matches) {
			params[name] = matches[i+1]
		}
	}

	return params, true
}

// WithMiddleware adds middleware to this route
func (r *Route) WithMiddleware(middleware ...Middleware) *Route {
	r.Middleware = append(r.Middleware, middleware...)
	return r
}

// WithName sets the name of this route
func (r *Route) WithName(name string) *Route {
	r.Name = name
	return r
}

// WithLoader sets the loader function for this route
func (r *Route) WithLoader(loader LoaderFunc) *Route {
	r.LoaderFn = loader
	return r
}

// WithMeta sets the metadata for this route
func (r *Route) WithMeta(meta *RouteMeta) *Route {
	r.Metadata = meta
	return r
}
