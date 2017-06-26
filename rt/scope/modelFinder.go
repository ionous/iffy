package scope

import (
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
)

type _ModelFinder struct {
	model ref.Model
}

func ModelFinder(model ref.Model) rt.ObjectFinder {
	return _ModelFinder{model}
}

func (mf _ModelFinder) FindObject(name string) (ref.Object, bool) {
	return mf.model.GetObject(name)
}
