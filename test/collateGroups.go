package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
)

var runCollateGroups = list.Reduce{
	FromList:     V("Settings"),
	IntoValue:    "Collation",
	UsingPattern: "collateGroups"}

var collateGroups = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "collateGroups",
		Params: []term.Preparer{
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
				As:   &list.AsRec{N("el")},
				Do: core.MakeActivity(
					&core.ChooseAction{
						If: &pattern.Determine{
							Pattern:   "matchGroups",
							Arguments: core.Args(V("settings"), &core.Unpack{V("el"), "Settings"})},
						Do: core.MakeActivity(
							&core.Assign{
								Var:  N("idx"),
								From: V("index"),
							},
							// implement a "break" for the each that returns a constant error?
						),
					},
				)}, // end go-each
			&core.ChooseAction{
				If: &core.CompareNum{
					A:  V("idx"),
					Is: &core.EqualTo{},
					B:  &core.Number{0},
				},
				// havent found a matching group?
				// pack the object and its settings into it,
				// push the group into the groups.
				Do: core.MakeActivity(
					&list.PutEdge{Into: &list.IntoTxtList{N("names")}, From: &core.Unpack{V("settings"), "Name"}},
					Put("group", "Objects", V("names")),
					Put("group", "Settings", V("settings")),
					&list.PutEdge{Into: &list.IntoRecList{N("groups")}, From: V("group")},
				), // end true
				// found a matching group?
				// unpack it, add the object to it, then pack it up again.
				Else: &core.ChooseNothingElse{core.MakeActivity(
					&core.Assign{N("group"), &core.FromRecord{&list.At{List: V("groups"), Index: V("idx")}}},
					&core.Assign{N("names"), &core.Unpack{V("group"), "Objects"}},
					&list.PutEdge{Into: &list.IntoTxtList{N("names")}, From: &core.Unpack{V("settings"), "Name"}},
					Put("group", "Objects", V("names")),
					&list.Set{List: "groups", Index: V("idx"), From: V("group")},
				), // end false
				},
			},
			Put("collation", "Groups", V("groups")),
		)},
	},
}

func V(n string) *core.Var {
	return &core.Var{Name: n}
}
func N(n string) core.Variable {
	return core.Variable{Str: n}
}
func Put(rec, field string, from core.Assignment) rt.Execute {
	return &core.PutAtField{Into: &core.IntoRec{Var: N(rec)}, AtField: field, From: from}
}
