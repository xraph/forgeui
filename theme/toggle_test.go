package theme_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/xraph/forgeui/theme"
)

func TestToggle(t *testing.T) {
	var buf bytes.Buffer

	node := theme.Toggle()

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for button element
	if !strings.Contains(html, "<button") {
		t.Error("Toggle should render a button element")
	}

	// Check for Alpine.js directives
	if !strings.Contains(html, "x-data") {
		t.Error("Toggle should have x-data directive")
	}

	// Alpine.js click handler can be x-on:click or @click
	if !strings.Contains(html, "x-on:click") && !strings.Contains(html, "@click") {
		t.Error("Toggle should have click event handler")
	}

	// Check for aria-label
	if !strings.Contains(html, "aria-label") {
		t.Error("Toggle should have aria-label for accessibility")
	}

	// Check for icons (svg)
	if !strings.Contains(html, "<svg") {
		t.Error("Toggle should contain SVG icons")
	}
}

func TestToggleWithLabel(t *testing.T) {
	var buf bytes.Buffer

	node := theme.Toggle(theme.WithLabel(true))

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for x-text directive (label)
	if !strings.Contains(html, "x-text") {
		t.Error("Toggle with label should have x-text directive")
	}
}

func TestToggleWithCustomLabels(t *testing.T) {
	var buf bytes.Buffer

	node := theme.Toggle(theme.WithLabels("☀️ Light", "🌙 Dark"))

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for custom labels
	if !strings.Contains(html, "Light") {
		t.Error("Toggle should contain custom light label")
	}

	if !strings.Contains(html, "Dark") {
		t.Error("Toggle should contain custom dark label")
	}
}

func TestToggleWithSize(t *testing.T) {
	sizes := []string{"sm", "md", "lg"}

	for _, size := range sizes {
		t.Run(size, func(t *testing.T) {
			var buf bytes.Buffer

			node := theme.Toggle(theme.WithToggleSize(size))

			if err := node.Render(context.Background(), &buf); err != nil {
				t.Fatalf("render error: %v", err)
			}

			html := buf.String()

			if !strings.Contains(html, "<button") {
				t.Error("Toggle should render a button")
			}
		})
	}
}

func TestSimpleToggle(t *testing.T) {
	var buf bytes.Buffer

	node := theme.SimpleToggle()

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "<button") {
		t.Error("SimpleToggle should render a button")
	}

	// Should be icon-only
	if !strings.Contains(html, "x-data") {
		t.Error("SimpleToggle should have Alpine.js directives")
	}
}

func TestToggleWithSystemOption(t *testing.T) {
	var buf bytes.Buffer

	node := theme.ToggleWithSystemOption()

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for three buttons
	if strings.Count(html, "<button") != 3 {
		t.Error("ToggleWithSystemOption should have three buttons")
	}

	// Check for theme options
	expectedOptions := []string{"light", "dark", "system"}
	for _, opt := range expectedOptions {
		if !strings.Contains(html, opt) {
			t.Errorf("Toggle should contain %s option", opt)
		}
	}
}

func TestToggleIcons(t *testing.T) {
	var buf bytes.Buffer

	node := theme.Toggle()

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for sun and moon icons
	if !strings.Contains(html, "x-show") {
		t.Error("Icons should use x-show for conditional display")
	}

	if !strings.Contains(html, "x-cloak") {
		t.Error("Icons should use x-cloak to prevent flash")
	}
}

func TestTogglePersistence(t *testing.T) {
	t.Run("default localStorage persistence", func(t *testing.T) {
		var buf bytes.Buffer

		node := theme.Toggle()

		if err := node.Render(context.Background(), &buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()

		// Check for localStorage usage
		if !strings.Contains(html, "localStorage") {
			t.Error("Toggle should use localStorage by default")
		}

		// Check for default storage key (HTML-escaped quotes)
		if !strings.Contains(html, "&#39;theme&#39;") && !strings.Contains(html, "'theme'") {
			t.Error("Toggle should use 'theme' as default storage key")
		}
	})

	t.Run("custom storage key", func(t *testing.T) {
		var buf bytes.Buffer

		node := theme.Toggle(theme.WithToggleStorageKey("my-app-theme"))

		if err := node.Render(context.Background(), &buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()

		if !strings.Contains(html, "my-app-theme") {
			t.Error("Toggle should use custom storage key")
		}
	})

	t.Run("sessionStorage type", func(t *testing.T) {
		var buf bytes.Buffer

		node := theme.Toggle(theme.WithStorageType(theme.SessionStorage))

		if err := node.Render(context.Background(), &buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()

		if !strings.Contains(html, "sessionStorage") {
			t.Error("Toggle should use sessionStorage when specified")
		}
	})

	t.Run("no storage type", func(t *testing.T) {
		var buf bytes.Buffer

		node := theme.Toggle(theme.WithStorageType(theme.NoStorage))

		if err := node.Render(context.Background(), &buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()

		// Should not have storage calls
		if strings.Contains(html, "localStorage.setItem") || strings.Contains(html, "sessionStorage.setItem") {
			t.Error("Toggle with NoStorage should not persist")
		}
	})
}

func TestToggleTabSync(t *testing.T) {
	t.Run("tab sync enabled by default", func(t *testing.T) {
		var buf bytes.Buffer

		node := theme.Toggle()

		if err := node.Render(context.Background(), &buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()

		// Check for storage event listener
		if !strings.Contains(html, "storage") {
			t.Error("Toggle should listen for storage events for tab sync")
		}
	})

	t.Run("tab sync disabled", func(t *testing.T) {
		var buf bytes.Buffer

		node := theme.Toggle(theme.WithTabSync(false))

		if err := node.Render(context.Background(), &buf); err != nil {
			t.Fatalf("render error: %v", err)
		}

		html := buf.String()

		// Should not have x-init with storage listener
		if strings.Contains(html, "x-init") && strings.Contains(html, "addEventListener('storage'") {
			t.Error("Toggle with tab sync disabled should not add storage listener")
		}
	})
}

func TestToggleSystemPreference(t *testing.T) {
	var buf bytes.Buffer

	node := theme.Toggle(theme.WithSystemPreference(true))

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for system preference detection
	if !strings.Contains(html, "prefers-color-scheme") {
		t.Error("Toggle should detect system color scheme preference")
	}
}

func TestPersistentToggle(t *testing.T) {
	var buf bytes.Buffer

	node := theme.PersistentToggle("custom-key")

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for custom storage key
	if !strings.Contains(html, "custom-key") {
		t.Error("PersistentToggle should use provided storage key")
	}

	// Check for localStorage
	if !strings.Contains(html, "localStorage") {
		t.Error("PersistentToggle should use localStorage")
	}
}

func TestDropdownToggle(t *testing.T) {
	var buf bytes.Buffer

	node := theme.DropdownToggle()

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for dropdown container
	if !strings.Contains(html, "relative") {
		t.Error("DropdownToggle should have relative positioning")
	}

	// Check for trigger button with aria-haspopup
	if !strings.Contains(html, "aria-haspopup") {
		t.Error("DropdownToggle should have aria-haspopup attribute")
	}

	// Check for three menu items (light, dark, system)
	if !strings.Contains(html, "menuitem") {
		t.Error("DropdownToggle should have menu items")
	}

	// Check for transitions
	if !strings.Contains(html, "x-transition") {
		t.Error("DropdownToggle should have transition animations")
	}
}

func TestToggleWithSystemOptionPersistence(t *testing.T) {
	var buf bytes.Buffer

	node := theme.ToggleWithSystemOption(theme.WithToggleStorageKey("app-theme"))

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for custom storage key
	if !strings.Contains(html, "app-theme") {
		t.Error("ToggleWithSystemOption should use custom storage key")
	}

	// Check for three buttons (light, dark, system)
	if strings.Count(html, "<button") != 3 {
		t.Error("ToggleWithSystemOption should have three buttons")
	}

	// Check for aria-pressed attribute
	if !strings.Contains(html, "aria-pressed") {
		t.Error("ToggleWithSystemOption buttons should have aria-pressed")
	}
}

func TestThemeTransitionCSS(t *testing.T) {
	var buf bytes.Buffer

	node := theme.ThemeTransitionCSS()

	if err := node.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render error: %v", err)
	}

	html := buf.String()

	// Check for style tag
	if !strings.Contains(html, "<style") {
		t.Error("ThemeTransitionCSS should render a style tag")
	}

	// Check for transition class
	if !strings.Contains(html, "theme-transitioning") {
		t.Error("ThemeTransitionCSS should define theme-transitioning class")
	}

	// Check for transition properties
	if !strings.Contains(html, "background-color") {
		t.Error("ThemeTransitionCSS should include background-color transition")
	}
}
