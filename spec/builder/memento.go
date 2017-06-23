package builder

import (
	"github.com/ionous/iffy/spec"
)

type Memento struct {
	Builder
	key   string
	spec  spec.Spec  // cmd
	specs spec.Specs // array
	val   interface{}
}

func (n *Memento) Interface() (ret interface{}) {
	if n.spec != nil {
		ret = n.spec
	} else if n.specs != nil {
		ret = n.specs
	} else {
		ret = n.val
	}
	return
}
