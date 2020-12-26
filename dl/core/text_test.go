package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/safe"
	"github.com/kr/pretty"
)

func TestText(t *testing.T) {
	var run baseRuntime

	t.Run("is", func(t *testing.T) {
		if e := testTrue(t, &run, &IsTrue{&Bool{true}}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsNotTrue{&Bool{false}}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("isEmpty", func(t *testing.T) {
		if e := testTrue(t, &run, &IsEmpty{&Text{}}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsNotTrue{&IsEmpty{&Text{"xxx"}}}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("includes", func(t *testing.T) {
		if e := testTrue(t, &run, &Includes{
			&Text{"full"},
			&Text{"ll"},
		}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsNotTrue{&Includes{
			&Text{"full"},
			&Text{"bull"},
		}}); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("join", func(t *testing.T) {
		if e := testTrue(t, &run, &CompareText{
			&Join{Parts: []rt.TextEval{
				&Text{"one"}, &Text{"two"}, &Text{"three"},
			}},
			&EqualTo{},
			&Text{"onetwothree"},
		}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{
			&Join{&Text{" "}, []rt.TextEval{
				&Text{"one"}, &Text{"two"}, &Text{"three"},
			}},
			&EqualTo{},
			&Text{"one two three"},
		}); e != nil {
			t.Fatal(e)
		}
	})
}

func testTrue(t *testing.T, run rt.Runtime, eval rt.BoolEval) (err error) {
	if ok, e := safe.GetBool(run, eval); e != nil {
		err = e
	} else if !ok.Bool() {
		err = errutil.New("expected true", pretty.Sprint(eval))
	}
	return
}
