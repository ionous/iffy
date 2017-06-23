package builder

import (
	"github.com/ionous/iffy/spec"
)

// Builder builds commands.
type Builder struct {
	internal *Factory
}

func NewBuilder(sf spec.SpecFactory, spec spec.Spec) Builder {
	internal := &Factory{SpecFactory: sf}
	b := Builder{internal}
	if _, e := internal.NewBlock(&Memento{
		Builder: b,
		spec:    spec,
	}); e != nil {
		panic(e)
	}
	return Builder{internal}
}

// Builder starts a new block of commands. Usually used as:
//  if c.Cmd("name").Block() {
//    c.End()
//  }
func (b Builder) Block() (ret bool) {
	if block, ok := b.internal.Current(); !ok {
		panic("Block can only be used inside of another block.")
	} else if _, e := b.internal.NewBlock(block.Memento()); e != nil {
		panic(e)
	} else {
		ret = true
	}
	return
}

// End terminates a block. See also Builder.Builder()
// Returns true at the very end.
func (b Builder) End() (ret bool) {
	if r, e := b.internal.EndBlock(); e != nil {
		panic(e)
	} else {
		ret = r
	}
	return
}

// Cmd adds a new command of name with the passed set of positional args. Args can contain Mementos and literals. Returns a memento which can be passed to arrays or commands, or chained.
// To add data to the new command, pass them via args or follow this call with a call to Builder.Block().
func (b Builder) Cmd(name string, args ...interface{}) (ret *Memento) {
	if spec, e := b.internal.NewSpec(name); e != nil {
		panic(e)
	} else if e := b.internal.Position(spec, args); e != nil {
		panic(e)
	} else if n, e := b.internal.AddMemento(&Memento{
		Builder: b,
		spec:    spec,
	}); e != nil {
		panic(e)
	} else {
		ret = n
	}
	return
}

// Cmds specifies a new array of commands. Additional elements can be added to the array using Builder.Block().
func (b Builder) Cmds(cmds ...*Memento) (ret *Memento) {
	if specs, e := b.internal.NewSpecs(); e != nil {
		panic(e)
	} else if e := b.internal.AddElements(specs, cmds); e != nil {
		panic(e)
	} else if n, e := b.internal.AddMemento(&Memento{
		Builder: b,
		specs:   specs,
	}); e != nil {
		panic(e)
	} else {
		ret = n
	}
	return
}

// Val specifies a single literal value: whether one primitive value or one array of primitive values.
func (b Builder) Val(val interface{}) (ret *Memento) {
	if _, isMemento := val.(*Memento); isMemento {
		panic("Only primitive values should be passed to Builder.Val")
	} else if n, e := b.internal.AddMemento(&Memento{
		Builder: b,
		val:     val,
	}); e != nil {
		panic(e)
	} else {
		ret = n
	}
	return
}

// Param adds a key-value parameter to the spec.
// The passed name is the key; the return value is used to specify the value.
func (b Builder) Param(field string) Param {
	// verify that it makes sense to call param.
	if block, ok := b.internal.Current(); !ok {
		panic("Param can only be used inside of a valid block.")
	} else if block.parent.spec == nil {
		panic("Only commands can have named parameters.")
	}
	return Param{b, field}
}
