package tests

// BaseClass provides a simple object with every common type.
type BaseClass struct {
	Name    string `if:"id,plural:baseType classes"`
	Num     float64
	Text    string
	Object  *BaseClass
	Nums    []float64
	Texts   []string
	Objects []*BaseClass
	State   TriState
	Labeled bool
}

// DerivedClass does nothing but extend BaseClass
type DerivedClass struct {
	BaseClass `if:"parent,plural:derives"`
}

//go:generate stringer -type=ListeningState
type ListeningState int

//go:generate stringer -type=StandingState
type StandingState int

const (
	Listening ListeningState = iota
	Laughing
)

const (
	Standing StandingState = iota
	Sitting
)

type Person struct {
	Name      string `if:"id"`
	Listening ListeningState
	Standing  StandingState
}

type Kind struct {
	Name              string `if:"id"`
	PrintedName       string
	IndefiniteArticle string
	SingularPlural
	CommonProper
}

//go:generate stringer -type=SingularPlural
type SingularPlural int

const (
	SingularNamed SingularPlural = iota
	PluralNamed
)

//go:generate stringer -type=CommonProper
type CommonProper int

const (
	CommonNamed CommonProper = iota
	ProperNamed
)

// TriState provides an enum with three choices for testing.
//go:generate stringer -type=TriState
type TriState int

const (
	No TriState = iota
	Yes
	Maybe
)
