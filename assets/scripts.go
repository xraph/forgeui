package assets

import (
	"sort"
	"sync"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ScriptEntry represents a script with metadata for ordering and rendering.
// Scripts can be inline or external, and are prioritized for optimal loading.
type ScriptEntry struct {
	// Path is the URL or file path to the script
	Path string

	// Priority determines load order (0-100, lower loads first)
	// Typical values:
	//   0-19: Critical framework scripts (Alpine.js, etc.)
	//   20-49: Library dependencies
	//   50-79: Application scripts
	//   80-100: Analytics, third-party widgets
	Priority int

	// Position determines where the script is rendered ("head" or "body")
	Position string

	// Inline indicates if this is an inline script (Content is used instead of Path)
	Inline bool

	// Content is the inline script content (only used if Inline is true)
	Content string

	// Attrs are additional HTML attributes (defer, async, type, etc.)
	Attrs map[string]string
}

// ScriptManager manages script loading order and rendering.
// It ensures scripts are loaded in the correct order based on priority
// and prevents duplicate script tags.
type ScriptManager struct {
	scripts []ScriptEntry
	mu      sync.RWMutex
}

// NewScriptManager creates a new script manager
func NewScriptManager() *ScriptManager {
	return &ScriptManager{
		scripts: make([]ScriptEntry, 0),
	}
}

// Add adds a script entry to the manager.
// Duplicate paths are ignored (first occurrence wins).
func (sm *ScriptManager) Add(entry ScriptEntry) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check for duplicates by path (ignore inline scripts in duplicate check)
	if !entry.Inline && entry.Path != "" {
		for _, existing := range sm.scripts {
			if existing.Path == entry.Path {
				return // Already exists
			}
		}
	}

	// Set defaults
	if entry.Position == "" {
		entry.Position = "body"
	}

	if entry.Attrs == nil {
		entry.Attrs = make(map[string]string)
	}

	sm.scripts = append(sm.scripts, entry)
}

// AddWithPriority is a convenience method to add a script with a specific priority
func (sm *ScriptManager) AddWithPriority(path string, priority int, attrs map[string]string) {
	sm.Add(ScriptEntry{
		Path:     path,
		Priority: priority,
		Position: "body",
		Attrs:    attrs,
	})
}

// AddInline adds an inline script with the given content
func (sm *ScriptManager) AddInline(content string, priority int, position string) {
	sm.Add(ScriptEntry{
		Inline:   true,
		Content:  content,
		Priority: priority,
		Position: position,
	})
}

// Render generates script tags for the specified position ("head" or "body").
// Scripts are sorted by priority (lower numbers first) before rendering.
func (sm *ScriptManager) Render(position string) []g.Node {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Filter scripts for this position
	var positionScripts []ScriptEntry

	for _, script := range sm.scripts {
		if script.Position == position {
			positionScripts = append(positionScripts, script)
		}
	}

	// Sort by priority (lower first)
	sort.Slice(positionScripts, func(i, j int) bool {
		return positionScripts[i].Priority < positionScripts[j].Priority
	})

	// Generate script nodes
	nodes := make([]g.Node, 0, len(positionScripts))
	for _, script := range positionScripts {
		nodes = append(nodes, sm.renderScript(script))
	}

	return nodes
}

// renderScript creates a single script node from an entry
func (sm *ScriptManager) renderScript(entry ScriptEntry) g.Node {
	if entry.Inline {
		// Inline script
		attrs := make([]g.Node, 0, len(entry.Attrs)+1)

		// Add custom attributes
		for key, value := range entry.Attrs {
			if value == "" {
				attrs = append(attrs, g.Attr(key))
			} else {
				attrs = append(attrs, g.Attr(key, value))
			}
		}

		// Add content
		attrs = append(attrs, g.Raw(entry.Content))

		return html.Script(attrs...)
	}

	// External script
	attrs := make([]g.Node, 0, len(entry.Attrs)+1)

	// Add src attribute
	attrs = append(attrs, g.Attr("src", entry.Path))

	// Add custom attributes
	for key, value := range entry.Attrs {
		if value == "" {
			attrs = append(attrs, g.Attr(key))
		} else {
			attrs = append(attrs, g.Attr(key, value))
		}
	}

	return html.Script(attrs...)
}

// Clear removes all scripts from the manager
func (sm *ScriptManager) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.scripts = make([]ScriptEntry, 0)
}

// Count returns the total number of scripts
func (sm *ScriptManager) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.scripts)
}

// CountByPosition returns the number of scripts in a specific position
func (sm *ScriptManager) CountByPosition(position string) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	count := 0

	for _, script := range sm.scripts {
		if script.Position == position {
			count++
		}
	}

	return count
}
