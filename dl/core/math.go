package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"math"
)

type Pair struct {
	A, B rt.NumberEval
}

func (cmd Pair) Get(run rt.Runtime) (reta, retb float64, err error) {
	if a, e := cmd.A.GetNumber(run); e != nil {
		err = errutil.New("a", e)
	} else if b, e := cmd.B.GetNumber(run); e != nil {
		err = errutil.New("b", e)
	} else {
		reta, retb = a, b
	}
	return
}

type Add Pair
type Sub Pair
type Mul Pair
type Div Pair

func (cmd *Add) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := Pair(*cmd).Get(run); e != nil {
		err = errutil.New("Add", e)
	} else {
		ret = a + b
	}
	return
}

func (cmd *Sub) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := Pair(*cmd).Get(run); e != nil {
		err = errutil.New("Sub", e)
	} else {
		ret = a - b
	}
	return
}

func (cmd *Mul) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := Pair(*cmd).Get(run); e != nil {
		err = errutil.New("Mul", e)
	} else {
		ret = a * b
	}
	return
}

func (cmd *Div) GetNumber(run rt.Runtime) (ret float64, err error) {
	if a, b, e := Pair(*cmd).Get(run); e != nil {
		err = errutil.New("Div", e)
	} else if math.Abs(b) <= 1e-5 {
		err = errutil.New("Div", "b is too small")
	} else {
		ret = a / b
	}
	return
}
