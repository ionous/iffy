package rules

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
)

// Rule is the base class for implementing patterns.
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
type ListNumbers struct {
	Rule
	Decide   rt.NumListEval // List to return if filter passes
	Continue PatternTiming
}
type ListText struct {
	Rule
	Decide   rt.TextListEval // List to return if filter passes
	Continue PatternTiming
}
type ListObjects struct {
	Rule
	Decide   rt.ObjListEval // List to return if filter passes
	Continue PatternTiming
}
type RunRule struct {
	Rule
	Decide   rt.ExecuteList // List to return if filter passes
	Continue PatternTiming
}

func (p *Rule) Init(pt unique.Types) (ret ident.Id, filters []rt.BoolEval, err error) {
	pid := ident.IdOf(p.Name)
	if _, ok := pt[pid]; !ok {
		err = errutil.New("unknown pattern", p.Name)
	} else {
		ret, filters = pid, expandFilters(p.If)
	}
	return
}

func (p *BoolRule) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddBoolRule(pid, filters, p.Decide)
	}
	return
}
func (p *NumberRule) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddNumberRule(pid, filters, p.Decide)
	}
	return
}
func (p *TextRule) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddTextRule(pid, filters, p.Decide)
	}
	return
}
func (p *ObjectRule) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddObjectRule(pid, filters, p.Decide)
	}
	return
}
func (p *ListNumbers) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddNumListRule(pid, filters, p.Decide, flag(p.Continue))
	}
	return
}
func (p *ListText) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddTextListRule(pid, filters, p.Decide, flag(p.Continue))
	}
	return
}
func (p *ListObjects) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddObjListRule(pid, filters, p.Decide, flag(p.Continue))
	}
	return
}
func (p *RunRule) Mandate(ps pat.Contract) (err error) {
	if pid, filters, e := p.Init(ps.Types); e != nil {
		err = e
	} else {
		err = ps.AddExecuteRule(pid, filters, p.Decide, flag(p.Continue))
	}
	return
}
func flag(time PatternTiming) pat.Flags {
	flags := pat.Infix
	if time != nil {
		flags = time.Flags()
	}
	return flags
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
