package bridge

import (
	"fmt"
	"reflect"
	"time"
)

// Function represents a registered bridge function
type Function struct {
	// Name is the function's registered name
	Name string

	// Handler is the actual function to call
	Handler reflect.Value

	// InputType is the type of the input parameter
	InputType reflect.Type

	// OutputType is the type of the output value
	OutputType reflect.Type

	// Description is the function's documentation
	Description string

	// RequireAuth indicates if authentication is required
	RequireAuth bool

	// RequireRoles specifies required user roles
	RequireRoles []string

	// Timeout is the maximum execution time
	Timeout time.Duration

	// RateLimit is the max requests per minute (0 = no limit)
	RateLimit int

	// Cacheable indicates if results can be cached
	Cacheable bool

	// CacheTTL is the cache time-to-live
	CacheTTL time.Duration
}

// FunctionOption configures a Function
type FunctionOption func(*Function)

// RequireAuth marks a function as requiring authentication
func RequireAuth() FunctionOption {
	return func(f *Function) {
		f.RequireAuth = true
	}
}

// RequireRoles specifies required user roles
func RequireRoles(roles ...string) FunctionOption {
	return func(f *Function) {
		f.RequireRoles = roles
	}
}

// WithFunctionTimeout sets a custom timeout for a function
func WithFunctionTimeout(d time.Duration) FunctionOption {
	return func(f *Function) {
		f.Timeout = d
	}
}

// WithRateLimit sets a rate limit (requests per minute)
func WithRateLimit(rpm int) FunctionOption {
	return func(f *Function) {
		f.RateLimit = rpm
	}
}

// WithFunctionCache enables result caching for a function
func WithFunctionCache(ttl time.Duration) FunctionOption {
	return func(f *Function) {
		f.Cacheable = true
		f.CacheTTL = ttl
	}
}

// WithDescription sets the function description
func WithDescription(desc string) FunctionOption {
	return func(f *Function) {
		f.Description = desc
	}
}

// validateFunction validates a function signature
// Expected signature: func(Context, InputType) (OutputType, error)
func validateFunction(fn any) error {
	fnType := reflect.TypeOf(fn)

	// Must be a function
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("handler must be a function, got %s", fnType.Kind())
	}

	// Must have exactly 2 parameters
	if fnType.NumIn() != 2 {
		return fmt.Errorf("handler must have 2 parameters (Context, Input), got %d", fnType.NumIn())
	}

	// First parameter must implement Context interface
	contextType := reflect.TypeFor[Context]()
	if !fnType.In(0).Implements(contextType) {
		return fmt.Errorf("first parameter must implement bridge.Context, got %s", fnType.In(0))
	}

	// Must return exactly 2 values
	if fnType.NumOut() != 2 {
		return fmt.Errorf("handler must return 2 values (result, error), got %d", fnType.NumOut())
	}

	// Second return value must be error
	errorType := reflect.TypeFor[error]()
	if !fnType.Out(1).Implements(errorType) {
		return fmt.Errorf("second return value must be error, got %s", fnType.Out(1))
	}

	return nil
}

// analyzeFunction extracts type information from a function
func analyzeFunction(fn any) (*Function, error) {
	if err := validateFunction(fn); err != nil {
		return nil, err
	}

	fnType := reflect.TypeOf(fn)

	f := &Function{
		Handler:    reflect.ValueOf(fn),
		InputType:  fnType.In(1),
		OutputType: fnType.Out(0),
		// Timeout is 0 by default - executeWithTimeout will use bridge config timeout
	}

	return f, nil
}

// TypeInfo returns information about the function's types
type TypeInfo struct {
	Name       string      `json:"name"`
	InputType  string      `json:"inputType"`
	OutputType string      `json:"outputType"`
	Fields     []FieldInfo `json:"fields,omitempty"`
}

// FieldInfo describes a struct field
type FieldInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	JSONName string `json:"jsonName"`
	Required bool   `json:"required"`
}

// GetTypeInfo returns type information for the function
func (f *Function) GetTypeInfo() TypeInfo {
	info := TypeInfo{
		Name:       f.Name,
		InputType:  f.InputType.String(),
		OutputType: f.OutputType.String(),
	}

	// Extract field information from input struct
	if f.InputType.Kind() == reflect.Struct {
		info.Fields = extractFields(f.InputType)
	} else if f.InputType.Kind() == reflect.Ptr && f.InputType.Elem().Kind() == reflect.Struct {
		info.Fields = extractFields(f.InputType.Elem())
	}

	return info
}

// extractFields extracts field information from a struct type
func extractFields(t reflect.Type) []FieldInfo {
	var fields []FieldInfo

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		jsonName := field.Name
		required := true

		if jsonTag != "" {
			// Parse json tag
			if jsonTag == "-" {
				continue
			}
			// Extract name from tag (before comma)
			commaIdx := -1

			for idx, c := range jsonTag {
				if c == ',' {
					commaIdx = idx
					break
				}
			}

			if commaIdx >= 0 {
				jsonName = jsonTag[:commaIdx]
				// Check for omitempty
				if len(jsonTag) > commaIdx+1 && jsonTag[commaIdx+1:] == "omitempty" {
					required = false
				}
			} else {
				// No comma, the whole tag is the name
				jsonName = jsonTag
			}
			// If name part is empty, use field name
			if jsonName == "" {
				jsonName = field.Name
			}
		}

		fields = append(fields, FieldInfo{
			Name:     field.Name,
			Type:     field.Type.String(),
			JSONName: jsonName,
			Required: required,
		})
	}

	return fields
}
