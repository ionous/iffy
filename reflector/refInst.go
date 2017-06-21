package reflector

import (
	"github.com/ionous/errutil"
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
	id := MakeId(name)
	if ret, ok := pv.(*bool); ok {
		if _, path, idx := n.cls.getPropertyByChoice(id); !ok {
			err = errutil.New("choice not found", name)
		} else if v := n.rval.FieldByIndex(path); !v.IsValid() {
			err = errutil.New("field not found", name)
		} else if v.Kind() == r.Bool {
			*ret = v.Bool()
		} else {
			*ret = v.Int() == int64(idx)
		}

	} else {
		panic("not implemented")
	} //else if pid == pluralId {
	return
}

// GetValue can return error when the value violates a property constraint,
// if the value is not of the requested type, or if the targeted property holder is read-only. Read-only values include the "many" side of a relation.
func (n *RefInst) SetValue(name string, v interface{}) (err error) {
	id := MakeId(name)
	if val, ok := v.(bool); ok {
		if p, path, idx := n.cls.getPropertyByChoice(id); !ok {
			err = errutil.New("choice not found", name)
		} else {
			if field := n.rval.FieldByIndex(path); !field.IsValid() {
				err = errutil.New("field not found", name)
			} else {
				err = p.setValue(field, idx, val)
			}
		}
	} else {
		panic("not implemented")
	} //else if pid == pluralId {
	return
}
