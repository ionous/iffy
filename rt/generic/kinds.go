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
		ret = BoolFrom(false, subtype)

	case affine.Number:
		ret = FloatFrom(0, subtype)

	case affine.NumList:
		ret = FloatsFrom(nil, subtype)

	case affine.Text:
		ret = StringFrom("", subtype)

	case affine.TextList:
		ret = StringsFrom(nil, subtype)

	case affine.Record:
		if k, e := ks.GetKindByName(subtype); e != nil {
			err = errutil.New("unknown kind", subtype, e)
		} else {
			ret = RecordFrom(k.NewRecord(), subtype)
		}

	case affine.RecordList:
		if _, e := ks.GetKindByName(subtype); e != nil {
			err = errutil.New("unknown kind", subtype, e)
		} else {
			ret = RecordsFrom(nil, subtype)
		}

	default:
		err = errutil.New("default value requested for unhandled affinity", a)
	}
	return
}
