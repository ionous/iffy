package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/test/testutil"
)

func TestAllTrue(t *testing.T) {
	run := &testutil.PanicRuntime{}
	var l boolList
	evals := []rt.BoolEval{}
	for i := 0; i < 3; i++ {
		test := &AllTrue{evals}
		if ok, e := safe.GetBool(run, test); e != nil {
			t.Fatal(e)
		} else if !ok.Bool() {
			t.Fatal("expected success")
		} else if l.asks != len(l.vals) {
			t.Fatal("expected all got tested")
		}
		//
		l.asks, l.vals = 0, append(l.vals, true)
		evals = append(evals, &l)
	}
	// turn one false.
	l.vals[1] = false
	test := &AllTrue{evals}
	if ok, e := safe.GetBool(run, test); e != nil {
		t.Fatal(e)
	} else if ok.Bool() {
		t.Fatal("expected failure")
	} else if l.asks != 2 {
		t.Fatal("expected only two got tested", l.asks)
	}
}

func TestAnyTrue(t *testing.T) {
	run := &testutil.PanicRuntime{}
	var l boolList
	evals := []rt.BoolEval{}
	for i := 0; i < 3; i++ {
		test := &AnyTrue{evals}
		if ok, e := safe.GetBool(run, test); e != nil {
			t.Fatal(e)
		} else if ok.Bool() {
			t.Fatal("expected failure")
		} else if l.asks != i {
			t.Fatal("expected all got tested", l.asks)
		}
		//
		l.asks, l.vals = 0, append(l.vals, false)
		evals = append(evals, &l)
	}
	// turn one true.
	l.vals[1] = true
	test := &AnyTrue{evals}
	if ok, e := safe.GetBool(run, test); e != nil {
		t.Fatal(e)
	} else if !ok.Bool() {
		t.Fatal("expected success")
	} else if l.asks != 2 {
		t.Fatal("expected two got tested", l.asks)
	}
}

type boolList struct {
	vals []bool
	asks int
}

func (b *boolList) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if a := b.asks; a >= len(b.vals) {
		err = errutil.New("out of range")
	} else {
		ok := b.vals[a]
		b.asks = a + 1
		ret = g.BoolOf(ok)
	}
	return
}
