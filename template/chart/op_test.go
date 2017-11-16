package chart

import (
	"github.com/ionous/iffy/template/types"
	"github.com/kr/pretty"
	"sort"
	"testing"
)

func TestOpListIsSorted(t *testing.T) {
	test := make([]_Match, len(list))
	copy(test, list)
	sort.Slice(test, func(i, j int) bool {
		return test[i].Text < test[j].Text
	})
	if diff := pretty.Diff(test, list); len(diff) > 0 {
		t.Log(pretty.Sprint(test))
		t.FailNow()
	}
}

func TestOps(t *testing.T) {
	type Test struct {
		str    string
		op     types.Operator
		errors bool
	}
	m := []Test{
		{str: "andy", op: types.LAND},
		{str: ">>>>", op: types.GTR},
		{str: "*", op: types.MUL},
		{str: "<=", op: types.LEQ},
		// this is not an operator error, it simply doesnt parse anything.
		{str: "#", op: types.Operator(-1)},
		// this provisionally matches, which has to generate an error because it consumes runes
		// without actually turning into a result.
		{str: "!!", errors: true},
	}
	for _, n := range m {
		str := n.str
		t.Log("test:", str)
		p := OperatorParser{}
		_ = Parse(&p, str) // this returns an error, lets ignore it.
		r, e := p.GetOperator()
		ok := e == nil
		if ok == n.errors {
			t.Fatalf("unexpected result %v for '%s'", ok, str)
		} else if !ok {
			t.Log("ok expected mismatch")
		} else if r != n.op {
			t.Fatalf("mismatch %s != %s", r, n.op)
		} else {
			t.Log("matched", n.op)
		}
	}
}
