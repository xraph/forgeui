package alpine

// rawJS is a helper type to represent raw JavaScript code that should not be quoted.
type rawJS struct {
	code string
}

// RawJS creates a raw JavaScript value for use in Alpine x-data.
// This is useful for defining functions and getters in Alpine state.
func RawJS(code string) any {
	return rawJS{code: code}
}

// String returns the raw JavaScript code.
func (r rawJS) String() string {
	return r.code
}
