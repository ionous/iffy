package rt

import (
	"io"

	"github.com/ionous/iffy/scope"
)

// Panic implements Runtime throwing a panic for every method
type Panic struct{}

func (Panic) GetObject(obj, field string, pv interface{}) error {
	panic("Runtime panic")
}
func (Panic) SetObject(obj, field string, v interface{}) error {
	panic("Runtime panic")
}
func (Panic) Writer() io.Writer {
	panic("Runtime panic")
}
func (Panic) PushWriter(io.Writer) {
	panic("Runtime panic")
}
func (Panic) PopWriter() {
	panic("Runtime panic")
}
func (Panic) GetVariable(name string, pv interface{}) error {
	panic("Runtime panic")
}
func (Panic) SetVariable(name string, v interface{}) error {
	panic("Runtime panic")
}
func (Panic) PushScope(scope.VariableScope) {
	panic("Runtime panic")
}
func (Panic) PopScope() {
	panic("Runtime panic")
}
func (Panic) Random(inclusiveMin, exclusiveMax int) int {
	panic("Runtime panic")
}

// WriteText evaluates t and outputs the results to w.
func WriteText(run Runtime, w io.Writer, eval TextEval) (err error) {
	if t, e := eval.GetText(run); e != nil {
		err = e
	} else {
		io.WriteString(w, t)
	}
	return
}

// Run executes the passed statement using the passed runtime.
// It's helpful especially for testing.
func Run(run Runtime, exec Execute) (err error) {
	return exec.Execute(run)
}
