package parser

import (
	"github.com/ionous/iffy/ident"
)

type Context interface {
	Plurals
	// ex. "held", for objects held by the player.
	GetPlayerScope(ident.Id) (Scope, error)
	// ex. take from a container
	GetOtherScope(ident.Id) (Scope, error)
}

type Plurals interface {
	IsPlural(word string) bool
}

// Scope encapsulates some set of objects.
// note: we use a visitor to support map traversal without copying keys if need be.
type Scope interface {
	// SearchScope should visit every object in the set defined by the scope by calling the passed function. If the visitor function returns true, the search should terminate and return a true value; otherwise the search should return false.
	SearchScope(NounVisitor) bool
}

type NounVisitor func(NounInstance) bool

type NounInstance interface {
	Id() ident.Id
	HasPlural(string) bool
	HasName(string) bool
	HasClass(string) bool
	HasAttribute(string) bool
}
