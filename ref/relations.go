package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	// "github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Relations maps ids to RefReleation.
// Compatible with unique.TypeRegistry.
type Relations struct {
	*Classes          // our own relation classes
	classes  *Classes // object classes
	objects  *Objects
	cache    RelationCache
}

func NewRelations(classes *Classes, objects *Objects) *Relations {
	return &Relations{
		NewClasses(),
		classes,
		objects,
		make(RelationCache),
	}
}

// RelationCache builds dynamically as relations are accessed.
type RelationCache map[string]*RefRelation

func (reg *Relations) GetRelation(name string) (ret rt.Relation, okay bool) {
	id := id.MakeId(name)
	if ref, ok := reg.cache[id]; ok {
		ret, okay = ref, true
	} else if cls, ok := reg.ClassMap[id]; ok {
		if ref, e := reg.newRelation(id, cls); e != nil {
			println(e.Error())
		} else {
			reg.cache[id] = ref
			ret, okay = ref, true
		}
	}
	return
}

func CountRelation(rtype r.Type) (one, many int, err error) {
OutOfLoop:
	for fw := unique.Fields(rtype); fw.HasNext(); {
		field := fw.GetNext()
		tag := unique.Tag(field.Tag)
		if rel, ok := tag.Find("rel"); ok {
			switch rel {
			case "one":
				one++
			case "many":
				many++
			default:
				err = errutil.New("unknown relation", rel)
				break OutOfLoop
			}
		}
	}
	return
}

// RegisterType compatible with unique.TypeRegistry
func (reg *Relations) RegisterType(rtype r.Type) (err error) {
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
			err = reg.Classes.RegisterType(rtype)
		}
	}
	return
}
