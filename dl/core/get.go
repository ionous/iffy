package core

import (
	"github.com/ionous/iffy/rt"
)

// Get retrieves a value from an object.
// FIX: test all forms of Get/Set
type Get struct {
	Obj  rt.ObjectEval
	Prop string
}

func (p *Get) GetBool(run rt.Runtime) (ret bool, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Prop, &ret)
	}
	return
}

func (p *Get) GetNumber(run rt.Runtime) (ret float64, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Prop, &ret)
	}
	return
}

func (p *Get) GetText(run rt.Runtime) (ret string, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Prop, &ret)
	}
	return
}

func (p *Get) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		err = obj.GetValue(p.Prop, &ret)
	}
	return
}

func (p *Get) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	var values []float64
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else if e := obj.GetValue(p.Prop, &values); e != nil {
		err = e
	} else {
		ret = NewNumberStream(values)
	}
	return
}

func (p *Get) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	var values []string
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else if e := obj.GetValue(p.Prop, &values); e != nil {
		err = e
	} else {
		ret = NewTextStream(values)
	}
	return
}

// func (p *Get) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
// 	var values []ident.Id
// 	if obj, e := p.Obj.GetObject(run); e != nil {
// 		err = e
// 	} else if e := obj.GetValue(p.Prop, &values); e != nil {
// 		err = e
// 	} else {
// 		ret = NewObjectStream(run, values)
// 	}
// 	return
// }
