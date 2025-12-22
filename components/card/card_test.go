package card

import (
	"bytes"
	"strings"
	"testing"

	g "github.com/maragudk/gomponents"
)

func TestCard(t *testing.T) {
	t.Run("renders basic card", func(t *testing.T) {
		card := Card(g.Text("content"))

		var buf bytes.Buffer
		if err := card.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "rounded-lg") {
			t.Error("expected rounded-lg class")
		}
		if !strings.Contains(html, "border") {
			t.Error("expected border class")
		}
		if !strings.Contains(html, "content") {
			t.Error("expected content")
		}
	})

	t.Run("renders with custom class", func(t *testing.T) {
		card := CardWithOptions(
			[]Option{WithClass("custom-card")},
			g.Text("content"),
		)

		var buf bytes.Buffer
		card.Render(&buf)

		if !strings.Contains(buf.String(), "custom-card") {
			t.Error("expected custom-card class")
		}
	})
}

func TestCard_Header(t *testing.T) {
	header := Header(g.Text("header content"))

	var buf bytes.Buffer
	header.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "grid") {
		t.Error("expected grid class")
	}
	if !strings.Contains(html, "px-6") {
		t.Error("expected padding class")
	}
	if !strings.Contains(html, "header content") {
		t.Error("expected header content")
	}
}

func TestCard_Title(t *testing.T) {
	title := Title("My Title")

	var buf bytes.Buffer
	title.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "<h3") {
		t.Error("expected <h3 tag")
	}
	if !strings.Contains(html, "leading-none") {
		t.Error("expected leading-none class")
	}
	if !strings.Contains(html, "font-semibold") {
		t.Error("expected font-semibold class")
	}
	if !strings.Contains(html, "My Title") {
		t.Error("expected title text")
	}
}

func TestCard_Description(t *testing.T) {
	desc := Description("Card description")

	var buf bytes.Buffer
	desc.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "<p") {
		t.Error("expected <p tag")
	}
	if !strings.Contains(html, "text-sm") {
		t.Error("expected text-sm class")
	}
	if !strings.Contains(html, "text-muted-foreground") {
		t.Error("expected text-muted-foreground class")
	}
	if !strings.Contains(html, "Card description") {
		t.Error("expected description text")
	}
}

func TestCard_Content(t *testing.T) {
	content := Content(g.Text("main content"))

	var buf bytes.Buffer
	content.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "px-6") {
		t.Error("expected px-6 class")
	}
	if !strings.Contains(html, "main content") {
		t.Error("expected content text")
	}
}

func TestCard_Footer(t *testing.T) {
	footer := Footer(g.Text("footer content"))

	var buf bytes.Buffer
	footer.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "flex") {
		t.Error("expected flex class")
	}
	if !strings.Contains(html, "items-center") {
		t.Error("expected items-center class")
	}
	if !strings.Contains(html, "footer content") {
		t.Error("expected footer content")
	}
}

func TestCard_Complete(t *testing.T) {
	// Test a complete card with all sections
	card := Card(
		Header(
			Title("Card Title"),
			Description("Card description text"),
		),
		Content(
			g.Text("Main card content"),
		),
		Footer(
			g.Text("Footer actions"),
		),
	)

	var buf bytes.Buffer
	if err := card.Render(&buf); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()

	expected := []string{
		"rounded-lg",
		"border",
		"Card Title",
		"Card description text",
		"Main card content",
		"Footer actions",
	}

	for _, exp := range expected {
		if !strings.Contains(html, exp) {
			t.Errorf("expected %v in complete card", exp)
		}
	}
}

func TestCard_TitleWithOptions(t *testing.T) {
	title := Title("Title", WithClass("custom-title"))

	var buf bytes.Buffer
	title.Render(&buf)

	if !strings.Contains(buf.String(), "custom-title") {
		t.Error("expected custom-title class")
	}
}

func TestCard_DescriptionWithOptions(t *testing.T) {
	desc := Description("Description", WithClass("custom-desc"))

	var buf bytes.Buffer
	desc.Render(&buf)

	if !strings.Contains(buf.String(), "custom-desc") {
		t.Error("expected custom-desc class")
	}
}
