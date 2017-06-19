package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type CompareType int

type CompareTo interface {
	Compare() CompareType
}

type EqualTo struct{}
type GreaterThan struct{}
type LesserThan struct{}
type NotEqualTo struct{}

func (EqualTo) Compare() CompareType     { return Compare_EqualTo }
func (GreaterThan) Compare() CompareType { return Compare_GreaterThan }
func (LesserThan) Compare() CompareType  { return Compare_LesserThan }
func (NotEqualTo) Compare() CompareType  { return Compare_NotEqualTo }

//go:generate stringer -type=CompareType
const (
	Compare_EqualTo CompareType = 1 << iota
	Compare_GreaterThan
	Compare_LesserThan
	Compare_NotEqualTo = Compare_GreaterThan | Compare_LesserThan
)

type CompareNum struct {
	A  rt.NumberEval
	Is CompareTo
	B  rt.NumberEval
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
			ret = (cmp & Compare_LesserThan) != 0
		case d > 0:
			ret = (cmp & Compare_GreaterThan) != 0
		}
	}
	return
}

// CompareText
type CompareText struct {
	A  rt.TextEval
	Is CompareTo
	B  rt.TextEval
}

func (comp *CompareText) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := comp.A.GetText(run); e != nil {
		err = errutil.New("CompareText.A", e)
	} else if tgt, e := comp.B.GetText(run); e != nil {
		err = errutil.New("CompareText.B", e)
	} else {
		switch cmp := comp.Is.Compare(); cmp {
		case Compare_EqualTo:
			ret = src == tgt
		case Compare_NotEqualTo:
			ret = src != tgt
		case Compare_LesserThan:
			ret = src < tgt
		case Compare_GreaterThan:
			ret = src > tgt
		case Compare_GreaterThan | Compare_EqualTo:
			ret = src >= tgt
		case Compare_LesserThan | Compare_EqualTo:
			ret = src <= tgt
		default:
			err = errutil.New("CompareText.Is", cmp, "unknown operand")
		}
	}
	return
}
