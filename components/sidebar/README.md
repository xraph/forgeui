# Sidebar Component

A collapsible sidebar navigation component with provider pattern integration for state management. Fully responsive with mobile support, animations, and customizable styling.

## Features

- ✅ **Provider Pattern**: Uses ForgeUI's provider pattern for state management
- ✅ **Responsive**: Adapts to mobile and desktop with different behaviors
- ✅ **Collapsible**: Can be collapsed to icon-only mode on desktop
- ✅ **Mobile Drawer**: Full drawer experience on mobile with backdrop
- ✅ **Animations**: Smooth transitions and animations
- ✅ **Nested Menus**: Support for submenus and collapsible groups
- ✅ **Customizable**: Extensive options for styling and behavior
- ✅ **Accessible**: ARIA labels and keyboard navigation
- ✅ **Keyboard Shortcuts**: Cmd/Ctrl+B to toggle sidebar (customizable)
- ✅ **State Persistence**: Remembers collapsed state across sessions
- ✅ **Multiple Variants**: Sidebar, floating, and inset styles
- ✅ **Collapsible Modes**: Offcanvas, icon-only, or non-collapsible
- ✅ **Tooltips**: Auto-show tooltips when sidebar is collapsed
- ✅ **Focus Management**: Comprehensive focus rings and keyboard navigation
- ✅ **Menu Variants**: Default and outline button styles
- ✅ **Menu Sizes**: Small, default, and large button sizes

## Basic Usage

```go
import "github.com/xraph/forgeui/components/sidebar"

sidebar.Sidebar(
    sidebar.SidebarHeader(
        g.Text("My App"),
    ),
    sidebar.SidebarContent(
        sidebar.SidebarGroup(
            sidebar.SidebarGroupLabel("Navigation"),
            sidebar.SidebarMenu(
                sidebar.SidebarMenuItem(
                    sidebar.SidebarMenuButton(
                        "Dashboard",
                        sidebar.WithMenuHref("/dashboard"),
                        sidebar.WithMenuIcon(icons.LayoutDashboard()),
                        sidebar.WithMenuActive(),
                    ),
                ),
            ),
        ),
    ),
    sidebar.SidebarFooter(
        g.Text("© 2024"),
    ),
    sidebar.SidebarToggle(),
)
```

## Provider Integration

The sidebar uses ForgeUI's provider pattern to manage state. This allows child components and even components outside the sidebar to access its state without prop drilling.

### Sidebar State

The sidebar provider exposes the following state:

- `collapsed` (bool): Whether the sidebar is collapsed (desktop only)
- `mobileOpen` (bool): Whether the sidebar is open on mobile
- `collapsible` (bool): Whether the sidebar can be collapsed
- `isMobile` (bool): Whether the viewport is mobile-sized

### Sidebar Methods

The sidebar provider exposes these methods:

- `toggle()`: Toggle collapsed state
- `openMobile()`: Open sidebar on mobile
- `closeMobile()`: Close sidebar on mobile
- `setCollapsed(value)`: Set collapsed state directly

### Accessing Sidebar State

```go
// In a component that needs to react to sidebar state
html.Div(
    g.Attr(":class", `{
        'ml-64': `+primitives.ProviderValue("sidebar", "!isMobile && !collapsed")+`,
        'ml-16': `+primitives.ProviderValue("sidebar", "!isMobile && collapsed")+`,
        'ml-0': `+primitives.ProviderValue("sidebar", "isMobile")+`
    }`),
    // Your content
)
```

## Components

### Sidebar

Main sidebar container with provider integration.

```go
sidebar.Sidebar(children...)
sidebar.SidebarWithOptions(options, children...)
```

**Options:**
- `WithDefaultCollapsed(bool)`: Set initial collapsed state
- `WithCollapsible(bool)`: Allow sidebar to be collapsed
- `WithSide(forgeui.Side)`: Position (left or right)
- `WithSidebarClass(string)`: Custom CSS classes
- `WithSidebarAttrs(...g.Node)`: Custom attributes
- `WithVariant(SidebarVariant)`: Visual variant (sidebar, floating, inset)
- `WithCollapsibleMode(SidebarCollapsibleMode)`: Collapse behavior (offcanvas, icon, none)
- `WithKeyboardShortcut(string)`: Set keyboard shortcut key (default: "b")
- `WithKeyboardShortcutEnabled(bool)`: Enable/disable keyboard shortcuts
- `WithPersistState(bool)`: Enable state persistence
- `WithStorageKey(string)`: Custom storage key for persistence
- `WithStorageType(string)`: Storage type ("localStorage" or "sessionStorage")

### SidebarHeader

Header section of the sidebar (typically for logo/branding).

```go
sidebar.SidebarHeader(
    html.Img(g.Attr("src", "/logo.svg")),
    g.Text("My App"),
)
```

### SidebarContent

Main content area containing navigation items.

```go
sidebar.SidebarContent(
    // Navigation groups and menus
)
```

### SidebarFooter

Footer section of the sidebar (typically for user profile or settings).

```go
sidebar.SidebarFooter(
    // Footer content
)
```

### SidebarToggle

Toggle button for collapsing/expanding the sidebar.

```go
sidebar.SidebarToggle()
```

### SidebarGroup

Groups related menu items together.

```go
sidebar.SidebarGroup(
    sidebar.SidebarGroupLabel("Settings"),
    sidebar.SidebarMenu(...),
)
```

### SidebarGroupCollapsible

Collapsible group with animation.

```go
sidebar.SidebarGroupCollapsible(
    []sidebar.SidebarGroupOption{
        sidebar.WithGroupKey("projects"),
        sidebar.WithGroupDefaultOpen(true),
    },
    sidebar.SidebarGroupLabelCollapsible("projects", "Projects", icon),
    sidebar.SidebarGroupContent("projects",
        sidebar.SidebarMenu(...),
    ),
)
```

### SidebarMenu

Container for menu items.

```go
sidebar.SidebarMenu(
    sidebar.SidebarMenuItem(...),
    sidebar.SidebarMenuItem(...),
)
```

### SidebarMenuItem

Individual menu item container.

```go
sidebar.SidebarMenuItem(
    sidebar.SidebarMenuButton(...),
)
```

### SidebarMenuButton

Menu button/link with icon, label, and badge support.

```go
sidebar.SidebarMenuButton(
    "Dashboard",
    sidebar.WithMenuHref("/dashboard"),
    sidebar.WithMenuIcon(icons.LayoutDashboard()),
    sidebar.WithMenuActive(),
    sidebar.WithMenuBadge(sidebar.SidebarMenuBadge("5")),
)
```

**Options:**
- `WithMenuHref(string)`: Link URL
- `WithMenuActive()`: Mark as active
- `WithMenuIcon(g.Node)`: Add icon
- `WithMenuBadge(g.Node)`: Add badge
- `WithMenuTooltip(string)`: Tooltip text
- `WithMenuAsButton()`: Render as button instead of link
- `WithMenuClass(string)`: Custom classes
- `WithMenuAttrs(...g.Node)`: Custom attributes

### SidebarMenuSub

Submenu container for nested items.

```go
sidebar.SidebarMenuSub(
    sidebar.SidebarMenuSubItem(
        sidebar.SidebarMenuSubButton("Overview", "/overview", false),
    ),
)
```

## Layout Components

### SidebarLayout

Root layout container.

```go
sidebar.SidebarLayout(
    sidebar.Sidebar(...),
    html.Main(...),
)
```

### SidebarLayoutContent

Content area that adjusts margins based on sidebar state.

```go
sidebar.SidebarLayoutContent(
    // Your page content
)
```

### SidebarInset

Alternative layout with inset content area.

```go
sidebar.SidebarInset(
    sidebar.SidebarInsetHeader(
        sidebar.SidebarTrigger(),
        breadcrumb.Breadcrumb(...),
    ),
    html.Main(...),
)
```

## Complete Example

```go
package main

import (
    "github.com/xraph/forgeui/components/sidebar"
    "github.com/xraph/forgeui/icons"
    g "maragu.dev/gomponents"
    "maragu.dev/gomponents/html"
)

func DashboardLayout() g.Node {
    return sidebar.SidebarLayout(
        // Sidebar
        sidebar.SidebarWithOptions(
            []sidebar.SidebarOption{
                sidebar.WithDefaultCollapsed(false),
                sidebar.WithCollapsible(true),
            },
            // Header
            sidebar.SidebarHeader(
                html.Div(
                    html.Class("flex items-center gap-2"),
                    html.Img(
                        g.Attr("src", "/logo.svg"),
                        html.Class("h-8 w-8"),
                    ),
                    html.Span(
                        html.Class("font-bold text-lg"),
                        g.Text("My App"),
                    ),
                ),
            ),
            
            // Content
            sidebar.SidebarContent(
                // Main navigation
                sidebar.SidebarGroup(
                    sidebar.SidebarGroupLabel("Main"),
                    sidebar.SidebarMenu(
                        sidebar.SidebarMenuItem(
                            sidebar.SidebarMenuButton(
                                "Dashboard",
                                sidebar.WithMenuHref("/dashboard"),
                                sidebar.WithMenuIcon(icons.LayoutDashboard()),
                                sidebar.WithMenuActive(),
                            ),
                        ),
                        sidebar.SidebarMenuItem(
                            sidebar.SidebarMenuButton(
                                "Analytics",
                                sidebar.WithMenuHref("/analytics"),
                                sidebar.WithMenuIcon(icons.Activity()),
                            ),
                        ),
                        sidebar.SidebarMenuItem(
                            sidebar.SidebarMenuButton(
                                "Messages",
                                sidebar.WithMenuHref("/messages"),
                                sidebar.WithMenuIcon(icons.Mail()),
                                sidebar.WithMenuBadge(sidebar.SidebarMenuBadge("12")),
                            ),
                        ),
                    ),
                ),
                
                // Collapsible projects section
                sidebar.SidebarGroupCollapsible(
                    []sidebar.SidebarGroupOption{
                        sidebar.WithGroupKey("projects"),
                        sidebar.WithGroupDefaultOpen(true),
                    },
                    sidebar.SidebarGroupLabelCollapsible(
                        "projects",
                        "Projects",
                        icons.FolderKanban(),
                    ),
                    sidebar.SidebarGroupContent("projects",
                        sidebar.SidebarMenu(
                            sidebar.SidebarMenuItem(
                                sidebar.SidebarMenuButton(
                                    "Website Redesign",
                                    sidebar.WithMenuHref("/projects/website"),
                                    sidebar.WithMenuIcon(icons.FolderKanban()),
                                ),
                                // Submenu
                                sidebar.SidebarMenuSub(
                                    sidebar.SidebarMenuSubItem(
                                        sidebar.SidebarMenuSubButton(
                                            "Overview",
                                            "/projects/website/overview",
                                            false,
                                        ),
                                    ),
                                    sidebar.SidebarMenuSubItem(
                                        sidebar.SidebarMenuSubButton(
                                            "Tasks",
                                            "/projects/website/tasks",
                                            true,
                                        ),
                                    ),
                                ),
                            ),
                        ),
                    ),
                ),
                
                // Settings
                sidebar.SidebarGroup(
                    sidebar.SidebarGroupLabel("Settings"),
                    sidebar.SidebarMenu(
                        sidebar.SidebarMenuItem(
                            sidebar.SidebarMenuButton(
                                "Account",
                                sidebar.WithMenuHref("/settings/account"),
                                sidebar.WithMenuIcon(icons.User()),
                            ),
                        ),
                    ),
                ),
            ),
            
            // Footer
            sidebar.SidebarFooter(
                html.Div(
                    html.Class("flex items-center gap-2"),
                    html.Img(
                        g.Attr("src", "/avatar.jpg"),
                        html.Class("h-8 w-8 rounded-full"),
                    ),
                    html.Span(g.Text("John Doe")),
                ),
            ),
            
            // Toggle button
            sidebar.SidebarToggle(),
        ),
        
        // Main content
        sidebar.SidebarInset(
            sidebar.SidebarInsetHeader(
                sidebar.SidebarTrigger(),
                html.H1(g.Text("Dashboard")),
            ),
            html.Main(
                html.Class("p-6"),
                // Your page content
            ),
        ),
    )
}
```

## Styling

The sidebar uses Tailwind CSS with custom theme colors. Make sure your theme includes:

```css
:root {
  --sidebar: 0 0% 98%;
  --sidebar-foreground: 240 5.3% 26.1%;
  --sidebar-primary: 240 5.9% 10%;
  --sidebar-primary-foreground: 0 0% 98%;
  --sidebar-accent: 240 4.8% 95.9%;
  --sidebar-accent-foreground: 240 5.9% 10%;
  --sidebar-border: 220 13% 91%;
}

.dark {
  --sidebar: 240 5.9% 10%;
  --sidebar-foreground: 240 4.8% 95.9%;
  --sidebar-primary: 0 0% 98%;
  --sidebar-primary-foreground: 240 5.9% 10%;
  --sidebar-accent: 240 3.7% 15.9%;
  --sidebar-accent-foreground: 240 4.8% 95.9%;
  --sidebar-border: 240 3.7% 15.9%;
}
```

## Migration from Menu Package

If you're migrating from the old `components/menu` sidebar:

### Before

```go
import "github.com/xraph/forgeui/components/menu"

menu.Sidebar(...)
menu.SidebarHeader(...)
```

### After

```go
import "github.com/xraph/forgeui/components/sidebar"

sidebar.Sidebar(...)
sidebar.SidebarHeader(...)
```

**Key Changes:**
1. Import path changed from `components/menu` to `components/sidebar`
2. All functions prefixed with `menu.` now use `sidebar.`
3. State management moved from Alpine stores to provider pattern
4. Backward compatibility maintained via Alpine store sync

## Backward Compatibility

The sidebar maintains backward compatibility with Alpine stores. If you have code that accesses `$store.sidebar`, it will continue to work:

```javascript
// Still works
$store.sidebar.collapsed
$store.sidebar.mobileOpen

// But prefer the new provider pattern
$el.closest('[data-provider="sidebar"]').__x.$data.collapsed
```

## Enhanced Features

### Keyboard Shortcuts

The sidebar supports keyboard shortcuts for quick toggling. By default, pressing `Cmd+B` (Mac) or `Ctrl+B` (Windows/Linux) will toggle the sidebar.

```go
// Default - Cmd/Ctrl+B
sidebar.Sidebar(...)

// Custom key
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        sidebar.WithKeyboardShortcut("k"), // Cmd/Ctrl+K
    },
    // children...
)

// Disable keyboard shortcuts
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        sidebar.WithKeyboardShortcutEnabled(false),
    },
    // children...
)
```

### State Persistence

The sidebar can remember its collapsed state across browser sessions using localStorage or sessionStorage.

```go
// Enable persistence with default settings (localStorage)
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        sidebar.WithPersistState(true),
    },
    // children...
)

// Custom storage key and type
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        sidebar.WithPersistState(true),
        sidebar.WithStorageKey("my_app_sidebar"),
        sidebar.WithStorageType("sessionStorage"),
    },
    // children...
)
```

### Sidebar Variants

Choose from different visual styles for your sidebar:

**Sidebar Variant (Default)**
- Standard sidebar flush with viewport edge
- No rounded corners or shadows

**Floating Variant**
- Rounded corners and shadow
- Margin from viewport edges
- Modern, card-like appearance

**Inset Variant**
- Similar to floating but more subtle
- Works well with inset content layouts

```go
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        sidebar.WithVariant(sidebar.SidebarVariantFloating),
    },
    // children...
)
```

### Collapsible Modes

Control how the sidebar collapses:

**Offcanvas Mode (Default)**
- Slides completely off screen when collapsed
- Similar to mobile drawer behavior

**Icon Mode**
- Shrinks to icon-only width (64px)
- Labels hidden, only icons visible
- Tooltips appear on hover

**None Mode**
- Sidebar cannot be collapsed
- Toggle button hidden

```go
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        sidebar.WithCollapsibleMode(sidebar.CollapsibleIcon),
    },
    // children...
)
```

### Tooltips

Tooltips automatically appear when hovering over menu items in collapsed/icon mode:

```go
sidebar.SidebarMenuButton(
    "Dashboard",
    sidebar.WithMenuHref("/dashboard"),
    sidebar.WithMenuIcon(icons.LayoutDashboard()),
    sidebar.WithMenuTooltip("Go to Dashboard"), // Shows when collapsed
)

// Or wrap any element with a tooltip
sidebar.SidebarMenuTooltip(
    "Dashboard",
    sidebar.SidebarMenuButton(...),
)
```

### Menu Button Variants

Choose from different button styles:

```go
// Default variant
sidebar.SidebarMenuButton(
    "Dashboard",
    sidebar.WithMenuVariant(sidebar.MenuButtonDefault),
)

// Outline variant (with border and shadow)
sidebar.SidebarMenuButton(
    "Settings",
    sidebar.WithMenuVariant(sidebar.MenuButtonOutline),
)
```

### Menu Button Sizes

Adjust button sizes for different hierarchy levels:

```go
// Small (28px height)
sidebar.SidebarMenuButton(
    "Sub Item",
    sidebar.WithMenuSize(sidebar.MenuButtonSizeSmall),
)

// Default (32px height)
sidebar.SidebarMenuButton(
    "Main Item",
    sidebar.WithMenuSize(sidebar.MenuButtonSizeDefault),
)

// Large (48px height)
sidebar.SidebarMenuButton(
    "Featured",
    sidebar.WithMenuSize(sidebar.MenuButtonSizeLarge),
)
```

### Additional Components

**SidebarInput** - Search/filter input for sidebar:

```go
sidebar.SidebarInput("Search...", "sidebar-search")
```

**SidebarSeparator** - Visual separator:

```go
sidebar.SidebarSeparator()
```

### Complete Example with All Features

```go
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        // Visual style
        sidebar.WithVariant(sidebar.SidebarVariantFloating),
        sidebar.WithCollapsibleMode(sidebar.CollapsibleIcon),
        
        // Behavior
        sidebar.WithDefaultCollapsed(false),
        sidebar.WithCollapsible(true),
        sidebar.WithSide(forgeui.SideLeft),
        
        // Keyboard shortcuts
        sidebar.WithKeyboardShortcut("b"),
        sidebar.WithKeyboardShortcutEnabled(true),
        
        // State persistence
        sidebar.WithPersistState(true),
        sidebar.WithStorageKey("myapp_sidebar"),
        sidebar.WithStorageType("localStorage"),
    },
    sidebar.SidebarHeader(
        html.Img(g.Attr("src", "/logo.svg")),
        g.Text("My App"),
    ),
    sidebar.SidebarContent(
        // Search input
        sidebar.SidebarInput("Search...", "search"),
        sidebar.SidebarSeparator(),
        
        // Main navigation
        sidebar.SidebarGroup(
            sidebar.SidebarGroupLabel("Main"),
            sidebar.SidebarMenu(
                sidebar.SidebarMenuItem(
                    sidebar.SidebarMenuButton(
                        "Dashboard",
                        sidebar.WithMenuHref("/dashboard"),
                        sidebar.WithMenuIcon(icons.LayoutDashboard()),
                        sidebar.WithMenuActive(),
                        sidebar.WithMenuTooltip("Go to Dashboard"),
                        sidebar.WithMenuVariant(sidebar.MenuButtonDefault),
                        sidebar.WithMenuSize(sidebar.MenuButtonSizeDefault),
                    ),
                ),
            ),
        ),
        
        sidebar.SidebarSeparator(),
        
        // Settings
        sidebar.SidebarGroup(
            sidebar.SidebarGroupLabel("Settings"),
            sidebar.SidebarMenu(
                sidebar.SidebarMenuItem(
                    sidebar.SidebarMenuButton(
                        "Preferences",
                        sidebar.WithMenuHref("/settings"),
                        sidebar.WithMenuIcon(icons.Settings()),
                        sidebar.WithMenuTooltip("Configure Settings"),
                        sidebar.WithMenuVariant(sidebar.MenuButtonOutline),
                        sidebar.WithMenuSize(sidebar.MenuButtonSizeSmall),
                    ),
                ),
            ),
        ),
    ),
    sidebar.SidebarFooter(
        g.Text("© 2024 My App"),
    ),
    sidebar.SidebarToggle(),
)
```

## Accessibility

The sidebar includes comprehensive accessibility features:

### ARIA Attributes
- `role="complementary"` on sidebar container
- `role="menuitem"` on menu buttons
- `role="navigation"` on menu containers
- `role="separator"` on separators
- `role="tooltip"` on tooltips
- `aria-label` on all interactive elements
- `aria-expanded` on collapsible elements
- `aria-current="page"` on active menu items
- `aria-hidden` on decorative elements

### Keyboard Navigation
- **Tab**: Navigate through menu items
- **Enter/Space**: Activate menu items
- **Arrow Up/Down**: Navigate within menu
- **Home**: Jump to first menu item
- **End**: Jump to last menu item
- **Cmd/Ctrl+B**: Toggle sidebar (customizable)

### Focus Management
- Visible focus rings on all interactive elements
- Focus trapped in mobile drawer
- Focus restored after sidebar toggle
- Keyboard-accessible tooltips

### Screen Reader Support
- Live region announces sidebar state changes
- Descriptive labels on all buttons
- Semantic HTML structure
- Proper heading hierarchy

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- Requires Alpine.js 3.x
- Requires Tailwind CSS 3.x

## Migration Guide

### Upgrading from Previous Version

If you're upgrading from a previous version of the sidebar component, here are the changes:

**New Features (Backward Compatible):**
- All new options are opt-in
- Existing sidebars will work without changes
- New variants and modes available
- Enhanced accessibility (automatic)

**Breaking Changes:**
- None - fully backward compatible

**Recommended Updates:**
1. Enable keyboard shortcuts (enabled by default)
2. Consider adding tooltips for collapsed state
3. Enable state persistence for better UX
4. Update to newer variant styles if desired
5. Add ARIA labels to custom menu items

**Example Migration:**

```go
// Old (still works)
sidebar.Sidebar(
    sidebar.SidebarHeader(g.Text("App")),
    sidebar.SidebarContent(...),
)

// New (with enhancements)
sidebar.SidebarWithOptions(
    []sidebar.SidebarOption{
        sidebar.WithVariant(sidebar.SidebarVariantFloating),
        sidebar.WithPersistState(true),
    },
    sidebar.SidebarHeader(g.Text("App")),
    sidebar.SidebarContent(
        sidebar.SidebarInput("Search...", "search"),
        sidebar.SidebarSeparator(),
        // ... menu items with tooltips
    ),
)
```

## See Also

- [Provider Pattern Documentation](../../primitives/PROVIDER_PATTERN.md)
- [Menu Component](../menu/menu.go) - For standalone menus
- [Icons](../../icons/README.md) - For sidebar icons

