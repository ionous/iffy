package list_test

// import (
// 	"testing"

// 	"github.com/ionous/iffy/dl/core"
// 	"github.com/ionous/iffy/dl/list"
// 	"github.com/ionous/iffy/dl/pattern"
// 	"github.com/ionous/iffy/dl/term"
// 	"github.com/ionous/iffy/rt"
// 	g "github.com/ionous/iffy/rt/generic"
// 	"github.com/kr/pretty"
// )

// func TestReduce(t *testing.T) {
// 	type Fruit struct {
// 		Name string
// 	}
// 	kinds.AddKinds((*Fruit)(nil))
// 	lt := listTime{
// 		kinds: &kinds,
// 		PatternMap: testutil.PatternMap{
// 			"reduce": &reduceRecords,
// 		},
// 		vals: values{},
// 	}
// 	if k, e := lt.GetKindByName("Fruit"); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		var fruits []*g.Record
// 		for _, f := range []string{"Orange", "Lemon", "Mango", "Banana", "Lime"} {
// 			one := k.NewRecord()
// 			if e := one.SetNamedField("Name", g.StringOf(f)); e != nil {
// 				t.Fatal(e)
// 			}
// 			fruits = append(fruits, one)
// 		}
// 		lt.vals["Source"] = g.RecordsOf(k, fruits)
// 		lt.vals["res"] = g.Empty
// 	}

// 	if e := reduce.Execute(&lt); e != nil {
// 		t.Fatal(e)
// 	} else if out, e := lt.vals["res"].GetText(); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		expected := "Orange, Lemon, Mango, Banana, Lime"
// 		if expected != out {
// 			t.Fatal(out)
// 		} else {
// 			pretty.Println("ok", out) // emiL
// 		}
// 	}
// }

// var reduce = list.Reduce{FromList: "Source", IntoValue: "res", UsingPattern: "reduce"}

// // join each record in turn
// var reduceRecords = pattern.Pattern{
// 		Name: "reduce",
// 		Params: []term.Preparer{
// 			&term.Record{Name: "in", Kind: "Fruit"},
// 			&term.Text{Name: "out"},
// 		},
// 	Rules: []*pattern.Rule{
// 		&pattern.Rule{
// 			Execute: &core.Assign{
// 				Name: "out",
// 				From: &core.FromText{&core.Join{Sep: T(", "), Parts: []rt.TextEval{
// 					V("out"),
// 					&core.Field{V("in"), "Name"},
// 				}}},
// 			},
// 		},
// 	},
// }
