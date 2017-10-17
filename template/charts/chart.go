package chart

// an interface so we can mock
type ChartMarker interface {
	MakeDigit() DigitState
	MakeDirective() DirectiveState
	MakeExpression() ExpressionState
	MakeFunction() FunctionState
	MakeNumber() NumberState
	MakeReference() ReferenceState
}

type ErrorState struct {
	error
}

func (d ErrorState) NewRune(rune) ParseState {
	return d
}

type DigitState interface {
	ParseState
	IsaDigit() bool
}

// type DirectiveState interface {
// 	ParseState
// 	GetDirective() (SpecFun, bool)
// }
// type ExpressionState interface {
// 	ParseState
// 	GetExpression() (string, bool)
// }
// type FunctionState interface {
// 	ParseState
// 	GetFunction() (SpecFun, bool)
// }
// type NumberState interface {
// 	ParseState
// 	GetNumber() (float64, bool)
// }
// type ReferenceState interface {
// 	ParseState
// 	GetReference() (string, bool)
// }

// ParseChart implements ChartMaker
type ParseChart struct{}

func (char ParseChart) MakeDigit() DigitState {
	return amiaDigit
}
func (char ParseChart) MakeDirective() DirectiveState {
	return &amiaDirective{}
}
func (char ParseChart) MakeExpression() ParseState /*ExpressionState */ {
	// return &amiaNumber{char.MakeDigit(), ""}
	return nil
}
func (char ParseChart) MakeFunction() FunctionState {
	// return &amiaNumber{char.MakeDigit(), ""}
	return nil
}
func (char ParseChart) MakeNumber() NumberState {
	return &amiaNumber{char.MakeDigit(), ""}
}
func (char ParseChart) MakeReference() ReferenceState {
	// return &amiaNumber{char.MakeDigit(), ""}
	return nil
}
