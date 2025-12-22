package bridge

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if len(email) < 3 || len(email) > 254 {
		return errors.New("email length must be between 3 and 254 characters")
	}

	// Basic email pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateRange validates that a number is within a range
func ValidateRange(value, minValue, maxValue int) error {
	if value < minValue {
		return fmt.Errorf("value %d is less than minimum %d", value, minValue)
	}

	if value > maxValue {
		return fmt.Errorf("value %d is greater than maximum %d", value, maxValue)
	}

	return nil
}

// ValidateFloatRange validates that a float is within a range
func ValidateFloatRange(value, minValue, maxValue float64) error {
	if value < minValue {
		return fmt.Errorf("value %f is less than minimum %f", value, minValue)
	}

	if value > maxValue {
		return fmt.Errorf("value %f is greater than maximum %f", value, maxValue)
	}

	return nil
}

// ValidatePattern validates that a string matches a pattern
func ValidatePattern(value, pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	if !regex.MatchString(value) {
		return fmt.Errorf("value does not match pattern %s", pattern)
	}

	return nil
}

// ValidateLength validates string length
func ValidateLength(value string, minLength, maxLength int) error {
	length := len(value)

	if length < minLength {
		return fmt.Errorf("string length %d is less than minimum %d", length, minLength)
	}

	if length > maxLength {
		return fmt.Errorf("string length %d is greater than maximum %d", length, maxLength)
	}

	return nil
}

// ValidateRequired validates that a value is not empty
func ValidateRequired(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("field is required")
	}

	return nil
}

// ValidateURL validates a URL
func ValidateURL(url string) error {
	if url == "" {
		return errors.New("URL is required")
	}

	// Basic URL pattern
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+(?:\.[a-zA-Z]{2,})+(?:/.*)?$`)
	if !urlRegex.MatchString(url) {
		return errors.New("invalid URL format")
	}

	return nil
}

// ValidateAlphanumeric validates that a string contains only alphanumeric characters
func ValidateAlphanumeric(value string) error {
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumericRegex.MatchString(value) {
		return errors.New("value must contain only alphanumeric characters")
	}

	return nil
}

// ValidateSlug validates a URL-friendly slug
func ValidateSlug(slug string) error {
	if slug == "" {
		return errors.New("slug is required")
	}

	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	if !slugRegex.MatchString(slug) {
		return errors.New("invalid slug format (use lowercase letters, numbers, and hyphens)")
	}

	return nil
}

// ValidateOneOf validates that a value is one of the allowed values
func ValidateOneOf(value string, allowed []string) error {
	if slices.Contains(allowed, value) {
		return nil
	}

	return fmt.Errorf("value must be one of: %s", strings.Join(allowed, ", "))
}
