package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
)

type Pop struct {
	List     string // variable name
	With     string // counter name
	Front    FrontOrBack
	Go, Else *core.Activity
	k        *g.Kind
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
	} else if elAffinity := affine.Element(vs.Affinity()); len(elAffinity) == 0 {
		err = errutil.Fmt("Variable %q is %q, pop expected a list", op.List, vs.Affinity())
	} else if cnt, e := vs.GetLen(); e != nil {
		err = e
	} else if cnt == 0 && op.Else != nil {
		err = op.Else.Execute(run)
	} else {
		const el = 0
		if op.k == nil || op.k.IsStaleKind(run) {
			op.k = g.NewKind(run, "", []g.Field{
				{Name: op.With, Affinity: elAffinity, Type: vs.Type()},
			})
		}
		ls := op.k.NewRecord()
		//
		var at, start, end int
		if op.Front {
			at, start, end = 0, 1, cnt
		} else {
			at, start, end = cnt-1, 0, cnt-1
		}
		if popped, e := vs.GetIndex(at); e != nil {
			err = e
		} else if newEls, e := vs.Slice(start, end); e != nil {
			err = e
		} else if e := run.SetField(object.Variables, op.List, newEls); e != nil {
			err = e
		} else if e := ls.SetFieldByIndex(el, popped); e != nil {
			err = e
		} else if op.Go != nil {
			run.PushScope(&scope.TargetRecord{object.Variables, ls})
			err = op.Go.Execute(run)
			run.PopScope()
		}
	}
	return
}
