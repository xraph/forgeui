package bridge

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// parseParams parses JSON parameters into the target type
func parseParams(params json.RawMessage, targetType reflect.Type) (reflect.Value, error) {
	if len(params) == 0 {
		// No params provided, return zero value
		return reflect.New(targetType).Elem(), nil
	}

	// Create a new instance of the target type
	valuePtr := reflect.New(targetType)

	// Unmarshal JSON into the value
	if err := json.Unmarshal(params, valuePtr.Interface()); err != nil {
		return reflect.Value{}, &Error{
			Code:    ErrCodeInvalidParams,
			Message: "Failed to parse parameters",
			Data:    err.Error(),
		}
	}

	// Return the dereferenced value
	return valuePtr.Elem(), nil
}

// validateParams validates parameters against the target type
// This provides additional validation beyond JSON unmarshaling
func validateParams(value reflect.Value, targetType reflect.Type) error {
	// Check if value matches target type
	if value.Type() != targetType {
		return &Error{
			Code:    ErrCodeInvalidParams,
			Message: fmt.Sprintf("Parameter type mismatch: expected %s, got %s", targetType, value.Type()),
		}
	}

	// For struct types, validate required fields
	if targetType.Kind() == reflect.Struct {
		return validateStruct(value, targetType)
	}

	return nil
}

// validateStruct validates struct fields
func validateStruct(value reflect.Value, structType reflect.Type) error {
	numFields := structType.NumField()
	for i := range numFields {
		field := structType.Field(i)
		fieldValue := value.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Check if field has validation tags
		if err := validateField(field, fieldValue); err != nil {
			return err
		}
	}

	return nil
}

// validateField validates a single field
func validateField(field reflect.StructField, value reflect.Value) error {
	// Get JSON tag
	jsonTag := field.Tag.Get("json")

	// Check if field is required (no omitempty)
	isRequired := true

	if jsonTag != "" {
		// Parse tag to check for omitempty
		for i, c := range jsonTag {
			if c == ',' {
				if len(jsonTag) > i+1 {
					rest := jsonTag[i+1:]
					if rest == "omitempty" || rest[:9] == "omitempty" {
						isRequired = false
					}
				}

				break
			}
		}
	}

	// If required, check if value is zero
	if isRequired && isZero(value) {
		fieldName := field.Name

		if jsonTag != "" && jsonTag != "-" {
			// Extract name from JSON tag
			for i, c := range jsonTag {
				if c == ',' {
					fieldName = jsonTag[:i]
					break
				}
			}

			if fieldName == jsonTag {
				fieldName = jsonTag
			}
		}

		return &Error{
			Code:    ErrCodeInvalidParams,
			Message: fmt.Sprintf("Required field '%s' is missing or zero", fieldName),
			Data:    map[string]string{"field": fieldName},
		}
	}

	// Additional validation based on validate tag
	validateTag := field.Tag.Get("validate")
	if validateTag != "" {
		return validateWithTag(field.Name, value, validateTag)
	}

	return nil
}

// isZero checks if a value is the zero value for its type
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return false
}

// validateWithTag validates based on the validate tag
func validateWithTag(fieldName string, value reflect.Value, tag string) error {
	// Basic validation rules
	switch tag {
	case "required":
		if isZero(value) {
			return &Error{
				Code:    ErrCodeInvalidParams,
				Message: fmt.Sprintf("Field '%s' is required", fieldName),
				Data:    map[string]string{"field": fieldName},
			}
		}
	case "email":
		if value.Kind() == reflect.String {
			email := value.String()
			if !isValidEmail(email) {
				return &Error{
					Code:    ErrCodeInvalidParams,
					Message: fmt.Sprintf("Field '%s' must be a valid email", fieldName),
					Data:    map[string]string{"field": fieldName},
				}
			}
		}
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	atIndex := -1

	for i, c := range email {
		if c == '@' {
			if atIndex != -1 {
				// Multiple @ symbols
				return false
			}

			atIndex = i
		}
	}

	if atIndex <= 0 || atIndex >= len(email)-1 {
		return false
	}

	// Check for dot after @
	hasDot := false

	for i := atIndex + 1; i < len(email); i++ {
		if email[i] == '.' {
			hasDot = true
			break
		}
	}

	return hasDot
}
