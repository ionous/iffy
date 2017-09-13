// Package coerce copies values and slices of values.
package coerce

import (
	"github.com/ionous/errutil"
	r "reflect"
)

// Value moves src into dst.
// For example:
// . ([])string <=> ([])string
// . ([])NumericType <=> ([])NumericType
// . bool <=> bool, eval<=>eval
// . string <=> enumerated value
// Note: if src and dst are slices of the same type then dst will acquire the same slice as src; that is: changing elements of src will change elements of dst.
//
// Coercion of ([])ident.Id <=> ([])rt.Object happens elsewhere.
// Coercion of ([])eval <=> ([])primitive happens elsewhere.
//
// In the future, iffy will likely support enumerated value of bool. Currently, iffy only supports int.
func Value(dst, src r.Value) (err error) {
	if !dst.CanSet() {
		err = errutil.New("destination not settable")
	} else if dt := dst.Type(); dt == src.Type() {
		dst.Set(src)
	} else if dst.Kind() == r.Slice && src.Kind() == r.Slice {
		err = Slice(dst, src, coerceValue)
	} else {
		err = coerceValue(dst, src)
	}
	return
}

// Slice moves a src slice into a dst slice.
// Panics if src or dst are not slices, errors if element coercion failed.
func Slice(dst, src r.Value, elFn ElementFn) (err error) {
	cnt := src.Len()
	if dt := dst.Type(); dt == src.Type() {
		dst.Set(src)
	} else {
		x := r.MakeSlice(dt, cnt, cnt)
		for i := 0; i < cnt; i++ {
			dst, src := x.Index(i), src.Index(i)
			if e := elFn(dst, src); e != nil {
				err = e
				break
			}
		}
		if err == nil {
			dst.Set(x)
		}
	}
	return
}

type ElementFn func(dst, src r.Value) error

func coerceValue(dst, src r.Value) (err error) {
	if st, dt := src.Type(), dst.Type(); !st.ConvertibleTo(dt) {
		err = errutil.New("cant convert to", dt, "from", st)
	} else {
		v := src.Convert(dt)
		dst.Set(v)
	}
	return
}
