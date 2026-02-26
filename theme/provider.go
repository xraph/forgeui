package theme

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

// ProviderProps defines configuration for the theme provider.
type ProviderProps struct {
	LightTheme       Theme  // Light theme (required)
	DarkTheme        Theme  // Dark theme (required)
	DefaultTheme     string // Default theme if no preference ("light" or "dark")
	EnableSystemSync bool   // Sync with system preference
	StorageKey       string // localStorage key for theme preference
}

// ProviderOption is a functional option for configuring the Provider.
type ProviderOption func(*ProviderProps)

// WithDefaultTheme sets the default theme when no preference is stored.
func WithDefaultTheme(theme string) ProviderOption {
	return func(p *ProviderProps) {
		p.DefaultTheme = theme
	}
}

// WithStorageKey sets a custom localStorage key for theme persistence.
func WithStorageKey(key string) ProviderOption {
	return func(p *ProviderProps) {
		p.StorageKey = key
	}
}

// DisableSystemSync disables automatic syncing with system theme preference.
func DisableSystemSync() ProviderOption {
	return func(p *ProviderProps) {
		p.EnableSystemSync = false
	}
}

// Provider wraps the application with theme support.
// It injects CSS custom properties and dark mode script into the page.
func Provider(children ...templ.Component) templ.Component {
	return ProviderWithThemes(DefaultLight(), DefaultDark(), children...)
}

// ProviderWithThemes wraps the application with custom light and dark themes.
func ProviderWithThemes(lightTheme, darkTheme Theme, children ...templ.Component) templ.Component {
	css := GenerateCSS(lightTheme, darkTheme)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<style data-forgeui-theme>%s</style>`, css); err != nil {
			return err
		}
		if err := DarkModeScript().Render(ctx, w); err != nil {
			return err
		}
		if err := ThemeStyleCloak().Render(ctx, w); err != nil {
			return err
		}
		for _, child := range children {
			if child != nil {
				if err := child.Render(ctx, w); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ProviderWithOptions provides fine-grained control over theme configuration.
func ProviderWithOptions(children []templ.Component, opts ...ProviderOption) templ.Component {
	props := &ProviderProps{
		LightTheme:       DefaultLight(),
		DarkTheme:        DefaultDark(),
		DefaultTheme:     "light",
		EnableSystemSync: true,
		StorageKey:       "theme",
	}

	for _, opt := range opts {
		opt(props)
	}

	css := GenerateCSS(props.LightTheme, props.DarkTheme)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<style data-forgeui-theme>%s</style>`, css); err != nil {
			return err
		}
		if err := DarkModeScriptWithDefault(props.DefaultTheme).Render(ctx, w); err != nil {
			return err
		}
		if err := ThemeStyleCloak().Render(ctx, w); err != nil {
			return err
		}
		for _, child := range children {
			if child != nil {
				if err := child.Render(ctx, w); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// HTMLWrapper wraps the entire HTML document with proper theme attributes.
func HTMLWrapper(lightTheme, darkTheme Theme, children ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<html lang="en" class="">`); err != nil {
			return err
		}
		for _, child := range children {
			if child != nil {
				if err := child.Render(ctx, w); err != nil {
					return err
				}
			}
		}
		_, err := io.WriteString(w, `</html>`)
		return err
	})
}

// HeadContent returns common <head> content for theme support.
func HeadContent(lightTheme, darkTheme Theme) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<meta charset="utf-8">`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<meta name="viewport" content="width=device-width, initial-scale=1">`); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, `<meta name="theme-color" content="hsl(%s)" media="(prefers-color-scheme: light)">`, lightTheme.Colors.Background); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, `<meta name="theme-color" content="hsl(%s)" media="(prefers-color-scheme: dark)">`, darkTheme.Colors.Background); err != nil {
			return err
		}
		_, err := io.WriteString(w, `<meta name="color-scheme" content="light dark">`)
		return err
	})
}

// StyleTag returns just the theme CSS as a style tag.
func StyleTag(lightTheme, darkTheme Theme) templ.Component {
	css := GenerateCSS(lightTheme, darkTheme)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<style data-forgeui-theme>%s</style>`, css)
		return err
	})
}
