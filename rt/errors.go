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

func (e UnknownObject) Error() string {
	return errutil.Sprintf("Unknown object %q", string(e))
}

func (e UnknownTarget) Error() string {
	return errutil.Sprintf("field not found %q", e.Target)
}

func (e UnknownField) Error() string {
	return errutil.Sprintf(`field not found "%s.%s"`, e.Target, e.Field)
}
