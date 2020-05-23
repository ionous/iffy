package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

// used for converting requests
type objectName struct {
	name rt.TextEval
}

func (on objectName) getTextName() rt.TextEval {
	return on.name

}
func (on objectName) getPrintedName() rt.TextEval {
	return &core.DetermineText{
		Pattern: "print name",
		Parameters: &core.Parameters{[]*core.Parameter{
			&core.Parameter{
				Name: "$1",
				From: &core.FromText{on.name},
			}}}}
}
