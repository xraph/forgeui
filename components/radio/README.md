# Radio Component

Radio button component with group support for single-selection inputs.

## Installation

```go
import "github.com/xraph/forgeui/components/radio"
```

## Basic Usage

```go
radioBtn := radio.Radio(
    radio.WithID("option1"),
    radio.WithName("choice"),
    radio.WithValue("1"),
)
```

## Props

- `ID` - Radio button ID
- `Name` - Radio button name (groups radios together)
- `Value` - Radio button value
- `Checked` - Checked state
- `Required` - Required state
- `Disabled` - Disabled state
- `Class` - Custom CSS classes
- `Attrs` - Additional HTML attributes

## Options

```go
radio.WithID(id string)
radio.WithName(name string)
radio.WithValue(value string)
radio.Checked()
radio.Required()
radio.Disabled()
radio.WithClass(class string)
radio.WithAttrs(attrs ...g.Node)
```

## Examples

### Basic Radio Button

```go
radioBtn := radio.Radio(
    radio.WithID("small"),
    radio.WithName("size"),
    radio.WithValue("small"),
)
```

### Checked Radio Button

```go
radioBtn := radio.Radio(
    radio.WithID("medium"),
    radio.WithName("size"),
    radio.WithValue("medium"),
    radio.Checked(),
)
```

### Disabled Radio Button

```go
radioBtn := radio.Radio(
    radio.WithID("disabled"),
    radio.WithName("size"),
    radio.WithValue("disabled"),
    radio.Disabled(),
)
```

## Radio Group

Use the `Group` component to create a group of radio buttons:

```go
sizeGroup := radio.Group(
    []radio.GroupOption{
        radio.WithGroupName("size"),
        radio.WithGroupValue("medium"),
    },
    radio.GroupItem("size-small", "Small", radio.WithValue("small"), radio.WithName("size")),
    radio.GroupItem("size-medium", "Medium", radio.WithValue("medium"), radio.WithName("size"), radio.Checked()),
    radio.GroupItem("size-large", "Large", radio.WithValue("large"), radio.WithName("size")),
)
```

### Group Options

```go
radio.WithGroupName(name string)        // Set group name
radio.WithGroupValue(value string)      // Set selected value
radio.WithGroupDisabled()               // Disable all radios
radio.WithGroupClass(class string)      // Custom classes
radio.WithGroupAttrs(attrs ...g.Node)   // Custom attributes
```

## With Label

Use with the `label` component for proper accessibility:

```go
html.Div(
    html.Class("flex items-center space-x-2"),
    radio.Radio(
        radio.WithID("option1"),
        radio.WithName("choice"),
        radio.WithValue("1"),
    ),
    label.Label(
        label.WithFor("option1"),
        label.WithText("Option 1"),
    ),
)
```

## With Form Field

Use with `form.Field` for complete form integration:

```go
field := form.Field(
    "Choose Size",
    radio.Group(
        []radio.GroupOption{radio.WithGroupName("size")},
        radio.GroupItem("size-small", "Small", radio.WithValue("small"), radio.WithName("size")),
        radio.GroupItem("size-medium", "Medium", radio.WithValue("medium"), radio.WithName("size")),
        radio.GroupItem("size-large", "Large", radio.WithValue("large"), radio.WithName("size")),
    ),
    form.WithID("size"),
    form.WithRequired(),
    form.WithDescription("Select your preferred size"),
)
```

## Styling

The radio component uses Tailwind CSS classes:

```go
// Base classes
"aspect-square size-4 rounded-full border border-input"
"text-primary shadow-xs transition-shadow"
"focus-visible:border-ring focus-visible:ring-ring/50"
"focus-visible:ring-[3px] outline-none"
"disabled:cursor-not-allowed disabled:opacity-50"
"data-[state=checked]:bg-primary"
"data-[state=checked]:text-primary-foreground"
"data-[state=checked]:border-primary"
```

### Custom Styling

```go
customRadio := radio.Radio(
    radio.WithID("custom"),
    radio.WithClass("size-6 border-2"),
)
```

## GroupItem Component

The `GroupItem` component creates a radio button with label in a flex layout:

```go
item := radio.GroupItem(
    "size-small",  // ID
    "Small",       // Label
    radio.WithValue("small"),
    radio.WithName("size"),
)
```

## Accessibility

- Uses semantic HTML `<input type="radio">` element
- Proper label association via `id` attribute
- Keyboard accessible (Arrow keys to navigate, Space to select)
- Focus visible with ring outline
- Disabled state prevents interaction
- Group has `role="radiogroup"` for screen readers
- Supports `aria-*` attributes via `WithAttrs`

## States

### Default

```go
radio.Radio(
    radio.WithID("default"),
    radio.WithName("state"),
)
```

### Checked

```go
radio.Radio(
    radio.WithID("checked"),
    radio.WithName("state"),
    radio.Checked(),
)
```

### Disabled

```go
radio.Radio(
    radio.WithID("disabled"),
    radio.WithName("state"),
    radio.Disabled(),
)
```

### Disabled + Checked

```go
radio.Radio(
    radio.WithID("disabled-checked"),
    radio.WithName("state"),
    radio.Checked(),
    radio.Disabled(),
)
```

## Complete Example

```go
func ShippingForm() g.Node {
    return form.Form(
        []form.Option{
            form.WithAction("/checkout"),
            form.WithMethod("POST"),
        },
        
        // Shipping method
        form.Field(
            "Shipping Method",
            radio.Group(
                []radio.GroupOption{
                    radio.WithGroupName("shipping"),
                    radio.WithGroupValue("standard"),
                },
                radio.GroupItem(
                    "shipping-standard",
                    "Standard Shipping (5-7 days) - Free",
                    radio.WithValue("standard"),
                    radio.WithName("shipping"),
                    radio.Checked(),
                ),
                radio.GroupItem(
                    "shipping-express",
                    "Express Shipping (2-3 days) - $10",
                    radio.WithValue("express"),
                    radio.WithName("shipping"),
                ),
                radio.GroupItem(
                    "shipping-overnight",
                    "Overnight Shipping (1 day) - $25",
                    radio.WithValue("overnight"),
                    radio.WithName("shipping"),
                ),
            ),
            form.WithID("shipping"),
            form.WithRequired(),
            form.WithDescription("Choose your preferred shipping method"),
        ),
        
        // Payment method
        form.Field(
            "Payment Method",
            radio.Group(
                []radio.GroupOption{
                    radio.WithGroupName("payment"),
                },
                radio.GroupItem(
                    "payment-card",
                    "Credit Card",
                    radio.WithValue("card"),
                    radio.WithName("payment"),
                ),
                radio.GroupItem(
                    "payment-paypal",
                    "PayPal",
                    radio.WithValue("paypal"),
                    radio.WithName("payment"),
                ),
                radio.GroupItem(
                    "payment-bank",
                    "Bank Transfer",
                    radio.WithValue("bank"),
                    radio.WithName("payment"),
                ),
            ),
            form.WithID("payment"),
            form.WithRequired(),
        ),
        
        button.Button(
            g.Text("Continue to Payment"),
            button.WithType("submit"),
        ),
    )
}
```

## Legacy RadioGroup

The package also includes a legacy `RadioGroup` function that uses primitives:

```go
group := radio.RadioGroup("size", []radio.RadioGroupOption{
    {ID: "size-small", Value: "small", Label: "Small"},
    {ID: "size-medium", Value: "medium", Label: "Medium", Checked: true},
    {ID: "size-large", Value: "large", Label: "Large"},
})
```

**Note**: Prefer using the new `Group` and `GroupItem` components for better consistency.

## Dark Mode

The radio component automatically supports dark mode:

- Dark background: `dark:bg-input/30`
- Dark checked state: `dark:data-[state=checked]:bg-primary`
- Dark border colors adjust automatically

## Related Components

- [Form](../form/README.md) - Form wrapper and validation
- [Label](../label/README.md) - Label component
- [Checkbox](../checkbox/README.md) - Checkbox (multiple selection)
- [Select](../select/README.md) - Dropdown selection

