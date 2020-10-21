package list_test

import (
	"strconv"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type listTime struct {
	rt.Panic
	strings []string
}

func I(i int) rt.NumberEval {
	return &core.Number{float64(i)}
}

func FromTs(vs []string) (ret core.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{&core.Text{vs[0]}}
	} else {
		ret = &core.FromTextList{&core.Texts{vs}}
	}
	return
}

func getNum(run rt.Runtime, op rt.NumberEval) (ret string) {
	if v, e := op.GetNumber(run); e != nil {
		ret = e.Error()
	} else {
		ret = strconv.FormatFloat(v, 'g', -1, 64)
	}
	return
}

func joinText(run rt.Runtime, op rt.TextListEval) (ret string) {
	if vs, e := op.GetTextList(run); e != nil {
		ret = e.Error()
	} else {
		ret = joinStrings(vs)
	}
	return
}

func joinStrings(vs []string) (ret string) {
	if len(vs) > 0 {
		ret = strings.Join(vs, ", ")
	} else {
		ret = "-"
	}
	return
}

func (g *listTime) GetField(target, field string) (ret rt.Value, err error) {
	if target != object.Variables {
		err = errutil.New("unexpected target", target)
	} else if field != "strings" {
		err = errutil.New("unexpected field", field)
	} else {
		ret = &generic.StringSlice{Values: g.strings}
	}
	return
}

func (g *listTime) SetField(target, field string, value rt.Value) (err error) {
	if target != object.Variables {
		err = errutil.New("unexpected target", target)
	} else if field != "strings" {
		err = errutil.New("unexpected field", field)
	} else if vs, e := value.GetTextList(); e != nil {
		err = e
	} else {
		g.strings = vs
	}
	return
}
