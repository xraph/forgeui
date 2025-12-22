package theme_test

import (
	"testing"

	"github.com/xraph/forgeui/theme"
)

func TestNew(t *testing.T) {
	th := theme.New()

	if th.Radius.MD == "" {
		t.Error("New() should set default radius")
	}

	if th.Spacing.MD == "" {
		t.Error("New() should set default spacing")
	}

	if th.FontSize.Base == "" {
		t.Error("New() should set default font size")
	}

	if th.Shadow.MD == "" {
		t.Error("New() should set default shadow")
	}
}

func TestColorTokens(t *testing.T) {
	tests := []struct {
		name   string
		theme  theme.Theme
		field  string
		expect string
	}{
		{
			name:   "DefaultLight background",
			theme:  theme.DefaultLight(),
			field:  "Background",
			expect: "1 0 0",
		},
		{
			name:   "DefaultDark background",
			theme:  theme.DefaultDark(),
			field:  "Background",
			expect: "0.145 0 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string

			if tt.field == "Background" {
				got = tt.theme.Colors.Background
			}

			if got != tt.expect {
				t.Errorf("expected %q, got %q", tt.expect, got)
			}
		})
	}
}

func TestNewColorTokens(t *testing.T) {
	th := theme.DefaultLight()

	// Test new Success token
	if th.Colors.Success == "" {
		t.Error("DefaultLight should have Success color")
	}

	// Test new Sidebar tokens
	if th.Colors.Sidebar == "" {
		t.Error("DefaultLight should have Sidebar color")
	}

	if th.Colors.SidebarForeground == "" {
		t.Error("DefaultLight should have SidebarForeground color")
	}

	if th.Colors.SidebarPrimary == "" {
		t.Error("DefaultLight should have SidebarPrimary color")
	}

	// Test extended Chart tokens
	if th.Colors.Chart6 == "" {
		t.Error("DefaultLight should have Chart6 color")
	}

	if th.Colors.Chart12 == "" {
		t.Error("DefaultLight should have Chart12 color")
	}
}

func TestRadiusTokens(t *testing.T) {
	th := theme.New()

	if th.Radius.SM != "calc(0.5rem - 2px)" {
		t.Errorf("expected SM radius to be calc(0.5rem - 2px), got %s", th.Radius.SM)
	}

	if th.Radius.MD != "0.5rem" {
		t.Errorf("expected MD radius to be 0.5rem, got %s", th.Radius.MD)
	}

	if th.Radius.LG != "calc(0.5rem + 2px)" {
		t.Errorf("expected LG radius to be calc(0.5rem + 2px) (base), got %s", th.Radius.LG)
	}

	if th.Radius.XL != "0.75rem" {
		t.Errorf("expected XL radius to be 0.75rem, got %s", th.Radius.XL)
	}

	if th.Radius.Full != "9999px" {
		t.Errorf("expected Full radius to be 9999px, got %s", th.Radius.Full)
	}
}

func TestSpacingTokens(t *testing.T) {
	th := theme.New()

	if th.Spacing.XS != "0.25rem" {
		t.Errorf("expected XS spacing to be 0.25rem, got %s", th.Spacing.XS)
	}

	if th.Spacing.MD != "1rem" {
		t.Errorf("expected MD spacing to be 1rem, got %s", th.Spacing.MD)
	}

	if th.Spacing.XXL != "3rem" {
		t.Errorf("expected XXL spacing to be 3rem, got %s", th.Spacing.XXL)
	}
}

func TestFontSizeTokens(t *testing.T) {
	th := theme.New()

	if th.FontSize.Base != "1rem" {
		t.Errorf("expected Base font size to be 1rem, got %s", th.FontSize.Base)
	}

	if th.FontSize.SM != "0.875rem" {
		t.Errorf("expected SM font size to be 0.875rem, got %s", th.FontSize.SM)
	}
}

func TestShadowTokens(t *testing.T) {
	th := theme.New()

	if th.Shadow.SM == "" {
		t.Error("Shadow.SM should not be empty")
	}

	if th.Shadow.MD == "" {
		t.Error("Shadow.MD should not be empty")
	}
}
