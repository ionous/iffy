package list_test

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/kr/pretty"
)

func TestMapStrings(t *testing.T) {
	fruit := []string{"Orange", "Lemon", "Mango", "Banana", "Lime"}
	g.NewKind(nil, "Record", []g.Field{
		{Name: "Fruit", Affinity: affine.Text, Type: "string"},
	})
	lt := listTime{
		src:   g.StringsOf(fruit),
		res:   g.StringsOf(nil),
		remap: &reverseStrings,
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

func TestMapRecords(t *testing.T) {
	errutil.Panic = true
	kind := g.NewKind(nil, "Record", []g.Field{
		{Name: "Fruit", Affinity: affine.Text, Type: "string"},
	})
	var fruits []*g.Record
	for _, f := range []string{"Orange", "Lemon", "Mango", "Banana", "Lime"} {
		one := kind.NewRecord()
		if e := one.SetNamedField("Fruit", g.StringOf(f)); e != nil {
			t.Fatal(e)
		}
		fruits = append(fruits, one)
	}
	//
	lt := listTime{
		rec:   kind,
		src:   g.RecordsOf(kind, fruits),
		res:   g.RecordsOf(kind, nil),
		remap: &reverseRecords,
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
