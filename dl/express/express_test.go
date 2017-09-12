package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
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
		// paren      = "(A.num) * (B.num + 5)"
	)
	cmds := ops.NewOps(nil)
	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*Commands)(nil))

	t.Run("literal", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		nconvert(t, cmds, &root, nparse(t, literalStr))
		testEqual(t, literalFn(), root.NumberEval)
	})
	t.Run("no dot", func(t *testing.T) {
		var root struct{ rt.ObjectEval }
		nconvert(t, cmds, &root, nparse(t, noDotStr))
		testEqual(t, noDotFn(), root.ObjectEval)
	})
	t.Run("big dot", func(t *testing.T) {
		var root struct{ rt.ObjectEval }
		nconvert(t, cmds, &root, nparse(t, bigDotStr))
		testEqual(t, bigDotFn(), root.ObjectEval)
	})
	t.Run("little dot", func(t *testing.T) {
		var root struct{ rt.ObjectEval }
		nconvert(t, cmds, &root, nparse(t, littleDotStr))
		testEqual(t, littleDotFn(), root.ObjectEval)
	})
	t.Run("binary", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		nconvert(t, cmds, &root, nparse(t, binaryStr))
		testEqual(t, binaryFn(), root.NumberEval)
	})
	t.Run("chain", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		nconvert(t, cmds, &root, nparse(t, chainStr))
		testEqual(t, chainFn(), root.NumberEval)
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

func nconvert(t *testing.T, cmds *ops.Ops, dst interface{}, n ast.Expr) {
	c := cmds.NewBuilder(dst, core.Xform{})
	if e := ConvertExpr(c, n); e != nil {
		t.Fatal(e, pretty.Sprint(n))
	} else if e := c.Build(); e != nil {
		t.Fatal(e, pretty.Sprint(n))
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
			Obj:  &GetAt{"a"},
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
