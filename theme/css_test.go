package theme_test

import (
	"strings"
	"testing"

	"github.com/xraph/forgeui/theme"
)

func TestGenerateCSS(t *testing.T) {
	light := theme.DefaultLight()
	dark := theme.DefaultDark()

	css := theme.GenerateCSS(light, dark)

	// Check for :root selector
	if !strings.Contains(css, ":root {") {
		t.Error("CSS should contain :root selector")
	}

	// Check for .dark selector
	if !strings.Contains(css, ".dark {") {
		t.Error("CSS should contain .dark selector")
	}

	// Check for color variables
	expectedVars := []string{
		"--background:",
		"--foreground:",
		"--primary:",
		"--secondary:",
		"--muted:",
		"--accent:",
		"--destructive:",
		"--success:",
		"--border:",
		"--input:",
		"--ring:",
		"--sidebar:",
		"--sidebar-foreground:",
	}

	for _, v := range expectedVars {
		if !strings.Contains(css, v) {
			t.Errorf("CSS should contain %s variable", v)
		}
	}
}

func TestGenerateLightCSS(t *testing.T) {
	th := theme.DefaultLight()

	css := theme.GenerateLightCSS(th)

	if !strings.Contains(css, ":root {") {
		t.Error("CSS should contain :root selector")
	}

	if strings.Contains(css, ".dark {") {
		t.Error("Light-only CSS should not contain .dark selector")
	}

	// Check for radius variables
	if !strings.Contains(css, "--radius:") {
		t.Error("CSS should contain --radius variable")
	}

	// Check for spacing variables
	if !strings.Contains(css, "--spacing:") {
		t.Error("CSS should contain --spacing variable")
	}
}

func TestGenerateTailwindConfig(t *testing.T) {
	config := theme.GenerateTailwindConfig()

	if config == "" {
		t.Error("GenerateTailwindConfig should return non-empty string")
	}

	// Check for Tailwind config structure with OKLCH
	expectedParts := []string{
		"darkMode:",
		"theme:",
		"extend:",
		"colors:",
		"oklch(var(--background))",
		"oklch(var(--primary))",
		"oklch(var(--success))",
		"oklch(var(--sidebar))",
	}

	for _, part := range expectedParts {
		if !strings.Contains(config, part) {
			t.Errorf("Tailwind config should contain %s", part)
		}
	}
}

func TestCSSVariableFormat(t *testing.T) {
	th := theme.DefaultLight()
	css := theme.GenerateLightCSS(th)

	// Verify OKLCH format (no "oklch()" wrapper in variables)
	if strings.Contains(css, "oklch(1 0 0)") {
		t.Error("CSS variables should not contain oklch() wrapper")
	}

	// Should just be the values in OKLCH format
	if !strings.Contains(css, "1 0 0") {
		t.Error("CSS should contain OKLCH values without wrapper")
	}

	// Check for OKLCH-specific patterns (lightness 0-1, not 0-100%)
	if strings.Contains(css, "100%") && !strings.Contains(css, "100% 97.3%") {
		t.Error("CSS should use OKLCH format (0-1 range) not HSL (0-100%)")
	}
}

func TestGenerateInputCSS(t *testing.T) {
	light := theme.DefaultLight()
	dark := theme.DefaultDark()

	css := theme.GenerateInputCSS(light, dark)

	// Check Tailwind v4 import
	if !strings.Contains(css, `@import "tailwindcss"`) {
		t.Error("Input CSS should contain @import \"tailwindcss\"")
	}

	// Check dark mode custom variant
	if !strings.Contains(css, "@custom-variant dark") {
		t.Error("Input CSS should contain @custom-variant dark")
	}

	// Check @theme inline block
	if !strings.Contains(css, "@theme inline {") {
		t.Error("Input CSS should contain @theme inline block")
	}

	// Check color mappings in @theme inline
	themeMappings := []string{
		"--color-background: var(--background)",
		"--color-primary: var(--primary)",
		"--color-sidebar: var(--sidebar)",
		"--color-destructive: var(--destructive)",
	}
	for _, m := range themeMappings {
		if !strings.Contains(css, m) {
			t.Errorf("Input CSS @theme inline should contain %s", m)
		}
	}

	// Check radius mappings in @theme inline
	if !strings.Contains(css, "--radius-lg: var(--radius)") {
		t.Error("Input CSS should contain radius mappings")
	}

	// Check :root with oklch() wrapped values
	if !strings.Contains(css, ":root {") {
		t.Error("Input CSS should contain :root selector")
	}
	if !strings.Contains(css, "--background: oklch(") {
		t.Error("Input CSS :root should contain oklch() wrapped values")
	}

	// Check .dark with oklch() wrapped values
	if !strings.Contains(css, ".dark {") {
		t.Error("Input CSS should contain .dark selector")
	}

	// Check @layer base
	if !strings.Contains(css, "@layer base {") {
		t.Error("Input CSS should contain @layer base")
	}
	if !strings.Contains(css, "@apply border-border") {
		t.Error("Input CSS @layer base should apply border-border")
	}
	if !strings.Contains(css, "@apply bg-background text-foreground") {
		t.Error("Input CSS @layer base should apply bg-background text-foreground")
	}

	// Check --radius is set in :root
	if !strings.Contains(css, "--radius:") {
		t.Error("Input CSS :root should set --radius")
	}

	// Check all sidebar color mappings
	sidebarMappings := []string{
		"--color-sidebar-foreground: var(--sidebar-foreground)",
		"--color-sidebar-primary: var(--sidebar-primary)",
		"--color-sidebar-accent: var(--sidebar-accent)",
		"--color-sidebar-border: var(--sidebar-border)",
	}
	for _, m := range sidebarMappings {
		if !strings.Contains(css, m) {
			t.Errorf("Input CSS @theme inline should contain %s", m)
		}
	}
}

func TestGenerateInputCSSMatchesTestappStructure(t *testing.T) {
	light := theme.DefaultLight()
	dark := theme.DefaultDark()

	css := theme.GenerateInputCSS(light, dark)

	// Verify the sections appear in the correct order
	importIdx := strings.Index(css, `@import "tailwindcss"`)
	variantIdx := strings.Index(css, "@custom-variant dark")
	themeIdx := strings.Index(css, "@theme inline")
	rootIdx := strings.Index(css, ":root {")
	darkIdx := strings.Index(css, ".dark {")
	layerIdx := strings.Index(css, "@layer base")

	if importIdx >= variantIdx {
		t.Error("@import should come before @custom-variant")
	}
	if variantIdx >= themeIdx {
		t.Error("@custom-variant should come before @theme inline")
	}
	if themeIdx >= rootIdx {
		t.Error("@theme inline should come before :root")
	}
	if rootIdx >= darkIdx {
		t.Error(":root should come before .dark")
	}
	if darkIdx >= layerIdx {
		t.Error(".dark should come before @layer base")
	}
}

func TestChartVariables(t *testing.T) {
	th := theme.DefaultLight()
	css := theme.GenerateLightCSS(th)

	// Test all 12 chart colors
	chartVars := []string{
		"--chart-1:",
		"--chart-2:",
		"--chart-3:",
		"--chart-4:",
		"--chart-5:",
		"--chart-6:",
		"--chart-7:",
		"--chart-8:",
		"--chart-9:",
		"--chart-10:",
		"--chart-11:",
		"--chart-12:",
	}

	for _, v := range chartVars {
		if !strings.Contains(css, v) {
			t.Errorf("CSS should contain %s variable", v)
		}
	}
}
