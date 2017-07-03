package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	// "github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type Relations struct {
	Classes               // our own relation classes
	objectClasses Classes // the classes used by our relations
	cache         RelationCache
}
type RelationCache map[string]*RefRelation

func MakeRelations(objectClasses Classes) Relations {
	return Relations{make(Classes), objectClasses, make(RelationCache)}
}

func (r *Relations) GetRelation(name string) (ret rt.Relation, okay bool) {
	id := id.MakeId(name)
	if ref, ok := r.cache[id]; ok {
		ret, okay = ref, true
	} else if cls, ok := r.Classes[id]; ok {
		if ref, e := NewRelation(id, cls, r.objectClasses); e != nil {
			println(e.Error())
		} else {
			r.cache[id] = ref
			ret, okay = ref, true
		}
	}
	return
}

func CountRelation(rtype r.Type) (one, many int, err error) {
OutOfLoop:
	for fw := unique.Fields(rtype.Elem()); fw.HasNext(); {
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

func (r Relations) RegisterType(rtype r.Type) (err error) {
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
			err = r.Classes.RegisterType(rtype)
		}
	}
	return
}
