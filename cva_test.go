package forgeui

import (
	"strings"
	"testing"
)

func TestCVA_NewCVA(t *testing.T) {
	cva := NewCVA("base-class-1", "base-class-2")

	classes := cva.Classes(nil)
	if classes != "base-class-1 base-class-2" {
		t.Errorf("expected 'base-class-1 base-class-2', got '%s'", classes)
	}
}

func TestCVA_Variant(t *testing.T) {
	cva := NewCVA("base").
		Variant("size", map[string][]string{
			"sm": {"text-sm", "px-2"},
			"md": {"text-base", "px-4"},
			"lg": {"text-lg", "px-6"},
		}).
		Default("size", "md")

	tests := []struct {
		name     string
		variants map[string]string
		want     string
	}{
		{
			name:     "default variant",
			variants: map[string]string{},
			want:     "base text-base px-4",
		},
		{
			name:     "small variant",
			variants: map[string]string{"size": "sm"},
			want:     "base text-sm px-2",
		},
		{
			name:     "large variant",
			variants: map[string]string{"size": "lg"},
			want:     "base text-lg px-6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cva.Classes(tt.variants)
			if got != tt.want {
				t.Errorf("Classes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCVA_MultipleVariants(t *testing.T) {
	cva := NewCVA("btn").
		Variant("size", map[string][]string{
			"sm": {"h-8", "px-3"},
			"lg": {"h-12", "px-6"},
		}).
		Variant("variant", map[string][]string{
			"primary":   {"bg-blue-500", "text-white"},
			"secondary": {"bg-gray-500", "text-white"},
		}).
		Default("size", "sm").
		Default("variant", "primary")

	tests := []struct {
		name     string
		variants map[string]string
		want     []string // Check that all these are present
	}{
		{
			name:     "defaults",
			variants: map[string]string{},
			want:     []string{"btn", "h-8", "px-3", "bg-blue-500", "text-white"},
		},
		{
			name:     "large secondary",
			variants: map[string]string{"size": "lg", "variant": "secondary"},
			want:     []string{"btn", "h-12", "px-6", "bg-gray-500", "text-white"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cva.Classes(tt.variants)
			for _, class := range tt.want {
				if !strings.Contains(got, class) {
					t.Errorf("Classes() missing %v in %v", class, got)
				}
			}
		})
	}
}

func TestCVA_CompoundVariants(t *testing.T) {
	cva := NewCVA("btn").
		Variant("size", map[string][]string{
			"sm": {"h-8"},
			"lg": {"h-12"},
		}).
		Variant("variant", map[string][]string{
			"primary":   {"bg-blue-500"},
			"secondary": {"bg-gray-500"},
		}).
		Compound(map[string]string{
			"size":    "sm",
			"variant": "primary",
		}, "shadow-sm", "font-medium").
		Default("size", "sm").
		Default("variant", "primary")

	tests := []struct {
		name     string
		variants map[string]string
		want     []string
		notWant  []string
	}{
		{
			name:     "compound matches",
			variants: map[string]string{"size": "sm", "variant": "primary"},
			want:     []string{"btn", "h-8", "bg-blue-500", "shadow-sm", "font-medium"},
		},
		{
			name:     "compound doesn't match",
			variants: map[string]string{"size": "lg", "variant": "primary"},
			want:     []string{"btn", "h-12", "bg-blue-500"},
			notWant:  []string{"shadow-sm", "font-medium"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cva.Classes(tt.variants)
			for _, class := range tt.want {
				if !strings.Contains(got, class) {
					t.Errorf("Classes() missing %v in %v", class, got)
				}
			}
			for _, class := range tt.notWant {
				if strings.Contains(got, class) {
					t.Errorf("Classes() should not contain %v in %v", class, got)
				}
			}
		})
	}
}

func TestCVA_EmptyVariants(t *testing.T) {
	cva := NewCVA("base")

	got := cva.Classes(nil)
	if got != "base" {
		t.Errorf("Classes() = %v, want 'base'", got)
	}

	got = cva.Classes(map[string]string{})
	if got != "base" {
		t.Errorf("Classes() = %v, want 'base'", got)
	}
}

func TestCVA_UnknownVariant(t *testing.T) {
	cva := NewCVA("base").
		Variant("size", map[string][]string{
			"sm": {"text-sm"},
		}).
		Default("size", "sm")

	// Unknown variant value should use default
	got := cva.Classes(map[string]string{"size": "unknown"})
	// Since "unknown" is not in the variant options, no classes from size should be added
	if got != "base" {
		t.Errorf("Classes() = %v, want 'base' for unknown variant", got)
	}
}
