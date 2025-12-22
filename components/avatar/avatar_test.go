package avatar

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xraph/forgeui"
)

func TestAvatar_WithImage(t *testing.T) {
	avatar := Avatar(
		WithSrc("/avatar.jpg"),
		WithAlt("User"),
	)

	var buf bytes.Buffer
	avatar.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "rounded-full") {
		t.Error("expected rounded-full class")
	}

	if !strings.Contains(html, "/avatar.jpg") {
		t.Error("expected image src")
	}
}

func TestAvatar_WithFallback(t *testing.T) {
	avatar := Avatar(WithFallback("AB"))

	var buf bytes.Buffer
	avatar.Render(&buf)
	html := buf.String()

	if !strings.Contains(html, "AB") {
		t.Error("expected fallback text")
	}

	if !strings.Contains(html, "bg-muted") {
		t.Error("expected bg-muted for fallback")
	}
}

func TestAvatar_Sizes(t *testing.T) {
	tests := []struct {
		size forgeui.Size
		want string
	}{
		{forgeui.SizeSM, "size-8"},
		{forgeui.SizeMD, "size-10"},
		{forgeui.SizeLG, "size-12"},
		{forgeui.SizeXL, "size-16"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			avatar := Avatar(WithSize(tt.size), WithFallback("A"))

			var buf bytes.Buffer
			avatar.Render(&buf)

			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("expected %v size classes", tt.want)
			}
		})
	}
}
