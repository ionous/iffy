package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type Push struct {
	List   string // variable name
	Front  FrontOrBack
	Insert core.Assignment
}

func (op *Push) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_push",
		Group: "list",
		Spec:  "push {list:text} {front} {inserting%insert?assignment}",
		Desc: `Push into list: Add elements to the front or back of a list.
Returns the new length of the list.`,
	}
}

func (op *Push) Execute(run rt.Runtime) (err error) {
	if _, e := op.push(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

// returns the new size
func (op *Push) GetNumber(run rt.Runtime) (ret float64, err error) {
	if cnt, e := op.push(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = float64(cnt)
	}
	return
}

func (op *Push) push(run rt.Runtime) (ret int, err error) {
	if orig, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else if res, e := op.pushList(run, orig); e != nil {
		err = e
	} else if cnt, e := res.GetLen(); e != nil {
		err = e
	} else if e := run.SetField(object.Variables, op.List, res); e != nil {
		err = e
	} else {
		ret = cnt
	}
	return
}

func (op *Push) pushList(run rt.Runtime, orig rt.Value) (ret rt.Value, err error) {
	switch a := orig.Affinity(); a {
	case affine.NumList:
		if add, e := getNewFloats(run, op.Insert); e != nil {
			err = e
		} else if res, e := pushNumbers(orig, add, bool(op.Front)); e != nil {
			err = e
		} else {
			ret  = generic.FloatsOf(res)
		}
	case affine.TextList:
		if add, e := getNewStrings(run, op.Insert); e != nil {
			err = e
		} else if res, e := pushText(orig, add, bool(op.Front)); e != nil {
			err = e
		} else {
			ret = generic.StringsOf(res)
		}
	default:
		err = errutil.Fmt("variable '%s(%s)' isn't a list", op.List, a)
	}
	return
}

func pushNumbers(orig rt.Value, add []float64, front bool) (ret []float64, err error) {
	if orig, e := orig.GetNumList(); e != nil {
		err = e
	} else if front {
		ret = append(add, orig...)
	} else {
		ret = append(orig, add...)
	}
	return
}

func pushText(orig rt.Value, add []string, front bool) (ret []string, err error) {
	if orig, e := orig.GetTextList(); e != nil {
		err = e
	} else if front {
		ret = append(add, orig...)
	} else {
		ret = append(orig, add...)
	}
	return
}
