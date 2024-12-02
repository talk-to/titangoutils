package goutilstesting

import (
	"testing"

	"github.com/talk-to/titangoutils"
)

func TestHello(t *testing.T) {
	got := titangoutils.Hello("World")
	want := "Hello, World!"
	if got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
