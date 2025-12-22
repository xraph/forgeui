package badge

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui"
)

func TestBadge(t *testing.T) {
	badge := Badge("New")

	var buf bytes.Buffer
	badge.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "inline-flex") {
		t.Error("expected inline-flex class")
	}

	if !strings.Contains(html, "New") {
		t.Error("expected badge text")
	}
}

func TestBadge_Variants(t *testing.T) {
	tests := []struct {
		variant forgeui.Variant
		want    string
	}{
		{forgeui.VariantDefault, "bg-primary"},
		{forgeui.VariantSecondary, "bg-secondary"},
		{forgeui.VariantDestructive, "bg-destructive"},
		{forgeui.VariantOutline, "text-foreground"},
	}

	for _, tt := range tests {
		t.Run(string(tt.variant), func(t *testing.T) {
			badge := Badge("Test", WithVariant(tt.variant))

			var buf bytes.Buffer
			badge.Render(&buf)

			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("expected %v class", tt.want)
			}
		})
	}
}
