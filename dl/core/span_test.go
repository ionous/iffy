package core

import (
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/sliceOf"
	"io"
	"testing"
)

func TestSpan(t *testing.T) {
	var p printer.Span
	words := sliceOf.String("hello", "there,", "world.")
	for _, w := range words {
		io.WriteString(&p, w)
	}
	expect := "hello there, world."
	if res := p.String(); res != expect {
		t.Fatalf("p should be '%s' was '%s'", expect, res)
	}
}
