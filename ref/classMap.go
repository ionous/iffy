package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Classes maps ids to rt.Class.
// Compatible with unique.TypeRegistry
type ClassMap map[string]rt.Class

// GetClass compatible with rt.Runtime
func (cm ClassMap) GetClass(name string) (ret rt.Class, okay bool) {
	id := id.MakeId(name)
	ret, okay = cm[id]
	return
}

// GetByType for cache usage
func (cm ClassMap) GetByType(rtype r.Type) (ret rt.Class, err error) {
	name := rtype.Name()
	id := id.MakeId(name)
	if cls, ok := cm[id]; !ok {
		err = errutil.New("class not found", name)
	} else if cls != rtype {
		err = errutil.New("class conflict", name, cls, rtype)
	} else {
		ret = cls
	}
	return
}
