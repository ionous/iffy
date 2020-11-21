package list_test

import (
	"strconv"
	"strings"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/test"
	"github.com/ionous/iffy/rt/writer"
)

type panicTime struct {
	test.PanicRuntime
}
type listTime struct {
	panicTime
	objs map[string]*g.Record
	scope.ScopeStack
	pattern.PatternMap
	kinds *test.Kinds
}

func (lt *listTime) Writer() writer.Output {
	return writer.NewStdout()
}
func newListTime(src []string, p pattern.PatternMap) (ret rt.Runtime, vals *g.Record, err error) {
	var kinds test.Kinds
	type Values struct{ Source []string }
	kinds.AddKinds((*Values)(nil))
	values := kinds.New("Values")
	lt := listTime{
		kinds:      &kinds,
		PatternMap: p,
		ScopeStack: scope.ScopeStack{
			Scopes: []rt.Scope{
				&scope.TargetRecord{object.Variables, values},
			},
		},
	}
	if e := values.SetNamedField("Source", g.StringsOf(src)); e != nil {
		err = e
	} else {
		ret = &lt
		vals = values
	}
	return
}

func B(i bool) rt.BoolEval    { return &core.Bool{i} }
func I(i int) rt.NumberEval   { return &core.Number{float64(i)} }
func T(i string) rt.TextEval  { return &core.Text{i} }
func V(i string) *core.GetVar { return &core.GetVar{Name: i} }

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
		(*op.out) = append((*op.out), t.String())
	}
	return
}

func getNum(run rt.Runtime, op rt.NumberEval) (ret string) {
	if v, e := op.GetNumber(run); e != nil {
		ret = e.Error()
	} else {
		ret = strconv.FormatFloat(v.Float(), 'g', -1, 64)
	}
	return
}

func joinText(run rt.Runtime, op rt.TextListEval) (ret string) {
	if vs, e := op.GetTextList(run); e != nil {
		ret = e.Error()
	} else {
		ret = joinStrings(vs.Strings())
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

func (lt *listTime) GetField(target, field string) (ret g.Value, err error) {
	if obj, ok := lt.objs[field]; target == object.Value && ok {
		ret = g.RecordOf(obj)
	} else {
		ret, err = lt.ScopeStack.GetField(target, field)
	}
	return
}

func (lt *listTime) SetField(target, field string, value g.Value) (err error) {
	return lt.ScopeStack.SetField(target, field, value)
}

func (lt *listTime) GetKindByName(name string) (*g.Kind, error) {
	return lt.kinds.GetKindByName(name)
}
