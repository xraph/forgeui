package router

import "net/http"

// HTTP method constants
const (
	MethodGet     = http.MethodGet
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodOptions = http.MethodOptions
	MethodHead    = http.MethodHead
)

// ParamKey is the type for parameter keys
type ParamKey string

// Params stores route parameters extracted from the URL path
type Params map[string]string

// Get retrieves a parameter value by key
func (p Params) Get(key string) string {
	return p[key]
}

// Has checks if a parameter exists
func (p Params) Has(key string) bool {
	_, ok := p[key]
	return ok
}

