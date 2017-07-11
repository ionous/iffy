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

type ChangeState struct {
	Ref   rt.ObjectEval
	State string
}

func (p *SetBool) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the orginal object.
func (p *SetBool) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetBool) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("SetBool.Obj", e)
	} else if val, e := p.Val.GetBool(run); e != nil {
		err = errutil.New("SetBool.Val", e)
	} else if e := obj.SetValue(p.Prop, val); e != nil {
		err = errutil.New("SetBool.Set", e)
	} else {
		ret = obj
	}
	return
}

func (p *SetNum) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the orginal object.
func (p *SetNum) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetNum) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("SetNum.Obj", e)
	} else if val, e := p.Val.GetNumber(run); e != nil {
		err = errutil.New("SetNum.Val", e)
	} else if e := obj.SetValue(p.Prop, val); e != nil {
		err = errutil.New("SetNum.Set", e)
	} else {
		ret = obj
	}
	return
}

func (p *SetText) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the orginal object.
func (p *SetText) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetText) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("SetText.Obj", e)
	} else if val, e := p.Val.GetText(run); e != nil {
		err = errutil.New("SetText.Val", e)
	} else if e := obj.SetValue(p.Prop, val); e != nil {
		err = errutil.New("SetText.Set", e)
	} else {
		ret = obj
	}
	return
}

func (p *SetObj) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

// GetObject executes the command, and returns a reference to the orginal object.
func (p *SetObj) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *SetObj) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Obj.GetObject(run); e != nil {
		err = errutil.New("SetObj.Obj", e)
	} else if val, e := p.Val.GetObject(run); e != nil {
		err = errutil.New("SetObj.Val", e)
	} else if e := obj.SetValue(p.Prop, val); e != nil {
		err = errutil.New("SetObj.Set", e)
	} else {
		ret = obj
	}
	return
}

func (p *ChangeState) Execute(run rt.Runtime) error {
	_, err := p.exec(run)
	return err
}

func (p *ChangeState) GetObject(run rt.Runtime) (rt.Object, error) {
	return p.exec(run)
}

func (p *ChangeState) exec(run rt.Runtime) (ret rt.Object, err error) {
	if obj, e := p.Ref.GetObject(run); e != nil {
		err = errutil.New("ChangeState.Ref", e)
	} else if prop, ok := obj.GetClass().GetPropertyByChoice(p.State); !ok {
		err = errutil.New("ChangeState", obj, "does not have choice", p.State)
	} else if e := obj.SetValue(prop.GetId(), p.State); e != nil {
		err = errutil.New("ChangeState", e)
	} else {
		ret = obj
	}
	return
}
