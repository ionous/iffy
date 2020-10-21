package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type Splice struct {
	List          string        // variable name
	Start, Delete rt.NumberEval // from start
	Append        core.Assignment
}

// if start is negative, it will begin that many elements from the end of the array.
//  If array.length + start is less than 0, it will begin from index 0.

// If deleteCount is 0 or negative, no elements are removed.

func (*Splice) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_splice",
		Group: "list",
		Spec:  "slice {list:text} {from entry%start?number} {ending before entry%eend?number}",
		Desc: `Splice of List: Modify a list by adding and removing elements.
Note: the type of the elements being added must match the type of the list. 
Text cant be added to a list of numbers, numbers cant be added to a list of text, 
and true/false values can't be added to a list.`,
	}
}

// returns the removed elements
func (op *Splice) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = cmdError(op, e)
	} else if els, e := vs.GetNumList(); e != nil {
		err = cmdError(op, e)
	} else if i, j, e := op.getIndices(run, len(els)); e != nil {
		err = cmdError(op, e)
	} else if add, e := op.getNewFloats(run); e != nil {
		err = cmdError(op, e)
	} else {
		if i >= 0 && j >= i {
			// cut out i to j, then i becomes the insertion point.
			// as long as the range was valid we take the result and set it back...
			// even if no elements are cut or inserted.
			// ( that bakes any evaluation that might have been in the target )
			ret = copyfloats(els[i:j]) // before we start altering the memory of else, copy out
			newVals := append(els[:i], append(add, els[j:]...)...)
			if e := run.SetField(object.Variables, op.List, &generic.FloatSlice{Values: newVals}); e != nil {
				err = cmdError(op, e)
			}
		}
	}
	return
}

// returns the removed elements
func (op *Splice) GetTextList(run rt.Runtime) (ret []string, err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = cmdError(op, e)
	} else if els, e := vs.GetTextList(); e != nil {
		err = cmdError(op, e)
	} else if i, j, e := op.getIndices(run, len(els)); e != nil {
		err = cmdError(op, e)
	} else if add, e := op.getNewStrings(run); e != nil {
		err = cmdError(op, e)
	} else {
		// ... mirrors GetNumList()
		if i >= 0 && j >= i {
			ret = copystrings(els[i:j])
			newVals := append(els[:i], append(add, els[j:]...)...)
			if e := run.SetField(object.Variables, op.List, &generic.StringSlice{Values: newVals}); e != nil {
				err = cmdError(op, e)
			}
		}
	}
	return
}

func (op *Splice) getNewFloats(run rt.Runtime) (ret []float64, err error) {
	if assign := op.Append; assign != nil {
		if v, e := assign.GetAssignedValue(run); e != nil {
			err = e
		} else {
			switch a := v.Affinity(); a {
			case affine.Number:
				if one, e := v.GetNumber(); e != nil {
					err = e
				} else {
					ret = []float64{one}
				}
			case affine.NumList:
				if many, e := v.GetNumList(); e != nil {
					err = e
				} else {
					ret = many
				}
			default:
				err = errutil.New("cant add %s to a num list", a)
			}
		}
	}
	return
}

func (op *Splice) getNewStrings(run rt.Runtime) (ret []string, err error) {
	if assign := op.Append; assign != nil {
		if v, e := assign.GetAssignedValue(run); e != nil {
			err = e
		} else {
			switch a := v.Affinity(); a {
			case affine.Text:
				if one, e := v.GetText(); e != nil {
					err = e
				} else {
					ret = []string{one}
				}
			case affine.TextList:
				if many, e := v.GetTextList(); e != nil {
					err = e
				} else {
					ret = many
				}
			default:
				err = errutil.New("cant add %s to a text list", a)
			}
		}
	}
	return
}

// reti is < 0 to indicate an empty list
func (op *Splice) getIndices(run rt.Runtime, cnt int) (reti, retj int, err error) {
	if i, e := rt.GetOptionalNumber(run, op.Start, 0); e != nil {
		err = e
	} else if rng, e := rt.GetOptionalNumber(run, op.Delete, 0); e != nil {
		err = e
	} else {
		reti = clipStart(int(i), cnt)
		retj = clipRange(reti, int(rng), cnt)
	}
	return
}
