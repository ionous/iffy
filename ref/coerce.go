package ref

import (
	"github.com/ionous/errutil"
	r "reflect"
)

func (or *Objects) coerce(dst, src r.Value) (err error) {
	switch dstType := dst.Type(); dstType.Kind() {
	case r.Interface: // dst is probably rt.Object
		if srcobj, e := or.getByValue(src); e != nil {
			err = e
		} else {
			src := r.ValueOf(srcobj) // get the value ( of *RefObj )
			if srcType := src.Type(); !srcType.AssignableTo(dstType) {
				err = errutil.New("cant assign", srcType, "to", dstType)
			} else {
				dst.Set(src)
			}
		}
	case r.Ptr: // dst is probably *Something
		err = copyPtr(src, func(ptr r.Value) (okay bool) {
			if dstType == ptr.Type() {
				dst.Set(ptr)
				okay = true
			}
			return
		})
	case r.Slice:
		if k := dstType.Elem().Kind(); k == r.Ptr || k == r.Interface {
			// dst is probably []*Something, or []rt.Object
			err = copyObjects(dst, src)
		} else {
			// dst is probably []primitive
			err = coerceValue(dst, src)
		}

	default: // dst is probably primitive or []primitive
		err = coerceValue(dst, src)
	}
	return
}

// give a value, which might be either an interface or a ptr, return the reflected value of *RefObject
func (or *Objects) getByValue(src r.Value) (ret *RefObject, err error) {
	switch k := src.Kind(); k {
	case r.Ptr, r.Interface:
		if src.IsNil() {
			err = errutil.New("nil pointers return error")
		} else if obj, ok := src.Interface().(*RefObject); ok {
			ret = obj
		} else {
			ret, err = or.GetByValue(src)
		}
	default:
		err = errutil.New("unexpected type", k)
	}
	return
}

// CoerceValue moves the src into the dst pointer
func CoerceValue(dst, src interface{}) error {
	return coerceValue(valueOf(dst), valueOf(src))
}

func coerceValue(dst, src r.Value) (err error) {
	// act as if dst is probably a primitive or []primitive
	if !dst.CanSet() {
		err = errutil.New("destination not settable")
	} else if !src.Type().ConvertibleTo(dst.Type()) {
		err = errutil.New("src", src.Type(), "not compatible with dst", dst.Type())
	} else {
		v := src.Convert(dst.Type())
		dst.Set(v)
	}
	return
}

func valueOf(i interface{}) (ret r.Value) {
	if v, ok := i.(r.Value); !ok {
		ret = r.ValueOf(i)
	} else {
		ret = v
	}
	return
}

//
func copyObjects(dst, src r.Value) (err error) {
	if dst.Kind() != r.Slice {
		err = errutil.New("dst is not a slice")
	} else if src.Kind() != r.Slice {
		err = errutil.New("src is not a slice")
	} else if dst.Type() == src.Type() {
		// example: []rt.Object to []rt.Object
		err = coerceValue(dst, src)
	} else {
		slice, elType := dst, dst.Type().Elem()
		if elType.Kind() != r.Ptr {
			err = errutil.New("unknown dst", dst.Type(), "for slice", src.Type())
		} else {
			// assume user is asking for []*Something
			for i := 0; i < src.Len(); i++ {
				el := src.Index(i)
				if e := copyPtr(el, func(ptr r.Value) (okay bool) {
					if ptr.Type() == elType {
						slice = r.Append(slice, ptr)
						okay = true
					}
					return
				}); e != nil {
					err = e
					break
				}
			}
			if err == nil {
				dst.Set(slice)
			}
		}
	}
	return
}

// convert src into a *Something pointer,
// then try setting to dst, walking up the pointer hierarchy.
// note: we cant getByValue bc src could be anonymous.
// ( that's going to be especially true of names. )
func copyPtr(src r.Value, set func(r.Value) bool) (err error) {
	if obj, ok := src.Interface().(*RefObject); ok {
		src := obj.Addr()
		err = Upcast(src, set)
	} else if src.Kind() == r.Ptr {
		err = Upcast(src, set)
	} else {
		err = errutil.New("unknown src", src.Type())
	}
	return
}

// Upcast src is a *Something, for each version of the pointer, we call set()
func Upcast(src r.Value, set func(r.Value) bool) (err error) {
Upcast:
	for !set(src) {
		el := src.Elem()
		walk := el.Type()
		for i := 0; i < walk.NumField(); i++ {
			field := walk.Field(i)
			if IsParentField(&field) {
				src = el.Field(i).Addr()
				continue Upcast
			}
		}
		err = errutil.New("couldnt assign pointer", walk)
		break
	}
	return
}
