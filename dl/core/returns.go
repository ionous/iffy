package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// Returns allows a local variable to be used as an output.
type Returns struct {
	Name  string
	Using *Activity
}

func (*Returns) Compose() composer.Spec {
	return composer.Spec{
		Name:  "returns",
		Spec:  "Return the variable {name:text} {?using}",
		Group: "variables",
		Desc:  "Return: Return the value of the named variable computed during using.",
	}
}

const returnNotImplemented = errutil.Error("return not implemented")

func (op *Returns) GetBool(run rt.Runtime) (ret g.Value, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetText(run rt.Runtime) (ret g.Value, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	err = returnNotImplemented
	return
}
