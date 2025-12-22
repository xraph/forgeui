package dropdown

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"
)

// ContextMenu creates a context menu triggered by right-click.
//
// Context menus appear at the cursor position when right-clicking
// on the trigger element. They share the same item components as
// DropdownMenu (items, separators, checkboxes, etc.).
//
// Example:
//
//	contextMenu.ContextMenu(
//	    html.Div(
//	        html.Class("p-4 border border-border"),
//	        g.Text("Right-click me"),
//	    ),
//	    contextMenu.ContextMenuContent(
//	        contextMenu.ContextMenuItem(g.Text("Copy")),
//	        contextMenu.ContextMenuItem(g.Text("Paste")),
//	        contextMenu.ContextMenuSeparator(),
//	        contextMenu.ContextMenuItem(g.Text("Delete")),
//	    ),
//	)
func ContextMenu(trigger g.Node, content ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{
			"open": false,
			"x":    0,
			"y":    0,
		}),
		alpine.XOn("keydown.escape.window", "open = false"),
		html.Class("inline-block"),

		// Trigger with right-click handler
		html.Div(
			alpine.XOn("contextmenu.prevent", "open = true; x = $event.clientX; y = $event.clientY"),
			trigger,
		),

		// Context menu content (positioned at cursor)
		g.Group(content),
	)
}

// ContextMenuContent creates the context menu content panel.
func ContextMenuContent(children ...g.Node) g.Node {
	return html.Div(
		alpine.XShow("open"),
		alpine.XCloak(),
		g.Group(alpine.XTransition(animation.ScaleIn())),
		alpine.XOn("click.outside", "open = false"),
		g.Attr(":style", "`position: fixed; left: ${x}px; top: ${y}px`"),
		html.Class("z-[80] bg-popover text-popover-foreground border border-border rounded-md shadow-md p-1 min-w-[200px]"),
		g.Attr("role", "menu"),
		g.Attr("aria-orientation", "vertical"),
		g.Group(children),
	)
}

// ContextMenuItem creates a context menu item.
// Same as DropdownMenuItem but closes the context menu on click.
func ContextMenuItem(children ...g.Node) g.Node {
	return html.Div(
		g.Attr("role", "menuitem"),
		g.Attr("tabindex", "-1"),
		html.Class("relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"),
		alpine.XOn("click", "open = false"),
		g.Group(children),
	)
}

// ContextMenuCheckboxItem creates a checkable context menu item.
func ContextMenuCheckboxItem(id, label string, checked bool) g.Node {
	return DropdownMenuCheckboxItem(id, label, checked)
}

// ContextMenuRadioGroup creates a radio group container.
func ContextMenuRadioGroup(value string, children ...g.Node) g.Node {
	return DropdownMenuRadioGroup(value, children...)
}

// ContextMenuRadioItem creates a radio context menu item.
func ContextMenuRadioItem(value, label string) g.Node {
	return html.Div(
		g.Attr("role", "menuitemradio"),
		g.Attr("tabindex", "-1"),
		alpine.XOn("click", fmt.Sprintf("selected = '%s'; open = false", value)),
		html.Class("relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"),

		// Radio indicator
		html.Span(
			html.Class("mr-2 flex h-4 w-4 items-center justify-center border border-border rounded-full"),
			alpine.XShow(fmt.Sprintf("selected === '%s'", value)),
			// Dot
			html.Span(
				html.Class("h-2 w-2 rounded-full bg-current"),
			),
		),
		g.Text(label),
	)
}

// ContextMenuSeparator creates a visual separator.
func ContextMenuSeparator() g.Node {
	return DropdownMenuSeparator()
}

// ContextMenuLabel creates a label for grouping items.
func ContextMenuLabel(text string) g.Node {
	return DropdownMenuLabel(text)
}

// ContextMenuShortcut displays a keyboard shortcut hint.
func ContextMenuShortcut(text string) g.Node {
	return DropdownMenuShortcut(text)
}

// ContextMenuSub creates a submenu.
func ContextMenuSub(trigger g.Node, content ...g.Node) g.Node {
	return DropdownMenuSub(trigger, content...)
}

// ContextMenuItemWithIcon creates a context menu item with an icon.
func ContextMenuItemWithIcon(icon, label g.Node) g.Node {
	return html.Div(
		g.Attr("role", "menuitem"),
		g.Attr("tabindex", "-1"),
		html.Class("relative flex cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"),
		alpine.XOn("click", "open = false"),
		html.Span(html.Class("flex-shrink-0"), icon),
		html.Span(label),
	)
}

