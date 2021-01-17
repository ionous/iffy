package pattern_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/writer"
	"github.com/ionous/iffy/test/testutil"
)

func B(i bool) rt.BoolEval     { return &core.Bool{i} }
func I(i int) rt.NumberEval    { return &core.Number{float64(i)} }
func T(i string) rt.TextEval   { return &core.Text{i} }
func V(i string) *core.Var     { return &core.Var{Name: i} }
func N(n string) core.Variable { return core.Variable{Str: n} }

type baseRuntime struct {
	testutil.PanicRuntime
}

type patternRuntime struct {
	baseRuntime
	scope.ScopeStack    // parameters are pushed onto the stack.
	testutil.PatternMap // holds pointers to patterns
}

func (patternRuntime) Writer() writer.Output {
	return writer.NewStdout()
}
