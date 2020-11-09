package generic

import "github.com/ionous/errutil"

// error constant for iterators
const StreamExceeded errutil.Error = "stream exceeded"

type Overflow struct {
	Index, Bounds int
}

type Underflow struct {
	Index, Bounds int
}

// error for GetField, SetField
type UnknownField struct {
	Target, Field string
}

func (e Underflow) Error() string {
	return errutil.Sprint(e.Index, "below range", e.Bounds)
}

func (e Overflow) Error() string {
	return errutil.Sprint(e.Index, "above range", e.Bounds)
}

func (e UnknownField) Error() string {
	return errutil.Sprintf(`field not found "%s.%s"`, e.Target, e.Field)
}
