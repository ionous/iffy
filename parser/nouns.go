package parser

type Noun interface {
	GetId() string
	HasName(string) bool
	HasClass(string) bool
	HasAttribute(string) bool
}

type NounVisitor func(n Noun) bool

// note: we use a visitor to support map traversal without copying keys if need be.
type Scope interface {
	SearchScope(NounVisitor) bool
}
