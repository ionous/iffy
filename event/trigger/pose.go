package trigger

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Pose provides runtime access to the results of parsed action.
// It mimics an object with properties for the noun and second nouns.
type Pose struct {
	objs []rt.Object
}

// Id is invalid. Anonymous.
func (po Pose) Id() (none ident.Id) {
	return
}

// Type returns a description of Pose properties.
func (po Pose) Type() r.Type {
	type Properties struct {
		First, Second rt.Object
	}
	return r.TypeOf((*Properties)(nil)).Elem()
}

// GetValue for noun and second nouns.
func (po Pose) GetValue(name string, pv interface{}) (err error) {
	switch name {
	case "noun":
		err = po.getObject(0, pv)
	case "second noun":
		err = po.getObject(1, pv)
		// FUTURE: case text/label.
	}
	return
}

// SetValue for Pose always returns error.
func (po Pose) SetValue(string, interface{}) error {
	return errutil.New("parser objects are read only")
}

func (po Pose) getObject(i int, pv interface{}) (err error) {
	if cnt := len(po.objs); i >= cnt {
		err = errutil.Fmt("too few nouns. wanted %d have %d", i+1, cnt)
	} else {
		obj := po.objs[i]
		err = ref.CoerceValue(pv, obj)
	}
	return
}
