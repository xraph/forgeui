package input

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestInputGroup(t *testing.T) {
	t.Run("renders basic input group", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupInput(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "group/input-group") {
			t.Error("expected input group to have group/input-group class")
		}

		if !strings.Contains(output, "relative flex w-full") {
			t.Error("expected input group to have flex container classes")
		}

		if !strings.Contains(output, "data-slot=\"input-group\"") {
			t.Error("expected input group to have data-slot attribute")
		}

		if !strings.Contains(output, "role=\"group\"") {
			t.Error("expected input group to have role attribute")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{WithGroupClass("custom-class")},
			InputGroupInput(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-class") {
			t.Error("expected input group to contain custom class")
		}
	})

	t.Run("renders with custom attributes", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{WithGroupAttrs(g.Attr("data-testid", "group"))},
			InputGroupInput(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-testid=\"group\"") {
			t.Error("expected input group to contain custom attribute")
		}
	})

	t.Run("renders with disabled state", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{GroupDisabled()},
			InputGroupInput(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-disabled=\"true\"") {
			t.Error("expected input group to have disabled state")
		}
	})

	t.Run("has focus state classes", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupInput(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "has-[[data-slot=input-group-control]:focus-visible]:border-ring") {
			t.Error("expected input group to have focus state classes")
		}
	})

	t.Run("has error state classes", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupInput(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "has-[[data-slot][aria-invalid=true]]:border-destructive") {
			t.Error("expected input group to have error state classes")
		}
	})

	t.Run("styles direct child inputs via selectors", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupInput(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		// InputGroup should have selectors to style child inputs
		// Note: & gets HTML-encoded as &amp; in the output, ! becomes !important
		if !strings.Contains(output, "&gt;input]:!border-0") {
			t.Error("expected input group to have child input border styling")
		}

		if !strings.Contains(output, "&gt;input]:!shadow-none") {
			t.Error("expected input group to have child input shadow styling")
		}
	})

	t.Run("styles direct child textareas via selectors", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupTextarea(),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		// InputGroup should have selectors to style child textareas
		// Note: & gets HTML-encoded as &amp; in the output, ! becomes !important
		if !strings.Contains(output, "&gt;textarea]:!border-0") {
			t.Error("expected input group to have child textarea border styling")
		}

		if !strings.Contains(output, "&gt;textarea]:resize-none") {
			t.Error("expected input group to have child textarea resize styling")
		}
	})
}

func TestInputGroupAddon(t *testing.T) {
	t.Run("renders addon with default inline-start alignment", func(t *testing.T) {
		addon := InputGroupAddon(
			[]AddonOption{},
			g.Text("https://"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "https://") {
			t.Error("expected addon to contain text")
		}

		if !strings.Contains(output, "data-align=\"inline-start\"") {
			t.Error("expected addon to have default inline-start alignment")
		}

		if !strings.Contains(output, "order-first") {
			t.Error("expected inline-start addon to have order-first class")
		}
	})

	t.Run("renders addon with inline-end alignment", func(t *testing.T) {
		addon := InputGroupAddon(
			[]AddonOption{WithAddonAlign(AlignInlineEnd)},
			g.Text("USD"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-align=\"inline-end\"") {
			t.Error("expected addon to have inline-end alignment")
		}

		if !strings.Contains(output, "order-last") {
			t.Error("expected inline-end addon to have order-last class")
		}
	})

	t.Run("renders addon with block-start alignment", func(t *testing.T) {
		addon := InputGroupAddon(
			[]AddonOption{WithAddonAlign(AlignBlockStart)},
			g.Text("Label"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-align=\"block-start\"") {
			t.Error("expected addon to have block-start alignment")
		}

		if !strings.Contains(output, "w-full") {
			t.Error("expected block-start addon to be full width")
		}
	})

	t.Run("renders addon with block-end alignment", func(t *testing.T) {
		addon := InputGroupAddon(
			[]AddonOption{WithAddonAlign(AlignBlockEnd)},
			g.Text("Helper text"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-align=\"block-end\"") {
			t.Error("expected addon to have block-end alignment")
		}

		if !strings.Contains(output, "order-last") {
			t.Error("expected block-end addon to have order-last class")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		addon := InputGroupAddon(
			[]AddonOption{WithAddonClass("custom-addon")},
			g.Text("@"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-addon") {
			t.Error("expected addon to contain custom class")
		}
	})

	t.Run("has correct data-slot attribute", func(t *testing.T) {
		addon := InputGroupAddon(
			[]AddonOption{},
			g.Text("$"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-slot=\"input-group-addon\"") {
			t.Error("expected addon to have data-slot attribute")
		}
	})
}

func TestInputGroupButton(t *testing.T) {
	t.Run("renders button with default xs size", func(t *testing.T) {
		btn := InputGroupButton(g.Text("Submit"))

		var buf bytes.Buffer
		btn.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Submit") {
			t.Error("expected button to contain text")
		}

		if !strings.Contains(output, "h-6") {
			t.Error("expected button to have xs size (h-6)")
		}
	})

	t.Run("renders button with sm size", func(t *testing.T) {
		btn := InputGroupButton(
			g.Text("Search"),
			WithGroupButtonSize(GroupButtonSizeSM),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "h-8") {
			t.Error("expected button to have sm size (h-8)")
		}
	})

	t.Run("renders icon button with icon-xs size", func(t *testing.T) {
		btn := InputGroupButton(
			g.Text("×"),
			WithGroupButtonSize(GroupButtonSizeIconXS),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "size-6") {
			t.Error("expected button to have icon-xs size (size-6)")
		}
	})

	t.Run("renders icon button with icon-sm size", func(t *testing.T) {
		btn := InputGroupButton(
			g.Text("×"),
			WithGroupButtonSize(GroupButtonSizeIconSM),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "size-8") {
			t.Error("expected button to have icon-sm size (size-8)")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		btn := InputGroupButton(
			g.Text("Click"),
			WithGroupButtonClass("custom-btn"),
		)

		var buf bytes.Buffer
		btn.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-btn") {
			t.Error("expected button to contain custom class")
		}
	})
}

func TestInputGroupText(t *testing.T) {
	t.Run("renders text element", func(t *testing.T) {
		text := InputGroupText("https://")

		var buf bytes.Buffer
		text.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "https://") {
			t.Error("expected text element to contain text")
		}

		if !strings.Contains(output, "text-muted-foreground") {
			t.Error("expected text to have muted foreground color")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		text := InputGroupText("$", WithTextClass("custom-text"))

		var buf bytes.Buffer
		text.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-text") {
			t.Error("expected text to contain custom class")
		}
	})
}

func TestInputGroupInput(t *testing.T) {
	t.Run("renders input with data-slot", func(t *testing.T) {
		input := InputGroupInput()

		var buf bytes.Buffer
		input.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-slot=\"input-group-control\"") {
			t.Error("expected input to have data-slot attribute")
		}
	})

	t.Run("renders input with semantic styling", func(t *testing.T) {
		input := InputGroupInput()

		var buf bytes.Buffer
		input.Render(&buf)
		output := buf.String()

		// Border/shadow styling is handled by parent InputGroup via [&>input] selectors
		if !strings.Contains(output, "placeholder:text-muted-foreground") {
			t.Error("expected input to have placeholder styling")
		}

		if !strings.Contains(output, "text-sm") {
			t.Error("expected input to have text-sm")
		}
	})

	t.Run("renders with placeholder", func(t *testing.T) {
		input := InputGroupInput(WithGroupInputPlaceholder("Enter value"))

		var buf bytes.Buffer
		input.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "placeholder=\"Enter value\"") {
			t.Error("expected input to have placeholder")
		}
	})

	t.Run("renders with name", func(t *testing.T) {
		input := InputGroupInput(WithGroupInputName("email"))

		var buf bytes.Buffer
		input.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "name=\"email\"") {
			t.Error("expected input to have name attribute")
		}
	})

	t.Run("renders with invalid state", func(t *testing.T) {
		input := InputGroupInput(GroupInputInvalid())

		var buf bytes.Buffer
		input.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "aria-invalid=\"true\"") {
			t.Error("expected input to have aria-invalid attribute")
		}
	})

	t.Run("renders with disabled state", func(t *testing.T) {
		input := InputGroupInput(GroupInputDisabled())

		var buf bytes.Buffer
		input.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "disabled") {
			t.Error("expected input to be disabled")
		}
	})

	t.Run("renders with required state", func(t *testing.T) {
		input := InputGroupInput(GroupInputRequired())

		var buf bytes.Buffer
		input.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "required") {
			t.Error("expected input to be required")
		}
	})
}

func TestInputGroupTextarea(t *testing.T) {
	t.Run("renders textarea with data-slot", func(t *testing.T) {
		textarea := InputGroupTextarea()

		var buf bytes.Buffer
		textarea.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-slot=\"input-group-control\"") {
			t.Error("expected textarea to have data-slot attribute")
		}
	})

	t.Run("renders textarea with semantic styling", func(t *testing.T) {
		textarea := InputGroupTextarea()

		var buf bytes.Buffer
		textarea.Render(&buf)
		output := buf.String()

		// Border/shadow/resize styling is handled by parent InputGroup via [&>textarea] selectors
		if !strings.Contains(output, "placeholder:text-muted-foreground") {
			t.Error("expected textarea to have placeholder styling")
		}

		if !strings.Contains(output, "text-sm") {
			t.Error("expected textarea to have text-sm")
		}
	})

	t.Run("renders with placeholder", func(t *testing.T) {
		textarea := InputGroupTextarea(WithGroupTextareaPlaceholder("Enter message"))

		var buf bytes.Buffer
		textarea.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "placeholder=\"Enter message\"") {
			t.Error("expected textarea to have placeholder")
		}
	})

	t.Run("renders with invalid state", func(t *testing.T) {
		textarea := InputGroupTextarea(GroupTextareaInvalid())

		var buf bytes.Buffer
		textarea.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "aria-invalid=\"true\"") {
			t.Error("expected textarea to have aria-invalid attribute")
		}
	})
}

// Legacy compatibility tests

func TestInputLeftAddon(t *testing.T) {
	t.Run("renders as inline-start addon", func(t *testing.T) {
		addon := InputLeftAddon(
			[]AddonOption{},
			g.Text("https://"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "https://") {
			t.Error("expected addon to contain text")
		}

		if !strings.Contains(output, "data-align=\"inline-start\"") {
			t.Error("expected left addon to have inline-start alignment")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		addon := InputLeftAddon(
			[]AddonOption{WithAddonClass("custom-addon")},
			g.Text("@"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-addon") {
			t.Error("expected addon to contain custom class")
		}
	})
}

func TestInputRightAddon(t *testing.T) {
	t.Run("renders as inline-end addon", func(t *testing.T) {
		addon := InputRightAddon(
			[]AddonOption{},
			g.Text("USD"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "USD") {
			t.Error("expected addon to contain text")
		}

		if !strings.Contains(output, "data-align=\"inline-end\"") {
			t.Error("expected right addon to have inline-end alignment")
		}
	})

	t.Run("renders with muted styling", func(t *testing.T) {
		addon := InputRightAddon(
			[]AddonOption{},
			g.Text(".com"),
		)

		var buf bytes.Buffer
		addon.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "text-muted-foreground") {
			t.Error("expected addon to have muted foreground")
		}
	})
}

func TestInputLeftElement(t *testing.T) {
	t.Run("renders left element", func(t *testing.T) {
		element := InputLeftElement(
			[]ElementOption{},
			g.Text("Icon"),
		)

		var buf bytes.Buffer
		element.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "absolute left-0") {
			t.Error("expected element to be positioned absolutely on left")
		}

		if !strings.Contains(output, "pointer-events-none") {
			t.Error("expected element to have pointer-events-none")
		}

		if !strings.Contains(output, "Icon") {
			t.Error("expected element to contain content")
		}
	})

	t.Run("renders with flex centering", func(t *testing.T) {
		element := InputLeftElement(
			[]ElementOption{},
			g.Text("🔍"),
		)

		var buf bytes.Buffer
		element.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "flex items-center justify-center") {
			t.Error("expected element to have flex centering classes")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		element := InputLeftElement(
			[]ElementOption{WithElementClass("custom-element")},
			g.Text("Icon"),
		)

		var buf bytes.Buffer
		element.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-element") {
			t.Error("expected element to contain custom class")
		}
	})
}

func TestInputRightElement(t *testing.T) {
	t.Run("renders right element", func(t *testing.T) {
		element := InputRightElement(
			[]ElementOption{},
			g.Text("✕"),
		)

		var buf bytes.Buffer
		element.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "absolute right-0") {
			t.Error("expected element to be positioned absolutely on right")
		}

		if !strings.Contains(output, "✕") {
			t.Error("expected element to contain content")
		}
	})

	t.Run("renders without pointer-events-none", func(t *testing.T) {
		element := InputRightElement(
			[]ElementOption{},
			g.Text("Button"),
		)

		var buf bytes.Buffer
		element.Render(&buf)
		output := buf.String()

		// Right elements shouldn't have pointer-events-none by default
		// as they often contain interactive elements like clear buttons
		if strings.Contains(output, "pointer-events-none") {
			t.Error("expected right element to NOT have pointer-events-none")
		}
	})
}

func TestInputGroupIntegration(t *testing.T) {
	t.Run("renders group with inline-start addon and input", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupAddon(
				[]AddonOption{WithAddonAlign(AlignInlineStart)},
				InputGroupText("https://"),
			),
			InputGroupInput(WithGroupInputPlaceholder("example.com")),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "https://") {
			t.Error("expected group to contain left addon")
		}

		if !strings.Contains(output, "example.com") {
			t.Error("expected group to contain input")
		}
	})

	t.Run("renders group with input and inline-end addon", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupInput(WithGroupInputType("number")),
			InputGroupAddon(
				[]AddonOption{WithAddonAlign(AlignInlineEnd)},
				InputGroupText("USD"),
			),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "type=\"number\"") {
			t.Error("expected group to contain number input")
		}

		if !strings.Contains(output, "USD") {
			t.Error("expected group to contain right addon")
		}
	})

	t.Run("renders group with both addons", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupAddon(
				[]AddonOption{WithAddonAlign(AlignInlineStart)},
				InputGroupText("$"),
			),
			InputGroupInput(WithGroupInputType("number")),
			InputGroupAddon(
				[]AddonOption{WithAddonAlign(AlignInlineEnd)},
				InputGroupText(".00"),
			),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "$") {
			t.Error("expected group to contain left addon")
		}

		if !strings.Contains(output, ".00") {
			t.Error("expected group to contain right addon")
		}
	})

	t.Run("renders group with button addon", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupInput(WithGroupInputPlaceholder("Search...")),
			InputGroupAddon(
				[]AddonOption{WithAddonAlign(AlignInlineEnd)},
				InputGroupButton(g.Text("Search")),
			),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Search") {
			t.Error("expected group to contain search button")
		}
	})

	t.Run("renders group with textarea", func(t *testing.T) {
		group := InputGroup(
			[]GroupOption{},
			InputGroupTextarea(WithGroupTextareaPlaceholder("Message...")),
			InputGroupAddon(
				[]AddonOption{WithAddonAlign(AlignBlockEnd)},
				InputGroupButton(g.Text("Send")),
			),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Message...") {
			t.Error("expected group to contain textarea placeholder")
		}
	})
}

func TestSearchInput(t *testing.T) {
	t.Run("renders search input", func(t *testing.T) {
		search := SearchInput()

		var buf bytes.Buffer
		search.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "Search...") {
			t.Error("expected search input to have search placeholder")
		}

		if !strings.Contains(output, "svg") {
			t.Error("expected search input to contain SVG icon")
		}
	})
}
