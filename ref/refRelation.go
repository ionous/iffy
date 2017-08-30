package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/rt"
)

// RefRelation describes a single relationship "archetype"
type RefRelation struct {
	Id    ident.Id
	Table *index.Table
}

// GetId returns the unique identifier for this types.
func (rel *RefRelation) GetId() ident.Id {
	return rel.Id
}

// GetType of the relation: one-to-one to many-to-many.
func (rel *RefRelation) GetType() index.Type {
	return rel.Table.Type()
}

func (rel *RefRelation) GetRelative(src, dst rt.Object) (ret interface{}, okay bool) {
	ret, okay = rel.Table.Data[index.Row{src.Id().Name, dst.Id().Name}]
	return
}

func (rel *RefRelation) GetTable() *index.Table {
	return rel.Table
}

// Relate defines a connection between two objects.
func (rel *RefRelation) Relate(src, dst rt.Object, onInsert index.OnInsert) (changed bool, err error) {
	if s, ok := reduce(src); !ok {
		err = errutil.Fmt("primary object is anonymous", src.Type())
	} else if d, ok := reduce(dst); !ok {
		err = errutil.Fmt("secondary object is anonymous", dst.Type())
	} else {
		changed, err = rel.Table.RelatePair(s, d, onInsert)
	}
	return
}

func reduce(obj rt.Object) (id string, okay bool) {
	if obj == nil {
		okay = true
	} else {
		id = obj.Id().Name
		okay = len(id) > 0
	}
	return
}
