package reflector

import (
	"github.com/ionous/iffy/ref"
	r "reflect"
)

type RefInst struct {
	id   string
	rval r.Value
	cls  *RefClass
}

// GetId returns the unique identifier for this Object.
func (n *RefInst) GetId() string {
	return n.id
}

// GetClass returns the variety of object.
func (n *RefInst) GetClass() (ret ref.Class) {
	return n.cls
}

// GetValue stores the value into the pointer pv.
// Values include ref.Objects for relations and pointers, numbers, and text. For numbers, pv can be any numberic type: float64, int, etc.
func (n *RefInst) GetValue(name string, pv interface{}) (err error) {
	//	id := MakeId(name)
	panic("not implemented")
	return
}

// GetValue can return error when the value violates a property constraint,
// if the value is not of the requested type, or if the targeted property holder is read-only. Read-only values include the "many" side of a relation.
func (n *RefInst) SetValue(name string, v interface{}) (err error) {
	panic("not implemented")
	return
}
