package clog_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/exp/slog"

	"github.com/endobit/clog"
	"github.com/endobit/clog/ansi"
)

func TestDefaults(t *testing.T) {
	buf := new(bytes.Buffer)
	log := slog.New(clog.NewHandler(buf))
	log.Info("foo")

	fields := strings.Fields(buf.String())
	if len(fields) != 3 {
		t.Fatalf("got %d fields, want 3", len(fields))
	}

	tm := stripANSI(fields[0])
	_, err := time.Parse(time.Kitchen, tm)
	if err != nil {
		t.Fatalf("wrong format for time %q", tm)
	}

	level := stripANSI(fields[1])
	if level != "INF" {
		t.Fatalf("got %q, want INF", level)
	}

	if fields[2] != "foo" {
		t.Fatalf("got %q, want foo", fields[2])
	}
}

func TestOptions(t *testing.T) {
	buf := new(bytes.Buffer)

	opts := clog.HandlerOptions{
		AddSource: true,
	}

	log := slog.New(opts.NewHandler(buf, clog.WithColor(nocolor())))

	log.Info("info", "a", 1, "b b", "2 2")

	fields := strings.Fields(buf.String())
	if len(fields) < 4 {
		t.Fatalf("got %d fields, want at least 4", len(fields))
	}

	// Already tested time and level in TestDefaults.

	if fields[2] != "info" {
		t.Fatalf("got %q, want info", fields[2])
	}

	file, _, err := parseSource(fields[3])
	if err != nil {
		t.Fatalf("got %q, want source", fields[3])
	}

	if file != "clog_test.go" {
		t.Fatalf("got %q, want clog_test.go", file)
	}

	attrs := strings.Split(strings.Join(fields[4:], " "), "=")

	want := []string{"a", "1 \"b b\"", "\"2 2\""}
	for i, w := range want {
		if attrs[i] != w {
			t.Fatalf("got %q, want %q", attrs[i], w)
		}
	}
}

func parseSource(s string) (file string, line int, err error) {
	source := strings.ReplaceAll(s, ":", " ")
	source = strings.ReplaceAll(source, "[", "")
	source = strings.ReplaceAll(source, "]", "")
	fields := strings.Fields(source)

	if len(fields) != 2 {
		return "", 0, fmt.Errorf("wrong number of fields")
	}

	line, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, err
	}

	file = fields[0]
	return
}

func nocolor() clog.ColorOptions {
	return clog.ColorOptions{
		Colorer: ansi.NewColorer(ansi.Colorable(func() bool { return false })),
	}
}

func stripANSI(s string) string {
	var stripped string

	skip := false

	for _, c := range s {
		switch c {
		case '\x1b':
			skip = true
		case 'm':
			skip = false
			continue
		}
		if !skip {
			stripped += string(c)
		}
	}

	return stripped
}
