package clog

import (
	"fmt"
	"os"
)

// Color is a ANSI escape sequence color number.
type Color int

// Colors can be bolded or dimmed.
const (
	Bold Color = iota + 1
	Faint
)

// Standard colors.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Bright colors.
const (
	BrightRed Color = iota + 91
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

type colorer struct {
	NoColor func() bool
}

// colorreturns the string representation of s wrapped in color, or unwrapped
// if the NoColor function returns true.
func (c colorer) color(s interface{}, color Color) string {
	if c.NoColor != nil && c.NoColor() {
		return fmt.Sprint(s)
	}

	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, s)
}

// noColor returns true if the [NO_COLOR] environment variable is set.
//
// [NO_COLOR]: https://no-color.org
func noColor() bool {
	_, found := os.LookupEnv("NO_COLOR")
	return found
}
