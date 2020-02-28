package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
	r "reflect"
)

// Express converts a postfix expression into iffy commands.
func Convert(cmds *ops.Factory, xs template.Expression, gen names) (ret *ops.Command, err error) {
	c := converter{cmds: cmds, gen: gen}
	if e := c.convert(xs); e != nil {
		err = e
	} else if len(c.stack) == 0 {
		err = errutil.New("empty output")
	} else if len(c.stack) > 1 {
		err = errutil.New("unparsed output")
	} else if cmd := c.stack[0]; cmd == nil {
		err = errutil.New("convert returned nil")
	} else {
		ret = cmd
	}
	return
}

type converter struct {
	cmds  *ops.Factory
	stack []*ops.Command
	gen   names
}

func (c *converter) convert(xs template.Expression) (err error) {
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
func (c *converter) pop(cnt int) (ret []*ops.Command, err error) {
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
	} else if cmp, e := c.cmds.EmplaceCommand(cmp); e != nil {
		err = e
	} else {
		a, b := args[0], args[1]
		if cmd := compare(a.Target(), b.Target()); cmd == nil {
			err = errutil.Fmt("cant compare %T to %T", a, b)
		} else if cmd, e := c.cmds.EmplaceCommand(cmd); e != nil {
			err = e
		} else {
			// as an alternative to this custom code, we could name arguments for every command.
			// ex. MUL: {cmd:Mul, args: "A", "B" }, CompareText[ A, B ]
			// wed still have to have an initalizer or something for "cmp".
			err = c.pushCommand(cmd, a, cmp, b)

		}
	}
	return
}

func compare(a, b r.Value) (ret interface{}) {
	try := []struct {
		Test func(rtype r.Type) bool
		Res  func() interface{}
	}{
		{kindOf.NumberEval, func() interface{} { return &core.CompareNum{} }},
		{kindOf.ObjectEval, func() interface{} { return &core.CompareObj{} }},
		{kindOf.TextEval, func() interface{} { return &core.CompareText{} }},
	}
	at, bt := a.Type(), b.Type()
	for _, x := range try {
		if x.Test(at) && x.Test(bt) {
			ret = x.Res()
			break
		}
	}
	return
}

// add the passed args to the passed command in order; then push the command.
func (c *converter) pushCommand(cmd *ops.Command, args ...*ops.Command) (err error) {
	if e := assign(cmd, args); e != nil {
		err = e
	} else if cmd, e := c.determine(cmd); e != nil {
		err = e
	} else {
		c.push(cmd)
	}
	return
}

// FIX? itd probably be better to know which commands are real and which are shadow/patterns;
// and just create as neede.
func (c *converter) determine(cmd *ops.Command) (ret *ops.Command, err error) {
	if tgt := cmd.Target(); tgt.Type() != r.TypeOf((*ops.ShadowClass)(nil)) {
		ret = cmd
	} else if det, e := c.cmds.CreateCommand("determine"); e != nil {
		err = e
	} else if e := det.Position(cmd); e != nil {
		err = e
	} else {
		ret = det
	}
	return
}

// add the passed args to the passed command in order.
func assign(cmd *ops.Command, args []*ops.Command) (err error) {
	for _, arg := range args {
		if e := cmd.Position(arg); e != nil {
			err = errutil.Fmt("couldnt assign %s to %s, because %s", arg, cmd, e)
			break
		}
	}
	return
}

func (c *converter) create(name string, args int) (err error) {
	if cmd, e := c.cmds.CreateCommand(name); e != nil {
		err = e
	} else if args, e := c.pop(args); e != nil {
		err = e
	} else {
		err = c.pushCommand(cmd, args...)
	}
	return
}

// name is "shuffle", etc.
// c.Cmd("say", c.Cmd("shuffle text",
// 		gen.NewName("shuffle counter"),
// 		sliceOf.String("a", "b", "c"),
// 	))
func (c *converter) sequence(name string, args int) (err error) {
	text := name + " text"
	counter := name + " counter"
	if cmd, e := c.cmds.CreateCommand(text); e != nil {
		err = e
	} else {
		name := c.gen.NewName(counter)
		if e := cmd.Position(name); e != nil {
			err = e
		} else if args, e := c.pop(args); e != nil {
			err = e
		} else {
			err = c.pushCommand(cmd, args...)
		}
	}
	return
}

// convert the passed function into iffy commands.
func (c *converter) addFunction(fn postfix.Function) (err error) {
	switch fn := fn.(type) {
	case types.Quote:
		if cmd, e := c.cmds.EmplaceCommand(&core.TextValue{fn.Value()}); e != nil {
			err = e
		} else {
			c.push(cmd)
		}

	case types.Number:
		if cmd, e := c.cmds.EmplaceCommand(&core.NumValue{fn.Value()}); e != nil {
			err = e
		} else {
			c.push(cmd)
		}

	case types.Reference:
		if fields := fn.Value(); len(fields) == 0 {
			err = errutil.New("empty reference")
		} else {
			// obj.a.b.c => Get{c Get{b Get{a GetAt{obj}}}}
			var op rt.ObjectEval
			if name := fields[0]; lang.IsCapitalized(name) {
				op = &core.ObjectName{name}
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

	case types.Command:
		err = c.create(fn.CommandName, fn.CommandArity)

	case types.Builtin:
		switch k := fn.Type; k {
		case types.Span:
			err = c.create("join", fn.ParameterCount)
		case types.Stopping:
			err = c.sequence("stopping", fn.ParameterCount)
		case types.Shuffle:
			err = c.sequence("shuffle", fn.ParameterCount)
		case types.Cycle:
			err = c.sequence("cycle", fn.ParameterCount)
		case types.IfStatement:
			// choose by hint?
			err = c.create("choose text", fn.ParameterCount)
		case types.UnlessStatement:
			if args, e := c.pop(fn.ParameterCount); e != nil {
			} else if cmd, e := c.cmds.CreateCommand("is not"); e != nil {
				err = e
			} else if e := cmd.Position(args[0]); e != nil {
				err = e
			} else {
				args[0] = cmd
				if cmd, e := c.cmds.CreateCommand("choose text"); e != nil {
					err = e
				} else {
					err = c.pushCommand(cmd, args...)
				}
			}
		}

	case types.Operator:
		switch fn {
		case types.MUL:
			err = c.binary(&core.ProductOf{})
		case types.QUO:
			err = c.binary(&core.QuotientOf{})
		case types.REM:
			err = c.binary(&core.RemainderOf{})
		case types.ADD:
			err = c.binary(&core.SumOf{})
		case types.SUB:
			err = c.binary(&core.DiffOf{})
		case types.EQL:
			err = c.compare(&core.EqualTo{})
		case types.NEQ:
			err = c.compare(&core.NotEqualTo{})
		case types.LSS:
			err = c.compare(&core.LessThan{})
		case types.LEQ:
			err = c.compare(&core.LessOrEqual{})
		case types.GTR:
			err = c.compare(&core.GreaterThan{})
		case types.GEQ:
			err = c.compare(&core.GreaterOrEqual{})
		case types.LAND:
			err = c.binary(&core.AllTrue{})
		case types.LOR:
			err = c.binary(&core.AnyTrue{})
		default:
			err = errutil.Fmt("unknown operator %s", fn)
		}
	default:
		err = errutil.Fmt("unknown function %T", fn)
	}
	return
}
