package core

import "github.com/ionous/iffy/dl/composer"

type CompareType int

// Comparator generates comparison flags.
// FIX: a combo-box of enumeration options should be possible.
type Comparator interface {
	Compare() CompareType
}

type EqualTo struct{}
type NotEqualTo struct{}
type GreaterThan struct{}
type LessThan struct{}
type GreaterOrEqual struct{}
type LessOrEqual struct{}

func (*EqualTo) Compare() CompareType {
	return Compare_EqualTo
}
func (*NotEqualTo) Compare() CompareType {
	return 0
}
func (*GreaterThan) Compare() CompareType {
	return Compare_GreaterThan
}
func (*LessThan) Compare() CompareType {
	return Compare_LessThan
}
func (*GreaterOrEqual) Compare() CompareType {
	return Compare_LessThan | Compare_EqualTo
}
func (*LessOrEqual) Compare() CompareType {
	return Compare_GreaterThan | Compare_EqualTo
}

//go:generate stringer -type=CompareType
const (
	Compare_EqualTo CompareType = 1 << iota
	Compare_GreaterThan
	Compare_LessThan
)

func (*EqualTo) Compose() composer.Spec {
	return composer.Spec{
		Name:  "equal",
		Spec:  "=",
		Group: "comparison",
		Desc:  "Equal: Two values exactly match.",
	}
}

func (*NotEqualTo) Compose() composer.Spec {
	return composer.Spec{
		Name:  "unequal",
		Spec:  "<>",
		Group: "comparison",
		Desc:  "Not Equal To: Two values don't match exactly.",
	}
}

func (*GreaterThan) Compose() composer.Spec {
	return composer.Spec{
		Name:  "greater_than",
		Spec:  ">",
		Group: "comparison",
		Desc:  "Greater Than: The first value is larger than the second value.",
	}
}

func (*LessThan) Compose() composer.Spec {
	return composer.Spec{
		Name:  "less_than",
		Spec:  "<",
		Group: "comparison",
		Desc:  "Less Than: The first value is less than the second value.",
	}
}

func (*GreaterOrEqual) Compose() composer.Spec {
	return composer.Spec{
		Name:  "at_least",
		Spec:  ">=",
		Group: "comparison",
		Desc:  "Greater Than or Equal To: The first value is larger than the second value.",
	}
}

func (*LessOrEqual) Compose() composer.Spec {
	return composer.Spec{
		Name:  "at_most",
		Spec:  "<=",
		Group: "comparison",
		Desc:  "Less Than or Equal To: The first value is larger than the second value.",
	}
}
