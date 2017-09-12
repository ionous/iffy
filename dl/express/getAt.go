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
	Prop string
}

func (p *GetAt) GetBool(run rt.Runtime) (ret bool, err error) {
	err = p.get(run, &ret)
	return
}

func (p *GetAt) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = p.get(run, &ret)
	return
}

func (p *GetAt) GetText(run rt.Runtime) (ret string, err error) {
	err = p.get(run, &ret)
	return
}

func (p *GetAt) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if e := p.get(run, &ret); e != nil {
		if obj, ok := run.FindObject(p.Prop); !ok {
			err = e
		} else {
			ret = obj
		}
	}
	return
}

func (p *GetAt) get(run rt.Runtime, pv interface{}) (err error) {
	var found bool
	if obj, ok := run.FindObject("@"); ok {
		if path := class.PropertyPath(obj.Type(), p.Prop); len(path) > 0 {
			found = true
			err = obj.GetValue(p.Prop, pv)
		}
	}
	if !found {
		if obj, ok := run.FindObject(p.Prop); !ok {
			err = errutil.New("couldnt find an object or property named", p.Prop)
		} else {
			err = run.Pack(r.ValueOf(pv), r.ValueOf(obj))
		}
	}
	return
}

func (p *GetAt) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	var values []float64
	if e := p.get(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewNumberStream(values)
	}
	return
}

func (p *GetAt) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	var values []string
	if e := p.get(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewTextStream(values)
	}
	return
}

func (p *GetAt) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	var values []rt.Object
	if e := p.get(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewObjectStream(values)
	}
	return
}
