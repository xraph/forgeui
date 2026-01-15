package toast

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui"
)

func TestToast(t *testing.T) {
	t.Run("renders with default props", func(t *testing.T) {
		toast := Toast(ToastProps{
			Title: "Test notification",
		})

		var buf bytes.Buffer
		if err := toast.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "Test notification") {
			t.Error("expected toast to contain title")
		}

		if !strings.Contains(html, "role=\"alert\"") {
			t.Error("expected toast to have alert role")
		}
	})

	t.Run("renders with description", func(t *testing.T) {
		toast := Toast(ToastProps{
			Title:       "Success",
			Description: "Operation completed",
		})

		var buf bytes.Buffer
		if err := toast.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "Operation completed") {
			t.Error("expected toast to contain description")
		}
	})

	t.Run("renders success variant", func(t *testing.T) {
		toast := Toast(ToastProps{
			Title:   "Success",
			Variant: "success",
		})

		var buf bytes.Buffer
		if err := toast.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "bg-card") {
			t.Error("expected success variant styling")
		}
	})

	t.Run("renders error variant", func(t *testing.T) {
		toast := Toast(ToastProps{
			Title:   "Error",
			Variant: forgeui.VariantDestructive,
		})

		var buf bytes.Buffer
		if err := toast.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "bg-destructive") {
			t.Error("expected error variant styling")
		}
	})

	t.Run("renders with progress bar", func(t *testing.T) {
		toast := Toast(ToastProps{
			Title:        "Loading",
			Duration:     3000,
			ShowProgress: true,
		})

		var buf bytes.Buffer
		if err := toast.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "progress") {
			t.Error("expected toast to show progress bar")
		}
	})

	t.Run("renders with action button", func(t *testing.T) {
		toast := Toast(ToastProps{
			Title: "Update available",
			Action: &ToastAction{
				Label:   "Update now",
				OnClick: "updateApp()",
			},
		})

		var buf bytes.Buffer
		if err := toast.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "Update now") {
			t.Error("expected toast to contain action button")
		}
	})

	t.Run("renders close button", func(t *testing.T) {
		toast := Toast(ToastProps{
			Title: "Dismissable",
		})

		var buf bytes.Buffer
		if err := toast.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "aria-label=\"Close\"") {
			t.Error("expected toast to have close button")
		}

		if !strings.Contains(html, "@click") || !strings.Contains(html, "close()") {
			t.Error("expected close button to have click handler")
		}
	})
}

func TestToaster(t *testing.T) {
	t.Run("renders toaster container", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "fixed") {
			t.Error("expected toaster to be fixed positioned")
		}

		if !strings.Contains(html, "z-[100]") {
			t.Error("expected toaster to have high z-index")
		}

		if !strings.Contains(html, "aria-live=\"polite\"") {
			t.Error("expected toaster to have aria-live attribute")
		}
	})

	t.Run("renders with default bottom-right position", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "bottom-0") || !strings.Contains(html, "right-0") {
			t.Error("expected default bottom-right position")
		}
	})

	t.Run("renders with custom position", func(t *testing.T) {
		toaster := Toaster(WithPosition("top-left"))

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "top-0") || !strings.Contains(html, "left-0") {
			t.Error("expected top-left position")
		}
	})

	t.Run("supports max toasts limit", func(t *testing.T) {
		toaster := Toaster(WithMaxToasts(3))

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "maxToasts: 3") {
			t.Error("expected maxToasts configuration")
		}

		if !strings.Contains(html, "visibleToasts") {
			t.Error("expected visibleToasts computed property")
		}
	})

	t.Run("renders with escape key handler", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "@keydown.escape.window") {
			t.Error("expected escape key handler")
		}

		if !strings.Contains(html, "$store.toast.clear()") {
			t.Error("expected clear() call on escape")
		}
	})

	t.Run("renders toast list template", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "x-for") {
			t.Error("expected x-for loop for toasts")
		}

		if !strings.Contains(html, "visibleToasts") {
			t.Error("expected loop over visibleToasts")
		}
	})

	t.Run("includes toast lifecycle management", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "init()") {
			t.Error("expected init lifecycle hook")
		}

		if !strings.Contains(html, "close()") {
			t.Error("expected close method")
		}

		if !strings.Contains(html, "$store.toast.remove") {
			t.Error("expected toast removal logic")
		}
	})

	t.Run("renders with smooth transitions", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "x-show") {
			t.Error("expected x-show for transitions")
		}

		if !strings.Contains(html, "transition") {
			t.Error("expected transition classes")
		}

		if !strings.Contains(html, "scale") {
			t.Error("expected scale animation")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		toaster := Toaster(WithToasterClass("custom-toaster"))

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "custom-toaster") {
			t.Error("expected custom toaster class")
		}
	})
}

func TestRegisterToastStore(t *testing.T) {
	t.Run("registers Alpine store", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "Alpine.store") {
			t.Error("expected Alpine.store registration")
		}

		if !strings.Contains(html, "toast") {
			t.Error("expected toast store name")
		}
	})

	t.Run("includes items array", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "items") {
			t.Error("expected items state property")
		}
	})

	t.Run("includes add method", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "add(") {
			t.Error("expected add method")
		}

		if !strings.Contains(html, "this.items.push") {
			t.Error("expected push to items array")
		}
	})

	t.Run("includes remove method", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "remove(") {
			t.Error("expected remove method")
		}

		if !strings.Contains(html, "filter") {
			t.Error("expected filter operation for removal")
		}
	})

	t.Run("includes clear method", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "clear()") {
			t.Error("expected clear method")
		}
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "Date.now()") && !strings.Contains(html, "Math.random()") {
			t.Error("expected unique ID generation")
		}
	})

	t.Run("sets default duration", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "duration: toast.duration || 5000") {
			t.Error("expected default duration of 5000ms")
		}
	})
}

func TestToastHelpers(t *testing.T) {
	t.Run("ToastSuccess generates success toast", func(t *testing.T) {
		expr := ToastSuccess("Saved successfully")

		if !strings.Contains(expr, "Saved successfully") {
			t.Error("expected success message")
		}

		if !strings.Contains(expr, "variant: 'success'") {
			t.Error("expected success variant")
		}

		if !strings.Contains(expr, "$store.toast.add") {
			t.Error("expected store.toast.add call")
		}
	})

	t.Run("ToastError generates error toast", func(t *testing.T) {
		expr := ToastError("Something went wrong")

		if !strings.Contains(expr, "Something went wrong") {
			t.Error("expected error message")
		}

		if !strings.Contains(expr, "variant: 'error'") {
			t.Error("expected error variant")
		}
	})

	t.Run("ToastWarning generates warning toast", func(t *testing.T) {
		expr := ToastWarning("Be careful")

		if !strings.Contains(expr, "Be careful") {
			t.Error("expected warning message")
		}

		if !strings.Contains(expr, "variant: 'warning'") {
			t.Error("expected warning variant")
		}
	})

	t.Run("ToastInfo generates info toast", func(t *testing.T) {
		expr := ToastInfo("FYI")

		if !strings.Contains(expr, "FYI") {
			t.Error("expected info message")
		}

		if !strings.Contains(expr, "variant: 'default'") {
			t.Error("expected default variant")
		}
	})
}

func TestToastVariantClasses(t *testing.T) {
	tests := []struct {
		variant forgeui.Variant
		want    string
	}{
		{"success", "bg-card"},
		{"error", "bg-destructive"},
		{forgeui.VariantDestructive, "bg-destructive"},
		{"warning", "bg-card"},
		{forgeui.VariantDefault, "bg-card"},
	}

	for _, tt := range tests {
		t.Run(string(tt.variant), func(t *testing.T) {
			classes := getToastVariantClasses(tt.variant)
			if !strings.Contains(classes, tt.want) {
				t.Errorf("getToastVariantClasses(%s) = %s, want to contain %s", tt.variant, classes, tt.want)
			}
		})
	}
}

func TestToasterPositionClasses(t *testing.T) {
	tests := []struct {
		position string
		wantTop  bool
		wantLeft bool
	}{
		{"top-left", true, true},
		{"top-right", true, false},
		{"top-center", true, false},
		{"bottom-left", false, true},
		{"bottom-right", false, false},
		{"bottom-center", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.position, func(t *testing.T) {
			classes := getToasterPositionClass(tt.position)

			hasTop := strings.Contains(classes, "top-0")
			hasBottom := strings.Contains(classes, "bottom-0")
			hasLeft := strings.Contains(classes, "left-0")
			hasRight := strings.Contains(classes, "right-0")

			if tt.wantTop && !hasTop {
				t.Errorf("expected top-0 for position %s", tt.position)
			}

			if !tt.wantTop && !hasBottom {
				t.Errorf("expected bottom-0 for position %s", tt.position)
			}

			if tt.wantLeft && !hasLeft {
				t.Errorf("expected left-0 for position %s", tt.position)
			}

			if !tt.wantLeft && !tt.wantTop && !hasRight && !strings.Contains(classes, "left-1/2") {
				// Should have right-0 or be centered
				if !strings.Contains(classes, "left-1/2") {
					t.Errorf("expected right-0 or centering for position %s", tt.position)
				}
			}
		})
	}
}

func TestMultiToastSupport(t *testing.T) {
	t.Run("toaster supports multiple toasts via x-for", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		// Verify it uses x-for for multiple toasts
		if !strings.Contains(html, "x-for") {
			t.Error("expected x-for directive for multiple toasts")
		}

		// Verify it iterates over items
		if !strings.Contains(html, "visibleToasts") {
			t.Error("expected iteration over visibleToasts")
		}

		// Verify each toast has unique key
		if !strings.Contains(html, ":key") {
			t.Error("expected :key for unique toast identification")
		}
	})

	t.Run("toaster respects max toasts limit", func(t *testing.T) {
		toaster := Toaster(WithMaxToasts(5))

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		if !strings.Contains(html, "maxToasts: 5") {
			t.Error("expected maxToasts limit of 5")
		}

		if !strings.Contains(html, "slice(0, this.maxToasts)") {
			t.Error("expected slice to enforce max toasts limit")
		}
	})

	t.Run("store supports adding multiple toasts", func(t *testing.T) {
		store := RegisterToastStore()

		var buf bytes.Buffer
		if err := store.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		// Verify items state exists
		if !strings.Contains(html, "items") {
			t.Error("expected items state property")
		}

		// Verify push adds to array
		if !strings.Contains(html, ".push(") {
			t.Error("expected push to add toasts to array")
		}
	})

	t.Run("each toast has independent lifecycle", func(t *testing.T) {
		toaster := Toaster()

		var buf bytes.Buffer
		if err := toaster.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		html := buf.String()

		// Each toast should have its own x-data
		if !strings.Contains(html, "x-data") {
			t.Error("expected x-data for independent toast state")
		}

		// Each toast manages its own show state
		if !strings.Contains(html, "show: false") {
			t.Error("expected independent show state for each toast")
		}

		// Each toast has its own timer
		if !strings.Contains(html, "timer: null") {
			t.Error("expected independent timer for each toast")
		}
	})
}
