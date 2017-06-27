package core

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
)

// Get retrieves a value from an object.
// FIX: test all forms of Get
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

func (p *Get) GetNumberStream(r rt.Runtime) (ret rt.NumberStream, err error) {
	var values []float64
	if obj, e := p.Obj.GetObject(r); e != nil {
		err = e
	} else if e := obj.GetValue(p.Value, &values); e != nil {
		err = e
	} else {
		ret = NewNumberStream(values)
	}
	return
}

func (p *Get) GetTextStream(r rt.Runtime) (ret rt.TextStream, err error) {
	var values []string
	if obj, e := p.Obj.GetObject(r); e != nil {
		err = e
	} else if e := obj.GetValue(p.Value, &values); e != nil {
		err = e
	} else {
		ret = NewTextStream(values)
	}
	return
}

// func (p *Get) GetObjectStream(r rt.Runtime) (ret rt.ObjectStream, err error) {
// 	var values []ident.Id
// 	if obj, e := p.Obj.GetObject(r); e != nil {
// 		err = e
// 	} else if e := obj.GetValue(p.Value, &values); e != nil {
// 		err = e
// 	} else {
// 		ret = rNewObjectStream(r, values)
// 	}
// 	return
// }
