package dropdown

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
)

// DropdownMenu creates a dropdown menu container with compound API.
//
// DropdownMenu provides a full-featured menu system with:
//   - Menu items with icons
//   - Keyboard navigation (arrows, enter, escape)
//   - Separators and labels
//   - Checkable items
//   - Radio items
//   - Submenu support
//   - Keyboard shortcuts display
//
// Example:
//
//	dropdown.DropdownMenu(
//	    dropdown.DropdownMenuTrigger(button.Button(g.Text("Options"))),
//	    dropdown.DropdownMenuContent(
//	        dropdown.DropdownMenuItem(g.Text("Profile")),
//	        dropdown.DropdownMenuItem(g.Text("Settings")),
//	        dropdown.DropdownMenuSeparator(),
//	        dropdown.DropdownMenuItem(g.Text("Logout")),
//	    ),
//	)
func DropdownMenu(children ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{
			"open":          false,
			"focusedIndex":  -1,
			"itemCount":     0,
			"checkedItems":  map[string]bool{},
			"selectedRadio": "",
		}),
		alpine.XOn("keydown.escape", "open = false; focusedIndex = -1"),
		alpine.XOn("keydown.down.prevent", "focusedIndex = (focusedIndex + 1) % itemCount"),
		alpine.XOn("keydown.up.prevent", "focusedIndex = focusedIndex <= 0 ? itemCount - 1 : focusedIndex - 1"),
		alpine.XOn("keydown.enter", "$refs['item-' + focusedIndex]?.click()"),
		html.Class("relative inline-block"),
		g.Group(children),
	)
}

// DropdownMenuTrigger wraps the trigger element.
func DropdownMenuTrigger(child g.Node) g.Node {
	return html.Div(
		alpine.XOn("click", "open = !open"),
		child,
	)
}

// DropdownMenuContent creates the menu content panel.
func DropdownMenuContent(children ...g.Node) g.Node {
	return DropdownMenuContentWithAlign(forgeui.AlignStart, children...)
}

// DropdownMenuContentWithAlign creates menu content with custom alignment
func DropdownMenuContentWithAlign(align forgeui.Align, children ...g.Node) g.Node {
	alignClass := ""
	switch align {
	case forgeui.AlignStart:
		alignClass = "left-0"
	case forgeui.AlignCenter:
		alignClass = "left-1/2 -translate-x-1/2"
	case forgeui.AlignEnd:
		alignClass = "right-0"
	}

	return html.Div(
		alpine.XShow("open"),
		alpine.XCloak(),
		g.Group(alpine.XTransition(getDropdownTransition(forgeui.PositionBottom))),
		alpine.XOn("click.outside", "open = false; focusedIndex = -1"),
		alpine.XInit("itemCount = $el.querySelectorAll('[role=menuitem]').length"),
		html.Class("absolute z-50 top-full mt-2 "+alignClass+" bg-popover text-popover-foreground border border-border rounded-md shadow-md p-1 min-w-[200px]"),
		g.Attr("role", "menu"),
		g.Attr("aria-orientation", "vertical"),
		g.Group(children),
	)
}

// DropdownMenuItem creates a menu item.
func DropdownMenuItem(children ...g.Node) g.Node {
	return html.Div(
		g.Attr("role", "menuitem"),
		g.Attr("tabindex", "-1"),
		html.Class("relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"),
		alpine.XOn("click", "open = false"),
		g.Group(children),
	)
}

// DropdownMenuCheckboxItem creates a checkable menu item.
func DropdownMenuCheckboxItem(id, label string, checked bool) g.Node {
	return html.Div(
		g.Attr("role", "menuitemcheckbox"),
		g.Attr("tabindex", "-1"),
		g.Attr("x-data", fmt.Sprintf(`{checked: %t}`, checked)),
		alpine.XOn("click", "checked = !checked"),
		html.Class("relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"),
		
		// Checkbox indicator
		html.Span(
			html.Class("mr-2 flex h-4 w-4 items-center justify-center border border-border rounded-sm"),
			alpine.XShow("checked"),
			// Checkmark SVG
			html.SVG(
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("width", "12"),
				g.Attr("height", "12"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("fill", "none"),
				g.Attr("stroke", "currentColor"),
				g.Attr("stroke-width", "3"),
				g.Attr("stroke-linecap", "round"),
				g.Attr("stroke-linejoin", "round"),
				g.El("path", g.Attr("d", "M20 6 9 17l-5-5")),
			),
		),
		g.Text(label),
	)
}

// DropdownMenuRadioGroup creates a radio group container.
func DropdownMenuRadioGroup(value string, children ...g.Node) g.Node {
	return html.Div(
		g.Attr("x-data", fmt.Sprintf(`{selected: '%s'}`, value)),
		g.Attr("role", "group"),
		g.Group(children),
	)
}

// DropdownMenuRadioItem creates a radio menu item.
func DropdownMenuRadioItem(value, label string) g.Node {
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

// DropdownMenuSeparator creates a visual separator.
func DropdownMenuSeparator() g.Node {
	return html.Div(
		g.Attr("role", "separator"),
		html.Class("my-1 h-px bg-border"),
	)
}

// DropdownMenuLabel creates a label for grouping items.
func DropdownMenuLabel(text string) g.Node {
	return html.Div(
		html.Class("px-2 py-1.5 text-xs font-semibold text-muted-foreground"),
		g.Text(text),
	)
}

// DropdownMenuShortcut displays a keyboard shortcut hint.
func DropdownMenuShortcut(text string) g.Node {
	return html.Span(
		html.Class("ml-auto text-xs tracking-widest text-muted-foreground"),
		g.Text(text),
	)
}

// DropdownMenuSub creates a submenu.
func DropdownMenuSub(trigger g.Node, content ...g.Node) g.Node {
	return html.Div(
		alpine.XData(map[string]any{"submenuOpen": false}),
		html.Class("relative"),
		
		// Submenu trigger
		html.Div(
			g.Attr("role", "menuitem"),
			g.Attr("tabindex", "-1"),
			alpine.XOn("mouseenter", "submenuOpen = true"),
			alpine.XOn("mouseleave", "submenuOpen = false"),
			html.Class("relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"),
			trigger,
			// Arrow indicator
			html.Span(
				html.Class("ml-auto"),
				html.SVG(
					g.Attr("xmlns", "http://www.w3.org/2000/svg"),
					g.Attr("width", "12"),
					g.Attr("height", "12"),
					g.Attr("viewBox", "0 0 24 24"),
					g.Attr("fill", "none"),
					g.Attr("stroke", "currentColor"),
					g.Attr("stroke-width", "2"),
					g.El("path", g.Attr("d", "m9 18 6-6-6-6")),
				),
			),
		),
		
		// Submenu content
		html.Div(
			alpine.XShow("submenuOpen"),
			alpine.XCloak(),
			g.Group(alpine.XTransition(getDropdownTransition(forgeui.PositionRight))),
			html.Class("absolute left-full top-0 ml-1 bg-popover text-popover-foreground border border-border rounded-md shadow-md p-1 min-w-[200px]"),
			g.Attr("role", "menu"),
			g.Group(content),
		),
	)
}

// DropdownMenuItemWithIcon creates a menu item with an icon.
func DropdownMenuItemWithIcon(icon, label g.Node) g.Node {
	return html.Div(
		g.Attr("role", "menuitem"),
		g.Attr("tabindex", "-1"),
		html.Class("relative flex cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"),
		alpine.XOn("click", "open = false"),
		html.Span(html.Class("flex-shrink-0"), icon),
		html.Span(label),
	)
}

