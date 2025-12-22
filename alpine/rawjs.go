package alpine

import (
	"io"
)

// rawJS is a helper type to represent raw JavaScript code that should not be quoted.
type rawJS struct {
	code string
}

// RawJS creates a raw JavaScript value for use in Alpine x-data.
// This is useful for defining functions and getters in Alpine state.
func RawJS(code string) any {
	return rawJS{code: code}
}

// Render implements g.Node interface (even though we don't want to render it).
func (r rawJS) Render(w io.Writer) error {
	// This should not be called directly
	_, err := w.Write([]byte(r.code))
	return err
}

// Type returns the type name.
func (r rawJS) Type() string {
	return "alpine.rawJS"
}

