package obj

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
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
	cls  r.Type
}
type queue map[ident.Id]queued

func NewObjects() *ObjBuilder {
	return &ObjBuilder{make(queue), make(infoMap)}
}

// Build returns an ObjectMap
func (b *ObjBuilder) Build(p rt.Runtime) ObjectMap {
	objs := make(ObjectMap)
	for id, q := range b.queue {
		objs[id] = RefObject{id, q.rval, p}
	}
	return objs
}

// RegisterValue the passed value as a future object.
// Compatible with unique.ValueRegistry
func (b *ObjBuilder) RegisterValue(rval r.Value) (err error) {
	if id, e := b.makeId(rval); e != nil {
		err = e
	} else if obj, ok := b.queue[id]; ok {
		err = errutil.New("duplicate object", id, obj, rval)
	} else {
		b.queue[id] = queued{rval, rval.Type()}
	}
	return
}

func (b *ObjBuilder) makeId(rval r.Value) (ret ident.Id, err error) {
	rtype := rval.Type()

	// cache the id path for this type
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
		ret = ident.IdOf(name)
	}
	return
}
