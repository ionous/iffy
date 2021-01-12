package rel

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/ephemera/reader"
)

// Relation requires a user-specified string.
type Relation struct {
	At  reader.Position `if:"internal"`
	Str string
}

// String returns the name of a user defined relation.
func (op *Relation) String() string { return op.Str }

// Choices: relation names have no predetermined choices.
func (*Relation) Choices() map[string]string { return nil }

// Compose returns info for the modeling language.
func (*Relation) Compose() composer.Spec {
	return composer.Spec{
		Name:        "relation_name",
		OpenStrings: true,
	}
}
