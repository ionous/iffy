package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

// create a new generic value capable of supporting the passed affinity
// from a snapshot of the passed value; errors if the two types are not compatible.
func CopyValue(run rt.Runtime, a affine.Affinity, val rt.Value) (ret rt.Value, err error) {
	switch a {
	case affine.Bool:
		if v, e := rt.GetBool(run, val); e != nil {
			err = e
		} else {
			ret = &Bool{Value: v}
		}
	case affine.Number:
		if v, e := rt.GetNumber(run, val); e != nil {
			err = e
		} else {
			ret = &Float{Value: v}
		}
	case affine.Text:
		if v, e := rt.GetText(run, val); e != nil {
			err = e
		} else {
			ret = &String{Value: v}
		}
	default:
		err = errutil.New("failed to copy value, expected %s got %v(%T)", a, val, val)
	}
	return
}
