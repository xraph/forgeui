# Command Component

A fast, composable command menu component for search and quick actions, inspired by shadcn/ui's command palette with ⌘K interaction pattern.

## Features

- 🔍 **Fast search** - Real-time filtering with fuzzy matching
- ⌨️ **Keyboard navigation** - Full keyboard support (arrows, enter, escape, home, end)
- 🎯 **⌘K Dialog** - Command palette that opens with Cmd+K / Ctrl+K
- 🎨 **Shadcn styling** - Beautiful design matching shadcn/ui patterns
- 🧩 **Composable API** - Build complex command menus with simple components
- ♿ **Accessible** - Proper ARIA attributes and keyboard navigation
- 🎭 **Alpine.js powered** - Reactive filtering and selection

## Installation

The command component is part of ForgeUI. Import it in your Go code:

```go
import "github.com/xraph/forgeui/components/command"
```

## Basic Usage

### Simple Command Menu

```go
command.Command(
    command.CommandInput(),
    command.CommandList(
        command.CommandEmpty("No results found."),
        command.CommandGroup("Suggestions",
            command.CommandItem("Calendar"),
            command.CommandItem("Search Emoji"),
            command.CommandItem("Calculator"),
        ),
        command.CommandSeparator(),
        command.CommandGroup("Settings",
            command.CommandItem("Profile"),
            command.CommandItem("Billing"),
            command.CommandItem("Settings"),
        ),
    ),
)
```

### Command Dialog (⌘K Style)

```go
command.CommandDialog(
    button.Button(g.Text("Open ⌘K")),
    command.CommandInput(),
    command.CommandList(
        command.CommandEmpty("No results found."),
        command.CommandGroup("Actions",
            command.CommandItem("New File"),
            command.CommandItem("New Folder"),
            command.CommandItem("New Document"),
        ),
    ),
)
```

## Components

### Command

Main container component that initializes the command menu state and handles keyboard navigation.

```go
command.Command(children ...g.Node)
```

**Features:**
- Initializes Alpine.js state for search and selection
- Handles keyboard navigation logic
- Filters items based on search query
- Tracks visible items and selected index

### CommandDialog

Wraps the command menu in a dialog that can be triggered with ⌘K/Ctrl+K.

```go
command.CommandDialog(trigger g.Node, children ...g.Node)
```

**Parameters:**
- `trigger` - Element that opens the dialog when clicked
- `children` - CommandInput, CommandList, and other command components

**Features:**
- Opens with ⌘K (Mac) or Ctrl+K (Windows/Linux)
- Closes with Escape key
- Auto-focuses input when opened
- Backdrop blur effect
- Focus trap for accessibility

### CommandInput

Search input field with icon and optional clear button.

```go
command.CommandInput(opts ...CommandInputOption)
```

**Options:**
- `WithPlaceholder(text)` - Custom placeholder text
- `WithIcon(icon)` - Custom search icon
- `WithShowClear()` - Show clear button when input has value
- `WithInputClass(class)` - Additional CSS classes
- `WithInputAttrs(attrs...)` - Additional HTML attributes

**Example:**
```go
command.CommandInput(
    command.WithPlaceholder("Search commands..."),
    command.WithShowClear(),
)
```

### CommandList

Scrollable container for command items with keyboard navigation support.

```go
command.CommandList(children ...g.Node)
```

**Keyboard Navigation:**
- **↑/↓** - Navigate through items
- **Enter** - Select current item
- **Home** - Jump to first item
- **End** - Jump to last item
- **Escape** - Close dialog (when in dialog)

### CommandEmpty

Empty state displayed when no items match the search query.

```go
command.CommandEmpty(message string)
```

**Example:**
```go
command.CommandEmpty("No results found.")
```

### CommandGroup

Groups related command items with an optional heading.

```go
command.CommandGroup(heading string, children ...g.Node)
```

**Example:**
```go
command.CommandGroup("File Operations",
    command.CommandItem("New File"),
    command.CommandItem("Open File"),
    command.CommandItem("Save File"),
)
```

### CommandItem

Individual selectable command item.

```go
command.CommandItem(text string, opts ...CommandItemOption)
```

**Options:**
- `WithValue(value)` - Custom search value (defaults to text)
- `WithDisabled(bool)` - Disable the item
- `WithOnSelect(expr)` - Alpine.js expression to execute on selection
- `WithItemIcon(icon)` - Leading icon
- `WithShortcut(keys...)` - Keyboard shortcut display
- `WithItemClass(class)` - Additional CSS classes
- `WithItemAttrs(attrs...)` - Additional HTML attributes

**Example:**
```go
command.CommandItem("New File",
    command.WithItemIcon(icons.FilePlus()),
    command.WithShortcut("⌘", "N"),
    command.WithOnSelect("window.location = '/new'"),
)
```

### CommandSeparator

Visual separator between command groups.

```go
command.CommandSeparator()
```

### CommandShortcut

Displays keyboard shortcuts (typically used inside CommandItem).

```go
command.CommandShortcut(keys ...string)
```

**Example:**
```go
command.CommandShortcut("⌘", "K")  // ⌘K
command.CommandShortcut("Ctrl", "Shift", "P")  // Ctrl Shift P
```

## Advanced Examples

### With Icons and Shortcuts

```go
command.Command(
    command.CommandInput(command.WithPlaceholder("Type a command...")),
    command.CommandList(
        command.CommandEmpty("No commands found."),
        command.CommandGroup("File",
            command.CommandItem("New File",
                command.WithItemIcon(icons.FilePlus()),
                command.WithShortcut("⌘", "N"),
            ),
            command.CommandItem("Open File",
                command.WithItemIcon(icons.FolderOpen()),
                command.WithShortcut("⌘", "O"),
            ),
            command.CommandItem("Save",
                command.WithItemIcon(icons.Save()),
                command.WithShortcut("⌘", "S"),
            ),
        ),
        command.CommandSeparator(),
        command.CommandGroup("Edit",
            command.CommandItem("Copy",
                command.WithItemIcon(icons.Copy()),
                command.WithShortcut("⌘", "C"),
            ),
            command.CommandItem("Paste",
                command.WithItemIcon(icons.Clipboard()),
                command.WithShortcut("⌘", "V"),
            ),
        ),
    ),
)
```

### With Router Integration

```go
command.CommandDialog(
    button.Button(
        g.Text("Search"),
        button.WithVariant(forgeui.VariantOutline),
    ),
    command.CommandInput(),
    command.CommandList(
        command.CommandEmpty("No pages found."),
        command.CommandGroup("Pages",
            command.CommandItem("Dashboard",
                command.WithItemIcon(icons.LayoutDashboard()),
                command.WithOnSelect("window.location = '/dashboard'"),
            ),
            command.CommandItem("Users",
                command.WithItemIcon(icons.Users()),
                command.WithOnSelect("window.location = '/users'"),
            ),
            command.CommandItem("Settings",
                command.WithItemIcon(icons.Settings()),
                command.WithOnSelect("window.location = '/settings'"),
            ),
        ),
    ),
)
```

### With HTMX Actions

```go
command.Command(
    command.CommandInput(),
    command.CommandList(
        command.CommandEmpty("No actions available."),
        command.CommandGroup("Quick Actions",
            command.CommandItem("Create Task",
                command.WithItemIcon(icons.Plus()),
                command.WithItemAttrs(
                    htmx.HXPost("/api/tasks"),
                    htmx.HXSwap("afterbegin"),
                    htmx.HXTarget("#task-list"),
                ),
            ),
            command.CommandItem("Refresh Data",
                command.WithItemIcon(icons.RefreshCw()),
                command.WithItemAttrs(
                    htmx.HXGet("/api/data"),
                    htmx.HXSwap("innerHTML"),
                    htmx.HXTarget("#data-container"),
                ),
            ),
        ),
    ),
)
```

### Disabled Items

```go
command.CommandItem("Premium Feature",
    command.WithItemIcon(icons.Crown()),
    command.WithDisabled(true),
)
```

### Custom Search Values

```go
// Search for "delete" or "remove" will both match this item
command.CommandItem("Delete Item",
    command.WithValue("delete remove trash"),
    command.WithItemIcon(icons.Trash()),
)
```

## Styling

The command component uses CVA (Class Variance Authority) for consistent styling. All components follow shadcn/ui design patterns.

### Customization

You can add custom classes to any component:

```go
command.CommandInput(
    command.WithInputClass("border-2 border-primary"),
)

command.CommandItem("Custom",
    command.WithItemClass("text-primary font-bold"),
)
```

### Theme Variables

The component uses these CSS custom properties:

- `--popover` - Background color
- `--popover-foreground` - Text color
- `--border` - Border color
- `--accent` - Hover/selected background
- `--accent-foreground` - Hover/selected text
- `--muted-foreground` - Muted text color

## Accessibility

The command component is built with accessibility in mind:

- ✅ Proper ARIA roles and attributes
- ✅ Full keyboard navigation
- ✅ Focus management in dialog mode
- ✅ Screen reader support
- ✅ Focus trap when in dialog
- ✅ Visible focus indicators

## Browser Compatibility

Works in all modern browsers that support:
- Alpine.js v3+
- CSS Grid and Flexbox
- ES6+ JavaScript

## Related Components

- [Dialog](../modal/dialog.go) - Base dialog component
- [Input](../input/input.go) - Form input component
- [Separator](../separator/separator.go) - Visual separator
- [Icons](../../icons/icons.go) - Icon system

## License

Part of ForgeUI - MIT License
