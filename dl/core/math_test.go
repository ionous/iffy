package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/test"
)

func TestMath(t *testing.T) {
	match := func(v float64, eval rt.NumberEval) (err error) {
		var run test.PanicRuntime
		if n, e := eval.GetNumber(&run); e != nil {
			err = e
		} else if n.Float() != v {
			err = errutil.Fmt("%v != %v (have != want)", n, v)
		}
		return
	}
	t.Run("Add", func(t *testing.T) {
		if e := match(11, &SumOf{&Number{1}, &Number{10}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Sub", func(t *testing.T) {
		if e := match(-9, &DiffOf{&Number{1}, &Number{10}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Mul", func(t *testing.T) {
		if e := match(200, &ProductOf{&Number{20}, &Number{10}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div", func(t *testing.T) {
		if e := match(2, &QuotientOf{&Number{20}, &Number{10}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("Div By Zero", func(t *testing.T) {
		if e := match(0, &QuotientOf{&Number{20}, &Number{0}}); e == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("Mod", func(t *testing.T) {
		if e := match(1, &RemainderOf{&Number{3}, &Number{2}}); e != nil {
			t.Fatal(e)
		}
	})
}
