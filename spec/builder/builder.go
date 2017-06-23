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
	internal.NewBlock(&Memento{
		Builder: b,
		spec:    spec,
	})
	return Builder{internal}
}

// Builder starts a new block of commands. Usually used as:
//  if c.Cmd("name").Block() {
//    c.End()
//  }
func (b Builder) Block() bool {
	b.internal.NewBlock(b.internal.Current().Memento())
	return true
}

// End terminates a block. See also Builder.Builder()
// Returns true at the very end.
func (b Builder) End() bool {
	return b.internal.EndBlock()
}

// Cmd creates a new command of name with the passed set of positional args.
// Returns a memento which can be passed to arrays or commands, or chained to add new data to the current block.
// To add data to *this* command, pass them via args or follow this immediately with Builder().
func (b Builder) Cmd(name string, args ...*Memento) (ret *Memento) {
	if spec, e := b.internal.NewSpec(name); e != nil {
		panic(e)
	} else {
		b.internal.Pull(args)
		for _, arg := range args {
			if e := spec.Position(arg.Interface()); e != nil {
				panic(e)
			}
		}
		ret = b.internal.AddMemento(&Memento{
			Builder: b,
			spec:    spec,
		})
	}
	return
}

// Array specifies a new array parameter.
func (b Builder) Array(vals ...*Memento) (ret *Memento) {
	if specs, e := b.internal.NewSpecs(); e != nil {
		panic(e)
	} else {
		b.internal.Pull(vals)
		ret = b.internal.AddMemento(&Memento{
			Builder: b,
			specs:   specs,
		})
	}
	return
}

// Val specifies a single literal: whether one primitive value or one array of primitive values. It does not start a new block, because primitive values have no additional parameters.
func (b Builder) Val(val interface{}) *Memento {
	return b.internal.AddMemento(&Memento{
		Builder: b,
		val:     val,
	})
}

// Param adds a key-value parameter to the spec.
// The passed name is the key; the return value is used to specify the value.
func (b Builder) Param(field string) Param {
	m := b.internal.Current()

	if m.parent.specs != nil {
		panic("arrays cant have named parameters")
	}

	return Param{b, field}
}
