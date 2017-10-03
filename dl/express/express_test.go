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
	r "reflect"
	"testing"
)

func TestExpr(t *testing.T) {
	bigDotFn := func() rt.NumberEval {
		return &core.Get{
			Obj:  &core.Object{"A"},
			Prop: "num",
		}
	}
	binaryFn := func() rt.NumberEval {
		return &core.Mul{
			bigDotFn(),
			&core.Get{
				Obj:  &core.Object{"B"},
				Prop: "num",
			},
		}
	}
	//
	tests := map[string]struct {
		str string
		res interface{}
	}{
		"literal": {
			"5", &core.Num{5},
		},
		"no dot": {
			"A", &core.Object{"A"},
		},
		"little dot": {
			"a.b.c", &core.Get{
				Obj: &core.Get{
					Obj:  &GetAt{"a"},
					Prop: "b",
				},
				Prop: "c",
			},
		},
		"big dot": {
			"A.num", bigDotFn(),
		},
		"binary": {
			"A.num * B.num", binaryFn(),
		},
		"chain": {
			"5 + A.num * B.num", &core.Add{
				&core.Num{5},
				binaryFn(),
			},
		},
		"logic": {
			"a && (b || !c)", &core.AllTrue{[]rt.BoolEval{
				&GetAt{"a"},
				&core.AnyTrue{[]rt.BoolEval{
					&GetAt{"b"},
					&core.IsNot{
						&GetAt{"c"},
					},
				}},
			}},
		},
	}
	cmds := ops.NewOps(nil)
	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*Commands)(nil))

	t.Run("literal", func(t *testing.T) {
		test := tests["literal"]
		var root struct{ rt.NumberEval }
		nconvert(t, cmds, &root, nparse(t, test.str))
		testEqual(t, test.res, root.NumberEval)
	})
	t.Run("no dot", func(t *testing.T) {
		test := tests["no dot"]
		var root struct{ rt.ObjectEval }
		nconvert(t, cmds, &root, nparse(t, test.str))
		testEqual(t, test.res, root.ObjectEval)
	})
	t.Run("big dot", func(t *testing.T) {
		test := tests["big dot"]
		var root struct{ rt.ObjectEval }
		nconvert(t, cmds, &root, nparse(t, test.str))
		testEqual(t, test.res, root.ObjectEval)
	})
	t.Run("little dot", func(t *testing.T) {
		test := tests["little dot"]
		var root struct{ rt.ObjectEval }
		nconvert(t, cmds, &root, nparse(t, test.str))
		testEqual(t, test.res, root.ObjectEval)
	})
	t.Run("binary", func(t *testing.T) {
		test := tests["binary"]
		var root struct{ rt.NumberEval }
		nconvert(t, cmds, &root, nparse(t, test.str))
		testEqual(t, test.res, root.NumberEval)
	})
	t.Run("chain", func(t *testing.T) {
		test := tests["chain"]
		var root struct{ rt.NumberEval }
		nconvert(t, cmds, &root, nparse(t, test.str))
		testEqual(t, test.res, root.NumberEval)
	})
	t.Run("logic", func(t *testing.T) {
		test := tests["logic"]
		var root struct{ rt.BoolEval }
		nconvert(t, cmds, &root, nparse(t, test.str))
		testEqual(t, test.res, root.BoolEval)
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
	if e := astExpr(c, n, r.TypeOf(dst)); e != nil {
		t.Fatal(e, pretty.Sprint(n))
	} else if e := c.Build(); e != nil {
		t.Fatal(e, pretty.Sprint(n))
	}
	return
}
