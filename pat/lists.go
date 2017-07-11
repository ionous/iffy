package pat

import (
	"github.com/ionous/iffy/rt"
)

type BoolPattern struct {
	Filters Filters
	rt.BoolEval
}
type NumberPattern struct {
	Filters Filters
	rt.NumberEval
}
type TextPattern struct {
	Filters Filters
	rt.TextEval
}
type ObjectPattern struct {
	Filters Filters
	rt.ObjectEval
}
type NumListPattern struct {
	Filters Filters
	rt.NumListEval
}
type TextListPattern struct {
	Filters Filters
	rt.TextListEval
}
type ObjListPattern struct {
	Filters Filters
	rt.ObjListEval
}

type BoolPatterns []BoolPattern
type NumberPatterns []NumberPattern
type TextPatterns []TextPattern
type ObjectPatterns []ObjectPattern
type NumListPatterns []NumListPattern
type TextListPatterns []TextListPattern
type ObjListPatterns []ObjListPattern

func (p BoolPatterns) GetBool(run rt.Runtime) (ret bool, err error) {
	for i, cnt := 0, len(p); i < cnt; i++ {
		v := p[cnt-i-1]
		if matched, e := v.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = v.BoolEval.GetBool(run)
			break
		}
	}
	return
}
func (p NumberPatterns) GetNumber(run rt.Runtime) (ret float64, err error) {
	for i, cnt := 0, len(p); i < cnt; i++ {
		v := p[cnt-i-1]
		if matched, e := v.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = v.NumberEval.GetNumber(run)
			break
		}
	}
	return
}
func (p TextPatterns) GetText(run rt.Runtime) (ret string, err error) {
	for i, cnt := 0, len(p); i < cnt; i++ {
		v := p[cnt-i-1]
		if matched, e := v.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = v.TextEval.GetText(run)
			break
		}
	}
	return
}
func (p ObjectPatterns) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	for i, cnt := 0, len(p); i < cnt; i++ {
		v := p[cnt-i-1]
		if matched, e := v.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = v.ObjectEval.GetObject(run)
			break
		}
	}
	return
}
func (p NumListPatterns) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	for i, cnt := 0, len(p); i < cnt; i++ {
		v := p[cnt-i-1]
		if matched, e := v.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = v.NumListEval.GetNumberStream(run)
			break
		}
	}
	return
}
func (p TextListPatterns) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	for i, cnt := 0, len(p); i < cnt; i++ {
		v := p[cnt-i-1]
		if matched, e := v.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = v.TextListEval.GetTextStream(run)
			break
		}
	}
	return
}
func (p ObjListPatterns) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	for i, cnt := 0, len(p); i < cnt; i++ {
		v := p[cnt-i-1]
		if matched, e := v.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = v.ObjListEval.GetObjectStream(run)
			break
		}
	}
	return
}
