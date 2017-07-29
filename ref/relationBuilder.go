package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/index"
)

type RelBuilder struct {
	types map[string]index.Type
}

func (b *RelBuilder) Build() Relations {
	r := make(Relations)
	for id, t := range b.types {
		r[id] = &RefRelation{id, index.NewTable(t)}
	}
	return r
}

func NewRelations() *RelBuilder {
	return &RelBuilder{make(map[string]index.Type)}
}

// RegisterType compatible with unique.TypeRegistry
func (b *RelBuilder) NewRelation(name string, kind index.Type) (err error) {
	id := id.MakeId(name)
	if k, ok := b.types[id]; !ok {
		b.types[id] = kind
	} else if k != kind {
		err = errutil.New("mismatched relations", k, kind)
	}
	return
}
