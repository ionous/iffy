package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
)

// Convert converts a postfix expression into iffy commands.
func Convert(ops *ops.Factory, expression postfix.Expression) (ret interface{}, err error) {
	c := converter{ops: ops}
	if e := c.convert(expression); e != nil {
	} else if len(c.stack) == 0 {
		err = errutil.New("empty output")
	} else if len(c.stack) > 1 {
		err = errutil.New("unparsed output")
	} else {
		ret = c.stack[0]
	}
	return
}

type converter struct {
	ops   *ops.Factory
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
func (c *converter) push(cmd interface{}) {
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
	} else if cmd, e := c.ops.CmdFromPointer(i); e != nil {
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
	} else if cmd, e := c.ops.CmdFromPointer(&core.CompareText{}); e != nil {
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
		c.push(cmd.Target().Interface())
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
	case chart.Quote:
		op := &core.Text{fn.Value()}
		c.push(op)

	case chart.Number:
		op := &core.Num{fn.Value()}
		c.push(op)

	case chart.Reference:
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
			c.push(op)
		}

	case chart.Command:
		if cmd, e := c.ops.CmdFromName(fn.CommandName); e != nil {
			err = e
		} else if args, e := c.pop(fn.CommandArity); e != nil {
			err = e
		} else {
			err = c.pushCommand(cmd, args...)
		}

	case chart.Operator:
		switch fn {
		case chart.MUL:
			err = c.binary(&core.Mul{})
		case chart.QUO:
			err = c.binary(&core.Div{})
		case chart.REM:
			err = c.binary(&core.Mod{})
		case chart.ADD:
			err = c.binary(&core.Add{})
		case chart.SUB:
			err = c.binary(&core.Sub{})
		case chart.EQL:
			err = c.compare(&core.EqualTo{})
		case chart.NEQ:
			err = c.compare(&core.NotEqualTo{})
		case chart.LSS:
			err = c.compare(&core.LesserThan{})
		case chart.LEQ:
			err = c.compare(&core.LesserThanOrEqualTo{})
		case chart.GTR:
			err = c.compare(&core.GreaterThan{})
		case chart.GEQ:
			err = c.compare(&core.GreaterThanOrEqualTo{})
		case chart.LAND:
			err = c.binary(&core.AllTrue{})
		case chart.LOR:
			err = c.binary(&core.AnyTrue{})
		default:
			err = errutil.Fmt("unknown operator %s", fn)
		}
	default:
		err = errutil.Fmt("unknown function %T", fn)
	}
	return
}
