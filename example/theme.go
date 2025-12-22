package main

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/xraph/forgeui/alpine"
	"github.com/xraph/forgeui/components/button"
	"github.com/xraph/forgeui/components/card"
	"github.com/xraph/forgeui/primitives"
	"github.com/xraph/forgeui/theme"
)

// ThemeDemo returns the theme system demonstration page.
func ThemeDemo() g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			theme.HeadContent(theme.DefaultLight(), theme.DefaultDark()),
			html.TitleEl(g.Text("ForgeUI - Theme System")),
			html.Script(html.Src("https://cdn.tailwindcss.com")),
			theme.TailwindConfigScript(),
			alpine.CloakCSS(),
			theme.StyleTag(theme.DefaultLight(), theme.DefaultDark()),
		),
		html.Body(
			html.Class("min-h-screen bg-background text-foreground"),
			g.Attr("x-data", ""), // Initialize Alpine.js
			theme.DarkModeScript(),

			// Header
			primitives.Box(
				primitives.WithBackground("bg-card border-b border-border"),
				primitives.WithPadding("py-4"),
				primitives.WithChildren(
					primitives.Container(
						primitives.HStack("4",
							html.Div(
								html.Class("flex-1"),
								primitives.Text(
									primitives.TextAs("span"),
									primitives.TextSize("text-lg"),
									primitives.TextWeight("font-semibold"),
									primitives.TextChildren(g.Text("ForgeUI Theme System")),
								),
							),
							html.A(
								html.Href("/"),
								html.Class("text-sm text-muted-foreground hover:text-foreground transition-colors"),
								g.Text("← Back"),
							),
							theme.Toggle(theme.WithStorageType(theme.SessionStorage), theme.WithLabel(true)),
						),
					),
				),
			),

			// Main Content
			primitives.Box(
				primitives.WithPadding("py-12"),
				primitives.WithChildren(
					primitives.Container(
						primitives.VStack("12",
							// Introduction
							primitives.VStack("4",
								primitives.Text(
									primitives.TextAs("h1"),
									primitives.TextSize("text-4xl"),
									primitives.TextWeight("font-bold"),
									primitives.TextChildren(g.Text("Theme System")),
								),
								primitives.Text(
									primitives.TextSize("text-lg"),
									primitives.TextColor("text-muted-foreground"),
									primitives.TextChildren(g.Text("Comprehensive theming with dark mode support and color presets")),
								),
							),

							// Theme Toggles
							card.Card(
								card.Header(
									card.Title("Theme Controls"),
									card.Description("Toggle between themes and test the switcher components"),
								),
								card.Content(
									primitives.VStack("6",
										primitives.VStack("2",
											primitives.Text(
												primitives.TextWeight("font-semibold"),
												primitives.TextChildren(g.Text("Simple Toggle")),
											),
											theme.SimpleToggle(),
										),
										primitives.VStack("2",
											primitives.Text(
												primitives.TextWeight("font-semibold"),
												primitives.TextChildren(g.Text("Toggle with Label")),
											),
											theme.Toggle(theme.WithStorageType(theme.SessionStorage), theme.WithLabel(true)),
										),
										primitives.VStack("2",
											primitives.Text(
												primitives.TextWeight("font-semibold"),
												primitives.TextChildren(g.Text("Toggle with System Option")),
											),
											theme.ToggleWithSystemOption(theme.WithStorageType(theme.SessionStorage)),
										),
									),
								),
							),

							// Color Tokens
							card.Card(
								card.Header(
									card.Title("Color Tokens"),
									card.Description("All theme color tokens in action"),
								),
								card.Content(
									primitives.Grid(
										primitives.GridCols(2),
										primitives.GridColsMD(3),
										primitives.GridColsLG(4),
										primitives.GridGap("4"),
										primitives.GridChildren(
											colorToken("Background", "bg-background text-foreground", "Primary background color"),
											colorToken("Card", "bg-card text-card-foreground", "Card background"),
											colorToken("Popover", "bg-popover text-popover-foreground", "Popover background"),
											colorToken("Primary", "bg-primary text-primary-foreground", "Primary brand color"),
											colorToken("Secondary", "bg-secondary text-secondary-foreground", "Secondary color"),
											colorToken("Muted", "bg-muted text-muted-foreground", "Muted background"),
											colorToken("Accent", "bg-accent text-accent-foreground", "Accent color"),
											colorToken("Destructive", "bg-destructive text-destructive-foreground", "Danger/error color"),
										),
									),
								),
							),

							// Component Examples
							card.Card(
								card.Header(
									card.Title("Components with Theme"),
									card.Description("All components automatically adapt to the current theme"),
								),
								card.Content(
									primitives.VStack("4",
										primitives.HStack("2",
											button.Primary(g.Text("Primary")),
											button.Secondary(g.Text("Secondary")),
											button.Destructive(g.Text("Destructive")),
											button.Outline(g.Text("Outline")),
											button.Ghost(g.Text("Ghost")),
										),
										html.Div(
											html.Class("p-4 border border-border rounded-lg bg-muted/50"),
											primitives.Text(
												primitives.TextChildren(g.Text("This box uses muted background with border-border")),
											),
										),
										html.Div(
											html.Class("p-4 border border-border rounded-lg bg-card"),
											primitives.Text(
												primitives.TextChildren(g.Text("This box uses card background - perfect for dark mode")),
											),
										),
									),
								),
							),

							// Color Presets Info
							card.Card(
								card.Header(
									card.Title("Available Color Presets"),
									card.Description("ForgeUI comes with multiple theme presets"),
								),
								card.Content(
									primitives.VStack("3",
										presetInfo("Default (Neutral)", "Neutral gray tones suitable for most applications"),
										presetInfo("Rose", "Rose/pink tones for feminine or creative aesthetics"),
										presetInfo("Blue", "Blue tones for professional and corporate applications"),
										presetInfo("Green", "Green tones for eco-friendly and health applications"),
										presetInfo("Orange", "Orange tones for energetic and creative applications"),
									),
								),
							),

							// Usage Instructions
							card.Card(
								card.Header(
									card.Title("Usage"),
									card.Description("How to use the theme system in your application"),
								),
								card.Content(
									primitives.VStack("4",
										html.Pre(
											html.Class("p-4 bg-muted rounded-lg text-sm overflow-x-auto"),
											html.Code(
												g.Raw(`// Wrap your app with the theme provider
theme.Provider(
    yourAppContent...,
)

// Use custom themes
theme.ProviderWithThemes(
    theme.BlueLight(),
    theme.BlueDark(),
    yourAppContent...,
)

// Add theme toggle
theme.SimpleToggle()  // Icon only
theme.Toggle(theme.WithLabel(true))  // With label`),
											),
										),
									),
								),
							),
						),
					),
				),
			),

			// Alpine.js scripts
			alpine.Scripts(),
		),
	)
}

// colorToken renders a color token preview box.
func colorToken(name, classes, description string) g.Node {
	return html.Div(
		html.Class("border border-border rounded-lg overflow-hidden"),
		html.Div(
			html.Class(classes+" p-6 h-24 flex items-center justify-center"),
			primitives.Text(
				primitives.TextWeight("font-semibold"),
				primitives.TextChildren(g.Text(name)),
			),
		),
		html.Div(
			html.Class("p-3 bg-card"),
			primitives.Text(
				primitives.TextSize("text-xs"),
				primitives.TextColor("text-muted-foreground"),
				primitives.TextChildren(g.Text(description)),
			),
		),
	)
}

// presetInfo renders information about a theme preset.
func presetInfo(name, description string) g.Node {
	return html.Div(
		html.Class("flex flex-col gap-1"),
		primitives.Text(
			primitives.TextWeight("font-semibold"),
			primitives.TextChildren(g.Text(name)),
		),
		primitives.Text(
			primitives.TextSize("text-sm"),
			primitives.TextColor("text-muted-foreground"),
			primitives.TextChildren(g.Text(description)),
		),
	)
}
