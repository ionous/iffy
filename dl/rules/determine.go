package rules

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type Determine struct{ Obj rt.ObjectEval }

// GetBool implements rt.BoolEval
func (p *Determine) GetBool(run rt.Runtime) (ret bool, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine bool, because", e)
	} else if r, e := run.GetBoolMatching(data); e != nil {
		err = errutil.New("couldnt determine bool, because", e)
	} else {
		ret = r
	}
	return
}

// GetNumber implements rt.NumberEval
func (p *Determine) GetNumber(run rt.Runtime) (ret float64, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine num, because", e)
	} else if r, e := run.GetNumMatching(data); e != nil {
		err = errutil.New("couldnt determine num, because", e)
	} else {
		ret = r
	}
	return
}

// GetText implements rt.TextEval
func (p *Determine) GetText(run rt.Runtime) (ret string, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine text, because", e)
	} else if r, e := run.GetTextMatching(data); e != nil {
		err = errutil.New("couldnt determine text, because", e)
	} else {
		ret = r
	}
	return
}

// GetObject implements rt.ObjectEval
func (p *Determine) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine obj, because", e)
	} else if r, e := run.GetObjectMatching(data); e != nil {
		err = errutil.New("couldnt determine obj, because", e)
	} else {
		ret = r
	}
	return
}

// GetNumberStream implements rt.NumListEval
func (p *Determine) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine nums, because", e)
	} else if r, e := run.GetNumStreamMatching(data); e != nil {
		err = errutil.New("couldnt determine nums, because", e)
	} else {
		ret = r
	}
	return
}

// GetTextStream implements rt.TextListEval
func (p *Determine) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine texts, because", e)
	} else if r, e := run.GetTextStreamMatching(data); e != nil {
		err = errutil.New("couldnt determine texts, because", e)
	} else {
		ret = r
	}
	return
}

// GetObjectStream implements rt.ObjListEval
func (p *Determine) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine objs, because", e)
	} else if r, e := run.GetObjStreamMatching(data); e != nil {
		err = errutil.New("couldnt determine objs, because", e)
	} else {
		ret = r
	}
	return
}
func (p *Determine) Execute(run rt.Runtime) (err error) {
	if data, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("couldnt determine exec, because", e)
	} else if e := run.ExecuteMatching(data); e != nil {
		err = errutil.New("couldnt determine exec, because", e)
	}
	return
}
