package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
)

var runCollateGroups = list.Reduce{FromList: "Settings", IntoValue: "Collation", UsingPattern: "collateGroups"}

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
			&core.Assign{"groups", &core.GetField{&core.Var{Name: "out"}, "Groups"}},
			&list.Each{
				List: "groups",
				With: "el",
				Go: core.NewActivity(
					&core.Choose{
						If: &pattern.DetermineBool{
							Pattern:   "isMatchingGroup",
							Arguments: core.Args(&core.Var{Name: "in"}, &core.GetField{&core.Var{Name: "el"}, "Settings"})},
						True: core.NewActivity(
							&core.Assign{
								Name: "idx",
								From: &core.Var{Name: "index"},
							},
							// implement a "break" for the each that returns a constant error?
						),
					},
				)}, // end go-each
			&core.Choose{
				If: &core.CompareNum{
					A:  &core.Var{Name: "idx"},
					Is: &core.EqualTo{},
					B:  &core.Number{0},
				},
				// havent found a matching group?
				// pack the object and its settings into it,
				// push the group into the groups.
				True: core.NewActivity(
					&list.Push{List: "names", Insert: &core.GetField{&core.Var{Name: "in"}, "Name"}},
					&core.SetField{&core.Var{Name: "group"}, "Objects", &core.Var{Name: "names"}},
					&core.SetField{&core.Var{Name: "group"}, "Settings", &core.Var{Name: "in"}},
					&list.Push{List: "groups", Insert: &core.Var{Name: "group"}},
				), // end true
				// found a matching group?
				// unpack it, add the object to it, then pack it up again.
				False: core.NewActivity(
					&core.Assign{"group", &core.FromObject{&list.At{"groups", &core.Var{Name: "idx"}}}},
					&core.Assign{"names", &core.GetField{&core.Var{Name: "group"}, "Objects"}},
					&list.Push{List: "names", Insert: &core.GetField{&core.Var{Name: "in"}, "Name"}},
					&core.SetField{&core.Var{Name: "group"}, "Objects", &core.Var{Name: "names"}},
					&list.Set{"groups", &core.Var{Name: "idx"}, &core.Var{Name: "group"}},
				), // end false
			},
			&core.SetField{&core.Var{Name: "out"}, "Groups", &core.Var{Name: "groups"}},
		)},
	},
}
