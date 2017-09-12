package express

// import (
// 	"github.com/ionous/iffy/dl/core"
// 	"github.com/ionous/iffy/ref/unique"
// 	"github.com/ionous/iffy/rt"
// 	"github.com/ionous/iffy/spec/ops"
// 	"github.com/kr/pretty"
// 	testify "github.com/stretchr/testify/assert"
// 	"go/ast"
// 	"go/parser"
// 	"testing"
// )

// func TestAssign(t *testing.T) {
// 	const (
// 		shortStr = "A.num= 5"
// 		longStr  = "B.num= 5 + A.num * B.num"
// 		// paren      = "(A.num) * (B.num + 5)"
// 	)
// 	cmds := ops.NewOps(nil)
// 	unique.PanicBlocks(cmds,
// 		(*core.Commands)(nil))

// 	t.Run("short", func(t *testing.T) {
// 		var root struct{ rt.Execute }
// 		sconvert(t, cmds, &root, sparse(t, shortStr))
// 		testEqual(t, shortAssignmentFn(), root.Execute)
// 	})
// 	t.Run("long", func(t *testing.T) {
// 		var root struct{ rt.Execute }
// 		sconvert(t, cmds, &root, sparse(t, longStr))
// 		testEqual(t, longAssigmentFn(), root.Execute)
// 	})
// }

// func sparse(t *testing.T, s string) (ret ast.Stmt) {
// 	wrap := "func(){" + s + "}"
// 	if a, e := parser.ParseExpr(wrap); e != nil {
// 		t.Fatal(e)
// 	} else if l := a.(*ast.FuncLit).Body.List; len(l) == 0 {
// 		t.Fatal("statement is empty")
// 	} else {
// 		ret = l[0]
// 	}
// 	return
// }

// func sconvert(t *testing.T, n ast.Stmt) (ret rt.Execute) {
// 	if e := ConvertStmt(n); e != nil {
// 		t.Fatal(e, pretty.Sprint(n))
// 	} else {
// 		ret = v
// 	}
// 	return
// }

// func shortAssignmentFn() rt.Execute {
// 	return &core.SetNum{
// 		Obj:  &core.Global{"A"},
// 		Prop: "num",
// 		Val:  literalFn(),
// 	}
// }

// func longAssigmentFn() rt.Execute {
// 	return &core.SetNum{
// 		Obj:  &core.Global{"B"},
// 		Prop: "num",
// 		Val:  chainFn(),
// 	}
// }
