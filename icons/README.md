# Icons Package

The `icons` package provides a flexible icon system with comprehensive Lucide icon integration for ForgeUI.

## Features

- 🎨 **Customizable**: Size, color, and stroke width options
- 📦 **Lucide Icons**: 1600+ pre-built beautiful icons (auto-generated from latest Lucide)
- 🔧 **Flexible**: Easy to add custom icons
- ♿ **Accessible**: Inline SVG with proper attributes
- 🔄 **Auto-generated**: Easy to update with new Lucide releases

## Basic Usage

```go
import (
    "github.com/xraph/forgeui/icons"
    g "github.com/maragudk/gomponents"
)

// Use a pre-built icon
checkIcon := icons.Check()

// Customize icon properties
largeIcon := icons.Search(
    icons.WithSize(32),
    icons.WithColor("blue"),
)

// Use in a button
btn := button.Button(
    g.Group([]g.Node{
        icons.Plus(icons.WithSize(16)),
        g.Text("Add Item"),
    }),
)
```

## Available Icons

This package includes **1666 icons** from [Lucide Icons](https://lucide.dev/). All icons are auto-generated from the latest Lucide release.

### Icon Discovery

To find available icons:
1. Browse the [Lucide Icons website](https://lucide.dev/icons/)
2. Convert the icon name from kebab-case to PascalCase
   - Example: `arrow-right` → `ArrowRight()`
   - Example: `chevron-down` → `ChevronDown()`
   - Example: `circle-check` → `CircleCheck()`

### Commonly Used Icons

**Navigation:**
- `ChevronDown`, `ChevronUp`, `ChevronLeft`, `ChevronRight`
- `Menu`, `House`, `ExternalLink`, `ArrowLeft`, `ArrowRight`

**Actions:**
- `Plus`, `Minus`, `Check`, `X`
- `Edit`, `Trash`, `Copy`, `Save`
- `Download`, `Upload`, `Share`
- `Search`, `Filter`, `Settings`

**Status:**
- `CircleAlert`, `Info`, `CircleCheck`, `CircleX`
- `Loader`, `AlertTriangle`, `HelpCircle`

**User & Communication:**
- `User`, `Mail`, `Phone`, `MessageCircle`
- `Eye`, `EyeOff`, `Bell`, `Send`

**Media & Files:**
- `File`, `Folder`, `Image`, `Video`, `Music`
- `FileText`, `FilePlus`, `FolderOpen`

### Backward Compatibility

For compatibility with previous versions, common aliases are provided:
- `Home()` → maps to `House()`
- `AlertCircle()` → maps to `CircleAlert()`
- `CheckCircle()` → maps to `CircleCheck()`
- `XCircle()` → maps to `CircleX()`

## Customization Options

### Size

```go
// Default size is 24px
icons.Check(icons.WithSize(16))  // Small
icons.Check(icons.WithSize(32))  // Large
```

### Color

```go
// Default is "currentColor" (inherits text color)
icons.Check(icons.WithColor("red"))
icons.Check(icons.WithColor("#3b82f6"))
```

### Stroke Width

```go
// Default is 2
icons.Check(icons.WithStrokeWidth(1.5))  // Thinner
icons.Check(icons.WithStrokeWidth(3))    // Thicker
```

### Custom Classes

```go
icons.Check(icons.WithClass("text-green-500 hover:text-green-600"))
```

### Custom Attributes

```go
icons.Check(
    icons.WithAttrs(
        g.Attr("aria-label", "Success"),
        g.Attr("data-testid", "check-icon"),
    ),
)
```

## Creating Custom Icons

### Single Path Icon

```go
customIcon := icons.Icon("M5 12h14")  // SVG path data
```

### Multi-Path Icon

```go
customIcon := icons.MultiPathIcon([]string{
    "M18 6 6 18",
    "m6 6 12 12",
})
```

## Examples

### Icon in Button

```go
saveButton := button.Primary(
    g.Group([]g.Node{
        icons.Check(icons.WithSize(16)),
        g.Text("Save"),
    }),
)
```

### Icon-Only Button

```go
closeButton := button.IconButton(
    icons.X(icons.WithSize(16)),
    button.WithVariant(forgeui.VariantGhost),
)
```

### Status Indicator

```go
successMessage := alert.Alert(
    g.Group([]g.Node{
        icons.CheckCircle(
            icons.WithSize(20),
            icons.WithColor("green"),
        ),
        alert.AlertTitle("Success"),
        alert.AlertDescription("Your changes have been saved."),
    }),
    alert.WithVariant(forgeui.VariantSuccess),
)
```

### Loading Spinner

```go
loadingButton := button.Button(
    g.Group([]g.Node{
        icons.Loader(
            icons.WithSize(16),
            icons.WithClass("animate-spin"),
        ),
        g.Text("Loading..."),
    }),
    button.Disabled(),
)
```

## Regenerating Icons

The icon library is auto-generated from Lucide's official icon set. To update to the latest Lucide release:

```bash
# Option 1: Using go generate
cd icons
go generate

# Option 2: Run the generator directly
cd icons/internal/generate
go run main.go
```

The generator will:
1. Fetch the latest icon metadata from the Lucide CDN
2. Parse all SVG elements (paths, circles, rects, polygons, etc.)
3. Convert icon names from kebab-case to PascalCase
4. Generate Go functions for all icons
5. Write to `lucide_generated.go`

### Generator Details

- **Data Source**: `https://cdn.jsdelivr.net/npm/lucide-static@latest/`
- **Icon Count**: 1666 icons (as of generation)
- **File Size**: ~16,000 lines of generated Go code
- **Naming Convention**: 
  - `arrow-right` → `ArrowRight()`
  - `3d-box` → `ThreeDBox()`
  - Reserved names get `Icon` suffix: `option` → `OptionIcon()`

## Icon Reference

All icons are from [Lucide](https://lucide.dev/) - a beautiful, consistent icon set.

## Best Practices

1. **Use appropriate sizes**: 16px for buttons, 20-24px for standalone icons
2. **Inherit color**: Use `currentColor` (default) to match text color
3. **Add aria-labels**: For icon-only buttons, add descriptive labels
4. **Consistent stroke**: Stick to default stroke width (2) for consistency
5. **Semantic usage**: Choose icons that clearly represent their action

## Performance

Icons are rendered as inline SVG elements, which:
- ✅ Load instantly (no HTTP requests)
- ✅ Scale perfectly at any size
- ✅ Support CSS styling and animations
- ✅ Work in all browsers

