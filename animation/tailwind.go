package animation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TailwindAnimation represents a Tailwind animation configuration.
type TailwindAnimation struct {
	Name      string
	Animation string
}

// TailwindKeyframe represents a Tailwind keyframe configuration.
type TailwindKeyframe struct {
	Name   string
	Frames map[string]map[string]string
}

// TailwindConfig generates Tailwind CSS configuration for animations.
// This can be merged into your tailwind.config.js theme.extend section.
//
// Example usage:
//
//	config := animation.TailwindConfig()
//	// Add to tailwind.config.js:
//	// module.exports = {
//	//   theme: {
//	//     extend: {
//	//       keyframes: config.Keyframes,
//	//       animation: config.Animations,
//	//     }
//	//   }
//	// }
type TailwindConfig struct {
	Keyframes  map[string]map[string]map[string]string `json:"keyframes"`
	Animations map[string]string                       `json:"animation"`
}

// GenerateTailwindConfig creates a TailwindConfig with predefined animations.
func GenerateTailwindConfig() *TailwindConfig {
	config := &TailwindConfig{
		Keyframes:  make(map[string]map[string]map[string]string),
		Animations: make(map[string]string),
	}

	// Add keyframes
	keyframes := PredefinedKeyframes()
	for _, k := range keyframes {
		frames := make(map[string]map[string]string)
		for selector, properties := range k.Frames {
			// Parse properties string into map
			props := make(map[string]string)
			// Simple parsing - in real usage, properties would be structured
			props["raw"] = properties
			frames[selector] = props
		}

		config.Keyframes[k.Name] = frames
	}

	// Add animations
	config.Animations["fade-in"] = "fadeIn 0.2s ease-out"
	config.Animations["fade-out"] = "fadeOut 0.15s ease-in"
	config.Animations["slide-in-from-top"] = "slideInFromTop 0.3s ease-out"
	config.Animations["slide-in-from-bottom"] = "slideInFromBottom 0.3s ease-out"
	config.Animations["slide-in-from-left"] = "slideInFromLeft 0.3s ease-out"
	config.Animations["slide-in-from-right"] = "slideInFromRight 0.3s ease-out"
	config.Animations["scale-in"] = "scaleIn 0.2s ease-out"
	config.Animations["scale-out"] = "scaleOut 0.15s ease-in"
	config.Animations["spin"] = "spin 1s linear infinite"
	config.Animations["pulse"] = "pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite"
	config.Animations["bounce"] = "bounce 1s infinite"
	config.Animations["shimmer"] = "shimmer 2s infinite"
	config.Animations["shake"] = "shake 0.5s"

	return config
}

// ToJSON converts the Tailwind config to JSON format.
func (tc *TailwindConfig) ToJSON() (string, error) {
	data, err := json.MarshalIndent(tc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	return string(data), nil
}

// ToJavaScript generates JavaScript code for tailwind.config.js.
func (tc *TailwindConfig) ToJavaScript() string {
	var js strings.Builder
	js.WriteString("// ForgeUI Animation Configuration for Tailwind CSS\n")
	js.WriteString("// Add this to your tailwind.config.js theme.extend section\n\n")
	js.WriteString("module.exports = {\n")
	js.WriteString("  theme: {\n")
	js.WriteString("    extend: {\n")

	// Keyframes
	js.WriteString("      keyframes: {\n")

	var jsSb98 strings.Builder

	var jsSb102 strings.Builder

	for name := range tc.Keyframes {
		jsSb98.WriteString(fmt.Sprintf("        '%s': {\n", name))

		var jsSb100 strings.Builder

		var jsSb105 strings.Builder

		for selector, props := range tc.Keyframes[name] {
			jsSb100.WriteString(fmt.Sprintf("          '%s': {\n", selector))

			var jsSb102 strings.Builder
			for key, value := range props {
				jsSb102.WriteString(fmt.Sprintf("            %s: '%s',\n", key, value))
			}

			jsSb105.WriteString(jsSb102.String())

			jsSb100.WriteString("          },\n")
		}

		js.WriteString(jsSb105.String())

		jsSb102.WriteString(jsSb100.String())

		jsSb98.WriteString("        },\n")
	}

	js.WriteString(jsSb102.String())

	js.WriteString(jsSb98.String())

	js.WriteString("      },\n")

	// Animations
	js.WriteString("      animation: {\n")

	var jsSb113 strings.Builder
	for name, value := range tc.Animations {
		jsSb113.WriteString(fmt.Sprintf("        '%s': '%s',\n", name, value))
	}

	js.WriteString(jsSb113.String())

	js.WriteString("      },\n")

	js.WriteString("    },\n")
	js.WriteString("  },\n")
	js.WriteString("};\n")

	return js.String()
}

// TailwindClasses returns utility class names that can be used in ForgeUI components.
func TailwindClasses() []string {
	return []string{
		"animate-fade-in",
		"animate-fade-out",
		"animate-slide-in-from-top",
		"animate-slide-in-from-bottom",
		"animate-slide-in-from-left",
		"animate-slide-in-from-right",
		"animate-scale-in",
		"animate-scale-out",
		"animate-spin",
		"animate-pulse",
		"animate-bounce",
		"animate-shimmer",
		"animate-shake",
	}
}
