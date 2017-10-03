package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/spec"
	"github.com/kr/pretty"
	"go/ast"
	"go/parser"
	"go/token"
	r "reflect"
	"strconv"
)

// ParseExpression interprets a string as a golang expression, converting it into iffy commands.
func ParseExpression(c spec.Block, s string, hint r.Type) (err error) {
	if a, e := parser.ParseExpr(s); e != nil {
		err = e
	} else {
		err = astExpr(c, a, hint)
	}
	return
}

// ParseDirective handles multi-part expressions:
// namely "go <return the value of a function>", and
// and, "directive <return the value of a pattern>".
func ParseDirective(c spec.Block, parts []string, hint r.Type) (err error) {
	// not wild about this being here --
	// but they are technically expressions and not templates.
	if len(parts) == 1 {
		err = ParseExpression(c, parts[0], hint)
	} else {
		switch op, rest := parts[0], parts[1:]; op {
		case "go":
			op, rest := rest[0], rest[1:]
			if c.Cmd(op).Begin() {
				if len(rest) > 0 {
					err = ParseDirective(c, rest, hint)
				}
				c.End()
			}
		case "determine":
			if c.Cmd(op).Begin() {
				pat, rest := rest[0], rest[1:]
				if c.Cmd(pat).Begin() {
					if len(rest) > 0 {
						err = ParseDirective(c, rest, hint)
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

type Hint r.Type

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

	case *ast.ParenExpr:
		err = astParenExpr(c, n, hint)

	case *ast.UnaryExpr:
		err = astUnaryExpr(c, n, hint)

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
	if op, ok := binaryMath[n.Op]; ok {
		if c := c.Cmd(op); c.Begin() {
			err = binaryPair(c, n.X, n.Y, kindOf.TypeNumEval)
			c.End()
		}
	} else if op, ok := anyAll[n.Op]; !ok {
		err = errutil.New("unsupported operation", n.Op)
	} else {
		if c := c.Cmd(op); c.Begin() {
			if c.Cmds().Begin() {
				err = binaryPair(c, n.X, n.Y, kindOf.TypeObjEval)
				c.End()
			}
			c.End()
		}
	}
	return
}

func binaryPair(c spec.Slot, x, y ast.Expr, hint Hint) (err error) {
	if e := astExpr(c, x, hint); e != nil {
		err = e
	} else if e := astExpr(c, y, hint); e != nil {
		err = e
	}
	return
}

func astParenExpr(c spec.Slot, n *ast.ParenExpr, hint Hint) error {
	return astExpr(c, n.X, hint)
}

func astUnaryExpr(c spec.Slot, n *ast.UnaryExpr, hint Hint) (err error) {
	switch n.Op {
	case token.NOT:
		if c := c.Cmd("is not"); c.Begin() {
			err = astExpr(c, n.X, hint)
			c.End()
		}
	default:
		err = errutil.New("unsupported unary expression", n.Op)
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

var anyAll = map[token.Token]string{
	token.LAND: "all true",
	token.LOR:  "any true",
}

func astSelectorExpr(c spec.Slot, n *ast.SelectorExpr, hint Hint) (err error) {
	var cmd string
	if kindOf.TextEval(hint) {
		cmd = "render"
	} else {
		cmd = "get"
	}
	if c := c.Cmd(cmd); c.Begin() {
		astExpr(c.Param("obj"), n.X, kindOf.TypeObjEval)
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
