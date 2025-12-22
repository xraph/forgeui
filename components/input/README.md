# Input Component

Text input component with variants, sizes, and validation states.

## Installation

```go
import "github.com/xraph/forgeui/components/input"
```

## Basic Usage

```go
textInput := input.Input(
    input.WithType("text"),
    input.WithPlaceholder("Enter text..."),
)
```

## Props

- `Type` - Input type (text, email, password, etc.)
- `ID` - Input ID
- `Name` - Input name
- `Value` - Input value
- `Placeholder` - Placeholder text
- `Disabled` - Disabled state
- `Required` - Required state
- `ReadOnly` - Read-only state
- `Class` - Custom CSS classes
- `Attrs` - Additional HTML attributes

## Options

```go
input.WithType(t string)
input.WithID(id string)
input.WithName(name string)
input.WithValue(value string)
input.WithPlaceholder(placeholder string)
input.Disabled()
input.Required()
input.ReadOnly()
input.WithClass(class string)
input.WithAttrs(attrs ...g.Node)
```

## Examples

### Text Input

```go
textInput := input.Input(
    input.WithType("text"),
    input.WithName("username"),
    input.WithPlaceholder("Enter username"),
)
```

### Email Input

```go
emailInput := input.Input(
    input.WithType("email"),
    input.WithName("email"),
    input.WithPlaceholder("you@example.com"),
    input.Required(),
)
```

### Password Input

```go
passwordInput := input.Input(
    input.WithType("password"),
    input.WithName("password"),
    input.WithPlaceholder("••••••••"),
    input.Required(),
)
```

### Disabled Input

```go
disabledInput := input.Input(
    input.WithType("text"),
    input.WithValue("Cannot edit"),
    input.Disabled(),
)
```

### Read-Only Input

```go
readOnlyInput := input.Input(
    input.WithType("text"),
    input.WithValue("Read only value"),
    input.ReadOnly(),
)
```

## Input Types

Supports all HTML5 input types:

- `text` - Single-line text
- `email` - Email address
- `password` - Password (hidden text)
- `number` - Numeric input
- `tel` - Telephone number
- `url` - URL
- `search` - Search input
- `date` - Date picker
- `time` - Time picker
- `datetime-local` - Date and time
- `month` - Month picker
- `week` - Week picker
- `color` - Color picker
- `file` - File upload
- `range` - Slider
- `hidden` - Hidden input

## With Form Field

Use with the `form.Field` component for complete form fields:

```go
field := form.Field(
    "Email Address",
    input.Input(
        input.WithType("email"),
        input.WithID("email"),
        input.WithName("email"),
        input.WithPlaceholder("you@example.com"),
    ),
    form.WithID("email"),
    form.WithRequired(),
    form.WithDescription("We'll never share your email."),
)
```

## Validation States

The input component supports validation states through CSS classes:

```go
// Error state (use with form.Field)
field := form.Field(
    "Email",
    input.Input(
        input.WithType("email"),
        input.WithClass("border-destructive"),
    ),
    form.WithError("Invalid email address"),
)
```

## Styling

The input component uses Tailwind CSS classes:

```go
// Base classes
"flex h-10 w-full rounded-md border border-input bg-background px-3 py-2"
"text-sm ring-offset-background"
"file:border-0 file:bg-transparent file:text-sm file:font-medium"
"placeholder:text-muted-foreground"
"focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
"disabled:cursor-not-allowed disabled:opacity-50"
```

### Custom Styling

```go
customInput := input.Input(
    input.WithType("text"),
    input.WithClass("bg-blue-50 border-blue-300"),
)
```

## Accessibility

- Uses semantic HTML `<input>` element
- Supports `aria-*` attributes via `WithAttrs`
- Works with `<label>` elements via `id` attribute
- Proper focus states with visible outline
- Disabled state prevents interaction

## File Input

```go
fileInput := input.Input(
    input.WithType("file"),
    input.WithName("upload"),
    input.WithAttrs(g.Attr("accept", "image/*")),
)
```

## Number Input

```go
numberInput := input.Input(
    input.WithType("number"),
    input.WithName("quantity"),
    input.WithAttrs(
        g.Attr("min", "1"),
        g.Attr("max", "100"),
        g.Attr("step", "1"),
    ),
)
```

## Search Input

```go
searchInput := input.Input(
    input.WithType("search"),
    input.WithName("query"),
    input.WithPlaceholder("Search..."),
)
```

## Date Input

```go
dateInput := input.Input(
    input.WithType("date"),
    input.WithName("birthdate"),
)
```

## Complete Example

```go
func LoginForm() g.Node {
    return form.Form(
        []form.Option{
            form.WithAction("/auth/login"),
            form.WithMethod("POST"),
        },
        form.Field(
            "Email",
            input.Input(
                input.WithType("email"),
                input.WithID("email"),
                input.WithName("email"),
                input.WithPlaceholder("you@example.com"),
                input.Required(),
            ),
            form.WithID("email"),
            form.WithRequired(),
        ),
        form.Field(
            "Password",
            input.Input(
                input.WithType("password"),
                input.WithID("password"),
                input.WithName("password"),
                input.WithPlaceholder("••••••••"),
                input.Required(),
            ),
            form.WithID("password"),
            form.WithRequired(),
        ),
        button.Button(
            g.Text("Sign In"),
            button.WithType("submit"),
        ),
    )
}
```

## Related Components

- [Form](../form/README.md) - Form wrapper and validation
- [Textarea](../textarea/README.md) - Multi-line text input
- [Select](../select/README.md) - Dropdown selection
- [Checkbox](../checkbox/README.md) - Checkbox input
- [Radio](../radio/README.md) - Radio button input

