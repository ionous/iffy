package generic

import (
	"github.com/ionous/iffy/affine"
)

type Nothing struct{}

var _ Value = (*Nothing)(nil)

func (n Nothing) Affinity() affine.Affinity {
	return ""
}
func (n Nothing) Type() string {
	return ""
}
func (n Nothing) Bool() bool {
	panic("value is not a bool")
}
func (n Nothing) Float() float64 {
	panic("value is not a number")
}
func (n Nothing) Int() int {
	panic("value is not a number")
}
func (n Nothing) String() string {
	panic("value is not a text")
}
func (n Nothing) Record() *Record {
	panic("value is not a record")
}
func (n Nothing) Floats() []float64 {
	panic("value is not a number list")
}
func (n Nothing) Strings() []string {
	panic("value is not a text list")
}
func (n Nothing) Records() []*Record {
	panic("value is not a record list")
}
func (n Nothing) Index(int) Value {
	panic("value is not indexable")
}
func (n Nothing) Len() int {
	panic("value is not measurable")
}
func (n Nothing) FieldByName(string) (Value, error) {
	panic("value is not an object")
}
func (n Nothing) SetFieldByName(string, Value) error {
	panic("value is not field writable")
}
func (n Nothing) SetIndex(int, Value) {
	panic("value is not index writable")
}
func (n Nothing) Append(Value) {
	panic("value is not appendable")
}
func (n Nothing) Slice(i, j int) (Value, error) {
	panic("value is not sliceable")
}
func (n Nothing) Splice(start, end int, add Value) (Value, error) {
	panic("value is not spliceable")
}
