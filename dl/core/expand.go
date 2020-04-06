package core

import (
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

func GetBool(run rt.Runtime, i interface{}) (ret bool, err error) {
	if eval, ok := i.(rt.BoolEval); ok {
		ret, err = eval.GetBool(run)
	} else {
		ret, err = assign.ToBool(i)
	}
	return
}

func GetNumber(run rt.Runtime, i interface{}) (ret float64, err error) {
	if eval, ok := i.(rt.NumberEval); ok {
		ret, err = eval.GetNumber(run)
	} else {
		ret, err = assign.ToFloat(i)
	}
	return
}

func GetText(run rt.Runtime, i interface{}) (ret string, err error) {
	if eval, ok := i.(rt.TextEval); ok {
		ret, err = eval.GetText(run)
	} else {
		ret, err = assign.ToString(i)
	}
	return
}

func GetNumbers(run rt.Runtime, i interface{}) (ret rt.Iterator, err error) {
	if eval, ok := i.(rt.NumListEval); ok {
		ret, err = eval.GetNumberStream(run)
	} else if vs, ok := i.([]float64); ok {
		ret = stream.NewNumberList(vs)
	} else {
		err = assign.Mismatch("NumListEval", eval, i)
	}
	return
}

func GetTexts(run rt.Runtime, i interface{}) (ret rt.Iterator, err error) {
	if eval, ok := i.(rt.TextListEval); ok {
		ret, err = eval.GetTextStream(run)
	} else if vs, ok := i.([]string); ok {
		ret = stream.NewTextList(vs)
	} else {
		err = assign.Mismatch("TextListEval", eval, i)
	}
	return
}
