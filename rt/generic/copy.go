package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

// CopyValue: create a new value from a snapshot of the passed value
func CopyValue(val Value) (ret interface{}, err error) {
	if val == nil {
		err = errutil.New("failed to copy nil value")
	} else {
		switch a := val.Affinity(); a {
		case affine.Bool:
			ret = val.Bool()
		case affine.Number:
			ret = val.Float()
		case affine.Text:
			ret = val.String()
		case affine.Record:
			if next, e := copyRecord(val.Record()); e != nil {
				err = e
			} else {
				ret = next
			}
		case affine.NumList:
			ret = copyFloats(val.Floats())

		case affine.TextList:
			ret = copyStrings(val.Strings())

		case affine.RecordList:
			vs := val.Records()
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
