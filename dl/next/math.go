package next

import (
	"math"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

func getPair(run rt.Runtime, a, b rt.NumberEval) (reta, retb float64, err error) {
	if a, e := a.GetNumber(run); e != nil {
		err = errutil.New("couldnt get first operand, because", e)
	} else if b, e := b.GetNumber(run); e != nil {
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

func (cmd *SumOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("Add", e)
	} else {
		ret = a + b
	}
	return
}

func (cmd *DiffOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("Sub", e)
	} else {
		ret = a - b
	}
	return
}

func (cmd *ProductOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("Mul", e)
	} else {
		ret = a * b
	}
	return
}

func (cmd *QuotientOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("Div", e)
	} else if math.Abs(b) <= 1e-5 {
		err = errutil.New("Div second operand is too small")
	} else {
		ret = a / b
	}
	return
}

func (cmd *RemainderOf) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := getPair(run, cmd.A, cmd.B); e != nil {
		err = errutil.New("Mod", e)
	} else {
		ret = math.Mod(a, b)
	}
	return
}
