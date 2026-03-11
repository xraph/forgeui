package bridge

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseParams(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	tests := []struct {
		name       string
		params     json.RawMessage
		targetType reflect.Type
		wantErr    bool
	}{
		{
			name:       "valid params",
			params:     json.RawMessage(`{"name":"John","age":30,"email":"john@example.com"}`),
			targetType: reflect.TypeFor[testStruct](),
			wantErr:    false,
		},
		{
			name:       "empty params",
			params:     json.RawMessage(``),
			targetType: reflect.TypeFor[testStruct](),
			wantErr:    false,
		},
		{
			name:       "invalid JSON",
			params:     json.RawMessage(`{invalid}`),
			targetType: reflect.TypeFor[testStruct](),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseParams(tt.params, tt.targetType)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateParams(t *testing.T) {
	type testStruct struct {
		Required string `json:"required"`
		Optional string `json:"optional,omitempty"`
	}

	tests := []struct {
		name    string
		value   reflect.Value
		wantErr bool
	}{
		{
			name:    "all required fields present",
			value:   reflect.ValueOf(testStruct{Required: "value"}),
			wantErr: false,
		},
		{
			name:    "required field missing",
			value:   reflect.ValueOf(testStruct{}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParams(tt.value, tt.value.Type())
			if (err != nil) != tt.wantErr {
				t.Errorf("validateParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			got := isValidEmail(tt.email)
			if got != tt.want {
				t.Errorf("isValidEmail(%s) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name  string
		value reflect.Value
		want  bool
	}{
		{"empty string", reflect.ValueOf(""), true},
		{"non-empty string", reflect.ValueOf("hello"), false},
		{"zero int", reflect.ValueOf(0), true},
		{"non-zero int", reflect.ValueOf(42), false},
		{"false bool", reflect.ValueOf(false), true},
		{"true bool", reflect.ValueOf(true), false},
		{"empty slice", reflect.ValueOf([]int{}), true},
		{"non-empty slice", reflect.ValueOf([]int{1}), false},
		{"zero float", reflect.ValueOf(0.0), true},
		{"non-zero float", reflect.ValueOf(3.14), false},
		{"zero uint", reflect.ValueOf(uint(0)), true},
		{"non-zero uint", reflect.ValueOf(uint(5)), false},
		{"nil pointer", reflect.ValueOf((*string)(nil)), true},
		{"empty map", reflect.ValueOf(map[string]int{}), true},
		{"non-empty map", reflect.ValueOf(map[string]int{"a": 1}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isZero(tt.value)
			if got != tt.want {
				t.Errorf("isZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- validateField tests ---

func TestValidateField_RequiredNoOmitempty(t *testing.T) {
	type s struct {
		Name string `json:"name"`
	}

	st := reflect.TypeFor[s]()
	field := st.Field(0)
	val := reflect.ValueOf(s{Name: ""})

	err := validateField(field, val.Field(0))
	if err == nil {
		t.Error("expected error for empty required field, got nil")
	}
}

func TestValidateField_RequiredPresent(t *testing.T) {
	type s struct {
		Name string `json:"name"`
	}

	st := reflect.TypeFor[s]()
	field := st.Field(0)
	val := reflect.ValueOf(s{Name: "hello"})

	err := validateField(field, val.Field(0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateField_OmitemptyAllowsZero(t *testing.T) {
	type s struct {
		Count int `json:"count,omitempty"`
	}

	st := reflect.TypeFor[s]()
	field := st.Field(0)
	val := reflect.ValueOf(s{Count: 0})

	err := validateField(field, val.Field(0))
	if err != nil {
		t.Errorf("unexpected error: omitempty field should allow zero value: %v", err)
	}
}

func TestValidateField_ValidateTagRequired(t *testing.T) {
	type s struct {
		Email string `json:"email,omitempty" validate:"required"`
	}

	st := reflect.TypeFor[s]()
	field := st.Field(0)
	val := reflect.ValueOf(s{Email: ""})

	err := validateField(field, val.Field(0))
	// omitempty makes it not auto-required, but validate:"required" should still catch it
	// Note: validateField first checks omitempty (makes isRequired=false),
	// then checks validate tag. validate:"required" should still fail.
	if err == nil {
		t.Error("expected error for validate:required on empty field, got nil")
	}
}

func TestValidateField_ValidateTagEmail(t *testing.T) {
	type s struct {
		Email string `json:"email,omitempty" validate:"email"`
	}

	st := reflect.TypeFor[s]()
	field := st.Field(0)

	// Invalid email
	val := reflect.ValueOf(s{Email: "not-an-email"})
	err := validateField(field, val.Field(0))
	if err == nil {
		t.Error("expected error for invalid email, got nil")
	}

	// Valid email
	val2 := reflect.ValueOf(s{Email: "test@example.com"})
	err2 := validateField(field, val2.Field(0))
	if err2 != nil {
		t.Errorf("unexpected error for valid email: %v", err2)
	}
}

// --- validateStruct tests ---

func TestValidateStruct_MixedFields(t *testing.T) {
	type mixed struct {
		Required string `json:"required"`
		Optional string `json:"optional,omitempty"`
		Also     int    `json:"also"`
	}

	st := reflect.TypeFor[mixed]()

	// Missing both required fields
	val := reflect.ValueOf(mixed{Optional: "present"})
	err := validateStruct(val, st)
	if err == nil {
		t.Error("expected error for missing required field")
	}

	// All required present
	val2 := reflect.ValueOf(mixed{Required: "yes", Also: 1})
	err2 := validateStruct(val2, st)
	if err2 != nil {
		t.Errorf("unexpected error: %v", err2)
	}
}

func TestValidateStruct_UnexportedFieldsSkipped(t *testing.T) {
	type withPrivate struct {
		Public  string `json:"public"`
		private string //nolint:unused
	}

	st := reflect.TypeFor[withPrivate]()
	val := reflect.ValueOf(withPrivate{Public: "yes"})

	err := validateStruct(val, st)
	if err != nil {
		t.Errorf("unexpected error (private fields should be skipped): %v", err)
	}
}

// --- validateParamsLax tests ---

func TestValidateParamsLax_SkipsAutoRequired(t *testing.T) {
	type s struct {
		Name  string `json:"name"`          // auto-required in strict
		Value int    `json:"value"`         // auto-required in strict
		Opt   string `json:"opt,omitempty"` // optional
	}

	st := reflect.TypeFor[s]()
	val := reflect.ValueOf(s{}) // all zero

	// Strict should fail
	strictErr := validateParams(val, st)
	if strictErr == nil {
		t.Error("strict validation should fail for empty required fields")
	}

	// Lax should pass (no validate:"required" tags)
	laxErr := validateParamsLax(val, st)
	if laxErr != nil {
		t.Errorf("lax validation should pass without validate tags: %v", laxErr)
	}
}

func TestValidateParamsLax_EnforcesExplicitRequired(t *testing.T) {
	type s struct {
		Name  string `json:"name" validate:"required"`
		Value int    `json:"value"`
	}

	st := reflect.TypeFor[s]()
	val := reflect.ValueOf(s{Value: 42}) // Name is empty

	err := validateParamsLax(val, st)
	if err == nil {
		t.Error("lax validation should fail for validate:required on empty field")
	}
}

func TestValidateParamsLax_PassesWithRequiredPresent(t *testing.T) {
	type s struct {
		Name  string `json:"name" validate:"required"`
		Value int    `json:"value"`
	}

	st := reflect.TypeFor[s]()
	val := reflect.ValueOf(s{Name: "present"})

	err := validateParamsLax(val, st)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateParamsLax_TypeMismatch(t *testing.T) {
	type s1 struct {
		Name string `json:"name"`
	}
	type s2 struct {
		Name string `json:"name"`
	}

	val := reflect.ValueOf(s1{Name: "test"})
	err := validateParamsLax(val, reflect.TypeFor[s2]())
	if err == nil {
		t.Error("expected type mismatch error")
	}
}

func TestValidateParamsLax_EmailValidation(t *testing.T) {
	type s struct {
		Email string `json:"email" validate:"email"`
	}

	st := reflect.TypeFor[s]()

	// Invalid email
	val := reflect.ValueOf(s{Email: "bad"})
	err := validateParamsLax(val, st)
	if err == nil {
		t.Error("lax validation should still enforce email validation")
	}

	// Valid email
	val2 := reflect.ValueOf(s{Email: "user@example.com"})
	err2 := validateParamsLax(val2, st)
	if err2 != nil {
		t.Errorf("unexpected error for valid email: %v", err2)
	}
}

// --- validateStructLax tests ---

func TestValidateStructLax_SkipsFieldsWithoutValidateTag(t *testing.T) {
	type s struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	st := reflect.TypeFor[s]()
	val := reflect.ValueOf(s{}) // all zero

	err := validateStructLax(val, st)
	if err != nil {
		t.Errorf("unexpected error (no validate tags): %v", err)
	}
}

func TestValidateStructLax_UnexportedSkipped(t *testing.T) {
	type s struct {
		private string `validate:"required"` //nolint:unused
		Public  string `json:"public"`
	}

	st := reflect.TypeFor[s]()
	val := reflect.ValueOf(s{Public: "yes"})

	err := validateStructLax(val, st)
	if err != nil {
		t.Errorf("unexpected error (private should be skipped): %v", err)
	}
}

// --- validateWithTag tests ---

func TestValidateWithTag_Required(t *testing.T) {
	// Zero value should fail
	err := validateWithTag("Field", reflect.ValueOf(""), "required")
	if err == nil {
		t.Error("expected error for empty required field")
	}

	// Non-zero should pass
	err2 := validateWithTag("Field", reflect.ValueOf("value"), "required")
	if err2 != nil {
		t.Errorf("unexpected error: %v", err2)
	}
}

func TestValidateWithTag_Email(t *testing.T) {
	// Invalid
	err := validateWithTag("Email", reflect.ValueOf("invalid"), "email")
	if err == nil {
		t.Error("expected error for invalid email")
	}

	// Valid
	err2 := validateWithTag("Email", reflect.ValueOf("a@b.com"), "email")
	if err2 != nil {
		t.Errorf("unexpected error: %v", err2)
	}

	// Non-string (should not error for email tag on non-string)
	err3 := validateWithTag("Count", reflect.ValueOf(42), "email")
	if err3 != nil {
		t.Errorf("unexpected error for email tag on non-string: %v", err3)
	}
}

func TestValidateWithTag_Unknown(t *testing.T) {
	// Unknown tag should not error
	err := validateWithTag("Field", reflect.ValueOf(""), "custom_unknown")
	if err != nil {
		t.Errorf("unexpected error for unknown validation tag: %v", err)
	}
}

// --- validateParams type mismatch ---

func TestValidateParams_TypeMismatch(t *testing.T) {
	type s1 struct{ Name string }
	type s2 struct{ Name string }

	val := reflect.ValueOf(s1{Name: "test"})
	err := validateParams(val, reflect.TypeFor[s2]())
	if err == nil {
		t.Error("expected type mismatch error")
	}
}

func TestValidateParams_NonStruct(t *testing.T) {
	val := reflect.ValueOf("hello")
	err := validateParams(val, reflect.TypeFor[string]())
	if err != nil {
		t.Errorf("unexpected error for non-struct type: %v", err)
	}
}
