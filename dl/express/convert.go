package express

import (
	"github.com/ionous/errutil"
	"strconv"
	// "github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	// "github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec"
	"github.com/kr/pretty"
	"go/ast"
	"go/token"
)

type Hint bool

func ConvertExpr(c spec.Slot, n ast.Expr) (err error) {
	return convertExpr(c, n, false)
}

func convertExpr(c spec.Slot, n ast.Expr, hint Hint) (err error) {
	switch n := n.(type) {
	case *ast.BasicLit:
		err = BasicLit(c, n)

	case *ast.BinaryExpr:
		err = BinaryExpr(c, n)

	case *ast.Ident:
		makeObject(c, n)

	case *ast.SelectorExpr:
		err = SelectorExpr(c, n, hint)

	default:
		err = errutil.New("unsupported node", pretty.Sprint(n))
	}
	return
}

func BasicLit(c spec.Slot, n *ast.BasicLit) (err error) {
	switch t, v := n.Kind, n.Value; t {
	case token.STRING:
		c.Val(v)
	case token.FLOAT, token.INT:
		// "literally" doesnt translate text to numbers, so we have to do it manually:
		if v, e := strconv.ParseFloat(v, 64); e != nil {
			err = e
		} else {
			c.Cmd("num", v)
		}
	default:
		//token.IMAG, token.CHAR,
		err = errutil.New("unsupported literal token", t)
	}
	return
}

func BinaryExpr(c spec.Slot, n *ast.BinaryExpr) (err error) {
	if op, ok := binaryMath[n.Op]; !ok {
		err = errutil.New("unsupported operation", n.Op)
	} else if c := c.Cmd(op); c.Begin() {
		ConvertExpr(c, n.X)
		ConvertExpr(c, n.Y)
		c.End()
	}
	return
}

var binaryMath = map[token.Token]string{
	token.ADD: "add",
	token.SUB: "sub",
	token.MUL: "mul",
	token.QUO: "div",
	token.REM: "mod",
}

func SelectorExpr(c spec.Slot, n *ast.SelectorExpr, hint Hint) (err error) {
	var cmd string
	if !hint {
		cmd = "get"
	} else {
		cmd = "render"
	}
	if c := c.Cmd(cmd); c.Begin() {
		ConvertExpr(c.Param("obj"), n.X)
		c.Param("prop").Val(n.Sel.Name)
		c.End()
	}
	return
}

func makeObject(c spec.Slot, n *ast.Ident) {
	if name := n.Name; lang.IsCapitalized(name) {
		c.Cmd("global", name)
	} else {
		c.Cmd("get at", name)
	}
}
