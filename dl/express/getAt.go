package express

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
	r "reflect"
)

// GetAt retrieves a value from the scope object.
// For GetAt.GetObject, if the value is not found, it tries the value from globals.
type GetAt struct {
	Value string
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
		if obj, ok := run.GetObject(p.Value); !ok {
			err = e
		} else {
			ret = obj
		}
	}
	return
}

func (p *GetAt) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	var values []float64
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewNumberStream(values)
	}
	return
}

func (p *GetAt) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	var values []string
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewTextStream(values)
	}
	return
}

func (p *GetAt) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	var values []rt.Object
	if e := p.getValue(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewObjectStream(values)
	}
	return
}

func (p *GetAt) getValue(run rt.Runtime, pv interface{}) (err error) {
	if obj, path, e := getAt(run, p.Value); e != nil {
		err = e
	} else if len(path) > 0 {
		err = obj.GetValue(p.Value, pv)
	} else {
		err = run.Pack(r.ValueOf(pv), r.ValueOf(obj))
	}
	return
}

func (p *GetAt) GetText(run rt.Runtime) (ret string, err error) {
	if obj, path, e := getAt(run, p.Value); e != nil {
		err = e
	} else {
		ret, err = textConvert(run, obj, path)
	}
	return
}

// getAt looks at the top object for a property of the passed name,
// or failing that, an object of the passed name.
func getAt(run rt.Runtime, name string) (ret rt.Object, at []int, err error) {
	var found bool
	if obj, ok := run.TopObject(); ok {
		if path := class.PropertyPath(obj.Type(), name); len(path) > 0 {
			ret, at = obj, path
			found = true
		}
	}
	if !found {
		if obj, ok := run.GetObject(name); !ok {
			err = errutil.New("couldnt find an object or property named", name)
		} else {
			ret = obj
		}
	}
	return
}
