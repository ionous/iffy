package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Splice struct {
	List          string        // variable name
	Start, Remove rt.NumberEval // from start
	Insert        core.Assignment
}

// if start is negative, it will begin that many elements from the end of the array.
// If array.length + start is less than 0, it will begin from index 0.
// If deleteCount is 0 or negative, no elements are removed.
func (*Splice) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_splice",
		Group: "list",
		Spec:  "splice into {list:text} {at entry%start?number} {removing%remove?number} {inserting%insert?assignment}",
		Desc: `Splice into list: Modify a list by adding and removing elements.
Note: the type of the elements being added must match the type of the list. 
Text cant be added to a list of numbers, numbers cant be added to a list of text, 
and true/false values can't be added to a list.`,
	}
}

func (op *Splice) Execute(run rt.Runtime) (err error) {
	if vs, e := safe.GetList(run, op.List); e != nil {
		err = cmdError(op, e)
	} else {
		switch a := vs.Affinity(); a {
		case affine.NumList:
			_, err = op.spliceNumbers(run, vs)
		case affine.TextList:
			_, err = op.spliceText(run, vs)
		default:
			err = errutil.Fmt("variable '%s(%s)' isn't a list", op.List, a)
		}
	}
	return
}

// returns the removed elements
func (op *Splice) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := safe.GetList(run, op.List); e != nil {
		err = cmdError(op, e)
	} else if vals, e := op.spliceNumbers(run, vs); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatsOf(vals)
	}
	return
}

// returns the removed elements
func (op *Splice) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := safe.GetList(run, op.List); e != nil {
		err = cmdError(op, e)
	} else if vals, e := op.spliceText(run, vs); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.StringsOf(vals)
	}
	return
}

func (op *Splice) spliceNumbers(run rt.Runtime, vs g.Value) (ret []float64, err error) {
	els := vs.Floats()
	if i, j, e := op.getIndices(run, len(els)); e != nil {
		err = e
	} else if add, e := getNewFloats(run, op.Insert); e != nil {
		err = e
	} else {
		if i >= 0 && j >= i {
			// cut out i to j, then i becomes the insertion point.
			// as long as the range was valid we take the result and set it back...
			// even if no elements are cut or inserted.
			// ( that bakes any evaluation that might have been in the target )
			out := g.CopyFloats(els[i:j]) // before we start altering the memory of els, copy out
			newVals := append(els[:i], append(add, els[j:]...)...)
			if e := run.SetField(object.Variables, op.List, g.FloatsOf(newVals)); e != nil {
				err = e
			} else {
				ret = out
			}
		}
	}
	return
}

// returns the removed elements
func (op *Splice) spliceText(run rt.Runtime, vs g.Value) (ret []string, err error) {
	els := vs.Strings()
	if i, j, e := op.getIndices(run, len(els)); e != nil {
		err = e
	} else if add, e := getNewStrings(run, op.Insert); e != nil {
		err = e
	} else {
		// ... mirrors GetNumList()
		if i >= 0 && j >= i {
			out := g.CopyStrings(els[i:j])
			newVals := append(els[:i], append(add, els[j:]...)...)
			if e := run.SetField(object.Variables, op.List, g.StringsOf(newVals)); e != nil {
				err = e
			} else {
				ret = out
			}
		}
	}
	return
}

// reti is < 0 to indicate an empty list
func (op *Splice) getIndices(run rt.Runtime, cnt int) (reti, retj int, err error) {
	if i, e := safe.GetOptionalNumber(run, op.Start, 0); e != nil {
		err = e
	} else if rng, e := safe.GetOptionalNumber(run, op.Remove, 0); e != nil {
		err = e
	} else {
		reti = clipStart(i.Int(), cnt)
		retj = clipRange(reti, rng.Int(), cnt)
	}
	return
}
