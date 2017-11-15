package core

import (
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// Get a property value from an object by name.
// FIX: test all forms of Get/Set
type Get struct {
	Obj  rt.ObjectEval
	Prop string
}

func (p *Get) GetBool(run rt.Runtime) (ret bool, err error) {
	err = p.get(run, &ret)
	return
}

func (p *Get) GetNumber(run rt.Runtime) (ret float64, err error) {
	err = p.get(run, &ret)
	return
}

func (p *Get) GetText(run rt.Runtime) (ret string, err error) {
	err = p.get(run, &ret)
	return
}

func (p *Get) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	err = p.get(run, &ret)
	return
}

func (p *Get) get(run rt.Runtime, pv interface{}) (err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Prop, pv)
	}
	return
}

func (p *Get) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	var values []float64
	if e := p.get(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewNumberStream(stream.FromList(values))
	}
	return
}

func (p *Get) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	var values []string
	if e := p.get(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewTextStream(stream.FromList(values))
	}
	return
}

func (p *Get) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	var values []rt.Object
	if e := p.get(run, &values); e != nil {
		err = e
	} else {
		ret = stream.NewObjectStream(stream.FromList(values))
	}
	return
}
