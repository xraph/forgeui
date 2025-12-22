package bridge

import (
	"encoding/json"
	"net/http"
)

// FunctionInfo contains information about a registered function
type FunctionInfo struct {
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	RequireAuth  bool       `json:"requireAuth"`
	RequireRoles []string   `json:"requireRoles,omitempty"`
	RateLimit    int        `json:"rateLimit,omitempty"`
	TypeInfo     TypeInfo   `json:"typeInfo"`
}

// IntrospectionHandler handles introspection requests
type IntrospectionHandler struct {
	bridge *Bridge
}

// NewIntrospectionHandler creates a new introspection handler
func NewIntrospectionHandler(bridge *Bridge) *IntrospectionHandler {
	return &IntrospectionHandler{bridge: bridge}
}

// ServeHTTP handles introspection HTTP requests
func (h *IntrospectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get all registered functions
	functionNames := h.bridge.ListFunctions()

	functions := make([]FunctionInfo, 0, len(functionNames))

	for _, name := range functionNames {
		fn, err := h.bridge.GetFunction(name)
		if err != nil {
			continue
		}

		info := FunctionInfo{
			Name:         fn.Name,
			Description:  fn.Description,
			RequireAuth:  fn.RequireAuth,
			RequireRoles: fn.RequireRoles,
			RateLimit:    fn.RateLimit,
			TypeInfo:     fn.GetTypeInfo(),
		}

		functions = append(functions, info)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"functions": functions,
		"count":     len(functions),
	})
}

// IntrospectionHandler returns an HTTP handler for function introspection
func (b *Bridge) IntrospectionHandler() http.Handler {
	return NewIntrospectionHandler(b)
}

// GetFunctionInfo returns information about a specific function
func (b *Bridge) GetFunctionInfo(name string) (*FunctionInfo, error) {
	fn, err := b.GetFunction(name)
	if err != nil {
		return nil, err
	}

	info := &FunctionInfo{
		Name:         fn.Name,
		Description:  fn.Description,
		RequireAuth:  fn.RequireAuth,
		RequireRoles: fn.RequireRoles,
		RateLimit:    fn.RateLimit,
		TypeInfo:     fn.GetTypeInfo(),
	}

	return info, nil
}

// ListFunctionInfo returns information about all registered functions
func (b *Bridge) ListFunctionInfo() []FunctionInfo {
	functionNames := b.ListFunctions()
	functions := make([]FunctionInfo, 0, len(functionNames))

	for _, name := range functionNames {
		info, err := b.GetFunctionInfo(name)
		if err != nil {
			continue
		}
		functions = append(functions, *info)
	}

	return functions
}

