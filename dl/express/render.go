package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
	"strconv"
)

// Render command extracts the value represented by an object's property eval.
// As a special case, templates can almost always be rendered to text,
// allowing most eval properties to be used where a text eval is needed.
type Render struct {
	Obj  rt.ObjectEval
	Prop string
}

// GetText
func (p *Render) GetText(run rt.Runtime) (ret string, err error) {
	return getText(run, p.Obj, p.Prop)
}

func (p *Render) GetBool(run rt.Runtime) (ret bool, err error) {
	err = p.getValue(run, &ret)
	return
}

func (p *Render) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = p.getValue(run, &ret)
	return
}

func (p *Render) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	err = p.getValue(run, &ret)
	return
}

func (p *Render) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	var values []float64
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewNumberStream(stream.FromList(values))
	}
	return
}

func (p *Render) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	var values []string
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewTextStream(stream.FromList(values))
	}
	return
}

func (p *Render) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	var values []rt.Object
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewObjectStream(stream.FromList(values))
	}
	return
}

// get the object property
func (p *Render) getValue(run rt.Runtime, pv interface{}) (err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New(p.Prop, "failed rendering because,", e)
	} else {
		err = obj.GetValue(p.Prop, pv)
	}
	return
}

func getText(run rt.Runtime, obj rt.ObjectEval, prop string) (ret string, err error) {
	if obj, e := obj.GetObject(run); e != nil {
		err = errutil.New(prop, "failed rendering because,", e)
	} else {
		cls := obj.Type()
		if path := class.PropertyPath(cls, prop); len(path) == 0 {
			err = errutil.New("property not found", prop)
		} else {
			ret, err = textConvert(run, obj, path)
		}
	}
	return
}

func textConvert(run rt.Runtime, obj rt.Object, path []int) (ret string, err error) {
	if len(path) == 0 {
		ret, err = getName(run, obj.Id())
	} else {
		field := obj.Type().FieldByIndex(path)
		switch ft := field.Type; {
		default:
			err = obj.GetValue(field.Name, &ret)

		case kindOf.BoolLike(ft):
			var v bool
			if e := obj.GetValue(field.Name, &v); e != nil {
				err = e
			} else {
				ret = strconv.FormatBool(v)
			}

		case kindOf.NumberLike(ft):
			var v float64
			if e := obj.GetValue(field.Name, &v); e != nil {
				err = e
			} else {
				ret = strconv.FormatFloat(v, 'g', -1, 64)
			}

		case kindOf.ObjectLike(ft):
			var v ident.Id
			if e := obj.GetValue(field.Name, &v); e != nil {
				err = e
			} else {
				ret, err = getName(run, v)
			}

		}
	}
	return
}

// this would be nice -- but express shouldnt depend on std.
// XXX - getName returns the printed name of an object.
func getName(run rt.Runtime, id ident.Id) (ret string, err error) {
	ret = id.String()
	// 	var span printer.Span
	// 	if e := rt.WritersBlock(run, &span, func() error {
	// 		return rt.Determine(run, &std.PrintName{id})
	// 	}); e != nil {
	// 		err = e
	// 	} else {
	// 		ret = span.String()
	// 	}
	return
}
