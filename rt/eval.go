package rt

import (
	g "github.com/ionous/iffy/rt/generic"
)

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

type TextEval interface {
	GetText(Runtime) (string, error)
}

type RecordEval interface {
	GetRecord(Runtime) (*g.Record, error)
}

// NumListEval returns or generates a series of numbers.
type NumListEval interface {
	GetNumList(Runtime) ([]float64, error)
}

// TextListEval returns or generates a series of strings.
type TextListEval interface {
	GetTextList(Runtime) ([]string, error)
}

type RecordListEval interface {
	GetRecordList(Runtime) (string, []*g.Record, error)
}
