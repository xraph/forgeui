package primitives

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// Helper function to render a component to string
func renderComponent(comp templ.Component) string {
	var buf bytes.Buffer
	_ = comp.Render(context.Background(), &buf)

	return buf.String()
}

// Helper function to check if string contains substring
func assertContains(t *testing.T, str, substr string) {
	t.Helper()

	if !strings.Contains(str, substr) {
		t.Errorf("Expected string to contain %q, got: %s", substr, str)
	}
}

func TestProvider(t *testing.T) {
	t.Run("creates provider with name and state", func(t *testing.T) {
		provider := Provider(
			WithProviderName("test"),
			WithProviderState(map[string]any{
				"count": 0,
				"open":  false,
			}),
		)

		rendered := renderComponent(provider)

		// Should have data-provider attribute
		assertContains(t, rendered, `data-provider="test"`)

		// Should have x-data with state (HTML entities are expected)
		assertContains(t, rendered, `x-data=`)
		assertContains(t, rendered, `count`)
		assertContains(t, rendered, `open`)
	})

	t.Run("creates provider with methods", func(t *testing.T) {
		provider := Provider(
			WithProviderName("counter"),
			WithProviderState(map[string]any{"count": 0}),
			WithProviderMethods(`
				increment() {
					this.count++;
				}
			`),
		)

		rendered := renderComponent(provider)

		assertContains(t, rendered, `data-provider="counter"`)
		assertContains(t, rendered, `increment()`)
		assertContains(t, rendered, `this.count++`)
	})

	t.Run("creates provider with init code", func(t *testing.T) {
		provider := Provider(
			WithProviderName("test"),
			WithProviderState(map[string]any{"ready": false}),
			WithProviderInit("this.ready = true"),
		)

		rendered := renderComponent(provider)

		assertContains(t, rendered, `x-init=`)
		assertContains(t, rendered, `this.ready = true`)
	})

	t.Run("creates provider with custom class", func(t *testing.T) {
		provider := Provider(
			WithProviderName("test"),
			WithProviderClass("custom-class"),
		)

		rendered := renderComponent(provider)

		assertContains(t, rendered, `class="custom-class"`)
	})

	t.Run("creates provider with children", func(t *testing.T) {
		provider := Provider(
			WithProviderName("test"),
			WithProviderChildren(
				templ.Raw("<div>child content</div>"),
			),
		)

		rendered := renderComponent(provider)

		assertContains(t, rendered, "child content")
	})

	t.Run("creates provider with debug mode", func(t *testing.T) {
		provider := Provider(
			WithProviderName("test"),
			WithProviderState(map[string]any{"value": 1}),
			WithProviderDebug(true),
		)

		rendered := renderComponent(provider)

		assertContains(t, rendered, `x-init=`)
		assertContains(t, rendered, `console.log`)
		assertContains(t, rendered, `[Provider:test]`)
	})

	t.Run("creates provider with lifecycle hooks", func(t *testing.T) {
		provider := Provider(
			WithProviderName("test"),
			WithProviderHook("onMount", "console.log('mounted')"),
			WithProviderHook("onUpdate", "console.log('updated')"),
		)

		rendered := renderComponent(provider)

		assertContains(t, rendered, `x-init=`)
		assertContains(t, rendered, `console.log`)
		assertContains(t, rendered, `mounted`)
		assertContains(t, rendered, `$watch`)
		assertContains(t, rendered, `updated`)
	})

	t.Run("returns children without wrapper when no name", func(t *testing.T) {
		provider := Provider(
			WithProviderChildren(
				templ.Raw("<div>content</div>"),
			),
		)

		if provider == nil {
			t.Error("Expected provider to be non-nil")
		}
	})
}

func TestProviderValue(t *testing.T) {
	t.Run("generates correct Alpine expression", func(t *testing.T) {
		expr := ProviderValue("sidebar", "collapsed")

		expected := `$el.closest('[data-provider="sidebar"]').__x.$data.collapsed`
		if expr != expected {
			t.Errorf("Expected %q, got %q", expected, expr)
		}
	})

	t.Run("works with nested properties", func(t *testing.T) {
		expr := ProviderValue("theme", "colors.primary")

		assertContains(t, expr, `data-provider="theme"`)
		assertContains(t, expr, `colors.primary`)
	})
}

func TestProviderMethod(t *testing.T) {
	t.Run("generates method call without args", func(t *testing.T) {
		expr := ProviderMethod("sidebar", "toggle")

		expected := `$el.closest('[data-provider="sidebar"]').__x.$data.toggle()`
		if expr != expected {
			t.Errorf("Expected %q, got %q", expected, expr)
		}
	})

	t.Run("generates method call with args", func(t *testing.T) {
		expr := ProviderMethod("form", "setValue", "name", "'John'")

		assertContains(t, expr, `data-provider="form"`)
		assertContains(t, expr, `setValue(name, 'John')`)
	})
}

func TestProviderDispatch(t *testing.T) {
	t.Run("generates dispatch expression", func(t *testing.T) {
		expr := ProviderDispatch("sidebar", "toggled", "{ open: true }")

		expected := `$dispatch('provider:sidebar:toggled', { open: true })`
		if expr != expected {
			t.Errorf("Expected %q, got %q", expected, expr)
		}
	})

	t.Run("uses empty object when no data", func(t *testing.T) {
		expr := ProviderDispatch("test", "event", "")

		assertContains(t, expr, `$dispatch('provider:test:event', {})`)
	})
}

func TestProviderScriptUtilities(t *testing.T) {
	t.Run("generates script with utilities", func(t *testing.T) {
		script := ProviderScriptUtilities()

		rendered := renderComponent(script)

		// Should be a script tag
		assertContains(t, rendered, "<script>")
		assertContains(t, rendered, "</script>")

		// Should have utility functions
		assertContains(t, rendered, "window.forgeui")
		assertContains(t, rendered, "getProvider")
		assertContains(t, rendered, "getValue")
		assertContains(t, rendered, "call")
		assertContains(t, rendered, "dispatch")

		// Should register Alpine magic property
		assertContains(t, rendered, "Alpine.magic('provider'")
	})
}

func TestProviderStack(t *testing.T) {
	t.Run("creates nested providers", func(t *testing.T) {
		stack := ProviderStack(
			Provider(
				WithProviderName("theme"),
				WithProviderState(map[string]any{"mode": "light"}),
			),
			Provider(
				WithProviderName("sidebar"),
				WithProviderState(map[string]any{"open": true}),
			),
		)

		rendered := renderComponent(stack)

		// Should have both providers
		assertContains(t, rendered, `data-provider="theme"`)
		assertContains(t, rendered, `data-provider="sidebar"`)
	})

	t.Run("returns empty group when no providers", func(t *testing.T) {
		stack := ProviderStack()

		// Just verify it doesn't panic
		if stack == nil {
			t.Error("Expected stack to be non-nil")
		}
	})
}

func TestProviderIntegration(t *testing.T) {
	t.Run("complete provider with state, methods, and children", func(t *testing.T) {
		provider := Provider(
			WithProviderName("counter"),
			WithProviderState(map[string]any{
				"count": 0,
				"step":  1,
			}),
			WithProviderMethods(`
				increment() {
					this.count += this.step;
					this.$dispatch('counter:changed', { count: this.count });
				},
				decrement() {
					this.count -= this.step;
					this.$dispatch('counter:changed', { count: this.count });
				},
				reset() {
					this.count = 0;
				}
			`),
			WithProviderInit("console.log('Counter initialized')"),
			WithProviderClass("counter-provider"),
			WithProviderChildren(
				templ.Raw("<button>Increment</button>"),
				templ.Raw("<span>Count</span>"),
			),
		)

		rendered := renderComponent(provider)

		// Verify all parts are present
		assertContains(t, rendered, `data-provider="counter"`)
		assertContains(t, rendered, `class="counter-provider"`)
		assertContains(t, rendered, `count`)
		assertContains(t, rendered, `step`)
		assertContains(t, rendered, `increment()`)
		assertContains(t, rendered, `decrement()`)
		assertContains(t, rendered, `reset()`)
		assertContains(t, rendered, `Counter initialized`)
		assertContains(t, rendered, "Increment")
		assertContains(t, rendered, "Count")
	})
}
