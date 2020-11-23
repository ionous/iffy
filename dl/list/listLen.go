package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Len struct {
	List string // variable name
}

func (*Len) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_len",
		Group: "list",
		Spec:  "length of {list:text}",
		Desc:  "Length of List: Determines the number of values in a list.",
	}
}

func (op *Len) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.List(run, op.List); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(v.Len())
	}
	return
}
