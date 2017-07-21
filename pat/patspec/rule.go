package patspec

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

func (p *BoolRule) Generate(ps PatternFactory) (err error) {
	return ps.AddBool(p.Name, p.If, p.Decide)
}
func (p *NumberRule) Generate(ps PatternFactory) (err error) {
	return ps.AddNumber(p.Name, p.If, p.Decide)
}
func (p *TextRule) Generate(ps PatternFactory) (err error) {
	return ps.AddText(p.Name, p.If, p.Decide)
}
func (p *ObjectRule) Generate(ps PatternFactory) (err error) {
	return ps.AddObject(p.Name, p.If, p.Decide)
}
func (p *NumListRule) Generate(ps PatternFactory) (err error) {
	return ps.AddNumList(p.Name, p.If, p.Decide)
}
func (p *TextListRule) Generate(ps PatternFactory) (err error) {
	return ps.AddTextList(p.Name, p.If, p.Decide)
}
func (p *ObjListRule) Generate(ps PatternFactory) (err error) {
	return ps.AddObjList(p.Name, p.If, p.Decide)
}
func (p *RunRule) Generate(ps PatternFactory) (err error) {
	flags := pat.Infix
	if p.Continue != nil {
		flags = p.Continue.Flags()
	}
	return ps.AddExecList(p.Name, p.If, p.Decide, flags)
}
