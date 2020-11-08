package list_test

import (
	"strconv"
	"strings"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
)

type panicTime struct {
	rt.Panic
}
type listTime struct {
	panicTime
	strings []string
	scope.ScopeStack
	sort *pattern.BoolPattern
}

func B(i bool) rt.BoolEval   { return &core.Bool{i} }
func I(i int) rt.NumberEval  { return &core.Number{float64(i)} }
func T(i string) rt.TextEval { return &core.Text{i} }

func FromTs(vs []string) (ret core.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{&core.Text{vs[0]}}
	} else {
		ret = &core.FromTextList{&core.Texts{vs}}
	}
	return
}

// cmd to collect some text into a list of strings.
type Write struct {
	out  *[]string
	Text rt.TextEval
}

// Execute writes text to the runtime's current writer.
func (op *Write) Execute(run rt.Runtime) (err error) {
	if t, e := op.Text.GetText(run); e != nil {
		err = e
	} else {
		(*op.out) = append((*op.out), t)
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
	if target != object.Variables || field != "strings" {
		ret, err = g.ScopeStack.GetField(target, field)
	} else {
		ret  = generic.StringsOf(g.strings)
	}
	return
}

func (g *listTime) SetField(target, field string, value rt.Value) (err error) {
	if target != object.Variables || field != "strings" {
		err = g.ScopeStack.SetField(target, field, value)
	} else if vs, e := value.GetTextList(); e != nil {
		err = e
	} else {
		g.strings = vs
	}
	return
}

func (g *listTime) GetEvalByName(name string, pv interface{}) (err error) {
	if name == "sort" {
		ptr := pv.(*pattern.BoolPattern)
		(*ptr) = *g.sort
	}
	return
}
