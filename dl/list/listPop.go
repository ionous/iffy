package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type Pop struct {
	List     string // variable name
	Front    FrontOrBack
	Go, Else *core.Activity
}

func (op *Pop) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_pop",
		Group: "list",
		Desc: `Pop from list: Remove an element from the front or back of a list.
Runs an activity with the popped value, or runs the 'else' activity if the list was empty.`,
		Locals: []string{"text", "num"},
	}
}

func (op *Pop) Execute(run rt.Runtime) (err error) {
	if e := op.execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Pop) execute(run rt.Runtime) (err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else if cnt, e := vs.GetLen(); e != nil {
		err = e
	} else if cnt == 0 && op.Else != nil {
		err = op.Else.Execute(run)
	} else if cnt > 0 {
		var terms term.Terms
		switch a := vs.Affinity(); a {
		case affine.NumList:
			err = op.popNumbers(run, &terms, vs)
		case affine.TextList:
			err = op.popText(run, &terms, vs)
		default:
			err = errutil.Fmt("variable '%s(%s)' is an unknown list", op.List, a)
		}
		if err == nil && op.Go != nil {
			run.PushScope(&terms)
			err = op.Go.Execute(run)
			run.PopScope()
		}
	}
	return
}

func (op *Pop) popNumbers(run rt.Runtime, terms *term.Terms, vs rt.Value) (err error) {
	if els, e := vs.GetNumList(); e != nil {
		err = e
	} else {
		var remove float64
		var remain []float64
		if op.Front {
			remove, remain = els[0], els[1:]
		} else {
			last := len(els) - 1
			remove, remain = els[last], els[:last]
		}
		if e := run.SetField(object.Variables, op.List,
			&generic.FloatSlice{Values: remain}); e != nil {
			err = e
		} else {
			terms.AddTerm("num", &generic.Float{Value: remove})
		}
	}
	return
}

func (op *Pop) popText(run rt.Runtime, terms *term.Terms, vs rt.Value) (err error) {
	if els, e := vs.GetTextList(); e != nil {
		err = e
	} else {
		var remove string
		var remain []string
		if op.Front {
			remove, remain = els[0], els[1:]
		} else {
			last := len(els) - 1
			remove, remain = els[last], els[:last]
		}
		if e := run.SetField(object.Variables, op.List,
			&generic.StringSlice{Values: remain}); e != nil {
			err = e
		} else {
			terms.AddTerm("text", &generic.String{Value: remove})
		}
	}
	return
}
