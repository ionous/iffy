package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
	r "reflect"
)

// GetAt a named property in the current top object;
// or, if the property isn't found, look for an object of that name instead.
// Acts similar to variable name resolution in other languages.
type GetAt struct {
	Name string
}

func (p *GetAt) GetBool(run rt.Runtime) (ret bool, err error) {
	err = p.getValue(run, &ret)
	return
}

func (p *GetAt) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = p.getValue(run, &ret)
	return
}

func (p *GetAt) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if e := p.getValue(run, &ret); e != nil {
		err = e
	} else if obj, ok := run.GetObject(p.Name); !ok {
		err = e
	} else {
		ret = obj
	}
	return
}

func (p *GetAt) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	var values []float64
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewNumberStream(stream.FromList(values))
	}
	return
}

func (p *GetAt) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	var values []string
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewTextStream(stream.FromList(values))
	}
	return
}

func (p *GetAt) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	var values []rt.Object
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewObjectStream(stream.FromList(values))
	}
	return
}

func (p *GetAt) getValue(run rt.Runtime, pv interface{}) (err error) {
	if obj, path, e := getAt(run, p.Name); e != nil {
		err = e
	} else if len(path) > 0 {
		err = obj.GetValue(p.Name, pv)
	} else {
		err = run.Pack(r.ValueOf(pv), r.ValueOf(obj))
	}
	return
}

func (p *GetAt) GetText(run rt.Runtime) (ret string, err error) {
	if obj, path, e := getAt(run, p.Name); e != nil {
		err = e
	} else {
		ret, err = textConvert(run, obj, path)
	}
	return
}

// getAt looks at the top object for a property of the passed name,
// or failing that, an object of the passed name.
func getAt(run rt.Runtime, name string) (ret rt.Object, at []int, err error) {
	top, topOk := run.TopObject()
	if topOk {
		if path := class.PropertyPath(top.Type(), name); len(path) > 0 {
			ret, at = top, path
		}
	}
	if ret == nil {
		if obj, ok := run.GetObject(name); !ok {
			if !topOk {
				err = errutil.New("couldnt find an object named", name)
			} else {
				err = errutil.Fmt("couldnt find a property or object named '%s' in context %v.", name, top)
			}

		} else {
			ret = obj
		}
	}
	return
}
