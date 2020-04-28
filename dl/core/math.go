package core

import (
	"math"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

func getPair(run rt.Runtime, a, b rt.NumberEval) (reta, retb float64, err error) {
	if a, e := rt.GetNumber(run, a); e != nil {
		err = errutil.New("couldnt get first operand, because", e)
	} else if b, e := rt.GetNumber(run, b); e != nil {
		err = errutil.New("couldnt get second operand, because", e)
	} else {
		reta, retb = a, b
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
		Name:  "sum_of",
		Group: "math",
		Desc:  "Add Numbers: Add two numbers.",
		Spec:  "( {a:number_eval} + {b:number_eval} )",
	}
}

func (cmd *SumOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("SumOf", e)
	} else {
		ret = a + b
	}
	return
}

func (*DiffOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "diff_of",
		Group: "math",
		Spec:  "( {a:number_eval} - {b:number_eval} )",
		Desc:  "Subtract Numbers: Subtract two numbers.",
	}
}

func (cmd *DiffOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("DiffOf", e)
	} else {
		ret = a - b
	}
	return
}

func (*ProductOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "product_of",
		Group: "math",
		Spec:  "( {a:number_eval} * {b:number_eval} )",
		Desc:  "Multiply Numbers: Multiply two numbers.",
	}
}

func (cmd *ProductOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("ProductOf", e)
	} else {
		ret = a * b
	}
	return
}

func (*QuotientOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "quotient_of",
		Group: "math",
		Spec:  "( {a:number_eval} / {b:number_eval} )",
		Desc:  "Divide Numbers: Divide one number by another.",
	}
}

func (cmd *QuotientOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("QuotientOf", e)
	} else if math.Abs(b) <= 1e-5 {
		err = errutil.New("QuotientOf second operand is too small")
	} else {
		ret = a / b
	}
	return
}

func (*RemainderOf) Compose() composer.Spec {
	return composer.Spec{
		Name:  "remainder_of",
		Group: "math",
		Spec:  "( {a:number_eval} % {b:number_eval} )",
		Desc:  "Modulus Numbers: Divide one number by another, and return the remainder.",
	}
}

func (cmd *RemainderOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("RemainderOf", e)
	} else {
		ret = math.Mod(a, b)
	}
	return
}
