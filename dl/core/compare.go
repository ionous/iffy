package core

import (
	"strings"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type CompareNum struct {
	A  rt.NumberEval `if:"selector=num"`
	Is Comparator    `if:"selector,compact"`
	B  rt.NumberEval `if:"selector"`
	// fix: add optional epsilon?
}

type CompareText struct {
	A  rt.TextEval `if:"selector=txt"`
	Is Comparator  `if:"selector,compact"`
	B  rt.TextEval `if:"selector"`
}

func (*CompareNum) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "is", Role: composer.Function},
		Group:  "logic",
		Desc:   "Compare Numbers: True if eq,ne,gt,lt,ge,le two numbers.",
	}
}

func (op *CompareNum) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if src, e := safe.GetNumber(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if tgt, e := safe.GetNumber(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else if is := op.Is; is == nil {
		err = cmdErrorCtx(op, "comparator is nil", nil)
	} else {
		res := compare(is, src.Float()-tgt.Float(), 1e-3)
		ret = g.BoolOf(res)
	}
	return
}

func (*CompareText) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "is", Role: composer.Function},
		Group:  "logic",
		Desc:   "Compare Text: True if eq,ne,gt,lt,ge,le two strings ( lexical. )",
	}
}

func (op *CompareText) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if src, e := safe.GetText(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if tgt, e := safe.GetText(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else if is := op.Is; is == nil {
		err = cmdErrorCtx(op, "comparator is nil", nil)
	} else {
		c := strings.Compare(src.String(), tgt.String())
		res := compare(is, float64(c), 0.5)
		ret = g.BoolOf(res)
	}
	return
}

func compare(is Comparator, d, epsilon float64) (ret bool) {
	switch cmp := is.Compare(); {
	default:
		ret = (cmp & Compare_EqualTo) != 0
	case d < -epsilon:
		ret = (cmp & Compare_LessThan) != 0
	case d > epsilon:
		ret = (cmp & Compare_GreaterThan) != 0
	}
	return
}
