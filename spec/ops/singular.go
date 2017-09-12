package ops

import (
	r "reflect"
)

// NewValue of the passed type which can be used as a Target for building.
func NewValue(v r.Type) Single {
	return Single{r.New(v).Elem()}
}

// Single value masquerading as a struct.
type Single struct {
	r.Value
}

func (r Single) NumField() int {
	return 1
}

func (r Single) Field(i int) (ret r.Value) {
	if i == 0 {
		ret = r.Value
	}
	return
}

func (r Single) FieldByName(string) (ret r.Value) {
	return
}
