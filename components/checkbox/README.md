# Checkbox Component

Checkbox input component with label support and customization options.

## Installation

```go
import "github.com/xraph/forgeui/components/checkbox"
```

## Basic Usage

```go
check := checkbox.Checkbox(
    checkbox.WithID("terms"),
    checkbox.WithName("terms"),
)
```

## Props

- `ID` - Checkbox ID
- `Name` - Checkbox name
- `Value` - Checkbox value
- `Checked` - Checked state
- `Required` - Required state
- `Disabled` - Disabled state
- `Class` - Custom CSS classes
- `Attrs` - Additional HTML attributes

## Options

```go
checkbox.WithID(id string)
checkbox.WithName(name string)
checkbox.WithValue(value string)
checkbox.Checked()
checkbox.Required()
checkbox.Disabled()
checkbox.WithClass(class string)
checkbox.WithAttrs(attrs ...g.Node)
```

## Examples

### Basic Checkbox

```go
check := checkbox.Checkbox(
    checkbox.WithID("agree"),
    checkbox.WithName("agree"),
)
```

### Checked Checkbox

```go
check := checkbox.Checkbox(
    checkbox.WithID("subscribe"),
    checkbox.WithName("subscribe"),
    checkbox.Checked(),
)
```

### Disabled Checkbox

```go
check := checkbox.Checkbox(
    checkbox.WithID("disabled"),
    checkbox.Disabled(),
)
```

### Required Checkbox

```go
check := checkbox.Checkbox(
    checkbox.WithID("terms"),
    checkbox.WithName("terms"),
    checkbox.Required(),
)
```

## With Label

Use with the `label` component for proper accessibility:

```go
html.Div(
    html.Class("flex items-center space-x-2"),
    checkbox.Checkbox(
        checkbox.WithID("terms"),
        checkbox.WithName("terms"),
    ),
    label.Label(
        label.WithFor("terms"),
        label.WithText("I agree to the terms and conditions"),
    ),
)
```

## With Form Field

Use with `form.Field` for complete form integration:

```go
field := form.Field(
    "Terms and Conditions",
    html.Div(
        html.Class("flex items-center space-x-2"),
        checkbox.Checkbox(
            checkbox.WithID("terms"),
            checkbox.WithName("terms"),
            checkbox.Required(),
        ),
        label.Label(
            label.WithFor("terms"),
            label.WithText("I agree to the terms"),
        ),
    ),
    form.WithID("terms"),
    form.WithRequired(),
    form.WithError("You must agree to the terms"),
)
```

## Styling

The checkbox component uses Tailwind CSS classes:

```go
// Base classes
"aspect-square size-4 rounded border border-input"
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
customCheck := checkbox.Checkbox(
    checkbox.WithID("custom"),
    checkbox.WithClass("size-6 border-2"),
)
```

## Checkbox Group

Create a group of checkboxes for multiple selections:

```go
html.Div(
    html.Class("space-y-2"),
    html.Div(
        html.Class("flex items-center space-x-2"),
        checkbox.Checkbox(
            checkbox.WithID("option1"),
            checkbox.WithName("options"),
            checkbox.WithValue("1"),
        ),
        label.Label(
            label.WithFor("option1"),
            label.WithText("Option 1"),
        ),
    ),
    html.Div(
        html.Class("flex items-center space-x-2"),
        checkbox.Checkbox(
            checkbox.WithID("option2"),
            checkbox.WithName("options"),
            checkbox.WithValue("2"),
        ),
        label.Label(
            label.WithFor("option2"),
            label.WithText("Option 2"),
        ),
    ),
    html.Div(
        html.Class("flex items-center space-x-2"),
        checkbox.Checkbox(
            checkbox.WithID("option3"),
            checkbox.WithName("options"),
            checkbox.WithValue("3"),
        ),
        label.Label(
            label.WithFor("option3"),
            label.WithText("Option 3"),
        ),
    ),
)
```

## Accessibility

- Uses semantic HTML `<input type="checkbox">` element
- Supports proper label association via `id` attribute
- Keyboard accessible (Space to toggle)
- Focus visible with ring outline
- Disabled state prevents interaction
- Supports `aria-*` attributes via `WithAttrs`

## States

### Default

```go
checkbox.Checkbox(checkbox.WithID("default"))
```

### Checked

```go
checkbox.Checkbox(
    checkbox.WithID("checked"),
    checkbox.Checked(),
)
```

### Disabled

```go
checkbox.Checkbox(
    checkbox.WithID("disabled"),
    checkbox.Disabled(),
)
```

### Disabled + Checked

```go
checkbox.Checkbox(
    checkbox.WithID("disabled-checked"),
    checkbox.Checked(),
    checkbox.Disabled(),
)
```

## Complete Example

```go
func SettingsForm() g.Node {
    return form.Form(
        []form.Option{
            form.WithAction("/settings"),
            form.WithMethod("POST"),
        },
        html.Div(
            html.Class("space-y-4"),
            
            // Email notifications
            html.Div(
                html.Class("flex items-center space-x-2"),
                checkbox.Checkbox(
                    checkbox.WithID("email-notifications"),
                    checkbox.WithName("notifications"),
                    checkbox.WithValue("email"),
                    checkbox.Checked(),
                ),
                label.Label(
                    label.WithFor("email-notifications"),
                    label.WithText("Email notifications"),
                ),
            ),
            
            // Push notifications
            html.Div(
                html.Class("flex items-center space-x-2"),
                checkbox.Checkbox(
                    checkbox.WithID("push-notifications"),
                    checkbox.WithName("notifications"),
                    checkbox.WithValue("push"),
                ),
                label.Label(
                    label.WithFor("push-notifications"),
                    label.WithText("Push notifications"),
                ),
            ),
            
            // Marketing emails
            html.Div(
                html.Class("flex items-center space-x-2"),
                checkbox.Checkbox(
                    checkbox.WithID("marketing"),
                    checkbox.WithName("marketing"),
                    checkbox.WithValue("yes"),
                ),
                label.Label(
                    label.WithFor("marketing"),
                    label.WithText("Receive marketing emails"),
                ),
            ),
        ),
        
        button.Button(
            g.Text("Save Settings"),
            button.WithType("submit"),
        ),
    )
}
```

## Dark Mode

The checkbox component automatically supports dark mode:

- Dark background: `dark:bg-input/30`
- Dark checked state: `dark:data-[state=checked]:bg-primary`
- Dark border colors adjust automatically

## Related Components

- [Form](../form/README.md) - Form wrapper and validation
- [Label](../label/README.md) - Label component
- [Radio](../radio/README.md) - Radio button (single selection)
- [Switch](../switch/README.md) - Toggle switch component

