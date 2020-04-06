package next

type CompareType int

// CompareTo generates comparison flags.
// FIX: a combo-box of enumeration options should be possible.
type CompareTo interface {
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
