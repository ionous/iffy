package rt

import g "github.com/ionous/iffy/rt/generic"

// Execute runs a bit of code that has no return value.
type Execute interface {
	Execute(Runtime) error
}

// BoolEval represents the result of some true or false expression.
type BoolEval interface {
	GetBool(Runtime) (g.Value, error)
}

// NumberEval ... some numeric expression.
type NumberEval interface {
	GetNumber(Runtime) (g.Value, error)
}

// RecordEval ... some expression which creates a string.
type TextEval interface {
	GetText(Runtime) (g.Value, error)
}

// RecordEval yields access to a set of fields and their values.
type RecordEval interface {
	GetRecord(Runtime) (g.Value, error)
}

// ObjectEval yields access to some thing in the game world.
// Objects acts a bit like records, and a bit like strings.
type ObjectEval interface {
	GetObject(Runtime) (g.Value, error)
}

// TextListEval represents the computation of a series of numeric values.
type NumListEval interface {
	GetNumList(Runtime) (g.Value, error)
}

// TextListEval ... a series of strings.
type TextListEval interface {
	GetTextList(Runtime) (g.Value, error)
}

// RecordListEval ... a series of a set of fields.
type RecordListEval interface {
	GetRecordList(Runtime) (g.Value, error)
}
