package command

import "github.com/xraph/forgeui"

// commandCVA defines the class variance authority configuration for command containers
// Following shadcn/ui command styling patterns
var commandCVA = forgeui.NewCVA(
	"flex",
	"flex-col",
	"w-full",
	"rounded-lg",
	"border",
	"border-input",
	"bg-popover",
	"text-popover-foreground",
	"overflow-hidden",
)

// commandInputCVA defines styling for command input field
var commandInputCVA = forgeui.NewCVA(
	"flex",
	"items-center",
	"border-b",
	"border-border",
	"px-3",
	"h-12",
	"text-sm",
	"[&_input]:flex-1",
	"[&_input]:bg-transparent",
	"[&_input]:py-3",
	"[&_input]:outline-none",
	"[&_input]:placeholder:text-muted-foreground",
	"[&_input]:disabled:cursor-not-allowed",
	"[&_input]:disabled:opacity-50",
)

// commandListCVA defines styling for scrollable command list
var commandListCVA = forgeui.NewCVA(
	"max-h-[300px]",
	"overflow-y-auto",
	"overflow-x-hidden",
	"p-1",
	"[&::-webkit-scrollbar]:w-2",
	"[&::-webkit-scrollbar-track]:bg-transparent",
	"[&::-webkit-scrollbar-thumb]:bg-border",
	"[&::-webkit-scrollbar-thumb]:rounded-full",
)

// commandEmptyCVA defines styling for empty state
var commandEmptyCVA = forgeui.NewCVA(
	"py-6",
	"text-center",
	"text-sm",
	"text-muted-foreground",
)

// commandGroupCVA defines styling for command groups
var commandGroupCVA = forgeui.NewCVA(
	"overflow-hidden",
	"p-1",
	"text-foreground",
	"[&_[data-command-group-heading]]:px-2",
	"[&_[data-command-group-heading]]:py-1.5",
	"[&_[data-command-group-heading]]:text-xs",
	"[&_[data-command-group-heading]]:font-medium",
	"[&_[data-command-group-heading]]:text-muted-foreground",
)

// commandItemCVA defines styling for individual command items with states
var commandItemCVA = forgeui.NewCVA(
	"relative",
	"flex",
	"cursor-default",
	"select-none",
	"items-center",
	"rounded-sm",
	"px-2",
	"py-1.5",
	"text-sm",
	"outline-none",
	"gap-2",
	"transition-colors",
	"data-[selected=true]:bg-accent",
	"data-[selected=true]:text-accent-foreground",
	"data-[disabled=true]:pointer-events-none",
	"data-[disabled=true]:opacity-50",
	"hover:bg-accent",
	"hover:text-accent-foreground",
	"[&_svg]:size-4",
	"[&_svg]:shrink-0",
)

// commandSeparatorCVA defines styling for command separator
var commandSeparatorCVA = forgeui.NewCVA(
	"-mx-1",
	"h-px",
	"bg-border",
	"my-1",
)

// commandShortcutCVA defines styling for keyboard shortcuts display
var commandShortcutCVA = forgeui.NewCVA(
	"ml-auto",
	"text-xs",
	"tracking-widest",
	"text-muted-foreground",
	"flex",
	"gap-0.5",
)

// commandDialogCVA defines styling for command dialog overlay
var commandDialogCVA = forgeui.NewCVA(
	"[&_[data-command]]:max-w-[450px]",
	"[&_[data-command]]:w-full",
)
