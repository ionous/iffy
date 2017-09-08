package prop

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/enum"
	r "reflect"
)

// State of a single enumerated value for a field which can hold many different values of some particular enumerated type.
type State struct {
	Field
	index int
}

// Value returns true if the parent object is in the state represented by this property; false otherwise.
func (x State) Value() interface{} {
	c, idx := x.fieldValue().Int(), x.index
	match := c == int64(idx)
	return match
}

// SetValue implements rt.Property, enabling or disabling the state represented by this property; v must be a boolean value.
func (x State) SetValue(v interface{}) (err error) {
	if b, ok := v.(bool); !ok {
		err = errutil.New("expected a true/false value")
	} else {
		err = x.EnableState(b)
	}
	return
}

// Enable (or disable) the state represented by this property,
func (x State) EnableState(b bool) (err error) {
	dst := x.fieldValue()
	if b {
		// when setting to true, we are asking for this index.
		err = CoerceValue(dst, r.ValueOf(x.index))
	} else {
		// when setting to false, we have to try to generate an opposite value.
		if choices := enum.Enumerate(x.Type()); len(choices) == 0 {
			err = errutil.New("not an enumerated field")
		} else if cnt := len(choices); cnt > 2 {
			err = errutil.New("no opposite value. too many choices")
		} else {
			// idx= 0; 2-(0+1)=1
			// idx= 1; 2-(1+1)=0
			// ret can be out of range for 1 length enums
			idx := 2 - (x.index + 1)
			err = CoerceValue(dst, r.ValueOf(idx))
		}
	}
	return
}
