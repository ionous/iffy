package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

// Express converts a postfix expression into iffy commands.
func Convert(cmds *ops.Factory, xs postfix.Expression) (ret *ops.Command, err error) {
	c := converter{cmds: cmds}
	if e := c.convert(xs); e != nil {
	} else if len(c.stack) == 0 {
		err = errutil.New("empty output")
	} else if len(c.stack) > 1 {
		err = errutil.New("unparsed output")
	} else {
		ret = c.stack[0].(*ops.Command)
	}
	return
}

type converter struct {
	cmds  *ops.Factory
	stack []interface{}
}

func (c *converter) convert(xs postfix.Expression) (err error) {
	for _, fn := range xs {
		if e := c.addFunction(fn); e != nil {
			err = e
			break
		}
	}
	return
}

// add a new command pointer to the output stack.
func (c *converter) push(cmd *ops.Command) {
	c.stack = append(c.stack, cmd)
}

// extract cnt commands from the converter stack.
func (c *converter) pop(cnt int) (ret []interface{}, err error) {
	if end := len(c.stack) - cnt; end < 0 {
		err = errutil.New("stack underflow")
	} else {
		ret, c.stack = c.stack[end:], c.stack[:end]
	}
	return
}

func (c *converter) binary(i interface{}) (err error) {
	if args, e := c.pop(2); e != nil {
		err = e
	} else if cmd, e := c.cmds.EmplaceCommand(i); e != nil {
		err = e
	} else {
		err = c.pushCommand(cmd, args...)
	}
	return
}

// adds a comparision statement to the output.
// FIX FIX we have to know the type: num,text,obj of the properties in question
// need more information on properties.
func (c *converter) compare(cmp core.CompareTo) (err error) {
	if args, e := c.pop(2); e != nil {
		err = e
	} else if cmd, e := c.cmds.EmplaceCommand(&core.CompareText{}); e != nil {
		err = e
	} else {
		// as an alternative to this custom code, we could name arguments for every command.
		// ex. MUL: {cmd:Mul, args: "A", "B" }, CompareText[ A, B ]
		// wed still have to have an initalizer or something for "cmp".
		err = c.pushCommand(cmd, args[0], cmp, args[1])
	}
	return
}

// add the passed args to the passed command in order; then push the command.
func (c *converter) pushCommand(cmd *ops.Command, args ...interface{}) (err error) {
	if e := assign(cmd, args); e != nil {
		err = e
	} else {
		c.push(cmd)
	}
	return
}

// add the passed args to the passed command in order.
func assign(cmd *ops.Command, args []interface{}) (err error) {
	for _, arg := range args {
		if e := cmd.Position(arg); e != nil {
			err = e
			break
		}
	}
	return
}

// convert the passed function into iffy commands.
func (c *converter) addFunction(fn postfix.Function) (err error) {
	switch fn := fn.(type) {
	case template.Quote:
		if cmd, e := c.cmds.EmplaceCommand(&core.Text{fn.Value()}); e != nil {
			err = e
		} else {
			c.push(cmd)
		}

	case template.Number:
		if cmd, e := c.cmds.EmplaceCommand(&core.Num{fn.Value()}); e != nil {
			err = e
		} else {
			c.push(cmd)
		}

	case template.Reference:
		if fields := fn.Value(); len(fields) == 0 {
			err = errutil.New("empty reference")
		} else {
			// obj.a.b.c => Get{c Get{b Get{a GetAt{obj}}}}
			var op rt.ObjectEval
			if name := fields[0]; lang.IsCapitalized(name) {
				op = &core.Object{name}
			} else {
				op = &GetAt{name}
			}
			for _, field := range fields[1:] {
				op = &Render{op, field}
			}
			if cmd, e := c.cmds.EmplaceCommand(op); e != nil {
				err = e
			} else {
				c.push(cmd)
			}
		}

	case template.Command:
		if cmd, e := c.cmds.CreateCommand(fn.CommandName); e != nil {
			err = e
		} else if args, e := c.pop(fn.CommandArity); e != nil {
			err = e
		} else {
			err = c.pushCommand(cmd, args...)
		}

	case template.Operator:
		switch fn {
		case template.MUL:
			err = c.binary(&core.Mul{})
		case template.QUO:
			err = c.binary(&core.Div{})
		case template.REM:
			err = c.binary(&core.Mod{})
		case template.ADD:
			err = c.binary(&core.Add{})
		case template.SUB:
			err = c.binary(&core.Sub{})
		case template.EQL:
			err = c.compare(&core.EqualTo{})
		case template.NEQ:
			err = c.compare(&core.NotEqualTo{})
		case template.LSS:
			err = c.compare(&core.LesserThan{})
		case template.LEQ:
			err = c.compare(&core.LesserThanOrEqualTo{})
		case template.GTR:
			err = c.compare(&core.GreaterThan{})
		case template.GEQ:
			err = c.compare(&core.GreaterThanOrEqualTo{})
		case template.LAND:
			err = c.binary(&core.AllTrue{})
		case template.LOR:
			err = c.binary(&core.AnyTrue{})
		default:
			err = errutil.Fmt("unknown operator %s", fn)
		}
	default:
		err = errutil.Fmt("unknown function %T", fn)
	}
	return
}
