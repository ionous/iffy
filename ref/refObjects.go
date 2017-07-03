package ref

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/rt"
)

type Objects map[string]*RefObject

func (m Objects) GetObject(name string) (ret rt.Object, okay bool) {
	id := id.MakeId(name)
	ret, okay = m[id]
	return
}
