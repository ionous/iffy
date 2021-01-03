package list

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// SortNumbers implements Sorter
type SortNumbers struct {
	Var     core.Variable `if:"label:numbers"`
	ByField *SortByField  `if:"unlabeled"`
	Order   `if:"unlabeled"`
}

// SortText implements Sorter
type SortText struct {
	Var     core.Variable `if:"label:text"`
	ByField *SortByField  `if:"unlabeled"`
	Order   `if:"unlabeled"`
	Case    `if:"unlabeled"`
}

// SortRecords implements Sorter
type SortRecords struct {
	Var   core.Variable `if:"label:records"`
	Using pattern.PatternName
}

type SortByField struct {
	Name string `if:"unlabeled"`
}

func (op *SortText) Sort(rt.Runtime, g.Value) (ret error) {
	return
}

func (op *SortNumbers) Sort(rt.Runtime, g.Value) (ret error) {
	return
}

func (op *SortRecords) Sort(rt.Runtime, g.Value) (ret error) {
	return
}

func (op *SortByField) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_by_field",
		Group:  "list",
		Desc:   `Sort numbers: .`,
		Fluent: &composer.Fluid{Name: "byField", Role: composer.Selector},
	}
}

func (op *SortNumbers) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_numbers",
		Group:  "list",
		Desc:   `Sort numbers: .`,
		Fluent: &composer.Fluid{Name: "sort", Role: composer.Command},
	}
}
func (op *SortText) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_text",
		Group:  "list",
		Desc:   `Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.`,
		Fluent: &composer.Fluid{Name: "sort", Role: composer.Command},
	}
}
func (op *SortRecords) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_using",
		Group:  "list",
		Desc:   `Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.`,
		Fluent: &composer.Fluid{Name: "sort", Role: composer.Command},
	}
}

// func (op *Sort) sort(run rt.Runtime) (err error) {
// 	if vs, e := safe.List(run, op.List); e != nil {
// 		err = e
// 	} else {
// 		var newVals g.Value
// 		switch a := vs.Affinity(); a {
// 		case affine.NumList:
// 			els := vs.Floats()
// 			if e := op.sortNumbers(run, els); e != nil {
// 				err = e
// 			} else {
// 				newVals = g.FloatsOf(els)
// 			}
// 		case affine.TextList:
// 			els := vs.Strings()
// 			if e := op.sortText(run, els); e != nil {
// 				err = e
// 			} else {
// 				newVals = g.StringsOf(els)
// 			}
// 		default:
// 			err = errutil.Fmt("variable '%s(%s)' isn't a list", op.List, a)
// 		}
// 		if err == nil {
// 			err = run.SetField(object.Variables, op.List, newVals)
// 		}
// 	}
// 	return
// }

// func (op *Sort) sortNumbers(run rt.Runtime, src []float64) (err error) {
// 	var one, two core.FromValue
// 	det := makeDet(op.Pattern, &one, &two)
// 	sort.Slice(src, func(i, j int) (ret bool) {
// 		one.Value, two.Value = g.FloatOf(src[i]), g.FloatOf(src[j])
// 		if x, e := det.GetBool(run); e != nil {
// 			err = errutil.Append(err, e)
// 		} else {
// 			ret = x.Bool()
// 		}
// 		return
// 	})
// 	return
// }

// func (op *Sort) sortText(run rt.Runtime, src []string) (err error) {
// 	var one, two core.FromValue
// 	det := makeDet(op.Pattern, &one, &two)
// 	sort.Slice(src, func(i, j int) (ret bool) {
// 		one.Value, two.Value = g.StringOf(src[i]), g.StringOf(src[j])
// 		if x, e := det.GetBool(run); e != nil {
// 			err = errutil.Append(err, e)
// 		} else {
// 			ret = x.Bool()
// 		}
// 		return
// 	})
// 	return
// }

// // similar to express buildPattern
// func makeDet(name pattern.PatternName, first, second *core.FromValue) rt.BoolEval {
// 	return &pattern.DetermineBool{
// 		Pattern: name,
// 		Arguments: &core.Arguments{
// 			Args: []*core.Argument{{
// 				Name: "$1",
// 				From: first,
// 			}, {
// 				Name: "$2",
// 				From: second,
// 			}},
// 		},
// 	}
// }
