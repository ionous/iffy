package unique

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	r "reflect"
	"strconv"
)

// Objects (unique.Objects) generates a set of unique instances.
type Objects struct {
	ids map[r.Type]int
	Types
}

func NewObjectGenerator() *Objects {
	return &Objects{
		make(map[r.Type]int),
		make(Types),
	}
}

func makeId(scope string, i int) ident.Id {
	return ident.IdOf(scope + "#" + strconv.Itoa(i))
}

func (u *Objects) Id(cls string) (ret string) {
	if rtype, ok := u.Types.FindType(cls); !ok {
		e := errutil.New("unknown class", cls)
		panic(e)
	} else {
		cnt := u.ids[rtype] + 1
		u.ids[rtype] = cnt
		ret = makeId(rtype.Name(), cnt).Name
	}
	return
}

func (u *Objects) Generate() (ret []interface{}, err error) {
	for rtype, cnt := range u.ids {
		if idpath, ok := PathOf(rtype, "id"); !ok {
			err = errutil.New("couldn't find id field in", rtype)
			break
		} else {
			scope := rtype.Name()
			for i := 0; i < cnt; i++ {
				n := r.New(rtype)
				idea := n.Elem().FieldByIndex(idpath)
				id := makeId(scope, i+1) // ids are 1 based
				idea.SetString(id.Name)  // the id path points to the string field which generates an id
				ret = append(ret, n.Interface())
			}
		}
	}

	return
}
