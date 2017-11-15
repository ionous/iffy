package obj

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Registry with ids, findable by the game.
type Registry struct {
	classes   unique.Types   // set of allowed types
	counters  ident.Counters // id generation for unnamed objets
	uniqueIds []uniqueId     // NewName queue
	values    []interface{}  // RegisterValue queue
}

// record of requested class and name we generatd for it.
type uniqueId struct {
	clsName string
	objName string
}

func (b *Registry) addName(cls string) string {
	if b.counters == nil {
		b.counters = make(ident.Counters)
	}
	return b.counters.NewName(cls)
}

// CreateName given the passed class, generate a unique name, which will generate a new blank object at Build().
func (b *Registry) NewName(cls string) string {
	name := b.addName(cls)
	b.uniqueIds = append(b.uniqueIds, uniqueId{cls, name})
	return name
}

// RegisterValue the passed value as a future object.
func (b *Registry) RegisterValues(vs []interface{}) {
	b.values = append(b.values, vs...)
}

// CreateObjects returns an ObjectMap
func (b *Registry) CreateObjects(run rt.Runtime, classes unique.Types) (ret ObjectMap, err error) {
	om := make(ObjectMap)
	if e := b.createFromIds(run, classes, om); e != nil {
		err = e
	} else if e := b.createFromValues(run, om); e != nil {
		err = e
	} else {
		ret = om
	}
	return
}

func (b *Registry) createFromValues(run rt.Runtime, om ObjectMap) (err error) {
	for i, v := range b.values {
		if rval, e := unique.ValuePtr(v); e != nil {
			err = errutil.New("couldnt get pointer for value", i, e)
			break
		} else if id := b.makeId(rval); !id.IsValid() {
			err = errutil.New("couldnt create id for value", i, rval)
			break
		} else if was, ok := om.addObject(run, id, rval); !ok {
			err = errutil.New("couldnt create id for value", i, rval, "was", id, was)
			break
		}
	}
	return
}

func (b *Registry) createFromIds(run rt.Runtime, classes unique.Types, om ObjectMap) (err error) {
	for _, u := range b.uniqueIds {
		if cls, ok := classes.FindType(u.clsName); !ok {
			err = errutil.New("unknown class for", u)
			break
		} else if path, ok := unique.PathOf(cls, "id"); !ok {
			err = errutil.New("no id path for", cls)
			break
		} else {
			n := r.New(cls).Elem()
			id := ident.IdOf(u.objName)
			if was, ok := om.addObject(run, id, n); !ok {
				err = errutil.New("duplicate id", u, was)
				break
			} else {
				field := n.FieldByIndex(path)
				field.SetString(u.objName)
			}
		}
	}
	return
}

func (om ObjectMap) addObject(run rt.Runtime, id ident.Id, rval r.Value) (ret RefObject, okay bool) {
	if prev, exists := om[id]; exists {
		ret = prev
	} else {
		ret = RefObject{id, rval, run}
		om[id] = ret
		okay = true
	}
	return
}

func (b *Registry) makeId(rval r.Value) (ret ident.Id) {
	rtype := rval.Type()
	if path, ok := unique.PathOf(rtype, "id"); ok {
		field := rval.FieldByIndex(path)
		name := field.String()
		// handle unnamed objects:
		if len(name) == 0 {
			name = b.addName(rtype.Name())
			field.SetString(name)
		}
		ret = ident.IdOf(name)
	}
	return
}
