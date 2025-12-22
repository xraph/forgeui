package form

import (
	"strings"
	"testing"
)

func TestRequired(t *testing.T) {
	validator := Required()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"valid value", "test", false},
		{"value with spaces", "  test  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	validator := Email()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"valid email", "user@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", false},
		{"invalid email - no @", "userexample.com", true},
		{"invalid email - no domain", "user@", true},
		{"invalid email - no TLD", "user@example", true},
		{"invalid email - spaces", "user @example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	validator := MinLength(5)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"too short", "test", true},
		{"exact length", "tests", false},
		{"longer", "testing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinLength() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && !strings.Contains(err.Error(), "at least 5") {
				t.Errorf("MinLength() error message = %v, want 'at least 5'", err)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	validator := MaxLength(10)

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false},
		{"short", "test", false},
		{"exact length", "1234567890", false},
		{"too long", "12345678901", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxLength() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && !strings.Contains(err.Error(), "not exceed 10") {
				t.Errorf("MaxLength() error message = %v, want 'not exceed 10'", err)
			}
		})
	}
}

func TestPattern(t *testing.T) {
	validator := Pattern(`^\d{3}-\d{3}-\d{4}$`, "Invalid phone number")

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"valid pattern", "123-456-7890", false},
		{"invalid pattern - no dashes", "1234567890", true},
		{"invalid pattern - wrong format", "12-345-6789", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pattern() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURL(t *testing.T) {
	validator := URL()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"valid http URL", "http://example.com", false},
		{"valid https URL", "https://example.com", false},
		{"valid URL with path", "https://example.com/path", false},
		{"invalid URL - no protocol", "example.com", true},
		{"invalid URL - spaces", "https://example .com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNumeric(t *testing.T) {
	validator := Numeric()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"valid numeric", "12345", false},
		{"invalid - letters", "abc", true},
		{"invalid - mixed", "123abc", true},
		{"invalid - spaces", "123 456", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Numeric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlpha(t *testing.T) {
	validator := Alpha()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"valid alpha", "abc", false},
		{"valid mixed case", "AbCdEf", false},
		{"invalid - numbers", "abc123", true},
		{"invalid - spaces", "abc def", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Alpha() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlphaNumeric(t *testing.T) {
	validator := AlphaNumeric()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"valid alphanumeric", "abc123", false},
		{"valid letters only", "abc", false},
		{"valid numbers only", "123", false},
		{"invalid - special chars", "abc-123", true},
		{"invalid - spaces", "abc 123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlphaNumeric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIn(t *testing.T) {
	validator := In([]string{"small", "medium", "large"})

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty string", "", false}, // Allow empty
		{"valid - small", "small", false},
		{"valid - medium", "medium", false},
		{"valid - large", "large", false},
		{"invalid - xlarge", "xlarge", true},
		{"invalid - wrong case", "Small", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("In() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCombine(t *testing.T) {
	validator := Combine(
		Required(),
		Email(),
		MaxLength(50),
	)

	tests := []struct {
		name    string
		value   string
		wantErr bool
		errMsg  string
	}{
		{"empty - fails Required", "", true, "required"},
		{"invalid email", "notanemail", true, "email"},
		{"too long", "verylongemailaddress@verylongdomainname.verylongtld", true, "exceed"},
		{"valid", "user@example.com", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Combine() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errMsg != "" && !strings.Contains(strings.ToLower(err.Error()), tt.errMsg) {
				t.Errorf("Combine() error message = %v, want to contain '%s'", err, tt.errMsg)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	t.Run("validates with multiple validators", func(t *testing.T) {
		err := Validate("test@example.com", Required(), Email())
		if err != nil {
			t.Errorf("Validate() unexpected error = %v", err)
		}
	})

	t.Run("returns first error", func(t *testing.T) {
		err := Validate("", Required(), Email())
		if err == nil {
			t.Error("Validate() expected error for empty required field")
		}

		if !strings.Contains(err.Error(), "required") {
			t.Errorf("Validate() error = %v, want 'required'", err)
		}
	})
}

func TestValidateField(t *testing.T) {
	t.Run("validates field with name", func(t *testing.T) {
		err := ValidateField("email", "test@example.com", Required(), Email())
		if err != nil {
			t.Errorf("ValidateField() unexpected error = %v", err)
		}
	})

	t.Run("includes field name in error", func(t *testing.T) {
		err := ValidateField("email", "", Required())
		if err == nil {
			t.Error("ValidateField() expected error")
		}

		if !strings.Contains(err.Error(), "email") {
			t.Errorf("ValidateField() error = %v, want to contain 'email'", err)
		}
	})
}

func TestValidationError(t *testing.T) {
	t.Run("formats error with field", func(t *testing.T) {
		err := NewValidationError("email", "Invalid email")

		expected := "email: Invalid email"
		if err.Error() != expected {
			t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("formats error without field", func(t *testing.T) {
		err := NewValidationError("", "Invalid input")

		expected := "Invalid input"
		if err.Error() != expected {
			t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
		}
	})
}
