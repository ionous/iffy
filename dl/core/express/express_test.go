package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	r "reflect"
	"testing"
)

func TestExpr(t *testing.T) {
	const (
		literalStr = "5"
		unaryStr   = "A.num"
		binaryStr  = "A.num * B.num"
		chainStr   = "5 + A.num * B.num"
		shortStr   = "A.num= 5"
		longStr    = "B.num= 5 + A.num * B.num"
		// paren      = "(A.num) * (B.num + 5)"
	)
	t.Run("literal", func(t *testing.T) {
		testEqual(t, literalFn(),
			nconvert(t, nparse(t, literalStr), nil))
	})
	t.Run("unary", func(t *testing.T) {
		testEqual(t, unaryFn(),
			nconvert(t, nparse(t, unaryStr), nil))
	})
	t.Run("binary", func(t *testing.T) {
		testEqual(t, binaryFn(),
			nconvert(t, nparse(t, binaryStr), nil))
	})
	t.Run("chain", func(t *testing.T) {
		testEqual(t, chainFn(),
			nconvert(t, nparse(t, chainStr), nil))
	})
	t.Run("short", func(t *testing.T) {
		testEqual(t, shortAssignmentFn(),
			sconvert(t, sparse(t, shortStr)))
	})
	t.Run("long", func(t *testing.T) {
		testEqual(t, longAssigmentFn(),
			sconvert(t, sparse(t, longStr)))
	})
}

func testEqual(t *testing.T, expect, res interface{}) {
	if !testify.ObjectsAreEqualValues(expect, res) {
		// res != expect
		t.Log(pretty.Diff(res, expect))
		t.Log("got:", pretty.Sprint(res))
		t.Log("want:", pretty.Sprint(expect))
		t.FailNow()
	}
}

func nparse(t *testing.T, s string) (ret ast.Expr) {
	if a, e := parser.ParseExpr(s); e != nil {
		t.Fatal(e)
	} else {
		ret = a
	}
	return
}

func nconvert(t *testing.T, n ast.Expr, hint r.Type) (ret interface{}) {
	if v, e := ConvertExpr(n, hint); e != nil {
		t.Fatal(e, pretty.Sprint(n))
	} else {
		ret = v
	}
	return
}

func sparse(t *testing.T, s string) (ret ast.Stmt) {
	wrap := "func(){" + s + "}"
	if a, e := parser.ParseExpr(wrap); e != nil {
		t.Fatal(e)
	} else if l := a.(*ast.FuncLit).Body.List; len(l) == 0 {
		t.Fatal("statement is empty")
	} else {
		ret = l[0]
	}
	return
}

func sconvert(t *testing.T, n ast.Stmt) (ret rt.Execute) {
	if v, e := ConvertStmt(n); e != nil {
		t.Fatal(e, pretty.Sprint(n))
	} else {
		ret = v
	}
	return
}

func literalFn() rt.NumberEval {
	return &core.Num{5}
}

func unaryFn() rt.NumberEval {
	return &core.Get{
		Obj:  &core.Object{"A"},
		Prop: "num",
	}
}

func binaryFn() rt.NumberEval {
	return &core.Mul{
		unaryFn(),
		&core.Get{
			Obj:  &core.Object{"B"},
			Prop: "num",
		},
	}
}

func chainFn() rt.NumberEval {
	return &core.Add{
		literalFn(),
		binaryFn(),
	}
}

func shortAssignmentFn() rt.Execute {
	return &core.SetNum{
		Obj:  &core.Object{"A"},
		Prop: "num",
		Val:  literalFn(),
	}
}

func longAssigmentFn() rt.Execute {
	return &core.SetNum{
		Obj:  &core.Object{"B"},
		Prop: "num",
		Val:  chainFn(),
	}
}
