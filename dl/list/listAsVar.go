package list

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

// ListIterator defines a variable name
type ListIterator interface {
	// ? maybe we could be the thing that genates the record?
	Name() string
	Affinity() affine.Affinity
}

type AsNum struct {
	Var core.Variable `if:"selector"`
}
type AsTxt struct {
	Var core.Variable `if:"selector"`
}
type AsRec struct {
	Var core.Variable `if:"selector"`
}

func (*AsNum) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Define the name of a number variable.",
	}
}

func (*AsTxt) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Define the name of a text variable.",
	}
}

func (*AsRec) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Selector},
		Desc:   "Define the name of a record variable.",
	}
}

func (op *AsNum) Name() string {
	return op.Var.String()
}

func (op *AsRec) Name() string {
	return op.Var.String()
}

func (op *AsTxt) Name() string {
	return op.Var.String()
}

func (op *AsNum) Affinity() affine.Affinity {
	return affine.Number
}

func (op *AsRec) Affinity() affine.Affinity {
	return affine.Record
}

func (op *AsTxt) Affinity() affine.Affinity {
	return affine.Text
}
