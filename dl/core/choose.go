package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// Choose to execute one of two blocks based on a boolean test.
type Choose struct {
	If          rt.BoolEval
	True, False []rt.Execute
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

func (*Choose) Compose() composer.Spec {
	return composer.Spec{
		Name: "choose",
		Spec: "if {choose%if:bool_eval} then: {true?execute|ghost} else: {false?execute|ghost}",
	}
}

// Execute evals, eats the returns
func (op *Choose) Execute(run rt.Runtime) (err error) {
	if b, e := rt.GetBool(run, op.If); e != nil {
		err = e
	} else {
		var next []rt.Execute
		if b {
			next = op.True
		} else {
			next = op.False
		}
		err = rt.Run(run, next)
	}
	return
}

func (*ChooseNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "choose_num",
		Group: "math",
		Desc:  "Choose Number: Pick one of two numbers based on a boolean test.",
	}
}

func (op *ChooseNum) GetNumber(run rt.Runtime) (ret float64, err error) {
	if b, e := rt.GetBool(run, op.If); e != nil {
		err = e
	} else {
		var next rt.NumberEval
		if b {
			next = op.True
		} else {
			next = op.False
		}
		if next != nil {
			ret, err = rt.GetNumber(run, next)
		}
	}
	return
}

func (*ChooseText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "choose_text",
		Group: "format",
		Desc:  "Choose Text: Pick one of two strings based on a boolean test.",
	}
}

func (op *ChooseText) GetText(run rt.Runtime) (ret string, err error) {
	if b, e := rt.GetBool(run, op.If); e != nil {
		err = e
	} else {
		var next rt.TextEval
		if b {
			next = op.True
		} else {
			next = op.False
		}
		if next != nil {
			ret, err = rt.GetText(run, next)
		}
	}
	return
}
