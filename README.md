# ForgeUI

ForgeUI is an SSR-first UI component library for Go, built on gomponents with Tailwind CSS styling and shadcn-inspired design patterns.

## Features

- ✅ **SSR-First**: Pure Go component rendering with no client-side dependencies
- ✅ **Type-Safe**: Full Go type safety with functional options pattern
- ✅ **CVA**: Class Variance Authority for flexible variant management
- ✅ **Tailwind CSS**: Utility-first CSS styling
- ✅ **Comprehensive Components**: 25+ production-ready components

## Installation

```bash
go get github.com/xraph/forgeui
```

## Quick Start

```go
package main

import (
    g "github.com/maragudk/gomponents"
    "github.com/maragudk/gomponents/html"
    
    "github.com/xraph/forgeui"
    "github.com/xraph/forgeui/components/button"
    "github.com/xraph/forgeui/components/card"
    "github.com/xraph/forgeui/primitives"
)

func main() {
    // Create a simple card with button
    page := html.Html(
        html.Head(
            html.Link(html.Rel("stylesheet"), html.Href("/static/styles.css")),
        ),
        html.Body(
            primitives.Container(
                card.Card(
                    card.Header(
                        card.Title("Welcome to ForgeUI"),
                        card.Description("Build beautiful UIs with Go"),
                    ),
                    card.Content(
                        primitives.VStack("4",
                            primitives.Text(
                                primitives.TextChildren(g.Text("ForgeUI provides type-safe, composable UI components.")),
                            ),
                            button.Primary(
                                g.Text("Get Started"),
                                button.WithSize(forgeui.SizeLG),
                            ),
                        ),
                    ),
                ),
            ),
        ),
    )
}
```

## Components

### Phase 01: Foundation
- ✅ Core types (Size, Variant, Radius)
- ✅ CVA (Class Variance Authority)
- ✅ Node builder with fluent API
- ✅ Props system with functional options
- ✅ App configuration

### Phase 02: Primitives
- ✅ **Box** - Polymorphic container
- ✅ **Flex** - Flexbox layout
- ✅ **Grid** - CSS Grid layout
- ✅ **Stack** - VStack/HStack helpers
- ✅ **Center** - Centered container
- ✅ **Container** - Responsive container
- ✅ **Text** - Typography primitive
- ✅ **Spacer** - Flexible spacer

### Phase 03: Core Components
- ✅ **Button** - All variants with button groups
- ✅ **Card** - Compound component (Header, Title, Content, Footer)
- ✅ **Badge** - Status and labels
- ✅ **Avatar** - User avatars with fallback
- ✅ **Alert** - Alert messages
- ✅ **Separator** - Horizontal/vertical separators
- ✅ **Spinner** - Loading indicator
- ✅ **Skeleton** - Loading placeholders
- ✅ **Progress** - Progress bars

### Phase 04: Form Components
- ✅ **Label** - Form labels
- ✅ **Input** - Text inputs with variants
- ✅ **Textarea** - Multi-line text input
- ✅ **Checkbox** - Checkbox inputs
- ✅ **Radio** - Radio buttons with groups
- ✅ **Switch** - Toggle switches
- ✅ **Select** - Native select dropdowns
- ✅ **Slider** - Range sliders
- ✅ **Form** - Form wrapper with validation helpers

## Usage Examples

### Button Variants

```go
// Primary button
button.Primary(g.Text("Save"))

// Destructive button
button.Destructive(
    g.Text("Delete"),
    button.WithSize(forgeui.SizeLG),
)

// Icon button
button.IconButton(
    g.Text("×"),
    button.WithVariant(forgeui.VariantGhost),
)

// Button group
button.Group(
    []button.GroupOption{button.WithGap("2")},
    button.Primary(g.Text("Save")),
    button.Secondary(g.Text("Cancel")),
)
```

### Card Component

```go
card.Card(
    card.Header(
        card.Title("Settings"),
        card.Description("Manage your account settings"),
    ),
    card.Content(
        // Your content here
    ),
    card.Footer(
        button.Primary(g.Text("Save Changes")),
    ),
)
```

### Form Example

```go
form.Form(
    []form.Option{
        form.WithAction("/submit"),
        form.WithMethod("POST"),
    },
    input.Field(
        "Email",
        []input.Option{
            input.WithType("email"),
            input.WithName("email"),
            input.WithPlaceholder("Enter your email"),
            input.Required(),
        },
        input.FormDescription("We'll never share your email."),
    ),
    checkbox.Checkbox(
        checkbox.WithName("subscribe"),
        checkbox.WithID("subscribe"),
    ),
    button.Primary(
        g.Text("Submit"),
        button.WithType("submit"),
    ),
)
```

### Layout with Primitives

```go
primitives.Container(
    primitives.VStack("8",
        primitives.Text(
            primitives.TextAs("h1"),
            primitives.TextSize("text-4xl"),
            primitives.TextWeight("font-bold"),
            primitives.TextChildren(g.Text("Dashboard")),
        ),
        primitives.Grid(
            primitives.GridCols(1),
            primitives.GridColsMD(2),
            primitives.GridColsLG(3),
            primitives.GridGap("6"),
            primitives.GridChildren(
                card.Card(/* ... */),
                card.Card(/* ... */),
                card.Card(/* ... */),
            ),
        ),
    ),
)
```

## Architecture

ForgeUI follows a strict layered architecture to avoid cyclic dependencies:

```
Phase 01 (Foundation)
    ↓
Phase 02 (Primitives)
    ↓
Phase 03 (Core Components)
    ↓
Phase 04 (Form Components)
```

- **No sibling imports**: Components cannot import from each other (except form → button)
- **Unidirectional flow**: Each phase only imports from phases below it
- **Type safety**: Full Go type checking at compile time

## Testing

ForgeUI includes comprehensive test coverage:

- **Phase 01**: 92.4% coverage (Foundation)
- **Phase 02**: 95.9% coverage (Primitives)
- **Phase 03**: 66.7%-94.3% coverage (Core Components)

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Configuration

```go
app := forgeui.New(
    forgeui.WithDebug(true),
    forgeui.WithThemeName("default"),
    forgeui.WithStaticPath("/static"),
    forgeui.WithDefaultSize(forgeui.SizeMD),
)

if err := app.Initialize(context.Background()); err != nil {
    log.Fatal(err)
}
```

## Contributing

ForgeUI is designed to be extended with additional components and features. Follow the existing patterns:

1. Use CVA for variant management
2. Follow functional options pattern
3. Write comprehensive tests
4. Document all exported functions
5. Respect the architecture layers

## License

MIT License

## Credits

- Built with [gomponents](https://github.com/maragudk/gomponents)
- Inspired by [shadcn/ui](https://ui.shadcn.com/)
- Styled with [Tailwind CSS](https://tailwindcss.com/)

