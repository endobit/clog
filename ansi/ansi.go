// Package ansi provides ANSI escape codes for setting the foreground color of
// text.
package ansi

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

type (
	// Colorer is an interface for colorizing text.
	Colorer interface {
		Color(interface{}, Color) string
	}

	colorWrapper struct{}
	nopWrapper   struct{}
)

// Colorable is a function that returns true if the Colorer should apply ANSI
// color codes.
var Colorable = colorable

// NewColorer returns a Colorer based on the Colorable function.
func NewColorer() Colorer {
	if Colorable() {
		return colorWrapper{}
	}

	return nopWrapper{}
}

// Color implements the Colorer interface for c.
func (c colorWrapper) Color(s interface{}, color Color) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, s)
}

// Color implements the Colorer interface for n. It always returns
// the raw string with no color applied.
func (n nopWrapper) Color(s interface{}, _ Color) string {
	return fmt.Sprintf("%v", s)
}

// colorable returns true if the [NO_COLOR] environment variable is not set.
//
// [NO_COLOR]: https://no-color.org
func colorable() bool {
	_, found := os.LookupEnv("NO_COLOR")
	return !found
}
