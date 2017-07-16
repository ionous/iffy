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
				err = errutil.New("choice not found", name)
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
			rtype := field.Type()
			if k := rtype.Kind(); k == r.Ptr {
				// the field is a pointer, we want to give back object ref.
				if field.IsNil() {
					err = errutil.New("field in nil", name)
				} else if obj, e := n.objects.GetByValue(field); e != nil {
					err = e
				} else if dst.Kind() == r.Interface {
					// assume user is asking for rt.Object
					v := r.ValueOf(obj)
					if vt, dt := v.Type(), dst.Type(); !vt.AssignableTo(dt) {
						err = errutil.New("cant assign", vt, "to", dt)
					} else {
						dst.Set(v)
					}
				} else if dst.Kind() == r.Ptr {
					// assume user is asking for *Something
					v := obj.rval.Addr()
					if vt, dt := v.Type(), dst.Type(); !vt.AssignableTo(dt) {
						err = errutil.New("cant assign", vt, "to", dt)
					} else {
						dst.Set(v)
					}
				} else {
					err = errutil.New("unknown out value for pointer field", name, dst.Type().Name())
				}
			} else if k == r.Slice && rtype.Elem().Kind() == r.Ptr {
				if dst.Kind() != r.Slice {
					err = errutil.New("expected (pointer to a) slice outvalue")
				} else {
					slice, dt := dst, dst.Type().Elem()
					if dt.Kind() == r.Interface {
						err = errutil.New("output to []rt.Object not supported yet")
					} else if dt.Kind() == r.Ptr {
						// assume user is asking for []*Something
						for i := 0; i < field.Len(); i++ {
							field := field.Index(i)
							if obj, e := n.objects.GetByValue(field); e != nil {
								err = e
								break
							} else {
								v := obj.rval.Addr()
								if vt := v.Type(); !vt.AssignableTo(dt) {
									err = errutil.New("cant assign output element", vt, "to", dt)
									break
								} else {
									slice = r.Append(slice, v)
								}
							}
						}
					} else {
						err = errutil.New("unknown out value for pointer slice", name, dst.Type().Name())
					}
					if err == nil {
						dst.Set(slice)
					}
				}
			} else {
				err = CoerceValue(dst, field)
			}
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
			err = errutil.New("choice not found", id)
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
	} else if ref, ok := v.(*RefObject); ok {
		// FIX: handle child pointers?
		// maybe walk parent class hierarchy searching for a pointer match?
		rvalue, elType := ref.rval.Addr(), field.Type()
		if from := rvalue.Type(); !from.AssignableTo(elType) {
			err = errutil.Fmt("incompatible pointer type. from: %v to: %v", from, elType)
		} else {
			field.Set(rvalue)
		}
	} else if objects, ok := v.([]rt.Object); ok {
		// this is very much like setField ( tho we may need to walk child pointers )
		// ( what if we supported aggregate ops? )
		// "field" is an rvalue -- in this case an array of some pointers..
		if field.Kind() != r.Slice {
			err = errutil.New("field is not a slice", id)
		} else {
			slice, elType := field, field.Type().Elem()
			for _, obj := range objects {
				if ref, ok := obj.(*RefObject); !ok {
					err = errutil.Fmt("object in list not a ref object %T", obj)
					break
				} else {
					rvalue := ref.rval.Addr()
					if from := rvalue.Type(); !from.AssignableTo(elType) {
						err = errutil.Fmt("incompatible element type. from: %v to: %v", from, elType)
						break
					} else {
						slice = r.Append(slice, rvalue)
					}
				}
			}
			if err == nil {
				field.Set(slice)
			}
		}
	} else {
		err = CoerceValue(field, v)
	}
	return
}
