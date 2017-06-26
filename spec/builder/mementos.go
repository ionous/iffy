package builder

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
