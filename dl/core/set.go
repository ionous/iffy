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

func (p *SetBool) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the original object.
func (p *SetBool) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetBool) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("set bool owner error", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("set bool property error", e)
	} else {
		ret = obj
	}
	return
}

func (p *SetNum) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the original object.
func (p *SetNum) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetNum) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("set num owner error", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("set num property error", e)
	} else {
		ret = obj
	}
	return
}

func (p *SetText) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the original object.
func (p *SetText) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetText) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("set text owner error", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("set text property error", e)
	} else {
		ret = obj
	}
	return
}

func (p *SetObj) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the original object.
func (p *SetObj) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetObj) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("set obj owner error", e)
	} else if e := obj.SetValue(p.Prop, p.Val); e != nil {
		err = errutil.New("set obj property error", e)
	} else {
		ret = obj
	}
	return
}

func (p *SetState) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

func (p *SetState) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetState) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Ref.GetObject(run); e != nil {
		err = errutil.New("set state owner error", e)
	} else if e := obj.SetValue(p.State, true); e != nil {
		err = errutil.New("set state property error", e)
	} else {
		ret = obj
	}
	return
}
