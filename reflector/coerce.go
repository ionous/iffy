package reflector

import (
	"github.com/ionous/errutil"
	r "reflect"
)

// Coerce moves the srcValue into the outValue pointer
func Coerce(outValue, srcValue interface{}) (err error) {
	if dst := r.ValueOf(outValue); dst.Kind() != r.Ptr {
		err = errutil.Fmt("destination not a pointer, %s", dst.Kind())
	} else {
		err = CoerceToValue(dst.Elem(), srcValue)
	}
	return
}

func CoerceToValue(dst r.Value, srcValue interface{}) (err error) {
	src := r.ValueOf(srcValue)
	if !src.Type().ConvertibleTo(dst.Type()) {
		err = errutil.New("incompatible types", src.Type(), dst.Type())
	} else {
		v := src.Convert(dst.Type())
		dst.Set(v)
	}
	return
}
