package list_test

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/test"
	"github.com/kr/pretty"
)

func TestReduce(t *testing.T) {
	errutil.Panic = true
	var kinds test.Kinds
	kinds.AddKinds((*Fruit)(nil))
	lt := listTime{
		kinds: &kinds,
		PatternMap: pattern.PatternMap{
			"reduce": &reduceRecords,
		},
		vals: values{},
	}
	if k, e := lt.GetKindByName("Fruit"); e != nil {
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
		lt.vals["src"] = g.RecordsOf(k, fruits)
		lt.vals["res"] = g.StringOf("")
	}

	if e := reduce.Execute(&lt); e != nil {
		t.Fatal(e)
	} else if out, e := lt.vals["res"].GetText(); e != nil {
		t.Fatal(e)
	} else {
		expected := "Orange, Lemon, Mango, Banana, Lime"
		if expected != out {
			t.Fatal(out)
		} else {
			pretty.Println("ok", out) // emiL
		}
	}
}

var reduce = list.Reduce{FromList: "src", IntoValue: "res", UsingPattern: "reduce"}

// join each record in turn
var reduceRecords = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "reduce",
		Prologue: []term.Preparer{
			&term.Record{Name: "in", Kind: "Fruit"},
			&term.Text{Name: "out"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		&pattern.ExecuteRule{
			Execute: &core.Assign{
				Name: "out",
				From: &core.FromText{&core.Join{Sep: T(", "), Parts: []rt.TextEval{
					V("out"),
					&core.GetField{V("in"), "Name"},
				}}},
			},
		},
	},
}
