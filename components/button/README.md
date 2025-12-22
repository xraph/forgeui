# Button Component

Versatile button component with multiple variants, sizes, and states.

## Features

- 🎨 **6 Variants**: Default, Secondary, Destructive, Outline, Ghost, Link
- 📏 **4 Sizes**: SM, MD (default), LG, Icon
- ⚡ **States**: Disabled, Loading
- 🔧 **Flexible**: Custom classes, attributes, and children
- ♿ **Accessible**: Proper ARIA attributes and keyboard navigation

## Basic Usage

```go
import (
    "github.com/xraph/forgeui/components/button"
    g "maragu.dev/gomponents"
)

// Simple button
btn := button.Button(g.Text("Click me"))

// With options
btn := button.Button(
    g.Text("Save"),
    button.WithVariant(forgeui.VariantDefault),
    button.WithSize(forgeui.SizeLG),
)
```

## Variants

### Default (Primary)
```go
button.Button(g.Text("Primary Action"))
// or
button.Primary(g.Text("Primary Action"))
```

### Secondary
```go
button.Secondary(g.Text("Secondary Action"))
```

### Destructive
```go
button.Destructive(g.Text("Delete"))
```

### Outline
```go
button.Outline(g.Text("Cancel"))
```

### Ghost
```go
button.Ghost(g.Text("Skip"))
```

### Link
```go
button.Link(g.Text("Learn More"))
```

## Sizes

```go
// Small
button.Button(g.Text("Small"), button.WithSize(forgeui.SizeSM))

// Medium (default)
button.Button(g.Text("Medium"))

// Large
button.Button(g.Text("Large"), button.WithSize(forgeui.SizeLG))

// Icon (square)
button.IconButton(icons.X())
```

## States

### Disabled
```go
button.Button(
    g.Text("Disabled"),
    button.Disabled(),
)
```

### Loading
```go
button.Button(
    g.Text("Saving..."),
    button.Loading(),
)
```

## Button Types

```go
// Submit button (for forms)
button.Button(
    g.Text("Submit"),
    button.WithType("submit"),
)

// Reset button
button.Button(
    g.Text("Reset"),
    button.WithType("reset"),
)
```

## With Icons

```go
import "github.com/xraph/forgeui/icons"

// Icon + Text
button.Primary(
    g.Group([]g.Node{
        icons.Plus(icons.WithSize(16)),
        g.Text("Add Item"),
    }),
)

// Icon only
button.IconButton(
    icons.Settings(),
    button.WithVariant(forgeui.VariantGhost),
)
```

## Button Group

Group related buttons together:

```go
btnGroup := button.Group(
    []button.GroupOption{
        button.WithGap("2"),
    },
    button.Primary(g.Text("Save")),
    button.Secondary(g.Text("Cancel")),
)
```

## Custom Styling

```go
button.Button(
    g.Text("Custom"),
    button.WithClass("shadow-lg hover:shadow-xl"),
)
```

## Custom Attributes

```go
button.Button(
    g.Text("Track Me"),
    button.WithAttrs(
        g.Attr("data-analytics", "signup-button"),
        g.Attr("id", "signup-btn"),
    ),
)
```

## Complete Example

```go
func SaveForm() g.Node {
    return primitives.Flex(
        primitives.FlexGap("4"),
        primitives.FlexChildren(
            button.Primary(
                g.Group([]g.Node{
                    icons.Check(icons.WithSize(16)),
                    g.Text("Save Changes"),
                }),
                button.WithType("submit"),
            ),
            button.Outline(
                g.Text("Cancel"),
                button.WithType("button"),
            ),
            button.Destructive(
                g.Group([]g.Node{
                    icons.Trash(icons.WithSize(16)),
                    g.Text("Delete"),
                }),
            ),
        ),
    )
}
```

## Best Practices

1. **Use semantic variants**: Primary for main actions, Destructive for dangerous actions
2. **Add loading states**: Show feedback during async operations
3. **Icon + Text**: Combine icons with text for clarity
4. **Proper types**: Use `type="submit"` for form submissions
5. **Disabled vs Loading**: Use Loading() for async operations, Disabled() for unavailable actions
6. **Accessibility**: Add `aria-label` for icon-only buttons

## Styling

The button uses CVA (Class Variance Authority) for consistent styling:
- Base classes for layout and transitions
- Variant classes for colors and styles
- Size classes for dimensions
- Focus states for keyboard navigation
- Dark mode support built-in

