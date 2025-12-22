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
