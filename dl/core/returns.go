package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
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

func (op *Returns) GetBool(run rt.Runtime) (ret bool, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetText(run rt.Runtime) (ret string, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = returnNotImplemented
	return
}

func (op *Returns) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	err = returnNotImplemented
	return
}
