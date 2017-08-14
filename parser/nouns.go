package parser

type Noun interface {
	GetId() string
	HasPlural(string) bool
	HasName(string) bool
	HasClass(string) bool
	HasAttribute(string) bool
}

type NounVisitor func(n Noun) bool

type Plurals interface {
	IsPlural(name string) bool
}

// note: we use a visitor to support map traversal without copying keys if need be.
type Context interface {
	Plurals
	GetScope() Scope
}

type Scope interface {
	SearchScope(NounVisitor) bool
}
