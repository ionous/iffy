package express

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// Express converts a postfix expression into iffy commands.
func Convert(xs template.Expression) (ret interface{}, err error) {
	c := converter{}
	if e := c.convert(xs); e != nil {
		err = e
	} else {
		ret, err = c.stack.flush()
	}
	return
}

type converter struct {
	stack rstack
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
func (c *converter) buildCommand(cmd interface{}, arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		ptr := r.ValueOf(cmd)
		if e := assignProps(ptr.Elem(), args); e != nil {
			err = e
		} else {
			c.stack.push(ptr)
		}
	}
	return
}

func (c *converter) buildBinary(cmd interface{}) (err error) {
	return c.buildCommand(cmd, 2)
}

var typeNumEval = r.TypeOf((*rt.NumberEval)(nil)).Elem()
var typeTextEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var compareNum = r.TypeOf((*core.CompareNum)(nil)).Elem()
var compareText = r.TypeOf((*core.CompareText)(nil)).Elem()

func implements(a, b r.Value, t r.Type) bool {
	return a.Type().Implements(t) && b.Type().Implements(t)
}

// fix? this is where a Scalar value could come in handy.
func (c *converter) buildCompare(cmp core.Comparator) (err error) {
	if args, e := c.stack.pop(2); e != nil {
		err = e
	} else {
		var ptr r.Value
		switch a, b := args[0], args[1]; {
		case implements(a, b, typeNumEval):
			ptr = r.New(compareNum)
		case implements(a, b, typeTextEval):
			ptr = r.New(compareText)
		default:
			err = errutil.New("unknown commands")
		}
		if err == nil {
			cmp := r.ValueOf(cmp)
			args = []r.Value{args[0], cmp, args[1]}
			if e := assignProps(ptr.Elem(), args); e != nil {
				err = e
			} else {
				c.stack.push(ptr)
			}
		}

	}
	return
}

func (c *converter) buildCmd(cmd interface{}) {
	c.stack.push(r.ValueOf(cmd))
}

// convert the passed function into iffy commands.
func (c *converter) addFunction(fn postfix.Function) (err error) {
	switch fn := fn.(type) {
	case types.Quote:
		txt := fn.Value()
		c.buildCmd(&core.Text{txt})

	case types.Number:
		num := fn.Value()
		c.buildCmd(&core.Number{num})

	case types.Operator:
		switch fn {
		case types.MUL:
			err = c.buildBinary(&core.ProductOf{})
		case types.QUO:
			err = c.buildBinary(&core.QuotientOf{})
		case types.REM:
			err = c.buildBinary(&core.RemainderOf{})
		case types.ADD:
			err = c.buildBinary(&core.SumOf{})
		case types.SUB:
			err = c.buildBinary(&core.DiffOf{})

		case types.EQL:
			err = c.buildCompare(&core.EqualTo{})
		case types.NEQ:
			err = c.buildCompare(&core.NotEqualTo{})
		case types.LSS:
			err = c.buildCompare(&core.LessThan{})
		case types.LEQ:
			err = c.buildCompare(&core.LessOrEqual{})
		case types.GTR:
			err = c.buildCompare(&core.GreaterThan{})
		case types.GEQ:
			err = c.buildCompare(&core.GreaterOrEqual{})

		case types.LAND:
			err = c.buildBinary(&core.AllTrue{})
		case types.LOR:
			err = c.buildBinary(&core.AnyTrue{})
		default:
			err = errutil.Fmt("unknown operator %s", fn)
		}
	default:
		err = errutil.Fmt("unknown function %T", fn)
	}
	return
}
