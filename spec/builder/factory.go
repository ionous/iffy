package builder

import (
	"github.com/ionous/iffy/spec"
)

type Factory struct {
	spec.SpecFactory
	scopes []*Block // stack of scopes
}

func (x *Factory) AddMemento(n *Memento) *Memento {
	curr := x.Current()
	curr.mementos = append(curr.mementos, n)
	return n
}

func (x *Factory) NewBlock(parent *Memento) *Block {
	if parent.val != nil {
		panic("can only use blocks with commands or arrays of commands")
	}
	m := &Block{parent: parent}
	x.scopes = append(x.scopes, m)
	return m
}

func (x *Factory) Current() *Block {
	return x.scopes[len(x.scopes)-1]
}

// pull removes the passed mementos from the current scope
// ns should contain a bunch of simple (not a block) commands, arrays of commands, and values
// because parameters are evaluated left to right the most recent is at the end
// panic if its not the most recent memento in the scope
func (x *Factory) Pull(ns []*Memento) {
	curr := x.Current()
	for i := len(ns) - 1; i >= 0; i-- {
		if mostRecent := curr.PopMemento(); ns[i] != mostRecent {
			panic("chained calls used? arguments should use standalone cmds, arrays, and values.")
		}
	}
}

func (x *Factory) EndBlock() bool {
	pop := len(x.scopes) - 1
	scope := x.scopes[pop]
	x.scopes = x.scopes[:pop]
	// take all the remaining elements (those that werent pull'd by into arguments)
	// and put them into spec of the memento
	if spec := scope.parent.spec; spec != nil {
		// println("finalizing a command block", len(scope.mementos))
		// the elements go into the spec
		var keys bool
		for _, n := range scope.mementos {
			arg := n.Interface()
			if k := n.key; len(k) > 0 {
				keys = true
				if e := spec.Assign(k, arg); e != nil {
					panic(e)
				}
			} else if !keys {
				if e := spec.Position(arg); e != nil {
					panic(e)
				}
			} else {
				panic("positional arguments cant follow key-value arguments")
			}
		}
	} else if specs := scope.parent.specs; specs != nil {
		// println("finalizing an array block", len(scope.mementos))
		// the elements go into the array
		for _, n := range scope.mementos {
			if n.spec == nil {
				panic("array element in not a spec?!")
			} else if e := specs.AddElement(n.spec); e != nil {
				panic(e)
			}
		}
	} else {
		panic("cant assign commands to a single value")
	}
	return pop == 0 // nothing left.
}
