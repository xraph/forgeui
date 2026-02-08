package command

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/xraph/forgeui/icons"
)

func TestCommand(t *testing.T) {
	t.Run("renders command container", func(t *testing.T) {
		cmd := Command(
			CommandInput(),
			CommandList(),
		)

		var buf bytes.Buffer
		if err := cmd.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "data-command") {
			t.Error("expected data-command attribute")
		}

		if !strings.Contains(html, `role="application"`) {
			t.Error("expected role=application")
		}

		if !strings.Contains(html, "x-data") {
			t.Error("expected Alpine x-data")
		}

		if !strings.Contains(html, "flex-col") {
			t.Error("expected flex-col class from CVA")
		}

		if !strings.Contains(html, "bg-popover") {
			t.Error("expected bg-popover class from CVA")
		}
	})

	t.Run("initializes Alpine state", func(t *testing.T) {
		cmd := Command()

		var buf bytes.Buffer
		if err := cmd.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		// Check for presence of state properties (may be HTML-encoded)
		if !strings.Contains(html, "search") {
			t.Error("expected search state in Alpine data")
		}

		if !strings.Contains(html, "selectedIndex") {
			t.Error("expected selectedIndex state in Alpine data")
		}

		if !strings.Contains(html, "visibleItems") {
			t.Error("expected visibleItems state in Alpine data")
		}

		if !strings.Contains(html, "filter") {
			t.Error("expected filter function in Alpine data")
		}

		if !strings.Contains(html, "updateVisible") {
			t.Error("expected updateVisible function in Alpine data")
		}
	})
}

func TestCommandInput(t *testing.T) {
	t.Run("renders with default props", func(t *testing.T) {
		input := CommandInput()

		var buf bytes.Buffer
		if err := input.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `<input`) {
			t.Error("expected input element")
		}

		if !strings.Contains(html, `type="text"`) {
			t.Error("expected text input type")
		}

		if !strings.Contains(html, "Type a command or search...") {
			t.Error("expected default placeholder")
		}

		if !strings.Contains(html, "x-model=\"search\"") {
			t.Error("expected x-model binding")
		}

		if !strings.Contains(html, "border-b") {
			t.Error("expected border-b class from CVA")
		}
	})

	t.Run("renders with custom placeholder", func(t *testing.T) {
		input := CommandInput(
			WithPlaceholder("Search commands..."),
		)

		var buf bytes.Buffer
		if err := input.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "Search commands...") {
			t.Error("expected custom placeholder")
		}
	})

	t.Run("renders search icon", func(t *testing.T) {
		input := CommandInput()

		var buf bytes.Buffer
		if err := input.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "<svg") {
			t.Error("expected SVG icon")
		}
	})

	t.Run("renders clear button", func(t *testing.T) {
		input := CommandInput(WithShowClear())

		var buf bytes.Buffer
		if err := input.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "x-show") && !strings.Contains(html, "search.length") {
			t.Error("expected x-show binding on clear button")
		}

		if !strings.Contains(html, `aria-label="Clear search"`) {
			t.Error("expected aria-label on clear button")
		}
	})

	t.Run("renders with custom icon", func(t *testing.T) {
		customIcon := icons.Command()
		input := CommandInput(
			WithIcon(customIcon),
		)

		var buf bytes.Buffer
		if err := input.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "<svg") {
			t.Error("expected custom icon")
		}
	})

	t.Run("applies custom class", func(t *testing.T) {
		input := CommandInput(
			WithInputClass("custom-class"),
		)

		var buf bytes.Buffer
		if err := input.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "custom-class") {
			t.Error("expected custom class")
		}
	})
}

func TestCommandList(t *testing.T) {
	t.Run("renders command list", func(t *testing.T) {
		list := CommandList(
			CommandItem("Item 1"),
			CommandItem("Item 2"),
		)

		var buf bytes.Buffer
		if err := list.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `role="listbox"`) {
			t.Error("expected role=listbox")
		}

		if !strings.Contains(html, `aria-label="Commands"`) {
			t.Error("expected aria-label")
		}

		if !strings.Contains(html, "max-h-[300px]") {
			t.Error("expected max-height class from CVA")
		}

		if !strings.Contains(html, "overflow-y-auto") {
			t.Error("expected overflow-y-auto class")
		}
	})

	t.Run("has keyboard navigation", func(t *testing.T) {
		list := CommandList()

		var buf bytes.Buffer
		if err := list.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "keydown.arrow-down.prevent") {
			t.Error("expected arrow-down handler")
		}

		if !strings.Contains(html, "keydown.arrow-up.prevent") {
			t.Error("expected arrow-up handler")
		}

		if !strings.Contains(html, "keydown.home.prevent") {
			t.Error("expected home key handler")
		}

		if !strings.Contains(html, "keydown.end.prevent") {
			t.Error("expected end key handler")
		}

		if !strings.Contains(html, "keydown.enter.prevent") {
			t.Error("expected enter key handler")
		}
	})
}

func TestCommandEmpty(t *testing.T) {
	t.Run("renders with default message", func(t *testing.T) {
		empty := CommandEmpty("")

		var buf bytes.Buffer
		if err := empty.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "No results found.") {
			t.Error("expected default message")
		}

		if !strings.Contains(html, "x-show=\"visibleItems.length === 0\"") {
			t.Error("expected x-show binding")
		}

		if !strings.Contains(html, "text-muted-foreground") {
			t.Error("expected muted text class")
		}
	})

	t.Run("renders with custom message", func(t *testing.T) {
		empty := CommandEmpty("Nothing here!")

		var buf bytes.Buffer
		if err := empty.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "Nothing here!") {
			t.Error("expected custom message")
		}
	})
}

func TestCommandGroup(t *testing.T) {
	t.Run("renders group with heading", func(t *testing.T) {
		group := CommandGroup("Settings",
			CommandItem("Profile"),
			CommandItem("Billing"),
		)

		var buf bytes.Buffer
		if err := group.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "Settings") {
			t.Error("expected group heading")
		}

		if !strings.Contains(html, `role="group"`) {
			t.Error("expected role=group")
		}

		if !strings.Contains(html, `aria-label="Settings"`) {
			t.Error("expected aria-label with heading")
		}

		if !strings.Contains(html, "data-command-group-heading") {
			t.Error("expected data-command-group-heading attribute")
		}

		if !strings.Contains(html, "Profile") {
			t.Error("expected child item")
		}
	})

	t.Run("renders group without heading", func(t *testing.T) {
		group := CommandGroup("",
			CommandItem("Item 1"),
		)

		var buf bytes.Buffer
		if err := group.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `role="group"`) {
			t.Error("expected role=group")
		}

		// The data-command-group-heading div still exists but should be empty
		// This is OK since g.If renders empty content when condition is false
		if !strings.Contains(html, "Item 1") {
			t.Error("expected child items to render")
		}
	})
}

func TestCommandItem(t *testing.T) {
	t.Run("renders command item", func(t *testing.T) {
		item := CommandItem("Profile")

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "Profile") {
			t.Error("expected item text")
		}

		if !strings.Contains(html, "data-command-item") {
			t.Error("expected data-command-item attribute")
		}

		if !strings.Contains(html, "data-command-item-text=\"Profile\"") {
			t.Error("expected data-command-item-text attribute")
		}

		if !strings.Contains(html, "data-command-item-value=\"Profile\"") {
			t.Error("expected data-command-item-value attribute with default value")
		}

		if !strings.Contains(html, `role="option"`) {
			t.Error("expected role=option")
		}

		if !strings.Contains(html, "rounded-sm") {
			t.Error("expected CVA classes")
		}
	})

	t.Run("renders with custom value", func(t *testing.T) {
		item := CommandItem("User Profile", WithValue("profile-settings"))

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "data-command-item-value=\"profile-settings\"") {
			t.Error("expected custom value attribute")
		}
	})

	t.Run("renders with icon", func(t *testing.T) {
		item := CommandItem("Settings",
			WithItemIcon(icons.Settings()),
		)

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "<svg") {
			t.Error("expected SVG icon")
		}
	})

	t.Run("renders with shortcut", func(t *testing.T) {
		item := CommandItem("Settings",
			WithShortcut("⌘", "S"),
		)

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "⌘") {
			t.Error("expected command key symbol")
		}

		if !strings.Contains(html, "S") {
			t.Error("expected S key")
		}

		if !strings.Contains(html, "ml-auto") {
			t.Error("expected ml-auto class for shortcut")
		}
	})

	t.Run("renders disabled item", func(t *testing.T) {
		item := CommandItem("Disabled", WithDisabled(true))

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "data-disabled=\"true\"") {
			t.Error("expected data-disabled attribute")
		}

		if !strings.Contains(html, `aria-disabled="true"`) {
			t.Error("expected aria-disabled attribute")
		}
	})

	t.Run("renders with onSelect handler", func(t *testing.T) {
		item := CommandItem("Action",
			WithOnSelect("alert('clicked')"),
		)

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		// Check for HTML-encoded or regular version
		if !strings.Contains(html, "alert") || !strings.Contains(html, "clicked") {
			t.Error("expected custom onSelect handler")
		}
	})

	t.Run("applies custom class", func(t *testing.T) {
		item := CommandItem("Custom",
			WithItemClass("custom-item-class"),
		)

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "custom-item-class") {
			t.Error("expected custom class")
		}
	})

	t.Run("has Alpine bindings", func(t *testing.T) {
		item := CommandItem("Test")

		var buf bytes.Buffer
		if err := item.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "x-show") {
			t.Error("expected x-show binding for filtering")
		}

		if !strings.Contains(html, "x-effect") {
			t.Error("expected x-effect for selected state")
		}

		if !strings.Contains(html, "@mouseenter") && !strings.Contains(html, "mouseenter") {
			t.Error("expected mouseenter handler")
		}

		if !strings.Contains(html, "@mouseleave") && !strings.Contains(html, "mouseleave") {
			t.Error("expected mouseleave handler")
		}
	})
}

func TestCommandSeparator(t *testing.T) {
	t.Run("renders separator", func(t *testing.T) {
		sep := CommandSeparator()

		var buf bytes.Buffer
		if err := sep.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `role="separator"`) {
			t.Error("expected role=separator")
		}

		if !strings.Contains(html, "bg-border") {
			t.Error("expected bg-border class")
		}

		if !strings.Contains(html, "h-px") {
			t.Error("expected h-px class")
		}
	})
}

func TestCommandShortcut(t *testing.T) {
	t.Run("renders single key", func(t *testing.T) {
		shortcut := CommandShortcut("K")

		var buf bytes.Buffer
		if err := shortcut.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "K") {
			t.Error("expected key text")
		}

		if !strings.Contains(html, "text-xs") {
			t.Error("expected text-xs class")
		}

		if !strings.Contains(html, "text-muted-foreground") {
			t.Error("expected muted foreground color")
		}
	})

	t.Run("renders multiple keys", func(t *testing.T) {
		shortcut := CommandShortcut("⌘", "Shift", "K")

		var buf bytes.Buffer
		if err := shortcut.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "⌘") {
			t.Error("expected command symbol")
		}

		if !strings.Contains(html, "Shift") {
			t.Error("expected Shift key")
		}

		if !strings.Contains(html, "K") {
			t.Error("expected K key")
		}
	})
}

func TestCommandDialog(t *testing.T) {
	t.Run("renders dialog with trigger", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Open")),
			CommandInput(),
			CommandList(),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "Open") {
			t.Error("expected trigger text")
		}

		if !strings.Contains(html, "x-data") {
			t.Error("expected Alpine x-data")
		}

		if !strings.Contains(html, `"open"`) {
			t.Error("expected open state")
		}
	})

	t.Run("has keyboard shortcut binding", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Trigger")),
			CommandInput(),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "keydown.window") {
			t.Error("expected global keydown handler")
		}

		if !strings.Contains(html, "metaKey") {
			t.Error("expected metaKey check for ⌘K")
		}

		if !strings.Contains(html, "ctrlKey") {
			t.Error("expected ctrlKey check for Ctrl+K")
		}

		// May be HTML-encoded as event.key === &#39;k&#39;
		if !strings.Contains(html, "event.key") || !strings.Contains(html, "k") {
			t.Error("expected k key check")
		}
	})

	t.Run("has escape handler", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Trigger")),
			CommandInput(),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "keydown.escape.window") {
			t.Error("expected escape key handler")
		}
	})

	t.Run("has backdrop and overlay", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Trigger")),
			CommandInput(),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "backdrop-blur-sm") {
			t.Error("expected backdrop blur")
		}

		if !strings.Contains(html, "fixed inset-0") {
			t.Error("expected fixed inset overlay")
		}

		if !strings.Contains(html, `role="dialog"`) {
			t.Error("expected role=dialog")
		}

		if !strings.Contains(html, `aria-modal="true"`) {
			t.Error("expected aria-modal=true")
		}
	})

	t.Run("auto-focuses input on open", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Trigger")),
			CommandInput(),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "$nextTick") {
			t.Error("expected $nextTick for focus")
		}

		if !strings.Contains(html, "querySelector") || !strings.Contains(html, "input") {
			t.Error("expected input selector")
		}

		if !strings.Contains(html, "focus()") {
			t.Error("expected focus() call")
		}
	})

	t.Run("has focus trap", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Trigger")),
			CommandInput(),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "x-trap.noscroll") {
			t.Error("expected focus trap with noscroll")
		}
	})

	t.Run("centers dialog with max-width", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Trigger")),
			CommandInput(),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "max-w-[450px]") {
			t.Error("expected max-width constraint")
		}

		if !strings.Contains(html, "justify-center") {
			t.Error("expected centered layout")
		}
	})
}

func TestCommandIntegration(t *testing.T) {
	t.Run("renders complete command menu", func(t *testing.T) {
		cmd := Command(
			CommandInput(WithPlaceholder("Search...")),
			CommandList(
				CommandEmpty("No results."),
				CommandGroup("Suggestions",
					CommandItem("Calendar"),
					CommandItem("Search Emoji"),
				),
				CommandSeparator(),
				CommandGroup("Settings",
					CommandItem("Profile", WithShortcut("⌘", "P")),
					CommandItem("Settings", WithShortcut("⌘", "S")),
				),
			),
		)

		var buf bytes.Buffer
		if err := cmd.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()

		// Check structure
		if !strings.Contains(html, "data-command") {
			t.Error("expected command container")
		}

		if !strings.Contains(html, "Search...") {
			t.Error("expected input placeholder")
		}

		if !strings.Contains(html, "No results.") {
			t.Error("expected empty state")
		}

		if !strings.Contains(html, "Suggestions") {
			t.Error("expected first group heading")
		}

		if !strings.Contains(html, "Calendar") {
			t.Error("expected first item")
		}

		if !strings.Contains(html, "Settings") {
			t.Error("expected second group heading")
		}

		if !strings.Contains(html, "⌘") {
			t.Error("expected keyboard shortcuts")
		}

		if !strings.Contains(html, `role="separator"`) {
			t.Error("expected separator")
		}
	})

	t.Run("renders command dialog with complete structure", func(t *testing.T) {
		dialog := CommandDialog(
			g.El("button", g.Text("Open ⌘K")),
			CommandInput(),
			CommandList(
				CommandEmpty("No results found."),
				CommandGroup("Actions",
					CommandItem("New File", WithItemIcon(icons.FilePlus())),
					CommandItem("New Folder", WithItemIcon(icons.FolderPlus())),
				),
			),
		)

		var buf bytes.Buffer
		if err := dialog.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()

		// Check dialog structure
		if !strings.Contains(html, "Open ⌘K") {
			t.Error("expected trigger button")
		}

		if !strings.Contains(html, "keydown.window") {
			t.Error("expected keyboard shortcut handler")
		}

		if !strings.Contains(html, "backdrop-blur-sm") {
			t.Error("expected backdrop")
		}

		if !strings.Contains(html, "data-command") {
			t.Error("expected command inside dialog")
		}

		if !strings.Contains(html, "Actions") {
			t.Error("expected group heading")
		}

		if !strings.Contains(html, "New File") {
			t.Error("expected items")
		}
	})
}
