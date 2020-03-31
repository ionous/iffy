package next

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type CompareType int

// CompareTo generates comparison flags.
// FIX: a combo-box of enumeration options should be possible.
type CompareTo interface {
	Compare() CompareType
}

type EqualTo struct{}
type NotEqualTo struct{}
type GreaterThan struct{}
type LessThan struct{}
type GreaterOrEqual struct{}
type LessOrEqual struct{}

func (EqualTo) Compare() CompareType {
	return Compare_EqualTo
}
func (NotEqualTo) Compare() CompareType {
	return 0
}
func (GreaterThan) Compare() CompareType {
	return Compare_GreaterThan
}
func (LessThan) Compare() CompareType {
	return Compare_LessThan
}
func (GreaterOrEqual) Compare() CompareType {
	return Compare_LessThan | Compare_EqualTo
}
func (LessOrEqual) Compare() CompareType {
	return Compare_GreaterThan | Compare_EqualTo
}

//go:generate stringer -type=CompareType
const (
	Compare_EqualTo CompareType = 1 << iota
	Compare_GreaterThan
	Compare_LessThan
)

type CompareNum struct {
	A  rt.NumberEval
	Is CompareTo
	B  rt.NumberEval
}

type CompareText struct {
	A  rt.TextWriter
	Is CompareTo
	B  rt.TextWriter
}

type CompareObj struct {
	A  rt.ObjectEval
	Is CompareTo
	B  rt.ObjectEval
}

func (comp *CompareNum) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := comp.A.GetNumber(run); e != nil {
		err = errutil.New("CompareNum.A", e)
	} else if tgt, e := comp.B.GetNumber(run); e != nil {
		err = errutil.New("CompareNum.B", e)
	} else {
		d := src - tgt
		switch cmp := comp.Is.Compare(); {
		case d == 0:
			ret = (cmp & Compare_EqualTo) != 0
		case d < 0:
			ret = (cmp & Compare_LessThan) != 0
		case d > 0:
			ret = (cmp & Compare_GreaterThan) != 0
		}
	}
	return
}

func (comp *CompareText) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := rt.GetText(run, comp.A); e != nil {
		err = errutil.New("CompareText.A", e)
	} else if tgt, e := rt.GetText(run, comp.B); e != nil {
		err = errutil.New("CompareText.B", e)
	} else {
		ret, err = compareStr(comp.Is.Compare(), src, tgt)
	}
	return
}

func (comp *CompareObj) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := comp.A.GetObject(run); e != nil {
		err = errutil.New("CompareObject.A", e)
	} else if tgt, e := comp.B.GetObject(run); e != nil {
		err = errutil.New("CompareObject.B", e)
	} else {
		ret, err = compareStr(comp.Is.Compare(), src, tgt)
	}
	return
}

func compareStr(cmp CompareType, src, tgt string) (ret bool, err error) {
	switch cmp {
	case Compare_EqualTo:
		ret = src == tgt
	case Compare_LessThan:
		ret = src < tgt
	case Compare_GreaterThan:
		ret = src > tgt
	case Compare_GreaterThan | Compare_EqualTo:
		ret = src >= tgt
	case Compare_LessThan | Compare_EqualTo:
		ret = src <= tgt
	default:
		err = errutil.New("CompareText.Is", cmp, "unknown operand")
	}

	return
}
