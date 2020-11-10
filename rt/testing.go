package rt

import (
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/writer"
)

// Panic implements Runtime throwing a panic for every method
type Panic struct{}

var _ Runtime = (*Panic)(nil)

func (Panic) ActivateDomain(name string, enable bool) {
	panic("Runtime panic")
}
func (Panic) GetEvalByName(string, interface{}) error {
	panic("Runtime panic")
}
func (Panic) GetKindByName(string) (*g.Kind, error) {
	panic("Runtime panic")
}
func (Panic) GetField(target, field string) (g.Value, error) {
	panic("Runtime panic")
}
func (Panic) SetField(target, field string, v g.Value) error {
	panic("Runtime panic")
}
func (Panic) Writer() writer.Output {
	panic("Runtime panic")
}
func (Panic) SetWriter(writer.Output) writer.Output {
	panic("Runtime panic")
}
func (Panic) PushScope(Scope) {
	panic("Runtime panic")
}
func (Panic) PopScope() {
	panic("Runtime panic")
}
func (Panic) MakeRecord(kind string) (g.Value, error) {
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
