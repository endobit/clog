package clog

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

type Handler struct {
	mutex      sync.Mutex
	opts       HandlerOptions
	colorOpts  ColorOptions
	formatOpts FormatOptions
	writer     io.Writer
	attrs      []slog.Attr
	groups     []string
}

type HandlerOptions slog.HandlerOptions

type FormatOptions struct {
	Time  string
	Level map[slog.Level]string
}

type ColorOptions struct {
	NoColor func() bool
	Time    Color
	Field   Color
	Level   map[slog.Level]Color
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
	NoColor: noColor,
	Time:    Faint,
	Field:   Cyan,
	Level: map[slog.Level]Color{
		slog.DebugLevel: Yellow,
		slog.InfoLevel:  Green,
		slog.WarnLevel:  Red,
		slog.ErrorLevel: BrightRed,
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
func (o HandlerOptions) NewHandler(w io.Writer, opts ...func(*Handler)) *Handler {
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
func NewHandler(w io.Writer) *Handler {
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

	c := colorer{NoColor: h.colorOpts.NoColor}

	key = c.color(key+"=", h.colorOpts.Field)

	if level >= slog.ErrorLevel && attr.Key == "err" {
		val = c.color(val, h.colorOpts.levelColor(level))
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

func (c ColorOptions) levelColor(l slog.Level) Color {
	if color, ok := c.Level[l]; ok {
		return color // exact match
	}

	var color Color

	switch {
	case l < slog.InfoLevel:
		color = c.Level[slog.DebugLevel]
	case l < slog.WarnLevel:
		color = c.Level[slog.InfoLevel]
	case l < slog.ErrorLevel:
		color = c.Level[slog.WarnLevel]
	default:
		color = c.Level[slog.ErrorLevel]
	}

	return color
}
