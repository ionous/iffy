package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
	"strconv"
)

// ObjBuilder with ids, findable by the game.
type ObjBuilder struct {
	queue queue
	info  infoMap
}

type infoMap map[r.Type]classInfo

type classInfo struct {
	pathOfId []int
	unnamed  int
}

type queued struct {
	rval r.Value
	cls  rt.Class
}
type queue map[string]queued

func NewObjects() *ObjBuilder {
	return &ObjBuilder{make(queue), make(infoMap)}
}

func (b *ObjBuilder) Build() *Objects {
	objs := &Objects{make(ObjectMap)}
	for id, q := range b.queue {
		objs.ObjectMap[id] = &RefObject{q.rval, objs}
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
	} else {
		b.queue[id] = queued{rval, rval.Type()}
	}
	return
}

func (b *ObjBuilder) MakeId(rval r.Value) (ret string, err error) {
	rtype := rval.Type()

	// cache the id path
	info, ok := b.info[rtype]
	if !ok {
		if path, ok := unique.PathOf(rtype, "id"); !ok {
			err = errutil.New("make id couldnt find id for", rtype)
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
