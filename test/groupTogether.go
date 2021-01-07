package test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
)

var runGroupTogther = list.Map{
	FromList:     &core.Var{Name: "Objects"},
	ToList:       "Settings",
	UsingPattern: "assignGrouping"}

// from a list of object names, build a list of group settings
var assignGrouping = pattern.ActivityPattern{
	CommonPattern: pattern.CommonPattern{
		Name: "assignGrouping",
		Prologue: []term.Preparer{
			&term.Text{Name: "in"},
			&term.Record{Name: "out", Kind: "GroupSettings"},
		},
	},
	Rules: []*pattern.ExecuteRule{
		{Execute: &core.Activity{[]rt.Execute{
			Put("out", "Name", V("in")),
			&core.ChooseAction{
				If: &core.Matches{
					Text:    &core.Var{Name: "in"},
					Pattern: "^thing"},
				Do: core.MakeActivity(
					Put("out", "Label", &core.FromText{&core.Text{"thingies"}}),
				),
			},
		}}},
	},
}
