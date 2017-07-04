package unique

import (
	"github.com/ionous/errutil"
	r "reflect"
)

type TypeRegistry interface {
	// expects a pointer to a struct
	RegisterType(r.Type) error
	FindType(name string) (r.Type, bool)
}

type ValueRegistry interface {
	RegisterValue(r.Value) error
	FindValue(name string) (r.Value, bool)
}

// RegisterBlock registers a structure containing pointers to commands.
func RegisterBlock(reg TypeRegistry, block interface{}) (err error) {
	if blockType := r.TypeOf(block); blockType.Kind() != r.Ptr {
		err = errutil.New("expected (nil) pointer (to a struct).")
	} else if structType := blockType.Elem(); structType.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer.")
	} else {
		for i, cnt := 0, structType.NumField(); i < cnt; i++ {
			field := structType.Field(i)
			if e := registerType(reg, field.Type); e != nil {
				err = errutil.New(field.Name, e)
				break
			}
		}
	}
	return
}

// RegisterType registers pointers to types.
func RegisterTypes(reg TypeRegistry, ptr ...interface{}) (err error) {
	for i, t := range ptr {
		if e := registerType(reg, r.TypeOf(t)); e != nil {
			err = errutil.New("RegisterType", i, e)
		}
	}
	return
}

func registerType(reg TypeRegistry, rtype r.Type) (err error) {
	if rtype.Kind() != r.Ptr {
		err = errutil.New("expected (nil) pointer (to a struct)", rtype)
	} else if rtype := rtype.Elem(); rtype.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer", rtype)
	} else {
		err = reg.RegisterType(rtype)
	}
	return
}

// RegisterValues registers multiple pointer values.
func RegisterValues(reg ValueRegistry, value ...interface{}) (err error) {
	for i, v := range value {
		if e := registerValue(reg, r.ValueOf(v)); e != nil {
			err = errutil.New("RegisterValue", i, e)
			break
		}
	}
	return
}

func registerValue(reg ValueRegistry, rval r.Value) (err error) {
	if rval.Kind() != r.Ptr {
		err = errutil.New("expected pointer (to a struct)", rval.Type())
	} else if rval := rval.Elem(); rval.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer", rval.Type())
	} else {
		err = reg.RegisterValue(rval)
	}
	return
}
