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

// StructPtr turns an interface pointer into a struct r.Type.
func StructPtr(ptr interface{}) (r.Type, error) {
	return TypePtr(r.Struct, ptr)
}

func TypePtr(kind r.Kind, ptr interface{}) (r.Type, error) {
	return typePtr(kind, r.TypeOf(ptr))
}

func typePtr(kind r.Kind, rtype r.Type) (ret r.Type, err error) {
	if rtype.Kind() != r.Ptr {
		err = errutil.Fmt("expected (nil) pointer (to a %s) %s", kind, rtype)
	} else if rtype := rtype.Elem(); rtype.Kind() != kind {
		err = errutil.Fmt("expected a %s pointer %s", kind, rtype)
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
	for _, block := range block {
		if structType, e := StructPtr(block); e != nil {
			err = e
			break
		} else if e := registerBlock(reg, structType); e != nil {
			err = e
			break
		}
	}
	return
}

func registerBlock(reg TypeRegistry, structType r.Type) (err error) {
	for i, cnt := 0, structType.NumField(); i < cnt; i++ {
		field := structType.Field(i)
		if field.Type.Kind() == r.Struct {
			if e := registerBlock(reg, field.Type); e != nil {
				err = e
				break
			}
		} else if rtype, e := typePtr(r.Struct, field.Type); e != nil {
			err = errutil.New("RegisterType", i, e)
			break
		} else if e := reg.RegisterType(rtype); e != nil {
			err = errutil.New("RegisterType", i, e)
			break
		}
	}
	return
}

// RegisterType registers pointers to types.
func RegisterTypes(reg TypeRegistry, ptr ...interface{}) (err error) {
	for i, t := range ptr {
		if rtype, e := StructPtr(t); e != nil {
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
