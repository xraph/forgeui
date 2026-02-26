package theme

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// TailwindConfigScript returns an inline script that configures Tailwind CSS
// to use OKLCH color format with CSS variables.
// This is needed when using the Tailwind Play CDN, as it doesn't natively
// understand that CSS variables contain OKLCH values that need to be wrapped
// with the oklch() function.
//
// Usage:
//
//	html.Head(
//	    html.Script(html.Src("https://cdn.tailwindcss.com")),
//	    theme.TailwindConfigScript(), // Add after Tailwind CDN
//	    theme.StyleTag(lightTheme, darkTheme),
//	)
func TailwindConfigScript() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script>tailwind.config = {
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        background: 'oklch(var(--background) / <alpha-value>)',
        foreground: 'oklch(var(--foreground) / <alpha-value>)',
        card: 'oklch(var(--card) / <alpha-value>)',
        'card-foreground': 'oklch(var(--card-foreground) / <alpha-value>)',
        popover: 'oklch(var(--popover) / <alpha-value>)',
        'popover-foreground': 'oklch(var(--popover-foreground) / <alpha-value>)',
        primary: 'oklch(var(--primary) / <alpha-value>)',
        'primary-foreground': 'oklch(var(--primary-foreground) / <alpha-value>)',
        secondary: 'oklch(var(--secondary) / <alpha-value>)',
        'secondary-foreground': 'oklch(var(--secondary-foreground) / <alpha-value>)',
        muted: 'oklch(var(--muted) / <alpha-value>)',
        'muted-foreground': 'oklch(var(--muted-foreground) / <alpha-value>)',
        accent: 'oklch(var(--accent) / <alpha-value>)',
        'accent-foreground': 'oklch(var(--accent-foreground) / <alpha-value>)',
        destructive: 'oklch(var(--destructive) / <alpha-value>)',
        'destructive-foreground': 'oklch(var(--destructive-foreground) / <alpha-value>)',
        success: 'oklch(var(--success) / <alpha-value>)',
        border: 'oklch(var(--border) / <alpha-value>)',
        input: 'oklch(var(--input) / <alpha-value>)',
        ring: 'oklch(var(--ring) / <alpha-value>)',
        'chart-1': 'oklch(var(--chart-1) / <alpha-value>)',
        'chart-2': 'oklch(var(--chart-2) / <alpha-value>)',
        'chart-3': 'oklch(var(--chart-3) / <alpha-value>)',
        'chart-4': 'oklch(var(--chart-4) / <alpha-value>)',
        'chart-5': 'oklch(var(--chart-5) / <alpha-value>)',
        'chart-6': 'oklch(var(--chart-6) / <alpha-value>)',
        'chart-7': 'oklch(var(--chart-7) / <alpha-value>)',
        'chart-8': 'oklch(var(--chart-8) / <alpha-value>)',
        'chart-9': 'oklch(var(--chart-9) / <alpha-value>)',
        'chart-10': 'oklch(var(--chart-10) / <alpha-value>)',
        'chart-11': 'oklch(var(--chart-11) / <alpha-value>)',
        'chart-12': 'oklch(var(--chart-12) / <alpha-value>)',
        sidebar: 'oklch(var(--sidebar) / <alpha-value>)',
        'sidebar-foreground': 'oklch(var(--sidebar-foreground) / <alpha-value>)',
        'sidebar-primary': 'oklch(var(--sidebar-primary) / <alpha-value>)',
        'sidebar-primary-foreground': 'oklch(var(--sidebar-primary-foreground) / <alpha-value>)',
        'sidebar-accent': 'oklch(var(--sidebar-accent) / <alpha-value>)',
        'sidebar-accent-foreground': 'oklch(var(--sidebar-accent-foreground) / <alpha-value>)',
        'sidebar-border': 'oklch(var(--sidebar-border) / <alpha-value>)',
        'sidebar-ring': 'oklch(var(--sidebar-ring) / <alpha-value>)'
      }
    }
  }
}</script>`)
		return err
	})
}
