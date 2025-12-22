package primitives

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestBox(t *testing.T) {
	t.Run("renders default div", func(t *testing.T) {
		box := Box(WithChildren(g.Text("content")))

		var buf bytes.Buffer
		if err := box.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "<div") {
			t.Error("expected <div tag")
		}

		if !strings.Contains(html, "content") {
			t.Error("expected content")
		}
	})

	t.Run("renders custom element", func(t *testing.T) {
		box := Box(
			WithAs("section"),
			WithChildren(g.Text("content")),
		)

		var buf bytes.Buffer
		box.Render(&buf)

		if !strings.Contains(buf.String(), "<section") {
			t.Error("expected <section tag")
		}
	})

	t.Run("applies utility classes", func(t *testing.T) {
		box := Box(
			WithPadding("p-4"),
			WithMargin("m-2"),
			WithBackground("bg-blue-500"),
			WithRounded("rounded-lg"),
			WithShadow("shadow-md"),
		)

		var buf bytes.Buffer
		box.Render(&buf)
		html := buf.String()

		classes := []string{"p-4", "m-2", "bg-blue-500", "rounded-lg", "shadow-md"}
		for _, class := range classes {
			if !strings.Contains(html, class) {
				t.Errorf("expected class %v", class)
			}
		}
	})
}

func TestFlex(t *testing.T) {
	t.Run("renders flex container", func(t *testing.T) {
		flex := Flex(FlexChildren(g.Text("content")))

		var buf bytes.Buffer
		flex.Render(&buf)

		html := buf.String()
		if !strings.Contains(html, "flex") {
			t.Error("expected flex class")
		}
	})

	t.Run("applies direction classes", func(t *testing.T) {
		tests := []struct {
			direction string
			want      string
		}{
			{"col", "flex-col"},
			{"row-reverse", "flex-row-reverse"},
			{"col-reverse", "flex-col-reverse"},
		}

		for _, tt := range tests {
			t.Run(tt.direction, func(t *testing.T) {
				flex := Flex(FlexDirection(tt.direction))

				var buf bytes.Buffer
				flex.Render(&buf)

				if !strings.Contains(buf.String(), tt.want) {
					t.Errorf("expected %v class", tt.want)
				}
			})
		}
	})

	t.Run("applies justify classes", func(t *testing.T) {
		flex := Flex(FlexJustify("center"))

		var buf bytes.Buffer
		flex.Render(&buf)

		if !strings.Contains(buf.String(), "justify-center") {
			t.Error("expected justify-center class")
		}
	})

	t.Run("applies gap", func(t *testing.T) {
		flex := Flex(FlexGap("4"))

		var buf bytes.Buffer
		flex.Render(&buf)

		if !strings.Contains(buf.String(), "gap-4") {
			t.Error("expected gap-4 class")
		}
	})
}

func TestGrid(t *testing.T) {
	t.Run("renders grid container", func(t *testing.T) {
		grid := Grid(GridChildren(g.Text("content")))

		var buf bytes.Buffer
		grid.Render(&buf)

		if !strings.Contains(buf.String(), "grid") {
			t.Error("expected grid class")
		}
	})

	t.Run("applies column count", func(t *testing.T) {
		grid := Grid(GridCols(3))

		var buf bytes.Buffer
		grid.Render(&buf)

		if !strings.Contains(buf.String(), "grid-cols-3") {
			t.Error("expected grid-cols-3 class")
		}
	})

	t.Run("applies responsive columns", func(t *testing.T) {
		grid := Grid(
			GridCols(1),
			GridColsSM(2),
			GridColsMD(3),
			GridColsLG(4),
		)

		var buf bytes.Buffer
		grid.Render(&buf)
		html := buf.String()

		classes := []string{"grid-cols-1", "sm:grid-cols-2", "md:grid-cols-3", "lg:grid-cols-4"}
		for _, class := range classes {
			if !strings.Contains(html, class) {
				t.Errorf("expected %v class", class)
			}
		}
	})
}

func TestVStack(t *testing.T) {
	stack := VStack("4", g.Text("item1"), g.Text("item2"))

	var buf bytes.Buffer
	stack.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "flex-col") {
		t.Error("expected flex-col class")
	}

	if !strings.Contains(html, "gap-4") {
		t.Error("expected gap-4 class")
	}

	if !strings.Contains(html, "item1") || !strings.Contains(html, "item2") {
		t.Error("expected children content")
	}
}

func TestHStack(t *testing.T) {
	stack := HStack("2", g.Text("item1"), g.Text("item2"))

	var buf bytes.Buffer
	stack.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "flex") {
		t.Error("expected flex class")
	}

	if !strings.Contains(html, "items-center") {
		t.Error("expected items-center class")
	}

	if !strings.Contains(html, "gap-2") {
		t.Error("expected gap-2 class")
	}
}

func TestCenter(t *testing.T) {
	center := Center(g.Text("centered content"))

	var buf bytes.Buffer
	center.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "justify-center") {
		t.Error("expected justify-center class")
	}

	if !strings.Contains(html, "items-center") {
		t.Error("expected items-center class")
	}

	if !strings.Contains(html, "centered content") {
		t.Error("expected content")
	}
}

func TestSpacer(t *testing.T) {
	spacer := Spacer()

	var buf bytes.Buffer
	spacer.Render(&buf)

	if !strings.Contains(buf.String(), "flex-1") {
		t.Error("expected flex-1 class")
	}
}

func TestContainer(t *testing.T) {
	container := Container(g.Text("content"))

	var buf bytes.Buffer
	container.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "container") {
		t.Error("expected container class")
	}

	if !strings.Contains(html, "mx-auto") {
		t.Error("expected mx-auto class")
	}

	if !strings.Contains(html, "content") {
		t.Error("expected content")
	}
}

func TestText(t *testing.T) {
	t.Run("renders default paragraph", func(t *testing.T) {
		text := Text(TextChildren(g.Text("Hello")))

		var buf bytes.Buffer
		text.Render(&buf)
		html := buf.String()

		if !strings.Contains(html, "<p") {
			t.Error("expected <p tag")
		}

		if !strings.Contains(html, "Hello") {
			t.Error("expected content")
		}
	})

	t.Run("renders custom element", func(t *testing.T) {
		text := Text(
			TextAs("h1"),
			TextChildren(g.Text("Title")),
		)

		var buf bytes.Buffer
		text.Render(&buf)

		if !strings.Contains(buf.String(), "<h1") {
			t.Error("expected <h1 tag")
		}
	})

	t.Run("applies typography classes", func(t *testing.T) {
		text := Text(
			TextSize("text-lg"),
			TextWeight("font-bold"),
			TextColor("text-blue-500"),
			TextAlign("text-center"),
		)

		var buf bytes.Buffer
		text.Render(&buf)
		html := buf.String()

		classes := []string{"text-lg", "font-bold", "text-blue-500", "text-center"}
		for _, class := range classes {
			if !strings.Contains(html, class) {
				t.Errorf("expected %v class", class)
			}
		}
	})
}
