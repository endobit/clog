// Package main provides a sample program for the clog package.
package main

import (
	"errors"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"github.com/endobit/clog"
)

var errReset = errors.New("connection reset by peer")

func main() {
	level := slog.LevelVar{}
	level.Set(slog.LevelDebug)

	opts := clog.HandlerOptions{Level: &level, AddSource: false}

	log := slog.New(opts.NewHandler(os.Stdout))

	// This is directly from the zerolog.ConsoleWriter README

	l := log.With(slog.Int("pid", 37556))
	l.Info("Starting listener", "listen", ":8080")
	l.Debug("Access", "database", "myapp", "host", "localhost:4932")

	l.Info("Access", "method", "GET", "path", "/users", "resp_time", 23*time.Millisecond)

	{
		l := l.With("method", "POST", "path", "/posts", "resp_time", 532*time.Millisecond)

		l.Info("Access")
		l.Warn("Slow request")
	}

	l.Info("Access", "method", "GET", "path", "/users", "resp_time", 10*time.Millisecond)

	l.Error("Database connection lost", clog.ErrorFieldName, errReset, "database", "myapp")
}
