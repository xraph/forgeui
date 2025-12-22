# Form Components

Complete form system with validation, error handling, and accessibility features.

## Components

- **Form** - Form wrapper with action and method
- **Field** - Complete form field with label, control, error, and description
- **Error** - Standalone error message component
- **Description** - Standalone helper text component
- **Validators** - Validation helper functions

## Installation

```go
import (
    "github.com/xraph/forgeui/components/form"
    "github.com/xraph/forgeui/components/input"
)
```

## Basic Usage

### Simple Form

```go
myForm := form.Form(
    []form.Option{
        form.WithAction("/submit"),
        form.WithMethod("POST"),
    },
    form.Field(
        "Email",
        input.Input(
            input.WithType("email"),
            input.WithName("email"),
        ),
    ),
    button.Button(
        g.Text("Submit"),
        button.WithType("submit"),
    ),
)
```

### Form Field with All Features

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
    form.WithName("email"),
    form.WithRequired(),
    form.WithDescription("We'll never share your email with anyone else."),
    form.WithError("Please enter a valid email address"),
)
```

## Form Component

### Props

- `Action` - Form action URL
- `Method` - HTTP method (default: "POST")
- `Class` - Custom CSS classes
- `Attrs` - Additional HTML attributes

### Options

```go
form.WithAction(action string)
form.WithMethod(method string)
form.WithClass(class string)
form.WithAttrs(attrs ...g.Node)
```

### Example

```go
loginForm := form.Form(
    []form.Option{
        form.WithAction("/auth/login"),
        form.WithMethod("POST"),
        form.WithClass("max-w-md mx-auto"),
    },
    // Form fields...
)
```

## Field Component

### Props

- `ID` - Field ID (for label association)
- `Name` - Field name
- `Required` - Mark as required (adds asterisk)
- `Disabled` - Mark as disabled
- `Error` - Error message to display
- `Description` - Helper text to display
- `Class` - Custom CSS classes

### Options

```go
form.WithID(id string)
form.WithName(name string)
form.WithRequired()
form.WithDisabled()
form.WithError(err string)
form.WithDescription(desc string)
form.WithFieldClass(class string)
```

### Example

```go
emailField := form.Field(
    "Email Address",
    input.Input(
        input.WithType("email"),
        input.WithID("email"),
        input.WithName("email"),
    ),
    form.WithID("email"),
    form.WithRequired(),
    form.WithDescription("Enter your work email"),
)
```

## Error Component

Standalone error message component with ARIA attributes.

### Options

```go
form.WithErrorClass(class string)
form.WithErrorAttrs(attrs ...g.Node)
```

### Example

```go
errorMsg := form.Error(
    "Invalid email address",
    form.WithErrorClass("mt-2"),
)
```

## Description Component

Standalone helper text component.

### Options

```go
form.WithDescriptionClass(class string)
form.WithDescriptionAttrs(attrs ...g.Node)
```

### Example

```go
helpText := form.Description(
    "Password must be at least 8 characters",
    form.WithDescriptionClass("mt-2"),
)
```

## Validation

### Built-in Validators

```go
form.Required()                          // Field is required
form.RequiredWithMessage(msg string)     // Required with custom message
form.Email()                             // Valid email address
form.MinLength(min int)                  // Minimum length
form.MaxLength(max int)                  // Maximum length
form.Pattern(regex, msg string)          // Regex pattern match
form.URL()                               // Valid URL
form.Numeric()                           // Only numbers
form.Alpha()                             // Only letters
form.AlphaNumeric()                      // Letters and numbers
form.In(allowed []string)                // Value in list
```

### Combining Validators

```go
validator := form.Combine(
    form.Required(),
    form.Email(),
    form.MaxLength(100),
)

if err := validator("user@example.com"); err != nil {
    // Handle validation error
}
```

### Validating Fields

```go
// Validate with field name for better errors
err := form.ValidateField(
    "email",
    "user@example.com",
    form.Required(),
    form.Email(),
)

if err != nil {
    // Error includes field name: "email: Invalid email address"
}
```

### Server-Side Validation Example

```go
func handleSubmit(ctx *router.PageContext) g.Node {
    email := ctx.FormValue("email")
    password := ctx.FormValue("password")

    // Validate email
    if err := form.ValidateField("email", email, form.Required(), form.Email()); err != nil {
        return renderForm(err.Error(), "")
    }

    // Validate password
    if err := form.ValidateField("password", password, form.Required(), form.MinLength(8)); err != nil {
        return renderForm("", err.Error())
    }

    // Process form...
}

func renderForm(emailErr, passwordErr string) g.Node {
    return form.Form(
        []form.Option{form.WithAction("/submit")},
        form.Field(
            "Email",
            input.Input(input.WithType("email"), input.WithName("email")),
            form.WithError(emailErr),
        ),
        form.Field(
            "Password",
            input.Input(input.WithType("password"), input.WithName("password")),
            form.WithError(passwordErr),
        ),
    )
}
```

## Complete Form Example

```go
func RegistrationForm() g.Node {
    return form.Form(
        []form.Option{
            form.WithAction("/auth/register"),
            form.WithMethod("POST"),
            form.WithClass("space-y-6 max-w-md mx-auto"),
        },
        // Name field
        form.Field(
            "Full Name",
            input.Input(
                input.WithType("text"),
                input.WithID("name"),
                input.WithName("name"),
                input.WithPlaceholder("John Doe"),
            ),
            form.WithID("name"),
            form.WithRequired(),
        ),
        
        // Email field
        form.Field(
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
        ),
        
        // Password field
        form.Field(
            "Password",
            input.Input(
                input.WithType("password"),
                input.WithID("password"),
                input.WithName("password"),
            ),
            form.WithID("password"),
            form.WithRequired(),
            form.WithDescription("Must be at least 8 characters."),
        ),
        
        // Submit button
        button.Button(
            g.Text("Create Account"),
            button.WithType("submit"),
            button.WithVariant(forgeui.VariantDefault),
            button.WithClass("w-full"),
        ),
    )
}
```

## Accessibility

The form components include proper ARIA attributes:

- **Field**: Proper label association with `for` attribute
- **Error**: `role="alert"` and `aria-live="polite"` for screen readers
- **Required**: Visual indicator (asterisk) and semantic HTML
- **Description**: Associated with input via `aria-describedby`

## Styling

All components use Tailwind CSS classes and support dark mode:

- **Error**: `text-destructive` color
- **Description**: `text-muted-foreground` color
- **Field**: `space-y-2` spacing between elements
- **Form**: `space-y-6` spacing between fields

## Custom Validators

Create custom validators for your specific needs:

```go
func CustomValidator(message string) form.Validator {
    return func(value string) error {
        if !myCustomLogic(value) {
            return form.NewValidationError("", message)
        }
        return nil
    }
}

// Usage
validator := form.Combine(
    form.Required(),
    CustomValidator("Custom validation failed"),
)
```

## Related Components

- [Input](../input/README.md) - Text input component
- [Checkbox](../checkbox/README.md) - Checkbox component
- [Radio](../radio/README.md) - Radio button component
- [Select](../select/README.md) - Select dropdown component
- [Textarea](../textarea/README.md) - Textarea component

