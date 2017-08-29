package rule

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
)

type Rule struct {
	Name string      // Name of pattern
	If   rt.BoolEval // Filter
}

type BoolRule struct {
	Rule
	Decide rt.BoolEval // Boolean to return if filter passes
}
type NumberRule struct {
	Rule
	Decide rt.NumberEval // Number to return if filter passes
}
type TextRule struct {
	Rule
	Decide rt.TextEval // String to return if filter passes
}
type ObjectRule struct {
	Rule
	Decide rt.ObjectEval // Obj to return if filter passes
}
type NumListRule struct {
	Rule
	Decide rt.NumListEval // List to return if filter passes
}
type TextListRule struct {
	Rule
	Decide rt.TextListEval // List to return if filter passes
}
type ObjListRule struct {
	Rule
	Decide rt.ObjListEval // List to return if filter passes
}
type RunRule struct {
	Rule
	Decide   rt.ExecuteList // List to return if filter passes
	Continue PatternTiming
}

func (p *Rule) Init(pt unique.Types) (ret string, filters []rt.BoolEval, err error) {
	pid := id.MakeId(p.Name)
	if _, ok := pt[pid]; !ok {
		err = errutil.New("unknown pattern", p.Name)
	} else {
		ret, filters = pid, expandFilters(p.If)
	}
	return
}

func (p *BoolRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		ps.Bools.AddRule(pid, filters, p.Decide)
	}
	return
}
func (p *NumberRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		ps.Numbers.AddRule(pid, filters, p.Decide)
	}
	return
}
func (p *TextRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		ps.Text.AddRule(pid, filters, p.Decide)
	}
	return
}
func (p *ObjectRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		ps.Objects.AddRule(pid, filters, p.Decide)
	}
	return
}
func (p *NumListRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		ps.NumLists.AddRule(pid, filters, p.Decide)
	}
	return
}
func (p *TextListRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		ps.TextLists.AddRule(pid, filters, p.Decide)
	}
	return
}
func (p *ObjListRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		ps.ObjLists.AddRule(pid, filters, p.Decide)
	}
	return
}
func (p *RunRule) Mandate(pt unique.Types, ps Rules) (err error) {
	if pid, filters, e := p.Init(pt); e != nil {
		err = e
	} else {
		flags := pat.Infix
		if p.Continue != nil {
			flags = p.Continue.Flags()
		}
		ps.Executes.AddRule(pid, filters, p.Decide, flags)
	}
	return
}

// expandFilters turns a single bool eval into an array by looking at its type.
// FIX? includes core for AllTrue introspection
// but could we do this at the spec level? its kind of an odd dependency here.
func expandFilters(eval rt.BoolEval) (ret []rt.BoolEval) {
	if eval != nil {
		if multi, ok := eval.(*core.AllTrue); ok {
			ret = multi.Test
		} else {
			ret = append(ret, eval)
		}
	}
	return
}
