package parser

import (
	"github.com/ionous/iffy/ident"
)

type Context interface {
	Plurals
	GetPlayerScope(ident.Id) (Scope, error)
	GetOtherScope(ident.Id) (Scope, error)
}

type Plurals interface {
	IsPlural(word string) bool
}

// note: we use a visitor to support map traversal without copying keys if need be.
type Scope interface {
	SearchScope(func(n NounVisitor) bool) bool
}

type NounVisitor interface {
	GetId() ident.Id
	HasPlural(string) bool
	HasName(string) bool
	HasClass(string) bool
	HasAttribute(string) bool
}
