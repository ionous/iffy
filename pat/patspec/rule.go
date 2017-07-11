package patspec

import (
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

func (p *BoolRule) Generate(ps PatternFactory) (err error) {
	return ps.Bool(p.Name, p.If, p.Decide)
}
func (p *NumberRule) Generate(ps PatternFactory) (err error) {
	return ps.Number(p.Name, p.If, p.Decide)
}
func (p *TextRule) Generate(ps PatternFactory) (err error) {
	return ps.Text(p.Name, p.If, p.Decide)
}
func (p *ObjectRule) Generate(ps PatternFactory) (err error) {
	return ps.Object(p.Name, p.If, p.Decide)
}
func (p *NumListRule) Generate(ps PatternFactory) (err error) {
	return ps.NumList(p.Name, p.If, p.Decide)
}
func (p *TextListRule) Generate(ps PatternFactory) (err error) {
	return ps.TextList(p.Name, p.If, p.Decide)
}
func (p *ObjListRule) Generate(ps PatternFactory) (err error) {
	return ps.ObjList(p.Name, p.If, p.Decide)
}
