package forgeui

// Size represents component size variants
type Size string

const (
	SizeXS   Size = "xs"
	SizeSM   Size = "sm"
	SizeMD   Size = "md"
	SizeLG   Size = "lg"
	SizeXL   Size = "xl"
	SizeFull Size = "full" // Full width/height for modals
	SizeIcon Size = "icon" // Special size for icon-only buttons
)

// Variant represents visual style variants
type Variant string

const (
	VariantDefault     Variant = "default"
	VariantPrimary     Variant = "primary"
	VariantSecondary   Variant = "secondary"
	VariantDestructive Variant = "destructive"
	VariantOutline     Variant = "outline"
	VariantGhost       Variant = "ghost"
	VariantLink        Variant = "link"
)

// Radius represents border radius options
type Radius string

const (
	RadiusNone Radius = "none"
	RadiusSM   Radius = "sm"
	RadiusMD   Radius = "md"
	RadiusLG   Radius = "lg"
	RadiusXL   Radius = "xl"
	RadiusFull Radius = "full"
)

// Side represents positioning side for drawers, popovers, etc.
type Side string

const (
	SideTop    Side = "top"
	SideRight  Side = "right"
	SideBottom Side = "bottom"
	SideLeft   Side = "left"
)

// Position represents positioning for floating elements (popovers, tooltips, dropdowns)
type Position string

const (
	PositionTop    Position = "top"
	PositionRight  Position = "right"
	PositionBottom Position = "bottom"
	PositionLeft   Position = "left"
)

// Align represents alignment for floating elements
type Align string

const (
	AlignStart  Align = "start"
	AlignCenter Align = "center"
	AlignEnd    Align = "end"
)
