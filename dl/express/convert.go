package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/spec"
	"github.com/kr/pretty"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// ParseExpression interprets a string as a golang expression, converting it into iffy commands.
func ParseExpression(c spec.Block, s string) (err error) {
	if a, e := parser.ParseExpr(s); e != nil {
		err = e
	} else {
		err = astExpr(c, a, true)
	}
	return
}

// ParseDirective handles multi-part expressions:
// namely "go <return the value of a function>", and
// and, "directive <return the value of a pattern>".
func ParseDirective(c spec.Block, parts []string) (err error) {
	// not wild about this being here --
	// but they are technically expressions and not templates.
	if len(parts) == 1 {
		err = ParseExpression(c, parts[0])
	} else {
		switch op, rest := parts[0], parts[1:]; op {
		case "go":
			op, rest := rest[0], rest[1:]
			if c.Cmd(op).Begin() {
				if len(rest) > 0 {
					err = ParseDirective(c, rest)
				}
				c.End()
			}
		case "determine":
			if c.Cmd(op).Begin() {
				pat, rest := rest[0], rest[1:]
				if c.Cmd(pat).Begin() {
					if len(rest) > 0 {
						err = ParseDirective(c, rest)
					}
					c.End()
				}
				c.End()
			}
		default:
			err = errutil.New("unknown multi-part expression", parts)
		}
	}
	return
}

type Hint bool

func astExpr(c spec.Slot, n ast.Expr, hint Hint) (err error) {
	switch n := n.(type) {
	case *ast.BasicLit:
		err = astBasicLit(c, n)

	case *ast.BinaryExpr:
		err = astBinaryExpr(c, n)

	case *ast.Ident:
		err = astIdent(c, n)

	case *ast.SelectorExpr:
		err = astSelectorExpr(c, n, hint)

	default:
		err = errutil.New("unsupported node", pretty.Sprint(n))
	}
	return
}

func astBasicLit(c spec.Slot, n *ast.BasicLit) (err error) {
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

func astBinaryExpr(c spec.Slot, n *ast.BinaryExpr) (err error) {
	if op, ok := binaryMath[n.Op]; !ok {
		err = errutil.New("unsupported operation", n.Op)
	} else if c := c.Cmd(op); c.Begin() {
		astExpr(c, n.X, false)
		astExpr(c, n.Y, false)
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

func astSelectorExpr(c spec.Slot, n *ast.SelectorExpr, hint Hint) (err error) {
	var cmd string
	if !hint {
		cmd = "get"
	} else {
		cmd = "render"
	}
	if c := c.Cmd(cmd); c.Begin() {
		astExpr(c.Param("obj"), n.X, false)
		c.Param("prop").Val(n.Sel.Name)
		c.End()
	}
	return
}

func astIdent(c spec.Slot, n *ast.Ident) (err error) {
	if name := n.Name; lang.IsCapitalized(name) {
		c.Cmd("object", name)
	} else {
		c.Cmd("get at", name)
	}
	return
}
