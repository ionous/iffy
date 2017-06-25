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
	Builder             // for chaining
	pos     Location    // source of the memento
	key     string      // the target of this memento in its parent
	spec    spec.Spec   // cmd data
	specs   spec.Specs  // array data
	val     interface{} // primitive data
	kids    Mementos    // child data, either array elements or command parameters
}

// Mementos as a stack.
type Mementos struct {
	list []*Memento // args
}

// Top memento on the stack.
func (ns *Mementos) Top() (ret *Memento, okay bool) {
	if cnt := len(ns.list); cnt > 0 {
		ret, okay = ns.list[cnt-1], true
	}
	return
}

// IsEmpty if nothing is on the stack.
func (ns *Mementos) IsEmpty() bool {
	return len(ns.list) == 0
}

// Push and return the passed child.
func (ns *Mementos) Push(child *Memento) *Memento {
	ns.list = append(ns.list, child)
	return child
}

// Pop and return the popped Memento if successful
func (ns *Mementos) Pop() (ret *Memento, okay bool) {
	if cnt := len(ns.list); cnt > 0 {
		pop := cnt - 1
		ret, okay = ns.list[pop], true
		ns.list = ns.list[:pop]
	}
	return
}
