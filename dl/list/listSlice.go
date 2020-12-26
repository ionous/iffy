package list

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Slice struct {
	List       core.Assignment
	Start, End rt.NumberEval // from start to end (end not included)
}

// Start is optional, if omitted slice starts at the first element.
// If start is greater the length, an empty array is returned.

// Slice doesnt include the ending index.
// Negatives indices indicates an offset from the end.

// When end is omitted, copy up to and including the last element;
// and do the same if the end is greater than the length

func (*Slice) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_slice",
		Group: "list",
		Spec:  "slice {list:assignment} {from entry%start?number} {ending before entry%end?number}",
		Desc:  "Slice of List: Create a new list from a section of another list.",
	}
}

func (op *Slice) Execute(run rt.Runtime) (err error) {
	if _, _, e := op.sliceList(run, ""); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Slice) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.sliceList(run, affine.NumList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.FloatsOf(nil)
	} else {
		ret = v
	}
	return
}

func (op *Slice) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.sliceList(run, affine.TextList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.StringsOf(nil)
	} else {
		ret = v
	}
	return
}

func (op *Slice) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	if v, t, e := op.sliceList(run, affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.RecordsOf(t, nil)
	} else {
		ret = v
	}
	return
}

func (op *Slice) sliceList(run rt.Runtime, aff affine.Affinity) (retVal g.Value, retType string, err error) {
	if els, e := core.GetAssignedValue(run, op.List); e != nil {
		err = e
	} else if e := safe.Check(els, aff); e != nil {
		err = e
	} else if i, j, e := op.getIndices(run, els.Len()); e != nil {
		err = e
	} else {
		if i >= 0 && j >= i {
			retVal, err = els.Slice(i, j)
		}
		if err == nil {
			retType = els.Type()
		}
	}
	return
}

// reti is < 0 to indicate an empty list
func (op *Slice) getIndices(run rt.Runtime, cnt int) (reti, retj int, err error) {
	if i, e := safe.GetOptionalNumber(run, op.Start, 0); e != nil {
		err = e
	} else if j, e := safe.GetOptionalNumber(run, op.End, 0); e != nil {
		err = e
	} else {
		reti = clipStart(i.Int(), cnt)
		retj = clipEnd(j.Int(), cnt)
	}
	return
}
