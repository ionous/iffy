package generic

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

// ValueOf returns a new generic value wrapper based on analyzing the passed value.
func ValueOf(a affine.Affinity, i interface{}) (rt.Value, error) {
	return ValueFor(i, a, defaultType)
}

// ValueFor adds an optional subtype to values, for example:
// marking text as specifically intended for aspects, traits, etc.
func ValueFor(i interface{}, a affine.Affinity, subtype string) (ret rt.Value, err error) {
	switch a {
	case affine.Bool:
		ret, err = newBoolValue(i, subtype)

	case affine.Number:
		ret, err = newNumValue(i, subtype)
	case affine.NumList:
		ret, err = newNumList(i, subtype)

	case affine.Text:
		ret, err = newTextValue(i, subtype)
	case affine.TextList:
		ret, err = newTextList(i, subtype)

	case affine.Record:
		if r, ok := i.(*Record); !ok {
			err = errutil.Fmt("unknown %s %T", a, i)
		} else if t := r.Type(); len(subtype) > 0 && t != subtype {
			err = errutil.Fmt("mismatched record types", a, t, subtype)
		} else {
			ret = r // record implements value
		}
	case affine.RecordList:
		if r, ok := i.(*RecordSlice); !ok {
			err = errutil.Fmt("unknown %s %T", a, i)
		} else if t := r.Type(); len(subtype) > 0 && t != subtype {
			err = errutil.Fmt("mismatched record types", a, t, subtype)
		} else {
			ret = r // record list implements value
		}
	default:
		err = errutil.New("unhandled affinity", a)
	}
	return
}

func BoolOf(v bool) rt.Value {
	return makeValue(affine.Bool, r.ValueOf(v), defaultType)
}

func StringOf(v string) rt.Value {
	return makeValue(affine.Text, r.ValueOf(v), defaultType)
}

func FloatOf(v float64) rt.Value {
	return makeValue(affine.Number, r.ValueOf(v), defaultType)
}

func StringsOf(vs []string) rt.Value {
	return makeValue(affine.TextList, r.ValueOf(vs), defaultType)
}

func FloatsOf(vs []float64) rt.Value {
	return makeValue(affine.NumList, r.ValueOf(vs), defaultType)
}

func makeValue(a affine.Affinity, v r.Value, subtype string) (ret refValue) {
	if len(subtype) == 0 {
		t := v.Type()
		if v.Kind() == r.Slice {
			t = t.Elem()
		}
		subtype = t.String()
	}
	return refValue{a: a, v: v, t: subtype}
}

func newBoolValue(i interface{}, subtype string) (ret rt.Value, err error) {
	a := affine.Bool
	switch v := i.(type) {
	case nil:
		// zero value for unhandled defaults in sqlite
		ret = False
	case bool:
		ret = makeValue(a, r.ValueOf(v), subtype)
	case int64:
		// sqlite, boolean values can be represented as 1/0
		ret = makeValue(a, r.ValueOf(v == 0), subtype)
	case *bool:
		// creates a dynamic value
		ret = makeValue(a, r.ValueOf(v).Elem(), subtype)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newNumValue(i interface{}, subtype string) (ret rt.Value, err error) {
	a := affine.Number
	switch v := i.(type) {
	case nil:
		// zero value for unhandled defaults in sqlite
		ret = Zero
	case int, int64, float64:
		ret = makeValue(a, r.ValueOf(v), subtype)
	case *float64:
		// creates a dynamic value
		ret = makeValue(a, r.ValueOf(v).Elem(), subtype)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newTextValue(i interface{}, subtype string) (ret rt.Value, err error) {
	a := affine.Text
	switch v := i.(type) {
	case nil:
		// zero value for unhandled defaults in sqlite
		ret = Empty
	case string:
		ret = makeValue(a, r.ValueOf(v), subtype)
	case *string:
		// creates a dynamic value
		ret = makeValue(a, r.ValueOf(v).Elem(), subtype)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newNumList(i interface{}, subtype string) (ret rt.Value, err error) {
	a := affine.NumList
	switch v := i.(type) {
	case []float64:
		ret = makeValue(a, r.ValueOf(v), subtype)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newTextList(i interface{}, subtype string) (ret rt.Value, err error) {
	a := affine.TextList
	switch v := i.(type) {
	case []string:
		ret = makeValue(a, r.ValueOf(v), subtype)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}
