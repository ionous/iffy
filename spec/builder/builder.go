package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

// Builder builds commands.
type Builder struct {
	specs  spec.Factory
	blocks *Mementos
}

func NewBuilder(sf spec.Factory, spec spec.Spec) Builder {
	b := Builder{sf, new(Mementos)}
	b.blocks.Push(&Memento{
		Builder: b,
		spec:    spec,
		pos:     Capture(1),
	})
	return b
}

// Builder starts a new block of commands. Usually used as:
//  if c.Cmd("name").Block() {
//    c.End()
//  }
func (b Builder) Block() (okay bool) {
	if e := b.newBlock(); e != nil {
		panic(e)
	} else {
		okay = true
	}
	return
}

// Build computes a final result.
func (b Builder) Build() (ret interface{}, err error) {
	if root, ok := b.blocks.Pop(); !ok {
		err = errutil.New("nothing to build")
	} else if !b.blocks.IsEmpty() {
		err = errutil.New("not all blocks have ended", len(b.blocks.list), "remain")
	} else if res, e := Build(root); e != nil {
		err = e
	} else {
		ret = res
	}
	return
}

// End terminates a block. See also Builder.Builder()
func (b Builder) End() {
	if e := b.endBlock(); e != nil {
		panic(e)
	}
	return
}

// Cmd adds a new command of name with the passed set of positional args. Args can contain Mementos and literals. Returns a memento which can be passed to arrays or commands, or chained.
// To add data to the new command, pass them via args or follow this call with a call to Builder.Block().
func (b Builder) Cmd(name string, args ...interface{}) (ret *Memento) {
	if n, e := b.newCmd(name, args); e != nil {
		panic(e)
	} else {
		ret = n
	}
	return
}

// Cmds specifies a new array of commands. Additional elements can be added to the array using Builder.Block().
func (b Builder) Cmds(cmds ...*Memento) (ret *Memento) {
	if n, e := b.newCmds(cmds); e != nil {
		panic(e)
	} else {
		ret = n
	}
	return
}

// Val specifies a single literal value: whether one primitive value or one array of primitive values.
func (b Builder) Val(val interface{}) (ret *Memento) {
	if n, e := b.newVal(val); e != nil {
		panic(e)
	} else {
		ret = n
	}
	return
}

// Param adds a key-value parameter to the spec.
// The passed name is the key; the return value is used to specify the value.
func (b Builder) Param(field string) Param {
	return Param{b, field}
}

func (b *Builder) newBlock() (err error) {
	// another way of thinking about new/Block() is that it elevates the most recent^ memento to block status. parameters are excluded from most recent because they are pulled from the block
	if block, ok := b.blocks.Top(); !ok {
		err = errutil.New("block can only be used inside of another block.")
	} else if mostRecent, ok := block.kids.Top(); !ok {
		err = errutil.New("block encountered an unexpected error")
	} else {
		b.blocks.Push(mostRecent)
	}
	return
}

func (b *Builder) endBlock() (err error) {
	if _, ok := b.blocks.Pop(); !ok {
		err = errutil.New("mismatched begin/end")
	}
	return
}

func (b *Builder) newCmd(name string, args []interface{}) (ret *Memento, err error) {
	if spec, e := b.specs.NewSpec(name); e != nil {
		err = e
	} else {
		if n, e := b.zip(&Memento{
			Builder: *b,
			spec:    spec,
			pos:     Capture(2),
		}, args); e != nil {
			err = e
		} else {
			ret = n
		}
	}
	return
}

func (b *Builder) newCmds(cmds []*Memento) (ret *Memento, err error) {
	if specs, e := b.specs.NewSpecs(); e != nil {
		err = e
	} else {
		// normalize into an array of interfaces :(
		args := make([]interface{}, len(cmds))
		for i, c := range cmds {
			args[i] = c
		}
		if n, e := b.zip(&Memento{
			Builder: *b,
			specs:   specs,
			pos:     Capture(2),
		}, args); e != nil {
			err = e
		} else {
			ret = n
		}
	}
	return
}

func (b *Builder) newVal(val interface{}) (ret *Memento, err error) {
	if _, isMemento := val.(*Memento); isMemento {
		err = errutil.New("New value requested, but the value is not a primitive.")
	} else {
		if n, e := b.zip(&Memento{
			Builder: *b,
			val:     val,
			pos:     Capture(2),
		}, nil); e != nil {
			err = e
		} else {
			ret = n
		}
	}
	return
}

// Move the passed args to the targeted memento, then add the target to the current block.
// FIX: add a check against the most recent block to ensure it doesnt get pulled --
// that could happen if the user called .Block inside of a call to Cmd/s
func (b *Builder) zip(dst *Memento, args []interface{}) (ret *Memento, err error) {
	if block, ok := b.blocks.Top(); !ok {
		err = errutil.New("no active block.")
	} else {
		// Commands created as arguments will exist as mementos in this block's stack.
		// Go evaluates expression calls left to right, so the rightmost arg is on top.
		// Walk right-to-left, removing from the block and adding to the dst.
		for i := len(args); i > 0; i-- {
			arg := args[i-1]
			if n, ok := arg.(*Memento); !ok {
				// not a memento? then its a positional value.
				dst.kids.Push(&Memento{
					val: arg,
					pos: Capture(3),
				})
			} else if mostRecent, _ := block.kids.Pop(); n == mostRecent {
				dst.kids.Push(n)
			} else {
				err = errutil.New("unexpected argument")
				break
			}
		}
		if err == nil {
			ret = block.kids.Push(dst)
		}
	}
	return
}
