package next

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type SetField struct {
	Obj, Field rt.TextEval
}

type SetFieldBool struct {
	SetField
	Val rt.BoolEval
}

type SetFieldNum struct {
	SetField
	Val rt.NumberEval
}

type SetFieldText struct {
	SetField
	Val rt.TextEval
}

// type SetFieldState struct {
// 	Obj, State rt.TextEval
// }

func (op *SetField) SetValue(run rt.Runtime, v interface{}) (err error) {
	if obj, e := rt.GetText(run, op.Obj); e != nil {
		err = e
	} else if field, e := rt.GetText(run, op.Field); e != nil {
		err = e
	} else {
		err = run.SetField(obj, field, v)
	}
	return
}

func (*SetFieldBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_field_bool",
		Group: "objects",
		Desc:  "Set Boolean Field: Sets the named field to the passed boolean value.",
	}
}

func (op *SetFieldBool) Execute(run rt.Runtime) (err error) {
	if val, e := rt.GetBool(run, op.Val); e != nil {
		err = e
	} else {
		err = op.SetValue(run, val)
	}
	return
}

func (*SetFieldNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_field_num",
		Group: "objects",
		Desc:  "Set Number Field: Sets the named field to the passed number.",
	}
}

func (op *SetFieldNum) Execute(run rt.Runtime) (err error) {
	if val, e := rt.GetNumber(run, op.Val); e != nil {
		err = e
	} else {
		err = op.SetValue(run, val)
	}
	return
}

func (*SetFieldText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "set_field_text",
		Group: "objects",
		Desc:  "Set Text Field: Sets the named field to the passed text value."}
}

func (op *SetFieldText) Execute(run rt.Runtime) (err error) {
	if val, e := rt.GetText(run, op.Val); e != nil {
		err = e
	} else {
		err = op.SetValue(run, val)
	}
	return
}

// corresponding Get?
// func (op *SetState) Execute(run rt.Runtime) (err error) {
// 	if obj, e := op.Ref.GetObject(run); e != nil {
// 		err = errutil.New("cant SetFieldState, because get owner", e)
// 	} else if e := run.SetValue(obj, op.State, true); e != nil {
// 		err = errutil.New("cant SetFieldState, because property", e)
// 	}
// 	return
// }
