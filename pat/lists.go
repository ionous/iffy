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

type BoolRule struct {
	Filters Filters
	rt.BoolEval
}
type NumberRule struct {
	Filters Filters
	rt.NumberEval
}
type TextRule struct {
	Filters Filters
	rt.TextEval
}
type ObjectRule struct {
	Filters Filters
	rt.ObjectEval
}
type NumListRule struct {
	Filters Filters
	rt.NumListEval
}
type TextListRule struct {
	Filters Filters
	rt.TextListEval
}
type ObjListRule struct {
	Filters Filters
	rt.ObjListEval
}
type ExecuteRule struct {
	Filters Filters
	rt.Execute
	Flags
}

type BoolRules []BoolRule
type NumberRules []NumberRule
type TextRules []TextRule
type ObjectRules []ObjectRule
type NumListRules []NumListRule
type TextListRules []TextListRule
type ObjListRules []ObjListRule
type ExecuteRules []ExecuteRule

func (ps BoolRules) GetBool(run rt.Runtime) (ret bool, err error) {
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
func (ps NumberRules) GetNumber(run rt.Runtime) (ret float64, err error) {
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
func (ps TextRules) GetText(run rt.Runtime) (ret string, err error) {
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
func (ps ObjectRules) GetObject(run rt.Runtime) (ret rt.Object, err error) {
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

func (ps NumListRules) GetNumberStream(run rt.Runtime) (ret rt.NumberStream, err error) {
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
func (ps TextListRules) GetTextStream(run rt.Runtime) (ret rt.TextStream, err error) {
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
func (ps ObjListRules) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
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

func (ps ExecuteRules) Execute(run rt.Runtime) (ret bool, err error) {
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
