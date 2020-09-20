package generic

import (
	"bytes"
	"encoding/gob"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

// Value implements a read and writable version of rt.Value.
// It can store ( and degob ) program byte arrays.
type Value struct {
	typeName string
	value    interface{}
}

// fix: in theory can change pattern management to remove this.
func (q *Value) FixMe() interface{} {
	return q.value
}

// one of tables.PRIM_
// a prim, a byte array to degob and cache, or the cached eval.
func NewValue(typeName string, value interface{}) *Value {
	return &Value{typeName, value}
}

func (q *Value) checkStored(typeName string) (err error) {
	if typeName != q.typeName {
		err = errutil.Fmt("Value: cant convert %s to %s", q.typeName, typeName)
	}
	return
}

func (q *Value) unexpectedValue(typeName string) error {
	return errutil.Fmt("Value: unexpected value [%v](%T) while returning %s", q.value, q.value, typeName)
}

func (q *Value) SetValue(run rt.Runtime, val rt.Value) (err error) {
	switch q.typeName {
	case tables.PRIM_BOOL:
		if v, e := rt.GetBool(run, val); e != nil {
			err = e
		} else {
			q.value = v
		}
	case tables.PRIM_DIGI:
		if v, e := rt.GetNumber(run, val); e != nil {
			err = e
		} else {
			q.value = v
		}
	case tables.PRIM_TEXT:
		if v, e := rt.GetText(run, val); e != nil {
			err = e
		} else {
			q.value = v
		}
	case "num_list":
		if vs, e := rt.GetNumList(run, val); e != nil {
			err = e
		} else {
			q.value = vs
		}
	case "text_list":
		if vs, e := rt.GetNumList(run, val); e != nil {
			err = e
		} else {
			q.value = vs
		}
	default:
		err = q.unexpectedValue(q.typeName)
	}
	return
}

func (q *Value) GetBool(run rt.Runtime) (ret bool, err error) {
	typeName := tables.PRIM_BOOL
	if e := q.checkStored(typeName); e != nil {
		err = e
	} else {
		switch a := q.value.(type) {
		case nil:
			// zero value for unhandled defaults in sqlite
		case bool:
			ret = a
		case int64: // sqlite, boolean values can be represented as 1/0
			ret = a != 0
		case rt.BoolEval:
			ret, err = a.GetBool(run)
		case []byte:
			var eval rt.BoolEval
			if e := q.unpack(a, &eval); e != nil {
				err = e
			} else {
				ret, err = eval.GetBool(run)
			}
		default:
			err = q.unexpectedValue(typeName)
		}
	}
	return
}

func (q *Value) GetNumber(run rt.Runtime) (ret float64, err error) {
	typeName := tables.PRIM_DIGI
	if e := q.checkStored(typeName); e != nil {
		err = e
	} else {
		switch a := q.value.(type) {
		case nil:
			// zero value for unhandled defaults in sqlite
		case int:
			ret = float64(a)
		case int64:
			ret = float64(a)
		case float64:
			ret = a
		case rt.NumberEval:
			ret, err = a.GetNumber(run)
		case []byte:
			var eval rt.NumberEval
			if e := q.unpack(a, &eval); e != nil {
				err = e
			} else {
				ret, err = eval.GetNumber(run)
			}
		default:
			err = q.unexpectedValue(typeName)
		}
	}
	return
}
func (q *Value) GetText(run rt.Runtime) (ret string, err error) {
	typeName := tables.PRIM_TEXT
	if e := q.checkStored(typeName); e != nil {
		err = e
	} else {
		switch a := q.value.(type) {
		case nil:
			// zero value for unhandled defaults in sqlite
		case string:
			ret = a
		case rt.TextEval:
			ret, err = a.GetText(run)
		case []byte:
			var eval rt.TextEval
			if e := q.unpack(a, &eval); e != nil {
				err = e
			} else {
				ret, err = eval.GetText(run)
			}
		default:
			err = q.unexpectedValue(typeName)
		}
	}
	return
}

func (q *Value) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	typeName := "num_list"
	if e := q.checkStored(typeName); e != nil {
		err = e
	} else {
		switch vs := q.value.(type) {
		case nil:
			ret = rt.EmptyStream(true)
		case []float64:
			ret = SliceFloats(vs)
		case rt.NumListEval:
			ret, err = vs.GetNumberStream(run)
		case []byte:
			var eval rt.NumListEval
			if e := q.unpack(vs, &eval); e != nil {
				err = e
			} else {
				ret, err = eval.GetNumberStream(run)
			}
		default:
			err = q.unexpectedValue(typeName)
		}
	}
	return
}

func (q *Value) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	typeName := "text_list"
	if e := q.checkStored(typeName); e != nil {
		err = e
	} else {
		switch vs := q.value.(type) {
		case nil:
			ret = rt.EmptyStream(true)
		case []string:
			ret = SliceStrings(vs)
		case rt.TextListEval:
			ret, err = vs.GetTextStream(run)
		case []byte:
			var eval rt.TextListEval
			if e := q.unpack(vs, &eval); e != nil {
				err = e
			} else {
				ret, err = eval.GetTextStream(run)
			}
		default:
			err = q.unexpectedValue(typeName)
		}
	}
	return
}

// convert the bytes ( if any ) to an eval
func (q *Value) unpack(b []byte, iptr interface{}) (err error) {
	rptr := r.ValueOf(iptr)
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	if e := dec.DecodeValue(rptr); e != nil {
		err = e
	} else {
		// garbage collection returns the byte array back into the aether.
		q.value = rptr.Elem().Interface()
	}
	return
}
