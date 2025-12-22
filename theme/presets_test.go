package theme_test

import (
	"testing"

	"github.com/xraph/forgeui/theme"
)

func TestDefaultLight(t *testing.T) {
	th := theme.DefaultLight()

	if th.Colors.Background != "1 0 0" {
		t.Errorf("DefaultLight background should be white in OKLCH, got %s", th.Colors.Background)
	}

	if th.Colors.Primary == "" {
		t.Error("DefaultLight should have a primary color")
	}

	if th.Radius.LG == "" {
		t.Error("DefaultLight should have radius tokens")
	}

	if th.Colors.Success == "" {
		t.Error("DefaultLight should have Success color")
	}
}

func TestDefaultDark(t *testing.T) {
	th := theme.DefaultDark()

	if th.Colors.Background != "0.145 0 0" {
		t.Errorf("DefaultDark background should be dark in OKLCH, got %s", th.Colors.Background)
	}

	if th.Colors.Primary == "" {
		t.Error("DefaultDark should have a primary color")
	}

	if th.Colors.Success == "" {
		t.Error("DefaultDark should have Success color")
	}
}

func TestColorPresets(t *testing.T) {
	presets := []struct {
		name  string
		theme theme.Theme
	}{
		{"RoseLight", theme.RoseLight()},
		{"RoseDark", theme.RoseDark()},
		{"BlueLight", theme.BlueLight()},
		{"BlueDark", theme.BlueDark()},
		{"GreenLight", theme.GreenLight()},
		{"GreenDark", theme.GreenDark()},
		{"OrangeLight", theme.OrangeLight()},
		{"OrangeDark", theme.OrangeDark()},
		{"NeutralLight", theme.NeutralLight()},
		{"NeutralDark", theme.NeutralDark()},
	}

	for _, p := range presets {
		t.Run(p.name, func(t *testing.T) {
			if p.theme.Colors.Primary == "" {
				t.Errorf("%s should have a primary color", p.name)
			}

			if p.theme.Colors.Background == "" {
				t.Errorf("%s should have a background color", p.name)
			}

			if p.theme.Colors.Foreground == "" {
				t.Errorf("%s should have a foreground color", p.name)
			}
		})
	}
}

func TestChartColors(t *testing.T) {
	th := theme.DefaultLight()

	if th.Colors.Chart1 == "" {
		t.Error("DefaultLight should have Chart1 color")
	}

	if th.Colors.Chart2 == "" {
		t.Error("DefaultLight should have Chart2 color")
	}

	if th.Colors.Chart5 == "" {
		t.Error("DefaultLight should have Chart5 color")
	}
	// Test extended chart colors
	if th.Colors.Chart6 == "" {
		t.Error("DefaultLight should have Chart6 color")
	}

	if th.Colors.Chart12 == "" {
		t.Error("DefaultLight should have Chart12 color")
	}
}

func TestSidebarColors(t *testing.T) {
	th := theme.DefaultLight()

	sidebarColors := []struct {
		name  string
		value string
	}{
		{"Sidebar", th.Colors.Sidebar},
		{"SidebarForeground", th.Colors.SidebarForeground},
		{"SidebarPrimary", th.Colors.SidebarPrimary},
		{"SidebarPrimaryForeground", th.Colors.SidebarPrimaryForeground},
		{"SidebarAccent", th.Colors.SidebarAccent},
		{"SidebarAccentForeground", th.Colors.SidebarAccentForeground},
		{"SidebarBorder", th.Colors.SidebarBorder},
		{"SidebarRing", th.Colors.SidebarRing},
	}

	for _, sc := range sidebarColors {
		if sc.value == "" {
			t.Errorf("DefaultLight should have %s color", sc.name)
		}
	}
}
