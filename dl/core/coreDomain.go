package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type HasDominion struct {
	Name string
}

func (*HasDominion) Compose() composer.Spec {
	return composer.Spec{
		Name:  "has_dominion",
		Group: "logic",
	}
}

func (op *HasDominion) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return run.GetField(object.Domain, op.Name)
}
