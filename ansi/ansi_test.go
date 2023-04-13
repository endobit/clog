package ansi_test

import (
	"os"
	"testing"

	"github.com/endobit/clog/ansi"
)

func TestColorable(t *testing.T) {
	if err := os.Unsetenv("NO_COLOR"); err != nil { // just in case it is set
		t.Fatal(err)
	}

	c := ansi.NewColorer()

	s := c.Color("foo", ansi.Red)
	if s != "\x1b[31mfoo\x1b[0m" {
		t.Errorf("expected %q, got %q", "\x1b[31mfoo\x1b[0m", s)
	}

	if err := os.Setenv("NO_COLOR", "1"); err != nil {
		t.Fatal(err)
	}

	c = ansi.NewColorer()

	s = c.Color("foo", ansi.Red)
	if s != "foo" {
		t.Errorf("expected %q, got %q", "foo", s)
	}

	c = ansi.NewColorer(ansi.Colorable(func() bool { return false }))

	s = c.Color("foo", ansi.Red)
	if s != "foo" {
		t.Errorf("expected %q, got %q", "foo", s)
	}
}
