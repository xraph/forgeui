package forgeui

// Props is the base interface for component properties
// Components can implement this interface to add validation
type Props interface {
	// Marker interface - components define their own props
}

// Option is a generic functional option type for configuring component properties
type Option[T any] func(*T)

// ApplyOptions applies a slice of options to props
// This is a helper function for components to apply functional options
func ApplyOptions[T any](props *T, opts []Option[T]) {
	for _, opt := range opts {
		opt(props)
	}
}

// BaseProps contains common properties shared by many components
type BaseProps struct {
	Class    string
	Disabled bool
	ID       string
}
