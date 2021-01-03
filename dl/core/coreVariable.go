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

// String returns user defined variable name.
func (op *Variable) String() string { return op.Str }

// Choices: variable names have no predetermined choices.
func (*Variable) Choices() map[string]string { return nil }

// Compose returns info for the modeling language.
func (*Variable) Compose() composer.Spec {
	return composer.Spec{
		Name:        "variable_name",
		OpenStrings: true,
	}
}
