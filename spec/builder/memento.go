package builder

import (
	"github.com/ionous/iffy/spec"
)

// Memento is returned by Builder. It contains a Builder to allow chaining of calls. Each chained call targets the surrounding block. For example, in:
//  if c.Cmd("parent").Block() {
//    c.Cmd("some command", params).Cmds(els).Val(value).End()
//  }
// the command, the array, and the val are all considered members of "parent".
type Memento struct {
	Builder
	key   string
	spec  spec.Spec  // cmd
	specs spec.Specs // array
	val   interface{}
}

// Interface returns the spec.Spec, spec.Specs, or primitive value specified via the Builder.
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
