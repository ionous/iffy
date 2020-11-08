package rt

import "github.com/ionous/errutil"

// error constant for iterators
const StreamExceeded errutil.Error = "stream exceeded"

type UnknownObject string

// error for GetField, SetField
type UnknownTarget struct {
	Target string
}

// error for GetField, SetField
type UnknownField struct {
	Target, Field string
}

type OutOfRange struct {
	Index, Bounds int
}

func (e UnknownObject) Error() string {
	return errutil.Sprintf("unknown object %q", string(e))
}

func (e UnknownTarget) Error() string {
	return errutil.Sprintf("target not found %q", e.Target)
}

func (e UnknownField) Error() string {
	return errutil.Sprintf(`field not found "%s.%s"`, e.Target, e.Field)
}

func (e OutOfRange) Error() string {
	return errutil.Sprint(e.Index, "out of range of", e.Bounds)
}
