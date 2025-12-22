package form

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator is a function that validates a value and returns an error if invalid
type Validator func(value string) error

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// Required validates that a value is not empty
//
// Example:
//
//	validator := form.Required()
//	if err := validator(""); err != nil {
//	    // Handle error: "This field is required"
//	}
func Required() Validator {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return NewValidationError("", "This field is required")
		}
		return nil
	}
}

// RequiredWithMessage validates that a value is not empty with a custom message
func RequiredWithMessage(message string) Validator {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return NewValidationError("", message)
		}
		return nil
	}
}

// Email validates that a value is a valid email address
//
// Example:
//
//	validator := form.Email()
//	if err := validator("invalid-email"); err != nil {
//	    // Handle error: "Invalid email address"
//	}
func Email() Validator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		if !emailRegex.MatchString(value) {
			return NewValidationError("", "Invalid email address")
		}
		return nil
	}
}

// MinLength validates that a value has at least the specified length
//
// Example:
//
//	validator := form.MinLength(8)
//	if err := validator("short"); err != nil {
//	    // Handle error: "Must be at least 8 characters"
//	}
func MinLength(min int) Validator {
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		if len(value) < min {
			return NewValidationError("", fmt.Sprintf("Must be at least %d characters", min))
		}
		return nil
	}
}

// MaxLength validates that a value does not exceed the specified length
//
// Example:
//
//	validator := form.MaxLength(100)
//	if err := validator("very long text..."); err != nil {
//	    // Handle error: "Must not exceed 100 characters"
//	}
func MaxLength(max int) Validator {
	return func(value string) error {
		if len(value) > max {
			return NewValidationError("", fmt.Sprintf("Must not exceed %d characters", max))
		}
		return nil
	}
}

// Pattern validates that a value matches the specified regex pattern
//
// Example:
//
//	validator := form.Pattern(`^\d{3}-\d{3}-\d{4}$`, "Invalid phone number format")
//	if err := validator("123-456-7890"); err != nil {
//	    // Handle error
//	}
func Pattern(pattern, message string) Validator {
	regex := regexp.MustCompile(pattern)
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		if !regex.MatchString(value) {
			if message == "" {
				message = "Invalid format"
			}
			return NewValidationError("", message)
		}
		return nil
	}
}

// URL validates that a value is a valid URL
//
// Example:
//
//	validator := form.URL()
//	if err := validator("https://example.com"); err != nil {
//	    // Handle error: "Invalid URL"
//	}
func URL() Validator {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		if !urlRegex.MatchString(value) {
			return NewValidationError("", "Invalid URL")
		}
		return nil
	}
}

// Numeric validates that a value contains only numeric characters
//
// Example:
//
//	validator := form.Numeric()
//	if err := validator("12345"); err != nil {
//	    // Handle error: "Must contain only numbers"
//	}
func Numeric() Validator {
	numericRegex := regexp.MustCompile(`^\d+$`)
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		if !numericRegex.MatchString(value) {
			return NewValidationError("", "Must contain only numbers")
		}
		return nil
	}
}

// Alpha validates that a value contains only alphabetic characters
//
// Example:
//
//	validator := form.Alpha()
//	if err := validator("abc"); err != nil {
//	    // Handle error: "Must contain only letters"
//	}
func Alpha() Validator {
	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		if !alphaRegex.MatchString(value) {
			return NewValidationError("", "Must contain only letters")
		}
		return nil
	}
}

// AlphaNumeric validates that a value contains only alphanumeric characters
//
// Example:
//
//	validator := form.AlphaNumeric()
//	if err := validator("abc123"); err != nil {
//	    // Handle error: "Must contain only letters and numbers"
//	}
func AlphaNumeric() Validator {
	alphaNumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		if !alphaNumRegex.MatchString(value) {
			return NewValidationError("", "Must contain only letters and numbers")
		}
		return nil
	}
}

// In validates that a value is in the specified list
//
// Example:
//
//	validator := form.In([]string{"small", "medium", "large"})
//	if err := validator("xlarge"); err != nil {
//	    // Handle error: "Must be one of: small, medium, large"
//	}
func In(allowed []string) Validator {
	return func(value string) error {
		if value == "" {
			return nil // Allow empty unless combined with Required()
		}
		for _, a := range allowed {
			if value == a {
				return nil
			}
		}
		return NewValidationError("", fmt.Sprintf("Must be one of: %s", strings.Join(allowed, ", ")))
	}
}

// Combine combines multiple validators into one
//
// Example:
//
//	validator := form.Combine(
//	    form.Required(),
//	    form.Email(),
//	    form.MaxLength(100),
//	)
//	if err := validator("user@example.com"); err != nil {
//	    // Handle first validation error
//	}
func Combine(validators ...Validator) Validator {
	return func(value string) error {
		for _, validator := range validators {
			if err := validator(value); err != nil {
				return err
			}
		}
		return nil
	}
}

// Validate validates a value against multiple validators
func Validate(value string, validators ...Validator) error {
	return Combine(validators...)(value)
}

// ValidateField validates a field with a name for better error messages
func ValidateField(field, value string, validators ...Validator) error {
	err := Validate(value, validators...)
	if err != nil {
		if ve, ok := err.(*ValidationError); ok {
			ve.Field = field
			return ve
		}
		return NewValidationError(field, err.Error())
	}
	return nil
}

