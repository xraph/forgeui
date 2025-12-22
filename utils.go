package forgeui

import (
	"fmt"
	"strings"
)

// CN (ClassName) merges class names, filtering empty strings
// This is useful for combining multiple class strings with conditional classes
func CN(classes ...string) string {
	result := make([]string, 0, len(classes))
	for _, c := range classes {
		c = strings.TrimSpace(c)
		if c != "" {
			result = append(result, c)
		}
	}

	return strings.Join(result, " ")
}

// If returns the value if condition is true, empty string otherwise
// Useful for conditional class application
func If(condition bool, value string) string {
	if condition {
		return value
	}

	return ""
}

// IfElse returns trueVal if condition is true, falseVal otherwise
// Useful for binary conditional class application
func IfElse(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}

	return falseVal
}

// MapGet returns map value or default value if key doesn't exist
// Generic helper for safely accessing maps with defaults
func MapGet[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if v, ok := m[key]; ok {
		return v
	}

	return defaultVal
}

// ToString converts various types to string for HTML attributes
func ToString(v any) string {
	return fmt.Sprintf("%v", v)
}
