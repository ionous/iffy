package next

import (
	"io"

	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// Choose to execute one of two blocks based on a boolean test.
type Choose struct {
	If    rt.BoolEval
	True  rt.Block
	False rt.Block
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
	True, False rt.TextWriter
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

func (op *ChooseText) WriteText(run rt.Runtime, w io.Writer) (err error) {
	if b, e := op.If.GetBool(run); e != nil {
		err = e
	} else {
		var next rt.TextWriter
		if b {
			next = op.True
		} else {
			next = op.False
		}
		if next != nil {
			err = next.WriteText(run, w)
		}
	}
	return
}

func (op *ChooseObj) GetObject(run rt.Runtime) (ret string, err error) {
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
		var next rt.Block
		if b {
			next = op.True
		} else {
			next = op.False
		}
		err = next.Execute(run)
	}
	return
}
