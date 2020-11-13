package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/test"
	"github.com/kr/pretty"
)

func TestMapStrings(t *testing.T) {
	var kinds test.Kinds
	kinds.Add((*Record)(nil))
	fruit := []string{"Orange", "Lemon", "Mango", "Banana", "Lime"}
	lt := listTime{
		src:   g.StringsOf(fruit),
		res:   g.StringsOf(nil),
		remap: &reverseStrings,
		kinds: &kinds,
	}
	if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if res, e := lt.res.GetTextList(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(res, []string{
		"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
	}); len(diff) > 0 {
		t.Fatal(res)
	} else {
		t.Log("ok", res)
	}
}

type Record struct{ Fruit string }

func TestMapRecords(t *testing.T) {
	var kinds test.Kinds
	kinds.Add((*Record)(nil))
	lt := listTime{
		kinds: &kinds,
		remap: &reverseRecords,
	}
	if k, e := lt.GetKindByName("Record"); e != nil {
		t.Fatal(e)
	} else {
		var fruits []*g.Record
		for _, f := range []string{"Orange", "Lemon", "Mango", "Banana", "Lime"} {
			one := k.NewRecord()
			if e := one.SetNamedField("Fruit", g.StringOf(f)); e != nil {
				t.Fatal(e)
			}
			fruits = append(fruits, one)
		}
		lt.src = g.RecordsOf(k, fruits)
		lt.res = g.RecordsOf(k, nil)
	}

	if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if res, e := lt.res.GetRecordList(); e != nil {
		t.Fatal(e)
	} else {
		expect := []string{
			"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
		}
		for i, el := range res {
			if v, e := el.GetNamedField("Fruit"); e != nil {
				t.Fatal(e)
			} else if s, e := v.GetText(); e != nil {
				t.Fatal(e)
			} else if x := expect[i]; x != s {
				t.Fatal(x)
			}
		}
	}
}

var remap = list.Map{FromList: "src", ToList: "res", Pattern: "remap"}

var reverseRecords = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "remap",
		Prologue: []term.Preparer{
			&term.Record{Name: "in", Kind: "Record"},
			&term.Record{Name: "out", Kind: "Record"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		&pattern.ExecuteRule{
			Execute: &core.SetField{
				Obj:   V("out"),
				Field: T("Fruit"),
				From: &core.FromText{
					&core.MakeReversed{
						&core.GetField{
							Obj:   V("in"),
							Field: T("Fruit"),
						},
					},
				},
			},
		},
	},
}

var reverseStrings = pattern.ActivityPattern{
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
