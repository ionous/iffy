package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Gather struct {
	Var   core.Variable  `if:"unlabeled"`
	From  FromListSource `if:"unlabeled"`
	Using pattern.PatternName
}

func (*Gather) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_gather",
		Fluent: &composer.Fluid{Name: "gather", Role: composer.Command},
		Group:  "list",
		Desc: `Gather list: Transform the values from a list.
		The named pattern gets called once for each value in the list.
		It get called with two parameters: 'in' as each value from the list, 
		and 'out' as the var passed to the gather.`,
	}
}

func (op *Gather) Execute(rt.Runtime) (ret g.Value, err error) {
	return
}
