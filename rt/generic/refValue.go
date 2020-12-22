package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/lang"
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
		panic(n.a.String() + " is not a number")
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
		panic(n.a.String() + " is not a number")
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
		panic(n.a.String() + " is not measurable")
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
		panic(n.a.String() + " is not indexable")
	}
	return
}

func (n refValue) FieldByName(f string) (ret Value, err error) {
	name := lang.SpecialBreakcase(f)
	if v, e := n.Record().GetNamedField(name); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func (n refValue) SetFieldByName(f string, v Value) (err error) {
	name := lang.SpecialBreakcase(f)
	return n.Record().SetNamedField(name, v)
}

func (n refValue) SetIndex(i int, v Value) {
	switch vp := n.i.(type) {
	case *[]float64:
		(*vp)[i] = v.Float()
	case *[]string:
		(*vp)[i] = v.String()
	case *[]*Record:
		if n.t != v.Type() {
			panic("record types dont match")
		}
		(*vp)[i] = v.Record()
	default:
		panic(n.a.String() + " is not index writable")
	}
}

func (n refValue) Slice(i, j int) (ret Value, err error) {
	if i < 0 {
		err = Underflow{i, 0}
	} else if cnt := n.Len(); j > cnt {
		err = Overflow{j, cnt}
	} else if i > j {
		err = errutil.New("bad range", i, j)
	} else {
		switch n.a {
		case affine.NumList:
			vp := n.i.(*[]float64)
			ret = FloatsOf(copyFloats((*vp)[i:j]))

		case affine.TextList:
			vp := n.i.(*[]string)
			ret = StringsOf(copyStrings((*vp)[i:j]))

		case affine.RecordList:
			vp := n.i.(*[]*Record)
			ret = RecordsOf(n.Type(), copyRecords((*vp)[i:j]))

		default:
			panic(n.a.String() + " is not sliceable")
		}
	}
	return
}

func (n refValue) Splice(i, j int, add Value) (ret Value, err error) {
	if i < 0 {
		err = Underflow{i, 0}
	} else if cnt := n.Len(); j > cnt {
		err = Overflow{j, cnt}
	} else if i > j {
		err = errutil.New("bad range", i, j)
	} else {
		switch n.a {
		case affine.NumList:
			vp := n.i.(*[]float64)
			els := (*vp)
			cut := copyFloats(els[i:j])
			ins := normalizeFloats(add)
			(*vp) = append(els[:i], append(ins, els[j:]...)...)
			ret = FloatsOf(cut)

		case affine.TextList:
			vp := n.i.(*[]string)
			els := (*vp)
			cut := copyStrings(els[i:j])
			ins := normalizeStrings(add)
			(*vp) = append(els[:i], append(ins, els[j:]...)...)
			ret = StringsOf(cut)

		case affine.RecordList:
			vp := n.i.(*[]*Record)
			if n.t != add.Type() {
				panic("record types dont match")
			}
			els := (*vp)
			cut := copyRecords(els[i:j])
			ins := normalizeRecords(add)
			(*vp) = append(els[:i], append(ins, els[j:]...)...)
			ret = RecordsOf(n.t, cut)

		default:
			panic(n.a.String() + " is not spliceable")
		}
	}
	return
}

func (n refValue) Append(add Value) {
	switch n.a {
	case affine.NumList:
		vp := n.i.(*[]float64)
		ins := normalizeFloats(add)
		(*vp) = append((*vp), ins...)

	case affine.TextList:
		vp := n.i.(*[]string)
		ins := normalizeStrings(add)
		(*vp) = append((*vp), ins...)

	case affine.RecordList:
		vp := n.i.(*[]*Record)
		if n.t != add.Type() {
			panic("record types dont match")
		}
		ins := normalizeRecords(add)
		(*vp) = append((*vp), ins...)

	default:
		panic(n.a.String() + " is not appendable")
	}
}
