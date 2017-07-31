package unique

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
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

func makeId(scope string, i int) string {
	return id.MakeId(scope + "#" + strconv.Itoa(i))
}

func (u *Objects) Id(cls string) (ret string) {
	if rtype, ok := u.Types.FindType(cls); !ok {
		e := errutil.New("unknown class", cls)
		panic(e)
	} else {
		cnt := u.ids[rtype] + 1
		u.ids[rtype] = cnt
		ret = makeId(rtype.Name(), cnt)
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
				idea.SetString(makeId(scope, i+1)) //ids are 1 based
				ret = append(ret, n.Interface())
			}
		}
	}

	return
}
