package animation

// Keyframe represents a CSS @keyframes animation definition.
type Keyframe struct {
	Name   string
	Frames map[string]string // percentage -> CSS properties
}

// PredefinedKeyframes returns commonly used CSS keyframes for animations.
// These can be added to your CSS or generated for Tailwind config.
func PredefinedKeyframes() []Keyframe {
	return []Keyframe{
		{
			Name: "fadeIn",
			Frames: map[string]string{
				"from": "opacity: 0",
				"to":   "opacity: 1",
			},
		},
		{
			Name: "fadeOut",
			Frames: map[string]string{
				"from": "opacity: 1",
				"to":   "opacity: 0",
			},
		},
		{
			Name: "slideInFromTop",
			Frames: map[string]string{
				"from": "transform: translateY(-100%)",
				"to":   "transform: translateY(0)",
			},
		},
		{
			Name: "slideInFromBottom",
			Frames: map[string]string{
				"from": "transform: translateY(100%)",
				"to":   "transform: translateY(0)",
			},
		},
		{
			Name: "slideInFromLeft",
			Frames: map[string]string{
				"from": "transform: translateX(-100%)",
				"to":   "transform: translateX(0)",
			},
		},
		{
			Name: "slideInFromRight",
			Frames: map[string]string{
				"from": "transform: translateX(100%)",
				"to":   "transform: translateX(0)",
			},
		},
		{
			Name: "scaleIn",
			Frames: map[string]string{
				"from": "transform: scale(0.95); opacity: 0",
				"to":   "transform: scale(1); opacity: 1",
			},
		},
		{
			Name: "scaleOut",
			Frames: map[string]string{
				"from": "transform: scale(1); opacity: 1",
				"to":   "transform: scale(0.95); opacity: 0",
			},
		},
		{
			Name: "spin",
			Frames: map[string]string{
				"from": "transform: rotate(0deg)",
				"to":   "transform: rotate(360deg)",
			},
		},
		{
			Name: "pulse",
			Frames: map[string]string{
				"0%, 100%": "opacity: 1",
				"50%":      "opacity: 0.5",
			},
		},
		{
			Name: "bounce",
			Frames: map[string]string{
				"0%, 100%": "transform: translateY(0); animation-timing-function: cubic-bezier(0.8, 0, 1, 1)",
				"50%":      "transform: translateY(-25%); animation-timing-function: cubic-bezier(0, 0, 0.2, 1)",
			},
		},
		{
			Name: "shimmer",
			Frames: map[string]string{
				"0%":   "background-position: -1000px 0",
				"100%": "background-position: 1000px 0",
			},
		},
		{
			Name: "shake",
			Frames: map[string]string{
				"0%, 100%": "transform: translateX(0)",
				"10%, 30%, 50%, 70%, 90%": "transform: translateX(-10px)",
				"20%, 40%, 60%, 80%":      "transform: translateX(10px)",
			},
		},
	}
}

// ToCSS converts a keyframe to CSS @keyframes rule.
func (k *Keyframe) ToCSS() string {
	css := "@keyframes " + k.Name + " {\n"
	for selector, properties := range k.Frames {
		css += "  " + selector + " {\n"
		css += "    " + properties + ";\n"
		css += "  }\n"
	}
	css += "}"
	return css
}

// GenerateCSS generates CSS for all predefined keyframes.
func GenerateCSS() string {
	keyframes := PredefinedKeyframes()
	css := "/* ForgeUI Keyframe Animations */\n\n"
	for _, k := range keyframes {
		css += k.ToCSS() + "\n\n"
	}
	return css
}

