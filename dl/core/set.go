package core

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type SetBool struct {
	Obj  rt.ObjectEval
	Prop string
	Val  rt.BoolEval
}

type SetNum struct {
	Obj  rt.ObjectEval
	Prop string
	Val  rt.NumberEval
}

type SetText struct {
	Obj  rt.ObjectEval
	Prop string
	Val  rt.TextEval
}

type SetObj struct {
	Obj  rt.ObjectEval
	Prop string
	Val  rt.ObjectEval
}

type SetState struct {
	Ref   rt.ObjectEval
	State string
}

func (p *SetBool) Execute(run rt.Runtime) (err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("cant SetBool, because get owner", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("cant SetBool, because property", e)
	}
	return
}

func (p *SetNum) Execute(run rt.Runtime) (err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("cant SetNum, because get owner", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("cant SetNum, because property", e)
	}
	return
}

func (p *SetText) Execute(run rt.Runtime) (err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("cant SetText, because get owner", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("cant SetText, because property", e)
	}
	return
}

func (p *SetObj) Execute(run rt.Runtime) (err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("cant SetObj, because get owner", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("cant SetObj, because property", e)
	}
	return
}

func (p *SetState) Execute(run rt.Runtime) (err error) {
	if obj, e := p.Ref.GetObject(run); e != nil {
		err = errutil.New("cant SetState, because get owner", e)
	} else if e := obj.SetValue(p.State, true); e != nil {
		err = errutil.New("cant SetState, because property", e)
	}
	return
}
