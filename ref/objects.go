package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Objects with ids, findable by the game.
type Objects struct {
	ObjectMap
	classes *Classes
}
type ObjectMap map[string]*RefObject

func NewObjects(classes *Classes) *Objects {
	return &Objects{make(map[string]*RefObject), classes}
}

// NewObject from the passed class.
// Compatible with rt.Runtime.
func (or *Objects) NewObject(class string) (ret rt.Object, err error) {
	if cls, ok := or.classes.GetClass(class); !ok {
		err = errutil.New("no such class", class)
	} else {
		ret = or.newObject(cls.(*RefClass))
	}
	return
}

func (or *Objects) newObject(cls *RefClass) *RefObject {
	rval := r.New(cls.rtype).Elem()
	return &RefObject{"", rval, cls, or}
}

// RegisterValue wrapping the passed value.
// Compatible with unique.ValueRegistry
func (or *Objects) RegisterValue(rval r.Value) (err error) {
	if id, e := MakeId(rval); e != nil {
		err = e
	} else if obj, ok := or.ObjectMap[id]; ok {
		err = errutil.New("conflicting objects", id, obj, rval)
	} else if cls, e := or.classes.GetByType(rval.Type()); e != nil {
		err = e
	} else {
		or.ObjectMap[id] = &RefObject{id, rval, cls, or}
	}
	return
}

// Emplace wraps the passed value as an anonymous object.
// Compatible with rt.Runtime.
func (or *Objects) Emplace(i interface{}) (ret rt.Object, err error) {
	rval := valueOf(i)
	if cls, e := or.classes.GetByType(rval.Type()); e != nil {
		err = e
	} else {
		ret = &RefObject{"", rval, cls, or}
	}
	return
}

// FindValue returns the originally specified object, not the wrapper.
// Compatible with unique.ValueRegistry
func (or *Objects) FindValue(name string) (ret r.Value, okay bool) {
	id := id.MakeId(name)
	if obj, ok := or.ObjectMap[id]; ok {
		ret, okay = obj.rval, true
	}
	return
}

// GetObject is compatible with rt.Runtime. The map can also be used directly.
func (or *Objects) GetObject(name string) (ret rt.Object, okay bool) {
	id := id.MakeId(name)
	ret, okay = or.ObjectMap[id]
	return
}

// GetByValue expects a pointer to a value, and it returns the ref object which wraps it.
// WARNING: it can return nil without error
func (or *Objects) GetByValue(rval r.Value) (ret *RefObject, err error) {
	if !rval.IsNil() {
		rval := rval.Elem()
		if id, e := MakeId(rval); e != nil {
			err = e
		} else if obj, ok := or.ObjectMap[id]; !ok {
			err = errutil.New("object not found", id)
		} else if obj.rval.Interface() != rval.Interface() {
			err = errutil.New("conflicting objects", id, obj, rval)
		} else {
			ret = obj
		}
	}
	return
}

// FIX: this is going to be way too slow for *casual use.
// an mru might of type might help,
// better might be caching the id path in the class,
// best might be forcing all classes to carry an explict id field as their first member.
func MakeId(rval r.Value) (ret string, err error) {
	rtype := rval.Type()
	if idField := FieldPathOfId(rtype); len(idField) == 0 {
		err = errutil.New("couldnt find id for", rtype)
	} else if name := rval.FieldByIndex(idField); name.Kind() != r.String {
		err = errutil.New("object needs an valid id", rval, rtype)
	} else {
		ret = id.MakeId(name.String())
	}
	return
}

// FieldPathOfId extracts the index of the id field
func FieldPathOfId(rtype r.Type) (ret []int) {
	for fw := unique.Fields(rtype); fw.HasNext(); {
		field := fw.GetNext()
		tag := unique.Tag(field.Tag)
		if _, ok := tag.Find("id"); ok {
			ret = append(field.Path, field.Index)
			break
		}
	}
	return
}
