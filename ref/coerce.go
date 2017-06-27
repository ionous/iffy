package ref

import (
	"github.com/ionous/errutil"
	r "reflect"
)

// CoerceValue moves the src into the dst pointer
func CoerceValue(dst, src interface{}) (err error) {
	if dst := valueOf(dst); !dst.CanSet() {
		err = errutil.New("destination not settable")
	} else {
		if src := valueOf(src); !src.Type().ConvertibleTo(dst.Type()) {
			err = errutil.New("incompatible types", dst.Type(), src.Type())
		} else {
			v := src.Convert(dst.Type())
			dst.Set(v)
		}
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
