package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/ref/kindOf"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"strconv"
)

type Render struct {
	Obj  rt.ObjectEval
	Prop string
}

func (p *Render) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		cls, prop := obj.Type(), p.Prop
		if path := class.PropertyPath(cls, prop); len(path) == 0 {
			err = errutil.New("property not found", prop)
		} else {
			field := cls.FieldByIndex(path)

			switch ft := field.Type; {
			default:
				err = obj.GetValue(prop, &ret)

			case kindOf.BoolLike(ft):
				var v bool
				if e := obj.GetValue(prop, &v); e != nil {
					err = e
				} else {
					ret = strconv.FormatBool(v)
				}

			case kindOf.NumberLike(ft):
				var v float64
				if e := obj.GetValue(prop, &v); e != nil {
					err = e
				} else {
					ret = strconv.FormatFloat(v, 'g', -1, 64)
				}

			case kindOf.ObjectLike(ft):
				var v ident.Id
				if e := obj.GetValue(prop, &v); e != nil {
					err = e
				} else {
					var span printer.Span
					printName := run.Emplace(&PrintName{v})
					run := rt.Writer(run, &span)
					if e := run.ExecuteMatching(run, printName); e != nil {
						err = e
					} else {
						ret = span.String()
					}
				}
			}
		}
	}
	return
}
