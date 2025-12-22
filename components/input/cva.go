package input

import "github.com/xraph/forgeui"

// inputCVA defines the class variance authority configuration for inputs
// Following shadcn/ui patterns with modern focus states
var inputCVA = forgeui.NewCVA(
	// Base styles matching shadcn exactly
	"file:text-foreground",
	"placeholder:text-muted-foreground",
	"selection:bg-primary",
	"selection:text-primary-foreground",
	"dark:bg-input/30",
	"border-input",
	"h-9",
	"w-full",
	"min-w-0",
	"rounded-md",
	"border",
	"bg-transparent",
	"px-3",
	"py-1",
	"text-base",
	"md:text-sm",
	"shadow-xs",
	"transition-[color,box-shadow]",
	"outline-none",
	"file:inline-flex",
	"file:h-7",
	"file:border-0",
	"file:bg-transparent",
	"file:text-sm",
	"file:font-medium",
	"disabled:pointer-events-none",
	"disabled:cursor-not-allowed",
	"disabled:opacity-50",

	// Focus states
	"focus-visible:border-ring",
	"focus-visible:ring-ring/50",
	"focus-visible:ring-[3px]",

	// Error states via aria-invalid
	"aria-invalid:ring-destructive/20",
	"dark:aria-invalid:ring-destructive/40",
	"aria-invalid:border-destructive",
).Variant("variant", map[string][]string{
	"default": {},
	"error": {
		"ring-destructive/20",
		"dark:ring-destructive/40",
		"border-destructive",
	},
}).Default("variant", "default")
