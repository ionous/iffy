package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/ref/enum"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type RefObject struct {
	rval    r.Value // stores the concrete value. ex. Rock, not *Rock.
	objects *Objects
}

// Value exposed for testing
func (n RefObject) Value() r.Value {
	return n.rval
}

// GetId returns the unique identifier for this Object.
// Blank for anonymous and temporary objects.
func (n RefObject) GetId() (ret string) {
	if path, ok := unique.PathOf(n.rval.Type(), "id"); ok {
		field := n.rval.FieldByIndex(path)
		name := field.String()
		ret = id.MakeId(name)
	}
	return
}

func (n RefObject) String() string {
	return n.GetId()
}

// GetClass returns the variety of object.
func (n RefObject) GetClass() rt.Class {
	return MakeClass(n.rval.Type())
}

// GetValue stores the value into the pointer pv.
// Values include rt.Objects for relations and pointers, numbers, and text. For numbers, pv can be any numberic type: float64, int, etc.
func (n RefObject) GetValue(name string, pv interface{}) (err error) {
	if e := n.getValue(name, pv); e != nil {
		err = errutil.New(n, e, name)
	}
	return
}

func (n RefObject) getValue(name string, pv interface{}) (err error) {
	if pdst := r.ValueOf(pv); pdst.Kind() != r.Ptr {
		err = errutil.New("get expected a pointer outvalue")
	} else {
		pid, dst := id.MakeId(name), pdst.Elem()
		rtype := n.rval.Type()
		// a bool is an indicator of state lookup
		if dst.Kind() == r.Bool {
			if path, idx := enum.PropertyPath(rtype, pid); len(path) == 0 {
				err = errutil.New("get choice not found")
			} else {
				field := n.rval.FieldByIndex(path)
				if field.Kind() == r.Bool {
					dst.SetBool(field.Bool())
				} else {
					c := field.Int()
					match := c == int64(idx)
					dst.SetBool(match)
				}
			}
		} else {
			if path := class.PropertyPath(rtype, pid); len(path) == 0 {
				err = errutil.New("get property not found")
			} else {
				field := n.rval.FieldByIndex(path)
				err = n.objects.coerce(dst, field)
			}
		}
	}
	return
}

// GetValue can return error when the value violates a property constraint,
// if the value is not of the requested type, or if the targeted property holder is read-only. Read-only values include the "many" side of a relation.
func (n RefObject) SetValue(name string, v interface{}) (err error) {
	pid := id.MakeId(name)
	if e := n.setValue(pid, v); e != nil {
		err = errutil.New(e, name)
	}
	return
}

func (n RefObject) setValue(pid string, v interface{}) (err error) {
	if val, ok := v.(bool); ok {
		err = n.setBool(pid, val)
	} else {
		rtype := n.rval.Type()
		if path := class.PropertyPath(rtype, pid); len(path) == 0 {
			err = errutil.New("set property not found", n)
		} else {
			field := n.rval.FieldByIndex(path)
			src := r.ValueOf(v)
			enumish := field.Kind() == r.Int && src.Kind() == r.String
			if !enumish {
				err = n.objects.coerce(field, src)
			} else {
				if choices := enum.Enumerate(field.Type()); len(choices) == 0 {
					err = errutil.New("not an enumerated field", pid)
				} else {
					choice := id.MakeId(src.String())
					if i, ok := enum.ChoiceToIndex(choice, choices); !ok {
						err = errutil.New("set unknown choice", choice, choices)
					} else {
						err = coerceValue(field, r.ValueOf(i))
					}
				}
			}
		}
	}
	return
}

func (n RefObject) setBool(pid string, val bool) (err error) {
	rtype := n.rval.Type()
	if path, idx := enum.PropertyPath(rtype, pid); len(path) == 0 {
		err = errutil.New("set choice not found", pid)
	} else {
		field := n.rval.FieldByIndex(path)
		if field.Kind() == r.Bool {
			err = CoerceValue(field, val)
		} else if val {
			// if setting the choice to true, then we are setting the field to the choice.
			err = CoerceValue(field, idx)
		} else {
			// we have to try to generate an opposite value.
			if choices := enum.Enumerate(field.Type()); len(choices) == 0 {
				err = errutil.New("not an enumerated field", pid)
			} else if cnt := len(choices); cnt > 2 {
				err = errutil.New("no opposite value. too many choices", pid, cnt)
			} else {
				// idx= 0; 2-(0+1)=1
				// idx= 1; 2-(1+1)=0
				// ret can be out of range for 1 length enums
				idx := 2 - (idx + 1)
				err = CoerceValue(field, idx)
			}
		}
	}
	return
}
