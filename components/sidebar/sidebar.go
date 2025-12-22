// Package sidebar provides sidebar navigation components with provider pattern integration.
//
// The sidebar uses ForgeUI's provider pattern for state management, allowing
// child components to access sidebar state without prop drilling.
package sidebar

import (
	"fmt"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/animation"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// SidebarVariant defines the visual style of the sidebar
type SidebarVariant string

const (
	SidebarVariantSidebar  SidebarVariant = "sidebar"  // Default sidebar
	SidebarVariantFloating SidebarVariant = "floating" // Floating sidebar with rounded corners
	SidebarVariantInset    SidebarVariant = "inset"    // Inset sidebar with margins
)

// SidebarCollapsibleMode defines how the sidebar collapses
type SidebarCollapsibleMode string

const (
	CollapsibleOffcanvas SidebarCollapsibleMode = "offcanvas" // Slides off screen
	CollapsibleIcon      SidebarCollapsibleMode = "icon"      // Shrinks to icon only
	CollapsibleNone      SidebarCollapsibleMode = "none"      // Not collapsible
)

// SidebarProps defines sidebar configuration
type SidebarProps struct {
	DefaultCollapsed       bool
	Collapsible            bool
	CollapsibleMode        SidebarCollapsibleMode
	Variant                SidebarVariant
	Side                   forgeui.Side
	Class                  string
	Attrs                  []g.Node
	KeyboardShortcutKey    string
	KeyboardShortcutEnable bool
	PersistState           bool
	StorageKey             string
	StorageType            string // "localStorage" or "sessionStorage"
}

// SidebarOption is a functional option for configuring sidebar
type SidebarOption func(*SidebarProps)

// WithDefaultCollapsed sets initial collapsed state
func WithDefaultCollapsed(collapsed bool) SidebarOption {
	return func(p *SidebarProps) { p.DefaultCollapsed = collapsed }
}

// WithCollapsible allows sidebar to be collapsed
func WithCollapsible(collapsible bool) SidebarOption {
	return func(p *SidebarProps) { p.Collapsible = collapsible }
}

// WithSide sets sidebar position (left or right)
func WithSide(side forgeui.Side) SidebarOption {
	return func(p *SidebarProps) { p.Side = side }
}

// WithSidebarClass adds custom classes to sidebar
func WithSidebarClass(class string) SidebarOption {
	return func(p *SidebarProps) { p.Class = class }
}

// WithSidebarAttrs adds custom attributes to sidebar
func WithSidebarAttrs(attrs ...g.Node) SidebarOption {
	return func(p *SidebarProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// WithVariant sets the visual variant of the sidebar
func WithVariant(variant SidebarVariant) SidebarOption {
	return func(p *SidebarProps) { p.Variant = variant }
}

// WithCollapsibleMode sets how the sidebar collapses
func WithCollapsibleMode(mode SidebarCollapsibleMode) SidebarOption {
	return func(p *SidebarProps) { p.CollapsibleMode = mode }
}

// WithKeyboardShortcut enables keyboard shortcut and sets the key (default: "b" for Cmd/Ctrl+B)
func WithKeyboardShortcut(key string) SidebarOption {
	return func(p *SidebarProps) {
		p.KeyboardShortcutEnable = true
		p.KeyboardShortcutKey = key
	}
}

// WithKeyboardShortcutEnabled enables or disables keyboard shortcuts
func WithKeyboardShortcutEnabled(enabled bool) SidebarOption {
	return func(p *SidebarProps) { p.KeyboardShortcutEnable = enabled }
}

// WithPersistState enables state persistence to localStorage/sessionStorage
func WithPersistState(persist bool) SidebarOption {
	return func(p *SidebarProps) { p.PersistState = persist }
}

// WithStorageKey sets a custom storage key for state persistence
func WithStorageKey(key string) SidebarOption {
	return func(p *SidebarProps) { p.StorageKey = key }
}

// WithStorageType sets the storage type ("localStorage" or "sessionStorage")
func WithStorageType(storageType string) SidebarOption {
	return func(p *SidebarProps) { p.StorageType = storageType }
}

// defaultSidebarProps returns default sidebar property values
func defaultSidebarProps() *SidebarProps {
	return &SidebarProps{
		DefaultCollapsed:       false,
		Collapsible:            true,
		CollapsibleMode:        CollapsibleOffcanvas,
		Variant:                SidebarVariantSidebar,
		Side:                   forgeui.SideLeft,
		KeyboardShortcutKey:    "b",
		KeyboardShortcutEnable: false, // Opt-in feature for backward compatibility
		PersistState:           false,
		StorageKey:             "forgeui_sidebar_state",
		StorageType:            "localStorage",
	}
}

// Sidebar creates a collapsible sidebar navigation with Alpine.js stores.
//
// The sidebar uses Alpine.js stores for state management, allowing any
// component on the page to access sidebar state via $store.sidebar.
// Use SidebarLayoutContent or SidebarInset to create content areas that
// react to sidebar state changes.
//
// Example:
//
//	sidebar.Sidebar(
//	    sidebar.SidebarHeader(g.Text("My App")),
//	    sidebar.SidebarContent(
//	        menu.Section("Main",
//	            menu.Item("/dashboard", g.Text("Dashboard")),
//	            menu.Item("/users", g.Text("Users")),
//	        ),
//	    ),
//	    sidebar.SidebarFooter(g.Text("© 2024")),
//	)
func Sidebar(children ...g.Node) g.Node {
	return SidebarWithOptions(nil, children...)
}

// SidebarWithOptions creates sidebar with custom options and Alpine.js store integration
func SidebarWithOptions(opts []SidebarOption, children ...g.Node) g.Node {
	props := defaultSidebarProps()
	for _, opt := range opts {
		opt(props)
	}

	var sideClass string
	if props.Side == forgeui.SideLeft {
		sideClass = "left-0"
	} else {
		sideClass = "right-0"
	}

	// Base classes
	classes := "fixed top-0 bottom-0 z-30 flex flex-col bg-sidebar text-sidebar-foreground border-border border-r transition-all duration-300 group/sidebar " + sideClass

	// Add variant-specific classes
	switch props.Variant {
	case SidebarVariantFloating:
		classes += " m-2 rounded-lg shadow-lg"
	case SidebarVariantInset:
		classes += " m-2 rounded-xl shadow-sm"
	}

	// Add collapsible mode classes
	switch props.CollapsibleMode {
	case CollapsibleOffcanvas:
		classes += " transition-transform"
	case CollapsibleIcon:
		classes += " transition-width"
	}

	if props.Class != "" {
		classes += " " + props.Class
	}

	// Build width class logic based on collapsible mode using object-based :class binding
	// This is more reliable than ternary expressions with Alpine stores
	var widthClass string

	switch props.CollapsibleMode {
	case CollapsibleIcon:
		// Icon mode: collapse to w-16 (shows only icons)
		widthClass = `{
			'w-64': !$store.sidebar || (!$store.sidebar.isMobile && !$store.sidebar.collapsed) || ($store.sidebar.isMobile && $store.sidebar.mobileOpen),
			'w-16': $store.sidebar && !$store.sidebar.isMobile && $store.sidebar.collapsed,
			'w-0 -translate-x-full': $store.sidebar && $store.sidebar.isMobile && !$store.sidebar.mobileOpen
		}`
	case CollapsibleOffcanvas:
		// Offcanvas mode: slide completely out of view
		widthClass = `{
			'w-64': !$store.sidebar || !$store.sidebar.collapsed || ($store.sidebar.isMobile && $store.sidebar.mobileOpen),
			'-translate-x-full': $store.sidebar && !$store.sidebar.isMobile && $store.sidebar.collapsed,
			'w-0 -translate-x-full': $store.sidebar && $store.sidebar.isMobile && !$store.sidebar.mobileOpen
		}`
	default:
		// No collapse: always w-64
		widthClass = `'w-64'`
	}

	// Build sidebar x-init script to register store
	collapsedStr := "false"
	if props.DefaultCollapsed {
		collapsedStr = "true"
	}

	collapsibleStr := "true"
	if !props.Collapsible {
		collapsibleStr = "false"
	}

	// Config is already set via state, no logging needed in production
	_ = fmt.Sprintf("%s%s%s", props.Variant, props.CollapsibleMode, props.Side)

	// Build store registration script that runs before Alpine initializes
	storeScript := fmt.Sprintf(`
		document.addEventListener('alpine:init', function() {
			if (!Alpine.store('sidebar')) {
				Alpine.store('sidebar', {
					collapsed: %s,
					mobileOpen: false,
					collapsible: %s,
					isMobile: window.innerWidth < 768
				});
			}
		});
		window.addEventListener('resize', function() {
			if (window.Alpine && Alpine.store('sidebar')) {
				Alpine.store('sidebar').isMobile = window.innerWidth < 768;
			}
		});
	`, collapsedStr, collapsibleStr)

	// Add state persistence
	if props.PersistState {
		storeScript += fmt.Sprintf(`
		document.addEventListener('alpine:init', function() {
			var storageKey = '%s';
			var storage = window.%s;
			try {
				var savedState = storage.getItem(storageKey);
				if (savedState !== null && Alpine.store('sidebar')) {
					Alpine.store('sidebar').collapsed = savedState === 'true';
				}
			} catch (e) {}
		});
		`, props.StorageKey, props.StorageType)
	}

	// Add keyboard shortcut
	if props.KeyboardShortcutEnable {
		storeScript += fmt.Sprintf(`
		window.addEventListener('keydown', function(event) {
			if ((event.metaKey || event.ctrlKey) && event.key === '%s') {
				event.preventDefault();
				if (window.Alpine && Alpine.store('sidebar') && Alpine.store('sidebar').collapsible) {
					Alpine.store('sidebar').collapsed = !Alpine.store('sidebar').collapsed;
				}
			}
		});
		`, props.KeyboardShortcutKey)
	}

	// Build sidebar content
	sidebarContent := []g.Node{
		// Store registration script - must be before Alpine loads
		html.Script(g.Raw(storeScript)),

		// Sidebar container
		html.Aside(
			html.Class(classes),
			g.Attr("x-data", "{}"),
			g.Attr(":class", widthClass),
			g.Attr("role", "complementary"),
			g.Attr("aria-label", "Main navigation sidebar"),
			g.Attr(":aria-expanded", "$store.sidebar && $store.sidebar.collapsed ? 'false' : 'true'"),
			g.Attr("data-state", "$store.sidebar && $store.sidebar.collapsed ? 'collapsed' : 'expanded'"),
			g.Attr("data-collapsible", string(props.CollapsibleMode)),
			g.Attr("data-variant", string(props.Variant)),
			g.Attr("data-side", string(props.Side)),
			g.Attr("data-provider", "sidebar"),
			g.Group(append(props.Attrs, children...)),
		),

		// Screen reader live region for state changes
		html.Div(
			g.Attr("role", "status"),
			g.Attr("aria-live", "polite"),
			g.Attr("aria-atomic", "true"),
			html.Class("sr-only"),
			g.Attr("x-data", "{}"),
			g.Attr("x-text", "$store.sidebar && $store.sidebar.collapsed ? 'Sidebar collapsed' : 'Sidebar expanded'"),
		),

		// Mobile backdrop
		html.Div(
			g.Attr("x-data", "{}"),
			g.Attr("x-show", "$store.sidebar && $store.sidebar.isMobile && $store.sidebar.mobileOpen"),
			alpine.XOn("click", "if ($store.sidebar) $store.sidebar.mobileOpen = false"),
			g.Group(alpine.XTransition(animation.FadeIn())),
			html.Class("fixed inset-0 z-20 bg-background/80 backdrop-blur-sm md:hidden"),
			g.Attr("role", "presentation"),
			g.Attr("aria-hidden", "true"),
		),

		// Mobile toggle button
		html.Button(
			g.Attr("type", "button"),
			g.Attr("x-data", "{}"),
			g.Attr("x-show", "$store.sidebar && $store.sidebar.isMobile && !$store.sidebar.mobileOpen"),
			alpine.XOn("click", "if ($store.sidebar) $store.sidebar.mobileOpen = true"),
			html.Class("fixed top-4 left-4 z-10 md:hidden inline-flex items-center justify-center rounded-md p-2 bg-background border shadow-md text-muted-foreground hover:bg-accent hover:text-accent-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 outline-none"),
			g.Attr("aria-label", "Open sidebar"),
			g.El("svg",
				html.Class("h-6 w-6"),
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("fill", "none"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("stroke", "currentColor"),
				g.El("path",
					g.Attr("stroke-linecap", "round"),
					g.Attr("stroke-linejoin", "round"),
					g.Attr("stroke-width", "2"),
					g.Attr("d", "M4 6h16M4 12h16M4 18h16"),
				),
			),
		),
	}

	// Wrap in a div to make it renderable (g.Group can't be rendered directly)
	return html.Div(
		g.Attr("data-sidebar-wrapper", ""),
		g.Group(sidebarContent),
	)
}

// SidebarHeader creates the sidebar header section
//
// Example:
//
//	sidebar.SidebarHeader(
//	    html.Img(g.Attr("src", "/logo.svg")),
//	    g.Text("My App"),
//	)
func SidebarHeader(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex items-center border-b border-border py-4 font-semibold"),
		g.Attr("x-data", "{}"),
		// Dynamic padding and centering for collapsed state
		g.Attr(":class", "$store.sidebar && $store.sidebar.collapsed && !$store.sidebar.isMobile ? 'justify-center px-2' : 'gap-2 px-4'"),
		g.Group(children),
	)
}

// SidebarContent creates the main navigation area
// Matches shadcn sidebar-content styling
//
// Example:
//
//	sidebar.SidebarContent(
//	    menu.Section("Main",
//	        menu.Item("/", g.Text("Home")),
//	    ),
//	)
func SidebarContent(children ...g.Node) g.Node {
	return html.Div(
		html.Class("flex min-h-0 flex-1 flex-col gap-2 overflow-auto"),
		g.Attr("data-slot", "sidebar-content"),
		g.Attr("x-data", "{}"),
		g.Attr(":class", "$store.sidebar && $store.sidebar.collapsed && !$store.sidebar.isMobile ? '[&>*]:p-1' : ''"),
		g.Group(children),
	)
}

// SidebarFooter creates the sidebar footer section
//
// Example:
//
//	sidebar.SidebarFooter(
//	    menu.Item("#", g.Text("Settings")),
//	)
func SidebarFooter(children ...g.Node) g.Node {
	return html.Div(
		html.Class("border-t border-border"),
		g.Attr("x-data", "{}"),
		g.Attr(":class", "$store.sidebar && $store.sidebar.collapsed && !$store.sidebar.isMobile ? 'p-2' : 'p-4'"),
		g.Group(children),
	)
}

// SidebarToggle creates a collapse/expand toggle button
//
// Example:
//
//	sidebar.SidebarToggle()
func SidebarToggle() g.Node {
	return html.Button(
		g.Attr("type", "button"),
		g.Attr("x-data", "{}"),
		g.Attr("x-show", "$store.sidebar && $store.sidebar.collapsible && !$store.sidebar.isMobile"),
		alpine.XOn("click", "if ($store.sidebar && $store.sidebar.collapsible) { $store.sidebar.collapsed = !$store.sidebar.collapsed }"),
		html.Class("absolute -right-3 top-20 z-40 flex h-6 w-6 items-center justify-center rounded-full border border-border bg-sidebar shadow-md text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-transform focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 outline-none"),
		g.Attr("aria-label", "Toggle sidebar"),
		g.Attr(":aria-expanded", "$store.sidebar && !$store.sidebar.collapsed ? 'true' : 'false'"),

		// Single chevron that rotates
		g.El("svg",
			html.Class("h-4 w-4 transition-transform duration-200"),
			g.Attr(":class", "$store.sidebar && $store.sidebar.collapsed ? 'rotate-180' : ''"),
			g.Attr("xmlns", "http://www.w3.org/2000/svg"),
			g.Attr("fill", "none"),
			g.Attr("viewBox", "0 0 24 24"),
			g.Attr("stroke", "currentColor"),
			g.Attr("aria-hidden", "true"),
			g.El("path",
				g.Attr("stroke-linecap", "round"),
				g.Attr("stroke-linejoin", "round"),
				g.Attr("stroke-width", "2"),
				g.Attr("d", "m15 18-6-6 6-6"),
			),
		),
	)
}

// SidebarLayout creates a layout wrapper that adjusts content for sidebar width
//
// Example:
//
//	sidebar.SidebarLayout(
//	    sidebar.Sidebar(...),
//	    html.Main(...), // your content
//	)
func SidebarLayout(children ...g.Node) g.Node {
	return html.Div(
		html.Class("min-h-screen"),
		g.Group(children),
	)
}

// SidebarLayoutContent creates content area that respects sidebar width
// Uses Alpine store to access sidebar state (works as sibling to Sidebar)
//
// Example:
//
//	sidebar.SidebarLayoutContent(
//	    html.Div(...), // your page content
//	)
func SidebarLayoutContent(children ...g.Node) g.Node {
	return html.Div(
		html.Class("transition-all duration-300 min-h-screen"),
		g.Attr("x-data", "{}"),
		g.Attr(":class", `{
			'ml-64': $store.sidebar && !$store.sidebar.isMobile && !$store.sidebar.collapsed,
			'ml-16': $store.sidebar && !$store.sidebar.isMobile && $store.sidebar.collapsed,
			'ml-0': $store.sidebar && $store.sidebar.isMobile
		}`),
		g.Group(children),
	)
}

// SidebarGroup creates a collapsible group within the sidebar
//
// Example:
//
//	sidebar.SidebarGroup(
//	    sidebar.SidebarGroupLabel("Projects"),
//	    sidebar.SidebarMenu(...),
//	)
func SidebarGroup(children ...g.Node) g.Node {
	return html.Div(
		html.Class("relative flex w-full min-w-0 flex-col p-2"),
		g.Attr("data-slot", "sidebar-group"),
		g.Group(children),
	)
}

// SidebarGroupProps defines collapsible sidebar group configuration
type SidebarGroupProps struct {
	Collapsible    bool
	DefaultOpen    bool
	CollapsibleKey string
	Class          string
}

// SidebarGroupOption is a functional option for configuring collapsible sidebar groups
type SidebarGroupOption func(*SidebarGroupProps)

// WithGroupCollapsible makes the group collapsible
func WithGroupCollapsible() SidebarGroupOption {
	return func(p *SidebarGroupProps) { p.Collapsible = true }
}

// WithGroupDefaultOpen sets the default open state
func WithGroupDefaultOpen(open bool) SidebarGroupOption {
	return func(p *SidebarGroupProps) { p.DefaultOpen = open }
}

// WithGroupKey sets a unique key for the collapsible state
func WithGroupKey(key string) SidebarGroupOption {
	return func(p *SidebarGroupProps) { p.CollapsibleKey = key }
}

// WithGroupClass adds custom classes
func WithGroupClass(class string) SidebarGroupOption {
	return func(p *SidebarGroupProps) { p.Class = class }
}

// defaultSidebarGroupProps returns default group properties
func defaultSidebarGroupProps() *SidebarGroupProps {
	return &SidebarGroupProps{
		Collapsible:    false,
		DefaultOpen:    true,
		CollapsibleKey: "group",
	}
}

// SidebarGroupCollapsible creates a collapsible group with animation
//
// Example:
//
//	sidebar.SidebarGroupCollapsible(
//	    []SidebarGroupOption{sidebar.WithGroupKey("projects"), sidebar.WithGroupDefaultOpen(true)},
//	    sidebar.SidebarGroupLabel("Projects"),
//	    sidebar.SidebarMenu(...),
//	)
func SidebarGroupCollapsible(opts []SidebarGroupOption, children ...g.Node) g.Node {
	props := defaultSidebarGroupProps()
	for _, opt := range opts {
		opt(props)
	}

	classes := "relative flex w-full min-w-0 flex-col p-2"
	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	return html.Div(
		html.Class(classes),
		g.Attr("data-slot", "sidebar-group"),
		g.Attr("x-data", fmt.Sprintf("{ %s_open: %t }", props.CollapsibleKey, props.DefaultOpen)),
		g.Group(children),
	)
}

// SidebarGroupLabel creates a label for a sidebar group
// Matches shadcn sidebar-group-label styling
//
// Example:
//
//	sidebar.SidebarGroupLabel("Navigation")
func SidebarGroupLabel(text string, children ...g.Node) g.Node {
	nodes := []g.Node{g.Text(text)}
	if len(children) > 0 {
		nodes = append(nodes, children...)
	}

	return html.Div(
		html.Class("flex h-8 shrink-0 items-center rounded-md px-2 text-xs font-medium text-sidebar-foreground/70 outline-none ring-sidebar-ring transition-[margin,opacity] duration-200 ease-linear focus-visible:ring-2"),
		g.Attr("data-slot", "sidebar-group-label"),
		g.Attr("x-data", "{}"),
		g.Attr(":class", "$store.sidebar && $store.sidebar.collapsed && !$store.sidebar.isMobile ? 'opacity-0 h-0 overflow-hidden -mt-2' : ''"),
		g.Group(nodes),
	)
}

// SidebarGroupLabelCollapsible creates a clickable label for collapsible groups
// Matches shadcn sidebar-group-label styling with collapsible behavior
//
// Example:
//
//	sidebar.SidebarGroupLabelCollapsible("projects", "Projects", icon)
func SidebarGroupLabelCollapsible(key string, text string, icon g.Node) g.Node {
	return html.Button(
		g.Attr("type", "button"),
		g.Attr("data-slot", "sidebar-group-label"),
		alpine.XOn("click", fmt.Sprintf("%s_open = !%s_open", key, key)),
		html.Class("flex h-8 w-full shrink-0 items-center rounded-md px-2 text-xs font-medium text-sidebar-foreground/70 outline-none ring-sidebar-ring transition-[margin,opacity] duration-200 ease-linear hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2"),
		g.Attr(":class", "$store.sidebar && $store.sidebar.collapsed && !$store.sidebar.isMobile ? 'justify-center' : ''"),
		g.Attr("aria-label", fmt.Sprintf("Toggle %s section", text)),
		g.Attr(":aria-expanded", key+"_open ? 'true' : 'false'"),
		html.Span(
			g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
			g.Text(text),
		),
		html.Span(
			html.Class("ml-auto"),
			g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
			g.El("svg",
				html.Class("h-4 w-4 transition-transform duration-200"),
				g.Attr(":class", key+"_open ? '' : '-rotate-90'"),
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("fill", "none"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("stroke", "currentColor"),
				g.Attr("aria-hidden", "true"),
				g.El("path",
					g.Attr("stroke-linecap", "round"),
					g.Attr("stroke-linejoin", "round"),
					g.Attr("stroke-width", "2"),
					g.Attr("d", "m6 9 6 6 6-6"),
				),
			),
		),
	)
}

// SidebarGroupContent wraps collapsible group content
//
// Example:
//
//	sidebar.SidebarGroupContent("projects", sidebar.SidebarMenu(...))
func SidebarGroupContent(key string, children ...g.Node) g.Node {
	return html.Div(
		g.Attr("data-slot", "sidebar-group-content"),
		g.Attr("x-show", key+"_open"),
		g.Attr("x-collapse", ""),
		html.Class("w-full text-sm"),
		g.Group(children),
	)
}

// SidebarGroupAction creates an action button for a group (e.g., "Add Project")
//
// Example:
//
//	sidebar.SidebarGroupAction(icons.Plus(), "Add Project")
func SidebarGroupAction(icon g.Node, label string) g.Node {
	return html.Button(
		g.Attr("type", "button"),
		html.Class("ml-auto opacity-0 group-hover:opacity-100 transition-opacity p-1 rounded-md hover:bg-sidebar-accent text-sidebar-foreground/70 hover:text-sidebar-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-0 outline-none focus-visible:opacity-100"),
		g.Attr("aria-label", label),
		g.Attr("tabindex", "0"),
		icon,
	)
}

// SidebarMenu creates a menu list container
// Matches shadcn sidebar-menu styling
//
// Example:
//
//	sidebar.SidebarMenu(
//	    sidebar.SidebarMenuItem(...),
//	    sidebar.SidebarMenuItem(...),
//	)
func SidebarMenu(children ...g.Node) g.Node {
	return html.Ul(
		html.Class("flex w-full min-w-0 flex-col gap-1"),
		g.Attr("data-slot", "sidebar-menu"),
		g.Group(children),
	)
}

// SidebarMenuItem creates a menu item container
// Matches shadcn sidebar-menu-item styling
//
// Example:
//
//	sidebar.SidebarMenuItem(
//	    sidebar.SidebarMenuButton(...),
//	)
func SidebarMenuItem(children ...g.Node) g.Node {
	return html.Li(
		html.Class("group/menu-item relative"),
		g.Attr("data-slot", "sidebar-menu-item"),
		g.Group(children),
	)
}

// SidebarMenuButtonVariant defines menu button visual styles
type SidebarMenuButtonVariant string

const (
	MenuButtonDefault SidebarMenuButtonVariant = "default"
	MenuButtonOutline SidebarMenuButtonVariant = "outline"
)

// SidebarMenuButtonSize defines menu button sizes
type SidebarMenuButtonSize string

const (
	MenuButtonSizeSmall   SidebarMenuButtonSize = "sm"
	MenuButtonSizeDefault SidebarMenuButtonSize = "default"
	MenuButtonSizeLarge   SidebarMenuButtonSize = "lg"
)

// SidebarMenuButtonProps defines menu button configuration
type SidebarMenuButtonProps struct {
	Href     string
	Active   bool
	Icon     g.Node
	Badge    g.Node
	Tooltip  string
	AsButton bool
	Variant  SidebarMenuButtonVariant
	Size     SidebarMenuButtonSize
	Class    string
	Attrs    []g.Node
}

// SidebarMenuButtonOption is a functional option for menu buttons
type SidebarMenuButtonOption func(*SidebarMenuButtonProps)

// WithMenuHref sets the link URL
func WithMenuHref(href string) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Href = href }
}

// WithMenuActive marks the menu item as active
func WithMenuActive() SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Active = true }
}

// WithMenuIcon adds an icon to the menu button
func WithMenuIcon(icon g.Node) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Icon = icon }
}

// WithMenuBadge adds a badge to the menu button
func WithMenuBadge(badge g.Node) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Badge = badge }
}

// WithMenuTooltip adds a tooltip for collapsed state
func WithMenuTooltip(text string) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Tooltip = text }
}

// WithMenuAsButton renders as button instead of link
func WithMenuAsButton() SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.AsButton = true }
}

// WithMenuClass adds custom classes
func WithMenuClass(class string) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Class = class }
}

// WithMenuAttrs adds custom attributes
func WithMenuAttrs(attrs ...g.Node) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// WithMenuVariant sets the menu button variant
func WithMenuVariant(variant SidebarMenuButtonVariant) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Variant = variant }
}

// WithMenuSize sets the menu button size
func WithMenuSize(size SidebarMenuButtonSize) SidebarMenuButtonOption {
	return func(p *SidebarMenuButtonProps) { p.Size = size }
}

// defaultSidebarMenuButtonProps returns default menu button properties
func defaultSidebarMenuButtonProps() *SidebarMenuButtonProps {
	return &SidebarMenuButtonProps{
		Active:   false,
		AsButton: false,
		Variant:  MenuButtonDefault,
		Size:     MenuButtonSizeDefault,
	}
}

// getMenuButtonVariantClasses returns variant-specific classes
func getMenuButtonVariantClasses(variant SidebarMenuButtonVariant) string {
	switch variant {
	case MenuButtonOutline:
		return "bg-background border border-sidebar-border shadow-sm hover:shadow-md"
	default:
		return ""
	}
}

// getMenuButtonSizeClasses returns size-specific classes
func getMenuButtonSizeClasses(size SidebarMenuButtonSize) string {
	switch size {
	case MenuButtonSizeSmall:
		return "h-7 py-1.5 text-xs"
	case MenuButtonSizeLarge:
		return "h-12 py-3 text-sm"
	default:
		return "h-8 py-2 text-sm"
	}
}

// SidebarMenuButton creates a menu button/link with icon and label
//
// Example:
//
//	sidebar.SidebarMenuButton(
//	    "Dashboard",
//	    sidebar.WithMenuHref("/dashboard"),
//	    sidebar.WithMenuIcon(icons.LayoutDashboard()),
//	    sidebar.WithMenuActive(),
//	)
func SidebarMenuButton(label string, opts ...SidebarMenuButtonOption) g.Node {
	props := defaultSidebarMenuButtonProps()
	for _, opt := range opts {
		opt(props)
	}

	// Base classes with focus ring - removed px-3 to make it dynamic
	baseClasses := "flex items-center rounded-md font-medium transition-colors group-hover/menu-item:bg-sidebar-accent peer/menu-button focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 outline-none"

	// Add variant classes
	baseClasses += " " + getMenuButtonVariantClasses(props.Variant)

	// Add size classes
	baseClasses += " " + getMenuButtonSizeClasses(props.Size)

	// Add active/inactive classes
	if props.Active {
		baseClasses += " bg-sidebar-accent text-sidebar-accent-foreground"
	} else {
		baseClasses += " text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
	}

	if props.Class != "" {
		baseClasses = forgeui.CN(baseClasses, props.Class)
	}

	// Build content nodes
	contentNodes := []g.Node{}

	// Add icon if provided
	if props.Icon != nil {
		contentNodes = append(contentNodes, html.Span(
			html.Class("shrink-0"),
			props.Icon,
		))
	}

	// Add label (hidden when collapsed) - uses $store.sidebar
	contentNodes = append(contentNodes, html.Span(
		g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
		html.Class("flex-1"),
		g.Text(label),
	))

	// Add badge if provided
	if props.Badge != nil {
		contentNodes = append(contentNodes, html.Span(
			g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
			html.Class("ml-auto"),
			props.Badge,
		))
	}

	attrs := []g.Node{
		html.Class(baseClasses),
		g.Attr("x-data", "{}"),
		// Dynamic classes for collapsed state - center icon and remove horizontal gap
		g.Attr(":class", "$store.sidebar && $store.sidebar.collapsed && !$store.sidebar.isMobile ? 'justify-center px-0' : 'gap-3 px-3'"),
	}
	attrs = append(attrs, props.Attrs...)

	// Add ARIA attributes
	attrs = append(attrs,
		g.Attr("role", "menuitem"),
		g.Attr("tabindex", "0"),
		g.Attr("aria-label", label),
	)

	if props.Active {
		attrs = append(attrs, g.Attr("aria-current", "page"))
	}

	if props.AsButton {
		attrs = append(attrs, g.Attr("type", "button"))

		return html.Button(
			g.Group(attrs),
			g.Group(contentNodes),
		)
	}

	// Render as link
	attrs = append(attrs, g.Attr("href", props.Href))

	// Wrap with tooltip if provided
	button := html.A(
		g.Group(attrs),
		g.Group(contentNodes),
	)

	if props.Tooltip != "" {
		return SidebarMenuTooltip(props.Tooltip, button)
	}

	return button
}

// SidebarMenuSub creates a submenu container
// Matches shadcn sidebar-menu-sub styling with border indicator
//
// Example:
//
//	sidebar.SidebarMenuSub(
//	    sidebar.SidebarMenuSubItem(...),
//	    sidebar.SidebarMenuSubItem(...),
//	)
func SidebarMenuSub(children ...g.Node) g.Node {
	return html.Ul(
		html.Class("mx-3.5 flex min-w-0 translate-x-px flex-col gap-1 border-l border-sidebar-border px-2.5 py-0.5"),
		g.Attr("data-slot", "sidebar-menu-sub"),
		g.Attr("x-data", "{}"),
		g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
		g.Group(children),
	)
}

// SidebarMenuSubItem creates a submenu item
// Matches shadcn sidebar-menu-sub-item styling
//
// Example:
//
//	sidebar.SidebarMenuSubItem(
//	    sidebar.SidebarMenuSubButton(...),
//	)
func SidebarMenuSubItem(children ...g.Node) g.Node {
	return html.Li(
		g.Attr("data-slot", "sidebar-menu-sub-item"),
		g.Group(children),
	)
}

// SidebarMenuSubButton creates a submenu button/link
// Matches shadcn sidebar-menu-sub-button styling
//
// Example:
//
//	sidebar.SidebarMenuSubButton("Settings", "/settings", false)
func SidebarMenuSubButton(label string, href string, active bool) g.Node {
	classes := "flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground outline-none ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-sidebar-accent-foreground"
	if active {
		classes += " bg-sidebar-accent text-sidebar-accent-foreground font-medium"
	}

	attrs := []g.Node{
		g.Attr("href", href),
		g.Attr("data-slot", "sidebar-menu-sub-button"),
		html.Class(classes),
	}

	if active {
		attrs = append(attrs, g.Attr("data-active", "true"))
	}

	return html.A(
		g.Group(attrs),
		html.Span(g.Text(label)),
	)
}

// SidebarMenuBadge creates a badge for menu items
//
// Example:
//
//	sidebar.SidebarMenuBadge("12")
func SidebarMenuBadge(text string) g.Node {
	return html.Span(
		html.Class("inline-flex items-center justify-center rounded-md bg-sidebar-primary px-2 py-0.5 text-xs font-medium text-sidebar-primary-foreground"),
		g.Text(text),
	)
}

// SidebarMenuAction creates an action button that appears on hover
//
// Example:
//
//	sidebar.SidebarMenuAction(icons.MoreVertical(), "Options")
func SidebarMenuAction(icon g.Node, label string, attrs ...g.Node) g.Node {
	allAttrs := []g.Node{
		g.Attr("type", "button"),
		g.Attr("x-data", "{}"),
		html.Class("absolute right-1 top-1.5 opacity-0 group-hover/menu-item:opacity-100 peer-data-[active=true]/menu-button:opacity-100 transition-opacity p-1 rounded-md hover:bg-sidebar-accent text-sidebar-foreground/70 hover:text-sidebar-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-0 outline-none focus-visible:opacity-100"),
		g.Attr("aria-label", label),
		g.Attr("tabindex", "0"),
		g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
	}
	allAttrs = append(allAttrs, attrs...)

	return html.Button(
		g.Group(allAttrs),
		icon,
	)
}

// SidebarInset creates the main content area with proper layout
// Uses Alpine store to access sidebar state (works as sibling to Sidebar)
//
// Example:
//
//	sidebar.SidebarInset(
//	    sidebar.SidebarInsetHeader(...),
//	    html.Main(...),
//	)
func SidebarInset(children ...g.Node) g.Node {
	return html.Div(
		html.Class("relative flex min-h-screen flex-1 flex-col bg-background transition-all duration-300"),
		g.Attr("x-data", "{}"),
		g.Attr(":class", `{
			'md:ml-64': $store.sidebar && !$store.sidebar.collapsed,
			'md:ml-16': $store.sidebar && $store.sidebar.collapsed
		}`),
		g.Group(children),
	)
}

// SidebarInsetHeader creates a header with sidebar trigger and breadcrumbs
//
// Example:
//
//	sidebar.SidebarInsetHeader(
//	    sidebar.SidebarTrigger(),
//	    breadcrumb.Breadcrumb(...),
//	)
func SidebarInsetHeader(children ...g.Node) g.Node {
	return html.Header(
		html.Class("sticky top-0 z-10 flex h-16 shrink-0 items-center gap-2 border-border border-b bg-background px-4"),
		g.Group(children),
	)
}

// SidebarTrigger creates a trigger button for mobile/desktop
// Uses Alpine store for toggling (works from anywhere on the page)
//
// Example:
//
//	sidebar.SidebarTrigger()
func SidebarTrigger() g.Node {
	return html.Button(
		g.Attr("type", "button"),
		html.Class("inline-flex items-center justify-center rounded-md p-2 text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground md:hidden focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 outline-none"),
		g.Attr("aria-label", "Toggle sidebar"),
		alpine.XOn("click", "if ($store.sidebar) { $store.sidebar.isMobile ? ($store.sidebar.mobileOpen = true) : ($store.sidebar.collapsed = !$store.sidebar.collapsed) }"),
		g.El("svg",
			html.Class("h-6 w-6"),
			g.Attr("xmlns", "http://www.w3.org/2000/svg"),
			g.Attr("fill", "none"),
			g.Attr("viewBox", "0 0 24 24"),
			g.Attr("stroke", "currentColor"),
			g.Attr("aria-hidden", "true"),
			g.El("path",
				g.Attr("stroke-linecap", "round"),
				g.Attr("stroke-linejoin", "round"),
				g.Attr("stroke-width", "2"),
				g.Attr("d", "M4 6h16M4 12h16M4 18h16"),
			),
		),
	)
}

// SidebarTriggerDesktop creates a desktop-only trigger
// Uses Alpine store for toggling (works from anywhere on the page)
//
// Example:
//
//	sidebar.SidebarTriggerDesktop()
func SidebarTriggerDesktop() g.Node {
	return html.Button(
		g.Attr("type", "button"),
		html.Class("hidden md:inline-flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 outline-none"),
		g.Attr("aria-label", "Toggle sidebar"),
		alpine.XOn("click", "if ($store.sidebar) { $store.sidebar.collapsed = !$store.sidebar.collapsed }"),
		g.El("svg",
			html.Class("h-5 w-5"),
			g.Attr("xmlns", "http://www.w3.org/2000/svg"),
			g.Attr("fill", "none"),
			g.Attr("viewBox", "0 0 24 24"),
			g.Attr("stroke", "currentColor"),
			g.Attr("aria-hidden", "true"),
			g.El("path",
				g.Attr("stroke-linecap", "round"),
				g.Attr("stroke-linejoin", "round"),
				g.Attr("stroke-width", "2"),
				g.Attr("d", "M3 4h18M3 12h18M3 20h18"),
			),
		),
	)
}

// SidebarRail creates an invisible rail for better collapse interaction
// Uses Alpine store for toggling
//
// Example:
//
//	sidebar.SidebarRail()
func SidebarRail() g.Node {
	return html.Button(
		g.Attr("type", "button"),
		g.Attr("x-data", "{}"),
		html.Class("absolute inset-y-0 -right-4 z-20 hidden w-4 cursor-pointer md:block focus-visible:ring-2 focus-visible:ring-ring outline-none"),
		g.Attr("aria-label", "Toggle sidebar"),
		alpine.XOn("click", "if ($store.sidebar && $store.sidebar.collapsible) { $store.sidebar.collapsed = !$store.sidebar.collapsed }"),
		g.Attr("x-show", "$store.sidebar && !$store.sidebar.isMobile"),
	)
}

// SidebarMenuTooltip wraps a menu button with a tooltip that shows when sidebar is collapsed
//
// Example:
//
//	sidebar.SidebarMenuTooltip("Dashboard",
//	    sidebar.SidebarMenuButton("Dashboard", ...),
//	)
func SidebarMenuTooltip(label string, children ...g.Node) g.Node {
	return html.Div(
		html.Class("relative group/tooltip"),
		g.Attr("x-data", "{ showTooltip: false }"),

		// Trigger area
		html.Div(
			alpine.XOn("mouseenter", "showTooltip = $store.sidebar && $store.sidebar.collapsed && !$store.sidebar.isMobile"),
			alpine.XOn("mouseleave", "showTooltip = false"),
			g.Group(children),
		),

		// Tooltip popup
		html.Div(
			g.Attr("x-show", "showTooltip"),
			g.Attr("x-transition:enter", "transition ease-out duration-100"),
			g.Attr("x-transition:enter-start", "opacity-0 scale-95"),
			g.Attr("x-transition:enter-end", "opacity-100 scale-100"),
			g.Attr("x-transition:leave", "transition ease-in duration-75"),
			g.Attr("x-transition:leave-start", "opacity-100 scale-100"),
			g.Attr("x-transition:leave-end", "opacity-0 scale-95"),
			html.Class("absolute left-full ml-2 top-1/2 -translate-y-1/2 z-50 px-3 py-1.5 text-sm bg-popover text-popover-foreground rounded-md shadow-md border pointer-events-none whitespace-nowrap"),
			g.Attr("role", "tooltip"),
			g.Text(label),
		),
	)
}

// SidebarInput creates a search/filter input for the sidebar
//
// Example:
//
//	sidebar.SidebarInput("Search...", "search-query")
func SidebarInput(placeholder string, name string, attrs ...g.Node) g.Node {
	inputAttrs := []g.Node{
		g.Attr("type", "text"),
		g.Attr("name", name),
		g.Attr("placeholder", placeholder),
		html.Class("h-8 w-full bg-background border border-sidebar-border rounded-md px-3 text-sm shadow-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-0 outline-none"),
		g.Attr("x-data", "{}"),
		g.Attr("x-show", "$store.sidebar && (!$store.sidebar.collapsed || $store.sidebar.isMobile)"),
	}
	inputAttrs = append(inputAttrs, attrs...)

	return html.Div(
		html.Class("px-2 py-2"),
		html.Input(inputAttrs...),
	)
}

// SidebarSeparator creates a visual separator in the sidebar
//
// Example:
//
//	sidebar.SidebarSeparator()
func SidebarSeparator() g.Node {
	return html.Div(
		html.Class("mx-2 my-2 h-px bg-sidebar-border"),
		g.Attr("role", "separator"),
		g.Attr("aria-orientation", "horizontal"),
	)
}
