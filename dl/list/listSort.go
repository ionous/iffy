package list

import (
	"sort"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type Sort struct {
	List    string // variable name
	Pattern string // where the pattern should take a a pair of list elements, and return true if the first is less than the second
}

func (op *Sort) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_sort",
		Group: "list",
		Desc:  `Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.`,
	}
}

func (op *Sort) Execute(run rt.Runtime) (err error) {
	if e := op.sort(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Sort) sort(run rt.Runtime) (err error) {
	if vs, e := safe.List(run, op.List); e != nil {
		err = e
	} else {
		var newVals g.Value
		switch a := vs.Affinity(); a {
		case affine.NumList:
			els := vs.Floats()
			if e := op.sortNumbers(run, els); e != nil {
				err = e
			} else {
				newVals = g.FloatsOf(els)
			}
		case affine.TextList:
			els := vs.Strings()
			if e := op.sortText(run, els); e != nil {
				err = e
			} else {
				newVals = g.StringsOf(els)
			}
		default:
			err = errutil.Fmt("variable '%s(%s)' isn't a list", op.List, a)
		}
		if err == nil {
			err = run.SetField(object.Variables, op.List, newVals)
		}
	}
	return
}

func (op *Sort) sortNumbers(run rt.Runtime, src []float64) (err error) {
	var one, two core.FromValue
	det := makeDet(op.Pattern, &one, &two)
	sort.Slice(src, func(i, j int) (ret bool) {
		one.Value, two.Value = g.FloatOf(src[i]), g.FloatOf(src[j])
		if x, e := det.GetBool(run); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = x.Bool()
		}
		return
	})
	return
}

func (op *Sort) sortText(run rt.Runtime, src []string) (err error) {
	var one, two core.FromValue
	det := makeDet(op.Pattern, &one, &two)
	sort.Slice(src, func(i, j int) (ret bool) {
		one.Value, two.Value = g.StringOf(src[i]), g.StringOf(src[j])
		if x, e := det.GetBool(run); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = x.Bool()
		}
		return
	})
	return
}

// similar to express buildPattern
func makeDet(name string, first, second *core.FromValue) rt.BoolEval {
	return &pattern.DetermineBool{
		Pattern: name,
		Arguments: &core.Arguments{
			Args: []*core.Argument{{
				Name: "$1",
				From: first,
			}, {
				Name: "$2",
				From: second,
			}},
		},
	}
}
