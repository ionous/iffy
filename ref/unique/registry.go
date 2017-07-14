package unique

import (
	"github.com/ionous/errutil"
	r "reflect"
)

type TypeRegistry interface {
	// RegisterType expects a pointer to a (nil) struct.
	RegisterType(r.Type) error
}

type ValueRegistry interface {
	// RegisterValue expects a pointer to a (non-nil) struct.
	RegisterValue(r.Value) error
}

// TypePtr turns an interface pointer into a struct r.Type.
func TypePtr(ptr interface{}) (r.Type, error) {
	return typePtr(r.TypeOf(ptr))
}

func typePtr(rtype r.Type) (ret r.Type, err error) {
	if rtype.Kind() != r.Ptr {
		err = errutil.New("expected (nil) pointer (to a struct)", rtype)
	} else if rtype := rtype.Elem(); rtype.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer", rtype)
	} else {
		ret = rtype
	}
	return
}

// ValuePtr turns an interface pointer into a struct r.Value.
func ValuePtr(ptr interface{}) (ret r.Value, err error) {
	if rval := r.ValueOf(ptr); rval.Kind() != r.Ptr {
		err = errutil.New("expected pointer (to a struct)", rval.Type())
	} else if rval := rval.Elem(); rval.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer", rval.Type())
	} else {
		ret = rval
	}
	return
}

// RegisterBlock registers a structure containing pointers to commands.
func RegisterBlocks(reg TypeRegistry, block ...interface{}) (err error) {
OutOfLoop:
	for _, block := range block {
		if structType, e := TypePtr(block); e != nil {
			err = e
			break OutOfLoop
		} else {
			for i, cnt := 0, structType.NumField(); i < cnt; i++ {
				field := structType.Field(i)
				if rtype, e := typePtr(field.Type); e != nil {
					err = errutil.New("RegisterType", i, e)
					break OutOfLoop
				} else if e := reg.RegisterType(rtype); e != nil {
					err = errutil.New("RegisterType", i, e)
					break OutOfLoop
				}
			}
		}
	}
	return
}

// RegisterType registers pointers to types.
func RegisterTypes(reg TypeRegistry, ptr ...interface{}) (err error) {
	for i, t := range ptr {
		if rtype, e := TypePtr(t); e != nil {
			err = errutil.New("RegisterType", i, e)
			break
		} else if e := reg.RegisterType(rtype); e != nil {
			err = errutil.New("RegisterType", i, e)
			break
		}
	}
	return
}

// RegisterValues registers multiple pointer values.
func RegisterValues(reg ValueRegistry, value ...interface{}) (err error) {
	for i, v := range value {
		if rval, e := ValuePtr(v); e != nil {
			err = errutil.New("RegisterValue", i, e)
			break
		} else if e := reg.RegisterValue(rval); e != nil {
			err = errutil.New("RegisterValue", i, e)
			break
		}
	}
	return
}
