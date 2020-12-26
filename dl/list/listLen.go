package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Len struct {
	List core.Assignment
}

func (*Len) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_len",
		Group: "list",
		Spec:  "length of {list:assignment}",
		Desc:  "Length of List: Determines the number of values in a list.",
	}
}

func (op *Len) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := core.GetAssignedValue(run, op.List); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(v.Len())
	}
	return
}
