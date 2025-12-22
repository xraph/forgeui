package input

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/components/button"
)

// Alignment defines the position of addons in an input group
type Alignment string

const (
	AlignInlineStart Alignment = "inline-start"
	AlignInlineEnd   Alignment = "inline-end"
	AlignBlockStart  Alignment = "block-start"
	AlignBlockEnd    Alignment = "block-end"
)

// GroupProps defines input group configuration
type GroupProps struct {
	Disabled bool
	Class    string
	Attrs    []g.Node
}

// GroupOption is a functional option for configuring input groups
type GroupOption func(*GroupProps)

// WithGroupClass adds custom classes to the input group
func WithGroupClass(class string) GroupOption {
	return func(p *GroupProps) { p.Class = class }
}

// WithGroupAttrs adds custom attributes to the input group
func WithGroupAttrs(attrs ...g.Node) GroupOption {
	return func(p *GroupProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// GroupDisabled sets the disabled state on the input group
func GroupDisabled() GroupOption {
	return func(p *GroupProps) { p.Disabled = true }
}

// InputGroup creates a group of inputs with addons or positioned elements.
//
// This follows the shadcn/ui input-group pattern with support for:
// - Inline addons (inline-start, inline-end)
// - Block addons (block-start, block-end)
// - Focus and error states
// - Disabled state
//
// Example with inline addons:
//
//	input.InputGroup(
//	    []input.GroupOption{},
//	    input.InputGroupAddon(
//	        []input.AddonOption{input.WithAddonAlign(input.AlignInlineStart)},
//	        icons.Mail(),
//	    ),
//	    input.InputGroupInput(input.WithPlaceholder("Enter email")),
//	)
//
// Example with button:
//
//	input.InputGroup(
//	    []input.GroupOption{},
//	    input.InputGroupInput(input.WithPlaceholder("Search...")),
//	    input.InputGroupAddon(
//	        []input.AddonOption{input.WithAddonAlign(input.AlignInlineEnd)},
//	        input.InputGroupButton(
//	            g.Text("Search"),
//	            input.WithGroupButtonSize(input.GroupButtonSizeSM),
//	        ),
//	    ),
//	)
func InputGroup(opts []GroupOption, children ...g.Node) g.Node {
	props := &GroupProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		// Base styles
		"group/input-group",
		"border-input",
		"dark:bg-input/30",
		"relative",
		"flex",
		"w-full",
		"items-center",
		"rounded-md",
		"border",
		"shadow-xs",
		"transition-[color,box-shadow]",
		"outline-none",

		// Default height with textarea override
		"h-9",
		"min-w-0",
		"has-[>textarea]:h-auto",

		// Style direct child input elements (allows using regular Input inside InputGroup)
		// Using ! for important to override CVA styles from Input component
		"[&>input]:flex-1",
		"[&>input]:h-full",
		"[&>input]:w-full",
		"[&>input]:min-w-0",
		"[&>input]:!rounded-none",
		"[&>input]:!border-0",
		"[&>input]:!bg-transparent",
		"[&>input]:px-3",
		"[&>input]:py-2",
		"[&>input]:!shadow-none",
		"[&>input]:!outline-none",
		"[&>input]:!ring-0",
		"[&>input]:!ring-offset-0",
		"[&>input]:focus:!ring-0",
		"[&>input]:focus:!ring-offset-0",
		"[&>input]:focus:!outline-none",
		"[&>input]:focus:!border-0",
		"[&>input]:focus:!shadow-none",
		"[&>input]:focus-visible:!ring-0",
		"[&>input]:focus-visible:!ring-offset-0",
		"[&>input]:focus-visible:!outline-none",
		"[&>input]:focus-visible:!border-0",
		"[&>input]:focus-visible:!shadow-none",
		"dark:[&>input]:!bg-transparent",

		// Style direct child textarea elements
		// Using ! for important to override any textarea styles
		"[&>textarea]:flex-1",
		"[&>textarea]:w-full",
		"[&>textarea]:min-w-0",
		"[&>textarea]:resize-none",
		"[&>textarea]:!rounded-none",
		"[&>textarea]:!border-0",
		"[&>textarea]:!bg-transparent",
		"[&>textarea]:px-3",
		"[&>textarea]:py-3",
		"[&>textarea]:!shadow-none",
		"[&>textarea]:!outline-none",
		"[&>textarea]:!ring-0",
		"[&>textarea]:!ring-offset-0",
		"[&>textarea]:focus:!ring-0",
		"[&>textarea]:focus:!ring-offset-0",
		"[&>textarea]:focus:!outline-none",
		"[&>textarea]:focus:!border-0",
		"[&>textarea]:focus:!shadow-none",
		"[&>textarea]:focus-visible:!ring-0",
		"[&>textarea]:focus-visible:!ring-offset-0",
		"[&>textarea]:focus-visible:!outline-none",
		"[&>textarea]:focus-visible:!border-0",
		"[&>textarea]:focus-visible:!shadow-none",
		"dark:[&>textarea]:!bg-transparent",

		// Variants based on alignment
		"has-[>[data-align=inline-start]]:[&>input]:pl-2",
		"has-[>[data-align=inline-end]]:[&>input]:pr-2",
		"has-[>[data-align=block-start]]:h-auto",
		"has-[>[data-align=block-start]]:flex-col",
		"has-[>[data-align=block-start]]:[&>input]:pb-3",
		"has-[>[data-align=block-end]]:h-auto",
		"has-[>[data-align=block-end]]:flex-col",
		"has-[>[data-align=block-end]]:[&>input]:pt-3",

		// Focus state - works with both data-slot and regular inputs
		// Using [3px] for ring width to match shadcn exactly
		"has-[[data-slot=input-group-control]:focus-visible]:border-ring",
		"has-[[data-slot=input-group-control]:focus-visible]:ring-[3px]",
		"has-[[data-slot=input-group-control]:focus-visible]:ring-ring/50",
		"has-[>input:focus-visible]:border-ring",
		"has-[>input:focus-visible]:ring-[3px]",
		"has-[>input:focus-visible]:ring-ring/50",
		"has-[>textarea:focus-visible]:border-ring",
		"has-[>textarea:focus-visible]:ring-[3px]",
		"has-[>textarea:focus-visible]:ring-ring/50",

		// Error state
		"has-[[data-slot][aria-invalid=true]]:ring-destructive/20",
		"has-[[data-slot][aria-invalid=true]]:border-destructive",
		"has-[>input[aria-invalid=true]]:ring-destructive/20",
		"has-[>input[aria-invalid=true]]:border-destructive",
		"has-[>textarea[aria-invalid=true]]:ring-destructive/20",
		"has-[>textarea[aria-invalid=true]]:border-destructive",
		"dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40",
		"dark:has-[>input[aria-invalid=true]]:ring-destructive/40",
		"dark:has-[>textarea[aria-invalid=true]]:ring-destructive/40",

		props.Class,
	)

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("data-slot", "input-group"),
		g.Attr("role", "group"),
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("data-disabled", "true"))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// inputGroupAddonCVA defines the class variance authority for addon alignment
var inputGroupAddonCVA = forgeui.NewCVA(
	// Base classes
	"text-muted-foreground",
	"flex",
	"h-auto",
	"cursor-text",
	"items-center",
	"justify-center",
	"gap-2",
	"py-1.5",
	"text-sm",
	"font-medium",
	"select-none",
	"[&>svg:not([class*='size-'])]:size-4",
	"[&>kbd]:rounded-[calc(var(--radius)-5px)]",
	"group-data-[disabled=true]/input-group:opacity-50",
).Variant("align", map[string][]string{
	"inline-start": {
		"order-first",
		"pl-3",
		"has-[>button]:ml-[-0.45rem]",
		"has-[>kbd]:ml-[-0.35rem]",
	},
	"inline-end": {
		"order-last",
		"pr-3",
		"has-[>button]:mr-[-0.45rem]",
		"has-[>kbd]:mr-[-0.35rem]",
	},
	"block-start": {
		"order-first",
		"w-full",
		"justify-start",
		"px-3",
		"pt-3",
		"[.border-b]:pb-3",
		"group-has-[>input]/input-group:pt-2.5",
	},
	"block-end": {
		"order-last",
		"w-full",
		"justify-start",
		"px-3",
		"pb-3",
		"[.border-t]:pt-3",
		"group-has-[>input]/input-group:pb-2.5",
	},
}).Default("align", "inline-start")

// AddonProps defines addon configuration
type AddonProps struct {
	Align Alignment
	Class string
	Attrs []g.Node
}

// AddonOption is a functional option for configuring addons
type AddonOption func(*AddonProps)

// WithAddonClass adds custom classes to the addon
func WithAddonClass(class string) AddonOption {
	return func(p *AddonProps) { p.Class = class }
}

// WithAddonAttrs adds custom attributes to the addon
func WithAddonAttrs(attrs ...g.Node) AddonOption {
	return func(p *AddonProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// WithAddonAlign sets the alignment of the addon
func WithAddonAlign(align Alignment) AddonOption {
	return func(p *AddonProps) { p.Align = align }
}

// InputGroupAddon creates an addon container for icons, text, or buttons.
//
// Supports four alignment positions:
// - inline-start: Left side of the input (default)
// - inline-end: Right side of the input
// - block-start: Above the input (stacked layout)
// - block-end: Below the input (stacked layout)
//
// Example:
//
//	input.InputGroupAddon(
//	    []input.AddonOption{input.WithAddonAlign(input.AlignInlineStart)},
//	    icons.Search(),
//	)
func InputGroupAddon(opts []AddonOption, children ...g.Node) g.Node {
	props := &AddonProps{
		Align: AlignInlineStart,
	}
	for _, opt := range opts {
		opt(props)
	}

	classes := inputGroupAddonCVA.Classes(map[string]string{
		"align": string(props.Align),
	})

	if props.Class != "" {
		classes = forgeui.CN(classes, props.Class)
	}

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("role", "group"),
		g.Attr("data-slot", "input-group-addon"),
		g.Attr("data-align", string(props.Align)),
	}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// GroupButtonSize defines the size of buttons in input groups
type GroupButtonSize string

const (
	GroupButtonSizeXS     GroupButtonSize = "xs"
	GroupButtonSizeSM     GroupButtonSize = "sm"
	GroupButtonSizeIconXS GroupButtonSize = "icon-xs"
	GroupButtonSizeIconSM GroupButtonSize = "icon-sm"
)

// inputGroupButtonCVA defines the class variance authority for group buttons
var inputGroupButtonCVA = forgeui.NewCVA(
	// Base classes
	"text-sm",
	"shadow-none",
	"flex",
	"gap-2",
	"items-center",
).Variant("size", map[string][]string{
	"xs": {
		"h-6",
		"gap-1",
		"px-2",
		"rounded-[calc(var(--radius)-5px)]",
		"[&>svg:not([class*='size-'])]:size-3.5",
		"has-[>svg]:px-2",
	},
	"sm": {
		"h-8",
		"px-2.5",
		"gap-1.5",
		"rounded-md",
		"has-[>svg]:px-2.5",
	},
	"icon-xs": {
		"size-6",
		"rounded-[calc(var(--radius)-5px)]",
		"p-0",
		"has-[>svg]:p-0",
	},
	"icon-sm": {
		"size-8",
		"p-0",
		"has-[>svg]:p-0",
	},
}).Default("size", "xs")

// GroupButtonProps defines input group button configuration
type GroupButtonProps struct {
	Size    GroupButtonSize
	Variant forgeui.Variant
	Type    string
	Class   string
	Attrs   []g.Node
}

// GroupButtonOption is a functional option for configuring group buttons
type GroupButtonOption func(*GroupButtonProps)

// WithGroupButtonSize sets the button size
func WithGroupButtonSize(size GroupButtonSize) GroupButtonOption {
	return func(p *GroupButtonProps) { p.Size = size }
}

// WithGroupButtonVariant sets the button variant
func WithGroupButtonVariant(v forgeui.Variant) GroupButtonOption {
	return func(p *GroupButtonProps) { p.Variant = v }
}

// WithGroupButtonType sets the button type
func WithGroupButtonType(t string) GroupButtonOption {
	return func(p *GroupButtonProps) { p.Type = t }
}

// WithGroupButtonClass adds custom classes to the button
func WithGroupButtonClass(class string) GroupButtonOption {
	return func(p *GroupButtonProps) { p.Class = class }
}

// WithGroupButtonAttrs adds custom attributes to the button
func WithGroupButtonAttrs(attrs ...g.Node) GroupButtonOption {
	return func(p *GroupButtonProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// InputGroupButton creates a button styled for use inside input groups.
//
// Uses ghost variant by default with smaller sizing options appropriate for input groups.
//
// Example:
//
//	input.InputGroupButton(
//	    g.Text("Search"),
//	    input.WithGroupButtonSize(input.GroupButtonSizeXS),
//	)
func InputGroupButton(children g.Node, opts ...GroupButtonOption) g.Node {
	props := &GroupButtonProps{
		Size:    GroupButtonSizeXS,
		Variant: forgeui.VariantGhost,
		Type:    "button",
	}
	for _, opt := range opts {
		opt(props)
	}

	sizeClasses := inputGroupButtonCVA.Classes(map[string]string{
		"size": string(props.Size),
	})

	classes := forgeui.CN(sizeClasses, props.Class)

	buttonOpts := []button.Option{
		button.WithVariant(props.Variant),
		button.WithType(props.Type),
		button.WithClass(classes),
		button.WithAttrs(g.Attr("data-size", string(props.Size))),
	}
	buttonOpts = append(buttonOpts, button.WithAttrs(props.Attrs...))

	return button.Button(children, buttonOpts...)
}

// InputGroupText creates a text element for use inside addons.
//
// Example:
//
//	input.InputGroupAddon(
//	    []input.AddonOption{input.WithAddonAlign(input.AlignInlineStart)},
//	    input.InputGroupText("https://"),
//	)
func InputGroupText(text string, opts ...func(*textProps)) g.Node {
	props := &textProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		"text-muted-foreground",
		"flex",
		"items-center",
		"gap-2",
		"text-sm",
		"[&_svg]:pointer-events-none",
		"[&_svg:not([class*='size-'])]:size-4",
		props.class,
	)

	return html.Span(
		html.Class(classes),
		g.Text(text),
	)
}

type textProps struct {
	class string
}

// WithTextClass adds custom classes to the text element
func WithTextClass(class string) func(*textProps) {
	return func(p *textProps) { p.class = class }
}

// GroupInputProps defines input group input configuration
type GroupInputProps struct {
	Type        string
	Name        string
	ID          string
	Placeholder string
	Value       string
	Required    bool
	Disabled    bool
	Invalid     bool
	Class       string
	Attrs       []g.Node
}

// GroupInputOption is a functional option for configuring group inputs
type GroupInputOption func(*GroupInputProps)

// WithGroupInputType sets the input type
func WithGroupInputType(t string) GroupInputOption {
	return func(p *GroupInputProps) { p.Type = t }
}

// WithGroupInputName sets the input name
func WithGroupInputName(name string) GroupInputOption {
	return func(p *GroupInputProps) { p.Name = name }
}

// WithGroupInputID sets the input ID
func WithGroupInputID(id string) GroupInputOption {
	return func(p *GroupInputProps) { p.ID = id }
}

// WithGroupInputPlaceholder sets the input placeholder
func WithGroupInputPlaceholder(placeholder string) GroupInputOption {
	return func(p *GroupInputProps) { p.Placeholder = placeholder }
}

// WithGroupInputValue sets the input value
func WithGroupInputValue(value string) GroupInputOption {
	return func(p *GroupInputProps) { p.Value = value }
}

// GroupInputRequired sets the required attribute
func GroupInputRequired() GroupInputOption {
	return func(p *GroupInputProps) { p.Required = true }
}

// GroupInputDisabled sets the disabled attribute
func GroupInputDisabled() GroupInputOption {
	return func(p *GroupInputProps) { p.Disabled = true }
}

// GroupInputInvalid sets the aria-invalid attribute for error states
func GroupInputInvalid() GroupInputOption {
	return func(p *GroupInputProps) { p.Invalid = true }
}

// WithGroupInputClass adds custom classes to the input
func WithGroupInputClass(class string) GroupInputOption {
	return func(p *GroupInputProps) { p.Class = class }
}

// WithGroupInputAttrs adds custom attributes to the input
func WithGroupInputAttrs(attrs ...g.Node) GroupInputOption {
	return func(p *GroupInputProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// InputGroupInput creates an input styled for use inside input groups.
//
// This input removes its own border and shadow since the parent InputGroup handles those.
// Note: You can also use the regular Input() component inside InputGroup - it will be styled automatically.
//
// Example:
//
//	input.InputGroupInput(
//	    input.WithGroupInputPlaceholder("Enter value..."),
//	    input.WithGroupInputName("search"),
//	)
func InputGroupInput(opts ...GroupInputOption) g.Node {
	props := &GroupInputProps{
		Type: "text",
	}
	for _, opt := range opts {
		opt(props)
	}

	// Most styling is handled by the parent InputGroup via [&>input] selectors
	// We only add semantic classes here
	classes := forgeui.CN(
		"file:text-foreground",
		"placeholder:text-muted-foreground",
		"selection:bg-primary",
		"selection:text-primary-foreground",
		"text-sm",
		"file:inline-flex",
		"file:h-7",
		"file:border-0",
		"file:bg-transparent",
		"file:text-sm",
		"file:font-medium",
		"disabled:cursor-not-allowed",
		"disabled:opacity-50",
		props.Class,
	)

	attrs := []g.Node{
		html.Class(classes),
		html.Type(props.Type),
		g.Attr("data-slot", "input-group-control"),
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Placeholder != "" {
		attrs = append(attrs, html.Placeholder(props.Placeholder))
	}

	if props.Value != "" {
		attrs = append(attrs, html.Value(props.Value))
	}

	if props.Required {
		attrs = append(attrs, g.Attr("required", ""))
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("disabled", ""))
	}

	if props.Invalid {
		attrs = append(attrs, g.Attr("aria-invalid", "true"))
	}

	attrs = append(attrs, props.Attrs...)

	return html.Input(g.Group(attrs))
}

// GroupTextareaProps defines input group textarea configuration
type GroupTextareaProps struct {
	Name        string
	ID          string
	Placeholder string
	Value       string
	Rows        int
	Required    bool
	Disabled    bool
	Invalid     bool
	Class       string
	Attrs       []g.Node
}

// GroupTextareaOption is a functional option for configuring group textareas
type GroupTextareaOption func(*GroupTextareaProps)

// WithGroupTextareaName sets the textarea name
func WithGroupTextareaName(name string) GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Name = name }
}

// WithGroupTextareaID sets the textarea ID
func WithGroupTextareaID(id string) GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.ID = id }
}

// WithGroupTextareaPlaceholder sets the textarea placeholder
func WithGroupTextareaPlaceholder(placeholder string) GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Placeholder = placeholder }
}

// WithGroupTextareaValue sets the textarea value
func WithGroupTextareaValue(value string) GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Value = value }
}

// WithGroupTextareaRows sets the textarea rows
func WithGroupTextareaRows(rows int) GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Rows = rows }
}

// GroupTextareaRequired sets the required attribute
func GroupTextareaRequired() GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Required = true }
}

// GroupTextareaDisabled sets the disabled attribute
func GroupTextareaDisabled() GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Disabled = true }
}

// GroupTextareaInvalid sets the aria-invalid attribute for error states
func GroupTextareaInvalid() GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Invalid = true }
}

// WithGroupTextareaClass adds custom classes to the textarea
func WithGroupTextareaClass(class string) GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Class = class }
}

// WithGroupTextareaAttrs adds custom attributes to the textarea
func WithGroupTextareaAttrs(attrs ...g.Node) GroupTextareaOption {
	return func(p *GroupTextareaProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// InputGroupTextarea creates a textarea styled for use inside input groups.
//
// This textarea removes its own border and shadow since the parent InputGroup handles those.
// Note: You can also use the regular Textarea component inside InputGroup - it will be styled automatically.
//
// Example:
//
//	input.InputGroupTextarea(
//	    input.WithGroupTextareaPlaceholder("Enter message..."),
//	    input.WithGroupTextareaRows(4),
//	)
func InputGroupTextarea(opts ...GroupTextareaOption) g.Node {
	props := &GroupTextareaProps{
		Rows: 3,
	}
	for _, opt := range opts {
		opt(props)
	}

	// Most styling is handled by the parent InputGroup via [&>textarea] selectors
	// We only add semantic classes here
	classes := forgeui.CN(
		"placeholder:text-muted-foreground",
		"selection:bg-primary",
		"selection:text-primary-foreground",
		"min-h-[80px]",
		"text-sm",
		"disabled:cursor-not-allowed",
		"disabled:opacity-50",
		props.Class,
	)

	attrs := []g.Node{
		html.Class(classes),
		g.Attr("data-slot", "input-group-control"),
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Placeholder != "" {
		attrs = append(attrs, html.Placeholder(props.Placeholder))
	}

	if props.Rows > 0 {
		attrs = append(attrs, html.Rows(string(rune('0'+props.Rows))))
	}

	if props.Required {
		attrs = append(attrs, g.Attr("required", ""))
	}

	if props.Disabled {
		attrs = append(attrs, g.Attr("disabled", ""))
	}

	if props.Invalid {
		attrs = append(attrs, g.Attr("aria-invalid", "true"))
	}

	attrs = append(attrs, props.Attrs...)

	children := []g.Node{}
	if props.Value != "" {
		children = append(children, g.Text(props.Value))
	}

	return html.Textarea(
		g.Group(attrs),
		g.Group(children),
	)
}

// Legacy support - keeping backward compatible functions

// InputLeftAddon creates a left addon (prefix) attached to the input.
// Deprecated: Use InputGroupAddon with WithAddonAlign(AlignInlineStart) instead.
func InputLeftAddon(opts []AddonOption, children ...g.Node) g.Node {
	opts = append([]AddonOption{WithAddonAlign(AlignInlineStart)}, opts...)
	return InputGroupAddon(opts, children...)
}

// InputRightAddon creates a right addon (suffix) attached to the input.
// Deprecated: Use InputGroupAddon with WithAddonAlign(AlignInlineEnd) instead.
func InputRightAddon(opts []AddonOption, children ...g.Node) g.Node {
	opts = append([]AddonOption{WithAddonAlign(AlignInlineEnd)}, opts...)
	return InputGroupAddon(opts, children...)
}

// ElementProps defines positioned element configuration
type ElementProps struct {
	Class string
	Attrs []g.Node
}

// ElementOption is a functional option for configuring positioned elements
type ElementOption func(*ElementProps)

// WithElementClass adds custom classes to the element
func WithElementClass(class string) ElementOption {
	return func(p *ElementProps) { p.Class = class }
}

// WithElementAttrs adds custom attributes to the element
func WithElementAttrs(attrs ...g.Node) ElementOption {
	return func(p *ElementProps) { p.Attrs = append(p.Attrs, attrs...) }
}

// InputLeftElement creates a left element positioned absolutely inside the input.
// Deprecated: Use InputGroupAddon with WithAddonAlign(AlignInlineStart) instead.
func InputLeftElement(opts []ElementOption, children ...g.Node) g.Node {
	props := &ElementProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		"absolute left-0 top-0 bottom-0",
		"flex items-center justify-center",
		"w-10 pointer-events-none",
		props.Class,
	)

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// InputRightElement creates a right element positioned absolutely inside the input.
// Deprecated: Use InputGroupAddon with WithAddonAlign(AlignInlineEnd) instead.
func InputRightElement(opts []ElementOption, children ...g.Node) g.Node {
	props := &ElementProps{}
	for _, opt := range opts {
		opt(props)
	}

	classes := forgeui.CN(
		"absolute right-0 top-0 bottom-0",
		"flex items-center justify-center",
		"w-10",
		props.Class,
	)

	attrs := []g.Node{html.Class(classes)}
	attrs = append(attrs, props.Attrs...)

	return html.Div(
		g.Group(attrs),
		g.Group(children),
	)
}

// SearchInput creates a search input with a search icon
func SearchInput(opts ...Option) g.Node {
	return InputGroup(
		[]GroupOption{},
		InputGroupAddon(
			[]AddonOption{WithAddonAlign(AlignInlineStart)},
			html.SVG(
				g.Attr("xmlns", "http://www.w3.org/2000/svg"),
				g.Attr("width", "16"),
				g.Attr("height", "16"),
				g.Attr("viewBox", "0 0 24 24"),
				g.Attr("fill", "none"),
				g.Attr("stroke", "currentColor"),
				g.Attr("stroke-width", "2"),
				html.Class("text-muted-foreground"),
				g.El("circle", g.Attr("cx", "11"), g.Attr("cy", "11"), g.Attr("r", "8")),
				g.El("path", g.Attr("d", "m21 21-4.3-4.3")),
			),
		),
		InputGroupInput(WithGroupInputPlaceholder("Search...")),
	)
}
