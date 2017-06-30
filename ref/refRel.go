package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/unique"
	r "reflect"
)

type RelationRegistry struct {
	unique.Types
}

//for now -- where this is rel.Relate(temp)
// and probably everything is implemented via reflection
// someday you might add an implementation that uses sqlite
// or some orm.

// gremlinRocks.Relate{nil, rocky)
//     now rocky is owned by no one

// gremlinRocks.Relate{claire, nil)
//     now claire owns nothing.

// gremlinRocks.Relate{claire, rocky)
//     now claire owns rocky, and no one else does.

// gremlinRocks.Relate{claire, rocky, firstMet)
//     sets/updates the data for the relation.

// rel.Query() -> which returns an object stream, and do with it what you will regardless of speed, including poking in values.

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

func (r RelationRegistry) RegisterType(rtype r.Type) (err error) {
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
			if many == 2 {
				err = errutil.New("many-to-many relations not supported")
			} else {
				err = r.Types.RegisterType(rtype)
			}
		}
	}
	return
}
