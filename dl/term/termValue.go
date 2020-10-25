package term

import (
	"github.com/ionous/iffy/rt"
)

type Value struct {
	value rt.Value
	// if needed could point back to the term definition
}

func (v *Value) SetValue(nv rt.Value) {
	if v.value.Affinity() != nv.Affinity() {
		panic("invalid value")
	}
	v.value = nv
}

func (v *Value) GetValue() rt.Value {
	return v.value
}
