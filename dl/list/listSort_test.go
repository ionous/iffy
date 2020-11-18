package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/kr/pretty"
)

func TestSort(t *testing.T) {
	fruit := []string{"Orange", "Lemon", "Mango", "Banana", "Lime"}
	if res, e := sortTest(fruit); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(res, []string{
		"Banana", "Lemon", "Lime", "Mango", "Orange",
	}); len(diff) > 0 {
		t.Fatal(res)
	} else {
		t.Log("ok", res)
	}
}

func sortTest(src []string) (ret []string, err error) {
	if run, values, e := newListTime(src, pattern.PatternMap{
		"sort": &sortPattern,
	}); e != nil {
		err = e
	} else if e := rt.RunOne(run, &list.Sort{"Source", "sort"}); e != nil {
		err = e
	} else if res, e := g.Must(values.GetNamedField("Source")).GetTextList(); e != nil {
		err = e
	} else {
		ret = res
	}
	return
}

var sortPattern = pattern.BoolPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "sort",
		Prologue: []term.Preparer{
			&term.Text{Name: "a"},
			&term.Text{Name: "b"},
		},
	},
	Rules: []*pattern.BoolRule{
		&pattern.BoolRule{
			Filter: B(true),
			BoolEval: &core.CompareText{
				A:  V("a"),
				Is: &core.LessThan{},
				B:  V("b"),
			},
		},
	},
}
