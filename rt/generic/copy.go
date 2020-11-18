package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

// CopyFloats: duplicate the passed slice.
// ( b/c golang's built in copy doesnt allocate )
func CopyFloats(src []float64) []float64 {
	out := make([]float64, len(src))
	copy(out, src)
	return out
}

// CopyStrings: duplicate the passed slice.
// ( b/c golang's built in copy doesnt allocate )
func CopyStrings(src []string) []string {
	out := make([]string, len(src))
	copy(out, src)
	return out
}

// CopyValue: create a new value from a snapshot of the passed value
func CopyValue(val Value) (ret interface{}, err error) {
	if val == nil {
		err = errutil.New("failed to copy nil value")
	} else {
		switch a := val.Affinity(); a {
		case affine.Bool:
			if v, e := val.GetBool(); e != nil {
				err = e
			} else {
				ret = v
			}
		case affine.Number:
			if v, e := val.GetNumber(); e != nil {
				err = e
			} else {
				ret = v
			}
		case affine.Text:
			if v, e := val.GetText(); e != nil {
				err = e
			} else {
				ret = v
			}
		case affine.Record:
			if v, e := val.GetRecord(); e != nil {
				err = e
			} else if next, e := copyRecord(v); e != nil {
				err = e
			} else {
				ret = next
			}
		case affine.NumList:
			if vs, e := val.GetNumList(); e != nil {
				err = e
			} else {
				ret = CopyFloats(vs)
			}
		case affine.TextList:
			if vs, e := val.GetTextList(); e != nil {
				err = e
			} else {
				ret = CopyStrings(vs)
			}
		case affine.RecordList:
			if vs, e := val.GetRecordList(); e != nil {
				err = e
			} else {
				values := make([]*Record, len(vs))
				for i, el := range vs {
					if cpy, e := copyRecord(el); e != nil {
						err = e
						break
					} else {
						values[i] = cpy
					}
				}
				if err == nil {
					ret = values
				}
			}
		case affine.Object:
			// new nouns cant be dynamically added to the runtime.
			err = errutil.New("can't duplicate object values")

		default:
			err = errutil.Fmt("failed to copy value, expected %s got %v(%T)", a, val, val)
		}
	}
	return
}

// assumes in value is a record.
func copyRecord(v *Record) (ret *Record, err error) {
	cnt := v.kind.NumField()
	values := make([]interface{}, cnt, cnt)
	for i := 0; i < cnt; i++ {
		if el, e := v.GetFieldByIndex(i); e != nil {
			err = e
			break
		} else if cpy, e := CopyValue(el); e != nil {
			err = e
			break
		} else {
			values[i] = cpy
		}
	}
	if err == nil {
		ret = &Record{kind: v.kind, values: values}
	}
	return
}
