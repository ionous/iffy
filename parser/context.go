package parser

import (
	"github.com/ionous/iffy/ident"
)

type Context interface {
	Plurals
	// ex. "held", for objects held by the player.
	GetPlayerBounds(string) (Bounds, error)
	// ex. take from a container
	GetObjectBounds(ident.Id) (Bounds, error)
}

type Plurals interface {
	IsPlural(word string) bool
}

// Bounds encapsulates some set of objects.
// note: we use a visitor to support map traversal without copying keys if need be.
type Bounds interface {
	// SearchBounds should visit every object in the set defined by the bounds by calling the passed function. If the visitor function returns true, the search should terminate and return a true value; otherwise the search should return false.
	SearchBounds(NounVisitor) bool
}

type NounVisitor func(NounInstance) bool

type NounInstance interface {
	// Id for the noun. Returned via ResultList.Objects() on a successful match.
	Id() ident.Id
	HasPlural(string) bool
	HasName(string) bool
	// does the noun satisfy the passed named class
	HasClass(string) bool
	// does the noun satisfy the passed named class
	HasAttribute(string) bool
}
