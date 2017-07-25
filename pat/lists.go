package pat

import (
	"github.com/ionous/iffy/rt"
)

type Flags int

//go:generate stringer -type=Flags
const (
	Infix Flags = iota
	Prefix
	Postfix
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
type ExecutePattern struct {
	Filters Filters
	rt.Execute
	Flags
}

type BoolPatterns []BoolPattern
type NumberPatterns []NumberPattern
type TextPatterns []TextPattern
type ObjectPatterns []ObjectPattern
type NumListPatterns []NumListPattern
type TextListPatterns []TextListPattern
type ObjListPatterns []ObjListPattern
type ExecutePatterns []ExecutePattern

func (ps BoolPatterns) GetBool(run rt.Runtime) (ret bool, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.BoolEval.GetBool(run)
			break
		}
	}
	return
}
func (ps NumberPatterns) GetNumber(run rt.Runtime) (ret float64, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.NumberEval.GetNumber(run)
			break
		}
	}
	return
}
func (ps TextPatterns) GetText(run rt.Runtime) (ret string, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.TextEval.GetText(run)
			break
		}
	}
	return
}
func (ps ObjectPatterns) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.ObjectEval.GetObject(run)
			break
		}
	}
	return
}

func (ps NumListPatterns) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.NumListEval.GetNumberStream(run)
			break
		}
	}
	return
}
func (ps TextListPatterns) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.TextListEval.GetTextStream(run)
			break
		}
	}
	return
}
func (ps ObjListPatterns) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			ret, err = p.ObjListEval.GetObjectStream(run)
			break
		}
	}
	return
}

func (ps ExecutePatterns) Execute(run rt.Runtime) (ret bool, err error) {
	var post rt.ExecuteList // a stack
	var matches int
	for i, cnt := 0, len(ps); i < cnt; i++ {
		p := ps[cnt-i-1]
		if matched, e := p.Filters.GetBool(run); e != nil {
			err = e
			break
		} else if matched {
			matches++
			if p.Flags == Postfix {
				post = append(post, p.Execute)
			} else if e := p.Execute.Execute(run); e != nil {
				err = e
			} else if p.Flags != Prefix {
				break // Infix ends once its done.
			}
		}
	}
	if err == nil {
		// we want to run the most recently added thing first
		for i, cnt := 0, len(post); i < cnt; i++ {
			exec := post[cnt-i-1]
			if e := exec.Execute(run); e != nil {
				err = e
				break
			}
		}
		ret = err == nil && matches > 0
	}
	return
}
