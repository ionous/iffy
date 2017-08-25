package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

func (ops *Ops) Execute(fn func(c *Builder)) (ret rt.Execute, err error) {
	var root struct{ Eval rt.Execute }
	if c, ok := ops.NewBuilder(&root); !ok {
		err = errutil.New("unknown error")
	} else {
		fn(c)
		if e := c.Build(); e != nil {
			err = e
		} else {
			ret = root.Eval
		}
	}
	return
}

func (ops *Ops) ObjectEval(fn func(c *Builder)) (ret rt.ObjectEval, err error) {
	var root struct{ Obj rt.ObjectEval }
	if c, ok := ops.NewBuilder(&root); !ok {
		err = errutil.New("unknown error")
	} else {
		fn(c)
		if e := c.Build(); e != nil {
			err = e
		} else {
			ret = root.Obj
		}
	}
	return
}
