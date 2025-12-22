package animation

// Transition represents animation configuration for Alpine.js x-transition.
// It defines enter and leave animations with start, middle, and end states.
type Transition struct {
	// Enter classes applied during entire enter transition
	Enter string

	// EnterStart classes applied before enter, removed one frame after
	EnterStart string

	// EnterEnd classes applied one frame after enter starts, removed when transition ends
	EnterEnd string

	// Leave classes applied during entire leave transition
	Leave string

	// LeaveStart classes applied before leave, removed one frame after
	LeaveStart string

	// LeaveEnd classes applied one frame after leave starts, removed when transition ends
	LeaveEnd string
}

// TransitionBuilder provides a fluent API for building transitions.
type TransitionBuilder struct {
	transition *Transition
}

// NewTransition creates a new transition builder.
func NewTransition() *TransitionBuilder {
	return &TransitionBuilder{
		transition: &Transition{},
	}
}

// Enter sets the enter transition classes.
func (b *TransitionBuilder) Enter(classes string) *TransitionBuilder {
	b.transition.Enter = classes
	return b
}

// EnterStart sets the enter start classes.
func (b *TransitionBuilder) EnterStart(classes string) *TransitionBuilder {
	b.transition.EnterStart = classes
	return b
}

// EnterEnd sets the enter end classes.
func (b *TransitionBuilder) EnterEnd(classes string) *TransitionBuilder {
	b.transition.EnterEnd = classes
	return b
}

// Leave sets the leave transition classes.
func (b *TransitionBuilder) Leave(classes string) *TransitionBuilder {
	b.transition.Leave = classes
	return b
}

// LeaveStart sets the leave start classes.
func (b *TransitionBuilder) LeaveStart(classes string) *TransitionBuilder {
	b.transition.LeaveStart = classes
	return b
}

// LeaveEnd sets the leave end classes.
func (b *TransitionBuilder) LeaveEnd(classes string) *TransitionBuilder {
	b.transition.LeaveEnd = classes
	return b
}

// Build returns the completed transition.
func (b *TransitionBuilder) Build() *Transition {
	return b.transition
}

// Preset Transitions following shadcn/ui patterns

// FadeIn creates a fade-in transition (opacity: 0 → 1).
func FadeIn() *Transition {
	return &Transition{
		Enter:      "transition-opacity duration-200 ease-out",
		EnterStart: "opacity-0",
		EnterEnd:   "opacity-100",
		Leave:      "transition-opacity duration-150 ease-in",
		LeaveStart: "opacity-100",
		LeaveEnd:   "opacity-0",
	}
}

// FadeOut creates a fade-out transition (opacity: 1 → 0).
func FadeOut() *Transition {
	return &Transition{
		Enter:      "transition-opacity duration-150 ease-in",
		EnterStart: "opacity-100",
		EnterEnd:   "opacity-0",
		Leave:      "transition-opacity duration-200 ease-out",
		LeaveStart: "opacity-0",
		LeaveEnd:   "opacity-100",
	}
}

// ScaleIn creates a scale and fade-in transition (scale: 0.95 → 1, opacity: 0 → 1).
// Commonly used for modals and dialogs.
func ScaleIn() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 scale-95",
		EnterEnd:   "opacity-100 scale-100",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 scale-100",
		LeaveEnd:   "opacity-0 scale-95",
	}
}

// ScaleOut creates a scale and fade-out transition (scale: 1 → 0.95, opacity: 1 → 0).
func ScaleOut() *Transition {
	return &Transition{
		Enter:      "transition-all duration-150 ease-in",
		EnterStart: "opacity-100 scale-100",
		EnterEnd:   "opacity-0 scale-95",
		Leave:      "transition-all duration-200 ease-out",
		LeaveStart: "opacity-0 scale-95",
		LeaveEnd:   "opacity-100 scale-100",
	}
}

// SlideUp creates a slide-up transition (translateY: 10px → 0).
// Commonly used for dropdowns and popovers.
func SlideUp() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 translate-y-2",
		EnterEnd:   "opacity-100 translate-y-0",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 translate-y-0",
		LeaveEnd:   "opacity-0 translate-y-2",
	}
}

// SlideDown creates a slide-down transition (translateY: -10px → 0).
func SlideDown() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 -translate-y-2",
		EnterEnd:   "opacity-100 translate-y-0",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 translate-y-0",
		LeaveEnd:   "opacity-0 -translate-y-2",
	}
}

// SlideLeft creates a slide-left transition (translateX: 10px → 0).
func SlideLeft() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 translate-x-2",
		EnterEnd:   "opacity-100 translate-x-0",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 translate-x-0",
		LeaveEnd:   "opacity-0 translate-x-2",
	}
}

// SlideRight creates a slide-right transition (translateX: -10px → 0).
func SlideRight() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 -translate-x-2",
		EnterEnd:   "opacity-100 translate-x-0",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 translate-x-0",
		LeaveEnd:   "opacity-0 -translate-x-2",
	}
}

// SlideInFromBottom creates a large slide-up transition from the bottom.
// Commonly used for mobile menus and drawers.
func SlideInFromBottom() *Transition {
	return &Transition{
		Enter:      "transition-transform duration-300 ease-out",
		EnterStart: "translate-y-full",
		EnterEnd:   "translate-y-0",
		Leave:      "transition-transform duration-200 ease-in",
		LeaveStart: "translate-y-0",
		LeaveEnd:   "translate-y-full",
	}
}

// SlideInFromTop creates a large slide-down transition from the top.
func SlideInFromTop() *Transition {
	return &Transition{
		Enter:      "transition-transform duration-300 ease-out",
		EnterStart: "-translate-y-full",
		EnterEnd:   "translate-y-0",
		Leave:      "transition-transform duration-200 ease-in",
		LeaveStart: "translate-y-0",
		LeaveEnd:   "-translate-y-full",
	}
}

// SlideInFromLeft creates a large slide-right transition from the left.
// Commonly used for sidebars and side sheets.
func SlideInFromLeft() *Transition {
	return &Transition{
		Enter:      "transition-transform duration-300 ease-out",
		EnterStart: "-translate-x-full",
		EnterEnd:   "translate-x-0",
		Leave:      "transition-transform duration-200 ease-in",
		LeaveStart: "translate-x-0",
		LeaveEnd:   "-translate-x-full",
	}
}

// SlideInFromRight creates a large slide-left transition from the right.
func SlideInFromRight() *Transition {
	return &Transition{
		Enter:      "transition-transform duration-300 ease-out",
		EnterStart: "translate-x-full",
		EnterEnd:   "translate-x-0",
		Leave:      "transition-transform duration-200 ease-in",
		LeaveStart: "translate-x-0",
		LeaveEnd:   "translate-x-full",
	}
}

// Collapse creates a smooth height collapse transition.
// Note: This does NOT use x-transition. Instead, it's meant to be used with
// Alpine's x-collapse directive which requires the Collapse plugin.
//
// Example:
//
//	html.Div(
//	    alpine.XShow("expanded"),
//	    alpine.XCollapse(),
//	    g.Text("Collapsible content"),
//	)
func Collapse() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out origin-top",
		EnterStart: "opacity-0 scale-y-0",
		EnterEnd:   "opacity-100 scale-y-100",
		Leave:      "transition-all duration-150 ease-in origin-top",
		LeaveStart: "opacity-100 scale-y-100",
		LeaveEnd:   "opacity-0 scale-y-0",
	}
}

// ZoomIn creates a zoom-in transition with scale and opacity.
func ZoomIn() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 scale-50",
		EnterEnd:   "opacity-100 scale-100",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 scale-100",
		LeaveEnd:   "opacity-0 scale-50",
	}
}

// ZoomOut creates a zoom-out transition with scale and opacity.
func ZoomOut() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 scale-150",
		EnterEnd:   "opacity-100 scale-100",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 scale-100",
		LeaveEnd:   "opacity-0 scale-150",
	}
}

// RotateIn creates a rotation and fade-in transition.
func RotateIn() *Transition {
	return &Transition{
		Enter:      "transition-all duration-200 ease-out",
		EnterStart: "opacity-0 rotate-180 scale-50",
		EnterEnd:   "opacity-100 rotate-0 scale-100",
		Leave:      "transition-all duration-150 ease-in",
		LeaveStart: "opacity-100 rotate-0 scale-100",
		LeaveEnd:   "opacity-0 rotate-180 scale-50",
	}
}

