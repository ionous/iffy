package parser

type Noun interface {
	GetId() string
	HasPlural(string) bool
	HasName(string) bool
	HasClass(string) bool
	HasAttribute(string) bool
}

type Plurals interface {
	IsPlural(name string) bool
}

type Context interface {
	Plurals
	GetPlayerScope(name string) (Scope, error)
	GetOtherScope(name string) (Scope, error)
}

// note: we use a visitor to support map traversal without copying keys if need be.
type Scope interface {
	SearchScope(NounVisitor) bool
}

type NounVisitor func(n Noun) bool
