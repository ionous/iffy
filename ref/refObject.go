package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type RefObject struct {
	id      string    // unique id, blank for anonymous and temporary objects.
	rval    r.Value   // stores the concrete type. ex. Rock, not *Rock.
	cls     *RefClass // extracted up front to make sure it exists.
	objects *Objects
}

// Value is primarily for testing
func (n *RefObject) Value() r.Value {
	return n.rval
}

// GetId returns the unique identifier for this Object.
func (n *RefObject) GetId() string {
	return n.id
}

// GetClass returns the variety of object.
func (n *RefObject) GetClass() (ret rt.Class) {
	return n.cls
}

// GetValue stores the value into the pointer pv.
// Values include rt.Objects for relations and pointers, numbers, and text. For numbers, pv can be any numberic type: float64, int, etc.
func (n *RefObject) GetValue(name string, pv interface{}) (err error) {
	if pdst := r.ValueOf(pv); pdst.Kind() != r.Ptr {
		err = errutil.New("expected a pointer outvalue")
	} else {
		id, dst := id.MakeId(name), pdst.Elem()
		// a bool is an indicator of state lookup
		if dst.Kind() == r.Bool {
			if enum, path, idx := n.cls.getPropertyByChoice(id); enum == nil {
				err = errutil.New("get value, choice not found", name)
			} else if field := n.rval.FieldByIndex(path); !field.IsValid() {
				err = errutil.New("field not found", name)
			} else if field.Kind() == r.Bool {
				dst.SetBool(field.Bool())
			} else {
				match := field.Int() == int64(idx)
				dst.SetBool(match)
			}
		} else if _, path, ok := n.cls.getProperty(id); !ok {
			err = errutil.New("property not found", name)
		} else if field := n.rval.FieldByIndex(path); !field.IsValid() {
			err = errutil.New("field not found", name)
		} else {
			err = n.objects.coerce(dst, field)
		}
	}
	return
}

// GetValue can return error when the value violates a property constraint,
// if the value is not of the requested type, or if the targeted property holder is read-only. Read-only values include the "many" side of a relation.
func (n *RefObject) SetValue(name string, v interface{}) (err error) {
	id := id.MakeId(name)
	if e := n.setValue(id, v); e != nil {
		err = errutil.New(e, name)
	}
	return
}

func (n *RefObject) setValue(id string, v interface{}) (err error) {
	if val, ok := v.(bool); ok {
		if p, path, idx := n.cls.getPropertyByChoice(id); !ok {
			err = errutil.New("set value, choice not found", id)
		} else {
			if field := n.rval.FieldByIndex(path); !field.IsValid() {
				err = errutil.New("field not found", id)
			} else {
				err = p.setValue(field, idx, val)
			}
		}
	} else if _, path, ok := n.cls.getProperty(id); !ok {
		err = errutil.New("property not found", id)
	} else if field := n.rval.FieldByIndex(path); !field.IsValid() {
		err = errutil.New("field not found", id)
	} else {
		src := r.ValueOf(v)
		err = n.objects.coerce(field, src)
	}
	return
}
