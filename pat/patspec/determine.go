package patspec

import (
	"github.com/ionous/iffy/rt"
)

type Determine struct {
	Obj rt.ObjectEval
}

// GetBool implements rt.BoolEval
func (p *Determine) GetBool(run rt.Runtime) (ret bool, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = run.GetBoolMatching(data)
	}
	return
}

// GetNumber implements rt.NumberEval
func (p *Determine) GetNumber(run rt.Runtime) (ret float64, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = run.GetNumMatching(data)
	}
	return
}

// GetText implements rt.TextEval
func (p *Determine) GetText(run rt.Runtime) (ret string, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = run.GetTextMatching(data)
	}
	return
}

// GetObject implements rt.ObjectEval
func (p *Determine) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = run.GetObjectMatching(data)
	}
	return
}

// GetNumberStream implements rt.NumListEval
func (p *Determine) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = run.GetNumStreamMatching(data)
	}
	return
}

// GetTextStream implements rt.TextListEval
func (p *Determine) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = run.GetTextStreamMatching(data)
	}
	return
}

// GetObjectStream implements rt.ObjListEval
func (p *Determine) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = e
	} else {
		ret, err = run.GetObjStreamMatching(data)
	}
	return
}
