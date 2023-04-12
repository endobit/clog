package clog

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/exp/slog"

	"github.com/endobit/clog/ansi"
)

// Handler implements an slog.Handler.
type Handler struct {
	mutex      sync.Mutex
	opts       HandlerOptions
	colorOpts  ColorOptions
	formatOpts FormatOptions
	writer     io.Writer
	attrs      []slog.Attr
	groups     []string
}

var _ slog.Handler = new(Handler) // Handle implements the slog.Handler interface.

// Enabled implements the slog.Handler interface.
func (h *Handler) Enabled(_ context.Context, l slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}

	return l >= minLevel
}

// Handle implements the slog.Handler interface.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	c := ansi.NewColorer()

	message := new(strings.Builder)

	fmt.Fprint(message, c.Color(r.Time.Format(h.formatOpts.Time), h.colorOpts.Time),
		" ", c.Color(h.formatOpts.levelString(r.Level), h.colorOpts.levelColor(r.Level)),
		" ", r.Message)

	if h.opts.AddSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		fmt.Fprint(message, " ", "[", f.File, ":", f.Line, "]")
	}

	for i := range h.attrs {
		key, val := h.attrFmt(r.Level, h.attrs[i])
		fmt.Fprint(message, " ", key, val)
	}

	r.Attrs(func(attr slog.Attr) {
		key, val := h.attrFmt(r.Level, attr)
		fmt.Fprint(message, " ", key, val)
	})

	h.mutex.Lock()
	defer h.mutex.Unlock()

	fmt.Fprintln(h.writer, message.String())

	return nil
}

// WithAttrs implements the slog.Handler interface.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()
	h2.attrs = append(h2.attrs, attrs...)

	return h2
}

// WithGroup implements the slog.Handler interface.
func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)

	return h2
}
