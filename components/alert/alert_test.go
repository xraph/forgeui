package alert

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/xraph/forgeui"
)

func TestAlert(t *testing.T) {
	alert := Alert(nil,
		AlertTitle("Alert"),
		AlertDescription("This is an alert"),
	)

	var buf bytes.Buffer
	alert.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role attribute")
	}

	if !strings.Contains(html, "Alert") {
		t.Error("expected alert title")
	}

	if !strings.Contains(html, "This is an alert") {
		t.Error("expected alert description")
	}
}

func TestAlert_Variants(t *testing.T) {
	tests := []struct {
		variant forgeui.Variant
		want    string
	}{
		{forgeui.VariantDefault, "bg-card"},
		{forgeui.VariantDestructive, "text-destructive"},
	}

	for _, tt := range tests {
		t.Run(string(tt.variant), func(t *testing.T) {
			alert := Alert(
				[]Option{WithVariant(tt.variant)},
				g.Text("content"),
			)

			var buf bytes.Buffer
			alert.Render(&buf)

			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("expected %v class", tt.want)
			}
		})
	}
}
