// Package primitives provides low-level layout primitives for ForgeUI.
//
// Primitives are the building blocks for more complex components.
// They provide type-safe wrappers around common CSS layout patterns.
// All primitives return templ.Component and accept children as templ.Component.
//
// Layout primitives:
//   - Box: Polymorphic container element
//   - Flex: Flexbox container
//   - Grid: CSS Grid container
//
// Stack helpers:
//   - VStack: Vertical stack (flex column)
//   - HStack: Horizontal stack (flex row)
//
// Utility primitives:
//   - Center: Centers content both horizontally and vertically
//   - Container: Responsive container with max-width
//   - Spacer: Flexible spacer that fills available space
//   - Text: Typography primitive
package primitives

import (
	"context"
	"fmt"
	stdhtml "html"
	"io"

	"github.com/a-h/templ"
)

// writeOpenTag writes an HTML opening tag with optional class and attributes.
func writeOpenTag(w io.Writer, tag, classes string, attrs templ.Attributes) error {
	if _, err := fmt.Fprintf(w, "<%s", tag); err != nil {
		return err
	}
	if classes != "" {
		if _, err := fmt.Fprintf(w, ` class="%s"`, stdhtml.EscapeString(classes)); err != nil {
			return err
		}
	}
	if err := writeAttrs(w, attrs); err != nil {
		return err
	}
	_, err := io.WriteString(w, ">")
	return err
}

// writeCloseTag writes an HTML closing tag.
func writeCloseTag(w io.Writer, tag string) error {
	_, err := fmt.Fprintf(w, "</%s>", tag)
	return err
}

// writeAttrs writes templ.Attributes as HTML attributes.
func writeAttrs(w io.Writer, attrs templ.Attributes) error {
	for k, v := range attrs {
		switch val := v.(type) {
		case string:
			if val == "" {
				if _, err := fmt.Fprintf(w, " %s", k); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(w, ` %s="%s"`, k, stdhtml.EscapeString(val)); err != nil {
					return err
				}
			}
		case bool:
			if val {
				if _, err := fmt.Fprintf(w, " %s", k); err != nil {
					return err
				}
			}
		default:
			if _, err := fmt.Fprintf(w, ` %s="%s"`, k, stdhtml.EscapeString(fmt.Sprint(val))); err != nil {
				return err
			}
		}
	}
	return nil
}

// renderChildren renders a slice of templ.Component children.
func renderChildren(ctx context.Context, w io.Writer, children []templ.Component) error {
	for _, child := range children {
		if child == nil {
			continue
		}
		if err := child.Render(ctx, w); err != nil {
			return err
		}
	}
	return nil
}
