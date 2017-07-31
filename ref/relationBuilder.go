package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/index"
)

type RelationBuilder map[string]index.Type

func (b RelationBuilder) Build() Relations {
	r := make(Relations)
	for id, t := range b {
		r[id] = &RefRelation{id, index.NewTable(t)}
	}
	return r
}

func NewRelations() RelationBuilder {
	return make(map[string]index.Type)
}

type RelationDesc struct {
	Name string
	Type index.Type
}

func RegisterRelations(b RelationBuilder, desc ...RelationDesc) (err error) {
	for _, d := range desc {
		if e := b.NewRelation(d.Name, d.Type); e != nil {
			err = e
			break
		}
	}
	return
}

// RegisterType compatible with unique.TypeRegistry
func (b RelationBuilder) NewRelation(name string, kind index.Type) (err error) {
	id := id.MakeId(name)
	if k, ok := b[id]; !ok {
		b[id] = kind
	} else if k != kind {
		err = errutil.New("mismatched relations", k, kind)
	}
	return
}
