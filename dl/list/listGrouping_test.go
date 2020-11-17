package list_test

// import (
// 	"testing"

// 	"github.com/ionous/iffy/dl/core"
// 	"github.com/ionous/iffy/dl/list"
// 	"github.com/ionous/iffy/dl/pattern"
// 	"github.com/ionous/iffy/dl/term"
// 	"github.com/ionous/iffy/rt"
// 	g "github.com/ionous/iffy/rt/generic"
// 	"github.com/ionous/iffy/rt/test"
// 	"github.com/ionous/sliceOf"
// 	"github.com/kr/pretty"
// )

// // - make groups, take a list of object (names)
// // - map "in" (object?) to out "Group Together"
// // - then sort

// func TestGrouping(t *testing.T) {
// 	var kinds test.Kinds

// 	type Things struct{}
// 	kinds.Add((*Things)(nil), (*test.GroupCollation)(nil))
// 	//
// 	if objs, e := objects(kinds.Kind("Things"), "mildred", "apple", "pen", "thing#1", "thing#2"); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		lt := listTime{
// 			kinds: &kinds,
// 			objs:  objs,
// 			// fix: change test to a record
// 			vals: values{
// 				"objects":   g.StringsOf(sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2")),
// 				"groups":    g.RecordsOf(kinds.Kind("GroupSettings"), nil),
// 				"collation": g.RecordOf(kinds.Kind("GroupCollation").NewRecord()),
// 			},
// 		}

// 		if e := runGroupTogther.Execute(&lt); e != nil {
// 			t.Fatal(e)
// 		} else if e := runCollateGroups.Execute(&lt); e != nil {
// 			t.Fatal(e)
// 		} else {
// 			pretty.Println(lt.vals["collation"])
// 		}
// 	}
// }

// func objects(kind *g.Kind, names ...string) (ret map[string]*g.Record, err error) {
// 	out := make(map[string]*g.Record)
// 	for _, name := range names {
// 		// we'll use normal records for this test....
// 		out[name] = kind.NewRecord()
// 	}
// 	if err == nil {
// 		ret = out
// 	}
// 	return
// }

// var runGroupTogther = list.Map{FromList: "objects", ToList: "groups", UsingPattern: "groupTogether"}
// var runCollateGroups = list.Reduce{FromList: "groups", IntoValue: "collection", UsingPattern: "collateGroups"}

// func T(t string) rt.TextEval {
// 	return &core.Text{t}
// }

// var collateGroups = pattern.ActivityPattern{
// 	CommonPattern: pattern.CommonPattern{
// 		Name: "collateGroups",
// 		Prologue: []term.Preparer{
// 			&term.Record{Name: "in", Kind: "GroupSettings"},
// 			&term.Record{Name: "out", Kind: "GroupCollation"},
// 		},
// 		Locals: []term.Preparer{
// 			&term.Number{Name: "idx"},
// 			&term.RecordList{Name: "groups", Kind: "GroupObjects"},
// 			&term.Record{Name: "group", Kind: "GroupObjects"},
// 			&term.TextList{Name: "names"},
// 		},
// 	},
// 	Rules: []*pattern.ExecuteRule{
// 		&pattern.ExecuteRule{Execute: core.NewActivity(
// 			// // walk out.Groups for matching settings
// 			&core.Assign{"groups", &core.GetField{V("out"), "groups"}},
// 			&list.Each{
// 				List: "groups",
// 				With: "el",
// 				Go: core.NewActivity(
// 					&core.Choose{
// 						If: &pattern.DetermineBool{
// 							Pattern:   "isMatchingGroup",
// 							Arguments: core.NewArgs(V("in"), V("el"))},
// 						True: core.NewActivity(
// 							&core.Assign{
// 								Name: "idx",
// 								From: V("index"),
// 							},
// 							// implement a "break" for the each that returns a constant error?
// 						),
// 					},
// 				)}, // end go-each
// 			&core.Choose{
// 				If: &core.CompareNum{
// 					A:  V("idx"),
// 					Is: &core.EqualTo{},
// 					B:  &core.Num{0},
// 				},
// 				// havent found a matching group?
// 				// make sure the scratch group is empty,
// 				// pack the object and its settings into it,
// 				// push the group into the groups.
// 				True: core.NewActivity(
// 					&core.Assign{"group", &core.Make{"GroupObjects"}},
// 					&core.Assign{"names", &core.GetField{V("group"), "objects"}},
// 					&list.Push{"names", &core.GetField{V("in"), "Name"}, true},
// 					&core.SetField{V("group"), "objects", V("names")},
// 					&list.Push{"groups", V("group"), true},
// 				), // end true
// 				// found a matching group?
// 				// unpack it, add the object to it, then pack it up again.
// 				False: core.NewActivity(
// 					&core.Assign{"group", &list.At{"groups", V("idx")}},
// 					&core.Assign{"names", &core.GetField{V("group"), "objects"}},
// 					&list.Push{"names", &core.GetField{V("in"), "Name"}, true},
// 					&core.SetField{V("group"), "objects", V("names")},
// 					&list.Set{V("groups"), V("idx"), V("group")},
// 				), // end false
// 			},
// 		)},
// 	},
// }

// // a pattern for matching groups --
// // we add rules that if things arent equal we return false
// // we could have just one rule if we wanted
// var isMatchingGroup = pattern.BoolPattern{
// 	CommonPattern: pattern.CommonPattern{
// 		Name: "isMatchingGroup",
// 		Prologue: []term.Preparer{
// 			&term.Record{Name: "a", Kind: "GroupSettings"},
// 			&term.Record{Name: "b", Kind: "GroupSettings"},
// 		},
// 	},
// 	Rules: []*pattern.BoolRule{{
// 		&core.CompareText{
// 			&core.GetField{
// 				Obj:   V("a"),
// 				Field: "Label",
// 			},
// 			&core.NotEqualTo{},
// 			&core.GetField{
// 				Obj:   V("b"),
// 				Field: "Label",
// 			},
// 		}, &core.Bool{Bool: false}}, {
// 		&core.CompareText{
// 			&core.GetField{
// 				Obj:   V("a"),
// 				Field: "Innumerable",
// 			},
// 			&core.NotEqualTo{},
// 			&core.GetField{
// 				Obj:   V("b"),
// 				Field: "Innumerable",
// 			},
// 		}, &core.Bool{Bool: false}}, {
// 		&core.Bool{
// 			&core.GetField{
// 				Obj:   V("a"),
// 				Field: "Options",
// 			},
// 			&core.NotEqualTo{},
// 			&core.GetField{
// 				Obj:   V("b"),
// 				Field: "Options",
// 			},
// 		}, &core.Bool{Bool: false}},
// 	},
// }

// //
// var groupTogether = pattern.ActivityPattern{
// 	CommonPattern: pattern.CommonPattern{
// 		Name: "groupTogether",
// 		Prologue: []term.Preparer{
// 			&term.Text{Name: "in"},
// 			&term.Record{Name: "out", Kind: "GroupSettings"},
// 		},
// 	},
// 	Rules: []*pattern.ExecuteRule{
// 		{Execute: &core.Activity{[]rt.Execute{
// 			&core.SetField{
// 				Obj:   V("out"),
// 				Field: T("Name"),
// 				From: &core.FromText{
// 					Val: T("in"),
// 				},
// 			},
// 		}}},
// 	},
// }
