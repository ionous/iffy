package express

import (
	r "reflect"
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// Express converts a postfix expression into iffy commands.
func Convert(xs template.Expression) (ret interface{}, err error) {
	c := Converter{}
	return c.Convert(xs)
}

type Converter struct {
	stack cmdStack // the stack is empty initially, and we fill it with converted commands
	// ( to be used later by other commands )
	autoCounter int
}

func (c *Converter) Convert(xs template.Expression) (ret interface{}, err error) {
	if e := c.convert(xs); e != nil {
		err = e
	} else if op, e := c.stack.flush(); e != nil {
		err = e
	} else if on, ok := op.(*dottedName); ok {
		// if the entire template can be reduced to an dottedName
		// ex. {.lantern} then we treat it as a request for the friendly name of the object
		ret = on.getPrintedName()
	} else {
		ret = op
	}
	return
}

func (c *Converter) convert(xs template.Expression) (err error) {
	for _, fn := range xs {
		if e := c.addFunction(fn); e != nil {
			err = e
			break
		}
	}
	return
}

func (c *Converter) buildOne(cmd interface{}) {
	c.stack.push(r.ValueOf(cmd))
}

func (c *Converter) buildTwo(cmd interface{}) (err error) {
	return c.buildCommand(cmd, 2)
}

func (c *Converter) buildCommand(cmd interface{}, arity int) (err error) {
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

// fix? this is where a Scalar value could come in handy.
func (c *Converter) buildCompare(cmp core.Comparator) (err error) {
	if args, e := c.stack.pop(2); e != nil {
		err = e
	} else {
		var ptr r.Value
		a, b := unpackArg(args[0]), unpackArg(args[1])
		an, bn := a.String(), b.String() // here for debugging
		switch {
		case implements(a, b, typeNumEval):
			ptr = r.New(compareNum)
		case implements(a, b, typeTextEval):
			ptr = r.New(compareText)
		default:
			err = errutil.Fmt("unknown commands %v %v", an, bn)
		}
		if err == nil {
			cmp := r.ValueOf(cmp)
			args = []r.Value{a, cmp, b}
			if e := assignProps(ptr.Elem(), args); e != nil {
				err = e
			} else {
				c.stack.push(ptr)
			}
		}
	}
	return
}

func (c *Converter) buildSequence(cmd rt.TextEval, seq *core.Sequence, count int) (err error) {
	if args, e := c.stack.pop(count); e != nil {
		err = e
	} else {
		var parts []rt.TextEval
		for i, a := range args {
			a := unpackArg(a)
			if text, ok := a.Interface().(rt.TextEval); !ok {
				err = errutil.Fmt("couldn't convert sequence part %d to text", i)
				break
			} else {
				parts = append(parts, text)
			}
		}
		if err == nil {
			c.autoCounter++
			counter := "autoexp" + strconv.Itoa(c.autoCounter)
			// seq is part of cmd
			seq.Parts = parts
			seq.Seq = counter
			// after filling out the cmd, we push it for later processing
			c.buildOne(cmd)
		}
	}
	return
}

// build an command named in the export Slat
// names in templates are currently "mixedCase" rather than "underscore_case".
func (c *Converter) buildExport(name string, arity int) (err error) {
	if a, ok := coreCache.get(name); !ok {
		err = c.buildPattern(name, arity)
	} else if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		rtype := r.TypeOf(a).Elem()
		ptr := r.New(rtype)
		if e := assignProps(ptr.Elem(), args); e != nil {
			err = e
		} else {
			c.stack.push(ptr)
		}
	}
	return
}

func (c *Converter) buildPattern(name string, arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		var ps pattern.Arguments
		for i, arg := range args {
			if newa, e := newAssignment(arg); e != nil {
				err = errutil.Append(e)
			} else {
				newp := &pattern.Argument{
					Name: "$" + strconv.Itoa(i+1),
					From: newa,
				}
				ps.Args = append(ps.Args, newp)
			}
		}
		if err == nil {
			// printing is generally an activity b/c say is an activity
			// and we want the ability to say several things in series.
			// expressions are text patterns... so for now adapt via text
			// expressions would ideally adapt based on the pattern type
			// the assembler probably needs to work directly on tokens...
			c.buildOne(&core.Buffer{core.NewActivity(&pattern.DetermineAct{
				Pattern:   name,
				Arguments: &ps,
			})})
		}
	}
	return
}

// an eval has been passed to a pattern, return the command to assign the eval to an arg.
func newAssignment(arg r.Value) (ret core.Assignment, err error) {
	switch arg := arg.Interface().(type) {
	case *dottedName:
		ret = arg.getFromVar()
	case rt.BoolEval:
		ret = &core.FromBool{arg}
	case rt.NumberEval:
		ret = &core.FromNum{arg}
	case rt.TextEval:
		ret = &core.FromText{arg}
	case rt.NumListEval:
		ret = &core.FromNumList{arg}
	case rt.TextListEval:
		ret = &core.FromTextList{arg}
	default:
		err = errutil.Fmt("unknown pattern parameter type %T", arg)
	}
	return
}

func (c *Converter) buildUnless(cmd interface{}, arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else if len(args) > 0 {
		arg := unpackArg(args[0])
		if a, ok := arg.Interface().(rt.BoolEval); !ok {
			err = errutil.New("argument is not a bool")
		} else {
			args[0] = r.ValueOf(&core.IsNotTrue{a}) // rewrite the arg.
			c.stack.push(args...)                   //
			err = c.buildCommand(cmd, arity)
		}
	}
	return
}

func (c *Converter) buildSpan(arity int) (err error) {
	if args, e := c.stack.pop(arity); e != nil {
		err = e
	} else {
		var txts []rt.TextEval
		for _, el := range args {
			switch el := el.Interface().(type) {
			// in a list of text evaluations,
			// for example maybe "{.bennie} and the {.jets}"
			// single occurrences of dotted names are treated as requests for a friendly name
			case *dottedName:
				txts = append(txts, el.getPrintedName())
			case rt.TextEval:
				txts = append(txts, el)
			default:
				e := errutil.New("argument %T is not a text eval", el)
				err = errutil.Append(err, e)
			}
		}
		if err == nil {
			c.buildOne(&core.Join{Parts: txts})
		}
	}
	return
}

// convert the passed postfix template function into iffy commands.
func (c *Converter) addFunction(fn postfix.Function) (err error) {
	switch fn := fn.(type) {
	case types.Quote:
		txt := fn.Value()
		c.buildOne(T(txt))

	case types.Number:
		num := fn.Value()
		c.buildOne(N(num))

	case types.Bool:
		b := fn.Value()
		c.buildOne(B(b))

	case types.Command: // see decode
		err = c.buildExport(fn.CommandName, fn.CommandArity)

	case types.Reference:
		// fields are an array of strings .a.b.c
		if fields := fn.Value(); len(fields) == 0 {
			err = errutil.New("empty reference")
		} else {
			// fix: this should add ephemera that there's an object of name
			// fix: can this add ephemera that there's a local of name?
			firstField := fields[0]
			//
			var obj core.ObjectRef
			if lang.IsCapitalized(firstField) {
				obj = &core.ObjectName{T(firstField)}
			} else {
				// unboxing: get the object id from the variable named by the first dot
				// its possible though that the name requested is actually an object in the first place
				// ex. could be .ringBearer, or could be .samWise
				// lets assume for now, that if a variable holds text referring to an object....
				// then it needs to hold the object id, ie. we dont have to resolve the name of the object again.
				obj = &core.GetVar{Name: T(firstField), TryTextAsObject: true}
			}
			if len(fields) == 1 {
				// we dont know yet how { .name.... } is being used:
				// - a command arg, so the desired type is known.
				// - a pattern arg, so the desired type isn't known.
				// - a request to print an object name
				//
				// the name itself could refer to:
				// - the name of an object,
				// - the name of a pattern parameter,
				// - a loop counter,
				// - etc.
				c.buildOne(&dottedName{firstField})
			} else {
				// a chain of dots indicates we're getting one or more fields of objects
				// ex. for { .object.fieldContainingAnObject.otherField }
				var getField *core.GetField
				// .a.b: from the named object a, we want its field b
				// .a.b.c: after getting the object name in field b, get that object's field c
				for _, field := range fields[1:] {
					// the first time through, we already have an object id
					// on subsequent loops we turn the results of the previous GetField
					// into a request for that object's name.
					if getField != nil {
						obj = &core.ObjectName{getField}
					}
					//
					getField = &core.GetField{
						Obj:   obj,
						Field: T(field),
					}
				}
				c.buildOne(getField)
			}
		}

	case types.Builtin:
		switch k := fn.Type; k {
		case types.IfStatement:
			// it would be nice if this could be choose text or choose number based on context
			// choose scalar might simplify things....
			err = c.buildCommand(&core.ChooseText{}, fn.ParameterCount)
		case types.UnlessStatement:
			err = c.buildUnless(&core.ChooseText{}, fn.ParameterCount)

		case types.Stopping:
			var seq core.StoppingText
			err = c.buildSequence(&seq, &seq.Sequence, fn.ParameterCount)
		case types.Shuffle:
			var seq core.ShuffleText
			err = c.buildSequence(&seq, &seq.Sequence, fn.ParameterCount)
		case types.Cycle:
			var seq core.CycleText
			err = c.buildSequence(&seq, &seq.Sequence, fn.ParameterCount)
		case types.Span:
			err = c.buildSpan(fn.ParameterCount)

		default:
			// fix? span is supposed to join text sections.... but there were no tests or examples in the og code.
			err = errutil.New("unhandled builtin", k.String())
		}

	case types.Operator:
		switch fn {
		case types.MUL:
			err = c.buildTwo(&core.ProductOf{})
		case types.QUO:
			err = c.buildTwo(&core.QuotientOf{})
		case types.REM:
			err = c.buildTwo(&core.RemainderOf{})
		case types.ADD:
			err = c.buildTwo(&core.SumOf{})
		case types.SUB:
			err = c.buildTwo(&core.DiffOf{})

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
			err = c.buildTwo(&core.AllTrue{})
		case types.LOR:
			err = c.buildTwo(&core.AnyTrue{})
		default:
			err = errutil.Fmt("unknown operator %s", fn)
		}

	default:
		err = errutil.Fmt("unknown function %T", fn)
	}
	return
}
