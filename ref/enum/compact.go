package enum

import (
	"github.com/ionous/errutil"
	r "reflect"
)

type _Stringer interface {
	String() string
}

func Compact(rtype r.Type, name string, index []uint8) (ret []string, err error) {
	if cnt := len(index); cnt == 0 {
		err = errutil.New("empty stringer index", name)
	} else if last := index[cnt-1]; len(name) != int(last) {
		err = errutil.New("mismatched stringer index", name)
	} else {
		for i := 0; i < cnt-1; i++ {
			n := name[index[i]:index[i+1]]
			ret = append(ret, n)
		}
	}
	return
}
