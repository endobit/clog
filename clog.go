// Package clog provides a slog Handler that mimics the output of the zerolog.Logger.
package clog

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	"github.com/endobit/clog/ansi"
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
	Colorer ansi.Colorer
	Time    ansi.Color
	Field   ansi.Color
	Level   map[slog.Level]ansi.Color
}

var defaultFormatOptions = FormatOptions{
	Time: time.Kitchen,
	Level: map[slog.Level]string{
		slog.LevelDebug: "DBG",
		slog.LevelInfo:  "INF",
		slog.LevelWarn:  "WRN",
		slog.LevelError: "ERR",
	},
}

var defaultColorOptions = ColorOptions{
	Colorer: ansi.NewColorer(),
	Time:    ansi.Faint,
	Field:   ansi.Cyan,
	Level: map[slog.Level]ansi.Color{
		slog.LevelDebug: ansi.Yellow,
		slog.LevelInfo:  ansi.Green,
		slog.LevelWarn:  ansi.Red,
		slog.LevelError: ansi.BrightRed,
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

// NewHandler returns a Handler with the default options that writes to w.
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

	c := h.colorOpts.Colorer

	key = c.Color(key+"=", h.colorOpts.Field)

	if level >= slog.LevelError && attr.Key == "err" {
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
	case l < slog.LevelInfo:
		return str(slog.LevelDebug, l-slog.LevelDebug)
	case l < slog.LevelWarn:
		return str(slog.LevelInfo, l)
	case l < slog.LevelError:
		return str(slog.LevelWarn, l-slog.LevelWarn)
	default:
		return str(slog.LevelError, l-slog.LevelError)
	}
}

func (c ColorOptions) levelColor(l slog.Level) ansi.Color {
	if color, ok := c.Level[l]; ok {
		return color // exact match
	}

	switch {
	case l < slog.LevelInfo:
		return c.Level[slog.LevelDebug]
	case l < slog.LevelWarn:
		return c.Level[slog.LevelInfo]
	case l < slog.LevelError:
		return c.Level[slog.LevelWarn]
	default:
		return c.Level[slog.LevelError]
	}
}
