package forgeui

import "testing"

func TestCN(t *testing.T) {
	tests := []struct {
		name    string
		classes []string
		want    string
	}{
		{
			name:    "single class",
			classes: []string{"btn"},
			want:    "btn",
		},
		{
			name:    "multiple classes",
			classes: []string{"btn", "btn-primary", "text-white"},
			want:    "btn btn-primary text-white",
		},
		{
			name:    "with empty strings",
			classes: []string{"btn", "", "btn-primary", ""},
			want:    "btn btn-primary",
		},
		{
			name:    "with whitespace",
			classes: []string{"btn", "  ", "btn-primary", "\t"},
			want:    "btn btn-primary",
		},
		{
			name:    "empty input",
			classes: []string{},
			want:    "",
		},
		{
			name:    "all empty",
			classes: []string{"", "", ""},
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CN(tt.classes...); got != tt.want {
				t.Errorf("CN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIf(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		value     string
		want      string
	}{
		{
			name:      "true condition",
			condition: true,
			value:     "active",
			want:      "active",
		},
		{
			name:      "false condition",
			condition: false,
			value:     "active",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := If(tt.condition, tt.value); got != tt.want {
				t.Errorf("If() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIfElse(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		trueVal   string
		falseVal  string
		want      string
	}{
		{
			name:      "true condition",
			condition: true,
			trueVal:   "active",
			falseVal:  "inactive",
			want:      "active",
		},
		{
			name:      "false condition",
			condition: false,
			trueVal:   "active",
			falseVal:  "inactive",
			want:      "inactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IfElse(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
				t.Errorf("IfElse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapGet(t *testing.T) {
	m := map[string]int{
		"one": 1,
		"two": 2,
	}

	tests := []struct {
		name       string
		key        string
		defaultVal int
		want       int
	}{
		{
			name:       "existing key",
			key:        "one",
			defaultVal: 0,
			want:       1,
		},
		{
			name:       "missing key",
			key:        "three",
			defaultVal: 99,
			want:       99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapGet(m, tt.key, tt.defaultVal); got != tt.want {
				t.Errorf("MapGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
