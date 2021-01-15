package core

import (
	"errors"
	"testing"

	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

func TestLoopBreak(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			&Bool{true}, MakeActivity(
				&Assign{N("i"), &FromNum{&SumOf{V("i"), I(1)}}},
				&ChooseAction{
					If: &CompareNum{V("i"), &GreaterOrEqual{}, I(4)},
					Do: MakeActivity(
						&Break{},
					),
				},
				// &Next{},
				&Assign{N("j"), &FromNum{&SumOf{V("j"), I(1)}}},
			)},
	); e != nil {
		t.Fatal(e)
	} else if run.i != 4 && run.j != 3 {
		t.Fatal("bad counters", run.i, run.j)
	}
}

func TestLoopNext(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			&Bool{true}, MakeActivity(
				&Assign{N("i"), &FromNum{&SumOf{V("i"), I(1)}}},
				&ChooseAction{
					If: &CompareNum{V("i"), &GreaterOrEqual{}, I(4)},
					Do: MakeActivity(
						&Break{},
					),
				},
				&Next{},
				&Assign{N("j"), &FromNum{&SumOf{V("j"), I(1)}}},
			)},
	); e != nil {
		t.Fatal(e)
	} else if run.i != 4 && run.j != 0 {
		t.Fatal("bad counters", run.i, run.j)
	}
}

func TestLoopInfinite(t *testing.T) {
	MaxLoopIterations = 100
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			&Bool{true}, MakeActivity(
				&Assign{N("i"), &FromNum{&SumOf{V("i"), I(1)}}},
			)},
	); !errors.Is(e, MaxLoopIterations) {
		t.Fatal(e)
	} else if run.i != 100 {
		t.Fatal("bad counters", run.i, run.j)
	} else {
		t.Log("ok, error is expected:", e)
	}
}

type loopRuntime struct {
	baseRuntime
	i, j int
}

func (k *loopRuntime) GetField(target, field string) (ret g.Value, err error) {
	switch {
	case field == "i" && target == object.Variables:
		ret = g.IntOf(k.i)
	case field == "j" && target == object.Variables:
		ret = g.IntOf(k.j)
	default:
		panic("unexpected get")
	}
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *loopRuntime) SetField(target, field string, v g.Value) (err error) {
	switch {
	case field == "i" && target == object.Variables:
		k.i = v.Int()
	case field == "j" && target == object.Variables:
		k.j = v.Int()
	default:
		panic("unexpected set")
	}
	return
}

func V(i string) *Var        { return &Var{Name: i} }
func N(n string) Variable    { return Variable{Str: n} }
func T(s string) rt.TextEval { return &Text{s} }
func I(n int) rt.NumberEval  { return &Number{float64(n)} }
func B(b bool) rt.BoolEval   { return &Bool{b} }
