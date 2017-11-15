package core

import (
	"github.com/ionous/iffy/rt"
)

// Choose to execute one of two blocks based on a boolean test.
type Choose struct {
	If    rt.BoolEval
	True  rt.ExecuteList
	False rt.ExecuteList
}

// Choose one of two number evaluations based on a boolean test.
type ChooseNum struct {
	If          rt.BoolEval
	True, False rt.NumberEval
}

// Choose one of two text evaluations based on a boolean test.
type ChooseText struct {
	If          rt.BoolEval
	True, False rt.TextEval
}

// Choose one of two object evaluations based on a boolean test.
type ChooseObj struct {
	If          rt.BoolEval
	True, False rt.ObjectEval
}

func (op *ChooseNum) GetNumber(run rt.Runtime) (ret float64, err error) {
	if b, e := op.If.GetBool(run); e != nil {
		err = e
	} else {
		var next rt.NumberEval
		if b {
			next = op.True
		} else {
			next = op.False
		}
		if next != nil {
			ret, err = next.GetNumber(run)
		}
	}
	return
}

func (op *ChooseText) GetText(run rt.Runtime) (ret string, err error) {
	if b, e := op.If.GetBool(run); e != nil {
		err = e
	} else {
		var next rt.TextEval
		if b {
			next = op.True
		} else {
			next = op.False
		}
		if next != nil {
			ret, err = next.GetText(run)
		}
	}
	return
}

func (op *ChooseObj) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if b, e := op.If.GetBool(run); e != nil {
		err = e
	} else {
		var next rt.ObjectEval
		if b {
			next = op.True
		} else {
			next = op.False
		}
		if next != nil {
			ret, err = next.GetObject(run)
		}
	}
	return
}

// Execute evals, eats the returns
func (op *Choose) Execute(run rt.Runtime) (err error) {
	if b, e := op.If.GetBool(run); e != nil {
		err = e
	} else {
		var next rt.ExecuteList
		if b {
			next = op.True
		} else {
			next = op.False
		}
		err = next.Execute(run)
	}
	return
}
