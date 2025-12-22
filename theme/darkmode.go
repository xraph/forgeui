package theme

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// DarkModeScript returns an inline script that prevents flash of unstyled content (FOUC)
// by reading the theme preference from localStorage before the page renders.
//
// This script should be placed at the beginning of the <body> tag to execute immediately.
// It will:
//  1. Check localStorage for a saved theme preference
//  2. Fall back to system preference (prefers-color-scheme)
//  3. Apply the 'dark' class to <html> element if needed
func DarkModeScript() g.Node {
	script := `(function() {
  try {
    const stored = localStorage.getItem('theme');
    const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const theme = stored === 'dark' || (!stored && system) ? 'dark' : 'light';
    document.documentElement.classList.toggle('dark', theme === 'dark');
    document.documentElement.setAttribute('data-theme', theme);
  } catch (e) {
    // localStorage might be unavailable (private browsing)
    const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
    document.documentElement.classList.toggle('dark', system);
    document.documentElement.setAttribute('data-theme', system ? 'dark' : 'light');
  }
})();`

	return html.Script(
		g.Attr("data-theme-script", ""),
		g.Raw(script),
	)
}

// DarkModeScriptWithDefault returns a dark mode script with a specified default theme.
// If no preference is stored and system preference is unavailable, use the default.
func DarkModeScriptWithDefault(defaultTheme string) g.Node {
	script := `(function() {
  try {
    const stored = localStorage.getItem('theme');
    const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const defaultTheme = '` + defaultTheme + `';
    const theme = stored || (system ? 'dark' : 'light') || defaultTheme;
    document.documentElement.classList.toggle('dark', theme === 'dark');
    document.documentElement.setAttribute('data-theme', theme);
  } catch (e) {
    const defaultTheme = '` + defaultTheme + `';
    document.documentElement.classList.toggle('dark', defaultTheme === 'dark');
    document.documentElement.setAttribute('data-theme', defaultTheme);
  }
})();`

	return html.Script(
		g.Attr("data-theme-script", ""),
		g.Raw(script),
	)
}

// ThemeScript returns a complete theme management script with Alpine.js integration.
// This provides functions to get, set, and sync theme across tabs.
//
// The script exposes:
//  - getTheme(): Returns current theme ('light' or 'dark')
//  - setTheme(theme): Sets theme and persists to localStorage
//  - toggleTheme(): Toggles between light and dark
func ThemeScript() g.Node {
	script := `(function() {
  // Initialize theme on load
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

  // Apply theme to document
  function applyTheme(theme) {
    document.documentElement.classList.toggle('dark', theme === 'dark');
    document.documentElement.setAttribute('data-theme', theme);
  }

  // Get current theme
  function getTheme() {
    return document.documentElement.classList.contains('dark') ? 'dark' : 'light';
  }

  // Set theme
  function setTheme(theme) {
    try {
      localStorage.setItem('theme', theme);
    } catch (e) {
      // localStorage unavailable
    }
    applyTheme(theme);
    
    // Dispatch event for other components to react
    window.dispatchEvent(new CustomEvent('theme-change', { detail: { theme } }));
  }

  // Toggle theme
  function toggleTheme() {
    const current = getTheme();
    const next = current === 'dark' ? 'light' : 'dark';
    setTheme(next);
    return next;
  }

  // Sync theme across tabs
  window.addEventListener('storage', function(e) {
    if (e.key === 'theme' && e.newValue) {
      applyTheme(e.newValue);
    }
  });

  // Listen for system theme changes
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
    try {
      const stored = localStorage.getItem('theme');
      if (!stored) {
        // Only auto-switch if user hasn't set a preference
        setTheme(e.matches ? 'dark' : 'light');
      }
    } catch (err) {
      // localStorage unavailable
    }
  });

  // Expose functions globally
  window.forgeui = window.forgeui || {};
  window.forgeui.theme = {
    get: getTheme,
    set: setTheme,
    toggle: toggleTheme,
    init: initTheme,
  };

  // Initialize immediately
  initTheme();
})();`

	return html.Script(
		g.Attr("data-theme-manager", ""),
		g.Raw(script),
	)
}

// ThemeStyleCloak returns CSS to hide elements with x-cloak attribute.
// This prevents flashing of theme-dependent elements during initialization.
func ThemeStyleCloak() g.Node {
	css := `[x-cloak] { display: none !important; }
[data-theme-loading] { opacity: 0; }
[data-theme-loading].theme-ready { opacity: 1; transition: opacity 0.2s ease; }`

	return g.El("style",
		g.Attr("data-theme-cloak", ""),
		g.Raw(css),
	)
}

