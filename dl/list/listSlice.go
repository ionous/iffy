package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

type Slice struct {
	List       string        // variable name
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
		Spec:  "slice {list:text} {from entry%start?number} {ending before entry%eend?number}",
		Desc:  "Slice of List: Create a new list from a section of another list.",
	}
}

func (op *Slice) GetNumList(run rt.Runtime) (ret []float64, err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = cmdError(op, e)
	} else if els, e := vs.GetNumList(); e != nil {
		err = cmdError(op, e)
	} else if i, j, e := op.getIndices(run, len(els)); e != nil {
		err = cmdError(op, e)
	} else if i >= 0 && i <= j {
		ret = els[i:j]
	}
	return
}

func (op *Slice) GetTextList(run rt.Runtime) (ret []string, err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = cmdError(op, e)
	} else if els, e := vs.GetTextList(); e != nil {
		err = cmdError(op, e)
	} else if i, j, e := op.getIndices(run, len(els)); e != nil {
		err = cmdError(op, e)
	} else if i >= 0 && i <= j {
		ret = els[i:j]
	}
	return
}

// reti is < 0 to indicate an empty list
func (op *Slice) getIndices(run rt.Runtime, cnt int) (reti, retj int, err error) {
	if i, e := rt.GetNumber(run, op.Start); e != nil {
		err = e
	} else if j, e := rt.GetNumber(run, op.End); e != nil {
		err = e
	} else {
		i, j := int(i), int(j)
		if i == 0 {
			reti = 0 // unspecified: start at the front of the list
		} else if i > cnt {
			reti = -1 // negative return means an empty list
		} else if i > 0 {
			reti = i - 1 // one based indicies
		} else if ofs := cnt + i; ofs > 0 {
			// offset from the end: slice(-2) extracts the last two elements in the sequence.
			reti = ofs
		} else {
			reti = 0
		}

		if j > cnt {
			retj = cnt
		} else if j > 0 {
			retj = j - 1
		} else if ofs := cnt + j; ofs > 0 {
			retj = ofs
		} else {
			reti = -1 // negative return means an empty list
		}
	}
	return
}
