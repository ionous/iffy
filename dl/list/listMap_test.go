package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestMapStrings(t *testing.T) {
	var kinds testutil.Kinds
	type Fruit struct {
		Name string
	}
	type Values struct {
		Fruits, Results []string
	}
	kinds.AddKinds((*Fruit)(nil), (*Values)(nil))
	values := kinds.New("Values")
	lt := listTime{
		PatternMap: testutil.PatternMap{
			"remap": &reverseStrings,
		},
		ScopeStack: scope.ScopeStack{
			Scopes: []rt.Scope{
				&scope.TargetRecord{object.Variables, values},
			},
		},
		Kinds: &kinds,
	}
	if e := values.SetNamedField("Fruits", g.StringsOf([]string{"Orange", "Lemon", "Mango", "Banana", "Lime"})); e != nil {
		t.Fatal(e)
	} else if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if results, e := values.GetNamedField("Results"); e != nil {
		t.Fatal(e)
	} else {
		res := results.Strings()
		if diff := pretty.Diff(res, []string{
			"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
		}); len(diff) > 0 {
			t.Fatal(res)
		} else {
			t.Log("ok", res)
		}
	}
}

func TestMapRecords(t *testing.T) {
	var kinds testutil.Kinds
	type Fruit struct {
		Name string
	}
	type Values struct {
		Fruits  []Fruit
		Results []Fruit
	}
	kinds.AddKinds((*Fruit)(nil), (*Values)(nil))
	values := kinds.New("Values")
	if k, e := kinds.GetKindByName("Fruit"); e != nil {
		t.Fatal(e)
	} else {
		var fruits []*g.Record
		for _, f := range []string{"Orange", "Lemon", "Mango", "Banana", "Lime"} {
			one := k.NewRecord()
			if e := one.SetNamedField("Name", g.StringOf(f)); e != nil {
				t.Fatal(e)
			}
			fruits = append(fruits, one)
		}
		if e := values.SetNamedField("Fruits", g.RecordsOf(k.Name(), fruits)); e != nil {
			t.Fatal(e)
		}
	}
	lt := listTime{
		Kinds: &kinds,
		PatternMap: testutil.PatternMap{
			"remap": &reverseRecords,
		},
		ScopeStack: scope.ScopeStack{
			Scopes: []rt.Scope{
				&scope.TargetRecord{object.Variables, values},
			},
		},
	}
	if e := remap.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if val, e := values.GetNamedField("Results"); e != nil {
		t.Fatal(e)
	} else if res := val.Records(); len(res) != 5 {
		t.Fatal("missing results")
	} else {
		expect := []string{
			"egnarO", "nomeL", "ognaM", "ananaB", "emiL",
		}
		var got []string
		for _, el := range res {
			if v, e := el.GetNamedField("Name"); e != nil {
				t.Fatal(e)
			} else {
				got = append(got, v.String())
			}
		}
		if diff := pretty.Diff(expect, got); len(diff) > 0 {
			t.Fatal("error", got)
		}
	}
}

var remap = list.Map{FromList: "Fruits", ToList: "Results", UsingPattern: "remap"}

var reverseRecords = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "remap",
		Prologue: []term.Preparer{
			&term.Record{Name: "in", Kind: "Fruit"},
			&term.Record{Name: "out", Kind: "Fruit"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		&pattern.ExecuteRule{
			Execute: &core.SetField{
				Obj:   V("out"),
				Field: "Name",
				From: &core.FromText{
					&core.MakeReversed{
						&core.GetField{
							Obj:   V("in"),
							Field: "Name",
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
						&core.Var{
							Name: "in",
						},
					},
				},
			},
		},
	},
}
