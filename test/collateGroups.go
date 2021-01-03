package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
)

var runCollateGroups = list.Reduce{
	FromList:     V("Settings"),
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
			&core.Assign{core.Variable{Str: "groups"}, &core.Unpack{V("collation"), "Groups"}},
			&list.Each{
				List: V("groups"),
				With: "el",
				Go: core.NewActivity(
					&core.Choose{
						If: &pattern.DetermineBool{
							Pattern:   "matchGroups",
							Arguments: core.Args(V("settings"), &core.Unpack{V("el"), "Settings"})},
						True: core.NewActivity(
							&core.Assign{
								Name: N("idx"),
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
				// pack the object and its settings into it,
				// push the group into the groups.
				True: core.NewActivity(
					&list.Push{List: "names", Insert: &core.Unpack{V("settings"), "Name"}},
					&core.Pack{V("group"), "Objects", V("names")},
					&core.Pack{V("group"), "Settings", V("settings")},
					&list.Push{List: "groups", Insert: V("group")},
				), // end true
				// found a matching group?
				// unpack it, add the object to it, then pack it up again.
				False: core.NewActivity(
					&core.Assign{N("group"), &core.FromRecord{&list.At{List: V("groups"), Index: V("idx")}}},
					&core.Assign{N("names"), &core.Unpack{V("group"), "Objects"}},
					&list.Push{List: "names", Insert: &core.Unpack{V("settings"), "Name"}},
					&core.Pack{V("group"), "Objects", V("names")},
					&list.Set{List: "groups", Index: V("idx"), From: V("group")},
				), // end false
			},
			&core.Pack{V("collation"), "Groups", V("groups")},
		)},
	},
}

func V(n string) *core.Var {
	return &core.Var{Name: n}
}
func N(n string) core.Variable {
	return core.Variable{Str: n}
}
