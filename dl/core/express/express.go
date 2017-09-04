package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"go/ast"
	"go/token"
	r "reflect"
	"strconv"
)

func ConvertExpr(n ast.Expr, hint r.Type) (ret interface{}, err error) {
	switch n := n.(type) {
	case *ast.BasicLit:
		ret, err = BasicLit(n, hint)
	case *ast.BinaryExpr:
		ret, err = BinaryExpr(n, hint)
	case *ast.SelectorExpr:
		ret, err = SelectorExpr(n, hint)
	default:
		err = errutil.Fmt("unsupported node %T", n)
	}
	return
}

func ConvertStmt(l ast.Stmt) (ret rt.Execute, err error) {
	var list rt.ExecuteList
	switch l := l.(type) {
	case *ast.AssignStmt:
		if cnt := len(l.Lhs); cnt != len(l.Rhs) {
			err = errutil.New("left and right sides dont match")
		} else {
			for i := 0; i < cnt; i++ {
				lhs, rhs := l.Lhs[i], l.Rhs[i]
				if r, e := assign(lhs, rhs); e != nil {
					err = e
					break
				} else if ret != nil {
					list = append(list, r)
					ret = list
				} else {
					ret = r
				}
			}
		}
	}
	return
}

func assign(lhs, rhs ast.Expr) (ret rt.Execute, err error) {
	if n, ok := lhs.(*ast.SelectorExpr); !ok {
		err = errutil.New("error on left, expected object")
	} else if obj, e := identifyObject(n); e != nil {
		err = e
	} else if v, e := ConvertExpr(rhs, nil); e != nil {
		err = errutil.New("error on right", e)
	} else {
		switch v := v.(type) {
		case rt.NumberEval:
			ret = &core.SetNum{obj, n.Sel.Name, v}
		case rt.TextEval:
			ret = &core.SetText{obj, n.Sel.Name, v}
		default:
			err = errutil.New("unknown type %T", v)
		}
	}
	return
}

func BasicLit(n *ast.BasicLit, hint r.Type) (ret interface{}, err error) {
	switch t, v := n.Kind, n.Value; t {
	case token.FLOAT, token.INT:
		if v, e := strconv.ParseFloat(v, 64); e != nil {
			err = e
		} else {
			ret = &core.Num{v}
		}
	case token.STRING:
		ret = &core.Text{v}
	default:
		err = errutil.New("unsupported literal token", t)
		// token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	}
	return
}

func BinaryExpr(n *ast.BinaryExpr, hint r.Type) (ret interface{}, err error) {
	// could pass "rt.Number" to hint for binary math
	if x, e := ConvertExpr(n.X, hint); e != nil {
		err = e
	} else if y, e := ConvertExpr(n.Y, hint); e != nil {
		err = e
	} else if pair, ok := binaryMath[n.Op]; !ok {
		err = errutil.New("unsupported operation", n.Op)
	} else if nx, ok := x.(rt.NumberEval); !ok {
		err = errutil.Fmt("x is not a number %T", x)
	} else if ny, ok := y.(rt.NumberEval); !ok {
		err = errutil.Fmt("x is not a number %T", y)
	} else {
		ret = pair(nx, ny)
	}
	return
}

type pairFn func(a, b rt.NumberEval) rt.NumberEval

var binaryMath = map[token.Token]pairFn{
	token.ADD: func(a, b rt.NumberEval) rt.NumberEval {
		return &core.Add{a, b}
	},
	token.SUB: func(a, b rt.NumberEval) rt.NumberEval {
		return &core.Sub{a, b}
	},
	token.MUL: func(a, b rt.NumberEval) rt.NumberEval {
		return &core.Mul{a, b}
	},
	token.QUO: func(a, b rt.NumberEval) rt.NumberEval {
		return &core.Div{a, b}
	},
	token.REM: func(a, b rt.NumberEval) rt.NumberEval {
		return &core.Mod{a, b}
	},
}

func SelectorExpr(n *ast.SelectorExpr, hint r.Type) (ret interface{}, err error) {
	if obj, e := identifyObject(n); e != nil {
		err = e
	} else {
		ret = &core.Get{
			Obj:  obj,
			Prop: n.Sel.Name,
		}
	}
	return
}

func identifyObject(n *ast.SelectorExpr) (ret *core.Object, err error) {
	if x, ok := n.X.(*ast.Ident); !ok {
		// FIX? include position here? [ could possibly return an error object with position in it, then use that pos to pretty print an error ]
		err = errutil.Fmt("expected object identifer, got %T", n.X)
	} else {
		ret = &core.Object{x.Name}
	}
	// Ident includes:
	// Name, NamePos, and Obj ( which describes a named language entity such as a package, constant, type, variable, function (incl. methods), or label. )
	return
}
