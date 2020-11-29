package list

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/rt/scope"
)

type Pop struct {
	List     string // variable name
	With     string // counter name
	Front    Front
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
	if e := op.pop(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Pop) pop(run rt.Runtime) (err error) {
	if vs, e := safe.List(run, op.List); e != nil {
		err = e
	} else {
		if cnt := vs.Len(); cnt == 0 && op.Else != nil {
			err = op.Else.Execute(run)
		} else {
			const el = 0
			if op.k == nil || op.k.IsStaleKind(run) {
				elAff, elType := affine.Element(vs.Affinity()), vs.Type()
				op.k = g.NewKind(run, "", []g.Field{
					{Name: op.With, Affinity: elAff, Type: elType},
				})
			}
			ls := op.k.NewRecord()
			//
			var at int
			if !op.Front {
				at = cnt - 1
			}
			if popped, e := vs.Splice(at, at+1, nil); e != nil {
				err = e
			} else if e := ls.SetFieldByIndex(el, popped.Index(0)); e != nil {
				err = e
			} else if op.Go != nil {
				run.PushScope(&scope.TargetRecord{object.Variables, ls})
				err = op.Go.Execute(run)
				run.PopScope()
			}
		}
	}
	return
}
