package rtm

import (
	"github.com/ionous/iffy/ref"
)

type ModelFinder struct {
	model ref.Model
}

func (mf ModelFinder) FindObject(name string) (ref.Object, bool) {
	return mf.model.GetObject(name)
}
