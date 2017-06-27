package scope

import (
	"github.com/ionous/iffy/rt"
)

type _ModelFinder struct {
	run rt.Runtime
}

// ModelFinder "searches" for objects by name, asking for them from the runtime without interpretation.
func ModelFinder(run rt.Runtime) rt.ObjectFinder {
	return _ModelFinder{run}
}

func (mf _ModelFinder) FindObject(name string) (rt.Object, bool) {
	return mf.run.GetObject(name)
}
