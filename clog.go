// Package clog provides a slog Handler that mimics the output of the zerolog.Logger.
package clog

import (
	"clog/ansi"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// HandlerOptions is a set of options for a Handler.
type HandlerOptions slog.HandlerOptions

// FormatOptions is a set of options for formatting log messages.
type FormatOptions struct {
	Time  string
	Level map[slog.Level]string
}

// ColorOptions is a set of options for colorizing the output of a Handler.
type ColorOptions struct {
	Time  ansi.Color
	Field ansi.Color
	Level map[slog.Level]ansi.Color
}

var defaultFormatOptions = FormatOptions{
	Time: time.Kitchen,
	Level: map[slog.Level]string{
		slog.DebugLevel: "DBG",
		slog.InfoLevel:  "INF",
		slog.WarnLevel:  "WRN",
		slog.ErrorLevel: "ERR",
	},
}

var defaultColorOptions = ColorOptions{
	Time:  ansi.Faint,
	Field: ansi.Cyan,
	Level: map[slog.Level]ansi.Color{
		slog.DebugLevel: ansi.Yellow,
		slog.InfoLevel:  ansi.Green,
		slog.WarnLevel:  ansi.Red,
		slog.ErrorLevel: ansi.BrightRed,
	},
}

// WithColor is an option setting function for NewHandler. It sets the
// ColorOptions for the Handler.
func WithColor(c ColorOptions) func(*Handler) {
	return func(h *Handler) {
		h.colorOpts = c
	}
}

// WithFormat is an option setting function for NewHandler. It sets the
// FormatOptions for the Handler.
func WithFormat(f FormatOptions) func(*Handler) {
	return func(h *Handler) {
		h.formatOpts = f
	}
}

// NewHandler returns a Handler the writes to w and invokes any option setting
// functions.
func (o HandlerOptions) NewHandler(w io.Writer, opts ...func(*Handler)) slog.Handler {
	h := Handler{
		opts:       o,
		colorOpts:  defaultColorOptions,
		formatOpts: defaultFormatOptions,
		writer:     w,
	}

	for _, opt := range opts {
		opt(&h)
	}

	return &h
}

// NewHandler returns a Handler the writes to w and invokes any option setting
// functions.
func NewHandler(w io.Writer) slog.Handler {
	return (HandlerOptions{}).NewHandler(w)
}

func (h *Handler) clone() *Handler {
	return &Handler{
		opts:       h.opts,
		colorOpts:  h.colorOpts,
		formatOpts: h.formatOpts,
		writer:     h.writer,
		attrs:      h.attrs,
		groups:     h.groups,
	}
}

func (h *Handler) attrFmt(level slog.Level, attr slog.Attr) (key, val string) {
	var prefix string

	if len(h.groups) > 0 {
		prefix = strings.Join(h.groups, ".") + "."
	}

	key = prefix + attr.Key
	if needsQuoting(key) {
		key = strconv.Quote(key)
	}

	val = attr.Value.String()
	if needsQuoting(val) {
		val = strconv.Quote(val)
	}

	c := ansi.NewColorer()

	key = c.Color(key+"=", h.colorOpts.Field)

	if level >= slog.ErrorLevel && attr.Key == "err" {
		val = c.Color(val, h.colorOpts.levelColor(level))
	}

	return key, val
}

func (f FormatOptions) levelString(l slog.Level) string {
	if level, ok := f.Level[l]; ok {
		return level // exact match
	}

	str := func(base, offset slog.Level) string {
		level, ok := f.Level[base]
		if !ok {
			level = base.String()
		}

		if offset == 0 {
			return level
		}

		return fmt.Sprintf("%s%+d", level, offset)
	}

	switch {
	case l < slog.InfoLevel:
		return str(slog.DebugLevel, l-slog.DebugLevel)
	case l < slog.WarnLevel:
		return str(slog.InfoLevel, l)
	case l < slog.ErrorLevel:
		return str(slog.WarnLevel, l-slog.WarnLevel)
	default:
		return str(slog.ErrorLevel, l-slog.ErrorLevel)
	}
}

func (c ColorOptions) levelColor(l slog.Level) ansi.Color {
	if color, ok := c.Level[l]; ok {
		return color // exact match
	}

	switch {
	case l < slog.InfoLevel:
		return c.Level[slog.DebugLevel]
	case l < slog.WarnLevel:
		return c.Level[slog.InfoLevel]
	case l < slog.ErrorLevel:
		return c.Level[slog.WarnLevel]
	default:
		return c.Level[slog.ErrorLevel]
	}
}
