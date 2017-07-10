package evtbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
)

type BuildOps func(c *ops.Builder)

func (b BuildOps) Build(ops *ops.Ops) (ret rt.Execute, err error) {
	var root struct{ Eval rt.Execute }
	if c, ok := ops.NewBuilder(&root); !ok {
		err = errutil.New("unknown error")
	} else {
		b(c)
		if _, e := c.Build(); e != nil {
			err = e
		} else {
			ret = root.Eval
		}
	}
	return
}

func (b BuildOps) Eval(ops *ops.Ops) (ret rt.ObjectEval, err error) {
	var root struct{ Obj rt.ObjectEval }
	if c, ok := ops.NewBuilder(&root); !ok {
		err = errutil.New("unknown error")
	} else {
		b(c)
		if _, e := c.Build(); e != nil {
			err = e
		} else {
			ret = root.Obj
		}
	}
	return
}
