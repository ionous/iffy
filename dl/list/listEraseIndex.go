package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

/**
 * erase: numEval
 * from: varName,
 * atIndex: num,
 */
type EraseIndex struct {
	Count   rt.NumberEval `if:"unlabeled"`
	From    ListSource    `if:"unlabeled"`
	AtIndex rt.NumberEval
}

func (*EraseIndex) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "erase", Role: composer.Command},
		Desc:   "Erase: remove one or more values from a list",
	}
}

func (op *EraseIndex) Execute(run rt.Runtime) (err error) {
	if _, e := op.pop(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *EraseIndex) pop(run rt.Runtime) (ret g.Value, err error) {
	if rub, e := safe.GetOptionalNumber(run, op.Count, 0); e != nil {
		err = e
	} else if els, e := GetListSource(run, op.From); e != nil {
		err = e
	} else if startOne, e := safe.GetNumber(run, op.AtIndex); e != nil {
		err = e
	} else {
		start, listLen := startOne.Int(), els.Len()
		if start < 0 {
			start += listLen // wrap negative starts
		} else {
			start -= 1 // adjust to zero based
		}
		var end int
		if start >= listLen {
			start, end = 0, 0 // (still) out of bounds? do nothing.
		} else if rub := rub.Int(); rub <= 0 {
			start, end = 0, 0 // zero and negative removal means remove nothing
		} else {
			// If length + start is less than 0, begin from index 0.
			if start < 0 {
				start = 0
			}
			// too many elements means remove all.
			end = start + rub
			if end > listLen {
				end = listLen
			}
		}
		ret, err = els.Splice(start, end, nil)
	}
	return
}
