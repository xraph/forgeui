package theme

import (
	"context"
	"fmt"
	stdhtml "html"
	"io"

	"github.com/a-h/templ"

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
	ShowLabel         bool        // Show text label alongside icon
	LightLabel        string      // Label for light mode
	DarkLabel         string      // Label for dark mode
	SystemLabel       string      // Label for system mode
	Position          string      // Button position (if floating)
	Size              string      // Button size (sm, md, lg)
	CustomClass       string      // Additional CSS classes
	IconOnly          bool        // Show only icon, no label
	StorageKey        string      // Key used for storage (default: "theme")
	StorageType       StorageType // Storage mechanism (localStorage, sessionStorage, none)
	SyncAcrossTabs    bool        // Sync theme changes across browser tabs
	RespectSystem     bool        // Respect system preference when no stored value
	DefaultTheme      string      // Default theme when no preference (light, dark, system)
	AnimateTransition bool        // Enable smooth transition animation
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
		ShowLabel:         false,
		LightLabel:        "Light",
		DarkLabel:         "Dark",
		SystemLabel:       "System",
		Size:              "md",
		IconOnly:          true,
		StorageKey:        "theme",
		StorageType:       LocalStorage,
		SyncAcrossTabs:    true,
		RespectSystem:     true,
		DefaultTheme:      "light",
		AnimateTransition: false,
	}
}

// writeAlpineAttrs writes templ.Attributes as HTML attributes on the current element.
func writeAlpineAttrs(w io.Writer, attrs templ.Attributes) error {
	for k, v := range attrs {
		switch val := v.(type) {
		case string:
			if _, err := fmt.Fprintf(w, ` %s="%s"`, k, stdhtml.EscapeString(val)); err != nil {
				return err
			}
		case bool:
			if val {
				if _, err := fmt.Fprintf(w, " %s", k); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Toggle creates a theme switcher button that toggles between light and dark modes.
func Toggle(opts ...ToggleOption) templ.Component {
	props := defaultToggleProps()
	for _, opt := range opts {
		opt(&props)
	}

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

	initScript := buildToggleInitScript(props)
	clickScript := buildToggleClickScript(props)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<button type="button" class="%s" aria-label="Toggle theme" data-theme-toggle`, stdhtml.EscapeString(classes)); err != nil {
			return err
		}

		// Alpine.js data
		if err := writeAlpineAttrs(w, alpine.XData(map[string]any{
			"theme": props.DefaultTheme,
			"init":  initScript,
		})); err != nil {
			return err
		}

		// Click handler
		if err := writeAlpineAttrs(w, alpine.XOn("click", clickScript)); err != nil {
			return err
		}

		// Theme change listener
		if err := writeAlpineAttrs(w, alpine.XOn("theme-change.window", "theme = $event.detail.theme")); err != nil {
			return err
		}

		// Tab sync
		if props.SyncAcrossTabs && props.StorageType == LocalStorage {
			if _, err := fmt.Fprintf(w, ` x-init="%s"`, stdhtml.EscapeString(buildTabSyncScript(props))); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}

		// Sun icon (visible in dark mode)
		if _, err := io.WriteString(w, `<span x-show="theme === 'dark'" x-cloak class="inline-flex items-center">`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, sunIconSVG); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `</span>`); err != nil {
			return err
		}

		// Moon icon (visible in light mode)
		if _, err := io.WriteString(w, `<span x-show="theme === 'light'" x-cloak class="inline-flex items-center">`); err != nil {
			return err
		}
		if _, err := io.WriteString(w, moonIconSVG); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `</span>`); err != nil {
			return err
		}

		// Optional label
		if props.ShowLabel {
			if _, err := fmt.Fprintf(w, `<span x-text="theme === 'dark' ? '%s' : '%s'"></span>`,
				stdhtml.EscapeString(props.DarkLabel),
				stdhtml.EscapeString(props.LightLabel)); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, `</button>`)
		return err
	})
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

// SVG icon constants
const sunIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="shrink-0"><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2"></path><path d="M12 20v2"></path><path d="m4.93 4.93 1.41 1.41"></path><path d="m17.66 17.66 1.41 1.41"></path><path d="M2 12h2"></path><path d="M20 12h2"></path><path d="m6.34 17.66-1.41 1.41"></path><path d="m19.07 4.93-1.41 1.41"></path></svg>`

const moonIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="shrink-0"><path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path></svg>`

const systemIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="shrink-0"><rect x="2" y="3" width="20" height="14" rx="2"></rect><path d="M8 21h8"></path><path d="M12 17v4"></path></svg>`

const checkIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"></polyline></svg>`

const chevronDownSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="shrink-0 ml-1"><path d="m6 9 6 6 6-6"></path></svg>`

// ToggleWithSystemOption creates a theme toggle with three options: light, dark, and system.
func ToggleWithSystemOption(opts ...ToggleOption) templ.Component {
	props := defaultToggleProps()
	props.DefaultTheme = "system"
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

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="inline-flex items-center gap-1 p-1 rounded-lg bg-muted" data-theme-toggle="system"`); err != nil {
			return err
		}

		if err := writeAlpineAttrs(w, alpine.XData(map[string]any{
			"theme":      props.DefaultTheme,
			"applyTheme": nil,
			"init":       initScript,
		})); err != nil {
			return err
		}

		// Tab sync
		if props.SyncAcrossTabs && props.StorageType == LocalStorage {
			if _, err := fmt.Fprintf(w, ` x-init="%s"`, stdhtml.EscapeString(fmt.Sprintf(`
				window.addEventListener('storage', (e) => {
					if (e.key === '%s') {
						theme = e.newValue || 'system';
						applyTheme(theme);
					}
				});
				window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
					if (theme === 'system') { applyTheme('system'); }
				});
			`, props.StorageKey))); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}

		// Light, Dark, System buttons
		if err := writeThemeOptionButton(w, "light", props.LightLabel, sunIconSVG, props); err != nil {
			return err
		}
		if err := writeThemeOptionButton(w, "dark", props.DarkLabel, moonIconSVG, props); err != nil {
			return err
		}
		if err := writeThemeOptionButton(w, "system", props.SystemLabel, systemIconSVG, props); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)
		return err
	})
}

// writeThemeOptionButton writes a single option button for the three-way toggle.
func writeThemeOptionButton(w io.Writer, value, label, iconSVG string, props ToggleProps) error {
	storageSet := ""
	storageRemove := ""

	switch props.StorageType {
	case LocalStorage:
		storageSet = fmt.Sprintf("localStorage.setItem('%s', theme);", props.StorageKey)
		storageRemove = fmt.Sprintf("localStorage.removeItem('%s');", props.StorageKey)
	case SessionStorage:
		storageSet = fmt.Sprintf("sessionStorage.setItem('%s', theme);", props.StorageKey)
		storageRemove = fmt.Sprintf("sessionStorage.removeItem('%s');", props.StorageKey)
	}

	clickScript := fmt.Sprintf(`
		theme = '%s';
		if (theme === 'system') { %s } else { %s }
		applyTheme(theme);
		window.dispatchEvent(new CustomEvent('theme-change', { detail: { theme } }));
	`, value, storageRemove, storageSet)

	if _, err := fmt.Fprintf(w,
		`<button type="button" class="inline-flex items-center justify-center gap-2 px-3 py-1.5 text-sm font-medium rounded-md transition-colors" :class="theme === '%s' ? 'bg-background shadow-sm' : 'text-muted-foreground hover:text-foreground'" aria-label="%s" aria-pressed :aria-pressed="theme === '%s'"`,
		value, stdhtml.EscapeString(label), value); err != nil {
		return err
	}

	if err := writeAlpineAttrs(w, alpine.XOn("click", clickScript)); err != nil {
		return err
	}

	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	if _, err := io.WriteString(w, iconSVG); err != nil {
		return err
	}

	if !props.IconOnly {
		if _, err := fmt.Fprintf(w, `<span>%s</span>`, stdhtml.EscapeString(label)); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, `</button>`)
	return err
}

// SimpleToggle creates a minimal icon-only theme toggle with persistence.
func SimpleToggle() templ.Component {
	return Toggle(IconOnly())
}

// PersistentToggle creates a theme toggle with full persistence support.
func PersistentToggle(storageKey string) templ.Component {
	return Toggle(
		IconOnly(),
		WithToggleStorageKey(storageKey),
		WithStorageType(LocalStorage),
		WithTabSync(true),
		WithSystemPreference(true),
	)
}

// DropdownToggle creates a dropdown-style theme selector.
func DropdownToggle(opts ...ToggleOption) templ.Component {
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

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="relative inline-block text-left" data-theme-toggle="dropdown"`); err != nil {
			return err
		}

		if err := writeAlpineAttrs(w, alpine.XData(map[string]any{
			"theme":             props.DefaultTheme,
			"open":              false,
			"getEffectiveTheme": nil,
			"setTheme":          nil,
			"init":              initScript,
		})); err != nil {
			return err
		}

		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}

		// Trigger button
		if _, err := io.WriteString(w, `<button type="button" class="inline-flex items-center justify-center gap-2 h-9 px-3 rounded-md text-sm font-medium transition-colors bg-transparent border border-input hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" aria-label="Select theme" aria-haspopup="true" :aria-expanded="open"`); err != nil {
			return err
		}
		if err := writeAlpineAttrs(w, alpine.XOn("click", "open = !open")); err != nil {
			return err
		}
		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}

		// Dynamic icons
		if _, err := fmt.Fprintf(w, `<span x-show="getEffectiveTheme() === 'dark'" x-cloak>%s</span>`, moonIconSVG); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, `<span x-show="getEffectiveTheme() === 'light'" x-cloak>%s</span>`, sunIconSVG); err != nil {
			return err
		}
		if _, err := io.WriteString(w, chevronDownSVG); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `</button>`); err != nil {
			return err
		}

		// Dropdown menu
		if _, err := io.WriteString(w, `<div x-show="open" x-cloak x-transition:enter="transition ease-out duration-100" x-transition:enter-start="transform opacity-0 scale-95" x-transition:enter-end="transform opacity-100 scale-100" x-transition:leave="transition ease-in duration-75" x-transition:leave-start="transform opacity-100 scale-100" x-transition:leave-end="transform opacity-0 scale-95"`); err != nil {
			return err
		}
		if err := writeAlpineAttrs(w, alpine.XOn("click.outside", "open = false")); err != nil {
			return err
		}
		if _, err := io.WriteString(w, ` class="absolute right-0 z-50 mt-2 w-36 origin-top-right rounded-md bg-popover border border-border shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none" role="menu" aria-orientation="vertical">`); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `<div class="py-1" role="none">`); err != nil {
			return err
		}

		// Menu items
		if err := writeDropdownMenuItem(w, "light", props.LightLabel, sunIconSVG); err != nil {
			return err
		}
		if err := writeDropdownMenuItem(w, "dark", props.DarkLabel, moonIconSVG); err != nil {
			return err
		}
		if err := writeDropdownMenuItem(w, "system", props.SystemLabel, systemIconSVG); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `</div></div></div>`); err != nil {
			return err
		}

		return nil
	})
}

// writeDropdownMenuItem writes a single menu item for the dropdown toggle.
func writeDropdownMenuItem(w io.Writer, value, label, iconSVG string) error {
	if _, err := fmt.Fprintf(w,
		`<button type="button" class="w-full flex items-center gap-2 px-3 py-2 text-sm transition-colors hover:bg-accent" :class="theme === '%s' ? 'bg-accent text-accent-foreground' : 'text-foreground'" role="menuitem"`,
		value); err != nil {
		return err
	}

	if err := writeAlpineAttrs(w, alpine.XOn("click", "setTheme('"+value+"')")); err != nil {
		return err
	}

	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	if _, err := io.WriteString(w, iconSVG); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, `<span>%s</span>`, stdhtml.EscapeString(label)); err != nil {
		return err
	}

	// Checkmark for selected
	if _, err := fmt.Fprintf(w, `<span x-show="theme === '%s'" x-cloak class="ml-auto">%s</span>`, value, checkIconSVG); err != nil {
		return err
	}

	_, err := io.WriteString(w, `</button>`)
	return err
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
func ThemeTransitionCSS() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<style data-theme-transition>
.theme-transitioning,
.theme-transitioning * {
	transition: background-color 0.3s ease, color 0.3s ease, border-color 0.3s ease !important;
}
</style>`)
		return err
	})
}
