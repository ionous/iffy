package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

func TestAllTrue(t *testing.T) {
	run := &rt.Panic{}
	var l boolList
	evals := []rt.BoolEval{}
	for i := 0; i < 3; i++ {
		test := &AllTrue{evals}
		if ok, e := rt.GetBool(run, test); e != nil {
			t.Fatal(e)
		} else if !ok {
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
	if ok, e := rt.GetBool(run, test); e != nil {
		t.Fatal(e)
	} else if ok {
		t.Fatal("expected failure")
	} else if l.asks != 2 {
		t.Fatal("expected only two got tested", l.asks)
	}
}

func TestAnyTrue(t *testing.T) {
	run := &rt.Panic{}
	var l boolList
	evals := []rt.BoolEval{}
	for i := 0; i < 3; i++ {
		test := &AnyTrue{evals}
		if ok, e := rt.GetBool(run, test); e != nil {
			t.Fatal(e)
		} else if ok {
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
	if ok, e := rt.GetBool(run, test); e != nil {
		t.Fatal(e)
	} else if !ok {
		t.Fatal("expected success")
	} else if l.asks != 2 {
		t.Fatal("expected two got tested", l.asks)
	}
}

type boolList struct {
	vals []bool
	asks int
}

func (b *boolList) GetBool(run rt.Runtime) (okay bool, err error) {
	if a := b.asks; a >= len(b.vals) {
		err = errutil.New("out of range")
	} else {
		okay, b.asks = b.vals[a], a+1
	}
	return
}
