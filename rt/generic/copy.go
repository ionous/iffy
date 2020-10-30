package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

// create a new generic value capable of supporting the passed affinity
// from a snapshot of the passed value; errors if the two types are not compatible.
func CopyValue(a affine.Affinity, val rt.Value) (ret rt.Value, err error) {
	if val == nil {
		err = errutil.New("failed to copy nil value")
	} else {
		switch a {
		case affine.Bool:
			if v, e := val.GetBool(); e != nil {
				err = e
			} else {
				ret = &Bool{Value: v}
			}
		case affine.Number:
			if v, e := val.GetNumber(); e != nil {
				err = e
			} else {
				ret = &Float{Value: v}
			}
		case affine.Text:
			if v, e := val.GetText(); e != nil {
				err = e
			} else {
				ret = &String{Value: v}
			}
		case affine.Object:
			ret = val // fix: we shouldnt be trying to copy objects.
		default:
			err = errutil.Fmt("failed to copy value, expected %s got %v(%T)", a, val, val)
		}
	}
	return
}
