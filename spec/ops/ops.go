package ops

import (
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/builder"
)

type Ops struct {
	unique.Types
	ShadowTypes *unique.Stack
}

// Builder is the ops version of builder.Builder.
type Builder struct {
	builder.Builder
}

func NewOps(classes unique.TypeRegistry) *Ops {
	// FIX: this uses the class registry slightly backwards
	// it knows to stack it rather than being given the types it should use
	// moreover it exposes types as an aggregate -- which shouldnt be needed;
	// isnt done by other similar things -- ex. ObjectGenerator.
	return &Ops{
		make(unique.Types),
		unique.NewStack(classes),
	}
}

// NewBuilder with a value transform.
func (ops *Ops) NewBuilder(root interface{}, x Transform) *Builder {
	return ops.NewFromTarget(InPlace(root), x)
}

func (ops *Ops) NewFromTarget(tgt Target, x Transform) *Builder {
	if x == nil {
		x = DefaultXform{}
	}
	c := &Command{xform: x, target: tgt}
	return &Builder{
		builder.NewBuilder(&Factory{ops, x}, c),
	}
}

// Build generates data into the root passed via NewBuilder()
func (u *Builder) Build(specs ...func(spec.Block)) (err error) {
	for _, s := range specs {
		s(u)
	}
	if _, e := u.Builder.Build(); e != nil {
		err = e
	}
	return
}
