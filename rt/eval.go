package rt

// Execute runs a bit of code that has no return value.
type Execute interface {
	Execute(Runtime) error
}

// BoolEval represents some boolean logic expression.
type BoolEval interface {
	GetBool(Runtime) (bool, error)
}

// NumberEval represents some numeric expression.
type NumberEval interface {
	GetNumber(Runtime) (float64, error)
}

// TextEval runs a bit of code that writes into w.
type TextEval interface {
	GetText(Runtime) (string, error)
}

// NumListEval returns or generates a series of numbers.
type NumListEval interface {
	GetNumList(Runtime) ([]float64, error)
}

// TextListEval returns or generates a series of strings.
type TextListEval interface {
	GetTextList(Runtime) ([]string, error)
}
