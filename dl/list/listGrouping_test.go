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
	"github.com/ionous/iffy/rt/test"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
)

func TestMatching(t *testing.T) {
	var kinds test.Kinds
	type Things struct{}
	kinds.AddKinds((*Things)(nil), (*test.GroupSettings)(nil))
	k := kinds.Kind("GroupSettings")

	//
	lt := listTime{kinds: &kinds,
		PatternMap: pattern.PatternMap{
			"isMatchingGroup": &isMatchingGroup,
		},
	}

	a, b := k.NewRecord(), k.NewRecord()
	runMatching := &pattern.DetermineBool{
		Pattern: "isMatchingGroup", Arguments: core.NewArgs(
			&core.FromValue{g.RecordOf(a)},
			&core.FromValue{g.RecordOf(b)},
		)}
	// default should match
	{
		if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok != true {
			t.Fatal(e)
		}
	}
	// different labels shouldnt match
	{
		if e := test.SetRecord(a, "Label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok != false {
			t.Fatal(e)
		}
	}
	// same labels should match
	{
		if e := test.SetRecord(b, "Label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok != true {
			t.Fatal(e)
		}
	}
	// many fields should match
	{
		if e := test.SetRecord(a, "Innumerable", "IsInnumerable"); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(b, "IsInnumerable", true); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(a, "GroupOptions", "WithArticles"); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(b, "GroupOptions", "WithArticles"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok != true {
			t.Fatal(e)
		}
	}
	// names shouldnt be involved
	{
		if e := test.SetRecord(a, "Name", "hola"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok != true {
			t.Fatal(e)
		}
	}
}

func TestGrouping(t *testing.T) {
	var kinds test.Kinds
	type Things struct{}
	type Values struct {
		Objects   []string
		Settings  []test.GroupSettings
		Collation test.GroupCollation
	}
	kinds.AddKinds((*Things)(nil), (*Values)(nil))
	objectNames := sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2")
	//
	if objs, e := objects(kinds.Kind("Things"), objectNames...); e != nil {
		t.Fatal(e)
	} else {
		values := kinds.New("Values")
		lt := listTime{
			kinds: &kinds,
			objs:  objs,
			ScopeStack: scope.ScopeStack{
				Scopes: []rt.Scope{
					&scope.TargetRecord{object.Variables, values},
				},
			},
			PatternMap: pattern.PatternMap{
				"groupTogether":   &groupTogether,
				"collateGroups":   &collateGroups,
				"isMatchingGroup": &isMatchingGroup,
			},
		}
		if e := values.SetNamedField("Objects", g.StringsOf(objectNames)); e != nil {
			t.Fatal(e)
		} else if e := runGroupTogther.Execute(&lt); e != nil {
			t.Fatal("groupTogther", e)
		} else if e := runCollateGroups.Execute(&lt); e != nil {
			t.Fatal("collateGroups", e)
		} else if collation, e := values.GetNamedField("Collation"); e != nil {
			t.Fatal(e)
		} else if groups, e := g.Must(collation.GetNamedField("Groups")).GetRecordList(); e != nil {
			t.Fatal(e)
		} else {
			expect := []interface{}{
				map[string]interface{}{
					"Settings": map[string]interface{}{
						"Name":         "mildred",
						"Label":        "",
						"Innumerable":  "NotInnumerable",
						"GroupOptions": "WithoutObjects",
					},
					"Objects": []string{"mildred", "apple", "pen"},
				},
				map[string]interface{}{
					"Settings": map[string]interface{}{
						"Name":         "thing#1",
						"Label":        "thingies",
						"Innumerable":  "NotInnumerable",
						"GroupOptions": "WithoutObjects",
					},
					"Objects": []string{"thing#1", "thing#2"},
				},
			}
			if diff := pretty.Diff(expect, g.RecordsToValue(groups)); len(diff) > 0 {
				t.Fatal(diff)
			}
		}
	}
}

func logGroups(t *testing.T, groups []*g.Record) {
	t.Log("groups", len(groups), pretty.Sprint(g.RecordsToValue(groups)))
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

var runGroupTogther = list.Map{FromList: "Objects", ToList: "Settings", UsingPattern: "groupTogether"}
var runCollateGroups = list.Reduce{FromList: "Settings", IntoValue: "Collation", UsingPattern: "collateGroups"}

// from a list of object names, build a list of group settings
var groupTogether = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "groupTogether",
		Prologue: []term.Preparer{
			&term.Text{Name: "in"},
			&term.Record{Name: "out", Kind: "GroupSettings"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		{Execute: &core.Activity{[]rt.Execute{
			&core.SetField{
				Obj:   V("out"),
				Field: "Name",
				From:  V("in"), // fix? this is no way ensures that the object is valid.
			},
			&core.Choose{
				If: &core.Matches{
					Text:    V("in"),
					Pattern: "^thing"},
				True: core.NewActivity(
					&core.SetField{
						Obj:   V("out"),
						Field: "Label",
						From:  &core.FromText{T("thingies")},
					},
				),
			},
		}}},
	},
}

var collateGroups = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "collateGroups",
		Prologue: []term.Preparer{
			&term.Record{Name: "in", Kind: "GroupSettings"},
			&term.Record{Name: "out", Kind: "GroupCollation"},
		},
		Locals: []term.Preparer{
			&term.Number{Name: "idx"},
			&term.RecordList{Name: "groups", Kind: "GroupObjects"},
			&term.Record{Name: "group", Kind: "GroupObjects"},
			&term.TextList{Name: "names"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		&pattern.ExecuteRule{Execute: core.NewActivity(
			// walk out.Groups for matching settings
			&core.Assign{"groups", &core.GetField{V("out"), "Groups"}},
			&list.Each{
				List: "groups",
				With: "el",
				Go: core.NewActivity(
					&core.Choose{
						If: &pattern.DetermineBool{
							Pattern:   "isMatchingGroup",
							Arguments: core.NewArgs(V("in"), &core.GetField{V("el"), "Settings"})},
						True: core.NewActivity(
							&core.Assign{
								Name: "idx",
								From: V("index"),
							},
							// implement a "break" for the each that returns a constant error?
						),
					},
				)}, // end go-each
			&core.Choose{
				If: &core.CompareNum{
					A:  V("idx"),
					Is: &core.EqualTo{},
					B:  &core.Number{0},
				},
				// havent found a matching group?
				// make sure the scratch group is empty,
				// pack the object and its settings into it,
				// push the group into the groups.
				True: core.NewActivity(
					&list.Push{"names", &core.GetField{V("in"), "Name"}, false},
					&core.SetField{V("group"), "Objects", V("names")},
					&core.SetField{V("group"), "Settings", V("in")},
					&list.Push{"groups", V("group"), false},
				), // end true
				// found a matching group?
				// unpack it, add the object to it, then pack it up again.
				False: core.NewActivity(
					&core.Assign{"group", &core.FromObject{&list.At{"groups", V("idx")}}},
					&core.Assign{"names", &core.GetField{V("group"), "Objects"}},
					&list.Push{"names", &core.GetField{V("in"), "Name"}, false},
					&core.SetField{V("group"), "Objects", V("names")},
					&list.Set{"groups", V("idx"), V("group")},
				), // end false
			},
			&core.SetField{V("out"), "Groups", V("groups")},
		)},
	},
}

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var isMatchingGroup = pattern.BoolPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "isMatchingGroup",
		Prologue: []term.Preparer{
			&term.Record{Name: "a", Kind: "GroupSettings"},
			&term.Record{Name: "b", Kind: "GroupSettings"},
		},
	},
	// rules are evaluated in reverse order ( see splitRules )
	Rules: []*pattern.BoolRule{{
		&core.Always{}, &core.Bool{true}}, {
		&core.CompareText{
			&core.GetField{
				Obj:   V("a"),
				Field: "Label",
			},
			&core.NotEqualTo{},
			&core.GetField{
				Obj:   V("b"),
				Field: "Label",
			},
		}, &core.Bool{Bool: false}}, {
		&core.CompareText{
			&core.GetField{
				Obj:   V("a"),
				Field: "Innumerable",
			},
			&core.NotEqualTo{},
			&core.GetField{
				Obj:   V("b"),
				Field: "Innumerable",
			},
		}, &core.Bool{Bool: false}}, {
		&core.CompareText{
			&core.GetField{
				Obj:   V("a"),
				Field: "GroupOptions",
			},
			&core.NotEqualTo{},
			&core.GetField{
				Obj:   V("b"),
				Field: "GroupOptions",
			},
		}, &core.Bool{Bool: false}},
	},
}
