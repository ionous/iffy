package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
)

type Relative struct {
	*RefObject
	rel *RefRelation
}

// SetValue mirrors Object.SetValue.
// Changing the value of a relative property changes the status of the connection.
// note: the set of the pointer values is done in the relation class.
func (l *Relative) SetValue(name string, v interface{}) (err error) {
	id := id.MakeId(name)
	major, minor := &l.rel.props[0], &l.rel.props[1]
	if id == major.id {
		// decode v, call relate, and update our inner value
		if primary, ok := v.(rt.Object); !ok && v != nil {
			err = errutil.New("parameter is not an object")
		} else if other, e := l.getField(minor); e != nil {
			err = e
		} else if rec, e := l.rel.relate(primary, other); e != nil {
			err = e
		} else {
			l.RefObject = rec
		}
	} else if id == minor.id {
		// decode v, call relate, and update our inner value
		if secondary, ok := v.(rt.Object); !ok && v != nil {
			err = errutil.New("parameter is not an object")
		} else if other, e := l.getField(major); e != nil {
			err = e
		} else if rec, e := l.rel.relate(other, secondary); e != nil {
			err = e
		} else {
			l.RefObject = rec
		}
	} else {
		err = l.setValue(id, v)
	}
	return
}

func (l *Relative) getField(prop *RefInfo) (ret rt.Object, err error) {
	field := l.rval.FieldByIndex(prop.fieldPath)
	if obj, ok := field.Interface().(rt.Object); !ok && obj != nil {
		err = errutil.New("field is not an object!?", prop.id)
	} else {
		ret = obj
	}
	return
}
