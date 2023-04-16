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
	options struct {
		colorable func() bool
	}

	// Colorer is an interface for colorizing text.
	Colorer interface {
		Color(interface{}, Color) string
	}

	colorWrapper struct{}
	nopWrapper   struct{}
)

// Colorable is an option setting function for NewColorer. It replaces the
// default function that determines if color is enabled.
func Colorable(fn func() bool) func(*options) {
	return func(o *options) {
		o.colorable = fn
	}
}

// NewColorer returns a Colorer based on the Colorable function.
func NewColorer(opts ...func(*options)) Colorer {
	o := options{
		colorable: colorable,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.colorable() {
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
