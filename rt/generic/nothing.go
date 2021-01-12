package generic

import (
	"github.com/ionous/iffy/affine"
)

// Nothing is a PanicValue where every method panics except type and affinity.
type Nothing struct{}

var _ Value = (*Nothing)(nil)

// Affinity returns a blank string.
func (n Nothing) Affinity() affine.Affinity {
	return ""
}

// Type returns a blank string.
func (n Nothing) Type() string {
	return ""
}

// Bool panics
func (n Nothing) Bool() bool {
	panic("value is not a bool")
}

// Float panics
func (n Nothing) Float() float64 {
	panic("value is not a number")
}

// Int panics
func (n Nothing) Int() int {
	panic("value is not a number")
}

// String panics
func (n Nothing) String() string {
	panic("value is not a text")
}

// Record panics
func (n Nothing) Record() *Record {
	panic("value is not a record")
}

// Floats panics
func (Floats Nothing) Floats() []float64 {
	panic("value is not a number list")
}

// Strings panics
func (n Nothing) Strings() []string {
	panic("value is not a text list")
}

// Records panics
func (n Nothing) Records() []*Record {
	panic("value is not a record list")
}

// Index panics
func (n Nothing) Index(int) Value {
	panic("value is not indexable")
}

// Len panics
func (n Nothing) Len() int {
	panic("value is not measurable")
}

// FieldByName panics
func (n Nothing) FieldByName(string) (Value, error) {
	panic("value is not an object")
}

// SetFieldByName panics
func (n Nothing) SetFieldByName(string, Value) error {
	panic("value is not field writable")
}

// SetIndex panics
func (n Nothing) SetIndex(int, Value) {
	panic("value is not index writable")
}

// Append panics
func (n Nothing) Append(Value) {
	panic("value is not appendable")
}

// Slice panics
func (n Nothing) Slice(i, j int) (Value, error) {
	panic("value is not sliceable")
}

// Splice panics
func (n Nothing) Splice(start, end int, add Value) (Value, error) {
	panic("value is not spliceable")
}
