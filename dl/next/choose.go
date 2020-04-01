package next

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// Choose to execute one of two blocks based on a boolean test.
type Choose struct {
	If    rt.BoolEval
	True  rt.Execute
	False rt.Execute
}

func (*Choose) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose",
		Spec: "if {choose%if:bool_eval} then: {true*execute|ghost} else: {false*execute|ghost}",
	}
}

// Choose one of two number evaluations based on a boolean test.
type ChooseNum struct {
	If          rt.BoolEval
	True, False rt.NumberEval
}

// Choose one of two text phrases based on a boolean test.
type ChooseText struct {
	If          rt.BoolEval
	True, False rt.TextEval
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

// Execute evals, eats the returns
func (op *Choose) Execute(run rt.Runtime) (err error) {
	if b, e := op.If.GetBool(run); e != nil {
		err = e
	} else {
		var next rt.Execute
		if b {
			next = op.True
		} else {
			next = op.False
		}
		err = next.Execute(run)
	}
	return
}
