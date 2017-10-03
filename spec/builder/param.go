package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

// Param targets a key-value spec argument. It implements spec.Slot.
type Param struct {
	src *Memento
	key string
}

// Cmd creates a new command of the passed name for the parameter mentioned by Memento.Param(). Args can contain Mementos and literals.
func (p Param) Cmd(name string, args ...interface{}) (ret spec.Block) {
	if n, e := p.src.factory.newCmd(p.src, name, args); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		n.key = p.key
		ret = n
	}
	return
}

// Cmds creates a new array of commands for the parameter mentioned by Memento.Param().
func (p Param) Cmds() (ret spec.Block) {
	if n, e := p.src.factory.newCmds(p.src); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		n.key = p.key
		ret = n
	}
	return
}

// Val specifies a single literal value for the parameter mentioned by Memento.Param().
func (p Param) Val(val interface{}) (ret spec.Block) {
	if n, e := p.src.factory.newVal(p.src, val); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		n.key = p.key
		ret = n
	}
	return
}
