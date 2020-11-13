package list_test

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/test"
	"github.com/ionous/sliceOf"
)

// - make groups, take a list of object (names)
// - map "in" (object?) to out "Group Together"
// - then sort

func TestGrouping(t *testing.T) {
	var kinds test.Kinds

	type Things struct{}
	kinds.Add((*Things)(nil), (*test.Groupings)(nil))
	//
	if objs, e := objects(kinds.Kind("Things"), "mildred", "apple", "pen", "thing#1", "thing#2"); e != nil {
		t.Fatal(e)
	} else {
		lt := listTime{
			kinds: &kinds,
			objs:  objs,
			vals: values{
				"src": g.StringsOf(sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2")),
				"res": g.RecordsOf(kinds.Kind("GroupSettings"), nil),
			},
			// remap: &groupObjects,
		}

		if e := groupAllTogther.Execute(&lt); e != nil {
			t.Fatal(e)
		} else if res, e := lt.vals["res"].GetRecordList(); e != nil {
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
}

func objects(kind *g.Kind, names ...string) (ret map[string]*g.Record, err error) {
	out := make(map[string]*g.Record)
	for _, name := range names {
		// we'll use normal records for this test....
		out[name] = kind.NewRecord()
	}
	if err == nil {
		ret = out
	}
	return
}

var groupAllTogther = list.Map{FromList: "src", ToList: "res", Pattern: "groupTogether"}

// len(k.Label) > 0 || k.ObjectGrouping != WithoutObjects
var sortGroups = list.Sort{List: "src", Pattern: "sortGroups"}

var groupCompare = pattern.BoolPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "sortGroups",
		Prologue: []term.Preparer{
			&term.Text{Name: "a"},
			&term.Text{Name: "b"},
		},
	},
	Rules: []*pattern.BoolRule{
		&pattern.BoolRule{
			Filter: B(true),
			BoolEval: &core.CompareText{
				A:  &core.GetVar{Name: T("a")},
				Is: &core.LessThan{},
				B:  &core.GetVar{Name: T("b")},
			},
		},
	},
}

//
var groupTogether = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "groupTogether",
		Prologue: []term.Preparer{
			&term.Text{Name: "in"},
			&term.Record{Name: "out", Kind: "GroupSettings"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		&pattern.ExecuteRule{Execute: &core.Activity{[]rt.Execute{
			&core.SetField{
				Obj:   V("out"),
				Field: T("Name"),
				From: &core.FromText{
					Val: T("in"),
				},
			},
		}}},
	},
}
