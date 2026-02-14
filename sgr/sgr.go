// Package sgr is a minimalist package for setting ANSI color terminal escape
// sequences. It implements 8 colors and several styles of the Select Graphic
// Rendition [SGR].
//
// [SGR]: https://en.wikipedia.org/wiki/ANSI_escape_code#SGR
package sgr

import (
	"os"
	"strconv"
	"strings"
)

// Escape is the leading control sequence for SGR commands.
const Escape = "\x1b["

// Param is a formatted SGR parameter.
type Param int

// Style SGR parameters.
const (
	Reset Param = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Color is base terminal color.
type Color int

// Base values for ANSI terminal colors. Use the `FG` and `BG` methods to obtain
// a valid Param.
const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Default
)

// FG returns the foreground color SGR parameter for c.
func (c Color) FG() Param {
	return Param(30 + c)
}

// BG returns the background color SGR parameter for c.
func (c Color) BG() Param {
	return Param(40 + c)
}

// code returns the SGR escape sequence for the combined Params.
func code(p []Param) string {
	if len(p) == 0 || noColor {
		return ""
	}

	codes := make([]string, len(p))
	for i := range p {
		codes[i] = strconv.Itoa(int(p[i]))
	}

	return Escape + strings.Join(codes, ";") + "m"
}

var (
	noColor bool // https://no-color.org/

	reset       = []Param{Reset}
	resetString = code(reset)
)

func init() {
	if s := os.Getenv("NO_COLOR"); s != "" {
		noColor = true
	}
}

// DisableColor disables all SGR color output. This is useful for testing.
func DisableColor() {
	noColor = true
	resetString = ""
}
