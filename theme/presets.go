package theme

// DefaultLight returns the default light theme following shadcn/ui design patterns.
// This uses neutral colors in OKLCH format suitable for most applications.
func DefaultLight() Theme {
	return Theme{
		Colors: ColorTokens{
			Background:               "1 0 0",               // Pure white
			Foreground:               "0.09 0 0",            // Near black
			Card:                     "1 0 0",               // Pure white (same as background for clean look)
			CardForeground:           "0.09 0 0",            // Near black
			Popover:                  "1 0 0",               // Pure white
			PopoverForeground:        "0.09 0 0",            // Near black
			Primary:                  "0.09 0 0",            // Near black for strong contrast
			PrimaryForeground:        "0.985 0 0",           // Near white
			Secondary:                "0.96 0 0",            // Very light gray
			SecondaryForeground:      "0.09 0 0",            // Near black
			Muted:                    "0.96 0 0",            // Very light gray
			MutedForeground:          "0.52 0 0",            // Medium gray for better readability
			Accent:                   "0.96 0 0",            // Very light gray
			AccentForeground:         "0.09 0 0",            // Near black
			Destructive:              "0.577 0.245 27.325",  // Red
			DestructiveForeground:    "0.985 0 0",           // Near white for contrast
			Success:                  "0.508 0.118 165.612", // Green
			Border:                   "0.898 0 0",           // Lighter border for subtle separation
			Input:                    "0.898 0 0",           // Lighter input border
			Ring:                     "0.09 0 0",            // Dark ring for visibility
			Chart1:                   "0.646 0.222 41.116",
			Chart2:                   "0.6 0.118 184.704",
			Chart3:                   "0.398 0.07 227.392",
			Chart4:                   "0.828 0.189 84.429",
			Chart5:                   "0.769 0.188 70.08",
			Chart6:                   "0.646 0.222 41.116",
			Chart7:                   "0.6 0.118 184.704",
			Chart8:                   "0.398 0.07 227.392",
			Chart9:                   "0.828 0.189 84.429",
			Chart10:                  "0.769 0.188 70.08",
			Chart11:                  "0.646 0.222 41.116",
			Chart12:                  "0.6 0.118 184.704",
			Sidebar:                  "0.985 0 0",
			SidebarForeground:        "0.145 0 0",
			SidebarPrimary:           "0.205 0 0",
			SidebarPrimaryForeground: "0.985 0 0",
			SidebarAccent:            "0.97 0 0",
			SidebarAccentForeground:  "0.205 0 0",
			SidebarBorder:            "0.922 0 0",
			SidebarRing:              "0.708 0 0",
		},
		Radius:   defaultRadius(),
		Spacing:  defaultSpacing(),
		FontSize: defaultFontSize(),
		Shadow:   defaultShadow(),
	}
}

// DefaultDark returns the default dark theme following shadcn/ui design patterns.
// This uses neutral colors in OKLCH format suitable for dark mode.
func DefaultDark() Theme {
	return Theme{
		Colors: ColorTokens{
			Background:               "0.145 0 0",           // Deep dark matching shadcn exactly
			Foreground:               "0.985 0 0",           // Near white
			Card:                     "0.145 0 0",           // Same as background for consistency
			CardForeground:           "0.985 0 0",           // Near white
			Popover:                  "0.145 0 0",           // Deep dark
			PopoverForeground:        "0.985 0 0",           // Near white
			Primary:                  "0.985 0 0",           // Near white for contrast
			PrimaryForeground:        "0.205 0 0",           // Darker for better contrast
			Secondary:                "0.198 0 0",           // Subtle dark gray (shadcn 14.9%)
			SecondaryForeground:      "0.985 0 0",           // Near white
			Muted:                    "0.198 0 0",           // Subtle dark gray (shadcn 14.9%)
			MutedForeground:          "0.67 0 0",            // Medium gray (shadcn 63.9%)
			Accent:                   "0.198 0 0",           // Subtle dark gray (shadcn 14.9%)
			AccentForeground:         "0.985 0 0",           // Near white
			Destructive:              "0.396 0.141 25.723",  // Darker red for dark mode
			DestructiveForeground:    "0.637 0.237 25.331",  // Lighter destructive foreground
			Success:                  "0.508 0.118 165.612", // Green matching light mode
			Border:                   "0.198 0 0",           // Subtle border (shadcn 14.9%)
			Input:                    "0.198 0 0",           // Subtle input border (shadcn 14.9%)
			Ring:                     "0.85 0 0",            // Focus ring (shadcn 83.1%)
			Chart1:                   "0.488 0.243 264.376",
			Chart2:                   "0.696 0.17 162.48",
			Chart3:                   "0.769 0.188 70.08",
			Chart4:                   "0.627 0.265 303.9",
			Chart5:                   "0.645 0.246 16.439",
			Chart6:                   "0.488 0.243 264.376",
			Chart7:                   "0.696 0.17 162.48",
			Chart8:                   "0.769 0.188 70.08",
			Chart9:                   "0.627 0.265 303.9",
			Chart10:                  "0.645 0.246 16.439",
			Chart11:                  "0.488 0.243 264.376",
			Chart12:                  "0.696 0.17 162.48",
			Sidebar:                  "0.205 0 0",
			SidebarForeground:        "0.985 0 0",
			SidebarPrimary:           "0.488 0.243 264.376",
			SidebarPrimaryForeground: "0.985 0 0",
			SidebarAccent:            "0.198 0 0", // Subtle sidebar accent
			SidebarAccentForeground:  "0.985 0 0",
			SidebarBorder:            "0.198 0 0", // Subtle sidebar border
			SidebarRing:              "0.85 0 0",  // Sidebar focus ring
		},
		Radius:   defaultRadius(),
		Spacing:  defaultSpacing(),
		FontSize: defaultFontSize(),
		Shadow:   defaultShadow(),
	}
}

// NeutralLight returns a neutral-toned light theme.
// Alias for DefaultLight for semantic clarity.
func NeutralLight() Theme {
	return DefaultLight()
}

// NeutralDark returns a neutral-toned dark theme.
// Alias for DefaultDark for semantic clarity.
func NeutralDark() Theme {
	return DefaultDark()
}

// RoseLight returns a rose/pink-toned light theme in OKLCH format.
// Suitable for applications targeting feminine or creative aesthetics.
func RoseLight() Theme {
	theme := DefaultLight()
	theme.Colors.Primary = "0.577 0.245 27.325"
	theme.Colors.PrimaryForeground = "0.985 0.01 27.325"
	theme.Colors.Secondary = "0.97 0.005 264.376"
	theme.Colors.SecondaryForeground = "0.577 0.245 27.325"
	theme.Colors.Accent = "0.97 0.005 264.376"
	theme.Colors.AccentForeground = "0.577 0.245 27.325"
	theme.Colors.Ring = "0.577 0.245 27.325"
	theme.Colors.Chart1 = "0.577 0.245 27.325"
	theme.Colors.Chart2 = "0.627 0.265 303.9"
	theme.Colors.Chart3 = "0.488 0.243 264.376"
	theme.Colors.Chart4 = "0.769 0.188 16.439"
	theme.Colors.Chart5 = "0.696 0.17 16.439"

	return theme
}

// RoseDark returns a rose/pink-toned dark theme in OKLCH format.
// Suitable for applications targeting feminine or creative aesthetics.
func RoseDark() Theme {
	theme := DefaultDark()
	theme.Colors.Primary = "0.577 0.245 27.325"
	theme.Colors.PrimaryForeground = "0.985 0.01 27.325"
	theme.Colors.Secondary = "0.2 0.005 264.376"
	theme.Colors.SecondaryForeground = "0.577 0.245 27.325"
	theme.Colors.Accent = "0.2 0.005 264.376"
	theme.Colors.AccentForeground = "0.577 0.245 27.325"
	theme.Colors.Ring = "0.577 0.245 27.325"
	theme.Colors.Chart1 = "0.577 0.245 27.325"
	theme.Colors.Chart2 = "0.627 0.265 303.9"
	theme.Colors.Chart3 = "0.488 0.243 264.376"
	theme.Colors.Chart4 = "0.769 0.188 16.439"
	theme.Colors.Chart5 = "0.696 0.17 16.439"

	return theme
}

// BlueLight returns a blue-toned light theme in OKLCH format.
// Suitable for professional, corporate, or tech-focused applications.
func BlueLight() Theme {
	theme := DefaultLight()
	theme.Colors.Primary = "0.6 0.222 264.376"
	theme.Colors.PrimaryForeground = "0.985 0 0"
	theme.Colors.Secondary = "0.97 0 0"
	theme.Colors.SecondaryForeground = "0.6 0.222 264.376"
	theme.Colors.Accent = "0.97 0 0"
	theme.Colors.AccentForeground = "0.6 0.222 264.376"
	theme.Colors.Ring = "0.6 0.222 264.376"
	theme.Colors.Chart1 = "0.6 0.222 264.376"
	theme.Colors.Chart2 = "0.696 0.17 264.376"
	theme.Colors.Chart3 = "0.646 0.25 264.376"
	theme.Colors.Chart4 = "0.828 0.15 264.376"
	theme.Colors.Chart5 = "0.558 0.18 220.704"

	return theme
}

// BlueDark returns a blue-toned dark theme in OKLCH format.
// Suitable for professional, corporate, or tech-focused applications.
func BlueDark() Theme {
	theme := DefaultDark()
	theme.Colors.Primary = "0.646 0.25 264.376"
	theme.Colors.PrimaryForeground = "0.205 0 0"
	theme.Colors.Secondary = "0.269 0.07 264.376"
	theme.Colors.SecondaryForeground = "0.646 0.25 264.376"
	theme.Colors.Accent = "0.269 0.07 264.376"
	theme.Colors.AccentForeground = "0.646 0.25 264.376"
	theme.Colors.Ring = "0.646 0.25 264.376"
	theme.Colors.Chart1 = "0.646 0.25 264.376"
	theme.Colors.Chart2 = "0.696 0.17 264.376"
	theme.Colors.Chart3 = "0.646 0.25 264.376"
	theme.Colors.Chart4 = "0.828 0.15 264.376"
	theme.Colors.Chart5 = "0.558 0.18 220.704"

	return theme
}

// GreenLight returns a green-toned light theme in OKLCH format.
// Suitable for eco-friendly, health, or nature-focused applications.
func GreenLight() Theme {
	theme := DefaultLight()
	theme.Colors.Primary = "0.508 0.118 165.612"
	theme.Colors.PrimaryForeground = "0.985 0.01 165.612"
	theme.Colors.Secondary = "0.97 0.005 165.612"
	theme.Colors.SecondaryForeground = "0.508 0.118 165.612"
	theme.Colors.Accent = "0.97 0.005 165.612"
	theme.Colors.AccentForeground = "0.508 0.118 165.612"
	theme.Colors.Ring = "0.508 0.118 165.612"
	theme.Colors.Chart1 = "0.508 0.118 165.612"
	theme.Colors.Chart2 = "0.558 0.14 165.612"
	theme.Colors.Chart3 = "0.6 0.118 184.704"
	theme.Colors.Chart4 = "0.558 0.10 152.48"
	theme.Colors.Chart5 = "0.646 0.118 140.48"

	return theme
}

// GreenDark returns a green-toned dark theme in OKLCH format.
// Suitable for eco-friendly, health, or nature-focused applications.
func GreenDark() Theme {
	theme := DefaultDark()
	theme.Colors.Primary = "0.558 0.14 165.612"
	theme.Colors.PrimaryForeground = "0.15 0.12 165.612"
	theme.Colors.Secondary = "0.2 0.005 165.612"
	theme.Colors.SecondaryForeground = "0.558 0.14 165.612"
	theme.Colors.Accent = "0.2 0.005 165.612"
	theme.Colors.AccentForeground = "0.558 0.14 165.612"
	theme.Colors.Ring = "0.558 0.14 165.612"
	theme.Colors.Chart1 = "0.558 0.14 165.612"
	theme.Colors.Chart2 = "0.508 0.14 165.612"
	theme.Colors.Chart3 = "0.6 0.118 184.704"
	theme.Colors.Chart4 = "0.558 0.10 152.48"
	theme.Colors.Chart5 = "0.646 0.118 140.48"

	return theme
}

// OrangeLight returns an orange-toned light theme in OKLCH format.
// Suitable for energetic, creative, or fun applications.
func OrangeLight() Theme {
	theme := DefaultLight()
	theme.Colors.Primary = "0.646 0.222 41.116"
	theme.Colors.PrimaryForeground = "0.985 0.01 41.116"
	theme.Colors.Secondary = "0.97 0.005 41.116"
	theme.Colors.SecondaryForeground = "0.646 0.222 41.116"
	theme.Colors.Accent = "0.97 0.005 41.116"
	theme.Colors.AccentForeground = "0.646 0.222 41.116"
	theme.Colors.Ring = "0.646 0.222 41.116"
	theme.Colors.Chart1 = "0.646 0.222 41.116"
	theme.Colors.Chart2 = "0.6 0.22 35.116"
	theme.Colors.Chart3 = "0.558 0.21 30.116"
	theme.Colors.Chart4 = "0.696 0.20 45.116"
	theme.Colors.Chart5 = "0.558 0.23 50.116"

	return theme
}

// OrangeDark returns an orange-toned dark theme in OKLCH format.
// Suitable for energetic, creative, or fun applications.
func OrangeDark() Theme {
	theme := DefaultDark()
	theme.Colors.Primary = "0.6 0.22 35.116"
	theme.Colors.PrimaryForeground = "0.985 0.01 35.116"
	theme.Colors.Secondary = "0.2 0.005 35.116"
	theme.Colors.SecondaryForeground = "0.6 0.22 35.116"
	theme.Colors.Accent = "0.2 0.005 35.116"
	theme.Colors.AccentForeground = "0.6 0.22 35.116"
	theme.Colors.Ring = "0.6 0.22 35.116"
	theme.Colors.Chart1 = "0.6 0.22 35.116"
	theme.Colors.Chart2 = "0.646 0.222 41.116"
	theme.Colors.Chart3 = "0.558 0.21 30.116"
	theme.Colors.Chart4 = "0.696 0.20 45.116"
	theme.Colors.Chart5 = "0.558 0.23 50.116"

	return theme
}
