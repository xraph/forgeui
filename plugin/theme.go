package plugin

import (
	"github.com/xraph/forgeui/theme"
)

// ThemePlugin provides custom themes for ForgeUI.
//
// A ThemePlugin can provide:
//   - Multiple theme presets (light/dark variants)
//   - Additional CSS for custom properties or utilities
//   - Custom fonts to load
//
// Example:
//
//	type CorporateThemePlugin struct {
//	    *PluginBase
//	}
//
//	func (p *CorporateThemePlugin) Themes() map[string]theme.Theme {
//	    return map[string]theme.Theme{
//	        "corporate-light": corporateLightTheme,
//	        "corporate-dark":  corporateDarkTheme,
//	    }
//	}
//
//	func (p *CorporateThemePlugin) Fonts() []theme.Font {
//	    return []theme.Font{
//	        {Family: "Inter", Weights: []int{400, 600, 700}},
//	    }
//	}
type ThemePlugin interface {
	Plugin

	// Themes returns a map of theme names to Theme configurations.
	// Theme names should be unique and descriptive (e.g., "corporate-light").
	Themes() map[string]theme.Theme

	// DefaultTheme returns the name of the default theme to use.
	// This should match one of the keys from Themes().
	// Return empty string to let the system choose.
	DefaultTheme() string

	// CSS returns additional CSS to inject.
	// This can include custom properties, utility classes, or font-face rules.
	CSS() string

	// Fonts returns fonts to load.
	// These will be loaded automatically via Google Fonts or custom URLs.
	Fonts() []theme.Font
}

// ThemePluginBase provides default implementations for ThemePlugin.
// Embed this to implement only the methods you need.
type ThemePluginBase struct {
	*PluginBase
	themes       map[string]theme.Theme
	defaultTheme string
	fonts        []theme.Font
}

// NewThemePluginBase creates a new ThemePluginBase.
func NewThemePluginBase(
	info PluginInfo,
	themes map[string]theme.Theme,
	defaultTheme string,
) *ThemePluginBase {
	return &ThemePluginBase{
		PluginBase:   NewPluginBase(info),
		themes:       themes,
		defaultTheme: defaultTheme,
	}
}

// NewThemePluginBaseWithFonts creates a ThemePluginBase with fonts.
func NewThemePluginBaseWithFonts(
	info PluginInfo,
	themes map[string]theme.Theme,
	defaultTheme string,
	fonts []theme.Font,
) *ThemePluginBase {
	return &ThemePluginBase{
		PluginBase:   NewPluginBase(info),
		themes:       themes,
		defaultTheme: defaultTheme,
		fonts:        fonts,
	}
}

// Themes returns the theme configurations.
func (t *ThemePluginBase) Themes() map[string]theme.Theme {
	return t.themes
}

// DefaultTheme returns the default theme name.
func (t *ThemePluginBase) DefaultTheme() string {
	return t.defaultTheme
}

// CSS returns empty string by default.
// Override this method to provide custom CSS.
func (t *ThemePluginBase) CSS() string {
	return ""
}

// Fonts returns the font configurations.
func (t *ThemePluginBase) Fonts() []theme.Font {
	return t.fonts
}

// AddTheme adds a theme to the plugin.
// This can be called during plugin initialization.
func (t *ThemePluginBase) AddTheme(name string, th theme.Theme) {
	if t.themes == nil {
		t.themes = make(map[string]theme.Theme)
	}
	t.themes[name] = th
}

// AddFont adds a font to the plugin.
func (t *ThemePluginBase) AddFont(font theme.Font) {
	t.fonts = append(t.fonts, font)
}

