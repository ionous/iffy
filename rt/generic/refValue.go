package generic

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

type refValue struct {
	a affine.Affinity
	v r.Value
	t string
}

var _ Value = (*refValue)(nil)

func (n refValue) Affinity() (ret affine.Affinity) {
	return n.a
}

func (n refValue) Type() string {
	return n.t
}

func (n refValue) GetBool() (ret bool, err error) {
	if n.v.Kind() != r.Bool {
		err = errutil.New("value is not a bool")
	} else {
		ret = n.v.Bool()
	}
	return
}

func (n refValue) GetNumber() (ret float64, err error) {
	switch k := n.v.Kind(); k {
	case r.Float32, r.Float64:
		ret = n.v.Float()
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		ret = float64(n.v.Int())
	default:
		err = errutil.New("value is not a number")
	}
	return
}

func (n refValue) GetText() (ret string, err error) {
	if n.v.Kind() != r.String {
		err = errutil.New("value is not a text")
	} else {
		ret = n.v.String()
	}
	return
}

func (n refValue) GetRecord() (ret *Record, err error) {
	if v, ok := n.v.Interface().(*Record); !ok {
		err = errutil.New("value is not a record")
	} else {
		ret = v
	}
	return
}

func (n refValue) GetNumList() (ret []float64, err error) {
	if vs, ok := n.v.Interface().([]float64); !ok {
		err = errutil.New("value is not a number list")
	} else {
		ret = vs
	}
	return
}
func (n refValue) GetTextList() (ret []string, err error) {
	if vs, ok := n.v.Interface().([]string); !ok {
		err = errutil.New("value is not a text list")
	} else {
		ret = vs
	}
	return
}
func (n refValue) GetRecordList() (ret []*Record, err error) {
	if vs, ok := n.v.Interface().([]*Record); !ok {
		err = errutil.New("value is not a record list")
	} else {
		ret = vs
	}
	return
}
func (n refValue) GetLen() (ret int, err error) {
	if n.v.Kind() != r.Slice {
		err = errutil.New("value is not measurable")
	} else {
		ret = n.v.Len()
	}
	return
}
func (n refValue) GetIndex(i int) (ret Value, err error) {
	if e := n.validateIndex(i); e != nil {
		err = e
	} else if elAffinity := affine.Element(n.a); len(elAffinity) == 0 {
		err = errutil.New("unknown list affinity", n.a)
	} else {
		ret = makeValue(elAffinity, n.t, n.v.Index(i))
	}
	return
}

func (n refValue) validateIndex(i int) (err error) {
	if n.v.Kind() != r.Slice {
		err = errutil.New("value is not indexable")
	} else if i < 0 {
		err = Underflow{i, 0}
	} else if cnt := n.v.Len(); i >= cnt {
		err = Overflow{i, cnt}
	}
	return
}

func (n refValue) GetNamedField(string) (ret Value, err error) {
	err = errutil.New("value doesnt have fields")
	return
}

func (n refValue) SetNamedField(string, Value) (err error) {
	err = errutil.New("value is not field writable")
	return
}

func (n refValue) SetIndexedValue(i int, v Value) (err error) {
	if e := n.validateIndex(i); e != nil {
		err = e
	} else if va, elAffinity := v.Affinity(), affine.Element(n.a); va != elAffinity {
		err = errutil.Fmt("mismatched affinity %q for element %q", va, elAffinity)
	} else if refv, ok := v.(refValue); !ok {
		err = errutil.Fmt("unknown source value %v(%T)", v, v)
	} else {
		n.v.Index(i).Set(refv.v)
	}
	return
}

// note: this can grow record slices with nil values.
func (n refValue) Resize(newLen int) (ret Value, err error) {
	if vs := n.v; vs.Kind() != r.Slice {
		err = errutil.New("value is not indexable")
	} else if newLen < 0 {
		err = Underflow{newLen, 0}
	} else if cap := vs.Cap(); newLen <= cap {
		vs.SetLen(newLen) // shrinking
		ret = n           // the slice memory stays the same.
	} else if grow := newLen - n.v.Len(); grow > 0 {
		// grow using make, append ( versus make, copy )
		// to trigger go's grow padding
		blanks := r.MakeSlice(vs.Type().Elem(), grow, grow)
		ret = makeValue(n.a, n.t, r.AppendSlice(vs, blanks))
	}
	return
}

func (n refValue) Slice(i, j int) (ret Value, err error) {
	if vs := n.v; vs.Kind() != r.Slice {
		err = errutil.New("value is not indexable")
	} else if i < 0 {
		err = Underflow{i, 0}
	} else if cnt := vs.Len(); j > cnt {
		err = Overflow{j, cnt}
	} else if i > j {
		err = errutil.New("bad range", i, j)
	} else {
		ret = makeValue(n.a, n.t, vs.Slice(i, j))
	}
	return
}

func (n refValue) Append(v Value) (ret Value, err error) {
	if vs := n.v; vs.Kind() != r.Slice {
		err = errutil.New("value is not indexable")
	} else if elAffinity := affine.Element(v.Affinity()); len(elAffinity) == 0 {
		ret, err = n.appendOne(v)
	} else {
		ret, err = n.appendMany(v)
	}
	return
}

func (n refValue) appendOne(v Value) (ret Value, err error) {
	elAffinity := affine.Element(v.Affinity())
	compatible := n.a == elAffinity && (elAffinity != affine.Record || v.Type() == n.t)
	if !compatible {
		err = errutil.New("value is not compatible with list")
	} else if refv, ok := v.(refValue); !ok {
		err = errutil.Fmt("unknown source value %v(%T)", v, v)
	} else {
		ret = makeValue(n.a, n.t, r.Append(n.v, refv.v))
	}
	return
}

func (n refValue) appendMany(v Value) (ret Value, err error) {
	va := v.Affinity()
	compatible := n.a == va && (va != affine.RecordList || v.Type() == n.t)
	if !compatible {
		err = errutil.New("value is not compatible with list")
	} else if refv, ok := v.(refValue); !ok {
		err = errutil.Fmt("unknown source value %v(%T)", v, v)
	} else {
		ret = makeValue(n.a, n.t, r.AppendSlice(n.v, refv.v))
	}
	return
}

// func (n refValue) SetValue(from Value) (err error) {
// 	if val := n.v; !val.CanSet() {
// 		err = errutil.New("value is not writable")
// 	} else {
// 		switch k := val.Kind(); k {
// 		case r.Bool:
// 			if v, e := from.GetBool(); e != nil {
// 				err = e
// 			} else {
// 				val.SetBool(v)
// 			}
// 		case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
// 			if v, e := from.GetNumber(); e != nil {
// 				err = e
// 			} else {
// 				val.SetInt(int64(v))
// 			}
// 		case r.Float64, r.Float32:
// 			if v, e := from.GetNumber(); e != nil {
// 				err = e
// 			} else {
// 				val.SetFloat(v)
// 			}
// 		case r.String:
// 			if v, e := from.GetText(); e != nil {
// 				err = e
// 			} else {
// 				val.SetString(v)
// 			}
// 		case r.Interface:
// 			// test each
// 			// val.Implements

// 		case r.Slice:
// 			switch k := n.v.Type().Elem().Kind(); k {
// 			case r.Float32, r.Float64, r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
// 				if vs, e := from.GetNumList(); e != nil {
// 					err = e
// 				} else {
// 					val.Set(r.ValueOf(vs))
// 				}
// 			case r.String:
// 				if vs, e := from.GetNumList(); e != nil {
// 					err = e
// 				} else {
// 					val.Set(r.ValueOf(vs))
// 				}
// 			case r.Ptr:
// 				//ret = affine.RecordList
// 			default:
// 				panic("unknown list " + k.String())
// 			}

// 		default:
// 			panic("unknown kind " + k.String())
// 		}

// 	}
// 	return
// }
