package core

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
)

// Get retrieves a value from an object.
type Get struct {
	Obj   rt.ObjectEval
	Value string
}

func (p *Get) GetBool(r rt.Runtime) (ret bool, err error) {
	if obj, e := p.Obj.GetObject(r); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Value, &ret)
	}
	return
}

func (p *Get) GetNumber(r rt.Runtime) (ret float64, err error) {
	if obj, e := p.Obj.GetObject(r); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Value, &ret)
	}
	return
}

func (p *Get) GetText(r rt.Runtime) (ret string, err error) {
	if obj, e := p.Obj.GetObject(r); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Value, &ret)
	}
	return
}

func (p *Get) GetObject(r rt.Runtime) (ret ref.Object, err error) {
	if obj, e := p.Obj.GetObject(r); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Value, &ret)
	}
	return
}
