package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

// every primitive value is its own unique instance.
// records are as pointers, and lists as pointers to slices;
// their data is therefore shared across refValues.
// pointers to slices allows for in-place grow, append, etc.
type refValue struct {
	a affine.Affinity
	t string
	i interface{}
}

var _ Value = (*refValue)(nil)

func (n refValue) Affinity() affine.Affinity {
	return n.a
}

func (n refValue) Type() string {
	return n.t
}

func (n refValue) Bool() bool {
	return n.i.(bool)
}

func (n refValue) Float() (ret float64) {
	switch v := n.i.(type) {
	case float64:
		ret = v
	case float32:
		ret = float64(v)
	case int:
		ret = float64(v)
	case int64:
		ret = float64(v)
	default:
		panic(errutil.Sprintf("value %v(%T) is not a number", v, v))
	}
	return
}

func (n refValue) Int() (ret int) {
	switch v := n.i.(type) {
	case int:
		ret = v
	case int64:
		ret = int(v)
	case float32:
		ret = int(v)
	case float64:
		ret = int(v)
	default:
		panic("value is not a number")
	}
	return
}

func (n refValue) String() string {
	return n.i.(string)
}

func (n refValue) Record() *Record {
	return n.i.(*Record)
}

func (n refValue) Floats() (ret []float64) {
	vp := n.i.(*[]float64)
	return *vp
}

func (n refValue) Strings() (ret []string) {
	vp := n.i.(*[]string)
	return *vp
}

func (n refValue) Records() (ret []*Record) {
	vp := n.i.(*[]*Record)
	return *vp
}

func (n refValue) Len() (ret int) {
	switch vp := n.i.(type) {
	case *[]float64:
		ret = len(*vp)
	case *[]string:
		ret = len(*vp)
	case *[]*Record:
		ret = len(*vp)
	default:
		panic("value is not measurable")
	}
	return
}

func (n refValue) Index(i int) (ret Value) {
	switch vp := n.i.(type) {
	case *[]float64:
		ret = FloatOf((*vp)[i])
	case *[]string:
		ret = StringOf((*vp)[i])
	case *[]*Record:
		ret = RecordOf((*vp)[i])
	default:
		panic("value is not measurable")
	}
	return
}

func (n refValue) FieldByName(f string) (ret Value, err error) {
	if v, e := n.Record().GetNamedField(f); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func (n refValue) SetFieldByName(f string, v Value) (err error) {
	return n.Record().SetNamedField(f, v)
}

func (n refValue) SetIndex(i int, v Value) {
	switch vp := n.i.(type) {
	case *[]float64:
		(*vp)[i] = v.Float()
	case *[]string:
		(*vp)[i] = v.String()
	case *[]*Record:
		if n.Type() != v.Type() {
			panic("record types dont match")
		}
		(*vp)[i] = v.Record()
	default:
		panic("value is not measurable")
	}
}

// note: this can grow record slices with nil values.
// func (n refValue) Resize(newLen int) {
// 	if vs := n.v; vs.Kind() != r.Slice {
// 		panic("value is not indexable")
// 	} else if newLen < 0 {
// 		err = Underflow{newLen, 0}
// 	} else if cap := vs.Cap(); newLen <= cap {
// 		vs.SetLen(newLen) // shrinking; the slice memory stays the same.
// 	} else if grow := newLen - n.v.Len(); grow > 0 {
// 		// grow using make, append ( versus make, copy )
// 		// to trigger go's grow padding
// 		blanks := r.MakeSlice(vs.Type().Elem(), grow, grow)
// 		n.v = r.AppendSlice(vs, blanks)
// 	}
// 	return
// }

func (n refValue) Slice(i, j int) (ret Value, err error) {
	if i < 0 {
		err = Underflow{i, 0}
	} else if cnt := n.Len(); j > cnt {
		err = Overflow{j, cnt}
	} else if i > j {
		err = errutil.New("bad range", i, j)
	} else {
		switch vp := n.i.(type) {
		case *[]float64:
			vs := CopyFloats((*vp)[i:j])
			ret = FloatsOf(vs)
		case *[]string:
			vs := CopyStrings((*vp)[i:j])
			ret = StringsOf(vs)
		case *[]*Record:
			vs := CopyRecords((*vp)[i:j])
			ret = RecordsOf(n.Type(), vs)
		default:
			panic("value is not sliceable")
		}
	}
	return
}

func (n refValue) Append(v Value) {
	if !affine.IsList(v.Affinity()) {
		n.appendOne(v)
	} else {
		n.appendMany(v)
	}
	return
}

func (n refValue) appendOne(v Value) {
	switch vp := n.i.(type) {
	case *[]float64:
		(*vp) = append((*vp), v.Float())
	case *[]string:
		(*vp) = append((*vp), v.String())
	case *[]*Record:
		if n.Type() != v.Type() {
			panic("record types dont match")
		}
		(*vp) = append((*vp), v.Record())
	default:
		panic("value is not extensible")
	}
}

func (n refValue) appendMany(v Value) {
	switch vp := n.i.(type) {
	case *[]float64:
		(*vp) = append((*vp), v.Floats()...)
	case *[]string:
		(*vp) = append((*vp), v.Strings()...)
	case *[]*Record:
		if n.Type() != v.Type() {
			panic("record types dont match")
		}
		(*vp) = append((*vp), v.Records()...)
	default:
		panic("value is not extensible")
	}
	return
}
