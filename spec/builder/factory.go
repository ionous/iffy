package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

type Factory struct {
	spec.SpecFactory
	blocks []*Block // stack of blocks
}

// Current block. Returns false if no block exists.
func (x *Factory) Current() (ret *Block, okay bool) {
	if cnt := len(x.blocks); cnt > 0 {
		ret, okay = x.blocks[cnt-1], true
	}
	return
}

// NewBlock pushes a new block onto the stack. All operations in the block will target the passed memento.
func (x *Factory) NewBlock(parent *Memento) (ret *Block, err error) {
	if parent == nil {
		err = errutil.New("NewBlock requires a valid target memento.")
	} else if parent.val != nil {
		err = errutil.New("NewBlock works only with commands or arrays of commands.")
	} else {
		block := &Block{parent: parent}
		x.blocks = append(x.blocks, block)
		ret = block
	}
	return
}

// AddMemento to the blockent block. Returns the passed memento, or err if there is no active block.
func (x *Factory) AddMemento(n *Memento) (ret *Memento, err error) {
	if block, ok := x.Current(); !ok {
		err = errutil.New("AddMemento called, but no active block.")
	} else {
		block.mementos = append(block.mementos, n)
		ret = n
	}
	return
}

// Position the passed args in the passed spec, removing mementos from the stack so that they do not get added to the parent spec.
func (x *Factory) Position(spec spec.Spec, args []interface{}) (err error) {
	if block, ok := x.Current(); !ok {
		err = errutil.New("Position called, but no active block.")
	} else {
		// decode each argument. because function calls in a go expression are evaluated left to right, any commands created as arguements inside of a call to Builder.Cmd/s(), will exist as mementos on the blockent stack. The right most argument will be on top.
		for i := len(args) - 1; i >= 0; i-- {
			arg := args[i]
			var val interface{}
			if n, ok := arg.(*Memento); !ok {
				val = arg
			} else if mostRecent := block.PopMemento(); n == mostRecent {
				val = n.Interface()
			} else {
				err = errutil.New("chained calls used? arguments should use standalone cmds, arrays, and values.")
				break
			}
			if e := spec.Position(val); e != nil {
				err = e
				break
			}
		}
	}
	return
}
func (x *Factory) AddElements(specs spec.Specs, ns []*Memento) (err error) {
	if block, ok := x.Current(); !ok {
		err = errutil.New("AddElements called, but no active block.")
	} else {
		for i := len(ns) - 1; i >= 0; i-- {
			if mostRecent := block.PopMemento(); ns[i] != mostRecent {
				err = errutil.New("AddElements detected method chaining? Use standalone command calls for array parameters.")
				break
			} else if spec := mostRecent.spec; spec == nil {
				err = errutil.New("AddElements requires commands, not values or other arrays.")
				break
			} else if e := specs.AddElement(spec); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// EndBlock finalizes and pops the most recent NewBlock.
// Returns true if this was the final block.
func (x *Factory) EndBlock() (ret bool, err error) {
	if cnt := len(x.blocks); cnt == 0 {
		err = errutil.New("EndBlock called, but no active block")
	} else {
		pop := cnt - 1
		if e := x.flush(x.blocks[pop]); e != nil {
			err = e
		} else {
			// if no error, pop the block
			x.blocks = x.blocks[:pop]
			ret = pop == 0 // true if nothing left.
		}
	}
	return
}

func (x *Factory) flush(block *Block) (err error) {
	// take all the remaining elements (those that werent pull'd by into arguments)
	// and put them into spec of the memento
	if spec := block.parent.spec; spec != nil {
		// println("finalizing a command block", len(block.mementos))
		// the elements go into the spec
		var keys bool
		for _, n := range block.mementos {
			arg := n.Interface()
			if k := n.key; len(k) > 0 {
				keys = true
				if e := spec.Assign(k, arg); e != nil {
					err = e
					break
				}
			} else if !keys {
				if e := spec.Position(arg); e != nil {
					err = e
					break
				}
			} else {
				err = errutil.New("positional arguments cant follow key-value arguments")
				break
			}
		}
	} else if specs := block.parent.specs; specs != nil {
		// println("finalizing an array block", len(block.mementos))
		// the elements go into the array
		for _, n := range block.mementos {
			if n.spec == nil {
				err = errutil.New("array element in not a spec?!")
			} else if e := specs.AddElement(n.spec); e != nil {
				err = e
			}
		}
	} else {
		err = errutil.New("cant assign commands to a single value")
	}
	return
}
