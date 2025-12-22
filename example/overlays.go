package main

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/badge"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/components/dropdown"
	"github.com/xraph/forgeui/components/modal"
	"github.com/xraph/forgeui/components/popover"
	"github.com/xraph/forgeui/components/toast"
	"github.com/xraph/forgeui/components/tooltip"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

// OverlaysDemo demonstrates all overlay components
func OverlaysDemo() g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("ForgeUI - Overlay Components")),
			// Tailwind Play CDN with JIT
			html.Script(g.Attr("src", "https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			// Theme CSS variables
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),
			// Alpine x-cloak CSS
			alpine.CloakCSS(),
		),
		html.Body(
			html.Class("bg-background text-foreground"),
			g.Attr("x-data", ""), // Initialize Alpine.js
			theme.DarkModeScript(),

			// Register toast store
			toast.RegisterToastStore(),

			// Toast container
			toast.Toaster(),

			// Main content
			primitives.Box(
				primitives.WithClass("container mx-auto px-4 py-12"),
				primitives.WithChildren(
					// Header
					html.Div(
						html.Class("space-y-6 mb-12"),
						html.H1(
							html.Class("text-4xl font-bold"),
							g.Text("Overlay Components"),
						),
						html.P(
							html.Class("text-lg text-muted-foreground"),
							g.Text("Interactive overlays with Alpine.js and smooth animations"),
						),
					),

					// Components Grid
					html.Div(
						html.Class("space-y-12"),

						// Modal & Dialog Section
						overlaySection(
							"Modals & Dialogs",
							"Modal overlays for focused interactions",
							html.Div(
								html.Class("flex gap-4 flex-wrap"),

								// Basic Modal
								modal.Modal(
									modal.ModalProps{
										Size:                forgeui.SizeMD,
										CloseOnEscape:       true,
										CloseOnOutsideClick: false,
										ShowClose:           true,
									},
									button.Button(g.Text("Basic Modal")),
									modal.ModalHeader("Modal Title", "This is a basic modal with header and content"),
									html.Div(
										html.Class("px-6 py-4"),
										html.P(g.Text("Modal body content goes here. You can add any elements.")),
									),
									modal.ModalFooter(
										modal.ModalClose(button.Button(g.Text("Cancel"), button.WithVariant(forgeui.VariantOutline))),
										button.Button(g.Text("Confirm"), button.WithVariant(forgeui.VariantPrimary)),
									),
								),

								// Dialog
								modal.Dialog(
									modal.DialogTrigger(button.Button(g.Text("Dialog"), button.WithVariant(forgeui.VariantPrimary))),
									modal.DialogContent(
										modal.DialogHeader(
											modal.DialogTitle("Edit Profile"),
											modal.DialogDescription("Make changes to your profile here."),
										),
										modal.DialogBody(
											html.P(g.Text("Dialog content with semantic API.")),
										),
										modal.DialogFooter(
											modal.DialogClose(button.Button(g.Text("Cancel"))),
											button.Button(g.Text("Save"), button.WithVariant(forgeui.VariantPrimary)),
										),
									),
								),

								// Alert Dialog
								modal.AlertDialog(
									modal.AlertDialogTrigger(
										button.Button(g.Text("Delete Account"), button.WithVariant(forgeui.VariantDestructive)),
									),
									modal.AlertDialogContent(
										modal.AlertDialogHeader(
											modal.AlertDialogTitle("Are you absolutely sure?"),
											modal.AlertDialogDescription("This action cannot be undone. This will permanently delete your account."),
										),
										modal.AlertDialogFooter(
											modal.AlertDialogCancel(button.Button(g.Text("Cancel"), button.WithVariant(forgeui.VariantOutline))),
											modal.AlertDialogAction(button.Button(g.Text("Delete"), button.WithVariant(forgeui.VariantDestructive))),
										),
									),
								),
							),
						),

						// Drawers & Sheets Section
						overlaySection(
							"Drawers & Sheets",
							"Side panels that slide in from screen edges",
							html.Div(
								html.Class("flex gap-4 flex-wrap"),

								// Drawer from Right
								modal.Drawer(
									modal.DrawerProps{
										Side: forgeui.SideRight,
										Size: forgeui.SizeMD,
									},
									button.Button(g.Text("Drawer Right")),
									modal.DrawerHeader("Settings", "Configure your preferences"),
									modal.DrawerBody(
										html.P(g.Text("Drawer content goes here.")),
									),
									modal.DrawerFooter(
										modal.DrawerClose(button.Button(g.Text("Close"))),
									),
								),

								// Sheet from Left
								modal.Sheet(
									modal.SheetTrigger(button.Button(g.Text("Sheet Left"), button.WithVariant(forgeui.VariantPrimary))),
									modal.SheetLeft(
										modal.SheetHeader(
											modal.SheetTitle("Navigation"),
											modal.SheetDescription("Browse menu items"),
										),
										modal.SheetBody(
											html.P(g.Text("Sheet content.")),
										),
									),
								),

								// Sheet from Bottom
								modal.Sheet(
									modal.SheetTrigger(button.Button(g.Text("Sheet Bottom"), button.WithVariant(forgeui.VariantSecondary))),
									modal.SheetBottom(
										modal.SheetHeader(
											modal.SheetTitle("Actions"),
										),
										modal.SheetBody(
											html.P(g.Text("Bottom sheet for mobile-style actions.")),
										),
									),
								),
							),
						),

						// Popovers & Tooltips Section
						overlaySection(
							"Popovers & Tooltips",
							"Floating content anchored to elements",
							html.Div(
								html.Class("flex gap-4 flex-wrap items-center"),

								// Popover
								popover.Popover(
									popover.PopoverProps{
										Position: forgeui.PositionBottom,
										Align:    forgeui.AlignStart,
									},
									button.Button(g.Text("Open Popover")),
									html.Div(
										html.H4(html.Class("font-semibold mb-2"), g.Text("Settings")),
										html.P(html.Class("text-sm text-muted-foreground"), g.Text("Quick settings panel")),
									),
								),

								// Tooltips
								tooltip.Tooltip(
									tooltip.TooltipProps{Position: forgeui.PositionTop},
									button.Button(g.Text("Hover me")),
									"This is a helpful tooltip",
								),

								tooltip.Tooltip(
									tooltip.TooltipProps{Position: forgeui.PositionRight},
									badge.Badge("Info"),
									"Right-positioned tooltip",
								),
							),
						),

						// Dropdowns & Menus Section
						overlaySection(
							"Dropdowns & Menus",
							"Click-triggered menus with items and navigation",
							html.Div(
								html.Class("flex gap-4 flex-wrap"),

								// Basic Dropdown
								dropdown.Dropdown(
									dropdown.DropdownProps{
										Position: forgeui.PositionBottom,
										Align:    forgeui.AlignStart,
									},
									button.Button(g.Text("Basic Dropdown")),
									html.Div(
										html.Class("p-2"),
										html.P(g.Text("Custom dropdown content")),
									),
								),

								// Dropdown Menu
								dropdown.DropdownMenu(
									dropdown.DropdownMenuTrigger(
										button.Button(g.Text("Menu"), button.WithVariant(forgeui.VariantPrimary)),
									),
									dropdown.DropdownMenuContent(
										dropdown.DropdownMenuLabel("My Account"),
										dropdown.DropdownMenuSeparator(),
										dropdown.DropdownMenuItem(g.Text("Profile")),
										dropdown.DropdownMenuItem(g.Text("Settings")),
										dropdown.DropdownMenuSeparator(),
										dropdown.DropdownMenuItem(g.Text("Logout")),
									),
								),

								// Dropdown with Checkboxes
								dropdown.DropdownMenu(
									dropdown.DropdownMenuTrigger(
										button.Button(g.Text("Preferences"), button.WithVariant(forgeui.VariantSecondary)),
									),
									dropdown.DropdownMenuContent(
										dropdown.DropdownMenuLabel("Display"),
										dropdown.DropdownMenuSeparator(),
										dropdown.DropdownMenuCheckboxItem("notifications", "Notifications", true),
										dropdown.DropdownMenuCheckboxItem("emails", "Email alerts", false),
									),
								),
							),
						),

						// Context Menu Section
						overlaySection(
							"Context Menu",
							"Right-click triggered menus",
							dropdown.ContextMenu(
								card.CardWithOptions(
									[]card.Option{card.WithClass("p-8 text-center border-2 border-dashed cursor-pointer")},
									html.Div(
										html.P(html.Class("text-muted-foreground"), g.Text("Right-click me")),
									),
								),
								dropdown.ContextMenuContent(
									dropdown.ContextMenuItem(g.Text("Copy")),
									dropdown.ContextMenuItem(g.Text("Paste")),
									dropdown.ContextMenuSeparator(),
									dropdown.ContextMenuItem(g.Text("Delete")),
								),
							),
						),

						// Toast Notifications Section
						overlaySection(
							"Toast Notifications",
							"Temporary notifications with auto-dismiss",
							html.Div(
								html.Class("flex gap-4 flex-wrap"),

								button.Button(
									g.Text("Success Toast"),
									button.WithVariant(forgeui.VariantPrimary),
									button.WithAttrs(
										alpine.XOn("click", toast.ToastSuccess("Operation completed successfully!")),
									),
								),

								button.Button(
									g.Text("Error Toast"),
									button.WithVariant(forgeui.VariantDestructive),
									button.WithAttrs(
										alpine.XOn("click", toast.ToastError("Something went wrong!")),
									),
								),

								button.Button(
									g.Text("Warning Toast"),
									button.WithVariant(forgeui.VariantOutline),
									button.WithAttrs(
										alpine.XOn("click", toast.ToastWarning("Please review your settings")),
									),
								),

								button.Button(
									g.Text("Info Toast"),
									button.WithVariant(forgeui.VariantGhost),
									button.WithAttrs(
										alpine.XOn("click", toast.ToastInfo("New update available")),
									),
								),

								button.Button(
									g.Text("Custom Toast"),
									button.WithVariant(forgeui.VariantSecondary),
									button.WithAttrs(
										alpine.XOn("click", `$store.toast.add({
										title: 'Custom Notification',
										description: 'This toast has a description and longer duration',
										variant: 'default',
										duration: 8000,
										showProgress: true
									})`),
									),
								),
							),
						),
					), // Close html.Div Components Grid (line 62)
				), // Close primitives.WithChildren (line 47)
			), // Close primitives.Box (line 45)

			// Alpine.js scripts at end of body (must be after store registration)
			// Use ScriptsImmediate (no defer) so they execute in document order
			alpine.ScriptsImmediate(alpine.PluginFocus, alpine.PluginCollapse),
		), // Close html.Body (line 35)
	) // Close html.HTML (line 22)
}

func overlaySection(title, description string, content g.Node) g.Node {
	return card.Card(
		card.Header(
			card.Title(title),
			card.Description(description),
		),
		card.Content(
			content,
		),
	)
}
