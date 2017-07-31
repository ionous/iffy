package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	r "reflect"
	"strconv"
)

// ObjBuilder with ids, findable by the game.
type ObjBuilder struct {
	queue   queue
	classes ClassMap
	info    infoMap
}

type infoMap map[r.Type]classInfo

type classInfo struct {
	pathOfId []int
	unnamed  int
}

type queued struct {
	rval r.Value
	cls  *RefClass
}
type queue map[string]queued

func NewObjects(classes *ClassBuilder) *ObjBuilder {
	return &ObjBuilder{make(queue), classes.ClassMap, make(infoMap)}
}

func (b *ObjBuilder) Build() *Objects {
	objs := &Objects{make(ObjectMap), b.classes}
	for id, q := range b.queue {
		objs.ObjectMap[id] = &RefObject{id, q.rval, q.cls, objs}
	}
	return objs
}

// RegisterValue the passed value as a future object.
// Compatible with unique.ValueRegistry
func (b *ObjBuilder) RegisterValue(rval r.Value) (err error) {
	if id, e := b.MakeId(rval); e != nil {
		err = e
	} else if obj, ok := b.queue[id]; ok {
		err = errutil.New("duplicate object", id, obj, rval)
	} else if cls, e := b.classes.GetByType(rval.Type()); e != nil {
		err = e
	} else {
		b.queue[id] = queued{rval, cls}
	}
	return
}

func (b *ObjBuilder) MakeId(rval r.Value) (ret string, err error) {
	rtype := rval.Type()

	info, ok := b.info[rtype]
	if !ok {
		if path, ok := unique.PathOf(rtype, "id"); !ok {
			err = errutil.New("couldnt find id for", rtype)
		} else {
			info.pathOfId = path
			b.info[rtype] = info
		}
	}
	if err == nil {
		field := rval.FieldByIndex(info.pathOfId)
		name := field.String()
		// handle unnamed objects:
		if len(name) == 0 {
			info.unnamed++
			b.info[rtype] = info
			name = rtype.Name() + "#" + strconv.Itoa(info.unnamed)
			field.SetString(name)
		}
		ret = id.MakeId(name)
	}

	return
}
