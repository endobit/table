// Package color provides default foreground and background color convenience
// variables. For simple color manipulation it is nice to be able to write:
//
//	fmt.Println(sgr.Wrap(color.Red, "error"))
//
// versus:
//
//	fmt.Println(sgr.Wrap(sgr.Red.FG(), "error"))
package color

import (
	"endobit.io/table/sgr"
)

// Standard ANSI foreground and background colors.
var (
	Black   = fg(sgr.Black)
	Red     = fg(sgr.Red)
	Green   = fg(sgr.Green)
	Yellow  = fg(sgr.Yellow)
	Blue    = fg(sgr.Blue)
	Magenta = fg(sgr.Magenta)
	Cyan    = fg(sgr.Cyan)
	White   = fg(sgr.White)

	BlackBG   = bg(sgr.Black)
	RedBG     = bg(sgr.Red)
	GreenBG   = bg(sgr.Green)
	YellowBG  = bg(sgr.Yellow)
	BlueBG    = bg(sgr.Blue)
	MagentaBG = bg(sgr.Magenta)
	CyanBG    = bg(sgr.Cyan)
	WhiteBG   = bg(sgr.White)
)

func fg(c sgr.Color) []sgr.Param {
	return []sgr.Param{c.FG()}
}

func bg(c sgr.Color) []sgr.Param {
	return []sgr.Param{c.BG()}
}
