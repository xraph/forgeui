// Package selectc provides select dropdown components matching shadcn/ui design.
package selectc

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/icons"
)

// SelectSize defines the size variants for the select trigger
type SelectSize string

const (
	SelectSizeDefault SelectSize = "default"
	SelectSizeSm      SelectSize = "sm"
)

// =============================================================================
// Select (Root Container)
// =============================================================================

// SelectProps defines the select root configuration
type SelectProps struct {
	Name         string
	ID           string
	DefaultValue string
	Disabled     bool
	Required     bool
	Class        string
	Attrs        []g.Node
}

// SelectOption is a functional option for Select
type SelectOption func(*SelectProps)

// WithName sets the select name for form submission
func WithName(name string) SelectOption {
	return func(p *SelectProps) { p.Name = name }
}

// WithID sets the select ID
func WithID(id string) SelectOption {
	return func(p *SelectProps) { p.ID = id }
}

// WithDefaultValue sets the default selected value
func WithDefaultValue(value string) SelectOption {
	return func(p *SelectProps) { p.DefaultValue = value }
}

// WithDisabled disables the select
func WithDisabled() SelectOption {
	return func(p *SelectProps) { p.Disabled = true }
}

// WithRequired makes the select required
func WithRequired() SelectOption {
	return func(p *SelectProps) { p.Required = true }
}

// WithClass adds custom classes to the select
func WithClass(class string) SelectOption {
	return func(p *SelectProps) { p.Class = class }
}

// WithAttrs adds custom attributes to the select
func WithAttrs(attrs ...g.Node) SelectOption {
	return func(p *SelectProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// Select creates an interactive select dropdown using Alpine.js
// This provides a custom styled select with full keyboard navigation
//
// Example:
//
//	selectc.Select(
//	    []selectc.SelectOption{selectc.WithName("country")},
//	    selectc.SelectTrigger(
//	        selectc.SelectValue("Select a country"),
//	    ),
//	    selectc.SelectContent(
//	        selectc.SelectItem("us", "United States"),
//	        selectc.SelectItem("uk", "United Kingdom"),
//	    ),
//	)
func Select(opts []SelectOption, children ...g.Node) g.Node {
	props := &SelectProps{}
	for _, opt := range opts {
		opt(props)
	}

	defaultValue := props.DefaultValue
	if defaultValue == "" {
		defaultValue = ""
	}

	// Alpine.js state for select behavior
	xData := fmt.Sprintf(`{
		open: false,
		value: '%s',
		label: '',
		disabled: %t,
		required: %t,
		toggle() {
			if (!this.disabled) {
				this.open = !this.open;
			}
		},
		close() {
			this.open = false;
		},
		selectItem(value, label) {
			this.value = value;
			this.label = label;
			this.close();
			this.$dispatch('change', { value: value, label: label });
		},
		handleKeydown(event) {
			if (event.key === 'Escape') {
				this.close();
			} else if (event.key === 'Enter' || event.key === ' ') {
				if (!this.open) {
					event.preventDefault();
					this.toggle();
				}
			} else if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
				event.preventDefault();
				if (!this.open) {
					this.toggle();
				}
			}
		}
	}`, defaultValue, props.Disabled, props.Required)

	classes := forgeui.CN("relative inline-block", props.Class)

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("data-slot", "select"),
		alpine.XData(nil),
		g.Attr("x-data", xData),
		alpine.XOn("keydown", "handleKeydown($event)"),
		alpine.XOn("click.away", "close()"),
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	attrs = append(attrs, props.Attrs...)

	// Hidden input for form submission
	var hiddenInput g.Node
	if props.Name != "" {
		hiddenAttrs := []g.Node{
			g.Attr("type", "hidden"),
			html.Name(props.Name),
			g.Attr(":value", "value"),
		}
		if props.Required {
			hiddenAttrs = append(hiddenAttrs, g.Attr("required", ""))
		}
		hiddenInput = html.Input(hiddenAttrs...)
	}

	return html.Div(
		g.Group(attrs),
		hiddenInput,
		g.Group(children),
	)
}

// =============================================================================
// SelectTrigger
// =============================================================================

// SelectTriggerProps defines the trigger button configuration
type SelectTriggerProps struct {
	Size  SelectSize
	Class string
	Attrs []g.Node
}

// SelectTriggerOption is a functional option for SelectTrigger
type SelectTriggerOption func(*SelectTriggerProps)

// WithTriggerSize sets the trigger size
func WithTriggerSize(size SelectSize) SelectTriggerOption {
	return func(p *SelectTriggerProps) { p.Size = size }
}

// WithTriggerClass adds custom classes to the trigger
func WithTriggerClass(class string) SelectTriggerOption {
	return func(p *SelectTriggerProps) { p.Class = class }
}

// WithTriggerAttrs adds custom attributes to the trigger
func WithTriggerAttrs(attrs ...g.Node) SelectTriggerOption {
	return func(p *SelectTriggerProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// SelectTrigger creates the button that opens the select dropdown
//
// Example:
//
//	selectc.SelectTrigger(
//	    []selectc.SelectTriggerOption{selectc.WithTriggerSize(selectc.SelectSizeSm)},
//	    selectc.SelectValue("Select an option"),
//	)
func SelectTrigger(opts []SelectTriggerOption, children ...g.Node) g.Node {
	props := &SelectTriggerProps{
		Size: SelectSizeDefault,
	}
	for _, opt := range opts {
		opt(props)
	}

	// Base classes matching shadcn
	baseClasses := forgeui.CN(
		"flex w-full items-center justify-between gap-2",
		"rounded-md border border-input bg-transparent",
		"px-3 py-2 text-sm whitespace-nowrap",
		"shadow-xs transition-[color,box-shadow]",
		"outline-none",
		// Focus states
		"focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50",
		// Disabled states
		"disabled:cursor-not-allowed disabled:opacity-50",
		// Dark mode
		"dark:bg-input/30 dark:hover:bg-input/50",
		// Invalid states
		"aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",
	)

	// Size classes
	var sizeClasses string
	switch props.Size {
	case SelectSizeSm:
		sizeClasses = "h-8"
	default:
		sizeClasses = "h-9"
	}

	classes := forgeui.CN(baseClasses, sizeClasses, props.Class)

	attrs := []g.Node{
		g.Attr("type", "button"),
		html.Class(classes),
		g.Attr("data-slot", "select-trigger"),
		g.Attr("data-size", string(props.Size)),
		g.Attr("role", "combobox"),
		g.Attr("aria-haspopup", "listbox"),
		g.Attr(":aria-expanded", "open"),
		g.Attr(":aria-disabled", "disabled"),
		g.Attr(":disabled", "disabled"),
		alpine.XOn("click", "toggle()"),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Button(
		g.Group(attrs),
		g.Group(children),
		// Chevron icon
		html.Span(
			html.Class("pointer-events-none shrink-0 opacity-50"),
			icons.ChevronDown(icons.WithSize(16)),
		),
	)
}

// =============================================================================
// SelectValue
// =============================================================================

// SelectValue creates the value display inside the trigger
//
// Example:
//
//	selectc.SelectValue("Select a fruit...")
func SelectValue(placeholder string, attrs ...g.Node) g.Node {
	allAttrs := []g.Node{
		html.Class("flex items-center gap-2 line-clamp-1"),
		g.Attr("data-slot", "select-value"),
		g.Attr("x-text", fmt.Sprintf("label || '%s'", placeholder)),
		g.Attr(":data-placeholder", "!label ? 'true' : null"),
		g.Attr(":class", "{ 'text-muted-foreground': !label }"),
	}
	allAttrs = append(allAttrs, attrs...)

	return html.Span(allAttrs...)
}

// =============================================================================
// SelectContent
// =============================================================================

// SelectContentProps defines the dropdown content configuration
type SelectContentProps struct {
	Position string // "item-aligned" or "popper"
	Align    string // "start", "center", "end"
	Class    string
	Attrs    []g.Node
}

// SelectContentOption is a functional option for SelectContent
type SelectContentOption func(*SelectContentProps)

// WithContentPosition sets the content position
func WithContentPosition(position string) SelectContentOption {
	return func(p *SelectContentProps) { p.Position = position }
}

// WithContentAlign sets the content alignment
func WithContentAlign(align string) SelectContentOption {
	return func(p *SelectContentProps) { p.Align = align }
}

// WithContentClass adds custom classes to the content
func WithContentClass(class string) SelectContentOption {
	return func(p *SelectContentProps) { p.Class = class }
}

// WithContentAttrs adds custom attributes to the content
func WithContentAttrs(attrs ...g.Node) SelectContentOption {
	return func(p *SelectContentProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// SelectContent creates the dropdown content container
//
// Example:
//
//	selectc.SelectContent(
//	    nil,
//	    selectc.SelectItem("apple", "Apple"),
//	    selectc.SelectItem("banana", "Banana"),
//	)
func SelectContent(opts []SelectContentOption, children ...g.Node) g.Node {
	props := &SelectContentProps{
		Position: "item-aligned",
		Align:    "center",
	}
	for _, opt := range opts {
		opt(props)
	}

	// Base classes matching shadcn
	baseClasses := forgeui.CN(
		"absolute z-50 mt-1 w-full min-w-[8rem]",
		"max-h-60 overflow-y-auto overflow-x-hidden",
		"rounded-md border bg-popover text-popover-foreground shadow-md",
		// Animation classes (using Alpine.js transitions)
	)

	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("data-slot", "select-content"),
		g.Attr("x-show", "open"),
		g.Attr("x-transition:enter", "transition ease-out duration-100"),
		g.Attr("x-transition:enter-start", "opacity-0 scale-95"),
		g.Attr("x-transition:enter-end", "opacity-100 scale-100"),
		g.Attr("x-transition:leave", "transition ease-in duration-75"),
		g.Attr("x-transition:leave-start", "opacity-100 scale-100"),
		g.Attr("x-transition:leave-end", "opacity-0 scale-95"),
		g.Attr("role", "listbox"),
		g.Attr("tabindex", "-1"),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		html.Div(
			html.Class("p-1"),
			g.Group(children),
		),
	)
}

// =============================================================================
// SelectGroup
// =============================================================================

// SelectGroup creates a group of related select items
//
// Example:
//
//	selectc.SelectGroup(
//	    selectc.SelectLabel("Fruits"),
//	    selectc.SelectItem("apple", "Apple"),
//	    selectc.SelectItem("banana", "Banana"),
//	)
func SelectGroup(children ...g.Node) g.Node {
	return html.Div(
		html.Class(""),
		g.Attr("data-slot", "select-group"),
		g.Attr("role", "group"),
		g.Group(children),
	)
}

// =============================================================================
// SelectLabel
// =============================================================================

// SelectLabel creates a label for a select group
//
// Example:
//
//	selectc.SelectLabel("Fruits")
func SelectLabel(text string, attrs ...g.Node) g.Node {
	allAttrs := []g.Node{
		html.Class("px-2 py-1.5 text-xs text-muted-foreground"),
		g.Attr("data-slot", "select-label"),
	}
	allAttrs = append(allAttrs, attrs...)

	return html.Div(
		g.Group(allAttrs),
		g.Text(text),
	)
}

// =============================================================================
// SelectItem
// =============================================================================

// SelectItemProps defines the item configuration
type SelectItemProps struct {
	Disabled bool
	Class    string
	Attrs    []g.Node
}

// SelectItemOption is a functional option for SelectItem
type SelectItemOption func(*SelectItemProps)

// WithItemDisabled disables the item
func WithItemDisabled() SelectItemOption {
	return func(p *SelectItemProps) { p.Disabled = true }
}

// WithItemClass adds custom classes to the item
func WithItemClass(class string) SelectItemOption {
	return func(p *SelectItemProps) { p.Class = class }
}

// WithItemAttrs adds custom attributes to the item
func WithItemAttrs(attrs ...g.Node) SelectItemOption {
	return func(p *SelectItemProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// SelectItem creates a selectable option within the dropdown
//
// Example:
//
//	selectc.SelectItem("us", "United States", nil)
//	selectc.SelectItem("uk", "United Kingdom", []selectc.SelectItemOption{selectc.WithItemDisabled()})
func SelectItem(value string, label string, opts []SelectItemOption) g.Node {
	props := &SelectItemProps{}
	for _, opt := range opts {
		opt(props)
	}

	// Base classes matching shadcn
	baseClasses := forgeui.CN(
		"relative flex w-full cursor-default select-none items-center gap-2",
		"rounded-sm py-1.5 pl-2 pr-8 text-sm outline-none",
		// Focus/hover states
		"focus:bg-accent focus:text-accent-foreground",
		"hover:bg-accent hover:text-accent-foreground",
		// Disabled states
		"data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
	)

	classes := forgeui.CN(baseClasses, props.Class)

	// Build click handler
	clickHandler := fmt.Sprintf("selectItem('%s', '%s')", value, label)
	if props.Disabled {
		clickHandler = ""
	}

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("data-slot", "select-item"),
		g.Attr("data-value", value),
		g.Attr("role", "option"),
		g.Attr(":aria-selected", fmt.Sprintf("value === '%s'", value)),
		g.Attr("tabindex", "0"),
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("data-disabled", "true"))
		attrs = append(attrs, g.Attr("aria-disabled", "true"))
	} else {
		attrs = append(attrs, alpine.XOn("click", clickHandler))
		attrs = append(attrs, alpine.XOn("keydown.enter.prevent", clickHandler))
		attrs = append(attrs, alpine.XOn("keydown.space.prevent", clickHandler))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		// Check indicator
		html.Span(
			html.Class("absolute right-2 flex h-3.5 w-3.5 items-center justify-center"),
			g.Attr("data-slot", "select-item-indicator"),
			g.Attr("x-show", fmt.Sprintf("value === '%s'", value)),
			icons.Check(icons.WithSize(16)),
		),
		// Item text
		html.Span(
			html.Class("flex-1"),
			g.Text(label),
		),
	)
}

// =============================================================================
// SelectSeparator
// =============================================================================

// SelectSeparator creates a visual separator between items
//
// Example:
//
//	selectc.SelectSeparator()
func SelectSeparator(attrs ...g.Node) g.Node {
	allAttrs := []g.Node{
		html.Class("pointer-events-none -mx-1 my-1 h-px bg-border"),
		g.Attr("data-slot", "select-separator"),
		g.Attr("role", "separator"),
	}
	allAttrs = append(allAttrs, attrs...)

	return html.Div(allAttrs...)
}

// =============================================================================
// Native Select (Simple Version)
// =============================================================================

// NativeSelectOption defines a native select option
type NativeSelectOption struct {
	Value    string
	Label    string
	Selected bool
	Disabled bool
}

// NativeSelect creates a native HTML select dropdown (for simpler use cases)
// This uses the browser's native select with custom styling
//
// Example:
//
//	selectc.NativeSelect(
//	    []selectc.NativeSelectOption{
//	        {Value: "us", Label: "United States"},
//	        {Value: "uk", Label: "United Kingdom", Selected: true},
//	    },
//	    selectc.WithName("country"),
//	)
func NativeSelect(options []NativeSelectOption, opts ...SelectOption) g.Node {
	props := &SelectProps{}
	for _, opt := range opts {
		opt(props)
	}

	// Base classes matching shadcn trigger styles
	baseClasses := forgeui.CN(
		"flex h-9 w-full items-center justify-between gap-2",
		"rounded-md border border-input bg-transparent",
		"px-3 py-2 text-sm whitespace-nowrap",
		"shadow-xs transition-[color,box-shadow]",
		"outline-none",
		// Focus states
		"focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50",
		// Disabled states
		"disabled:cursor-not-allowed disabled:opacity-50",
		// Dark mode
		"dark:bg-input/30 dark:hover:bg-input/50",
		// Appearance
		"appearance-none cursor-pointer",
		"bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2216%22%20height%3D%2216%22%20viewBox%3D%220%200%2024%2024%22%20fill%3D%22none%22%20stroke%3D%22currentColor%22%20stroke-width%3D%222%22%3E%3Cpath%20d%3D%22m6%209%206%206%206-6%22%2F%3E%3C%2Fsvg%3E')]",
		"bg-no-repeat bg-[right_0.5rem_center]",
		"pr-8",
	)

	classes := forgeui.CN(baseClasses, props.Class)

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("data-slot", "native-select"),
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Required {
		attrs = append(attrs, g.Attr("required", ""))
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("disabled", ""))
	}

	attrs = append(attrs, props.Attrs...)

	// Build option elements
	optionNodes := make([]g.Node, 0, len(options))
	for _, opt := range options {
		optAttrs := []g.Node{html.Value(opt.Value)}
		if opt.Selected {
			optAttrs = append(optAttrs, g.Attr("selected", ""))
		}
		if opt.Disabled {
			optAttrs = append(optAttrs, g.Attr("disabled", ""))
		}

		optionNodes = append(optionNodes, html.Option(
			g.Group(optAttrs),
			g.Text(opt.Label),
		))
	}

	return html.Select(
		g.Group(attrs),
		g.Group(optionNodes),
	)
}
