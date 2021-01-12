package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
)

func TestRange(t *testing.T) {
	compact := func(eval rt.NumListEval) (ret []float64, err error) {
		if vs, e := eval.GetNumList(nil); e != nil {
			err = e
		} else {
			for i, cnt := 0, vs.Len(); i < cnt; i++ {
				v := vs.Index(i)
				ret = append(ret, v.Float())
			}
		}
		return
	}
	t.Run("range(10)", func(t *testing.T) {
		want := []float64{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		}
		if have, e := compact(&list.Range{To: I(10)}); e != nil {
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
		if have, e := compact(&list.Range{From: I(2), To: I(11)}); e != nil {
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
		if have, e := compact(&list.Range{
			From:   I(0),
			To:     I(30),
			ByStep: I(5),
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
		if have, e := compact(&list.Range{
			From:   I(0),
			To:     I(9),
			ByStep: I(3),
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
		if have, e := compact(&list.Range{
			From:   I(0),
			To:     I(-10),
			ByStep: I(-1),
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
	t.Run("range(0, -9, -2)", func(t *testing.T) {
		want := []float64{
			0, -2, -4, -6, -8,
		}
		if have, e := compact(&list.Range{
			From:   I(0),
			To:     I(-9),
			ByStep: I(-2),
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
		if have, e := compact(&list.Range{
			To: I(1),
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
		if have, e := compact(&list.Range{
			To: I(0),
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
		if have, e := compact(&list.Range{
			From: I(1),
			To:   I(0),
		}); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(have, want); len(diff) > 0 {
			t.Fatal("have", have, "want", want, "diff", diff)
		} else {
			t.Log(have)
		}
	})
}
