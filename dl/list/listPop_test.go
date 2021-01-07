package list_test

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/kr/pretty"
)

func TestPopping(t *testing.T) {
	errutil.Panic = true
	// pop from the front of a list
	// front := popTest(true, 5, "Orange", "Lemon", "Mango")
	// if d := pretty.Diff(front, []string{"Orange", "Lemon", "Mango", "x", "x"}); len(d) > 0 {
	// 	t.Fatal("pop front", front)
	// }
	// pop from the back of a list
	back := popTest(false, 5, "Orange", "Lemon", "Mango")
	if d := pretty.Diff(back, []string{"Mango", "Lemon", "Orange", "x", "x"}); len(d) > 0 {
		t.Fatal("pop back", back)
	}
}

func popTest(front bool, amt int, src ...string) []string {
	var out []string
	var start int
	if front {
		start = 1
	} else {
		start = -1
	}
	pop := &list.Erasing{
		EraseIndex: list.EraseIndex{
			Count:   I(1),
			AtIndex: I(start),
			From:    &list.FromTxtList{N("Source")},
		},
		As: "text",
		Do: *core.NewActivity(&core.ChooseAction{
			If: &core.CompareNum{&list.Len{V("text")}, &core.EqualTo{}, I(0)},
			Do: core.MakeActivity(&Write{&out, T("x")}),
			Else: &core.ChooseNothingElse{
				Do: core.MakeActivity(&Write{&out, &list.At{V("text"), I(1)}}),
			},
		}),
	}
	if run, _, e := newListTime(src, nil); e != nil {
		panic(e)
	} else {
		for i := 0; i < amt; i++ {
			if e := pop.Execute(run); e != nil {
				panic(e)
			}
		}
	}
	return out
}
