package rt

import (
	"github.com/ionous/iffy/ref"
)

type ObjectFinder interface {
	FindObject(name string) (ref.Object, bool)
}

type Runtime interface {
	ref.Model
	ObjectFinder
}
