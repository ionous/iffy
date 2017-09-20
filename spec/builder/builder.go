package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

type Builder struct {
	*Memento
}

type SpecFactory interface {
	NewSpec(name string) (spec.Spec, error)
	NewSpecs() (spec.Specs, error)
}

func NewBuilder(sf SpecFactory, spec spec.Spec) Builder {
	factory := &_Factory{sf, new(Mementos)}
	return Builder{factory.blocks.Push(&Memento{
		factory: factory,
		spec:    spec,
		pos:     Capture(1),
	})}
}

func NewBuilderPos(sf SpecFactory, spec spec.Spec, i int) Builder {
	factory := &_Factory{sf, new(Mementos)}
	return Builder{factory.blocks.Push(&Memento{
		factory: factory,
		spec:    spec,
		pos:     Capture(i),
	})}
}

// Build computes a final result.
func (b Builder) Build() (ret interface{}, err error) {
	f := b.factory
	if root, ok := f.blocks.Pop(); !ok {
		err = errutil.New("nothing to build, possibly too many ends")
	} else if !f.blocks.IsEmpty() {
		err = errutil.New("not all blocks have ended", f.blocks.list[0].pos, len(f.blocks.list), "remain")
	} else if res, e := Build(root); e != nil {
		err = e
	} else {
		ret = res
	}
	return
}
