package ref

import (
	"github.com/ionous/errutil"
	r "reflect"
)

type RelBuilder struct {
	tables, objectClasses *Classes
}

func (b *RelBuilder) Build(objects *Objects) *Relations {
	return &Relations{
		b.tables,
		b.objectClasses,
		objects,
		make(RelationCache),
	}
}

func NewRelations(objectClasses *Classes) *RelBuilder {
	return &RelBuilder{
		NewClasses(),
		objectClasses,
	}
}

// Compatible with unique.TypeRegistry.
func (b *RelBuilder) FindType(name string) (r.Type, bool) {
	return b.tables.FindType(name)
}

// RegisterType compatible with unique.TypeRegistry
func (b *RelBuilder) RegisterType(rtype r.Type) (err error) {
	// filter then:
	if one, many, e := CountRelation(rtype); e != nil {
		err = e
	} else if err == nil {
		switch cnt := one + many; {
		case cnt < 2:
			err = errutil.New("too few relations specified")
		case cnt > 2:
			err = errutil.New("too many relations specified")
		default:
			err = b.tables.RegisterType(rtype)
		}
	}
	return
}
