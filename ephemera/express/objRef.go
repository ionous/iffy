package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
)

// used for converting requests for objects into automatic printing of an object's name
// ex. {.lantern}
type objRef struct {
	name rt.TextEval
}

func (on objRef) getTextName() rt.TextEval {
	return on.name
}

func (on objRef) getPrintedName() rt.TextEval {
	return &core.Buffer{core.NewActivity(
		&pattern.DetermineAct{
			Pattern: "printAName",
			Parameters: &pattern.Parameters{[]*pattern.Parameter{{
				Name: "$1",
				From: &core.FromText{on.name},
			}}}})}
}
