package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
)

var runCollateGroups = list.Reduce{
	FromList:     &core.Var{Name: "Settings"},
	IntoValue:    "Collation",
	UsingPattern: "collateGroups"}

var collateGroups = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "collateGroups",
		Prologue: []term.Preparer{
			&term.Record{Name: "settings", Kind: "GroupSettings"},
			&term.Record{Name: "collation", Kind: "GroupCollation"},
		},
		Locals: []term.Preparer{
			&term.Number{Name: "idx"},
			&term.RecordList{Name: "groups", Kind: "GroupedObjects"},
			&term.Record{Name: "group", Kind: "GroupedObjects"},
			&term.TextList{Name: "names"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		&pattern.ExecuteRule{Execute: core.NewActivity(
			// walk collation.Groups for matching settings
			&core.Assign{"groups", &core.Unpack{&core.Var{Name: "collation"}, "Groups"}},
			&list.Each{
				List: &core.Var{Name: "groups"},
				With: "el",
				Go: core.NewActivity(
					&core.Choose{
						If: &pattern.DetermineBool{
							Pattern:   "matchGroups",
							Arguments: core.Args(&core.Var{Name: "settings"}, &core.Unpack{&core.Var{Name: "el"}, "Settings"})},
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
					&list.Push{List: "names", Insert: &core.Unpack{&core.Var{Name: "settings"}, "Name"}},
					&core.Pack{&core.Var{Name: "group"}, "Objects", &core.Var{Name: "names"}},
					&core.Pack{&core.Var{Name: "group"}, "Settings", &core.Var{Name: "settings"}},
					&list.Push{List: "groups", Insert: &core.Var{Name: "group"}},
				), // end true
				// found a matching group?
				// unpack it, add the object to it, then pack it up again.
				False: core.NewActivity(
					&core.Assign{"group", &core.FromRecord{&list.At{List: &core.Var{Name: "groups"}, Index: &core.Var{Name: "idx"}}}},
					&core.Assign{"names", &core.Unpack{&core.Var{Name: "group"}, "Objects"}},
					&list.Push{List: "names", Insert: &core.Unpack{&core.Var{Name: "settings"}, "Name"}},
					&core.Pack{&core.Var{Name: "group"}, "Objects", &core.Var{Name: "names"}},
					&list.Set{List: "groups", Index: &core.Var{Name: "idx"}, From: &core.Var{Name: "group"}},
				), // end false
			},
			&core.Pack{&core.Var{Name: "collation"}, "Groups", &core.Var{Name: "groups"}},
		)},
	},
}
