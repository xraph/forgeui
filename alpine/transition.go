package alpine

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/animation"
)

// Transition represents Alpine.js transition configuration.
// Used with x-transition or x-show for smooth animations.
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

// Attrs converts the transition to Alpine.js x-transition attributes.
func (t *Transition) Attrs() templ.Attributes {
	attrs := templ.Attributes{}

	if t.Enter != "" {
		attrs["x-transition:enter"] = t.Enter
	}

	if t.EnterStart != "" {
		attrs["x-transition:enter-start"] = t.EnterStart
	}

	if t.EnterEnd != "" {
		attrs["x-transition:enter-end"] = t.EnterEnd
	}

	if t.Leave != "" {
		attrs["x-transition:leave"] = t.Leave
	}

	if t.LeaveStart != "" {
		attrs["x-transition:leave-start"] = t.LeaveStart
	}

	if t.LeaveEnd != "" {
		attrs["x-transition:leave-end"] = t.LeaveEnd
	}

	return attrs
}

// XTransition applies transition attributes to an element.
// Accepts both alpine.Transition and animation.Transition.
//
// Example (in .templ files):
//
//	<div { alpine.XShow("open")... } { alpine.XTransition(myTransition)... }>
func XTransition(t any) templ.Attributes {
	switch v := t.(type) {
	case *Transition:
		return v.Attrs()
	case *animation.Transition:
		return convertAnimationTransition(v)
	default:
		return templ.Attributes{}
	}
}

// convertAnimationTransition converts an animation.Transition to Alpine attributes.
func convertAnimationTransition(t *animation.Transition) templ.Attributes {
	attrs := templ.Attributes{}

	if t.Enter != "" {
		attrs["x-transition:enter"] = t.Enter
	}

	if t.EnterStart != "" {
		attrs["x-transition:enter-start"] = t.EnterStart
	}

	if t.EnterEnd != "" {
		attrs["x-transition:enter-end"] = t.EnterEnd
	}

	if t.Leave != "" {
		attrs["x-transition:leave"] = t.Leave
	}

	if t.LeaveStart != "" {
		attrs["x-transition:leave-start"] = t.LeaveStart
	}

	if t.LeaveEnd != "" {
		attrs["x-transition:leave-end"] = t.LeaveEnd
	}

	return attrs
}

// XTransitionSimple creates a simple x-transition attribute without custom classes.
// Alpine will use default transition behavior.
func XTransitionSimple() templ.Attributes {
	return templ.Attributes{"x-transition": ""}
}

// XTransitionDuration creates an x-transition with custom duration.
// Duration should be in milliseconds.
func XTransitionDuration(ms int) templ.Attributes {
	return templ.Attributes{"x-transition.duration": ""}
}

// XCollapse creates an x-collapse attribute for height transitions.
// Requires the Collapse plugin to be loaded.
func XCollapse() templ.Attributes {
	return templ.Attributes{"x-collapse": ""}
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
