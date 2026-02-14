package sgr

import "fmt"

// Wrapped represents SGR wrapped text.
type Wrapped struct {
	Text string // Text is the original text, not the colored value.
	head string
	tail string
}

// String implements the Stringer interface for w. This will be the colored
// text, to access the uncolored value use the `Text` field.
func (w Wrapped) String() string {
	return w.head + w.Text + w.tail
}

// Wrap applies the SGR parameters to wrap the formatted text.
func Wrap(p []Param, a ...any) Wrapped {
	return wrap(p, fmt.Sprint(a...))
}

// Wrapf applies the SGR parameters to wrap the formatted text.
func Wrapf(p []Param, format string, a ...any) Wrapped {
	return wrap(p, fmt.Sprintf(format, a...))
}

func wrap(p []Param, s string) Wrapped {
	return Wrapped{
		Text: s,
		head: code(p...),
		tail: resetString,
	}
}
