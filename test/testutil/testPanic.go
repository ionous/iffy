package testutil

import (
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/writer"
)

// PanicRuntime implements Runtime throwing a panic for every method
type PanicRuntime struct{}

var _ rt.Runtime = (*PanicRuntime)(nil)

func (PanicRuntime) ActivateDomain(name string, enable bool) {
	panic("Runtime panic")
}
func (PanicRuntime) GetEvalByName(string, interface{}) error {
	panic("Runtime panic")
}
func (PanicRuntime) GetKindByName(string) (*g.Kind, error) {
	panic("Runtime panic")
}
func (PanicRuntime) RelateTo(a, b, relation string) error {
	panic("Runtime panic")
}
func (PanicRuntime) RelativesOf(a, relation string) ([]string, error) {
	panic("Runtime panic")
}
func (PanicRuntime) ReciprocalsOf(a, relation string) ([]string, error) {
	panic("Runtime panic")
}
func (PanicRuntime) GetField(target, field string) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) SetField(target, field string, v g.Value) error {
	panic("Runtime panic")
}
func (PanicRuntime) Writer() writer.Output {
	panic("Runtime panic")
}
func (PanicRuntime) SetWriter(writer.Output) writer.Output {
	panic("Runtime panic")
}
func (PanicRuntime) PushScope(rt.Scope) {
	panic("Runtime panic")
}
func (PanicRuntime) PopScope() {
	panic("Runtime panic")
}
func (PanicRuntime) MakeRecord(kind string) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) Random(inclusiveMin, exclusiveMax int) int {
	panic("Runtime panic")
}
func (PanicRuntime) PluralOf(single string) string {
	panic("Runtime panic")
}
func (PanicRuntime) SingularOf(plural string) string {
	panic("Runtime panic")
}
