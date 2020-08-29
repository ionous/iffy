package core

import (
	"testing"

	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
)

func TestText(t *testing.T) {
	var run baseRuntime

	t.Run("is", func(t *testing.T) {
		testTrue(t, &run, &IsTrue{&Bool{true}})
		testTrue(t, &run, &IsNotTrue{&Bool{false}})
	})

	t.Run("isEmpty", func(t *testing.T) {
		testTrue(t, &run, &IsEmpty{&Text{}})
		testTrue(t, &run, &IsNotTrue{&IsEmpty{&Text{"xxx"}}})
	})

	t.Run("includes", func(t *testing.T) {
		testTrue(t, &run, &Includes{
			&Text{"full"},
			&Text{"ll"},
		})
		testTrue(t, &run, &IsNotTrue{&Includes{
			&Text{"full"},
			&Text{"bull"},
		}})
	})

	t.Run("join", func(t *testing.T) {
		testTrue(t, &run, &CompareText{
			&Join{Parts: []rt.TextEval{
				&Text{"one"}, &Text{"two"}, &Text{"three"},
			}},
			&EqualTo{},
			&Text{"onetwothree"},
		})
		testTrue(t, &run, &CompareText{
			&Join{&Text{" "}, []rt.TextEval{
				&Text{"one"}, &Text{"two"}, &Text{"three"},
			}},
			&EqualTo{},
			&Text{"one two three"},
		})
	})
}

func testTrue(t *testing.T, run rt.Runtime, eval rt.BoolEval) {
	if ok, e := rt.GetBool(run, eval); e != nil {
		t.Fatal(e)
	} else if !ok {
		t.Fatal("expected true", pretty.Sprint(eval))
	}
}
