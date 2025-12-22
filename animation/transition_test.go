package animation

import (
	"strings"
	"testing"
)

func TestNewTransition(t *testing.T) {
	builder := NewTransition()
	if builder == nil {
		t.Fatal("NewTransition() returned nil")
	}

	if builder.transition == nil {
		t.Fatal("NewTransition() created builder with nil transition")
	}
}

func TestTransitionBuilder(t *testing.T) {
	transition := NewTransition().
		Enter("enter-class").
		EnterStart("enter-start-class").
		EnterEnd("enter-end-class").
		Leave("leave-class").
		LeaveStart("leave-start-class").
		LeaveEnd("leave-end-class").
		Build()

	if transition.Enter != "enter-class" {
		t.Errorf("Enter = %v, want enter-class", transition.Enter)
	}

	if transition.EnterStart != "enter-start-class" {
		t.Errorf("EnterStart = %v, want enter-start-class", transition.EnterStart)
	}

	if transition.EnterEnd != "enter-end-class" {
		t.Errorf("EnterEnd = %v, want enter-end-class", transition.EnterEnd)
	}

	if transition.Leave != "leave-class" {
		t.Errorf("Leave = %v, want leave-class", transition.Leave)
	}

	if transition.LeaveStart != "leave-start-class" {
		t.Errorf("LeaveStart = %v, want leave-start-class", transition.LeaveStart)
	}

	if transition.LeaveEnd != "leave-end-class" {
		t.Errorf("LeaveEnd = %v, want leave-end-class", transition.LeaveEnd)
	}
}

func TestFadeIn(t *testing.T) {
	transition := FadeIn()
	if transition == nil {
		t.Fatal("FadeIn() returned nil")
	}

	if !strings.Contains(transition.Enter, "transition-opacity") {
		t.Errorf("FadeIn() Enter missing transition-opacity")
	}

	if !strings.Contains(transition.EnterStart, "opacity-0") {
		t.Errorf("FadeIn() EnterStart missing opacity-0")
	}

	if !strings.Contains(transition.EnterEnd, "opacity-100") {
		t.Errorf("FadeIn() EnterEnd missing opacity-100")
	}
}

func TestFadeOut(t *testing.T) {
	transition := FadeOut()
	if transition == nil {
		t.Fatal("FadeOut() returned nil")
	}

	if !strings.Contains(transition.EnterStart, "opacity-100") {
		t.Errorf("FadeOut() should start with opacity-100")
	}

	if !strings.Contains(transition.EnterEnd, "opacity-0") {
		t.Errorf("FadeOut() should end with opacity-0")
	}
}

func TestScaleIn(t *testing.T) {
	transition := ScaleIn()
	if transition == nil {
		t.Fatal("ScaleIn() returned nil")
	}

	if !strings.Contains(transition.EnterStart, "scale-95") {
		t.Errorf("ScaleIn() EnterStart missing scale-95")
	}

	if !strings.Contains(transition.EnterStart, "opacity-0") {
		t.Errorf("ScaleIn() EnterStart missing opacity-0")
	}

	if !strings.Contains(transition.EnterEnd, "scale-100") {
		t.Errorf("ScaleIn() EnterEnd missing scale-100")
	}

	if !strings.Contains(transition.EnterEnd, "opacity-100") {
		t.Errorf("ScaleIn() EnterEnd missing opacity-100")
	}
}

func TestScaleOut(t *testing.T) {
	transition := ScaleOut()
	if transition == nil {
		t.Fatal("ScaleOut() returned nil")
	}

	if !strings.Contains(transition.EnterStart, "scale-100") {
		t.Errorf("ScaleOut() should start with scale-100")
	}

	if !strings.Contains(transition.EnterEnd, "scale-95") {
		t.Errorf("ScaleOut() should end with scale-95")
	}
}

func TestSlideTransitions(t *testing.T) {
	tests := []struct {
		name       string
		transition *Transition
		wantStart  string
	}{
		{
			name:       "SlideUp",
			transition: SlideUp(),
			wantStart:  "translate-y-2",
		},
		{
			name:       "SlideDown",
			transition: SlideDown(),
			wantStart:  "-translate-y-2",
		},
		{
			name:       "SlideLeft",
			transition: SlideLeft(),
			wantStart:  "translate-x-2",
		},
		{
			name:       "SlideRight",
			transition: SlideRight(),
			wantStart:  "-translate-x-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.transition == nil {
				t.Fatalf("%s() returned nil", tt.name)
			}

			if !strings.Contains(tt.transition.EnterStart, tt.wantStart) {
				t.Errorf("%s() EnterStart = %v, want to contain %v",
					tt.name, tt.transition.EnterStart, tt.wantStart)
			}

			if !strings.Contains(tt.transition.EnterEnd, "translate-y-0") &&
				!strings.Contains(tt.transition.EnterEnd, "translate-x-0") {
				t.Errorf("%s() EnterEnd should contain translate to 0", tt.name)
			}
		})
	}
}

func TestLargeSlideTransitions(t *testing.T) {
	tests := []struct {
		name       string
		transition *Transition
		wantStart  string
		wantEnd    string
	}{
		{
			name:       "SlideInFromBottom",
			transition: SlideInFromBottom(),
			wantStart:  "translate-y-full",
			wantEnd:    "translate-y-0",
		},
		{
			name:       "SlideInFromTop",
			transition: SlideInFromTop(),
			wantStart:  "-translate-y-full",
			wantEnd:    "translate-y-0",
		},
		{
			name:       "SlideInFromLeft",
			transition: SlideInFromLeft(),
			wantStart:  "-translate-x-full",
			wantEnd:    "translate-x-0",
		},
		{
			name:       "SlideInFromRight",
			transition: SlideInFromRight(),
			wantStart:  "translate-x-full",
			wantEnd:    "translate-x-0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.transition == nil {
				t.Fatalf("%s() returned nil", tt.name)
			}

			if !strings.Contains(tt.transition.EnterStart, tt.wantStart) {
				t.Errorf("%s() EnterStart = %v, want to contain %v",
					tt.name, tt.transition.EnterStart, tt.wantStart)
			}

			if !strings.Contains(tt.transition.EnterEnd, tt.wantEnd) {
				t.Errorf("%s() EnterEnd = %v, want to contain %v",
					tt.name, tt.transition.EnterEnd, tt.wantEnd)
			}
		})
	}
}

func TestCollapse(t *testing.T) {
	transition := Collapse()
	if transition == nil {
		t.Fatal("Collapse() returned nil")
	}

	if !strings.Contains(transition.EnterStart, "scale-y-0") {
		t.Errorf("Collapse() should start with scale-y-0")
	}

	if !strings.Contains(transition.EnterEnd, "scale-y-100") {
		t.Errorf("Collapse() should end with scale-y-100")
	}
}

func TestZoomTransitions(t *testing.T) {
	zoomIn := ZoomIn()
	if !strings.Contains(zoomIn.EnterStart, "scale-50") {
		t.Errorf("ZoomIn() should start with scale-50")
	}

	zoomOut := ZoomOut()
	if !strings.Contains(zoomOut.EnterStart, "scale-150") {
		t.Errorf("ZoomOut() should start with scale-150")
	}
}

func TestRotateIn(t *testing.T) {
	transition := RotateIn()
	if transition == nil {
		t.Fatal("RotateIn() returned nil")
	}

	if !strings.Contains(transition.EnterStart, "rotate-180") {
		t.Errorf("RotateIn() should start with rotate-180")
	}

	if !strings.Contains(transition.EnterEnd, "rotate-0") {
		t.Errorf("RotateIn() should end with rotate-0")
	}
}

func TestTransitionChaining(t *testing.T) {
	// Test that builder methods return the builder for chaining
	builder := NewTransition()

	result := builder.Enter("test")
	if result != builder {
		t.Error("Enter() should return the builder for chaining")
	}

	result = builder.EnterStart("test")
	if result != builder {
		t.Error("EnterStart() should return the builder for chaining")
	}

	result = builder.EnterEnd("test")
	if result != builder {
		t.Error("EnterEnd() should return the builder for chaining")
	}

	result = builder.Leave("test")
	if result != builder {
		t.Error("Leave() should return the builder for chaining")
	}

	result = builder.LeaveStart("test")
	if result != builder {
		t.Error("LeaveStart() should return the builder for chaining")
	}

	result = builder.LeaveEnd("test")
	if result != builder {
		t.Error("LeaveEnd() should return the builder for chaining")
	}
}

func TestTransitionDurations(t *testing.T) {
	// Test that transitions have reasonable durations
	tests := []struct {
		name       string
		transition *Transition
	}{
		{"FadeIn", FadeIn()},
		{"ScaleIn", ScaleIn()},
		{"SlideUp", SlideUp()},
		{"SlideInFromBottom", SlideInFromBottom()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.transition.Enter, "duration-") {
				t.Errorf("%s() Enter missing duration", tt.name)
			}

			if !strings.Contains(tt.transition.Leave, "duration-") {
				t.Errorf("%s() Leave missing duration", tt.name)
			}
		})
	}
}
