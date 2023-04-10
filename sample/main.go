// Package main provides a sample program for the clog package.
package main

import (
	"clog"
	"net"
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	l := slog.LevelVar{}
	l.Set(slog.DebugLevel)

	opts := clog.HandlerOptions{Level: &l, AddSource: true}

	log := slog.New(opts.NewHandler(os.Stdout))

	log.Debug("hello world", "name", "Al")

	log.Error("oops", net.ErrClosed, "status", 500)

	x := log.WithGroup("my stuff")

	x.LogAttrs(slog.ErrorLevel, "oops",
		slog.Int("status", 500), slog.Any("err", net.ErrClosed))

	y := log.With(slog.Int("foo", 42), slog.String("bar", "foo"))
	y.Info("stuff", "count", "lots of it")
}
