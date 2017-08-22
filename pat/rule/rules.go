package rule

import (
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/rt"
)

type BoolRule struct {
	Name   string
	If     rt.BoolEval
	Decide rt.BoolEval
}
type NumberRule struct {
	Name   string
	If     rt.BoolEval
	Decide rt.NumberEval
}
type TextRule struct {
	Name   string
	If     rt.BoolEval
	Decide rt.TextEval
}
type ObjectRule struct {
	Name   string
	If     rt.BoolEval
	Decide rt.ObjectEval
}
type NumListRule struct {
	Name   string
	If     rt.BoolEval
	Decide rt.NumListEval
}
type TextListRule struct {
	Name   string
	If     rt.BoolEval
	Decide rt.TextListEval
}
type ObjListRule struct {
	Name   string
	If     rt.BoolEval
	Decide rt.ObjListEval
}
type RunRule struct {
	Name     string
	If       rt.BoolEval
	Decide   rt.ExecuteList
	Continue PatternTiming
}

func (p *BoolRule) Mandate(ps RuleFactory) (err error) {
	return ps.AddBool(p.Name, p.If, p.Decide)
}
func (p *NumberRule) Mandate(ps RuleFactory) (err error) {
	return ps.AddNumber(p.Name, p.If, p.Decide)
}
func (p *TextRule) Mandate(ps RuleFactory) (err error) {
	return ps.AddText(p.Name, p.If, p.Decide)
}
func (p *ObjectRule) Mandate(ps RuleFactory) (err error) {
	return ps.AddObject(p.Name, p.If, p.Decide)
}
func (p *NumListRule) Mandate(ps RuleFactory) (err error) {
	return ps.AddNumList(p.Name, p.If, p.Decide)
}
func (p *TextListRule) Mandate(ps RuleFactory) (err error) {
	return ps.AddTextList(p.Name, p.If, p.Decide)
}
func (p *ObjListRule) Mandate(ps RuleFactory) (err error) {
	return ps.AddObjList(p.Name, p.If, p.Decide)
}
func (p *RunRule) Mandate(ps RuleFactory) (err error) {
	flags := pat.Infix
	if p.Continue != nil {
		flags = p.Continue.Flags()
	}
	return ps.AddExecList(p.Name, p.If, p.Decide, flags)
}
