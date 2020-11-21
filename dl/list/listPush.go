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
)

type Push struct {
	List   string // variable name
	Insert core.Assignment
	Front  FrontOrBack
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
func (op *Push) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if cnt, e := op.push(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(cnt)
	}
	return
}

func (op *Push) push(run rt.Runtime) (ret int, err error) {
	if vs, e := safe.GetList(run, op.List); e != nil {
		err = e
	} else if res, e := op.pushList(run, vs); e != nil {
		err = e
	} else if e := run.SetField(object.Variables, op.List, res); e != nil {
		err = e
	} else {
		ret = res.Len()
	}
	return
}

func (op *Push) pushList(run rt.Runtime, vs g.Value) (ret g.Value, err error) {
	switch a := vs.Affinity(); a {
	case affine.NumList:
		if add, e := getNewFloats(run, op.Insert); e != nil {
			err = e
		} else if res, e := pushNumbers(vs, add, bool(op.Front)); e != nil {
			err = e
		} else {
			ret = g.FloatsOf(res)
		}
	case affine.TextList:
		if add, e := getNewStrings(run, op.Insert); e != nil {
			err = e
		} else if res, e := pushText(vs, add, bool(op.Front)); e != nil {
			err = e
		} else {
			ret = g.StringsOf(res)
		}
	case affine.RecordList:
		t := vs.Type()
		if add, e := getNewRecords(run, t, op.Insert); e != nil {
			err = e
		} else if res, e := pushRecords(vs, add, bool(op.Front)); e != nil {
			err = e
		} else {
			ret, err = g.ValueFrom(res, a, t)
		}
	default:
		err = errutil.Fmt("variable '%s(%s)' isn't a list", op.List, a)
	}
	return
}

func pushNumbers(v g.Value, add []float64, front bool) (ret []float64, err error) {
	vs := v.Floats()
	if front {
		ret = append(add, vs...)
	} else {
		ret = append(vs, add...)
	}
	return
}

func pushText(v g.Value, add []string, front bool) (ret []string, err error) {
	vs := v.Strings()
	if front {
		ret = append(add, vs...)
	} else {
		ret = append(vs, add...)
	}
	return
}

func pushRecords(v g.Value, add []*g.Record, front bool) (ret []*g.Record, err error) {
	vs := v.Records()
	if front {
		ret = append(add, vs...)
	} else {
		ret = append(vs, add...)
	}
	return
}
