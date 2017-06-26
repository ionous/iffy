package scope

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
)

type _MultiFinder []rt.ObjectFinder

func MultiFinder(finders ...rt.ObjectFinder) rt.ObjectFinder {
	return _MultiFinder(finders)
}

func (mf _MultiFinder) FindObject(name string) (ret ref.Object, okay bool) {
	for _, f := range mf {
		if obj, ok := f.FindObject(name); ok {
			ret, okay = obj, true
			break
		}
	}
	return
}
