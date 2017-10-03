package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

// Memento implements spec.Block. Each chained call targets the surrounding block. For example, in:
//  if c.Cmd("parent").Begin() {
//    c.Cmd("some command", params).Cmds().Val(value).End()
//  }
// "the command, the array, and the val are all considered members of "parent".
type Memento struct {
	factory *_Factory   // for chaining
	chain   *Memento    // for detecting bad chaining
	pos     Location    // source of the memento
	key     string      // the target of this memento in its parent
	spec    spec.Spec   // cmd data
	specs   spec.Specs  // array data
	val     interface{} // primitive data
	kids    Mementos    // child data, either array elements or command parameters
}

// Begin starts a new parameter block. Usually used as:
//  if c.Cmd("name").Begin() {
//    c.End()
//  }
func (n *Memento) Begin() (okay bool) {
	if e := n.factory.newBlock(); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		okay = true
	}
	return
}

// End terminates a block.
func (n *Memento) End() {
	if e := n.factory.endBlock(); e != nil {
		panic(errutil.New(e, Capture(1)))
	}
	return
}

// Cmd adds a new command of name with the passed set of positional args. Args can contain Mementos and literals. Returns a memento which can be passed to arrays or commands, or chained.
// To add data to the new command, pass them via args or follow this call with a call to Begin().
func (n *Memento) Cmd(name string, args ...interface{}) (ret spec.Block) {
	if next, e := n.factory.newCmd(n, name, args); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		ret = next
	}
	return
}

// Cmds specifies a new array of commands. Additional elements can be added to the array using Begin().
func (n *Memento) Cmds() (ret spec.Block) {
	if n, e := n.factory.newCmds(n); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		ret = n
	}
	return
}

// Val specifies a single literal value: whether one primitive value or one array of primitive values.
func (n *Memento) Val(val interface{}) (ret spec.Block) {
	if n, e := n.factory.newVal(n, val); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		ret = n
	}
	return
}

// Param adds a key-value parameter to the spec.
// The passed name is the key; the return value is used to specify the value.
func (n *Memento) Param(field string) spec.Slot {
	return Param{n, field}
}
