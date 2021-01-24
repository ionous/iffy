package generic

import (
	"github.com/ionous/iffy/affine"
)

// PanicValue is a PanicValue where every method panics except type and affinity.
type PanicValue struct{}

var _ Value = (*PanicValue)(nil)

// Affinity returns a blank string.
func (n PanicValue) Affinity() affine.Affinity {
	return ""
}

// Type returns a blank string.
func (n PanicValue) Type() string {
	return ""
}

// Bool panics
func (n PanicValue) Bool() bool {
	panic("value is not a bool")
}

// Float panics
func (n PanicValue) Float() float64 {
	panic("value is not a number")
}

// Int panics
func (n PanicValue) Int() int {
	panic("value is not a number")
}

// String panics
func (n PanicValue) String() string {
	panic("value is not a text")
}

// Record panics
func (n PanicValue) Record() *Record {
	panic("value is not a record")
}

// Floats panics
func (Floats PanicValue) Floats() []float64 {
	panic("value is not a number list")
}

// Strings panics
func (n PanicValue) Strings() []string {
	panic("value is not a text list")
}

// Records panics
func (n PanicValue) Records() []*Record {
	panic("value is not a record list")
}

// Index panics
func (n PanicValue) Index(int) Value {
	panic("value is not indexable")
}

// Len panics
func (n PanicValue) Len() int {
	panic("value is not measurable")
}

// FieldByName panics
func (n PanicValue) FieldByName(string) (Value, error) {
	panic("value is not an object")
}

// SetFieldByName panics
func (n PanicValue) SetFieldByName(string, Value) error {
	panic("value is not field writable")
}

// SetIndex panics
func (n PanicValue) SetIndex(int, Value) {
	panic("value is not index writable")
}

// Append panics
func (n PanicValue) Append(Value) {
	panic("value is not appendable")
}

// Slice panics
func (n PanicValue) Slice(i, j int) (Value, error) {
	panic("value is not sliceable")
}

// Splice panics
func (n PanicValue) Splice(start, end int, add Value) (Value, error) {
	panic("value is not spliceable")
}
