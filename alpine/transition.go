package alpine

import (
	g "maragu.dev/gomponents"
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
func (t *Transition) Attrs() []g.Node {
	var attrs []g.Node

	if t.Enter != "" {
		attrs = append(attrs, g.Attr("x-transition:enter", t.Enter))
	}
	if t.EnterStart != "" {
		attrs = append(attrs, g.Attr("x-transition:enter-start", t.EnterStart))
	}
	if t.EnterEnd != "" {
		attrs = append(attrs, g.Attr("x-transition:enter-end", t.EnterEnd))
	}
	if t.Leave != "" {
		attrs = append(attrs, g.Attr("x-transition:leave", t.Leave))
	}
	if t.LeaveStart != "" {
		attrs = append(attrs, g.Attr("x-transition:leave-start", t.LeaveStart))
	}
	if t.LeaveEnd != "" {
		attrs = append(attrs, g.Attr("x-transition:leave-end", t.LeaveEnd))
	}

	return attrs
}

// XTransition applies transition attributes to an element.
// Accepts both alpine.Transition and animation.Transition.
//
// Example:
//
//	html.Div(
//	    alpine.XShow("open"),
//	    g.Group(alpine.XTransition(myTransition)),
//	    g.Text("Content"),
//	)
func XTransition(t interface{}) []g.Node {
	switch v := t.(type) {
	case *Transition:
		return v.Attrs()
	case *animation.Transition:
		// Convert animation.Transition to alpine attributes
		return convertAnimationTransition(v)
	default:
		return nil
	}
}

// convertAnimationTransition converts an animation.Transition to Alpine attributes.
func convertAnimationTransition(t *animation.Transition) []g.Node {
	var attrs []g.Node

	if t.Enter != "" {
		attrs = append(attrs, g.Attr("x-transition:enter", t.Enter))
	}
	if t.EnterStart != "" {
		attrs = append(attrs, g.Attr("x-transition:enter-start", t.EnterStart))
	}
	if t.EnterEnd != "" {
		attrs = append(attrs, g.Attr("x-transition:enter-end", t.EnterEnd))
	}
	if t.Leave != "" {
		attrs = append(attrs, g.Attr("x-transition:leave", t.Leave))
	}
	if t.LeaveStart != "" {
		attrs = append(attrs, g.Attr("x-transition:leave-start", t.LeaveStart))
	}
	if t.LeaveEnd != "" {
		attrs = append(attrs, g.Attr("x-transition:leave-end", t.LeaveEnd))
	}

	return attrs
}

// XTransitionSimple creates a simple x-transition attribute without custom classes.
// Alpine will use default transition behavior.
func XTransitionSimple() g.Node {
	return g.Attr("x-transition", "")
}

// XTransitionDuration creates an x-transition with custom duration.
// Duration should be in milliseconds.
//
// Example:
//
//	alpine.XTransitionDuration(300) // 300ms transition
func XTransitionDuration(ms int) g.Node {
	return g.Attr("x-transition.duration", "")
}

// XCollapse creates an x-collapse attribute for height transitions.
// Requires the Collapse plugin to be loaded.
//
// Example:
//
//	html.Div(
//	    alpine.XShow("expanded"),
//	    alpine.XCollapse(),
//	    g.Text("Collapsible content"),
//	)
func XCollapse() g.Node {
	return g.Attr("x-collapse", "")
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

