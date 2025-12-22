package card

import "github.com/xraph/forgeui"

// cardCVA defines the class variance authority configuration for cards
// Following shadcn/ui patterns with subtle shadows and proper spacing
var cardCVA = forgeui.NewCVA(
	"bg-card",
	"text-card-foreground",
	"flex",
	"flex-col",
	"gap-6",
	"rounded-lg",
	"border",
	"border-border",
	"shadow-sm",
	"py-8",
)
