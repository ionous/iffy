package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/index"
	r "reflect"
)

// Relative represents a cell in a table relation.
type RefRelative struct {
	*RefObject
	rel *RefRelation
}

// SetValue mirrors Object.SetValue.
// Changing the value of a relative property changes the status of the connection.
// note: the set of the pointer values is done in the relation class.
func (l *RefRelative) SetValue(name string, v interface{}) (err error) {
	id := id.MakeId(name)
	major, minor := &l.rel.props[0], &l.rel.props[1]
	if id == major.id {
		l.onChanged(major, minor, v,
			func(minor, oldMajor, newMajor string) (kd *index.KeyData) {
				if len(oldMajor) > 0 {
					l.rel.table.DeletePair(oldMajor, minor)
				}
				if len(newMajor) > 0 {
					kd, _ = l.rel.table.Relate(newMajor, minor)
				}
				return
			})
	} else if id == minor.id {
		l.onChanged(minor, major, v,
			func(major, oldMinor, newMinor string) (kd *index.KeyData) {
				if len(oldMinor) > 0 {
					l.rel.table.DeletePair(major, oldMinor)
				}
				if len(newMinor) > 0 {
					kd, _ = l.rel.table.Relate(major, newMinor)
				}
				return
			})
	} else {
		err = l.setValue(id, v)
	}
	return
}

func setPointer(field r.Value, next *RefObject) {
	if next == nil {
		field.Set(r.Zero(field.Type()))
	} else {
		field.Set(next.rval.Addr())
	}
}

// k is the constant id from the far field; guaranteed to exist
// p is the previous value of the near field
// n is the next value of the near field
// returns keydata if there is a matching cell in the table.
type RelChanged func(k, p, n string) *index.KeyData

func (l *RefRelative) onChanged(near, far *RefInfo, v interface{}, cb RelChanged) (err error) {
	objects := l.rel.relations.objects
	// if we kill a pair, we are just killing the pair and nothing else.
	// we can even keep our memory so long as we alter the pointer
	if next, ok := v.(*RefObject); !ok && v != nil {
		err = errutil.New("parameter is not an object")
	} else {
		nearField := l.rval.FieldByIndex(near.fieldPath)
		if prev, e := objects.GetByValue(nearField); e != nil {
			err = e
		} else if prev != next {
			// set our field
			setPointer(nearField, next)
			// now, talk to the table
			// note: it only makes sense to do so if the far side exists.
			// ( because we need pairs of values to change the table )
			farField := l.rval.FieldByIndex(far.fieldPath)
			if farVal, e := objects.GetByValue(farField); e != nil {
				err = e
			} else if farVal != nil {
				var k, p, n string
				k = farVal.GetId()
				if prev != nil {
					p = prev.GetId()
				}
				if next != nil {
					n = next.GetId()
				}
				if kd := cb(k, p, n); kd != nil {
					// then we want to establish or connect to something new.
					if ref, ok := kd.Data.(*RefObject); !ok {
						kd.Data = l.RefObject // establish
					} else {
						l.RefObject = ref // connect
					}
				}
			}
		}
	}
	return
}
