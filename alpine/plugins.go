package alpine

// Plugin represents an official Alpine.js plugin.
type Plugin string

const (
	// PluginMask provides input masking functionality
	// https://alpinejs.dev/plugins/mask
	PluginMask Plugin = "mask"

	// PluginIntersect provides intersection observer functionality
	// https://alpinejs.dev/plugins/intersect
	PluginIntersect Plugin = "intersect"

	// PluginPersist provides localStorage persistence
	// https://alpinejs.dev/plugins/persist
	PluginPersist Plugin = "persist"

	// PluginFocus provides focus management utilities
	// https://alpinejs.dev/plugins/focus
	PluginFocus Plugin = "focus"

	// PluginCollapse provides smooth height transitions
	// https://alpinejs.dev/plugins/collapse
	PluginCollapse Plugin = "collapse"

	// PluginAnchor provides element positioning
	// https://alpinejs.dev/plugins/anchor
	PluginAnchor Plugin = "anchor"

	// PluginMorph provides DOM morphing capabilities
	// https://alpinejs.dev/plugins/morph
	PluginMorph Plugin = "morph"

	// PluginSort provides drag-and-drop sorting
	// https://alpinejs.dev/plugins/sort
	PluginSort Plugin = "sort"
)

// pluginURLs maps plugins to their CDN URLs.
// Using 3.x.x for automatic minor/patch updates while staying on v3.
var pluginURLs = map[Plugin]string{
	PluginMask:      "https://cdn.jsdelivr.net/npm/@alpinejs/mask@3.x.x/dist/cdn.min.js",
	PluginIntersect: "https://cdn.jsdelivr.net/npm/@alpinejs/intersect@3.x.x/dist/cdn.min.js",
	PluginPersist:   "https://cdn.jsdelivr.net/npm/@alpinejs/persist@3.x.x/dist/cdn.min.js",
	PluginFocus:     "https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js",
	PluginCollapse:  "https://cdn.jsdelivr.net/npm/@alpinejs/collapse@3.x.x/dist/cdn.min.js",
	PluginAnchor:    "https://cdn.jsdelivr.net/npm/@alpinejs/anchor@3.x.x/dist/cdn.min.js",
	PluginMorph:     "https://cdn.jsdelivr.net/npm/@alpinejs/morph@3.x.x/dist/cdn.min.js",
	PluginSort:      "https://cdn.jsdelivr.net/npm/@alpinejs/sort@3.x.x/dist/cdn.min.js",
}

// PluginURL returns the CDN URL for the given plugin.
// Returns empty string if plugin is not recognized.
func PluginURL(p Plugin) string {
	return pluginURLs[p]
}

// AllPlugins returns a list of all available plugins.
func AllPlugins() []Plugin {
	return []Plugin{
		PluginMask,
		PluginIntersect,
		PluginPersist,
		PluginFocus,
		PluginCollapse,
		PluginAnchor,
		PluginMorph,
		PluginSort,
	}
}

