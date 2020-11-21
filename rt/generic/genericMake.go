package generic

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

// ValueFrom adds an optional subtype to values, for example:
// marking text as specifically intended for aspects, traits, etc.
func ValueFrom(i interface{}, a affine.Affinity, subtype string) (ret Value, err error) {
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
		ret, err = newRecord(i, subtype)
	case affine.RecordList:
		ret, err = newRecordList(i, subtype)
	default:
		err = errutil.New("unhandled affinity", a)
	}
	return
}

func BoolOf(v bool) Value {
	return makeValue(affine.Bool, defaultType, v)
}
func StringOf(v string) Value {
	return makeValue(affine.Text, defaultType, v)
}
func FloatOf(v float64) Value {
	return makeValue(affine.Number, defaultType, v)
}
func IntOf(v int) Value {
	return makeValue(affine.Number, defaultType, v)
}
func RecordOf(v *Record) Value {
	return makeValue(affine.Record, v.Type(), v)
}
func RecordsOf(typeName string, vs []*Record) Value {
	return makeValue(affine.RecordList, typeName, &vs)
}
func StringsOf(vs []string) Value {
	return makeValue(affine.TextList, defaultType, &vs)
}
func FloatsOf(vs []float64) Value {
	return makeValue(affine.NumList, defaultType, &vs)
}

func makeValue(a affine.Affinity, subtype string, i interface{}) (ret refValue) {
	if len(subtype) == 0 {
		t := r.TypeOf(i)
		if t.Kind() == r.Ptr {
			t = t.Elem()
		}
		if t.Kind() == r.Slice {
			t = t.Elem()
		}
		subtype = t.String()
	}
	return refValue{a: a, i: i, t: subtype}
}

func newBoolValue(i interface{}, subtype string) (ret Value, err error) {
	a := affine.Bool
	switch v := i.(type) {
	case nil:
		// zero value for unhandled defaults in sqlite
		ret = False
	case bool:
		ret = makeValue(a, subtype, v)
	case int64:
		// sqlite, boolean values can be represented as 1/0
		ret = makeValue(a, subtype, v == 0)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newNumValue(i interface{}, subtype string) (ret Value, err error) {
	a := affine.Number
	switch v := i.(type) {
	case nil:
		// zero value for unhandled defaults in sqlite
		ret = Zero
	case int, int64, float64:
		ret = makeValue(a, subtype, v)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newTextValue(i interface{}, subtype string) (ret Value, err error) {
	a := affine.Text
	switch v := i.(type) {
	case nil:
		// zero value for unhandled defaults in sqlite
		ret = Empty
	case string:
		ret = makeValue(a, subtype, v)
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newNumList(i interface{}, subtype string) (ret Value, err error) {
	a := affine.NumList
	switch vs := i.(type) {
	case []float64:
		ret = makeValue(a, subtype, &vs)
	case *[]float64:
		ret = makeValue(a, subtype, vs)
	default:
		err = errutil.Fmt("unknown %s %T", a, vs)
	}
	return
}

func newTextList(i interface{}, subtype string) (ret Value, err error) {
	a := affine.TextList
	switch vs := i.(type) {
	case []string:
		ret = makeValue(a, subtype, &vs)
	case *[]string:
		ret = makeValue(a, subtype, vs)
	default:
		err = errutil.Fmt("unknown %s %T", a, vs)
	}
	return
}

func newRecord(i interface{}, subtype string) (ret Value, err error) {
	a := affine.Record
	if v, ok := i.(*Record); !ok {
		err = errutil.Fmt("unknown %s %T", a, i)
	} else if t := v.Type(); len(subtype) > 0 && t != subtype {
		err = errutil.New("mismatched record types", a, t, subtype)
	} else {
		ret = makeValue(a, t, v)
	}
	return
}

// note: doesnt check individual record types.
func newRecordList(i interface{}, subtype string) (ret Value, err error) {
	a := affine.RecordList
	switch vs := i.(type) {
	case []*Record:
		ret = makeValue(a, subtype, &vs)
	case *[]*Record:
		ret = makeValue(a, subtype, vs)
	default:
		err = errutil.Fmt("unknown %s %T", a, vs)
	}
	return
}
