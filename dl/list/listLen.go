package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
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

func (op *Len) GetNumber(run rt.Runtime) (ret float64, err error) {
	if v, e := run.GetField(object.Variables, op.List); e != nil {
		err = cmdError(op, e)
	} else if l, e := v.GetLen(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = float64(l)
	}
	return
}
