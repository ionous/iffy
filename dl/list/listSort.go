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
	"github.com/ionous/iffy/rt/generic"
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
	if e := op.execute(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Sort) execute(run rt.Runtime) (err error) {
	if vs, e := run.GetField(object.Variables, op.List); e != nil {
		err = e
	} else {
		var newVals rt.Value // write back in case vs was temporary storage.
		switch a := vs.Affinity(); a {
		case affine.NumList:
			if els, e := vs.GetNumList(); e != nil {
				err = e
			} else if e := op.sortNumbers(run, els); e != nil {
				err = e
			} else {
				newVals = &generic.FloatSlice{Values: els}
			}
		case affine.TextList:
			if els, e := vs.GetTextList(); e != nil {
				err = e
			} else if e := op.sortText(run, els); e != nil {
				err = e
			} else {
				newVals = &generic.StringSlice{Values: els}
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
	var one, two core.Number
	det := makeDet("sort", &core.FromNum{&one}, &core.FromNum{&two})
	sort.Slice(src, func(i, j int) (ret bool) {
		one.Num, two.Num = src[i], src[j]
		if x, e := det.GetBool(run); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = x
		}
		return
	})
	return
}

func (op *Sort) sortText(run rt.Runtime, src []string) (err error) {
	var one, two core.Text
	det := makeDet("sort", &core.FromText{&one}, &core.FromText{&two})
	sort.Slice(src, func(i, j int) (ret bool) {
		one.Text, two.Text = src[i], src[j]
		if x, e := det.GetBool(run); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = x
		}
		return
	})
	return
}

// similar to express buildPattern
func makeDet(name string, first, second core.Assignment) rt.BoolEval {
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
