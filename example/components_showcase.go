package main

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/accordion"
	"github.com/xraph/forgeui/components/alert"
	"github.com/xraph/forgeui/components/avatar"
	"github.com/xraph/forgeui/components/badge"
	"github.com/xraph/forgeui/components/breadcrumb"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/components/checkbox"
	"github.com/xraph/forgeui/components/dropdown"
	"github.com/xraph/forgeui/components/emptystate"
	"github.com/xraph/forgeui/components/input"
	"github.com/xraph/forgeui/components/label"
	"github.com/xraph/forgeui/components/list"
	"github.com/xraph/forgeui/components/menu"
	"github.com/xraph/forgeui/components/modal"
	"github.com/xraph/forgeui/components/pagination"
	"github.com/xraph/forgeui/components/popover"
	"github.com/xraph/forgeui/components/progress"
	"github.com/xraph/forgeui/components/radio"
	selectcomp "github.com/xraph/forgeui/components/select"
	"github.com/xraph/forgeui/components/separator"
	"github.com/xraph/forgeui/components/skeleton"
	"github.com/xraph/forgeui/components/slider"
	"github.com/xraph/forgeui/components/spinner"
	switchcomp "github.com/xraph/forgeui/components/switch"
	"github.com/xraph/forgeui/components/table"
	"github.com/xraph/forgeui/components/tabs"
	"github.com/xraph/forgeui/components/textarea"
	"github.com/xraph/forgeui/components/toast"
	"github.com/xraph/forgeui/components/tooltip"
	"github.com/xraph/forgeui/icons"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

func handleComponentsShowcase(w http.ResponseWriter, r *http.Request) {
	page := html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Complete Components Showcase")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			alpine.CloakCSS(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),
			html.StyleEl(g.Raw(`
				@layer base {
					* {
						@apply border-border;
					}
				}
			`)),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground"),
			alpine.XData(map[string]any{
				"activeTab":     "all",
				"searchQuery":   "",
				"showToast":     false,
				"selectedRadio": "option1",
				"checked":       false,
				"sliderValue":   50,
			}),

			// Sticky Header
			html.Div(
				html.Class("sticky top-0 z-50 bg-card/95 backdrop-blur border-b"),
				primitives.Container(
					primitives.Box(
						primitives.WithClass("py-4"),
						primitives.WithChildren(
							primitives.HStack("4",
							html.Div(
								html.Class("flex-1 flex items-center gap-3"),
								icons.Menu(icons.WithSize(24)),
								primitives.Text(
									primitives.TextAs("span"),
									primitives.TextSize("text-xl"),
									primitives.TextWeight("font-bold"),
									primitives.TextChildren(g.Text("ForgeUI Components")),
								),
								badge.Badge("50+ Components", badge.WithVariant(forgeui.VariantSecondary)),
							),
								primitives.HStack("2",
									html.A(
										html.Href("/"),
										button.Ghost(
											g.Group([]g.Node{
												icons.Home(icons.WithSize(16)),
												g.Text("Home"),
											}),
										),
									),
									html.A(
										html.Href("/interactive"),
										button.Ghost(
											g.Group([]g.Node{
												icons.Plus(icons.WithSize(16)),
												g.Text("Interactive"),
											}),
										),
									),
									html.A(
										html.Href("/icons"),
										button.Ghost(
											g.Group([]g.Node{
												icons.Search(icons.WithSize(16)),
												g.Text("Icons"),
											}),
										),
									),
									html.A(
										html.Href("/data"),
										button.Ghost(g.Text("Data")),
									),
									theme.SimpleToggle(),
								),
							),
						),
					),
				),
			),

			// Main Content
			primitives.Container(
				primitives.Box(
					primitives.WithClass("py-12 md:py-16 space-y-16"),
					primitives.WithChildren(
						// Hero Section
						showcaseHero(),

						// Components Grid
						showcaseButtonsAndActions(),
						showcaseInputsAndForms(),
						showcaseDisplayComponents(),
						showcaseLayoutComponents(),
						showcaseNavigationComponents(),
						showcaseOverlayComponents(),
						showcaseDataComponents(),
						showcaseFeedbackComponents(),
					),
				),
			),

			// Toast Container
			toast.Toaster(),

			// Alpine.js scripts
			alpine.Scripts(alpine.PluginCollapse),
		),
	)

	w.Header().Set("Content-Type", "text/html")
	_ = page.Render(w)
}

func showcaseHero() g.Node {
	return primitives.Box(
		primitives.WithClass("text-center space-y-6 py-12 md:py-16"),
		primitives.WithChildren(
			primitives.Text(
				primitives.TextAs("h1"),
				primitives.TextSize("text-4xl md:text-5xl lg:text-6xl"),
				primitives.TextWeight("font-bold"),
				primitives.TextClass("tracking-tight"),
				primitives.TextChildren(g.Text("Complete Component Library")),
			),
			primitives.Text(
				primitives.TextSize("text-lg md:text-xl"),
				primitives.TextColor("text-muted-foreground"),
				primitives.TextClass("max-w-3xl mx-auto leading-relaxed"),
				primitives.TextChildren(g.Text("Every ForgeUI component in one place. Production-ready, type-safe, and beautiful.")),
			),
		),
	)
}

func showcaseButtonsAndActions() g.Node {
	return showcaseSection(
		"Buttons & Actions",
		"Interactive elements for user actions",
		primitives.VStack("8",
			// Buttons
			componentDemo(
				"Buttons",
				"Multiple variants and sizes",
				primitives.VStack("4",
					// Variants
					primitives.HStack("2",
						button.Primary(
							g.Group([]g.Node{
								icons.Check(icons.WithSize(16)),
								g.Text("Primary"),
							}),
						),
						button.Secondary(g.Text("Secondary")),
						button.Destructive(
							g.Group([]g.Node{
								icons.Trash(icons.WithSize(16)),
								g.Text("Destructive"),
							}),
						),
						button.Outline(g.Text("Outline")),
						button.Ghost(g.Text("Ghost")),
						button.Link(g.Text("Link")),
					),
					// Sizes
					primitives.HStack("2",
						button.Primary(g.Text("Small"), button.WithSize(forgeui.SizeSM)),
						button.Primary(g.Text("Default")),
						button.Primary(g.Text("Large"), button.WithSize(forgeui.SizeLG)),
						button.IconButton(icons.Settings(icons.WithSize(16))),
					),
					// Button Group
					button.Group(
						[]button.GroupOption{button.WithGap("2")},
						button.Primary(
							g.Group([]g.Node{
								icons.Check(icons.WithSize(16)),
								g.Text("Save"),
							}),
						),
						button.Secondary(
							g.Group([]g.Node{
								icons.X(icons.WithSize(16)),
								g.Text("Cancel"),
							}),
						),
						button.Outline(g.Text("Reset")),
					),
				),
			),

			// Dropdown
			componentDemo(
				"Dropdown Menu",
				"Context menus and dropdowns",
				dropdown.DropdownMenu(
					dropdown.DropdownMenuTrigger(
						button.Outline(g.Text("Open Menu")),
					),
					dropdown.DropdownMenuContent(
						dropdown.DropdownMenuItem(g.Text("Profile")),
						dropdown.DropdownMenuItem(g.Text("Settings")),
						dropdown.DropdownMenuSeparator(),
						dropdown.DropdownMenuItem(g.Text("Logout")),
					),
				),
			),
		),
	)
}

func showcaseInputsAndForms() g.Node {
	return showcaseSection(
		"Inputs & Forms",
		"Form controls and validation",
		primitives.VStack("8",
			// Input
			componentDemo(
				"Input",
				"Text input with variants",
				primitives.VStack("4",
					primitives.VStack("2",
						label.Label("Email", label.WithFor("email")),
						input.Input(
							input.WithID("email"),
							input.WithType("email"),
							input.WithPlaceholder("you@example.com"),
						),
						input.FormDescription("We'll never share your email."),
					),
					primitives.VStack("2",
						label.Label("Disabled Input"),
						input.Input(
							input.WithPlaceholder("Disabled"),
							input.Disabled(),
						),
					),
				),
			),

			// Input Group
			componentDemo(
				"Input Group",
				"Inputs with addons and icons",
				primitives.VStack("4",
					input.InputGroup(
						nil,
						input.InputLeftAddon(nil, g.Text("https://")),
						input.Input(input.WithPlaceholder("example.com")),
					),
					input.InputGroup(
						nil,
						input.Input(input.WithPlaceholder("Search...")),
						input.InputRightAddon(nil, g.Text("🔍")),
					),
					input.InputGroup(
						nil,
						input.InputLeftAddon(nil, g.Text("$")),
						input.Input(
							input.WithType("number"),
							input.WithPlaceholder("0.00"),
						),
						input.InputRightAddon(nil, g.Text(".00")),
					),
				),
			),

			// Textarea
			componentDemo(
				"Textarea",
				"Multi-line text input",
				primitives.VStack("2",
					label.Label("Description", label.WithFor("description")),
					textarea.Textarea(
						textarea.WithID("description"),
						textarea.WithPlaceholder("Enter your description here..."),
						textarea.WithRows(4),
					),
				),
			),

			// Select
			componentDemo(
				"Select",
				"Dropdown selection",
				primitives.VStack("2",
					label.Label("Country", label.WithFor("country")),
					selectcomp.NativeSelect(
						[]selectcomp.NativeSelectOption{
							{Value: "", Label: "Select a country", Disabled: true},
							{Value: "us", Label: "United States"},
							{Value: "uk", Label: "United Kingdom"},
							{Value: "ca", Label: "Canada"},
						},
						selectcomp.WithID("country"),
					),
				),
			),

			// Checkbox & Switch
			componentDemo(
				"Checkbox & Switch",
				"Boolean inputs",
				primitives.VStack("4",
					primitives.HStack("2",
						checkbox.Checkbox(
							checkbox.WithID("terms"),
							checkbox.WithName("terms"),
						),
						label.Label("Accept terms and conditions", label.WithFor("terms")),
					),
					primitives.HStack("2",
						switchcomp.Switch(
							switchcomp.WithID("notifications"),
							switchcomp.WithName("notifications"),
							switchcomp.WithAttrs(alpine.XModel("checked")),
						),
						label.Label("Enable notifications", label.WithFor("notifications")),
					),
				),
			),

			// Radio Group
			componentDemo(
				"Radio Group",
				"Single selection from options",
				primitives.VStack("2",
					label.Label("Plan"),
					radio.RadioGroup("plan", []radio.RadioGroupOption{
						{ID: "free", Value: "free", Label: "Free", Checked: true},
						{ID: "pro", Value: "pro", Label: "Pro - $10/month"},
						{ID: "enterprise", Value: "enterprise", Label: "Enterprise - Custom pricing"},
					}),
				),
			),

			// Slider
			componentDemo(
				"Slider",
				"Numeric range input",
				primitives.VStack("4",
					label.Label("Volume"),
					slider.Slider(
						slider.WithMin(0),
						slider.WithMax(100),
						slider.WithValue(50),
						slider.WithAttrs(alpine.XModel("sliderValue")),
					),
					primitives.Text(
						primitives.TextSize("text-sm"),
						primitives.TextColor("text-muted-foreground"),
						primitives.TextChildren(
							g.Group([]g.Node{
								g.Text("Value: "),
								alpine.XText("sliderValue"),
							}),
						),
					),
				),
			),
		),
	)
}

func showcaseDisplayComponents() g.Node {
	return showcaseSection(
		"Display Components",
		"Visual elements for displaying content",
		primitives.VStack("8",
			// Badges
			componentDemo(
				"Badges",
				"Status indicators and labels",
				primitives.HStack("2",
					badge.Badge("Default"),
					badge.Badge("Secondary", badge.WithVariant(forgeui.VariantSecondary)),
					badge.Badge("Destructive", badge.WithVariant(forgeui.VariantDestructive)),
					badge.Badge("Outline", badge.WithVariant(forgeui.VariantOutline)),
				),
			),

			// Avatars
			componentDemo(
				"Avatars",
				"User representations",
				primitives.HStack("4",
					avatar.Avatar(
						avatar.WithFallback("AB"),
						avatar.WithSize(forgeui.SizeSM),
					),
					avatar.Avatar(
						avatar.WithFallback("CD"),
						avatar.WithSize(forgeui.SizeMD),
					),
					avatar.Avatar(
						avatar.WithFallback("EF"),
						avatar.WithSize(forgeui.SizeLG),
					),
					avatar.Avatar(
						avatar.WithFallback("GH"),
						avatar.WithSize(forgeui.SizeXL),
					),
				),
			),

			// Cards
			componentDemo(
				"Cards",
				"Content containers",
				primitives.Grid(
					primitives.GridCols(1),
					primitives.GridColsMD(2),
					primitives.GridGap("4"),
					primitives.GridChildren(
						card.Card(
							card.Header(
								card.Title("Card Title"),
								card.Description("Card description goes here"),
							),
							card.Content(
								primitives.Text(
									primitives.TextSize("text-sm"),
									primitives.TextChildren(g.Text("This is the card content area where you can add any components or text.")),
								),
							),
							card.Footer(
								button.Primary(g.Text("Action"), button.WithSize(forgeui.SizeSM)),
							),
						),
						card.Card(
							card.Header(
								card.Title("Minimal Card"),
							),
							card.Content(
								primitives.Text(
									primitives.TextSize("text-sm"),
									primitives.TextChildren(g.Text("A simpler card without footer.")),
								),
							),
						),
					),
				),
			),

			// Progress
			componentDemo(
				"Progress",
				"Progress indicators",
				primitives.VStack("4",
					progress.Progress(progress.WithValue(25)),
					progress.Progress(progress.WithValue(50)),
					progress.Progress(progress.WithValue(75)),
				),
			),

			// Skeleton
			componentDemo(
				"Skeleton",
				"Loading placeholders",
				primitives.VStack("3",
					skeleton.Skeleton(
						skeleton.WithHeight("h-4"),
						skeleton.WithWidth("w-full"),
					),
					skeleton.Skeleton(
						skeleton.WithHeight("h-4"),
						skeleton.WithWidth("w-3/4"),
					),
					skeleton.Skeleton(
						skeleton.WithHeight("h-4"),
						skeleton.WithWidth("w-1/2"),
					),
				),
			),

			// Spinner
			componentDemo(
				"Spinner",
				"Loading indicator",
				primitives.HStack("4",
					spinner.Spinner(spinner.WithSize(forgeui.SizeSM)),
					spinner.Spinner(spinner.WithSize(forgeui.SizeMD)),
					spinner.Spinner(spinner.WithSize(forgeui.SizeLG)),
				),
			),
		),
	)
}

func showcaseLayoutComponents() g.Node {
	return showcaseSection(
		"Layout Components",
		"Structure and organization",
		primitives.VStack("8",
			// Separator
			componentDemo(
				"Separator",
				"Visual divider",
				primitives.VStack("4",
					primitives.Text(
						primitives.TextChildren(g.Text("Section 1")),
					),
					separator.Separator(),
					primitives.Text(
						primitives.TextChildren(g.Text("Section 2")),
					),
				),
			),

			// Accordion
			componentDemo(
				"Accordion",
				"Collapsible content sections",
				accordion.Accordion(
					accordion.Item(
						"item-1",
						"What is ForgeUI?",
						g.Text("ForgeUI is a type-safe UI component library for Go that renders server-side HTML."),
					),
					accordion.Item(
						"item-2",
						"How do I install it?",
						g.Text("Install ForgeUI using go get: go get github.com/xraph/forgeui"),
					),
					accordion.Item(
						"item-3",
						"Is it production ready?",
						g.Text("Yes! ForgeUI is battle-tested and ready for production use."),
					),
				),
			),

			// Tabs
			componentDemo(
				"Tabs",
				"Tabbed content",
				tabs.Tabs(
					tabs.TabList(
						tabs.Tab("tab1", g.Text("Overview")),
						tabs.Tab("tab2", g.Text("Features")),
						tabs.Tab("tab3", g.Text("Documentation")),
					),
					tabs.TabPanel("tab1", g.Text("Overview content goes here.")),
					tabs.TabPanel("tab2", g.Text("Features content goes here.")),
					tabs.TabPanel("tab3", g.Text("Documentation content goes here.")),
				),
			),
		),
	)
}

func showcaseNavigationComponents() g.Node {
	return showcaseSection(
		"Navigation",
		"Wayfinding and menus",
		primitives.VStack("8",
			// Breadcrumb
			componentDemo(
				"Breadcrumb",
				"Navigation trail",
				breadcrumb.Breadcrumb(
					breadcrumb.Link("/", g.Text("Home")),
					breadcrumb.Link("/components", g.Text("Components")),
					breadcrumb.Page(g.Text("Showcase")),
				),
			),

			// Menu
			componentDemo(
				"Menu",
				"Navigation menu",
				menu.Menu(
					menu.Item("/dashboard", g.Text("Dashboard")),
					menu.Item("/projects", g.Text("Projects")),
					menu.Item("/settings", g.Text("Settings")),
				),
			),

			// Pagination
			componentDemo(
				"Pagination",
				"Page navigation",
				pagination.Pagination(
					pagination.WithCurrentPage(5),
					pagination.WithTotalPages(10),
					pagination.WithSiblingCount(1),
				),
			),
		),
	)
}

func showcaseOverlayComponents() g.Node {
	return showcaseSection(
		"Overlays",
		"Modals, dialogs, and tooltips",
		primitives.VStack("8",
			// Dialog
			componentDemo(
				"Dialog",
				"Modal dialog",
				modal.Dialog(
					modal.DialogTrigger(
						button.Outline(g.Text("Open Dialog")),
					),
					modal.DialogContent(
						modal.DialogHeader(
							modal.DialogTitle("Dialog Title"),
							modal.DialogDescription("This is a dialog description."),
						),
						modal.DialogBody(
							primitives.Text(
								primitives.TextChildren(g.Text("Dialog content goes here.")),
							),
						),
						modal.DialogFooter(
							modal.DialogClose(
								button.Secondary(g.Text("Cancel")),
							),
							button.Primary(g.Text("Continue")),
						),
					),
				),
			),

			// Alert Dialog
			componentDemo(
				"Alert Dialog",
				"Confirmation dialog",
				modal.AlertDialog(
					modal.AlertDialogTrigger(
						button.Destructive(g.Text("Delete Account")),
					),
					modal.AlertDialogContent(
						modal.AlertDialogHeader(
							modal.AlertDialogTitle("Are you absolutely sure?"),
							modal.AlertDialogDescription("This action cannot be undone. This will permanently delete your account."),
						),
						modal.AlertDialogFooter(
							modal.AlertDialogCancel(
								button.Secondary(g.Text("Cancel")),
							),
							modal.AlertDialogAction(
								button.Destructive(g.Text("Delete")),
							),
						),
					),
				),
			),

			// Popover
			componentDemo(
				"Popover",
				"Contextual popup",
				popover.Popover(
					popover.PopoverProps{Position: forgeui.PositionTop},
					button.Outline(g.Text("Open Popover")),
					html.Div(
						html.Class("space-y-2"),
						primitives.Text(
							primitives.TextWeight("font-medium"),
							primitives.TextChildren(g.Text("Popover Title")),
						),
						primitives.Text(
							primitives.TextSize("text-sm"),
							primitives.TextColor("text-muted-foreground"),
							primitives.TextChildren(g.Text("This is a popover with some content.")),
						),
					),
				),
			),

			// Tooltip
			componentDemo(
				"Tooltip",
				"Hover information",
				tooltip.Tooltip(
					tooltip.TooltipProps{Position: forgeui.PositionTop},
					button.Outline(g.Text("Hover me")),
					"This is a tooltip",
				),
			),
		),
	)
}

func showcaseDataComponents() g.Node {
	return showcaseSection(
		"Data Display",
		"Tables and lists",
		primitives.VStack("8",
			// Table
			componentDemo(
				"Table",
				"Data table",
				table.Table()(
					table.TableHeader()(
						table.TableRow()(
							table.TableHeaderCell()(g.Text("Name")),
							table.TableHeaderCell()(g.Text("Email")),
							table.TableHeaderCell()(g.Text("Role")),
						),
					),
					table.TableBody()(
						table.TableRow()(
							table.TableCell()(g.Text("John Doe")),
							table.TableCell()(g.Text("john@example.com")),
							table.TableCell()(badge.Badge("Admin")),
						),
						table.TableRow()(
							table.TableCell()(g.Text("Jane Smith")),
							table.TableCell()(g.Text("jane@example.com")),
							table.TableCell()(badge.Badge("User", badge.WithVariant(forgeui.VariantSecondary))),
						),
					),
				),
			),

			// List
			componentDemo(
				"List",
				"Unordered and ordered lists",
				primitives.Grid(
					primitives.GridCols(1),
					primitives.GridColsMD(2),
					primitives.GridGap("4"),
					primitives.GridChildren(
						primitives.VStack("2",
							primitives.Text(
								primitives.TextWeight("font-medium"),
								primitives.TextChildren(g.Text("Bullet List")),
							),
							list.List()(
								list.ListItem()(g.Text("First item")),
								list.ListItem()(g.Text("Second item")),
								list.ListItem()(g.Text("Third item")),
							),
						),
						primitives.VStack("2",
							primitives.Text(
								primitives.TextWeight("font-medium"),
								primitives.TextChildren(g.Text("Ordered List")),
							),
							list.OrderedList()(
								list.ListItem()(g.Text("Step 1")),
								list.ListItem()(g.Text("Step 2")),
								list.ListItem()(g.Text("Step 3")),
							),
						),
					),
				),
			),

			// Empty State
			componentDemo(
				"Empty State",
				"No data placeholder",
				emptystate.EmptyState(
					emptystate.WithTitle("No data found"),
					emptystate.WithDescription("Try adjusting your filters or add new data."),
					emptystate.WithAction(
						button.Primary(g.Text("Add Item"), button.WithSize(forgeui.SizeSM)),
					),
				),
			),
		),
	)
}

func showcaseFeedbackComponents() g.Node {
	return showcaseSection(
		"Feedback",
		"User feedback and alerts",
		primitives.VStack("8",
			// Alerts
			componentDemo(
				"Alerts",
				"Informational messages",
				primitives.VStack("4",
					alert.Alert(
						nil,
						alert.AlertTitle("Info"),
						alert.AlertDescription("This is an informational alert."),
					),
					alert.Alert(
						[]alert.Option{alert.WithVariant(forgeui.VariantDestructive)},
						alert.AlertTitle("Error"),
						alert.AlertDescription("This is an error alert."),
					),
				),
			),

			// Toast
			componentDemo(
				"Toast",
				"Temporary notifications",
				button.Primary(
					g.Text("Show Toast"),
					button.WithAttrs(alpine.XOn("click", "$dispatch('show-toast', { message: 'Hello from ForgeUI!', variant: 'default' })")),
				),
			),
		),
	)
}

func showcaseSection(title, description string, children g.Node) g.Node {
	return primitives.Box(
		primitives.WithClass("space-y-6"),
		primitives.WithChildren(
			html.Div(
				html.Class("border-l-4 border-primary pl-4"),
				primitives.VStack("1",
					primitives.Text(
						primitives.TextAs("h2"),
						primitives.TextSize("text-2xl"),
						primitives.TextWeight("font-bold"),
						primitives.TextChildren(g.Text(title)),
					),
					primitives.Text(
						primitives.TextSize("text-sm"),
						primitives.TextColor("text-muted-foreground"),
						primitives.TextChildren(g.Text(description)),
					),
				),
			),
			children,
		),
	)
}

func componentDemo(title, description string, demo g.Node) g.Node {
	return card.Card(
		card.Header(
			card.Title(title),
			card.Description(description),
		),
		card.Content(demo),
	)
}

