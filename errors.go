package forgeui

import "fmt"

// ComponentError represents an error in component rendering
type ComponentError struct {
	Component string
	Message   string
	Err       error
}

func (e *ComponentError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Component, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Component, e.Message)
}

func (e *ComponentError) Unwrap() error {
	return e.Err
}

// ValidationError represents a props validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// PluginError represents a plugin-related error
type PluginError struct {
	Plugin  string
	Message string
	Err     error
}

func (e *PluginError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("plugin %s: %s: %v", e.Plugin, e.Message, e.Err)
	}
	return fmt.Sprintf("plugin %s: %s", e.Plugin, e.Message)
}

func (e *PluginError) Unwrap() error {
	return e.Err
}
