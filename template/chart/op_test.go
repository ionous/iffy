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
		// {str: "(", op: types.LPAREN},
		{str: "<=", op: types.LEQ},
		{str: "#", errors: true},
	}
	for _, n := range m {
		str := n.str
		t.Log("test:", str)
		p := OperatorParser{}
		Parse(&p, str)
		if r, ok := p.GetOperator(); ok == n.errors {
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
