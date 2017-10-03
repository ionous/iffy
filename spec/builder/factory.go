package builder

import (
	"github.com/ionous/errutil"
)

type _Factory struct {
	specs  SpecFactory
	blocks *Mementos
}

func (b *_Factory) newBlock() (err error) {
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

func (b *_Factory) endBlock() (err error) {
	if _, ok := b.blocks.Pop(); !ok {
		err = errutil.New("mismatched begin/end")
	}
	return
}

func (b *_Factory) newCmd(src *Memento, name string, args []interface{}) (ret *Memento, err error) {
	if spec, e := b.specs.NewSpec(name); e != nil {
		err = e
	} else {
		ret, err = b.zip(&Memento{
			chain:   src,
			factory: b,
			spec:    spec,
			pos:     Capture(2),
		}, args)
	}
	return
}

func (b *_Factory) newCmds(src *Memento) (ret *Memento, err error) {
	if specs, e := b.specs.NewSpecs(); e != nil {
		err = e
	} else {
		ret, err = b.zip(&Memento{
			chain:   src,
			factory: b,
			specs:   specs,
			pos:     Capture(2),
		}, nil)
	}
	return
}

func (b *_Factory) newVal(src *Memento, val interface{}) (ret *Memento, err error) {
	if _, isMemento := val.(*Memento); isMemento {
		err = errutil.New("New value requested, but the value is not a primitive.")
	} else {
		ret, err = b.zip(&Memento{
			chain:   src,
			factory: b,
			val:     val,
			pos:     Capture(2),
		}, nil)
	}
	return
}

// Move the passed args to the targeted memento, then add the target to the current block.
// FIX: add a check against the most recent block to ensure it doesnt get pulled --
// that could happen if the user called .Block inside of a call to Cmd/s
func (b *_Factory) zip(dst *Memento, args []interface{}) (ret *Memento, err error) {
	if block, ok := b.blocks.Top(); !ok {
		err = errutil.New("no active block.")
	} else if !dst.kids.IsEmpty() {
		err = errutil.New("kids not empty")
	} else {
		if cnt := len(args); cnt > 0 {
			a := make([]*Memento, cnt)
			// Commands created as arguments will exist as mementos in this block's stack.
			// Go evaluates expression calls left to right, so the rightmost arg is on top.
			// Walk right-to-left, removing from the block and adding to the dst.
			for i := cnt - 1; i >= 0; i-- {
				arg := args[i]
				if n, ok := arg.(*Memento); !ok {
					// not a memento? then its a positional value.
					a[i] = &Memento{
						val: arg,
						pos: Capture(3),
					}
				} else if mostRecent, _ := block.kids.Pop(); n != mostRecent {
					err = errutil.New("unexpected argument")
					break
				} else if n.chain.chain != nil {
					// parameters always have a parent: the root node of the system
					// but a parent chain means we are in too deep.
					// we could, alternatively, run the chain out here.
					err = errutil.New("chaining calls while passing parameters is not permitted")
					break
				} else {
					a[i] = n
					n.chain = nil // clear for garbage collector
				}
			}
			dst.kids.list = a
		}
		if err == nil {
			ret = block.kids.Push(dst)
		}
	}
	return
}
