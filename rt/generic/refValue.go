package generic

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type refValue struct {
	a affine.Affinity
	v r.Value
	t string
}

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
func (n refValue) GetRecordList() (ret []rt.Value, err error) {
	if vs, ok := n.v.Interface().([]rt.Value); !ok {
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
func (n refValue) GetIndex(i int) (ret rt.Value, err error) {
	if n.v.Kind() != r.Slice {
		err = errutil.New("value is not indexable")
	} else if i < 0 {
		err = rt.OutOfRange{i, 0}
	} else if cnt := n.v.Len(); i >= cnt {
		err = rt.OutOfRange{i, cnt}
	} else if a, e := indexedAffinity(n.a); e != nil {
		err = e
	} else {
		ret = refValue{a: a, t: n.t, v: n.v.Index(i)}
	}
	return
}

func indexedAffinity(list affine.Affinity) (ret affine.Affinity, err error) {
	switch a := list; a {
	case affine.TextList:
		ret = affine.Text
	case affine.NumList:
		ret = affine.Number
	case affine.RecordList:
		ret = affine.Record
	default:
		err = errutil.New("unknown affinity", a)
	}
	return
}

func (n refValue) GetField(string) (ret rt.Value, err error) {
	err = errutil.New("value doesnt have fields")
	return
}

func (n refValue) SetField(string, rt.Value) (err error) {
	err = errutil.New("value is not writable")
	return
}

func (n refValue) SetValue(from rt.Value) (err error) {
	if val := n.v; !val.CanSet() {
		err = errutil.New("value is not writable")
	} else {
		switch k := val.Kind(); k {
		case r.Bool:
			if v, e := from.GetBool(); e != nil {
				err = e
			} else {
				val.SetBool(v)
			}
		case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
			if v, e := from.GetNumber(); e != nil {
				err = e
			} else {
				val.SetInt(int64(v))
			}
		case r.Float64, r.Float32:
			if v, e := from.GetNumber(); e != nil {
				err = e
			} else {
				val.SetFloat(v)
			}
		case r.String:
			if v, e := from.GetText(); e != nil {
				err = e
			} else {
				val.SetString(v)
			}
		case r.Interface:
			// test each
			// val.Implements

		case r.Slice:
			switch k := n.v.Type().Elem().Kind(); k {
			case r.Float32, r.Float64, r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
				if vs, e := from.GetNumList(); e != nil {
					err = e
				} else {
					val.Set(r.ValueOf(vs))
				}
			case r.String:
				if vs, e := from.GetNumList(); e != nil {
					err = e
				} else {
					val.Set(r.ValueOf(vs))
				}
			case r.Ptr:
				//ret = affine.RecordList
			default:
				panic("unknown list " + k.String())
			}

		default:
			panic("unknown kind " + k.String())
		}

	}
	return
}
