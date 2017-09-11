package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"testing"
)

func TestExpr(t *testing.T) {
	const (
		literalStr   = "5"
		noDotStr     = "A"
		bigDotStr    = "A.num"
		littleDotStr = "a.b.c"
		binaryStr    = "A.num * B.num"
		chainStr     = "5 + A.num * B.num"
		shortStr     = "A.num= 5"
		longStr      = "B.num= 5 + A.num * B.num"
		// paren      = "(A.num) * (B.num + 5)"
	)
	t.Run("literal", func(t *testing.T) {
		testEqual(t, literalFn(),
			nconvert(t, nparse(t, literalStr)))
	})
	t.Run("no dot", func(t *testing.T) {
		testEqual(t, noDotFn(),
			nconvert(t, nparse(t, noDotStr)))
	})
	t.Run("big dot", func(t *testing.T) {
		testEqual(t, bigDotFn(),
			nconvert(t, nparse(t, bigDotStr)))
	})
	t.Run("little dot", func(t *testing.T) {
		testEqual(t, littleDotFn(),
			nconvert(t, nparse(t, littleDotStr)))
	})
	t.Run("binary", func(t *testing.T) {
		testEqual(t, binaryFn(),
			nconvert(t, nparse(t, binaryStr)))
	})
	t.Run("chain", func(t *testing.T) {
		testEqual(t, chainFn(),
			nconvert(t, nparse(t, chainStr)))
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

func nconvert(t *testing.T, n ast.Expr) (ret interface{}) {
	if v, e := ConvertExpr(n); e != nil {
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
func noDotFn() rt.ObjectEval {
	return &core.Global{"A"}
}
func bigDotFn() rt.NumberEval {
	return &core.Get{
		Obj:  &core.Global{"A"},
		Prop: "num",
	}
}
func littleDotFn() rt.NumberEval {
	return &core.Get{
		Obj: &core.Get{
			Obj:  &core.GetAt{"a"},
			Prop: "b",
		},
		Prop: "c",
	}
}

func binaryFn() rt.NumberEval {
	return &core.Mul{
		bigDotFn(),
		&core.Get{
			Obj:  &core.Global{"B"},
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
		Obj:  &core.Global{"A"},
		Prop: "num",
		Val:  literalFn(),
	}
}

func longAssigmentFn() rt.Execute {
	return &core.SetNum{
		Obj:  &core.Global{"B"},
		Prop: "num",
		Val:  chainFn(),
	}
}
