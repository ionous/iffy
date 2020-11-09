package term

import g "github.com/ionous/iffy/rt/generic"

type Value struct {
	value g.Value
	// if needed could point back to the term definition
}

func (v *Value) SetValue(nv g.Value) {
	if v.value.Affinity() != nv.Affinity() {
		panic("invalid value")
	}
	v.value = nv
}

func (v *Value) GetValue() g.Value {
	return v.value
}
