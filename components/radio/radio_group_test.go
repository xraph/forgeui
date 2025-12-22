package radio

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestGroup(t *testing.T) {
	t.Run("renders basic group", func(t *testing.T) {
		group := Group(
			[]GroupOption{WithGroupName("size")},
			Radio(WithName("size"), WithValue("small")),
			Radio(WithName("size"), WithValue("medium")),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "role=\"radiogroup\"") {
			t.Error("expected group to have radiogroup role")
		}
		if !strings.Contains(output, "space-y-2") {
			t.Error("expected group to have spacing classes")
		}
	})

	t.Run("renders group with name", func(t *testing.T) {
		group := Group(
			[]GroupOption{WithGroupName("size")},
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-name=\"size\"") {
			t.Error("expected group to have data-name attribute")
		}
	})

	t.Run("renders group with value", func(t *testing.T) {
		group := Group(
			[]GroupOption{
				WithGroupName("size"),
				WithGroupValue("medium"),
			},
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-value=\"medium\"") {
			t.Error("expected group to have data-value attribute")
		}
	})

	t.Run("renders disabled group", func(t *testing.T) {
		group := Group(
			[]GroupOption{
				WithGroupName("size"),
				WithGroupDisabled(),
			},
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "aria-disabled=\"true\"") {
			t.Error("expected disabled group to have aria-disabled attribute")
		}
	})

	t.Run("renders group with custom class", func(t *testing.T) {
		group := Group(
			[]GroupOption{
				WithGroupName("size"),
				WithGroupClass("custom-group"),
			},
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "custom-group") {
			t.Error("expected group to contain custom class")
		}
	})

	t.Run("renders group with custom attributes", func(t *testing.T) {
		group := Group(
			[]GroupOption{
				WithGroupName("size"),
				WithGroupAttrs(g.Attr("data-testid", "size-group")),
			},
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "data-testid=\"size-group\"") {
			t.Error("expected group to contain custom attribute")
		}
	})
}

func TestGroupItem(t *testing.T) {
	t.Run("renders group item", func(t *testing.T) {
		item := GroupItem(
			"size-small",
			"Small",
			WithValue("small"),
			WithName("size"),
		)

		var buf bytes.Buffer
		item.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "flex") {
			t.Error("expected item to have flex layout")
		}
		if !strings.Contains(output, "items-center") {
			t.Error("expected item to have items-center class")
		}
		if !strings.Contains(output, "space-x-2") {
			t.Error("expected item to have spacing")
		}
	})

	t.Run("renders item with radio button", func(t *testing.T) {
		item := GroupItem(
			"size-small",
			"Small",
			WithValue("small"),
			WithName("size"),
		)

		var buf bytes.Buffer
		item.Render(&buf)
		output := buf.String()

		if !strings.Contains(output, "type=\"radio\"") {
			t.Error("expected item to contain radio input")
		}
		if !strings.Contains(output, "value=\"small\"") {
			t.Error("expected item to contain value")
		}
	})
}

func TestGroupIntegration(t *testing.T) {
	t.Run("renders complete radio group", func(t *testing.T) {
		group := Group(
			[]GroupOption{
				WithGroupName("size"),
				WithGroupValue("medium"),
			},
			GroupItem("size-small", "Small", WithValue("small"), WithName("size")),
			GroupItem("size-medium", "Medium", WithValue("medium"), WithName("size"), Checked()),
			GroupItem("size-large", "Large", WithValue("large"), WithName("size")),
		)

		var buf bytes.Buffer
		group.Render(&buf)
		output := buf.String()

		// Check group structure
		if !strings.Contains(output, "role=\"radiogroup\"") {
			t.Error("expected radiogroup role")
		}

		// Check all items are present
		if !strings.Contains(output, "value=\"small\"") {
			t.Error("expected small option")
		}
		if !strings.Contains(output, "value=\"medium\"") {
			t.Error("expected medium option")
		}
		if !strings.Contains(output, "value=\"large\"") {
			t.Error("expected large option")
		}

		// Check checked state
		if !strings.Contains(output, "checked") {
			t.Error("expected checked attribute on medium option")
		}
	})
}

