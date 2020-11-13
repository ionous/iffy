package generic

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

// ValueOf returns a new generic value wrapper based on analyzing the passed value.
func ValueOf(i interface{}) (ret Value, err error) {
	switch i.(type) {
	case bool, *bool:
		ret, err = newBoolValue(i, defaultType)
	case int, int64, float32, float64, *int, *int64, *float32, *float64:
		ret, err = newNumValue(i, defaultType)
	case string, *string:
		ret, err = newTextValue(i, defaultType)
	case []float64:
		ret, err = newNumList(i, defaultType)
	case []string:
		ret, err = newTextList(i, defaultType)
	case *Record:
		ret, err = newRecord(i, defaultType)
	// case []*Record:
	// 	ret, err = newRecordList(i, defaultType)
	default:
		err = errutil.New("unhandled value %v(%T)", i, i)
	}
	return
}

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
	return makeValue(affine.Bool, defaultType, r.ValueOf(v))
}
func StringOf(v string) Value {
	return makeValue(affine.Text, defaultType, r.ValueOf(v))
}
func FloatOf(v float64) Value {
	return makeValue(affine.Number, defaultType, r.ValueOf(v))
}
func RecordOf(v *Record) Value {
	return makeValue(affine.Record, v.Type(), r.ValueOf(v))
}
func RecordsOf(k *Kind, vs []*Record) Value {
	return makeValue(affine.RecordList, k.Name(), r.ValueOf(vs))
}

func StringsOf(vs []string) Value {
	return makeValue(affine.TextList, defaultType, r.ValueOf(vs))
}

func FloatsOf(vs []float64) Value {
	return makeValue(affine.NumList, defaultType, r.ValueOf(vs))
}

func makeValue(a affine.Affinity, subtype string, v r.Value) (ret refValue) {
	if len(subtype) == 0 {
		t := v.Type()
		if v.Kind() == r.Slice {
			t = t.Elem()
		}
		subtype = t.String()
	}
	return refValue{a: a, v: v, t: subtype}
}

func newBoolValue(i interface{}, subtype string) (ret Value, err error) {
	a := affine.Bool
	switch v := i.(type) {
	case nil:
		// zero value for unhandled defaults in sqlite
		ret = False
	case bool:
		ret = makeValue(a, subtype, r.ValueOf(v))
	case int64:
		// sqlite, boolean values can be represented as 1/0
		ret = makeValue(a, subtype, r.ValueOf(v == 0))
	case *bool:
		// creates a dynamic value
		ret = makeValue(a, subtype, r.ValueOf(v).Elem())
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
		ret = makeValue(a, subtype, r.ValueOf(v))
	case *float64:
		// creates a dynamic value
		ret = makeValue(a, subtype, r.ValueOf(v).Elem())
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
		ret = makeValue(a, subtype, r.ValueOf(v))
	case *string:
		// creates a dynamic value
		ret = makeValue(a, subtype, r.ValueOf(v).Elem())
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newNumList(i interface{}, subtype string) (ret Value, err error) {
	a := affine.NumList
	switch v := i.(type) {
	case []float64:
		ret = makeValue(a, subtype, r.ValueOf(v))
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newTextList(i interface{}, subtype string) (ret Value, err error) {
	a := affine.TextList
	switch v := i.(type) {
	case []string:
		ret = makeValue(a, subtype, r.ValueOf(v))
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}

func newRecord(i interface{}, subtype string) (ret Value, err error) {
	a := affine.Record
	if v, ok := i.(*Record); !ok {
		err = errutil.Fmt("unknown %s %T", a, i)
	} else if t := v.Type(); len(subtype) > 0 && t != subtype {
		err = errutil.Fmt("mismatched record types", a, t, subtype)
	} else {
		ret = makeValue(a, t, r.ValueOf(v))
	}
	return
}

// note: doesnt check individual record types.
func newRecordList(i interface{}, subtype string) (ret Value, err error) {
	a := affine.RecordList
	switch v := i.(type) {
	case []*Record:
		ret = makeValue(a, subtype, r.ValueOf(v))
	default:
		err = errutil.Fmt("unknown %s %T", a, v)
	}
	return
}
