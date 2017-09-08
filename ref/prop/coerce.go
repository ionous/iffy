package prop

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/enum"
	r "reflect"
)

// CoerceValue moves src into dst.
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
//
func CoerceValue(dst, src r.Value) (err error) {
	if !dst.CanSet() {
		err = errutil.New("destination not settable")
	} else if dst.Kind() == r.Slice && src.Kind() == r.Slice {
		err = coerceSlice(dst, src)
	} else {
		err = coerceValue(dst, src)
	}
	return
}

func coerceSlice(dst, src r.Value) (err error) {
	cnt := src.Len()
	if dt := dst.Type(); dt == src.Type() {
		dst.Set(src)
	} else {
		x := r.MakeSlice(dt, cnt, cnt)
		for i := 0; i < cnt; i++ {
			dst, src := x.Index(i), src.Index(i)
			if e := coerceValue(dst, src); e != nil {
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

func coerceValue(dst, src r.Value) (err error) {
	st, dt := src.Type(), dst.Type()
	// note, we have to check for enum conversions first, otherwise go will convert the integer as a rune into a string.
	if sk, dk := st.Kind(), dt.Kind(); sk == r.String && dk == r.Int {
		if !PackEnum(dst, src) {
			err = errutil.New("couldnt pack enum", dt, st)
		}
	} else if dk == r.String && sk == r.Int {
		if !UnpackEnum(dst, src) {
			err = errutil.New("couldnt unpack enum", dt, st)
		}
	} else if st.ConvertibleTo(dt) {
		v := src.Convert(dt)
		dst.Set(v)
	} else {
		err = errutil.New("couldnt convert", st, "to", dt)
	}
	return
}

// PackEnum changes a named choice into an indexed choice.
// Requires that dst be the value of an enumerated type.
func PackEnum(dst, src r.Value) (okay bool) {
	if choices := enum.Enumerate(dst.Type()); len(choices) > 0 {
		choice := src.String()
		if i, ok := enum.StringToIndex(choice, choices); ok {
			okay = coerceValue(dst, r.ValueOf(i)) == nil
		}
	}
	return
}

// UnpackEnum changes an indexed choice into a named choice.
// Requires that src be the value of an enumerated type.
func UnpackEnum(dst, src r.Value) (okay bool) {
	if choices := enum.Enumerate(src.Type()); len(choices) > 0 {
		if c, ok := enum.IndexToString(int(src.Int()), choices); ok {
			okay = coerceValue(dst, r.ValueOf(c)) == nil
		}
	}
	return
}
