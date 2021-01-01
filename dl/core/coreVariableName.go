package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
)

// VariableName requires a user-specified string.
type VariableName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *VariableName) String() string {
	return op.Str
}

func (*VariableName) Choices() (closed bool, choices map[string]string) {
	closed = false
	return
}

func (*VariableName) Compose() composer.Spec {
	return composer.Spec{
		Name: "variable_name",
		Uses: "str",
		Spec: "{variable_name}",
	}
}
