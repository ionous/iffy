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
	return &Ops{
		make(unique.Types),
		unique.NewStack(classes),
	}
}

// NewBuilder starts creating a call tree. Always returns true.
func (ops *Ops) NewBuilder(root interface{}) (*Builder, bool) {
	return ops.NewXBuilder(root, DefaultXform{})
}

// NewBuilder with a value transform.
func (ops *Ops) NewXBuilder(root interface{}, x Transform) (*Builder, bool) {
	c := &Command{xform: x, target: InPlace(root)}
	return &Builder{
		builder.NewBuilder(&Factory{ops, x}, c),
	}, true
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
