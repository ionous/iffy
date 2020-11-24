package test

import (
	"testing"

	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/test/testutil"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
)

// fix? i wonder if now "collation" can be the list of groups.
func TestGrouping(t *testing.T) {
	var kinds testutil.Kinds
	kinds.AddKinds((*Things)(nil), (*Values)(nil))
	objectNames := sliceOf.String("mildred", "apple", "pen", "thing#1", "thing#2")
	//
	if objs, e := objects(kinds.Kind("Things"), objectNames...); e != nil {
		t.Fatal(e)
	} else {
		values := kinds.New("Values", "Objects", objectNames)
		lt := testTime{
			Kinds: &kinds,
			objs:  objs,
			ScopeStack: scope.ScopeStack{
				Scopes: []rt.Scope{
					&scope.TargetRecord{object.Variables, values},
				},
			},
			PatternMap: testutil.PatternMap{
				"groupTogether":   &groupTogether,
				"collateGroups":   &collateGroups,
				"isMatchingGroup": &isMatchingGroup,
			},
		}
		if e := runGroupTogther.Execute(&lt); e != nil {
			t.Fatal("groupTogther", e)
		} else if e := runCollateGroups.Execute(&lt); e != nil {
			t.Fatal("collateGroups", e)
		} else if collation, e := values.GetNamedField("Collation"); e != nil {
			t.Fatal(e)
		} else if groups, e := collation.FieldByName("Groups"); e != nil {
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
			got := g.RecordsToValue(groups.Records())
			if diff := pretty.Diff(expect, got); len(diff) > 0 {
				t.Log(pretty.Sprint(got))
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
