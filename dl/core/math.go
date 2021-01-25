package core

import (
	"math"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

func getPair(run rt.Runtime, a, b rt.NumberEval) (reta, retb float64, err error) {
	if a, e := safe.GetNumber(run, a); e != nil {
		err = errutil.New("couldnt get first operand, because", e)
	} else if b, e := safe.GetNumber(run, b); e != nil {
		err = errutil.New("couldnt get second operand, because", e)
	} else {
		reta, retb = a.Float(), b.Float()
	}
	return
}

type SumOf struct{ A, B rt.NumberEval }
type DiffOf struct{ A, B rt.NumberEval }
type ProductOf struct{ A, B rt.NumberEval }
type QuotientOf struct{ A, B rt.NumberEval }
type RemainderOf struct{ A, B rt.NumberEval }

func (*SumOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "math",
		Desc:  "Add Numbers: Add two numbers.",
		Spec:  "( {a:number_eval} + {b:number_eval} )",
	}
}

func (op *SumOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a + b)
	}
	return
}

func (*DiffOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "math",
		Spec:  "( {a:number_eval} - {b:number_eval} )",
		Desc:  "Subtract Numbers: Subtract two numbers.",
	}
}

func (op *DiffOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a - b)
	}
	return
}

func (*ProductOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "math",
		Spec:  "( {a:number_eval} * {b:number_eval} )",
		Desc:  "Multiply Numbers: Multiply two numbers.",
	}
}

func (op *ProductOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a * b)
	}
	return
}

func (*QuotientOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "math",
		Spec:  "( {a:number_eval} / {b:number_eval} )",
		Desc:  "Divide Numbers: Divide one number by another.",
	}
}

func (op *QuotientOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else if math.Abs(b) <= 1e-5 {
		e := errutil.New("QuotientOf second operand is too small")
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a / b)
	}
	return
}

func (*RemainderOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "math",
		Spec:  "( {a:number_eval} % {b:number_eval} )",
		Desc:  "Modulus Numbers: Divide one number by another, and return the remainder.",
	}
}

func (op *RemainderOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		mod := math.Mod(a, b)
		ret = g.FloatOf(mod)
	}
	return
}
