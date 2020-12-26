package core

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// Assignment helps limit variable and parameter assignment to particular contexts.
type Assignment interface {
	// write the results of evaluating this into that.
	GetAssignedValue(rt.Runtime) (g.Value, error)
	Affinity() affine.Affinity
}

func GetAssignedValue(run rt.Runtime, a Assignment) (ret g.Value, err error) {
	if a == nil {
		err = safe.MissingEval("empty assignment")
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

type FromName struct {
	Val rt.TextEval
}

type FromRecord struct {
	Val rt.RecordEval
}

type FromObject struct {
	Val rt.ObjectEval
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
func (op *FromBool) Affinity() affine.Affinity {
	return affine.Bool
}
func (op *FromBool) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetBool(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (*FromNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_num",
		Spec:  "{val:number_eval}",
		Group: "variables",
		Desc:  "Assign Number: Assigns the passed number.",
	}
}
func (op *FromNum) Affinity() affine.Affinity {
	return affine.Number
}
func (op *FromNum) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetNumber(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (*FromText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_text",
		Group: "variables",
		Desc:  "Assign Text: Assigns the passed piece of text.",
	}
}
func (op *FromText) Affinity() affine.Affinity {
	return affine.Text
}
func (op *FromText) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetText(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (*FromName) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_name",
		Group: "variables",
		Desc:  "Assign Name: Assigns the passed piece of name.",
	}
}
func (op *FromName) Affinity() affine.Affinity {
	return affine.Object
}
func (op *FromName) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetText(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else if obj, e := getObjectNamed(run, val.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = obj
	}
	return
}

func (*FromRecord) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_record",
		Group: "variables",
		Desc:  "Assign Record: Assigns the passed record.",
	}
}
func (op *FromRecord) Affinity() affine.Affinity {
	return affine.Record
}
func (op *FromRecord) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.GetRecord(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = obj
	}
	return
}

func (*FromObject) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_object",
		Group: "variables",
		Desc:  "Assign Object: Assigns the passed object",
	}
}
func (op *FromObject) Affinity() affine.Affinity {
	return affine.Object
}
func (op *FromObject) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.GetObject(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = obj
	}
	return
}

func (*FromNumList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_num_list",
		Group: "variables",
		Desc:  "Assign Number List: Assigns the passed number list.",
	}
}
func (op *FromNumList) Affinity() affine.Affinity {
	return affine.NumList
}
func (op *FromNumList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := safe.GetNumList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vals
	}
	return
}

func (*FromTextList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_text_list",
		Group: "variables",
		Desc:  "Assign Text List: Assigns the passed text list.",
	}
}
func (op *FromTextList) Affinity() affine.Affinity {
	return affine.TextList
}
func (op *FromTextList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := safe.GetTextList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vals
	}
	return
}

func (*FromRecordList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "assign_record_list",
		Group: "variables",
		Desc:  "Assign Record List: Assigns the passed record list.",
	}
}
func (op *FromRecordList) Affinity() affine.Affinity {
	return affine.RecordList
}
func (op *FromRecordList) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if objs, e := safe.GetRecordList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = objs
	}
	return
}
