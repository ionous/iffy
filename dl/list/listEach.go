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

type Each struct {
	List     string // variable name
	Go, Else *core.Activity
}

func (op *Each) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_each",
		Group:  "list",
		Desc:   `For each in list: Loops over the elements in the passed list, or runs the 'else' activity if empty.`,
		Locals: []string{"index", "first", "last", "text", "num"},
	}
}

func (op *Each) Execute(run rt.Runtime) (err error) {
	if e := op.execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Each) execute(run rt.Runtime) (err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else if cnt, e := vs.GetLen(); e != nil {
		err = e
	} else if otherwise := op.Else; otherwise != nil && cnt == 0 {
		err = op.Else.Execute(run)
	} else if act := op.Go; act != nil && cnt > 0 {
		//
		var field string
		var zero rt.Value
		//
		switch a := vs.Affinity(); a {
		case affine.NumList:
			field = "num"
			zero = generic.Zero
		case affine.TextList:
			field = "text"
			zero = generic.Empty
		default:
			err = errutil.Fmt("variable '%s(%s)' is an unknown list", op.List, a)
		}
		if err == nil {
			var terms term.Terms
			el := terms.AddTerm(field, zero)
			index := terms.AddTerm("index", generic.Zero)
			first := terms.AddTerm("first", generic.True)
			last := terms.AddTerm("last", generic.False)
			run.PushScope(&terms)
			for i := 0; i < cnt; i++ {
				if at, e := vs.GetIndex(i); e != nil {
					err = e
					break
				} else {
					el.SetValue(at)
					next := i + 1
					index.SetValue(&generic.Int{Value: next})
					if hasNext := next < cnt; !hasNext {
						last.SetValue(generic.True)
					}
					if e := op.Go.Execute(run); e != nil {
						err = e
						break
					}
					first.SetValue(generic.False)
				}
			}
			run.PopScope()
		}
	}
	return
}
