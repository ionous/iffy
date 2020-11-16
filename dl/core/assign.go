package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// Assignment helps limit variable and parameter assignment to particular contexts.
type Assignment interface {
	// write the results of evaluating this into that.
	GetAssignedValue(rt.Runtime) (g.Value, error)
	// fix? for import so we can determine the eval type.
	// maybe at least change to affinity
	GetEval() interface{}
}

func GetAssignedValue(run rt.Runtime, a Assignment) (ret g.Value, err error) {
	if a == nil {
		err = rt.MissingEval("empty assignment")
	} else {
		ret, err = a.GetAssignedValue(run)
	}
	return
}

// Assign a value to a local variable.
type Assign struct {
	Name string // name of variable or parameter we are assigning to.
	From Assignment
}

type FromBool struct {
	Val rt.BoolEval
}

type FromNum struct {
	Val rt.NumberEval
}

type FromText struct {
	Val rt.TextEval
}

type FromRecord struct {
	Val rt.RecordEval
}

type FromNumList struct {
	Vals rt.NumListEval
}

type FromTextList struct {
	Vals rt.TextListEval
}

type FromRecordList struct {
	Vals rt.RecordListEval
}

func (*Assign) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign",
		Spec:  "let {name:variable_name} be {from:assignment}",
		Group: "variables",
		Desc:  "Assignment: Sets a variable to a value.",
	}
}

func (op *Assign) Execute(run rt.Runtime) (err error) {
	if v, e := GetAssignedValue(run, op.From); e != nil {
		err = cmdError(op, e)
	} else if e := run.SetField(object.Variables, op.Name, v); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (*FromBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_bool",
		Group: "variables",
		Desc:  "Assign Boolean: Assigns the passed boolean value.",
	}
}

func (op *FromBool) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := rt.GetBool(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = g.ValueOf(val)
	}
	return
}

func (op *FromBool) GetEval() interface{} {
	return op.Val
}

func (*FromNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_num",
		Spec:  "{val:number_eval}",
		Group: "variables",
		Desc:  "Assign Number: Assigns the passed number.",
	}
}

func (op *FromNum) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := rt.GetNumber(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = g.ValueOf(val)
	}
	return
}

func (op *FromNum) GetEval() interface{} {
	return op.Val
}

func (*FromText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_text",
		Group: "variables",
		Desc:  "Assign Text: Assigns the passed piece of text.",
	}
}

func (op *FromText) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := rt.GetText(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = g.ValueOf(val)
	}
	return
}

func (op *FromText) GetEval() interface{} {
	return op.Val
}

func (*FromRecord) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_record",
		Group: "variables",
		Desc:  "Assign Record: Assigns the passed record.",
	}
}

func (op *FromRecord) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := rt.GetRecord(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.RecordOf(val)
	}
	return
}

func (op *FromRecord) GetEval() interface{} {
	return op.Val
}

func (*FromNumList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_num_list",
		Group: "variables",
		Desc:  "Assign Number List: Assigns the passed number list.",
	}
}

func (op *FromNumList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := rt.GetNumList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = g.ValueOf(vals)
	}
	return
}

func (op *FromNumList) GetEval() interface{} {
	return op.Vals
}

func (*FromTextList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_text_list",
		Group: "variables",
		Desc:  "Assign Text List: Assigns the passed text list.",
	}
}

func (op *FromTextList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := rt.GetTextList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = g.ValueOf(vals)
	}
	return
}

func (op *FromTextList) GetEval() interface{} {
	return op.Vals
}

func (*FromRecordList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_record_list",
		Group: "variables",
		Desc:  "Assign Record List: Assigns the passed record list.",
	}
}

func (op *FromRecordList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if typeName, vals, e := rt.GetRecordList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = g.ValueFrom(vals, affine.RecordList, typeName)
	}
	return
}

func (op *FromRecordList) GetEval() interface{} {
	return op.Vals
}
