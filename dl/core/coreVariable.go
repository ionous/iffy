package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
)

// Variable requires a user-specified string.
type Variable struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *Variable) String() string { return op.Str }

func (*Variable) Choices() []string { return nil }

func (*Variable) Compose() composer.Spec {
	return composer.Spec{
		Name:        "variable_name",
		OpenStrings: true,
	}
}
