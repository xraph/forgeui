package theme

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
)

// StorageType defines the storage mechanism for theme persistence.
type StorageType string

const (
	// LocalStorage persists theme preference across sessions.
	LocalStorage StorageType = "localStorage"
	// SessionStorage persists theme only for the current session.
	SessionStorage StorageType = "sessionStorage"
	// NoStorage disables persistence (in-memory only).
	NoStorage StorageType = "none"
)

// ToggleProps defines configuration for the theme toggle button.
type ToggleProps struct {
	ShowLabel       bool        // Show text label alongside icon
	LightLabel      string      // Label for light mode
	DarkLabel       string      // Label for dark mode
	SystemLabel     string      // Label for system mode
	Position        string      // Button position (if floating)
	Size            string      // Button size (sm, md, lg)
	CustomClass     string      // Additional CSS classes
	IconOnly        bool        // Show only icon, no label
	StorageKey      string      // Key used for storage (default: "theme")
	StorageType     StorageType // Storage mechanism (localStorage, sessionStorage, none)
	SyncAcrossTabs  bool        // Sync theme changes across browser tabs
	RespectSystem   bool        // Respect system preference when no stored value
	DefaultTheme    string      // Default theme when no preference (light, dark, system)
	AnimateTransition bool      // Enable smooth transition animation
}

// ToggleOption is a functional option for configuring the Toggle.
type ToggleOption func(*ToggleProps)

// WithLabel enables the text label.
func WithLabel(show bool) ToggleOption {
	return func(p *ToggleProps) {
		p.ShowLabel = show
	}
}

// WithLabels sets custom labels for light and dark modes.
func WithLabels(light, dark string) ToggleOption {
	return func(p *ToggleProps) {
		p.LightLabel = light
		p.DarkLabel = dark
		p.ShowLabel = true
	}
}

// WithAllLabels sets custom labels for light, dark, and system modes.
func WithAllLabels(light, dark, system string) ToggleOption {
	return func(p *ToggleProps) {
		p.LightLabel = light
		p.DarkLabel = dark
		p.SystemLabel = system
		p.ShowLabel = true
	}
}

// WithToggleSize sets the button size.
func WithToggleSize(size string) ToggleOption {
	return func(p *ToggleProps) {
		p.Size = size
	}
}

// WithToggleClass adds custom CSS classes.
func WithToggleClass(class string) ToggleOption {
	return func(p *ToggleProps) {
		p.CustomClass = class
	}
}

// IconOnly makes the toggle button icon-only.
func IconOnly() ToggleOption {
	return func(p *ToggleProps) {
		p.IconOnly = true
		p.ShowLabel = false
	}
}

// WithToggleStorageKey sets a custom storage key for theme persistence.
func WithToggleStorageKey(key string) ToggleOption {
	return func(p *ToggleProps) {
		p.StorageKey = key
	}
}

// WithStorageType sets the storage mechanism for persistence.
func WithStorageType(storageType StorageType) ToggleOption {
	return func(p *ToggleProps) {
		p.StorageType = storageType
	}
}

// WithTabSync enables synchronization of theme changes across browser tabs.
func WithTabSync(enabled bool) ToggleOption {
	return func(p *ToggleProps) {
		p.SyncAcrossTabs = enabled
	}
}

// WithSystemPreference enables respecting system color scheme preference.
func WithSystemPreference(enabled bool) ToggleOption {
	return func(p *ToggleProps) {
		p.RespectSystem = enabled
	}
}

// WithToggleDefaultTheme sets the default theme when no preference is stored.
func WithToggleDefaultTheme(theme string) ToggleOption {
	return func(p *ToggleProps) {
		p.DefaultTheme = theme
	}
}

// WithTransition enables smooth theme transition animation.
func WithTransition(enabled bool) ToggleOption {
	return func(p *ToggleProps) {
		p.AnimateTransition = enabled
	}
}

// defaultToggleProps returns default toggle properties.
func defaultToggleProps() ToggleProps {
	return ToggleProps{
		ShowLabel:       false,
		LightLabel:      "Light",
		DarkLabel:       "Dark",
		SystemLabel:     "System",
		Size:            "md",
		IconOnly:        true,
		StorageKey:      "theme",
		StorageType:     LocalStorage,
		SyncAcrossTabs:  true,
		RespectSystem:   true,
		DefaultTheme:    "light",
		AnimateTransition: false,
	}
}

// Toggle creates a theme switcher button that toggles between light and dark modes.
// It uses Alpine.js for interactivity and persists the preference to browser storage.
//
// Features:
//   - Persists theme preference to localStorage (default) or sessionStorage
//   - Syncs theme changes across browser tabs
//   - Respects system color scheme preference
//   - Prevents flash of wrong theme on page load
//
// Usage:
//
//	Toggle()                                    // Icon-only toggle with localStorage
//	Toggle(WithLabel(true))                     // With label
//	Toggle(WithStorageKey("my-app-theme"))      // Custom storage key
//	Toggle(WithStorageType(SessionStorage))     // Use sessionStorage
//	Toggle(WithTabSync(false))                  // Disable cross-tab sync
func Toggle(opts ...ToggleOption) g.Node {
	props := defaultToggleProps()
	for _, opt := range opts {
		opt(&props)
	}

	// Button size classes
	sizeClasses := map[string]string{
		"sm": "h-8 px-2 text-sm",
		"md": "h-9 px-3 text-sm",
		"lg": "h-10 px-4 text-base",
	}

	baseClasses := "inline-flex items-center justify-center gap-2 rounded-md font-medium transition-colors " +
		"bg-transparent border border-input hover:bg-accent hover:text-accent-foreground " +
		"focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring " +
		"disabled:pointer-events-none disabled:opacity-50"

	classes := baseClasses + " " + sizeClasses[props.Size]
	if props.CustomClass != "" {
		classes += " " + props.CustomClass
	}

	// Build the Alpine.js init script with persistence
	initScript := buildToggleInitScript(props)
	clickScript := buildToggleClickScript(props)

	nodes := []g.Node{
		g.Attr("type", "button"),
		html.Class(classes),
		g.Attr("aria-label", "Toggle theme"),
		g.Attr("data-theme-toggle", ""),

		// Alpine.js data for theme state with persistence
		alpine.XData(map[string]any{
			"theme": props.DefaultTheme,
			"init":  initScript,
		}),

		// Click handler to toggle theme with persistence
		alpine.XOn("click", clickScript),

		// Listen for theme changes from other sources (cross-tab sync, system changes)
		alpine.XOn("theme-change.window", "theme = $event.detail.theme"),

		// Sun icon (visible in dark mode)
		html.Span(
			g.Attr("x-show", "theme === 'dark'"),
			g.Attr("x-cloak", ""),
			html.Class("inline-flex items-center"),
			sunIcon(),
		),

		// Moon icon (visible in light mode)
		html.Span(
			g.Attr("x-show", "theme === 'light'"),
			g.Attr("x-cloak", ""),
			html.Class("inline-flex items-center"),
			moonIcon(),
		),
	}

	// Optional label
	if props.ShowLabel {
		nodes = append(nodes, html.Span(
			g.Attr("x-text", fmt.Sprintf("theme === 'dark' ? '%s' : '%s'", props.DarkLabel, props.LightLabel)),
		))
	}

	// Add cross-tab sync listener if enabled
	if props.SyncAcrossTabs && props.StorageType == LocalStorage {
		nodes = append(nodes, g.Attr("x-init", buildTabSyncScript(props)))
	}

	return html.Button(nodes...)
}

// buildToggleInitScript creates the Alpine.js init script for theme initialization.
func buildToggleInitScript(props ToggleProps) string {
	storageGet := getStorageGetScript(props.StorageType, props.StorageKey)
	systemCheck := ""
	if props.RespectSystem {
		systemCheck = `
		if (!stored) {
			const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
			stored = system ? 'dark' : 'light';
		}`
	}

	return fmt.Sprintf(`
		let stored = %s;%s
		this.theme = stored || '%s';
	`, storageGet, systemCheck, props.DefaultTheme)
}

// buildToggleClickScript creates the click handler script for theme toggling.
func buildToggleClickScript(props ToggleProps) string {
	storageSet := getStorageSetScript(props.StorageType, props.StorageKey, "theme")

	transitionStart := ""
	transitionEnd := ""
	if props.AnimateTransition {
		transitionStart = `document.documentElement.classList.add('theme-transitioning');`
		transitionEnd = `setTimeout(() => document.documentElement.classList.remove('theme-transitioning'), 300);`
	}

	return fmt.Sprintf(`
		%s
		theme = theme === 'dark' ? 'light' : 'dark';
		%s
		document.documentElement.classList.toggle('dark', theme === 'dark');
		document.documentElement.setAttribute('data-theme', theme);
		window.dispatchEvent(new CustomEvent('theme-change', { detail: { theme } }));
		%s
	`, transitionStart, storageSet, transitionEnd)
}

// buildTabSyncScript creates the script for cross-tab synchronization.
func buildTabSyncScript(props ToggleProps) string {
	return fmt.Sprintf(`
		window.addEventListener('storage', (e) => {
			if (e.key === '%s' && e.newValue) {
				theme = e.newValue;
				document.documentElement.classList.toggle('dark', theme === 'dark');
				document.documentElement.setAttribute('data-theme', theme);
			}
		});
		if (%v) {
			window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
				const stored = %s;
				if (!stored) {
					theme = e.matches ? 'dark' : 'light';
					document.documentElement.classList.toggle('dark', theme === 'dark');
					document.documentElement.setAttribute('data-theme', theme);
				}
			});
		}
	`, props.StorageKey, props.RespectSystem, getStorageGetScript(props.StorageType, props.StorageKey))
}

// getStorageGetScript returns the JavaScript to get a value from storage.
func getStorageGetScript(storageType StorageType, key string) string {
	switch storageType {
	case LocalStorage:
		return fmt.Sprintf("localStorage.getItem('%s')", key)
	case SessionStorage:
		return fmt.Sprintf("sessionStorage.getItem('%s')", key)
	case NoStorage:
		return "null"
	default:
		return fmt.Sprintf("localStorage.getItem('%s')", key)
	}
}

// getStorageSetScript returns the JavaScript to set a value in storage.
func getStorageSetScript(storageType StorageType, key, valueVar string) string {
	switch storageType {
	case LocalStorage:
		return fmt.Sprintf("localStorage.setItem('%s', %s);", key, valueVar)
	case SessionStorage:
		return fmt.Sprintf("sessionStorage.setItem('%s', %s);", key, valueVar)
	case NoStorage:
		return ""
	default:
		return fmt.Sprintf("localStorage.setItem('%s', %s);", key, valueVar)
	}
}

// ToggleWithSystemOption creates a theme toggle with three options: light, dark, and system.
// This provides more granular control over theme preferences with full persistence support.
//
// Features:
//   - Three-way toggle: Light, Dark, System
//   - System option respects OS color scheme preference
//   - Persists preference to browser storage
//   - Syncs across browser tabs
//
// Usage:
//
//	ToggleWithSystemOption()                        // Default three-way toggle
//	ToggleWithSystemOption(WithStorageKey("pref")) // Custom storage key
//	ToggleWithSystemOption(IconOnly())              // Icons only, no labels
func ToggleWithSystemOption(opts ...ToggleOption) g.Node {
	props := defaultToggleProps()
	props.DefaultTheme = "system" // Default to system for three-way toggle
	for _, opt := range opts {
		opt(&props)
	}

	storageGet := getStorageGetScript(props.StorageType, props.StorageKey)

	initScript := fmt.Sprintf(`
		const stored = %s;
		this.theme = stored || 'system';
		this.applyTheme = (t) => {
			if (t === 'system') {
				const system = window.matchMedia('(prefers-color-scheme: dark)').matches;
				document.documentElement.classList.toggle('dark', system);
				document.documentElement.setAttribute('data-theme', system ? 'dark' : 'light');
			} else {
				document.documentElement.classList.toggle('dark', t === 'dark');
				document.documentElement.setAttribute('data-theme', t);
			}
		};
		this.applyTheme(this.theme);
	`, storageGet)

	nodes := []g.Node{
		html.Class("inline-flex items-center gap-1 p-1 rounded-lg bg-muted"),
		g.Attr("data-theme-toggle", "system"),

		// Alpine.js data for theme state with persistence
		alpine.XData(map[string]any{
			"theme":      props.DefaultTheme,
			"applyTheme": nil,
			"init":       initScript,
		}),

		// Light button
		themeOptionButton("light", props.LightLabel, sunIcon(), props),

		// Dark button
		themeOptionButton("dark", props.DarkLabel, moonIcon(), props),

		// System button
		themeOptionButton("system", props.SystemLabel, systemIcon(), props),
	}

	// Add cross-tab sync and system preference listener
	if props.SyncAcrossTabs && props.StorageType == LocalStorage {
		nodes = append(nodes, g.Attr("x-init", fmt.Sprintf(`
			window.addEventListener('storage', (e) => {
				if (e.key === '%s') {
					theme = e.newValue || 'system';
					applyTheme(theme);
				}
			});
			window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
				if (theme === 'system') {
					applyTheme('system');
				}
			});
		`, props.StorageKey)))
	}

	return html.Div(nodes...)
}

// themeOptionButton creates a single button for the three-option toggle.
func themeOptionButton(value, label string, icon g.Node, props ToggleProps) g.Node {
	storageSet := ""
	storageRemove := ""

	switch props.StorageType {
	case LocalStorage:
		storageSet = fmt.Sprintf("localStorage.setItem('%s', theme);", props.StorageKey)
		storageRemove = fmt.Sprintf("localStorage.removeItem('%s');", props.StorageKey)
	case SessionStorage:
		storageSet = fmt.Sprintf("sessionStorage.setItem('%s', theme);", props.StorageKey)
		storageRemove = fmt.Sprintf("sessionStorage.removeItem('%s');", props.StorageKey)
	case NoStorage:
		storageSet = ""
		storageRemove = ""
	}

	clickScript := fmt.Sprintf(`
		theme = '%s';
		if (theme === 'system') {
			%s
		} else {
			%s
		}
		applyTheme(theme);
		window.dispatchEvent(new CustomEvent('theme-change', { detail: { theme } }));
	`, value, storageRemove, storageSet)

	return html.Button(
		g.Attr("type", "button"),
		html.Class("inline-flex items-center justify-center gap-2 px-3 py-1.5 text-sm font-medium rounded-md transition-colors"),
		g.Attr(":class", "theme === '"+value+"' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'"),
		g.Attr("aria-label", label),
		g.Attr("aria-pressed", ""),
		g.Attr(":aria-pressed", "theme === '"+value+"'"),

		alpine.XOn("click", clickScript),

		icon,

		g.If(!props.IconOnly,
			html.Span(g.Text(label)),
		),
	)
}

// sunIcon returns a sun SVG icon for light mode.
func sunIcon() g.Node {
	return g.El("svg",
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("width", "16"),
		g.Attr("height", "16"),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", "currentColor"),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
		html.Class("shrink-0"),
		g.El("circle", g.Attr("cx", "12"), g.Attr("cy", "12"), g.Attr("r", "4")),
		g.El("path", g.Attr("d", "M12 2v2")),
		g.El("path", g.Attr("d", "M12 20v2")),
		g.El("path", g.Attr("d", "m4.93 4.93 1.41 1.41")),
		g.El("path", g.Attr("d", "m17.66 17.66 1.41 1.41")),
		g.El("path", g.Attr("d", "M2 12h2")),
		g.El("path", g.Attr("d", "M20 12h2")),
		g.El("path", g.Attr("d", "m6.34 17.66-1.41 1.41")),
		g.El("path", g.Attr("d", "m19.07 4.93-1.41 1.41")),
	)
}

// moonIcon returns a moon SVG icon for dark mode.
func moonIcon() g.Node {
	return g.El("svg",
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("width", "16"),
		g.Attr("height", "16"),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", "currentColor"),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
		html.Class("shrink-0"),
		g.El("path", g.Attr("d", "M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z")),
	)
}

// systemIcon returns a monitor SVG icon for system theme.
func systemIcon() g.Node {
	return g.El("svg",
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("width", "16"),
		g.Attr("height", "16"),
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", "currentColor"),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
		html.Class("shrink-0"),
		g.El("rect", g.Attr("x", "2"), g.Attr("y", "3"), g.Attr("width", "20"), g.Attr("height", "14"), g.Attr("rx", "2")),
		g.El("path", g.Attr("d", "M8 21h8")),
		g.El("path", g.Attr("d", "M12 17v4")),
	)
}

// SimpleToggle creates a minimal icon-only theme toggle with persistence.
// This is a convenience function for the most common use case.
func SimpleToggle() g.Node {
	return Toggle(IconOnly())
}

// PersistentToggle creates a theme toggle with full persistence support.
// This is a convenience function that enables all persistence features.
func PersistentToggle(storageKey string) g.Node {
	return Toggle(
		IconOnly(),
		WithToggleStorageKey(storageKey),
		WithStorageType(LocalStorage),
		WithTabSync(true),
		WithSystemPreference(true),
	)
}

// DropdownToggle creates a dropdown-style theme selector.
// This provides a more compact interface with a popover menu.
//
// Usage:
//
//	DropdownToggle()
//	DropdownToggle(WithStorageKey("my-theme"))
func DropdownToggle(opts ...ToggleOption) g.Node {
	props := defaultToggleProps()
	props.DefaultTheme = "system"
	for _, opt := range opts {
		opt(&props)
	}

	storageGet := getStorageGetScript(props.StorageType, props.StorageKey)

	initScript := fmt.Sprintf(`
		const stored = %s;
		this.theme = stored || 'system';
		this.open = false;
		this.getEffectiveTheme = () => {
			if (this.theme === 'system') {
				return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
			}
			return this.theme;
		};
		this.setTheme = (t) => {
			this.theme = t;
			this.open = false;
			if (t === 'system') {
				%s
			} else {
				%s
			}
			const effective = this.getEffectiveTheme();
			document.documentElement.classList.toggle('dark', effective === 'dark');
			document.documentElement.setAttribute('data-theme', effective);
			window.dispatchEvent(new CustomEvent('theme-change', { detail: { theme: t } }));
		};
		const effective = this.getEffectiveTheme();
		document.documentElement.classList.toggle('dark', effective === 'dark');
		document.documentElement.setAttribute('data-theme', effective);
	`, storageGet,
		getStorageRemoveScript(props.StorageType, props.StorageKey),
		getStorageSetScript(props.StorageType, props.StorageKey, "t"))

	return html.Div(
		html.Class("relative inline-block text-left"),
		g.Attr("data-theme-toggle", "dropdown"),
		alpine.XData(map[string]any{
			"theme":            props.DefaultTheme,
			"open":             false,
			"getEffectiveTheme": nil,
			"setTheme":         nil,
			"init":             initScript,
		}),

		// Trigger button
		html.Button(
			g.Attr("type", "button"),
			html.Class("inline-flex items-center justify-center gap-2 h-9 px-3 rounded-md text-sm font-medium transition-colors "+
				"bg-transparent border border-input hover:bg-accent hover:text-accent-foreground "+
				"focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"),
			g.Attr("aria-label", "Select theme"),
			g.Attr("aria-haspopup", "true"),
			g.Attr(":aria-expanded", "open"),
			alpine.XOn("click", "open = !open"),

			// Dynamic icon based on current theme
			html.Span(
				g.Attr("x-show", "getEffectiveTheme() === 'dark'"),
				g.Attr("x-cloak", ""),
				moonIcon(),
			),
			html.Span(
				g.Attr("x-show", "getEffectiveTheme() === 'light'"),
				g.Attr("x-cloak", ""),
				sunIcon(),
			),

			// Chevron down icon
			g.El("svg",
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("width", "12"),
				g.Attr("height", "12"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("fill", "none"),
				g.Attr("stroke", "currentColor"),
				g.Attr("stroke-width", "2"),
				html.Class("shrink-0 ml-1"),
				g.El("path", g.Attr("d", "m6 9 6 6 6-6")),
			),
		),

		// Dropdown menu
		html.Div(
			g.Attr("x-show", "open"),
			g.Attr("x-cloak", ""),
			g.Attr("x-transition:enter", "transition ease-out duration-100"),
			g.Attr("x-transition:enter-start", "transform opacity-0 scale-95"),
			g.Attr("x-transition:enter-end", "transform opacity-100 scale-100"),
			g.Attr("x-transition:leave", "transition ease-in duration-75"),
			g.Attr("x-transition:leave-start", "transform opacity-100 scale-100"),
			g.Attr("x-transition:leave-end", "transform opacity-0 scale-95"),
			alpine.XOn("click.outside", "open = false"),
			html.Class("absolute right-0 z-50 mt-2 w-36 origin-top-right rounded-md bg-popover border border-border shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"),
			g.Attr("role", "menu"),
			g.Attr("aria-orientation", "vertical"),

			html.Div(
				html.Class("py-1"),
				g.Attr("role", "none"),

				// Light option
				dropdownMenuItem("light", props.LightLabel, sunIcon()),

				// Dark option
				dropdownMenuItem("dark", props.DarkLabel, moonIcon()),

				// System option
				dropdownMenuItem("system", props.SystemLabel, systemIcon()),
			),
		),
	)
}

// dropdownMenuItem creates a single menu item for the dropdown toggle.
func dropdownMenuItem(value, label string, icon g.Node) g.Node {
	return html.Button(
		g.Attr("type", "button"),
		html.Class("w-full flex items-center gap-2 px-3 py-2 text-sm transition-colors hover:bg-accent"),
		g.Attr(":class", "theme === '"+value+"' ? 'bg-accent text-accent-foreground' : 'text-foreground'"),
		g.Attr("role", "menuitem"),
		alpine.XOn("click", "setTheme('"+value+"')"),
		icon,
		html.Span(g.Text(label)),
		// Checkmark for selected
		html.Span(
			g.Attr("x-show", "theme === '"+value+"'"),
			g.Attr("x-cloak", ""),
			html.Class("ml-auto"),
			g.El("svg",
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("width", "14"),
				g.Attr("height", "14"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("fill", "none"),
				g.Attr("stroke", "currentColor"),
				g.Attr("stroke-width", "2"),
				g.El("polyline", g.Attr("points", "20 6 9 17 4 12")),
			),
		),
	)
}

// getStorageRemoveScript returns the JavaScript to remove a value from storage.
func getStorageRemoveScript(storageType StorageType, key string) string {
	switch storageType {
	case LocalStorage:
		return fmt.Sprintf("localStorage.removeItem('%s');", key)
	case SessionStorage:
		return fmt.Sprintf("sessionStorage.removeItem('%s');", key)
	case NoStorage:
		return ""
	default:
		return fmt.Sprintf("localStorage.removeItem('%s');", key)
	}
}

// ThemeTransitionCSS returns CSS for smooth theme transitions.
// Add this to your page to enable animated theme changes.
func ThemeTransitionCSS() g.Node {
	css := `
.theme-transitioning,
.theme-transitioning * {
	transition: background-color 0.3s ease, color 0.3s ease, border-color 0.3s ease !important;
}
`
	return g.El("style",
		g.Attr("data-theme-transition", ""),
		g.Raw(css),
	)
}

