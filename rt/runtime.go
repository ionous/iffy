package rt

import "io"

// Pluralize turns single words into their plural variants.
type Pluralize interface {
	Pluralize(single string) string
}

// VariableScope reads from and writes to a pool of named variables;
// the variables, their names, and initial values depend on the implementation and its context.
// Often, scopes are arranged in a stack with the newest scope checked for variables first, the oldest last.
type VariableScope interface {
	// GetVariable writes the value at 'name' into the value pointed to by 'pv'.
	GetVariable(name string, pv interface{}) error
	// SetVariable writes the value of 'v' into the value at 'name'.
	SetVariable(name string, v interface{}) error
}

type VariableStack interface {
	VariableScope
	PushScope(VariableScope)
	PopScope()
}

type Fields interface {
	GetField(target, field string, pv interface{}) error
	SetField(target, field string, v interface{}) error
}

// // Model interacts with the predefined world.
type Model interface {
	// 	// GetObject with the passed name.
	// 	GetObject(name string) (string, bool)
	// 	// GetRelation with the passed name.
	// 	GetRelation(name string) (Relation, bool)
}

// Ancestors customizes the parent-child event hierarchy.
// type Ancestors interface {
// 	// GetAncestors returns a stream of objects starting with the parent of the passed object, walking up whatever hierarchy the particular runtime implementation has defined.
// 	GetAncestors(Runtime, string) (ObjectStream, error)
// }

type WriterStack interface {
	io.Writer
	PushWriter(io.Writer)
	PopWriter()
}

// Runtime environment for an in-progress game.
type Runtime interface {
	Model
	Fields
	WriterStack
	VariableStack
	Random(inclusiveMin, exclusiveMax int) int

	// Ancestors
	// Pattern
	// Pluralize

}

// WritersBlock applies a writer to the runtime for the duration of fn.
// If the writer also implements io.Closer and fn is free of errors,
// w.Close gets called and its result returned.
func WritersBlock(run Runtime, w io.Writer, fn func() error) (err error) {
	run.PushWriter(w)
	e := fn()
	run.PopWriter()
	if e != nil {
		err = e
	} else if closer, ok := w.(io.Closer); ok {
		err = closer.Close()
	}
	return
}

// ScopeBlock brings the names of an object's properties into scope for the duration of fn.
func ScopeBlock(run Runtime, scope VariableScope, block Execute) (err error) {
	run.PushScope(scope)
	err = block.Execute(run)
	run.PopScope()
	return
}

// // ScopeBlock brings the names of an object's properties into scope for the duration of fn.
// func ScopeBlock(run Runtime, scope VariableScope, fn func() error) (err error) {
// 	run.PushScope(scope)
// 	err = fn()
// 	run.PopScope()
// 	return
// }

// PORTING....
// func Determine(run Runtime, p interface{}) error {
// 	return run.ExecuteMatching(run.Emplace(p))
// }
