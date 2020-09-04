package rt

import (
	"github.com/ionous/iffy/rt/writer"
)

// Panic implements Runtime throwing a panic for every method
type Panic struct{}

func (Panic) ActivateDomain(name string, enable bool) {
	panic("Runtime panic")
}
func (Panic) GetEvalByName(string, interface{}) error {
	panic("Runtime panic")
}
func (Panic) GetFieldByIndex(target string, idx int) (string, error) {
	panic("Runtime panic")
}
func (Panic) GetField(target, field string) (interface{}, error) {
	panic("Runtime panic")
}
func (Panic) SetField(target, field string, v interface{}) error {
	panic("Runtime panic")
}
func (Panic) Writer() writer.Output {
	panic("Runtime panic")
}
func (Panic) SetWriter(writer.Output) writer.Output {
	panic("Runtime panic")
}
func (Panic) GetVariable(name string) (interface{}, error) {
	panic("Runtime panic")
}
func (Panic) SetVariable(name string, v interface{}) error {
	panic("Runtime panic")
}
func (Panic) PushScope(VariableScope) {
	panic("Runtime panic")
}
func (Panic) PopScope() {
	panic("Runtime panic")
}
func (Panic) Random(inclusiveMin, exclusiveMax int) int {
	panic("Runtime panic")
}
func (Panic) PluralOf(single string) string {
	panic("Runtime panic")
}
func (Panic) SingularOf(plural string) string {
	panic("Runtime panic")
}
