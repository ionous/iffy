package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

type Builder struct {
	*Memento
}

func NewBuilder(sf spec.Factory, spec spec.Spec) Builder {
	factory := &Factory{sf, new(Mementos)}
	return Builder{factory.blocks.Push(&Memento{
		factory: factory,
		spec:    spec,
		pos:     Capture(1),
	})}
}

// Build computes a final result.
func (b Builder) Build() (ret interface{}, err error) {
	f := b.factory
	if root, ok := f.blocks.Pop(); !ok {
		err = errutil.New("nothing to build")
	} else if !f.blocks.IsEmpty() {
		err = errutil.New("not all blocks have ended", len(f.blocks.list), "remain")
	} else if res, e := Build(root); e != nil {
		err = e
	} else {
		ret = res
	}
	return
}
