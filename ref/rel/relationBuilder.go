package rel

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/index"
)

type RelationBuilder map[ident.Id]*index.Table

func (b RelationBuilder) Build() Relations {
	r := make(Relations)
	for id, t := range b {
		r[id] = &RefRelation{id, t}
	}
	return r
}

func NewRelations() RelationBuilder {
	return make(RelationBuilder)
}

func (b RelationBuilder) NewRelation(name string, kind index.Type) (err error) {
	id := ident.IdOf(name)
	if t, exists := b[id]; !exists {
		b[id] = index.NewTable(kind)
	} else if k := t.Type(); k != kind {
		err = errutil.New("mismatched relations", k, kind)
	}
	return
}

func (b RelationBuilder) AddTable(name string, t *index.Table) (err error) {
	id := ident.IdOf(name)
	if _, exists := b[id]; !exists {
		b[id] = t
	} else {
		err = errutil.New("table already exists", name, id)
	}
	return
}
