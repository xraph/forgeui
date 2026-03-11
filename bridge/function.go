package bridge

import (
	"fmt"
	"reflect"
	"time"

	"github.com/a-h/templ"
)

// SignatureType describes the function signature shape
type SignatureType int

const (
	// SigInputOutput: func(Context, Input) (Output, error)
	SigInputOutput SignatureType = iota
	// SigOutput: func(Context) (Output, error) - no input
	SigOutput
	// SigInputOnly: func(Context, Input) error - no output value
	SigInputOnly
	// SigVoid: func(Context) error - no input, no output
	SigVoid
)

// Function represents a registered bridge function
type Function struct {
	// Name is the function's registered name
	Name string

	// Handler is the actual function to call
	Handler reflect.Value

	// InputType is the type of the input parameter (nil if no input)
	InputType reflect.Type

	// OutputType is the type of the output value (nil if no output)
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

	// SignatureType describes the handler's signature shape
	SignatureType SignatureType

	// HasInput indicates the handler accepts an input parameter
	HasInput bool

	// HasOutput indicates the handler returns an output value (beyond error)
	HasOutput bool

	// ReturnsHTML indicates the output type implements templ.Component
	ReturnsHTML bool

	// AllowedMethods restricts which HTTP methods the HTMX handler accepts
	AllowedMethods []string

	// Renderer converts output data to a templ.Component for HTMX rendering
	Renderer func(any) templ.Component

	// HTMXTriggers are event names to set in the HX-Trigger response header
	HTMXTriggers []string

	// HTMXRedirect is a URL to set in the HX-Redirect response header
	HTMXRedirect string

	// HTMXReswap overrides the swap method via HX-Reswap response header
	HTMXReswap string

	// HTMXRetarget overrides the target element via HX-Retarget response header
	HTMXRetarget string

	// LaxValidation skips auto-required field validation
	LaxValidation bool
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

// WithHTTPMethod restricts which HTTP methods the HTMX handler accepts.
// If not set, methods are auto-detected: GET+POST for functions without input,
// POST only for functions with input.
func WithHTTPMethod(methods ...string) FunctionOption {
	return func(f *Function) {
		f.AllowedMethods = methods
	}
}

// WithHTMXTrigger sets event names for the HX-Trigger response header
func WithHTMXTrigger(events ...string) FunctionOption {
	return func(f *Function) {
		f.HTMXTriggers = events
	}
}

// WithHTMXRedirect sets a URL for the HX-Redirect response header
func WithHTMXRedirect(url string) FunctionOption {
	return func(f *Function) {
		f.HTMXRedirect = url
	}
}

// WithHTMXReswap overrides the swap method via HX-Reswap response header
func WithHTMXReswap(swap string) FunctionOption {
	return func(f *Function) {
		f.HTMXReswap = swap
	}
}

// WithHTMXRetarget overrides the target element via HX-Retarget response header
func WithHTMXRetarget(target string) FunctionOption {
	return func(f *Function) {
		f.HTMXRetarget = target
	}
}

// WithLaxValidation skips auto-required field validation.
// In lax mode, only fields with an explicit validate:"required" tag are validated.
func WithLaxValidation() FunctionOption {
	return func(f *Function) {
		f.LaxValidation = true
	}
}

// WithRenderer sets a renderer function that converts output data to a templ.Component
// for HTMX HTML responses. The type parameter T must match the handler's output type.
func WithRenderer[T any](renderer func(T) templ.Component) FunctionOption {
	return func(f *Function) {
		f.Renderer = func(data any) templ.Component {
			return renderer(data.(T))
		}
	}
}

// validateFunction validates a function signature and returns its shape.
//
// Supported signatures:
//   - func(Context, Input) (Output, error)  → SigInputOutput
//   - func(Context) (Output, error)         → SigOutput
//   - func(Context, Input) error            → SigInputOnly
//   - func(Context) error                   → SigVoid
func validateFunction(fn any) (SignatureType, error) {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		return 0, fmt.Errorf("handler must be a function, got %s", fnType.Kind())
	}

	numIn := fnType.NumIn()
	numOut := fnType.NumOut()

	if numIn < 1 || numIn > 2 {
		return 0, fmt.Errorf("handler must have 1 or 2 parameters (Context[, Input]), got %d", numIn)
	}

	contextType := reflect.TypeFor[Context]()
	if !fnType.In(0).Implements(contextType) {
		return 0, fmt.Errorf("first parameter must implement bridge.Context, got %s", fnType.In(0))
	}

	if numOut < 1 || numOut > 2 {
		return 0, fmt.Errorf("handler must return 1 or 2 values ([Output, ]error), got %d", numOut)
	}

	errorType := reflect.TypeFor[error]()

	switch {
	case numIn == 2 && numOut == 2:
		// func(Context, Input) (Output, error)
		if !fnType.Out(1).Implements(errorType) {
			return 0, fmt.Errorf("second return value must be error, got %s", fnType.Out(1))
		}
		return SigInputOutput, nil

	case numIn == 1 && numOut == 2:
		// func(Context) (Output, error)
		if !fnType.Out(1).Implements(errorType) {
			return 0, fmt.Errorf("second return value must be error, got %s", fnType.Out(1))
		}
		return SigOutput, nil

	case numIn == 2 && numOut == 1:
		// func(Context, Input) error
		if !fnType.Out(0).Implements(errorType) {
			return 0, fmt.Errorf("return value must be error, got %s", fnType.Out(0))
		}
		return SigInputOnly, nil

	case numIn == 1 && numOut == 1:
		// func(Context) error
		if !fnType.Out(0).Implements(errorType) {
			return 0, fmt.Errorf("return value must be error, got %s", fnType.Out(0))
		}
		return SigVoid, nil
	}

	return 0, fmt.Errorf("unsupported signature: %d inputs, %d outputs", numIn, numOut)
}

// analyzeFunction extracts type information from a function
func analyzeFunction(fn any) (*Function, error) {
	sigType, err := validateFunction(fn)
	if err != nil {
		return nil, err
	}

	fnType := reflect.TypeOf(fn)
	templComponentType := reflect.TypeFor[templ.Component]()

	f := &Function{
		Handler:       reflect.ValueOf(fn),
		SignatureType: sigType,
	}

	switch sigType {
	case SigInputOutput:
		f.InputType = fnType.In(1)
		f.OutputType = fnType.Out(0)
		f.HasInput = true
		f.HasOutput = true
	case SigOutput:
		f.OutputType = fnType.Out(0)
		f.HasOutput = true
	case SigInputOnly:
		f.InputType = fnType.In(1)
		f.HasInput = true
	case SigVoid:
		// no input, no output
	}

	// Detect if output type implements templ.Component
	if f.HasOutput && f.OutputType.Implements(templComponentType) {
		f.ReturnsHTML = true
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
		Name: f.Name,
	}

	if f.InputType != nil {
		info.InputType = f.InputType.String()
		// Extract field information from input struct
		if f.InputType.Kind() == reflect.Struct {
			info.Fields = extractFields(f.InputType)
		} else if f.InputType.Kind() == reflect.Ptr && f.InputType.Elem().Kind() == reflect.Struct {
			info.Fields = extractFields(f.InputType.Elem())
		}
	}

	if f.OutputType != nil {
		info.OutputType = f.OutputType.String()
	}

	return info
}

// extractFields extracts field information from a struct type
func extractFields(t reflect.Type) []FieldInfo {
	var fields []FieldInfo

	numFields := t.NumField()
	for i := range numFields {
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
