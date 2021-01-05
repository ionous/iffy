package generic_test

import (
	"testing"

	"github.com/ionous/iffy/affine"
	g "github.com/ionous/iffy/rt/generic"
)

func TestSplices(t *testing.T) {
	zeroSplice := func(src g.Value) {
		if vs, e := src.Splice(0, 0, nil); e != nil {
			t.Fatal("empty splice should be legal")
		} else if vs == nil {
			t.Fatal("empty splice should return value")
		} else if a := vs.Affinity(); a != affine.TextList {
			t.Fatal("empty splice should return a text list not", a)
		} else if cnt := vs.Len(); cnt != 0 {
			t.Fatal("empty splice should return an empty list", cnt)
		}
	}
	zeroSplice(g.StringsOf(nil))
	zeroSplice(g.StringsOf([]string{"a"}))
	zeroSplice(g.StringsOf([]string{"a", "b", "c"}))
}
