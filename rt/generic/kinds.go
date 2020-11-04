package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

// Kinds database ( primarily for generating default values )
type Kinds interface {
	KindByName(name string) (*Kind, error)
}

func MakeDefault(ks Kinds, affinity affine.Affinity, typeName string) (ret rt.Value, err error) {
	// return the default value for the
	switch a := affinity; a {
	case affine.Bool:
		ret = &Bool{}
	case affine.Number:
		ret = &Float{}
	case affine.NumList:
		ret = &FloatSlice{}
	case affine.Text:
		ret = &String{}
	case affine.TextList:
		ret = &StringSlice{}
	case affine.Record:
		if n, e := ks.KindByName(typeName); e != nil {
			err = errutil.New("unknown kind", typeName, e)
		} else {
			ret = n.NewRecord()
		}
	case affine.RecordList:
		if n, e := ks.KindByName(typeName); e != nil {
			err = errutil.New("unknown kind", typeName, e)
		} else {
			ret = n.NewRecordSlice()
		}
	default:
		err = errutil.New("unhandled affinity", a)
	}
	return
}
