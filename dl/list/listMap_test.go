package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/kr/pretty"
)

func TestMap(t *testing.T) {
	fruit := []string{"Orange", "Lemon", "Mango", "Banana", "Lime"}
	if res, e := mapTest(fruit); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(res, []string{
		"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
	}); len(diff) > 0 {
		t.Fatal(res)
	} else {
		t.Log("ok", res)
	}
}

func mapTest(strs []string) (ret []string, err error) {
	remap := list.Map{FromList: "strings", ToList: "res", Pattern: "remap"}
	lt := listTime{strings: strs, remap: &remapPattern}
	if e := remap.Execute(&lt); e != nil {
		err = e
	} else {
		ret = lt.res
	}
	return
}

var remapPattern = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "remap",
		Prologue: []term.Preparer{
			&term.Text{Name: "in"},
			&term.Text{Name: "out"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		&pattern.ExecuteRule{
			Execute: &core.Assign{
				Name: "out",
				From: &core.FromText{
					&core.MakeReversed{
						&core.GetVar{
							Name: T("in"),
						},
					},
				},
			},
		},
	},
}
