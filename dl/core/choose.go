package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

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

func (*ChooseNum) Compose() composer.Spec {
	return composer.Spec{
		Group: "math",
		Desc:  "Choose Number: Pick one of two numbers based on a boolean test.",
	}
}

func (op *ChooseNum) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if b, e := safe.GetBool(run, op.If); e != nil {
		err = cmdError(op, e)
	} else {
		var next rt.NumberEval
		if b.Bool() {
			next = op.True
		} else {
			next = op.False
		}
		if v, e := safe.GetOptionalNumber(run, next, 0); e != nil {
			err = cmdError(op, e)
		} else {
			ret = v
		}
	}
	return
}

func (*ChooseText) Compose() composer.Spec {
	return composer.Spec{
		Group: "format",
		Desc:  "Choose Text: Pick one of two strings based on a boolean test.",
	}
}

func (op *ChooseText) GetText(run rt.Runtime) (ret g.Value, err error) {
	if b, e := safe.GetBool(run, op.If); e != nil {
		err = cmdError(op, e)
	} else {
		var next rt.TextEval
		if b.Bool() {
			next = op.True
		} else {
			next = op.False
		}
		if v, e := safe.GetOptionalText(run, next, ""); e != nil {
			err = cmdError(op, e)
		} else {
			ret = v
		}
	}
	return
}
