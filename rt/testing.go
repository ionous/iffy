package rt

import (
	"io"

	"github.com/ionous/iffy/scope"
)

// Panic implements Runtime throwing a panic for every method
type Panic struct{}

func (Panic) GetField(target, field string, pv interface{}) error {
	panic("Runtime panic")
}
func (Panic) SetField(target, field string, v interface{}) error {
	panic("Runtime panic")
}
func (Panic) Write(p []byte) (n int, err error) {
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
