package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path"
	"runtime"
	"strings"
	"sync"
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
	c := h.colorOpts.Colorer

	msg := new(strings.Builder)

	fmt.Fprint(msg, c.Color(r.Time.Format(h.formatOpts.Time), h.colorOpts.Time),
		" ", c.Color(h.formatOpts.levelString(r.Level), h.colorOpts.levelColor(r.Level)),
		" ", r.Message)

	if h.opts.AddSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		src := fmt.Sprint("[", path.Base(f.File), ":", f.Line, "]")
		fmt.Fprint(msg, " ", c.Color(src, h.colorOpts.Source))
	}

	for i := range h.attrs {
		key, val := h.attrFmt(r.Level, h.attrs[i])
		fmt.Fprint(msg, " ", key, val)
	}

	r.Attrs(func(attr slog.Attr) bool {
		key, val := h.attrFmt(r.Level, attr)
		fmt.Fprint(msg, " ", key, val)

		return true
	})

	h.mutex.Lock()
	defer h.mutex.Unlock()

	fmt.Fprintln(h.writer, msg.String())

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
