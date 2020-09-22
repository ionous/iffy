package core

import (
	"strings"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type CompareNum struct {
	A  rt.NumberEval
	Is Comparator
	B  rt.NumberEval
}

type CompareText struct {
	A  rt.TextEval
	Is Comparator
	B  rt.TextEval
}

func (*CompareNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "compare_num",
		Group: "logic",
		Spec:  "{a:number_eval} {is:comparator} {b:number_eval}",
		Desc:  "Compare Numbers: True if eq,ne,gt,lt,ge,le two numbers.",
	}
}

func (op *CompareNum) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := rt.GetNumber(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if tgt, e := rt.GetNumber(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else if is := op.Is; is == nil {
		err = cmdErrorCtx(op, "comparator is nil", nil)
	} else {
		ret = compare(is, int(src-tgt))
	}
	return
}

func (*CompareText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "compare_text",
		Group: "logic",
		Desc:  "Compare Text: True if eq,ne,gt,lt,ge,le two strings ( lexical. )",
		Spec:  "{a:text_eval} {is:comparator} {b:text_eval}",
	}
}

func (op *CompareText) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := rt.GetText(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if tgt, e := rt.GetText(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else if is := op.Is; is == nil {
		err = cmdErrorCtx(op, "comparator is nil", nil)
	} else {
		ret = compare(is, strings.Compare(src, tgt))
	}
	return
}

func compare(is Comparator, d int) (ret bool) {
	switch cmp := is.Compare(); {
	case d == 0:
		ret = (cmp & Compare_EqualTo) != 0
	case d < 0:
		ret = (cmp & Compare_LessThan) != 0
	case d > 0:
		ret = (cmp & Compare_GreaterThan) != 0
	}
	return
}
