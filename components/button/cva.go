package button

import "github.com/xraph/forgeui"

// buttonCVA defines the class variance authority configuration for buttons
// This follows shadcn/ui button styling patterns with modern focus states and dark mode support
var buttonCVA = forgeui.NewCVA(
	// Base classes - always applied
	"inline-flex",
	"items-center",
	"justify-center",
	"gap-2",
	"whitespace-nowrap",
	"rounded-md",
	"text-sm",
	"font-medium",
	"font-semibold",
	"transition-all",
	"duration-200",
	"shrink-0",
	"[&_svg]:pointer-events-none",
	"[&_svg:not([class*='size-'])]:size-4",
	"[&_svg]:shrink-0",
	"outline-none",
	"ring-offset-background",
	"focus-visible:outline-none",
	"focus-visible:ring-2",
	"focus-visible:ring-ring",
	"focus-visible:ring-offset-2",
	"disabled:pointer-events-none",
	"disabled:opacity-50",
).Variant("variant", map[string][]string{
	"default": {
		"bg-primary",
		"text-primary-foreground",
		"shadow-sm",
		"hover:bg-primary/90",
		"hover:shadow-md",
	},
	"destructive": {
		"bg-destructive",
		"text-destructive-foreground",
		"shadow-sm",
		"hover:bg-destructive/90",
		"hover:shadow-md",
	},
	"outline": {
		"border",
		"border-input",
		"bg-background",
		"hover:bg-accent",
		"hover:text-accent-foreground",
		"hover:border-accent-foreground/20",
	},
	"secondary": {
		"bg-secondary",
		"text-secondary-foreground",
		"shadow-sm",
		"hover:bg-secondary/80",
	},
	"ghost": {
		"hover:bg-accent",
		"hover:text-accent-foreground",
	},
	"link": {
		"text-primary",
		"underline-offset-4",
		"hover:underline",
	},
}).Variant("size", map[string][]string{
	"sm":      {"h-8", "rounded-md", "gap-1.5", "px-3", "text-xs", "has-[>svg]:px-2.5"},
	"md":      {"h-9", "px-4", "py-2", "has-[>svg]:px-3"},
	"lg":      {"h-10", "rounded-md", "px-8", "has-[>svg]:px-5"},
	"icon":    {"size-9"},
	"icon-sm": {"size-8"},
	"icon-lg": {"size-10"},
}).Default("variant", "default").Default("size", "md")
