package theme

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// DarkModeScript returns an inline script that prevents flash of unstyled content (FOUC)
// by reading the theme preference from localStorage before the page renders.
//
// This script should be placed at the beginning of the <body> tag to execute immediately.
// It will:
//  1. Check localStorage for a saved theme preference
//  2. Fall back to system preference (prefers-color-scheme)
//  3. Apply the 'dark' class to <html> element if needed
func DarkModeScript() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script data-theme-script>(function() {
  try {
    const stored = localStorage.getItem('theme');
    const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const theme = stored === 'dark' || (!stored && system) ? 'dark' : 'light';
    document.documentElement.classList.toggle('dark', theme === 'dark');
    document.documentElement.setAttribute('data-theme', theme);
  } catch (e) {
    const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
    document.documentElement.classList.toggle('dark', system);
    document.documentElement.setAttribute('data-theme', system ? 'dark' : 'light');
  }
})();</script>`)
		return err
	})
}

// DarkModeScriptWithDefault returns a dark mode script with a specified default theme.
// If no preference is stored and system preference is unavailable, use the default.
func DarkModeScriptWithDefault(defaultTheme string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script data-theme-script>(function() {
  try {
    const stored = localStorage.getItem('theme');
    const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const defaultTheme = '`+defaultTheme+`';
    const theme = stored || (system ? 'dark' : 'light') || defaultTheme;
    document.documentElement.classList.toggle('dark', theme === 'dark');
    document.documentElement.setAttribute('data-theme', theme);
  } catch (e) {
    const defaultTheme = '`+defaultTheme+`';
    document.documentElement.classList.toggle('dark', defaultTheme === 'dark');
    document.documentElement.setAttribute('data-theme', defaultTheme);
  }
})();</script>`)
		return err
	})
}

// ThemeScript returns a complete theme management script with Alpine.js integration.
// This provides functions to get, set, and sync theme across tabs.
//
// The script exposes:
//   - getTheme(): Returns current theme ('light' or 'dark')
//   - setTheme(theme): Sets theme and persists to localStorage
//   - toggleTheme(): Toggles between light and dark
func ThemeScript() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script data-theme-manager>(function() {
  function initTheme() {
    try {
      const stored = localStorage.getItem('theme');
      const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
      const theme = stored === 'dark' || (!stored && system) ? 'dark' : 'light';
      applyTheme(theme);
      return theme;
    } catch (e) {
      const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
      const theme = system ? 'dark' : 'light';
      applyTheme(theme);
      return theme;
    }
  }
  function applyTheme(theme) {
    document.documentElement.classList.toggle('dark', theme === 'dark');
    document.documentElement.setAttribute('data-theme', theme);
  }
  function getTheme() {
    return document.documentElement.classList.contains('dark') ? 'dark' : 'light';
  }
  function setTheme(theme) {
    try { localStorage.setItem('theme', theme); } catch (e) {}
    applyTheme(theme);
    window.dispatchEvent(new CustomEvent('theme-change', { detail: { theme } }));
  }
  function toggleTheme() {
    const current = getTheme();
    const next = current === 'dark' ? 'light' : 'dark';
    setTheme(next);
    return next;
  }
  window.addEventListener('storage', function(e) {
    if (e.key === 'theme' && e.newValue) { applyTheme(e.newValue); }
  });
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
    try {
      const stored = localStorage.getItem('theme');
      if (!stored) { setTheme(e.matches ? 'dark' : 'light'); }
    } catch (err) {}
  });
  window.forgeui = window.forgeui || {};
  window.forgeui.theme = { get: getTheme, set: setTheme, toggle: toggleTheme, init: initTheme };
  initTheme();
})();</script>`)
		return err
	})
}

// ThemeStyleCloak returns CSS to hide elements with x-cloak attribute.
// This prevents flashing of theme-dependent elements during initialization.
func ThemeStyleCloak() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<style data-theme-cloak>[x-cloak] { display: none !important; }
[data-theme-loading] { opacity: 0; }
[data-theme-loading].theme-ready { opacity: 1; transition: opacity 0.2s ease; }</style>`)
		return err
	})
}
