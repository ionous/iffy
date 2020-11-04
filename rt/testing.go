package rt

import (
	"github.com/ionous/iffy/affine"
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
func (Panic) Make(string) (Value, error) {
	panic("Runtime panic")
}
func (Panic) Copy(Value) (Value, error) {
	panic("Runtime panic")
}
func (Panic) GetField(target, field string) (Value, error) {
	panic("Runtime panic")
}
func (Panic) SetField(target, field string, v Value) error {
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
func (Panic) MakeRecord(kind string) (Value, error) {
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

type PanicValue struct{}

var _ Value = (*PanicValue)(nil)

func (PanicValue) Affinity() affine.Affinity {
	panic("Runtime panic")
}
func (PanicValue) Type() string {
	panic("Runtime panic")
}
func (PanicValue) GetBool() (bool, error) {
	panic("Runtime panic")
}
func (PanicValue) GetNumber() (float64, error) {
	panic("Runtime panic")
}
func (PanicValue) GetText() (string, error) {
	panic("Runtime panic")
}
func (PanicValue) GetNumList() ([]float64, error) {
	panic("Runtime panic")
}
func (PanicValue) GetTextList() ([]string, error) {
	panic("Runtime panic")
}
func (PanicValue) GetRecordList() ([]Value, error) {
	panic("Runtime panic")
}
func (PanicValue) GetLen() (int, error) {
	panic("Runtime panic")
}
func (PanicValue) GetIndex(int) (Value, error) {
	panic("Runtime panic")
}
func (PanicValue) GetField(string) (Value, error) {
	panic("Runtime panic")
}
func (PanicValue) SetField(string, Value) error {
	panic("Runtime panic")
}
