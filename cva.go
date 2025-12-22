package forgeui

import "strings"

// CVA (Class Variance Authority) manages component class variants
// It provides a type-safe way to handle component variants and their combinations
type CVA struct {
	base      []string
	variants  map[string]map[string][]string
	defaults  map[string]string
	compounds []CompoundVariant
}

// CompoundVariant applies classes when multiple conditions match
type CompoundVariant struct {
	Conditions map[string]string
	Classes    []string
}

// NewCVA creates a new CVA instance with base classes
// Base classes are always included in the output
func NewCVA(base ...string) *CVA {
	return &CVA{
		base:     base,
		variants: make(map[string]map[string][]string),
		defaults: make(map[string]string),
	}
}

// Variant adds a variant dimension with its possible values and corresponding classes
// Example: cva.Variant("size", map[string][]string{"sm": {"text-sm", "px-2"}, "lg": {"text-lg", "px-4"}})
func (c *CVA) Variant(name string, options map[string][]string) *CVA {
	c.variants[name] = options
	return c
}

// Default sets the default value for a variant dimension
// This value is used when no value is explicitly provided for that variant
func (c *CVA) Default(name, value string) *CVA {
	c.defaults[name] = value
	return c
}

// Compound adds a compound variant rule that applies classes when multiple conditions match
// Example: cva.Compound(map[string]string{"size": "sm", "variant": "primary"}, "extra-class")
func (c *CVA) Compound(conditions map[string]string, classes ...string) *CVA {
	c.compounds = append(c.compounds, CompoundVariant{
		Conditions: conditions,
		Classes:    classes,
	})
	return c
}

// Classes generates the final class string based on the provided variant values
// It combines base classes, variant-specific classes, and compound variant classes
func (c *CVA) Classes(variants map[string]string) string {
	// Pre-allocate with reasonable capacity
	result := make([]string, 0, len(c.base)+len(c.variants)*2+len(c.compounds))

	// Add base classes
	result = append(result, c.base...)

	// Apply variant classes
	for variantName, options := range c.variants {
		value, ok := variants[variantName]
		if !ok {
			// Use default value if not provided
			value = c.defaults[variantName]
		}
		if classes, exists := options[value]; exists {
			result = append(result, classes...)
		}
	}

	// Apply compound variants
	for _, cv := range c.compounds {
		if c.matchesCompound(cv, variants) {
			result = append(result, cv.Classes...)
		}
	}

	return strings.Join(result, " ")
}

// matchesCompound checks if all conditions of a compound variant are met
func (c *CVA) matchesCompound(cv CompoundVariant, variants map[string]string) bool {
	for key, expectedValue := range cv.Conditions {
		actualValue, ok := variants[key]
		if !ok {
			// Use default value if not provided
			actualValue = c.defaults[key]
		}
		if actualValue != expectedValue {
			return false
		}
	}
	return true
}
