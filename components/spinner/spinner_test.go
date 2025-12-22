package spinner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui"
)

func TestSpinner(t *testing.T) {
	spinner := Spinner()

	var buf bytes.Buffer
	spinner.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "animate-spin") {
		t.Error("expected animate-spin class")
	}

	if !strings.Contains(html, `role="status"`) {
		t.Error("expected role attribute")
	}
}

func TestSpinner_Sizes(t *testing.T) {
	tests := []struct {
		size forgeui.Size
		want string
	}{
		{forgeui.SizeSM, "size-4"},
		{forgeui.SizeMD, "size-6"},
		{forgeui.SizeLG, "size-8"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			spinner := Spinner(WithSize(tt.size))

			var buf bytes.Buffer
			spinner.Render(&buf)

			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("expected %v size classes", tt.want)
			}
		})
	}
}
