package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/kr/pretty"
)

func TestSort(t *testing.T) {
	fruit := []string{"Orange", "Lemon", "Mango", "Banana", "Lime"}
	if e := sortTest(fruit); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(fruit, []string{
		"Banana", "Lemon", "Lime", "Mango", "Orange",
	}); len(diff) > 0 {
		t.Fatal(diff)
	}
	t.Log("ok", fruit)
}

func sortTest(strs []string) (err error) {
	sort := list.Sort{"src", "sort"}
	return sort.Execute(&listTime{
		vals: values{"src": g.StringsOf(strs)},
		PatternMap: pattern.PatternMap{
			"sort": &sortPattern,
		}})
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
