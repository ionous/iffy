package next

import (
	"testing"

	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
)

func TestText(t *testing.T) {
	var run baseRuntime

	t.Run("is", func(t *testing.T) {
		testTrue(t, &run, &Is{&Bool{true}})
		testTrue(t, &run, &IsNot{&Bool{false}})
	})

	t.Run("isEmpty", func(t *testing.T) {
		testTrue(t, &run, &IsEmpty{&Text{}})
		testTrue(t, &run, &IsNot{&IsEmpty{&Text{"xxx"}}})
	})

	t.Run("includes", func(t *testing.T) {
		testTrue(t, &run, &Includes{
			&Text{"full"},
			&Text{"ll"},
		})
		testTrue(t, &run, &IsNot{&Includes{
			&Text{"full"},
			&Text{"bull"},
		}})
	})

	t.Run("join", func(t *testing.T) {
		testTrue(t, &run, &CompareText{
			&Join{Elems: &Texts{oneTwoThree}},
			&EqualTo{},
			&Text{"onetwothree"},
		})
		testTrue(t, &run, &CompareText{
			&Join{&Texts{oneTwoThree}, &Text{" "}},
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
