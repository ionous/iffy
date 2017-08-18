package parser

type Context interface {
	Plurals
	GetPlayerScope(name string) (Scope, error)
	GetOtherScope(name string) (Scope, error)
}

type Plurals interface {
	IsPlural(name string) bool
}

// note: we use a visitor to support map traversal without copying keys if need be.
type Scope interface {
	SearchScope(func(n NounVisitor) bool) bool
}

type NounVisitor interface {
	GetId() string
	HasPlural(string) bool
	HasName(string) bool
	HasClass(string) bool
	HasAttribute(string) bool
}
