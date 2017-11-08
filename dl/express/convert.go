package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
)

// ParseExpression converts a postfix expression into iffy commands.
func ParseExpression(ops *ops.Factory, expression postfix.Expression) (err error) {
	p := parser{ops: ops}
	for _, fn := range expression {
		if e := p.addFunction(fn); e != nil {
			err = e
			break
		}
	}
	return
}

type parser struct {
	ops   *ops.Factory
	stack []interface{}
}

// add a new command pointer to the output stack.
func (p *parser) push(cmd interface{}) {
	p.stack = append(p.stack, cmd)
}

// extract cnt commands from the parser stack.
func (p *parser) pop(cnt int) (ret []interface{}, err error) {
	if end := len(p.stack) - cnt; end < 0 {
		err = errutil.New("stack underflow")
	} else {
		ret, p.stack = p.stack[end:], p.stack[:end]
	}
	return
}

func (p *parser) binary(i interface{}) (err error) {
	if args, e := p.pop(2); e != nil {
		err = e
	} else if cmd, e := p.ops.CmdFromPointer(i); e != nil {
		err = e
	} else {
		err = p.pushCommand(cmd, args...)
	}
	return
}

// adds a comparision statement to the output.
// FIX FIX we have to know the type: num,text,obj of the properties in question
// need more information on properties.
func (p *parser) compare(cmp core.CompareTo) (err error) {
	if args, e := p.pop(2); e != nil {
		err = e
	} else if cmd, e := p.ops.CmdFromPointer(&core.CompareText{}); e != nil {
		err = e
	} else {
		// as an alternative to this custom code, we could name arguments for every command.
		// ex. MUL: {cmd:Mul, args: "A", "B" }, CompareText[ A, B ]
		// wed still have to have an initalizer or something for "cmp".
		err = p.pushCommand(cmd, args[0], cmp, args[1])
	}
	return
}

// add the passed args to the passed command in order; then push the command.
func (p *parser) pushCommand(cmd *ops.Command, args ...interface{}) (err error) {
	if e := assign(cmd, args); e != nil {
		err = e
	} else {
		p.push(cmd.Target().Interface())
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
func (p *parser) addFunction(fn postfix.Function) (err error) {
	switch fn := fn.(type) {
	case chart.Quote:
		op := &core.Text{fn.Value()}
		p.push(op)

	case chart.Number:
		op := &core.Num{fn.Value()}
		p.push(op)

	case chart.Reference:
		if fields := fn.Value(); len(fields) == 0 {
			err = errutil.New("empty reference")
		} else {
			// obj.a.b.c => Get{c Get{b Get{a GetAt{obj}}}}
			var op rt.ObjectEval = &GetAt{fields[0]}
			for _, field := range fields[1:] {
				op = &Render{op, field}
			}
			p.push(op)
		}

	case chart.Command:
		if cmd, e := p.ops.CmdFromName(fn.CommandName); e != nil {
			err = e
		} else if args, e := p.pop(fn.CommandArity); e != nil {
			err = e
		} else {
			err = p.pushCommand(cmd, args...)
		}

	case chart.Operator:
		switch fn {
		case chart.MUL:
			err = p.binary(&core.Mul{})
		case chart.QUO:
			err = p.binary(&core.Div{})
		case chart.REM:
			err = p.binary(&core.Mod{})
		case chart.ADD:
			err = p.binary(&core.Add{})
		case chart.SUB:
			err = p.binary(&core.Sub{})
		case chart.EQL:
			err = p.compare(&core.EqualTo{})
		case chart.NEQ:
			err = p.compare(&core.NotEqualTo{})
		case chart.LSS:
			err = p.compare(&core.LesserThan{})
		case chart.LEQ:
			err = p.compare(&core.LesserThanOrEqualTo{})
		case chart.GTR:
			err = p.compare(&core.GreaterThan{})
		case chart.GEQ:
			err = p.compare(&core.GreaterThanOrEqualTo{})
		case chart.LAND:
			err = p.binary(&core.AllTrue{})
		case chart.LOR:
			err = p.binary(&core.AnyTrue{})
		default:
			err = errutil.Fmt("unknown operator %s", fn)
		}
	default:
		err = errutil.Fmt("unknown function %T", fn)
	}
	return
}
