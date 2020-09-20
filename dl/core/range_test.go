package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
)

func TestRange(t *testing.T) {
	list := func(eval rt.NumListEval) (ret []float64, err error) {
		var run baseRuntime
		if it, e := rt.GetNumberStream(&run, eval); e != nil {
			err = e
		} else {
			cnt, total := it.(rt.StreamCount).Remaining(), 0
			for it.HasNext() {
				if n, e := it.GetNext(); e != nil {
					err = e
					break
				} else if n, e := n.GetNumber(&run); e != nil {
					err = e
					break
				} else {
					ret = append(ret, n)
				}
				total++
			}
			if err == nil && cnt != total {
				err = errutil.New("reported count", cnt, "actual count", total)
			}
		}
		return
	}
	t.Run("range(10)", func(t *testing.T) {
		want := []float64{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		}
		if have, e := list(&Range{Stop: &Number{10}}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(2, 11)", func(t *testing.T) {
		want := []float64{
			2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
		}
		if have, e := list(&Range{Start: &Number{2}, Stop: &Number{11}}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(0, 30, 5)", func(t *testing.T) {
		want := []float64{
			0, 5, 10, 15, 20, 25, 30,
		}
		if have, e := list(&Range{
			Start: &Number{0},
			Stop:  &Number{30},
			Step:  &Number{5},
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(0, 9, 3)", func(t *testing.T) {
		want := []float64{
			0, 3, 6, 9,
		}
		if have, e := list(&Range{
			Start: &Number{0},
			Stop:  &Number{9},
			Step:  &Number{3},
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(0, -10, -1)", func(t *testing.T) {
		want := []float64{
			0, -1, -2, -3, -4, -5, -6, -7, -8, -9, -10,
		}
		if have, e := list(&Range{
			Start: &Number{0},
			Stop:  &Number{-10},
			Step:  &Number{-1},
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(1)", func(t *testing.T) {
		want := []float64{1}
		if have, e := list(&Range{
			Stop: &Number{1},
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(0)", func(t *testing.T) {
		want := []float64{}
		if have, e := list(&Range{
			Stop: &Number{0},
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(1, 0)", func(t *testing.T) {
		want := []float64{}
		if have, e := list(&Range{
			Start: &Number{1},
			Stop:  &Number{0},
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})

}
