package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
)

// Kinds database ( primarily for generating default values )
type Kinds interface {
	GetKindByName(n string) (*Kind, error)
}

// NewDefaultValue generates a zero value for the specified affinity;
// uses the passed Kinds to generate empty records when necessary.
func NewDefaultValue(ks Kinds, a affine.Affinity, subtype string) (ret Value, err error) {
	// return the default value for the
	switch a {
	case affine.Bool:
		ret = False

	case affine.Number:
		ret = Zero
	case affine.NumList:
		ret, err = newNumList(new([]float64), subtype)

	case affine.Text:
		ret = Empty
	case affine.TextList:
		ret, err = newTextList(new([]string), subtype)

	case affine.Record:
		if k, e := ks.GetKindByName(subtype); e != nil {
			err = errutil.New("unknown kind", subtype, e)
		} else {
			ret, err = newRecord(k.NewRecord(), subtype)
		}
	case affine.RecordList:
		if _, e := ks.GetKindByName(subtype); e != nil {
			err = errutil.New("unknown kind", subtype, e)
		} else {
			ret, err = newRecordList(new([]*Record), subtype)
		}
	default:
		err = errutil.New("unhandled affinity", a)
	}
	return
}
