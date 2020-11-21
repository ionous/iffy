package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/rt/scope"
)

type Each struct {
	List     string // variable name
	With     string // counter name
	Go, Else *core.Activity
	k        *g.Kind
}

func (op *Each) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_each",
		Group:  "list",
		Desc:   `For each in list: Loops over the elements in the passed list, or runs the 'else' activity if empty.`,
		Locals: []string{"index", "first", "last"},
	}
}

func (op *Each) Execute(run rt.Runtime) (err error) {
	if e := op.execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Each) execute(run rt.Runtime) (err error) {

	if vs, e := safe.GetList(run, op.List); e != nil {
		err = e
	} else if elAffinity := affine.Element(vs.Affinity()); len(elAffinity) == 0 {
		err = errutil.Fmt("Variable %q is %q, each expected a list", op.List, vs.Affinity())
	} else {
		cnt := vs.Len()
		if otherwise := op.Else; otherwise != nil && cnt == 0 {
			err = op.Else.Execute(run)
		} else if act := op.Go; act != nil && cnt > 0 {
			const el, index, first, last = 0, 1, 2, 3
			if op.k == nil || op.k.IsStaleKind(run) {
				op.k = g.NewKind(run, "", []g.Field{
					{Name: op.With, Affinity: elAffinity, Type: vs.Type()},
					{Name: "index", Affinity: affine.Number},
					{Name: "first", Affinity: affine.Bool},
					{Name: "last", Affinity: affine.Bool},
				})
			}
			ls := op.k.NewRecord()
			run.PushScope(&scope.TargetRecord{object.Variables, ls})
			for i := 0; i < cnt; i++ {
				at := vs.Index(i)
				if e := ls.SetFieldByIndex(el, at); e != nil {
					err = e
					break
				} else if e := ls.SetFieldByIndex(index, g.IntOf(i+1)); e != nil {
					err = e
					break
				} else if e := ls.SetFieldByIndex(first, g.BoolOf(i == 0)); e != nil {
					err = e
					break
				} else if e := ls.SetFieldByIndex(last, g.BoolOf((i+1) == cnt)); e != nil {
					err = e
					break
				} else if e := op.Go.Execute(run); e != nil {
					err = e
					break
				}
			}
			run.PopScope()
		}
	}
	return
}
