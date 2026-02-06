// Package command provides a command menu component with search and keyboard navigation.
//
// The Command component is a fast, composable command menu for search and quick actions,
// inspired by shadcn/ui's command palette with ⌘K interaction pattern.
//
// Basic usage:
//
//	command.Command(
//	    command.CommandInput(),
//	    command.CommandList(
//	        command.CommandEmpty("No results found."),
//	        command.CommandGroup("Suggestions",
//	            command.CommandItem("Calendar"),
//	            command.CommandItem("Search Emoji"),
//	        ),
//	    ),
//	)
package command

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
	"github.com/xraph/forgeui/icons"
)

// CommandInputProps defines configuration for command input
type CommandInputProps struct {
	Placeholder string
	Icon        g.Node
	ShowClear   bool
	Class       string
	Attrs       []g.Node
}

// CommandInputOption is a functional option for configuring command input
type CommandInputOption func(*CommandInputProps)

// WithPlaceholder sets the input placeholder text
func WithPlaceholder(text string) CommandInputOption {
	return func(p *CommandInputProps) { p.Placeholder = text }
}

// WithIcon sets a custom icon for the input
func WithIcon(icon g.Node) CommandInputOption {
	return func(p *CommandInputProps) { p.Icon = icon }
}

// WithShowClear shows the clear button when input has value
func WithShowClear() CommandInputOption {
	return func(p *CommandInputProps) { p.ShowClear = true }
}

// WithInputClass adds custom classes to the input container
func WithInputClass(class string) CommandInputOption {
	return func(p *CommandInputProps) { p.Class = class }
}

// WithInputAttrs adds custom attributes to the input
func WithInputAttrs(attrs ...g.Node) CommandInputOption {
	return func(p *CommandInputProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// CommandItemProps defines configuration for command items
type CommandItemProps struct {
	Value    string
	Disabled bool
	OnSelect string
	Icon     g.Node
	Shortcut []string
	Class    string
	Attrs    []g.Node
}

// CommandItemOption is a functional option for configuring command items
type CommandItemOption func(*CommandItemProps)

// WithValue sets the search value for the item
func WithValue(value string) CommandItemOption {
	return func(p *CommandItemProps) { p.Value = value }
}

// WithDisabled disables the command item
func WithDisabled(disabled bool) CommandItemOption {
	return func(p *CommandItemProps) { p.Disabled = disabled }
}

// WithOnSelect sets the Alpine expression to execute on selection
func WithOnSelect(expr string) CommandItemOption {
	return func(p *CommandItemProps) { p.OnSelect = expr }
}

// WithItemIcon adds an icon to the command item
func WithItemIcon(icon g.Node) CommandItemOption {
	return func(p *CommandItemProps) { p.Icon = icon }
}

// WithShortcut adds keyboard shortcut display to the item
func WithShortcut(keys ...string) CommandItemOption {
	return func(p *CommandItemProps) { p.Shortcut = keys }
}

// WithItemClass adds custom classes to the item
func WithItemClass(class string) CommandItemOption {
	return func(p *CommandItemProps) { p.Class = class }
}

// WithItemAttrs adds custom attributes to the item
func WithItemAttrs(attrs ...g.Node) CommandItemOption {
	return func(p *CommandItemProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Command creates a command menu container with search and keyboard navigation.
//
// Example:
//
//	command.Command(
//	    command.CommandInput(),
//	    command.CommandList(
//	        command.CommandGroup("Actions",
//	            command.CommandItem("New File"),
//	        ),
//	    ),
//	)
func Command(children ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{
			"search":        "",
			"selectedIndex": 0,
			"visibleItems":  []any{},
			"filter": alpine.RawJS(`function(text, value) {
				if (!this.search) return true;
				const searchLower = this.search.toLowerCase();
				const textLower = (text || '').toLowerCase();
				const valueLower = (value || text || '').toLowerCase();
				return textLower.includes(searchLower) || valueLower.includes(searchLower);
			}`),
			"updateVisible": alpine.RawJS(`function() {
				const items = Array.from(this.$el.querySelectorAll('[data-command-item]'));
				this.visibleItems = items.filter(item => {
					const text = item.getAttribute('data-command-item-text');
					const value = item.getAttribute('data-command-item-value');
					return this.filter(text, value);
				});
				if (this.selectedIndex >= this.visibleItems.length) {
					this.selectedIndex = Math.max(0, this.visibleItems.length - 1);
				}
			}`),
		}),
		alpine.XInit("updateVisible()"),
		html.Class(commandCVA.Classes(map[string]string{})),
		g.Attr("data-command", ""),
		g.Attr("role", "application"),
		g.Group(children),
	)
}

// CommandInput creates a search input for the command menu.
//
// Example:
//
//	command.CommandInput(
//	    command.WithPlaceholder("Type a command..."),
//	)
func CommandInput(opts ...CommandInputOption) g.Node {
	props := &CommandInputProps{
		Placeholder: "Type a command or search...",
		ShowClear:   true,
	}

	for _, opt := range opts {
		opt(props)
	}

	// Use default search icon if not provided
	icon := props.Icon
	if icon == nil {
		icon = icons.Search(icons.WithSize(16))
	}

	classes := commandInputCVA.Classes(map[string]string{})
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),

		// Search icon
		html.Div(
			html.Class("shrink-0"),
			icon,
		),

		// Input field
		html.Input(
			html.Type("text"),
			html.Placeholder(props.Placeholder),
			alpine.XModel("search"),
			alpine.XOn("input", "updateVisible(); selectedIndex = 0"),
			g.Attr("autocomplete", "off"),
			g.Attr("autocorrect", "off"),
			g.Attr("spellcheck", "false"),
		),

		// Clear button
		g.If(props.ShowClear,
			html.Button(
				html.Type("button"),
				g.Attr("x-show", "search.length > 0"),
				alpine.XCloak(),
				g.Attr("@click", "search = ''; updateVisible(); $el.previousElementSibling.focus()"),
				html.Class("shrink-0 opacity-50 hover:opacity-100 transition-opacity"),
				g.Attr("aria-label", "Clear search"),
				icons.X(icons.WithSize(16)),
			),
		),
	)
}

// CommandList creates a scrollable container for command items.
//
// Example:
//
//	command.CommandList(
//	    command.CommandEmpty("No results."),
//	    command.CommandGroup("Items", ...),
//	)
func CommandList(children ...g.Node) g.Node {
	return html.Div(
		html.Class(commandListCVA.Classes(map[string]string{})),
		g.Attr("role", "listbox"),
		g.Attr("aria-label", "Commands"),
		alpine.XOn("keydown.arrow-down.prevent", `
			if (visibleItems.length > 0) {
				selectedIndex = Math.min(selectedIndex + 1, visibleItems.length - 1);
				visibleItems[selectedIndex]?.scrollIntoView({ block: 'nearest' });
			}
		`),
		alpine.XOn("keydown.arrow-up.prevent", `
			if (visibleItems.length > 0) {
				selectedIndex = Math.max(selectedIndex - 1, 0);
				visibleItems[selectedIndex]?.scrollIntoView({ block: 'nearest' });
			}
		`),
		alpine.XOn("keydown.home.prevent", `
			if (visibleItems.length > 0) {
				selectedIndex = 0;
				visibleItems[0]?.scrollIntoView({ block: 'nearest' });
			}
		`),
		alpine.XOn("keydown.end.prevent", `
			if (visibleItems.length > 0) {
				selectedIndex = visibleItems.length - 1;
				visibleItems[selectedIndex]?.scrollIntoView({ block: 'nearest' });
			}
		`),
		alpine.XOn("keydown.enter.prevent", `
			if (visibleItems.length > 0 && visibleItems[selectedIndex]) {
				visibleItems[selectedIndex].click();
			}
		`),
		g.Group(children),
	)
}

// CommandEmpty displays a message when no results are found.
//
// Example:
//
//	command.CommandEmpty("No results found.")
func CommandEmpty(message string) g.Node {
	if message == "" {
		message = "No results found."
	}

	return html.Div(
		g.Attr("x-show", "visibleItems.length === 0"),
		alpine.XCloak(),
		html.Class(commandEmptyCVA.Classes(map[string]string{})),
		g.Text(message),
	)
}

// CommandGroup creates a group of related command items with an optional heading.
//
// Example:
//
//	command.CommandGroup("Settings",
//	    command.CommandItem("Profile"),
//	    command.CommandItem("Billing"),
//	)
func CommandGroup(heading string, children ...g.Node) g.Node {
	return html.Div(
		html.Class(commandGroupCVA.Classes(map[string]string{})),
		g.Attr("role", "group"),
		g.If(heading != "", g.Attr("aria-label", heading)),
		
		// Add heading if provided
		g.If(heading != "",
			html.Div(
				g.Attr("data-command-group-heading", ""),
				html.Class("px-2 py-1.5 text-xs font-medium text-muted-foreground"),
				g.Text(heading),
			),
		),
		
		g.Group(children),
	)
}

// CommandItem creates a selectable command item.
//
// Example:
//
//	command.CommandItem("New File",
//	    command.WithItemIcon(icons.FilePlus()),
//	    command.WithShortcut("⌘", "N"),
//	    command.WithOnSelect("console.log('New file')"),
//	)
func CommandItem(text string, opts ...CommandItemOption) g.Node {
	props := &CommandItemProps{
		Value: text,
	}

	for _, opt := range opts {
		opt(props)
	}

	classes := commandItemCVA.Classes(map[string]string{})
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("data-command-item", ""),
		g.Attr("data-command-item-text", text),
		g.Attr("data-command-item-value", props.Value),
		g.Attr("role", "option"),
		g.Attr("x-show", "filter($el.getAttribute('data-command-item-text'), $el.getAttribute('data-command-item-value'))"),
		g.Attr("x-effect", "$el.setAttribute('data-selected', visibleItems[selectedIndex] === $el ? 'true' : 'false')"),
		g.Attr("@mouseenter", "$el.setAttribute('data-selected', 'true'); selectedIndex = visibleItems.indexOf($el)"),
		g.Attr("@mouseleave", "$el.setAttribute('data-selected', 'false')"),
	}

	// Add disabled state
	if props.Disabled {
		attrs = append(attrs,
			g.Attr("data-disabled", "true"),
			g.Attr("aria-disabled", "true"),
		)
	} else {
		// Add click handler if not disabled
		clickHandler := props.OnSelect
		if clickHandler == "" {
			clickHandler = "console.log('Selected:', $el.getAttribute('data-command-item-text'))"
		}
		attrs = append(attrs, g.Attr("@click", clickHandler))
	}

	attrs = append(attrs, props.Attrs...)

	// Build item content
	content := []g.Node{}

	// Add icon if provided
	if props.Icon != nil {
		content = append(content, props.Icon)
	}

	// Add text
	content = append(content, g.Text(text))

	// Add shortcut if provided
	if len(props.Shortcut) > 0 {
		content = append(content, CommandShortcut(props.Shortcut...))
	}

	return html.Div(
		g.Group(attrs),
		g.Group(content),
	)
}

// CommandSeparator creates a visual separator between command groups.
//
// Example:
//
//	command.CommandSeparator()
func CommandSeparator() g.Node {
	return html.Div(
		html.Class(commandSeparatorCVA.Classes(map[string]string{})),
		g.Attr("role", "separator"),
	)
}

// CommandShortcut displays keyboard shortcuts in a command item.
//
// Example:
//
//	command.CommandShortcut("⌘", "K")
func CommandShortcut(keys ...string) g.Node {
	keyNodes := make([]g.Node, len(keys))
	for i, key := range keys {
		keyNodes[i] = html.Span(g.Text(key))
	}

	return html.Span(
		html.Class(commandShortcutCVA.Classes(map[string]string{})),
		g.Group(keyNodes),
	)
}

// CommandDialog creates a command palette dialog that opens with ⌘K/Ctrl+K.
//
// Example:
//
//	command.CommandDialog(
//	    button.Button(g.Text("Open Command")),
//	    command.CommandInput(),
//	    command.CommandList(...),
//	)
func CommandDialog(trigger g.Node, children ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{
			"open": false,
		}),
		// Global keyboard shortcut listener
		g.Attr("@keydown.window", `
			if ((event.metaKey || event.ctrlKey) && event.key === 'k') {
				event.preventDefault();
				open = !open;
				if (open) {
					$nextTick(() => {
						$el.querySelector('[data-command] input')?.focus();
					});
				}
			}
		`),
		g.Attr("@keydown.escape.window", "open = false"),

		// Trigger button
		html.Div(
			g.Attr("@click", "open = true; $nextTick(() => { $el.parentElement.querySelector('[data-command] input')?.focus() })"),
			trigger,
		),

		// Dialog overlay
		html.Div(
			g.Attr("x-show", "open"),
			alpine.XCloak(),
			html.Class("fixed inset-0 z-50 overflow-y-auto"),
			g.Attr("aria-modal", "true"),
			g.Attr("role", "dialog"),

			// Backdrop
			html.Div(
				g.Group(alpine.XTransition(animation.FadeIn())),
				html.Class("fixed inset-0 bg-background/80 backdrop-blur-sm transition-all"),
				g.Attr("@click", "open = false"),
			),

			// Dialog content
			html.Div(
				html.Class("fixed inset-0 flex items-start justify-center pt-[20vh]"),

				// Command panel
				html.Div(
					g.Group(alpine.XTransition(animation.ScaleIn())),
					g.Attr("@click.stop", ""),
					g.Attr("x-trap.noscroll", "open"),
					html.Class(fmt.Sprintf("relative w-full max-w-[450px] %s", commandDialogCVA.Classes(map[string]string{}))),

					// Command container
					Command(children...),
				),
			),
		),
	)
}
