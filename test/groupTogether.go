package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
)

var runGroupTogther = list.Map{
	FromList: "Objects", ToList: "Settings",
	UsingPattern: "groupTogether"}

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
			&core.Pack{
				Record: &core.Var{Name: "out"},
				Field:  "Name",
				From:   &core.Var{Name: "in"}, // fix? this is no way ensures that the object is valid.
			},
			&core.Choose{
				If: &core.Matches{
					Text:    &core.Var{Name: "in"},
					Pattern: "^thing"},
				True: core.NewActivity(
					&core.Pack{
						Record: &core.Var{Name: "out"},
						Field:  "Label",
						From:   &core.FromText{&core.Text{"thingies"}},
					},
				),
			},
		}}},
	},
}
