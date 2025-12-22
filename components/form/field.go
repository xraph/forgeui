package form

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/components/label"
)

// FieldProps defines form field configuration
type FieldProps struct {
	ID          string
	Name        string
	Required    bool
	Disabled    bool
	Error       string
	Description string
	Class       string
}

// FieldOption is a functional option for configuring form fields
type FieldOption func(*FieldProps)

// WithID sets the field ID
func WithID(id string) FieldOption {
	return func(p *FieldProps) { p.ID = id }
}

// WithName sets the field name
func WithName(name string) FieldOption {
	return func(p *FieldProps) { p.Name = name }
}

// WithRequired marks the field as required
func WithRequired() FieldOption {
	return func(p *FieldProps) { p.Required = true }
}

// WithDisabled marks the field as disabled
func WithDisabled() FieldOption {
	return func(p *FieldProps) { p.Disabled = true }
}

// WithError sets an error message
func WithError(err string) FieldOption {
	return func(p *FieldProps) { p.Error = err }
}

// WithDescription sets helper text
func WithDescription(desc string) FieldOption {
	return func(p *FieldProps) { p.Description = desc }
}

// WithFieldClass adds custom classes
func WithFieldClass(class string) FieldOption {
	return func(p *FieldProps) { p.Class = class }
}

// Field creates a complete form field with label, control, error, and description
//
// Example:
//
//	field := form.Field(
//	    "Email",
//	    input.Input(
//	        input.WithType("email"),
//	        input.WithPlaceholder("you@example.com"),
//	    ),
//	    form.WithID("email"),
//	    form.WithRequired(),
//	    form.WithDescription("We'll never share your email."),
//	)
func Field(labelText string, control g.Node, opts ...FieldOption) g.Node {
	props := &FieldProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN("space-y-2", props.Class)
	
	children := []g.Node{
		FieldLabel(labelText, props.ID, props.Required),
		control,
	}

	if props.Description != "" {
		children = append(children, FieldDescription(props.Description))
	}

	if props.Error != "" {
		children = append(children, FieldError(props.Error))
	}

	return html.Div(
		html.Class(classes),
		g.Group(children),
	)
}

// FieldLabel creates a form field label
func FieldLabel(text, htmlFor string, required bool) g.Node {
	labelText := text
	if required {
		labelText = text + " *"
	}

	opts := []label.Option{}
	if htmlFor != "" {
		opts = append(opts, label.WithFor(htmlFor))
	}

	return label.Label(labelText, opts...)
}

// FieldControl is a simple wrapper for form controls (for API consistency)
func FieldControl(control g.Node) g.Node {
	return control
}

// FieldError creates a form field error message
func FieldError(text string) g.Node {
	if text == "" {
		return nil
	}

	return html.P(
		html.Class("text-sm font-medium text-destructive"),
		html.Role("alert"),
		g.Attr("aria-live", "polite"),
		g.Text(text),
	)
}

// FieldDescription creates a form field description/helper text
func FieldDescription(text string) g.Node {
	if text == "" {
		return nil
	}

	return html.P(
		html.Class("text-sm text-muted-foreground"),
		g.Text(text),
	)
}

