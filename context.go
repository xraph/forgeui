package forgeui

import "context"

// Context keys for ForgeUI
type contextKey string

const (
	themeKey   contextKey = "forgeui_theme"
	configKey  contextKey = "forgeui_config"
	requestKey contextKey = "forgeui_request"
)

// WithTheme adds theme to context
func WithTheme(ctx context.Context, theme any) context.Context {
	return context.WithValue(ctx, themeKey, theme)
}

// ThemeFromContext retrieves theme from context
func ThemeFromContext(ctx context.Context) (any, bool) {
	theme := ctx.Value(themeKey)
	return theme, theme != nil
}

// WithConfig adds config to context
func WithConfig(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, configKey, config)
}

// ConfigFromContext retrieves config from context
func ConfigFromContext(ctx context.Context) (*Config, bool) {
	config, ok := ctx.Value(configKey).(*Config)
	return config, ok
}
