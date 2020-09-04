package rt

import (
	"github.com/ionous/iffy/rt/writer"
)

// VariableScope reads from and writes to a pool of named variables;
// the variables, their names, and initial values depend on the implementation and its context.
// Often, scopes are arranged in a stack with the newest scope checked for variables first, the oldest last.
type VariableScope interface {
	// GetVariable writes the value at 'name' into the value pointed to by 'pv'.
	GetVariable(name string) (interface{}, error)
	// SetVariable writes the value of 'v' into the value at 'name'.
	SetVariable(name string, v interface{}) error
}

type VariableStack interface {
	VariableScope
	PushScope(VariableScope)
	PopScope()
}

type Fields interface {
	GetField(target, field string) (interface{}, error)
	SetField(target, field string, v interface{}) error
	// returns the name of the indexed field
	GetFieldByIndex(taget string, index int) (string, error)
}

// Runtime environment for an in-progress game.
type Runtime interface {
	// various runtime objects (ex. nouns, kinds, etc. ) store data addressed by name.
	Fields
	// find a function, test, or pattern addressed by name
	// pv should be a pointer to a concrete type.
	GetEvalByName(name string, pv interface{}) error
	// stacks of scopes for local variables
	VariableStack
	// nouns are grouped into potentially hierarchical "domains"
	// de/activating makes those groups hidden/visible to the runtime.
	ActivateDomain(name string, enable bool)
	// turn single words into their plural variants, and vice-versa
	PluralOf(single string) string
	SingularOf(plural string) string
	// return a pseudo-random number
	Random(inclusiveMin, exclusiveMax int) int
	// Return the built-in writer, or the current override.
	Writer() writer.Output
	// Override the current writer
	SetWriter(writer.Output) (prev writer.Output)
}

// WritersBlock applies a writer to the runtime for the duration of fn.
// If the writer also implements io.Closer and fn is free of errors,
// w.Close gets called and its result returned.
func WritersBlock(run Runtime, w writer.Output, fn func() error) (err error) {
	was := run.SetWriter(w)
	e := fn()
	run.SetWriter(was)
	if e != nil {
		err = e
	} else {
		err = writer.Close(w)
	}
	return
}
